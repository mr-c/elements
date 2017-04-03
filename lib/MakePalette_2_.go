// Generates instructions to make a pallette of all colours in an image
package lib

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/download"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/export"
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

//Colournames []string

// map of colour name (as index) to component name

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _MakePalette_2Requirements() {

}

// Conditions to run on startup
func _MakePalette_2Setup(_ctx context.Context, _input *MakePalette_2Input) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _MakePalette_2Steps(_ctx context.Context, _input *MakePalette_2Input, _output *MakePalette_2Output) {

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
	positiontocolourmap, img, _ := image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, &chosencolourpalette, _input.Rotate, _input.AutoRotate)

	newFilename, err := splitFilename(_input.Imagefilename, "_plateformat")

	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	_output.OutPutImage, err = image.Export(img, newFilename)

	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	// remove duplicates
	positiontocolourmap = image.RemoveDuplicatesValuesfromMap(positiontocolourmap)

	solutions := make([]*wtype.LHComponent, 0)
	colourtoComponentMap := make(map[string]string)

	counter := 0

	for _, colour := range positiontocolourmap {

		colourindex := chosencolourpalette.Index(colour)

		if colour != nil {

			cmyk := image.ColourtoCMYK(colour)

			var maxuint8 uint8 = 255

			if cmyk.C <= _input.LowerThreshold && cmyk.Y <= _input.LowerThreshold && cmyk.M <= _input.LowerThreshold && cmyk.K <= _input.LowerThreshold {

				continue

			} else {

				var solution *wtype.LHComponent
				var componentVols []wunit.Volume

				counter = counter + 1

				if cmyk.C > _input.LowerThreshold {

					cyanvol := wunit.NewVolume(((float64(cmyk.C) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

					if cyanvol.RawValue() < 0.5 && cyanvol.Unit().PrefixedSymbol() == "ul" {
						cyanvol.SetValue(0.5)
					}

					cyanSample := mixer.Sample(_input.Cyan, cyanvol)

					// if no other components are going to be added after change the liquidtype to Mix else change to pipette above
					if wunit.AddVolumes(componentVols).EqualTo(wunit.MultiplyVolume(_input.VolumeForFullcolour, 4)) {
						cyanSample.Type = wtype.LTMegaMix
					} else {
						cyanSample.Type = wtype.LTDISPENSEABOVE

					}

					solution = execute.MixNamed(_ctx, _input.PalettePlate.Type, "", "Palette", cyanSample)

				}

				if cmyk.Y > _input.LowerThreshold {

					yellowvol := wunit.NewVolume(((float64(cmyk.Y) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

					if yellowvol.RawValue() < 0.5 && yellowvol.Unit().PrefixedSymbol() == "ul" {
						yellowvol.SetValue(0.5)
					}

					yellowSample := mixer.Sample(_input.Yellow, yellowvol)

					componentVols = append(componentVols, yellowvol)

					// if no other components are going to be added after change the liquidtype to Mix else change to pipette above
					if wunit.AddVolumes(componentVols).EqualTo(wunit.MultiplyVolume(_input.VolumeForFullcolour, 4)) {
						yellowSample.Type = wtype.LTMegaMix
					} else {
						yellowSample.Type = wtype.LTDISPENSEABOVE

					}

					if solution != nil {
						solution = execute.Mix(_ctx, solution, yellowSample)
					} else {
						solution = execute.MixNamed(_ctx, _input.PalettePlate.Type, "", "Palette", yellowSample)
					}

				}

				if cmyk.M > _input.LowerThreshold {
					magentavol := wunit.NewVolume(((float64(cmyk.M) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

					if magentavol.RawValue() < 0.5 && magentavol.Unit().PrefixedSymbol() == "ul" {
						magentavol.SetValue(0.5)
					}

					magentaSample := mixer.Sample(_input.Magenta, magentavol)

					componentVols = append(componentVols, magentavol)

					// if no other components are going to be added after change the liquidtype to Mix else change to pipette above
					if wunit.AddVolumes(componentVols).EqualTo(wunit.MultiplyVolume(_input.VolumeForFullcolour, 4)) {
						magentaSample.Type = wtype.LTMegaMix
					} else {
						magentaSample.Type = wtype.LTDISPENSEABOVE

					}

					if solution != nil {
						solution = execute.Mix(_ctx, solution, magentaSample)
					} else {
						solution = execute.MixNamed(_ctx, _input.PalettePlate.Type, "", "Palette", magentaSample)
					}

				}

				if cmyk.K > _input.LowerThreshold {

					blackvol := wunit.NewVolume(((float64(cmyk.K) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

					if blackvol.RawValue() < 0.5 && blackvol.Unit().PrefixedSymbol() == "ul" {
						blackvol.SetValue(0.5)
					}

					blackSample := mixer.Sample(_input.Black, blackvol)

					componentVols = append(componentVols, blackvol)

					// if no other components are going to be added after change the liquidtype to Mix else change to pipette above
					if wunit.AddVolumes(componentVols).EqualTo(wunit.MultiplyVolume(_input.VolumeForFullcolour, 4)) {
						blackSample.Type = wtype.LTMegaMix
					} else {
						blackSample.Type = wtype.LTDISPENSEABOVE

					}

					if solution != nil {
						solution = execute.Mix(_ctx, solution, blackSample)
					} else {
						solution = execute.MixNamed(_ctx, _input.PalettePlate.Type, "", "Palette", blackSample)
					}

				}

				// top up colour to 4 x volumeforfullcolour with white to make the correct shade

				// calculate volume of white to add
				whitevol := wunit.SubtractVolumes(wunit.MultiplyVolume(_input.VolumeForFullcolour, 4), componentVols)

				// mix with white sample
				_input.White.Type = wtype.LTMegaMix

				whiteSample := mixer.Sample(_input.White, whitevol)

				if solution != nil {
					solution = execute.Mix(_ctx, solution, whiteSample)
				} else if _input.NotThisColour == "white" {
					// skip
				} else {
					solution = execute.MixNamed(_ctx, _input.PalettePlate.Type, "", "Palette", whiteSample)
				}
				// change name of component
				originalname := solution.CName
				solution.CName = originalname + "_colour_" + strconv.Itoa(colourindex)

				// add solution to be exported later
				solutions = append(solutions, solution)
				colourtoComponentMap[strconv.Itoa(colourindex)] = solution.CName

			}

		}
	}

	_output.Colours = solutions
	_output.Numberofcolours = len(chosencolourpalette)
	_output.Palette = chosencolourpalette
	_output.ColourtoComponentMap = colourtoComponentMap

	_output.PaletteFile, err = export.JSON(_output.Palette, "Palette.json")

	if err != nil {
		execute.Errorf(_ctx, "Error exporting palette to json")
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MakePalette_2Analysis(_ctx context.Context, _input *MakePalette_2Input, _output *MakePalette_2Output) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _MakePalette_2Validation(_ctx context.Context, _input *MakePalette_2Input, _output *MakePalette_2Output) {

}

func splitFilename(filename, addition string) (newfilename string, err error) {
	fields := strings.Split(filename, `.`)

	if len(fields) <= 1 {
		return filename + addition, fmt.Errorf("filename has no dot so added addition on to end")
	} else if len(fields) == 2 {
		// rename file
		newfilename = fields[0] + addition + `.` + fields[1]
	} else if len(fields) > 2 {
		// rename file
		firstPart := strings.Join(fields[0:len(fields)-2], `.`)

		newfilename = firstPart + addition + `.` + fields[len(fields)-1]
	}
	return
}
func _MakePalette_2Run(_ctx context.Context, input *MakePalette_2Input) *MakePalette_2Output {
	output := &MakePalette_2Output{}
	_MakePalette_2Setup(_ctx, input)
	_MakePalette_2Steps(_ctx, input, output)
	_MakePalette_2Analysis(_ctx, input, output)
	_MakePalette_2Validation(_ctx, input, output)
	return output
}

func MakePalette_2RunSteps(_ctx context.Context, input *MakePalette_2Input) *MakePalette_2SOutput {
	soutput := &MakePalette_2SOutput{}
	output := _MakePalette_2Run(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MakePalette_2New() interface{} {
	return &MakePalette_2Element{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MakePalette_2Input{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MakePalette_2Run(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MakePalette_2Input{},
			Out: &MakePalette_2Output{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type MakePalette_2Element struct {
	inject.CheckedRunner
}

type MakePalette_2Input struct {
	AutoRotate          bool
	Black               *wtype.LHComponent
	Cyan                *wtype.LHComponent
	Imagefilename       string
	LowerThreshold      uint8
	Magenta             *wtype.LHComponent
	NotThisColour       string
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

type MakePalette_2Output struct {
	Colours              []*wtype.LHComponent
	ColourtoComponentMap map[string]string
	Numberofcolours      int
	OutPutImage          wtype.File
	Palette              color.Palette
	PaletteFile          wtype.File
}

type MakePalette_2SOutput struct {
	Data struct {
		ColourtoComponentMap map[string]string
		Numberofcolours      int
		OutPutImage          wtype.File
		Palette              color.Palette
		PaletteFile          wtype.File
	}
	Outputs struct {
		Colours []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MakePalette_2",
		Constructor: MakePalette_2New,
		Desc: component.ComponentDesc{
			Desc: "Generates instructions to make a pallette of all colours in an image\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/PipetteImage/PipetteImage/fromPalette/MakePalette_2.an",
			Params: []component.ParamDesc{
				{Name: "AutoRotate", Desc: "", Kind: "Parameters"},
				{Name: "Black", Desc: "", Kind: "Inputs"},
				{Name: "Cyan", Desc: "", Kind: "Inputs"},
				{Name: "Imagefilename", Desc: "name of image file or if using URL use this field to set the desired filename\n", Kind: "Parameters"},
				{Name: "LowerThreshold", Desc: "", Kind: "Parameters"},
				{Name: "Magenta", Desc: "", Kind: "Inputs"},
				{Name: "NotThisColour", Desc: "", Kind: "Parameters"},
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
				{Name: "ColourtoComponentMap", Desc: "map of colour name (as index) to component name\n", Kind: "Data"},
				{Name: "Numberofcolours", Desc: "", Kind: "Data"},
				{Name: "OutPutImage", Desc: "", Kind: "Data"},
				{Name: "Palette", Desc: "Colournames []string\n", Kind: "Data"},
				{Name: "PaletteFile", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
