// This element mixes component B onto component A
package lib

import

// Place golang packages to import here
(
	"context"
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
func _MixSetup(_ctx context.Context, _input *MixInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _MixSteps(_ctx context.Context, _input *MixInput, _output *MixOutput) {
	_output.MixedComponent = execute.Mix(_ctx, _input.ComponentA, _input.ComponentB)
}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _MixAnalysis(_ctx context.Context, _input *MixInput, _output *MixOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _MixValidation(_ctx context.Context, _input *MixInput, _output *MixOutput) {

}
func _MixRun(_ctx context.Context, input *MixInput) *MixOutput {
	output := &MixOutput{}
	_MixSetup(_ctx, input)
	_MixSteps(_ctx, input, output)
	_MixAnalysis(_ctx, input, output)
	_MixValidation(_ctx, input, output)
	return output
}

func MixRunSteps(_ctx context.Context, input *MixInput) *MixSOutput {
	soutput := &MixSOutput{}
	output := _MixRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MixNew() interface{} {
	return &MixElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MixInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MixRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MixInput{},
			Out: &MixOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type MixElement struct {
	inject.CheckedRunner
}

type MixInput struct {
	ComponentA *wtype.LHComponent
	ComponentB *wtype.LHComponent
}

type MixOutput struct {
	MixedComponent *wtype.LHComponent
}

type MixSOutput struct {
	Data struct {
	}
	Outputs struct {
		MixedComponent *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Mix",
		Constructor: MixNew,
		Desc: component.ComponentDesc{
			Desc: "This element mixes component B onto component A\n",
			Path: "src/github.com/antha-lang/elements/an/Utility/Mix/Mix.an",
			Params: []component.ParamDesc{
				{Name: "ComponentA", Desc: "", Kind: "Inputs"},
				{Name: "ComponentB", Desc: "", Kind: "Inputs"},
				{Name: "MixedComponent", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
