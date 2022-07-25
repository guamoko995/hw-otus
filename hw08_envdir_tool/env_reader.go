package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	fileList, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed get file list into %s", dir)
	}

	pwd, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(pwd)

	fakeEnv := make(map[string]EnvValue)

	for _, file := range fileList {
		if !file.IsDir() {
			file, err := os.Open(file.Name())
			if err != nil {
				return nil, errors.Wrapf(err, "failed open file %q", file.Name())
			}
			defer file.Close()

			sc := bufio.NewScanner(file)
			if sc.Scan() {
				fakeEnv[file.Name()] = EnvValue{
					Value: strings.TrimRight(string(bytes.ReplaceAll(sc.Bytes(), []byte{0x00}, []byte("\n"))), " 	"),
				}
			} else {
				fakeEnv[file.Name()] = EnvValue{
					NeedRemove: true,
				}
			}
		}
	}

	return fakeEnv, nil
}

func EnvRead(fake Environment) []string {
	// дополнение карты окружения реальными переменными окружения
	for _, realEnv := range os.Environ() {
		keyVal := strings.SplitN(realEnv, "=", 2)
		key, val := keyVal[0], keyVal[1]
		if _, ok := fake[key]; ok {
			if fake[key].NeedRemove {
				delete(fake, key)
			}
		} else {
			fake[key] = EnvValue{
				Value: val,
			}
		}
	}

	// формирование полного фейкового окружения
	fakeEnv := make([]string, 0, len(fake))
	for name, val := range fake {
		fakeEnv = append(fakeEnv, fmt.Sprintf("%s=%s", name, val.Value))
	}
	return fakeEnv
}
