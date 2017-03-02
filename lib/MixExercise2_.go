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

func _MixExercise2Requirements() {

}

// Conditions to run on startup
func _MixExercise2Setup(_ctx context.Context, _input *MixExercise2Input) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _MixExercise2Steps(_ctx context.Context, _input *MixExercise2Input, _output *MixExercise2Output) {
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

		for j := 0; j < len(_input.SampleNames); j++ {
			sample := mixer.Sample(_input.SampleNames[j], _input.SampleVolumes[j])
			eachreaction = append(eachreaction, sample)
		}

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
func _MixExercise2Analysis(_ctx context.Context, _input *MixExercise2Input, _output *MixExercise2Output) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _MixExercise2Validation(_ctx context.Context, _input *MixExercise2Input, _output *MixExercise2Output) {

}
func _MixExercise2Run(_ctx context.Context, input *MixExercise2Input) *MixExercise2Output {
	output := &MixExercise2Output{}
	_MixExercise2Setup(_ctx, input)
	_MixExercise2Steps(_ctx, input, output)
	_MixExercise2Analysis(_ctx, input, output)
	_MixExercise2Validation(_ctx, input, output)
	return output
}

func MixExercise2RunSteps(_ctx context.Context, input *MixExercise2Input) *MixExercise2SOutput {
	soutput := &MixExercise2SOutput{}
	output := _MixExercise2Run(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MixExercise2New() interface{} {
	return &MixExercise2Element{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MixExercise2Input{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MixExercise2Run(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MixExercise2Input{},
			Out: &MixExercise2Output{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type MixExercise2Element struct {
	inject.CheckedRunner
}

type MixExercise2Input struct {
	ByRow             bool
	NumberOfReactions int
	Outplate          *wtype.LHPlate
	SampleNames       []*wtype.LHComponent
	SampleVolumes     []wunit.Volume
}

type MixExercise2Output struct {
	Reactions []*wtype.LHComponent
}

type MixExercise2SOutput struct {
	Data struct {
	}
	Outputs struct {
		Reactions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MixExercise2",
		Constructor: MixExercise2New,
		Desc: component.ComponentDesc{
			Desc: "Aliquot a solution into a specified plate.\noptionally premix the solution before aliquoting\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson3_Mix_Loops/MainExercise2/MixExercise2.an",
			Params: []component.ParamDesc{
				{Name: "ByRow", Desc: "", Kind: "Parameters"},
				{Name: "NumberOfReactions", Desc: "", Kind: "Parameters"},
				{Name: "Outplate", Desc: "", Kind: "Inputs"},
				{Name: "SampleNames", Desc: "", Kind: "Inputs"},
				{Name: "SampleVolumes", Desc: "", Kind: "Parameters"},
				{Name: "Reactions", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
