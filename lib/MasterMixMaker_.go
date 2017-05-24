// Make a general mastermix comprising of a list of components, list of volumes
// and specifying the number of reactions required.
// The output of this element can be wired into other elements such as AutoAssembly, AutoPCR or Aliquot.
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
	"strconv"
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

// Volume of mastermix made. This will account for the residual volume of the plate and add 20% extra to account for evaporation and transfer loss etc.

// Data output supplying the total volume added for each component.

// Physical Inputs to this protocol with types

// The plate which the mastermix will be made in.
// However, There is one scenario where it will not be used.
// 1. If OptimisePlateUsage is selected then the Antha scheduler will search for a suitable location to use to make the mastermix without adding an additional plate on to the deck.
// In either of these two cases if a plate is selected here then that plate's residual volume per well will be added to the total mastermix volume.
// If OptimisePlateUsage is selected it is therefore advisable to select the platetype of the most likely destination of the mastermix to be mixed to
// (i.e. one of the other plates used or the inplate, default inplate is usually pcrplate_skirted. You can simulate to check where the mastermix will be put.

// Physical outputs from this protocol with types

// The output of the protocol which can be wired into downstream elements such as Aliquot, AutoAssembly or AutoPCR

// The output plate containing the mastermix.

func _MasterMixMakerRequirements() {
}

// Conditions to run on startup
func _MasterMixMakerSetup(_ctx context.Context, _input *MasterMixMakerInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _MasterMixMakerSteps(_ctx context.Context, _input *MasterMixMakerInput, _output *MasterMixMakerOutput) {

	_output.VolumesUsed = make(map[string]wunit.Volume)

	// make up 20% extra to ensure reagents are sufficient accounting for dead volumes and evaporation
	extraReactions := float64(_input.Reactionspermastermix) * 1.2

	roundedReactions, err := wutil.RoundDown(extraReactions)

	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	roundedUpReactions := roundedReactions + 1

	if roundedUpReactions <= _input.Reactionspermastermix {
		_input.Reactionspermastermix = _input.Reactionspermastermix + 1
	} else {
		_input.Reactionspermastermix = roundedUpReactions
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

	var adjustedVols []wunit.Volume
	// adjust volumes
	for k := range _input.ComponentVolumesperReaction {

		// multiply volume of each component by number of reactions per mastermix
		adjustedVol := wunit.MultiplyVolume(_input.ComponentVolumesperReaction[k], float64(_input.Reactionspermastermix))

		adjustedVols = append(adjustedVols, adjustedVol)

	}

	prelimTotalVol := wunit.AddVolumes(adjustedVols)

	// now make mastermix

	eachmastermix := make([]*wtype.LHComponent, 0)

	for k, component := range lhComponents {

		if k == len(lhComponents)-1 {
			component.Type = wtype.LTMegaMix
		}

		// multiply volume of each component by number of reactions per mastermix
		adjustedvol := adjustedVols[k]

		var nilPlate *wtype.LHPlate

		if _input.OutPlate != nilPlate {
			residualVol := _input.OutPlate.Welltype.ResidualVolume()

			proportionOfresidualVol := wunit.MultiplyVolume(residualVol, float64(adjustedvol.SIValue()/prelimTotalVol.SIValue()))

			adjustedvol = wunit.AddVolumes([]wunit.Volume{adjustedvol, proportionOfresidualVol})

			if _input.OutPlate.Welltype.MaxVolume().LessThan(adjustedvol) {
				execute.Errorf(_ctx, "After accounting for residual well volume of %s, the Volume required of %s is too high for the %s well capacity of %s. Please select a plate with a large enough well capacity for this volume", residualVol, adjustedvol, _input.OutPlate.Name(), _input.OutPlate.Welltype.MaxVolume().ToString())
			}
			if _, found := _output.VolumesUsed[component.CName]; !found {
				_output.VolumesUsed[component.CName] = adjustedvol
			} else {
				var counter int
				for counter < 100 {
					name := component.CName + strconv.Itoa(counter)
					if _, found := _output.VolumesUsed[name]; !found {
						_output.VolumesUsed[component.CName] = adjustedvol
						break
					}
					counter++
				}
			}
		}

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

	_output.MasterMixVolume = _output.Mastermix.Volume()
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
	MasterMixVolume    wunit.Volume
	Mastermix          *wtype.LHComponent
	PlateWithMastermix *wtype.LHPlate
	VolumesUsed        map[string]wunit.Volume
}

type MasterMixMakerSOutput struct {
	Data struct {
		MasterMixVolume wunit.Volume
		VolumesUsed     map[string]wunit.Volume
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
			Desc: "Make a general mastermix comprising of a list of components, list of volumes\nand specifying the number of reactions required.\nThe output of this element can be wired into other elements such as AutoAssembly, AutoPCR or Aliquot.\n",
			Path: "src/github.com/antha-lang/elements/starter/MakeMasterMix_PCR/MasterMixMaker.an",
			Params: []component.ParamDesc{
				{Name: "CheckPartsInInventory", Desc: "If using the inventory system, select whether to check inventory for parts so missing parts may be ordered.\n", Kind: "Parameters"},
				{Name: "ComponentVolumesperReaction", Desc: "Specify volumes per component in same order of components.\nThe actual volume added will be multiplied by the number of Reactionspermastermix\n", Kind: "Parameters"},
				{Name: "Components", Desc: "List of names of components to be added\nThese will be used to look up components by name in the factory.\nIf not found in the factory, new components will be created using dna_mix as a template\nIf empty, the the ComponentIn will be returned as an output\n", Kind: "Parameters"},
				{Name: "OptimisePlateUsage", Desc: "If set to true the mix will be prepared on the next available position on the input plate\nOtherwise the mastermix will be added to OutPlate\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "The plate which the mastermix will be made in.\nHowever, There is one scenario where it will not be used.\n1. If OptimisePlateUsage is selected then the Antha scheduler will search for a suitable location to use to make the mastermix without adding an additional plate on to the deck.\nIn either of these two cases if a plate is selected here then that plate's residual volume per well will be added to the total mastermix volume.\nIf OptimisePlateUsage is selected it is therefore advisable to select the platetype of the most likely destination of the mastermix to be mixed to\n(i.e. one of the other plates used or the inplate, default inplate is usually pcrplate_skirted. You can simulate to check where the mastermix will be put.\n", Kind: "Inputs"},
				{Name: "Reactionspermastermix", Desc: "This specifies the multiplier of each of the Volumes for each component to add\ne.g. if \"glucose\" vol is \"1ul\" and Reactionspermastermix == 3 then 3ul glucose is added to mastermix\n", Kind: "Parameters"},
				{Name: "MasterMixVolume", Desc: "Volume of mastermix made. This will account for the residual volume of the plate and add 20% extra to account for evaporation and transfer loss etc.\n", Kind: "Data"},
				{Name: "Mastermix", Desc: "The output of the protocol which can be wired into downstream elements such as Aliquot, AutoAssembly or AutoPCR\n", Kind: "Outputs"},
				{Name: "PlateWithMastermix", Desc: "The output plate containing the mastermix.\n", Kind: "Outputs"},
				{Name: "VolumesUsed", Desc: "Data output supplying the total volume added for each component.\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
