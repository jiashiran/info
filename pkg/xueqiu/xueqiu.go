package xueqiu

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"math"
	"math/rand"
	"regexp"
	"spider/pkg"
	"spider/pkg/jijin"
	"strconv"
	"strings"
	"time"
)

var stocks = make([]Stock, 0)

func Run() {
	//jijin.DbInit()
	// 代码,净资产收益率,毛利率
	getDetailHTML() // 净资产收益率,毛利率 1
	//C() //合并，过滤

	//代码,总市值,流通市值,净资产收益率,毛利率,市盈率(静),TTM,市盈率差,总股本,当前股价,名称,服务,简介  //https://xueqiu.com/S/SH688677
	getDHTML() //详细数据， 2  //getDetail("",Stock{code: "123"})

	//代码,总市值,流通市值,净资产收益率,毛利率,市盈率(静),TTM,市盈率差,总股本,当前股价,名称,服务,简介,自由现金流  //getFFC() //自由现金流
	//getCaculateFFC("code-detail.csv") // 3 //代码,总市值,流通市值,净资产收益率,净利率,市盈率(静),TTM,市盈率(动),市盈率差,总股本,当前股价,现金流股价比率,自由现金流,每股自由现金流,上期自由现金流,上期每股自由现金流,名称,服务,简介
}

func caculate() {
	pkg.ReadLine("stock-result.csv", func(s string) {
		if s == "" {
			return
		}
		ss := strings.Split(s, `,`)
		TTM, _ := strconv.ParseFloat(ss[6], 64)
		Sum, _ := strconv.ParseFloat(ss[8], 64)
		Jzc, _ := strconv.ParseFloat(ss[3], 64)
		_ = Jzc
		Price, _ := strconv.ParseFloat(ss[9], 64)
		_ = Price
		flow := ss[13]
		unit := 0.0
		if !strings.Contains(flow, "-") && TTM > 0 {
			if strings.Contains(flow, "千亿") {
				unit = 100000000 * 1000
				flow = strings.ReplaceAll(flow, "千亿", "")

			} else if strings.Contains(flow, "百亿") {
				unit = 100000000 * 100
				flow = strings.ReplaceAll(flow, "百亿", "")

			} else if strings.Contains(flow, "亿") {
				unit = 100000000
				flow = strings.ReplaceAll(flow, "亿", "")

			} else if strings.Contains(flow, "千万") {
				unit = 10000 * 100
				flow = strings.ReplaceAll(flow, "千万", "")

			} else if strings.Contains(flow, "百万") {
				unit = 10000 * 100
				flow = strings.ReplaceAll(flow, "百万", "")

			} else if strings.Contains(flow, "万") {
				unit = 10000
				flow = strings.ReplaceAll(flow, "万", "")

			}

			Flow, err := strconv.ParseFloat(flow, 64)
			if err != nil {
				fmt.Println(err)
			}
			Flow = Flow * unit
			valuation := (Flow / Sum)
			if valuation == 0 {
				fmt.Println(ss[0], Flow, Sum, TTM)
			}
			ttmNew := Price / valuation
			value := ss[0] + "," + ss[1] + "," + ss[2] + "," + ss[3] + "," + ss[4] + "," + ss[5] + "," + ss[6] + "," + ss[7] + "," + ss[8] + "," +
				ss[9] + "," + ss[13] + "," + fmt.Sprint(valuation) + "," + fmt.Sprint(ttmNew) + "," + fmt.Sprint(ttmNew/TTM) + "," +
				ss[10] + "," + ss[11] + "," + ss[12]
			/*if Jzc > 15 && Jzc < 20{
				pkg.ToTxt(value,"stock-valuation-15.csv")
			}*/
			pkg.ToTxt(value, "stock-valuation.csv")

		}

	})
}

// 查询净资产收益率
func getFFCDetail(name, res string, ctx context.Context) {
	pkg.ReadLine(name, func(s string) {
		ss := strings.Split(s, `,`)
		code := ss[0]
		/*TTM,_ := strconv.ParseFloat(ss[5],64)
		Sum,_ := strconv.ParseFloat(ss[7],64)
		Price ,_ := strconv.ParseFloat(ss[8],64)
		stock := Stock{code: code,Sum: Sum,Price: Price,Ttm: TTM}*/
		runes := []rune(code)
		chromedp.Navigate("https://caibaoshuo.com/terms/" + string(runes[2:]) + "/free_cash_flow").Do(ctx)
		chromedp.OuterHTML(`body`, &res, chromedp.ByQuery).Do(ctx)
		fmt.Println(res)
		if strings.Contains(res, `您要访问的页面不存在`) || strings.Contains(res, `服务器内部错误`) {
			return
		}
		chromedp.InnerHTML(`div.scroll-container`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		//fmt.Println(res)
		res = strings.Split(res, `自由现金流(FCF)`)[1]
		ss1 := strings.Split(res, ` </td>`)
		res = strings.ReplaceAll(ss1[len(ss1)-2], `<td>`, ``)
		res = strings.ReplaceAll(res, `
`, ``)
		res = strings.ReplaceAll(res, ` `, ``)
		fmt.Println(res)
		value := s + "," + fmt.Sprint(res)
		pkg.ToTxt(value, "stock-result.csv")
		time.Sleep(time.Duration(rand.Int63n(2)) * time.Second)
	})
}

func getFFC() {
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
		chromedp.WindowSize(1280, 1024),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	var res string
	tasks := chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("step1")

			getFFCDetail("code-detail.csv", res, ctx)
			return nil
		}),
	}
	fmt.Println("start run chromedp tasks!")
	err := chromedp.Run(ctx,
		tasks,
	)
	if err != nil {
		log.Println(err)
	}
}

// 遍历csr,查询详细数据
func getDHTML() {
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
		chromedp.WindowSize(1280, 1024),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	var res string
	tasks := chromedp.Tasks{
		//network.Enable(),
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("step1")
			reCheckByCSRBaidu("code.csv", res, ctx)
			//reCheckByCSRDongfang("code.csv", res, ctx)
			return nil
		}),
	}
	fmt.Println("start run chromedp tasks!")
	err := chromedp.Run(ctx,
		tasks,
	)
	if err != nil {
		log.Println(err)
	}
}

func reCheckByCSRBaidu(name string, res string, ctx context.Context) {
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
		chromedp.Navigate("https://gushitong.baidu.com/stock/ab-" + codeR + ``).Do(ctx)
		time.Sleep(time.Duration(2) * time.Second)
		chromedp.InnerHTML(`div.quote-market-container`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		chromedp.InnerHTML(`div.trading-num.flex-align-center`, &html, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		chromedp.InnerHTML(`div.name.text-nowrap.c-gap-right.cos-text-headline-xl.cos-font-medium.cos-color-text`, &namehtml, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		res = strings.ReplaceAll(strings.ReplaceAll(res, `<span>`, ""), `</span>`, "")
		//fmt.Println(res)
		stock = getDetailBaidu(res, html, namehtml, stock)
		value := fmt.Sprint(stock.code) + "," + fmt.Sprint(stock.Zsz) + "," +
			fmt.Sprint(stock.Ltz) + "," + fmt.Sprint(stock.Jzcsyl) + "," + fmt.Sprint(stock.Mll) + "," +
			fmt.Sprint(stock.Sylj) + "," + fmt.Sprint(stock.Ttm) + "," + fmt.Sprint(stock.Syld) + "," +
			fmt.Sprint(math.Floor(stock.Sylj-stock.Ttm)) + "," +
			fmt.Sprint(stock.Sum) + "," +
			fmt.Sprint(stock.Price) + "," + stock.Name + "," + fmt.Sprint(stock.RevenueFrowth) + "," + stock.Service + "," + stock.Info
		//代码,总市值,流通市值,净资产收益率,毛利率,市盈率(静),TTM,市盈率(动),市盈率差,总股本,当前股价,名称,营收增长,服务,简介,
		_ = value
		pkg.ToTxt(value, "code-detail.csv")
		time.Sleep(time.Duration(rand.Int63n(15)) * time.Second)
	})
}

func reCheckByCSRSouhu(name string, res string, ctx context.Context) {
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
		html := ""
		namehtml := ""
		stock := Stock{code: code, Jzcsyl: jzc, Mll: mll}
		codeR := strings.ReplaceAll(code, "SZ", "")
		codeR = strings.ReplaceAll(codeR, "SH", "")
		chromedp.Navigate("https://q.stock.sohu.com/cn/" + codeR + `/index.shtml`).Do(ctx)
		time.Sleep(time.Duration(2) * time.Second)
		chromedp.InnerHTML(`div.quote-market-container`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		chromedp.InnerHTML(`div.trading-num.flex-align-center`, &html, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		chromedp.InnerHTML(`div.name.text-nowrap.c-gap-right.cos-text-headline-xl.cos-font-medium.cos-color-text`, &namehtml, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		res = strings.ReplaceAll(strings.ReplaceAll(res, `<span>`, ""), `</span>`, "")
		//fmt.Println(res)
		stock = getDetailBaidu(res, html, namehtml, stock)
		value := fmt.Sprint(stock.code) + "," + fmt.Sprint(stock.Zsz) + "," +
			fmt.Sprint(stock.Ltz) + "," + fmt.Sprint(stock.Jzcsyl) + "," + fmt.Sprint(stock.Mll) + "," +
			fmt.Sprint(stock.Sylj) + "," + fmt.Sprint(stock.Ttm) + "," +
			fmt.Sprint(math.Floor(stock.Sylj-stock.Ttm)) + "," +
			fmt.Sprint(stock.Sum) + "," +
			fmt.Sprint(stock.Price) + "," + stock.Name + "," + stock.Service + "," + stock.Info
		//代码,总市值,流通市值,净资产收益率,毛利率,市盈率(静),TTM,市盈率差,总股本,当前股价,名称,服务,简介
		_ = value
		pkg.ToTxt(value, "code-detail.csv")
		time.Sleep(time.Duration(rand.Int63n(10)) * time.Second)
	})
}

func reCheckByCSR(name string, res string, ctx context.Context) {
	exits := make(map[string]string)
	pkg.ReadLine("code-detail.csv", func(s string) {
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
		stock := Stock{code: code, Jzcsyl: jzc, Mll: mll}
		shtml := ""
		chromedp.Navigate("https://xueqiu.com/S/" + code + ``).Do(ctx)
		time.Sleep(time.Duration(2) * time.Second)
		chromedp.InnerHTML(`div.container-lg`, &shtml, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		chromedp.InnerHTML(`div.quote-container`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		//res = strings.ReplaceAll(strings.ReplaceAll(res,`<span>`,""),`</span>`,"")
		fmt.Println(res)
		stock = getDetail(res, shtml, stock)
		value := fmt.Sprint(stock.code) + "," + fmt.Sprint(stock.Zsz) + "," +
			fmt.Sprint(stock.Ltz) + "," + fmt.Sprint(stock.Jzcsyl) + "," + fmt.Sprint(stock.Mll) + "," +
			fmt.Sprint(stock.Sylj) + "," + fmt.Sprint(stock.Ttm) + "," +
			fmt.Sprint(math.Floor(stock.Sylj-stock.Ttm)) + "," +
			fmt.Sprint(stock.Sum) + "," +
			fmt.Sprint(stock.Price) + "," + stock.Name + "," + stock.Service + "," + stock.Info
		//代码,总市值,流通市值,净资产收益率,毛利率,市盈率(静),TTM,市盈率差,总股本,当前股价,名称,服务,简介
		pkg.ToTxt(value, "code-detail.csv")
		var ccs, _ = network.GetAllCookies().Do(ctx)
		for i, cookie := range ccs {
			log.Printf("chrome cookie %d: %+v", i, cookie)
		}
		time.Sleep(time.Duration(rand.Int63n(10)) * time.Second)
	})
}

func getDetail(html, shtml string, s Stock) Stock {
	fmt.Println(shtml)
	//html = `<div class="stock-info"><div class="stock-price stock-fall"><div class="stock-current"><strong>¥18.03</strong></div><div class="stock-change">-0.12  -0.66%</div></div><div class="stock-time"><div>&nbsp;53.55 万球友关注</div><div class="quote-market-status"><span>交易中<pan> 09-30 14:12:03 北京时间</span></div></div></div><table class="quote-info"><tbody><tr><td>最高：<span class="stock-rise">18.20</spa<td>今开：<span class="stock-fall">18.09</span></td><td>涨停：<span class="stock-rise">19.97</span></td><td>成交量：<span>66.24万手</sp><tr class="separateTop"><td>最低：<span class="stock-fall">17.71</span></td><td>昨收：<span>18.15</span></td><td>跌停：<span class="st>16.34</span></td><td>成交额：<span>11.85亿</span></td></tr><tr class="separateBottom"><td>量比：<span class="stock-fall">0.66</span></手：<span>0.34%</span></td><td>市盈率(动)：<span>9.95</span></td><td>市盈率(TTM)：<span>10.66</span></td></tr><tr><td>委比：<span class-40.71%</span></td><td>振幅：<span>2.70%</span></td><td>市盈率(静)：<span>12.10</span></td><td>市净率：<span>1.14</span></td></tr><tr><n>1.69</span></td><td>股息(TTM)：<span>0.18</span></td><td>总股本：<span>194.06亿</span></td><td>总市值：<span>3498.89亿</span></td></t资产：<span>15.83</span></td><td>股息率(TTM)：<span>1.00%</span></td><td>流通股：<span>194.06亿</span></td><td>流通值：<span>3498.86亿<tr><td>52周最高：<span>25.16</span></td><td>52周最低：<span>14.64</span></td><td>货币单位：<span>CNY</span></td></tr></tbody></table>`
	zsz := strings.ReplaceAll(getByTag(html, `总市值：`+`<span>`, ""), "亿", "")
	ltz := getByTag(html, `流通值：`+`<span>`, "")
	ttm := getByTag(html, `市盈率(TTM)：`+`<span>`, "")
	sylj := getByTag(html, `市盈率(静)：`+`<span>`, "")
	sum := getByTag(html, `总股本：<span>`, "")
	price := strings.ReplaceAll(getByTag(html, `<div class="stock-current"><strong>`, "</strong>"), "¥", "")
	Zsz, err := strconv.ParseFloat(zsz, 64)
	if err != nil {
		Zsz = 0
	}
	s.Zsz = Zsz

	re := regexp.MustCompile("[\u4e00-\u9fa5]{1,}")
	service := re.FindAllString(strings.ReplaceAll(getByTag(shtml, `<div class="title">业务</div>`, `<!---->`), ",", "，"), -1) //业务
	name := re.FindAllString(strings.ReplaceAll(getByTag(shtml, `<h1 class="stock-name">`, `(`), ",", "，"), -1)              //name
	info := re.FindAllString(strings.ReplaceAll(getByTag(shtml, `<div class="title">简介</div>`, `<!---->`), ",", "，"), -1)    //简介

	s.Name = strings.Join(name, "，")
	s.Service = strings.Join(service, "，")
	s.Info = strings.Join(info, "，")

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

type Stock struct {
	Name          string  `名称`
	Service       string  `业务`
	Info          string  `简介`
	code          string  `代码`
	Jzcsyl        float64 `净资产收益率`
	Mll           float64 `毛利率`
	Zsz           float64 `总市值`
	Ltz           string  `流通市值`
	Ttm           float64 `TTM`
	Syld          float64 `市盈率(动)`
	Sylj          float64 `市盈率(静)`
	Sum           float64 `总股本`
	Price         float64 `当前价格`
	RevenueFrowth float64 `营业收入增长`
}

func getByTag(src, tag, end string) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	ss := strings.Split(src, tag)[1]
	if end != "" {
		return strings.Split(ss, end)[0]
	}
	return strings.Split(ss, `</span>`)[0]
}

func getByTagBaidu(src, tag, end string) string {
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
		return strings.Split(ss, end)[0]
	}
	return strings.Split(ss, `</span>`)[0]
}

func getDetailBaidu(html, shtml, namehtml string, s Stock) Stock {
	fmt.Println(namehtml)
	//html = `<div class="stock-info"><div class="stock-price stock-fall"><div class="stock-current"><strong>¥18.03</strong></div><div class="stock-change">-0.12  -0.66%</div></div><div class="stock-time"><div>&nbsp;53.55 万球友关注</div><div class="quote-market-status"><span>交易中<pan> 09-30 14:12:03 北京时间</span></div></div></div><table class="quote-info"><tbody><tr><td>最高：<span class="stock-rise">18.20</spa<td>今开：<span class="stock-fall">18.09</span></td><td>涨停：<span class="stock-rise">19.97</span></td><td>成交量：<span>66.24万手</sp><tr class="separateTop"><td>最低：<span class="stock-fall">17.71</span></td><td>昨收：<span>18.15</span></td><td>跌停：<span class="st>16.34</span></td><td>成交额：<span>11.85亿</span></td></tr><tr class="separateBottom"><td>量比：<span class="stock-fall">0.66</span></手：<span>0.34%</span></td><td>市盈率(动)：<span>9.95</span></td><td>市盈率(TTM)：<span>10.66</span></td></tr><tr><td>委比：<span class-40.71%</span></td><td>振幅：<span>2.70%</span></td><td>市盈率(静)：<span>12.10</span></td><td>市净率：<span>1.14</span></td></tr><tr><n>1.69</span></td><td>股息(TTM)：<span>0.18</span></td><td>总股本：<span>194.06亿</span></td><td>总市值：<span>3498.89亿</span></td></t资产：<span>15.83</span></td><td>股息率(TTM)：<span>1.00%</span></td><td>流通股：<span>194.06亿</span></td><td>流通值：<span>3498.86亿<tr><td>52周最高：<span>25.16</span></td><td>52周最低：<span>14.64</span></td><td>货币单位：<span>CNY</span></td></tr></tbody></table>`
	zsz := strings.ReplaceAll(getByTagBaidu(html, `总市值`, "</div>"), "亿", "")
	ltz := getByTagBaidu(html, `流通值`, "</div>")
	ttm := getByTagBaidu(html, `市盈(TTM)`, "</div>")
	syld := getByTagBaidu(html, `市盈率(动)`, "</div>")
	sylj := getByTagBaidu(html, `市盈(静)`, "</div>")
	sum := getByTagBaidu(html, `总股本`, "</div>")
	shtm := getByTagBaidu(shtml, `price harmony-os-bold`, "</div>")
	price := getByTagBaidu(shtm, ">", "")
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
	Syld, err := strconv.ParseFloat(syld, 64)
	if err != nil {
		Syld = 0
	}
	s.Syld = Syld
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

func getListHTML(name string) {
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
		chromedp.WindowSize(1280, 1024),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	var res string
	tasks := chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("step1")

			var count int32 = 0
			jijin.Db.Table("info").Count(&count)
			var pageIndex int32 = 1
			var pageSize int32 = 100
			for ; pageIndex*pageSize < count; pageIndex++ {
				infos := jijin.SelectAllStock(pageIndex, pageSize)
				for _, stock := range infos {
					if stock.Proportion == "" {
						chromedp.Navigate("https://xueqiu.com/k?q=" + stock.Name).Do(ctx)

						chromedp.InnerHTML(`p.search__stock__bd__code`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)

						res = strings.ReplaceAll(strings.ReplaceAll(res, `<span>`, ""), `</span>`, "")

						fmt.Println(res)

						stock.Num = res

						chromedp.Navigate("https://xueqiu.com/snowman/S/" + res + "/detail#/JJCG").Do(ctx)
						time.Sleep(3 * time.Second)
						chromedp.InnerHTML(`div.container-md.float-left.stock__info__main`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
						if strings.Contains(res, "全部合计") {
							chromedp.InnerHTML(`table.brief-info tbody tr`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
							res = strings.ReplaceAll(res, `\n`, "")
							res = strings.Split(res, `<td>`)[3]
							res = strings.ReplaceAll(res, `</td>`, "")

							fmt.Println(res)

							stock.Proportion = res

							jijin.Db.Save(&stock).Commit()
						} else {
							stock.Proportion = "null"

							jijin.Db.Save(&stock).Commit()
						}
					}
				}
			}

			return nil
		}),
	}
	fmt.Println("start run chromedp tasks!")
	err := chromedp.Run(ctx,

		tasks,
	)
	if err != nil {
		log.Println(err)
	}

}

// 查询所以净资产收益率
func getDetailHTML() {
	url := "https://xueqiu.com/snowman/S/%s/detail#/ZYCWZB"
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
		chromedp.WindowSize(1280, 1024),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	var res string
	tasks := chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("step1")
			search300(res, url, ctx)
			search002(res, url, ctx)
			search600(res, url, ctx)
			search688(res, url, ctx)
			search000(res, url, ctx)
			return nil
		}),
	}
	fmt.Println("start run chromedp tasks!")
	err := chromedp.Run(ctx,
		tasks,
	)
	if err != nil {
		log.Println(err)
	}

}

// 查询所以净资产收益率
func getCaculateFFC(file string) {
	exits := make(map[string]string)
	pkg.ReadLine("result.csv", func(s string) {
		exits[strings.Split(s, ",")[0]] = "1"
	})
	url := "https://xueqiu.com/snowman/S/%s/detail#/XJLLB"
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
		chromedp.WindowSize(1280, 1024),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	var res string
	tasks := chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			line := 1
			pkg.ReadLine(file, func(s string) {
				defer func() {
					if err := recover(); err != nil {
						pkg.ToTxt(fmt.Sprint(err), "err.txt")
					}
				}()
				ss := strings.Split(s, `,`)
				code := ss[0]
				jll, _ := strconv.ParseFloat(ss[4], 64)
				if jll < 30 {
					return
				}
				zengzhang, _ := strconv.ParseFloat(ss[11], 64)
				if zengzhang < 5 {
					return
				}
				if _, ok := exits[code]; ok {
					return
				}
				chromedp.Navigate(strings.ReplaceAll(url, "%s", code)).Do(ctx)
				if line == 1 {
					time.Sleep(5 * time.Second)
				}
				chromedp.OuterHTML(`body`, &res, chromedp.ByQuery).Do(ctx)
				if !strings.Contains(res, `购建固定资产、无形资产和其他长期资产支付的现金`) {
					fmt.Println("no date code:" + code)
					return
				}
				chromedp.InnerHTML(`div.wrapper`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
				if strings.Contains(res, `购建固定资产、无形资产和其他长期资产支付的现金</td>`) {
					getFlow(res, ss, s)
				}
				line++
				if line > 1 {
					//os.Exit(0)
				}
				time.Sleep(time.Duration(rand.Int63n(4)) * time.Second)
			})

			return nil
		}),
	}

	fmt.Println("start run chromedp tasks!")
	err := chromedp.Run(ctx,
		tasks,
	)
	if err != nil {
		log.Println(err)
	}

}

func getFlow(res string, ss []string, s string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Errorf("getFlow err:%s", err)
		}
	}()
	//代码,总市值,流通市值,净资产收益率,毛利率,市盈率(静),TTM,市盈率差,总股本,当前股价,名称,服务,简介,自由现金流,每股自由现金流,上期自由现金流,上期每股自由现金流,
	res = strings.Split(res, `<td colspan="2">经营活动产生的现金流量净额</td>`)[1]
	je := strings.ReplaceAll(strings.Split(res, `<span></span>`)[0], `<td><p>`, ``)
	jeOld := strings.ReplaceAll(strings.Split(res, `<span></span>`)[2], `<td><p>`, ``)
	jeOld = strings.Split(jeOld, `</p></td>`)[1]
	res = strings.Split(res, `<td colspan="2">购建固定资产、无形资产和其他长期资产支付的现金</td>`)[1]
	zf := strings.ReplaceAll(strings.Split(res, `<span></span>`)[0], `<td><p>`, ``)
	zfOld := strings.ReplaceAll(strings.Split(res, `<span></span>`)[2], `<td><p>`, ``)
	zfOld = strings.Split(zfOld, `</p></td>`)[1]
	fmt.Println("净现金：", je, "支付现金：", zf)
	fmt.Println(jeOld, zfOld)
	flow := transferFlout(je) - transferFlout(zf)
	flowOld := transferFlout(jeOld) - transferFlout(zfOld)
	zgm, err := strconv.ParseFloat(ss[8], 64)
	if err != nil {
		fmt.Println(err)
		zgm = -1
	}
	ttm, err := strconv.ParseFloat(ss[6], 64)
	if err != nil {
		fmt.Println(err)
		ttm = 1
	}
	price, err := strconv.ParseFloat(ss[9], 64)
	if err != nil {
		fmt.Println(err)
		price = 1000
	}
	fmt.Println(flow)
	flowMeigu := flow / zgm
	flowMeiguOld := flowOld / zgm
	sort := flowMeigu * ttm / price
	lenn := len(ss)
	s = strings.Join(ss[:lenn-3], ",")
	pkg.ToTxt(s+","+fmt.Sprint(sort)+","+fmt.Sprint(flow)+","+fmt.Sprint(flowMeigu)+","+fmt.Sprint(flowOld)+","+fmt.Sprint(flowMeiguOld)+","+fmt.Sprint(ss[lenn-3])+","+fmt.Sprint(ss[lenn-2])+","+fmt.Sprint(ss[lenn-1]), "result.csv")
}

func transferFlout(flow string) float64 {
	var unit float64
	if strings.Contains(flow, "百亿") {
		unit = 100000000 * 100
		flow = strings.ReplaceAll(flow, "百亿", "")
	} else if strings.Contains(flow, "亿") {
		unit = 100000000
		flow = strings.ReplaceAll(flow, "亿", "")
	} else if strings.Contains(flow, "千万") {
		unit = 10000 * 100
		flow = strings.ReplaceAll(flow, "千万", "")
	} else if strings.Contains(flow, "百万") {
		unit = 10000 * 100
		flow = strings.ReplaceAll(flow, "百万", "")
	} else if strings.Contains(flow, "万") {
		unit = 10000
		flow = strings.ReplaceAll(flow, "万", "")
	}
	Flow, err := strconv.ParseFloat(flow, 64)
	if err != nil {
		fmt.Println(err)
	}
	Flow = Flow * unit
	return Flow
}

func search300(res, url string, ctx context.Context) {
	m := make(map[string]string)
	pkg.ReadLine("300.csv", func(s string) {
		c := strings.Split(s, ",")[0]
		m[c] = "1"
	})
	for i := 300000; i <= 300999; i++ {
		//SZ 300000  创业板
		//SZ00 2000  中小板
		//SH 600000  601000 603000 沪市A
		//SZ000999 递减  深市A
		cs := []rune(fmt.Sprint(i))
		code := "SZ" + string(cs)
		if _, ok := m[code]; ok {
			continue
		}
		chromedp.Navigate(strings.ReplaceAll(url, "%s", code)).Do(ctx)
		chromedp.InnerHTML(`div#app`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		if strings.Index(res, `净资产收益率`) > 0 {
			s := strings.Split(res, `净资产收益率</td>`)[1]
			s = strings.Split(s, `<span>%</span>`)[0]
			s = strings.ReplaceAll(s, `<td><p>`, "")
			fmt.Println(s)
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				f = 0.0
				s = ""
			}
			mll := strings.Split(res, `净利率</td>`)[1]
			mll = strings.Split(mll, `<span>%</span>`)[0]
			mll = strings.ReplaceAll(mll, `<td><p>`, "")
			//fmt.Println("毛利率：", mll)
			mllf, err := strconv.ParseFloat(mll, 64)
			if err != nil {
				mllf = 0.0
			}
			yyze := strings.Split(res, `营业收入同比增长</td>`)[1]
			yyze = strings.Split(yyze, `<span>%</span>`)[0]
			yyze = strings.ReplaceAll(yyze, `<td><p>`, "")
			fmt.Println("营业收入同比增长：", yyze)
			yyzef, err := strconv.ParseFloat(yyze, 64)
			if err != nil {
				yyzef = 0.0
			}
			stock := Stock{code: code, Jzcsyl: f, Mll: mllf, RevenueFrowth: yyzef}
			stocks = append(stocks, stock)
			pkg.ToTxt(code+","+s+","+mll+","+yyze, "300.csv")
		}
		time.Sleep(time.Duration(rand.Int63n(4)) * time.Second)

	}
}

func search002(res, url string, ctx context.Context) {
	m := make(map[string]string)
	pkg.ReadLine("002.csv", func(s string) {
		c := strings.Split(s, ",")[0]
		m[c] = "1"
	})
	for i := 2000; i <= 2999; i++ {
		//SZ 300000  创业板
		//SZ00 2000  中小板
		//SH 600000  601000 603000 沪市A
		//SZ000999 递减  深市A
		cs := []rune(fmt.Sprint(i))
		code := "SZ00" + string(cs)
		if _, ok := m[code]; ok {
			continue
		}
		chromedp.Navigate(strings.ReplaceAll(url, "%s", code)).Do(ctx)
		chromedp.InnerHTML(`div#app`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		if strings.Index(res, `净资产收益率`) > 0 {
			s := strings.Split(res, `净资产收益率</td>`)[1]
			s = strings.Split(s, `<span>%</span>`)[0]
			s = strings.ReplaceAll(s, `<td><p>`, "")
			fmt.Println(s)
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				f = 0.0
			}
			mll := strings.Split(res, `净利率</td>`)[1]
			mll = strings.Split(mll, `<span>%</span>`)[0]
			mll = strings.ReplaceAll(mll, `<td><p>`, "")
			//fmt.Println("毛利率：", mll)
			mllf, err := strconv.ParseFloat(mll, 64)
			if err != nil {
				mllf = 0.0
			}
			yyze := strings.Split(res, `营业收入同比增长</td>`)[1]
			yyze = strings.Split(yyze, `<span>%</span>`)[0]
			yyze = strings.ReplaceAll(yyze, `<td><p>`, "")
			fmt.Println("营业收入同比增长：", yyze)
			yyzef, err := strconv.ParseFloat(yyze, 64)
			if err != nil {
				yyzef = 0.0
			}
			/*if yyzef < 10 {
				return
			}*/
			stock := Stock{code: code, Jzcsyl: f, Mll: mllf, RevenueFrowth: yyzef}
			stocks = append(stocks, stock)
			pkg.ToTxt(code+","+s+","+mll+","+yyze, "002.csv")
		}
		time.Sleep(time.Duration(rand.Int63n(4)) * time.Second)
	}
}

func search600(res, url string, ctx context.Context) {
	m := make(map[string]string)
	pkg.ReadLine("600.csv", func(s string) {
		c := strings.Split(s, ",")[0]
		m[c] = "1"
	})
	for i := 600000; i <= 603999; i++ {
		//SZ 300000  创业板
		//SZ00 2000  中小板
		//SH 600000  601000 603000 沪市A
		if i >= 602000 && i < 603000 {
			continue
		}
		//SZ000999 递减  深市A
		cs := []rune(fmt.Sprint(i))
		code := "SH" + string(cs)
		if _, ok := m[code]; ok {
			continue
		}
		chromedp.Navigate(strings.ReplaceAll(url, "%s", code)).Do(ctx)
		chromedp.InnerHTML(`div#app`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		if strings.Index(res, `净资产收益率`) > 0 {
			s := strings.Split(res, `净资产收益率</td>`)[1]
			s = strings.Split(s, `<span>%</span>`)[0]
			s = strings.ReplaceAll(s, `<td><p>`, "")
			fmt.Println(s)
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				f = 0.0
			}
			mll := strings.Split(res, `净利率</td>`)[1]
			mll = strings.Split(mll, `<span>%</span>`)[0]
			mll = strings.ReplaceAll(mll, `<td><p>`, "")
			//fmt.Println("毛利率：", mll)
			mllf, err := strconv.ParseFloat(mll, 64)
			if err != nil {
				mllf = 0.0
			}
			yyze := strings.Split(res, `营业收入同比增长</td>`)[1]
			yyze = strings.Split(yyze, `<span>%</span>`)[0]
			yyze = strings.ReplaceAll(yyze, `<td><p>`, "")
			fmt.Println("营业收入同比增长：", yyze)
			yyzef, err := strconv.ParseFloat(yyze, 64)
			if err != nil {
				yyzef = 0.0
			}
			stock := Stock{code: code, Jzcsyl: f, Mll: mllf, RevenueFrowth: yyzef}
			stocks = append(stocks, stock)
			pkg.ToTxt(code+","+s+","+mll+","+yyze, "600.csv")
		}
		time.Sleep(time.Duration(rand.Int63n(4)) * time.Second)
	}
}

func search688(res, url string, ctx context.Context) {
	m := make(map[string]string)
	pkg.ReadLine("688.csv", func(s string) {
		c := strings.Split(s, ",")[0]
		m[c] = "1"
	})
	for i := 688000; i <= 688999; i++ {
		//SZ000999 递减  深市A
		cs := []rune(fmt.Sprint(i))
		code := "SH" + string(cs)
		if _, ok := m[code]; ok {
			continue
		}
		chromedp.Navigate(strings.ReplaceAll(url, "%s", code)).Do(ctx)
		chromedp.InnerHTML(`div#app`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		if strings.Index(res, `净资产收益率`) > 0 {
			s := strings.Split(res, `净资产收益率</td>`)[1]
			s = strings.Split(s, `<span>%</span>`)[0]
			s = strings.ReplaceAll(s, `<td><p>`, "")
			fmt.Println(s)
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				f = 0.0
			}
			mll := strings.Split(res, `净利率</td>`)[1]
			mll = strings.Split(mll, `<span>%</span>`)[0]
			mll = strings.ReplaceAll(mll, `<td><p>`, "")
			//fmt.Println("毛利率：", mll)
			mllf, err := strconv.ParseFloat(mll, 64)
			if err != nil {
				mllf = 0.0
			}
			yyze := strings.Split(res, `营业收入同比增长</td>`)[1]
			yyze = strings.Split(yyze, `<span>%</span>`)[0]
			yyze = strings.ReplaceAll(yyze, `<td><p>`, "")
			fmt.Println("营业收入同比增长：", yyze)
			yyzef, err := strconv.ParseFloat(yyze, 64)
			if err != nil {
				yyzef = 0.0
			}
			stock := Stock{code: code, Jzcsyl: f, Mll: mllf, RevenueFrowth: yyzef}
			stocks = append(stocks, stock)
			pkg.ToTxt(code+","+s+","+mll+","+yyze, "688.csv")
		}
		time.Sleep(time.Duration(rand.Int63n(4)) * time.Second)
	}
}

func search000(res, url string, ctx context.Context) {
	m := make(map[string]string)
	pkg.ReadLine("000.csv", func(s string) {
		c := strings.Split(s, ",")[0]
		m[c] = "1"
	})
	for i := 1000; i <= 1999; i++ {
		//SZ 300000  创业板
		//SZ00 2000  中小板
		//SH 600000  601000 603000 沪市A
		//SZ000999 递减  深市A
		cs := []rune(fmt.Sprint(i))
		code := "SZ000" + string(cs[1:])
		if _, ok := m[code]; ok {
			continue
		}
		chromedp.Navigate(strings.ReplaceAll(url, "%s", code)).Do(ctx)
		chromedp.InnerHTML(`div#app`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		if strings.Index(res, `净资产收益率`) > 0 {
			s := strings.Split(res, `净资产收益率</td>`)[1]
			s = strings.Split(s, `<span>%</span>`)[0]
			s = strings.ReplaceAll(s, `<td><p>`, "")
			fmt.Println(s)

			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				f = 0.0
			}
			mll := strings.Split(res, `净利率</td>`)[1]
			mll = strings.Split(mll, `<span>%</span>`)[0]
			mll = strings.ReplaceAll(mll, `<td><p>`, "")
			//fmt.Println("毛利率：", mll)
			mllf, err := strconv.ParseFloat(mll, 64)
			if err != nil {
				mllf = 0.0
			}
			yyze := strings.Split(res, `营业收入同比增长</td>`)[1]
			yyze = strings.Split(yyze, `<span>%</span>`)[0]
			yyze = strings.ReplaceAll(yyze, `<td><p>`, "")
			fmt.Println("营业收入同比增长：", yyze)
			yyzef, err := strconv.ParseFloat(yyze, 64)
			if err != nil {
				yyzef = 0.0
			}
			stock := Stock{code: code, Jzcsyl: f, Mll: mllf, RevenueFrowth: yyzef}
			stocks = append(stocks, stock)
			pkg.ToTxt(code+","+s+","+mll+","+yyze, "000.csv")
		}
		time.Sleep(time.Duration(rand.Int63n(4)) * time.Second)
	}
}

func C() {
	files := []string{"000.csv", "002.csv", "300.csv", "600.csv", "688.csv"}
	for _, file := range files {
		pkg.ReadLine(file, func(s string) {
			if s == "" {
				return
			}
			ss := strings.Split(s, `,`)
			v := ss[1]
			if strings.Contains(v, "-") {
				return
			}
			jzc, err := strconv.ParseFloat(v, 64)
			_ = jzc
			if err != nil {
				jzc = 0.0 //净资产收益率
			}
			v = ss[2]
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				f = 0.0 //净利率
			}
			v = ss[3]
			f, err = strconv.ParseFloat(v, 64)
			if err != nil {
				f = 0.0 //营收增速
			}
			if f < 16 {
				return
			}
			pkg.ToTxt(s, "code.csv")
		})
	}
}
