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

func _SimpleAliquotMixNamedManyPlatesRequirements() {

}

// Conditions to run on startup
func _SimpleAliquotMixNamedManyPlatesSetup(_ctx context.Context, _input *SimpleAliquotMixNamedManyPlatesInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _SimpleAliquotMixNamedManyPlatesSteps(_ctx context.Context, _input *SimpleAliquotMixNamedManyPlatesInput, _output *SimpleAliquotMixNamedManyPlatesOutput) {
	for _, platename := range _input.PlateNames {
		for i := 0; i < _input.NumberOfAliquots; i++ {
			sampleA := mixer.Sample(_input.SampleName, _input.AliquotVolume)
			aliquot := execute.MixNamed(_ctx, _input.Outplate, "", platename, sampleA)
			_output.Aliquots = append(_output.Aliquots, aliquot)
		}
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _SimpleAliquotMixNamedManyPlatesAnalysis(_ctx context.Context, _input *SimpleAliquotMixNamedManyPlatesInput, _output *SimpleAliquotMixNamedManyPlatesOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _SimpleAliquotMixNamedManyPlatesValidation(_ctx context.Context, _input *SimpleAliquotMixNamedManyPlatesInput, _output *SimpleAliquotMixNamedManyPlatesOutput) {

}
func _SimpleAliquotMixNamedManyPlatesRun(_ctx context.Context, input *SimpleAliquotMixNamedManyPlatesInput) *SimpleAliquotMixNamedManyPlatesOutput {
	output := &SimpleAliquotMixNamedManyPlatesOutput{}
	_SimpleAliquotMixNamedManyPlatesSetup(_ctx, input)
	_SimpleAliquotMixNamedManyPlatesSteps(_ctx, input, output)
	_SimpleAliquotMixNamedManyPlatesAnalysis(_ctx, input, output)
	_SimpleAliquotMixNamedManyPlatesValidation(_ctx, input, output)
	return output
}

func SimpleAliquotMixNamedManyPlatesRunSteps(_ctx context.Context, input *SimpleAliquotMixNamedManyPlatesInput) *SimpleAliquotMixNamedManyPlatesSOutput {
	soutput := &SimpleAliquotMixNamedManyPlatesSOutput{}
	output := _SimpleAliquotMixNamedManyPlatesRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SimpleAliquotMixNamedManyPlatesNew() interface{} {
	return &SimpleAliquotMixNamedManyPlatesElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SimpleAliquotMixNamedManyPlatesInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SimpleAliquotMixNamedManyPlatesRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SimpleAliquotMixNamedManyPlatesInput{},
			Out: &SimpleAliquotMixNamedManyPlatesOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type SimpleAliquotMixNamedManyPlatesElement struct {
	inject.CheckedRunner
}

type SimpleAliquotMixNamedManyPlatesInput struct {
	AliquotVolume    wunit.Volume
	NumberOfAliquots int
	Outplate         string
	PlateNames       []string
	SampleName       *wtype.LHComponent
}

type SimpleAliquotMixNamedManyPlatesOutput struct {
	Aliquots []*wtype.LHComponent
}

type SimpleAliquotMixNamedManyPlatesSOutput struct {
	Data struct {
	}
	Outputs struct {
		Aliquots []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "SimpleAliquotMixNamedManyPlates",
		Constructor: SimpleAliquotMixNamedManyPlatesNew,
		Desc: component.ComponentDesc{
			Desc: "Aliquot a solution into a specified plate.\noptionally premix the solution before aliquoting\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson3_Mix_Loops/SimpleAliquot/SimpleAliquotMixNamed/SimpleAliquotMixNamedManyPlates/SimpleAliquotMixNamedManyPlates.an",
			Params: []component.ParamDesc{
				{Name: "AliquotVolume", Desc: "", Kind: "Parameters"},
				{Name: "NumberOfAliquots", Desc: "", Kind: "Parameters"},
				{Name: "Outplate", Desc: "", Kind: "Parameters"},
				{Name: "PlateNames", Desc: "", Kind: "Parameters"},
				{Name: "SampleName", Desc: "", Kind: "Inputs"},
				{Name: "Aliquots", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
