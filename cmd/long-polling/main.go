package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"mTLS/services"
	"mime/multipart"
	"time"
)

type Envelope struct {
	XMLName  xml.Name `xml:"Envelope"`
	Document Document `xml:"Document"`
}

type Document struct {
	TransferPacs008 Transfer `xml:"FIToFICstmrCdtTrf"`
	TransferPacs004 Transfer `xml:"PmtRtr"`
}

type Transfer struct {
	Amount             Amount `xml:"CdtTrfTxInf>IntrBkSttlmAmt"`
	EndToEndId         string `xml:"CdtTrfTxInf>PmtId>EndToEndId"`
	DebtorName         string `xml:"CdtTrfTxInf>Dbtr>Nm"`
	ReturnID           string `xml:"TxInf>RtrId"`
	OriginalEndToEndId string `xml:"TxInf>OrgnlEndToEndId"`
}
type Amount struct {
	Currency string `xml:"Ccy,attr"`
	Value    string `xml:",chardata"`
}

const debug = false

func main() {
	conn := services.CreateConnection(debug)
	/*fmt.Println("=== TLS Connection Test Completed ===")
	fmt.Println("Sending HTTP/1.0 request...")*/

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
		connection := services.CreateConnection(debug)
		response, err = services.GetMessages(connection, response.PIPullNext)
		if err != nil {
			//			services.FinishStream(connection, response.PIPullNext)
			fmt.Printf("Error getting messages: %v\n", err)
			break
		}

		go handleIncomingMessage(response.Message)
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

	/*fmt.Printf("Parsed Message:\n")
	fmt.Println("Message Type: ", parseMessage.XMLName.Space)*/

	if parseMessage.XMLName.Space == "https://www.bcb.gov.br/pi/pacs.008/1.14" {
		respondPacs008(parseMessage)
	}

	if parseMessage.XMLName.Space == "https://www.bcb.gov.br/pi/pacs.004/1.5" {
		respondPacs004(parseMessage)
	}
}

func respondPacs008(message Envelope) {
	e2eID := message.Document.TransferPacs008.EndToEndId

	pacs002 := services.GeneratePacs002(e2eID)

	fmt.Println("E2E ID:", e2eID)

	var responseContent bytes.Buffer

	mw := multipart.NewWriter(&responseContent)

	err := services.AddXMLPart(mw, pacs002)
	if err != nil {
		panic("Error adding XML part: " + err.Error())
	}

	mw.Close()

	var compressedMessage bytes.Buffer

	err = services.CompressContentToGzip(responseContent.Bytes(), &compressedMessage)
	if err != nil {
		panic("Error adding XML part: " + err.Error())
	}

	conn := services.CreateConnection(debug)

	err = services.PostMessage(conn, string(compressedMessage.Bytes()), mw.Boundary())
	if err != nil {
		fmt.Println("Error posting message:", err)
	}
}

func respondPacs004(message Envelope) {
	e2eID := message.Document.TransferPacs004.OriginalEndToEndId
	returnId := message.Document.TransferPacs004.ReturnID

	pacs002 := services.GeneratePacs002ForPacs004(e2eID, returnId)

	var responseContent bytes.Buffer

	mw := multipart.NewWriter(&responseContent)

	err := services.AddXMLPart(mw, pacs002)
	if err != nil {
		panic("Error adding XML part: " + err.Error())
	}

	mw.Close()

	var compressedMessage bytes.Buffer

	err = services.CompressContentToGzip(responseContent.Bytes(), &compressedMessage)
	if err != nil {
		panic("Error adding XML part: " + err.Error())
	}

	conn := services.CreateConnection(debug)

	err = services.PostMessage(conn, string(compressedMessage.Bytes()), mw.Boundary())
	if err != nil {
		fmt.Println("Error posting message:", err)
	}
}
