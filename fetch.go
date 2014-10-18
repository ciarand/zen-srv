package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func fetchCmd(args []string) error {
	fs := flag.NewFlagSet("fetch", flag.ExitOnError)

	workers := fs.Int("w", 4, "the number of worker threads to use")
	number := fs.Int("n", 1, "the number of zens to fetch")
	warn := fs.Bool("warn", true, "when provided zen-srv will warn about rudeness")
	delayStr := fs.String("d", "1s", "the delay between fetch requests (e.g. 1s, 300ms, 2h45m)")

	fs.Parse(args)

	if *number > 10 && *warn {
		errLog.Println("That number's pretty high. You sure?")
		die()
	}

	delay, err := time.ParseDuration(*delayStr)
	if err != nil {
		errLog.Printf("Unable to parse duration string (%s): %s\n", *delayStr, err)
	}

	ch := make(chan string, *number)

	go fetchZens(ch, *workers, delay)

	for count := 0; count < *number; count++ {
		fmt.Println(<-ch)
	}

	return nil
}

func fetchZens(out chan string, numWorkers int, delay time.Duration) {
	// 4 threads
	for i := 0; i < numWorkers; i++ {
		go func() {
			for {
				if zen := fetchZen(); zen != "" {
					out <- zen
				}

				time.Sleep(delay)
			}
		}()
	}
}

func fetchZen() string {
	resp, err := http.Get("https://api.github.com/zen")
	if err != nil {
		errLog.Print(err)
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errLog.Print(err)
		return ""
	}

	if resp.StatusCode != 200 {
		errLog.Print(string(body))
		die()
	}

	return string(body)
}
