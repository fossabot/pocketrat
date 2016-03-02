package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// The Configuration struct represents the server configuration.
type Configuration struct {
	ListenAddr string
	ListenPort uint16
	ServerName string
	GopherRoot string
}

// LoadConfiguration reads configuration from config.json and returns a populated Configuration struct.
func LoadConfiguration(path string) (config Configuration, err error) {
	text, readErr := ioutil.ReadFile(path)

	var c Configuration

	if readErr == nil {
		jsonErr := json.Unmarshal(text, &c)

		if jsonErr != nil {
			panic(jsonErr)
		}

		if c.ListenPort == 0 {
			c.ListenPort = 70
		}

		if c.GopherRoot == "" {
			c.GopherRoot, err = filepath.Abs(".")
			if err != nil {
				panic(err)
			}
		} else {
			c.GopherRoot, err = filepath.Abs(c.GopherRoot)
			if err != nil {
				panic(err)
			}
		}

		if c.ServerName == "" {
			hostName, err := os.Hostname()
			if err == nil {
				c.ServerName = hostName
			} else {
				panic(err)
			}
		}

		return c, nil
	}

	panic(err)
}
