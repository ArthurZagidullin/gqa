package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gqa/config"
	"gqa/lib"
	"math/rand"
	"net/http"
	"time"
)

type Generator struct {
	conf       config.GeneratorConfig
	dataValues map[string]int
}

func (g *Generator) GetData() []lib.Message {
	var result []lib.Message
	for _, sconf := range g.conf.Data_sources {
		g.dataValues[sconf.Id] += rand.Intn(sconf.Max_change_step)
		result = append(result, lib.Message{Id: sconf.Id, Value: g.dataValues[sconf.Id]})
	}

	return result
}

func (g *Generator) Run() {
	period := time.Second * time.Duration(g.conf.Send_period_s)
	fmt.Printf("Run: Generator run with period %v ", period)
	ticker := time.NewTicker(period)
	for _ = range ticker.C {
		for _, d := range g.GetData() {
			go AddToQueue(d)
		}
	}

	ticker.Stop()
}

func AddToQueue(msg lib.Message) {

	byteMsg, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("AddToQueue: %v \n", err)
		return
	}

	resp, err := http.Post("http://localhost:7701/add", "application/json", bytes.NewBuffer(byteMsg))
	if err != nil {
		fmt.Printf("AddToQueue: %v \n", err)
		return
	}

	fmt.Printf("Message sended: %v. Status: %s \n", msg, resp.Status)

}

func NewGenerator(c config.GeneratorConfig) *Generator {
	values := make(map[string]int)
	for _, source := range c.Data_sources {
		values[source.Id] = source.Init_value
	}
	newGen := &Generator{
		conf:       c,
		dataValues: values,
	}
	fmt.Printf("NewGenerator: Created new Generator with %+v data sources \n", values)

	return newGen
}
