package state_test

import (
	"fmt"
	"github.com/FactomProject/factomd/testHelper"
	"testing"
)

var _ = fmt.Print

func Test_StateFER(t *testing.T) {
	FEREntries := make([]testHelper.FEREntryWithHeight, 0)
	FEREntries = append(FEREntries, *testHelper.MakeFEREntryWithHeightFromContent(5, 5, 777, 5, 1))

	fmt.Println("  EntriesWithHaeight seen as: ", FEREntries)

	aState := testHelper.CreateAndPopulateTestStateForFER(FEREntries, 10)
	FER := aState.GetPredictiveFER()

	aState.ValidatorLoop()

	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!  Factoids found to be ", FER)
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!  Chain id ", aState.FERChainId)
}
