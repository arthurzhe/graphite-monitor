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
	Endpoint  string
	Interval  string
	Target    string
	Threshold float64
	Rule      string
}

type Data struct {
	Target     string
	DataPoints [][2]float64
}

func main() {
	config := readConfig()
	data := getData(config)
	monitorData(data, config.Rule, config.Threshold)
}

func monitorData(d []Data, rule string, thres float64) {
	for i := range d {
		data := d[i]

		switch rule {
		case "==":
			for j := range data.DataPoints {
				if data.DataPoints[j][0] == thres {
					fmt.Printf("target : %s is ", data.Target)
					fmt.Printf("equal to the threshold of %f\n", thres)
					break
				}
			}
		case "!=":
			for j := range data.DataPoints {
				if data.DataPoints[j][0] != thres {
					fmt.Printf("target : %s is ", data.Target)
					fmt.Printf("not equal to the threshold of %f\n", thres)
					break
				}
			}
		case "<":
			for j := range data.DataPoints {
				if data.DataPoints[j][0] < thres {
					fmt.Printf("target : %s is ", data.Target)
					fmt.Printf("less than the threshold of %f\n", thres)
					break
				}
			}
		case "<=":
			for j := range data.DataPoints {
				if data.DataPoints[j][0] <= thres {
					fmt.Printf("target : %s is ", data.Target)
					fmt.Printf("less than or equal to the threshold of %f\n", thres)
					break
				}
			}
		case ">":
			for j := range data.DataPoints {
				if data.DataPoints[j][0] > thres {
					fmt.Printf("target : %s is ", data.Target)
					fmt.Printf("greater than the threshold of %f\n", thres)
					break
				}
			}
		case ">=":
			for j := range data.DataPoints {
				if data.DataPoints[j][0] >= thres {
					fmt.Printf("target : %s is ", data.Target)
					fmt.Printf("greater than or equal to the threshold of %f\n", thres)
					break
				}
			}
		}
	}

	// 	==    equal
	// !=    not equal
	// <     less
	// <=    less or equal
	// >     greater
	// >=    greater or equal
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
