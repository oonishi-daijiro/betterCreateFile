package main

import (
	"files"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"prototype"

	"github.com/fatih/color"
)

var red *color.Color = color.New(color.FgRed).Add(color.Underline)

func main() {

	initialPrototypeDirPath := flag.String("init", "", "To set the path of prototype directory")
	prototypeName := flag.String("p", "", "The name of prototype.")
	flagShowPrototypeList := flag.Bool("list", false, "Flag whether to show prototype list")

	flag.Parse()

	if *initialPrototypeDirPath != "" {
		if err := prototype.Init(*initialPrototypeDirPath); err != nil {
			red.Println(err.Error())
			return
		}
	}

	if *flagShowPrototypeList {
		basePath, errGetBaseDir := prototype.GetBaseDirectoryPath()
		if errGetBaseDir != nil {
			red.Println(errGetBaseDir.Error())
			return
		}
		directories, errReadDir := os.ReadDir(basePath)
		if errReadDir != nil {
			red.Println(errReadDir.Error())
		}
		for _, dir := range directories {
			info, errGetInfo := dir.Info()
			if errGetInfo != nil {
				red.Println(errGetInfo.Error())
				return
			}
			if info.IsDir() {
				fmt.Println(info.Name())
			}
		}
		return
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
		if string(path[0]) == "_" {
			path = previousDir + string(path[1:])
		}
		if files.IsExist(path) {
			if !files.IsRequiredOverwirte(path) {
				continue
			}
		}
		previousDir = filepath.Dir(path)
		err := files.Create(path)
		if err != nil {
			red.Println(err.Error())
		}
	}
}
