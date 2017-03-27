// Make a general mastermix comprising of a list of components, list of volumes
// and specifying the number of reactions required
package lib

import (
	"context"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/setup"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/factory"
)

// Input parameters for this protocol (data)

// This specifies the multiplier of each of the Volumes for each component to add
// e.g. if "glucose" vol is "1ul" and Reactionspermastermix == 3 then 3ul glucose is added to mastermix

// Specify volumes per component in same order of components.
// The actual volume added will be multiplied by the number of Reactionspermastermix

// List of names of components to be added
// These will be used to look up components by name in the factory.
// If not found in the factory, new components will be created using dna_mix as a template
// If empty, the the ComponentIn will be returned as an output

// If using the inventory system, select whether to check inventory for parts so missing parts may be ordered.

// If set to true the mix will be prepared on the next available position on the input plate
// Otherwise the mastermix will be added to OutPlate

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// if OptimisePlateUsage is set to false this will be the plate type which the mastermix will be transferred to.

// Physical outputs from this protocol with types

func _MasterMixMakerRequirements() {
}

// Conditions to run on startup
func _MasterMixMakerSetup(_ctx context.Context, _input *MasterMixMakerInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _MasterMixMakerSteps(_ctx context.Context, _input *MasterMixMakerInput, _output *MasterMixMakerOutput) {

	// make up 20% extra to ensure reagents are sufficient accounting for dead volumes and evaporation
	extraReactions := float64(_input.Reactionspermastermix) * 1.2

	roundedReactions, err := wutil.RoundDown(extraReactions)

	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	roundedUpReactions := roundedReactions + 1

	if roundedUpReactions <= _input.Reactionspermastermix {
		_input.Reactionspermastermix = _input.Reactionspermastermix + 1
	}

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
			defaultcomponent := factory.GetComponentByType("dna_part")
			defaultcomponent.Type = wtype.LTDNAMIX

			defaultcomponent.CName = component
			lhComponents = append(lhComponents, defaultcomponent)
		}

	}

	if _input.CheckPartsInInventory {

		// First specify some handles for UI interaction
		// Adds Ordering handle for the UI
		lhComponents[0] = execute.Handle(_ctx, setup.OrderInfo(lhComponents[0]))
		// we need a plate prep step
		lhComponents[0] = execute.Handle(_ctx, setup.PlatePrep(lhComponents[0]))

		// a setup step
		lhComponents[0] = execute.Handle(_ctx, setup.MarkForSetup(lhComponents[0]))
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
	if _input.OptimisePlateUsage {
		mastermix = execute.Mix(_ctx, eachmastermix...)
	} else {
		mastermix = execute.MixInto(_ctx, _input.OutPlate, "", eachmastermix...)
	}

	_output.Mastermix = mastermix
	_output.PlateWithMastermix = _input.OutPlate

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
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type MasterMixMakerElement struct {
	inject.CheckedRunner
}

type MasterMixMakerInput struct {
	CheckPartsInInventory       bool
	ComponentVolumesperReaction []wunit.Volume
	Components                  []string
	OptimisePlateUsage          bool
	OutPlate                    *wtype.LHPlate
	Reactionspermastermix       int
}

type MasterMixMakerOutput struct {
	Mastermix          *wtype.LHComponent
	PlateWithMastermix *wtype.LHPlate
	Status             string
}

type MasterMixMakerSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
		Mastermix          *wtype.LHComponent
		PlateWithMastermix *wtype.LHPlate
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MasterMixMaker",
		Constructor: MasterMixMakerNew,
		Desc: component.ComponentDesc{
			Desc: "Make a general mastermix comprising of a list of components, list of volumes\nand specifying the number of reactions required\n",
			Path: "src/github.com/antha-lang/elements/starter/MakeMasterMix_PCR/MasterMixMaker.an",
			Params: []component.ParamDesc{
				{Name: "CheckPartsInInventory", Desc: "If using the inventory system, select whether to check inventory for parts so missing parts may be ordered.\n", Kind: "Parameters"},
				{Name: "ComponentVolumesperReaction", Desc: "Specify volumes per component in same order of components.\nThe actual volume added will be multiplied by the number of Reactionspermastermix\n", Kind: "Parameters"},
				{Name: "Components", Desc: "List of names of components to be added\nThese will be used to look up components by name in the factory.\nIf not found in the factory, new components will be created using dna_mix as a template\nIf empty, the the ComponentIn will be returned as an output\n", Kind: "Parameters"},
				{Name: "OptimisePlateUsage", Desc: "If set to true the mix will be prepared on the next available position on the input plate\nOtherwise the mastermix will be added to OutPlate\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "if OptimisePlateUsage is set to false this will be the plate type which the mastermix will be transferred to.\n", Kind: "Inputs"},
				{Name: "Reactionspermastermix", Desc: "This specifies the multiplier of each of the Volumes for each component to add\ne.g. if \"glucose\" vol is \"1ul\" and Reactionspermastermix == 3 then 3ul glucose is added to mastermix\n", Kind: "Parameters"},
				{Name: "Mastermix", Desc: "", Kind: "Outputs"},
				{Name: "PlateWithMastermix", Desc: "", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
