# Gin Application


## Design

![alt text](https://github.com/ibh1127/assignment/blob/main/Design.png)


## API Usage

```
GET /urlinfo/1/{hostname_and_port}/{original_path_and_query_string}
```
This returns a json object that indicates weather or not a URL is allowed/secure:

{
“isAllowed”: bool
}

```
PUT /urlinfo/1/{hostname_and_port}/{original_path_and_query_string}
```
This PUT uses the same json object to insert/update a URL into database indicating whether or not the URL is allowed/secure:

{
“isAllowed”: bool
}

## Quickstart:
```
docker build -t ian/assignment:1 -t ian/assignment:latest .
```
```
docker run -p 127.0.0.1:6379:6379 -d redis
```
```
docker run -p 8080:8080 ian/assignment:latest
```

### GET Test (curl command)
```
curl http://localhost:8080/urlinfo/1/www.google.com:443/%2Fpath%2Fto%2Fthing%3Fa%3D5%0A
```

### PUT Test (python script)
```
import requests
requests.put('http://localhost:8080/urlinfo/1/www.google.com:443/%2Fpath%2Fto%2Fthing%3Fa%3D5%0A', json={"isAllowed": True})
```

## Technologies Utilized:

Languages: Go
- Go provides great performance and concurrency with its light-weight runtime and use of coroutines known as goroutines.
It should outperform other languages/frameworks that use event-loops or thread pools.

Web Service Framework: Gin
- Gin is a simple Go web framework that will perform well under high load. Each user request/connection will result
in spawning a goroutine. Gin will allow for a high level of concurrent users without consuming too much memory due to
the low stack size of goroutines.

Storage: Redis
- Redis is an in-memory key-value database that will store our URL data. Redis will provide low-latency access which is important for our use-case.
It can also scale out horizontally as a cluster to provide additional storage.

## Problems and solutions:
```
The size of the URL list could grow infinitely,
how might you scale this beyond the memory capacity of this VM?
```
I have implemented this by utilizing Redis as our database which can live on a separate VM/host.
Redis can scale horizontally to provide additional storage.

Given more time I would implement a secondary storage method such as DynamoDB since scaling Redis could be very expensive.
DynamoDB could give us cheaper storage and store less-frequently accessed URLs. This would increase latenecy, but reduce storage 
costs for a large URL list.

```
The number of requests may exceed the capacity of this VM, how
might you solve that? Bonus if you implement this.
```

This problem could be solved by horizontally scaling the docker container that I have created. The reason that this solution can scale horizontally 
is that the application is stateless, due to the fact that our state is stored on Redis.
The horizontal scaling could be provided by a variety of methods and cloud services such as AWS ECS.

``` 
What are some strategies you might use to update the service
with new URLs? Updates may be as much as 5 thousand URLs a day
with updates arriving every 10 minutes.
```

 I created a PUT handler that takes in the json object:

{
“isAllowed”: bool
}

This allows URLs to be added to the redis database through the web service. 
Although this handler only handles a single URL at a time, 5000 URLs a day should be no big deal for this design.
Given more time to work on this project, a PUT handler could be created to for bulk inserts/updates.
