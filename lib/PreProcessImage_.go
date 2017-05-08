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

	// if image is from url, download
	if _input.UseURL {
		//downloading image
		imgFile, err := download.File(_input.URL, _input.Imagefilename)
		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}

		//opening the image file
		img, err := image.OpenFile(imgFile)
		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}
	}

	if _input.PosterizeImage {
		posterizedImg, err = image.Posterize(img, _input.PosterizeLevels)
		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}
	}

	chosencolourpalette := image.AvailablePalettes()[_input.Palette]

	if _input.CheckAllResizeAlgorithms {
		image.CheckAllResizealgorithms(_input.Imagefilename, _input.OutPlate, _input.Rotate, image.AllResampleFilters)
	}
	_, plateImg, _ := image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, &chosencolourpalette, _input.Rotate, _input.AutoRotate)

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
	Imagefilename            string
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
	ProcessedImageFilename string
}

type PreProcessImageSOutput struct {
	Data struct {
		ProcessedImageFilename string
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
				{Name: "Imagefilename", Desc: "name of image file or if using URL use this field to set the desired filename\n", Kind: "Parameters"},
				{Name: "Negative", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Palette", Desc: "", Kind: "Parameters"},
				{Name: "PosterizeImage", Desc: "", Kind: "Parameters"},
				{Name: "PosterizeLevels", Desc: "", Kind: "Parameters"},
				{Name: "Rotate", Desc: "", Kind: "Parameters"},
				{Name: "URL", Desc: "enter URL link to the image file here if applicable\n", Kind: "Parameters"},
				{Name: "UseURL", Desc: "select this if getting the image from a URL\n", Kind: "Parameters"},
				{Name: "ProcessedImageFilename", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
