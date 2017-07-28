//This element will couple the given LHComponents with the Palette. It merges the arrays by index, so the first LHComponent
//will be the first color and so on. The resulting information is held in the AnthaPalette Object
package lib

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"image/color"
)

// Input parameters for this protocol (data)

//Palette to couple with LHComponents
//LHComponents to couple with the colors in the palette

// Data which is returned from this protocol, and data types

//Resulting AnthaPalette object, the LHComponents are coupled with the colors.

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _MakeAnthaPaletteRequirements() {
}

// Conditions to run on startup
func _MakeAnthaPaletteSetup(_ctx context.Context, _input *MakeAnthaPaletteInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _MakeAnthaPaletteSteps(_ctx context.Context, _input *MakeAnthaPaletteInput, _output *MakeAnthaPaletteOutput) {

	//checking if the recieved palette is empty
	if len(_input.Palette) == 0 {
		fmt.Println("Empty palette received")
	}

	if len(_input.Palette) != len(_input.LHComponents) {
		fmt.Println("different number of LHComponents and palette colors given")
	}

	//converting the palette to an AnthaPalette
	_output.AnthaPalette = image.MakeAnthaPalette(_input.Palette, _input.LHComponents)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MakeAnthaPaletteAnalysis(_ctx context.Context, _input *MakeAnthaPaletteInput, _output *MakeAnthaPaletteOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _MakeAnthaPaletteValidation(_ctx context.Context, _input *MakeAnthaPaletteInput, _output *MakeAnthaPaletteOutput) {

}
func _MakeAnthaPaletteRun(_ctx context.Context, input *MakeAnthaPaletteInput) *MakeAnthaPaletteOutput {
	output := &MakeAnthaPaletteOutput{}
	_MakeAnthaPaletteSetup(_ctx, input)
	_MakeAnthaPaletteSteps(_ctx, input, output)
	_MakeAnthaPaletteAnalysis(_ctx, input, output)
	_MakeAnthaPaletteValidation(_ctx, input, output)
	return output
}

func MakeAnthaPaletteRunSteps(_ctx context.Context, input *MakeAnthaPaletteInput) *MakeAnthaPaletteSOutput {
	soutput := &MakeAnthaPaletteSOutput{}
	output := _MakeAnthaPaletteRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MakeAnthaPaletteNew() interface{} {
	return &MakeAnthaPaletteElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MakeAnthaPaletteInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MakeAnthaPaletteRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MakeAnthaPaletteInput{},
			Out: &MakeAnthaPaletteOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type MakeAnthaPaletteElement struct {
	inject.CheckedRunner
}

type MakeAnthaPaletteInput struct {
	LHComponents []*wtype.LHComponent
	Palette      color.Palette
}

type MakeAnthaPaletteOutput struct {
	AnthaPalette *image.AnthaPalette
}

type MakeAnthaPaletteSOutput struct {
	Data struct {
		AnthaPalette *image.AnthaPalette
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MakeAnthaPalette",
		Constructor: MakeAnthaPaletteNew,
		Desc: component.ComponentDesc{
			Desc: "This element will couple the given LHComponents with the Palette. It merges the arrays by index, so the first LHComponent\nwill be the first color and so on. The resulting information is held in the AnthaPalette Object\n",
			Path: "src/github.com/antha-lang/elements/an/ImageHandling/MakeAnthaPalette/element.an",
			Params: []component.ParamDesc{
				{Name: "LHComponents", Desc: "LHComponents to couple with the colors in the palette\n", Kind: "Parameters"},
				{Name: "Palette", Desc: "Palette to couple with LHComponents\n", Kind: "Parameters"},
				{Name: "AnthaPalette", Desc: "Resulting AnthaPalette object, the LHComponents are coupled with the colors.\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
