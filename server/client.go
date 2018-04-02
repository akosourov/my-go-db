package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	serverURL string
}

func NewClient(host, port string) *Client {
	serverURL := "http://" + host + ":" + port
	return &Client{
		serverURL: serverURL,
	}
}

func (c *Client) GetValue(key string) (string, error) {
	url := fmt.Sprintf("%s/storage/%s", c.serverURL, key)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		log.Printf("Error :%v\n", err.Error())
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error :%v\n", err.Error())
		return "", err
	}
	r := Response{}
	if err = json.Unmarshal(body, &r); err != nil {
		log.Printf("Error: %v\n", err.Error())
		return "", err
	}
	return string(body), nil
}

func (c *Client) SetValue(key, value string, ttl int64) (string, error) {
	req := Payload{
		ItemText: value,
		TTL:      ttl,
	}
	reqB, err := json.Marshal(&req)
	if err != nil {
		log.Printf("Error :%v\n", err.Error())
		return "", err
	}

	url := fmt.Sprintf("%s/storage/%s", c.serverURL, key)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqB))
	defer resp.Body.Close()
	if err != nil {
		log.Printf("Error :%v\n", err.Error())
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error :%v\n", err.Error())
		return "", err
	}
	r := Response{}
	if err = json.Unmarshal(body, &r); err != nil {
		log.Printf("Error: %v\n", err.Error())
		return "", err
	}
	return string(body), nil
}
