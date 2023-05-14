package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Config struct {
	Routes map[string]string `yaml:"routes"`
}

func main() {
	// Load configuration file
	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	config := Config{}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create reverse proxy
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			dest, ok := config.Routes[req.URL.Path]
			if ok {
				target, _ := url.Parse(dest)
				req.URL.Scheme = target.Scheme
				req.URL.Host = target.Host
				req.URL.Path = target.Path
			}
		},
	}

	// Start server
	http.HandleFunc("/", proxy.ServeHTTP)
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
