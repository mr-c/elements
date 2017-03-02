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

func _SimpleAliquotMixIntoRequirements() {

}

// Conditions to run on startup
func _SimpleAliquotMixIntoSetup(_ctx context.Context, _input *SimpleAliquotMixIntoInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _SimpleAliquotMixIntoSteps(_ctx context.Context, _input *SimpleAliquotMixIntoInput, _output *SimpleAliquotMixIntoOutput) {

	for i := 0; i < _input.NumberOfAliquots; i++ {
		sampleA := mixer.Sample(_input.SampleName, _input.AliquotVolume)
		aliquot := execute.MixInto(_ctx, _input.Outplate, "", sampleA)
		_output.Aliquots = append(_output.Aliquots, aliquot)
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _SimpleAliquotMixIntoAnalysis(_ctx context.Context, _input *SimpleAliquotMixIntoInput, _output *SimpleAliquotMixIntoOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _SimpleAliquotMixIntoValidation(_ctx context.Context, _input *SimpleAliquotMixIntoInput, _output *SimpleAliquotMixIntoOutput) {

}
func _SimpleAliquotMixIntoRun(_ctx context.Context, input *SimpleAliquotMixIntoInput) *SimpleAliquotMixIntoOutput {
	output := &SimpleAliquotMixIntoOutput{}
	_SimpleAliquotMixIntoSetup(_ctx, input)
	_SimpleAliquotMixIntoSteps(_ctx, input, output)
	_SimpleAliquotMixIntoAnalysis(_ctx, input, output)
	_SimpleAliquotMixIntoValidation(_ctx, input, output)
	return output
}

func SimpleAliquotMixIntoRunSteps(_ctx context.Context, input *SimpleAliquotMixIntoInput) *SimpleAliquotMixIntoSOutput {
	soutput := &SimpleAliquotMixIntoSOutput{}
	output := _SimpleAliquotMixIntoRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SimpleAliquotMixIntoNew() interface{} {
	return &SimpleAliquotMixIntoElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SimpleAliquotMixIntoInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SimpleAliquotMixIntoRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SimpleAliquotMixIntoInput{},
			Out: &SimpleAliquotMixIntoOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type SimpleAliquotMixIntoElement struct {
	inject.CheckedRunner
}

type SimpleAliquotMixIntoInput struct {
	AliquotVolume    wunit.Volume
	NumberOfAliquots int
	Outplate         *wtype.LHPlate
	SampleName       *wtype.LHComponent
}

type SimpleAliquotMixIntoOutput struct {
	Aliquots []*wtype.LHComponent
}

type SimpleAliquotMixIntoSOutput struct {
	Data struct {
	}
	Outputs struct {
		Aliquots []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "SimpleAliquotMixInto",
		Constructor: SimpleAliquotMixIntoNew,
		Desc: component.ComponentDesc{
			Desc: "Aliquot a solution into a specified plate.\noptionally premix the solution before aliquoting\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson3_Mix_Loops/SimpleAliquot/SimpleAliquotMixInto/SimpleAliquotMixInto.an",
			Params: []component.ParamDesc{
				{Name: "AliquotVolume", Desc: "", Kind: "Parameters"},
				{Name: "NumberOfAliquots", Desc: "", Kind: "Parameters"},
				{Name: "Outplate", Desc: "", Kind: "Inputs"},
				{Name: "SampleName", Desc: "", Kind: "Inputs"},
				{Name: "Aliquots", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
