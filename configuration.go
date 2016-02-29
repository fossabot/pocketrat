package main

import (
    "encoding/json"
    "io/ioutil"
    "os"
)

// The Configuration ...
type Configuration struct {
    ListenAddr    string
    ListenPort    uint16
    ServerName    string
    GopherRoot    string
}

// LoadConfiguration ...
func LoadConfiguration(path string) (config Configuration, err error) {
    text, readErr := ioutil.ReadFile(path)
    
    var c Configuration
    
    if readErr == nil {
        jsonErr := json.Unmarshal(text, &c)
        
        if (jsonErr != nil) {
            panic(jsonErr)
        }
        
        if c.ListenPort == 0 {
            c.ListenPort = 70
        }
        
        if c.GopherRoot == "" {
            c.GopherRoot = "."
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