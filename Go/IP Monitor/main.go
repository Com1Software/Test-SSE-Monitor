package main

import (
	"bufio"
	"log"
	"net/http"
)

func main() {
	url := "http://com1software.com"
	// url = "http://192.168.1.105:8080"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, erra := reader.ReadBytes('\n')
		if erra != nil {
			log.Fatal(err)
		}

		log.Println(string(line))
	}

}
