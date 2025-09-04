package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"strings"
)

// XMLPart represents a single XML part from the multipart body
type XMLPart struct {
	Content    []byte
	ResourceID string
}

func ParseMultipartXMLWithBoundary(body io.Reader, boundary string) ([]XMLPart, error) {
	reader := multipart.NewReader(body, boundary)

	var xmlParts []XMLPart

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading multipart part: %v", err)
		}

		content, err := io.ReadAll(part)
		if err != nil {
			part.Close()
			return nil, fmt.Errorf("error reading part content: %v", err)
		}

		contentType := part.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/xml"
		}

		resourceId := part.Header.Get("PI-ResourceId")

		xmlPart := XMLPart{
			Content:    content,
			ResourceID: resourceId,
		}

		xmlParts = append(xmlParts, xmlPart)
		part.Close()
	}

	return xmlParts, nil
}

// Example usage
func main() {
	boundary := "599a8c95-608a-49b5-8c8e-7755790f0160"
	body := `--599a8c95-608a-49b5-8c8e-7755790f0160
PI-ResourceId: UEkBmRU0GT9Cq9AT+ZxLxKrdjEv5dL9j
Content-Type: application/xml;charset=utf-8

<?xml version="1.0" encoding="UTF-8" standalone="no"?><Envelope xmlns="https://www.bcb.gov.br/pi/pibr.002/1.3"><AppHdr><Fr><FIId><FinInstnId><Othr><Id>00038166</Id></Othr></FinInstnId></FIId></Fr><To><FIId><FinInstnId><Othr><Id>52833288</Id></Othr></FinInstnId></FIId></To><BizMsgIdr>M000381660MF5IV3GF21S17W2YJ2BRSP</BizMsgIdr><MsgDefIdr>pibr.002.spi.1.3</MsgDefIdr><CreDt>2025-09-04T14:49:19.935Z</CreDt><Sgntr><ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#"><ds:SignedInfo><ds:CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/><ds:SignatureMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"/><ds:Reference URI="#key-info-id"><ds:Transforms><ds:Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/></ds:Transforms><ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/><ds:DigestValue>YsJJMDNL6aFmCwNgnAjeTlQshWGxH+4IkHegixG7eAk=</ds:DigestValue></ds:Reference><ds:Reference URI=""><ds:Transforms><ds:Transform Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature"/><ds:Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/></ds:Transforms><ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/><ds:DigestValue>ABOx6122Pd00YbE92z6mY/w/FeBeEu5gqqdCHUFHevA=</ds:DigestValue></ds:Reference><ds:Reference><ds:Transforms><ds:Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/></ds:Transforms><ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/><ds:DigestValue>FUvoZyWnhuX/GP+cWYCgSMQOmvO2yoH4GjBIDAuDmYI=</ds:DigestValue></ds:Reference></ds:SignedInfo><ds:SignatureValue>xuwVNni1Y+W87zVRVOpgx9NVAPqHbAe+92xDZ4KSMZSUbKlT03Vuhj/2T5nOKPpdjMkfXwSqAb4Y&#13;
D6NXxss245jAfXTG0aiJOBolAiTFN94A/3XqASRGfpXZF98sDmP5C0bBx3gHmA22TMs51XFhSlfI&#13;
95JW4TD552hl3Bh7QMzXGOObrkB0qDDIFxqAwXuLP8tCYfKsqBsim9ASUA2guA2PEXM2/Fc4z8nk&#13;
tUiTCLwnkA/4b3bEkCIS0DO9b3mQpEBjYIdhfGagc3dbbOLjQM03bz2TDmmrKvCWwTYPPySGW+zZ&#13;
U0Z7rWSk+/aiDv/ZgsPcYhAbS9L1hYHjeo0VKg==</ds:SignatureValue><ds:KeyInfo Id="key-info-id"><ds:X509Data><ds:X509IssuerSerial><ds:X509IssuerName>CN=Autoridade Certificadora do SERPRO Final SSL, OU=Servico Federal de Processamento de Dados - SERPRO, OU=CSPB-1, O=ICP-Brasil, C=BR</ds:X509IssuerName><ds:X509SerialNumber>5529224565488886204649542255</ds:X509SerialNumber></ds:X509IssuerSerial></ds:X509Data></ds:KeyInfo></ds:Signature></Sgntr></AppHdr><Document><EchoRpt><GrpHdr><MsgId>M000381660MF5IV3GF21S17W2YJ2BRSP</MsgId><CreDtTm>2025-09-04T14:49:19.935Z</CreDtTm></GrpHdr><EchoTxInf><OrgnlData>Campo livre</OrgnlData></EchoTxInf></EchoRpt></Document></Envelope>
--599a8c95-608a-49b5-8c8e-7755790f0160
PI-ResourceId: UEkBmRU0GT+M9RarTFlOb7dEF9wJUHFa
Content-Type: application/xml;charset=utf-8

<?xml version="1.0" encoding="UTF-8" standalone="no"?><Envelope xmlns="https://www.bcb.gov.br/pi/pibr.002/1.3"><AppHdr><Fr><FIId><FinInstnId><Othr><Id>00038166</Id></Othr></FinInstnId></FIId></Fr><To><FIId><FinInstnId><Othr><Id>52833288</Id></Othr></FinInstnId></FIId></To><BizMsgIdr>M000381660MF5IV3GF21S17W2YR2XJPV</BizMsgIdr><MsgDefIdr>pibr.002.spi.1.3</MsgDefIdr><CreDt>2025-09-04T14:49:19.935Z</CreDt><Sgntr><ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#"><ds:SignedInfo><ds:CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/><ds:SignatureMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"/><ds:Reference URI="#key-info-id"><ds:Transforms><ds:Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/></ds:Transforms><ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/><ds:DigestValue>YsJJMDNL6aFmCwNgnAjeTlQshWGxH+4IkHegixG7eAk=</ds:DigestValue></ds:Reference><ds:Reference URI=""><ds:Transforms><ds:Transform Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature"/><ds:Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/></ds:Transforms><ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/><ds:DigestValue>TJHQXQVTRiCGqzf5Wyim16vWC15N6s5wVEDXLz2DKQw=</ds:DigestValue></ds:Reference><ds:Reference><ds:Transforms><ds:Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/></ds:Transforms><ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/><ds:DigestValue>W7AvyyHPz6lGD4gLRP2HHDZ5vFIjT++KRrXpc1OzhKQ=</ds:DigestValue></ds:Reference></ds:SignedInfo><ds:SignatureValue>0gTe7+K1gVQ8wnQry7Q0/l1WW+IhSYRv8KeYUjEQs8uTsE+DV1AhL5CjT9rFFaelXujRX91BxGnP&#13;
2lmz9Pbbp3QovzuTO+70zZnCz459jA0cdVkWsynNLSh14Vp4q8zQxyDAqxdZK9iYzjO0VpeNiUdE&#13;
LJpcefQZuAxFSJq2AnmKNn2Ekk0KF5AsNSB+WP5iToHJP2KLAmLddYBi8N/53McCkG9nixeC8HP6&#13;
qLbr+hORhg/4bkvOJRzDL1xWHw8qzaH7RnZr2Fie4ZMcST8cv/w6V++hZCHySTuploasRjLEQlj9&#13;
n9aWo4zCUYW+yhSWuQDhCCetKKvBsE1g8t8biw==</ds:SignatureValue><ds:KeyInfo Id="key-info-id"><ds:X509Data><ds:X509IssuerSerial><ds:X509IssuerName>CN=Autoridade Certificadora do SERPRO Final SSL, OU=Servico Federal de Processamento de Dados - SERPRO, OU=CSPB-1, O=ICP-Brasil, C=BR</ds:X509IssuerName><ds:X509SerialNumber>5529224565488886204649542255</ds:X509SerialNumber></ds:X509IssuerSerial></ds:X509Data></ds:KeyInfo></ds:Signature></Sgntr></AppHdr><Document><EchoRpt><GrpHdr><MsgId>M000381660MF5IV3GF21S17W2YR2XJPV</MsgId><CreDtTm>2025-09-04T14:49:19.935Z</CreDtTm></GrpHdr><EchoTxInf><OrgnlData>Campo livre</OrgnlData></EchoTxInf></EchoRpt></Document></Envelope>
--599a8c95-608a-49b5-8c8e-7755790f0160--`

	reader := strings.NewReader(body)

	parts, err := ParseMultipartXMLWithBoundary(reader, boundary)

	if err != nil {
		log.Fatalf("Error parsing multipart XML: %v", err)
	}

	for i, part := range parts {
		fmt.Println("Headers of part", i+1, ":", part.ResourceID)
		fmt.Println("Body of part", i+1, ":", string(part.Content))
	}
}
