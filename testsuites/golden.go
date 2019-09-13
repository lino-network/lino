package testsuites

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/testutils"
)

var update = flag.Bool("update", false, "update .golden files")

type dumperManager interface {
	NewDumper() *testutils.Dumper
}

type GoldenTestSuite struct {
	CtxTestSuite
	dm       dumperManager
	storeKey *sdk.KVStoreKey
}

func NewGoldenTestSuite(dm dumperManager, storeKey *sdk.KVStoreKey) GoldenTestSuite {
	return GoldenTestSuite{
		dm:       dm,
		storeKey: storeKey,
	}
}

func (suite *GoldenTestSuite) LoadState(enableSubTest bool, importJsonPaths ...string) {
	// Maybe we shouldn't SetupCtx but only reset the store here.
	suite.SetupCtx(0, time.Unix(0, 0), suite.storeKey)
	if suite.dm == nil {
		suite.FailNow("Should register a dumpManager before calling LoadState")
	}
	d := suite.dm.NewDumper()
	t := suite.T()
	name := t.Name()
	parentName := name[:strings.LastIndex(name, "/")]

	var fileName string
	if enableSubTest {
		fileName = name
	} else {
		fileName = parentName
	}

	gp := filepath.Join("input", fileName+".input")
	if fileExists(gp) {
		d.LoadFromFile(suite.Ctx, gp)
	}
	for _, path := range importJsonPaths {
		gp := filepath.Join("input", "common", path+".input")
		d.LoadFromFile(suite.Ctx, gp)
	}
}

// XXX(@ryanli): Maybe we use delta for golden, better readability
func (suite *GoldenTestSuite) Golden() {
	if suite.dm == nil {
		suite.FailNow("Should register a dumpManager before calling Golden")
	}
	d := suite.dm.NewDumper()
	b := d.ToJSON(suite.Ctx)
	// Add a newline character at the end of file
	b = append(b, '\n')
	t := suite.T()
	gp := filepath.Join("golden", t.Name()+".golden")
	dir, _ := filepath.Split(gp)
	err := ensureDir(dir)
	suite.NoError(err)
	if *update {
		err := ioutil.WriteFile(gp, b, 0644)
		suite.NoError(err)
	}
	g, err := ioutil.ReadFile(gp)
	suite.NoError(err)
	suite.Equal(b, g)
}

func (suite *GoldenTestSuite) AssertStateUnchanged(enableSubTest bool, importJsonPaths ...string) {
	t := suite.T()
	name := t.Name()
	parentName := name[:strings.LastIndex(name, "/")]
	var fileName string
	if enableSubTest {
		fileName = name
	} else {
		fileName = parentName
	}
	gp := filepath.Join("input", fileName+".input")
	states := make(testutils.JSONState, 0)
	if fileExists(gp) {
		input, err := ioutil.ReadFile(gp)
		suite.NoError(err)
		err = json.Unmarshal(input, &states)
		suite.NoError(err)
	}

	for _, p := range importJsonPaths {
		tmp := make(testutils.JSONState, 0)
		gp := filepath.Join("input", "common", p+".input")
		input, err := ioutil.ReadFile(gp)
		suite.NoError(err)
		err = json.Unmarshal(input, &tmp)
		suite.NoError(err)
		states = append(states, tmp...)
	}

	// sort by prefix and then key
	sort.Slice(states, func(i, j int) bool {
		return states[i].Prefix < states[j].Prefix ||
			(states[i].Prefix == states[j].Prefix && states[i].Key < states[j].Key)
	})

	inputStates, err := json.Marshal(states)
	suite.NoError(err)

	gp = filepath.Join("golden", t.Name()+".golden")
	suite.NoError(err)
	golden, err := ioutil.ReadFile(gp)
	suite.NoError(err)
	suite.JSONEq(string(inputStates), string(golden))
}

func ensureDir(dirName string) error {
	err := os.MkdirAll(dirName, 0700)
	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
