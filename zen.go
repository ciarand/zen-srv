package main

import (
	"flag"
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

var (
	portFlag   *string
	binaryName string
	errLog     *log.Logger
	regLog     *log.Logger
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	errLog = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime)
	regLog = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)

	// go run zen.go
	if len(os.Args) > 2 && os.Args[0] == "go" && os.Args[1] == "run" {
		binaryName = strings.Join(os.Args[0:2], " ")
	} else {
		binaryName = os.Args[0]
	}

	portFlag = flag.String("p", "8080", "the port to listen on")

	flag.Parse()
}

func main() {
	args := flag.Args()

	if len(args) == 0 {
		runServer()
		os.Exit(0)
	}

	if len(args) > 2 || args[0] != "fetch" {
		errLog.Println("unknown command structure:", strings.Join(args, " "))
		os.Exit(1)
	}

	// default is 1
	var num = 1
	var err error

	if len(args) == 2 {
		num, err = strconv.Atoi(args[1])
		if err != nil {
			errLog.Fatalf("couldn't translate %s into a number: %s", args[1], err)
		}
	}

	if num > 10 {
		errLog.Println("That number's pretty high. You sure?")
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
		errLog.Println("Error reading zens.txt:", err)
		os.Exit(1)
	}
	mux := http.NewServeMux()

	mux.Handle("/zen", handler)

	regLog.Printf("Beginning listening on port %s\n", *portFlag)
	errLog.Fatal(http.ListenAndServe("localhost:"+*portFlag, mux))

	return nil
}

// ZensBag is the holder of all our zens.
type ZensBag struct {
	// This is where we keep all our zens
	Messages []string
}

func (h *ZensBag) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg := h.Messages[rand.Intn(len(h.Messages))]

	// GET /zen ZEN HERE
	regLog.Printf(`%s %s "%s"`, r.Method, r.RequestURI, msg)

	// write the response out
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
	// 4 threads
	for i := 0; i < 4; i++ {
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
		os.Exit(2)
	}

	return string(body)
}
