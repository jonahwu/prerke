package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Address string
	Info    []string
}

type Configs struct {
	Cfgs []Config `nodes`
}

func main() {

	var config Configs

	filename := os.Args[1]
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	//	source := []byte(data)

	if err := yaml.Unmarshal(source, &config); err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("--- config:\n%v\n\n", config)
	fmt.Println("len of cfg", len(config.Cfgs))
	fmt.Println("len of value", len(config.Cfgs[0].Info))
	fmt.Println("first info value", config.Cfgs[0].Info[0])
	fmt.Println(config.Cfgs[0].Address)

}
