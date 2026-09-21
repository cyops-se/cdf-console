package devices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"server/logger"
	"time"
)

var Token string

func Post(url string, body interface{}, t ...int) ([]byte, error) {

	data, _ := json.Marshal(body)
	bodyReader := bytes.NewReader(data)

	request, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		logger.Error("Post", "post failed, newrequest returned error: %s", err.Error())
		return nil, err
	}

	if Token != "" {
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", Token))
		logger.Trace("Post", "authorization header: %v", request.Header["Authorization"])
	}

	client := http.Client{}
	if len(t) > 0 {
		client.Timeout = time.Duration(t[0]) * time.Second
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	result, _ := io.ReadAll(resp.Body)
	return result, nil
}

func Get(url string, t ...int) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("Get", "get failed, newrequest returned error: %s", err.Error())
		return nil, err
	}

	if Token != "" {
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", Token))
		logger.Trace("Get", "authorization header: %v", request.Header["Authorization"])
	}

	client := http.Client{}
	if len(t) > 0 {
		client.Timeout = time.Duration(t[0]) * time.Millisecond
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	result, _ := io.ReadAll(resp.Body)
	return result, nil
}
