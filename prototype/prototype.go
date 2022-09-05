package prototype

import (
	"config"
	"errors"
	"files"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func isExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func Init(prototypePath string) error {
	prototypePath, err := filepath.Abs(prototypePath)
	if err != nil {
		return err
	}
	config, err := config.Init[string]()
	if err != nil {
		return err
	}
	if err := checkProtrotypePath(prototypePath); err != nil {
		return err
	}
	config.Set("prototype", prototypePath)
	return nil
}

func getBaseDirPath() (string, error) {
	config, err := config.Init[string]()
	if err != nil {
		return "", err
	}
	path, err := config.Get("prototype")
	if err != nil {
		return "", err
	}
	return path, nil
}

func Create(prototypeName string) error {
	base, err := getBaseDirPath()
	if err != nil {
		return err
	}
	if !isExist(filepath.Join(base, prototypeName)) {
		return errors.New("no such prototype")
	}
	dir, err := files.ReadDirStructure(filepath.Join(base, prototypeName))
	if err != nil {
		return err
	}

	for _, path := range dir.Files {
		files.Create(path)
		err := files.CopyFile(path, filepath.Join(base, prototypeName, path))
		if err != nil {
			return err
		}
	}
	return nil
}

func checkProtrotypePath(specifiedPath string) error {
	if specifiedPath == "" {
		var ans string
		fmt.Println("Please input the path of prototype directory")
		fmt.Scan(&ans)
		if ans == "" {
			checkProtrotypePath("")
		}
		specifiedPath = ans
	}
	specifiedPath = strings.Replace(specifiedPath, "\\", "/", -1)
	if !isExist(specifiedPath) {
		notFnd := errors.New("no such directory or incorrect path")
		return notFnd
	}
	return nil
}
