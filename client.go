package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/ajg/form"
)

// Version represents the version of the CLI
const Version = "0.0.1"

// Client is the means of connecting to the Civo API service
type Client struct {
	BaseURL   *url.URL
	UserAgent string
	APIKey    string

	httpClient *http.Client
}

// HTTPError is the error returned when the API fails with an HTTP error
type HTTPError struct {
	Code   int
	Status string
}

// SimpleResponse is a structure that returns success and/or any error
type SimpleResponse struct {
	Result       string `json:"result"`
	ErrorCode    string `json:"code"`
	ErrorReason  string `json:"reason"`
	ErrorDetails string `json:"details"`
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Status)
}

// NewClientWithURL initializes a Client with a specific API URL
func NewClientWithURL(apiKey string, civoAPIURL string) (*Client, error) {
	parsedURL, err := url.Parse(civoAPIURL)
	if err != nil {
		return nil, err
	}

	client := &Client{
		BaseURL:    parsedURL,
		UserAgent:  "civogo/" + Version,
		APIKey:     apiKey,
		httpClient: &http.Client{},
	}
	return client, nil
}

// NewClient initializes a Client connecting to the production API
func NewClient(apiKey string) (*Client, error) {
	return NewClientWithURL(apiKey, "https://api.civo.com")
}

// NewAdvancedClientForTesting initializes a Client connecting to a local test server and allows for specifying methods
func NewAdvancedClientForTesting(responses map[string]map[string]string) (*Client, *httptest.Server, error) {
	var responseSent bool

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			return
		}

		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		for url, criteria := range responses {
			if strings.Contains(req.URL.String(), url) &&
				req.Method == criteria["method"] {
				if criteria["method"] == "PUT" || criteria["method"] == "POST" || criteria["method"] == "PATCH" {

					if string(body) == criteria["requestBody"] {
						responseSent = true
						rw.Write([]byte(criteria["responseBody"]))
					}
				} else {
					responseSent = true
					rw.Write([]byte(criteria["responseBody"]))
				}
			}
		}

		if responseSent == false {
			fmt.Println("Failed to find a matching request!")
			fmt.Println("Request body:", string(body))
			fmt.Println("Method:", req.Method)
			fmt.Println("URL:", req.URL.String())
			rw.Write([]byte(`{"result": "failed to find a matching request"}`))
		}
	}))

	client, err := NewClientForTestingWithServer(server)

	return client, server, err
}

// NewClientForTesting initializes a Client connecting to a local test server
func NewClientForTesting(responses map[string]string) (*Client, *httptest.Server, error) {
	var responseSent bool

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		for url, response := range responses {
			if strings.Contains(req.URL.String(), url) {
				responseSent = true
				rw.Write([]byte(response))
			}
		}

		if responseSent == false {
			fmt.Println("Failed to find a matching request!")
			fmt.Println("URL:", req.URL.String())

			rw.Write([]byte(`{"result": "failed to find a matching request"}`))
		}
	}))

	client, err := NewClientForTestingWithServer(server)

	return client, server, err
}

// NewClientForTestingWithServer initializes a Client connecting to a passed-in local test server
func NewClientForTestingWithServer(server *httptest.Server) (*Client, error) {
	client, err := NewClientWithURL("TEST-API-KEY", server.URL)
	if err != nil {
		return nil, err
	}
	client.httpClient = server.Client()
	return client, err
}

func (c *Client) prepareClientURL(requestURL string) *url.URL {
	u, _ := url.Parse(c.BaseURL.String() + requestURL)
	return u
}

func (c *Client) sendRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", c.APIKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, HTTPError{Code: resp.StatusCode, Status: resp.Status}
	}

	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

// SendGetRequest sends a correctly authenticated get request to the API server
func (c *Client) SendGetRequest(requestURL string) ([]byte, error) {
	u := c.prepareClientURL(requestURL)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return c.sendRequest(req)
}

// SendPostRequest sends a correctly authenticated post request to the API server
func (c *Client) SendPostRequest(requestURL string, params interface{}) ([]byte, error) {
	u := c.prepareClientURL(requestURL)
	values, err := form.EncodeToValues(params)
	if err != nil {
		return nil, err
	}

	body := values.Encode()
	if body == "=" {
		body = ""
	}

	req, err := http.NewRequest("POST", u.String(), strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	return c.sendRequest(req)
}

// SendPutRequest sends a correctly authenticated put request to the API server
func (c *Client) SendPutRequest(requestURL string, params interface{}) ([]byte, error) {
	u := c.prepareClientURL(requestURL)
	values, err := form.EncodeToValues(params)
	if err != nil {
		return nil, err
	}

	body := values.Encode()
	if body == "=" {
		body = ""
	}

	req, err := http.NewRequest("PUT", u.String(), strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	return c.sendRequest(req)
}

// SendDeleteRequest sends a correctly authenticated delete request to the API server
func (c *Client) SendDeleteRequest(requestURL string) ([]byte, error) {
	u := c.prepareClientURL(requestURL)
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return c.sendRequest(req)
}

// DecodeSimpleResponse parses a response body in to a SimpleResponse object
func (c *Client) DecodeSimpleResponse(resp []byte) (*SimpleResponse, error) {
	response := SimpleResponse{}
	err := json.NewDecoder(bytes.NewReader(resp)).Decode(&response)
	return &response, err
}
