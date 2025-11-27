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
		message := services.GeneratePacs004ForDict(
			"E49931906202511271738eT6eMqvatl9",
			"49931906",
			"0.10",
			services.USER_REQUEST_REASON,
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
