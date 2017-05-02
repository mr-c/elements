// Protocol to make a series of serial dilution sets. Each set targeting a series of specified setpoint concentrations.
// A series of input solutions are specified which must have the stock concentration specified, e.g. by NewLHComponents.
// A common diluent will be used for all.
package lib

import (
	"context"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol (data)

// specify a starting total volume per dilution, not accounting for the volume lost by using that component to make the next dilution
// a "default" may be specified which applies to all values with no explicit value set in this map

// specify target concentrations to make for each solution
// a "default" may be specified which applies to all values with no explicit value set in this map
// specify target concentrations for

// optionally choose whether to aliqout the serial dilutions by row instead of the default by column

// optionally start after a specified well position if wells are allready used in the plate

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// starting solutions. The names of the solutions will be used to set concentrations and starting volumes in the other parameters

// Use the same diluent for all component dilutions

// use the same outplate for all dilutions

// Physical outputs from this protocol with types

func _SerialDilution_Conc_multiRequirements() {

}

// Conditions to run on startup
func _SerialDilution_Conc_multiSetup(_ctx context.Context, _input *SerialDilution_Conc_multiInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _SerialDilution_Conc_multiSteps(_ctx context.Context, _input *SerialDilution_Conc_multiInput, _output *SerialDilution_Conc_multiOutput) {

	wellsused := _input.WellsAlreadyUsed
	_output.DilutionsByComponent = make(map[string][]*wtype.LHComponent)

	for _, solution := range _input.SolutionsWithConcentrations {

		var startVol wunit.Volume

		if vol, found := _input.StartVolumeperDilution[solution.CName]; found {
			startVol = vol
		} else if vol, found := _input.StartVolumeperDilution["default"]; found {
			startVol = vol
		} else {
			execute.Errorf(_ctx, "No volume specified for %s and no default volume specified", solution.CName)
		}

		var targetConcs []wunit.Concentration

		if concs, found := _input.TargetConcentrations[solution.CName]; found {
			targetConcs = concs
		} else if concs, found := _input.TargetConcentrations["default"]; found {
			targetConcs = concs
		} else {
			execute.Errorf(_ctx, "No target concentrations specified for %s and no default specified", solution.CName)
		}

		var solConc wunit.Concentration
		if solution.HasConcentration() {
			solConc = solution.Concentration()
		} else {
			execute.Errorf(_ctx, "no concentration found for %s, please set this. ", solution.CName)
		}

		// run SerialDilution_ForConcentration element
		result := SerialDilution_forConcentrationRunSteps(_ctx, &SerialDilution_forConcentrationInput{StartVolumeperDilution: startVol,
			StockConcentration:   solConc,
			TargetConcentrations: targetConcs,
			ByRow:                _input.ByRow,
			WellsAlreadyUsed:     wellsused,

			Solution: solution,
			Diluent:  _input.Diluent,
			OutPlate: _input.OutPlate},
		)

		// update wells used to carry on next set of dilutions to next available position
		wellsused = result.Data.WellsUsedPostRun

		// add all dilutions to output
		for _, dilution := range result.Outputs.Dilutions {
			_output.AllDilutions = append(_output.AllDilutions, dilution)
		}
		_output.DilutionsByComponent[solution.CName] = result.Outputs.Dilutions
	}

	_output.WellsUsedPostRun = wellsused

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _SerialDilution_Conc_multiAnalysis(_ctx context.Context, _input *SerialDilution_Conc_multiInput, _output *SerialDilution_Conc_multiOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _SerialDilution_Conc_multiValidation(_ctx context.Context, _input *SerialDilution_Conc_multiInput, _output *SerialDilution_Conc_multiOutput) {

}
func _SerialDilution_Conc_multiRun(_ctx context.Context, input *SerialDilution_Conc_multiInput) *SerialDilution_Conc_multiOutput {
	output := &SerialDilution_Conc_multiOutput{}
	_SerialDilution_Conc_multiSetup(_ctx, input)
	_SerialDilution_Conc_multiSteps(_ctx, input, output)
	_SerialDilution_Conc_multiAnalysis(_ctx, input, output)
	_SerialDilution_Conc_multiValidation(_ctx, input, output)
	return output
}

func SerialDilution_Conc_multiRunSteps(_ctx context.Context, input *SerialDilution_Conc_multiInput) *SerialDilution_Conc_multiSOutput {
	soutput := &SerialDilution_Conc_multiSOutput{}
	output := _SerialDilution_Conc_multiRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SerialDilution_Conc_multiNew() interface{} {
	return &SerialDilution_Conc_multiElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SerialDilution_Conc_multiInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SerialDilution_Conc_multiRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SerialDilution_Conc_multiInput{},
			Out: &SerialDilution_Conc_multiOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type SerialDilution_Conc_multiElement struct {
	inject.CheckedRunner
}

type SerialDilution_Conc_multiInput struct {
	ByRow                       bool
	Diluent                     *wtype.LHComponent
	OutPlate                    *wtype.LHPlate
	SolutionsWithConcentrations []*wtype.LHComponent
	StartVolumeperDilution      map[string]wunit.Volume
	TargetConcentrations        map[string][]wunit.Concentration
	WellsAlreadyUsed            int
}

type SerialDilution_Conc_multiOutput struct {
	AllDilutions         []*wtype.LHComponent
	DilutionsByComponent map[string][]*wtype.LHComponent
	WellsUsedPostRun     int
}

type SerialDilution_Conc_multiSOutput struct {
	Data struct {
		WellsUsedPostRun int
	}
	Outputs struct {
		AllDilutions         []*wtype.LHComponent
		DilutionsByComponent map[string][]*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "SerialDilution_Conc_multi",
		Constructor: SerialDilution_Conc_multiNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol to make a series of serial dilution sets. Each set targeting a series of specified setpoint concentrations.\nA series of input solutions are specified which must have the stock concentration specified, e.g. by NewLHComponents.\nA common diluent will be used for all.\n",
			Path: "src/github.com/antha-lang/elements/an/Utility/SerialDilution_Conc_multi/SerialDilution_Conc_multi.an",
			Params: []component.ParamDesc{
				{Name: "ByRow", Desc: "optionally choose whether to aliqout the serial dilutions by row instead of the default by column\n", Kind: "Parameters"},
				{Name: "Diluent", Desc: "Use the same diluent for all component dilutions\n", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "use the same outplate for all dilutions\n", Kind: "Inputs"},
				{Name: "SolutionsWithConcentrations", Desc: "starting solutions. The names of the solutions will be used to set concentrations and starting volumes in the other parameters\n", Kind: "Inputs"},
				{Name: "StartVolumeperDilution", Desc: "specify a starting total volume per dilution, not accounting for the volume lost by using that component to make the next dilution\na \"default\" may be specified which applies to all values with no explicit value set in this map\n", Kind: "Parameters"},
				{Name: "TargetConcentrations", Desc: "specify target concentrations to make for each solution\na \"default\" may be specified which applies to all values with no explicit value set in this map\n\nspecify target concentrations for\n", Kind: "Parameters"},
				{Name: "WellsAlreadyUsed", Desc: "optionally start after a specified well position if wells are allready used in the plate\n", Kind: "Parameters"},
				{Name: "AllDilutions", Desc: "", Kind: "Outputs"},
				{Name: "DilutionsByComponent", Desc: "", Kind: "Outputs"},
				{Name: "WellsUsedPostRun", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
