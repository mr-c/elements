// Generates instructions to make a pallette of all colours in an image
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
)

// Input parameters for this protocol (data)

// name of image file or if using URL use this field to set the desired filename
// select this if getting the image from a URL
// enter URL link to the image file here if applicable

// Data which is returned from this protocol, and data types

//Colournames []string

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _MakePalette_OneByOneRequirements() {

}

// Conditions to run on startup
func _MakePalette_OneByOneSetup(_ctx context.Context, _input *MakePalette_OneByOneInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _MakePalette_OneByOneSteps(_ctx context.Context, _input *MakePalette_OneByOneInput, _output *MakePalette_OneByOneOutput) {

	// if image is from url, download
	if _input.UseURL {
		err := download.File(_input.URL, _input.Imagefilename)
		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}
	}

	if _input.PosterizeImage {
		_, _input.Imagefilename = image.Posterize(_input.Imagefilename, _input.PosterizeLevels)
	}

	// make palette of colours from image
	chosencolourpalette := image.MakeSmallPalleteFromImage(_input.Imagefilename, _input.OutPlate, _input.Rotate)

	// make a map of colour to well coordinates
	positiontocolourmap, _, _ := image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, &chosencolourpalette, _input.Rotate, _input.AutoRotate)

	// remove duplicates
	positiontocolourmap = image.RemoveDuplicatesValuesfromMap(positiontocolourmap)

	solutions := make([]*wtype.LHComponent, 0)
	colourtoComponentMap := make(map[string]*wtype.LHComponent)

	counter := 0

	for _, colour := range positiontocolourmap {

		colourindex := chosencolourpalette.Index(colour)

		if colour != nil {
			components := make([]*wtype.LHComponent, 0)

			cmyk := image.ColourtoCMYK(colour)

			var maxuint8 uint8 = 255

			if cmyk.C == 0 && cmyk.Y == 0 && cmyk.M == 0 && cmyk.K == 0 {

				continue

			} else {

				counter = counter + 1

				if cmyk.C > 0 {

					cyanvol := wunit.NewVolume(((float64(cmyk.C) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

					if cyanvol.RawValue() < 0.5 && cyanvol.Unit().PrefixedSymbol() == "ul" {
						cyanvol.SetValue(0.5)
					}

					cyanSample := mixer.Sample(_input.Cyan, cyanvol)
					components = append(components, cyanSample)
				}

				if cmyk.Y > 0 {
					yellowvol := wunit.NewVolume(((float64(cmyk.Y) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

					if yellowvol.RawValue() < 0.5 && yellowvol.Unit().PrefixedSymbol() == "ul" {
						yellowvol.SetValue(0.5)
					}

					yellowSample := mixer.Sample(_input.Yellow, yellowvol)
					components = append(components, yellowSample)
				}

				if cmyk.M > 0 {
					magentavol := wunit.NewVolume(((float64(cmyk.M) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

					if magentavol.RawValue() < 0.5 && magentavol.Unit().PrefixedSymbol() == "ul" {
						magentavol.SetValue(0.5)
					}

					magentaSample := mixer.Sample(_input.Magenta, magentavol)
					components = append(components, magentaSample)
				}

				if cmyk.K > 0 {

					blackvol := wunit.NewVolume(((float64(cmyk.K) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

					if blackvol.RawValue() < 0.5 && blackvol.Unit().PrefixedSymbol() == "ul" {
						blackvol.SetValue(0.5)
					}

					blackSample := mixer.Sample(_input.Black, blackvol)
					components = append(components, blackSample)
				}

				// top up colour to 4 x volumeforfullcolour with white to make the correct shade

				// get all component volumes
				// and change liquid types
				var componentvols []wunit.Volume
				for _, component := range components {
					componentvols = append(componentvols, component.Volume())
					component.Type = wtype.LTDoNotMix
				}
				// calculate volume of white to add
				whitevol := wunit.SubtractVolumes(wunit.MultiplyVolume(_input.VolumeForFullcolour, 4), componentvols)

				// mix with white sample
				_input.White.Type = wtype.LTPostMix

				whiteSample := mixer.Sample(_input.White, whitevol)
				components = append(components, whiteSample)

				solution := execute.MixInto(_ctx, _input.PalettePlate, "", components...)

				// change name of component
				originalname := solution.CName
				solution.CName = originalname + "_colour_" + strconv.Itoa(colourindex)

				// add solution to be exported later
				solutions = append(solutions, solution)
				colourtoComponentMap[strconv.Itoa(colourindex)] = solution

			}

		}
	}

	_output.Colours = solutions
	_output.Numberofcolours = len(chosencolourpalette)
	_output.Palette = chosencolourpalette
	_output.ColourtoComponentMap = colourtoComponentMap

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MakePalette_OneByOneAnalysis(_ctx context.Context, _input *MakePalette_OneByOneInput, _output *MakePalette_OneByOneOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _MakePalette_OneByOneValidation(_ctx context.Context, _input *MakePalette_OneByOneInput, _output *MakePalette_OneByOneOutput) {

}
func _MakePalette_OneByOneRun(_ctx context.Context, input *MakePalette_OneByOneInput) *MakePalette_OneByOneOutput {
	output := &MakePalette_OneByOneOutput{}
	_MakePalette_OneByOneSetup(_ctx, input)
	_MakePalette_OneByOneSteps(_ctx, input, output)
	_MakePalette_OneByOneAnalysis(_ctx, input, output)
	_MakePalette_OneByOneValidation(_ctx, input, output)
	return output
}

func MakePalette_OneByOneRunSteps(_ctx context.Context, input *MakePalette_OneByOneInput) *MakePalette_OneByOneSOutput {
	soutput := &MakePalette_OneByOneSOutput{}
	output := _MakePalette_OneByOneRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MakePalette_OneByOneNew() interface{} {
	return &MakePalette_OneByOneElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MakePalette_OneByOneInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MakePalette_OneByOneRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MakePalette_OneByOneInput{},
			Out: &MakePalette_OneByOneOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type MakePalette_OneByOneElement struct {
	inject.CheckedRunner
}

type MakePalette_OneByOneInput struct {
	AutoRotate          bool
	Black               *wtype.LHComponent
	Cyan                *wtype.LHComponent
	Imagefilename       string
	Magenta             *wtype.LHComponent
	OutPlate            *wtype.LHPlate
	PalettePlate        *wtype.LHPlate
	PosterizeImage      bool
	PosterizeLevels     int
	Rotate              bool
	URL                 string
	UseURL              bool
	VolumeForFullcolour wunit.Volume
	White               *wtype.LHComponent
	Yellow              *wtype.LHComponent
}

type MakePalette_OneByOneOutput struct {
	Colours              []*wtype.LHComponent
	ColourtoComponentMap map[string]*wtype.LHComponent
	Numberofcolours      int
	Palette              color.Palette
}

type MakePalette_OneByOneSOutput struct {
	Data struct {
		ColourtoComponentMap map[string]*wtype.LHComponent
		Numberofcolours      int
		Palette              color.Palette
	}
	Outputs struct {
		Colours []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MakePalette_OneByOne",
		Constructor: MakePalette_OneByOneNew,
		Desc: component.ComponentDesc{
			Desc: "Generates instructions to make a pallette of all colours in an image\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/PipetteImage/MakePalette_OnebyOne.an",
			Params: []component.ParamDesc{
				{Name: "AutoRotate", Desc: "", Kind: "Parameters"},
				{Name: "Black", Desc: "", Kind: "Inputs"},
				{Name: "Cyan", Desc: "", Kind: "Inputs"},
				{Name: "Imagefilename", Desc: "name of image file or if using URL use this field to set the desired filename\n", Kind: "Parameters"},
				{Name: "Magenta", Desc: "", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "PalettePlate", Desc: "", Kind: "Inputs"},
				{Name: "PosterizeImage", Desc: "", Kind: "Parameters"},
				{Name: "PosterizeLevels", Desc: "", Kind: "Parameters"},
				{Name: "Rotate", Desc: "", Kind: "Parameters"},
				{Name: "URL", Desc: "enter URL link to the image file here if applicable\n", Kind: "Parameters"},
				{Name: "UseURL", Desc: "select this if getting the image from a URL\n", Kind: "Parameters"},
				{Name: "VolumeForFullcolour", Desc: "", Kind: "Parameters"},
				{Name: "White", Desc: "", Kind: "Inputs"},
				{Name: "Yellow", Desc: "", Kind: "Inputs"},
				{Name: "Colours", Desc: "", Kind: "Outputs"},
				{Name: "ColourtoComponentMap", Desc: "", Kind: "Data"},
				{Name: "Numberofcolours", Desc: "", Kind: "Data"},
				{Name: "Palette", Desc: "Colournames []string\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
