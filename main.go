package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"time"

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
		fmt.Printf("%s: %s\n", msg, err)
		os.Exit(1)
	}
}

//go:embed .env
var envStr string

func main() {
	var envMap map[string]string
	var err error
	envMap, err = godotenv.Read()
	if err != nil {
		envMap, err = godotenv.Unmarshal(envStr)
		failOnError(err, "Load Env From Embedded File")
	}

	c := http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(urlFormat, envMap["HOST"]), nil)
	failOnError(err, "Create Request")

	resp, err := c.Do(req)
	failOnError(err, "Call Api with Request")

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
