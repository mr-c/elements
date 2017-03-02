// Aliquot a solution into a specified plate.
// optionally premix the solution before aliquoting
package lib

import (
	"context"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _SimpleAliquotMixToRequirements() {

}

// Conditions to run on startup
func _SimpleAliquotMixToSetup(_ctx context.Context, _input *SimpleAliquotMixToInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _SimpleAliquotMixToSteps(_ctx context.Context, _input *SimpleAliquotMixToInput, _output *SimpleAliquotMixToOutput) {

	for i := 0; i < _input.NumberOfAliquots; i++ {
		sampleA := mixer.Sample(_input.SampleName, _input.AliquotVolume)
		aliquot := execute.MixTo(_ctx, _input.Outplate, "", _input.PlateNumber, sampleA)
		_output.Aliquots = append(_output.Aliquots, aliquot)
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _SimpleAliquotMixToAnalysis(_ctx context.Context, _input *SimpleAliquotMixToInput, _output *SimpleAliquotMixToOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _SimpleAliquotMixToValidation(_ctx context.Context, _input *SimpleAliquotMixToInput, _output *SimpleAliquotMixToOutput) {

}
func _SimpleAliquotMixToRun(_ctx context.Context, input *SimpleAliquotMixToInput) *SimpleAliquotMixToOutput {
	output := &SimpleAliquotMixToOutput{}
	_SimpleAliquotMixToSetup(_ctx, input)
	_SimpleAliquotMixToSteps(_ctx, input, output)
	_SimpleAliquotMixToAnalysis(_ctx, input, output)
	_SimpleAliquotMixToValidation(_ctx, input, output)
	return output
}

func SimpleAliquotMixToRunSteps(_ctx context.Context, input *SimpleAliquotMixToInput) *SimpleAliquotMixToSOutput {
	soutput := &SimpleAliquotMixToSOutput{}
	output := _SimpleAliquotMixToRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SimpleAliquotMixToNew() interface{} {
	return &SimpleAliquotMixToElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SimpleAliquotMixToInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SimpleAliquotMixToRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SimpleAliquotMixToInput{},
			Out: &SimpleAliquotMixToOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type SimpleAliquotMixToElement struct {
	inject.CheckedRunner
}

type SimpleAliquotMixToInput struct {
	AliquotVolume    wunit.Volume
	NumberOfAliquots int
	Outplate         string
	PlateNumber      int
	SampleName       *wtype.LHComponent
}

type SimpleAliquotMixToOutput struct {
	Aliquots []*wtype.LHComponent
}

type SimpleAliquotMixToSOutput struct {
	Data struct {
	}
	Outputs struct {
		Aliquots []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "SimpleAliquotMixTo",
		Constructor: SimpleAliquotMixToNew,
		Desc: component.ComponentDesc{
			Desc: "Aliquot a solution into a specified plate.\noptionally premix the solution before aliquoting\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson3_Mix_Loops/SimpleAliquot/SimpleAliquotMixTo/SimpleAliquotMixTo/SimpleAliquotMixTo.an",
			Params: []component.ParamDesc{
				{Name: "AliquotVolume", Desc: "", Kind: "Parameters"},
				{Name: "NumberOfAliquots", Desc: "", Kind: "Parameters"},
				{Name: "Outplate", Desc: "", Kind: "Parameters"},
				{Name: "PlateNumber", Desc: "", Kind: "Parameters"},
				{Name: "SampleName", Desc: "", Kind: "Inputs"},
				{Name: "Aliquots", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
