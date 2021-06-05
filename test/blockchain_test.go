package test

import (
	"BrunoCoin/pkg/block"
	"BrunoCoin/pkg/block/tx"
	"BrunoCoin/pkg/block/tx/txi"
	"BrunoCoin/pkg/block/tx/txo"
	"BrunoCoin/pkg/blockchain"
	"encoding/hex"
	"testing"
)

func TestAdd(t *testing.T) {
	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	testConf := blockchain.DefaultConfig()

	bc := blockchain.New(testConf)

	testTxi := &txi.TransactionInput {
		TransactionHash: "",
		OutputIndex:     0,
		UnlockingScript: "",
		Amount:          0,
	}

	testTxo := &txo.TransactionOutput{
		Amount:        0,
		LockingScript: "",
		Liminal:       false,
	}

	t1 := &tx.Transaction{
		Version:  0,
		Inputs:   []*txi.TransactionInput{testTxi},
		Outputs:  []*txo.TransactionOutput{testTxo},
		LockTime: 0,
	}

	transactions := []*tx.Transaction{t1}

	byte1 := byte(0)

	testBlk := block.New(genNd.Chain.LastBlock.Hash(), transactions, hex.EncodeToString([]byte{byte1}))

	bc.Add(testBlk)
}