package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"mTLS/services"
	ramdom "math/rand"
	"os/exec"
	"regexp"
	"strings"
	"time"
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
	// pacsRandReader is used to make testing easier by allowing us to replace the random source
	pacsRandReader io.Reader = rand.Reader

	// numericPattern validates that a string contains only digits
	numericPattern = regexp.MustCompile(`^[0-9]{8}$`)

	// alphanumericPattern validates that a string contains only alphanumeric characters
	alphanumericPattern = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

	// getCurrentTime is used to make testing easier by allowing us to replace the current time
	getCurrentTime = time.Now
)

var randReader io.Reader = rand.Reader

var pacs008 = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
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
        <BizMsgIdr>%s</BizMsgIdr>
        <MsgDefIdr>pacs.008.spi.1.13</MsgDefIdr>
        <CreDt>%s</CreDt>
        <Sgntr/>
    </AppHdr>
    <Document>
        <FIToFICstmrCdtTrf>
            <GrpHdr>
                <MsgId>%s</MsgId>
                <CreDtTm>%s</CreDtTm>
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
                    <EndToEndId>%s</EndToEndId>
                </PmtId>
                <IntrBkSttlmAmt Ccy="BRL">1.00</IntrBkSttlmAmt>
                <AccptncDtTm>%s</AccptncDtTm>
                <ChrgBr>SLEV</ChrgBr>
                <MndtRltdInf>
                    <Tp>
                        <LclInstrm>
                            <Prtry>MANU</Prtry>
                        </LclInstrm>
                    </Tp>
                </MndtRltdInf>
                <Dbtr>
                    <Nm>Rogerio Inacio</Nm>
                    <Id>
                        <PrvtId>
                            <Othr>
                                <Id>14811554744</Id>
                            </Othr>
                        </PrvtId>
                    </Id>
                </Dbtr>
                <DbtrAcct>
                    <Id>
                        <Othr>
                            <Id>0038952</Id>
                            <Issr>0001</Issr>
                        </Othr>
                    </Id>
                    <Tp>
                        <Cd>CACC</Cd>
                    </Tp>
                </DbtrAcct>
                <DbtrAgt>
                    <FinInstnId>
                        <ClrSysMmbId>
                            <MmbId>52833288</MmbId>
                        </ClrSysMmbId>
                    </FinInstnId>
                </DbtrAgt>
                <CdtrAgt>
                    <FinInstnId>
                        <ClrSysMmbId>
                            <MmbId>99999004</MmbId>
                        </ClrSysMmbId>
                    </FinInstnId>
                </CdtrAgt>
                <Cdtr>
                    <Id>
                        <PrvtId>
                            <Othr>
                                <Id>68163319000171</Id>
                            </Othr>
                        </PrvtId>
                    </Id>
                </Cdtr>
                <CdtrAcct>
                    <Id>
                        <Othr>
                            <Id>79396001</Id>
                            <Issr>0001</Issr>
                        </Othr>
                    </Id>
                    <Tp>
                        <Cd>CACC</Cd>
                    </Tp>
                </CdtrAcct>
                <Purp>
                    <Cd>IPAY</Cd>
                </Purp>
                <RmtInf>
                    <Ustrd>Teste 2</Ustrd>
                </RmtInf>
            </CdtTrfTxInf>
        </FIToFICstmrCdtTrf>
    </Document>
</Envelope>`

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
        <BizMsgIdr>M52833288bfefffd6533b49708ba8101</BizMsgIdr>
        <MsgDefIdr>pibr.001.spi.1.3</MsgDefIdr>
        <CreDt>2020-04-07T13:47:22.580Z</CreDt>
        <Sgntr/>
    </AppHdr>
    <Document>
        <EchoReq>
            <GrpHdr>
                <MsgId>M52833288bfefffd6533b49708ba8101</MsgId>
                <CreDtTm>2020-04-07T13:47:22.580Z</CreDtTm>
            </GrpHdr>
            <EchoTxInf>
                <Data>Campo livre</Data>
            </EchoTxInf>
        </EchoReq>
    </Document>
</Envelope>`

var camt060_saldo_momento = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<Envelope xmlns="https://www.bcb.gov.br/pi/camt.060/1.9">
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
        <BizMsgIdr>%s</BizMsgIdr>
        <MsgDefIdr>camt.060.spi.1.9</MsgDefIdr>
        <CreDt>%s</CreDt>
        <Sgntr/>
    </AppHdr>
    <Document>
        <AcctRptgReq>
            <GrpHdr>
                <MsgId>%s</MsgId>
                <CreDtTm>%s</CreDtTm>
            </GrpHdr>
            <RptgReq>
                <ReqdMsgNmId>camt.053</ReqdMsgNmId>
                <AcctOwnr>
                    <Agt>
                        <FinInstnId>
                            <ClrSysMmbId>
                                <MmbId>52833288</MmbId>
                            </ClrSysMmbId>
                        </FinInstnId>
                    </Agt>
                </AcctOwnr>
                <ReqdBalTp>
                    <CdOrPrtry>
                        <Prtry>CSA</Prtry>
                    </CdOrPrtry>
                </ReqdBalTp>
            </RptgReq>
        </AcctRptgReq>
    </Document>
</Envelope>`

var camt060_saldo_data_anterior = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<Envelope xmlns="https://www.bcb.gov.br/pi/camt.060/1.9">
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
        <BizMsgIdr>%s</BizMsgIdr>
        <MsgDefIdr>camt.060.spi.1.9</MsgDefIdr>
        <CreDt>%s</CreDt>
        <Sgntr/>
    </AppHdr>
    <Document>
        <AcctRptgReq>
            <GrpHdr>
                <MsgId>%s</MsgId>
                <CreDtTm>%s</CreDtTm>
            </GrpHdr>
            <RptgReq>
                <ReqdMsgNmId>camt.053</ReqdMsgNmId>
                <AcctOwnr>
                    <Agt>
                        <FinInstnId>
                            <ClrSysMmbId>
                                <MmbId>52833288</MmbId>
                            </ClrSysMmbId>
                        </FinInstnId>
                    </Agt>
                </AcctOwnr>
                <RptgPrd>
                    <FrToDt>
                        <FrDt>2025-08-20</FrDt>
                    </FrToDt>
                    <Tp>ALLL</Tp>
                </RptgPrd>
                <ReqdBalTp>
                    <CdOrPrtry>
                        <Prtry>CSA</Prtry>
                    </CdOrPrtry>
                </ReqdBalTp>
            </RptgReq>
            </RptgReq>
        </AcctRptgReq>
    </Document>
</Envelope>`

var camt060_arquivo_trd = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<Envelope xmlns="https://www.bcb.gov.br/pi/camt.060/1.9">
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
        <BizMsgIdr>%s</BizMsgIdr>
        <MsgDefIdr>camt.060.spi.1.9</MsgDefIdr>
        <CreDt>%s</CreDt>
        <Sgntr/>
    </AppHdr>
    <Document>
        <AcctRptgReq>
            <GrpHdr>
                <MsgId>%s</MsgId>
                <CreDtTm>%s</CreDtTm>
            </GrpHdr>
            <RptgReq>
                <ReqdMsgNmId>camt.052</ReqdMsgNmId>
                <AcctOwnr>
                    <Agt>
                        <FinInstnId>
                            <ClrSysMmbId>
                                <MmbId>52833288</MmbId>
                            </ClrSysMmbId>
                        </FinInstnId>
                    </Agt>
                </AcctOwnr>
                <RptgPrd>
                    <FrToDt>
                        <FrDt>2025-08-20</FrDt>
                    </FrToDt>
                    <Tp>ALLL</Tp>
                </RptgPrd>
                <ReqdBalTp>
                    <CdOrPrtry>
                        <Prtry>TRD</Prtry>
                    </CdOrPrtry>
                </ReqdBalTp>
            </RptgReq>
        </AcctRptgReq>
    </Document>
</Envelope>`

var pacs004 = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<Envelope xmlns="https://www.bcb.gov.br/pi/pacs.004/1.5">
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
        <BizMsgIdr>%s</BizMsgIdr>
        <MsgDefIdr>pacs.004.spi.1.5</MsgDefIdr>
        <CreDt>%s</CreDt>
        <Sgntr/>
    </AppHdr>
    <Document>
        <PmtRtr>
            <GrpHdr>
                <MsgId>%s</MsgId>
                <CreDtTm>%s</CreDtTm>
                <NbOfTxs>1</NbOfTxs>
                <SttlmInf>
                    <SttlmMtd>CLRG</SttlmMtd>
                </SttlmInf>
            </GrpHdr>
            <TxInf>
                <RtrId>%s</RtrId>
                <OrgnlEndToEndId>%s</OrgnlEndToEndId>
                <RtrdIntrBkSttlmAmt Ccy="BRL">1000.00</RtrdIntrBkSttlmAmt>
                <SttlmPrty>HIGH</SttlmPrty>
                <ChrgBr>SLEV</ChrgBr>
                <RtrRsnInf>
                    <Rsn>
                        <Cd>BE08</Cd>
                    </Rsn>
                </RtrRsnInf>
                <OrgnlTxRef>
                    <DbtrAgt>
                        <FinInstnId>
                            <ClrSysMmbId>
                                <MmbId>99999009</MmbId>
                            </ClrSysMmbId>
                        </FinInstnId>
                    </DbtrAgt>
                    <CdtrAgt>
                        <FinInstnId>
                            <ClrSysMmbId>
                                <MmbId>99999010</MmbId>
                            </ClrSysMmbId>
                        </FinInstnId>
                    </CdtrAgt>
                </OrgnlTxRef>
            </TxInf>
        </PmtRtr>
    </Document>
</Envelope>`

var pacs002 = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<Envelope xmlns="https://www.bcb.gov.br/pi/pacs.002/1.14">
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
        <BizMsgIdr>%s</BizMsgIdr>
        <MsgDefIdr>pacs.002.spi.1.14</MsgDefIdr>
        <CreDt>%s</CreDt>
        <Sgntr/>
    </AppHdr>
    <Document>
        <FIToFIPmtStsRpt>
            <GrpHdr>
                <MsgId>%s</MsgId>
                <CreDtTm>%s</CreDtTm>
            </GrpHdr>
            <TxInfAndSts>
                <OrgnlInstrId>%s</OrgnlInstrId>
                <OrgnlEndToEndId>%s</OrgnlEndToEndId>
                <TxSts>ACSP</TxSts>
            </TxInfAndSts>
        </FIToFIPmtStsRpt>
    </Document>
</Envelope>`

func callJavaFunction(message string) (string, error) {
	// Run Java program with message as argument
	cmd := exec.Command("java", "-jar", "/home/roger/projects/lb/signer-java/target/signer-java-1.0-SNAPSHOT.jar", "-a", message)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run Java: %v", err)
	}

	fmt.Print(string(output))

	return string(output), nil
}

func main() {
	//now := getCurrentTime().UTC()

	//endToEndID, _ := GenerateEndToEndId("52833288", now)
	//endToEndID2, _ := GenerateEndToEndId("00038166", now)

	//id, _ := GenerateMsgId("52833288")

	//returnId := GenerateReturnId("52833288")

	//ready := fmt.Sprintf(pacs008, id, now.Format("2006-01-02T15:04:05.000Z"), id, now.Format("2006-01-02T15:04:05.000Z"), endToEndID, now.Format("2006-01-02T15:04:05.000Z")) //pacs008

	//ready := fmt.Sprintf(pacs004, id, now.Format("2006-01-02T15:04:05.000Z"), id, now.Format("2006-01-02T15:04:05.000Z"), returnId, endToEndID) //pacs004
	//ready := fmt.Sprintf(pacs002, id, now.Format("2006-01-02T15:04:05.000Z"), id, now.Format("2006-01-02T15:04:05.000Z"), endToEndID, endToEndID2) //pacs002
	//ready := fmt.Sprintf(camt060_saldo_momento, id, now.Format("2006-01-02T15:04:05.000Z"), id, now.Format("2006-01-02T15:04:05.000Z")) //camt060_saldo_momento
	//ready := fmt.Sprintf(camt060_saldo_data_anterior, id, now.Format("2006-01-02T15:04:05.000Z"), id, now.Format("2006-01-02T15:04:05.000Z")) //camt060_saldo_data_anterior
	//ready := fmt.Sprintf(camt060_arquivo_trd, id, now.Format("2006-01-02T15:04:05.000Z"), id, now.Format("2006-01-02T15:04:05.000Z")) //camt060_arquivo_trd

	conn := services.CreateConnection()

	if conn == nil {
		panic("Failed to create TLS connection")
	}

	//str, err := callJavaFunction(ready)
	str, err := callJavaFunction(pibr001)

	fmt.Printf("Generated PACS.008 message:\n%s\n", str)

	if err != nil {
		fmt.Println("Error calling Java function:", err)
		return
	}

	err = services.PostMessage(conn, str)

	if err != nil {
		fmt.Println("Error posting message:", err)
	}
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

func GenerateReturnId(ispb string) string {
	// Data/hora atual em UTC no formato yyyyMMddHHmm
	now := time.Now().UTC()
	timestamp := now.Format("200601021504")

	// Sequencial único de 11 caracteres alfanuméricos
	sequential := GenerateRandomAlphanumeric(11)

	return fmt.Sprintf("%s%s%s%s", returnIdPrefix, ispb, timestamp, sequential)
}
