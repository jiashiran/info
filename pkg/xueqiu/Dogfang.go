package xueqiu

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"info/pkg"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func reCheckByCSRDongfang(name string, res string, ctx context.Context) {
	exits := make(map[string]string)
	pkg.ReadLine("code-detail.csv", func(s string) {
		exits[strings.Split(s, ",")[0]] = "1"
	})
	pkg.ReadLine("old", func(s string) {
		exits[strings.Split(s, ",")[0]] = "1"
	})
	pkg.ReadLine(name, func(s string) {
		defer func() {
			if err := recover(); err != nil {
				pkg.ToTxt(fmt.Sprint(err), "err.txt")
			}
		}()
		ss := strings.Split(s, `,`)
		code := ss[0]
		if _, ok := exits[code]; ok {
			return
		}
		jzc := 0.0
		if !strings.Contains(ss[1], `<`) {
			jzc, _ = strconv.ParseFloat(ss[1], 64)
		}
		mll := 0.0
		if !strings.Contains(ss[2], `<`) {
			mll, _ = strconv.ParseFloat(ss[2], 64)
		}
		RevenueFrowth := 0.0
		if !strings.Contains(ss[3], `<`) {
			RevenueFrowth, _ = strconv.ParseFloat(ss[3], 64)
		}
		html := ""
		namehtml := ""
		stock := Stock{code: code, Jzcsyl: jzc, Mll: mll, RevenueFrowth: RevenueFrowth}
		codeR := strings.ReplaceAll(code, "SZ", "")
		codeR = strings.ReplaceAll(codeR, "SH", "")
		chromedp.Navigate("https://quote.cfi.cn/quote_" + codeR + `.html`).Do(ctx)
		time.Sleep(time.Duration(2) * time.Second)
		chromedp.InnerHTML(`form#formsearch`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		chromedp.InnerHTML(`form#formsearch`, &html, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		if strings.Contains(res, `Lfont`) {
			chromedp.InnerHTML(`div.Lfont`, &namehtml, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		} else {
			return
		}
		//res = strings.ReplaceAll(strings.ReplaceAll(res, `<span>`, ""), `</span>`, "")
		fmt.Println(res)
		stock = getDetailDongfang(res, html, namehtml, stock)
		value := fmt.Sprint(stock.code) + "," + fmt.Sprint(stock.Zsz) + "," +
			fmt.Sprint(stock.Ltz) + "," + fmt.Sprint(stock.Jzcsyl) + "," + fmt.Sprint(stock.Mll) + "," +
			fmt.Sprint(stock.Sylj) + "," + fmt.Sprint(stock.Ttm) + "," +
			fmt.Sprint(math.Floor(stock.Sylj-stock.Ttm)) + "," +
			fmt.Sprint(stock.Sum) + "," +
			fmt.Sprint(stock.Price) + "," + stock.Name + "," + fmt.Sprint(stock.RevenueFrowth) + "," + stock.Service + "," + stock.Info
		//代码,总市值,流通市值,净资产收益率,毛利率,市盈率(静),TTM,市盈率差,总股本,当前股价,名称,营收增长,服务,简介,
		_ = value
		pkg.ToTxt(value, "code-detail.csv")
		time.Sleep(time.Duration(rand.Int63n(10)) * time.Second)
	})
}

func getDetailDongfang(html, shtml, namehtml string, s Stock) Stock {
	fmt.Println(namehtml)
	//html = `<div class="stock-info"><div class="stock-price stock-fall"><div class="stock-current"><strong>¥18.03</strong></div><div class="stock-change">-0.12  -0.66%</div></div><div class="stock-time"><div>&nbsp;53.55 万球友关注</div><div class="quote-market-status"><span>交易中<pan> 09-30 14:12:03 北京时间</span></div></div></div><table class="quote-info"><tbody><tr><td>最高：<span class="stock-rise">18.20</spa<td>今开：<span class="stock-fall">18.09</span></td><td>涨停：<span class="stock-rise">19.97</span></td><td>成交量：<span>66.24万手</sp><tr class="separateTop"><td>最低：<span class="stock-fall">17.71</span></td><td>昨收：<span>18.15</span></td><td>跌停：<span class="st>16.34</span></td><td>成交额：<span>11.85亿</span></td></tr><tr class="separateBottom"><td>量比：<span class="stock-fall">0.66</span></手：<span>0.34%</span></td><td>市盈率(动)：<span>9.95</span></td><td>市盈率(TTM)：<span>10.66</span></td></tr><tr><td>委比：<span class-40.71%</span></td><td>振幅：<span>2.70%</span></td><td>市盈率(静)：<span>12.10</span></td><td>市净率：<span>1.14</span></td></tr><tr><n>1.69</span></td><td>股息(TTM)：<span>0.18</span></td><td>总股本：<span>194.06亿</span></td><td>总市值：<span>3498.89亿</span></td></t资产：<span>15.83</span></td><td>股息率(TTM)：<span>1.00%</span></td><td>流通股：<span>194.06亿</span></td><td>流通值：<span>3498.86亿<tr><td>52周最高：<span>25.16</span></td><td>52周最低：<span>14.64</span></td><td>货币单位：<span>CNY</span></td></tr></tbody></table>`
	zsz := strings.ReplaceAll(getByTagDongfang(html, `总市值`, "</span>"), "亿", "")
	ltz := getByTagDongfang(html, `流通市值`, "</span>")
	ttm := getByTagDongfang(html, `市盈(TTM)`, "</div>")
	sylj := getByTagDongfang(html, `市盈率:`, "倍")
	sum := getByTagDongfang(html, `总股本`, "</div>")
	shtm := getByTagDongfang(shtml, `price harmony-os-bold`, "</div>")
	price := getByTagDongfang(shtm, ">", "")
	Zsz, err := strconv.ParseFloat(zsz, 64)
	if err != nil {
		Zsz = 0
	}
	s.Zsz = Zsz

	//re := regexp.MustCompile("[\u4e00-\u9fa5]{1,}")
	//service := re.FindAllString(strings.ReplaceAll(getByTag(shtml, `<div class="title">业务</div>`, `<!---->`), ",", "，"), -1) //业务
	s.Name = namehtml //getByTagBaidu(namehtml, `>`, "</div>") //name
	//info := re.FindAllString(strings.ReplaceAll(getByTag(shtml, `<div class="title">简介</div>`, `<!---->`), ",", "，"), -1)    //简介

	//s.Name = strings.Join(name, "，")
	//s.Service = strings.Join(service, "，")
	//s.Info = strings.Join(info, "，")

	s.Ltz = ltz
	Ttm, err := strconv.ParseFloat(ttm, 64)
	if err != nil {
		Ttm = 0
	}
	s.Ttm = Ttm
	Sylj, err := strconv.ParseFloat(sylj, 64)
	if err != nil {
		Sylj = 0
	}
	s.Sylj = Sylj
	c := 1.0
	if strings.Contains(sum, "万") {
		sum = strings.ReplaceAll(sum, "万", "")
		c = 10000
	} else if strings.Contains(sum, "亿") {
		c = 100000000
		sum = strings.ReplaceAll(sum, "亿", "")
	}
	Sum, err := strconv.ParseFloat(sum, 64)
	if err != nil {
		Sum = 0.0
	}
	s.Sum = Sum * c
	Price, err := strconv.ParseFloat(price, 64)
	if err != nil {
		Price = 0.0
	}
	s.Price = Price
	//sy := getByTag(html,`市盈`)
	fmt.Println(zsz, ltz, ttm, sylj, s.Sum, s.Price)
	return s
}

func getByTagDongfang(src, tag, end string) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err, tag, end)
		}
	}()
	ss := strings.Split(src, tag)[1]
	if strings.Contains(ss, `font-weight: 500;">`) {
		ss = strings.Split(ss, `font-weight: 500;">`)[1]
	}
	if end != "" {
		ss = strings.Split(ss, end)[0]
		return ss
	}
	return strings.Split(ss, `</span>`)[0]
}
