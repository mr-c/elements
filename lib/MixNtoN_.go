// Protocol MixNtoN performs mixing for two rows of liquid components of equal length

//This protocol takes in two arrays of components (A and B) of equal length and samples volume I of VolumeA of component i of A and transfers to a new Outplate where volume i of Volume B of component i of B is sampled into.
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
func _MixNtoNSetup(_ctx context.Context, _input *MixNtoNInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _MixNtoNSteps(_ctx context.Context, _input *MixNtoNInput, _output *MixNtoNOutput) {
	if len(_input.ComponentsA) == len(_input.ComponentsB) && len(_input.VolumesA) == len(_input.VolumesB) && len(_input.ComponentsA) == len(_input.VolumesA) {
		for i := 0; i < len(_input.ComponentsA); i++ {
			_output.MixedComponents[i] = execute.MixInto(_ctx, _input.OutPlate, "",
				mixer.Sample(_input.ComponentsA[i], _input.VolumesA[i]),
				mixer.Sample(_input.ComponentsB[i], _input.VolumesB[i]))
		}
	} else {
		execute.Errorf(_ctx, "The number of components specified in the two lists do not match! You have %s Volumes and %s Components for A and %s Volumes and %s Components for B.", len(_input.VolumesA), len(_input.ComponentsA), len(_input.ComponentsB), len(_input.VolumesB))
	}
}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _MixNtoNAnalysis(_ctx context.Context, _input *MixNtoNInput, _output *MixNtoNOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _MixNtoNValidation(_ctx context.Context, _input *MixNtoNInput, _output *MixNtoNOutput) {

}
func _MixNtoNRun(_ctx context.Context, input *MixNtoNInput) *MixNtoNOutput {
	output := &MixNtoNOutput{}
	_MixNtoNSetup(_ctx, input)
	_MixNtoNSteps(_ctx, input, output)
	_MixNtoNAnalysis(_ctx, input, output)
	_MixNtoNValidation(_ctx, input, output)
	return output
}

func MixNtoNRunSteps(_ctx context.Context, input *MixNtoNInput) *MixNtoNSOutput {
	soutput := &MixNtoNSOutput{}
	output := _MixNtoNRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MixNtoNNew() interface{} {
	return &MixNtoNElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MixNtoNInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MixNtoNRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MixNtoNInput{},
			Out: &MixNtoNOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type MixNtoNElement struct {
	inject.CheckedRunner
}

type MixNtoNInput struct {
	ComponentsA []*wtype.LHComponent
	ComponentsB []*wtype.LHComponent
	OutPlate    *wtype.LHPlate
	VolumesA    []wunit.Volume
	VolumesB    []wunit.Volume
}

type MixNtoNOutput struct {
	MixedComponents []*wtype.LHComponent
}

type MixNtoNSOutput struct {
	Data struct {
	}
	Outputs struct {
		MixedComponents []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MixNtoN",
		Constructor: MixNtoNNew,
		Desc: component.ComponentDesc{
			Desc: "This protocol takes in two arrays of components (A and B) of equal length and samples volume I of VolumeA of component i of A and transfers to a new Outplate where volume i of Volume B of component i of B is sampled into.\n",
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
