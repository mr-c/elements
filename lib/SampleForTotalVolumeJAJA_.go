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

func _SampleForTotalVolumeJAJARequirements() {

}

// Conditions to run on startup
func _SampleForTotalVolumeJAJASetup(_ctx context.Context, _input *SampleForTotalVolumeJAJAInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _SampleForTotalVolumeJAJASteps(_ctx context.Context, _input *SampleForTotalVolumeJAJAInput, _output *SampleForTotalVolumeJAJAOutput) {

	// Make an empty slice of LHComponents (i.e. of length 0) ready to sequentially add all samples to
	// See golangbook chapter 6 for more details on slices and arrays
	allsamples := make([]*wtype.LHComponent, 0)

	// SampleForTotalVolume will "top up" solution to the TotalVolume with Diluent.
	// In this case it will still add diluent first but calculates the volume to add by substracting the volumes of subsequent components
	diluentsample := mixer.SampleForTotalVolume(_input.Diluent, _input.TotalVolume) // i.e. if TotalVolume == 20ul and SolutionVolume == 2ul then 18ul of Diluent will be sampled here

	// append will add the diluent sample to the allsamples slice
	allsamples = append(allsamples, diluentsample)

	// Associate our SolutionVolume parameter to our Solution parameter and store the LHComponent information in our solutionsample variable
	solutionsample := mixer.Sample(_input.Solution, _input.SolutionVolume)

	// Add our solutionsample LHComponent information to our allsamples slice
	allsamples = append(allsamples, solutionsample)

	// Remember the Sample functions will not generate liquid handling instructions on their own
	// We would need to tell Antha what to do with samples
	// For this we would need to use one of the Mix functions
	_output.DilutedSample = allsamples
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _SampleForTotalVolumeJAJAAnalysis(_ctx context.Context, _input *SampleForTotalVolumeJAJAInput, _output *SampleForTotalVolumeJAJAOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _SampleForTotalVolumeJAJAValidation(_ctx context.Context, _input *SampleForTotalVolumeJAJAInput, _output *SampleForTotalVolumeJAJAOutput) {

}
func _SampleForTotalVolumeJAJARun(_ctx context.Context, input *SampleForTotalVolumeJAJAInput) *SampleForTotalVolumeJAJAOutput {
	output := &SampleForTotalVolumeJAJAOutput{}
	_SampleForTotalVolumeJAJASetup(_ctx, input)
	_SampleForTotalVolumeJAJASteps(_ctx, input, output)
	_SampleForTotalVolumeJAJAAnalysis(_ctx, input, output)
	_SampleForTotalVolumeJAJAValidation(_ctx, input, output)
	return output
}

func SampleForTotalVolumeJAJARunSteps(_ctx context.Context, input *SampleForTotalVolumeJAJAInput) *SampleForTotalVolumeJAJASOutput {
	soutput := &SampleForTotalVolumeJAJASOutput{}
	output := _SampleForTotalVolumeJAJARun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SampleForTotalVolumeJAJANew() interface{} {
	return &SampleForTotalVolumeJAJAElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SampleForTotalVolumeJAJAInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SampleForTotalVolumeJAJARun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SampleForTotalVolumeJAJAInput{},
			Out: &SampleForTotalVolumeJAJAOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type SampleForTotalVolumeJAJAElement struct {
	inject.CheckedRunner
}

type SampleForTotalVolumeJAJAInput struct {
	Diluent        *wtype.LHComponent
	Solution       *wtype.LHComponent
	SolutionVolume wunit.Volume
	TotalVolume    wunit.Volume
}

type SampleForTotalVolumeJAJAOutput struct {
	DilutedSample []*wtype.LHComponent
}

type SampleForTotalVolumeJAJASOutput struct {
	Data struct {
	}
	Outputs struct {
		DilutedSample []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "SampleForTotalVolumeJAJA",
		Constructor: SampleForTotalVolumeJAJANew,
		Desc: component.ComponentDesc{
			Desc: "example protocol demonstrating the use of the SampleForTotalVolume function\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson2_Sample/2B_SampleForTotalVolume/B_SampleForTotalVolume.an",
			Params: []component.ParamDesc{
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "Solution", Desc: "", Kind: "Inputs"},
				{Name: "SolutionVolume", Desc: "e.g. 2ul\n", Kind: "Parameters"},
				{Name: "TotalVolume", Desc: "e.g. 20ul\n", Kind: "Parameters"},
				{Name: "DilutedSample", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
