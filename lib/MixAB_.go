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

func _MixABRequirements() {

}

// Conditions to run on startup
func _MixABSetup(_ctx context.Context, _input *MixABInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _MixABSteps(_ctx context.Context, _input *MixABInput, _output *MixABOutput) {
	sampleA := mixer.Sample(_input.SampleNameA, _input.SampleVolumeA)
	sampleAMix := execute.Mix(_ctx, sampleA)
	sampleB := mixer.Sample(_input.SampleNameB, _input.SampleVolumeB)
	sampleABMix := execute.Mix(_ctx, sampleAMix, sampleB)
	_output.SampleABMix = sampleABMix

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MixABAnalysis(_ctx context.Context, _input *MixABInput, _output *MixABOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _MixABValidation(_ctx context.Context, _input *MixABInput, _output *MixABOutput) {

}
func _MixABRun(_ctx context.Context, input *MixABInput) *MixABOutput {
	output := &MixABOutput{}
	_MixABSetup(_ctx, input)
	_MixABSteps(_ctx, input, output)
	_MixABAnalysis(_ctx, input, output)
	_MixABValidation(_ctx, input, output)
	return output
}

func MixABRunSteps(_ctx context.Context, input *MixABInput) *MixABSOutput {
	soutput := &MixABSOutput{}
	output := _MixABRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MixABNew() interface{} {
	return &MixABElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MixABInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MixABRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MixABInput{},
			Out: &MixABOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type MixABElement struct {
	inject.CheckedRunner
}

type MixABInput struct {
	SampleNameA   *wtype.LHComponent
	SampleNameB   *wtype.LHComponent
	SampleVolumeA wunit.Volume
	SampleVolumeB wunit.Volume
}

type MixABOutput struct {
	SampleABMix *wtype.LHComponent
}

type MixABSOutput struct {
	Data struct {
	}
	Outputs struct {
		SampleABMix *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MixAB",
		Constructor: MixABNew,
		Desc: component.ComponentDesc{
			Desc: "Aliquot a solution into a specified plate.\noptionally premix the solution before aliquoting\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson3_Mix_Loops/SampleABMix/MixAB.an",
			Params: []component.ParamDesc{
				{Name: "SampleNameA", Desc: "", Kind: "Inputs"},
				{Name: "SampleNameB", Desc: "", Kind: "Inputs"},
				{Name: "SampleVolumeA", Desc: "", Kind: "Parameters"},
				{Name: "SampleVolumeB", Desc: "", Kind: "Parameters"},
				{Name: "SampleABMix", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
