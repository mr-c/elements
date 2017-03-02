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

func _SimpleAliquotMixToManyPlatesRequirements() {

}

// Conditions to run on startup
func _SimpleAliquotMixToManyPlatesSetup(_ctx context.Context, _input *SimpleAliquotMixToManyPlatesInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _SimpleAliquotMixToManyPlatesSteps(_ctx context.Context, _input *SimpleAliquotMixToManyPlatesInput, _output *SimpleAliquotMixToManyPlatesOutput) {
	for platenumber := 1; platenumber <= (_input.NumberOfPlates); platenumber++ {
		for i := 0; i < _input.NumberOfAliquots; i++ {
			sampleA := mixer.Sample(_input.SampleName, _input.AliquotVolume)
			aliquot := execute.MixTo(_ctx, _input.Outplate, "", platenumber, sampleA)
			_output.Aliquots = append(_output.Aliquots, aliquot)
		}
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _SimpleAliquotMixToManyPlatesAnalysis(_ctx context.Context, _input *SimpleAliquotMixToManyPlatesInput, _output *SimpleAliquotMixToManyPlatesOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _SimpleAliquotMixToManyPlatesValidation(_ctx context.Context, _input *SimpleAliquotMixToManyPlatesInput, _output *SimpleAliquotMixToManyPlatesOutput) {

}
func _SimpleAliquotMixToManyPlatesRun(_ctx context.Context, input *SimpleAliquotMixToManyPlatesInput) *SimpleAliquotMixToManyPlatesOutput {
	output := &SimpleAliquotMixToManyPlatesOutput{}
	_SimpleAliquotMixToManyPlatesSetup(_ctx, input)
	_SimpleAliquotMixToManyPlatesSteps(_ctx, input, output)
	_SimpleAliquotMixToManyPlatesAnalysis(_ctx, input, output)
	_SimpleAliquotMixToManyPlatesValidation(_ctx, input, output)
	return output
}

func SimpleAliquotMixToManyPlatesRunSteps(_ctx context.Context, input *SimpleAliquotMixToManyPlatesInput) *SimpleAliquotMixToManyPlatesSOutput {
	soutput := &SimpleAliquotMixToManyPlatesSOutput{}
	output := _SimpleAliquotMixToManyPlatesRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SimpleAliquotMixToManyPlatesNew() interface{} {
	return &SimpleAliquotMixToManyPlatesElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SimpleAliquotMixToManyPlatesInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SimpleAliquotMixToManyPlatesRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SimpleAliquotMixToManyPlatesInput{},
			Out: &SimpleAliquotMixToManyPlatesOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type SimpleAliquotMixToManyPlatesElement struct {
	inject.CheckedRunner
}

type SimpleAliquotMixToManyPlatesInput struct {
	AliquotVolume    wunit.Volume
	NumberOfAliquots int
	NumberOfPlates   int
	Outplate         string
	SampleName       *wtype.LHComponent
}

type SimpleAliquotMixToManyPlatesOutput struct {
	Aliquots []*wtype.LHComponent
}

type SimpleAliquotMixToManyPlatesSOutput struct {
	Data struct {
	}
	Outputs struct {
		Aliquots []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "SimpleAliquotMixToManyPlates",
		Constructor: SimpleAliquotMixToManyPlatesNew,
		Desc: component.ComponentDesc{
			Desc: "Aliquot a solution into a specified plate.\noptionally premix the solution before aliquoting\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson3_Mix_Loops/SimpleAliquot/SimpleAliquotMixTo/SimpleAliquotMixToManyPlates/SimpleAliquotMixToManyPlates.an",
			Params: []component.ParamDesc{
				{Name: "AliquotVolume", Desc: "", Kind: "Parameters"},
				{Name: "NumberOfAliquots", Desc: "", Kind: "Parameters"},
				{Name: "NumberOfPlates", Desc: "", Kind: "Parameters"},
				{Name: "Outplate", Desc: "", Kind: "Parameters"},
				{Name: "SampleName", Desc: "", Kind: "Inputs"},
				{Name: "Aliquots", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
