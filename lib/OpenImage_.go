// Opens an image file, returning an image object
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

//File object

// Data which is returned from this protocol, and data types

//Opened image
//error message

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _OpenImageRequirements() {
}

// Conditions to run on startup
func _OpenImageSetup(_ctx context.Context, _input *OpenImageInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _OpenImageSteps(_ctx context.Context, _input *OpenImageInput, _output *OpenImageOutput) {

	//-------------------------------------------------------------------------------------
	//opening the image file
	//-------------------------------------------------------------------------------------

	_output.Image, _output.Error = image.OpenFile(_input.ImgFile)
	if _output.Error != nil {
		execute.Errorf(_ctx, _output.Error.Error())
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _OpenImageAnalysis(_ctx context.Context, _input *OpenImageInput, _output *OpenImageOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _OpenImageValidation(_ctx context.Context, _input *OpenImageInput, _output *OpenImageOutput) {

}
func _OpenImageRun(_ctx context.Context, input *OpenImageInput) *OpenImageOutput {
	output := &OpenImageOutput{}
	_OpenImageSetup(_ctx, input)
	_OpenImageSteps(_ctx, input, output)
	_OpenImageAnalysis(_ctx, input, output)
	_OpenImageValidation(_ctx, input, output)
	return output
}

func OpenImageRunSteps(_ctx context.Context, input *OpenImageInput) *OpenImageSOutput {
	soutput := &OpenImageSOutput{}
	output := _OpenImageRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func OpenImageNew() interface{} {
	return &OpenImageElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &OpenImageInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _OpenImageRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &OpenImageInput{},
			Out: &OpenImageOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type OpenImageElement struct {
	inject.CheckedRunner
}

type OpenImageInput struct {
	ImgFile wtype.File
}

type OpenImageOutput struct {
	Error error
	Image *goimage.NRGBA
}

type OpenImageSOutput struct {
	Data struct {
		Error error
		Image *goimage.NRGBA
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "OpenImage",
		Constructor: OpenImageNew,
		Desc: component.ComponentDesc{
			Desc: "Opens an image file, returning an image object\n",
			Path: "src/github.com/antha-lang/elements/an/ImageHandling/OpenImage.an",
			Params: []component.ParamDesc{
				{Name: "ImgFile", Desc: "File object\n", Kind: "Parameters"},
				{Name: "Error", Desc: "error message\n", Kind: "Data"},
				{Name: "Image", Desc: "Opened image\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
