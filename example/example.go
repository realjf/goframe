package main

import "github.com/realjf/goframe"

func main() {
	server := goframe.NewServer("./config/config.yaml")
	server.Run()
}
