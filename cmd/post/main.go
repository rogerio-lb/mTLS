package main

import (
	"bytes"
	"fmt"
	"mTLS/services"
	"mime/multipart"
)

const debug = false

func main() {
	ISPB := "04902979"

	conn := services.CreateConnectionV2(debug, false)

	if conn == nil {
		panic("Failed to create TLS connection")
	}

	var responseContent bytes.Buffer

	mw := multipart.NewWriter(&responseContent)

	for i := 0; i < 1; i++ {
		//message := services.CreateMessage()

		message := services.GeneratePacs004ForDict(
			"E52833288202512191710eyaUu9DK5jg",
			"52833288",
			"10.00",
			services.USER_REQUEST_REASON,
			ISPB,
		)

		/*message := services.GeneratePacs008Manual(
			"04902979",
			"003816482",
			"0001",
			"43528405058",
			"CACC",
			"10.00",
			"Teste de envio",
		)*/

		//fmt.Println(message)

		/*message := services.GeneratePacs008Dict(
			"99999004",
			"003816482",
			"0001",
			"43528405058",
			"CACC",
			"10.00",
			"Teste de envio",
			"+5531982661780",
		)*/

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

	err = services.PostMessage(conn, string(compressedMessage.Bytes()), mw.Boundary(), ISPB)
	//err = services.PostMessage(conn, string(responseContent.Bytes()), mw.Boundary())
	//err = services.PostMessage(conn, message, "")
	if err != nil {
		fmt.Println("Error posting message:", err)
	}
}
