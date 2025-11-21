package main

import (
	"encoding/xml"
	"fmt"
	"mTLS/services"
	"time"
)

type Envelope struct {
	XMLName  xml.Name `xml:"Envelope"`
	Document Document `xml:"Document"`
}

type Document struct {
	Transfer Transfer `xml:"FIToFICstmrCdtTrf"`
}

type Transfer struct {
	Amount     Amount `xml:"CdtTrfTxInf>IntrBkSttlmAmt"`
	EndToEndId string `xml:"CdtTrfTxInf>PmtId>EndToEndId"`
	DebtorName string `xml:"CdtTrfTxInf>Dbtr>Nm"`
}

type Amount struct {
	Currency string `xml:"Ccy,attr"`
	Value    string `xml:",chardata"`
}

func main() {
	conn := services.CreateConnection()
	fmt.Println("=== TLS Connection Test Completed ===")
	fmt.Println("Sending HTTP/1.0 request...")

	var response *services.GetMessageResponse

	response, err := services.GetMessages(conn, "start")
	if err != nil {
		services.FinishStream(conn, response.PIPullNext)
		panic("Error getting messages: %v\n" + err.Error())
		return
	}

	handleIncomingMessage(response.Message)

	for {
		time.Sleep(2 * time.Second)
		connection := services.CreateConnection()
		response, err = services.GetMessages(connection, response.PIPullNext)
		if err != nil {
			//			services.FinishStream(connection, response.PIPullNext)
			fmt.Printf("Error getting messages: %v\n", err)
			break
		}

		handleIncomingMessage(response.Message)
	}

	//services.FinishStream(conn, "/api/v1/out/52833288/stream/11e5d52c-e302-4aec-af0a-83b3f068d26a")
}

func handleIncomingMessage(message string) {
	var parseMessage Envelope

	if message == "" {
		return
	}

	err := xml.Unmarshal([]byte(message), &parseMessage)
	if err != nil {
		fmt.Printf("Error parsing XML: %v\n", err)
		return
	}

	fmt.Printf("Parsed Message:\n")
	fmt.Println("Message Type: ", parseMessage.XMLName.Space)
}
