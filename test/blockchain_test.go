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

func TestGetUTXOForAmt(t *testing.T) {
	inputs := []*txi.TransactionInput{
		{
			TransactionHash: "123",
			OutputIndex:     0,
			UnlockingScript: "123",
			Amount:          100,
		}}

	outputs := []*txo.TransactionOutput{
		{
			Amount:        99,
			LockingScript: "123",
			Liminal:       false,
		},
	}

	transaction := tx.Transaction{
		Version:  1,
		Inputs:   inputs,
		Outputs:  outputs,
		LockTime: 0,
	}

	bc := blockchain.New(blockchain.DefaultConfig())
	bc.Add(&block.Block{
		Hdr: block.Header{
			PrvBlkHsh: bc.LastBlock.Hash(),
		},
		Transactions: []*tx.Transaction{&transaction},
	})

	a, b, c := bc.GetUTXOForAmt(100, "123")

	if c {
		t.Errorf("received: %v, %v, %v", a[0].Amt, b, c)
	}
}

func TestAdd(t *testing.T) {
		genNd := NewGenNd()
		genNd.Start()
		genNd.StartMiner()

		testConf := blockchain.DefaultConfig()

		bc := blockchain.New(testConf)

		testTxi := &txi.TransactionInput{
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

		if bc.LastBlock.Block != testBlk {
			t.Errorf("ERROR {blockchain.Add}: Failed to add block")
		}

}