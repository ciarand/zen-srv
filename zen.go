package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	portFlag   *string
	binaryName string
	errLog     *log.Logger
	regLog     *log.Logger
	commands   map[string]executer
)

type executer func([]string) error

func init() {
	commands = map[string]executer{
		"fetch": fetchCmd,
		"serve": serveCmd,
		"help":  usageCmd,
	}

	flag.Usage = func() {
		usageCmd([]string{})

		exit()
	}

	rand.Seed(time.Now().UTC().UnixNano())

	errLog = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime)
	regLog = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		usageCmd(args)
		exit()
	}

	cmd := commands[args[0]]
	if cmd == nil {
		errLog.Println("command not found")
		usageCmd(args)
		die()
	}

	if err := cmd(args[1:]); err != nil {
		errLog.Println(err)
		die()
	}

	exit()
}

func exit() {
	os.Exit(0)
}

func die() {
	os.Exit(1)
}
