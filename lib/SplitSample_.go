// Protocol SplitSample performs something.
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
func _SplitSampleSetup(_ctx context.Context, _input *SplitSampleInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _SplitSampleSteps(_ctx context.Context, _input *SplitSampleInput, _output *SplitSampleOutput) {
	sampleA := mixer.Sample(_input.InputSolution, _input.VolumeA)
	sampleB := mixer.Sample(_input.InputSolution, _input.VolumeB)
	_output.ComponentA = execute.MixNamed(_ctx, _input.Platetype, _input.WellA, _input.PlateNameA, sampleA)
	_output.ComponentB = execute.MixNamed(_ctx, _input.Platetype, _input.WellB, _input.PlateNameB, sampleB)
}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _SplitSampleAnalysis(_ctx context.Context, _input *SplitSampleInput, _output *SplitSampleOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _SplitSampleValidation(_ctx context.Context, _input *SplitSampleInput, _output *SplitSampleOutput) {

}
func _SplitSampleRun(_ctx context.Context, input *SplitSampleInput) *SplitSampleOutput {
	output := &SplitSampleOutput{}
	_SplitSampleSetup(_ctx, input)
	_SplitSampleSteps(_ctx, input, output)
	_SplitSampleAnalysis(_ctx, input, output)
	_SplitSampleValidation(_ctx, input, output)
	return output
}

func SplitSampleRunSteps(_ctx context.Context, input *SplitSampleInput) *SplitSampleSOutput {
	soutput := &SplitSampleSOutput{}
	output := _SplitSampleRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SplitSampleNew() interface{} {
	return &SplitSampleElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SplitSampleInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SplitSampleRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SplitSampleInput{},
			Out: &SplitSampleOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type SplitSampleElement struct {
	inject.CheckedRunner
}

type SplitSampleInput struct {
	InputSolution *wtype.LHComponent
	PlateNameA    string
	PlateNameB    string
	Platetype     string
	VolumeA       wunit.Volume
	VolumeB       wunit.Volume
	WellA         string
	WellB         string
}

type SplitSampleOutput struct {
	ComponentA *wtype.LHComponent
	ComponentB *wtype.LHComponent
}

type SplitSampleSOutput struct {
	Data struct {
	}
	Outputs struct {
		ComponentA *wtype.LHComponent
		ComponentB *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "SplitSample",
		Constructor: SplitSampleNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol SplitSample performs something.\n\nAll of this text should be used to describe what this protocol does.  It\nshould begin with a one sentence summary begining with \"Protocol X...\". If\nneccessary, a empty line with a detailed description can follow (like this\ndescription does).\n\nSpend some time thinking of a good protocol name as this is the name by\nwhich this protocol will be referred. It should convey the purpose and scope\nof the protocol to an outsider and should suggest an obvious\nparameterization.\n\nProtocol names are also case-sensitive, so try to use a consistent casing\nscheme.\n\nExamples of bad names:\n  - MyProtocol\n  - GeneAssembly\n  - WildCAPSsmallANDLARGE\n\nBetter names:\n  - Aliquot\n  - TypeIIsConstructAssembly\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/Exercises/LiquidHandlingExercises/SplitSample/SplitSample.an",
			Params: []component.ParamDesc{
				{Name: "InputSolution", Desc: "", Kind: "Inputs"},
				{Name: "PlateNameA", Desc: "", Kind: "Parameters"},
				{Name: "PlateNameB", Desc: "", Kind: "Parameters"},
				{Name: "Platetype", Desc: "", Kind: "Parameters"},
				{Name: "VolumeA", Desc: "", Kind: "Parameters"},
				{Name: "VolumeB", Desc: "", Kind: "Parameters"},
				{Name: "WellA", Desc: "", Kind: "Parameters"},
				{Name: "WellB", Desc: "", Kind: "Parameters"},
				{Name: "ComponentA", Desc: "", Kind: "Outputs"},
				{Name: "ComponentB", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
