package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	Password string
	Nodes    []TConfig
}

type TConfig struct {
	Address string
	Info    []string
}

func main() {
	filename := os.Args[1]
	var config Config
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Value: %#v\n", config.Nodes[0].Address)
	fmt.Printf("Value: %#v\n", config.Nodes[0].Info[0])
	fmt.Println("Value: ", config.Password)
	fmt.Println(config)
}
