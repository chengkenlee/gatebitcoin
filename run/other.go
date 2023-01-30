// @program:     gatebitcoin
// @file:        other.go
// @author:      chengkenlee
// @create:      2023-01-06 13:45
// @description:
package run

import (
	"fmt"
	"gatebitcoin/util"
	"github.com/antihax/optional"
	"github.com/gateio/gateapi-go/v6"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var w sync.WaitGroup

/*查询钱包总额*/
func amount() {
	balance, _, err := Gcli.WalletApi.GetTotalBalance(Gctx, nil)
	if err != nil {
		util.Logger.Error(err.Error())
	}
	amount, _ := strconv.ParseFloat(balance.Total.Amount, 64)
	Amount = amount
	w.Done()
}

/*查询货币单价*/
func (g *Gate) ticker(item gateapi.SpotAccount) {
	var c float64
	w.Add(1)
	go amount()
	w.Wait()
	f, _ := strconv.ParseFloat(item.Available, 64)
	ticker, _, err := Gcli.SpotApi.ListTickers(Gctx,
		&gateapi.ListTickersOpts{
			CurrencyPair: optional.NewString(fmt.Sprintf("%s_usdt", strings.ToLower(item.Currency))),
		},
	)
	if err != nil {
		util.Logger.Error(err.Error())
	}
	for _, i2 := range ticker {
		c, _ = strconv.ParseFloat(i2.Last, 64)
	}
	Ticker = c

	ava, _ := strconv.ParseFloat(item.Available, 64)
	g.T = &Tickers{
		Currency:    item.Currency,
		Available:   ava,
		Count:       Ticker,
		CurrencyCny: f * Ticker * 6.6962,
		Amount:      Amount,
		AmountCny:   Amount * 6.6962,
	}

	acc := fmt.Sprintf("%s：%f，COUNT：%f，CNY：%f，钱包总额(USDT)：%f，钱包总额(CNY)：%f",
		g.T.Currency, g.T.Available, g.T.Count, g.T.CurrencyCny, g.T.Amount, g.T.AmountCny)
	util.Logger.Info(acc)
	acJson = append(acJson, acc)
	wait.Done()
}

/*url触发*/
func (g *Gate) httpUrl(url string) []byte {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		util.Logger.Error(err.Error())
		return nil
	}
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	response, err := client.Do(request)
	if err != nil {
		util.Logger.Error(err.Error())
		return nil
	}
	defer response.Body.Close()
	bs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		util.Logger.Error(err.Error())
		return nil
	}
	return bs
}
