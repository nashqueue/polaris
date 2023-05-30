// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2023, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package state

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/core"

	"pkg.berachain.dev/polaris/eth/common"
)

// InitGenesis takes in a pointer to a genesis state object and populates the KV store.
func (p *plugin) InitGenesis(ctx sdk.Context, ethGen *core.Genesis) {
	p.Reset(ctx)

	// Iterate over the genesis accounts and set the balances.
	for address, account := range ethGen.Alloc {
		// TODO: technically wrong since its overriding / hacking the auth keeper and
		// we are using the nonce from the account keeper as well.
		p.CreateAccount(address)
		p.SetBalance(address, account.Balance)
		if account.Code != nil {
			p.SetCode(address, account.Code)
			for k, v := range account.Storage {
				p.SetState(address, k, v)
			}
		}
	}
	p.Finalize()
}

// Export genesis modifies a pointer to a genesis state object and populates it.
func (p *plugin) ExportGenesis(ctx sdk.Context, ethGen *core.Genesis) {
	p.Reset(ctx)
	for address, account := range ethGen.Alloc {
		account.Code = p.GetCode(address)
		account.Balance = p.GetBalance(address)
		account.Storage = make(map[common.Hash]common.Hash)
		err := p.ForEachStorage(address, func(key, value common.Hash) bool {
			account.Storage[key] = value
			return false
		})
		if err != nil {
			panic(err)
		}
	}
	p.Finalize()
}
