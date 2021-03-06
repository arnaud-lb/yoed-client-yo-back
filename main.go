package main

import (
	clientInterface "github.com/yoed/yoed-client-interface"
	"net/http"
	"net/url"
	"log"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type YoBackYoedClient struct {
	clientInterface.BaseYoedClient
	config *YoBackYoedClientConfig
}

type YoBackYoedClientConfig struct {
	ApiKey string `json:"apiKey"`
}

func (c *YoBackYoedClient) loadConfig(configPath string) (*YoBackYoedClientConfig, error) {
	configJson, err := clientInterface.ReadConfig(configPath)

	if err != nil {
		return nil, err
	}

	config := &YoBackYoedClientConfig{}

	if err := json.Unmarshal(configJson, config); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *YoBackYoedClient) Handle(username string) {
	resp, err := http.PostForm("http://api.justyo.co/yo/", url.Values{
		"api_token": {c.config.ApiKey},
		"username":  {username},
	})

	log.Printf("Yo-ing back %s", username)

	if err != nil {
		log.Printf("yobackHandler: %s", err)
		return
	}

	defer resp.Body.Close()

	log.Printf("yobackHandler: %s", resp.Status)

	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		log.Printf("yobackHandler: %s", err)
	} else {
		log.Printf("yobackHandler: %s", string(body))
	}
}

func NewYoBackYoedClient() (*YoBackYoedClient, error) {
	c := &YoBackYoedClient{}
	config, err := c.loadConfig("./config.json")

	if err != nil {
		panic(fmt.Sprintf("failed loading config: %s", err))
	}

	c.config = config
	baseClient, err := clientInterface.NewBaseYoedClient()

	if err != nil {
		return nil, err
	}
	c.BaseYoedClient = *baseClient

	return c, nil
}

func main() {
	c, _ := NewYoBackYoedClient()

	clientInterface.Run(c)
}