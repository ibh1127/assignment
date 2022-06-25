# assignment

Design

![alt text](https://github.com/ibh1127/assignment/blob/main/Design.png)

Usage:
```
docker build -t ian/assignment:1 -t ian/assignment:latest .
```
```
docker run -p 127.0.0.1:6379:6379 -d redis
```
```
docker run -p 8080:8080 ian/assignment:latest
```
```
curl http://localhost:8080/urlinfo/1/www.google.com:443/%2Fpath%2Fto%2Fthing%3Fa%3D5%0A
```
