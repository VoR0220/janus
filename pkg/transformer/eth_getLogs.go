package transformer

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
	"github.com/qtumproject/janus/pkg/utils"
)

// ProxyETHGetLogs implements ETHProxy
type ProxyETHGetLogs struct {
	*qtum.Qtum
}

func (p *ProxyETHGetLogs) Method() string {
	return "eth_getLogs"
}

func (p *ProxyETHGetLogs) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	var req eth.GetLogsRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	qtumreq, err := p.ToRequest(&req)
	if err != nil {
		return nil, err
	}

	return p.request(qtumreq)
}

func (p *ProxyETHGetLogs) request(req *qtum.SearchLogsRequest) (*eth.GetLogsResponse, error) {
	receipts, err := p.SearchLogs(req)
	if err != nil {
		return nil, err
	}

	logs := make([]eth.Log, 0)
	for _, receipt := range receipts {
		r := qtum.TransactionReceiptStruct(receipt)
		logs = append(logs, getEthLogs(&r)...)
	}

	resp := eth.GetLogsResponse(logs)
	return &resp, nil
}

func (p *ProxyETHGetLogs) ToRequest(ethreq *eth.GetLogsRequest) (*qtum.SearchLogsRequest, error) {
	from, err := getQtumBlockNumber(ethreq.FromBlock, 0)
	if err != nil {
		return nil, err
	}

	to, err := getQtumBlockNumber(ethreq.ToBlock, -1)
	if err != nil {
		return nil, err
	}

	var addresses []string
	if ethreq.Address != nil {
		if isString(ethreq.Address) {
			var addr string
			if err = json.Unmarshal(ethreq.Address, &addr); err != nil {
				return nil, err
			}
			addresses = append(addresses, addr)
		} else {
			if err = json.Unmarshal(ethreq.Address, &addresses); err != nil {
				return nil, err
			}
		}
		for i := range addresses {
			addresses[i] = utils.RemoveHexPrefix(addresses[i])
		}
	}

	return &qtum.SearchLogsRequest{
		Addresses: addresses,
		FromBlock: from,
		ToBlock:   to,
		Topics:    ethreq.Topics,
	}, nil
}

func getQtumBlockNumber(ethBlock json.RawMessage, defaultVal int64) (*big.Int, error) {
	if ethBlock == nil {
		return big.NewInt(defaultVal), nil
	}

	if isString(ethBlock) {
		var ethBlockStr string
		if err := json.Unmarshal(ethBlock, &ethBlockStr); err != nil {
			return nil, err
		}

		switch ethBlockStr {
		case "latest":
			return big.NewInt(-1), nil
		case "pending", "earliest":
			return nil, errors.New(`tags, "pending" and "earliest", are unsupported`)
		default:
			return utils.DecodeBig(ethBlockStr)
		}
	}

	var b int64
	if err := json.Unmarshal(ethBlock, &b); err != nil {
		return nil, err
	}

	return big.NewInt(b), nil
}

func isString(v json.RawMessage) bool {
	return v[0] == '"'
}

func getEthLogs(receipt *qtum.TransactionReceiptStruct) []eth.Log {
	logs := make([]eth.Log, 0, len(receipt.Log))
	for index, log := range receipt.Log {
		topics := make([]string, 0, len(log.Topics))
		for _, topic := range log.Topics {
			topics = append(topics, utils.AddHexPrefix(topic))
		}
		logs = append(logs, eth.Log{
			TransactionHash:  utils.AddHexPrefix(receipt.TransactionHash),
			TransactionIndex: hexutil.EncodeUint64(receipt.TransactionIndex),
			BlockHash:        utils.AddHexPrefix(receipt.BlockHash),
			BlockNumber:      hexutil.EncodeUint64(receipt.BlockNumber),
			Data:             utils.AddHexPrefix(log.Data),
			Address:          utils.AddHexPrefix(log.Address),
			Topics:           topics,
			LogIndex:         hexutil.EncodeUint64(uint64(index)),
		})
	}
	return logs
}
