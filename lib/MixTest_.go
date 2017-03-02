// Protocol MixTest performs something.
//
// All of this text should be used to describe what this protocol does.  It
// should begin with a one sentence summary begining with "Protocol X...". If
// neccessary, a empty line with a detailed description can follow (like this
// description does).
//
// Spend some time thinking of a good protocol name as this is the name by
// which this protocol will be referred. It should convey the purpose and scope
// of the protocol to an outsider and should suggest an obvious
// parameterization.
//
// Protocol names are also case-sensitive, so try to use a consistent casing
// scheme.
//
// Examples of bad names:
//   - MyProtocol
//   - GeneAssembly
//   - WildCAPSsmallANDLARGE
//
// Better names:
//   - Aliquot
//   - TypeIIsConstructAssembly
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
func _MixTestSetup(_ctx context.Context, _input *MixTestInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _MixTestSteps(_ctx context.Context, _input *MixTestInput, _output *MixTestOutput) {
	sampleA := mixer.Sample(_input.ComponentA, _input.VolumeA)
	_output.MixedComponent = execute.MixNamed(_ctx, _input.InputPlate, "A1", "sausages", sampleA)
	sampleB := mixer.Sample(_input.ComponentB, _input.VolumeA)
	_output.MixedComponentB = execute.MixNamed(_ctx, _input.InputPlate, "A1", "sausages", sampleB)
}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _MixTestAnalysis(_ctx context.Context, _input *MixTestInput, _output *MixTestOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _MixTestValidation(_ctx context.Context, _input *MixTestInput, _output *MixTestOutput) {

}
func _MixTestRun(_ctx context.Context, input *MixTestInput) *MixTestOutput {
	output := &MixTestOutput{}
	_MixTestSetup(_ctx, input)
	_MixTestSteps(_ctx, input, output)
	_MixTestAnalysis(_ctx, input, output)
	_MixTestValidation(_ctx, input, output)
	return output
}

func MixTestRunSteps(_ctx context.Context, input *MixTestInput) *MixTestSOutput {
	soutput := &MixTestSOutput{}
	output := _MixTestRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MixTestNew() interface{} {
	return &MixTestElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MixTestInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MixTestRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MixTestInput{},
			Out: &MixTestOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type MixTestElement struct {
	inject.CheckedRunner
}

type MixTestInput struct {
	ComponentA *wtype.LHComponent
	ComponentB *wtype.LHComponent
	InputPlate string
	VolumeA    wunit.Volume
}

type MixTestOutput struct {
	MixedComponent  *wtype.LHComponent
	MixedComponentB *wtype.LHComponent
	Sum             float64
}

type MixTestSOutput struct {
	Data struct {
		Sum float64
	}
	Outputs struct {
		MixedComponent  *wtype.LHComponent
		MixedComponentB *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MixTest",
		Constructor: MixTestNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol MixTest performs something.\n\nAll of this text should be used to describe what this protocol does.  It\nshould begin with a one sentence summary begining with \"Protocol X...\". If\nneccessary, a empty line with a detailed description can follow (like this\ndescription does).\n\nSpend some time thinking of a good protocol name as this is the name by\nwhich this protocol will be referred. It should convey the purpose and scope\nof the protocol to an outsider and should suggest an obvious\nparameterization.\n\nProtocol names are also case-sensitive, so try to use a consistent casing\nscheme.\n\nExamples of bad names:\n  - MyProtocol\n  - GeneAssembly\n  - WildCAPSsmallANDLARGE\n\nBetter names:\n  - Aliquot\n  - TypeIIsConstructAssembly\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson3_Mix_Loops/MixTestForSid/MixToNamed.an",
			Params: []component.ParamDesc{
				{Name: "ComponentA", Desc: "", Kind: "Inputs"},
				{Name: "ComponentB", Desc: "", Kind: "Inputs"},
				{Name: "InputPlate", Desc: "", Kind: "Inputs"},
				{Name: "VolumeA", Desc: "", Kind: "Parameters"},
				{Name: "MixedComponent", Desc: "", Kind: "Outputs"},
				{Name: "MixedComponentB", Desc: "", Kind: "Outputs"},
				{Name: "Sum", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
