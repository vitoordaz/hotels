#!/bin/bash

curl -XPUT '127.0.0.1:9200/app?pretty' -H 'Content-Type: application/json' -d'{
  "mappings": {
    "hotel": {
      "properties": {
        "location": {"type": "geo_point"}
      }
    },
    "hotel-info": {
      "_parent": {"type": "hotel"},
      "properties": {
        "language": {"type": "text"},
        "name": {
          "type": "text",
          "fields": {
            "en_US": {"type": "text", "analyzer": "english"},
            "ru_RU": {"type": "text", "analyzer": "russian"}
          }
        }
      }
    }
  }
}'

curl -XPUT '127.0.0.1:9200/app/hotel/1?pretty' -H 'Content-Type: application/json' -d'{
  "location": "47.6400733,-122.3362087"
}'

curl -XPUT '127.0.0.1:9200/app/hotel-info/ru_RU-1?parent=1&pretty' -H 'Content-Type: application/json' -d'{
  "language": "ru_RU",
  "name": "Гостиничный комплекс Тау-Таш",
  "address": "Горнолыжная Улица 33, Абзаково, Россия",
  "description": "К услугам гостей этого гостиничного комплекса на горнолыжном курорте Ново-Абзаково – аквапарк с крытым бассейном, финская сауна и боулинг. Также постояльцам предоставляется бесплатный трансфер на микроавтобусе до подъемника (800 метров)."
}'

curl -XPUT '127.0.0.1:9200/app/hotel-info/en_US-1?parent=1&pretty' -H 'Content-Type: application/json' -d'{
  "language": "en_US",
  "name": "Tau-Tash Hotel Complex",
  "address": "Gornolyzhnaya Street 33 , 453565 Abzakovo, Russia",
  "description": "Showcasing a seasonal outdoor pool and barbecue, Tau-Tash Hotel Complex is located in Abzakovo."
}'