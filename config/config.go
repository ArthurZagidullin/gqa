package config

import (
	"encoding/json"
	"fmt"
	"gqa/lib"
	"io/ioutil"
	"os"
)

type AgregatorConfig struct {
	Sub_ids           []string
	Agregate_period_s int
}

type GeneratorConfig struct {
	Timeout_s, Send_period_s int
	Data_sources             []struct {
		Id                          string
		Init_value, Max_change_step int
	}
}

type QueueConfig struct {
	Size int
}

type Config struct {
	Generators   []GeneratorConfig
	Agregators   []AgregatorConfig
	Queue        QueueConfig
	Storage_type int
}

func FromFile(filename string) *Config {
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	config := &Config{}

	if err := json.Unmarshal([]byte(byteValue), config); err != nil {
		fmt.Println(err)
	}

	return config
}

type Saver interface {
	Save(id string, avg int) error
}

func (c *Config) Storage() Saver {
	fmt.Printf("Storage type %d \n", c.Storage_type)
	if c.Storage_type == 0 {
		return &lib.SimpleStorage{}
	} else if c.Storage_type == 1 {
		return lib.NewFileStorage("./storage.db")
	}
	return nil
}
