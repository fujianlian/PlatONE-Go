package packet

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/PlatONEnetwork/PlatONE-Go/core/types"

	precompile "github.com/PlatONEnetwork/PlatONE-Go/cmd/platonecli/precompiled"
	"github.com/PlatONEnetwork/PlatONE-Go/cmd/platonecli/utils"
	"github.com/PlatONEnetwork/PlatONE-Go/common"
	"github.com/PlatONEnetwork/PlatONE-Go/common/byteutil"
	"github.com/PlatONEnetwork/PlatONE-Go/common/hexutil"
	"github.com/PlatONEnetwork/PlatONE-Go/crypto"
	"github.com/PlatONEnetwork/PlatONE-Go/rlp"
)

var (
	txReceiptSuccessCode = hexutil.EncodeUint64(types.ReceiptStatusSuccessful)
	txReceiptFailureCode = hexutil.EncodeUint64(types.ReceiptStatusFailed)
)

const (
	txReceiptSuccessMsg = "Operation Succeeded"
	txReceiptFailureMsg = "Operation Failed"
)

// Receipt, eth_getTransactionReceipt return data struct
type Receipt struct {
	BlockHash         string    `json:"blockHash"`          // hash of the block
	BlockNumber       string    `json:"blockNumber"`        // height of the block
	ContractAddress   string    `json:"contractAddress"`    // contract address of the contract deployment. otherwise null
	CumulativeGasUsed string    `json:"cumulativeGas_used"` //
	From              string    `json:"from"`               // the account address used to send the transaction
	GasUsed           string    `json:"gasUsed"`            // gas used by executing the transaction
	Root              string    `json:"root"`
	To                string    `json:"to"`               // the address the transaction is sent to
	TransactionHash   string    `json:"transactionHash"`  // the hash of the transaction
	TransactionIndex  string    `json:"transactionIndex"` // the index of the transaction
	Logs              RecptLogs `json:"logs"`
	Status            string    `json:"status"` // the execution status of the transaction, "0x1" for success
}

type Log struct {
	Address string   `json:"address"`
	Topics  []string `json:"topics"`
	Data    string   `json:"data"`
}

type RecptLogs []*Log

// ParseSysContractResult parsed the rpc response to Receipt object
func ParseTxReceipt(response interface{}) (*Receipt, error) {
	var receipt = &Receipt{}

	temp, _ := json.Marshal(response)
	err := json.Unmarshal(temp, receipt)
	if err != nil {
		// LogErr.Printf(ErrUnmarshalBytesFormat, "transaction receipt", err.Error())
		errStr := fmt.Sprintf(utils.ErrUnmarshalBytesFormat, "transaction receipt", err.Error())
		return nil, errors.New(errStr)
	}

	return receipt, nil
}

// EventParsing parsing all the events recorded in receipt log.
// The event should be written in the abiBytes provided.
// Otherwise, the event will not be parsed
func EventParsing(logs RecptLogs, abiBytes []byte) (result string) {
	var rlpList []interface{}
	var eventName string
	var topicTypes []string
	var cnsInvokeAddr = precompile.CnsInvokeAddress

	for i, eLog := range logs {
		// for cns invoking, to parse the events defined in sc_cns_invoke.go
		// we add if condition to check the address in the receipt logs
		if strings.EqualFold(eLog.Address, cnsInvokeAddr) {
			p := precompile.List[cnsInvokeAddr]
			cnsInvokeAbiBytes, _ := precompile.Asset(p)
			eventName, topicTypes = findLogTopic(eLog.Topics[0], cnsInvokeAbiBytes)
		} else {
			eventName, topicTypes = findLogTopic(eLog.Topics[0], abiBytes)
		}

		if len(topicTypes) == 0 {
			continue
		}

		dataBytes, _ := hexutil.Decode(eLog.Data)
		err := rlp.DecodeBytes(dataBytes, &rlpList)
		if err != nil {
			// todo: error handle
			fmt.Printf("the error is %v\n", err)
		}
		result += fmt.Sprintf("\nEvent[%d]: %s ", i, eventName)
		result += parseReceiptLogData(rlpList, topicTypes)
	}

	return
}

func ReceiptParsing(receipt *Receipt, abiBytes []byte) string {
	var result string

	switch {
	case len(receipt.Logs) != 0:
		result = EventParsing(receipt.Logs, abiBytes)

	case receipt.Status == txReceiptFailureCode:
		result = txReceiptFailureMsg

	case receipt.ContractAddress != "":
		result = receipt.ContractAddress

	case receipt.Status == txReceiptSuccessCode:
		result = txReceiptSuccessMsg
	}

	return result
}

func findLogTopic(topic string, abiBytes []byte) (string, []string) {
	var types []string
	var name string

	abiFunc, _ := ParseAbiFromJson(abiBytes)

	for _, data := range abiFunc {
		if data.Type != "event" {
			continue
		}

		if strings.EqualFold(logTopicEncode(data.Name), topic) {
			name = data.Name
			for _, v := range data.Inputs {
				types = append(types, v.Type)
			}
			break
		}
	}

	return name, types
}

func parseReceiptLogData(data []interface{}, types []string) string {
	var str string

	for i, v := range data {
		result := ConvertRlpBytesTo(v.([]uint8), types[i])
		str += fmt.Sprintf("%v ", result)
	}

	return str
}

func logTopicEncode(name string) string {
	return common.BytesToHash(crypto.Keccak256([]byte(name))).String()
}

func ConvertRlpBytesTo(input []byte, targetType string) interface{} {
	v, ok := Bytes2X_CMD[targetType]
	if !ok {
		panic("unsupported type")
	}

	return reflect.ValueOf(v).Call([]reflect.Value{reflect.ValueOf(input)})[0].Interface()
}

var Bytes2X_CMD = map[string]interface{}{
	"string": byteutil.BytesToString,

	// "uint8":  RlpBytesToUint,
	"uint16": RlpBytesToUint16,
	"uint32": RlpBytesToUint32,
	"uint64": RlpBytesToUint64,

	// "uint8":  RlpBytesToUint,
	"int16": RlpBytesToUint16,
	"int32": RlpBytesToUint32,
	"int64": RlpBytesToUint64,

	"bool": RlpBytesToBool,
}

func RlpBytesToUint16(b []byte) uint16 {
	b = common.LeftPadBytes(b, 32)
	result := common.CallResAsUint32(b)
	return uint16(result)
}

func RlpBytesToUint32(b []byte) uint32 {
	b = common.LeftPadBytes(b, 32)
	return common.CallResAsUint32(b)
}

func RlpBytesToUint64(b []byte) uint64 {
	b = common.LeftPadBytes(b, 32)
	return common.CallResAsUint64(b)
}

func RlpBytesToBool(b []byte) bool {
	if bytes.Compare(b, []byte{1}) == 0 {
		return true
	}
	return false
}

/*
func RlpBytesToUintV2(b []byte) interface{} {
	var val interface{}

	for _, v := range b {
		val = val << 8
		val |= uint(v)
	}

	return val
}*/
