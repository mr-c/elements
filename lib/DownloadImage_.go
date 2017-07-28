// Download and open an image file given a URL
package lib

import (
	"context"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/download"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	goimage "image"
)

// Input parameters for this protocol (data)

//URL from which the image will be downloaded

// Data which is returned from this protocol, and data types

//Opened image
//error message

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _DownloadImageRequirements() {
}

// Conditions to run on startup
func _DownloadImageSetup(_ctx context.Context, _input *DownloadImageInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _DownloadImageSteps(_ctx context.Context, _input *DownloadImageInput, _output *DownloadImageOutput) {

	//-------------------------------------------------------------------------------------
	//Globals
	//-------------------------------------------------------------------------------------

	var imgFile wtype.File

	//-------------------------------------------------------------------------------------
	//Fetching image file
	//-------------------------------------------------------------------------------------

	imgFile, _output.Error = download.File(_input.URL, "Downloaded file")
	if _output.Error != nil {
		execute.Errorf(_ctx, _output.Error.Error())
	}

	//-------------------------------------------------------------------------------------
	//Opening image
	//-------------------------------------------------------------------------------------

	_output.Image, _output.Error = image.OpenFile(imgFile)
	if _output.Error != nil {
		execute.Errorf(_ctx, _output.Error.Error())
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _DownloadImageAnalysis(_ctx context.Context, _input *DownloadImageInput, _output *DownloadImageOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _DownloadImageValidation(_ctx context.Context, _input *DownloadImageInput, _output *DownloadImageOutput) {

}
func _DownloadImageRun(_ctx context.Context, input *DownloadImageInput) *DownloadImageOutput {
	output := &DownloadImageOutput{}
	_DownloadImageSetup(_ctx, input)
	_DownloadImageSteps(_ctx, input, output)
	_DownloadImageAnalysis(_ctx, input, output)
	_DownloadImageValidation(_ctx, input, output)
	return output
}

func DownloadImageRunSteps(_ctx context.Context, input *DownloadImageInput) *DownloadImageSOutput {
	soutput := &DownloadImageSOutput{}
	output := _DownloadImageRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func DownloadImageNew() interface{} {
	return &DownloadImageElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &DownloadImageInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _DownloadImageRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &DownloadImageInput{},
			Out: &DownloadImageOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type DownloadImageElement struct {
	inject.CheckedRunner
}

type DownloadImageInput struct {
	URL string
}

type DownloadImageOutput struct {
	Error error
	Image *goimage.NRGBA
}

type DownloadImageSOutput struct {
	Data struct {
		Error error
		Image *goimage.NRGBA
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "DownloadImage",
		Constructor: DownloadImageNew,
		Desc: component.ComponentDesc{
			Desc: "Download and open an image file given a URL\n",
			Path: "src/github.com/antha-lang/elements/an/ImageHandling/DownloadImage/element.an",
			Params: []component.ParamDesc{
				{Name: "URL", Desc: "URL from which the image will be downloaded\n", Kind: "Parameters"},
				{Name: "Error", Desc: "error message\n", Kind: "Data"},
				{Name: "Image", Desc: "Opened image\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
