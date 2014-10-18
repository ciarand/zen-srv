package main

import "fmt"

func usageCmd(args []string) error {
	fmt.Printf(`zen-srv: A tiny but inspirational web server

USAGE:
    zen-srv fetch [x]  # fetches x new zens and prints them to stdout
    zen-srv serve -p y # starts the server running on port y
    zen-srv help       # prints this help
`)

	return nil
}
