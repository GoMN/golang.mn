package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var (
	config conf
)

func init() {
	log.Println("loading configuration")
	b, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Fatal("application configuration (conf.json) not found", err)
	}
	json.Unmarshal(b, &config)
	log.Println("configuration loaded")
}

type conf struct{
	Cache  conf_cache
	Meetup conf_meetup
	Maps   conf_maps
}

type conf_cache struct {
	LocalTimeout  int64
	RemoteTimeout int64
}

type conf_meetup struct {
	BaseUrl string
	Key     string
}

type conf_maps struct {
	Key string
}
