package services

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

func ConfigureMTLS() *tls.Config {
	caCert, _ := os.ReadFile("D:\\LbCerts\\Cadeia_Oficial.p7b")
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, _ := tls.LoadX509KeyPair("D:\\LbCerts\\certificado-25065760.cer", "D:\\LbCerts\\spb_hm_private_unencrypted.key")

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
