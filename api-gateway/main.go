package main

import (
	"log"
	"fmt"
	// "net/url"
	"net/http"
	// "net/http/httputil"
)

func main() {
	// Define services URLs
	// userServiceURL, _ := url.Parse("http://user-service:8081")
	// gameServiceURL, _ := url.Parse("http://game-service:8082")
	// lobbyServiceURL, _ := url.Parse("http://lobby-service:8083")
	// friendServiceURL, _ := url.Parse("http://friend-circle-service:8084")
	// chatServiceURL, _ := url.Parse("http://game-chat-service:8085")

	// // create reverse proxies
	// userProxy := httputil.NewSingleHostReverseProxy(userServiceURL)


	// Router handler
	http.HandleFunc("/", func(w http.ResponseWriter, r * http.Request){
		switch{
		case r.URL.Path == "/health":
			fmt.Fprintln(w, "API Gateway is healthy")
		}
	})
	


	log.Println("API Gateway running on: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))


}