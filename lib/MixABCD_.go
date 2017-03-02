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

func _MixABCDRequirements() {

}

// Conditions to run on startup
func _MixABCDSetup(_ctx context.Context, _input *MixABCDInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _MixABCDSteps(_ctx context.Context, _input *MixABCDInput, _output *MixABCDOutput) {
	sampleA := mixer.Sample(_input.SampleNameA, _input.SampleVolumeA)
	sampleAMix := execute.Mix(_ctx, sampleA)
	sampleB := mixer.Sample(_input.SampleNameB, _input.SampleVolumeB)
	sampleABMix := execute.Mix(_ctx, sampleAMix, sampleB)
	sampleC := mixer.Sample(_input.SampleNameC, _input.SampleVolumeC)
	sampleABCMix := execute.Mix(_ctx, sampleABMix, sampleC)
	sampleD := mixer.Sample(_input.SampleNameD, _input.SampleVolumeD)
	sampleABCDMix := execute.Mix(_ctx, sampleABCMix, sampleD)
	_output.SampleABCDMix = sampleABCDMix
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MixABCDAnalysis(_ctx context.Context, _input *MixABCDInput, _output *MixABCDOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _MixABCDValidation(_ctx context.Context, _input *MixABCDInput, _output *MixABCDOutput) {

}
func _MixABCDRun(_ctx context.Context, input *MixABCDInput) *MixABCDOutput {
	output := &MixABCDOutput{}
	_MixABCDSetup(_ctx, input)
	_MixABCDSteps(_ctx, input, output)
	_MixABCDAnalysis(_ctx, input, output)
	_MixABCDValidation(_ctx, input, output)
	return output
}

func MixABCDRunSteps(_ctx context.Context, input *MixABCDInput) *MixABCDSOutput {
	soutput := &MixABCDSOutput{}
	output := _MixABCDRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MixABCDNew() interface{} {
	return &MixABCDElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MixABCDInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MixABCDRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MixABCDInput{},
			Out: &MixABCDOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type MixABCDElement struct {
	inject.CheckedRunner
}

type MixABCDInput struct {
	SampleNameA   *wtype.LHComponent
	SampleNameB   *wtype.LHComponent
	SampleNameC   *wtype.LHComponent
	SampleNameD   *wtype.LHComponent
	SampleVolumeA wunit.Volume
	SampleVolumeB wunit.Volume
	SampleVolumeC wunit.Volume
	SampleVolumeD wunit.Volume
}

type MixABCDOutput struct {
	SampleABCDMix *wtype.LHComponent
}

type MixABCDSOutput struct {
	Data struct {
	}
	Outputs struct {
		SampleABCDMix *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MixABCD",
		Constructor: MixABCDNew,
		Desc: component.ComponentDesc{
			Desc: "Aliquot a solution into a specified plate.\noptionally premix the solution before aliquoting\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson3_Mix_Loops/SampleABCDMix/MixABCD.an",
			Params: []component.ParamDesc{
				{Name: "SampleNameA", Desc: "", Kind: "Inputs"},
				{Name: "SampleNameB", Desc: "", Kind: "Inputs"},
				{Name: "SampleNameC", Desc: "", Kind: "Inputs"},
				{Name: "SampleNameD", Desc: "", Kind: "Inputs"},
				{Name: "SampleVolumeA", Desc: "", Kind: "Parameters"},
				{Name: "SampleVolumeB", Desc: "", Kind: "Parameters"},
				{Name: "SampleVolumeC", Desc: "", Kind: "Parameters"},
				{Name: "SampleVolumeD", Desc: "", Kind: "Parameters"},
				{Name: "SampleABCDMix", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
