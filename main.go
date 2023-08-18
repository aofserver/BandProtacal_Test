package main

import (
	"fmt"
	"time"
	"encoding/json"
	"github.com/valyala/fasthttp"
)


type BodyBroadcast struct {
	Symbol  string `json:"symbol"`
	Price   int    `json:"price"`
	Timestamp int `json:"timestamp"`
}



func main() {
	unixTimestamp := time.Now().Unix()
	fmt.Println("UnixTimestamp:", int(unixTimestamp))
	broadcastData := BodyBroadcast{
		Symbol:  "BTC",
		Price:   100000,
		Timestamp: int(unixTimestamp),
	}

	tx := BroadcastTransaction(broadcastData)
	fmt.Println("TX:", tx)

	for i := 0; i < 60; i++ {
		tx_status := GetTransaction(tx)
		fmt.Println("TX STATUS:", tx_status)
		if (tx_status != "PENDING"){
			break
		}
		time.Sleep(1*time.Second)
	}
}


func BroadcastTransaction(bodyBroadcast BodyBroadcast) string {
	url := "https://mock-node-wgqbnxruha-as.a.run.app/broadcast"
	jsonPayload, err_p := json.Marshal(bodyBroadcast)
	if err_p != nil {
		fmt.Println("Error:", err_p)
		return "Error"
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetRequestURI(url)
	req.SetBody(jsonPayload)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	if err := fasthttp.Do(req, resp); err != nil {
		fmt.Println("Error:", err)
		return "Error"
	}

	type Response struct {
		Tx_hash  string `json:"tx_hash"`
	}
	body := resp.Body()
	var response Response
	err := json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error:", err)
		return "Error"
	}
	return response.Tx_hash
}


func GetTransaction(tx string) string {
	url := "https://mock-node-wgqbnxruha-as.a.run.app/check/" + tx
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := fasthttp.Do(req, resp); err != nil {
		fmt.Println("Error:", err)
		return "Error"
	}

	type Response struct {
		Tx_status  string `json:"tx_status"`
	}
	body := resp.Body()
	var response Response
	err := json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error:", err)
		return "Error"
	}
	return response.Tx_status
}