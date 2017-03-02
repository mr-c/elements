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

func _MixExercise3Requirements() {

}

// Conditions to run on startup
func _MixExercise3Setup(_ctx context.Context, _input *MixExercise3Input) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _MixExercise3Steps(_ctx context.Context, _input *MixExercise3Input, _output *MixExercise3Output) {

	// make empty slice of LHComponents (i.e. of length 0) ready to sequentially add all samples to
	// See golangbook chapter 6 for more details on slices and arrays
	dilutedsample1 := make([]*wtype.LHComponent, 0)
	dilutedsample := make([]*wtype.LHComponent, 0)

	// SampleForTotalVolume will "top up" solution to the TotalVolume with Diluent.
	// In this case it will still add diluent first but calculates the volume to add by substracting the volumes of subsequent components

	diluentsample1 := mixer.SampleForTotalVolume(_input.Diluent, _input.TotalVolume) // i.e. if TotalVolume == 20ul and SolutionVolume == 2ul then 18ul of Diluent will be sampled here
	dilutedsample1 = append(dilutedsample1, diluentsample1)
	solutionsample1 := mixer.Sample(_input.Solution, _input.DilutionVolume)
	dilutedsample1 = append(dilutedsample1, solutionsample1)
	dilutedsamplemix := execute.MixInto(_ctx, _input.OutPlateType, "", dilutedsample1...)
	dilutedSamples := append(_output.DilutedSamples, dilutedsamplemix)

	for i := 0; i < _input.NumberOfSerialDilutions; i++ {
		nextdiluentsample := mixer.SampleForTotalVolume(_input.Diluent, _input.TotalVolume)
		dilutedsample = append(dilutedsample, nextdiluentsample)
		solutionsample := mixer.Sample(dilutedsamplemix, _input.DilutionVolume)
		dilutedsample = append(dilutedsample, solutionsample)
		nextdilutedsamplemix := execute.Mix(_ctx, dilutedsample...)
		dilutedSamples = append(_output.DilutedSamples, nextdilutedsamplemix)
		dilutedsamplemix = nextdilutedsamplemix
	}

	_output.DilutedSamples = dilutedSamples

	// The Sample functions will not generate liquid handling instructions on their own
	// We need to tell Antha what to do with samples
	// For this we need to use one of the Mix functions
	// therefore finally we use Mix to combine samples into a new component

	// Now we have an antha element which will generate liquid handling instructions
	// let's see how to actually run the protocol
	// open the terminal and
	// work your way through the lessons there showing how to specify parameters and different types of workflow

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MixExercise3Analysis(_ctx context.Context, _input *MixExercise3Input, _output *MixExercise3Output) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _MixExercise3Validation(_ctx context.Context, _input *MixExercise3Input, _output *MixExercise3Output) {

}
func _MixExercise3Run(_ctx context.Context, input *MixExercise3Input) *MixExercise3Output {
	output := &MixExercise3Output{}
	_MixExercise3Setup(_ctx, input)
	_MixExercise3Steps(_ctx, input, output)
	_MixExercise3Analysis(_ctx, input, output)
	_MixExercise3Validation(_ctx, input, output)
	return output
}

func MixExercise3RunSteps(_ctx context.Context, input *MixExercise3Input) *MixExercise3SOutput {
	soutput := &MixExercise3SOutput{}
	output := _MixExercise3Run(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MixExercise3New() interface{} {
	return &MixExercise3Element{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MixExercise3Input{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MixExercise3Run(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MixExercise3Input{},
			Out: &MixExercise3Output{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type MixExercise3Element struct {
	inject.CheckedRunner
}

type MixExercise3Input struct {
	Diluent                 *wtype.LHComponent
	DilutionVolume          wunit.Volume
	NumberOfSerialDilutions int
	OutPlateType            *wtype.LHPlate
	Solution                *wtype.LHComponent
	TotalVolume             wunit.Volume
}

type MixExercise3Output struct {
	DilutedSamples []*wtype.LHComponent
}

type MixExercise3SOutput struct {
	Data struct {
	}
	Outputs struct {
		DilutedSamples []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MixExercise3",
		Constructor: MixExercise3New,
		Desc: component.ComponentDesc{
			Desc: "example protocol demonstrating the use of the SampleForTotalVolume function\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson3_Mix_Loops/MainExercise3/MixExercise3.an",
			Params: []component.ParamDesc{
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "DilutionVolume", Desc: "e.g. 2ul\n", Kind: "Parameters"},
				{Name: "NumberOfSerialDilutions", Desc: "", Kind: "Parameters"},
				{Name: "OutPlateType", Desc: "", Kind: "Inputs"},
				{Name: "Solution", Desc: "", Kind: "Inputs"},
				{Name: "TotalVolume", Desc: "e.g. 20ul\n", Kind: "Parameters"},
				{Name: "DilutedSamples", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
