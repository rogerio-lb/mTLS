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
	TransferPacs002 Transfer `xml:"FIToFIPmtStsRpt"`
}

type Transfer struct {
	Amount                    Amount `xml:"CdtTrfTxInf>IntrBkSttlmAmt"`
	EndToEndId                string `xml:"CdtTrfTxInf>PmtId>EndToEndId"`
	MessageId                 string `xml:"GrpHdr>MsgId"`
	DebtorName                string `xml:"CdtTrfTxInf>Dbtr>Nm"`
	DebtorISPB                string `xml:"CdtTrfTxInf>DbtrAgt>FinInstnId>ClrSysMmbId>MmbId"`
	DebtorAccount             string `xml:"CdtTrfTxInf>DbtrAcct>Id>Othr>Id"`
	DebtorBranch              string `xml:"CdtTrfTxInf>DbtrAcct>Id>Othr>Issr"`
	ReturnID                  string `xml:"TxInf>RtrId"`
	OriginalEndToEndId        string `xml:"TxInf>OrgnlEndToEndId"`
	OriginalEndToEndIdPacs002 string `xml:"TxInfAndSts>OrgnlEndToEndId"`
	OriginalInstructionId     string `xml:"TxInfAndSts>OrgnlInstrId"`
	TransactionStatus         string `xml:"TxInfAndSts>TxSts"`
	ReturnReason              string `xml:"TxInf>RtrRsnInf>Rsn>Cd"`
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
		_ = services.FinishStream(conn, response.PIPullNext)
		panic("Error getting messages: %v\n" + err.Error())
		return
	}

	handleIncomingMessage(response.Message)

	for {
		time.Sleep(1 * time.Second)
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

	handleMessage(parseMessage, message)
}

func respondPacs008(message Envelope) {
	e2eID := message.Document.TransferPacs008.EndToEndId

	pacs002 := services.GeneratePacs002(e2eID)

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

func handleMessage(message Envelope, rawMessage string) {
	if message.XMLName.Space == "https://www.bcb.gov.br/pi/pacs.002/1.15" {
		fmt.Println()
		fmt.Println("==================================== Received PACS002 ====================================")
		fmt.Println("Original E2E ID:", message.Document.TransferPacs002.OriginalEndToEndIdPacs002)
		fmt.Println("Original Instruction ID:", message.Document.TransferPacs002.OriginalInstructionId)
		fmt.Println("Transaction Status:", message.Document.TransferPacs002.TransactionStatus)
		fmt.Println("==================================== End of PACS002 ====================================")

		if message.Document.TransferPacs002.TransactionStatus == "RJCT" {
			fmt.Println("Transaction was rejected. Reason:", rawMessage)
		}

		return
	}

	if message.XMLName.Space == "https://www.bcb.gov.br/pi/pacs.008/1.14" {
		fmt.Println()
		fmt.Println("==================================== Received PACS008 ====================================")
		fmt.Println("E2E ID:", message.Document.TransferPacs008.EndToEndId)
		fmt.Println("Message ID:", message.Document.TransferPacs008.MessageId)
		fmt.Println("Amount:", message.Document.TransferPacs008.Amount.Value, message.Document.TransferPacs008.Amount.Currency)
		fmt.Println("Debtor Name: ", message.Document.TransferPacs008.DebtorName)
		fmt.Println("Debtor ISPB: ", message.Document.TransferPacs008.DebtorISPB)
		fmt.Println("Debtor Account: ", message.Document.TransferPacs008.DebtorAccount)
		fmt.Println("Debtor Branch: ", message.Document.TransferPacs008.DebtorBranch)
		fmt.Println("==================================== End of PACS008 ====================================")

		respondPacs008(message)

		return
	}

	if message.XMLName.Space == "https://www.bcb.gov.br/pi/pacs.004/1.5" {
		fmt.Println()
		fmt.Println("==================================== Received PACS004 ====================================")
		fmt.Println("Original E2E ID:", message.Document.TransferPacs004.EndToEndId)
		fmt.Println("Return ID:", message.Document.TransferPacs004.ReturnID)
		fmt.Println("Message ID:", message.Document.TransferPacs004.MessageId)
		fmt.Println("Debtor ISPB: ", message.Document.TransferPacs004.DebtorISPB)
		fmt.Println("Reason for Return: ", message.Document.TransferPacs004.ReturnID)
		fmt.Println("==================================== End of PACS004 ====================================")

		respondPacs004(message)
	}

	fmt.Println("Received: ", message.XMLName.Space)
	fmt.Println("Full Message: ", rawMessage)
}
