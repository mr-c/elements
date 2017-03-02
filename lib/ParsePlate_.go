// example protocol for parsing the contents of a plate from a csv file
package lib

import (
	"bytes"
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	inplate "github.com/antha-lang/antha/target/mixer"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

//Gel

// No special requirements on inputs
func _ParsePlateRequirements() {

}

// Condititions run on startup
// Including configuring an controls required, and the blocking level needed
// for them (in this case, per plate of samples processed)
func _ParsePlateSetup(_ctx context.Context, _input *ParsePlateInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _ParsePlateSteps(_ctx context.Context, _input *ParsePlateInput, _output *ParsePlateOutput) {

	// Read file
	filecontents, err := _input.InputCSVfile.ReadAll()

	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	// create a reader from filecontents
	reader := bytes.NewReader(filecontents)

	// Parse plate from reader
	plateresult, err := inplate.ParsePlateCSV(reader)

	if err != nil {
		_output.Error = err
		execute.Errorf(_ctx, err.Error())
	}

	// assign LHPlate variable
	inputplate := plateresult.Plate

	// return any warnings from plate parsing process
	_output.Warnings = plateresult.Warnings

	components := make([]*wtype.LHComponent, 0)
	_output.ComponentMap = make(map[string]*wtype.LHComponent)

	// get all plate components and return into both a slice and a map
	for _, wellcontents := range inputplate.AllWellPositions(wtype.BYCOLUMN) {

		if !inputplate.WellMap()[wellcontents].Empty() {

			component := inputplate.WellMap()[wellcontents].WContents
			components = append(components, component)
			_output.ComponentMap[component.CName] = component

		}
	}
	if len(_output.AllComponents) == 0 {
		err = fmt.Errorf("No Components found when parsing plate" + _input.InputCSVfile.Name)
		_output.Warnings = append(_output.Warnings, err.Error())
		_output.Error = err
	}
	_output.AllComponents = components
	_output.PlatewithComponents = inputplate
	execute.SetInputPlate(_ctx, _output.PlatewithComponents)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _ParsePlateAnalysis(_ctx context.Context, _input *ParsePlateInput, _output *ParsePlateOutput) {

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _ParsePlateValidation(_ctx context.Context, _input *ParsePlateInput, _output *ParsePlateOutput) {

}
func _ParsePlateRun(_ctx context.Context, input *ParsePlateInput) *ParsePlateOutput {
	output := &ParsePlateOutput{}
	_ParsePlateSetup(_ctx, input)
	_ParsePlateSteps(_ctx, input, output)
	_ParsePlateAnalysis(_ctx, input, output)
	_ParsePlateValidation(_ctx, input, output)
	return output
}

func ParsePlateRunSteps(_ctx context.Context, input *ParsePlateInput) *ParsePlateSOutput {
	soutput := &ParsePlateSOutput{}
	output := _ParsePlateRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ParsePlateNew() interface{} {
	return &ParsePlateElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ParsePlateInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ParsePlateRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ParsePlateInput{},
			Out: &ParsePlateOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type ParsePlateElement struct {
	inject.CheckedRunner
}

type ParsePlateInput struct {
	InputCSVfile wtype.File
}

type ParsePlateOutput struct {
	AllComponents       []*wtype.LHComponent
	ComponentMap        map[string]*wtype.LHComponent
	Error               error
	PlatewithComponents *wtype.LHPlate
	Warnings            []string
}

type ParsePlateSOutput struct {
	Data struct {
		Error    error
		Warnings []string
	}
	Outputs struct {
		AllComponents       []*wtype.LHComponent
		ComponentMap        map[string]*wtype.LHComponent
		PlatewithComponents *wtype.LHPlate
	}
}

func init() {
	if err := addComponent(component.Component{Name: "ParsePlate",
		Constructor: ParsePlateNew,
		Desc: component.ComponentDesc{
			Desc: "example protocol for parsing the contents of a plate from a csv file\n",
			Path: "src/github.com/antha-lang/elements/an/Utility/ParsePlate.an",
			Params: []component.ParamDesc{
				{Name: "InputCSVfile", Desc: "", Kind: "Parameters"},
				{Name: "AllComponents", Desc: "Gel\n", Kind: "Outputs"},
				{Name: "ComponentMap", Desc: "", Kind: "Outputs"},
				{Name: "Error", Desc: "", Kind: "Data"},
				{Name: "PlatewithComponents", Desc: "", Kind: "Outputs"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
