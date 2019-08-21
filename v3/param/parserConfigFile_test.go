package param_test

import (
	"fmt"
	"testing"

	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/param.mod/v3/param"
	"github.com/nickwells/param.mod/v3/param/paramset"
	"github.com/nickwells/param.mod/v3/param/psetter"
)

var CFValExample1 bool
var CFValExample2 int64

func TestConfigFileStrict(t *testing.T) {
	CFValExample1 = false
	CFValExample2 = 0

	ps, err := paramset.NewNoHelpNoExitNoErrRpt(CFAddParams1, CFAddParams2)
	if err != nil {
		t.Fatal("TestConfigFile : couldn't construct the PSet: ", err)
	}
	const fname = "./testdata/config-strict.test"
	ps.SetConfigFileStrict(fname, filecheck.MustExist)

	ps.Parse([]string{})

	errs := ps.Errors()["config file: "+fname]
	if len(errs) != 0 {
		t.Logf("Unexpected error in config file:\n")
		t.Errorf("\t: %s\n", fname)
		t.Errorf("\t: got: %v\n", errs)
	}

	if CFValExample2 != 5 {
		t.Errorf("CFValExample2 should be 5 but is: %d\n", CFValExample2)
	}
}

func TestConfigFile(t *testing.T) {
	CFValExample1 = false
	CFValExample2 = 0

	ps, err := paramset.NewNoHelpNoExitNoErrRpt(CFAddParams1, CFAddParams2)
	if err != nil {
		t.Fatal("TestConfigFile : couldn't construct the PSet: ", err)
	}
	const mustExistDoes = "./testdata/config.test"
	ps.SetConfigFile(mustExistDoes, filecheck.MustExist)

	const mustExistDoesNot = "./testdata/config.nosuch"
	ps.AddConfigFile(mustExistDoesNot, filecheck.MustExist)

	const mayExistDoes = "./testdata/config.opt"
	ps.AddConfigFile(mayExistDoes, filecheck.Optional)

	const mayExistDoesNot = "./testdata/config.opt.nosuch"
	ps.AddConfigFile(mayExistDoesNot, filecheck.Optional)

	ps.Parse([]string{})

	if errs, ok := ps.Errors()["config file: "+mustExistDoes]; ok {
		t.Logf("Unexpected problem with config file:\n")
		t.Errorf("\t: %s\n", mustExistDoes)
		t.Errorf("\t: got: %v\n", errs)
	}

	if _, ok := ps.Errors()["config file: "+mustExistDoesNot]; !ok {
		t.Logf("A problem was expected with missing, must-exist config file:\n")
		t.Errorf("\t: %s\n", mustExistDoesNot)
		t.Errorf("\t: none found\n")
	}

	if errs, ok := ps.Errors()["config file: "+mayExistDoes]; ok {
		t.Logf("Unexpected problem with config file:\n")
		t.Errorf("\t: %s\n", mayExistDoes)
		t.Errorf("\t: got: %v\n", errs)
	}

	if errs, ok := ps.Errors()["config file: "+mayExistDoesNot]; ok {
		t.Logf("Unexpected problem with missing, optional config file:\n")
		t.Errorf("\t: %s\n", mayExistDoesNot)
		t.Errorf("\t: got: %v\n", errs)
	}

	if CFValExample2 != 5 {
		t.Errorf("CFValExample2 should be 5 but is: %d\n", CFValExample2)
	}
}

// CFAddParams1 will set the "example1" parameter in the PSet
func CFAddParams1(ps *param.PSet) error {
	ps.Add("example1",
		psetter.Bool{Value: &CFValExample1},
		"here is where you would describe the parameter",
		param.AltName("e1"))

	return nil
}

// CFAddParams2 will set the "example2" parameter in the PSet
func CFAddParams2(ps *param.PSet) error {
	ps.Add("example2",
		psetter.Int64{Value: &CFValExample2},
		"the description of the parameter",
		param.AltName("e2"))

	return nil
}

var groupCFName1 = "grp1"
var groupCFName2 = "grp2"
var paramInt1 int64
var paramInt2 int64
var paramBool1 bool
var paramBool2 bool

type expVals struct {
	pi1Val int64
	pi2Val int64
	pb1Val bool
	pb2Val bool
}

func TestGroupConfigFile(t *testing.T) {
	configFileNameA := "testdata/groupConfigFile.A"
	configFileNameB := "testdata/groupConfigFile.B"
	configFileNameC := "testdata/groupConfigFile.C"
	configFileNameNonesuch := "testdata/groupConfigFile.nonesuch"
	testCases := []struct {
		name         string
		gName        string
		fileName     string
		check        filecheck.Exists
		errsExpected map[string][]string
		valsExpected expVals
	}{
		{
			name:         "all good - file must exist and does",
			gName:        groupCFName1,
			fileName:     configFileNameA,
			check:        filecheck.MustExist,
			valsExpected: expVals{pi1Val: 42, pb1Val: true},
		},
		{
			name:     "config file exists but has an unknown parameter",
			gName:    groupCFName1,
			fileName: configFileNameB,
			check:    filecheck.MustExist,
			errsExpected: map[string][]string{
				"unknown-param": {
					"this is not a parameter of this program",
					"config file for " + groupCFName1,
					configFileNameB,
				},
			},
			valsExpected: expVals{pi1Val: 42, pb1Val: true},
		},
		{
			name:     "config file exists but has a parameter from another group",
			gName:    groupCFName1,
			fileName: configFileNameC,
			check:    filecheck.MustExist,
			errsExpected: map[string][]string{
				"pi2": {
					"this parameter is not a member of group: " + groupCFName1,
					"config file for " + groupCFName1,
					configFileNameC,
				},
			},
			valsExpected: expVals{pi1Val: 42, pb1Val: true},
		},
		{
			name:     "missing file",
			gName:    groupCFName1,
			fileName: configFileNameNonesuch,
			check:    filecheck.MustExist,
			errsExpected: map[string][]string{
				"config file: " + configFileNameNonesuch: {
					"no such file or directory",
					configFileNameNonesuch,
				},
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)

		ps, err := paramset.NewNoHelpNoExitNoErrRpt()
		if err != nil {
			t.Fatal(tcID, " : couldn't construct the PSet: ", err)
		}
		addParamsForGroupCF(ps)
		ps.AddGroupConfigFile(tc.gName, tc.fileName, tc.check)

		resetParamVals()
		errMap := ps.Parse([]string{})

		errMapCheck(t, tcID, errMap, tc.errsExpected)
		valsCheck(t, tcID, tc.valsExpected)
	}

}

// resetParamVals resets the param values to their initial state
func resetParamVals() {
	paramInt1 = 0
	paramInt2 = 0
	paramBool1 = false
	paramBool2 = false
}

// valsCheck checks that the values match the expected values
func valsCheck(t *testing.T, testID string, vals expVals) {
	t.Helper()

	var nameLogged bool
	if paramInt1 != vals.pi1Val {
		nameLogged = logName(t, nameLogged, testID)
		t.Errorf("\t: unexpected values: paramInt1 = %d, should be %d\n",
			paramInt1, vals.pi1Val)
	}

	if paramInt2 != vals.pi2Val {
		nameLogged = logName(t, nameLogged, testID)
		t.Errorf("\t: unexpected values: paramInt2 = %d, should be %d\n",
			paramInt2, vals.pi2Val)
	}

	if paramBool1 != vals.pb1Val {
		nameLogged = logName(t, nameLogged, testID)
		t.Errorf("\t: unexpected values: paramBool1 = %v, should be %v\n",
			paramBool1, vals.pb1Val)
	}

	if paramBool2 != vals.pb2Val {
		logName(t, nameLogged, testID)
		t.Errorf("\t: unexpected values: paramBool2 = %v, should be %v\n",
			paramBool2, vals.pb2Val)
	}
}

func addParamsForGroupCF(ps *param.PSet) {
	ps.SetGroupDescription(groupCFName1, "blah blah blah - 1")
	ps.SetGroupDescription(groupCFName2, "blah blah blah - 2")
	ps.Add("pi1", psetter.Int64{Value: &paramInt1},
		"param int val 1",
		param.GroupName(groupCFName1))
	ps.Add("pi2", psetter.Int64{Value: &paramInt2},
		"param int val 2",
		param.GroupName(groupCFName2))
	ps.Add("pb1", psetter.Bool{Value: &paramBool1},
		"param bool val 1",
		param.GroupName(groupCFName1))
	ps.Add("pb2", psetter.Bool{Value: &paramBool2},
		"param bool val 2",
		param.GroupName(groupCFName2))
}
