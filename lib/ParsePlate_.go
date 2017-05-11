// Protocol ParsePlate can take in a plate template file (.csv) containing liquid components (e.g. PCR reactions)
// and be used as an input element into another protocol (e.g. DNA_gel.an)
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
	"github.com/antha-lang/antha/target/mixer"
)

// Input parameters for this protocol (data)

//the file containg the plate layout

// Data which is returned from this protocol, and data types

//Error
//Warnings slice to store errors

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

//Slice of the liquid components obtained from the input file to be linked into another element
//output of the plate layout obtained from the input file
//map of all the components and volumes

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

	//read input .csv file
	filecontents, err := _input.InputCSVfile.ReadAll()

	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	//create a reader from the filecontents (contens of the .csv file)
	reader := bytes.NewReader(filecontents)

	//parse plate information from the reader
	plateresult, err := mixer.ParsePlateCSV(reader)

	if err != nil {
		_output.Error = err
		execute.Errorf(_ctx, "Errors in input plate layour .csv file: %s", err.Error())
		_output.Warnings = append(_output.Warnings, err.Error())

	}

	//assign LHPlate variable to the plate information
	inputplate := plateresult.Plate

	//return any warnings from plate parsing process and append to the Warnings slice
	if len(plateresult.Warnings) == 0 {
		plateWarnings := plateresult.Warnings
		_output.Warnings = append(_output.Warnings, plateWarnings...)
	}

	//set up a slice and a map to fill with the components
	components := make([]*wtype.LHComponent, 0)
	_output.ComponentMap = make(map[string]*wtype.LHComponent)

	//get all plate components from the well columns and return into both a slice and a map
	for _, wellcontents := range inputplate.AllWellPositions(wtype.BYCOLUMN) {

		//fill the slice and map with the components and report back errors if no compoenents found
		if !inputplate.WellMap()[wellcontents].Empty() {
			component := inputplate.WellMap()[wellcontents].WContents
			components = append(components, component)
			_output.ComponentMap[component.CName] = component
		} else {
			err = fmt.Errorf("No Components found when parsing plate: " + _input.InputCSVfile.Name)
			_output.Warnings = append(_output.Warnings, err.Error())
			_output.Error = err
		}
	}

	//update the output variables
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
			Desc: "Protocol ParsePlate can take in a plate template file (.csv) containing liquid components (e.g. PCR reactions)\nand be used as an input element into another protocol (e.g. DNA_gel.an)\n",
			Path: "src/github.com/antha-lang/elements/an/Utility/ParsePlate/ParsePlate.an",
			Params: []component.ParamDesc{
				{Name: "InputCSVfile", Desc: "the file containg the plate layout\n", Kind: "Parameters"},
				{Name: "AllComponents", Desc: "Slice of the liquid components obtained from the input file to be linked into another element\n", Kind: "Outputs"},
				{Name: "ComponentMap", Desc: "map of all the components and volumes\n", Kind: "Outputs"},
				{Name: "Error", Desc: "Error\n", Kind: "Data"},
				{Name: "PlatewithComponents", Desc: "output of the plate layout obtained from the input file\n", Kind: "Outputs"},
				{Name: "Warnings", Desc: "Warnings slice to store errors\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
