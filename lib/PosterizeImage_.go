// Posterizing refers to expressing an image using a defined number of different tones.
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
//Posterizing level. This refers to the number of colors used to express the image. It cannot be set to 1.

// Data which is returned from this protocol, and data types

//Opened image
//error message

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _PosterizeImageRequirements() {
}

// Conditions to run on startup
func _PosterizeImageSetup(_ctx context.Context, _input *PosterizeImageInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _PosterizeImageSteps(_ctx context.Context, _input *PosterizeImageInput, _output *PosterizeImageOutput) {

	var err error

	_output.OutputImage, err = image.Posterize(_input.InputImage, _input.Level)
	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PosterizeImageAnalysis(_ctx context.Context, _input *PosterizeImageInput, _output *PosterizeImageOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PosterizeImageValidation(_ctx context.Context, _input *PosterizeImageInput, _output *PosterizeImageOutput) {

}
func _PosterizeImageRun(_ctx context.Context, input *PosterizeImageInput) *PosterizeImageOutput {
	output := &PosterizeImageOutput{}
	_PosterizeImageSetup(_ctx, input)
	_PosterizeImageSteps(_ctx, input, output)
	_PosterizeImageAnalysis(_ctx, input, output)
	_PosterizeImageValidation(_ctx, input, output)
	return output
}

func PosterizeImageRunSteps(_ctx context.Context, input *PosterizeImageInput) *PosterizeImageSOutput {
	soutput := &PosterizeImageSOutput{}
	output := _PosterizeImageRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PosterizeImageNew() interface{} {
	return &PosterizeImageElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PosterizeImageInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PosterizeImageRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PosterizeImageInput{},
			Out: &PosterizeImageOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type PosterizeImageElement struct {
	inject.CheckedRunner
}

type PosterizeImageInput struct {
	InputImage *goimage.NRGBA
	Level      int
}

type PosterizeImageOutput struct {
	Error       error
	OutputImage *goimage.NRGBA
}

type PosterizeImageSOutput struct {
	Data struct {
		Error       error
		OutputImage *goimage.NRGBA
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PosterizeImage",
		Constructor: PosterizeImageNew,
		Desc: component.ComponentDesc{
			Desc: "Posterizing refers to expressing an image using a defined number of different tones.\n",
			Path: "src/github.com/antha-lang/elements/an/ImageHandling/LowLevel/PosterizeImage.an",
			Params: []component.ParamDesc{
				{Name: "InputImage", Desc: "Input image for this element\n", Kind: "Parameters"},
				{Name: "Level", Desc: "Posterizing level. This refers to the number of colors used to express the image. It cannot be set to 1.\n", Kind: "Parameters"},
				{Name: "Error", Desc: "error message\n", Kind: "Data"},
				{Name: "OutputImage", Desc: "Opened image\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
