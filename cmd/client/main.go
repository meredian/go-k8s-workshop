package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	var n int
	flag.IntVar(&n, "n", 10, "-n=10")
	flag.Parse()
	fmt.Printf("Iterations: %v\n", n)

	ticker := time.NewTicker(100 * time.Millisecond)

	for i := 0; i < n; i++ {
		go makeRequest(i)
		<-ticker.C
	}
}

func makeRequest(i int) {
	client := http.Client{}
	response, err := client.Get("https://austin.k8s.gromnsk.ru/get-status?message_id=ABCDE&user_id=user")
	if err != nil {
		log.Fatalf("failed to get a http request: %v\n", err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("failed to read a response: %v\n", err)
	}
	defer response.Body.Close()

	log.Printf("Request # %d, Status: %v, response: %v\n", i, response.Status, string(body))
}
