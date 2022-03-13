package main

import (
	"log"
	"os"

	"github.com/frankrap/bybit-api/rest"
)

const (
	testNetBaseURL = "https://api-testnet.bybit.com/" // 测试网络
	mainNetBaseURL = "https://api.bybit.com/"         // 主网络
)

func main() {
	baseURL := mainNetBaseURL
	apiKey := os.Getenv("API_KEY")
	secretKey := os.Getenv("SECRET_KEY")
	debugMode := true

	b := rest.New(nil, baseURL, apiKey, secretKey, debugMode)

	// get all positions:
	getPositionsQuery, getPositionsResp, positionData, err := b.GetPositions()
	if err != nil {
		log.Printf("%v", err)

		return
	}

	// log.Printf("getPositionsQuery: %#v \n", getPositionsQuery)
	// log.Printf("getPositionsResp: %#v \n", getPositionsResp)
	_ = getPositionsQuery
	_ = getPositionsResp

	for _, position := range positionData {
		log.Printf("position: %#v \n", position)
	}

	// create order:
	symbol := "BTCUSD"
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
