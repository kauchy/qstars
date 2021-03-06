// Copyright 2018 The QOS Authors

package kvstore

import (
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/client/context"
	"github.com/QOSGroup/qstars/client/utils"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/wire"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

// ResultSendKV result of send kv
type ResultSendKV struct {
	Hash string `json:"hash"`
}

// SendKVOption option param ofr send kv
type SendKVOption struct {
	chainID  string
	sequence string
}

type SetSendKVOption func(*SendKVOption) error

// NewSendKVOption new and set option param
func NewSendKVOption(fs ...SetSendKVOption) (*SendKVOption, error) {
	sopt := &SendKVOption{
		chainID: "chainid",
	}

	if fs != nil {
		for _, f := range fs {
			if err := f(sopt); err != nil {
				return nil, err
			}
		}
	}

	return sopt, nil
}

// SendKVOptionChainID set chain id
func SendKVOptionChainID(chainID string) SetSendKVOption {
	return func(opt *SendKVOption) error {
		opt.chainID = chainID
		return nil
	}
}

// SendKVOptionSequence
func SendKVOptionSequence(sequence string) SetSendKVOption {
	return func(opt *SendKVOption) error {
		opt.sequence = sequence
		return nil
	}
}

func wrapToStdTx(key string, value string, chainid string) *txs.TxStd {
	kv := NewKvstoreTx([]byte(key), []byte(value))
	return txs.NewTxStd(kv, chainid, qbasetypes.NewInt(int64(10000)))
}

// SendKV process of set kv
func SendKV(cliCtx context.CLIContext, cdc *wire.Codec, privateKey, key, value string, option *SendKVOption) (*ResultSendKV, error) {
	//get addr from private key
	var priv ed25519.PrivKeyEd25519
	bz := utility.Decbase64(privateKey)
	copy(priv[:], bz)

	txStd := wrapToStdTx(key, value, option.chainID)

	hash, err := utils.SendTx(cliCtx, cdc, txStd, priv)
	if err != nil {
		return nil, err
	}
	result := &ResultSendKV{}
	result.Hash = hash

	return result, nil
}

// ResultGetKV result of get kv
type ResultGetKV struct {
	Value string `json:"value"`
}

// GetKVOption option param ofr get kv
type GetKVOption struct {
	chainID  string
	sequence string
}

type SetGetKVOption func(*GetKVOption) error

// NewGetKVOption new and set option param
func NewGetKVOption(fs ...SetGetKVOption) (*GetKVOption, error) {
	sopt := &GetKVOption{
		chainID: "chainid",
	}

	if fs != nil {
		for _, f := range fs {
			if err := f(sopt); err != nil {
				return nil, err
			}
		}
	}

	return sopt, nil
}

// GetKVOptionChainID set chain id
func GetKVOptionChainID(chainID string) SetGetKVOption {
	return func(opt *GetKVOption) error {
		opt.chainID = chainID
		return nil
	}
}

// GetKV process of get kv
func GetKV(cliCtx context.CLIContext, cdc *wire.Codec, key string, opt *GetKVOption) (*ResultGetKV, error) {
	value, err := cliCtx.QueryKV([]byte(key))
	if err != nil {
		return nil, err
	}
	result := &ResultGetKV{}
	result.Value = string(value)

	return result, nil
}
