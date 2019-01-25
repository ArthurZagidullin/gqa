package main

import (
	"fmt"
	"gqa/config"
	"gqa/generator"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	conf := config.FromFile("config.json")
	fmt.Printf("main: %v\n", conf)

	for _, gconf := range conf.Generators {
		go generator.NewGenerator(gconf).Run()
	}

	fmt.Scanln()

}
