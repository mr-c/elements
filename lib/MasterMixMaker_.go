// Make a general mastermix comprising of a list of components, list of volumes
// and specifying the number of reactions required
package lib

import (
	"context"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/factory"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _MasterMixMakerRequirements() {
}

// Conditions to run on startup
func _MasterMixMakerSetup(_ctx context.Context, _input *MasterMixMakerInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _MasterMixMakerSteps(_ctx context.Context, _input *MasterMixMakerInput, _output *MasterMixMakerOutput) {

	var mastermix *wtype.LHComponent

	if len(_input.Components) != len(_input.ComponentVolumesperReaction) {
		execute.Errorf(_ctx, "len(Components) != len(OtherComponentVolumes)")
	}

	// get components from factory and if not present use default dna component

	lhComponents := make([]*wtype.LHComponent, 0)

	for _, component := range _input.Components {

		if factory.ComponentInFactory(component) {
			lhComponents = append(lhComponents, factory.GetComponentByType(component))
		} else {
			// if component not in factory use dna as default component type
			defaultcomponent := factory.GetComponentByType("dna")
			defaultcomponent.CName = component
			lhComponents = append(lhComponents, defaultcomponent)
		}

	}

	// now make mastermix

	eachmastermix := make([]*wtype.LHComponent, 0)

	for k, component := range lhComponents {
		if k == len(lhComponents)-1 {
			component.Type = wtype.LTPostMix
		}

		// multiply volume of each component by number of reactions per mastermix
		adjustedvol := wunit.MultiplyVolume(_input.ComponentVolumesperReaction[k], float64(_input.Reactionspermastermix))

		componentSample := mixer.Sample(component, adjustedvol)

		eachmastermix = append(eachmastermix, componentSample)

	}
	mastermix = execute.MixInto(_ctx, _input.OutPlate, "", eachmastermix...)

	_output.Mastermix = mastermix

	_output.Status = "Mastermix Made"

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MasterMixMakerAnalysis(_ctx context.Context, _input *MasterMixMakerInput, _output *MasterMixMakerOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _MasterMixMakerValidation(_ctx context.Context, _input *MasterMixMakerInput, _output *MasterMixMakerOutput) {
}
func _MasterMixMakerRun(_ctx context.Context, input *MasterMixMakerInput) *MasterMixMakerOutput {
	output := &MasterMixMakerOutput{}
	_MasterMixMakerSetup(_ctx, input)
	_MasterMixMakerSteps(_ctx, input, output)
	_MasterMixMakerAnalysis(_ctx, input, output)
	_MasterMixMakerValidation(_ctx, input, output)
	return output
}

func MasterMixMakerRunSteps(_ctx context.Context, input *MasterMixMakerInput) *MasterMixMakerSOutput {
	soutput := &MasterMixMakerSOutput{}
	output := _MasterMixMakerRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MasterMixMakerNew() interface{} {
	return &MasterMixMakerElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MasterMixMakerInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MasterMixMakerRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MasterMixMakerInput{},
			Out: &MasterMixMakerOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type MasterMixMakerElement struct {
	inject.CheckedRunner
}

type MasterMixMakerInput struct {
	ComponentVolumesperReaction []wunit.Volume
	Components                  []string
	OutPlate                    *wtype.LHPlate
	Reactionspermastermix       int
}

type MasterMixMakerOutput struct {
	Mastermix *wtype.LHComponent
	Status    string
}

type MasterMixMakerSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
		Mastermix *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MasterMixMaker",
		Constructor: MasterMixMakerNew,
		Desc: component.ComponentDesc{
			Desc: "Make a general mastermix comprising of a list of components, list of volumes\nand specifying the number of reactions required\n",
			Path: "src/github.com/antha-lang/elements/starter/MakeMasterMix_PCR/MasterMixMaker.an",
			Params: []component.ParamDesc{
				{Name: "ComponentVolumesperReaction", Desc: "", Kind: "Parameters"},
				{Name: "Components", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Reactionspermastermix", Desc: "", Kind: "Parameters"},
				{Name: "Mastermix", Desc: "", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
