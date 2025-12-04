package main

import (
	"bytes"
	"fmt"
	"mTLS/services"
	"mime/multipart"
)

const debug = false

func main() {
	conn := services.CreateConnection(debug)

	if conn == nil {
		panic("Failed to create TLS connection")
	}

	var responseContent bytes.Buffer

	mw := multipart.NewWriter(&responseContent)

	for i := 0; i < 1; i++ {
		//message := services.CreateMessage()
		/*message := services.GeneratePacs004ForDict(
			"E54811417202512041739bJrrfAFhvUy",
			"54811417",
			"10.00",
			services.FRAUD_REASON,
		)*/
		message := services.GeneratePacs008Manual(
			"99999004",
			"003816482",
			"0001",
			"43528405058",
			"CACC",
			"10.00",
			"Teste de envio",
		)

		err := services.AddXMLPart(mw, message)

		if err != nil {
			fmt.Println("Error adding XML part:", err)
			continue
		}
	}

	mw.Close()

	var compressedMessage bytes.Buffer

	err := services.CompressContentToGzip(responseContent.Bytes(), &compressedMessage)
	if err != nil {
		fmt.Println("Error compressing message:", err)
		return
	}

	err = services.PostMessage(conn, string(compressedMessage.Bytes()), mw.Boundary())
	//err = services.PostMessage(conn, string(responseContent.Bytes()), mw.Boundary())
	//err = services.PostMessage(conn, message, "")
	if err != nil {
		fmt.Println("Error posting message:", err)
	}
}
