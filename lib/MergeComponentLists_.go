// Protocol MergeComponentLists merges two lists of components together into a single list.
// A new list of MergedComponents will be made by appending ComponentsA with all entries in ComponentsB.
// The order of components will be preserved
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
func _MergeComponentListsSetup(_ctx context.Context, _input *MergeComponentListsInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _MergeComponentListsSteps(_ctx context.Context, _input *MergeComponentListsInput, _output *MergeComponentListsOutput) {

	for _, component := range _input.ComponentsA {
		_output.MergedComponents = append(_output.MergedComponents, component)
	}

	for _, component := range _input.ComponentsB {
		_output.MergedComponents = append(_output.MergedComponents, component)
	}

}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _MergeComponentListsAnalysis(_ctx context.Context, _input *MergeComponentListsInput, _output *MergeComponentListsOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _MergeComponentListsValidation(_ctx context.Context, _input *MergeComponentListsInput, _output *MergeComponentListsOutput) {

}
func _MergeComponentListsRun(_ctx context.Context, input *MergeComponentListsInput) *MergeComponentListsOutput {
	output := &MergeComponentListsOutput{}
	_MergeComponentListsSetup(_ctx, input)
	_MergeComponentListsSteps(_ctx, input, output)
	_MergeComponentListsAnalysis(_ctx, input, output)
	_MergeComponentListsValidation(_ctx, input, output)
	return output
}

func MergeComponentListsRunSteps(_ctx context.Context, input *MergeComponentListsInput) *MergeComponentListsSOutput {
	soutput := &MergeComponentListsSOutput{}
	output := _MergeComponentListsRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MergeComponentListsNew() interface{} {
	return &MergeComponentListsElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MergeComponentListsInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MergeComponentListsRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MergeComponentListsInput{},
			Out: &MergeComponentListsOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type MergeComponentListsElement struct {
	inject.CheckedRunner
}

type MergeComponentListsInput struct {
	ComponentsA []*wtype.LHComponent
	ComponentsB []*wtype.LHComponent
}

type MergeComponentListsOutput struct {
	MergedComponents []*wtype.LHComponent
}

type MergeComponentListsSOutput struct {
	Data struct {
	}
	Outputs struct {
		MergedComponents []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MergeComponentLists",
		Constructor: MergeComponentListsNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol MergeComponentLists merges two lists of components together into a single list.\nA new list of MergedComponents will be made by appending ComponentsA with all entries in ComponentsB.\nThe order of components will be preserved\n",
			Path: "src/github.com/antha-lang/elements/an/Utility/MergeComponentLists/MergeComponentLists.an",
			Params: []component.ParamDesc{
				{Name: "ComponentsA", Desc: "", Kind: "Inputs"},
				{Name: "ComponentsB", Desc: "", Kind: "Inputs"},
				{Name: "MergedComponents", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
