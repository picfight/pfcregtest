// Copyright (c) 2018 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package pfcregtest

import (
	"fmt"
	"github.com/picfight/pfcd/rpcclient"
	"strconv"
)

// Create a test chain with the desired number of mature coinbase outputs
func generateTestChain(numToGenerate uint32, node *rpcclient.Client) error {
	fmt.Printf("Generating %v blocks...\n", numToGenerate)
	_, err := node.Generate(numToGenerate)
	if err != nil {
		return err
	}
	fmt.Println("Block generation complete.")
	return nil
}

// Waits for wallet to sync to the target height
func syncWalletTo(rpcClient *rpcclient.Client, desiredHeight int64) (int64, error) {
	var count int64
	var err error
	for count != desiredHeight {
		//rpctest.Sleep(100)
		count, err = rpcClient.GetBlockCount()
		if err != nil {
			return -1, err
		}
		fmt.Println("   sync to: " + strconv.FormatInt(count, 10))
	}
	return count, nil
}
