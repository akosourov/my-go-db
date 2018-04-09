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

func (c *Client) GetValue(key string) (interface{}, error) {
	url := fmt.Sprintf("%s/%s", c.storageURL, key)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error :%v\n", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	log.Println("Response:", string(respBytes))
	if err != nil {
		log.Printf("Error :%v\n", err.Error())
		return nil, err
	}

	respBody := ResponseBody{}
	if err = json.Unmarshal(respBytes, &respBody); err != nil {
		log.Printf("Error: %v\n", err.Error())
		return nil, err
	}

	if respBody.String != "" {
		return respBody.String, nil
	}
	if respBody.Int > 0 {
		return respBody.Int, nil
	}
	if respBody.StringList != nil {
		return respBody.StringList, nil
	}
	if respBody.IntList != nil {
		return respBody.IntList, nil
	}
	if respBody.StringDict != nil {
		return respBody.StringDict, nil
	}
	if respBody.IntDict != nil {
		return respBody.IntDict, nil
	}
	return nil, nil
}


func (c *Client) SetInt(key string, value, ttl int) error {
	reqBody := new(RequestBody)
	reqBody.Int = value
	reqBody.TTL = ttl
	_, err := c.doPost(key, reqBody)
	return err
}

func (c *Client) SetString(key, value string, ttl int) error {
	reqBody := new(RequestBody)
	reqBody.String = value
	reqBody.TTL = ttl
	_, err := c.doPost(key, reqBody)
	return err
}

func (c *Client) SetIntSlice(key string, value []int, ttl int) error {
	reqBody := new(RequestBody)
	reqBody.IntList = value
	reqBody.TTL = ttl
	_, err := c.doPost(key, reqBody)
	return err
}

func (c *Client) SetStringSlice(key string, value []string, ttl int) error {
	reqBody := new(RequestBody)
	reqBody.StringList = value
	reqBody.TTL = ttl
	_, err := c.doPost(key, reqBody)
	return err
}

func (c *Client) SetStringMap(key string, value map[string]string, ttl int) error {
	reqBody := new(RequestBody)
	reqBody.StringDict = value
	reqBody.TTL = ttl
	_, err := c.doPost(key, reqBody)
	return err
}

func (c *Client) SetIntMap(key string, value map[string]int, ttl int) error {
	reqBody := new(RequestBody)
	reqBody.IntDict = value
	reqBody.TTL = ttl
	_, err := c.doPost(key, reqBody)
	return err
}

func (c *Client) GetKeys() []string {
	resp, err := http.Get(c.storageURL + "/")
	if err != nil {
		log.Println("GetKeys() error: ", err.Error())
		return nil
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	log.Printf("respBytes: %v", string(respBytes))
	if err != nil {
		log.Println("ioutil.ReadAll(resp.Body) error: ", err.Error())
		return nil
	}

	respBody := new(ResponseBody)
	if err = json.Unmarshal(respBytes, respBody); err != nil {
		log.Println("Unmarhal error", err.Error())
		return nil
	}
	return respBody.Keys
}

func (c *Client) Remove(key string) error {
	url := c.getKeyUrl(key)

	client := new(http.Client)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Println("Remove error", err.Error())
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Remove error", err.Error())
		return err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Remove error", err.Error())
		return err
	}

	log.Println("Response:", string(respBytes))

	respBody := new(ResponseBody)
	if err := json.Unmarshal(respBytes, respBody); err != nil {
		log.Println("Remove error", err.Error())
		return err
	}
	if !respBody.Success {
		return fmt.Errorf(respBody.Message)
	}
	return nil
}

func (c *Client) getKeyUrl(key string) string {
	return fmt.Sprintf("%s/%s", c.storageURL, key)
}

func (c *Client) doPost(key string, reqBody *RequestBody) (*ResponseBody, error) {
	reqBytes, err := json.Marshal(reqBody)
	log.Printf("RequestBody: %v", string(reqBytes))
	if err != nil {
		log.Println("Marshal error", err.Error())
		return nil, err
	}

	url := c.getKeyUrl(key)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		log.Println("doPost error:", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ioutil.ReadAll error:", err.Error())
		return nil, err
	}

	respBody := new(ResponseBody)
	log.Printf("ResponseBody :%v", string(respBytes))
	if err = json.Unmarshal(respBytes, respBody); err != nil {
		log.Println("Unmarhsal error:", err.Error())
		return nil, err
	}

	return respBody, nil
}

