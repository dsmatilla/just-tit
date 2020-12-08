# Just-Tit [![Build Status](https://travis-ci.org/dsmatilla/just-tit.svg?branch=v4)](https://travis-ci.org/dsmatilla/just-tit) [![GoDoc](http://godoc.org/github.com/dsmatilla/just-tit?status.svg)](http://godoc.org/github.com/dsmatilla/just-tit) [![Foundation](https://img.shields.io/badge/Golang-Foundation-green.svg)](http://golangfoundation.org) [![Go Report Card](https://goreportcard.com/badge/github.com/dsmatilla/just-tit)](https://goreportcard.com/report/github.com/dsmatilla/just-tit)

# Just-tit

Just-tit is an adult video search engine. It uses goroutines to fetch results from several video providers in parallel and shows the results of your queries in a convenient way for both mobile users and desktop users.

Just-tit can be used locally by compiling this project or via docker since it doesn't depend on any external service. Redis caching is optional and recommended for high traffic sites.

## Usage

### Compiling

    git clone https://github.com/dsmatilla/just-tit.git
    cd just-tit/
    go build -o just-tit *go
    ./just-tit
    Navigate to http://localhost:8080

### Docker

    docker run -d -p8080:8080 dsmatilla/just-tit:v4.0

## REDIS (optional)
Just-tit can be configured to use a Redis server for caching, in order to avoid hitting the limits of the provider's APIs.

    REDISHOST=IP_ADDRESS:6379
    REDISNAME=REDIS_NAME
    REDISDBNUM=0
    REDISPASSWD=YOUR_PASSWORD


