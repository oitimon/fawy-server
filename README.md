## Word of Wisdom (WoW) Server

Project includes Server and Client applications working with Challenge-Response and Proof Of Work patterns.

@see https://en.wikipedia.org/wiki/Challenge%E2%80%93response_authentication

@see https://en.wikipedia.org/wiki/Proof_of_work

You can switch POW Challenge driver changing **WOWS_CHALLENGE** env.

### POW drivers

* **HASHBASED**: Implemented Hash-Based Challenge Algorithm (with Leading Targets). Some key advantages are: easy to implement, security, resource intensiveness, protection against replay attacks (dynamic challenges and target prefixes), adaptability (adjustable difficulty), resistance to collusion (each participant needs to individually solve the challenge).
* **GO-POW**: Adapter that connects external https://github.com/bwesterb/go-pow Challenge-Request library. To show how project can use external libraries.
* **Interface** https://github.com/oitimon/fawy-server/blob/main/pkg/pow/challenge.go, any driver can be implemented. **NUMERIC** is just for example and test purposes. 

### Configuration

Project works with .env and environment variables (they overload .env files).

### Server config
* WOWS_NETWORK - network (for example tcp4)
* WOWS_HOST - network host
* WOWS_PORT - network port
* WOWS_TIMEOUT - timeout in seconds for each handler
* WOWS_MAXHANDLERS - size of pipeline for handlers (maximum open TCP handler-connections)
* WOWS_CHALLENGE - name of Challenge-Request driver (HASHBASED, GO-POW, NUMERIC)
* WOWS_DIFFICULTY - POW difficulty

### Client config
* WOWS_HOST - network host
* WOWS_PORT - network port
* WOWS_MAXREQUESTS - size of pipeline for requests (maximum open request-connections)
* WOWS_TIMEOUT - timeout in seconds for request
* WOWS_CHALLENGE - name of Challenge-Request driver, should be the same as Server has 

## Building and running

For Server run:
```shell
docker build -f Dockerfile_server -t wow-server .
docker run -it -p 8888:8888 --rm --name wow-server wow-server
```

For Client run:
```shell
docker build -f Dockerfile_client -t wow-client .
docker run -it --rm --network host --name wow-client wow-client
```

## Testing

On local env:
```shell
go test ./...
go test -bench=. ./...
```

For basic stress test you can use (https://github.com/wg/wrk):
```shell
wrk -t6 -c100 -d10s --timeout 2s tcp://127.0.0.1:8888
```
