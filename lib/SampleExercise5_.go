// example protocol demonstrating the use of the SampleForTotalVolume function
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

// e.g. 2ul
// e.g. 20ul

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _SampleExercise5Requirements() {

}

// Conditions to run on startup
func _SampleExercise5Setup(_ctx context.Context, _input *SampleExercise5Input) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _SampleExercise5Steps(_ctx context.Context, _input *SampleExercise5Input, _output *SampleExercise5Output) {

	// make empty slice of LHComponents (i.e. of length 0) ready to sequentially add all samples to
	// See golangbook chapter 6 for more details on slices and arrays
	dilutedsample1 := make([]*wtype.LHComponent, 0)
	dilutedsample2 := make([]*wtype.LHComponent, 0)
	dilutedsample3 := make([]*wtype.LHComponent, 0)
	dilutedsample4 := make([]*wtype.LHComponent, 0)

	diluentsample1 := mixer.SampleForTotalVolume(_input.Diluent, _input.TotalVolume)
	dilutedsample1 = append(dilutedsample1, diluentsample1)
	solutionsample1 := mixer.Sample(_input.Solution, _input.DilutionVolume)
	dilutedsample1 = append(dilutedsample1, solutionsample1)
	dilutedsample1mix := execute.Mix(_ctx, dilutedsample1...)
	_output.DilutedSample1 = dilutedsample1mix

	diluentsample2 := mixer.SampleForTotalVolume(_input.Diluent, _input.TotalVolume)
	dilutedsample2 = append(dilutedsample2, diluentsample2)
	solutionsample2 := mixer.Sample(dilutedsample1mix, _input.DilutionVolume)
	dilutedsample2 = append(dilutedsample2, solutionsample2)
	dilutedsample2mix := execute.Mix(_ctx, dilutedsample2...)
	_output.DilutedSample2 = dilutedsample2mix

	diluentsample3 := mixer.SampleForTotalVolume(_input.Diluent, _input.TotalVolume)
	dilutedsample3 = append(dilutedsample3, diluentsample3)
	solutionsample3 := mixer.Sample(dilutedsample2mix, _input.DilutionVolume)
	dilutedsample3 = append(dilutedsample3, solutionsample3)
	dilutedsample3mix := execute.Mix(_ctx, dilutedsample3...)
	_output.DilutedSample3 = dilutedsample3mix

	diluentsample4 := mixer.SampleForTotalVolume(_input.Diluent, _input.TotalVolume)
	dilutedsample4 = append(dilutedsample4, diluentsample4)
	solutionsample4 := mixer.Sample(dilutedsample3mix, _input.DilutionVolume)
	dilutedsample4 = append(dilutedsample4, solutionsample4)
	dilutedsample4mix := execute.Mix(_ctx, dilutedsample4...)
	_output.DilutedSample4 = dilutedsample4mix

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _SampleExercise5Analysis(_ctx context.Context, _input *SampleExercise5Input, _output *SampleExercise5Output) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _SampleExercise5Validation(_ctx context.Context, _input *SampleExercise5Input, _output *SampleExercise5Output) {

}
func _SampleExercise5Run(_ctx context.Context, input *SampleExercise5Input) *SampleExercise5Output {
	output := &SampleExercise5Output{}
	_SampleExercise5Setup(_ctx, input)
	_SampleExercise5Steps(_ctx, input, output)
	_SampleExercise5Analysis(_ctx, input, output)
	_SampleExercise5Validation(_ctx, input, output)
	return output
}

func SampleExercise5RunSteps(_ctx context.Context, input *SampleExercise5Input) *SampleExercise5SOutput {
	soutput := &SampleExercise5SOutput{}
	output := _SampleExercise5Run(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SampleExercise5New() interface{} {
	return &SampleExercise5Element{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SampleExercise5Input{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SampleExercise5Run(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SampleExercise5Input{},
			Out: &SampleExercise5Output{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type SampleExercise5Element struct {
	inject.CheckedRunner
}

type SampleExercise5Input struct {
	Diluent        *wtype.LHComponent
	DilutionVolume wunit.Volume
	Solution       *wtype.LHComponent
	TotalVolume    wunit.Volume
}

type SampleExercise5Output struct {
	DilutedSample1 *wtype.LHComponent
	DilutedSample2 *wtype.LHComponent
	DilutedSample3 *wtype.LHComponent
	DilutedSample4 *wtype.LHComponent
}

type SampleExercise5SOutput struct {
	Data struct {
	}
	Outputs struct {
		DilutedSample1 *wtype.LHComponent
		DilutedSample2 *wtype.LHComponent
		DilutedSample3 *wtype.LHComponent
		DilutedSample4 *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "SampleExercise5",
		Constructor: SampleExercise5New,
		Desc: component.ComponentDesc{
			Desc: "example protocol demonstrating the use of the SampleForTotalVolume function\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson2_Sample/SampleExercise5/SampleExercise5.an",
			Params: []component.ParamDesc{
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "DilutionVolume", Desc: "e.g. 2ul\n", Kind: "Parameters"},
				{Name: "Solution", Desc: "", Kind: "Inputs"},
				{Name: "TotalVolume", Desc: "e.g. 20ul\n", Kind: "Parameters"},
				{Name: "DilutedSample1", Desc: "", Kind: "Outputs"},
				{Name: "DilutedSample2", Desc: "", Kind: "Outputs"},
				{Name: "DilutedSample3", Desc: "", Kind: "Outputs"},
				{Name: "DilutedSample4", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
