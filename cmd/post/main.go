package main

import (
	"bytes"
	"fmt"
	"mTLS/services"
	"mime/multipart"
)

func main() {
	conn := services.CreateConnection()

	if conn == nil {
		panic("Failed to create TLS connection")
	}

	var responseContent bytes.Buffer

	mw := multipart.NewWriter(&responseContent)

	for i := 0; i < 5; i++ {
		message := services.CreateMessage()
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
