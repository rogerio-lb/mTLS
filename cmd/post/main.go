package main

import (
	"fmt"
	"mTLS/services"
)

var pibr001 = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<Envelope xmlns="https://www.bcb.gov.br/pi/pibr.001/1.3">
    <AppHdr>
        <Fr>
            <FIId>
                <FinInstnId>
                    <Othr>
                        <Id>52833288</Id>
                    </Othr>
                </FinInstnId>
            </FIId>
        </Fr>
        <To>
            <FIId>
                <FinInstnId>
                    <Othr>
                        <Id>00038166</Id>
                    </Othr>
                </FinInstnId>
            </FIId>
        </To>
        <BizMsgIdr>M528332880884965a4f82c427ae16bcf</BizMsgIdr>
        <MsgDefIdr>pibr.001.spi.1.3</MsgDefIdr>
        <CreDt>2025-07-03T18:04:00.000Z</CreDt>
		<Sgntr>
            <ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
                <ds:SignedInfo>
                    <ds:CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/>
                    <ds:SignatureMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"/>
                    <ds:Reference URI="#key-info-id">
                        <ds:Transforms>
                            <ds:Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/>
                        </ds:Transforms>
                        <ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/>
                        <ds:DigestValue>YsJJMDNL6aFmCwNgnAjeTlQshWGxH+4IkHegixG7eAk=</ds:DigestValue>
                    </ds:Reference>
                    <ds:Reference URI="">
                        <ds:Transforms>
                            <ds:Transform Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature"/>
                            <ds:Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/>
                        </ds:Transforms>
                        <ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/>
                        <ds:DigestValue>CUV14SRYW6vbJnIutRRt9zg9haO6c0O+hTdBQakSfFA=</ds:DigestValue>
                    </ds:Reference>
                    <ds:Reference>
                        <ds:Transforms>
                            <ds:Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/>
                        </ds:Transforms>
                        <ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/>
                        <ds:DigestValue>KIYsuJEgu6CFWEtI7dVxIxdZZhAx6aqdblI3rIoZ8xc=</ds:DigestValue>
                    </ds:Reference>
                </ds:SignedInfo>
                <ds:SignatureValue>cwuC3G6afQlkjdquk6ot5tsZQ21xxX+wK/rbbF6hfhBJqnGimWkObuNwZS/GTKa71qZbY/3Ndyk9&#13;
MipqlbSQBKL5PHatgrowqn/Ow7Vtu8NpEhXMeD0nI3f3ceqyb89PMCwo+fug59Dst+zQ47PEtuIt&#13;
X5WVnNlluZSpRzIJfeBQOp/hRmp1NySr35rTmsZ067klTSoli4B2ZKnWd5WTCeijzvvUIIKwh0kl&#13;
tkPRgaOYpy9QRWZ45d42+kcGvndWH+XRiBFtVWLyKwR2OLwBJQdur3ff0OSpLmYS1PyN25yH8cD0&#13;
syeVgETjwdiGUNeZrsRHqFqYauwZW0HH+o9nNQ==</ds:SignatureValue>
                <ds:KeyInfo Id="key-info-id">
                    <ds:X509Data>
                        <ds:X509IssuerSerial>
                            <ds:X509IssuerName>CN=Autoridade Certificadora do SERPRO Final SSL, OU=Servico Federal de Processamento de Dados - SERPRO, OU=CSPB-1, O=ICP-Brasil, C=BR</ds:X509IssuerName>
                            <ds:X509SerialNumber>5529224565488886204649542255</ds:X509SerialNumber>
                        </ds:X509IssuerSerial>
                    </ds:X509Data>
                </ds:KeyInfo>
            </ds:Signature>
        </Sgntr>
    </AppHdr>
    <Document>
        <EchoReq>
            <GrpHdr>
                <MsgId>M528332880884965a4f82c427ae16bcf</MsgId>
                <CreDtTm>2025-07-01T20:15:00.000Z</CreDtTm>
            </GrpHdr>
            <EchoTxInf>
                <Data>Teste 1</Data>
            </EchoTxInf>
        </EchoReq>
    </Document>
</Envelope>`

func main() {
	conn := services.CreateConnection()

	if conn == nil {
		panic("Failed to create TLS connection")
	}

	err := services.PostMessage(conn, pibr001)

	if err != nil {
		fmt.Println("Error posting message:", err)
	}
}
