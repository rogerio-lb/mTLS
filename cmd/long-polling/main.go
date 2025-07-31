package main

import (
	"fmt"
	"mTLS/services"
	"time"
)

func main() {
	conn := services.CreateConnection()
	fmt.Println("=== TLS Connection Test Completed ===")
	fmt.Println("Sending HTTP/1.0 request...")

	var response *services.GetMessageResponse

	response, err := services.GetMessages(conn, "start")

	fmt.Println("Messages received:")
	fmt.Println(response.Message)

	if err != nil {
		services.FinishStream(conn, response.PIPullNext)
		panic("Error getting messages: %v\n" + err.Error())
		return
	}

	for {
		time.Sleep(2 * time.Second)
		connection := services.CreateConnection()
		response, err = services.GetMessages(connection, response.PIPullNext)
		if err != nil {
			services.FinishStream(connection, response.PIPullNext)
			fmt.Printf("Error getting messages: %v\n", err)
			break
		}

		fmt.Println("Messages received:")
		fmt.Println(response.Message)
	}

	//services.FinishStream(conn, "/api/v1/out/52833288/stream/11e5d52c-e302-4aec-af0a-83b3f068d26a")
}
