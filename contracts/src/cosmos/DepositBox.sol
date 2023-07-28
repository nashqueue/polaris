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

pragma solidity ^0.8.4;

import {IERC20} from "../../lib/IERC20.sol";

contract DepositBox {
    function deposit(IERC20 token, address owner, uint256 amount) external {
        // transfer tokens first
        token.transferFrom(owner, address(this), amount);

        // read the sender's balance
        uint256 slot = slotFor(token, owner);
        uint256 balance = readFromSlot(slot);

        // balance + amount will revert on overflow.
        writeToSlot(slot, balance + amount);
    }

    function withdraw(IERC20 token, address owner, uint256 amount) external {
        // read the sender's balance
        uint256 slot = slotFor(token, owner);
        uint256 balance = readFromSlot(slot);

        // balance - amount will revert on underflow.
        writeToSlot(slot, balance - amount);

        // transfer tokens last
        token.transfer(owner, amount);
    }

    // writeToSlot
    function writeToSlot(uint256 slot, uint256 value) public {
        assembly {
            sstore(slot, value) // Store 'value' in the specified slot
        }
    }

    // readFromSlot
    function readFromSlot(uint256 slot) public view returns (uint256) {
        uint256 value;
        assembly {
            value := sload(slot) // Load the value from the specified slot
        }
        return value;
    }

    // slotFor gets the l
    function slotFor(
        IERC20 token,
        address owner
    ) public pure returns (uint256) {
        return uint256(keccak256(abi.encode(token, owner)));
    }
}
