// Exports an image object to a file.
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

//desired name for the output file
//image object to export to File

// Data which is returned from this protocol, and data types

//Image as a file object

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _ExportImageRequirements() {
}

// Conditions to run on startup
func _ExportImageSetup(_ctx context.Context, _input *ExportImageInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _ExportImageSteps(_ctx context.Context, _input *ExportImageInput, _output *ExportImageOutput) {

	var err error

	//--------------------------------------------------------------
	//Exporting resulting images
	//--------------------------------------------------------------

	_output.ImageFile, err = image.Export(_input.Image, _input.OutputImageFileName)
	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _ExportImageAnalysis(_ctx context.Context, _input *ExportImageInput, _output *ExportImageOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _ExportImageValidation(_ctx context.Context, _input *ExportImageInput, _output *ExportImageOutput) {

}
func _ExportImageRun(_ctx context.Context, input *ExportImageInput) *ExportImageOutput {
	output := &ExportImageOutput{}
	_ExportImageSetup(_ctx, input)
	_ExportImageSteps(_ctx, input, output)
	_ExportImageAnalysis(_ctx, input, output)
	_ExportImageValidation(_ctx, input, output)
	return output
}

func ExportImageRunSteps(_ctx context.Context, input *ExportImageInput) *ExportImageSOutput {
	soutput := &ExportImageSOutput{}
	output := _ExportImageRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ExportImageNew() interface{} {
	return &ExportImageElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ExportImageInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ExportImageRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ExportImageInput{},
			Out: &ExportImageOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type ExportImageElement struct {
	inject.CheckedRunner
}

type ExportImageInput struct {
	Image               *goimage.NRGBA
	OutputImageFileName string
}

type ExportImageOutput struct {
	ImageFile wtype.File
}

type ExportImageSOutput struct {
	Data struct {
		ImageFile wtype.File
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "ExportImage",
		Constructor: ExportImageNew,
		Desc: component.ComponentDesc{
			Desc: "Exports an image object to a file.\n",
			Path: "src/github.com/antha-lang/elements/an/ImageHandling/LowLevel/ExportImage.an",
			Params: []component.ParamDesc{
				{Name: "Image", Desc: "image object to export to File\n", Kind: "Parameters"},
				{Name: "OutputImageFileName", Desc: "desired name for the output file\n", Kind: "Parameters"},
				{Name: "ImageFile", Desc: "Image as a file object\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
