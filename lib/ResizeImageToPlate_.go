// Resize an image to fit the specified plate. We use the resizing package at:
// https://github.com/disintegration/imaging
// We use Lanczos resampling to resize since this is the best but slowest method.
package lib

import (
	"context"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	goimage "image"
)

// Input parameters for this protocol (data)

//Input image for this element
//Plate to which the image will be resized

// Data which is returned from this protocol, and data types

//Opened image
//error message

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _ResizeImageToPlateRequirements() {
}

// Conditions to run on startup
func _ResizeImageToPlateSetup(_ctx context.Context, _input *ResizeImageToPlateInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _ResizeImageToPlateSteps(_ctx context.Context, _input *ResizeImageToPlateInput, _output *ResizeImageToPlateOutput) {

	_output.OutputImage = image.ResizeImagetoPlateMin(_input.InputImage, _input.InputPlate)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _ResizeImageToPlateAnalysis(_ctx context.Context, _input *ResizeImageToPlateInput, _output *ResizeImageToPlateOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _ResizeImageToPlateValidation(_ctx context.Context, _input *ResizeImageToPlateInput, _output *ResizeImageToPlateOutput) {

}
func _ResizeImageToPlateRun(_ctx context.Context, input *ResizeImageToPlateInput) *ResizeImageToPlateOutput {
	output := &ResizeImageToPlateOutput{}
	_ResizeImageToPlateSetup(_ctx, input)
	_ResizeImageToPlateSteps(_ctx, input, output)
	_ResizeImageToPlateAnalysis(_ctx, input, output)
	_ResizeImageToPlateValidation(_ctx, input, output)
	return output
}

func ResizeImageToPlateRunSteps(_ctx context.Context, input *ResizeImageToPlateInput) *ResizeImageToPlateSOutput {
	soutput := &ResizeImageToPlateSOutput{}
	output := _ResizeImageToPlateRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ResizeImageToPlateNew() interface{} {
	return &ResizeImageToPlateElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ResizeImageToPlateInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ResizeImageToPlateRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ResizeImageToPlateInput{},
			Out: &ResizeImageToPlateOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type ResizeImageToPlateElement struct {
	inject.CheckedRunner
}

type ResizeImageToPlateInput struct {
	InputImage *goimage.NRGBA
	InputPlate *wtype.LHPlate
}

type ResizeImageToPlateOutput struct {
	Error       error
	OutputImage *goimage.NRGBA
}

type ResizeImageToPlateSOutput struct {
	Data struct {
		Error       error
		OutputImage *goimage.NRGBA
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "ResizeImageToPlate",
		Constructor: ResizeImageToPlateNew,
		Desc: component.ComponentDesc{
			Desc: "Resize an image to fit the specified plate. We use the resizing package at:\nhttps://github.com/disintegration/imaging\nWe use Lanczos resampling to resize since this is the best but slowest method.\n",
			Path: "src/github.com/antha-lang/elements/an/ImageHandling/LowLevel/ResizeImageToPlate.an",
			Params: []component.ParamDesc{
				{Name: "InputImage", Desc: "Input image for this element\n", Kind: "Parameters"},
				{Name: "InputPlate", Desc: "Plate to which the image will be resized\n", Kind: "Parameters"},
				{Name: "Error", Desc: "error message\n", Kind: "Data"},
				{Name: "OutputImage", Desc: "Opened image\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
