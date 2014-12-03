package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
)

func serveCmd(args []string) error {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	portFlag = fs.String("p", "8080", "the port to listen on")
	fs.Parse(args)

	handler, err := NewZenBag("zens.txt")
	if err != nil {
		errLog.Println("Error reading zens.txt:", err)
		die()
	}
	mux := http.NewServeMux()

	mux.Handle("/zen", handler)
	mux.HandleFunc("/", missingHandler)

	regLog.Printf("Beginning listening on port %s\n", *portFlag)
	errLog.Fatal(http.ListenAndServe("localhost:"+*portFlag, mux))

	return nil
}

// ZensBag is the holder of all our zens.
type ZensBag struct {
	// This is where we keep all our zens
	Messages []string
}

// NewZenBag creates a new ZenBag from the provided filename
func NewZenBag(file string) (*ZensBag, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return &ZensBag{strings.Split(string(bytes), "\n")}, nil
}

func (h *ZensBag) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg := h.Messages[rand.Intn(len(h.Messages))]

	// GET /zen ZEN HERE
	regLog.Printf(`%s %s "%s"`, r.Method, r.RequestURI, msg)

	// write the response out
	fmt.Fprintf(w, msg)
}

func missingHandler(w http.ResponseWriter, r *http.Request) {
	regLog.Printf("%s %s", r.Method, r.RequestURI)

	http.NotFound(w, r)
}
