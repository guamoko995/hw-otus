package main

import (
	"log"
	"os"
)

func main() {
	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatalln(err)
		return
	}
	fullEnv := EnvRead(env)
	os.Exit(RunCmd(os.Args[2:], fullEnv))
}
