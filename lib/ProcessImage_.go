//Aggregation of image manipulation functions.
package lib

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/anthalib/wtype"

	"context"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	goimage "image"
)

// Input parameters for this protocol (data)

//Image to use for his element
//Rotating image to fit plate
//Rotating image to fit plate
//Posterize the image (espress it with fewer colors)
//Posterizing value (number of colors to express the image with)
//Palette name with which to the image will be changed
//Iterate over every type of resizing algorithm to see the different image they output

// Data which is returned from this protocol, and data types

//Resulting image

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _ProcessImageRequirements() {

}

// Conditions to run on startup
func _ProcessImageSetup(_ctx context.Context, _input *ProcessImageInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _ProcessImageSteps(_ctx context.Context, _input *ProcessImageInput, _output *ProcessImageOutput) {

	//-------------------------------------------------------------------------------------
	//Globals
	//-------------------------------------------------------------------------------------

	//Placeholders for image and errors
	imgBase := _input.InputImage
	var err error

	//--------------------------------------------------------------
	//Posterize image
	//--------------------------------------------------------------

	if _input.PosterizeImage {
		imgBase, err = image.Posterize(imgBase, _input.PosterizeLevels)
		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}
	}

	//--------------------------------------------------------------
	//Iterate over resizing Algorithms types
	//--------------------------------------------------------------

	if _input.CheckAllResizeAlgorithms {
		image.CheckAllResizealgorithms(imgBase, _input.OutPlate, _input.Rotate, image.AllResampleFilters)
	}

	//--------------------------------------------------------------
	//Choosing Palette
	//--------------------------------------------------------------

	chosencolourpalette := image.AvailablePalettes()[_input.Palette]

	//--------------------------------------------------------------
	//Fitting image to plate
	//--------------------------------------------------------------

	_, imgBase = image.ImagetoPlatelayout(imgBase, _input.OutPlate, &chosencolourpalette, _input.Rotate, _input.AutoRotate)

	//--------------------------------------------------------------
	//Returning image
	//--------------------------------------------------------------

	_output.ProcessedImage = imgBase

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _ProcessImageAnalysis(_ctx context.Context, _input *ProcessImageInput, _output *ProcessImageOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _ProcessImageValidation(_ctx context.Context, _input *ProcessImageInput, _output *ProcessImageOutput) {

}
func _ProcessImageRun(_ctx context.Context, input *ProcessImageInput) *ProcessImageOutput {
	output := &ProcessImageOutput{}
	_ProcessImageSetup(_ctx, input)
	_ProcessImageSteps(_ctx, input, output)
	_ProcessImageAnalysis(_ctx, input, output)
	_ProcessImageValidation(_ctx, input, output)
	return output
}

func ProcessImageRunSteps(_ctx context.Context, input *ProcessImageInput) *ProcessImageSOutput {
	soutput := &ProcessImageSOutput{}
	output := _ProcessImageRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ProcessImageNew() interface{} {
	return &ProcessImageElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ProcessImageInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ProcessImageRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ProcessImageInput{},
			Out: &ProcessImageOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type ProcessImageElement struct {
	inject.CheckedRunner
}

type ProcessImageInput struct {
	AutoRotate               bool
	CheckAllResizeAlgorithms bool
	InputImage               *goimage.NRGBA
	OutPlate                 *wtype.LHPlate
	Palette                  string
	PosterizeImage           bool
	PosterizeLevels          int
	Rotate                   bool
}

type ProcessImageOutput struct {
	ProcessedImage *goimage.NRGBA
}

type ProcessImageSOutput struct {
	Data struct {
		ProcessedImage *goimage.NRGBA
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "ProcessImage",
		Constructor: ProcessImageNew,
		Desc: component.ComponentDesc{
			Desc: "Aggregation of image manipulation functions.\n",
			Path: "src/github.com/antha-lang/elements/an/ImageHandling/ProcessImage/element.an",
			Params: []component.ParamDesc{
				{Name: "AutoRotate", Desc: "Rotating image to fit plate\n", Kind: "Parameters"},
				{Name: "CheckAllResizeAlgorithms", Desc: "Iterate over every type of resizing algorithm to see the different image they output\n", Kind: "Parameters"},
				{Name: "InputImage", Desc: "Image to use for his element\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Palette", Desc: "Palette name with which to the image will be changed\n", Kind: "Parameters"},
				{Name: "PosterizeImage", Desc: "Posterize the image (espress it with fewer colors)\n", Kind: "Parameters"},
				{Name: "PosterizeLevels", Desc: "Posterizing value (number of colors to express the image with)\n", Kind: "Parameters"},
				{Name: "Rotate", Desc: "Rotating image to fit plate\n", Kind: "Parameters"},
				{Name: "ProcessedImage", Desc: "Resulting image\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
