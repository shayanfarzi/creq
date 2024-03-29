package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var method string
var body string

func init() {
	const defaultMethod = "GET"
	flag.StringVar(&method, "method", defaultMethod, "HTTP method")
	flag.StringVar(&method, "m", defaultMethod, "HTTP method (shorthand)")
	flag.StringVar(&body, "body", "", "HTTP body (TODO: support json) under development")
	flag.StringVar(&body, "b", "", "HTTP body (shorthand) (TODO: support json) under development")
}

func main() {
	firstArg := os.Args[1]
	switch firstArg {
	case "-h":
		fmt.Print("\nSet url of request at the end of commands\n\n")
	case "--help":
		fmt.Print("\nSet url of request at the end of commands\n\n")
	case "help":
		fmt.Print("\nSet url of request at the end of commands\n\n")
		flag.Usage()
		return
	}
	flag.Parse()
	fmt.Println("body: ", body)
	jsonBody, err := json.Marshal(body)

	fmt.Println("After: ", string(jsonBody))

	lenOfArgs := len(os.Args)

	request := os.Args[lenOfArgs-1]

	if !strings.Contains(request, "http") {
		fmt.Println("url should contain http:// or https://")
		os.Exit(1)
	}

	resp, err := makeRequest(request, method, nil)

	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	fmt.Printf("Status: %v\n", resp.Status)
	fmt.Printf("URL: %v\n", resp.Request.URL)
	fmt.Println("Response: ")
	if resp.StatusCode == 200 {
		body := make([]byte, 4096*25)
		resp.Body.Read(body)
		fmt.Printf("%s\n", body)
	}
}

func makeRequest(request_path string, method string, body json.RawMessage) (*http.Response, error) {
	url_path, err := url.Parse(request_path)
	if err != nil {
		return nil, fmt.Errorf("sing par URL: %s", err)
	}

	var response *http.Response
	methodLower := strings.ToLower(method)
	switch methodLower {
	case "get":
		response, err = http.Get(url_path.String())
		if err != nil {
			return nil, fmt.Errorf("getting response from URL: %s", err)
		}
	case "delete":
		req, err := http.NewRequest("DELETE", url_path.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("creating DELETE request: %s", err)
		}
		client := http.Client{}

		response, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("getting DELETE response from URL: %s", err)
		}
	default:
		return nil, fmt.Errorf("unsupported method %s", methodLower)
	}

	return response, nil
}
