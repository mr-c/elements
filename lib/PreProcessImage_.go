package lib

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/download"
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

// name of image file or if using URL use this field to set the desired filename
// select this if getting the image from a URL
// enter URL link to the image file here if applicable

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _PreProcessImageRequirements() {

}

// Conditions to run on startup
func _PreProcessImageSetup(_ctx context.Context, _input *PreProcessImageInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _PreProcessImageSteps(_ctx context.Context, _input *PreProcessImageInput, _output *PreProcessImageOutput) {

	//-------------------------------------------------------------------------------------
	//Globals
	//-------------------------------------------------------------------------------------

	var imgFile wtype.File
	var imgBase *goimage.NRGBA
	var err error

	//-------------------------------------------------------------------------------------
	//Fetching image
	//-------------------------------------------------------------------------------------

	// if image is from url, download
	if _input.UseURL {
		//downloading image
		imgFile, err = download.File(_input.URL, _input.ImageFileName)
		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}

		//opening the image file
		imgBase, err = image.OpenFile(imgFile)
		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}
	}

	//opening the image file
	imgBase, err = image.OpenFile(_input.InputFile)
	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	//--------------------------------------------------------------
	//Image Processing
	//--------------------------------------------------------------

	if _input.PosterizeImage {

		imgBase, err = image.Posterize(imgBase, _input.PosterizeLevels)
		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}
	}

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

	_output.ProcessedImage, err = image.Export(imgBase, _input.ImageFileName)
	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PreProcessImageAnalysis(_ctx context.Context, _input *PreProcessImageInput, _output *PreProcessImageOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PreProcessImageValidation(_ctx context.Context, _input *PreProcessImageInput, _output *PreProcessImageOutput) {

}
func _PreProcessImageRun(_ctx context.Context, input *PreProcessImageInput) *PreProcessImageOutput {
	output := &PreProcessImageOutput{}
	_PreProcessImageSetup(_ctx, input)
	_PreProcessImageSteps(_ctx, input, output)
	_PreProcessImageAnalysis(_ctx, input, output)
	_PreProcessImageValidation(_ctx, input, output)
	return output
}

func PreProcessImageRunSteps(_ctx context.Context, input *PreProcessImageInput) *PreProcessImageSOutput {
	soutput := &PreProcessImageSOutput{}
	output := _PreProcessImageRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PreProcessImageNew() interface{} {
	return &PreProcessImageElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PreProcessImageInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PreProcessImageRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PreProcessImageInput{},
			Out: &PreProcessImageOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type PreProcessImageElement struct {
	inject.CheckedRunner
}

type PreProcessImageInput struct {
	AutoRotate               bool
	CheckAllResizeAlgorithms bool
	ImageFileName            string
	InputFile                wtype.File
	Negative                 bool
	OutPlate                 *wtype.LHPlate
	Palette                  string
	PosterizeImage           bool
	PosterizeLevels          int
	Rotate                   bool
	URL                      string
	UseURL                   bool
}

type PreProcessImageOutput struct {
	ProcessedImage wtype.File
}

type PreProcessImageSOutput struct {
	Data struct {
		ProcessedImage wtype.File
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PreProcessImage",
		Constructor: PreProcessImageNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/PipetteImage/PreProcessImage.an",
			Params: []component.ParamDesc{
				{Name: "AutoRotate", Desc: "", Kind: "Parameters"},
				{Name: "CheckAllResizeAlgorithms", Desc: "", Kind: "Parameters"},
				{Name: "ImageFileName", Desc: "name of image file or if using URL use this field to set the desired filename\n", Kind: "Parameters"},
				{Name: "InputFile", Desc: "", Kind: "Parameters"},
				{Name: "Negative", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Palette", Desc: "", Kind: "Parameters"},
				{Name: "PosterizeImage", Desc: "", Kind: "Parameters"},
				{Name: "PosterizeLevels", Desc: "", Kind: "Parameters"},
				{Name: "Rotate", Desc: "", Kind: "Parameters"},
				{Name: "URL", Desc: "enter URL link to the image file here if applicable\n", Kind: "Parameters"},
				{Name: "UseURL", Desc: "select this if getting the image from a URL\n", Kind: "Parameters"},
				{Name: "ProcessedImage", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
