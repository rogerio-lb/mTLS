package services

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

type GetMessageResponse struct {
	StatusCode int
	Headers    http.Header
	Message    string
	PIPullNext string
	resourceID string
}

func GetMessages(conn *http.Client, step string, ispb string) (*GetMessageResponse, error) {
	if step == "start" {
		return GetMessage(conn, "", ispb)
	}

	return GetMessage(conn, step, ispb)
}

func GetMessage(conn *http.Client, pullnext string, ispb string) (*GetMessageResponse, error) {
	if pullnext == "" {
		fmt.Println("Starting stream...")
	}

	if pullnext == "" {
		//pullnext = "/api/v1/out/52833288/stream/start"
		pullnext = "start"
	}

	//url := fmt.Sprintf("%s%s", "http://localhost:3000", pullnext)
	url := fmt.Sprintf("%s/api/v1/out/%s/stream/%s", "http://localhost:3000", ispb, pullnext)

	httpRequest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return nil, err
	}

	httpRequest.Host = "icom-h.pi.rsfn.net.br"
	httpRequest.Header.Set("Accept", "application/xml")
	httpRequest.Header.Set("Accept-Encoding", "gzip")
	httpRequest.Header.Set("User-Agent", "Go-http-client/1.2")
	httpRequest.Header.Set("Connection", "Close")

	resp, err := conn.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("Failed to send request: %v\n", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read body: %v\n", err)
	}

	/*
		request := "GET " + pullnext + " HTTP/1.2\r\n" +
			"Host: icom-h.pi.rsfn.net.br\r\n" +
			//"Accept: multipart/mixed\r\n" +
			"Accept: application/xml\r\n" +
			"Accept-Encoding: gzip\r\n" +
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
		}*/

	var decompressedMessage []byte

	if resp.StatusCode == http.StatusOK {
		decompressedMessage, err = decompressContentFromGzip(body)
		if err != nil {
			fmt.Printf("Failed to decompress body: %v\n", err)
			return nil, err
		}

		contentType := resp.Header.Get("Content-Type")
		boundary := contentType[strings.Index(contentType, ";")+1:]
		boundary = strings.TrimSpace(boundary)
		boundary = strings.TrimPrefix(boundary, "boundary=")

		/*err = parseMultipartFromString(string(decompressedMessage), boundary)
		if err != nil {
			return nil, err
		}*/
	}

	/*if resp.StatusCode != http.StatusNoContent {
		fmt.Printf("Status: %s\n", resp.Status)
		fmt.Printf("Body received (%d bytes)\n", len(body))
		fmt.Printf("Body: %s\n", decompressedMessage)
		fmt.Printf("Headers: %v\n", resp.Header)
	}*/

	if pullnext == "start" {
		fmt.Println("Stream started successfully.")
	}

	rawPiPullNext := resp.Header.Get("PI-Pull-Next")
	piPullNextContent := strings.Split(rawPiPullNext, "/")
	piPullNextId := piPullNextContent[len(piPullNextContent)-1]

	return &GetMessageResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Message:    string(decompressedMessage),
		PIPullNext: piPullNextId,
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

	/*body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read body: %v\n", err)
		return err
	}*/

	/*if resp.StatusCode != http.StatusNoContent {
		fmt.Printf("Status: %s\n", resp.Status)
		fmt.Printf("Status Code: %d\n", resp.StatusCode)
		fmt.Printf("Body received (%d bytes)\n", len(body))
		fmt.Printf("Body: %s\n", body)
		fmt.Printf("Headers: %v\n", resp.Header)

		fmt.Println("Stream finished successfully.")
	}*/

	return nil
}

func PostMessage(conn *http.Client, content string, boundary string) error {
	fmt.Println("Sending message...")

	contentLength := len(content)

	/*var request string

	if boundary != "" {
		request = "POST /api/v1/in/52833288/msgs HTTP/1.1\r\n" +
			"Host: icom-h.pi.rsfn.net.br\r\n" +
			"Content-Type: multipart/mixed; boundary=" + boundary + "\r\n" +
			"Content-Encoding: gzip\r\n" +
			"Content-Length: " + strconv.Itoa(contentLength) + "\r\n" +
			"User-Agent: Go-http-client/1.1\r\n" +
			"\r\n" +
			content
	} else {
		request = "POST /api/v1/in/52833288/msgs HTTP/1.2\r\n" +
			"Host: icom-h.pi.rsfn.net.br\r\n" +
			"Content-Type: application/xml; charset=utf-8\r\n" +
			"Content-Encoding: gzip\r\n" +
			"Content-Length: " + strconv.Itoa(contentLength) + "\r\n" +
			"User-Agent: Go-http-client/1.1\r\n" +
			"Connection: close\r\n" +
			"\r\n" +
			content
	}

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
	}*/

	bodyReader := bytes.NewReader([]byte(content))

	request, err := http.NewRequest("POST", "http://localhost:3000/api/v1/in/52833288/msgs", bodyReader)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return err
	}

	request.Host = "icom-h.pi.rsfn.net.br"
	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Content-Length", strconv.Itoa(contentLength))
	request.Header.Set("User-Agent", "Go-http-client/1.1")
	if boundary != "" {
		request.Header.Set("Content-Type", "multipart/mixed; boundary="+boundary)
	} else {
		request.Header.Set("Content-Type", "application/xml; charset=utf-8")
	}

	resp, err := conn.Do(request)
	if err != nil {
		fmt.Printf("Failed to send request: %v\n", err)
		return err
	}

	defer resp.Body.Close()

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Headers: %v\n", resp.Header)

	return nil
}

func decompressContentFromGzip(body []byte) ([]byte, error) {
	gzipReader, err := gzip.NewReader(bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	decompressedBody, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, err
	}

	return decompressedBody, nil
}

func parseMultipartFromString(responseBody string, boundary string) error {
	reader := multipart.NewReader(strings.NewReader(responseBody), boundary)

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read multipart: %v", err)
		}

		/*contentType := part.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/xml") {
			xmlData, err := io.ReadAll(part)
			if err != nil {
				return fmt.Errorf("failed to read XML part: %v", err)
			}

			fmt.Printf("XML Part: %s\n", xmlData)
		}*/

		part.Close()
	}

	return nil
}
