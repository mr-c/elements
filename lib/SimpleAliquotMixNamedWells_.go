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

func _SimpleAliquotMixNamedWellsRequirements() {

}

// Conditions to run on startup
func _SimpleAliquotMixNamedWellsSetup(_ctx context.Context, _input *SimpleAliquotMixNamedWellsInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _SimpleAliquotMixNamedWellsSteps(_ctx context.Context, _input *SimpleAliquotMixNamedWellsInput, _output *SimpleAliquotMixNamedWellsOutput) {

	for i := 0; i < _input.NumberOfAliquots; i++ {
		sampleA := mixer.Sample(_input.SampleName, _input.AliquotVolume)
		aliquot := execute.MixInto(_ctx, _input.Outplate, _input.WellLocation[i], sampleA)
		_output.Aliquots = append(_output.Aliquots, aliquot)
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _SimpleAliquotMixNamedWellsAnalysis(_ctx context.Context, _input *SimpleAliquotMixNamedWellsInput, _output *SimpleAliquotMixNamedWellsOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _SimpleAliquotMixNamedWellsValidation(_ctx context.Context, _input *SimpleAliquotMixNamedWellsInput, _output *SimpleAliquotMixNamedWellsOutput) {

}
func _SimpleAliquotMixNamedWellsRun(_ctx context.Context, input *SimpleAliquotMixNamedWellsInput) *SimpleAliquotMixNamedWellsOutput {
	output := &SimpleAliquotMixNamedWellsOutput{}
	_SimpleAliquotMixNamedWellsSetup(_ctx, input)
	_SimpleAliquotMixNamedWellsSteps(_ctx, input, output)
	_SimpleAliquotMixNamedWellsAnalysis(_ctx, input, output)
	_SimpleAliquotMixNamedWellsValidation(_ctx, input, output)
	return output
}

func SimpleAliquotMixNamedWellsRunSteps(_ctx context.Context, input *SimpleAliquotMixNamedWellsInput) *SimpleAliquotMixNamedWellsSOutput {
	soutput := &SimpleAliquotMixNamedWellsSOutput{}
	output := _SimpleAliquotMixNamedWellsRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SimpleAliquotMixNamedWellsNew() interface{} {
	return &SimpleAliquotMixNamedWellsElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SimpleAliquotMixNamedWellsInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SimpleAliquotMixNamedWellsRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SimpleAliquotMixNamedWellsInput{},
			Out: &SimpleAliquotMixNamedWellsOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type SimpleAliquotMixNamedWellsElement struct {
	inject.CheckedRunner
}

type SimpleAliquotMixNamedWellsInput struct {
	AliquotVolume    wunit.Volume
	NumberOfAliquots int
	Outplate         *wtype.LHPlate
	PlateName        string
	SampleName       *wtype.LHComponent
	WellLocation     []string
}

type SimpleAliquotMixNamedWellsOutput struct {
	Aliquots []*wtype.LHComponent
}

type SimpleAliquotMixNamedWellsSOutput struct {
	Data struct {
	}
	Outputs struct {
		Aliquots []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "SimpleAliquotMixNamedWells",
		Constructor: SimpleAliquotMixNamedWellsNew,
		Desc: component.ComponentDesc{
			Desc: "Aliquot a solution into a specified plate.\noptionally premix the solution before aliquoting\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson3_Mix_Loops/SimpleAliquot/SimpleAliquotMixNamedWells/SimpleAliquotMixNamedWells.an",
			Params: []component.ParamDesc{
				{Name: "AliquotVolume", Desc: "", Kind: "Parameters"},
				{Name: "NumberOfAliquots", Desc: "", Kind: "Parameters"},
				{Name: "Outplate", Desc: "", Kind: "Inputs"},
				{Name: "PlateName", Desc: "", Kind: "Parameters"},
				{Name: "SampleName", Desc: "", Kind: "Inputs"},
				{Name: "WellLocation", Desc: "", Kind: "Parameters"},
				{Name: "Aliquots", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
