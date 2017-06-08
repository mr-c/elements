//This element will convert a digital image to its physical representation, an AnthaImg object with well position information
package lib

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/microArch/factory"

	"context"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	goimage "image"
)

// Input parameters for this protocol (data)

//Image to use for this element
//Palette to use for this element
//Name of the plate type to use to print the image

// Data which is returned from this protocol, and data types

//converted image to anthaImage

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _MakeAnthaImageRequirements() {

}

// Conditions to run on startup
func _MakeAnthaImageSetup(_ctx context.Context, _input *MakeAnthaImageInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _MakeAnthaImageSteps(_ctx context.Context, _input *MakeAnthaImageInput, _output *MakeAnthaImageOutput) {

	//TODO: There should be an emlement to do this from parameters, but for now we have to use a string to refer to it
	plate := factory.GetPlateByType(_input.PlateName)

	//Checking if the image boundaries are the same as the

	//This function will create an AnthaImage object from a digital image.
	_output.AnthaImage = image.MakeAnthaImg(_input.InputImage, _input.AnthaPalette, plate)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MakeAnthaImageAnalysis(_ctx context.Context, _input *MakeAnthaImageInput, _output *MakeAnthaImageOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _MakeAnthaImageValidation(_ctx context.Context, _input *MakeAnthaImageInput, _output *MakeAnthaImageOutput) {

}
func _MakeAnthaImageRun(_ctx context.Context, input *MakeAnthaImageInput) *MakeAnthaImageOutput {
	output := &MakeAnthaImageOutput{}
	_MakeAnthaImageSetup(_ctx, input)
	_MakeAnthaImageSteps(_ctx, input, output)
	_MakeAnthaImageAnalysis(_ctx, input, output)
	_MakeAnthaImageValidation(_ctx, input, output)
	return output
}

func MakeAnthaImageRunSteps(_ctx context.Context, input *MakeAnthaImageInput) *MakeAnthaImageSOutput {
	soutput := &MakeAnthaImageSOutput{}
	output := _MakeAnthaImageRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MakeAnthaImageNew() interface{} {
	return &MakeAnthaImageElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MakeAnthaImageInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MakeAnthaImageRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MakeAnthaImageInput{},
			Out: &MakeAnthaImageOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type MakeAnthaImageElement struct {
	inject.CheckedRunner
}

type MakeAnthaImageInput struct {
	AnthaPalette image.AnthaPalette
	InputImage   *goimage.NRGBA
	PlateName    string
}

type MakeAnthaImageOutput struct {
	AnthaImage *image.AnthaImg
}

type MakeAnthaImageSOutput struct {
	Data struct {
		AnthaImage *image.AnthaImg
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MakeAnthaImage",
		Constructor: MakeAnthaImageNew,
		Desc: component.ComponentDesc{
			Desc: "This element will convert a digital image to its physical representation, an AnthaImg object with well position information\n",
			Path: "src/github.com/antha-lang/elements/an/ImageHandling/HighLevel/MakeAnthaImage.an",
			Params: []component.ParamDesc{
				{Name: "AnthaPalette", Desc: "Palette to use for this element\n", Kind: "Parameters"},
				{Name: "InputImage", Desc: "Image to use for this element\n", Kind: "Parameters"},
				{Name: "PlateName", Desc: "Name of the plate type to use to print the image\n", Kind: "Parameters"},
				{Name: "AnthaImage", Desc: "converted image to anthaImage\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
