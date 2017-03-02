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

func _MixIntoABRequirements() {

}

// Conditions to run on startup
func _MixIntoABSetup(_ctx context.Context, _input *MixIntoABInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _MixIntoABSteps(_ctx context.Context, _input *MixIntoABInput, _output *MixIntoABOutput) {
	sampleA := mixer.Sample(_input.SampleNameA, _input.SampleVolumeA)
	sampleAMix := execute.MixInto(_ctx, _input.Outplate, "A1", sampleA)
	sampleB := mixer.Sample(_input.SampleNameB, _input.SampleVolumeB)
	sampleABMix := execute.Mix(_ctx, sampleAMix, sampleB)
	_output.SampleABMix = sampleABMix

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MixIntoABAnalysis(_ctx context.Context, _input *MixIntoABInput, _output *MixIntoABOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _MixIntoABValidation(_ctx context.Context, _input *MixIntoABInput, _output *MixIntoABOutput) {

}
func _MixIntoABRun(_ctx context.Context, input *MixIntoABInput) *MixIntoABOutput {
	output := &MixIntoABOutput{}
	_MixIntoABSetup(_ctx, input)
	_MixIntoABSteps(_ctx, input, output)
	_MixIntoABAnalysis(_ctx, input, output)
	_MixIntoABValidation(_ctx, input, output)
	return output
}

func MixIntoABRunSteps(_ctx context.Context, input *MixIntoABInput) *MixIntoABSOutput {
	soutput := &MixIntoABSOutput{}
	output := _MixIntoABRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MixIntoABNew() interface{} {
	return &MixIntoABElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MixIntoABInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MixIntoABRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MixIntoABInput{},
			Out: &MixIntoABOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type MixIntoABElement struct {
	inject.CheckedRunner
}

type MixIntoABInput struct {
	Outplate      *wtype.LHPlate
	SampleNameA   *wtype.LHComponent
	SampleNameB   *wtype.LHComponent
	SampleVolumeA wunit.Volume
	SampleVolumeB wunit.Volume
}

type MixIntoABOutput struct {
	SampleABMix *wtype.LHComponent
}

type MixIntoABSOutput struct {
	Data struct {
	}
	Outputs struct {
		SampleABMix *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MixIntoAB",
		Constructor: MixIntoABNew,
		Desc: component.ComponentDesc{
			Desc: "Aliquot a solution into a specified plate.\noptionally premix the solution before aliquoting\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson3_Mix_Loops/SampleABMixInto/MixIntoAB.an",
			Params: []component.ParamDesc{
				{Name: "Outplate", Desc: "", Kind: "Inputs"},
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
