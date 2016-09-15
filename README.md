Consul Watch
===========

Consume Consul watch API, uses checks option to detect health checks and send to dynamodb.  

## Setup
Create and set table name, default is `dashboard`.  
Set AWS region

## Build

```sh
glide install // Install dependencies

go build
```

##  Running

```sh
consul watch -http-addr=your.consul.address -type checks consul_watch
```
