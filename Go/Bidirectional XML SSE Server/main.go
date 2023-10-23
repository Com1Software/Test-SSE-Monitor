package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

// -------------------------------------------------------------------------
func main() {
	agent := SSE()
	xip := fmt.Sprintf("%s", GetOutboundIP())
	port := "8080"
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	//
	//--- tctl 0 = normal mode test
	//         1 = high speed mode test
	//
	tctl := 0
	tc := 0
	fmt.Println("Test SSE Server")
	fmt.Printf("Operating System : %s\n", runtime.GOOS)
	fmt.Printf("Outbound IP  : %s Port : %s\n", xip, port)
	if runtime.GOOS == "windows" {
		xip = "http://localhost"
	}

	go func() {
		for {
			switch {
			case tctl == 0:
				time.Sleep(time.Second * 1)
			case tctl == 1:
				time.Sleep(time.Second * -1)
				tc++
				fmt.Printf("loop count = %d\n", tc)
			}
			dtime := fmt.Sprintf("%s", time.Now())
			msg := "<message>"
			msg = msg + "<controller>" + fmt.Sprint(GetOutboundIP()) + "</controller>"
			msg = msg + "<date_time>" + dtime[0:24] + "</date_time>"
			msg = msg + "<rand_num>" + fmt.Sprintf("%d", r1.Intn(100)) + "</rand_num>"
			msg = msg + "/<message>\n"
			event := msg
			//		event := fmt.Sprintf("Controller=%s Time=%v\n", GetOutboundIP(), dtime[0:24])
			agent.Notifier <- []byte(event)
		}
	}()
	if tctl == 0 {
		Openbrowser(xip + ":" + port)
	}
	fmt.Printf("Listening at  : %s Port : %s\n", xip, port)
	if runtime.GOOS == "windows" {
		http.ListenAndServe(":"+port, agent)
	} else {
		http.ListenAndServe(xip+":"+port, agent)
	}
}

// Openbrowser : Opens default web browser to specified url
func Openbrowser(url string) error {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "linux":
		cmd = "chromium-browser"
		args = []string{""}

	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

type Agent struct {
	Notifier    chan []byte
	newuser     chan chan []byte
	closinguser chan chan []byte
	user        map[chan []byte]bool
}

func SSE() (agent *Agent) {
	agent = &Agent{
		Notifier:    make(chan []byte, 1),
		newuser:     make(chan chan []byte),
		closinguser: make(chan chan []byte),
		user:        make(map[chan []byte]bool),
	}
	go agent.listen()
	return
}

func (agent *Agent) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "Error ", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	mChan := make(chan []byte)
	agent.newuser <- mChan
	defer func() {
		agent.closinguser <- mChan
	}()
	notify := req.Context().Done()
	go func() {
		<-notify
		agent.closinguser <- mChan
	}()
	for {
		fmt.Fprintf(rw, "%s", <-mChan)
		flusher.Flush()
	}

}

func (agent *Agent) listen() {
	for {
		select {
		case s := <-agent.newuser:
			agent.user[s] = true
		case s := <-agent.closinguser:
			delete(agent.user, s)
		case event := <-agent.Notifier:
			for userMChan, _ := range agent.user {
				userMChan <- event
			}
		}
	}

}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
