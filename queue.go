package main

import (
	"encoding/json"
	"fmt"
	"gqa/config"
	"gqa/lib"
	"io/ioutil"
	"net"
	"net/http"
	"sync/atomic"
)

type agregator chan<- lib.Message

var (
	conf          = config.FromFile("config.json").Queue
	subscribing   = make(chan agregator)
	leaving       = make(chan agregator)
	messages      = make(chan lib.Message, conf.Size)
	messagesCount int64
)

func broadcaster() {
	defer panic("broadcaster is closing")
	agregators := make(map[agregator]bool)
	for {
		select {
		case ag := <-leaving:
			fmt.Printf("Aggregator: %s disconnected \n", ag)
			delete(agregators, ag)
			close(ag)
		case msg := <-messages:
			atomic.AddInt64(&messagesCount, -1)
			fmt.Printf("broadcaster: Queue(%d) get message %v \n", messagesCount, msg)
			//Широковещательное сообщение во все
			//каналы аггрегаторв
			for ag := range agregators {
				ag <- msg
			}
		case ag := <-subscribing:
			agregators[ag] = true

		}
	}
}

func handleAgConn(conn net.Conn) {
	ch := make(chan lib.Message)

	agregator := conn.RemoteAddr().String()
	fmt.Printf("Aggregator: %s connected \n", agregator)

	subscribing <- ch

	for msg := range ch {
		byteJson, _ := json.Marshal(msg)
		if _, err := fmt.Fprintln(conn, string(byteJson)); err != nil {
			fmt.Printf("sender2ag: Err %v \n", err)
			break
		}
		fmt.Printf("sender2ag: %v \n", string(byteJson))
	}

	leaving <- ch
	conn.Close()
}

func main() {
	// Прием сообщений
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		var msg lib.Message
		if err := json.Unmarshal(b, &msg); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		atomic.AddInt64(&messagesCount, 1)
		messages <- msg
	})

	// Подключение аггрегаторов
	agConns, err := net.Listen("tcp", ":7700")
	if err != nil {
		panic(err)
	}
	defer agConns.Close()

	go http.ListenAndServe("localhost:7701", nil)

	go broadcaster()

	for {
		conn, err := agConns.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleAgConn(conn)
	}

}
