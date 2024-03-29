package pfcregtest
//
//import (
//	"encoding/hex"
//	"github.com/picfight/pfcd/blockchain/stake"
//	"github.com/picfight/pfcd/chaincfg"
//	"github.com/picfight/pfcd/chaincfg/chainhash"
//	"github.com/picfight/pfcd/dcrjson"
//	"github.com/picfight/pfcd/dcrutil"
//	"github.com/picfight/pfcd/rpcclient"
//	"github.com/picfight/pfcd/txscript"
//	"github.com/decred/pfcwallet/wallet"
//	"github.com/google/go-cmp/cmp"
//	"github.com/jfixby/pin"
//	"math"
//	"math/big"
//	"reflect"
//	"strconv"
//	"strings"
//	"testing"
//	"time"
//)
//
//const defaultWalletPassphrase = "password"
//
//func TestGetNewAddress(t *testing.T) {
//
//	r := ObtainWalletHarness(mainWalletHarnessName)
//	// Wallet RPC client
//	wcl := r.Wallet
//
//	err := wcl.WalletUnlock(defaultWalletPassphrase, 0)
//	if err != nil {
//		t.Fatal("Failed to unlock wallet:", err)
//	}
//
//	// Get a new address from "default" account
//	// This is the first GetNewAddress call
//	// in this test for the "default" account
//	addr, err := wcl.GetNewAddress("default")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// Verify that address is for current network
//	if !addr.IsForNet(r.Node.Network()) {
//		t.Fatalf("Address not for active network (%s)", r.Node.Network())
//	}
//
//	// ValidateAddress
//	validRes, err := wcl.ValidateAddress(addr)
//	if err != nil {
//		t.Fatalf("Unable to validate address %s: %v", addr, err)
//	}
//	if !validRes.IsValid {
//		t.Fatalf("Address not valid: %s", addr)
//	}
//
//	// Create new account
//	accountName := "newAddressTest"
//	err = wcl.CreateNewAccount(accountName)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// Get a new address from new "newAddressTest" account
//	addrA, err := r.WalletRPCClient().GetNewAddress(accountName)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// Verify that address is for current network
//	if !addrA.IsForNet(r.Node.Network()) {
//		t.Fatalf("Address not for active network (%s)", r.Node.Network())
//	}
//
//	validRes, err = wcl.ValidateAddress(addrA)
//	if err != nil {
//		t.Fatalf("Unable to validate address %s: %v", addrA, err)
//	}
//	if !validRes.IsValid {
//		t.Fatalf("Address not valid: %s", addr)
//	}
//
//	// respect DefaultGapLimit
//	// -1 because of the first GetNewAddress("default") call above
//	for i := 0; i < wallet.DefaultGapLimit-1; i++ {
//		addr, err = wcl.GetNewAddress("default")
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		validRes, err = wcl.ValidateAddress(addr)
//		if err != nil {
//			t.Fatalf(
//				"Unable to validate address %s: %v",
//				addr,
//				err,
//			)
//		}
//		if !validRes.IsValid {
//			t.Fatalf("Address not valid: %s", addr)
//		}
//	}
//
//	// Expecting error:
//	// "policy violation: generating next address violates
//	// the unused address gap limit policy"
//	addr, err = wcl.GetNewAddress("default")
//	if err == nil {
//		t.Fatalf(
//			"Should report gap policy violation (%d)",
//			wallet.DefaultGapLimit,
//		)
//	}
//
//	// gap policy with wrapping
//	// reuse each address numOfReusages times
//	numOfReusages := 3
//	addrCounter := make(map[string]int)
//	for i := 0; i < wallet.DefaultGapLimit*numOfReusages; i++ {
//		addr, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetNewAddressGapPolicy(
//			"default", rpcclient.GapPolicyWrap)
//
//		// count address
//		num := addrCounter[addr.String()]
//		num++
//		addrCounter[addr.String()] = num
//
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		validRes, err := r.WalletRPCClient().Internal().(*rpcclient.Client).ValidateAddress(addr)
//		if err != nil {
//			t.Fatalf(
//				"Unable to validate address %s: %v",
//				addr,
//				err,
//			)
//		}
//		if !validRes.IsValid {
//			t.Fatalf("Address not valid: %s", addr)
//		}
//	}
//
//	// check reusages
//	for _, reused := range addrCounter {
//		if reused != numOfReusages {
//			t.Fatalf(
//				"Each address is expected to be reused: %d times, actual %d",
//				numOfReusages,
//				reused,
//			)
//		}
//	}
//
//	// ignore gap policy
//	for i := 0; i < wallet.DefaultGapLimit*2; i++ {
//		addr, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetNewAddressGapPolicy(
//			"default", rpcclient.GapPolicyIgnore)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		validRes, err := r.WalletRPCClient().Internal().(*rpcclient.Client).ValidateAddress(addr)
//		if err != nil {
//			t.Fatalf(
//				"Unable to validate address %s: %v",
//				addr,
//				err,
//			)
//		}
//		if !validRes.IsValid {
//			t.Fatalf("Address not valid: %s", addr)
//		}
//	}
//
//}
//
//func TestValidateAddress(t *testing.T) {
//
//	r := ObtainWalletHarness(mainWalletHarnessName)
//	// Wallet RPC client
//	wcl := r.Wallet
//
//	// Check that wallet is now unlocked
//	walletInfo, err := wcl.WalletInfo()
//	if err != nil {
//		t.Fatal("walletinfo failed.")
//	}
//	if !walletInfo.Unlocked {
//		t.Fatal("WalletPassphrase failed to unlock the wallet with the correct passphrase")
//	}
//
//	//-----------------------------------------
//	newAccountName := "testValidateAddress"
//	// Create a non-default account
//	err = wcl.CreateNewAccount(newAccountName)
//	if err != nil {
//		t.Fatalf("Unable to create account %s: %v", newAccountName, err)
//	}
//	accounts := []string{"default", newAccountName}
//	//-----------------------------------------
//	addrStr := "SsqvxBX8MZC5iiKCgBscwt69jg4u4hHhDKU"
//	// Try to validate an address that is not owned by wallet
//	otherAddress, err := dcrutil.DecodeAddress(addrStr)
//	if err != nil {
//		t.Fatalf("Unable to decode address %v: %v", otherAddress, err)
//	}
//	validRes, err := r.WalletRPCClient().Internal().(*rpcclient.Client).ValidateAddress(otherAddress)
//	if err != nil {
//		t.Fatalf("Unable to validate address %s with secondary wallet: %v",
//			addrStr, err)
//	}
//	if !validRes.IsValid {
//		t.Fatalf("Address not valid: %s", addrStr)
//	}
//	if validRes.IsMine {
//		t.Fatalf("Address incorrectly identified as mine: %s", addrStr)
//	}
//	if validRes.IsScript {
//		t.Fatalf("Address incorrectly identified as script: %s", addrStr)
//	}
//	//-----------------------------------------
//	// Validate simnet dev subsidy address
//	devSubPkScript := chaincfg.SimNetParams.OrganizationPkScript // "ScuQxvveKGfpG1ypt6u27F99Anf7EW3cqhq"
//	devSubPkScrVer := chaincfg.SimNetParams.OrganizationPkScriptVersion
//	_, addrs, _, err := txscript.ExtractPkScriptAddrs(
//		devSubPkScrVer, devSubPkScript, r.Node.Network().Params().(*chaincfg.Params))
//	if err != nil {
//		t.Fatal("Failed to extract addresses from PkScript:", err)
//	}
//	devSubAddrStr := addrs[0].String()
//
//	DevAddr, err := dcrutil.DecodeAddress(devSubAddrStr)
//	if err != nil {
//		t.Fatalf("Unable to decode address %s: %v", devSubAddrStr, err)
//	}
//
//	validRes, err = r.WalletRPCClient().Internal().(*rpcclient.Client).ValidateAddress(DevAddr)
//	if err != nil {
//		t.Fatalf("Unable to validate address %s: ", devSubAddrStr)
//	}
//	if !validRes.IsValid {
//		t.Fatalf("Address not valid: %s", devSubAddrStr)
//	}
//	if validRes.IsMine {
//		t.Fatalf("Address incorrectly identified as mine: %s", devSubAddrStr)
//	}
//	// final address overflow check for each account
//	for _, acct := range accounts {
//		// let's overflow DefaultGapLimit
//		for i := 0; i < wallet.DefaultGapLimit+5; i++ {
//			// Get a new address from current account
//			addr, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetNewAddressGapPolicy(
//				acct, rpcclient.GapPolicyIgnore)
//			if err != nil {
//				t.Fatal(err)
//			}
//			// Verify that address is for current network
//			if !addr.IsForNet(r.Node.Network().Params().(*chaincfg.Params)) {
//				t.Fatalf(
//					"Address[%d] not for active network (%s), <%s>",
//					i,
//					r.Node.Network(),
//					acct,
//				)
//			}
//			// ValidateAddress
//			addrStr := addr.String()
//			validRes, err := r.WalletRPCClient().Internal().(*rpcclient.Client).ValidateAddress(addr)
//			if err != nil {
//				t.Fatalf(
//					"Unable to validate address[%d] %s: %v for <%s>",
//					i,
//					addrStr,
//					err,
//					acct,
//				)
//			}
//			if !validRes.IsValid {
//				t.Fatalf(
//					"Address[%d] not valid: %s for <%s>",
//					i,
//					addrStr,
//					acct,
//				)
//			}
//			if !validRes.IsMine {
//				t.Fatalf(
//					"Address[%d] incorrectly identified as NOT mine: %s for <%s>",
//					i,
//					addrStr,
//					acct,
//				)
//			}
//			if validRes.IsScript {
//				t.Fatalf(
//					"Address[%d] incorrectly identified as script: %s for <%s>",
//					i,
//					addrStr,
//					acct,
//				)
//			}
//			// Address is "mine", so we can check account
//			if strings.Compare(acct, validRes.Account) != 0 {
//				t.Fatalf("Address[%d] %s reported as not from <%s> account",
//					i,
//					addrStr,
//					acct,
//				)
//			}
//			// Decode address
//			_, err = dcrutil.DecodeAddress(addrStr)
//			if err != nil {
//				t.Fatalf("Unable to decode address[%d] %s: %v for <%s>",
//					i,
//					addr.String(),
//					err,
//					acct,
//				)
//			}
//		}
//
//	}
//
//}
//
//func TestGetBalance(t *testing.T) {
//	r := ObtainWalletHarness(mainWalletHarnessName)
//	wcl := r.Wallet
//
//	list, err := wcl.ListAccounts()
//	if err != nil {
//		t.Fatalf("ListAccounts failed: %v", err)
//	}
//	pin.D("ListAccounts", list)
//
//	err = wcl.WalletUnlock(defaultWalletPassphrase, 0)
//	if err != nil {
//		t.Fatal("Failed to unlock wallet:", err)
//	}
//
//	_, err = wcl.GetBalance()
//	if err != nil {
//		t.Fatalf("GetBalance failed: %v", err)
//	}
//	//pin.D("balance", balance)
//	//pin.S("balance", balance)
//
//	accountName := "getBalanceTest"
//	err = wcl.CreateNewAccount(accountName)
//	if err != nil {
//		t.Fatalf("CreateNewAccount failed: %v", err)
//	}
//
//	// Grab a fresh address from the test account
//	addr, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetNewAddressGapPolicy(
//		accountName,
//		rpcclient.GapPolicyWrap,
//	)
//	if err != nil {
//		t.Fatalf("GetNewAddress failed: %v", err)
//	}
//
//	// Check invalid account name
//	_, err = r.WalletRPCClient().Internal().(*rpcclient.Client).GetBalanceMinConf("invalid account", 0)
//	// -4: account name 'invalid account' not found
//	if err == nil {
//		t.Fatalf("GetBalanceMinConfType failed to return non-nil error for invalid account name: %v", err)
//	}
//
//	// Check invalid minconf
//	_, err = r.WalletRPCClient().Internal().(*rpcclient.Client).GetBalanceMinConf("default", -1)
//	if err == nil {
//		t.Fatalf("GetBalanceMinConf failed to return non-nil error for invalid minconf (-1)")
//	}
//
//	preBalances, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetBalanceMinConf("*", 0)
//	if err != nil {
//		t.Fatalf("GetBalanceMinConf(\"*\", 0) failed: %v", err)
//	}
//
//	preAccountBalanceSpendable := 0.0
//	preAccountBalances := make(map[string]dcrjson.GetAccountBalanceResult)
//	for _, bal := range preBalances.Balances {
//		preAccountBalanceSpendable += bal.Spendable
//		preAccountBalances[bal.AccountName] = bal
//	}
//
//	// Send from default to test account
//	sendAmount := dcrutil.Amount(700000000)
//	if _, err = r.WalletRPCClient().Internal().(*rpcclient.Client).SendFromMinConf("default", addr, sendAmount, 1); err != nil {
//		t.Fatalf("SendFromMinConf failed: %v", err)
//	}
//
//	// Check invalid minconf
//	postBalances, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetBalanceMinConf("*", 0)
//	if err != nil {
//		t.Fatalf("GetBalanceMinConf failed: %v", err)
//	}
//
//	postAccountBalanceSpendable := 0.0
//	postAccountBalances := make(map[string]dcrjson.GetAccountBalanceResult)
//	for _, bal := range postBalances.Balances {
//		postAccountBalanceSpendable += bal.Spendable
//		postAccountBalances[bal.AccountName] = bal
//	}
//
//	// Fees prevent easy exact comparison
//	if preAccountBalances["default"].Spendable <= postAccountBalances["default"].Spendable {
//		t.Fatalf("spendable balance of account 'default' not decreased: %v <= %v",
//			preAccountBalances["default"].Spendable,
//			postAccountBalances["default"].Spendable)
//	}
//
//	if sendAmount.ToCoin() != (postAccountBalances[accountName].Spendable - preAccountBalances[accountName].Spendable) {
//		t.Fatalf("spendable balance of account '%s' not increased: %v >= %v",
//			accountName,
//			preAccountBalances[accountName].Spendable,
//			postAccountBalances[accountName].Spendable)
//	}
//
//	// Make sure "*" account balance has decreased (fees)
//	if postAccountBalanceSpendable >= preAccountBalanceSpendable {
//		t.Fatalf("Total balance over all accounts not decreased after send.")
//	}
//
//	// Test vanilla GetBalance()
//	amtGetBalance, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetBalance("default")
//	if err != nil {
//		t.Fatalf("GetBalance failed: %v", err)
//	}
//
//	// For GetBalance(), default minconf=1.
//	defaultBalanceMinConf1, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetBalanceMinConf("default", 1)
//	if err != nil {
//		t.Fatalf("GetBalanceMinConfType failed: %v", err)
//	}
//
//	if amtGetBalance.Balances[0].Spendable != defaultBalanceMinConf1.Balances[0].Spendable {
//		t.Fatalf(`Balance from GetBalance("default") does not equal amount `+
//			`from GetBalanceMinConf: %v != %v`,
//			amtGetBalance.Balances[0].Spendable,
//			defaultBalanceMinConf1.Balances[0].Spendable)
//	}
//
//	// Verify minconf=1 balances of receiving account before/after new block
//	// Before, getbalance minconf=1
//	amtTestMinconf1BeforeBlock, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetBalanceMinConf(accountName, 1)
//	if err != nil {
//		t.Fatalf("GetBalanceMinConf failed: %v", err)
//	}
//
//	// Mine 2 new blocks to validate tx
//	newBestBlock(r, t)
//	newBestBlock(r, t)
//
//	// After, getbalance minconf=1
//	amtTestMinconf1AfterBlock, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetBalanceMinConf(accountName, 1)
//	if err != nil {
//		t.Fatalf("GetBalanceMinConf failed: %v", err)
//	}
//
//	// Verify that balance (minconf=1) has increased
//	if sendAmount.ToCoin() != (amtTestMinconf1AfterBlock.Balances[0].Spendable - amtTestMinconf1BeforeBlock.Balances[0].Spendable) {
//		t.Fatalf(`Balance (minconf=1) not increased after new block: %v - %v != %v`,
//			amtTestMinconf1AfterBlock.Balances[0].Spendable,
//			amtTestMinconf1BeforeBlock.Balances[0].Spendable,
//			sendAmount)
//	}
//}
//
//func TestListAccounts(t *testing.T) {
//
//	r := ObtainWalletHarness(mainWalletHarnessName)
//	// Wallet RPC client
//	wcl := r.Wallet
//
//	// Create a new account and verify that we can see it
//	listBeforeCreateAccount, err := r.WalletRPCClient().Internal().(*rpcclient.Client).ListAccounts()
//	if err != nil {
//		t.Fatal("Failed to create new account ", err)
//	}
//
//	// New account
//	accountName := "listaccountsTestAcct"
//	err = wcl.CreateNewAccount(accountName)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// Account list after creating new
//	accountsBalancesDefault1, err := r.WalletRPCClient().Internal().(*rpcclient.Client).ListAccounts()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// Verify that new account is in the list, with zero balance
//	foundNewAcct := false
//	for acct, amt := range accountsBalancesDefault1 {
//		if _, ok := listBeforeCreateAccount[acct]; !ok {
//			// Found new account.  Now check name and balance
//			if amt != 0 {
//				t.Fatalf("New account (%v) found with non-zero balance: %v",
//					acct, amt)
//			}
//			if accountName == acct {
//				foundNewAcct = true
//				break
//			}
//			t.Fatalf("Found new account, %v; Expected %v", acct, accountName)
//		}
//	}
//	if !foundNewAcct {
//		t.Fatalf("Failed to find newly created account, %v.", accountName)
//	}
//
//	// Grab a fresh address from the test account
//	addr, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetNewAddressGapPolicy(
//		accountName,
//		rpcclient.GapPolicyWrap,
//	)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// For ListAccountsCmd: MinConf *int `jsonrpcdefault:"1"`
//	// Let's test that ListAccounts() is equivalent to explicit minconf=1
//	accountsBalancesMinconf1, err := r.WalletRPCClient().Internal().(*rpcclient.Client).ListAccountsMinConf(1)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if !reflect.DeepEqual(accountsBalancesDefault1, accountsBalancesMinconf1) {
//		t.Fatal("ListAccounts() returned different result from ListAccountsMinConf(1): ",
//			accountsBalancesDefault1, accountsBalancesMinconf1)
//	}
//
//	// Get accounts with minconf=0 pre-send
//	accountsBalancesMinconf0PreSend, err := r.WalletRPCClient().Internal().(*rpcclient.Client).ListAccountsMinConf(0)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// Get balance of test account prior to a send
//	acctBalancePreSend := accountsBalancesMinconf0PreSend[accountName]
//
//	// Send from default to test account
//	sendAmount := dcrutil.Amount(700000000)
//	if _, err = r.WalletRPCClient().Internal().(*rpcclient.Client).SendFromMinConf("default", addr, sendAmount, 1); err != nil {
//		t.Fatal("SendFromMinConf failed.", err)
//	}
//
//	// Get accounts with minconf=0 post-send
//	accountsBalancesMinconf0PostSend, err := r.WalletRPCClient().Internal().(*rpcclient.Client).ListAccountsMinConf(0)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// Get balance of test account prior to a send
//	acctBalancePostSend := accountsBalancesMinconf0PostSend[accountName]
//
//	// Check if reported balances match expectations
//	if sendAmount != (acctBalancePostSend - acctBalancePreSend) {
//		t.Fatalf("Test account balance not changed by expected amount after send: "+
//			"%v -%v != %v", acctBalancePostSend, acctBalancePreSend, sendAmount)
//	}
//
//	// Verify minconf>0 works: list, mine, list
//
//	// List BEFORE mining a block
//	accountsBalancesMinconf1PostSend, err := r.WalletRPCClient().Internal().(*rpcclient.Client).ListAccountsMinConf(1)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// Get balance of test account prior to a send
//	acctBalanceMin1PostSend := accountsBalancesMinconf1PostSend[accountName]
//
//	// Mine 2 new blocks to validate tx
//	newBestBlock(r, t)
//	newBestBlock(r, t)
//
//	// List AFTER mining a block
//	accountsBalancesMinconf1PostMine, err := r.WalletRPCClient().Internal().(*rpcclient.Client).ListAccountsMinConf(1)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// Get balance of test account prior to a send
//	acctBalanceMin1PostMine := accountsBalancesMinconf1PostMine[accountName]
//
//	// Check if reported balances match expectations
//	if sendAmount != (acctBalanceMin1PostMine - acctBalanceMin1PostSend) {
//		t.Fatalf("Test account balance (minconf=1) not changed by expected "+
//			"amount after new block: %v - %v != %v", acctBalanceMin1PostMine,
//			acctBalanceMin1PostSend, sendAmount)
//	}
//
//	// Note that ListAccounts uses Store.balanceFullScan to handle a UTXO scan
//	// for each specific account. We can compare against GetBalanceMinConfType.
//	// Also, I think there is the same bug that allows negative minconf values,
//	// but does not handle unconfirmed outputs the same way as minconf=0.
//
//	GetBalancePostSend, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetBalanceMinConf(accountName, 0)
//	if err != nil {
//		t.Fatal(err)
//	}
//	// Note that BFBalanceSpendable is used with GetBalanceMinConf (not Type),
//	// which uses BFBalanceFullScan when a single account is specified.
//	// Recall thet fullscan is used by listaccounts.
//
//	if GetBalancePostSend.Balances[0].Spendable != acctBalancePostSend.ToCoin() {
//		t.Fatalf("Balance for default account from GetBalanceMinConf does not "+
//			"match balance from ListAccounts: %v != %v", GetBalancePostSend,
//			acctBalancePostSend)
//	}
//
//	// Mine 2 blocks to validate the tx and clean up UTXO set
//	newBestBlock(r, t)
//	newBestBlock(r, t)
//}
//
//func TestListUnspent(t *testing.T) {
//
//	r := ObtainWalletHarness(mainWalletHarnessName)
//	// Wallet RPC client
//	wcl := r.Wallet
//
//	// New account
//	accountName := "listUnspentTestAcct"
//	err := wcl.CreateNewAccount(accountName)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// Grab an address from the test account
//	addr, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetNewAddressGapPolicy(
//		accountName,
//		rpcclient.GapPolicyWrap,
//	)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// UTXOs before send
//	list, err := r.WalletRPCClient().Internal().(*rpcclient.Client).ListUnspent()
//	if err != nil {
//		t.Fatalf("failed to get utxos")
//	}
//	utxosBeforeSend := make(map[string]float64)
//	for _, utxo := range list {
//		// Get a OutPoint string in the form of hash:index
//		outpointStr, err := getOutPointString(&utxo)
//		if err != nil {
//			t.Fatal(err)
//		}
//		// if utxo.Spendable ...
//		utxosBeforeSend[outpointStr] = utxo.Amount
//	}
//
//	// Check Min/Maxconf arguments
//	defaultMaxConf := 9999999
//
//	listMin1MaxBig, err := r.WalletRPCClient().Internal().(*rpcclient.Client).ListUnspentMinMax(1, defaultMaxConf)
//	if err != nil {
//		t.Fatalf("failed to get utxos")
//	}
//	if !reflect.DeepEqual(list, listMin1MaxBig) {
//		t.Fatal("Outputs from ListUnspent() and ListUnspentMinMax() do not match.")
//	}
//
//	// Grab an address from known unspents to test the filter
//	refOut := list[0]
//	PkScript, err := hex.DecodeString(refOut.ScriptPubKey)
//	if err != nil {
//		t.Fatalf("Failed to decode ScriptPubKey into PkScript.")
//	}
//	// The Address field is broken, including only one address, so don't use it
//	_, addrs, _, err := txscript.ExtractPkScriptAddrs(
//		txscript.DefaultScriptVersion, PkScript, r.Node.Network().Params().(*chaincfg.Params))
//	if err != nil {
//		t.Fatal("Failed to extract addresses from PkScript:", err)
//	}
//
//	// List with all of the above address
//	listAddressesKnown, err := r.WalletRPCClient().Internal().(*rpcclient.Client).ListUnspentMinMaxAddresses(1, defaultMaxConf, addrs)
//	if err != nil {
//		t.Fatalf("Failed to get utxos with addresses argument.")
//	}
//
//	// Check that there is at least one output for the input addresses
//	if len(listAddressesKnown) == 0 {
//		t.Fatalf("Failed to find expected UTXOs with addresses.")
//	}
//
//	// Make sure each found output's txid:vout is in original list
//	var foundTxID = false
//	for _, listRes := range listAddressesKnown {
//		// Get a OutPoint string in the form of hash:index
//		outpointStr, err := getOutPointString(&listRes)
//		if err != nil {
//			t.Fatal(err)
//		}
//		if _, ok := utxosBeforeSend[outpointStr]; !ok {
//			t.Fatalf("Failed to find TxID")
//		}
//		// Also verify that the txid of the reference output is in the list
//		if listRes.TxID == refOut.TxID {
//			foundTxID = true
//		}
//	}
//	if !foundTxID {
//		t.Fatal("Original TxID not found in list by addresses.")
//	}
//
//	// SendFromMinConf to addr
//	amountToSend := dcrutil.Amount(700000000)
//	txid, err := r.WalletRPCClient().Internal().(*rpcclient.Client).SendFromMinConf("default", addr, amountToSend, 0)
//	if err != nil {
//		t.Fatalf("sendfromminconf failed: %v", err)
//	}
//
//	newBestBlock(r, t)
//	time.Sleep(1 * time.Second)
//	// New block is necessary for GetRawTransaction to give a tx with sensible
//	// MsgTx().TxIn[:].ValueIn values.
//
//	// Get *dcrutil.Tx of send to check the inputs
//	rawTx, err := r.NodeRPCClient().Internal().(*rpcclient.Client).GetRawTransaction(txid)
//	if err != nil {
//		t.Fatalf("getrawtransaction failed: %v", err)
//	}
//
//	// Get previous OutPoint of each TxIn for send transaction
//	txInIDs := make(map[string]float64)
//	for _, txIn := range rawTx.MsgTx().TxIn {
//		prevOut := &txIn.PreviousOutPoint
//		// Outpoint.String() appends :index to the hash
//		txInIDs[prevOut.String()] = dcrutil.Amount(txIn.ValueIn).ToCoin()
//	}
//
//	// First check to make sure we see these in the UTXO list prior to send,
//	// then not in the UTXO list after send.
//	for txinID, amt := range txInIDs {
//		if _, ok := utxosBeforeSend[txinID]; !ok {
//			t.Fatalf("Failed to find txid %v (%v PFC) in list of UTXOs",
//				txinID, amt)
//		}
//	}
//
//	// Validate the send Tx with 2 new blocks
//	newBestBlock(r, t)
//	newBestBlock(r, t)
//
//	// Make sure these txInIDS are not in the new UTXO set
//	time.Sleep(2 * time.Second)
//	list, err = r.WalletRPCClient().Internal().(*rpcclient.Client).ListUnspent()
//	if err != nil {
//		t.Fatalf("Failed to get UTXOs")
//	}
//	for _, utxo := range list {
//		// Get a OutPoint string in the form of hash:index
//		outpointStr, err := getOutPointString(&utxo)
//		if err != nil {
//			t.Fatal(err)
//		}
//		if amt, ok := txInIDs[outpointStr]; ok {
//			t.Fatalf("Found PreviousOutPoint of send still in UTXO set: %v, "+
//				"%v PFC", outpointStr, amt)
//		}
//	}
//}
//
//func TestSendToAddress(t *testing.T) {
//
//	r := ObtainWalletHarness(mainWalletHarnessName)
//
//	// Grab a fresh address from the wallet.
//	addr, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetNewAddressGapPolicy(
//		"default", rpcclient.GapPolicyIgnore)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// Check balance of default account
//	_, err = r.WalletRPCClient().Internal().(*rpcclient.Client).GetBalanceMinConf("default", 1)
//	if err != nil {
//		t.Fatalf("GetBalanceMinConfType failed: %v", err)
//	}
//
//	// SendToAddress
//	txid, err := r.WalletRPCClient().Internal().(*rpcclient.Client).SendToAddress(addr, 1000000)
//	if err != nil {
//		t.Fatalf("SendToAddress failed: %v", err)
//	}
//
//	// Generate a single block, in which the transaction the wallet created
//	// should be found.
//	_, block, _ := newBestBlock(r, t)
//
//	if len(block.Transactions()) <= 1 {
//		t.Fatalf("expected transaction not included in block")
//	}
//	// Confirm that the expected tx was mined into the block.
//	minedTx := block.Transactions()[1]
//	txHash := minedTx.Hash()
//	if *txid != *txHash {
//		t.Fatalf("txid's don't match, %v vs %v", txHash, txid)
//	}
//
//	// We should now check to confirm that the utxo that wallet used to create
//	// that sendfrom was properly marked as spent and removed from utxo set. Use
//	// GetTxOut to tell if the outpoint is spent.
//	//
//	// The spending transaction has to be off the tip block for the previous
//	// outpoint to be spent, out of the UTXO set. Generate another block.
//	_, err = GenerateBlock(r, block.MsgBlock().Header.Height)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// Check each PreviousOutPoint for the sending tx.
//	time.Sleep(1 * time.Second)
//	// Get the sending Tx
//	rawTx, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetRawTransaction(txid)
//	if err != nil {
//		t.Fatalf("Unable to get raw transaction %v: %v", txid, err)
//	}
//	// txid is rawTx.MsgTx().TxIn[0].PreviousOutPoint.Hash
//
//	// Check all inputs
//	for i, txIn := range rawTx.MsgTx().TxIn {
//		prevOut := &txIn.PreviousOutPoint
//
//		// If a txout is spent (not in the UTXO set) GetTxOutResult will be nil
//		res, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetTxOut(&prevOut.Hash, prevOut.Index, false)
//		if err != nil {
//			t.Fatal("GetTxOut failure:", err)
//		}
//		if res != nil {
//			t.Fatalf("Transaction output %v still unspent.", i)
//		}
//	}
//}
//
//func TestSendFrom(t *testing.T) {
//	r := ObtainWalletHarness(mainWalletHarnessName)
//
//	err := r.Wallet.WalletUnlock(defaultWalletPassphrase, 0)
//	if err != nil {
//		t.Fatal("Failed to unlock wallet:", err)
//	}
//
//	accountName := "sendFromTest"
//	err = r.WalletRPCClient().Internal().(*rpcclient.Client).CreateNewAccount(accountName)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// Grab a fresh address from the wallet.
//	addr, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetNewAddressGapPolicy(
//		accountName,
//		rpcclient.GapPolicyWrap,
//	)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	amountToSend := dcrutil.Amount(1000000)
//	// Check spendable balance of default account
//	defaultBalanceBeforeSend, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetBalanceMinConf("default", 0)
//	if err != nil {
//		t.Fatalf("GetBalanceMinConf failed: %v", err)
//	}
//
//	// Get utxo list before send
//	list, err := r.WalletRPCClient().Internal().(*rpcclient.Client).ListUnspent()
//	if err != nil {
//		t.Fatalf("failed to get utxos")
//	}
//	utxosBeforeSend := make(map[string]float64)
//	for _, utxo := range list {
//		// Get a OutPoint string in the form of hash:index
//		outpointStr, err := getOutPointString(&utxo)
//		if err != nil {
//			t.Fatal(err)
//		}
//		// if utxo.Spendable ...
//		utxosBeforeSend[outpointStr] = utxo.Amount
//	}
//
//	// SendFromMinConf 1000 to addr
//	txid, err := r.WalletRPCClient().Internal().(*rpcclient.Client).SendFromMinConf("default", addr, amountToSend, 0)
//	if err != nil {
//		t.Fatalf("sendfromminconf failed: %v", err)
//	}
//
//	// Check spendable balance of default account
//	defaultBalanceAfterSendNoBlock, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetBalanceMinConf("default", 0)
//	if err != nil {
//		t.Fatalf("GetBalanceMinConf failed: %v", err)
//	}
//
//	// Check balance of sendfrom account
//	sendFromBalanceAfterSendNoBlock, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetBalanceMinConf(accountName, 0)
//	if err != nil {
//		t.Fatalf("GetBalanceMinConf failed: %v", err)
//	}
//	if sendFromBalanceAfterSendNoBlock.Balances[0].Spendable != amountToSend.ToCoin() {
//		t.Fatalf("balance for %s account incorrect:  want %v got %v",
//			accountName, amountToSend, sendFromBalanceAfterSendNoBlock.Balances[0].Spendable)
//	}
//
//	// Generate a single block, the transaction the wallet created should
//	// be found in this block.
//	_, block, _ := newBestBlock(r, t)
//
//	// Check to make sure the transaction that was sent was included in the block
//	if len(block.Transactions()) <= 1 {
//		t.Fatalf("expected transaction not included in block")
//	}
//	minedTx := block.Transactions()[1]
//	txHash := minedTx.Hash()
//	if *txid != *txHash {
//		t.Fatalf("txid's don't match, %v vs. %v (actual vs. expected)",
//			txHash, txid)
//	}
//
//	// Generate another block, since it takes 2 blocks to validate a tx
//	newBestBlock(r, t)
//
//	// Get rawTx of sent txid so we can calculate the fee that was used
//	time.Sleep(1 * time.Second)
//	rawTx, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetRawTransaction(txid)
//	if err != nil {
//		t.Fatalf("getrawtransaction failed: %v", err)
//	}
//
//	var totalSpent int64
//	for _, txIn := range rawTx.MsgTx().TxIn {
//		totalSpent += txIn.ValueIn
//	}
//
//	var totalSent int64
//	for _, txOut := range rawTx.MsgTx().TxOut {
//		totalSent += txOut.Value
//	}
//
//	feeAtoms := dcrutil.Amount(totalSpent - totalSent)
//
//	// Calculate the expected balance for the default account after the tx was sent
//	sentAtoms := uint64(amountToSend + feeAtoms)
//
//	m1 := new(big.Float)
//	m1.SetFloat64(-1)
//
//	E8 := new(big.Float)
//	E8.SetFloat64(math.Pow10(int(8)))
//
//	sentAtomsNegative := new(big.Float)
//	sentAtomsNegative.SetUint64(sentAtoms)
//	sentAtomsNegative = sentAtomsNegative.Mul(sentAtomsNegative, m1)
//
//	oldBalanceCoins := new(big.Float)
//	oldBalanceCoins.SetFloat64(defaultBalanceBeforeSend.Balances[0].Spendable)
//	oldBalanceAtoms := new(big.Float)
//	oldBalanceAtoms = oldBalanceAtoms.Mul(oldBalanceCoins, E8)
//
//	expectedBalanceAtoms := new(big.Float)
//	expectedBalanceAtoms.Add(oldBalanceAtoms, sentAtomsNegative)
//
//	currentBalanceCoinsNegative := new(big.Float)
//	currentBalanceCoinsNegative.SetFloat64(defaultBalanceAfterSendNoBlock.Balances[0].Spendable)
//	currentBalanceCoinsNegative = currentBalanceCoinsNegative.Mul(currentBalanceCoinsNegative, m1)
//
//	currentBalanceAtomsNegative := new(big.Float)
//	currentBalanceAtomsNegative = currentBalanceAtomsNegative.Mul(currentBalanceCoinsNegative, E8)
//
//	diff := new(big.Float)
//	diff.Add(currentBalanceAtomsNegative, expectedBalanceAtoms)
//
//	zero := new(big.Float)
//	zero.SetFloat64(0)
//
//	if diff.Cmp(zero) != 0 {
//		t.Fatalf("balance for %s account incorrect: want %v got %v, diff %V",
//			"default",
//			currentBalanceAtomsNegative,
//			defaultBalanceAfterSendNoBlock.Balances[0].Spendable,
//			diff,
//		)
//	}
//
//	// Check balance of sendfrom account
//	sendFromBalanceAfterSend1Block, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetBalanceMinConf(accountName, 1)
//	if err != nil {
//		t.Fatalf("getbalanceminconftype failed: %v", err)
//	}
//
//	if sendFromBalanceAfterSend1Block.Balances[0].Total != amountToSend.ToCoin() {
//		t.Fatalf("balance for %s account incorrect:  want %v got %v",
//			accountName, amountToSend, sendFromBalanceAfterSend1Block.Balances[0].Total)
//	}
//
//	// We have confirmed that the expected tx was mined into the block.
//	// We should now check to confirm that the utxo that wallet used to create
//	// that sendfrom was properly marked to spent and removed from utxo set.
//
//	// Get the sending Tx
//	rawTx, err = r.WalletRPCClient().Internal().(*rpcclient.Client).GetRawTransaction(txid)
//	if err != nil {
//		t.Fatalf("Unable to get raw transaction %v: %v", txid, err)
//	}
//
//	// Check all inputs
//	for i, txIn := range rawTx.MsgTx().TxIn {
//		prevOut := &txIn.PreviousOutPoint
//
//		// If a txout is spent (not in the UTXO set) GetTxOutResult will be nil
//		res, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetTxOut(&prevOut.Hash, prevOut.Index, false)
//		if err != nil {
//			t.Fatal("GetTxOut failure:", err)
//		}
//		if res != nil {
//			t.Fatalf("Transaction output %v still unspent.", i)
//		}
//	}
//}
//
//func TestSendMany(t *testing.T) {
//
//	r := ObtainWalletHarness(mainWalletHarnessName)
//	// Wallet RPC client
//	wcl := r.Wallet
//
//	err := wcl.WalletUnlock(defaultWalletPassphrase, 0)
//	if err != nil {
//		t.Fatal("Failed to unlock wallet:", err)
//	}
//
//	// Create 2 accounts to receive funds
//	accountNames := []string{"sendManyTestA", "sendManyTestB"}
//	amountsToSend := []dcrutil.Amount{700000000, 1400000000}
//	addresses := []dcrutil.Address{}
//
//	for _, acct := range accountNames {
//		err = wcl.CreateNewAccount(acct)
//		if err != nil {
//			t.Fatal(err)
//		}
//	}
//
//	// Grab new addresses from the wallet, under each account.
//	// Set corresponding amount to send to each address.
//	addressAmounts := make(map[dcrutil.Address]dcrutil.Amount)
//	totalAmountToSend := dcrutil.Amount(0)
//
//	for i, acct := range accountNames {
//		addr, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetNewAddressGapPolicy(
//			acct,
//			rpcclient.GapPolicyWrap,
//		)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		// Set the amounts to send to each address
//		addresses = append(addresses, addr)
//		addressAmounts[addr] = amountsToSend[i]
//		totalAmountToSend += amountsToSend[i]
//	}
//
//	// Check spendable balance of default account
//	defaultBalanceBeforeSend, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetBalanceMinConf("default", 0)
//	if err != nil {
//		t.Fatalf("GetBalanceMinConf default failed: %v", err)
//	}
//
//	// SendMany to two addresses
//	txid, err := r.WalletRPCClient().Internal().(*rpcclient.Client).SendMany("default", addressAmounts)
//	if err != nil {
//		t.Fatalf("SendMany failed: %v", err)
//	}
//
//	// XXX
//	time.Sleep(250 * time.Millisecond)
//
//	// Check spendable balance of default account
//	defaultBalanceAfterSendUnmined, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetBalanceMinConf("default", 0)
//	if err != nil {
//		t.Fatalf("GetBalanceMinConf failed: %v", err)
//	}
//
//	// Check balance of each receiving account
//	for i, acct := range accountNames {
//		bal, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetBalanceMinConf(acct, 0)
//		if err != nil {
//			t.Fatalf("GetBalanceMinConf '%s' failed: %v", acct, err)
//		}
//		addr := addresses[i]
//		if bal.Balances[0].Total != addressAmounts[addr].ToCoin() {
//			t.Fatalf("Balance for %s account incorrect:  want %v got %v",
//				acct, addressAmounts[addr], bal)
//		}
//	}
//
//	// Get rawTx of sent txid so we can calculate the fee that was used
//	rawTx, err := r.NodeRPCClient().Internal().(*rpcclient.Client).GetRawTransaction(txid)
//	if err != nil {
//		t.Fatalf("getrawtransaction failed: %v", err)
//	}
//	fee := getWireMsgTxFee(rawTx)
//
//	// Generate a single block, the transaction the wallet created should be
//	// found in this block.
//	_, block, _ := newBestBlock(r, t)
//
//	rawTx, err = r.NodeRPCClient().Internal().(*rpcclient.Client).GetRawTransaction(txid)
//	if err != nil {
//		t.Fatalf("getrawtransaction failed: %v", err)
//	}
//	fee = getWireMsgTxFee(rawTx)
//
//	// Calculate the expected balance for the default account after the tx was sent
//
//	sentAtoms := totalAmountToSend + fee
//	sentCoinsFloat := sentAtoms.ToCoin()
//
//	sentCoinsNegative := new(big.Float)
//	sentCoinsNegative.SetFloat64(-sentCoinsFloat)
//
//	oldBalanceCoins := new(big.Float)
//	oldBalanceCoins.SetFloat64(defaultBalanceBeforeSend.Balances[0].Spendable)
//
//	expectedBalanceCoins := new(big.Float)
//	expectedBalanceCoins.Add(oldBalanceCoins, sentCoinsNegative)
//
//	currentBalanceCoinsNegative := new(big.Float)
//	currentBalanceCoinsNegative.SetFloat64(defaultBalanceAfterSendUnmined.Balances[0].Spendable)
//
//	f64A, _ := currentBalanceCoinsNegative.Float64()
//	f64B, _ := expectedBalanceCoins.Float64()
//
//	if f64A != f64B {
//		t.Fatalf("Balance for %s account (sender) incorrect: want %v got %v",
//			"default",
//			f64B,
//			f64A,
//		)
//	}
//
//	// Check to make sure the transaction that was sent was included in the block
//	if !includesTx(txid, block) {
//		t.Fatalf("Expected transaction not included in block")
//	}
//
//	// Validate
//	newBestBlock(r, t)
//
//	// Check balance after confirmations
//	for i, acct := range accountNames {
//		balanceAcctValidated, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetBalanceMinConf(acct, 1)
//		if err != nil {
//			t.Fatalf("GetBalanceMinConf '%s' failed: %v", acct, err)
//		}
//
//		addr := addresses[i]
//		if balanceAcctValidated.Balances[0].Total != addressAmounts[addr].ToCoin() {
//			t.Fatalf("Balance for %s account incorrect:  want %v got %v",
//				acct, addressAmounts[addr].ToCoin(), balanceAcctValidated.Balances[0].Total)
//		}
//	}
//
//	// Check all inputs
//	for i, txIn := range rawTx.MsgTx().TxIn {
//		prevOut := &txIn.PreviousOutPoint
//
//		// If a txout is spent (not in the UTXO set) GetTxOutResult will be nil
//		res, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetTxOut(&prevOut.Hash, prevOut.Index, false)
//		if err != nil {
//			t.Fatal("GetTxOut failure:", err)
//		}
//		if res != nil {
//			t.Fatalf("Transaction output %v still unspent.", i)
//		}
//	}
//}
//
//func TestListTransactions(t *testing.T) {
//	r := ObtainWalletHarness(t.Name())
//	// Wallet RPC client
//	wcl := r.Wallet
//
//	err := wcl.WalletUnlock(defaultWalletPassphrase, 0)
//	if err != nil {
//		t.Fatal("Failed to unlock wallet:", err)
//	}
//
//	// List latest transaction
//	txList1, err := r.WalletRPCClient().Internal().(*rpcclient.Client).ListTransactionsCount("*", 1)
//	if err != nil {
//		t.Fatal("ListTransactionsCount failed:", err)
//	}
//
//	// Verify that only one returned (a PoW coinbase since this is a fresh
//	// harness with only blocks generated and no other transactions).
//	if len(txList1) != 1 {
//		t.Fatalf("Transaction list not len=1: %d", len(txList1))
//	}
//
//	// Verify paid to MiningAddress
//	if txList1[0].Address != r.MiningAddress.String() {
//		t.Fatalf("Unexpected address in latest transaction: %v",
//			txList1[0].Address)
//	}
//
//	// Verify that it is a coinbase
//	if !txList1[0].Generated {
//		t.Fatal("Latest transaction output not a coinbase output.")
//	}
//
//	// Not "generate" category until mature
//	if txList1[0].Category != "immature" {
//		t.Fatalf("Latest transaction not immature. Category: %v",
//			txList1[0].Category)
//	}
//
//	// Verify blockhash is non-nil and valid
//	hash, err := chainhash.NewHashFromStr(txList1[0].BlockHash)
//	if err != nil {
//		t.Fatal("Blockhash not valid")
//	}
//	_, err = r.WalletRPCClient().Internal().(*rpcclient.Client).GetBlock(hash)
//	if err != nil {
//		t.Fatal("Blockhash does not refer to valid block")
//	}
//
//	// "regular" not "stake" txtype
//	if *txList1[0].TxType != dcrjson.LTTTRegular {
//		t.Fatal(`txtype not "regular".`)
//	}
//
//	// ListUnspent only shows validated (confirmations>=1) coinbase tx, so the
//	// first result should have 2 confirmations.
//	if txList1[0].Confirmations != 1 {
//		t.Fatalf("Latest coinbase tx listed has %v confirmations, expected 1.",
//			txList1[0].Confirmations)
//	}
//
//	// Check txid
//	txid, err := chainhash.NewHashFromStr(txList1[0].TxID)
//	if err != nil {
//		t.Fatal("Invalid Txid: ", err)
//	}
//
//	rawTx, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetRawTransaction(txid)
//	if err != nil {
//		t.Fatal("Invalid Txid: ", err)
//	}
//
//	// Use Vout from listtransaction to index []TxOut from getrawtransaction.
//	if len(rawTx.MsgTx().TxOut) <= int(txList1[0].Vout) {
//		t.Fatal("Too few vouts.")
//	}
//	txOut := rawTx.MsgTx().TxOut[txList1[0].Vout]
//	voutAmt := dcrutil.Amount(txOut.Value).ToCoin()
//	// Verify amounts agree
//	if txList1[0].Amount != voutAmt {
//		t.Fatalf("Listed amount %v does not match expected vout amount %v",
//			txList1[0].Amount, voutAmt)
//	}
//
//	// Test number of transactions (count).  With only coinbase in this harness,
//	// length of result slice should be equal to number requested.
//	txList2, err := r.WalletRPCClient().Internal().(*rpcclient.Client).ListTransactionsCount("*", 2)
//	if err != nil {
//		t.Fatal("ListTransactionsCount failed:", err)
//	}
//
//	// With only coinbase transactions, there will only be one result per tx
//	if len(txList2) != 2 {
//		t.Fatalf("Expected 2 transactions, got %v", len(txList2))
//	}
//
//	// List all transactions
//	txListAllInit, err := r.WalletRPCClient().Internal().(*rpcclient.Client).ListTransactionsCount("*", 9999999)
//	if err != nil {
//		t.Fatal("ListTransactionsCount failed:", err)
//	}
//	initNumTx := len(txListAllInit)
//
//	// Send within wallet, and check for both send and receive parts of tx.
//	accountName := "listTransactionsTest"
//	err = wcl.CreateNewAccount(accountName)
//	if err != nil {
//		t.Fatalf("Failed to create account for listtransactions test, %v", err)
//	}
//
//	addr, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetNewAddressGapPolicy(
//		accountName,
//		rpcclient.GapPolicyWrap,
//	)
//	if err != nil {
//		t.Fatal("Failed to get new address.")
//	}
//
//	atomsInCoin := dcrutil.AtomsPerCoin
//	sendAmount := dcrutil.Amount(2400 * atomsInCoin)
//	txHash, err := r.WalletRPCClient().Internal().(*rpcclient.Client).SendFromMinConf("default", addr, sendAmount, 6)
//	if err != nil {
//		t.Fatal("Failed to send:", err)
//	}
//
//	// Mine next block
//	mineBlock(t, r)
//
//	// Number of results should be +3 now
//	txListAll, err := r.WalletRPCClient().Internal().(*rpcclient.Client).ListTransactionsCount("*", 9999999)
//	txListAll = reverse(txListAll)
//	if err != nil {
//		t.Fatal("ListTransactionsCount failed:", err)
//	}
//	// Expect 3 more results in the list: a receive for the owned address in
//	// the amount sent, a send in the amount sent, and the a send from the
//	// original outpoint for the mined coins.
//	expectedAdditional := 3
//	if len(txListAll) != initNumTx+expectedAdditional {
//		t.Fatalf("Expected %v listtransactions results, got %v", initNumTx+expectedAdditional,
//			len(txListAll))
//	}
//
//	// The top of the list should be one send and one receive.  The coinbase
//	// spend should be lower in the list.
//	var sendResult, recvResult dcrjson.ListTransactionsResult
//	if txListAll[0].Category == txListAll[1].Category {
//		t.Fatal("Expected one send and one receive, got two", txListAll[0].Category)
//	}
//	// Use a map since order doesn't matter, and keys are not duplicate
//	rxtxResults := map[string]dcrjson.ListTransactionsResult{
//		txListAll[0].Category: txListAll[0],
//		txListAll[1].Category: txListAll[1],
//	}
//	var ok bool
//	if sendResult, ok = rxtxResults["send"]; !ok {
//		t.Fatal("Expected send transaction not found.")
//	}
//	if recvResult, ok = rxtxResults["receive"]; !ok {
//		t.Fatal("Expected receive transaction not found.")
//	}
//
//	// Verify send result amount
//	if sendResult.Amount != -sendAmount.ToCoin() {
//		t.Fatalf("Listed send tx amount incorrect. Expected %v, got %v",
//			-sendAmount.ToCoin(), sendResult.Amount)
//	}
//
//	// Verify send result fee
//	if sendResult.Fee == nil {
//		t.Fatal("Fee in send tx result is nil.")
//	}
//
//	// last transactions:
//	// ...
//	//  [4] coinbase of block 40
//	//  [3] coinbase of block 41
//	//  [2] new coinbase
//	//  [1] send
//	//  [0] receive
//	//
//	txList1New, err := r.WalletRPCClient().Internal().(*rpcclient.Client).ListTransactionsCount("*", 3)
//	if err != nil {
//		t.Fatal("Failed to listtransactions:", err)
//	}
//	txList1New = reverse(txList1New)
//	// txList1New is:
//	//  [3] coinbase of block 41
//	//  [2] new coinbase
//	//  [1] send
//	//  [0] receive
//	//
//
//	//coinbase of block 41
//	cb1 := txList1[0]
//	//one block passed, so update to match
//	cb1.Confirmations = cb1.Confirmations + 1
//	cb1n := txList1New[3]
//
//	// Should be equal to earlier result
//	if !cmp.Equal(cb1, cb1n) {
//		t.Fatal("Listtransaction results not equal. " + cmp.Diff(cb1, cb1n))
//	}
//
//	// Get rawTx of sent txid so we can calculate the fee that was used
//	newBestBlock(r, t) // or getrawtransaction is wrong
//	rawTx, err = r.NodeRPCClient().Internal().(*rpcclient.Client).GetRawTransaction(txHash)
//	if err != nil {
//		t.Fatalf("getrawtransaction failed: %v", err)
//	}
//
//	expectedFee := getWireMsgTxFee(rawTx).ToCoin()
//	gotFee := -*sendResult.Fee
//	if gotFee != expectedFee {
//		t.Fatalf("Expected fee %v, got %v", expectedFee, gotFee)
//	}
//
//	// Verify receive results amount
//	if recvResult.Amount != sendAmount.ToCoin() {
//		t.Fatalf("Listed send tx amount incorrect. Expected %v, got %v",
//			sendAmount.ToCoin(), recvResult.Amount)
//	}
//
//	// Verify TxID in both send and receive results
//	txstr := txHash.String()
//	if sendResult.TxID != txstr {
//		t.Fatalf("TxID in send tx result was %v, expected %v.",
//			sendResult.TxID, txstr)
//	}
//	if recvResult.TxID != txstr {
//		t.Fatalf("TxID in receive tx result was %v, expected %v.",
//			recvResult.TxID, txstr)
//	}
//
//	// Should only accept "*" account
//	_, err = r.WalletRPCClient().Internal().(*rpcclient.Client).ListTransactions("default")
//	if err == nil {
//		t.Fatal(`Listtransactions should only work on "*" account. "default" succeeded.`)
//	}
//
//	txList0, err := r.WalletRPCClient().Internal().(*rpcclient.Client).ListTransactionsCount("*", 0)
//	if err != nil {
//		t.Fatal("listtransactions failed:", err)
//	}
//	if len(txList0) != 0 {
//		t.Fatal("Length of listransactions result not zero:", len(txList0))
//	}
//
//	txListAll, err = r.WalletRPCClient().Internal().(*rpcclient.Client).ListTransactionsCount("*", 99999999)
//	if err != nil {
//		t.Fatal("ListTransactionsCount failed:", err)
//	}
//
//	// Create 2 accounts to receive funds
//	accountNames := []string{"listTxA", "listTxB"}
//	amountsToSend := []dcrutil.Amount{
//		dcrutil.Amount(7 * atomsInCoin),
//		dcrutil.Amount(14 * atomsInCoin),
//	}
//
//	for _, acct := range accountNames {
//		err := wcl.CreateNewAccount(acct)
//		if err != nil {
//			t.Fatal(err)
//		}
//	}
//
//	// Grab new addresses from the wallet, under each account.
//	// Set corresponding amount to send to each address.
//	addressAmounts := make(map[dcrutil.Address]dcrutil.Amount)
//
//	for i, acct := range accountNames {
//		addr, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetNewAddressGapPolicy(
//			acct,
//			rpcclient.GapPolicyWrap,
//		)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		// Set the amounts to send to each address
//		addressAmounts[addr] = amountsToSend[i]
//	}
//
//	// SendMany to two addresses
//	_, err = r.WalletRPCClient().Internal().(*rpcclient.Client).SendMany("default", addressAmounts)
//	if err != nil {
//		t.Fatalf("sendmany failed: %v", err)
//	}
//
//	// Mine next block
//	mineBlock(t, r)
//
//	// This should add 5 results: coinbase send, 2 receives, 2 sends
//	listSentMany, err := r.WalletRPCClient().Internal().(*rpcclient.Client).ListTransactionsCount("*", 99999999)
//	if err != nil {
//		t.Fatalf("ListTransactionsCount failed: %v", err)
//	}
//	if len(listSentMany) != len(txListAll)+5 {
//		t.Fatalf("Expected %v tx results, got %v", len(txListAll)+5,
//			len(listSentMany))
//	}
//}
//
//func TestGetSetRelayFee(t *testing.T) {
//
//	r := ObtainWalletHarness(mainWalletHarnessName)
//
//	// dcrrpcclient does not have a getwalletfee or any direct method, so we
//	// need to use walletinfo to get.  SetTxFee can be used to set.
//
//	// Wallet RPC client
//	wcl := r.Wallet
//
//	// Increase the ticket fee so these SSTx get mined first
//	walletInfo, err := r.WalletRPCClient().Internal().(*rpcclient.Client).WalletInfo()
//	if err != nil {
//		t.Fatal("WalletInfo failed:", err)
//	}
//	// Save the original fee
//	origTxFee, err := dcrutil.NewAmount(walletInfo.TxFee)
//	if err != nil {
//		t.Fatalf("Invalid Amount %f. %v", walletInfo.TxFee, err)
//	}
//	// Increase fee by 50%
//	newTxFeeCoin := walletInfo.TxFee * 1.5
//	newTxFee, err := dcrutil.NewAmount(newTxFeeCoin)
//	if err != nil {
//		t.Fatalf("Invalid Amount %f. %v", newTxFeeCoin, err)
//	}
//
//	err = r.WalletRPCClient().Internal().(*rpcclient.Client).SetTxFee(newTxFee)
//	if err != nil {
//		t.Fatal("SetTxFee failed:", err)
//	}
//
//	// Check that wallet thinks the fee is as expected
//	walletInfo, err = r.WalletRPCClient().Internal().(*rpcclient.Client).WalletInfo()
//	if err != nil {
//		t.Fatal("WalletInfo failed:", err)
//	}
//	newTxFeeActual, err := dcrutil.NewAmount(walletInfo.TxFee)
//	if err != nil {
//		t.Fatalf("Invalid Amount %f. %v", walletInfo.TxFee, err)
//	}
//	if newTxFee != newTxFeeActual {
//		t.Fatalf("Expected tx fee %v, got %v.", newTxFee, newTxFeeActual)
//	}
//
//	// Create a transaction and compute the effective fee
//	accountName := "testGetSetRelayFee"
//	err = wcl.CreateNewAccount(accountName)
//	if err != nil {
//		t.Fatal("Failed to create account.")
//	}
//
//	// Grab a fresh address from the test account
//	addr, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetNewAddressGapPolicy(
//		accountName,
//		rpcclient.GapPolicyWrap,
//	)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// SendFromMinConf to addr
//	amountToSend := dcrutil.Amount(700000000)
//	txid, err := r.WalletRPCClient().Internal().(*rpcclient.Client).SendFromMinConf("default", addr, amountToSend, 0)
//	if err != nil {
//		t.Fatalf("sendfromminconf failed: %v", err)
//	}
//
//	newBestBlock(r, t)
//
//	// Compute the fee
//	rawTx, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetRawTransaction(txid)
//	if err != nil {
//		t.Fatalf("getrawtransaction failed: %v", err)
//	}
//
//	fee := getWireMsgTxFee(rawTx)
//	feeRate := fee.ToCoin() / float64(rawTx.MsgTx().SerializeSize()) * 1000
//
//	// Ensure actual fee is at least nominal
//	if feeRate < walletInfo.TxFee {
//		t.Errorf("Regular tx fee rate difference (actual-set) too high: %v",
//			walletInfo.TxFee-feeRate)
//	}
//
//	// Negative fee should throw an error
//	err = r.WalletRPCClient().Internal().(*rpcclient.Client).SetTxFee(dcrutil.Amount(-1))
//	if err == nil {
//		t.Fatal("SetTxFee accepted negative fee")
//	}
//
//	// Set it back
//	err = r.WalletRPCClient().Internal().(*rpcclient.Client).SetTxFee(origTxFee)
//	if err != nil {
//		t.Fatal("SetTxFee failed:", err)
//	}
//
//	// Validate last tx before we complete
//	newBestBlock(r, t)
//}
//
//func TestGetSetTicketFee(t *testing.T) {
//
//	r := ObtainWalletHarness(mainWalletHarnessName)
//	// dcrrpcclient does not have a getticketee or any direct method, so we
//	// need to use walletinfo to get.  SetTicketFee can be used to set.
//
//	// Get the current ticket fee
//	walletInfo, err := r.WalletRPCClient().Internal().(*rpcclient.Client).WalletInfo()
//	if err != nil {
//		t.Fatal("WalletInfo failed:", err)
//	}
//	nominalTicketFee := walletInfo.TicketFee
//	origTicketFee, err := dcrutil.NewAmount(nominalTicketFee)
//	if err != nil {
//		t.Fatal("Invalid Amount:", nominalTicketFee)
//	}
//
//	// Increase the ticket fee to ensure the SSTx in ths test gets mined
//	newTicketFeeCoin := nominalTicketFee * 1.5
//	newTicketFee, err := dcrutil.NewAmount(newTicketFeeCoin)
//	if err != nil {
//		t.Fatal("Invalid Amount:", newTicketFeeCoin)
//	}
//
//	err = r.WalletRPCClient().Internal().(*rpcclient.Client).SetTicketFee(newTicketFee)
//	if err != nil {
//		t.Fatal("SetTicketFee failed:", err)
//	}
//
//	// Check that wallet is set to use the new fee
//	walletInfo, err = r.WalletRPCClient().Internal().(*rpcclient.Client).WalletInfo()
//	if err != nil {
//		t.Fatal("WalletInfo failed:", err)
//	}
//	nominalTicketFee = walletInfo.TicketFee
//	newTicketFeeActual, err := dcrutil.NewAmount(nominalTicketFee)
//	if err != nil {
//		t.Fatalf("Invalid Amount %f. %v", nominalTicketFee, err)
//	}
//	if newTicketFee != newTicketFeeActual {
//		t.Fatalf("Expected ticket fee %v, got %v.", newTicketFee,
//			newTicketFeeActual)
//	}
//
//	// Purchase ticket
//	minConf, numTickets := 0, 1
//	priceLimit, err := dcrutil.NewAmount(2 * mustGetStakeDiffNext(r, t))
//	if err != nil {
//		t.Fatal("Invalid Amount. ", err)
//	}
//	noSplitTransactions := false
//	hashes, err := r.WalletRPCClient().Internal().(*rpcclient.Client).PurchaseTicket("default", priceLimit,
//		&minConf, nil, &numTickets, nil, nil, nil, &noSplitTransactions, nil)
//	if err != nil {
//		t.Fatal("Unable to purchase ticket:", err)
//	}
//	if len(hashes) != numTickets {
//		t.Fatalf("Number of returned hashes does not equal expected."+
//			"got %v, want %v", len(hashes), numTickets)
//	}
//
//	// Need 2 blocks or the vin is incorrect in getrawtransaction
//	// Not yet at StakeValidationHeight, so no voting.
//	newBestBlock(r, t)
//	newBestBlock(r, t)
//
//	// Compute the actual fee for the ticket purchase
//	rawTx, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetRawTransaction(hashes[0])
//	if err != nil {
//		t.Fatal("Invalid Txid:", err)
//	}
//
//	fee := getWireMsgTxFee(rawTx)
//	feeRate := fee.ToCoin() / float64(rawTx.MsgTx().SerializeSize()) * 1000
//
//	// Ensure actual fee is at least nominal
//	if feeRate < nominalTicketFee {
//		t.Errorf("Ticket fee rate difference (actual-set) too high: %v",
//			nominalTicketFee-feeRate)
//	}
//
//	// Negative fee should throw and error
//	err = r.WalletRPCClient().Internal().(*rpcclient.Client).SetTicketFee(dcrutil.Amount(-1))
//	if err == nil {
//		t.Fatal("SetTicketFee accepted negative fee")
//	}
//
//	// Set it back
//	err = r.WalletRPCClient().Internal().(*rpcclient.Client).SetTicketFee(origTicketFee)
//	if err != nil {
//		t.Fatal("SetTicketFee failed:", err)
//	}
//
//	// Validate last tx before we complete
//	newBestBlock(r, t)
//}
//
//func TestGetTickets(t *testing.T) {
//
//	r := ObtainWalletHarness(mainWalletHarnessName)
//	// Wallet.purchaseTicket() in wallet/createtx.go
//
//	// Initial number of mature (live) tickets
//	ticketHashes, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetTickets(false)
//	if err != nil {
//		t.Fatal("GetTickets failed:", err)
//	}
//	numTicketsInitLive := len(ticketHashes)
//
//	// Initial number of immature (not live) and unconfirmed (unmined) tickets
//	ticketHashes, err = r.WalletRPCClient().Internal().(*rpcclient.Client).GetTickets(true)
//	if err != nil {
//		t.Fatal("GetTickets failed:", err)
//	}
//
//	numTicketsInit := len(ticketHashes)
//
//	// Purchase a full blocks worth of tickets
//	minConf, numTicketsPurchased := 1, int(chaincfg.SimNetParams.MaxFreshStakePerBlock)
//	priceLimit, err := dcrutil.NewAmount(2 * mustGetStakeDiffNext(r, t))
//	if err != nil {
//		t.Fatal("Invalid Amount. ", err)
//	}
//	noSplitTransactions := false
//	hashes, err := r.WalletRPCClient().Internal().(*rpcclient.Client).PurchaseTicket("default", priceLimit,
//		&minConf, nil, &numTicketsPurchased, nil, nil, nil, &noSplitTransactions, nil)
//	if err != nil {
//		t.Fatal("Unable to purchase tickets:", err)
//	}
//	if len(hashes) != numTicketsPurchased {
//		t.Fatalf("Expected %v ticket hashes, got %v.", numTicketsPurchased,
//			len(hashes))
//	}
//
//	// Verify GetTickets(true) sees these unconfirmed SSTx
//	ticketHashes, err = r.WalletRPCClient().Internal().(*rpcclient.Client).GetTickets(true)
//	if err != nil {
//		t.Fatal("GetTickets failed:", err)
//	}
//
//	if numTicketsInit+numTicketsPurchased != len(ticketHashes) {
//		t.Fatal("GetTickets(true) did not include unmined tickets")
//	}
//
//	// Compare GetTickets(includeImmature = false) before the purchase with
//	// GetTickets(includeImmature = true) after the purchase. This tests that
//	// the former does exclude unconfirmed tickets, which we now have following
//	// the above purchase.
//	if len(ticketHashes) <= numTicketsInitLive {
//		t.Fatalf("Number of live tickets (%d) not less than total tickets (%d).",
//			numTicketsInitLive, len(ticketHashes))
//	}
//
//	// Mine the split tx and THEN stake submission itself
//	newBestBlock(r, t)
//	_, block, _ := newBestBlock(r, t)
//
//	// Verify stake submissions were mined
//	for _, hash := range hashes {
//		if !includesStakeTx(hash, block) {
//			t.Errorf("SSTx expected, not found in block %v.", block.Height())
//		}
//	}
//
//	// Verify each SSTx hash
//	for _, hash := range ticketHashes {
//		tx, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetRawTransaction(hash)
//		if err != nil {
//			t.Fatalf("Invalid transaction %v: %v", tx, err)
//		}
//
//		// Ensure result is a SSTx
//		if !stake.IsSStx(tx.MsgTx()) {
//			t.Fatal("Ticket hash is not for a SSTx.")
//		}
//	}
//}
//
//func TestPurchaseTickets(t *testing.T) {
//	t.SkipNow()
//	r := ObtainWalletHarness(mainWalletHarnessName)
//	// Wallet.purchaseTicket() in wallet/createtx.go
//
//	// Grab a fresh address from the wallet.
//	addr, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetNewAddressGapPolicy(
//		"default",
//		rpcclient.GapPolicyWrap,
//	)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// Set various variables for the test
//	minConf := 0
//	expiry := 0
//	priceLimit, err := dcrutil.NewAmount(2 * mustGetStakeDiffNext(r, t))
//	if err != nil {
//		t.Fatal("Invalid Amount.", err)
//	}
//
//	// Test nil ticketAddress
//	oneTix := 1
//	noSplitTransactions := false
//	hashes, err := r.WalletRPCClient().Internal().(*rpcclient.Client).PurchaseTicket("default", priceLimit,
//		&minConf, nil, &oneTix, nil, nil, &expiry, &noSplitTransactions, nil)
//	if err != nil {
//		t.Fatal("Unable to purchase with nil ticketAddr:", err)
//	}
//	if len(hashes) != 1 {
//		t.Fatal("More than one tx hash returned purchasing single ticket.")
//	}
//	_, err = r.WalletRPCClient().Internal().(*rpcclient.Client).GetRawTransaction(hashes[0])
//	if err != nil {
//		t.Fatal("Invalid Txid:", err)
//	}
//
//	// test numTickets == nil
//	hashes, err = r.WalletRPCClient().Internal().(*rpcclient.Client).PurchaseTicket("default", priceLimit,
//		&minConf, nil, nil, nil, nil, &expiry, &noSplitTransactions, nil)
//	if err != nil {
//		t.Fatal("Unable to purchase with nil numTickets:", err)
//	}
//	if len(hashes) != 1 {
//		t.Fatal("More than one tx hash returned. Expected one.")
//	}
//	_, err = r.WalletRPCClient().Internal().(*rpcclient.Client).GetRawTransaction(hashes[0])
//	if err != nil {
//		t.Fatal("Invalid Txid:", err)
//	}
//
//	// Get current blockheight to make sure chain is at the desiredHeight
//	curBlockHeight := getBestBlockHeight(r, t)
//
//	// Test expiry - earliest is next height + 1
//	// invalid
//	expiry = int(curBlockHeight)
//	_, err = r.WalletRPCClient().Internal().(*rpcclient.Client).PurchaseTicket("default", priceLimit,
//		&minConf, nil, nil, nil, nil, &expiry, &noSplitTransactions, nil)
//	if err == nil {
//		t.Fatal("Invalid expiry used to purchase tickets")
//	}
//	// invalid
//	expiry = int(curBlockHeight) + 1
//	_, err = r.WalletRPCClient().Internal().(*rpcclient.Client).PurchaseTicket("default", priceLimit,
//		&minConf, nil, nil, nil, nil, &expiry, &noSplitTransactions, nil)
//	if err == nil {
//		t.Fatal("Invalid expiry used to purchase tickets")
//	}
//
//	// valid expiry
//	expiry = int(curBlockHeight) + 2
//	hashes, err = r.WalletRPCClient().Internal().(*rpcclient.Client).PurchaseTicket("default", priceLimit,
//		&minConf, nil, nil, nil, nil, &expiry, &noSplitTransactions, nil)
//	if err != nil {
//		t.Fatal("Unable to purchase tickets:", err)
//	}
//	if len(hashes) != 1 {
//		t.Fatal("More than one tx hash returned. Expected one.")
//	}
//	ticketWithExpiry := hashes[0]
//	_, err = r.WalletRPCClient().Internal().(*rpcclient.Client).GetRawTransaction(ticketWithExpiry)
//	if err != nil {
//		t.Fatal("Invalid Txid:", err)
//	}
//
//	// Now purchase 2 blocks worth of tickets to be mined before the above
//	// ticket with an expiry 2 blocks away.
//
//	// Increase the ticket fee so these SSTx get mined first
//	walletInfo, err := r.WalletRPCClient().Internal().(*rpcclient.Client).WalletInfo()
//	if err != nil {
//		t.Fatal("WalletInfo failed.", err)
//	}
//	origTicketFee, err := dcrutil.NewAmount(walletInfo.TicketFee)
//	if err != nil {
//		t.Fatalf("Invalid Amount %f. %v", walletInfo.TicketFee, err)
//	}
//	newTicketFee, err := dcrutil.NewAmount(walletInfo.TicketFee * 1.5)
//	if err != nil {
//		t.Fatalf("Invalid Amount %f. %v", walletInfo.TicketFee, err)
//	}
//
//	if err = r.WalletRPCClient().Internal().(*rpcclient.Client).SetTicketFee(newTicketFee); err != nil {
//		t.Fatalf("SetTicketFee failed for Amount %v: %v", newTicketFee, err)
//	}
//
//	expiry = 0
//	numTicket := 2 * int(chaincfg.SimNetParams.MaxFreshStakePerBlock)
//	_, err = r.WalletRPCClient().Internal().(*rpcclient.Client).PurchaseTicket("default", priceLimit,
//		&minConf, addr, &numTicket, nil, nil, &expiry, &noSplitTransactions, nil)
//	if err != nil {
//		t.Fatal("Unable to purchase tickets:", err)
//	}
//
//	if err = r.WalletRPCClient().Internal().(*rpcclient.Client).SetTicketFee(origTicketFee); err != nil {
//		t.Fatalf("SetTicketFee failed for Amount %v: %v", origTicketFee, err)
//	}
//
//	// Check for the ticket
//	_, err = r.WalletRPCClient().Internal().(*rpcclient.Client).GetTransaction(ticketWithExpiry)
//	if err != nil {
//		t.Fatal("Ticket not found:", err)
//	}
//
//	// Mine 2 blocks, should include the higher fee tickets with no expiry
//	curBlockHeight, _, _ = newBlockAt(curBlockHeight, r, t)
//	curBlockHeight, _, _ = newBlockAt(curBlockHeight, r, t)
//
//	// Ticket with expiry set should now be expired (unmined and removed from
//	// mempool).  An unmined and expired tx should have been removed/pruned
//	txRawVerbose, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetRawTransactionVerbose(ticketWithExpiry)
//	if err == nil {
//		t.Fatalf("Found transaction that should be expired (height %v): %v",
//			txRawVerbose.BlockHeight, err)
//	}
//
//	// Test too low price
//	lowPrice := dcrutil.Amount(1)
//	hashes, err = r.WalletRPCClient().Internal().(*rpcclient.Client).PurchaseTicket("default", lowPrice,
//		&minConf, nil, nil, nil, nil, nil, &noSplitTransactions, nil)
//	if err == nil {
//		t.Fatalf("PurchaseTicket succeeded with limit of %f, but diff was %f.",
//			lowPrice.ToCoin(), mustGetStakeDiff(r, t))
//	}
//	if len(hashes) > 0 {
//		t.Fatal("At least one tickets hash returned. Expected none.")
//	}
//
//	// NOTE: ticket maturity = 16 (spendable at 17), stakeenabled height = 144
//	// Must have tickets purchased before block 128
//
//	// Keep generating blocks until desiredHeight is achieved
//	desiredHeight := uint32(150)
//	numTicket = int(chaincfg.SimNetParams.MaxFreshStakePerBlock)
//	for curBlockHeight < desiredHeight {
//		priceLimit, err = dcrutil.NewAmount(2 * mustGetStakeDiffNext(r, t))
//		if err != nil {
//			t.Fatal("Invalid Amount.", err)
//		}
//		_, err = r.WalletRPCClient().Internal().(*rpcclient.Client).PurchaseTicket("default", priceLimit,
//			&minConf, addr, &numTicket, nil, nil, nil, &noSplitTransactions, nil)
//
//		// Do not allow even ErrSStxPriceExceedsSpendLimit since price is set
//		if err != nil {
//			t.Fatal("Failed to purchase tickets:", err)
//		}
//		curBlockHeight, _, _ = newBlockAtQuick(curBlockHeight, r, t)
//		time.Sleep(100 * time.Millisecond)
//	}
//
//	// Validate last tx
//	newBestBlock(r, t)
//
//	// TODO: test pool fees
//
//}
//
//// testGetStakeInfo gets a FRESH harness
//func TestGetStakeInfo(t *testing.T) {
//	t.SkipNow()
//	r := ObtainWalletHarness(t.Name() + "-harness")
//
//	// Compare stake difficulty from getstakeinfo with getstakeinfo
//	sdiff, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetStakeDifficulty()
//	if err != nil {
//		t.Fatal("GetStakeDifficulty failed: ", err)
//	}
//
//	stakeinfo, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetStakeInfo()
//	if err != nil {
//		t.Fatal("GetStakeInfo failed: ", err)
//	}
//	// Ensure we are starting with a fresh harness
//	if stakeinfo.AllMempoolTix != 0 || stakeinfo.Immature != 0 ||
//		stakeinfo.Live != 0 {
//		t.Fatalf("GetStakeInfo reported active tickets. Expected 0, got:\n"+
//			"%d/%d/%d (allmempooltix/immature/live)",
//			stakeinfo.AllMempoolTix, stakeinfo.Immature, stakeinfo.Live)
//	}
//	// At the expected block height
//	height, block, _ := getBestBlock(r, t)
//	if stakeinfo.BlockHeight != int64(height) {
//		t.Fatalf("Block height reported by GetStakeInfo incorrect. Expected %d, got %d.",
//			height, stakeinfo.BlockHeight)
//	}
//	poolSize := block.MsgBlock().Header.PoolSize
//	if stakeinfo.PoolSize != poolSize {
//		t.Fatalf("Reported pool size incorrect. Expected %d, got %d.",
//			poolSize, stakeinfo.PoolSize)
//	}
//
//	// Ticket fate values should also be zero
//	if stakeinfo.Voted != 0 || stakeinfo.Missed != 0 ||
//		stakeinfo.Revoked != 0 {
//		t.Fatalf("GetStakeInfo reported spent tickets:\n"+
//			"%d/%d/%d (voted/missed/revoked/pct. missed)", stakeinfo.Voted,
//			stakeinfo.Missed, stakeinfo.Revoked)
//	}
//	if stakeinfo.ProportionLive != 0 {
//		t.Fatalf("ProportionLive incorrect. Expected %f, got %f.", 0.0,
//			stakeinfo.ProportionLive)
//	}
//	if stakeinfo.ProportionMissed != 0 {
//		t.Fatalf("ProportionMissed incorrect. Expected %f, got %f.", 0.0,
//			stakeinfo.ProportionMissed)
//	}
//
//	// Verify getstakeinfo.difficulty == getstakedifficulty
//	if sdiff.CurrentStakeDifficulty != stakeinfo.Difficulty {
//		t.Fatalf("Stake difficulty mismatch: %f vs %f (getstakedifficulty, getstakeinfo)",
//			sdiff.CurrentStakeDifficulty, stakeinfo.Difficulty)
//	}
//
//	// Buy tickets to check that they shows up in ownmempooltix/allmempooltix
//	minConf := 1
//	priceLimit, err := dcrutil.NewAmount(2 * mustGetStakeDiffNext(r, t))
//	if err != nil {
//		t.Fatal("Invalid Amount.", err)
//	}
//	numTickets := int(chaincfg.SimNetParams.MaxFreshStakePerBlock)
//	noSplitTransactions := false
//	tickets, err := r.WalletRPCClient().Internal().(*rpcclient.Client).PurchaseTicket("default", priceLimit,
//		&minConf, nil, &numTickets, nil, nil, nil, &noSplitTransactions, nil)
//	if err != nil {
//		t.Fatal("Failed to purchase tickets:", err)
//	}
//
//	// Before mining a block allmempooltix and ownmempooltix should be equal to
//	// the number of tickets just purchesed in this fresh harness
//	stakeinfo = mustGetStakeInfo(r.WalletRPCClient().Internal().(*rpcclient.Client), t)
//	if stakeinfo.AllMempoolTix != uint32(numTickets) {
//		t.Fatalf("getstakeinfo AllMempoolTix mismatch: %d vs %d",
//			stakeinfo.AllMempoolTix, numTickets)
//	}
//	if stakeinfo.AllMempoolTix != stakeinfo.OwnMempoolTix {
//		t.Fatalf("getstakeinfo AllMempoolTix/OwnMempoolTix mismatch: %d vs %d",
//			stakeinfo.AllMempoolTix, stakeinfo.OwnMempoolTix)
//	}
//
//	// Mine the split tx, which creates the correctly-sized outpoints for the
//	// actual SSTx
//	newBestBlock(r, t)
//	// Mine SSTx
//	newBestBlock(r, t)
//
//	// Compute the height at which these tickets mature
//	ticketsTx, err := r.WalletRPCClient().Internal().(*rpcclient.Client).GetRawTransactionVerbose(tickets[0])
//	if err != nil {
//		t.Fatalf("Unable to gettransaction for ticket.")
//	}
//	maturityHeight := ticketsTx.BlockHeight + int64(chaincfg.SimNetParams.TicketMaturity)
//
//	// After mining tickets, immature should be the number of tickets
//	stakeinfo = mustGetStakeInfo(r.WalletRPCClient().Internal().(*rpcclient.Client), t)
//	if stakeinfo.Immature != uint32(numTickets) {
//		t.Fatalf("Tickets not reported as immature (got %d, expected %d)",
//			stakeinfo.Immature, numTickets)
//	}
//	// mempool tickets should be zero
//	if stakeinfo.OwnMempoolTix != 0 {
//		t.Fatalf("Tickets reported in mempool (got %d, expected %d)",
//			stakeinfo.OwnMempoolTix, 0)
//	}
//	// mempool tickets should be zero
//	if stakeinfo.AllMempoolTix != 0 {
//		t.Fatalf("Tickets reported in mempool (got %d, expected %d)",
//			stakeinfo.AllMempoolTix, 0)
//	}
//
//	// Advance to maturity height
//	t.Logf("Advancing to maturity height %d for tickets in block %d", maturityHeight,
//		ticketsTx.BlockHeight)
//	advanceToHeight(r, t, uint32(maturityHeight))
//	// NOTE: voting does not begin until TicketValidationHeight
//
//	// mature should be number of tickets now
//	stakeinfo = mustGetStakeInfo(r.WalletRPCClient().Internal().(*rpcclient.Client), t)
//	if stakeinfo.Live != uint32(numTickets) {
//		t.Fatalf("Tickets not reported as live (got %d, expected %d)",
//			stakeinfo.Live, numTickets)
//	}
//	// immature tickets should be zero
//	if stakeinfo.Immature != 0 {
//		t.Fatalf("Tickets reported as immature (got %d, expected %d)",
//			stakeinfo.Immature, 0)
//	}
//
//	// Buy some more tickets (4 blocks worth) so chain doesn't stall when voting
//	// burns through the batch purchased above
//	for i := 0; i < 4; i++ {
//		priceLimit, err := dcrutil.NewAmount(2 * mustGetStakeDiffNext(r, t))
//		if err != nil {
//			t.Fatal("Invalid Amount.", err)
//		}
//		numTickets := int(chaincfg.SimNetParams.MaxFreshStakePerBlock)
//		_, err = r.WalletRPCClient().Internal().(*rpcclient.Client).PurchaseTicket("default", priceLimit,
//			&minConf, nil, &numTickets, nil, nil, nil, &noSplitTransactions, nil)
//		if err != nil {
//			t.Fatal("Failed to purchase tickets:", err)
//		}
//
//		newBestBlock(r, t)
//	}
//
//	// Advance to voting height and votes should happen right away
//	votingHeight := chaincfg.SimNetParams.StakeValidationHeight
//	advanceToHeight(r, t, uint32(votingHeight))
//	time.Sleep(250 * time.Millisecond)
//
//	// voted should be TicketsPerBlock
//	stakeinfo = mustGetStakeInfo(r.WalletRPCClient().Internal().(*rpcclient.Client), t)
//	expectedVotes := chaincfg.SimNetParams.TicketsPerBlock
//	if stakeinfo.Voted != uint32(expectedVotes) {
//		t.Fatalf("Tickets not reported as voted (got %d, expected %d)",
//			stakeinfo.Voted, expectedVotes)
//	}
//
//	newBestBlock(r, t)
//	// voted should be 2*TicketsPerBlock
//	stakeinfo = mustGetStakeInfo(r.WalletRPCClient().Internal().(*rpcclient.Client), t)
//	expectedVotes = 2 * chaincfg.SimNetParams.TicketsPerBlock
//	if stakeinfo.Voted != uint32(expectedVotes) {
//		t.Fatalf("Tickets not reported as voted (got %d, expected %d)",
//			stakeinfo.Voted, expectedVotes)
//	}
//
//	// ProportionLive
//	proportionLive := float64(stakeinfo.Live) / float64(stakeinfo.PoolSize)
//	if stakeinfo.ProportionLive != proportionLive {
//		t.Fatalf("ProportionLive mismatch.  Expected %f, got %f",
//			proportionLive, stakeinfo.ProportionLive)
//	}
//
//	// ProportionMissed
//	proportionMissed := float64(stakeinfo.Missed) /
//		(float64(stakeinfo.Voted) + float64(stakeinfo.Missed))
//	if stakeinfo.ProportionMissed != proportionMissed {
//		t.Fatalf("ProportionMissed mismatch.  Expected %f, got %f",
//			proportionMissed, stakeinfo.ProportionMissed)
//	}
//}
//
//// testWalletInfo
//func TestWalletInfo(t *testing.T) {
//
//	r := ObtainWalletHarness(mainWalletHarnessName)
//
//	// WalletInfo is tested exhaustively in other test, so only do some basic
//	// checks here
//	walletInfo, err := r.WalletRPCClient().Internal().(*rpcclient.Client).WalletInfo()
//	if err != nil {
//		t.Fatal("walletinfo failed.")
//	}
//	if !walletInfo.DaemonConnected {
//		t.Fatal("WalletInfo indicates that daemon is not connected.")
//	}
//}
//
//func TestWalletPassphrase(t *testing.T) {
//	r := ObtainWalletHarness(mainWalletHarnessName)
//	// Wallet RPC client
//	wcl := r.Wallet
//
//	// Remember to leave the wallet unlocked for any subsequent tests
//
//	// Lock the wallet since test wallet is unlocked by default
//	err := wcl.WalletLock()
//	if err != nil {
//		t.Fatal("Unable to lock wallet.")
//	}
//
//	// Check that wallet is locked
//	walletInfo, err := wcl.WalletInfo()
//	if err != nil {
//		t.Fatal("walletinfo failed.")
//	}
//	if walletInfo.Unlocked {
//		t.Fatal("WalletLock failed to lock the wallet")
//	}
//
//	// Try incorrect password
//	err = wcl.WalletUnlock("Wrong Password", 0)
//	// Check for "-14: invalid passphrase for master private key"
//	if err != nil && err.(*dcrjson.RPCError).Code !=
//		dcrjson.ErrRPCWalletPassphraseIncorrect {
//		// dcrjson.ErrWalletPassphraseIncorrect.Code
//		t.Fatalf("WalletPassphrase with INCORRECT passphrase exited with: %v",
//			err)
//	}
//
//	// Check that wallet is still locked
//	walletInfo, err = wcl.WalletInfo()
//	if err != nil {
//		t.Fatal("walletinfo failed.")
//	}
//	if walletInfo.Unlocked {
//		t.Fatal("WalletPassphrase unlocked the wallet with the wrong passphrase")
//	}
//
//	// Verify that a restricted operation like createnewaccount fails
//	accountName := "cannotCreateThisAccount"
//	err = wcl.CreateNewAccount(accountName)
//	if err == nil {
//		t.Fatal("createnewaccount succeeded on a locked wallet.")
//	}
//	// dcrjson.ErrRPCWalletUnlockNeeded
//	if !strings.HasPrefix(err.Error(),
//		strconv.Itoa(int(dcrjson.ErrRPCWalletUnlockNeeded))) {
//		t.Fatalf("createnewaccount returned error (%v) instead of %v",
//			err, dcrjson.ErrRPCWalletUnlockNeeded)
//	}
//
//	// Unlock with correct passphrase
//	err = r.WalletRPCClient().Internal().(*rpcclient.Client).WalletPassphrase(defaultWalletPassphrase, 0)
//	if err != nil {
//		t.Fatalf("WalletPassphrase failed: %v", err)
//	}
//
//	// Check that wallet is now unlocked
//	walletInfo, err = wcl.WalletInfo()
//	if err != nil {
//		t.Fatal("walletinfo failed.")
//	}
//	if !walletInfo.Unlocked {
//		t.Fatal("WalletPassphrase failed to unlock the wallet with the correct passphrase")
//	}
//
//	// Check for ErrRPCWalletAlreadyUnlocked
//	err = wcl.WalletUnlock(defaultWalletPassphrase, 0)
//	// Check for "-17: Wallet is already unlocked"
//	if err != nil && err.(*dcrjson.RPCError).Code !=
//		dcrjson.ErrRPCWalletAlreadyUnlocked {
//		t.Fatalf("WalletPassphrase failed: %v", err)
//	}
//
//	// Re-lock wallet
//	err = wcl.WalletLock()
//	if err != nil {
//		t.Fatal("Unable to lock wallet.")
//	}
//
//	// Unlock with timeout
//	timeOut := int64(6)
//	err = wcl.WalletUnlock(defaultWalletPassphrase, timeOut)
//	if err != nil {
//		t.Fatalf("WalletPassphrase failed: %v", err)
//	}
//
//	// Check that wallet is now unlocked
//	walletInfo, err = wcl.WalletInfo()
//	if err != nil {
//		t.Fatal("walletinfo failed.")
//	}
//	if !walletInfo.Unlocked {
//		t.Fatal("WalletPassphrase failed to unlock the wallet with the correct passphrase")
//	}
//
//	time.Sleep(time.Duration(timeOut+2) * time.Second)
//
//	// Check that wallet is now locked
//	walletInfo, err = wcl.WalletInfo()
//	if err != nil {
//		t.Fatal("walletinfo failed.")
//	}
//	if walletInfo.Unlocked {
//		t.Fatal("Wallet still unlocked after timeout")
//	}
//
//	if err := wcl.WalletUnlock(defaultWalletPassphrase, 0); err != nil {
//		t.Fatal("Unable to unlock wallet:", err)
//	}
//
//	// TODO: Watching-only error?
//}
