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

func _MixIntoABCDCByCRequirements() {

}

// Conditions to run on startup
func _MixIntoABCDCByCSetup(_ctx context.Context, _input *MixIntoABCDCByCInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _MixIntoABCDCByCSteps(_ctx context.Context, _input *MixIntoABCDCByCInput, _output *MixIntoABCDCByCOutput) {
	var reaction []*wtype.LHComponent
	for i := 0; i < _input.NumberOfReactions; i++ {
		sampleA := mixer.Sample(_input.SampleNameA, _input.SampleVolumeA)
		sampleAMix := execute.MixInto(_ctx, _input.Outplate, "", sampleA)
		sampleB := mixer.Sample(_input.SampleNameB, _input.SampleVolumeB)
		sampleABMix := execute.Mix(_ctx, sampleAMix, sampleB)
		sampleC := mixer.Sample(_input.SampleNameC, _input.SampleVolumeC)
		sampleABCMix := execute.Mix(_ctx, sampleABMix, sampleC)
		sampleD := mixer.Sample(_input.SampleNameD, _input.SampleVolumeD)
		sampleABCDMix := execute.Mix(_ctx, sampleABCMix, sampleD)
		reaction = append(reaction, sampleABCDMix)
	}
	_output.Reactions = reaction
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MixIntoABCDCByCAnalysis(_ctx context.Context, _input *MixIntoABCDCByCInput, _output *MixIntoABCDCByCOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _MixIntoABCDCByCValidation(_ctx context.Context, _input *MixIntoABCDCByCInput, _output *MixIntoABCDCByCOutput) {

}
func _MixIntoABCDCByCRun(_ctx context.Context, input *MixIntoABCDCByCInput) *MixIntoABCDCByCOutput {
	output := &MixIntoABCDCByCOutput{}
	_MixIntoABCDCByCSetup(_ctx, input)
	_MixIntoABCDCByCSteps(_ctx, input, output)
	_MixIntoABCDCByCAnalysis(_ctx, input, output)
	_MixIntoABCDCByCValidation(_ctx, input, output)
	return output
}

func MixIntoABCDCByCRunSteps(_ctx context.Context, input *MixIntoABCDCByCInput) *MixIntoABCDCByCSOutput {
	soutput := &MixIntoABCDCByCSOutput{}
	output := _MixIntoABCDCByCRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MixIntoABCDCByCNew() interface{} {
	return &MixIntoABCDCByCElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MixIntoABCDCByCInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MixIntoABCDCByCRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MixIntoABCDCByCInput{},
			Out: &MixIntoABCDCByCOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type MixIntoABCDCByCElement struct {
	inject.CheckedRunner
}

type MixIntoABCDCByCInput struct {
	NumberOfReactions int
	Outplate          *wtype.LHPlate
	SampleNameA       *wtype.LHComponent
	SampleNameB       *wtype.LHComponent
	SampleNameC       *wtype.LHComponent
	SampleNameD       *wtype.LHComponent
	SampleVolumeA     wunit.Volume
	SampleVolumeB     wunit.Volume
	SampleVolumeC     wunit.Volume
	SampleVolumeD     wunit.Volume
}

type MixIntoABCDCByCOutput struct {
	Reactions []*wtype.LHComponent
}

type MixIntoABCDCByCSOutput struct {
	Data struct {
	}
	Outputs struct {
		Reactions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MixIntoABCDCByC",
		Constructor: MixIntoABCDCByCNew,
		Desc: component.ComponentDesc{
			Desc: "Aliquot a solution into a specified plate.\noptionally premix the solution before aliquoting\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson3_Mix_Loops/SampleABCDMixInto/MixIntoABCDComponentByComponent/MixIntoABCDCByC.an",
			Params: []component.ParamDesc{
				{Name: "NumberOfReactions", Desc: "", Kind: "Parameters"},
				{Name: "Outplate", Desc: "", Kind: "Inputs"},
				{Name: "SampleNameA", Desc: "", Kind: "Inputs"},
				{Name: "SampleNameB", Desc: "", Kind: "Inputs"},
				{Name: "SampleNameC", Desc: "", Kind: "Inputs"},
				{Name: "SampleNameD", Desc: "", Kind: "Inputs"},
				{Name: "SampleVolumeA", Desc: "", Kind: "Parameters"},
				{Name: "SampleVolumeB", Desc: "", Kind: "Parameters"},
				{Name: "SampleVolumeC", Desc: "", Kind: "Parameters"},
				{Name: "SampleVolumeD", Desc: "", Kind: "Parameters"},
				{Name: "Reactions", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
