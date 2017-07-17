//This will take an AnthaImg object and generate the instructions for the robot to print it on a plate.
package lib

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"

	"image/color"

	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol (data)

//AnthaImage to print on a plate
//Volume of LHComponent needed to make a pixel
//Background color which we do not want to print, leave blank if you want to print all colors.

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _PrintAnthaImageRequirements() {

}

// Conditions to run on startup
func _PrintAnthaImageSetup(_ctx context.Context, _input *PrintAnthaImageInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _PrintAnthaImageSteps(_ctx context.Context, _input *PrintAnthaImageInput, _output *PrintAnthaImageOutput) {
	//------------------------------------------------------------------
	//Globals
	//------------------------------------------------------------------

	//placeholders
	var pixelSolution *wtype.LHComponent
	var wellLocation string
	var colorToExclude color.Color

	if _input.ColorToExclude == "" {
		colorToExclude = image.SelectColor(_input.ColorToExclude)
	} else {
		colorToExclude = color.NRGBA{}
	}

	//------------------------------------------------------------------
	//Iterating through each pixels in the image and pipetting them
	//------------------------------------------------------------------

	counter := 0
	for _, pix := range _input.AnthaImage.Pix {

		//Skipping pixel if it is of the background color
		if colorToExclude == pix.Color.Color {
			continue
		}

		//Getting the LHComponent of this pixel
		pixelSolution = pix.Color.Component

		//formatting the well coordinates to A1 format
		wellLocation = pix.Location.FormatA1()

		//initiating the LHComponent with the volume
		pixelSolution = mixer.Sample(pixelSolution, _input.PixVolume)

		//Executing the liquidHandling action
		execute.MixNamed(_ctx, _input.AnthaImage.Plate.Type, wellLocation, _input.AnthaImage.Plate.ID, pixelSolution)

		fmt.Println(counter)

		counter++
	}

	fmt.Println("AnthaImage printed")
	//getting the number of pipette actions

	//------------------------------------------------------------------
	//Returning printed AnthaImage
	//------------------------------------------------------------------

	_output.OutAnthaImage = _input.AnthaImage

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PrintAnthaImageAnalysis(_ctx context.Context, _input *PrintAnthaImageInput, _output *PrintAnthaImageOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PrintAnthaImageValidation(_ctx context.Context, _input *PrintAnthaImageInput, _output *PrintAnthaImageOutput) {

}
func _PrintAnthaImageRun(_ctx context.Context, input *PrintAnthaImageInput) *PrintAnthaImageOutput {
	output := &PrintAnthaImageOutput{}
	_PrintAnthaImageSetup(_ctx, input)
	_PrintAnthaImageSteps(_ctx, input, output)
	_PrintAnthaImageAnalysis(_ctx, input, output)
	_PrintAnthaImageValidation(_ctx, input, output)
	return output
}

func PrintAnthaImageRunSteps(_ctx context.Context, input *PrintAnthaImageInput) *PrintAnthaImageSOutput {
	soutput := &PrintAnthaImageSOutput{}
	output := _PrintAnthaImageRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PrintAnthaImageNew() interface{} {
	return &PrintAnthaImageElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PrintAnthaImageInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PrintAnthaImageRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PrintAnthaImageInput{},
			Out: &PrintAnthaImageOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type PrintAnthaImageElement struct {
	inject.CheckedRunner
}

type PrintAnthaImageInput struct {
	AnthaImage     *image.AnthaImg
	ColorToExclude string
	PixVolume      wunit.Volume
}

type PrintAnthaImageOutput struct {
	OutAnthaImage *image.AnthaImg
}

type PrintAnthaImageSOutput struct {
	Data struct {
	}
	Outputs struct {
		OutAnthaImage *image.AnthaImg
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PrintAnthaImage",
		Constructor: PrintAnthaImageNew,
		Desc: component.ComponentDesc{
			Desc: "This will take an AnthaImg object and generate the instructions for the robot to print it on a plate.\n",
			Path: "src/github.com/antha-lang/elements/an/ImageHandling/PrintAnthaImage.an",
			Params: []component.ParamDesc{
				{Name: "AnthaImage", Desc: "AnthaImage to print on a plate\n", Kind: "Parameters"},
				{Name: "ColorToExclude", Desc: "Background color which we do not want to print, leave blank if you want to print all colors.\n", Kind: "Parameters"},
				{Name: "PixVolume", Desc: "Volume of LHComponent needed to make a pixel\n", Kind: "Parameters"},
				{Name: "OutAnthaImage", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
