package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/PumpkinSeed/container-invoke/server"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	host string
	port int
	commands string
)

func init() {
	flag.StringVar(&host, "host", "localhost", "")
	flag.IntVar(&port, "port", server.DefaultPort, "")
	flag.StringVar(&commands, "commands", "", "")
	flag.Parse()
}

func main() {
	resp, err := http.Get(geturl())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		print(bodyBytes)
	}
}

func print(data []byte) {
	var response server.Response
	if err := json.Unmarshal(data, &response); err != nil {
		log.Fatal(err)
	}
	if response.Error != "" {
		fmt.Println(response.Error)
		return
	}
	for _, cr := range response.CommandResponses {
		fmt.Printf("%s>%d: %s\n",cr.Name, cr.Iteration,  cr.Command)
		if cr.Error != "" {
			fmt.Println(cr.Error)
			continue
		}
		fmt.Print(cr.StdOutput)
		fmt.Print(cr.StdError)
	}
}

func geturl() string {
	urlpath := strings.Replace(commands, ",", "/", -1)
	return fmt.Sprintf("http://%s:%d/%s", host, port, urlpath)
}
