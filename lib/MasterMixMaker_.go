// Make a general mastermix comprising of a list of components, list of volumes
// and specifying the number of reactions required
package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

//TopUpBuffer *wtype.LHComponent // optional if nil this is ignored

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
		panic("len(Components) != len(OtherComponentVolumes)")
	}

	eachmastermix := make([]*wtype.LHComponent, 0)

	for k, component := range _input.Components {
		if k == len(_input.Components) {
			component.Type = wtype.LTNeedToMix //"NeedToMix"
		}

		// multiply volume of each component by number of reactions per mastermix
		adjustedvol := wunit.NewVolume(float64(_input.Reactionspermastermix)*_input.ComponentVolumesperReaction[k].SIValue()*1000000, "ul")

		componentSample := mixer.Sample(component, adjustedvol)
		component.CName = "component" + fmt.Sprint(k+1)
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
	Components                  []*wtype.LHComponent
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
				{Name: "Components", Desc: "TopUpBuffer *wtype.LHComponent // optional if nil this is ignored\n", Kind: "Inputs"},
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
