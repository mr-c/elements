// Perform multiple PCR reactions with common default parameters using a mastermix
package lib

import (
	"context"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

//"github.com/antha-lang/antha/antha/anthalib/thermocycle"

// Input parameters for this protocol (data)

// map of which reaction uses which template e.g. ["left homology arm"]:"templatename"
// map of which reaction uses which primer pair e.g. ["left homology arm"]:"fwdprimer","revprimer"
// Volume of template in each reaction
// e.g. for  10X Q5 buffer this would be 10

// Data which is returned from this protocol, and data types

// return an error message if an error is encountered

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _AutoPCR_mmx_multiRequirements() {
}

// Conditions to run on startup
func _AutoPCR_mmx_multiSetup(_ctx context.Context, _input *AutoPCR_mmx_multiInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _AutoPCR_mmx_multiSteps(_ctx context.Context, _input *AutoPCR_mmx_multiInput, _output *AutoPCR_mmx_multiOutput) {

	// set up a counter to use as an index for increasing well position
	var counter int

	// set up some empty slices to fill as we iterate through the reactions
	_output.Reactions = make([]*wtype.LHComponent, 0)
	volumes := make([]wunit.Volume, 0)
	welllocations := make([]string, 0)
	// initialise map
	_output.ReactionMap = make(map[string]*wtype.LHComponent)

	// To allow using defaults in either the template map or the primer map,
	// evaluate which is longer and make a slice or reactions to range through
	var reactions []string

	if len(_input.Reactiontotemplate) >= len(_input.Reactiontoprimerpair) {
		for reactionname := range _input.Reactiontotemplate {
			reactions = append(reactions, reactionname)
		}
	} else {
		for reactionname := range _input.Reactiontoprimerpair {
			reactions = append(reactions, reactionname)
		}
	}

	for _, reactionname := range reactions {

		// look up template from map
		var template string

		if templateName, found := _input.Reactiontotemplate[reactionname]; found {
			template = templateName
		} else if templateName, found := _input.Reactiontotemplate["default"]; found {
			template = templateName
		} else {
			execute.Errorf(_ctx, `No template set for %s and no "default" primers set`, reactionname)
		}

		// look up primers from map
		var fwdPrimer string
		var revPrimer string

		if primers, found := _input.Reactiontoprimerpair[reactionname]; found {
			fwdPrimer, revPrimer = primers[0], primers[1]
		} else if primers, found := _input.Reactiontoprimerpair["default"]; found {
			fwdPrimer, revPrimer = primers[0], primers[1]
		} else {
			execute.Errorf(_ctx, `No primers set for %s and no "default" primers set`, reactionname)
		}

		// use counter to find next available well position in plate

		var allwellpositionsforplate []string

		allwellpositionsforplate = _input.Plate.AllWellPositions(wtype.BYCOLUMN)

		wellposition := allwellpositionsforplate[counter]

		// handle to set up thermocycler
		//MasterMix = Handle(thermocycle.SetUp(MasterMix))

		// Run PCR_vol element
		result := PCR_vol_mmxRunSteps(_ctx, &PCR_vol_mmxInput{MasterMixVolume: _input.DefaultMasterMixVolume,
			PrimersalreadyAddedtoMasterMix:    _input.PrimersalreadyAddedtoMasterMix,
			PolymeraseAlreadyaddedtoMastermix: _input.PolymeraseAlreadyaddedtoMastermix,
			FwdPrimerName:                     fwdPrimer,
			RevPrimerName:                     revPrimer,
			TemplateName:                      template,
			ReactionName:                      reactionname,
			FwdPrimerVol:                      _input.DefaultPrimerVolume,
			RevPrimerVol:                      _input.DefaultPrimerVolume,
			PolymeraseVolume:                  _input.DefaultPolymeraseVolume,
			Templatevolume:                    _input.DefaultTemplateVol,
			Numberofcycles:                    1,
			InitDenaturationtime:              wunit.NewTime(30, "s"),
			Denaturationtime:                  wunit.NewTime(5, "s"),
			Annealingtime:                     wunit.NewTime(10, "s"),
			AnnealingTemp:                     wunit.NewTemperature(72, "C"), // Should be calculated from primer and template binding
			Extensiontime:                     wunit.NewTime(60, "s"),        // should be calculated from template length and polymerase rate
			Finalextensiontime:                wunit.NewTime(180, "s"),
			WellPosition:                      wellposition,

			FwdPrimer:     _input.FwdPrimertype,
			RevPrimer:     _input.RevPrimertype,
			PCRPolymerase: _input.DefaultPolymerase,
			MasterMix:     _input.MasterMix,
			Template:      _input.Templatetype,
			OutPlate:      _input.Plate},
		)

		// add result to reactions slice
		_output.Reactions = append(_output.Reactions, result.Outputs.Reaction)
		volumes = append(volumes, result.Outputs.Reaction.Volume())
		welllocations = append(welllocations, wellposition)
		_output.ReactionMap[reactionname] = result.Outputs.Reaction

		// increase counter by 1 ready for next iteration of loop
		counter++

	}

	// once all values of loop have been completed, export the plate contents as a csv file
	_output.Errors = append(_output.Errors, wtype.ExportPlateCSV(_input.Projectname+".csv", _input.Plate, _input.Projectname+"outputPlate", welllocations, _output.Reactions, volumes))

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AutoPCR_mmx_multiAnalysis(_ctx context.Context, _input *AutoPCR_mmx_multiInput, _output *AutoPCR_mmx_multiOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AutoPCR_mmx_multiValidation(_ctx context.Context, _input *AutoPCR_mmx_multiInput, _output *AutoPCR_mmx_multiOutput) {
}
func _AutoPCR_mmx_multiRun(_ctx context.Context, input *AutoPCR_mmx_multiInput) *AutoPCR_mmx_multiOutput {
	output := &AutoPCR_mmx_multiOutput{}
	_AutoPCR_mmx_multiSetup(_ctx, input)
	_AutoPCR_mmx_multiSteps(_ctx, input, output)
	_AutoPCR_mmx_multiAnalysis(_ctx, input, output)
	_AutoPCR_mmx_multiValidation(_ctx, input, output)
	return output
}

func AutoPCR_mmx_multiRunSteps(_ctx context.Context, input *AutoPCR_mmx_multiInput) *AutoPCR_mmx_multiSOutput {
	soutput := &AutoPCR_mmx_multiSOutput{}
	output := _AutoPCR_mmx_multiRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AutoPCR_mmx_multiNew() interface{} {
	return &AutoPCR_mmx_multiElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AutoPCR_mmx_multiInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AutoPCR_mmx_multiRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AutoPCR_mmx_multiInput{},
			Out: &AutoPCR_mmx_multiOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type AutoPCR_mmx_multiElement struct {
	inject.CheckedRunner
}

type AutoPCR_mmx_multiInput struct {
	DefaultBufferConcinX              int
	DefaultMasterMixVolume            wunit.Volume
	DefaultPolymerase                 *wtype.LHComponent
	DefaultPolymeraseVolume           wunit.Volume
	DefaultPrimerVolume               wunit.Volume
	DefaultTemplateVol                wunit.Volume
	FwdPrimertype                     *wtype.LHComponent
	MasterMix                         *wtype.LHComponent
	Plate                             *wtype.LHPlate
	PolymeraseAlreadyaddedtoMastermix bool
	PrimersalreadyAddedtoMasterMix    bool
	Projectname                       string
	Reactiontoprimerpair              map[string][]string
	Reactiontotemplate                map[string]string
	RevPrimertype                     *wtype.LHComponent
	Templatetype                      *wtype.LHComponent
}

type AutoPCR_mmx_multiOutput struct {
	Errors      []error
	ReactionMap map[string]*wtype.LHComponent
	Reactions   []*wtype.LHComponent
}

type AutoPCR_mmx_multiSOutput struct {
	Data struct {
		Errors []error
	}
	Outputs struct {
		ReactionMap map[string]*wtype.LHComponent
		Reactions   []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AutoPCR_mmx_multi",
		Constructor: AutoPCR_mmx_multiNew,
		Desc: component.ComponentDesc{
			Desc: "Perform multiple PCR reactions with common default parameters using a mastermix\n",
			Path: "src/github.com/antha-lang/elements/starter/MakeMasterMix_PCR/AutoPCR_mmx.an",
			Params: []component.ParamDesc{
				{Name: "DefaultBufferConcinX", Desc: "e.g. for  10X Q5 buffer this would be 10\n", Kind: "Parameters"},
				{Name: "DefaultMasterMixVolume", Desc: "", Kind: "Parameters"},
				{Name: "DefaultPolymerase", Desc: "", Kind: "Inputs"},
				{Name: "DefaultPolymeraseVolume", Desc: "", Kind: "Parameters"},
				{Name: "DefaultPrimerVolume", Desc: "", Kind: "Parameters"},
				{Name: "DefaultTemplateVol", Desc: "Volume of template in each reaction\n", Kind: "Parameters"},
				{Name: "FwdPrimertype", Desc: "", Kind: "Inputs"},
				{Name: "MasterMix", Desc: "", Kind: "Inputs"},
				{Name: "Plate", Desc: "", Kind: "Inputs"},
				{Name: "PolymeraseAlreadyaddedtoMastermix", Desc: "", Kind: "Parameters"},
				{Name: "PrimersalreadyAddedtoMasterMix", Desc: "", Kind: "Parameters"},
				{Name: "Projectname", Desc: "", Kind: "Parameters"},
				{Name: "Reactiontoprimerpair", Desc: "map of which reaction uses which primer pair e.g. [\"left homology arm\"]:\"fwdprimer\",\"revprimer\"\n", Kind: "Parameters"},
				{Name: "Reactiontotemplate", Desc: "map of which reaction uses which template e.g. [\"left homology arm\"]:\"templatename\"\n", Kind: "Parameters"},
				{Name: "RevPrimertype", Desc: "", Kind: "Inputs"},
				{Name: "Templatetype", Desc: "", Kind: "Inputs"},
				{Name: "Errors", Desc: "return an error message if an error is encountered\n", Kind: "Data"},
				{Name: "ReactionMap", Desc: "", Kind: "Outputs"},
				{Name: "Reactions", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
