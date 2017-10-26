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
            "en": {"type": "text", "analyzer": "english"},
            "ru": {"type": "text", "analyzer": "russian"}
          }
        }
      }
    }
  }
}'

curl -XPUT '127.0.0.1:9200/app/hotel/1?pretty' -H 'Content-Type: application/json' -d'{
  "location": {"lat": 53.806803, "lon": 58.635771}
}'

curl -XPUT '127.0.0.1:9200/app/hotel-info/ru-1?parent=1&pretty' -H 'Content-Type: application/json' -d'{
  "language": "ru_RU",
  "name": "Гостиничный комплекс Тау-Таш",
  "address": "Горнолыжная Улица 33, Абзаково, Россия",
  "description": "К услугам гостей этого гостиничного комплекса на горнолыжном курорте Ново-Абзаково – аквапарк с крытым бассейном, финская сауна и боулинг. Также постояльцам предоставляется бесплатный трансфер на микроавтобусе до подъемника (800 метров)."
}'

curl -XPUT '127.0.0.1:9200/app/hotel-info/en-1?parent=1&pretty' -H 'Content-Type: application/json' -d'{
  "language": "en_US",
  "name": "Tau-Tash Hotel Complex",
  "address": "Gornolyzhnaya Street 33 , 453565 Abzakovo, Russia",
  "description": "Showcasing a seasonal outdoor pool and barbecue, Tau-Tash Hotel Complex is located in Abzakovo."
}'

curl -XPUT '127.0.0.1:9200/app/hotel/2?pretty' -H 'Content-Type: application/json' -d'{
  "location": {"lat": 53.827882, "lon": 58.593006}
}'

curl -XPUT '127.0.0.1:9200/app/hotel-info/ru-2?parent=2&pretty' -H 'Content-Type: application/json' -d'{
  "language": "ru_RU",
  "name": "Гостиница Abzakovo Weekend",
  "address": "Улица Лесная, 39, Абзаково, Россия",
  "description": "Гостиница «Абзаково - Weekend» расположена в поселке Абзаково."
}'

curl -XPUT '127.0.0.1:9200/app/hotel-info/en-2?parent=2&pretty' -H 'Content-Type: application/json' -d'{
  "language": "en_US",
  "name": "Abzakovo Weekend",
  "address": "Ulitsa Lesnaya 39, 455039 Abzakovo, Russia",
  "description": "Featuring free WiFi and a sauna, Abzakovo Weekend offers accommodations in Abzakovo. Free private parking is available on site."
}'

curl -XPUT '127.0.0.1:9200/app/hotel/3?pretty' -H 'Content-Type: application/json' -d'{
  "location": {"lat": 53.811105, "lon": 58.636158}
}'

curl -XPUT '127.0.0.1:9200/app/hotel-info/ru-3?parent=3&pretty' -H 'Content-Type: application/json' -d'{
  "language": "ru_RU",
  "name": "Эдельвейс",
  "address": "Горный проезд 1A, Абзаково, Россия",
  "description": "Гостиничный комплекс Эдельвейс расположен на территории лыжного курорта Абзаково, в 1 км от горнолыжного центра."
}'

curl -XPUT '127.0.0.1:9200/app/hotel-info/en-3?parent=3&pretty' -H 'Content-Type: application/json' -d'{
  "language": "en_US",
  "name": "Edelweiss Hotel",
  "address": "Gorniy Proezd 1A, 453565 Abzakovo, Russia",
  "description": "Featuring free WiFi throughout the property, Edelweiss Hotel offers accommodations in Abzakovo. Free private parking is available on site."
}'

curl -XPUT '127.0.0.1:9200/app/hotel/4?pretty' -H 'Content-Type: application/json' -d'{
  "location": {"lat": 53.810764, "lon": 58.641739}
}'

curl -XPUT '127.0.0.1:9200/app/hotel-info/ru-4?parent=4&pretty' -H 'Content-Type: application/json' -d'{
  "language": "ru_RU",
  "name": "Гостиница Усадьба Еловое",
  "address": "Переулок Ёлочка 3, Абзаково, Россия",
  "description": "Отель Усадьба Еловое расположен в сосновом бору на территории популярного лыжного курорта Новоабзаково на Южном Урале в Республике Башкортостан."
}'

curl -XPUT '127.0.0.1:9200/app/hotel-info/en-4?parent=4&pretty' -H 'Content-Type: application/json' -d'{
  "language": "en_US",
  "name": "Usadba Elovoe Hotel",
  "address": "Pereulok Elochka 3, 455000 Abzakovo, Russia",
  "description": "Located in Abzakovo in the region of Bashkortostan, 30 miles from Magnitogorsk, Usadba Elovoe Hotel features a playground and views of the mountains."
}'