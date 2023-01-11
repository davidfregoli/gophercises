package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
)

func main() {
	server := flag.String("server", "localhost", "The server hostname")
	port := flag.String("port", "8080", "The server port")
	flag.Parse()

	location := "/"
	for location != "/quit" {
		resp, err := http.Get("http://" + *server + ":" + *port + location + "?cli=true")
		if err != nil {
			panic(err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(body))

		var answer string
		fmt.Scanf("%s\n", &answer)
		location = "/" + answer

		fmt.Println("----------------------")
	}
}
