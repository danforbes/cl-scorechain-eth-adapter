package main

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/linkpoolio/bridges/bridge"
)

func init() {
	godotenv.Load()
}

// Scorechain is the Scorechain Ethereum API
type Scorechain struct {
	url   *url.URL
	token string
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
	case "get-scoring":
		a := h.GetParam("address")
		d := h.GetParam("direction")
		if "" == a || "" == d {
			return nil, errors.New("get-scoring endpoint requires address and direction parameters")
		}

		depth := 3
		depthStr := h.GetParam("depth")
		if "" != depthStr {
			depth, _ = strconv.Atoi(depthStr)
		}

		r := &Scoring{}
		e := h.HTTPCall(http.MethodGet, sc.url.String()+"/scoring/address/"+a+"/"+d+"?depth="+strconv.Itoa(depth)+"&pretty=false&token="+sc.token, r)
		return r, e
	case "get-status":
		r := &Status{}
		e := h.HTTPCall(http.MethodGet, sc.url.String()+"/status?token="+sc.token, r)
		return r, e
	case "get-trx":
		hash := h.GetParam("hash")
		if "" == hash {
			return nil, errors.New("get-transaction endpoint requires hash parameter")
		}

		r := &Transaction{}
		e := h.HTTPCall(http.MethodGet, sc.url.String()+"/tx/"+h.GetParam("hash")+"?token="+sc.token, r)
		return r, e
	default:
		return nil, errors.New("unrecognized or unsupported Scorechain Ethereum API endpoint")
	}
}

func getScorechainURL() *url.URL {
	rawURL, ok := os.LookupEnv("SC_ETH_URL")
	if !ok {
		rawURL = "https://api.ethereum.scorechain.com"
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
	Hash   string
	Height int64
}

// StatusBlock is a block from the Scorechain Ethereum API status endpoint
type StatusBlock struct {
	Block
	Nbtx int64
	Date string
}

// Status is the status of the Scorechain Ethereum API
type Status struct {
	Success bool
	Result  StatusResult
}

// State is the state of the Scoreschain Ethereum API
type State string

const (
	// Synced is a Scorechain Ethereum API that is in sync
	Synced State = "synced"

	// Syncing is a Scorechain Ethereum API that is syncing
	Syncing State = "syncing"

	// OutOfSync is a Scorechain Ethereum API that is out-of-sync
	OutOfSync State = "out-of-sync"
)

// StatusResult is the result of a call to the Scorechain Ethereum API status endpoint
type StatusResult struct {
	Timestamp      int64
	State          State
	UnsyncedBlocks int
	LastBlock      StatusBlock
}

// Transaction is an Ethereum transaction
type Transaction struct {
	Success bool
	Result  TransactionResult
}

// TransactionResult is the result of a call to the Scorechain Ethereum API transaction endpoint
type TransactionResult struct {
	Hash              string
	Block             Block
	From              TransactionSender
	To                TransactionReceiver
	Gas               Gas
	Value             Values
	Timestamp         int64
	Date              string
	Confirmations     int
	InternalTransfers []Transfer
	TokenTransfers    []TokenTransfer
}

// TransactionReceiver is the receiver of an Ethereum transaction
type TransactionReceiver TransactionSender

// TransactionSender is the sender of an Ethereum transaction
type TransactionSender struct {
	Address    string
	Label      string
	Type       AddressType
	IsContract bool
}

// AddressType is the type of an Ethereum address
type AddressType struct {
	ID       int
	Label    string
	Score    int
	Custom   bool
	ParentID int
}

// Gas is gas used for an Ethereum transaction
type Gas struct {
	Price       int
	Quantity    int
	Used        int
	UsedPercent float64
	Cost        float64
}

// Values is values in different currencies
type Values struct {
	ETH int
	AUD float64
	BRL float64
	CAD float64
	CHF float64
	CNY float64
	EUR float64
	GBP float64
	HKD float64
	ILS float64
	JPY float64
	KRW float64
	MXN float64
	NOK float64
	NZD float64
	PLN float64
	RUB float64
	SEK float64
	SGD float64
	TRY float64
	USD float64
	BTC float64
}

// Transfer is an Ethereum transfer
type Transfer struct {
	From  TransactionSender
	To    TransactionReceiver
	Value Values
}

// TokenTransfer is an Ethereum token transfer
type TokenTransfer struct {
	Transfer
	CounterValues Values
}

// Token is an Ethereum token
type Token struct {
	Address string
	Name    string
	Symbol  string
}

// Direction is the direction of scoring in the Scorechain Ethereum API
type Direction string

const (
	// Incoming is incoming scoring
	Incoming Direction = "incoming"

	// Outgoing is outgoing scoring
	Outgoing Direction = "outgoing"
)

// Scoring is the scoring of an address from the Scorechain Ethereum API
type Scoring struct {
	Success bool
	Result  ScoringResult
}

// ScoringResult is the result of a call to the Scorechain Ethereum API scoring endpoint
type ScoringResult struct {
	Scx     int
	Details []ScoringDetail
}

// ScoringDetail is the detail of a score
type ScoringDetail struct {
	Address    string
	Amount     float64
	Percentage float64
	Tag        string
	Type       string
	Scx        int
}
