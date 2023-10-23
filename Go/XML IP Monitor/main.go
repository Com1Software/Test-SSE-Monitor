package main

import (
	"bufio"
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type message struct {
	Controller string `xml:"controller"`
	DateTime   string `xml:"date_time"`
	RandNum    string `xml:"rand_num"`
}

func main() {
	url := "http://com1software.com"
	url = "http://192.168.1.105:8080"
	tc := 0
	tctl := 0
	msg := &message{}
	switch {
	case tctl == 0:
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
			xml.Unmarshal(line, &msg)
			fmt.Printf("%s  -  %s  -  %s\n", string(msg.Controller), string(msg.DateTime), string(msg.RandNum))

		}
	case tctl == 1:
		ctx := cancelCtxOnSigterm(context.Background())
		startWork(ctx, url)
		fmt.Println("test")
		fmt.Println(tc)

	}
}

func cancelCtxOnSigterm(ctx context.Context) context.Context {
	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-exitCh
		cancel()
	}()
	return ctx
}
func startWork(ctx context.Context, url string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	// resp, err := http.Get(url)
	// if err != nil {
	//	log.Fatal(err)
	// }
	//reader := bufio.NewReader(resp.Body)

	for {
		fmt.Println("xx")
		if err := work(ctx, url); err != nil {
			fmt.Printf("failed to do work: %s", err)
		}
		select {
		case <-ticker.C:
			fmt.Println("test")
			continue
		case <-ctx.Done():
			fmt.Println("oo")
			return
		}
	}
}

func work(ctx context.Context, url string) error {
	fmt.Println("doing work")

	return nil
}
