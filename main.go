package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TerminusDeus/bybit-api/rest"
)

const (
	testNetBaseURL = "https://api-testnet.bybit.com/" // 测试网络
	mainNetBaseURL = "https://api.bybit.com/"         // 主网络
)

var client = &http.Client{Timeout: 10 * time.Second}

func main() {
	_, _ = getSpotSymbolList()

	baseURL := mainNetBaseURL
	apiKey := os.Getenv("API_KEY")
	secretKey := os.Getenv("SECRET_KEY")
	debugMode := false

	b := rest.New(client, baseURL, apiKey, secretKey, debugMode)

	// get all positions:
	getPositionsQuery, getPositionsResp, positions, err := b.GetPositions()
	if err != nil {
		log.Printf("GetPositions failed with an error: %s", err.Error())

		return
	}

	// log.Printf("getPositionsQuery: %#v \n", getPositionsQuery)
	// log.Printf("getPositionsResp: %#v \n", getPositionsResp)
	_ = getPositionsQuery
	_ = getPositionsResp

	for _, position := range positions {
		log.Printf("position: %#v \n\n", position)
	}

	// get list of futures symbols:
	getSymbolsQuery, getSymbolsResp, symbols, err := b.GetSymbols()
	if err != nil {
		log.Printf("GetSymbols failed with an error: %s", err.Error())

		return
	}

	// log.Printf("getSymbolsQuery: %#v \n", getSymbolsQuery)
	// log.Printf("getSymbolsResp: %#v \n", getSymbolsResp)
	_ = getSymbolsQuery
	_ = getSymbolsResp

	for i, symbol := range symbols {
		if i == len(symbols)-1 || i == len(symbols)-2 {

			log.Printf("symbol: Base Quote: %s %s \n\n", symbol.BaseCurrency, symbol.QuoteCurrency)
		}
	}

	symbol := "ADAUSDT"

	// get order book for given futures symbol:
	getOrderBookQuery, getOrderBookResp, orderBook, err := b.GetOrderBook(symbol)
	if err != nil {
		log.Printf("GetOrderBook failed with an error: %s", err.Error())

		return
	}

	getOrderBookResponse := struct {
		RetCode int         `json:"ret_code"`
		RetMsg  string      `json:"ret_msg"`
		ExtCode string      `json:"ext_code"`
		ExtInfo string      `json:"ext_info"`
		Result  interface{} `json:"result"`
		TimeNow string      `json:"time_now"`
	}{}

	if err = json.Unmarshal(getOrderBookResp, &getOrderBookResponse); err != nil {
		fmt.Printf("unmarshalling GetOrderBook resp body as bytes failed: %s", err.Error())

		return
	}
	if getOrderBookResponse.RetMsg != "OK" {
		fmt.Printf("getOrderBookResponse.RetMsg is non empty: %s", getOrderBookResponse.RetMsg)

		return
	}

	// log.Printf("getOrderBookQuery: %#v \n", getOrderBookQuery)
	// log.Printf("getOrderBookResp: %#v \n", getOrderBookResp)
	_ = getOrderBookQuery
	_ = getOrderBookResp

	for i, orderBookAsk := range orderBook.Asks {
		if i == len(orderBook.Asks)-1 || i == len(orderBook.Asks)-2 {
			log.Printf("orderBookAsk = %+v\n\n", orderBookAsk)
		}
	}

	// create order:
	symbol = "BTCUSD"
	side := "Buy"
	orderType := "Limit"
	qty := 30
	price := 7000.0
	timeInForce := "GoodTillCancel"
	reduceOnly := false
	closeOnTrigger := false
	orderLinkID := ""
	var stopLoss, takeProfit float64
	createOrderQuery, createOrderResp, order, err := b.CreateOrder(side, orderType, price, qty, timeInForce, takeProfit, stopLoss, reduceOnly, closeOnTrigger, orderLinkID, symbol)
	if err != nil {
		log.Println(err)

		return
	}

	// log.Printf("createOrderQuery: %#v \n", createOrderQuery)
	// log.Printf("createOrderResp: %#v \n", createOrderResp)
	_ = createOrderQuery
	_ = createOrderResp

	log.Printf("created order: %#v \n", order)
}

func getSpotSymbolList() (symbolsList []string, err error) {
	url := "https://api2.bybit.com/spot/api/basic/symbol_list"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("getSpotSymbolList: init request to bybit service failed: %s", err.Error())

		return
	}

	req.Header.Add("Content-Type", "application/json")

	newResp, err := client.Do(req)
	if err != nil {
		fmt.Printf("getSpotSymbolList: request to bybit service failed: %s", err.Error())

		return
	}

	defer newResp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(newResp.Body)
	if err != nil {
		fmt.Printf("getSpotSymbolList: ReadAll response body failed: %s", err.Error())

		return
	}

	if newResp.StatusCode != http.StatusOK {
		fmt.Printf("getSpotSymbolList: request to bybit service status is not OK: status code: %v, status: %v", newResp.StatusCode, newResp.Status)

		return
	}

	spotSymbolListResp := struct{}{}

	if err = json.Unmarshal(bodyBytes, &spotSymbolListResp); err != nil {
		fmt.Printf("getSpotSymbolList: unmarshalling resp body bytes as structured data failed: %s", err.Error())

		return
	}

	return
}
