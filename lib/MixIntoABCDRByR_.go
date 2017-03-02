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

func _MixIntoABCDRByRRequirements() {

}

// Conditions to run on startup
func _MixIntoABCDRByRSetup(_ctx context.Context, _input *MixIntoABCDRByRInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _MixIntoABCDRByRSteps(_ctx context.Context, _input *MixIntoABCDRByRInput, _output *MixIntoABCDRByROutput) {

	sampleABCDMixreactions := make([]*wtype.LHComponent, 0)
	for i := 0; i < _input.NumberOfReactions; i++ {
		eachreaction := make([]*wtype.LHComponent, 0)
		sampleA := mixer.Sample(_input.SampleNameA, _input.SampleVolumeA)
		eachreaction = append(eachreaction, sampleA)
		sampleB := mixer.Sample(_input.SampleNameB, _input.SampleVolumeB)
		eachreaction = append(eachreaction, sampleB)
		sampleC := mixer.Sample(_input.SampleNameC, _input.SampleVolumeC)
		eachreaction = append(eachreaction, sampleC)
		sampleD := mixer.Sample(_input.SampleNameD, _input.SampleVolumeD)
		eachreaction = append(eachreaction, sampleD)
		sampleABCDMixreaction := execute.MixInto(_ctx, _input.Outplate, "", eachreaction...)
		sampleABCDMixreactions = append(sampleABCDMixreactions, sampleABCDMixreaction)
	}
	_output.SampleABCDMixreactions = sampleABCDMixreactions
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MixIntoABCDRByRAnalysis(_ctx context.Context, _input *MixIntoABCDRByRInput, _output *MixIntoABCDRByROutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _MixIntoABCDRByRValidation(_ctx context.Context, _input *MixIntoABCDRByRInput, _output *MixIntoABCDRByROutput) {

}
func _MixIntoABCDRByRRun(_ctx context.Context, input *MixIntoABCDRByRInput) *MixIntoABCDRByROutput {
	output := &MixIntoABCDRByROutput{}
	_MixIntoABCDRByRSetup(_ctx, input)
	_MixIntoABCDRByRSteps(_ctx, input, output)
	_MixIntoABCDRByRAnalysis(_ctx, input, output)
	_MixIntoABCDRByRValidation(_ctx, input, output)
	return output
}

func MixIntoABCDRByRRunSteps(_ctx context.Context, input *MixIntoABCDRByRInput) *MixIntoABCDRByRSOutput {
	soutput := &MixIntoABCDRByRSOutput{}
	output := _MixIntoABCDRByRRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MixIntoABCDRByRNew() interface{} {
	return &MixIntoABCDRByRElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MixIntoABCDRByRInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MixIntoABCDRByRRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MixIntoABCDRByRInput{},
			Out: &MixIntoABCDRByROutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type MixIntoABCDRByRElement struct {
	inject.CheckedRunner
}

type MixIntoABCDRByRInput struct {
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

type MixIntoABCDRByROutput struct {
	SampleABCDMixreactions []*wtype.LHComponent
}

type MixIntoABCDRByRSOutput struct {
	Data struct {
	}
	Outputs struct {
		SampleABCDMixreactions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MixIntoABCDRByR",
		Constructor: MixIntoABCDRByRNew,
		Desc: component.ComponentDesc{
			Desc: "Aliquot a solution into a specified plate.\noptionally premix the solution before aliquoting\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson3_Mix_Loops/SampleABCDMixInto/MixIntoABCDReactionByReaction/MixIntoABCDRByR.an",
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
				{Name: "SampleABCDMixreactions", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
