// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

func consistentChainIDTest(t *TestEnv) {
	var (
		expectedChainID = big.NewInt(7) //nolint:gomnd // TODO: REFACTOR.
	)

	cID, err := t.Eth.ChainID(t.Ctx())
	if err != nil {
		t.Fatalf("could not get chain ID: %v", err)
	}

	if expectedChainID.Cmp(cID) != 0 {
		t.Fatalf("expected chain ID %d, got %d", expectedChainID, cID)
	}
}

func eth_gasPriceTest(t *TestEnv) {

	initialBaseFee := params.InitialBaseFee

	gasPrice, err := t.Eth.SuggestGasPrice(t.Ctx())
	if err != nil {
		t.Fatalf("could not get gas price: %v", err)
	}

	if gasPrice.Cmp(big.NewInt(int64(initialBaseFee))) != 0 {
		t.Fatalf("expected gas price %d, got %d", initialBaseFee, gasPrice)
	}

}

func eth_blockNumberTest(t *TestEnv) {

	_, err := t.Eth.BlockNumber(t.Ctx())

	if err != nil {
		t.Fatalf("could not get block number: %v", err)
	}

}
func eth_getBalance(t *TestEnv) {

	// fix
	balance, err := t.Eth.BalanceAt(t.Ctx(), common.Address{}, nil)
	if err != nil {
		t.Fatalf("could not get balance: %v", err)
	}
	if balance.Cmp(big.NewInt(0)) != 0 {
		t.Fatalf("expected balance 0, got %d", balance)
	}

}
func eth_estimateGas(t *TestEnv) {

}
func deployContract(t *TestEnv) {

}
func interactWithContract(t *TestEnv) {

}
func eth_getTransactionByHash(t *TestEnv) {

}
