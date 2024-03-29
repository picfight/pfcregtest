package pfcregtest

import (
	"flag"
	"github.com/jfixby/coinharness"
	"github.com/jfixby/pin"
	"os"
	"testing"
)

// ObtainHarness manages access to the Pool for test cases
func ObtainHarness(tag string) *coinharness.Harness {
	s := testSetup.harnessPool.ObtainSpawnableConcurrentSafe(tag)
	return s.(*coinharness.Harness)
}

func ObtainWalletHarness(tag string) *coinharness.Harness {
	s := testSetup.harnessWalletPool.ObtainSpawnableConcurrentSafe(tag)
	return s.(*coinharness.Harness)
}

var testSetup *SimpleTestSetup

// TestMain is executed by go-test, and is
// responsible for setting up and disposing test environment.
func TestMain(m *testing.M) {
	flag.Parse()

	testSetup = Setup()

	if !testing.Short() {
		// Initialize harnesses before running any tests
		// otherwise they will be created on request.
		{
			tagsList := []string{
				//mainHarnessName,
			}
			testSetup.harnessPool.InitTags(tagsList)
		}
		{
			tagsList := []string{
				//mainWalletHarnessName,
			}
			testSetup.harnessWalletPool.InitTags(tagsList)
		}
	}

	// Run tests
	exitCode := m.Run()

	testSetup.TearDown()

	pin.VerifyNoAssetsLeaked()

	os.Exit(exitCode)
}
