// Aliquot a solution into a specified plate.
// optionally premix the solution before aliquoting
package lib

import (
	"context"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"strconv"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _AliquotByRowRequirements() {

}

// Conditions to run on startup
func _AliquotByRowSetup(_ctx context.Context, _input *AliquotByRowInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _AliquotByRowSteps(_ctx context.Context, _input *AliquotByRowInput, _output *AliquotByRowOutput) {

	var wellpositionarray []string
	for y := 0; y < _input.Outplate.WlsY; y++ {
		for x := 0; x < _input.Outplate.WlsX; x++ {
			wellposition := wutil.NumToAlpha(y+1) + strconv.Itoa(x+1)
			wellpositionarray = append(wellpositionarray, wellposition)
		}
	}

	var aliquots []*wtype.LHComponent
	counter := 0
	for i := 0; i < _input.NumberOfAliquots; i++ {
		sampleA := mixer.Sample(_input.SampleName, _input.AliquotVolume)
		aliquot := execute.MixInto(_ctx, _input.Outplate, wellpositionarray[counter], sampleA)
		aliquots = append(_output.Aliquots, aliquot)
		counter++
	}

	_output.Aliquots = aliquots

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AliquotByRowAnalysis(_ctx context.Context, _input *AliquotByRowInput, _output *AliquotByRowOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _AliquotByRowValidation(_ctx context.Context, _input *AliquotByRowInput, _output *AliquotByRowOutput) {

}
func _AliquotByRowRun(_ctx context.Context, input *AliquotByRowInput) *AliquotByRowOutput {
	output := &AliquotByRowOutput{}
	_AliquotByRowSetup(_ctx, input)
	_AliquotByRowSteps(_ctx, input, output)
	_AliquotByRowAnalysis(_ctx, input, output)
	_AliquotByRowValidation(_ctx, input, output)
	return output
}

func AliquotByRowRunSteps(_ctx context.Context, input *AliquotByRowInput) *AliquotByRowSOutput {
	soutput := &AliquotByRowSOutput{}
	output := _AliquotByRowRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AliquotByRowNew() interface{} {
	return &AliquotByRowElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AliquotByRowInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AliquotByRowRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AliquotByRowInput{},
			Out: &AliquotByRowOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type AliquotByRowElement struct {
	inject.CheckedRunner
}

type AliquotByRowInput struct {
	AliquotVolume    wunit.Volume
	NumberOfAliquots int
	Outplate         *wtype.LHPlate
	SampleName       *wtype.LHComponent
}

type AliquotByRowOutput struct {
	Aliquots []*wtype.LHComponent
}

type AliquotByRowSOutput struct {
	Data struct {
	}
	Outputs struct {
		Aliquots []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AliquotByRow",
		Constructor: AliquotByRowNew,
		Desc: component.ComponentDesc{
			Desc: "Aliquot a solution into a specified plate.\noptionally premix the solution before aliquoting\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson3_Mix_Loops/SimpleAliquot/AliquotByRowOrCol/AliquotByRow/AliquotByRow.an",
			Params: []component.ParamDesc{
				{Name: "AliquotVolume", Desc: "", Kind: "Parameters"},
				{Name: "NumberOfAliquots", Desc: "", Kind: "Parameters"},
				{Name: "Outplate", Desc: "", Kind: "Inputs"},
				{Name: "SampleName", Desc: "", Kind: "Inputs"},
				{Name: "Aliquots", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
