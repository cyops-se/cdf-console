package devices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"server/logger"
	"server/messages"
)

func LoginRESTAPI(ip string, username string, password string) messages.LoginRequestResponse {

	var response messages.LoginRequestResponse
	response.Success = false

	body := &messages.LoginRequest{Username: username, Password: password}
	data, _ := json.Marshal(body)
	bodyReader := bytes.NewReader(data)

	request, err := http.NewRequest("POST", fmt.Sprintf("http://%s/api/login", ip), bodyReader)
	if err != nil {
		logger.Error("LoginRESTAPI", "login failed, newrequest returned error: %s", err.Error())
		return response
	}

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		logger.Error("LoginRESTAPI", "login failed, http POST returned error: %s", err.Error())
		return response
	}

	result, _ := io.ReadAll(resp.Body)
	json.Unmarshal(result, &response)

	return response
}
