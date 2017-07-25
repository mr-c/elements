// Adds a list of components to a mastermix.
// Volumes of each component are specified by a map.
// A default volume may be specified which applies to all which are not present explicitely in the map.
// 20% extra volume is made up to ensure sufficient volume is made i.e. accounting for dead volumes etc.
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
	"github.com/antha-lang/antha/inventory"
	"strconv"
)

// Input parameters for this protocol (data)

// Specify volume per component name per reaction or specify a "default" to apply to all.
// The actual volume added will be multiplied by the number of Reactionspermastermix

// List of names of components to be added
// These will be used to look up components by name in the factory.
// If not found in the factory, new components will be created using dna_mix as a template
// If empty, the the ComponentIn will be returned as an output

// This specifies the multiplier of each of the Volumes for each component to add
// e.g. if "glucose": "1ul" and Reactionspermastermix == 3 then 3ul glucose is added to mastermix

// If using the inventory system, select whether to check inventory for parts so missing parts may be ordered.

// If this is selected the mastermix will be moved to a new location specified by OutPlate type
// If not selected, the components will be added to the ComponentIn

// Factor of total volume to make up as extra volume to account for evaporation.
// Default is 0.1 (10%)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// The component to add all new components to.

// if MixToNewLocation is set to true this will be the plate type which the mastermix will be transferred to.
// Otherwise the residual volume will be used to added to the total volume to add,
// so this plate should be selected as teh plate in which the mastermix is already in.

// Physical outputs from this protocol with types

func _AddToMastermixRequirements() {
}

// Conditions to run on startup
func _AddToMastermixSetup(_ctx context.Context, _input *AddToMastermixInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _AddToMastermixSteps(_ctx context.Context, _input *AddToMastermixInput, _output *AddToMastermixOutput) {

	_output.VolumesUsed = make(map[string]wunit.Volume)

	var factor float64

	if _input.MakeExtraPercentage == 0.0 {
		factor = 0.1
	}

	if _input.MakeExtraPercentage >= 0.0 && _input.MakeExtraPercentage <= 1 {
		factor = 1.0 + _input.MakeExtraPercentage
	} else {
		execute.Errorf(_ctx, "MakeExtraPercentage needs to be specified between 0.01 (1%) and 1 (100%) but specified as %f. Leave as zero to use default of 0.1 (10%)", _input.MakeExtraPercentage)
	}

	// make up extra volume based on MakeExtraPercentage parameter to ensure reagents are sufficient accounting for dead volumes and evaporation
	extraReactions := float64(_input.Reactionspermastermix) * factor

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

	// if no components to add, return original component in as output
	if len(_input.ComponentsToAdd) == 0 {

		if _input.MixToNewLocation {
			_output.Mastermix = execute.MixInto(_ctx, _input.OutPlate, "", _input.ComponentIn)
		} else {
			_output.Mastermix = _input.ComponentIn
		}
		_output.Status = "No Components added to Mastermix"

	} else {

		var mastermix *wtype.LHComponent

		// get components from factory and if not present use default dna component

		lhComponents := make([]*wtype.LHComponent, 0)

		if _input.MixToNewLocation {
			lhComponents = append(lhComponents, _input.ComponentIn)
		}
		for _, component := range _input.ComponentsToAdd {
			comp, err := inventory.NewComponent(_ctx, component)
			if err == inventory.ErrUnknownType {
				// if component not in factory use dna as default component type
				comp, err = inventory.NewComponent(_ctx, "dna")
				if err == nil {
					comp.CName = component
				}
			}
			if err != nil {
				execute.Errorf(_ctx, "cannot make component: %s", err)
			}

			lhComponents = append(lhComponents, comp)

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
		for _, component := range lhComponents {

			var volToUse wunit.Volume

			if vol, found := _input.VolumesToAdd[component.CName]; found {
				volToUse = vol
			} else if vol, found := _input.VolumesToAdd["default"]; found {
				volToUse = vol
			} else {
				execute.Errorf(_ctx, "No volume for %s or default volume specified.", component.CName)
			}

			// multiply volume of each component by number of reactions per mastermix
			adjustedVol := wunit.MultiplyVolume(volToUse, float64(_input.Reactionspermastermix))

			adjustedVols = append(adjustedVols, adjustedVol)

		}

		// Add volumes to work out total volume
		prelimTotalVol := wunit.AddVolumes(adjustedVols)

		// include volume of original component
		prelimTotalVol = wunit.AddVolumes([]wunit.Volume{_input.ComponentIn.Volume(), prelimTotalVol})

		// now make mastermix

		eachmastermix := make([]*wtype.LHComponent, 0)

		for k, component := range lhComponents {

			if k == len(lhComponents)-1 {
				component.Type = wtype.LTMegaMix
			}

			adjustedvol := adjustedVols[k]

			// now account for residual volume
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
		if _input.MixToNewLocation {
			mastermix = execute.MixInto(_ctx, _input.OutPlate, "", eachmastermix...)
		} else {
			for i := range eachmastermix {

				if i == 0 {
					mastermix = _input.ComponentIn
				}

				mastermix = execute.Mix(_ctx, mastermix, eachmastermix[i])
			}
		}
		_output.Mastermix = mastermix

		_output.Status = "Mastermix Made"
	}
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AddToMastermixAnalysis(_ctx context.Context, _input *AddToMastermixInput, _output *AddToMastermixOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AddToMastermixValidation(_ctx context.Context, _input *AddToMastermixInput, _output *AddToMastermixOutput) {
}
func _AddToMastermixRun(_ctx context.Context, input *AddToMastermixInput) *AddToMastermixOutput {
	output := &AddToMastermixOutput{}
	_AddToMastermixSetup(_ctx, input)
	_AddToMastermixSteps(_ctx, input, output)
	_AddToMastermixAnalysis(_ctx, input, output)
	_AddToMastermixValidation(_ctx, input, output)
	return output
}

func AddToMastermixRunSteps(_ctx context.Context, input *AddToMastermixInput) *AddToMastermixSOutput {
	soutput := &AddToMastermixSOutput{}
	output := _AddToMastermixRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AddToMastermixNew() interface{} {
	return &AddToMastermixElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AddToMastermixInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AddToMastermixRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AddToMastermixInput{},
			Out: &AddToMastermixOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type AddToMastermixElement struct {
	inject.CheckedRunner
}

type AddToMastermixInput struct {
	CheckPartsInInventory bool
	ComponentIn           *wtype.LHComponent
	ComponentsToAdd       []string
	MakeExtraPercentage   float64
	MixToNewLocation      bool
	OutPlate              *wtype.LHPlate
	Reactionspermastermix int
	VolumesToAdd          map[string]wunit.Volume
}

type AddToMastermixOutput struct {
	Mastermix   *wtype.LHComponent
	Status      string
	VolumesUsed map[string]wunit.Volume
}

type AddToMastermixSOutput struct {
	Data struct {
		Status      string
		VolumesUsed map[string]wunit.Volume
	}
	Outputs struct {
		Mastermix *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AddToMastermix",
		Constructor: AddToMastermixNew,
		Desc: component.ComponentDesc{
			Desc: "Adds a list of components to a mastermix.\nVolumes of each component are specified by a map.\nA default volume may be specified which applies to all which are not present explicitely in the map.\n20% extra volume is made up to ensure sufficient volume is made i.e. accounting for dead volumes etc.\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/MakeMastermix/AddToMasterMix/AddToMastermix.an",
			Params: []component.ParamDesc{
				{Name: "CheckPartsInInventory", Desc: "If using the inventory system, select whether to check inventory for parts so missing parts may be ordered.\n", Kind: "Parameters"},
				{Name: "ComponentIn", Desc: "The component to add all new components to.\n", Kind: "Inputs"},
				{Name: "ComponentsToAdd", Desc: "List of names of components to be added\nThese will be used to look up components by name in the factory.\nIf not found in the factory, new components will be created using dna_mix as a template\nIf empty, the the ComponentIn will be returned as an output\n", Kind: "Parameters"},
				{Name: "MakeExtraPercentage", Desc: "Factor of total volume to make up as extra volume to account for evaporation.\nDefault is 0.1 (10%)\n", Kind: "Parameters"},
				{Name: "MixToNewLocation", Desc: "If this is selected the mastermix will be moved to a new location specified by OutPlate type\nIf not selected, the components will be added to the ComponentIn\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "if MixToNewLocation is set to true this will be the plate type which the mastermix will be transferred to.\nOtherwise the residual volume will be used to added to the total volume to add,\nso this plate should be selected as teh plate in which the mastermix is already in.\n", Kind: "Inputs"},
				{Name: "Reactionspermastermix", Desc: "This specifies the multiplier of each of the Volumes for each component to add\ne.g. if \"glucose\": \"1ul\" and Reactionspermastermix == 3 then 3ul glucose is added to mastermix\n", Kind: "Parameters"},
				{Name: "VolumesToAdd", Desc: "Specify volume per component name per reaction or specify a \"default\" to apply to all.\nThe actual volume added will be multiplied by the number of Reactionspermastermix\n", Kind: "Parameters"},
				{Name: "Mastermix", Desc: "", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
				{Name: "VolumesUsed", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
