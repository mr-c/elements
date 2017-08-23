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

// This is the final volume that you will achieve after the dilutions have been performed.

// The dilution factor to be applied to the serial dilution, e.g. 10 would take 1 part solution to 9 parts diluent for each dilution.

// The number of dilutions you wish to make.

// An optional parameter to define whether you want your dilutions to be made in rows or columns in your plate.

// If using a plate that already has solutions in other wells you can specify from which well you want your serial dilution to start from.

// Option to remove the excess solution volume in the last dilution to the input plate if equal volumes across all dilutions are desired.

// This name will be used as the identifier for the specific plate you are making your Serial Dilutions to. If left blank it will default to DilutionPlate. If running more than one instance of the Serial Dilution element in parallel and want the dilutions to all be made on the same plate be sure this name is the same across all instance parameter sets, and be sure to wire WellsAlreadyUsed into WellsUsed between the different instances.

// Data which is returned from this protocol, and data types

// How many wells were used by this element in your output plate.

// Physical Inputs to this protocol with types

// The physical solution you wish to serially dilute, e.g. BSA, DNA, Glucose.

// The physical solution you want to make your dilutions into, e.g. water, Buffer.

//The physical plate where your serial dilutions will be made.

// Physical outputs from this protocol with types

// The physical dilutions made by this element.

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

	// Create a slice to store our liquid handling instructions for all our serial dilutions.
	dilutions := make([]*wtype.LHComponent, 0)

	// Create a variable to store the liquid handling instructions for a single dilution
	var firstDilution *wtype.LHComponent

	// calculate solution volume

	// create copy of TotalVolumeperDilution
	solutionVolume := (wunit.CopyVolume(_input.TotalVolumeperDilution))

	// use divideby method
	solutionVolume.DivideBy(_input.DilutionFactor - 1.00)

	// create a slice to hold all the volumes
	var volumes []wunit.Volume

	// add all the volumes to the slice
	volumes = append(volumes, solutionVolume)
	volumes = append(volumes, _input.TotalVolumeperDilution)

	// calculate the total volume for the intermediate dilution
	totalVolume := wunit.AddVolumes(volumes)

	// check the volume of the plate type selected for each well against the intermediate dilution total volume. If it is exceeded report an error identifying the discrepancy
	if _input.OutPlate.Welltype.MaxVolume().LessThanRounded(totalVolume, 7) {
		execute.Errorf(_ctx, "Intermediate dilution volume calculated: (%s) too high for well capacity (%s) of current plate (%s)", totalVolume.ToString(), _input.OutPlate.Welltype.MaxVolume().ToString(), _input.OutPlate.Type)
	}

	// sample diluent
	diluentSample := mixer.Sample(_input.Diluent, _input.TotalVolumeperDilution)

	// Ensure liquid type set to Pre and Post Mix
	_input.Solution.Type = wtype.LTNeedToMix

	// sample solution
	solutionSample := mixer.Sample(_input.Solution, solutionVolume)

	// Check to see if a DilutionPlateName has been specified by the user and if it hasnt then set to default "DilutionPlate"
	if _input.DilutionPlateName == "" {
		_input.DilutionPlateName = "DilutionPlate"
	}

	// mix both diluent and sample to OutPlate
	firstDilution = execute.MixNamed(_ctx, _input.OutPlate.Type, allwellpositions[counter], _input.DilutionPlateName, diluentSample, solutionSample)

	// Create a variable to store the solution name in
	var solutionname string

	// Create a variable to store the calculated concentration per solution
	var newconcentration wunit.Concentration

	// If the Stock Solution has a concentration associated with it then rename the diluted solution with its new concentration
	if _input.Solution.HasConcentration() {

		// extract the solution name as a string and store in our solution name variable
		solutionname = _input.Solution.CName

		// extract the solution concentration as a string and store in our solution concentration variable
		solutionConcentration := _input.Solution.Concentration()

		// Calculate the new concentration after dilution by dividing the solution concentration by the dilution factor
		newconcentration = wunit.DivideConcentration(solutionConcentration, float64(_input.DilutionFactor))

		// Rename the first dilution sample to contain the concentration in its name
		firstDilution.CName = newconcentration.ToString() + " " + solutionname

		// Normalise the name to a format that can be parsed for DOE elements
		firstDilution.CName = normalise(firstDilution.CName)
	}

	// add to dilutions array
	dilutions = append(dilutions, firstDilution)

	// add 1 to our counter to keep track of the number of wells that have been used
	counter++

	// loop through NumberOfDilutions until all serial dilutions are made
	for k := 1; k < _input.NumberOfDilutions; k++ {

		// take next sample of diluent
		nextDiluentSample := mixer.Sample(_input.Diluent, _input.TotalVolumeperDilution)

		// Ensure liquid type set to Pre and Post Mix
		firstDilution.Type = wtype.LTNeedToMix

		// sample from previous dilution sample
		nextSample := mixer.Sample(firstDilution, solutionVolume)

		// Mix sample into nextdiluent sample
		nextDilution := execute.MixNamed(_ctx, _input.OutPlate.Type, allwellpositions[counter], _input.DilutionPlateName, nextDiluentSample, nextSample)

		if newconcentration.RawValue() > 0 {
			// Calculate the conc entration for the next dilution based on the concentration of the previous dilution
			nextconcentration := wunit.DivideConcentration(newconcentration, float64(_input.DilutionFactor))

			// Rename the next dilution sample to show its new concentration
			nextDilution.CName = nextconcentration.ToString() + " " + solutionname

			// Normalise the name to a format that can be parsed for DOE elements
			nextDilution.CName = normalise(nextDilution.CName)

			// Set the new concentration to the next concentration calculated ready for the next round of the loop
			newconcentration = nextconcentration
		}

		// add to dilutions array
		dilutions = append(dilutions, nextDilution)

		// reset aliquot
		firstDilution = nextDilution

		// add 1 to the counter to keep track of the wells used in our output plate
		counter++
	}

	// Option to remove the excess solution volume in the last dilution to the input plate if equal volumes across all dilutions are desired.
	if _input.RemoveExcessSolution {
		// Remove the aditional solution volume from the final dilution and move it to the input plate such that the final dilution volume equals the user defined final volume.
		disposeSample := mixer.Sample(firstDilution, solutionVolume)
		execute.Mix(_ctx, disposeSample)
	}
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
	DilutionFactor         float64
	DilutionPlateName      string
	NumberOfDilutions      int
	OutPlate               *wtype.LHPlate
	RemoveExcessSolution   bool
	Solution               *wtype.LHComponent
	TotalVolumeperDilution wunit.Volume
	WellsAlreadyUsed       int
}

type SerialDilutionOutput struct {
	Dilutions []*wtype.LHComponent
	WellsUsed int
}

type SerialDilutionSOutput struct {
	Data struct {
		WellsUsed int
	}
	Outputs struct {
		Dilutions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "SerialDilution",
		Constructor: SerialDilutionNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol to make a serial dilution series from a solution and diluent\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/SerialDilution/SerialDilution.an",
			Params: []component.ParamDesc{
				{Name: "ByRow", Desc: "An optional parameter to define whether you want your dilutions to be made in rows or columns in your plate.\n", Kind: "Parameters"},
				{Name: "Diluent", Desc: "The physical solution you want to make your dilutions into, e.g. water, Buffer.\n", Kind: "Inputs"},
				{Name: "DilutionFactor", Desc: "The dilution factor to be applied to the serial dilution, e.g. 10 would take 1 part solution to 9 parts diluent for each dilution.\n", Kind: "Parameters"},
				{Name: "DilutionPlateName", Desc: "This name will be used as the identifier for the specific plate you are making your Serial Dilutions to. If left blank it will default to DilutionPlate. If running more than one instance of the Serial Dilution element in parallel and want the dilutions to all be made on the same plate be sure this name is the same across all instance parameter sets, and be sure to wire WellsAlreadyUsed into WellsUsed between the different instances.\n", Kind: "Parameters"},
				{Name: "NumberOfDilutions", Desc: "The number of dilutions you wish to make.\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "The physical plate where your serial dilutions will be made.\n", Kind: "Inputs"},
				{Name: "RemoveExcessSolution", Desc: "Option to remove the excess solution volume in the last dilution to the input plate if equal volumes across all dilutions are desired.\n", Kind: "Parameters"},
				{Name: "Solution", Desc: "The physical solution you wish to serially dilute, e.g. BSA, DNA, Glucose.\n", Kind: "Inputs"},
				{Name: "TotalVolumeperDilution", Desc: "This is the final volume that you will achieve after the dilutions have been performed.\n", Kind: "Parameters"},
				{Name: "WellsAlreadyUsed", Desc: "If using a plate that already has solutions in other wells you can specify from which well you want your serial dilution to start from.\n", Kind: "Parameters"},
				{Name: "Dilutions", Desc: "The physical dilutions made by this element.\n", Kind: "Outputs"},
				{Name: "WellsUsed", Desc: "How many wells were used by this element in your output plate.\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
