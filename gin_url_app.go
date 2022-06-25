package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type URLInfo struct {
	IsAllowed bool `json:"isAllowed"`
}

func main() {
	router := gin.Default()

	router.UseRawPath = true
	router.UnescapePathValues = true
	router.RemoveExtraSlash = false

	set_handlers(router)

	router.Run()
}

func set_handlers(router *gin.Engine) {
	// FYI path must be URL encoded i.e. http://localhost:8080/urlinfo/1/www.google.com:443/%2Fpath%2Fto%2Fthing%3Fa%3D5%0A
	router.GET("/urlinfo/1/:hostname_and_port/:path", func(c *gin.Context) {

		rdb := get_new_redis_client()

		hostname_and_port := c.Param("hostname_and_port")
		path := c.Param("path")
		key := get_url_key(hostname_and_port, path)

		is_url_allowed := get_is_url_allowed_redis(rdb, key)

		c.JSON(http.StatusOK, gin.H{"isAllowed": is_url_allowed})
	})

	router.POST("/urlinfo/1/:hostname_and_port/:path", func(c *gin.Context) {
		var json URLInfo

		rdb := get_new_redis_client()

		hostname_and_port := c.Param("hostname_and_port")
		path := c.Param("path")
		key := get_url_key(hostname_and_port, path)

		fmt.Println(key)

		if err := c.ShouldBindJSON(&json); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		set_url_is_allowed_redis(rdb, key, json.IsAllowed)

		c.String(http.StatusOK, "Success")
	})
}

func get_is_url_allowed_redis(rdb *redis.Client, key string) bool {
	lookup_result, err := rdb.Get(ctx, key).Result()
	var is_url_allowed bool

	if err == redis.Nil {
		is_url_allowed = false // if no redis result for URL, assume URL is not allowed (return false)
	} else if err != nil {
		panic(err)
	} else {
		is_url_allowed, err = strconv.ParseBool(lookup_result)
		if err != nil {
			panic(err)
		}
	}

	return is_url_allowed
}

func set_url_is_allowed_redis(rdb *redis.Client, urlKey string, isAllowed bool) {
	err := rdb.Set(ctx, urlKey, isAllowed, 0).Err()
	if err != nil {
		panic(err)
	}
}

func get_url_key(hostname_and_port string, path string) string {
	key := hostname_and_port + path
	key = strings.TrimSuffix(key, "\n") // sometimes the path ends up with a trailing new line, so we're stripping it out

	return key
}

func get_new_redis_client() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "host.docker.internal:6379", // For development purposes this connection string assumes redis is running on the same host
		Password: "",                          // no password set
		DB:       0,                           // use default DB
	})
}
