package phelp_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/phelp"
	"github.com/nickwells/param.mod/v5/param/psetter"
	"github.com/nickwells/testhelper.mod/testhelper"
)

const (
	testDataDir = "testdata"
	helpSubDir  = "help"
	cfgFileDir  = "configFiles"

	paramGroupName = "test-group1"
)

var gfc = testhelper.GoldenFileCfg{
	DirNames:               []string{testDataDir, helpSubDir},
	Sfx:                    "txt",
	UpdFlagName:            "upd-help-files",
	KeepBadResultsFlagName: "keep-bad-results",
}

func init() {
	gfc.AddUpdateFlag()
	gfc.AddKeepBadResultsFlag()
}

var int64ValPos1 int64 = 101
var int64ValPos2 int64 = 102

var int64Val1 int64 = 1
var int64Val2 int64 = 2
var float64Val3 float64 = 3.333
var boolVal4 bool
var str5 string = "v1"
var str6 string = "v2"

// setInitialValues sets the parameters to their initial values - resetting
// any values overwritten by previous tests
func setInitialValues() {
	int64ValPos1 = 101
	int64ValPos2 = 102

	int64Val1 = 1
	int64Val2 = 2
	float64Val3 = 3.333
	boolVal4 = false
	str5 = "v1"
	str6 = "v2"
}

// addByPosParams will add positional parameters to the passed ParamSet
func addByPosParams(ps *param.PSet) error {
	ps.AddByPos("pos1", psetter.Int64{Value: &int64ValPos1},
		"help text for first positional parameter")
	ps.AddByPos("pos2", psetter.Int64{Value: &int64ValPos2},
		"help text for second positional parameter")

	return nil
}

// addByNameParams will add named parameters to the passed ParamSet
func addByNameParams(ps *param.PSet) error {
	ps.AddGroup(paramGroupName, "test parameters.")

	ps.Add("param1", psetter.Int64{Value: &int64Val1},
		"help text for param1",
		param.GroupName(paramGroupName),
		param.AltName("param1-alt1"),
		param.Attrs(param.CommandLineOnly),
	)

	ps.Add("param2", psetter.Int64{Value: &int64Val2},
		"help text for param2.\nWith an embedded new line and a lot of"+
			" text to demonstrate the behaviour when text is wrapped"+
			" across multiple lines",
		param.GroupName(paramGroupName),
		param.AltName("param2-alt2"),
		param.Attrs(param.MustBeSet),
	)

	ps.Add("param3", psetter.Float64{Value: &float64Val3},
		"help...",
		param.GroupName(paramGroupName),
		param.AltName("p3"),
		param.Attrs(param.DontShowInStdUsage),
	)

	ps.Add("param4", psetter.Bool{Value: &boolVal4},
		"help...",
		param.GroupName(paramGroupName),
		param.Attrs(param.SetOnlyOnce),
	)

	ps.Add("param5", psetter.Enum{
		AllowedVals: psetter.AllowedVals{
			"v1": "a value",
			"v2": "another value",
		},
		Value: &str5,
	},
		"help...",
		param.GroupName(paramGroupName),
	)

	ps.Add("param6", psetter.Enum{
		AllowedVals: psetter.AllowedVals{
			"v1": "a value",
			"v2": "another value",
		},
		Value: &str6,
	},
		"help...",
		param.GroupName(paramGroupName),
	)

	return nil
}

// configFileDetails records details about the type of config file to be set
// up for the param set
type configFileDetails struct {
	name      string
	groupName string
	strictCF  bool
	mustExist bool
}

// addConfigFiles works through the slice of config file details and adds
// them to the param set
func addConfigFiles(ps *param.PSet, configFiles []configFileDetails) {
	firstConfigFile := true
	groupConfigFile := map[string]bool{}

	for _, cfd := range configFiles {
		fc := filecheck.Optional
		if cfd.mustExist {
			fc = filecheck.MustExist
		}

		if cfd.groupName != "" {
			if groupConfigFile[cfd.groupName] {
				ps.AddGroupConfigFile(cfd.groupName, cfd.name, fc)
			} else {
				groupConfigFile[cfd.groupName] = true
				ps.SetGroupConfigFile(cfd.groupName, cfd.name, fc)
			}
		} else if cfd.strictCF {
			if firstConfigFile {
				firstConfigFile = false
				ps.SetConfigFileStrict(cfd.name, fc)
			} else {
				ps.AddConfigFileStrict(cfd.name, fc)
			}
		} else {
			if firstConfigFile {
				firstConfigFile = false
				ps.SetConfigFile(cfd.name, fc)
			} else {
				ps.AddConfigFile(cfd.name, fc)
			}
		}
	}
}

func TestHelp(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		params       []string
		progDesc     string
		configFiles  []configFileDetails
		envPrefixes  []string
		errsExpected bool
		paramAdder   []param.PSetOptFunc
	}{
		{
			ID:         testhelper.MkID("help"),
			progDesc:   "a description of what the program does (help)",
			params:     []string{"-help", "-param2=99"},
			paramAdder: []param.PSetOptFunc{addByNameParams},
		},
		{
			ID:         testhelper.MkID("help-a"),
			progDesc:   "a description of what the program does (help-a)",
			params:     []string{"-help-a", "-param2=99"},
			paramAdder: []param.PSetOptFunc{addByNameParams},
		},
		{
			ID:         testhelper.MkID("help-s"),
			progDesc:   "a description of what the program does (help-s)",
			params:     []string{"-help-s", "-param2=99"},
			paramAdder: []param.PSetOptFunc{addByNameParams},
		},
		{
			ID:       testhelper.MkID("help-params"),
			progDesc: "a description of what the program does",
			params: []string{
				"-help-params",
				"help-groups,help,help",
				"-param2=99",
			},
			paramAdder: []param.PSetOptFunc{addByNameParams},
		},
		{
			ID:       testhelper.MkID("help-params-bad-param"),
			progDesc: "a description of what the program does",
			params: []string{
				"-help-params",
				"help-groups,help,no-such-param",
				"-param2=99",
			},
			errsExpected: true,
			paramAdder:   []param.PSetOptFunc{addByNameParams},
		},
		{
			ID:       testhelper.MkID("help-params-multi-bad-param"),
			progDesc: "a description of what the program does",
			params: []string{
				"-help-params",
				"help-groups,help,no-such-param,not-a-param",
				"-param2=99",
			},
			errsExpected: true,
			paramAdder:   []param.PSetOptFunc{addByNameParams},
		},
		{
			ID:         testhelper.MkID("help-show-groups"),
			progDesc:   "a description of what the program does",
			params:     []string{"-help-show=groups", "-param2=99"},
			paramAdder: []param.PSetOptFunc{addByNameParams},
		},
		{
			ID:       testhelper.MkID("help-groups"),
			progDesc: "a description of what the program does",
			params: []string{
				"-help-groups",
				paramGroupName,
				"-param2=99"},
			paramAdder: []param.PSetOptFunc{addByNameParams},
		},
		{
			ID:       testhelper.MkID("help-groups-all"),
			progDesc: "a description of what the program does",
			params: []string{
				"-help-all",
				"-help-groups",
				paramGroupName,
				"-param2=99"},
			paramAdder: []param.PSetOptFunc{addByNameParams},
		},
		{
			ID:       testhelper.MkID("help-show-intro"),
			progDesc: "a description of what the program does",
			params: []string{
				"-help-show=intro",
				"-param2=99",
			},
			paramAdder: []param.PSetOptFunc{addByNameParams},
		},
		{
			ID:       testhelper.MkID("params-show"),
			progDesc: "a description of what the program does",
			params: []string{
				"-params-show-where-set",
				"-params-show-unused",
				"-help-s",
			},
			configFiles: []configFileDetails{
				{
					name:      filepath.Join(testDataDir, cfgFileDir, "cfg-with-param"),
					mustExist: true,
				},
			},
			errsExpected: true,
			paramAdder:   []param.PSetOptFunc{addByNameParams},
		},
		{
			ID:       testhelper.MkID("help-show-sources"),
			progDesc: "a description of what the program does",
			params: []string{
				"-help-show=sources",
				"-param2=99",
			},
			paramAdder: []param.PSetOptFunc{addByNameParams},
		},
		{
			ID:       testhelper.MkID("params-file-cmdline-param"),
			progDesc: "a description of what the program does",
			params: []string{
				"-params-file",
				"testdata/configFiles/param-cmdline.cfg",
				"-param2=99"},
			errsExpected: false,
			paramAdder:   []param.PSetOptFunc{addByNameParams},
		},
		{
			ID: testhelper.MkID("badParams-multi"),
			progDesc: "a description of what the program does" +
				" (badParams)",
			params: []string{
				"-params-file=testdata/nonesuch",
				"-help-groups=notAGroup",
				"-help-groups",
			},
			errsExpected: true,
			paramAdder:   []param.PSetOptFunc{addByNameParams},
		},
		{
			ID:       testhelper.MkID("badParams-blank-filename"),
			progDesc: "a description of what the program does",
			params: []string{
				"-params-file",
				"",
				"-param2=99"},
			errsExpected: true,
			paramAdder:   []param.PSetOptFunc{addByNameParams},
		},
		{
			ID:       testhelper.MkID("badParams-dup-filename"),
			progDesc: "a description of what the program does",
			params: []string{
				"-params-file",
				"testdata/configFiles/param.cfg",
				"-params-file",
				"testdata/configFiles/param.cfg",
				"-param2=99"},
			errsExpected: true,
			paramAdder:   []param.PSetOptFunc{addByNameParams},
		},
		{
			ID:       testhelper.MkID("badParams-groups"),
			progDesc: "a description of what the program does",
			params: []string{
				"-help-groups",
				"nonesuch1,nonesuch2,nonesuch3",
				"-param2=99"},
			errsExpected: true,
			paramAdder:   []param.PSetOptFunc{addByNameParams},
		},
		{
			ID:       testhelper.MkID("badParams-multi-style"),
			progDesc: "a description of what the program does",
			params: []string{
				"-help",
				"-help-groups",
				paramGroupName,
				"-param2=99"},
			paramAdder: []param.PSetOptFunc{addByNameParams},
		},
		{
			ID: testhelper.MkID("help-with-config"),
			progDesc: "a description of what the program" +
				" does\n\nWith embedded new lines" +
				" (help-with-config)",
			params: []string{"-help", "-param2=99"},
			configFiles: []configFileDetails{
				{
					name:      filepath.Join(testDataDir, cfgFileDir, "cfg1"),
					groupName: paramGroupName,
					mustExist: true,
				},
				{
					name:      filepath.Join(testDataDir, cfgFileDir, "cfg2"),
					groupName: paramGroupName,
				},
				{
					name:      filepath.Join(testDataDir, cfgFileDir, "cfg3"),
					mustExist: true,
				},
				{
					name: filepath.Join(testDataDir, cfgFileDir, "cfg4"),
				},
				{
					name:      filepath.Join(testDataDir, cfgFileDir, "cfg5"),
					strictCF:  true,
					mustExist: true,
				},
				{
					name:     filepath.Join(testDataDir, cfgFileDir, "cfg6"),
					strictCF: true,
				},
			},
			envPrefixes: []string{"A_", "B_", "C_"},
			paramAdder:  []param.PSetOptFunc{addByNameParams},
		},
		{
			ID: testhelper.MkID("help-with-positional-params"),
			progDesc: "a description of what the program" +
				" does\n\nWith embedded new lines" +
				" (help-with-positional-params)",
			params:     []string{"123", "456", "-help", "-param2=99"},
			paramAdder: []param.PSetOptFunc{addByNameParams, addByPosParams},
		},
	}

	for _, tc := range testCases {
		idStr := tc.IDStr()
		setInitialValues()
		helper := phelp.NewStdHelp()
		helper.SetExitAfterHelp(false)
		helper.SetDontExitOnErrors(true)

		var stdoutBuf bytes.Buffer
		var stderrBuf bytes.Buffer

		ps, err := param.NewSet(
			param.SetHelper(helper),
			param.SetStdWriter(&stdoutBuf),
			param.SetErrWriter(&stderrBuf),
			param.SetProgramDescription(tc.progDesc),
		)
		if err != nil {
			t.Log(idStr)
		}
		for _, psof := range tc.paramAdder {
			err = psof(ps)
			if err != nil {
				t.Log(idStr)
				t.Fatal("\t: Unexpected failure to add parameters:",
					err)
			}
		}
		addConfigFiles(ps, tc.configFiles)
		for _, ep := range tc.envPrefixes {
			ps.AddEnvPrefix(ep)
		}

		errMap := ps.Parse(tc.params)
		if len(errMap) != 0 {
			if !tc.errsExpected {
				t.Log(idStr)
				t.Errorf("\t: Unexpected errors: %s", stderrBuf.String())
			}
		} else if tc.errsExpected {
			t.Log(idStr)
			t.Errorf("\t: Errors expected but not seen")
		}

		gfc.Check(t, idStr+" [stdout]", tc.ID.Name+".stdout", stdoutBuf.Bytes())
		gfc.Check(t, idStr+" [stderr]", tc.ID.Name+".stderr", stderrBuf.Bytes())
	}
}

func TestHelpWithMessage(t *testing.T) {
	helper := phelp.NewStdHelp()
	helper.SetExitAfterHelp(false)
	helper.SetDontExitOnErrors(true)

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	ps, err := param.NewSet(
		param.SetHelper(helper),
		param.SetStdWriter(&stdoutBuf),
		param.SetErrWriter(&stderrBuf),
		param.SetProgramDescription("program description"))
	if err != nil {
		t.Fatal("Unexpected failure to build the parameter set:", err)
	}

	errMap := ps.Parse([]string{})
	if len(errMap) != 0 {
		t.Fatal("Unexpected errors")
	}

	ps.Help("message1", "message2")

	gfc.Check(t, t.Name()+" [stdout]", t.Name()+".stdout", stdoutBuf.Bytes())
	gfc.Check(t, t.Name()+" [stderr]", t.Name()+".stderr", stderrBuf.Bytes())
}
