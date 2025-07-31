package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
)

type GetMessageResponse struct {
	Messages   string
	PIPullNext string
	resourceID string
}

func main() {
	conn := CreateConnection()
	fmt.Println("=== TLS Connection Test Completed ===")

	fmt.Println("Sending HTTP/1.0 request...")

	response, err := GetMessages(conn, "start")

	if err != nil {
		panic(err)
	}
}

func ConfigureMTLS() *tls.Config {
	caCert, _ := os.ReadFile("/home/roger/projects/lb/mTLS/certs/Cadeia_Oficial.p7b")
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, _ := tls.LoadX509KeyPair("/home/roger/projects/lb/mTLS/certs/certificado-25065760.cer", "/home/roger/projects/lb/mTLS/certs/spb_hm_private_unencrypted.key")

	// Configure TLS
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		MinVersion:         tls.VersionTLS12,                                    // --tlsv1.2
		MaxVersion:         tls.VersionTLS12,                                    // Force TLS 1.2
		CipherSuites:       []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256}, // --ciphers
		InsecureSkipVerify: true,                                                // -k flag (skip verification)
	}

	return tlsConfig
}

func CreateConnection() *tls.Conn {
	fmt.Println("=== Testing TLS Connection ===")

	tlsConfig := ConfigureMTLS()

	conn, err := tls.Dial("tcp", "127.0.0.1:16522", tlsConfig)
	if err != nil {
		fmt.Printf("TLS connection failed: %v\n", err)
		return nil
	}

	fmt.Println("TLS connection successful!")
	fmt.Printf("TLS Version: %x\n", conn.ConnectionState().Version)
	fmt.Printf("Cipher Suite: %x\n", conn.ConnectionState().CipherSuite)
	fmt.Printf("Server Certificates: %d\n", len(conn.ConnectionState().PeerCertificates))

	return conn
}

func StartStream(conn *tls.Conn) (*GetMessageResponse, error) {
	fmt.Println("Starting stream...")

	request := "GET /api/v1/out/52833288/stream/start HTTP/1.2\r\n" +
		"Host: icom-h.pi.rsfn.net.br\r\n" +
		//"Accept: multipart/mixed\r\n" +
		"User-Agent: Go-http-client/1.2\r\n" +
		"Connection: Keep-Alive\r\n\r\n"

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
		Messages:   string(body),
		PIPullNext: resp.Header.Get("PIPull-Next"),
		resourceID: resp.Header.Get("Resource-ID"),
	}, nil
}

func GetMessages(conn *tls.Conn, step string) (*GetMessageResponse, error) {
	if step == "start" {
		return StartStream(conn)
	}

	fmt.Println("Invalid step. Only 'start' is supported for now.")
	return nil, nil
}
