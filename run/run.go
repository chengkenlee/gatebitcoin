// @program:     gatebitcoin
// @file:        run.go
// @author:      chengkenlee
// @create:      2023-01-05 21:50
// @description:
package run

import (
	"context"
	"fmt"
	"gatebitcoin/util"
	"github.com/blinkbean/dingtalk"
	"github.com/gateio/gateapi-go/v6"
	"github.com/patrickmn/go-cache"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func init() {
	apiKey, _ := util.AesDecrypt2(util.Config.GetString("gateio.apikey"), util.ENCKEY)
	apiSecret, _ := util.AesDecrypt2(util.Config.GetString("gateio.secretkey"), util.ENCKEY)
	Gcli = gateapi.NewAPIClient(gateapi.NewConfiguration())
	Gctx = context.WithValue(context.Background(),
		gateapi.ContextGateAPIV4,
		gateapi.GateAPIV4{
			Key:    apiKey,
			Secret: apiSecret,
		})
}

func Run() {
	g.total()
	g.sum()
}

/*主体*/
func (g *Gate) sum() {
	tt := time.NewTicker(time.Second * time.Duration(util.Config.GetInt("gateio.cache.time.newticker")))
	c := cache.New(12*time.Hour, 24*time.Hour)
	i := 0
	for {
		util.Logger.Info(fmt.Sprintf("**************************** %d ****************************", i))
		select {
		case <-tt.C:
			var title string
			g.total()
			now := time.Now()
			k := g.T.Currency + now.Format("20060102T150405")
			v := g.T.AmountCny
			/*设置当前的数值*/
			c.Set(k, v, cache.DefaultExpiration)
			/*获取当前数值*/
			f0, _ := c.Get(k)
			/*获取一分钟之前的数值*/
			t, err := time.ParseDuration(util.Config.GetString("gateio.cache.time.parseduration"))
			if err != nil {
				util.Logger.Error(err.Error())
				continue
			}
			k2 := g.T.Currency + now.Add(t).Format("20060102T150405")
			f1, b1 := c.Get(k2)
			if b1 {
				_, b := c.Get("song")
				if i == 100 {
					i = 0
					c.Delete("song")
				}
				if b {
					util.Logger.Info("发现已有提示，不做信息发送")
					i++
					continue
				}
				util.Logger.Info(fmt.Sprintf("key：%s now：%f，%s(%s)：%f", now.Add(t).Format("20060102T150405"), f0, k2, util.Config.GetString("gateio.cache.time.parseduration"), f1))
				if f1.(float64) >= f0.(float64) {
					/*正值*/
					if f1.(float64)-f0.(float64) >= util.Config.GetFloat64("gateio.cache.loss") {
						title = fmt.Sprintf("亏损(%f)元", f1.(float64)-f0.(float64))
					}
				} else {
					/*负值*/
					if f0.(float64)-f1.(float64) >= util.Config.GetFloat64("gateio.cache.loss") {
						title = fmt.Sprintf("盈利(%f)元", f0.(float64)-f1.(float64))
					}
				}
				g.dingSend(title, f0.(float64), f1.(float64))
				g.wechatSend(title, f0.(float64), f1.(float64))
				c.Set("song", "done", cache.DefaultExpiration)
			}
			/*获取10分钟，30分钟，60分钟的数据与现在的对比差距*/
			g.gettime(c, now, f0)
		}
		i++
	}
}

/*发送dingding*/
func (g *Gate) dingSend(title string, f0, f1 float64) {
	dm := dingtalk.DingMap()
	dm.Set(title, dingtalk.H1)
	dm.Set("---", dingtalk.N)
	dm.Set(fmt.Sprintf("%f 现在", f0), dingtalk.GREEN)
	dm.Set(fmt.Sprintf("%f %s", f1, util.Config.GetString("gateio.cache.time.parseduration")), dingtalk.GREEN)
	dm.Set(strings.Join(acJson, "\n"), dingtalk.H6)
	dm.Set(strings.Join(Arr, "\n"), dingtalk.H6)
	err := util.Dingcli.SendMarkDownMessageBySlice(fmt.Sprintf("¥ %f", g.T.CurrencyCny), dm.Slice())
	if err != nil {
		util.Logger.Error(err.Error())
		return
	}
}
func (g *Gate) wechatSend(title string, f0, f1 float64) {
	m := fmt.Sprintf(`hey~ ChengKen，
现在是: %s，
我发现账号 %s
-----------
【%f】 现在
【%f】 %s
-----------
%s
%s
`, time.Now().Format("2006-01-02 15:04:05"), title, f0, f1, util.Config.GetString("gateio.cache.time.parseduration"), strings.Join(acJson, "\n"), strings.Join(Arr, "\n"))
	escapeUrl := url.QueryEscape(m)
	url := fmt.Sprintf("http://localhost:3001/send/%s?msg=%s", util.AesDecrypt(util.Config.GetString("wechat.token"), util.ENCKEY), escapeUrl)
	res := g.httpUrl(url)
	util.Logger.Info(string(res))
}

/*统计金额*/
func (g *Gate) total() {
	r, _, err := Gcli.SpotApi.ListSpotAccounts(Gctx, nil)
	if err != nil {
		util.Logger.Error(err.Error())
		return
	}
	acJson = nil
	for _, item := range r {
		f, _ := strconv.ParseFloat(item.Available, 64)
		if f < 1 {
			continue
		}
		wait.Add(1)
		go g.ticker(item)
	}
	wait.Wait()
}

/*获取10分钟，30分钟，60分钟值*/
func (g *Gate) gettime(c *cache.Cache, now time.Time, f0 interface{}) {
	_, b := c.Get("song")
	if b {
		util.Logger.Info("发现已有提示，不做信息发送")
		return
	}
	util.Logger.Warn(fmt.Sprintf("【%s】对比历史数据。。。", now.Format("2006-01-02 15:04")))
	var title string

	for _, s := range util.TicksTime {
		t, err := time.ParseDuration(s)
		if err != nil {
			util.Logger.Error(err.Error())
			return
		}
		k := g.T.Currency + now.Add(t).Format("20060102T150405")
		f1, b1 := c.Get(k)
		if b1 {
			util.Logger.Info(fmt.Sprintf("【##############获取到%s到现在的值有变化差异，key：%s now：%f，%s(%s)：%f】", s, now.Add(t).Format("20060102T150405"), f0, k, util.Config.GetString("gateio.cache.time.parseduration"), f1))
			if f1.(float64) >= f0.(float64) {
				/*正值*/
				if f1.(float64)-f0.(float64) >= util.Config.GetFloat64("gateio.cache.loss") {
					title = fmt.Sprintf("哇~ 发现%s到现在 亏损(%f)元", s, f1.(float64)-f0.(float64))
				}
			} else {
				/*负值*/
				if f0.(float64)-f1.(float64) >= util.Config.GetFloat64("gateio.cache.loss") {
					title = fmt.Sprintf("咳~ tmd%s到现在已经 盈利了(%f)元", s, f0.(float64)-f1.(float64))
				}
			}
			g.dingSend(title, f0.(float64), f1.(float64))
			g.wechatSend(title, f0.(float64), f1.(float64))
		}
	}
}
