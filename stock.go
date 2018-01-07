package gofin

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var client = http.DefaultClient

type HSStockData struct {
	Name         string  /*[1]股票名称*/
	Gid          string  /*[2]股票编号*/
	NowPri       float64 /*[3]当前价格*/
	YestClosePri float64 /*[4]昨日收盘价*/
	OpeningPri   float64 /*[5]今日开盘价*/
	TraNumber    int64   /*[6]成交量*/
	Outter       int64   /*[7]外盘*/
	Inner        int64   /*[8]内盘*/
	BuyOne       int64   /*[9]买一报价*/
	BuyOnePri    float64 /*[10]买一*/
	BuyTwo       int64   /*[11]买二*/
	BuyTwoPri    float64 /*[12]买二报价*/
	BuyThree     int64   /*[13]买三*/
	BuyThreePri  float64 /*[14]买三报价*/
	BuyFour      int64   /*[15]买四*/
	BuyFourPri   float64 /*[16]买四报价*/
	BuyFive      int64   /*[17]买五*/
	BuyFivePri   float64 /*[18]买五报价*/
	SellOne      int64   /*[19]卖一*/
	SellOnePri   float64 /*[20]卖一报价*/
	SellTwo      int64   /*[21]卖二*/
	SellTwoPri   float64 /*[22]卖二报价*/
	SellThree    int64   /*[23]卖三*/
	SellThreePri float64 /*[24]卖三报价*/
	SellFour     int64   /*[25]卖四*/
	SellFourPri  float64 /*[26]卖四报价*/
	SellFive     int64   /*[27]卖五*/
	SellFivePri  float64 /*[28]卖五报价*/
	/*[29]最近逐笔成交*/
	Time      string  /*[30]时间*/
	Change    float64 /*[31]涨跌*/
	ChangePer float64 /*[32]涨跌%*/
	YodayMax  float64 /*[33]今日最高价*/
	YodayMin  float64 /*[34]今日最低价*/
	/*[35]价格/成交量（手）/成交额*/
	TradeCount int64   /*[36]成交量*/
	TradeAmont int64   /*[37]成交额*/
	ChangeRate float64 /*[38]换手率*/
	PERatio    float64 /*[39]市盈率*/
	/*[40]*/
	/*[41]最高*/
	/*[42]最低*/
	MaxMinChange float64 /*[43]振幅*/
	MarketAmont  float64 /*[44]流通市值*/
	TotalAmont   float64 /*[45]总市值*/
	PBRatio      float64 /*[46]市净率*/
	HighPri      float64 /*[47]涨停价*/
	LowPri       float64 /*[48]跌停价*/
}

const (
	URL_Last_Price = "http://sqt.gtimg.cn/utf8/q=%s" // 最后一个交易日信息
)

var ErrMalformedData = errors.New("malformed data")

// http://sqt.gtimg.cn/utf8/q=sh600000
func parseStockData(s string) (*HSStockData, error) {
	kv := strings.Split(s, "=")
	if len(kv) != 2 {
		return nil, ErrMalformedData
	}

	rawData := strings.Trim(kv[1], "\";")
	cols := strings.Split(rawData, "~")
	if len(cols) != 54 {
		return nil, ErrMalformedData
	}

	d := new(HSStockData)
	d.Name = cols[1]
	d.Gid = cols[2]
	d.NowPri, _ = strconv.ParseFloat(cols[3], 64)
	d.YestClosePri, _ = strconv.ParseFloat(cols[4], 64)
	d.OpeningPri, _ = strconv.ParseFloat(cols[5], 64)
	d.TraNumber, _ = strconv.ParseInt(cols[6], 10, 64)
	d.Outter, _ = strconv.ParseInt(cols[7], 10, 64)
	d.Inner, _ = strconv.ParseInt(cols[8], 10, 64)
	d.BuyOnePri, _ = strconv.ParseFloat(cols[9], 64)
	d.BuyOne, _ = strconv.ParseInt(cols[10], 10, 64)
	d.BuyTwoPri, _ = strconv.ParseFloat(cols[11], 64)
	d.BuyTwo, _ = strconv.ParseInt(cols[12], 10, 64)
	d.BuyThreePri, _ = strconv.ParseFloat(cols[13], 64)
	d.BuyThree, _ = strconv.ParseInt(cols[14], 10, 64)
	d.BuyFourPri, _ = strconv.ParseFloat(cols[15], 64)
	d.BuyFour, _ = strconv.ParseInt(cols[16], 10, 64)
	d.BuyFivePri, _ = strconv.ParseFloat(cols[17], 64)
	d.BuyFive, _ = strconv.ParseInt(cols[18], 10, 64)
	d.SellOnePri, _ = strconv.ParseFloat(cols[19], 64)
	d.SellOne, _ = strconv.ParseInt(cols[20], 10, 64)
	d.SellTwoPri, _ = strconv.ParseFloat(cols[21], 64)
	d.SellTwo, _ = strconv.ParseInt(cols[22], 10, 64)
	d.SellThreePri, _ = strconv.ParseFloat(cols[23], 64)
	d.SellThree, _ = strconv.ParseInt(cols[24], 10, 64)
	d.SellFourPri, _ = strconv.ParseFloat(cols[25], 64)
	d.SellFour, _ = strconv.ParseInt(cols[26], 10, 64)
	d.SellFivePri, _ = strconv.ParseFloat(cols[27], 64)
	d.SellFive, _ = strconv.ParseInt(cols[28], 10, 64)
	d.Time = cols[30]
	d.Change, _ = strconv.ParseFloat(cols[31], 64)
	d.ChangePer, _ = strconv.ParseFloat(cols[32], 64)
	d.YodayMax, _ = strconv.ParseFloat(cols[33], 64)
	d.YodayMin, _ = strconv.ParseFloat(cols[34], 64)
	d.TradeCount, _ = strconv.ParseInt(cols[36], 10, 64)
	d.TradeAmont, _ = strconv.ParseInt(cols[37], 10, 64)
	d.ChangeRate, _ = strconv.ParseFloat(cols[38], 64)
	d.PERatio, _ = strconv.ParseFloat(cols[39], 64)
	d.MaxMinChange, _ = strconv.ParseFloat(cols[43], 64)
	d.MarketAmont, _ = strconv.ParseFloat(cols[44], 64)
	d.TotalAmont, _ = strconv.ParseFloat(cols[45], 64)
	d.PBRatio, _ = strconv.ParseFloat(cols[46], 64)
	d.HighPri, _ = strconv.ParseFloat(cols[47], 64)
	d.LowPri, _ = strconv.ParseFloat(cols[48], 64)
	return d, nil
}

func GetLastPrice(code string) (*HSStockData, error) {
	resp, err := client.Get(fmt.Sprintf(URL_Last_Price, code))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseStockData(string(body))
}
