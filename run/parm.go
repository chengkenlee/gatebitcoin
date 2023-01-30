// @program:     gatebitcoin
// @file:        parm.go
// @author:      chengkenlee
// @create:      2023-01-05 21:55
// @description:
package run

import (
	"context"
	"github.com/gateio/gateapi-go/v6"
	"sync"
)

var (
	wait   sync.WaitGroup
	Amount float64
	Ticker float64
	acJson []string
	Arr    []string
	g      Gate
	Gcli   *gateapi.APIClient
	Gctx   context.Context
)

type Gate struct {
	E *Exahange
	T *Tickers
}

type Exahange struct {
	Reason string `json:"reason"`
	Result struct {
		Update string     `json:"update"`
		List   [][]string `json:"list"`
	} `json:"result"`
	ErrorCode int `json:"error_code"`
}

type Tickers struct {
	Currency    string  `bson:"Currency"`
	Available   float64 `bson:"Available"`
	Count       float64 `bson:"Count"`
	CurrencyCny float64 `bson:"CurrencyCny"`
	Amount      float64 `bson:"Amount"`
	AmountCny   float64 `bson:"AmountCny"`
}
