// Protocol to make a serial dilution series from a solution and diluent
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

// e.g. 10 would take 1 part solution to 9 parts diluent for each dilution

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _SerialDilutionRequirements() {

}

// Conditions to run on startup
func _SerialDilutionSetup(_ctx context.Context, _input *SerialDilutionInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _SerialDilutionSteps(_ctx context.Context, _input *SerialDilutionInput, _output *SerialDilutionOutput) {

	// This code allows the user to specify how the Serial Dilutions should be made in order, by row or by column.

	allwellpositions := _input.OutPlate.AllWellPositions(_input.ByRow)
	var counter int = _input.WellsAlreadyUsed

	dilutions := make([]*wtype.LHComponent, 0)

	var aliquot *wtype.LHComponent

	// calculate solution volume

	// create copy of TotalVolumeperDilution
	solutionVolume := (wunit.CopyVolume(_input.TotalVolumeperDilution))

	// use divideby method
	solutionVolume.DivideBy(float64(_input.DilutionFactor) - 1.00)

	// sample diluent
	diluentSample := mixer.Sample(_input.Diluent, _input.TotalVolumeperDilution)

	// Ensure liquid type set to Pre and Post Mix
	_input.Solution.Type = wtype.LTNeedToMix

	// sample solution
	solutionSample := mixer.Sample(_input.Solution, solutionVolume)

	// mix both diluent and sample to OutPlate
	aliquot = execute.MixNamed(_ctx, _input.OutPlate.Type, allwellpositions[counter], "DilutionPlate", diluentSample, solutionSample)

	// add to dilutions array
	dilutions = append(dilutions, aliquot)

	counter++
	// loop through NumberOfDilutions until all serial dilutions are made
	for k := 1; k < _input.NumberOfDilutions; k++ {

		// take next sample of diluent
		nextDiluentSample := mixer.Sample(_input.Diluent, _input.TotalVolumeperDilution)

		// Ensure liquid type set to Pre and Post Mix
		aliquot.Type = wtype.LTNeedToMix

		// sample from previous dilution sample
		nextSample := mixer.Sample(aliquot, solutionVolume)

		// Mix sample into nextdiluent sample
		nextaliquot := execute.MixNamed(_ctx, _input.OutPlate.Type, allwellpositions[counter], "DilutionPlate", nextDiluentSample, nextSample)

		// add to dilutions array
		dilutions = append(dilutions, nextaliquot)
		// reset aliquot
		aliquot = nextaliquot
		counter++
	}

	disposeSample := mixer.Sample(aliquot, solutionVolume)
	_output.Waste = execute.MixNamed(_ctx, "reservoir", "", "SolutionWaste", disposeSample)

	// export as Output
	_output.Dilutions = dilutions

	// Output the number of wells that have been used on this plate
	_output.WellsUsed = counter

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _SerialDilutionAnalysis(_ctx context.Context, _input *SerialDilutionInput, _output *SerialDilutionOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _SerialDilutionValidation(_ctx context.Context, _input *SerialDilutionInput, _output *SerialDilutionOutput) {

}
func _SerialDilutionRun(_ctx context.Context, input *SerialDilutionInput) *SerialDilutionOutput {
	output := &SerialDilutionOutput{}
	_SerialDilutionSetup(_ctx, input)
	_SerialDilutionSteps(_ctx, input, output)
	_SerialDilutionAnalysis(_ctx, input, output)
	_SerialDilutionValidation(_ctx, input, output)
	return output
}

func SerialDilutionRunSteps(_ctx context.Context, input *SerialDilutionInput) *SerialDilutionSOutput {
	soutput := &SerialDilutionSOutput{}
	output := _SerialDilutionRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SerialDilutionNew() interface{} {
	return &SerialDilutionElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SerialDilutionInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SerialDilutionRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SerialDilutionInput{},
			Out: &SerialDilutionOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type SerialDilutionElement struct {
	inject.CheckedRunner
}

type SerialDilutionInput struct {
	ByRow                  bool
	Diluent                *wtype.LHComponent
	DilutionFactor         int
	NumberOfDilutions      int
	OutPlate               *wtype.LHPlate
	Solution               *wtype.LHComponent
	TotalVolumeperDilution wunit.Volume
	WellsAlreadyUsed       int
}

type SerialDilutionOutput struct {
	Dilutions []*wtype.LHComponent
	Waste     *wtype.LHComponent
	WellsUsed int
}

type SerialDilutionSOutput struct {
	Data struct {
		WellsUsed int
	}
	Outputs struct {
		Dilutions []*wtype.LHComponent
		Waste     *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "SerialDilution",
		Constructor: SerialDilutionNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol to make a serial dilution series from a solution and diluent\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/SerialDilution/SerialDilution.an",
			Params: []component.ParamDesc{
				{Name: "ByRow", Desc: "", Kind: "Parameters"},
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "DilutionFactor", Desc: "e.g. 10 would take 1 part solution to 9 parts diluent for each dilution\n", Kind: "Parameters"},
				{Name: "NumberOfDilutions", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Solution", Desc: "", Kind: "Inputs"},
				{Name: "TotalVolumeperDilution", Desc: "", Kind: "Parameters"},
				{Name: "WellsAlreadyUsed", Desc: "", Kind: "Parameters"},
				{Name: "Dilutions", Desc: "", Kind: "Outputs"},
				{Name: "Waste", Desc: "", Kind: "Outputs"},
				{Name: "WellsUsed", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
