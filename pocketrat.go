package main

import (
	"fmt"
)

// Version string representing the current version of PocketRat.
const Version = "0.1.1"

func main() {
	fmt.Println("PocketRat version", Version)
	config, err := LoadConfiguration("config.json")
	if err == nil {
		fmt.Println("listening on", config.ListenAddr, "port", config.ListenPort)
		fmt.Println("server name is", config.ServerName, "and gopher root is", config.GopherRoot)
		RunServer(config.ListenAddr, config.ListenPort, config.ServerName, config.GopherRoot)
	} else {
		panic(err)
	}
}
