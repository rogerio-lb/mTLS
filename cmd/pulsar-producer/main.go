package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	ramdom "math/rand"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

const (
	msgIdPrefix    = "M"
	msgIdISPBLen   = 8
	msgIdSuffixLen = 23
	// All alphanumeric characters as per SPI specification [a-z|A-Z|0-9]
	msgIdAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	returnIdPrefix    = "D"
	returnIdISPBLen   = 8
	returnIdSuffixLen = 11
	// All alphanumeric characters as per SPI specification [a-z|A-Z|0-9]
	returnIdAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	pacsRandReader      io.Reader = rand.Reader
	numericPattern                = regexp.MustCompile(`^[0-9]{8}$`)
	alphanumericPattern           = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	getCurrentTime                = time.Now
	randReader          io.Reader = rand.Reader
)

var message = `<Envelope xmlns="https://www.bcb.gov.br/pi/pacs.002/1.14">
    <AppHdr>
        <Fr>
            <FIId>
                <FinInstnId>
                    <Othr>
                        <Id>00038166</Id>
                    </Othr>
                </FinInstnId>
            </FIId>
        </Fr>
        <To>
            <FIId>
                <FinInstnId>
                    <Othr>
                        <Id>52833288</Id>
                    </Othr>
                </FinInstnId>
            </FIId>
        </To>
        <BizMsgIdr>M5283328820250829144521906IxC8GD</BizMsgIdr>
        <MsgDefIdr>pacs.002.spi.1.14</MsgDefIdr>
        <CreDt>2025-08-29T14:45:21.906Z</CreDt>
        <Sgntr/>
    </AppHdr>
    <Document>
        <FIToFIPmtStsRpt>
            <GrpHdr>
                <MsgId>M5283328820250829144521906IxC8GD</MsgId>
                <CreDtTm>2025-08-29T14:45:21.906Z</CreDtTm>
            </GrpHdr>
            <TxInfAndSts>
                <OrgnlInstrId>E71027866202508272059568098cbddQ</OrgnlInstrId>
                <OrgnlEndToEndId>E71027866202508272059568098cbddQ</OrgnlEndToEndId>
                <TxSts>ACSP</TxSts>
                <FctvIntrBkSttlmDt>
                    <DtTm>2025-08-29T14:45:21.906Z</DtTm>
                </FctvIntrBkSttlmDt>
                <OrgnlTxRef>
                    <IntrBkSttlmDt>2025-08-29</IntrBkSttlmDt>
                </OrgnlTxRef>
            </TxInfAndSts>
        </FIToFIPmtStsRpt>
    </Document>
</Envelope>`

var pacs008 = `<?xml version="1.0" encoding="UTF-8" standalone="no"?><Envelope xmlns="https://www.bcb.gov.br/pi/pacs.008/1.13"><AppHdr><Fr><FIId><FinInstnId><Othr><Id>00038166</Id></Othr></FinInstnId></FIId></Fr><To><FIId><FinInstnId><Othr><Id>52833288</Id></Othr></FinInstnId></FIId></To><BizMsgIdr>M000381660MEUGLGGA21S3D8M1PYVR8J</BizMsgIdr><MsgDefIdr>pacs.008.spi.1.13</MsgDefIdr><CreDt>2025-08-27T21:00:23.050Z</CreDt><Sgntr><ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#"><ds:SignedInfo><ds:CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/><ds:SignatureMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"/><ds:Reference URI="#key-info-id"><ds:Transforms><ds:Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/></ds:Transforms><ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/><ds:DigestValue>YsJJMDNL6aFmCwNgnAjeTlQshWGxH+4IkHegixG7eAk=</ds:DigestValue></ds:Reference><ds:Reference URI=""><ds:Transforms><ds:Transform Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature"/><ds:Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/></ds:Transforms><ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/><ds:DigestValue>j925kvK+5KCg2PefQ+P4s9EVQD2x5fG02gymmjfJXLk=</ds:DigestValue></ds:Reference><ds:Reference><ds:Transforms><ds:Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/></ds:Transforms><ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/><ds:DigestValue>TZcTo5jy4t4xdRh5YpzhjrDSeZCHqjm1Hwzs5hLcAfo=</ds:DigestValue></ds:Reference></ds:SignedInfo><ds:SignatureValue>lkQ3tp4ARpQXGdVN0fb6RBIQl/Gp/abXkcBLSW1KzUMImwhmT8gE3dUFg9iN08q/5kXutZ9Fw1qI&#13;
1oYJFw8KkO69olbCsw/iiSG0cYpWtPOYrpStuYdreMjmlYzJJG0UPbVpza7pBHnRd0fzCbtRM/cC&#13;
VG0B/GFph6TvR922TjQ1iN+knlrgWhg+1+6FN3pwRcviWl+wLbhj0cZGQxaFsgeiM1PkfVlUCEPx&#13;
+Y+KtaZMhCWdhMJRmxMilJJuvgGd7/DDm1SDCmVgZ/cOZVG6RtrnZP8KQHUeS+/BVS6FaH1yyDhC&#13;
BndHr1BDqJyTRPWrU1XsZ9vChUgzgkeaG/RBOA==</ds:SignatureValue><ds:KeyInfo Id="key-info-id"><ds:X509Data><ds:X509IssuerSerial><ds:X509IssuerName>CN=Autoridade Certificadora do SERPRO Final SSL, OU=Servico Federal de Processamento de Dados - SERPRO, OU=CSPB-1, O=ICP-Brasil, C=BR</ds:X509IssuerName><ds:X509SerialNumber>5529224565488886204649542255</ds:X509SerialNumber></ds:X509IssuerSerial></ds:X509Data></ds:KeyInfo></ds:Signature></Sgntr></AppHdr><Document><FIToFICstmrCdtTrf><GrpHdr><MsgId>M000381660MEUGLGGA21S3D8M1PYVR8J</MsgId><CreDtTm>2025-08-27T21:00:23.050Z</CreDtTm><NbOfTxs>1</NbOfTxs><SttlmInf><SttlmMtd>CLRG</SttlmMtd></SttlmInf><PmtTpInf><InstrPrty>HIGH</InstrPrty><SvcLvl><Prtry>PAGPRI</Prtry></SvcLvl></PmtTpInf></GrpHdr><CdtTrfTxInf><PmtId><EndToEndId>E71027866202508272059568098cbddP</EndToEndId></PmtId><IntrBkSttlmAmt Ccy="BRL">6.00</IntrBkSttlmAmt><AccptncDtTm>2025-08-27T21:00:20.700Z</AccptncDtTm><ChrgBr>SLEV</ChrgBr><MndtRltdInf><Tp><LclInstrm><Prtry>MANU</Prtry></LclInstrm></Tp></MndtRltdInf><Dbtr><Nm>TeofiloBS2Pay</Nm><Id><PrvtId><Othr><Id>84825873000123</Id></Othr></PrvtId></Id></Dbtr><DbtrAcct><Id><Othr><Id>90530802</Id><Issr>0001</Issr></Othr></Id><Tp><Cd>CACC</Cd></Tp></DbtrAcct><DbtrAgt><FinInstnId><ClrSysMmbId><MmbId>71027866</MmbId></ClrSysMmbId></FinInstnId></DbtrAgt><CdtrAgt><FinInstnId><ClrSysMmbId><MmbId>52833288</MmbId></ClrSysMmbId></FinInstnId></CdtrAgt><Cdtr><Id><PrvtId><Othr><Id>11909227000170</Id></Othr></PrvtId></Id></Cdtr><CdtrAcct><Id><Othr><Id>1169414</Id><Issr>0001</Issr></Othr></Id><Tp><Cd>TRAN</Cd></Tp></CdtrAcct><Purp><Cd>IPAY</Cd></Purp></CdtTrfTxInf></FIToFICstmrCdtTrf></Document></Envelope>`

var pacs002 = `<Envelope xmlns="https://www.bcb.gov.br/pi/pacs.002/1.14">
    <AppHdr>
        <Fr>
            <FIId>
                <FinInstnId>
                    <Othr>
                        <Id>00038166</Id>
                    </Othr>
                </FinInstnId>
            </FIId>
        </Fr>
        <To>
            <FIId>
                <FinInstnId>
                    <Othr>
                        <Id>52833288</Id>
                    </Othr>
                </FinInstnId>
            </FIId>
        </To>
        <BizMsgIdr>M5283328820250829144521906IxC8GD</BizMsgIdr>
        <MsgDefIdr>pacs.002.spi.1.14</MsgDefIdr>
        <CreDt>2025-08-29T14:45:21.906Z</CreDt>
        <Sgntr/>
    </AppHdr>
    <Document>
        <FIToFIPmtStsRpt>
            <GrpHdr>
                <MsgId>M5283328820250829144521906IxC8GD</MsgId>
                <CreDtTm>2025-08-29T14:45:21.906Z</CreDtTm>
            </GrpHdr>
            <TxInfAndSts>
                <OrgnlInstrId>E71027866202508272059568098cbddP</OrgnlInstrId>
                <OrgnlEndToEndId>E71027866202508272059568098cbddP</OrgnlEndToEndId>
                <TxSts>ACSP</TxSts>
                <FctvIntrBkSttlmDt>
                    <DtTm>2025-08-29T14:45:21.906Z</DtTm>
                </FctvIntrBkSttlmDt>
                <OrgnlTxRef>
                    <IntrBkSttlmDt>2025-08-29</IntrBkSttlmDt>
                </OrgnlTxRef>
            </TxInfAndSts>
        </FIToFIPmtStsRpt>
    </Document>
</Envelope>`

func main() {
	// Create Pulsar client with SSL
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:            "pulsar+ssl://pc-a5bec094.aws-use2-production-snci-pool-kid.streamnative.aws.snio.cloud",
		Authentication: pulsar.NewAuthenticationToken("eyJhbGciOiJSUzI1NiIsImtpZCI6IjE0NjNhODQ5LTNkNzUtNTlmMi1hMTgyLTVjNzE0ODY4YjBhMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsidXJuOnNuOnB1bHNhcjpvLWU1NWNwOmxiLXBheSJdLCJodHRwczovL3N0cmVhbW5hdGl2ZS5pby9zY29wZSI6WyJhZG1pbiIsImFjY2VzcyJdLCJodHRwczovL3N0cmVhbW5hdGl2ZS5pby91c2VybmFtZSI6ImxiLXN0Z0BvLWU1NWNwLmF1dGguc3RyZWFtbmF0aXZlLmNsb3VkIiwiaWF0IjoxNzUzMTMxMTcxLCJpc3MiOiJodHRwczovL3BjLWE1YmVjMDk0LmF3cy11c2UyLXByb2R1Y3Rpb24tc25jaS1wb29sLWtpZC5zdHJlYW1uYXRpdmUuYXdzLnNuaW8uY2xvdWQvYXBpa2V5cy8iLCJqdGkiOiJhYTk3YTVhN2YwNzA0N2FjYTI3MzQ4ODdlOTI1ZDMyMyIsInBlcm1pc3Npb25zIjpbXSwic3ViIjoiVWNDbGVyaENVRFh6S1NUZ21WbHFPVkF5b1R0aDlUT1lAY2xpZW50cyJ9.UfLKDZNusPHya-xgdWHoSNXbp6nhBMaEyizzULkWCsriY4VKdfkJ6OqnrPXP9xOi0aVKzCL-9ObgzxBklKoFObguZJ1MrIgzeiQfp0FUfmylwWz_jb-zbPZ5cbclvrbMXojKJte1lk9GxmmggBf-zUpuRGDiVGV42ZnU2AVJ-1PXx5frQ5SUbfJkfIRDp566b6PoF9r80gYc594CCo_Z0nUUMjHR_1molD5BDBYoK3O71yy-kEf-_J_nOdfMBQdHZbQGitpo5BzLvE-kdpHg0JZ392IZPeWhoZCEyGfLaNp6aNi8tyCMe--NQFm78h6bQ4L3VHge7BVR7dOjMRiBQQ"),
	})
	if err != nil {
		log.Fatalf("Could not instantiate Pulsar client: %v", err)
	}
	defer client.Close()

	// Create producer
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic:                   "persistent://lb-core/spi/bridge-to-psti-topic-partition-0",
		Name:                    "my-producer",
		SendTimeout:             10 * time.Second,
		DisableBatching:         false,
		BatchingMaxPublishDelay: 100 * time.Millisecond,
		BatchingMaxMessages:     1000,
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	/*now := getCurrentTime().UTC()

	endToEndID, _ := GenerateEndToEndId("52833288", now)

	id, _ := GenerateMsgId("52833288")

	ready := fmt.Sprintf(pacs008, id, now.Format("2006-01-02T15:04:05.000Z"), id, now.Format("2006-01-02T15:04:05.000Z"), endToEndID, now.Format("2006-01-02T15:04:05.000Z")) //pacs008

	m, _ := callJavaFunction(ready)*/

	// Send messages
	ctx := context.Background()

	for i := 0; i < 1; i++ {
		message := &pulsar.ProducerMessage{
			Payload: []byte(pacs002),
			Key:     "message-key-" + string(rune(i+'0')),
			Properties: map[string]string{
				"timestamp": time.Now().Format(time.RFC3339),
				"sender":    "go-producer",
			},
		}

		messageID, err := producer.Send(ctx, message)
		if err != nil {
			log.Printf("Failed to send message %d: %v", i, err)
			continue
		}

		log.Printf("Successfully sent message %d with ID: %s", i, messageID)
		time.Sleep(1 * time.Second)
	}

	log.Println("All messages sent successfully!")
}

func callJavaFunction(message string) (string, error) {
	// Run Java program with message as argument
	cmd := exec.Command("java", "-jar", "/home/roger/projects/lb/signer-java/target/signer-java-1.0-SNAPSHOT.jar", "-a", message)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run Java: %v", err)
	}

	return string(output), nil
}

func GenerateMsgId(ispb string) (string, error) {
	// Generate random bytes
	randomBytes := make([]byte, msgIdSuffixLen)
	if _, err := randReader.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Convert random bytes to alphanumeric characters
	suffix := make([]byte, msgIdSuffixLen)
	for i := 0; i < msgIdSuffixLen; i++ {
		suffix[i] = msgIdAlphabet[randomBytes[i]%byte(len(msgIdAlphabet))]
	}

	return msgIdPrefix + ispb + string(suffix), nil
}

func GenerateEndToEndId(ident string, t time.Time) (string, error) {
	// Validate identifier format
	if !numericPattern.MatchString(ident) {
		return "", fmt.Errorf("identifier must be exactly 8 numeric digits")
	}

	// Validate timestamp is within 12 hours of current time
	now := getCurrentTime().UTC()
	diff := t.UTC().Sub(now)
	if diff < -12*time.Hour || diff > 12*time.Hour {
		return "", fmt.Errorf("timestamp must be within 12 hours of current time")
	}

	timestamp := t.UTC().Format("200601021504")

	randomBytes := make([]byte, 8)
	if _, err := pacsRandReader.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	suffix := base64.RawURLEncoding.EncodeToString(randomBytes)
	suffix = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, suffix)

	if len(suffix) > 11 {
		suffix = suffix[:11]
	}
	for len(suffix) < 11 {
		suffix += "0" // padding if needed
	}

	return fmt.Sprintf("E%s%s%s", ident, timestamp, suffix), nil
}

func GenerateRandomAlphanumeric(length int) string {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = returnIdAlphabet[ramdom.Intn(len(returnIdAlphabet))]
	}
	return string(result)
}
