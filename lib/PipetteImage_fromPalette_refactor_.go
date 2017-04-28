// Generates instructions to pipette out a defined image onto a defined plate using a defined palette of colours
package lib

import (
	"context"
	"fmt"
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

func _PipetteImage_fromPalette_refactorRequirements() {

}

// Conditions to run on startup
func _PipetteImage_fromPalette_refactorSetup(_ctx context.Context, _input *PipetteImage_fromPalette_refactorInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _PipetteImage_fromPalette_refactorSteps(_ctx context.Context, _input *PipetteImage_fromPalette_refactorInput, _output *PipetteImage_fromPalette_refactorOutput) {

	// if image is from url, download
	if _input.UseURL {
		_, err := download.File(_input.URL, _input.Imagefilename)
		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}
	}

	if _input.PosterizeImage {
		_, _input.Imagefilename = image.Posterize(_input.Imagefilename, _input.PosterizeLevels)
	}

	positiontocolourmap, _, _ := image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, &_input.Palette, _input.Rotate, _input.AutoRotate)

	// use position to colour map to make a colour to []well positions map
	colourtoWellLocationMap := make(map[string][]string)

	for well, colour := range positiontocolourmap {

		colourindex := _input.Palette.Index(colour)

		colourstring := strconv.Itoa(colourindex)

		if locations, found := colourtoWellLocationMap[colourstring]; !found {
			colourtoWellLocationMap[colourstring] = []string{well}
		} else {
			locations = append(locations, well)
			colourtoWellLocationMap[colourstring] = locations
		}
	}

	image.CheckAllResizealgorithms(_input.Imagefilename, _input.OutPlate, _input.Rotate, image.AllResampleFilters)

	solutions := make([]*wtype.LHComponent, 0)

	counter := 0

	for colourindex, wells := range colourtoWellLocationMap {

		component, componentpresent := _input.ColourIndextoComponentMap[colourindex]

		/*
			if !componentpresent {

				var foundthese []string

				for key, _ := range ColourIndextoComponentMap {
					foundthese = append(foundthese,key)
				}
				sort.Strings(foundthese)
				Errorf("Component ", colourindex, "not found in ColourIndextoComponentMap.", "Found these entries only: ", strings.Join(foundthese, ","))

			}
		*/
		if componentpresent {

			if _input.LiquidType != "" {
				fmt.Println("liquidtype", _input.LiquidType)
				liquidtype, err := wtype.LiquidTypeFromString(_input.LiquidType)

				if err != nil {
					execute.Errorf(_ctx, "for component", component.CName, err.Error())
				}

				component.Type = liquidtype

			}

			for _, locationkey := range wells {

				// due to trilution error, temporarily skip any wells with x in the well coordinates, e.g. x1,x2,x12
				if !strings.Contains(locationkey, "x") && !strings.Contains(locationkey, "X") {

					// use index of colour in palette to retrieve real colour

					colourint, err := strconv.Atoi(colourindex)

					if err != nil {
						execute.Errorf(_ctx, err.Error())
					}

					actualcolour := _input.Palette[colourint]

					if _input.OnlythisColour != "" && image.Colourcomponentmap[actualcolour] == _input.OnlythisColour {

						pixelSample := mixer.Sample(component, _input.VolumePerWell)
						solution := execute.MixTo(_ctx, _input.OutPlate.Type, locationkey, 1, pixelSample)
						solutions = append(solutions, solution)
						counter++

					} else if component.CName != _input.NotthisColour {
						pixelSample := mixer.Sample(component, _input.VolumePerWell)
						solution := execute.MixTo(_ctx, _input.OutPlate.Type, locationkey, 1, pixelSample)
						solutions = append(solutions, solution)
						counter++
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
func _PipetteImage_fromPalette_refactorAnalysis(_ctx context.Context, _input *PipetteImage_fromPalette_refactorInput, _output *PipetteImage_fromPalette_refactorOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PipetteImage_fromPalette_refactorValidation(_ctx context.Context, _input *PipetteImage_fromPalette_refactorInput, _output *PipetteImage_fromPalette_refactorOutput) {

}
func _PipetteImage_fromPalette_refactorRun(_ctx context.Context, input *PipetteImage_fromPalette_refactorInput) *PipetteImage_fromPalette_refactorOutput {
	output := &PipetteImage_fromPalette_refactorOutput{}
	_PipetteImage_fromPalette_refactorSetup(_ctx, input)
	_PipetteImage_fromPalette_refactorSteps(_ctx, input, output)
	_PipetteImage_fromPalette_refactorAnalysis(_ctx, input, output)
	_PipetteImage_fromPalette_refactorValidation(_ctx, input, output)
	return output
}

func PipetteImage_fromPalette_refactorRunSteps(_ctx context.Context, input *PipetteImage_fromPalette_refactorInput) *PipetteImage_fromPalette_refactorSOutput {
	soutput := &PipetteImage_fromPalette_refactorSOutput{}
	output := _PipetteImage_fromPalette_refactorRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PipetteImage_fromPalette_refactorNew() interface{} {
	return &PipetteImage_fromPalette_refactorElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PipetteImage_fromPalette_refactorInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PipetteImage_fromPalette_refactorRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PipetteImage_fromPalette_refactorInput{},
			Out: &PipetteImage_fromPalette_refactorOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type PipetteImage_fromPalette_refactorElement struct {
	inject.CheckedRunner
}

type PipetteImage_fromPalette_refactorInput struct {
	AutoRotate                bool
	ColourIndextoComponentMap map[string]*wtype.LHComponent
	Colourcomponents          []*wtype.LHComponent
	ColourtoWellLocationMap   map[string][]string
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

type PipetteImage_fromPalette_refactorOutput struct {
	Numberofpixels int
	Pixels         []*wtype.LHComponent
}

type PipetteImage_fromPalette_refactorSOutput struct {
	Data struct {
		Numberofpixels int
	}
	Outputs struct {
		Pixels []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PipetteImage_fromPalette_refactor",
		Constructor: PipetteImage_fromPalette_refactorNew,
		Desc: component.ComponentDesc{
			Desc: "Generates instructions to pipette out a defined image onto a defined plate using a defined palette of colours\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/PipetteImage/PipetteImage_fromPalette_refactor.an",
			Params: []component.ParamDesc{
				{Name: "AutoRotate", Desc: "", Kind: "Parameters"},
				{Name: "ColourIndextoComponentMap", Desc: "", Kind: "Parameters"},
				{Name: "Colourcomponents", Desc: "", Kind: "Inputs"},
				{Name: "ColourtoWellLocationMap", Desc: "", Kind: "Parameters"},
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
