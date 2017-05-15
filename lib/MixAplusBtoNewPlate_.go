// Protocol MixA+BtoNewPlate performs mixing for two rows of liquid components of equal length.
// This protocol takes in two rows of components (ComponentsA and ComponentsB) and two rows of Volumes (VolumesA and VolumesB). All four rows have to be the same length.
// The element takes a specified volume of component A and transfers it to a new plate where a specified volume of component B is added.
// The order of the rows specifies which volume and which component is added per reaction.
// e.g. VolumesA [1ul,2ul] for ComponentsA [DNA1, DNA2] mixed with VolumesB [4ul,3ul] for ComponentsB [water,water] will mix 1ul of DNA1 with 4ul of water into well position 1 on a new plate. Then 2ul of DNA2 are mixed with 3ul of water into well position2 on a new plate.
package lib

import

// Place golang packages to import here
(
	"context"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Parameters to this protocol

// Output data of this protocol

// Physical inputs to this protocol

// Physical outputs to this protocol

// Conditions to run on startup
func _MixAplusBtoNewPlateSetup(_ctx context.Context, _input *MixAplusBtoNewPlateInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _MixAplusBtoNewPlateSteps(_ctx context.Context, _input *MixAplusBtoNewPlateInput, _output *MixAplusBtoNewPlateOutput) {
	if len(_input.ComponentsA) == len(_input.ComponentsB) && len(_input.VolumesA) == len(_input.VolumesB) && len(_input.ComponentsA) == len(_input.VolumesA) {
		for i := 0; i < len(_input.ComponentsA); i++ {
			_output.MixedComponents = append(_output.MixedComponents, execute.MixInto(_ctx, _input.OutPlate, "",
				mixer.Sample(_input.ComponentsA[i], _input.VolumesA[i]),
				mixer.Sample(_input.ComponentsB[i], _input.VolumesB[i])))
		}
	} else {
		execute.Errorf(_ctx, "The number of components specified in the two lists do not match! You have %s Volumes and %s Components for A and %s Volumes and %s Components for B.", len(_input.VolumesA), len(_input.ComponentsA), len(_input.ComponentsB), len(_input.VolumesB))
	}
}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _MixAplusBtoNewPlateAnalysis(_ctx context.Context, _input *MixAplusBtoNewPlateInput, _output *MixAplusBtoNewPlateOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _MixAplusBtoNewPlateValidation(_ctx context.Context, _input *MixAplusBtoNewPlateInput, _output *MixAplusBtoNewPlateOutput) {

}
func _MixAplusBtoNewPlateRun(_ctx context.Context, input *MixAplusBtoNewPlateInput) *MixAplusBtoNewPlateOutput {
	output := &MixAplusBtoNewPlateOutput{}
	_MixAplusBtoNewPlateSetup(_ctx, input)
	_MixAplusBtoNewPlateSteps(_ctx, input, output)
	_MixAplusBtoNewPlateAnalysis(_ctx, input, output)
	_MixAplusBtoNewPlateValidation(_ctx, input, output)
	return output
}

func MixAplusBtoNewPlateRunSteps(_ctx context.Context, input *MixAplusBtoNewPlateInput) *MixAplusBtoNewPlateSOutput {
	soutput := &MixAplusBtoNewPlateSOutput{}
	output := _MixAplusBtoNewPlateRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MixAplusBtoNewPlateNew() interface{} {
	return &MixAplusBtoNewPlateElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MixAplusBtoNewPlateInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MixAplusBtoNewPlateRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MixAplusBtoNewPlateInput{},
			Out: &MixAplusBtoNewPlateOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type MixAplusBtoNewPlateElement struct {
	inject.CheckedRunner
}

type MixAplusBtoNewPlateInput struct {
	ComponentsA []*wtype.LHComponent
	ComponentsB []*wtype.LHComponent
	OutPlate    *wtype.LHPlate
	VolumesA    []wunit.Volume
	VolumesB    []wunit.Volume
}

type MixAplusBtoNewPlateOutput struct {
	MixedComponents []*wtype.LHComponent
}

type MixAplusBtoNewPlateSOutput struct {
	Data struct {
	}
	Outputs struct {
		MixedComponents []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MixAplusBtoNewPlate",
		Constructor: MixAplusBtoNewPlateNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol MixA+BtoNewPlate performs mixing for two rows of liquid components of equal length.\nThis protocol takes in two rows of components (ComponentsA and ComponentsB) and two rows of Volumes (VolumesA and VolumesB). All four rows have to be the same length.\nThe element takes a specified volume of component A and transfers it to a new plate where a specified volume of component B is added.\nThe order of the rows specifies which volume and which component is added per reaction.\ne.g. VolumesA [1ul,2ul] for ComponentsA [DNA1, DNA2] mixed with VolumesB [4ul,3ul] for ComponentsB [water,water] will mix 1ul of DNA1 with 4ul of water into well position 1 on a new plate. Then 2ul of DNA2 are mixed with 3ul of water into well position2 on a new plate.\n",
			Path: "src/github.com/antha-lang/elements/an/Utility/MixNtoN.an",
			Params: []component.ParamDesc{
				{Name: "ComponentsA", Desc: "", Kind: "Inputs"},
				{Name: "ComponentsB", Desc: "", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "VolumesA", Desc: "", Kind: "Parameters"},
				{Name: "VolumesB", Desc: "", Kind: "Parameters"},
				{Name: "MixedComponents", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
