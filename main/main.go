package main

import (
	"buptncov"
	"fmt"
	"os"
)

func main() {
	username := os.Getenv("BUPT_USERNAME")
	password := os.Getenv("BUPT_PASSWORD")
	if username == "" {
		fmt.Println("未设置USERNAME")
		os.Exit(1)
	}
	if password == "" {
		fmt.Println("未设置PASSWORD")
		os.Exit(1)
	}
	user := buptncov.New()
	err1 := user.Login(username, password)
	if err1 != nil {
		fmt.Println("error:", err1)
		os.Exit(1)
	}
	err2 := user.GetAndPostInfo()
	if err2 != nil {
		fmt.Println("error:", err2)
		os.Exit(1)
	}
}
