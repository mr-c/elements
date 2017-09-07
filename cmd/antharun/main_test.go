package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	"github.com/antha-lang/antha/antha/anthalib/wtype"

	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/execute/executeutil"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/inventory/testinventory"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/human"
	"github.com/antha-lang/antha/workflowtest"
)

const (
	testTimeout = 30 * time.Second
)

func makeContext() (context.Context, error) {
	ctx := inject.NewContext(context.Background())
	for _, desc := range library {
		obj := desc.Constructor()
		runner, ok := obj.(inject.Runner)
		if !ok {
			return nil, fmt.Errorf("component %q has unexpected type %T", desc.Name, obj)
		}
		if err := inject.Add(ctx, inject.Name{Repo: desc.Name}, runner); err != nil {
			return nil, err
		}
	}

	return testinventory.NewContext(ctx), nil
}

/// relevant to copying defaults over to repos.antha.com
//var path string = anthapath.Path()

var reposPath string = "/Users/theukshowdown/go/src/repos.antha.com/antha-ninja/elements-westeros/"

var anthalangElementsExamples string = "/Users/theukshowdown/go/src/github.com/antha-lang/"

var synthaceElementsExamples string = "/Users/theukshowdown/go/src/github.com/Synthace/"

func oldRepoToNew(path string) (outPath string, err error) {

	if strings.Contains(path, "antha-lang/elements") {
		return filepath.Join(reposPath, "antha-lang-elements"), nil
	} else if strings.Contains(path, "Synthace/elements") {
		return filepath.Join(reposPath, "synthace-elements"), nil
	}

	return reposPath, fmt.Errorf("unexpected path %s", path)
}

func multiElementPath(path string) (outPath string, err error) {

	if strings.Contains(path, "antha-lang/elements") {
		return filepath.Join(reposPath, "antha-lang-elements", strings.Trim(path, anthalangElementsExamples)), nil
	} else if strings.Contains(path, "Synthace/elements") {
		return filepath.Join(reposPath, "synthace-elements", strings.Trim(path, synthaceElementsExamples)), nil
	}

	return reposPath, fmt.Errorf("unexpected path %s", path)
}

func addBytesToFiles(params map[string]map[string]json.RawMessage) map[string]map[string]json.RawMessage {
	for _, param := range params {
		for key, value := range param {
			if newFile, err := addBytes(value); err == nil {
				param[key] = newFile
			}
		}
	}
	return params
}

func addBytes(param json.RawMessage) (fileWithbytes json.RawMessage, err error) {

	var file wtype.File

	bts, err := param.MarshalJSON()

	if err != nil {
		return
	}
	err = json.Unmarshal(bts, &file)

	if err != nil {
		return
	}

	contents, err := ioutil.ReadFile(file.Name)

	if err != nil {
		return
	}

	err = file.WriteAll(contents)

	if err != nil {
		return
	}

	marshalled, err := json.Marshal(file)

	if err != nil {
		return
	}

	err = fileWithbytes.UnmarshalJSON(marshalled)

	return
}

// use this to copy bundles over into correct folders
func runTestInput(t *testing.T, ctx context.Context, input *executeutil.TestInput) {
	defer func() {
		if res := recover(); res != nil {
			t.Error(res)
		}
	}()
	tgt := target.New()
	tgt.AddDevice(human.New(human.Opt{CanIncubate: true, CanHandle: true, CanMix: true}))

	errs := make(chan error)
	go func() {
		defer close(errs)
		// HACK(ddn): Sink chdir inside goroutine to "improve" chances that
		// golang scheduler puts this goroutine on the os thread
		// corresponding to the chdir call.
		//
		// Until elements are refactored to not know their working
		// directory we can't "go test parallel" these tests
		if len(input.Dir) != 0 {
			if err := os.Chdir(input.Dir); err != nil {
				errs <- err
				return
			}
		}
		results, err := execute.Run(ctx, execute.Opt{
			Workflow: input.Workflow,
			Params:   input.Params,
			Target:   tgt,
			TransitionalReadLocalFiles: true,
		})

		// This code is added to copy bundle files as defaults
		if err == nil {
			var bundleProcesses []string
			for _, v := range input.Workflow.Processes {
				bundleProcesses = append(bundleProcesses, v.Component)
			}

			sort.Strings(bundleProcesses)

			outPath, _ := oldRepoToNew(input.Dir)

			originalPath, _ := multiElementPath(input.Dir)

			var bundleName string

			var overwriteExistingFiles bool //= true

			if len(bundleProcesses) > 1 {
				bundleName = filepath.Join(outPath, originalPath, strings.Join(bundleProcesses, "_"), strings.Join(bundleProcesses, "_")+".bundle.json")
			} else {
				bundleName = filepath.Join(outPath, bundleProcesses[0], bundleProcesses[0]+".bundle.json")
			}

			if overwriteExistingFiles {
				input.ExportBundle(bundleName)
			} else if _, err := os.Stat(bundleName); !os.IsNotExist(err) {

				fmt.Printf("file %s already exists so appending name", bundleName)
				// 9 example files permitted per process
				for i := 2; i < 10; i++ {
					if len(bundleProcesses) > 1 {
						bundleName = filepath.Join(outPath, "multi-element-bundles", strings.Join(bundleProcesses, "_"), strings.Join(bundleProcesses, "_")+fmt.Sprint("_", i)+".bundle.json")
					} else {
						bundleName = filepath.Join(outPath, bundleProcesses[0], bundleProcesses[0]+fmt.Sprint("_", i)+".bundle.json")
					}

					if _, err := os.Stat(bundleName); os.IsNotExist(err) {
						input.ExportBundle(bundleName)
						break
					}
				}

			} else {
				input.ExportBundle(bundleName)
			}

			if len(bundleProcesses) == 1 {
				parametersName := filepath.Join(outPath, bundleProcesses[0], "metadata.json")

				if _, err := os.Stat(parametersName); os.IsNotExist(err) {
					input.ExportDefaults(parametersName)
				}
			}
		}

		if err == nil && input.Expected != nil {
			err = workflowtest.CompareTestResults(results, *input.Expected)
		}

		errs <- err
	}()

	var err error

	select {
	case err = <-errs:
	case <-time.After(testTimeout):
		err = fmt.Errorf("timeout after %ds", testTimeout/time.Second)
		if inputMatches(input, string(filepath.Separator)+"long") {
			err = nil
		}
	}

	if err == nil {
	} else {
		t.Errorf("error running %s: %s", inputLabel(input), err)
	}
}

//

func inputLabel(input *executeutil.TestInput) string {
	if len(input.BundlePath) != 0 {
		return fmt.Sprintf("bundle %q", input.BundlePath)
	}
	return fmt.Sprintf("workflow %q with parameters %q", input.WorkflowPath, input.ParamsPath)
}

func inputMatches(in *executeutil.TestInput, xs ...string) bool {
	if len(xs) == 0 {
		return true
	}

	for _, x := range xs {
		for _, p := range in.Paths() {
			if strings.Contains(p, x) {
				return true
			}
		}
	}
	return false
}

func runElements(t *testing.T, ctx context.Context, inputs []*executeutil.TestInput) {
	args := flag.Args()

	odir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	for _, input := range inputs {
		in := input
		if !inputMatches(in, args...) {
			continue
		}

		t.Run(inputLabel(in), func(t *testing.T) { runTestInput(t, ctx, in) })
	}

	if err := os.Chdir(odir); err != nil {
		t.Fatal(err)
	}
}

func findInputs(basePaths ...string) ([]*executeutil.TestInput, error) {
	var inputDirs []string
	for _, c := range basePaths {
		_, err := os.Stat(c)
		if err == nil {
			inputDirs = append(inputDirs, c)
		}
	}

	if len(inputDirs) == 0 {
		return nil, fmt.Errorf("could not find example inputs in %v", basePaths)
	}

	var inputs []*executeutil.TestInput
	for _, dir := range inputDirs {
		ins, err := executeutil.FindTestInputs(dir)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, ins...)
	}

	return inputs, nil
}

var (
	anthaExamples = []string{
		"../../antha-lang/elements/an", "../../../../antha-lang/elements/an", /* "../../../../antha-lang/elements/workflows", "../../antha-lang/elements/workflows", "../../../../antha-lang/elements/defaults", "../../antha-lang/elements/defaults"*/
	}

	synthaceExamples = []string{
		"an",
		"../../examples",
		"examples",
		"../../defaults",
		"defaults",
	}
)

func TestElementsWithExampleInputs(t *testing.T) {
	flag.Parse()

	ctx, err := makeContext()
	if err != nil {
		t.Fatal(err)
	}

	inputs, err := findInputs(anthaExamples...)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("found %d test inputs\n", len(inputs))

	runElements(t, ctx, inputs)
}
