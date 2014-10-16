package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal(runServer())
	}

	if len(os.Args) > 3 || os.Args[1] != "fetch" {
		log.Println("ERROR: unknown command structure:", strings.Join(os.Args, " "))
		os.Exit(1)
	}

	var num = 100
	var err error

	if len(os.Args) == 3 {
		num, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("ERROR: couldn't translate %s into a number: %s", os.Args[2], err)
		}
	}

	ch := make(chan string, num)

	go fetchZens(ch)

	for count := 0; count < num; count++ {
		fmt.Println(<-ch)
	}
}

func runServer() error {
	handler, err := NewZenBag("zens.txt")
	if err != nil {
		fmt.Print("Error reading zens.txt:", err)
		os.Exit(1)
	}
	mux := http.NewServeMux()

	mux.Handle("/zen", handler)

	log.Fatal(http.ListenAndServe("localhost:8080", mux))

	return nil
}

// ZensBag is the holder of all our zens.
type ZensBag struct {
	// This is where we keep all our zens
	Messages []string
}

func (h *ZensBag) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg := h.Messages[rand.Intn(len(h.Messages))]

	log.Println(msg)
	fmt.Fprintf(w, msg)
}

// NewZenBag creates a new ZenBag from the provided filename
func NewZenBag(file string) (*ZensBag, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return &ZensBag{strings.Split(string(bytes), "\n")}, nil
}

func fetchZens(out chan string) {
	// 10 threads
	for i := 0; i < 10; i++ {
		go func() {
			for {
				if zen := fetchZen(); zen != "" {
					out <- zen
				}
			}
		}()
	}
}

func fetchZen() string {
	resp, err := http.Get("https://api.github.com/zen")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ""
	}

	if resp.StatusCode != 200 {
		log.Fatal(string(body))
	}

	return string(body)
}
