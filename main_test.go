package main

import (
	"os"
	"testing"

	"github.com/linkpoolio/bridges/bridge"
	"github.com/stretchr/testify/assert"
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
	_, err := sc.Run(bridge.NewHelper(query))
	assert.Nil(t, err)
}

func TestScorechain_Transaction(t *testing.T) {
	sc := &Scorechain{getScorechainURL(), os.Getenv("SC_ETH_TOKEN")}
	data := map[string]interface{}{
		"endpoint": "get-trx",
		"hash":     "0x0000000000000000000000000000000000000000000000000000000000000000",
	}

	query, _ := bridge.ParseInterface(data)
	_, err := sc.Run(bridge.NewHelper(query))
	assert.Nil(t, err)
}
