package main

import (
	"files"
	"flag"
	"fmt"
	"path/filepath"
	"prototype"

	"github.com/fatih/color"
)

var red *color.Color = color.New(color.FgRed).Add(color.Underline)

func main() {

	initialPrototypeDirPath := flag.String("init", "", "To set the path of prototype directory")
	prototypeName := flag.String("p", "", "The name of prototype.")

	flag.Parse()

	if *initialPrototypeDirPath != "" {
		if err := prototype.Init(*initialPrototypeDirPath); err != nil {
			red.Println(err.Error())
			return
		}
	}

	if *prototypeName == "" && len(flag.Args()) == 0 && *initialPrototypeDirPath == "" {
		fmt.Println("Please set argument.")
		return
	}

	if *prototypeName != "" {
		if err := prototype.Create(*prototypeName); err != nil {
			red.Println(err)
			return
		}
	}
	previousDir := ""
	for _, path := range flag.Args() {
		if files.IsExist(path) {
			if !files.IsRequiredOverwirte(path) {
				continue
			}
		}
		if string(path[0]) == "_" {
			path = previousDir + string(path[1:])
		}
		previousDir = filepath.Dir(path)
		err := files.Create(path)
		if err != nil {
			red.Println(err.Error())
		}
	}
}
