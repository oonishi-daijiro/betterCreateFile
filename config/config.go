package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type JsonTypes interface {
	int | bool | string | []string | []int | []bool
}

type Config[T JsonTypes] struct {
	path     string
	rawJson  map[string]T
	fileName string
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func Init[T JsonTypes](fileName string) (Config[T], error) {
	exePath, exePathErr := os.Executable()
	if exePathErr != nil {
		return Config[T]{}, exePathErr
	}
	exePath = filepath.Dir(exePath)

	if !isExist(exePath + "\\" + fileName) {
		if initErr := initFile(exePath + "\\" + fileName); initErr != nil {
			return Config[T]{}, initErr
		}
	}
	config := Config[T]{}
	config.path = exePath + "\\" + fileName
	config.rawJson = make(map[string]T)
	config.fileName = fileName
	return config, nil
}

func initFile(path string) error {
	_, err := os.Create(path)
	if err != nil {
		return err
	}
	return nil
}

func (p *Config[T]) Get(key string) (T, error) {
	c, err := os.ReadFile(p.path)
	if err != nil {
		var i T
		return i, err
	}
	json.Unmarshal(c, &p.rawJson)
	return p.rawJson[key], nil
}

func (p *Config[T]) Set(key string, value T) error {
	p.rawJson[key] = value
	c, err := json.Marshal(&p.rawJson)
	if err != nil {
		return err
	}
	if err := os.WriteFile(p.path, c, 0664); err != nil {
		return err
	}
	return nil
}
