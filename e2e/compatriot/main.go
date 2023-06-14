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

package main

import (
	"fmt"

	"github.com/fatih/color"
)

const CACHED = "./cached.json"
const NONCACHED = "./noncached.json"

func main() {
	setup()

	// make queries and save results to file 1
	makeCalls(CACHED)

	// kill the chain

	// make queries adn save results to file 2
	makeCalls(NONCACHED)

	// compare file 1 and file 2

	// run sanity checks

	// print results

	color.Set(color.FgGreen)
	fmt.Println("The following JSON-RPC methods are likely supported in your EVM chain:")
	for _, val := range supportedMethods {
		fmt.Println(val)
	}
	fmt.Println()

	color.Set(color.FgYellow)
	fmt.Println("The following JSON-RPC methods may or may not be supported in your EVM chain:")
	for _, val := range possiblySupportedMethods {
		fmt.Println(val)
	}
	fmt.Println()

	color.Set(color.FgRed)
	fmt.Println("The following JSON-RPC methods are likely unsupported in your EVM chain:")
	for _, val := range unsupportedMethods {
		fmt.Println(val)
	}
	fmt.Println()
}
