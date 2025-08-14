package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"mTLS/services"
	"os/exec"
)

const (
	msgIdPrefix    = "M"
	msgIdISPBLen   = 8
	msgIdSuffixLen = 23
	// All alphanumeric characters as per SPI specification [a-z|A-Z|0-9]
	msgIdAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var randReader io.Reader = rand.Reader

var pibr001 = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<Envelope xmlns="https://www.bcb.gov.br/pi/pacs.008/1.13">
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
        <BizMsgIdr>M52833288202508141640Mq7fgRaca2N</BizMsgIdr>
        <MsgDefIdr>pacs.008.spi.1.13</MsgDefIdr>
        <CreDt>2020-01-01T08:30:12.000Z</CreDt>
        <Sgntr/>
    </AppHdr>
    <Document>
        <FIToFICstmrCdtTrf>
            <GrpHdr>
                <MsgId>M5283328820250814121142</MsgId>
                <CreDtTm>2020-01-01T08:30:12.000Z</CreDtTm>
                <NbOfTxs>1</NbOfTxs>
                <SttlmInf>
                    <SttlmMtd>CLRG</SttlmMtd>
                </SttlmInf>
                <PmtTpInf>
                    <InstrPrty>HIGH</InstrPrty>
                    <SvcLvl>
                        <Prtry>PAGPRI</Prtry>
                    </SvcLvl>
                </PmtTpInf>
            </GrpHdr>
            <CdtTrfTxInf>
                <PmtId>
                    <EndToEndId>E9999901012341234123412345678900</EndToEndId>
                </PmtId>
                <IntrBkSttlmAmt Ccy="BRL">1000.00</IntrBkSttlmAmt>
                <AccptncDtTm>2020-01-01T08:30:00.000Z</AccptncDtTm>
                <ChrgBr>SLEV</ChrgBr>
                <MndtRltdInf>
                    <Tp>
                        <LclInstrm>
                            <Prtry>MANU</Prtry>
                        </LclInstrm>
                    </Tp>
                </MndtRltdInf>
                <Dbtr>
                    <Nm>Fulano da Silva</Nm>
                    <Id>
                        <PrvtId>
                            <Othr>
                                <Id>70000000000</Id>
                            </Othr>
                        </PrvtId>
                    </Id>
                </Dbtr>
                <DbtrAcct>
                    <Id>
                        <Othr>
                            <Id>500000</Id>
                            <Issr>3000</Issr>
                        </Othr>
                    </Id>
                    <Tp>
                        <Cd>CACC</Cd>
                    </Tp>
                </DbtrAcct>
                <DbtrAgt>
                    <FinInstnId>
                        <ClrSysMmbId>
                            <MmbId>10000000</MmbId>
                        </ClrSysMmbId>
                    </FinInstnId>
                </DbtrAgt>
                <CdtrAgt>
                    <FinInstnId>
                        <ClrSysMmbId>
                            <MmbId>20000000</MmbId>
                        </ClrSysMmbId>
                    </FinInstnId>
                </CdtrAgt>
                <Cdtr>
                    <Id>
                        <PrvtId>
                            <Othr>
                                <Id>80000000000</Id>
                            </Othr>
                        </PrvtId>
                    </Id>
                </Cdtr>
                <CdtrAcct>
                    <Id>
                        <Othr>
                            <Id>600000</Id>
                            <Issr>4000</Issr>
                        </Othr>
                    </Id>
                    <Tp>
                        <Cd>SVGS</Cd>
                    </Tp>
                </CdtrAcct>
                <Purp>
                    <Cd>IPAY</Cd>
                </Purp>
                <RmtInf>
                    <Ustrd>Campo livre [0]</Ustrd>
                </RmtInf>
            </CdtTrfTxInf>
        </FIToFICstmrCdtTrf>
    </Document>
</Envelope>`

func callJavaFunction(message string) (string, error) {
	// Run Java program with message as argument
	cmd := exec.Command("java", "-jar", "java/signer.jar", "-a", message)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run Java: %v", err)
	}

	fmt.Print(string(output))

	return string(output), nil
}

func main() {
	conn := services.CreateConnection()

	if conn == nil {
		panic("Failed to create TLS connection")
	}

	str, err := callJavaFunction(pibr001)

	if err != nil {
		fmt.Println("Error calling Java function:", err)
		return
	}

	err = services.PostMessage(conn, str)

	if err != nil {
		fmt.Println("Error posting message:", err)
	}

	/*id, _ := GenerateMsgId("52833288")

	fmt.Println("Generated Message ID:", id)*/

	//message, _ := document.New()
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
