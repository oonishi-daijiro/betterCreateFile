package main

import (
	"files"
	"flag"
	"fmt"
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
		prototype.Create(*prototypeName)
	}

	for _, path := range flag.Args() {
		if files.IsExist(path) {
			if !files.IsRequiredOverwirte(path) {
				continue
			}
		}
		err := files.Create(path)
		if err != nil {
			red.Println(err.Error())
		}
	}
}
