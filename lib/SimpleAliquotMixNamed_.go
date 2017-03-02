// Aliquot a solution into a specified plate.
// optionally premix the solution before aliquoting
package lib

import (
	"context"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _SimpleAliquotMixNamedRequirements() {

}

// Conditions to run on startup
func _SimpleAliquotMixNamedSetup(_ctx context.Context, _input *SimpleAliquotMixNamedInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _SimpleAliquotMixNamedSteps(_ctx context.Context, _input *SimpleAliquotMixNamedInput, _output *SimpleAliquotMixNamedOutput) {

	for i := 0; i < _input.NumberOfAliquots; i++ {
		sampleA := mixer.Sample(_input.SampleName, _input.AliquotVolume)
		aliquot := execute.MixNamed(_ctx, _input.Outplate, "", _input.PlateName, sampleA)
		_output.Aliquots = append(_output.Aliquots, aliquot)
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _SimpleAliquotMixNamedAnalysis(_ctx context.Context, _input *SimpleAliquotMixNamedInput, _output *SimpleAliquotMixNamedOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _SimpleAliquotMixNamedValidation(_ctx context.Context, _input *SimpleAliquotMixNamedInput, _output *SimpleAliquotMixNamedOutput) {

}
func _SimpleAliquotMixNamedRun(_ctx context.Context, input *SimpleAliquotMixNamedInput) *SimpleAliquotMixNamedOutput {
	output := &SimpleAliquotMixNamedOutput{}
	_SimpleAliquotMixNamedSetup(_ctx, input)
	_SimpleAliquotMixNamedSteps(_ctx, input, output)
	_SimpleAliquotMixNamedAnalysis(_ctx, input, output)
	_SimpleAliquotMixNamedValidation(_ctx, input, output)
	return output
}

func SimpleAliquotMixNamedRunSteps(_ctx context.Context, input *SimpleAliquotMixNamedInput) *SimpleAliquotMixNamedSOutput {
	soutput := &SimpleAliquotMixNamedSOutput{}
	output := _SimpleAliquotMixNamedRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SimpleAliquotMixNamedNew() interface{} {
	return &SimpleAliquotMixNamedElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SimpleAliquotMixNamedInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SimpleAliquotMixNamedRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SimpleAliquotMixNamedInput{},
			Out: &SimpleAliquotMixNamedOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type SimpleAliquotMixNamedElement struct {
	inject.CheckedRunner
}

type SimpleAliquotMixNamedInput struct {
	AliquotVolume    wunit.Volume
	NumberOfAliquots int
	Outplate         string
	PlateName        string
	SampleName       *wtype.LHComponent
}

type SimpleAliquotMixNamedOutput struct {
	Aliquots []*wtype.LHComponent
}

type SimpleAliquotMixNamedSOutput struct {
	Data struct {
	}
	Outputs struct {
		Aliquots []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "SimpleAliquotMixNamed",
		Constructor: SimpleAliquotMixNamedNew,
		Desc: component.ComponentDesc{
			Desc: "Aliquot a solution into a specified plate.\noptionally premix the solution before aliquoting\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson3_Mix_Loops/SimpleAliquot/SimpleAliquotMixNamed/SimpleAliquotMixNamed/SimpleAliquotMixNamed.an",
			Params: []component.ParamDesc{
				{Name: "AliquotVolume", Desc: "", Kind: "Parameters"},
				{Name: "NumberOfAliquots", Desc: "", Kind: "Parameters"},
				{Name: "Outplate", Desc: "", Kind: "Parameters"},
				{Name: "PlateName", Desc: "", Kind: "Parameters"},
				{Name: "SampleName", Desc: "", Kind: "Inputs"},
				{Name: "Aliquots", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
