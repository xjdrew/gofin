package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
)

const DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"

type Quote struct {
	Symbol        string `json:"symbol"`        // 股票代号
	Exchange      string `json:"exchange"`      // 交易所代码
	Code          string `json:"code"`          // 代码
	Name          string `json:"name"`          // 名称
	Current       string `json:"current"`       // 当前价
	Percentage    string `json:"percentage"`    // 当日涨跌百分比
	Change        string `json:"change"`        // 当日涨跌
	Open          string `json:"open"`          // 开盘价
	High          string `json:"high"`          // 最高价
	Low           string `json:"low"`           // 最低价
	Close         string `json:"close"`         // 收盘价
	LastClose     string `json:"last_close"`    // 上一日收盘价
	High52week    string `json:"high52week"`    // 52周最高价
	Low52week     string `json:"low52week"`     // 52周最低价
	Volume        string `json:"volume"`        // 成交量
	MarketCapital string `json:"marketCapital"` // 总市值
	EPS           string `json:"eps"`           // 每股收益
	PE_TTM        string `json:"pe_ttm"`        // 动态市盈率
	PE_LYR        string `json:"pe_lyr"`        // 静态市盈率
	Time          string `json:"time"`          // 更新时间
	UpdateAt      string `json:"updateAt"`      // 更新时间UTC
	TurnoverRate  string `json:"turnover_rate"` // 换手率
	CurrencyUnit  string `json:"currency_unit"` // 货币单位
	Amount        string `json:"amount"`        // 成交额
	NetAssets     string `json:"net_assets"`    // 净资产
	PB            string `json:"pb"`            // 市净率
}

type QuoteResponse map[string]*Quote

type XQClient struct {
	*http.Client
}

func (xqc *XQClient) newRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", DefaultUserAgent)
	return req, err
}

func (xqc *XQClient) GetQuote(code string) (*Quote, error) {
	url := fmt.Sprintf("https://xueqiu.com/v4/stock/quote.json?code=%s", code)
	req, err := xqc.newRequest(url)
	if err != nil {
		return nil, err
	}

	resp, err := xqc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	quoteResponse := make(QuoteResponse)

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&quoteResponse); err != nil {
		return nil, err
	}
	return quoteResponse[code], nil
}

// 加载雪球首页，为client设置正确的cookie
func (xqc *XQClient) Init() error {
	req, err := xqc.newRequest("https://xueqiu.com")
	if err != nil {
		return err
	}

	resp, err := xqc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func newClient() *XQClient {
	jar, _ := cookiejar.New(nil)
	xqc := &XQClient{&http.Client{Jar: jar}}
	return xqc
}

func parseCodes() ([]string, error) {
	var err error
	var input *os.File
	f := flag.Arg(0)
	if f != "" {
		input, err = os.Open(f)
		if err != nil {
			return nil, err
		}
		defer input.Close()
	} else {
		input = os.Stdin
	}

	chunk, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(chunk), "\n")
	var codes []string
	for _, line := range lines {
		code := strings.TrimSpace(line)
		if code != "" && !strings.HasPrefix(code, "#") {
			codes = append(codes, code)
		}
	}
	return codes, nil
}

func main() {
	flag.Parse()

	codes, err := parseCodes()
	if err != nil {
		log.Fatal("parseCodes failed:", err)
	}

	client := newClient()
	if err := client.Init(); err != nil {
		log.Fatal("Init failed:", err)
	}

	fmt.Println("名称, 市净率")
	for _, code := range codes {
		quote, err := client.GetQuote(code)
		if err != nil {
			log.Fatal(err)
		}
		name := quote.Name
		if quote.Exchange == "HK" {
			name = name + "H"
		}
		fmt.Printf("%s, %s\n", name, quote.PB)
	}
}
