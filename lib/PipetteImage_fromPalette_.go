// Generates instructions to pipette out a defined image onto a defined plate using a defined palette of colours
package lib

import (
	"context"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/download"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"image/color"
	"strconv"
	"strings"
)

// Input parameters for this protocol (data)

// name of image file or if using URL use this field to set the desired filename
// select this if getting the image from a URL
// enter URL link to the image file here if applicable

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _PipetteImage_fromPaletteRequirements() {

}

// Conditions to run on startup
func _PipetteImage_fromPaletteSetup(_ctx context.Context, _input *PipetteImage_fromPaletteInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _PipetteImage_fromPaletteSteps(_ctx context.Context, _input *PipetteImage_fromPaletteInput, _output *PipetteImage_fromPaletteOutput) {

	// if image is from url, download
	if _input.UseURL {
		err := download.File(_input.URL, _input.Imagefilename)
		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}
	}

	var err error

	if _input.PosterizeImage {
		_, _input.Imagefilename = image.Posterize(_input.Imagefilename, _input.PosterizeLevels)
	}

	positiontocolourmap, _, _ := image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, &_input.Palette, _input.Rotate, _input.AutoRotate)

	image.CheckAllResizealgorithms(_input.Imagefilename, _input.OutPlate, _input.Rotate, image.AllResampleFilters)

	solutions := make([]*wtype.LHComponent, 0)

	counter := 0

	for locationkey, colour := range positiontocolourmap {

		// temp skip of wells with x, e.g. x1,x2,x12
		if !strings.Contains(locationkey, "x") && !strings.Contains(locationkey, "X") {

			colourindex := strconv.Itoa(_input.Palette.Index(colour))

			component, componentpresent := _input.ColourIndextoComponentMap[colourindex]

			/*
				if !componentpresent {


				for key, _ := range ColourIndextoComponentMap {
					foundthese = append(foundthese,key)
				}
				sort.Strings(foundthese)
				Errorf("Component ", colourindex, "not found in ColourIndextoComponentMap.", "Found these entries only: ", strings.Join(foundthese, ","))
				}
			*/

			if componentpresent {

				if _input.LiquidType != "" {
					component.Type, err = wtype.LiquidTypeFromString(_input.LiquidType)

					if err != nil {
						execute.Errorf(_ctx, "for component", component.CName, err.Error())
					}
				}
				if _input.OnlythisColour != "" {

					if image.Colourcomponentmap[colour] == _input.OnlythisColour {
						counter = counter + 1

						pixelSample := mixer.Sample(component, _input.VolumePerWell)
						solution := execute.MixTo(_ctx, _input.OutPlate.Type, locationkey, 1, pixelSample)
						solutions = append(solutions, solution)
					}

				} else {
					if component.CName != _input.NotthisColour {
						counter = counter + 1
						pixelSample := mixer.Sample(component, _input.VolumePerWell)
						solution := execute.MixTo(_ctx, _input.OutPlate.Type, locationkey, 1, pixelSample)
						solutions = append(solutions, solution)
					}
				}

			}
		}
	}
	_output.Pixels = solutions
	_output.Numberofpixels = len(_output.Pixels)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PipetteImage_fromPaletteAnalysis(_ctx context.Context, _input *PipetteImage_fromPaletteInput, _output *PipetteImage_fromPaletteOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PipetteImage_fromPaletteValidation(_ctx context.Context, _input *PipetteImage_fromPaletteInput, _output *PipetteImage_fromPaletteOutput) {

}
func _PipetteImage_fromPaletteRun(_ctx context.Context, input *PipetteImage_fromPaletteInput) *PipetteImage_fromPaletteOutput {
	output := &PipetteImage_fromPaletteOutput{}
	_PipetteImage_fromPaletteSetup(_ctx, input)
	_PipetteImage_fromPaletteSteps(_ctx, input, output)
	_PipetteImage_fromPaletteAnalysis(_ctx, input, output)
	_PipetteImage_fromPaletteValidation(_ctx, input, output)
	return output
}

func PipetteImage_fromPaletteRunSteps(_ctx context.Context, input *PipetteImage_fromPaletteInput) *PipetteImage_fromPaletteSOutput {
	soutput := &PipetteImage_fromPaletteSOutput{}
	output := _PipetteImage_fromPaletteRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PipetteImage_fromPaletteNew() interface{} {
	return &PipetteImage_fromPaletteElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PipetteImage_fromPaletteInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PipetteImage_fromPaletteRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PipetteImage_fromPaletteInput{},
			Out: &PipetteImage_fromPaletteOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type PipetteImage_fromPaletteElement struct {
	inject.CheckedRunner
}

type PipetteImage_fromPaletteInput struct {
	AutoRotate                bool
	ColourIndextoComponentMap map[string]*wtype.LHComponent
	Colourcomponents          []*wtype.LHComponent
	Imagefilename             string
	LiquidType                string
	NotthisColour             string
	OnlythisColour            string
	OutPlate                  *wtype.LHPlate
	Palette                   color.Palette
	PosterizeImage            bool
	PosterizeLevels           int
	Rotate                    bool
	URL                       string
	UseURL                    bool
	VolumePerWell             wunit.Volume
}

type PipetteImage_fromPaletteOutput struct {
	Numberofpixels int
	Pixels         []*wtype.LHComponent
}

type PipetteImage_fromPaletteSOutput struct {
	Data struct {
		Numberofpixels int
	}
	Outputs struct {
		Pixels []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PipetteImage_fromPalette",
		Constructor: PipetteImage_fromPaletteNew,
		Desc: component.ComponentDesc{
			Desc: "Generates instructions to pipette out a defined image onto a defined plate using a defined palette of colours\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/PipetteImage/PipetteImage/fromPalette/PipetteImage_fromPalette.an",
			Params: []component.ParamDesc{
				{Name: "AutoRotate", Desc: "", Kind: "Parameters"},
				{Name: "ColourIndextoComponentMap", Desc: "", Kind: "Parameters"},
				{Name: "Colourcomponents", Desc: "", Kind: "Inputs"},
				{Name: "Imagefilename", Desc: "name of image file or if using URL use this field to set the desired filename\n", Kind: "Parameters"},
				{Name: "LiquidType", Desc: "", Kind: "Parameters"},
				{Name: "NotthisColour", Desc: "", Kind: "Parameters"},
				{Name: "OnlythisColour", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Palette", Desc: "", Kind: "Parameters"},
				{Name: "PosterizeImage", Desc: "", Kind: "Parameters"},
				{Name: "PosterizeLevels", Desc: "", Kind: "Parameters"},
				{Name: "Rotate", Desc: "", Kind: "Parameters"},
				{Name: "URL", Desc: "enter URL link to the image file here if applicable\n", Kind: "Parameters"},
				{Name: "UseURL", Desc: "select this if getting the image from a URL\n", Kind: "Parameters"},
				{Name: "VolumePerWell", Desc: "", Kind: "Parameters"},
				{Name: "Numberofpixels", Desc: "", Kind: "Data"},
				{Name: "Pixels", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
