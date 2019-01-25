package main

import (
	"fmt"
	"gqa/aggregator"
	"gqa/config"
	"sync"
	//"gqa/lib"
	"net"
)

func main() {
	conf := config.FromFile("config.json")
	storage := conf.Storage()

	if storage == nil {
		panic("Storage creating error")
	}

	fmt.Printf("main: %v\n", conf)

	group := sync.WaitGroup{}

	for _, aconf := range conf.Agregators {
		conn, err := net.Dial("tcp", "localhost:7700")
		if err != nil {
			panic(err)
		}
		group.Add(1)
		go aggregator.NewAggregator(&group, aconf, storage, conn).Run()
	}
	group.Wait()
	fmt.Println("Aggregation ended")
}
