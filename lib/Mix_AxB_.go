// Protocol Mix_AxB mixes all combinations of two arrays of components or samples to a new location.
// The components may be samples if the Sample_multi element was used.
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
func _Mix_AxBSetup(_ctx context.Context, _input *Mix_AxBInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _Mix_AxBSteps(_ctx context.Context, _input *Mix_AxBInput, _output *Mix_AxBOutput) {

	for _, componentA := range _input.SampleAs {

		var sampleA *wtype.LHComponent
		var sampleVol wunit.Volume

		if vol, found := _input.SampleAVolumes[componentA.CName]; found {
			sampleVol = vol
		} else if vol, found := _input.SampleAVolumes["default"]; found {
			sampleVol = vol
		}

		if sampleVol.RawValue() > 0.0 {
			sampleA = mixer.Sample(componentA, sampleVol)
		} else {
			sampleA = mixer.SampleAll(componentA)
		}

		for _, componentB := range _input.SampleBs {
			var sampleB *wtype.LHComponent
			var sampleVol wunit.Volume

			if vol, found := _input.SampleBVolumes[componentB.CName]; found {
				sampleVol = vol
			} else if vol, found := _input.SampleBVolumes["default"]; found {
				sampleVol = vol
			}

			if sampleVol.RawValue() > 0.0 {
				sampleB = mixer.Sample(componentB, sampleVol)
			} else {
				sampleB = mixer.SampleAll(componentB)
			}
			mixedComponent := execute.MixInto(_ctx, _input.OutPlate, "", sampleA, sampleB)
			_output.MixedComponents = append(_output.MixedComponents, mixedComponent)
		}
	}

}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _Mix_AxBAnalysis(_ctx context.Context, _input *Mix_AxBInput, _output *Mix_AxBOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _Mix_AxBValidation(_ctx context.Context, _input *Mix_AxBInput, _output *Mix_AxBOutput) {

}
func _Mix_AxBRun(_ctx context.Context, input *Mix_AxBInput) *Mix_AxBOutput {
	output := &Mix_AxBOutput{}
	_Mix_AxBSetup(_ctx, input)
	_Mix_AxBSteps(_ctx, input, output)
	_Mix_AxBAnalysis(_ctx, input, output)
	_Mix_AxBValidation(_ctx, input, output)
	return output
}

func Mix_AxBRunSteps(_ctx context.Context, input *Mix_AxBInput) *Mix_AxBSOutput {
	soutput := &Mix_AxBSOutput{}
	output := _Mix_AxBRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Mix_AxBNew() interface{} {
	return &Mix_AxBElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Mix_AxBInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Mix_AxBRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Mix_AxBInput{},
			Out: &Mix_AxBOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type Mix_AxBElement struct {
	inject.CheckedRunner
}

type Mix_AxBInput struct {
	OutPlate       *wtype.LHPlate
	SampleAVolumes map[string]wunit.Volume
	SampleAs       []*wtype.LHComponent
	SampleBVolumes map[string]wunit.Volume
	SampleBs       []*wtype.LHComponent
}

type Mix_AxBOutput struct {
	MixedComponents []*wtype.LHComponent
}

type Mix_AxBSOutput struct {
	Data struct {
	}
	Outputs struct {
		MixedComponents []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Mix_AxB",
		Constructor: Mix_AxBNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol Mix_AxB mixes all combinations of two arrays of components or samples to a new location.\nThe components may be samples if the Sample_multi element was used.\n",
			Path: "src/github.com/antha-lang/elements/an/Utility/Mix_AxB/Mix_AxB.an",
			Params: []component.ParamDesc{
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "SampleAVolumes", Desc: "", Kind: "Parameters"},
				{Name: "SampleAs", Desc: "", Kind: "Inputs"},
				{Name: "SampleBVolumes", Desc: "", Kind: "Parameters"},
				{Name: "SampleBs", Desc: "", Kind: "Inputs"},
				{Name: "MixedComponents", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
