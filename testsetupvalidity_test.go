package pfcregtest

import (
	"fmt"
	"github.com/picfight/pfcd/dcrutil"
	"testing"
)

func TestSetupValidity(t *testing.T) {
	coins50 := dcrutil.Amount(50 /*PFC*/ * 1e8)
	stringVal := fmt.Sprintf("%v", coins50)
	expectedStringVal := "50 PFC"
	//pin.D("stringVal", stringVal)
	if expectedStringVal != stringVal {
		t.Fatalf("Incorrect coin: "+
			"expected %v, got %v", expectedStringVal, stringVal)
	}
}
