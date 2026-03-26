package param_test

import (
	"testing"

	"github.com/nickwells/param.mod/v7/param"
	"github.com/nickwells/param.mod/v7/paramset"
	"github.com/nickwells/param.mod/v7/psetter"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestSuggestVals(t *testing.T) {
	const (
		paramName = "palbatross"
		groupName = "galbatross"
		noteName  = "nalbatross"
		trialName = "albatross"
		diffName  = "something-completely-different"
	)

	ps := paramset.NewNoHelpNoExit()

	ps.Add(paramName, psetter.Nil{}, "param desc")
	ps.Add(diffName, psetter.Nil{}, "different param desc")

	ps.AddGroup(groupName, "group desc")
	ps.AddGroup(diffName, "different group desc")

	ps.AddNote(noteName, "note text")
	ps.AddNote(diffName, "different note text")

	testCases := []struct {
		testhelper.ID
		sf        param.SuggestionFunc
		expResult []string
	}{
		{
			ID:        testhelper.MkID("parameter name alternatives"),
			sf:        param.SuggestParams,
			expResult: []string{paramName},
		},
		{
			ID:        testhelper.MkID("group name alternatives"),
			sf:        param.SuggestGroups,
			expResult: []string{groupName},
		},
		{
			ID:        testhelper.MkID("note name alternatives"),
			sf:        param.SuggestNotes,
			expResult: []string{noteName},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			results := tc.sf(ps, trialName)
			testhelper.DiffStringSlice(t,
				tc.IDStr(), "suggested vals",
				results, tc.expResult)
		})
	}
}
