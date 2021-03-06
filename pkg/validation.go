package pkg

import (
	"BrunoCoin/pkg/block"
	"BrunoCoin/pkg/block/tx"
	"fmt"
)

/*
 *  Brown University, CS1951L, Summer 2021
 *  Designed by: Colby Anderson, John Roy,
 *	Parker Ljung
 *
 */

// ChkBlk (CheckBlock) validates a block based on multiple
// conditions.
// To be valid:
// The block must be syntactically (ChkBlkSyn), semantically
// (ChkBlkSem), and configurally (ChkBlkConf) valid.
// Each transaction on the block must be syntactically (ChkTxSyn),
// semantically (ChkTxSem), and configurally (ChkTxConf) valid.
// Each transaction on the block must reference UTXO on the same
// chain (main or forked chain) and not be a double spend on that
// chain.
// Inputs:
// b *block.Block the block to be checked for validity
// Returns:
// bool True if the block is valid. false
// otherwise
// TODO:
// to be valid

// Each transaction on the block must reference UTXO on the same
// chain (main or forked chain) and not be a double spend on that
// chain.
// The block's size must be less than or equal to the largest
// allowed block size.
// The block hash must be less than the difficulty target.
// The block's first transaction must be of type Coinbase.

// Some helpful functions/methods/fields:
// note: let t be a transaction object
// note: let b be a block object
// t.IsCoinbase()
// b.SatisfiesPOW(...)
// n.Conf.MxBlkSz
// b.Sz()
// n.Chain.ChkChainsUTXO(...)
func (n *Node) ChkBlk(b *block.Block) bool {

	if b != nil {

		sizeLessThanMax := b.Sz() <= n.Conf.MxBlkSz
		firstIsCoinbase := false
		outputsExist := false
		txsChecked := true

		satisfiesPOW := b.SatisfiesPOW(b.Hdr.DiffTarg)

		chkChains := n.Chain.ChkChainsUTXO(b.Transactions, b.Hdr.PrvBlkHsh)

		for i, t := range b.Transactions {
			if t != nil {
				if i == 0 {
					firstIsCoinbase = t.IsCoinbase()
					outputsExist = t.SumOutputs() > 0
				} else {
					txsChecked = txsChecked && n.ChkTx(t)
				}
			} else {
				fmt.Printf("ERROR {Node.ChkBlk}: nil transaction")
				return false
			}
		}
		return sizeLessThanMax && firstIsCoinbase && outputsExist && txsChecked && chkChains && satisfiesPOW
	} else {
		fmt.Printf("ERROR {Miner.GenCBTx}: nil block")
		return false
	}
}


// ChkTx (CheckTransaction) validates a transaction.
// Inputs:
// t *tx.Transaction the transaction to be checked for validity
// Returns:
// bool True if the transaction is syntactically valid. false
// otherwise
// TODO:
// to be valid:

// The transaction's inputs and outputs must not be empty.
// The transaction's output amounts must be larger than 0.
// The sum of the transaction's inputs must be larger
// than the sum of the transaction's outputs.
// The transaction must not double spend any UTXO.
// The unlocking script on each of the transaction's
// inputs must successfully unlock each of the corresponding
// UTXO.
// The transaction must not be larger than the
// maximum allowed block size.

// Some helpful functions/methods/fields:
// note: let t be a transaction object
// note: let b be a block object
// note: let u be a transaction output object
// n.Conf.MxBlkSz
// t.Sz()
// u.IsUnlckd(...)
// n.Chain.GetUTXO(...)
// n.Chain.IsInvalidInput(...)
// t.SumInputs()
// t.SumOutputs()

func containsHash(s string, hashes []string) bool{
	for _, h := range hashes {
		if h == s {
			return true
		}
	}
	return false
}

func (n *Node) ChkTx(t *tx.Transaction) bool {

	if t != nil {

		inputsNonEmpty := len(t.Inputs) != 0
		outputsNonEmpty := len(t.Outputs) != 0
		outputsGreaterThanZero := t.SumOutputs() > 0
		inputsLargerThanOutputs := t.SumInputs() >= t.SumOutputs()
		smallerThanMaxSize := t.Sz() <= n.Conf.MxBlkSz

		seenHashes := make([]string, 0)

		isUnlockedAndValid := true
		notDoubleSpent := true
		validInputsAndOutputs := len(t.Inputs) > 0 && len(t.Outputs) > 0

		for _, i := range t.Inputs {
			if containsHash(i.Hash(), seenHashes) {
				notDoubleSpent = false
			}
			seenHashes = append(seenHashes, i.Hash())
			u := n.Chain.GetUTXO(i)
			if u != nil {
				isUnlockedAndValid =
					isUnlockedAndValid &&
						!n.Chain.IsInvalidInput(i) &&
						n.Chain.GetUTXO(i).IsUnlckd(i.UnlockingScript)
			} else {
				fmt.Printf("nil utxo")
				isUnlockedAndValid = false
			}
		}

		for _, i := range t.Outputs {
			if i.Amount <= 0  {
				validInputsAndOutputs = false
			}
		}

		return inputsNonEmpty &&
			outputsNonEmpty &&
			outputsGreaterThanZero &&
			inputsLargerThanOutputs &&
			smallerThanMaxSize &&
			isUnlockedAndValid &&
			notDoubleSpent &&
			validInputsAndOutputs
	} else {
		fmt.Printf("ERROR {Node.ChkTx}: nil transaction")
		return false
	}
}
