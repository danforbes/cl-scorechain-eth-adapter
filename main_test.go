package main

import (
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/linkpoolio/bridges/bridge"
	"github.com/stretchr/testify/assert"
)

const (
	addressRegex = "^0x[a-z0-9]{40}$"
	dateRegex    = "^[0-9]{4}-[01][0-9]-[0-3][0-9]T[0-2][0-9]:[0-5][0-9]:[0-5][0-9]\\+[0-2][0-9]:[0-5][0-9]$"
	hashRegex    = "^0x[a-z0-9]{64}$"
	iso8601Regex = "^[0-9]{4}-[01][0-9]-[0-3][0-9]T[0-2][0-9]:[0-5][0-9]:[0-5][0-9]\\.[0-9]{3}Z$"
)

func TestScorechain_BadEndpoint(t *testing.T) {
	sc := &Scorechain{getScorechainURL(), os.Getenv("SC_ETH_TOKEN")}
	data := map[string]interface{}{
		"endpoint": "this is not a valid endpoint",
	}

	query, _ := bridge.ParseInterface(data)
	_, err := sc.Run(bridge.NewHelper(query))
	assert.NotNil(t, err)
	assert.Equal(t, "unrecognized or unsupported Scorechain Ethereum API endpoint", err.Error())
}

func TestScorechain_Status(t *testing.T) {
	sc := &Scorechain{getScorechainURL(), os.Getenv("SC_ETH_TOKEN")}
	data := map[string]interface{}{
		"endpoint": "get-status",
	}

	query, _ := bridge.ParseInterface(data)
	r, err := sc.Run(bridge.NewHelper(query))
	assert.Nil(t, err)

	s := r.(*Status)
	assert.True(t, s.Success)
	assert.LessOrEqual(t, s.Result.Timestamp, time.Now().Unix())
	assert.NotEmpty(t, s.Result.State)
	assert.NotEmpty(t, s.Result.LastBlock.Hash)
	assert.Greater(t, s.Result.LastBlock.Height, int64(0))
	assert.Greater(t, s.Result.LastBlock.Nbtx, int64(0))

	d, _ := regexp.Match(iso8601Regex, []byte(s.Result.LastBlock.Date))
	assert.True(t, d)
}

func TestScorechain_Transaction(t *testing.T) {
	sc := &Scorechain{getScorechainURL(), os.Getenv("SC_ETH_TOKEN")}
	hash := "0x7ead327c1d8d3ccc5ae619522ffec76334f4d786f84633c1462b968a71a4e8dd"
	data := map[string]interface{}{
		"endpoint": "get-trx",
		"hash":     hash,
	}

	query, _ := bridge.ParseInterface(data)
	r, err := sc.Run(bridge.NewHelper(query))
	assert.Nil(t, err)

	txn := r.(*Transaction)
	assert.True(t, txn.Success)
	assert.Equal(t, hash, txn.Result.Hash)
	assert.Greater(t, txn.Result.Block.Height, int64(0))

	bh, _ := regexp.Match(hashRegex, []byte(txn.Result.Block.Hash))
	assert.True(t, bh)

	fa, _ := regexp.Match(addressRegex, []byte(txn.Result.From.Address))
	assert.True(t, fa)
	assert.Empty(t, txn.Result.From.Label)
	assert.Empty(t, txn.Result.From.Type)
	assert.False(t, txn.Result.From.IsContract)

	ta, _ := regexp.Match(addressRegex, []byte(txn.Result.To.Address))
	assert.True(t, ta)
	assert.Empty(t, txn.Result.To.Label)
	assert.Empty(t, txn.Result.To.Type)
	assert.True(t, txn.Result.To.IsContract)
	assert.Equal(t, 8, txn.Result.Gas.Price)
	assert.Equal(t, 133705, txn.Result.Gas.Quantity)
	assert.Equal(t, 33705, txn.Result.Gas.Used)
	assert.Equal(t, 25.21, txn.Result.Gas.UsedPercent)
	assert.Equal(t, 0.00026964, txn.Result.Gas.Cost)
	assert.Equal(t, 500, txn.Result.Value.ETH)
	assert.IsType(t, 0.0, txn.Result.Value.AUD)
	assert.LessOrEqual(t, txn.Result.Timestamp, time.Now().Unix())

	d, _ := regexp.Match(dateRegex, []byte(txn.Result.Date))
	assert.True(t, d)
	assert.GreaterOrEqual(t, txn.Result.Confirmations, 3017326)
	assert.Equal(t, 2, len(txn.Result.InternalTransfers))
	assert.Equal(t, "Bitstamp", txn.Result.InternalTransfers[0].To.Label)
	assert.Equal(t, 23, txn.Result.InternalTransfers[0].To.Type.ID)
	assert.Equal(t, "Exchange", txn.Result.InternalTransfers[0].To.Type.Label)
	assert.Equal(t, 90, txn.Result.InternalTransfers[0].To.Type.Score)
	assert.Empty(t, txn.Result.TokenTransfers)
}

func TestScorechain_Scoring(t *testing.T) {
	sc := &Scorechain{getScorechainURL(), os.Getenv("SC_ETH_TOKEN")}
	data := map[string]interface{}{
		"endpoint":  "get-scoring",
		"address":   "0x0000000000000000000000000000000000000000",
		"direction": "outgoing",
	}

	query, _ := bridge.ParseInterface(data)
	r, err := sc.Run(bridge.NewHelper(query))
	assert.Nil(t, err)

	scx := r.(*Scoring)
	assert.True(t, scx.Success)
	assert.Greater(t, scx.Result.Scx, 0)
	assert.Empty(t, scx.Result.Details)
}

func TestScorechain_ScoringOpts(t *testing.T) {
	sc := &Scorechain{getScorechainURL(), os.Getenv("SC_ETH_TOKEN")}
	data := map[string]interface{}{
		"endpoint":  "get-scoring",
		"address":   "0x0000000000000000000000000000000000000000",
		"depth":     2,
		"direction": "incoming",
	}

	query, _ := bridge.ParseInterface(data)
	r, err := sc.Run(bridge.NewHelper(query))
	assert.Nil(t, err)

	scx := r.(*Scoring)
	assert.True(t, scx.Success)
	assert.Greater(t, scx.Result.Scx, 0)
	assert.NotEmpty(t, scx.Result.Details)
}
