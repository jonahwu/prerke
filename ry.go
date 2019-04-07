package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Foo string
	Bar []string
}

type Configs struct {
	Cfgs []Config `foobars`
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
	fmt.Println("len of value", len(config.Cfgs[0].Bar))
	fmt.Println(config.Cfgs[1].Bar[0])

}
