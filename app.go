package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/text/language"
	"gopkg.in/olivere/elastic.v5"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const ES_URL = "http://127.0.0.1:9200/"
const MAX_LIMIT = 25

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Hotel struct {
	Id          string   `json:"id"`
	Location    Location `json:"location"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Address     string   `json:"address,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

var NotFound = ErrorResponse{Success: false, Error: "Not found"}
var InternalServerError = ErrorResponse{Success: false, Error: "Internal server error"}

type ListHotelResponse struct {
	Success bool     `json:"success"`
	Total   int64    `json:"total"`
	Data    []*Hotel `json:"data,omitempty"`
}

var client *elastic.Client

var LANGUAGE_MATCHER = language.NewMatcher([]language.Tag{
	language.English,
	language.Russian,
})

func writeResponse(w http.ResponseWriter, status int, resp interface{}) {
	if body, err := json.Marshal(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"success\": false, \"errors\":[\"Internal Server Error\"]}"))
	} else {
		w.WriteHeader(status)
		w.Write(body)
	}
}

func getRequestLang(r *http.Request) language.Tag {
	lang, _ := r.Cookie("lang")
	accept := r.Header.Get("Accept-Language")
	tag, _ := language.MatchStrings(LANGUAGE_MATCHER, lang.String(), accept)
	return tag
}

func getRequestOffsetLimit(r *http.Request) (int, int) {
	var offset uint64
	var limit uint64
	var err error

	o := strings.TrimSpace(r.URL.Query().Get("offset"))
	if offset, err = strconv.ParseUint(o, 10, 32); err != nil {
		offset = 0
	}

	l := strings.TrimSpace(r.URL.Query().Get("limit"))
	if limit, err = strconv.ParseUint(l, 10, 32); err != nil || limit > MAX_LIMIT {
		limit = MAX_LIMIT
	}

	return int(offset), int(limit)
}

func addDistanceFilter(l string, r string, query *elastic.BoolQuery) {
	l = strings.TrimSpace(l)
	r = strings.TrimSpace(r)

	if l == "" || r == "" {
		return
	}

	radius, err := strconv.ParseFloat(r, 64)
	if err != nil {
		return
	}

	lat_lon := strings.SplitN(l, ",", 2)
	if len(lat_lon) != 2 {
		return
	}

	lat, err := strconv.ParseFloat(lat_lon[0], 64)
	if err != nil {
		return
	} else if lat < -90 || lat > 90 {
		return
	}

	lon, err := strconv.ParseFloat(lat_lon[1], 64)
	if err != nil {
		return
	} else if lon < -180 || lon > 180 {
		return
	}

	query.Filter(elastic.NewGeoDistanceQuery("location").Lat(lat).Lon(lon).Distance(fmt.Sprintf("%.3fkm", radius)))
}

func HotelsList(w http.ResponseWriter, req *http.Request) {
	lang := getRequestLang(req)

	searchQuery := elastic.NewBoolQuery()
	name := strings.TrimSpace(req.URL.Query().Get("name"))
	if name != "" {
		searchQuery.Must(
			elastic.NewHasChildQuery(
				"hotel-info",
				elastic.NewBoolQuery().Should(
					elastic.NewFuzzyQuery("name.en", name),
					elastic.NewFuzzyQuery("name.ru", name),
				),
			),
		)
	}
	addDistanceFilter(req.URL.Query().Get("l"), req.URL.Query().Get("r"), searchQuery)

	ctx := context.Background()

	offset, limit := getRequestOffsetLimit(req)
	searchResult, err := client.Search("app").Type("hotel").Query(searchQuery).From(offset).Size(limit).Do(ctx)
	if err != nil {
		log.Printf("elastic search error: %s\n", err)
		writeResponse(w, http.StatusInternalServerError, InternalServerError)
		return
	}

	data := []*Hotel{}
	if len(searchResult.Hits.Hits) > 0 {
		multiGetItems := []*elastic.MultiGetItem{}

		for _, hit := range searchResult.Hits.Hits {
			h := Hotel{Id: hit.Id}
			multiGetItems = append(
				multiGetItems,
				elastic.NewMultiGetItem().Index("app").Type("hotel-info").Routing(hit.Id).Id(fmt.Sprintf("%s-%s", lang, hit.Id)),
			)
			if err := json.Unmarshal(*hit.Source, &h); err != nil {
				log.Printf("json unmarshal error: %s\n", err)
				writeResponse(w, http.StatusInternalServerError, InternalServerError)
				return
			}
			data = append(data, &h)
		}

		multiGetResult, err := client.MultiGet().Add(multiGetItems...).Do(ctx)
		if err != nil {
			log.Printf("elastic multi get error: %s\n", err)
			writeResponse(w, http.StatusInternalServerError, InternalServerError)
			return
		}
		for i, d := range multiGetResult.Docs {
			if err := json.Unmarshal(*d.Source, data[i]); err != nil {
				log.Printf("json unmarshal error: %s\n", err)
				writeResponse(w, http.StatusInternalServerError, InternalServerError)
				return
			}
		}
	}
	writeResponse(w, http.StatusOK, ListHotelResponse{
		Success: true,
		Total:   searchResult.Hits.TotalHits,
		Data:    data,
	})
}

func HotelsDetail(w http.ResponseWriter, req *http.Request) {
	lang := getRequestLang(req)

	vars := mux.Vars(req)
	id := vars["id"]

	ctx := context.Background()

	hotel := Hotel{Id: id}

	d, err := client.Get().Index("app").Type("hotel").Id(id).Do(ctx)
	if err != nil {
		if elastic.IsStatusCode(err, http.StatusNotFound) {
			writeResponse(w, http.StatusNotFound, NotFound)
		} else {
			log.Printf("elastic get error: %s\n", err)
			writeResponse(w, http.StatusInternalServerError, InternalServerError)
		}
		return
	}
	if err := json.Unmarshal(*d.Source, &hotel); err != nil {
		log.Fatal(err)
		writeResponse(w, http.StatusInternalServerError, InternalServerError)
		return
	}

	d, err = client.Get().Index("app").Type("hotel-info").Routing(id).Id(fmt.Sprintf("%s-%s", lang, id)).Do(ctx)
	if err != nil {
		if elastic.IsStatusCode(err, http.StatusNotFound) {
			writeResponse(w, http.StatusNotFound, NotFound)
		} else {
			log.Printf("elastic get error: %s\n", err)
			writeResponse(w, http.StatusInternalServerError, InternalServerError)
		}
		return
	}
	if err := json.Unmarshal(*d.Source, &hotel); err != nil {
		log.Fatal(err)
		writeResponse(w, http.StatusInternalServerError, InternalServerError)
		return
	}

	writeResponse(w, http.StatusOK, hotel)
}

func main() {
	ctx := context.Background()

	var err error

	var options []elastic.ClientOptionFunc
	options = append(options, elastic.SetURL(ES_URL))
	options = append(options, elastic.SetInfoLog(log.New(os.Stderr, "", log.LstdFlags)))
	options = append(options, elastic.SetTraceLog(log.New(os.Stderr, "", log.LstdFlags)))
	options = append(options, elastic.SetSniff(false))
	client, err = elastic.NewClient(options...)
	if err != nil {
		panic(err)
	}
	defer client.Stop()

	log.Printf("Checking if index app exists")
	if exists, err := client.IndexExists("app").Do(ctx); err != nil {
		panic(err)
	} else if !exists {
		log.Fatalln("app index does not exist")
		os.Exit(1)
	}
	log.Printf("Index app exists")

	r := mux.NewRouter()
	r.HandleFunc("/hotels", HotelsList).Methods("GET")
	r.HandleFunc("/hotels/{id}", HotelsDetail).Methods("GET")
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	log.Printf("Listening on 127.0.0.1:8080")
	srv := &http.Server{
		Handler:      loggedRouter,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
