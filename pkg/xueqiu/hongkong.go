package xueqiu

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"math/rand"
	"spider/pkg"
	"strings"
	"time"
)

func RunX() {
	//getFFCX()
	pkg.ReadLine("hongkong-result-1.csv", func(s string) {
		ss := strings.Split(s, ",")
		code := strings.ReplaceAll(ss[0], "H ", "0")
		l := 5 - len(code)
		for i := 1; i <= l; i++ {
			code = "0" + code
		}

		pkg.ToTxt("h"+code+","+ss[1]+","+ss[2], "hongkong-result-2.csv")
	})
}

func getFFCX() {
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
			getHongKongFcf("hongkang.csv", res, ctx)
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

func getHongKongFcf(name, res string, ctx context.Context) {
	exits := make(map[string]string)
	pkg.ReadLine("hongkong-result.csv", func(s string) {
		exits[strings.Split(s, ",")[0]] = "1"
	})
	count := 0
	pkg.ReadLine(name, func(s string) {
		if count > 0 {
			fmt.Println(s)
			//return
		}
		count++
		if _, ok := exits[s]; ok {
			return
		}
		chromedp.Navigate("https://www.gurufocus.cn/stock/HKSE:" + s + "/term/price_to_free_cash_flow").Do(ctx)
		chromedp.OuterHTML(`body`, &res, chromedp.ByQuery).Do(ctx)
		//fmt.Println(res)
		if strings.Contains(res, `没有找到对应的股票`) || strings.Contains(res, `服务器内部错误`) {
			pkg.ToTxt(s+",10000000"+",0", "hongkong-result.csv")
			return
		}
		chromedp.InnerHTML(`div#term-page-description`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		fmt.Println(res)
		ss := strings.Split(res, `每股自由现金流 </a>    为 HK$ `)
		ss = strings.Split(ss[1], `，`)
		fcf := ss[0]
		ss = strings.Split(ss[1], ` 为 `)
		priceFcf := strings.ReplaceAll(ss[1], `。`, "")
		fmt.Println("priceFcf:", priceFcf, "fcf:", fcf)
		pkg.ToTxt(s+","+priceFcf+","+fcf, "hongkong-result.csv")
		time.Sleep(time.Duration(rand.Int63n(5)) * time.Second)

	})
}
