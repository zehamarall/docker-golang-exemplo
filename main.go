package main

import (
	"bytes"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"

	"github.com/garyburd/redigo/redis"
)

func main() {

	var bufferConnection bytes.Buffer
	bufferConnection.WriteString(os.Getenv("HOST_REDIS"))
	bufferConnection.WriteString(":")
	bufferConnection.WriteString(os.Getenv("PORT_REDIS"))

	c, err := redis.Dial("tcp", bufferConnection.String())

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c.Do("INCR", "hits")
		redisHits, err := redis.String(c.Do("GET", "hits"))
		if err != nil {
			panic(err)
		}

		hostName, err := os.Hostname()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Hello World! %q times %s  hostname %s", html.EscapeString(r.URL.Path), redisHits, hostName)
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})

	log.Fatal(http.ListenAndServe(":5000", nil))

}
