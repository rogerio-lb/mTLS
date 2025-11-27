package services

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	ramdom "math/rand"
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
	returnIdAlphabet     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	FRAUD_REASON         = "FR01"
	USER_REQUEST_REASON  = "MD06"
	BANK_ERROR_REASON    = "BE08"
	SERVICE_CAUSE_REASON = "SL02"
)

var (
	pacsRandReader      io.Reader = rand.Reader
	numericPattern                = regexp.MustCompile(`^[0-9]{8}$`)
	alphanumericPattern           = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	getCurrentTime                = time.Now
	randReader          io.Reader = rand.Reader
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
                <RtrdIntrBkSttlmAmt Ccy="BRL">50</RtrdIntrBkSttlmAmt>
                <SttlmPrty>HIGH</SttlmPrty>
                <ChrgBr>SLEV</ChrgBr>
                <RtrRsnInf>
                    <Rsn>
                        <Cd>FR01</Cd>
                    </Rsn>
                </RtrRsnInf>
                <OrgnlTxRef>
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
                                <MmbId>58160789</MmbId>
                            </ClrSysMmbId>
                        </FinInstnId>
                    </CdtrAgt>
                </OrgnlTxRef>
            </TxInf>
        </PmtRtr>
    </Document>
</Envelope>`

var pacs004_dict = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
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
                <RtrdIntrBkSttlmAmt Ccy="BRL">%s</RtrdIntrBkSttlmAmt>
                <SttlmPrty>HIGH</SttlmPrty>
                <ChrgBr>SLEV</ChrgBr>
                <RtrRsnInf>
                    <Rsn>
                        <Cd>%s</Cd>
                    </Rsn>
                </RtrRsnInf>
                <OrgnlTxRef>
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
                                <MmbId>%s</MmbId>
                            </ClrSysMmbId>
                        </FinInstnId>
                    </CdtrAgt>
                </OrgnlTxRef>
            </TxInf>
        </PmtRtr>
    </Document>
</Envelope>`

var pacs002 = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<Envelope xmlns="https://www.bcb.gov.br/pi/pacs.002/1.15">
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
        <MsgDefIdr>pacs.002.spi.1.15</MsgDefIdr>
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

var pacs002_pacs004 = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<Envelope xmlns="https://www.bcb.gov.br/pi/pacs.002/1.15">
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
        <MsgDefIdr>pacs.002.spi.1.15</MsgDefIdr>
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

var pacs008_511_manu = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<Envelope xmlns="https://www.bcb.gov.br/pi/pacs.008/1.14">
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
        <MsgDefIdr>pacs.008.spi.1.14</MsgDefIdr>
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
                <IntrBkSttlmAmt Ccy="BRL">5.00</IntrBkSttlmAmt>
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
                                <Id>61363314000143</Id>
                            </Othr>
                        </PrvtId>
                    </Id>
                </Cdtr>
                <CdtrAcct>
                    <Id>
                        <Othr>
                            <Id>90570189</Id>
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
                    <Ustrd>Teste 1</Ustrd>
                </RmtInf>
            </CdtTrfTxInf>
        </FIToFICstmrCdtTrf>
    </Document>
</Envelope>`

var pacs008_511_dict = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<Envelope xmlns="https://www.bcb.gov.br/pi/pacs.008/1.14">
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
        <MsgDefIdr>pacs.008.spi.1.14</MsgDefIdr>
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
                <IntrBkSttlmAmt Ccy="BRL">50</IntrBkSttlmAmt>
                <AccptncDtTm>%s</AccptncDtTm>
                <ChrgBr>SLEV</ChrgBr>
                <MndtRltdInf>
                    <Tp>
                        <LclInstrm>
                            <Prtry>DICT</Prtry>
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
                            <MmbId>58160789</MmbId>
                        </ClrSysMmbId>
                    </FinInstnId>
                </CdtrAgt>
                <Cdtr>
                    <Id>
                        <PrvtId>
                            <Othr>
                                <Id>00011098000182</Id>
                            </Othr>
                        </PrvtId>
                    </Id>
                </Cdtr>
                <CdtrAcct>
                    <Id>
                        <Othr>
                            <Id>008947835</Id>
                            <Issr>0097</Issr>
                        </Othr>
                    </Id>
                    <Tp>
                        <Cd>CACC</Cd>
                    </Tp>
                    <Prxy>
                        <Id>+5511981047272</Id>
                    </Prxy>
                </CdtrAcct>
                <Purp>
                    <Cd>IPAY</Cd>
                </Purp>
                <RmtInf>
                    <Ustrd>Teste</Ustrd>
                </RmtInf>
            </CdtTrfTxInf>
        </FIToFICstmrCdtTrf>
    </Document>
</Envelope>`

func CreateMessage() string {
	now := getCurrentTime().UTC()
	endToEndID, _ := GenerateEndToEndId("52833288", now)
	id, _ := GenerateMsgId("52833288")

	//returnId := GenerateReturnId("52833288")

	//ready := fmt.Sprintf(pacs008, id, now.Format("2006-01-02T15:04:05.000Z"), id, now.Format("2006-01-02T15:04:05.000Z"), endToEndID, now.Format("2006-01-02T15:04:05.000Z")) //pacs008
	ready := fmt.Sprintf(pacs008_511_manu, id, now.Format("2006-01-02T15:04:05.000Z"), id, now.Format("2006-01-02T15:04:05.000Z"), endToEndID, now.Format("2006-01-02T15:04:05.000Z")) //pacs008
	//ready := fmt.Sprintf(pacs008_511_dict, id, now.Format("2006-01-02T15:04:05.000Z"), id, now.Format("2006-01-02T15:04:05.000Z"), endToEndID, now.Format("2006-01-02T15:04:05.000Z")) //pacs008

	//ready := fmt.Sprintf(pacs004, id, now.Format("2006-01-02T15:04:05.000Z"), id, now.Format("2006-01-02T15:04:05.000Z"), returnId, "E581607892025112619453PDvqKHESKg") //pacs004
	//ready := fmt.Sprintf(pacs002, id, now.Format("2006-01-02T15:04:05.000Z"), id, now.Format("2006-01-02T15:04:05.000Z"), endToEndID, endToEndID2) //pacs002
	//ready := fmt.Sprintf(camt060_saldo_momento, id, now.Format("2006-01-02T15:04:05.000Z"), id, now.Format("2006-01-02T15:04:05.000Z")) //camt060_saldo_momento
	//ready := fmt.Sprintf(camt060_saldo_data_anterior, id, now.Format("2006-01-02T15:04:05.000Z"), id, now.Format("2006-01-02T15:04:05.000Z")) //camt060_saldo_data_anterior
	//ready := fmt.Sprintf(camt060_arquivo_trd, id, now.Format("2006-01-02T15:04:05.000Z"), id, now.Format("2006-01-02T15:04:05.000Z")) //camt060_arquivo_trd

	str, err := SignMessage(ready)
	//str, err := services.SignMessage(pibr001)
	if err != nil {
		fmt.Println("Error calling Java function:", err)
		return ""
	}

	fmt.Println("EndToEndID:", endToEndID)

	return str
}

func GeneratePacs002(e2eID string) string {
	now := getCurrentTime().UTC()
	//endToEndID, _ := GenerateEndToEndId("52833288", now)
	id, _ := GenerateMsgId("52833288")

	ready := fmt.Sprintf(pacs002, id, now.Format("2006-01-02T15:04:05.000Z"), id, now.Format("2006-01-02T15:04:05.000Z"), e2eID, e2eID) //pacs002
	str, err := SignMessage(ready)
	if err != nil {
		fmt.Println("Error calling Java function:", err)
		return ""
	}

	return str
}

func GeneratePacs004ForDict(e2eID, ispb, value, reason string) string {
	now := getCurrentTime().UTC()
	id, _ := GenerateMsgId("52833288")
	returnId := GenerateReturnId("52833288")

	ready := fmt.Sprintf(
		pacs004_dict,
		id,
		now.Format("2006-01-02T15:04:05.000Z"),
		id,
		now.Format("2006-01-02T15:04:05.000Z"),
		returnId,
		e2eID,
		value,
		reason,
		ispb,
	) //pacs004
	str, err := SignMessage(ready)
	if err != nil {
		fmt.Println("Error calling Java function:", err)
		return ""
	}

	return str
}

func GeneratePacs002ForPacs004(e2eID, returnId string) string {
	now := getCurrentTime().UTC()
	//endToEndID, _ := GenerateEndToEndId("52833288", now)
	id, _ := GenerateMsgId("52833288")

	ready := fmt.Sprintf(pacs002_pacs004, id, now.Format("2006-01-02T15:04:05.000Z"), id, now.Format("2006-01-02T15:04:05.000Z"), returnId, e2eID) //pacs002
	str, err := SignMessage(ready)
	if err != nil {
		fmt.Println("Error calling Java function:", err)
		return ""
	}

	return str
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

func CompressContentToGzip(body []byte, buffer *bytes.Buffer) error {
	gzipWriter := gzip.NewWriter(buffer)
	defer gzipWriter.Close()

	if _, err := gzipWriter.Write(body); err != nil {
		return err
	}

	if err := gzipWriter.Close(); err != nil {
		return err
	}

	return nil
}
