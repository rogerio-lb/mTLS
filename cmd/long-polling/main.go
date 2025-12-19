package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"mTLS/services"
	"mime/multipart"
	"os"
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

const message = `<Envelope xmlns="https://www.bcb.gov.br/pi/pacs.008/1.14"><AppHdr><Fr><FIId><FinInstnId><Othr><Id>52833288</Id></Othr></FinInstnId></FIId></Fr><To><FIId><FinInstnId><Othr><Id>00038166</Id></Othr></FinInstnId></FIId></To><BizMsgIdr>M52833288YO8s5ar29g4nAH9QdIetTbw</BizMsgIdr><MsgDefIdr>pacs.008.spi.1.14</MsgDefIdr><CreDt>2025-12-18T17:17:59.054Z</CreDt><Sgntr><Signature xmlns=""><SignedInfo><CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></CanonicalizationMethod><SignatureMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"></SignatureMethod><Reference URI=""><Transforms><Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></Transform><Transform Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature"></Transform><Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></Transform><Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></Transform></Transforms><DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"></DigestMethod><DigestValue>hQbAz2HJPNk5p1AnyCkVYfzu9M4qKwDlq1JaTENf6N0=</DigestValue></Reference></SignedInfo><SignatureValue>U58TmfqttDGD09uxAU9XbTeQgAuyUsNI0VctWRoqWcS5mOZTd6hhPNmNXkcVFSf37Ijo5tyGTeEZ&#xD;&#xA;O+tkhWfv3TnV95XjC3bxBOGEYfqmoSzYpUSPClB+kTJ3UPMWEQTbzSgZdCkrLNjCwW34jT/3vRoH&#xD;&#xA;EHQSGc4hPyfnOfh8ByEZDQbYpaljH5D0WE3Eqab6c6wclK0NX/Hu7Buoeh3d5aW9kthWuP0urV/1&#xD;&#xA;aEPajnJ89KYE8f0vR5phi2mY7DjWsrX6C1glt59pSrTgF+J6XzL2aMbULl1pMbVi76Qz+FDK6lLZ&#xD;&#xA;oCHNYRGPASOF9kwH4rlX/iXCi9Ay4P8KVEgUhw==</SignatureValue><KeyInfo><X509Data><X509Certificate></X509Certificate></X509Data></KeyInfo></Signature></Sgntr></AppHdr><Document><FIToFICstmrCdtTrf><GrpHdr><MsgId>M52833288YO8s5ar29g4nAH9QdIetTbw</MsgId><CreDtTm>2025-12-18T17:17:59.054Z</CreDtTm><NbOfTxs>1</NbOfTxs><SttlmInf><SttlmMtd>CLRG</SttlmMtd></SttlmInf><PmtTpInf><InstrPrty>HIGH</InstrPrty><SvcLvl><Prtry>PAGPRI</Prtry></SvcLvl></PmtTpInf></GrpHdr><CdtTrfTxInf><PmtId><EndToEndId>E52833288202512181717wkl9T080i3Y</EndToEndId></PmtId><IntrBkSttlmAmt Ccy="BRL">10.00</IntrBkSttlmAmt><AccptncDtTm>2025-12-18T17:17:59.054Z</AccptncDtTm><ChrgBr>SLEV</ChrgBr><MndtRltdInf><Tp><LclInstrm><Prtry>MANU</Prtry></LclInstrm></Tp></MndtRltdInf><Dbtr><Nm>Rogerio Inacio</Nm><Id><PrvtId><Othr><Id>14811554744</Id></Othr></PrvtId></Id></Dbtr><DbtrAcct><Id><Othr><Id>0038952</Id><Issr>0001</Issr></Othr></Id><Tp><Cd>CACC</Cd></Tp></DbtrAcct><DbtrAgt><FinInstnId><ClrSysMmbId><MmbId>52833288</MmbId></ClrSysMmbId></FinInstnId></DbtrAgt><CdtrAgt><FinInstnId><ClrSysMmbId><MmbId>04902979</MmbId></ClrSysMmbId></FinInstnId></CdtrAgt><Cdtr><Id><PrvtId><Othr><Id>43528405058</Id></Othr></PrvtId></Id></Cdtr><CdtrAcct><Id><Othr><Id>003816482</Id><Issr>0001</Issr></Othr></Id><Tp><Cd>CACC</Cd></Tp></CdtrAcct><Purp><Cd>IPAY</Cd></Purp><RmtInf><Ustrd>Teste de envio</Ustrd></RmtInf></CdtTrfTxInf></FIToFICstmrCdtTrf></Document></Envelope>`

const debug = false

func main() {
	ISPB := os.Getenv("ISPB")
	//conn := services.CreateConnection(debug)
	conn := services.CreateConnectionV2(debug, false)

	var response *services.GetMessageResponse

	response, err := services.GetMessages(conn, "start", ISPB)
	if err != nil {
		//_ = services.FinishStream(conn, response.PIPullNext)
		panic("Error getting messages: %v\n" + err.Error())
		return
	}

	handleIncomingMessage(response.Message, ISPB)

	for {
		time.Sleep(1 * time.Second)
		connection := services.CreateConnectionV2(debug, false)
		response, err = services.GetMessages(connection, response.PIPullNext, ISPB)
		if err != nil {
			//			services.FinishStream(connection, response.PIPullNext)
			fmt.Printf("Error getting messages: %v\n", err)
			break
		}

		go handleIncomingMessage(response.Message, ISPB)
	}

	//services.FinishStream(conn, "/api/v1/out/52833288/stream/11e5d52c-e302-4aec-af0a-83b3f068d26a")
}

/*func main() {
	handleIncomingMessage(message)
}*/

func handleIncomingMessage(message, ispb string) {
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

	handleMessage(parseMessage, message, ispb)
}

func respondPacs008(message Envelope, ispb string) {
	e2eID := message.Document.TransferPacs008.EndToEndId

	pacs002 := services.GeneratePacs002(e2eID, ispb)

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

	conn := services.CreateConnectionV2(debug, false)

	err = services.PostMessage(conn, string(compressedMessage.Bytes()), mw.Boundary(), ispb)
	if err != nil {
		fmt.Println("Error posting message:", err)
	}
}

func respondPacs004(message Envelope, ispb string) {
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

	conn := services.CreateConnectionV2(debug, false)

	err = services.PostMessage(conn, string(compressedMessage.Bytes()), mw.Boundary(), ispb)
	if err != nil {
		fmt.Println("Error posting message:", err)
	}
}

func handleMessage(message Envelope, rawMessage, ispb string) {
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

		fmt.Println("Raw Message: ", rawMessage)

		respondPacs008(message, ispb)

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

		respondPacs004(message, ispb)
	}

	fmt.Println("Received: ", message.XMLName.Space)
	fmt.Println("Full Message: ", rawMessage)
}
