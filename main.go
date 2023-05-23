package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

// Config represents the structure of the configuration file (config.yaml)
type Config struct {
	Rules []Rule `yaml:"rules"`
}

// Rule represents the redirection rule
type Rule struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}

// ProxyHandler redirects incoming requests based on the configuration rules
func ProxyHandler(config Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, rule := range config.Rules {
			sourceURL, err := url.Parse(rule.Source)
			if err != nil {
				log.Printf("Invalid source URL: %s\n", rule.Source)
				continue
			}

			if r.Host == sourceURL.Host {
				destinationURL, err := url.Parse(rule.Destination)
				if err != nil {
					log.Printf("Invalid destination URL: %s\n", rule.Destination)
					continue
				}

				destinationURL.Path = r.URL.Path
				destinationURL.RawQuery = r.URL.RawQuery

				log.Printf("Redirecting request from %s to %s\n", r.Host, destinationURL.String())
				http.Redirect(w, r, destinationURL.String(), http.StatusMovedPermanently)
				return
			}
		}
		log.Printf("No redirection rule found for host: %s\n", r.Host)
		http.NotFound(w, r)
	}
}

func main() {
	// Set up access and error logging
	accessLog, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open access log file: %v", err)
	}
	defer accessLog.Close()

	errorLog, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open error log file: %v", err)
	}
	defer errorLog.Close()

	// Set the log output to both the log files and stdout
	log.SetOutput(io.MultiWriter(accessLog, os.Stdout))
	log.SetErrorOutput(io.MultiWriter(errorLog, os.Stderr))

	// Read the configuration file
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	// Parse the configuration
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	// Create a new HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(ProxyHandler(config)),
		ErrorLog: log.New(errorLog, "", 0),
	}

	// Start the server
	log.Println("Proxy server started on port 8080")
	err = server.ListenAndServeTLS("cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("Failed to start proxy server: %v", err)
	}
}

