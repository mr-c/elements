// Protocol to make a serial dilution series targeting a series of specified concentrations from a solution of known Stock concentration and a diluent.
// The next dilution in the series will always be made from the previous dilution and not from the original stock solution.
package lib

import (
	"context"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"strings"
)

// Input parameters for this protocol (data)

// specify a starting concentration
// e.g. 10 would take 1 part solution to 9 parts diluent for each dilution
// optionally choose whether to aliqout the serial dilutions by row instead of the default by column
// optionally start after a specified well position if wells are allready used in the plate

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// dilutions only
// original solution + dilutions

func _SerialDilution_forConcentrationRequirements() {

}

// Conditions to run on startup
func _SerialDilution_forConcentrationSetup(_ctx context.Context, _input *SerialDilution_forConcentrationInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _SerialDilution_forConcentrationSteps(_ctx context.Context, _input *SerialDilution_forConcentrationInput, _output *SerialDilution_forConcentrationOutput) {

	allwellpositions := _input.OutPlate.AllWellPositions(_input.ByRow)

	var counter int = _input.WellsAlreadyUsed
	var dilutionPosition int = 0

	dilutions := make([]*wtype.LHComponent, 0)

	var aliquot *wtype.LHComponent

	// calculate solution volume
	solutionVolume, err := wunit.VolumeForTargetConcentration(_input.TargetConcentrations[dilutionPosition], _input.StockConcentration, _input.StartVolumeperDilution)

	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	// use same approach to work out diluent volume to add
	diluentVolume := (wunit.CopyVolume(_input.StartVolumeperDilution))

	// this time using the substract method
	diluentVolume.Subtract(solutionVolume)

	// sample diluent
	diluentSample := mixer.Sample(_input.Diluent, diluentVolume)

	// Ensure liquid type set to Pre and Post Mix
	_input.Solution.Type = wtype.LTNeedToMix

	// sample solution
	solutionSample := mixer.Sample(_input.Solution, solutionVolume)

	// mix both samples to OutPlate
	aliquot = execute.MixNamed(_ctx, _input.OutPlate.Type, allwellpositions[counter], "DilutionPlate", diluentSample, solutionSample)

	var solutionname string

	// rename sample to include concentration
	containsconc, _, componentNameOnly := wunit.ParseConcentration(_input.Solution.CName)

	if containsconc {
		solutionname = componentNameOnly
	} else {
		solutionname = _input.Solution.CName
	}

	aliquot.CName = _input.TargetConcentrations[dilutionPosition].ToString() + " " + solutionname
	aliquot.CName = normalise(aliquot.CName)
	aliquot.SetConcentration(_input.TargetConcentrations[dilutionPosition])

	// add to dilutions array
	dilutions = append(dilutions, aliquot)

	// loop through NumberOfDilutions until all serial dilutions are made

	counter++
	dilutionPosition++

	for dilutionPosition < len(_input.TargetConcentrations) {

		// calculate new solution volume
		solutionVolume, err := wunit.VolumeForTargetConcentration(_input.TargetConcentrations[dilutionPosition], _input.TargetConcentrations[dilutionPosition-1], _input.StartVolumeperDilution)

		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}

		// use same approach to work out diluent volume to add
		diluentVolume = (wunit.CopyVolume(_input.StartVolumeperDilution))

		// this time using the substract method
		diluentVolume.Subtract(solutionVolume)

		// take next sample of diluent
		nextdiluentSample := mixer.Sample(_input.Diluent, diluentVolume)

		nextdiluentSample = execute.MixNamed(_ctx, _input.OutPlate.Type, allwellpositions[counter], "DilutionPlate", nextdiluentSample)

		// Ensure liquid type set to Pre and Post Mix
		aliquot.Type = wtype.LTNeedToMix

		// sample from previous dilution sample
		nextSample := mixer.Sample(aliquot, solutionVolume)

		// Mix sample into nextdiluent sample
		nextaliquot := execute.Mix(_ctx, nextdiluentSample, nextSample)

		// rename sample to include concentration
		nextaliquot.CName = _input.TargetConcentrations[dilutionPosition].ToString() + " " + solutionname

		nextaliquot.CName = normalise(nextaliquot.CName)

		nextaliquot.SetConcentration(_input.TargetConcentrations[dilutionPosition])
		// add to dilutions array
		dilutions = append(dilutions, nextaliquot)
		// reset aliquot
		aliquot = nextaliquot

		counter++
		dilutionPosition++
	}

	// export as Output
	_output.Dilutions = dilutions

	// export all concentrations used as export
	_output.AllDilutions = append(_output.AllDilutions, _input.Solution)
	_output.AllConcentrations = append(_output.AllConcentrations, _input.StockConcentration)

	_input.Solution.CName = _input.StockConcentration.ToString() + " " + solutionname
	_input.Solution.CName = normalise(_input.Solution.CName)
	_output.ComponentNames = append(_output.ComponentNames, _input.Solution.CName)
	for i, dilution := range _output.Dilutions {

		_output.AllDilutions = append(_output.AllDilutions, dilution)
		_output.ComponentNames = append(_output.ComponentNames, dilution.CName)
		_output.AllConcentrations = append(_output.AllConcentrations, _input.TargetConcentrations[i])

	}

	_output.WellsUsedPostRun = counter

}

// if the component name contains a concentration the concentration name will be normalised
// e.g. 10ng/ul glucose will be normalised to 10 mg/l glucose or 10mM glucose to 10 mM/l glucose or 10mM/l glucose to 10 mM/l glucose or glucose 10mM/l to 10 mM/l glucose
func normalise(name string) (normalised string) {

	if strings.Contains(name, wtype.MIXDELIMITER) {
		return name
	}

	containsConc, conc, nameonly := wunit.ParseConcentration(name)

	if containsConc {
		return conc.ToString() + " " + nameonly
	} else {
		return nameonly
	}
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _SerialDilution_forConcentrationAnalysis(_ctx context.Context, _input *SerialDilution_forConcentrationInput, _output *SerialDilution_forConcentrationOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _SerialDilution_forConcentrationValidation(_ctx context.Context, _input *SerialDilution_forConcentrationInput, _output *SerialDilution_forConcentrationOutput) {

}
func _SerialDilution_forConcentrationRun(_ctx context.Context, input *SerialDilution_forConcentrationInput) *SerialDilution_forConcentrationOutput {
	output := &SerialDilution_forConcentrationOutput{}
	_SerialDilution_forConcentrationSetup(_ctx, input)
	_SerialDilution_forConcentrationSteps(_ctx, input, output)
	_SerialDilution_forConcentrationAnalysis(_ctx, input, output)
	_SerialDilution_forConcentrationValidation(_ctx, input, output)
	return output
}

func SerialDilution_forConcentrationRunSteps(_ctx context.Context, input *SerialDilution_forConcentrationInput) *SerialDilution_forConcentrationSOutput {
	soutput := &SerialDilution_forConcentrationSOutput{}
	output := _SerialDilution_forConcentrationRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SerialDilution_forConcentrationNew() interface{} {
	return &SerialDilution_forConcentrationElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SerialDilution_forConcentrationInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SerialDilution_forConcentrationRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SerialDilution_forConcentrationInput{},
			Out: &SerialDilution_forConcentrationOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type SerialDilution_forConcentrationElement struct {
	inject.CheckedRunner
}

type SerialDilution_forConcentrationInput struct {
	ByRow                  bool
	Diluent                *wtype.LHComponent
	OutPlate               *wtype.LHPlate
	Solution               *wtype.LHComponent
	StartVolumeperDilution wunit.Volume
	StockConcentration     wunit.Concentration
	TargetConcentrations   []wunit.Concentration
	WellsAlreadyUsed       int
}

type SerialDilution_forConcentrationOutput struct {
	AllConcentrations []wunit.Concentration
	AllDilutions      []*wtype.LHComponent
	ComponentNames    []string
	Dilutions         []*wtype.LHComponent
	WellsUsedPostRun  int
}

type SerialDilution_forConcentrationSOutput struct {
	Data struct {
		AllConcentrations []wunit.Concentration
		ComponentNames    []string
		WellsUsedPostRun  int
	}
	Outputs struct {
		AllDilutions []*wtype.LHComponent
		Dilutions    []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "SerialDilution_forConcentration",
		Constructor: SerialDilution_forConcentrationNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol to make a serial dilution series targeting a series of specified concentrations from a solution of known Stock concentration and a diluent.\nThe next dilution in the series will always be made from the previous dilution and not from the original stock solution.\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/SerialDilution/SerialDilution_forConcentration.an",
			Params: []component.ParamDesc{
				{Name: "ByRow", Desc: "optionally choose whether to aliqout the serial dilutions by row instead of the default by column\n", Kind: "Parameters"},
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Solution", Desc: "", Kind: "Inputs"},
				{Name: "StartVolumeperDilution", Desc: "", Kind: "Parameters"},
				{Name: "StockConcentration", Desc: "specify a starting concentration\n", Kind: "Parameters"},
				{Name: "TargetConcentrations", Desc: "e.g. 10 would take 1 part solution to 9 parts diluent for each dilution\n", Kind: "Parameters"},
				{Name: "WellsAlreadyUsed", Desc: "optionally start after a specified well position if wells are allready used in the plate\n", Kind: "Parameters"},
				{Name: "AllConcentrations", Desc: "", Kind: "Data"},
				{Name: "AllDilutions", Desc: "original solution + dilutions\n", Kind: "Outputs"},
				{Name: "ComponentNames", Desc: "", Kind: "Data"},
				{Name: "Dilutions", Desc: "dilutions only\n", Kind: "Outputs"},
				{Name: "WellsUsedPostRun", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
