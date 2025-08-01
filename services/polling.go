package services

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type GetMessageResponse struct {
	Message    string
	PIPullNext string
	resourceID string
}

func StartStream(conn *tls.Conn) (*GetMessageResponse, error) {
	fmt.Println("Starting stream...")

	request := "GET /api/v1/out/52833288/stream/start HTTP/1.2\r\n" +
		"Host: icom-h.pi.rsfn.net.br\r\n" +
		//"Accept: multipart/mixed\r\n" +
		"User-Agent: Go-http-client/1.2\r\n" +
		"Connection: Close\r\n\r\n"

	_, err := conn.Write([]byte(request))
	if err != nil {
		fmt.Printf("Failed to write request: %v\n", err)
		return nil, err
	}

	reader := bufio.NewReader(conn)
	resp, err := http.ReadResponse(reader, nil)
	if err != nil {
		fmt.Printf("Failed to parse response: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read body: %v\n", err)
		return nil, err
	}

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Body received (%d bytes)\n", len(body))
	fmt.Printf("Body: %s\n", body)
	fmt.Printf("Headers: %v\n", resp.Header)

	fmt.Println("Stream started successfully.")

	return &GetMessageResponse{
		Message:    string(body),
		PIPullNext: resp.Header.Get("Pi-Pull-Next"),
		resourceID: resp.Header.Get("Pi-Resourceid"),
	}, nil
}

func GetMessages(conn *tls.Conn, step string) (*GetMessageResponse, error) {
	if step == "start" {
		return StartStream(conn)
	}

	return GetMessage(conn, step)
}

func GetMessage(conn *tls.Conn, pullnext string) (*GetMessageResponse, error) {
	fmt.Println("Starting stream...")

	request := "GET " + pullnext + " HTTP/1.2\r\n" +
		"Host: icom-h.pi.rsfn.net.br\r\n" +
		//"Accept: multipart/mixed\r\n" +
		"User-Agent: Go-http-client/1.2\r\n" +
		"Connection: Close\r\n\r\n"

	_, err := conn.Write([]byte(request))
	if err != nil {
		fmt.Printf("Failed to write request: %v\n", err)
		return nil, err
	}

	reader := bufio.NewReader(conn)
	resp, err := http.ReadResponse(reader, nil)
	if err != nil {
		fmt.Printf("Failed to parse response: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read body: %v\n", err)
		return nil, err
	}

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Body received (%d bytes)\n", len(body))
	fmt.Printf("Body: %s\n", body)
	fmt.Printf("Headers: %v\n", resp.Header)

	fmt.Println("Stream started successfully.")

	return &GetMessageResponse{
		Message:    string(body),
		PIPullNext: resp.Header.Get("PI-Pull-Next"),
		resourceID: resp.Header.Get("Pi-Resourceid"),
	}, nil
}

func FinishStream(conn *tls.Conn, pullNext string) error {
	fmt.Println("Finishing stream...")

	request := "DELETE " + pullNext + " HTTP/1.2\r\n" +
		"Host: icom-h.pi.rsfn.net.br\r\n" +
		"User-Agent: Go-http-client/1.2\r\n" +
		"Connection: Close\r\n\r\n"

	_, err := conn.Write([]byte(request))
	if err != nil {
		fmt.Printf("Failed to write request: %v\n", err)
		return err
	}

	reader := bufio.NewReader(conn)
	resp, err := http.ReadResponse(reader, nil)
	if err != nil {
		fmt.Printf("Failed to parse response: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read body: %v\n", err)
		return err
	}

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Body received (%d bytes)\n", len(body))
	fmt.Printf("Body: %s\n", body)
	fmt.Printf("Headers: %v\n", resp.Header)

	fmt.Println("Stream finished successfully.")

	return nil
}

func PostMessage(conn *tls.Conn, content string) error {
	fmt.Println("Sending message...")

	contentLength := len(content)

	request := "POST /api/v1/in/52833288/msgs HTTP/1.2\r\n" +
		"Host: icom-h.pi.rsfn.net.br\r\n" +
		"Content-Type: application/xml; charset=utf-8\r\n" +
		"Content-Length: " + strconv.Itoa(contentLength) + "\r\n" +
		"User-Agent: Go-http-client/1.1\r\n" +
		"Connection: close\r\n\r\n" +
		content

	_, err := conn.Write([]byte(request))
	if err != nil {
		fmt.Printf("Failed to write request: %v\n", err)
		return err
	}

	_, err = conn.Write([]byte(content))

	if err != nil {
		fmt.Printf("Failed to write content: %v\n", err)
		return err
	}

	reader := bufio.NewReader(conn)
	resp, err := http.ReadResponse(reader, nil)
	if err != nil {
		fmt.Printf("Failed to parse response: %v\n", err)
		return err
	}

	defer resp.Body.Close()

	fmt.Printf("Status: %s\n", resp.Status)
	return nil
}
