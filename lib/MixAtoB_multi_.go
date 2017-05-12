// Mix a list of ComponentsA to ComponentsB.
// By default, each component in ComponentsB will have the component from the corresponding position in ComponentsA added to it
// e.g.
// ComponentsA [dna1,dna2] + ComponentsB [water, pbs] = MixedComponents [water + dna1, pbs + dna2]
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

// If a sample volume is specifed for a sample name contained  in ComponentsB, that volume of that component will be sampled.
// if a "default" volume is specified that will be used as the sample volume for all components which do not have a value explicitely specified.
// If no sample volume is specified for a component and no default set then the entire contents will be sampled.

// Output data of this protocol

// Physical inputs to this protocol

// List of components to add to all components in ComponentsB

// Each component in ComponentsB will have the component from the corresponding position in ComponentsA added to it

// Physical outputs to this protocol

// Conditions to run on startup
func _MixAtoB_multiSetup(_ctx context.Context, _input *MixAtoB_multiInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _MixAtoB_multiSteps(_ctx context.Context, _input *MixAtoB_multiInput, _output *MixAtoB_multiOutput) {

	if len(_input.ComponentsA) < len(_input.ComponentsB) {
		execute.Errorf(_ctx, "not enough samples in ComponentsA to mix to all of ComponentsB: Found %d in ComponentsA, %d in ComponentsB", len(_input.ComponentsA), len(_input.ComponentsB))
	}

	for i := range _input.ComponentsB {
		var sample *wtype.LHComponent

		var sampleVol wunit.Volume
		if vol, found := _input.SampleVolumes[_input.ComponentsB[i].CName]; found {
			sampleVol = vol
		} else if vol, found := _input.SampleVolumes["default"]; found {
			sampleVol = vol
		}

		if sampleVol.RawValue() > 0.0 {
			sample = mixer.Sample(_input.ComponentsA[i], sampleVol)
		} else {
			sample = mixer.SampleAll(_input.ComponentsA[i])
		}
		mixedComponent := execute.Mix(_ctx, _input.ComponentsB[i], sample)
		_output.MixedComponents = append(_output.MixedComponents, mixedComponent)

	}

}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _MixAtoB_multiAnalysis(_ctx context.Context, _input *MixAtoB_multiInput, _output *MixAtoB_multiOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _MixAtoB_multiValidation(_ctx context.Context, _input *MixAtoB_multiInput, _output *MixAtoB_multiOutput) {

}
func _MixAtoB_multiRun(_ctx context.Context, input *MixAtoB_multiInput) *MixAtoB_multiOutput {
	output := &MixAtoB_multiOutput{}
	_MixAtoB_multiSetup(_ctx, input)
	_MixAtoB_multiSteps(_ctx, input, output)
	_MixAtoB_multiAnalysis(_ctx, input, output)
	_MixAtoB_multiValidation(_ctx, input, output)
	return output
}

func MixAtoB_multiRunSteps(_ctx context.Context, input *MixAtoB_multiInput) *MixAtoB_multiSOutput {
	soutput := &MixAtoB_multiSOutput{}
	output := _MixAtoB_multiRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MixAtoB_multiNew() interface{} {
	return &MixAtoB_multiElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MixAtoB_multiInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MixAtoB_multiRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MixAtoB_multiInput{},
			Out: &MixAtoB_multiOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type MixAtoB_multiElement struct {
	inject.CheckedRunner
}

type MixAtoB_multiInput struct {
	ComponentsA   []*wtype.LHComponent
	ComponentsB   []*wtype.LHComponent
	SampleVolumes map[string]wunit.Volume
}

type MixAtoB_multiOutput struct {
	MixedComponents []*wtype.LHComponent
}

type MixAtoB_multiSOutput struct {
	Data struct {
	}
	Outputs struct {
		MixedComponents []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MixAtoB_multi",
		Constructor: MixAtoB_multiNew,
		Desc: component.ComponentDesc{
			Desc: "Mix a list of ComponentsA to ComponentsB.\nBy default, each component in ComponentsB will have the component from the corresponding position in ComponentsA added to it\ne.g.\nComponentsA [dna1,dna2] + ComponentsB [water, pbs] = MixedComponents [water + dna1, pbs + dna2]\n",
			Path: "src/github.com/antha-lang/elements/an/Utility/Mix/MixAtoB_Multi.an",
			Params: []component.ParamDesc{
				{Name: "ComponentsA", Desc: "List of components to add to all components in ComponentsB\n", Kind: "Inputs"},
				{Name: "ComponentsB", Desc: "Each component in ComponentsB will have the component from the corresponding position in ComponentsA added to it\n", Kind: "Inputs"},
				{Name: "SampleVolumes", Desc: "If a sample volume is specifed for a sample name contained  in ComponentsB, that volume of that component will be sampled.\nif a \"default\" volume is specified that will be used as the sample volume for all components which do not have a value explicitely specified.\nIf no sample volume is specified for a component and no default set then the entire contents will be sampled.\n", Kind: "Parameters"},
				{Name: "MixedComponents", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
