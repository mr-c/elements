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

func _MixIntoABCDRequirements() {

}

// Conditions to run on startup
func _MixIntoABCDSetup(_ctx context.Context, _input *MixIntoABCDInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _MixIntoABCDSteps(_ctx context.Context, _input *MixIntoABCDInput, _output *MixIntoABCDOutput) {
	sampleA := mixer.Sample(_input.SampleNameA, _input.SampleVolumeA)
	sampleAMix := execute.MixInto(_ctx, _input.Outplate, "A1", sampleA)
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
func _MixIntoABCDAnalysis(_ctx context.Context, _input *MixIntoABCDInput, _output *MixIntoABCDOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _MixIntoABCDValidation(_ctx context.Context, _input *MixIntoABCDInput, _output *MixIntoABCDOutput) {

}
func _MixIntoABCDRun(_ctx context.Context, input *MixIntoABCDInput) *MixIntoABCDOutput {
	output := &MixIntoABCDOutput{}
	_MixIntoABCDSetup(_ctx, input)
	_MixIntoABCDSteps(_ctx, input, output)
	_MixIntoABCDAnalysis(_ctx, input, output)
	_MixIntoABCDValidation(_ctx, input, output)
	return output
}

func MixIntoABCDRunSteps(_ctx context.Context, input *MixIntoABCDInput) *MixIntoABCDSOutput {
	soutput := &MixIntoABCDSOutput{}
	output := _MixIntoABCDRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MixIntoABCDNew() interface{} {
	return &MixIntoABCDElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MixIntoABCDInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MixIntoABCDRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MixIntoABCDInput{},
			Out: &MixIntoABCDOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type MixIntoABCDElement struct {
	inject.CheckedRunner
}

type MixIntoABCDInput struct {
	Outplate      *wtype.LHPlate
	SampleNameA   *wtype.LHComponent
	SampleNameB   *wtype.LHComponent
	SampleNameC   *wtype.LHComponent
	SampleNameD   *wtype.LHComponent
	SampleVolumeA wunit.Volume
	SampleVolumeB wunit.Volume
	SampleVolumeC wunit.Volume
	SampleVolumeD wunit.Volume
}

type MixIntoABCDOutput struct {
	SampleABCDMix *wtype.LHComponent
}

type MixIntoABCDSOutput struct {
	Data struct {
	}
	Outputs struct {
		SampleABCDMix *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MixIntoABCD",
		Constructor: MixIntoABCDNew,
		Desc: component.ComponentDesc{
			Desc: "Aliquot a solution into a specified plate.\noptionally premix the solution before aliquoting\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson3_Mix_Loops/SampleABCDMixInto/MixIntoABCD/MixIntoABCD.an",
			Params: []component.ParamDesc{
				{Name: "Outplate", Desc: "", Kind: "Inputs"},
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
