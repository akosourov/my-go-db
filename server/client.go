package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Client struct {
	serverURL   string
	storageURL  string
}

func NewClient(host, port string) *Client {
	serverURL := fmt.Sprintf("http://%s:%s", host, port)
	return &Client{
		serverURL: serverURL,
		storageURL: serverURL + "/storage",
	}
}


//func (c *Client) Execute


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
	r := ResponseBody{}
	if err = json.Unmarshal(body, &r); err != nil {
		log.Printf("Error: %v\n", err.Error())
		return "", err
	}
	if r.ValueStr != "" {
		return r.ValueStr, nil
	} else if r.ValueInt > 0 {
		return strconv.Itoa(r.ValueInt), nil
	}
	return r.Message, nil
}

func (c *Client) SetValue(key, value string, ttl int) (string, error) {
	req := new(RequestBody)

	// is int?
	fmt.Printf("value :%v", value)
	if valueInt, err := strconv.Atoi(value); err == nil {
		req.ValueInt = valueInt
	} else {
		req.ValueStr = value
	}

	// is string


	reqB, err := json.Marshal(&req)
	if err != nil {
		log.Printf("Error :%v\n", err.Error())
		return "", err
	}

	url := fmt.Sprintf("%s/%s", c.storageURL, key)
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
	r := ResponseBody{}
	if err = json.Unmarshal(body, &r); err != nil {
		log.Printf("Error: %v\n", err.Error())
		return "", err
	}
	if r.Success {
		return "", nil
	}
	return r.Message, nil
}
