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

func _AliquotByRowColBool2Requirements() {

}

// Conditions to run on startup
func _AliquotByRowColBool2Setup(_ctx context.Context, _input *AliquotByRowColBool2Input) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _AliquotByRowColBool2Steps(_ctx context.Context, _input *AliquotByRowColBool2Input, _output *AliquotByRowColBool2Output) {

	var wellpositionarray []string
	var aliquots []*wtype.LHComponent

	if _input.ByRow {
		for y := 0; y < _input.Outplate.WlsY; y++ {
			for x := 0; x < _input.Outplate.WlsX; x++ {
				wellposition := wutil.NumToAlpha(y+1) + strconv.Itoa(x+1)
				wellpositionarray = append(wellpositionarray, wellposition)
			}
		}
	} else {
		for x := 0; x < _input.Outplate.WlsX; x++ {
			for y := 0; y < _input.Outplate.WlsY; y++ {
				wellposition := wutil.NumToAlpha(y+1) + strconv.Itoa(x+1)
				wellpositionarray = append(wellpositionarray, wellposition)
			}
		}
	}
	platenumber := 1
	counter := 0
	for i := 0; i < _input.NumberOfAliquots; i++ {
		sampleA := mixer.Sample(_input.SampleName, _input.AliquotVolume)
		aliquot := execute.MixTo(_ctx, _input.Outplate.Type, wellpositionarray[counter], platenumber, sampleA)
		aliquots = append(_output.Aliquots, aliquot)

		if counter+1 == len(wellpositionarray) {
			platenumber++
			counter = 0
		} else {
			counter++
		}
	}
	_output.Aliquots = aliquots
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AliquotByRowColBool2Analysis(_ctx context.Context, _input *AliquotByRowColBool2Input, _output *AliquotByRowColBool2Output) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _AliquotByRowColBool2Validation(_ctx context.Context, _input *AliquotByRowColBool2Input, _output *AliquotByRowColBool2Output) {

}
func _AliquotByRowColBool2Run(_ctx context.Context, input *AliquotByRowColBool2Input) *AliquotByRowColBool2Output {
	output := &AliquotByRowColBool2Output{}
	_AliquotByRowColBool2Setup(_ctx, input)
	_AliquotByRowColBool2Steps(_ctx, input, output)
	_AliquotByRowColBool2Analysis(_ctx, input, output)
	_AliquotByRowColBool2Validation(_ctx, input, output)
	return output
}

func AliquotByRowColBool2RunSteps(_ctx context.Context, input *AliquotByRowColBool2Input) *AliquotByRowColBool2SOutput {
	soutput := &AliquotByRowColBool2SOutput{}
	output := _AliquotByRowColBool2Run(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AliquotByRowColBool2New() interface{} {
	return &AliquotByRowColBool2Element{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AliquotByRowColBool2Input{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AliquotByRowColBool2Run(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AliquotByRowColBool2Input{},
			Out: &AliquotByRowColBool2Output{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type AliquotByRowColBool2Element struct {
	inject.CheckedRunner
}

type AliquotByRowColBool2Input struct {
	AliquotVolume    wunit.Volume
	ByRow            bool
	NumberOfAliquots int
	Outplate         *wtype.LHPlate
	SampleName       *wtype.LHComponent
}

type AliquotByRowColBool2Output struct {
	Aliquots []*wtype.LHComponent
}

type AliquotByRowColBool2SOutput struct {
	Data struct {
	}
	Outputs struct {
		Aliquots []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AliquotByRowColBool2",
		Constructor: AliquotByRowColBool2New,
		Desc: component.ComponentDesc{
			Desc: "Aliquot a solution into a specified plate.\noptionally premix the solution before aliquoting\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson3_Mix_Loops/SimpleAliquot/AliquotByRowColBool/AliquotByRowColBool2/AliquotByRowColBool2.an",
			Params: []component.ParamDesc{
				{Name: "AliquotVolume", Desc: "", Kind: "Parameters"},
				{Name: "ByRow", Desc: "", Kind: "Parameters"},
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
