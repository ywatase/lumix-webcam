package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"sync"
	"time"

	"github.com/mattn/go-mjpeg"
)

var (
	camera   = flag.String("camera", "", "Camera host or IP")
	addr     = flag.String("addr", ":8080", "Server address")
	interval = flag.Duration("interval", 10*time.Millisecond, "interval")
)

func calcOffset(buf []byte, length int) int {
	for i := 130; i < 320 && i < length; i++ {
		if buf[i] == 0xFF && buf[i+1] == 0xD8 {
			return i
		}
	}
	return 130
}

func proxy(wg *sync.WaitGroup, stream *mjpeg.Stream) {
	defer wg.Done()

	udpAddr, err := net.ResolveUDPAddr("udp4", ":49199")
	if err != nil {
		log.Fatal(err)
	}
    conn, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	log.Printf("open socket %v", udpAddr)

	for {
		buf := make([]byte, 35000)
		length, _, err := conn.ReadFromUDP(buf[0:])
		start := calcOffset(buf[0:], length)
		buf2 := buf[start:]
		err = stream.Update(buf2)
		if err != nil {
			break
		}
	}
}

func main() {

	flag.Parse()

	stream := mjpeg.NewStreamWithInterval(*interval)

	var wg sync.WaitGroup
	wg.Add(1)
	go proxy(&wg, stream)

	http.HandleFunc("/jpeg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(stream.Current())
	})

	http.HandleFunc("/mjpeg", stream.ServeHTTP)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<img src="/mjpeg" />`))
	})

	server := &http.Server{Addr: *addr}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	go func() {
		<-sc
		server.Shutdown(context.Background())
	}()
	log.Printf("Listening on %s", *addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	stream.Close()

	wg.Wait()
}
