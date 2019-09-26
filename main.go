package main

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/linkpoolio/bridges/bridge"
)

// Scorechain is the Scorechain Ethereum API
type Scorechain struct {
	url *url.URL
	key string
}

// Opts is the options for the Scorechain bridge
func (sc *Scorechain) Opts() *bridge.Opts {
	return &bridge.Opts{
		Name:   "Scorechain",
		Lambda: true,
	}
}

// Run is the main Scorechain Ethereum API adapter implementation
func (sc *Scorechain) Run(h *bridge.Helper) (interface{}, error) {
	switch h.GetParam("endpoint") {
	case "get-status":
		r := &Status{}
		e := h.HTTPCall(http.MethodGet, sc.url.String()+"/status", r)
		return r, e
	case "get-trx":
		hash := h.GetParam("hash")
		if "" == hash {
			return nil, errors.New("get-transaction endpoint requires hash parameter")
		}

		r := &Transaction{}
		e := h.HTTPCall(http.MethodGet, sc.url.String()+"/tx/"+h.GetParam("hash"), r)
		return r, e
	default:
		return nil, errors.New("unrecognized or unsupported Scorechain Ethereum API endpoint")
	}
}

func getScorechainURL() *url.URL {
	rawURL, ok := os.LookupEnv("SC_ETH_URL")
	if !ok {
		rawURL = "https://ethereum.scorechain.com"
	}

	url, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal(err.Error())
	}

	return url
}

func main() {
	bridge.NewServer(&Scorechain{getScorechainURL(), os.Getenv("SC_ETH_TOKEN")}).Start(8080)
}

// Block is an Ethereum block
type Block struct {
	hash   string
	height int
	nbtx   int
	date   string
}

// Status is the status of the Scorechain Ethereum API
type Status struct {
	Timestamp      int
	State          string
	UnsyncedBlocks int
	LastBlock      Block
}

// Transaction is an Ethereum transaction
type Transaction struct {
}
