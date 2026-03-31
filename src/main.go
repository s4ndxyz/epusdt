package main

import (
	"github.com/assimon/luuu/command"
	"github.com/gookit/color"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			color.Error.Println("[Start Server Err!!!] ", err)
		}
	}()
	if err := command.Execute(); err != nil {
		panic(err)
	}
}
