package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Config struct {
	Endpoint string
	Interval string
	Target   string
}

type Data struct {
	Target     string
	DataPoints [][2]float64
}

func main() {
	config := readConfig()
	data := getData(config)
	fmt.Printf("got %d data targets\n", len(data))
	fmt.Println("they are: ")
	for i := range data {
		fmt.Printf("target: %s\n", data[i].Target)
	}
}

func readConfig() Config {
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration
}

func getData(config Config) []Data {
	ep, _ := url.Parse(config.Endpoint)
	values := url.Values{}
	values.Set("target", config.Target)
	values.Add("from", config.Interval)
	actualurl := ep.String() + "/render" + "?" + values.Encode() + "&format=json"

	resp, err := http.Get(actualurl)
	if err != nil {
		log.Fatal(err)
	}
	dec := json.NewDecoder(resp.Body)
	var m []Data
	for {
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}
	return m
}
