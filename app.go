package main

import (
  "context"
  "github.com/gorilla/mux"
  "golang.org/x/text/language"
  "gopkg.in/olivere/elastic.v5"
  "os"
  "io"
  "log"
  "net/http"
  "time"
)

const DEFAULT_LOCALE = "en_US"

const ES_URL = "http://127.0.0.1:9200/"

var client *elastic.Client

var LANGUAGE_MATCHER = language.NewMatcher([]language.Tag{
    language.English,
    language.Russian,
})

func getRequestLang(r *http.Request) language.Tag {
  lang, _ := r.Cookie("lang")
  accept := r.Header.Get("Accept-Language")
  tag, _ := language.MatchStrings(LANGUAGE_MATCHER, lang.String(), accept)
  return tag
}

func HotelsList(w http.ResponseWriter, req *http.Request) {
  // lang := getRequestLang(req)

  name := req.URL.Query().Get("name")
  position := req.URL.Query().Get("p")
  radius := req.URL.Query().Get("r")
  log.Println("name = " + name)
  log.Println("position = " + position)
  log.Println("radius = " + radius)

  query := elastic.NewTermQuery("user", "olivere")
  searchResult, err := client.Search("app").Type("hotel").Query(query).Pretty(true).Do(context.Background())
  if err != nil {
    panic(err)
  }
  log.Println(searchResult)

  // {"query": {"has_child": {"type": "hotel-info", "query": {"match": {"name": "Tau-Tash"}}}}}
  // query := map[string]interface{}{"query": map[string]interface{}{}}

  // if esc, err := es.New(es.WithHosts([]string{ES_HOST})); err != nil {
  //   // TODO: return internal server error
  //   panic(err)
  // }
  // resp, err := esc.Search(query)

  io.WriteString(w, "hello, world!\n")
}

func HotelsDetail(w http.ResponseWriter, req *http.Request) {
  vars := mux.Vars(req)
  id := vars["id"]
  log.Println(id)
  io.WriteString(w, "hello, world!\n")
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

  if exists, err := client.IndexExists("app").Do(ctx); err != nil {
    panic(err)
  } else if !exists {
    log.Fatalln("app index does not exist")
    os.Exit(1)
  }

  r := mux.NewRouter()
  r.HandleFunc("/hotels", HotelsList).Methods("GET")
  r.HandleFunc("/hotels/{id}", HotelsDetail).Methods("GET")
  srv := &http.Server{
    Handler:      r,
    Addr:         "127.0.0.1:8080",
    WriteTimeout: 15 * time.Second,
    ReadTimeout:  15 * time.Second,
  }
  log.Fatal(srv.ListenAndServe())
}