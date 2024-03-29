package pfcregtest

import (
	"fmt"
	"github.com/jfixby/coinharness"
	"github.com/picfight/pfcharness"
	"github.com/jfixby/pin"
	"github.com/jfixby/pin/commandline"
	"github.com/jfixby/pin/gobuilder"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/picfight/pfcd/chaincfg"
)

// Default harness name
const mainHarnessName = "main"
const mainWalletHarnessName = "main-wallet"

// SimpleTestSetup harbours:
// - rpctest setup
// - csf-fork test setup
// - and bip0009 test setup
type SimpleTestSetup struct {
	// harnessPool stores and manages harnesses
	// multiple harness instances may be run concurrently, to allow for testing
	// complex scenarios involving multiple nodes.
	harnessPool *pin.Pool

	harnessWalletPool *pin.Pool

	// Mainnet creates a mainnet test harness
	Mainnet0 *coinharness.ChainWithMatureOutputsSpawner

	// Regnet25 creates a regnet test harness
	// with 25 mature outputs.
	Regnet25 *coinharness.ChainWithMatureOutputsSpawner

	// Simnet25 creates a simnet test harness
	// with 25 mature outputs.
	Simnet25 *coinharness.ChainWithMatureOutputsSpawner

	// Regnet5 creates a regnet test harness
	// with 5 mature outputs.
	Regnet5 *coinharness.ChainWithMatureOutputsSpawner

	// Regnet1 creates a regnet test harness
	// with 1 mature output.
	Regnet1 *coinharness.ChainWithMatureOutputsSpawner

	// Simnet1 creates a simnet test harness
	// with 1 mature output.
	Simnet1 *coinharness.ChainWithMatureOutputsSpawner

	// Regnet0 creates a regnet test harness
	// with only the genesis block.
	Regnet0 *coinharness.ChainWithMatureOutputsSpawner

	// Simnet0 creates a simnet test harness
	// with only the genesis block.
	Simnet0 *coinharness.ChainWithMatureOutputsSpawner

	// WorkingDir defines test setup working dir
	WorkingDir *pin.TempDirHandler
}

// TearDown all harnesses in test Pool.
// This includes removing all temporary directories,
// and shutting down any created processes.
func (setup *SimpleTestSetup) TearDown() {
	setup.harnessPool.DisposeAll()
	setup.harnessWalletPool.DisposeAll()
	//setup.nodeGoBuilder.Dispose()
	setup.WorkingDir.Dispose()
}

// Setup deploys this test setup
func Setup() *SimpleTestSetup {
	setup := &SimpleTestSetup{
		WorkingDir: pin.NewTempDir(setupWorkingDir(), "simpleregtest").MakeDir(),
	}

	//memWalletFactory := &memwallet.WalletFactory{}

	wEXE := &commandline.ExplicitExecutablePathString{
		PathString: "pfcwallet",
	}
	consoleWalletFactory := &pfcharness.ConsoleWalletFactory{
		WalletExecutablePathProvider: wEXE,
	}

	regnetWalletFactory := consoleWalletFactory
	mainnetWalletFactory := consoleWalletFactory
	simnetWalletFactory := consoleWalletFactory

	dEXE := &commandline.ExplicitExecutablePathString{
		PathString: "pfcd",
	}
	nodeFactory := &pfcharness.ConsoleNodeFactory{
		NodeExecutablePathProvider: dEXE,
	}

	portManager := &coinharness.LazyPortManager{
		BasePort: 30000,
	}

	testSeed := pfcharness.NewTestSeed

	// Deploy harness spawner with generated
	// test chain of 25 mature outputs
	setup.Regnet25 = &coinharness.ChainWithMatureOutputsSpawner{
		WorkingDir:        setup.WorkingDir.Path(),
		DebugNodeOutput:   true,
		DebugWalletOutput: true,
		NumMatureOutputs:  25,
		NetPortManager:    portManager,
		WalletFactory:     regnetWalletFactory,
		NodeFactory:       nodeFactory,
		ActiveNet:         &pfcharness.Network{&chaincfg.RegNetParams},
		CreateTempWallet:  true,
		NewTestSeed:       testSeed,
	}

	setup.Mainnet0 = &coinharness.ChainWithMatureOutputsSpawner{
		WorkingDir:        setup.WorkingDir.Path(),
		DebugNodeOutput:   true,
		DebugWalletOutput: true,
		NumMatureOutputs:  0,
		NetPortManager:    portManager,
		WalletFactory:     mainnetWalletFactory,
		NodeFactory:       nodeFactory,
		ActiveNet:         &pfcharness.Network{&chaincfg.PicFightCoinNetParams},
		CreateTempWallet:  true,
		NewTestSeed:       testSeed,
	}

	// Deploy harness spawner with generated
	// test chain of 5 mature outputs
	setup.Regnet5 = &coinharness.ChainWithMatureOutputsSpawner{
		WorkingDir:        setup.WorkingDir.Path(),
		DebugNodeOutput:   true,
		DebugWalletOutput: true,
		NumMatureOutputs:  5,
		NetPortManager:    portManager,
		WalletFactory:     regnetWalletFactory,
		NodeFactory:       nodeFactory,
		ActiveNet:         &pfcharness.Network{&chaincfg.RegNetParams},
		CreateTempWallet:  true,
		NewTestSeed:       testSeed,
	}

	setup.Regnet1 = &coinharness.ChainWithMatureOutputsSpawner{
		WorkingDir:        setup.WorkingDir.Path(),
		DebugNodeOutput:   true,
		DebugWalletOutput: true,
		NumMatureOutputs:  1,
		NetPortManager:    portManager,
		WalletFactory:     regnetWalletFactory,
		NodeFactory:       nodeFactory,
		ActiveNet:         &pfcharness.Network{&chaincfg.RegNetParams},
		CreateTempWallet:  true,
		NewTestSeed:       testSeed,
		NodeStartExtraArguments: map[string]interface{}{
			"rejectnonstd": commandline.NoArgumentValue,
		},
	}

	setup.Simnet1 = &coinharness.ChainWithMatureOutputsSpawner{
		WorkingDir:        setup.WorkingDir.Path(),
		DebugNodeOutput:   true,
		DebugWalletOutput: true,
		NumMatureOutputs:  1,
		NetPortManager:    portManager,
		WalletFactory:     simnetWalletFactory,
		NodeFactory:       nodeFactory,
		ActiveNet:         &pfcharness.Network{&chaincfg.SimNetParams},
		CreateTempWallet:  true,
		NewTestSeed:       testSeed,
		NodeStartExtraArguments: map[string]interface{}{
			"rejectnonstd": commandline.NoArgumentValue,
		},
	}

	setup.Simnet25 = &coinharness.ChainWithMatureOutputsSpawner{
		WorkingDir:        setup.WorkingDir.Path(),
		DebugNodeOutput:   true,
		DebugWalletOutput: true,
		NumMatureOutputs:  25,
		NetPortManager:    portManager,
		WalletFactory:     simnetWalletFactory,
		NodeFactory:       nodeFactory,
		ActiveNet:         &pfcharness.Network{&chaincfg.SimNetParams},
		CreateTempWallet:  true,
		NewTestSeed:       testSeed,
		NodeStartExtraArguments: map[string]interface{}{
			"rejectnonstd": commandline.NoArgumentValue,
		},
	}

	// Deploy harness spawner with empty test chain
	setup.Regnet0 = &coinharness.ChainWithMatureOutputsSpawner{
		WorkingDir:        setup.WorkingDir.Path(),
		DebugNodeOutput:   true,
		DebugWalletOutput: true,
		NumMatureOutputs:  0,
		NetPortManager:    portManager,
		WalletFactory:     regnetWalletFactory,
		NodeFactory:       nodeFactory,
		ActiveNet:         &pfcharness.Network{&chaincfg.RegNetParams},
		CreateTempWallet:  true,
		NewTestSeed:       testSeed,
	}
	// Deploy harness spawner with empty test chain
	setup.Simnet0 = &coinharness.ChainWithMatureOutputsSpawner{
		WorkingDir:        setup.WorkingDir.Path(),
		DebugNodeOutput:   true,
		DebugWalletOutput: true,
		NumMatureOutputs:  0,
		NetPortManager:    portManager,
		WalletFactory:     simnetWalletFactory,
		NodeFactory:       nodeFactory,
		ActiveNet:         &pfcharness.Network{&chaincfg.SimNetParams},
		CreateTempWallet:  true,
		NewTestSeed:       testSeed,
	}

	setup.harnessPool = pin.NewPool(setup.Regnet25)
	setup.harnessWalletPool = pin.NewPool(setup.Simnet25)

	return setup
}

func setupWorkingDir() string {
	testWorkingDir, err := ioutil.TempDir("", "integrationtest")
	if err != nil {
		fmt.Println("Unable to create working dir: ", err)
		os.Exit(-1)
	}
	return testWorkingDir
}

func setupBuild(buildName string, workingDir string, nodeProjectGoPath string) *gobuilder.GoBuider {

	tempBinDir := filepath.Join(workingDir, "bin")
	pin.MakeDirs(tempBinDir)

	nodeGoBuilder := &gobuilder.GoBuider{
		GoProjectPath:    nodeProjectGoPath,
		OutputFolderPath: tempBinDir,
		BuildFileName:    buildName,
	}
	return nodeGoBuilder
}
