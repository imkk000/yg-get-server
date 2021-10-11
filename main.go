package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"

	"github.com/joho/godotenv"
)

const (
	urlFormat = "http://%s/users/getserverlist"
)

type (
	ServerList struct {
		Servers []struct {
			Id      int    `json:"serverid"`
			Name    string `json:"name"`
			Traffic int    `json:"traffic"`
		} `json:"servers"`
	}
)

func (s ServerList) Len() int {
	return len(s.Servers)
}

func (s ServerList) Less(i, j int) bool {
	return s.Servers[i].Id < s.Servers[j].Id
}

func (s ServerList) Swap(i, j int) {
	s.Servers[i], s.Servers[j] = s.Servers[j], s.Servers[i]
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

func main() {
	err := godotenv.Load()
	failOnError(err, "Load Env File")

	resp, err := http.Get(fmt.Sprintf(urlFormat, os.Getenv("HOST")))
	failOnError(err, "Call Api")

	data, err := ioutil.ReadAll(resp.Body)
	failOnError(err, "Get Body")

	var svList ServerList
	err = json.Unmarshal(data, &svList)
	failOnError(err, "Unmarshal Body")

	sort.Sort(svList)

	for _, sv := range svList.Servers {
		fmt.Printf("%3d%% %s\n", sv.Traffic, sv.Name)
	}
}
