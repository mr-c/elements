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

func _AliquotByRowColBool3Requirements() {

}

// Conditions to run on startup
func _AliquotByRowColBool3Setup(_ctx context.Context, _input *AliquotByRowColBool3Input) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _AliquotByRowColBool3Steps(_ctx context.Context, _input *AliquotByRowColBool3Input, _output *AliquotByRowColBool3Output) {

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

	for _, samplename := range _input.SampleNames {
		for i := 0; i < _input.NumberOfAliquots; i++ {
			sampleA := mixer.Sample(samplename, _input.AliquotVolume)
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
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AliquotByRowColBool3Analysis(_ctx context.Context, _input *AliquotByRowColBool3Input, _output *AliquotByRowColBool3Output) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _AliquotByRowColBool3Validation(_ctx context.Context, _input *AliquotByRowColBool3Input, _output *AliquotByRowColBool3Output) {

}
func _AliquotByRowColBool3Run(_ctx context.Context, input *AliquotByRowColBool3Input) *AliquotByRowColBool3Output {
	output := &AliquotByRowColBool3Output{}
	_AliquotByRowColBool3Setup(_ctx, input)
	_AliquotByRowColBool3Steps(_ctx, input, output)
	_AliquotByRowColBool3Analysis(_ctx, input, output)
	_AliquotByRowColBool3Validation(_ctx, input, output)
	return output
}

func AliquotByRowColBool3RunSteps(_ctx context.Context, input *AliquotByRowColBool3Input) *AliquotByRowColBool3SOutput {
	soutput := &AliquotByRowColBool3SOutput{}
	output := _AliquotByRowColBool3Run(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AliquotByRowColBool3New() interface{} {
	return &AliquotByRowColBool3Element{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AliquotByRowColBool3Input{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AliquotByRowColBool3Run(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AliquotByRowColBool3Input{},
			Out: &AliquotByRowColBool3Output{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type AliquotByRowColBool3Element struct {
	inject.CheckedRunner
}

type AliquotByRowColBool3Input struct {
	AliquotVolume    wunit.Volume
	ByRow            bool
	NumberOfAliquots int
	Outplate         *wtype.LHPlate
	SampleNames      []*wtype.LHComponent
}

type AliquotByRowColBool3Output struct {
	Aliquots []*wtype.LHComponent
}

type AliquotByRowColBool3SOutput struct {
	Data struct {
	}
	Outputs struct {
		Aliquots []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AliquotByRowColBool3",
		Constructor: AliquotByRowColBool3New,
		Desc: component.ComponentDesc{
			Desc: "Aliquot a solution into a specified plate.\noptionally premix the solution before aliquoting\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson3_Mix_Loops/SimpleAliquot/AliquotByRowColBool/AliquotByRowColBool3/AliquotByRowColBool3.an",
			Params: []component.ParamDesc{
				{Name: "AliquotVolume", Desc: "", Kind: "Parameters"},
				{Name: "ByRow", Desc: "", Kind: "Parameters"},
				{Name: "NumberOfAliquots", Desc: "", Kind: "Parameters"},
				{Name: "Outplate", Desc: "", Kind: "Inputs"},
				{Name: "SampleNames", Desc: "", Kind: "Inputs"},
				{Name: "Aliquots", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
