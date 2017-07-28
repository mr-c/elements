//Select the colors you want to form a palette object. This element uses names defined in the standard library. You
//can either select a library, or specific colors.
package lib

import (
	"context"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"image/color"
)

// Input parameters for this protocol (data)

//Name of the Palette to use. The names are in the pixelToPlate package in the standard library.
//ID of the available colors. Leave blank if you want to use the palette

// Data which is returned from this protocol, and data types

//Selected palette
//error message

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _ChoosePaletteRequirements() {
}

// Conditions to run on startup
func _ChoosePaletteSetup(_ctx context.Context, _input *ChoosePaletteInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _ChoosePaletteSteps(_ctx context.Context, _input *ChoosePaletteInput, _output *ChoosePaletteOutput) {

	//-----------------------------------------------------------
	//Globals
	//-----------------------------------------------------------

	var palette color.Palette
	var tempColor color.Color

	//-----------------------------------------------------------
	//Creating palette
	//-----------------------------------------------------------

	switch {
	case _input.PaletteID != "":
		palette = image.SelectLibrary(_input.PaletteID)
	case len(_input.AvailableColors) > 0:
		for i := range _input.AvailableColors {
			tempColor = image.SelectColor(_input.AvailableColors[i])
			palette = append(palette, tempColor)
		}
	default:
		panic("no option selected in ChoosePalette")
	}

	if len(_input.AvailableColors) > 0 && _input.PaletteID != "" {
		panic("Choose to either select from library or from color ID")
	}

	//-----------------------------------------------------------
	//Returning retrieved palette
	//-----------------------------------------------------------

	_output.Palette = palette

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _ChoosePaletteAnalysis(_ctx context.Context, _input *ChoosePaletteInput, _output *ChoosePaletteOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _ChoosePaletteValidation(_ctx context.Context, _input *ChoosePaletteInput, _output *ChoosePaletteOutput) {

}
func _ChoosePaletteRun(_ctx context.Context, input *ChoosePaletteInput) *ChoosePaletteOutput {
	output := &ChoosePaletteOutput{}
	_ChoosePaletteSetup(_ctx, input)
	_ChoosePaletteSteps(_ctx, input, output)
	_ChoosePaletteAnalysis(_ctx, input, output)
	_ChoosePaletteValidation(_ctx, input, output)
	return output
}

func ChoosePaletteRunSteps(_ctx context.Context, input *ChoosePaletteInput) *ChoosePaletteSOutput {
	soutput := &ChoosePaletteSOutput{}
	output := _ChoosePaletteRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ChoosePaletteNew() interface{} {
	return &ChoosePaletteElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ChoosePaletteInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ChoosePaletteRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ChoosePaletteInput{},
			Out: &ChoosePaletteOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type ChoosePaletteElement struct {
	inject.CheckedRunner
}

type ChoosePaletteInput struct {
	AvailableColors []string
	PaletteID       string
}

type ChoosePaletteOutput struct {
	Error   error
	Palette color.Palette
}

type ChoosePaletteSOutput struct {
	Data struct {
		Error   error
		Palette color.Palette
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "ChoosePalette",
		Constructor: ChoosePaletteNew,
		Desc: component.ComponentDesc{
			Desc: "Select the colors you want to form a palette object. This element uses names defined in the standard library. You\ncan either select a library, or specific colors.\n",
			Path: "src/github.com/antha-lang/elements/an/ImageHandling/ChoosePalette/element.an",
			Params: []component.ParamDesc{
				{Name: "AvailableColors", Desc: "ID of the available colors. Leave blank if you want to use the palette\n", Kind: "Parameters"},
				{Name: "PaletteID", Desc: "Name of the Palette to use. The names are in the pixelToPlate package in the standard library.\n", Kind: "Parameters"},
				{Name: "Error", Desc: "error message\n", Kind: "Data"},
				{Name: "Palette", Desc: "Selected palette\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
