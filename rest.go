package godevino

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	DEFAULT_URL      = "https://integrationapi.net/rest"
	DEFAULT_USERNAME = ""
	DEFAULT_PASSWORD = ""
)

var (
	sessionId string                    // Session ID
	Url       string = DEFAULT_URL      // API Url
	Username  string = DEFAULT_USERNAME // Username
	Password  string = DEFAULT_PASSWORD // Password

	err error
)

type Client string

func urlencode(data map[string]interface{}) string {
	var buf bytes.Buffer
	for k, v := range data {
		if k == "destinationAddresses" {
			for _, d := range data["destinationAddresses"].([]string) {
				buf.WriteString(url.QueryEscape(k))
				buf.WriteByte('=')
				buf.WriteString(fmt.Sprintf("%s", d))
				buf.WriteByte('&')
			}
		} else {
			buf.WriteString(url.QueryEscape(k))
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(fmt.Sprintf("%v", v)))
			buf.WriteByte('&')
		}
	}

	mData := buf.String()

	return mData[0 : len(mData)-1]
}

func Initialize() error {
	sessionId, err = getSessionId()
	if err != nil {
		return err
	}

	return nil
}

func getSessionId() (string, error) {
	var params = map[string]interface{}{
		"login":    Username,
		"password": Password,
	}
	dst := fmt.Sprintf("%s%s%s", Url, "/user/sessionId?", urlencode(params))
	req, err := http.NewRequest("GET", dst, bytes.NewBuffer([]byte("")))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www/form-urlencoded; charset=utf-8")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	return string(body)[1 : len(body)-1], nil
}

func (c *Client) GetBalance() (string, error) {
	var params = map[string]interface{}{}

	balance, err := request("GET", "/user/balance?", params)
	if err != nil {
		return "", err
	}

	return balance, nil
}

func (c *Client) SendMessage(src_address, dst_address, data string, validity int, send_date_utc string) (string, error) {
	var params = map[string]interface{}{
		"sourceAddress":      src_address,
		"destinationAddress": dst_address,
		"data":               data,
		"validity":           validity,
		"sendDate":           send_date_utc,
	}

	req, err := request("POST", "/sms/send?", params)
	if err != nil {
		return "", err
	}

	return req, nil
}

func (c *Client) SendMessageBulk(src_address string, dst_address []string, data string, validity int, send_date_utc string) (string, error) {
	var params = map[string]interface{}{
		"sourceAddress":        src_address,
		"destinationAddresses": dst_address,
		"data":                 data,
		"validity":             validity,
		"sendDate":             send_date_utc,
	}

	req, err := request("POST", "/sms/sendbulk?", params)
	if err != nil {
		return "", err
	}

	return req, nil
}

func (c *Client) SendMessageByTimezone(src_address, dst_address, data string, validity int, send_date string) (string, error) {
	var params = map[string]interface{}{
		"sourceAddress":      src_address,
		"destinationAddress": dst_address,
		"data":               data,
		"validity":           validity,
		"sendDate":           send_date,
	}

	req, err := request("POST", "/sms/sendbytimezone?", params)
	if err != nil {
		return "", err
	}

	return req, nil
}

func (c *Client) GetMessageState(message_id string) (string, error) {
	var params = map[string]interface{}{
		"messageId": message_id,
	}

	req, err := request("GET", "/sms/state?", params)
	if err != nil {
		return "", err
	}

	return req, nil
}

func (c *Client) GetStatistics(start_date, end_date string) (string, error) {
	var params = map[string]interface{}{
		"startDateTime": start_date,
		"endDateTime":   end_date,
	}

	req, err := request("GET", "/sms/statistics?", params)
	if err != nil {
		return "", err
	}

	return req, nil
}

func (c *Client) GetIncomingMessages(start_date_utc, end_date_utc string) (string, error) {
	var params = map[string]interface{}{
		"minDateUTC": start_date_utc,
		"maxDateUTC": end_date_utc,
	}

	req, err := request("GET", "/sms/in?", params)
	if err != nil {
		return "", err
	}

	return req, nil
}

func request(method, path string, params map[string]interface{}) (string, error) {
RequestStart:
	params["sessionId"] = sessionId
	dst := fmt.Sprintf("%s%s%s", Url, path, urlencode(params))
	// fmt.Printf("[INF] Request: %s\n", dst)
	req, err := http.NewRequest(method, dst, bytes.NewBuffer([]byte("")))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode == 401 {
		sessionId, err = getSessionId()
		if err != nil {
			return "", err
		}
		response.Body.Close()

		goto RequestStart
	}

	body, _ := ioutil.ReadAll(response.Body)

	return string(body), nil
}
