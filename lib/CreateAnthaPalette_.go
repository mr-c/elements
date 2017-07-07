//This element will work out and perform the dilutions necessary to make an antha palette.
package lib

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/factory"
	"image/color"
)

// Input parameters for this protocol (data)

//AnthaPalette to be generated
//Volume of each color on the palette you want to generate
//Plate on which the palette will be generated
//RGB value below which we do not create a color (we consider it nil). The number is between 0 and 65535.

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

//The palette with physical location information added to the LHComponents.

func _CreateAnthaPaletteRequirements() {

}

// Conditions to run on startup
func _CreateAnthaPaletteSetup(_ctx context.Context, _input *CreateAnthaPaletteInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _CreateAnthaPaletteSteps(_ctx context.Context, _input *CreateAnthaPaletteInput, _output *CreateAnthaPaletteOutput) {
	//----------------------------------------------------------------------------------
	//Globals
	//----------------------------------------------------------------------------------

	//This is the maximum reached by summing RGB values. We use this as a reference to get percentages values for each channels
	var maxCMYK uint8 = 255

	//This is the minimum volume below which the sample() functions would return nil.
	var minVolume = wunit.NewVolume(0.5, "ul")

	//Placeholder for the current number of colors added to a solution.
	colorsAdded := 1

	//----------------------------------------------------------------------------------
	//Generating colors
	//----------------------------------------------------------------------------------

	//iterate over each colors in the palette
	for i, AnthaColor := range _input.AnthaPalette.AnthaColors {

		//PlaceHolder for the LHComponent solution to be pipetted into the well
		var solution *wtype.LHComponent

		//PlaceHolder for the total well volume
		//var totalVolume []Volume
		//Boolean checking if the solution being made for the current color is initialized
		solutionInitialized := false

		//extract RGB values. We do not use the a (alpha) value
		r, g, b, _ := AnthaColor.Color.RGBA()

		//convert to the CMYK model since that is what is used to generate colors from paint
		c, m, y, k := color.RGBToCMYK(uint8(r), uint8(g), uint8(b))

		//figure out the volumes needed for the color and adding them to a map for easier iteration
		volumes := map[string]wunit.Volume{
			"c": wunit.NewVolume((float64(c)/float64(maxCMYK))*_input.ColorVolume.RawValue(), _input.ColorVolume.Unit().PrefixedSymbol()),
			"m": wunit.NewVolume((float64(m)/float64(maxCMYK))*_input.ColorVolume.RawValue(), _input.ColorVolume.Unit().PrefixedSymbol()),
			"y": wunit.NewVolume((float64(y)/float64(maxCMYK))*_input.ColorVolume.RawValue(), _input.ColorVolume.Unit().PrefixedSymbol()),
			"k": wunit.NewVolume((float64(k)/float64(maxCMYK))*_input.ColorVolume.RawValue(), _input.ColorVolume.Unit().PrefixedSymbol()),
		}

		//Check if the values for each colors is lower than the lowerThreshold defined, if they are we do not produce this palette color
		if r <= uint32(_input.LowerThreshold) && g <= uint32(_input.LowerThreshold) && b <= uint32(_input.LowerThreshold) {
			continue
		} else {

			//placeholders for the CMYK indidual color samples
			var sample *wtype.LHComponent

			//We range over the calculated RGB volumes and produce each of the desired colors.
			for _, volume := range volumes {
				switch solutionInitialized {
				case false:
					//Initiating the LHComponents to pipette with volume information
					//we need to check if the volumes are lower than 0.05ul because otherwise the sample() function would return nil.
					if minVolume.GreaterThan(volume) {
						sample = mixer.Sample(AnthaColor.Component, minVolume)
					} else {
						sample = mixer.Sample(AnthaColor.Component, volume)
					}

					//since this is the first color to be pipetted, we use MixNamed to instantiate the LHComponent
					solution = execute.MixNamed(_ctx, _input.PalettePlate.Type, "", "Palette", sample)

					//change to pipette above to signal LHComponents are still to be added
					sample.Type = wtype.LTDISPENSEABOVE
					//adding to the added color counter
					colorsAdded++

					//signalling that the solution is intitialized
					solutionInitialized = true

				default:
					//Initiating the LHComponents to pipette with volume information
					//we need to check if the volumes are lower than 0.05ul because otherwise the sample() function would return nil.
					if minVolume.GreaterThan(volume) {
						sample = mixer.Sample(AnthaColor.Component, minVolume)
					} else {
						sample = mixer.Sample(AnthaColor.Component, volume)
					}

					//adding the sample to the solution
					solution = execute.Mix(_ctx, solution, sample)

					//checking if no more colors need to be added, and changing the information appropriately. Since we use the CMYK model this is four.
					if colorsAdded != 4 {
						//change to pipette above to signal LHComponents are still to be added
						sample.Type = wtype.LTDISPENSEABOVE
						//adding to the added color counter
						colorsAdded++
					} else {
						sample.Type = wtype.LTMegaMix

						//adding the final created LHComponent to the AnthaColor (since it has the added mixing information)
						_input.AnthaPalette.AnthaColors[i].Component = solution

						//resetting counters
						colorsAdded = 1
						solutionInitialized = false
					}
				}
			}
		}
	}

	//returning the AnthaPalette with updated LHComponents
	_output.MixedAnthaPalette = _input.AnthaPalette

	fmt.Println("Palette Made")
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _CreateAnthaPaletteAnalysis(_ctx context.Context, _input *CreateAnthaPaletteInput, _output *CreateAnthaPaletteOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _CreateAnthaPaletteValidation(_ctx context.Context, _input *CreateAnthaPaletteInput, _output *CreateAnthaPaletteOutput) {

}
func _CreateAnthaPaletteRun(_ctx context.Context, input *CreateAnthaPaletteInput) *CreateAnthaPaletteOutput {
	output := &CreateAnthaPaletteOutput{}
	_CreateAnthaPaletteSetup(_ctx, input)
	_CreateAnthaPaletteSteps(_ctx, input, output)
	_CreateAnthaPaletteAnalysis(_ctx, input, output)
	_CreateAnthaPaletteValidation(_ctx, input, output)
	return output
}

func CreateAnthaPaletteRunSteps(_ctx context.Context, input *CreateAnthaPaletteInput) *CreateAnthaPaletteSOutput {
	soutput := &CreateAnthaPaletteSOutput{}
	output := _CreateAnthaPaletteRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func CreateAnthaPaletteNew() interface{} {
	return &CreateAnthaPaletteElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &CreateAnthaPaletteInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _CreateAnthaPaletteRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &CreateAnthaPaletteInput{},
			Out: &CreateAnthaPaletteOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type CreateAnthaPaletteElement struct {
	inject.CheckedRunner
}

type CreateAnthaPaletteInput struct {
	AnthaPalette   *image.AnthaPalette
	ColorVolume    wunit.Volume
	LowerThreshold int
	PalettePlate   wtype.LHPlate
}

type CreateAnthaPaletteOutput struct {
	MixedAnthaPalette *image.AnthaPalette
}

type CreateAnthaPaletteSOutput struct {
	Data struct {
	}
	Outputs struct {
		MixedAnthaPalette *image.AnthaPalette
	}
}

func init() {
	if err := addComponent(component.Component{Name: "CreateAnthaPalette",
		Constructor: CreateAnthaPaletteNew,
		Desc: component.ComponentDesc{
			Desc: "This element will work out and perform the dilutions necessary to make an antha palette.\n",
			Path: "src/github.com/antha-lang/elements/an/ImageHandling/CreateAnthaPalette.an",
			Params: []component.ParamDesc{
				{Name: "AnthaPalette", Desc: "AnthaPalette to be generated\n", Kind: "Parameters"},
				{Name: "ColorVolume", Desc: "Volume of each color on the palette you want to generate\n", Kind: "Parameters"},
				{Name: "LowerThreshold", Desc: "RGB value below which we do not create a color (we consider it nil). The number is between 0 and 65535.\n", Kind: "Parameters"},
				{Name: "PalettePlate", Desc: "Plate on which the palette will be generated\n", Kind: "Parameters"},
				{Name: "MixedAnthaPalette", Desc: "The palette with physical location information added to the LHComponents.\n", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
