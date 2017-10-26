# 1. Setup elasticsearch

Run elasticsearch docker container

```
sudo ./run_elasticsearch.sh
```

Once elasticsearch running completes run 

```
./init_elasticsearch
```

This will create mappings and initial documents

# 2. Run application

Install [godep](https://github.com/golang/dep)

```
brew install dep
```

Install dependencies

```
dep ensure
```

Build app

```
go build
```

Application is ready to start

```
./hotels
```

Here is a sample requests using httpie

```
# Accept-Language header controls response language. If this header is not given than default language is english.

# Returns all hotels that are in radious on 1 km to location lat: 53.806803 lon: 58.63577
http 'http://127.0.0.1:8080/hotels?l=53.806803,58.635771&r=1' Accept-Language:ru-RU

# Same as above but in english
http 'http://127.0.0.1:8080/hotels?l=53.806803,58.635771&r=1' Accept-Language:en-US

# Returns default information about hotel with id 1
http 'http://127.0.0.1:8080/hotels/1' Accept-Language:en-US

# Searches for hotels with name Abzakovo
http 'http://127.0.0.1:8080/hotels?name=Abzakovo' Accept-Language:en-US

# Same as above but in radious of one 1km from lat: 53.811105 lon: 58.636158
http 'http://127.0.0.1:8080/hotels?name=Abzakovo&l=53.811105,58.636158&r=1' Accept-Language:ru-RU
```