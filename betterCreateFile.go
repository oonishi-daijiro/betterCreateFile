package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fatih/color"
)

type configFromJSON struct {
	PrototypePath string `json:"prototypePath"`
}

var red *color.Color = color.New(color.FgRed).Add(color.Underline)

func main() {
	exePath, exePathErr := os.Executable()
	if exePathErr != nil {
		red.Println("Error :", exePathErr.Error())
		return
	}
	exePath = filepath.Dir(exePath)
	initPrototypeDir := flag.String("init", "", "To set the path of prototype directory")
	prototypeName := flag.String("p", "", "The name of prototype.")
	flag.Parse()
	if !isExist(exePath+"\\config.json") || *initPrototypeDir != "" {
		if initErr := initConfig(*initPrototypeDir, exePath); initErr != nil {
			red.Println("Error :", initErr)
			return
		}
		return
	}
	prototypePath, getErr := getProtoTypeDir(exePath)
	if getErr != nil {
		red.Println("Error :", getErr)
		return
	}
	requirePrptotypePath := prototypePath + "\\" + filepath.Base(*prototypeName)
	if !isExist(requirePrptotypePath) {
		fmt.Println("No such prototype.")
		return
	}

	if *prototypeName == "" && len(flag.Args()) == 0 && *initPrototypeDir == "" {
		fmt.Println("Please set argument.")
		return
	}
	if len(flag.Args()) != 0 {
		if err := createFileAndDir(flag.Arg(0)); err != nil {
			red.Println("Error :", err)
			return
		}
		return
	}
	filePath := make([]string, 0)
	dirPath := make([]string, 0)
	readDirStruct(requirePrptotypePath, &filePath, &dirPath, prototypePath)
	if len(filePath) == 0 && len(dirPath) != 0 {
		createPrototypesEachDir(&dirPath, ".\\", prototypeName)
		return
	}
	if createDirErr := createPrototypesEachDir(&dirPath, ".\\", prototypeName); createDirErr != nil {
		red.Println(createDirErr.Error())
		return
	}
	waitCopying := new(sync.WaitGroup)
	limitRoutine := make(chan struct{}, 100)
	for _, c := range filePath {
		waitCopying.Add(1)
		go func(wg *sync.WaitGroup, c string) {
			limitRoutine <- struct{}{}
			cpErr := copyFile(c, prototypeName, prototypePath)
			if cpErr != nil {
				red.Println(cpErr)
				return
			}
			<-limitRoutine
			wg.Done()
		}(waitCopying, c)
	}
	waitCopying.Wait()
}

func getProtoTypeDir(exePath string) (string, error) {
	if tf := isExist(exePath + "\\config.json"); !tf {
		notExist := errors.New("Cannot find the config.json")
		return "", notExist
	}
	raw, err := ioutil.ReadFile(exePath + "\\config.json")
	if err != nil {
		return "", err
	}
	var config configFromJSON
	json.Unmarshal(raw, &config)
	return config.PrototypePath, nil
}

func initConfig(specifiedPath string, exePath string) error {
	if specifiedPath == "" {
		var ans string
		fmt.Println("Please input the path of prototype directory")
		fmt.Scan(&ans)
		if ans == "" {
			initConfig("", exePath)
		}
		specifiedPath = ans
	}
	specifiedPath = strings.Replace(specifiedPath, "\\", "/", -1)
	if !isExist(specifiedPath) {
		notFnd := errors.New("No such directory or incorrect path")
		return notFnd
	}
	fileContent := "{\n\"prototypePath\":" + "\"" + specifiedPath + "\"" + "\n}"
	file, err := os.Create(exePath + "\\config.json")
	if err != nil {
		return err
	}
	file.Write([]byte(fileContent))
	return nil
}

func createPrototypesEachDir(path *[]string, target string, prototypeName *string) error {
	for _, c := range *path {
		relPathForCrruentDir, relErr := filepath.Rel(*prototypeName, c)
		if relErr != nil {
			return relErr
		}
		if err := os.MkdirAll(relPathForCrruentDir, 0777); err != nil {
			red.Println("Error: ", err.Error())
			return err
		}
	}
	return nil
}

func copyFile(distTo string, copyFrom *string, protoPath string) error {
	relPathForCrruentDir, relErr := filepath.Rel(*copyFrom, distTo)
	if relErr != nil {
		return relErr
	}
	file, createErr := os.Create(relPathForCrruentDir)
	src, openErr := os.Open(protoPath + "\\" + distTo)
	if createErr != nil {
		return createErr
	}
	if openErr != nil {
		return openErr
	}
	_, copyErr := io.Copy(file, src)
	if copyErr != copyErr {
		return copyErr
	}
	return nil
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func readDirStruct(path string, filePath *[]string, dirPath *[]string, protoPath string) {
	contents, _ := ioutil.ReadDir(path)
	for _, c := range contents {
		if c.IsDir() {
			rel, relErr := filepath.Rel(protoPath, path+"\\"+c.Name())
			if relErr != nil {
				red.Println("Error: ", relErr.Error())
				return
			}
			readDirStruct(path+"\\"+c.Name(), filePath, dirPath, protoPath)
			*dirPath = append(*dirPath, rel)
			continue
		}
		rel, relErr := filepath.Rel(protoPath, path+"\\"+c.Name())
		if relErr != nil {
			red.Println(relErr.Error())
			return
		}
		*filePath = append(*filePath, rel)
	}
}

func createFileAndDir(path string) error {
	if isExist(".\\" + path) {
		for {
			var ans string
			fmt.Println("\"", path, "\"", "is already exist. Overwrite? [y/n]")
			fmt.Scan(&ans)
			if ans == "y" {
				break
			} else if ans == "n" {
				return nil
			} else {
				continue
			}
		}
	}
	if filepath.Ext(".\\"+path) == "" && isExist(".\\"+path) {
		fStat, statErr := os.Stat(".\\" + path)
		if statErr != nil {
			return statErr
		}
		if fStat.IsDir() {
			if mkdirErr := os.MkdirAll(filepath.Dir(".\\"+path), 0666); mkdirErr != nil {
				return mkdirErr
			}
		} else {
			f, openErr := os.OpenFile(".\\"+path, os.O_CREATE, 0666)
			if openErr != nil {
				return openErr
			}
			c := make([]byte, 0)
			f.Write(c)
			return nil
		}
	}
	if mkdirErr := os.MkdirAll(filepath.Dir(".\\"+path), 0666); mkdirErr != nil {
		return mkdirErr
	}
	f, openErr := os.OpenFile(".\\"+path, os.O_CREATE, 0666)
	if openErr != nil {
		return openErr
	}
	c := make([]byte, 0)
	f.Write(c)
	return nil
}
