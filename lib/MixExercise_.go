// Aliquot a solution into a specified plate.
// optionally premix the solution before aliquoting
package lib

import (
	"context"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"strconv"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _MixExerciseRequirements() {

}

// Conditions to run on startup
func _MixExerciseSetup(_ctx context.Context, _input *MixExerciseInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _MixExerciseSteps(_ctx context.Context, _input *MixExerciseInput, _output *MixExerciseOutput) {
	counter := 0
	platenumber := 1
	reactions := make([]*wtype.LHComponent, 0)
	var wellpositionarray []string

	if _input.ByRow {
		for y := 0; y < _input.Outplate.WlsY; y++ {
			for x := 0; x < _input.Outplate.WlsX; x++ {
				wellposition := wutil.NumToAlpha(y+1) + strconv.Itoa(x+1)
				wellpositionarray = append(wellpositionarray, wellposition)
			}
		}
	} else {
		for x := 0; x < _input.Outplate.WlsX; x++ {
			for y := 0; y < _input.Outplate.WlsY; y++ {
				wellposition := wutil.NumToAlpha(y+1) + strconv.Itoa(x+1)
				wellpositionarray = append(wellpositionarray, wellposition)
			}
		}
	}

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
		sampleE := mixer.Sample(_input.SampleNameE, _input.SampleVolumeE)
		eachreaction = append(eachreaction, sampleE)
		sampleABCDMixreaction := execute.MixTo(_ctx, _input.Outplate.Type, wellpositionarray[counter], platenumber, eachreaction...)
		reactions = append(reactions, sampleABCDMixreaction)

		if counter+1 == len(wellpositionarray) {
			platenumber++
			counter = 0
		} else {
			counter++
		}
	}
	_output.Reactions = reactions
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MixExerciseAnalysis(_ctx context.Context, _input *MixExerciseInput, _output *MixExerciseOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _MixExerciseValidation(_ctx context.Context, _input *MixExerciseInput, _output *MixExerciseOutput) {

}
func _MixExerciseRun(_ctx context.Context, input *MixExerciseInput) *MixExerciseOutput {
	output := &MixExerciseOutput{}
	_MixExerciseSetup(_ctx, input)
	_MixExerciseSteps(_ctx, input, output)
	_MixExerciseAnalysis(_ctx, input, output)
	_MixExerciseValidation(_ctx, input, output)
	return output
}

func MixExerciseRunSteps(_ctx context.Context, input *MixExerciseInput) *MixExerciseSOutput {
	soutput := &MixExerciseSOutput{}
	output := _MixExerciseRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MixExerciseNew() interface{} {
	return &MixExerciseElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MixExerciseInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MixExerciseRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MixExerciseInput{},
			Out: &MixExerciseOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type MixExerciseElement struct {
	inject.CheckedRunner
}

type MixExerciseInput struct {
	ByRow             bool
	NumberOfReactions int
	Outplate          *wtype.LHPlate
	SampleNameA       *wtype.LHComponent
	SampleNameB       *wtype.LHComponent
	SampleNameC       *wtype.LHComponent
	SampleNameD       *wtype.LHComponent
	SampleNameE       *wtype.LHComponent
	SampleVolumeA     wunit.Volume
	SampleVolumeB     wunit.Volume
	SampleVolumeC     wunit.Volume
	SampleVolumeD     wunit.Volume
	SampleVolumeE     wunit.Volume
}

type MixExerciseOutput struct {
	Reactions []*wtype.LHComponent
}

type MixExerciseSOutput struct {
	Data struct {
	}
	Outputs struct {
		Reactions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MixExercise",
		Constructor: MixExerciseNew,
		Desc: component.ComponentDesc{
			Desc: "Aliquot a solution into a specified plate.\noptionally premix the solution before aliquoting\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson3_Mix_Loops/MainExercise1/MixExercise.an",
			Params: []component.ParamDesc{
				{Name: "ByRow", Desc: "", Kind: "Parameters"},
				{Name: "NumberOfReactions", Desc: "", Kind: "Parameters"},
				{Name: "Outplate", Desc: "", Kind: "Inputs"},
				{Name: "SampleNameA", Desc: "", Kind: "Inputs"},
				{Name: "SampleNameB", Desc: "", Kind: "Inputs"},
				{Name: "SampleNameC", Desc: "", Kind: "Inputs"},
				{Name: "SampleNameD", Desc: "", Kind: "Inputs"},
				{Name: "SampleNameE", Desc: "", Kind: "Inputs"},
				{Name: "SampleVolumeA", Desc: "", Kind: "Parameters"},
				{Name: "SampleVolumeB", Desc: "", Kind: "Parameters"},
				{Name: "SampleVolumeC", Desc: "", Kind: "Parameters"},
				{Name: "SampleVolumeD", Desc: "", Kind: "Parameters"},
				{Name: "SampleVolumeE", Desc: "", Kind: "Parameters"},
				{Name: "Reactions", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
