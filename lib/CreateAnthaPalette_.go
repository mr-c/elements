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
)

// Input parameters for this protocol (data)

//AnthaPalette to be generated
//ColorVolume		Volume				//Volume of each color on the palette you want to generate
//PalettePlate	wtype.LHPlate		//Plate on which the palette will be generated
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
	var maxRGB = 196605

	//TODO Those globals are set up for testing
	ColorVolume := wunit.NewVolume(10, "ul")

	PalettePlate := factory.GetPlateByType("pcrplate_skirted")

	//----------------------------------------------------------------------------------
	//Generating colors
	//----------------------------------------------------------------------------------

	//iterate over each colors in the palette
	for i := range _input.AnthaPalette.Palette {

		//placeholder for the current color we are iterating over
		var currentAnthaColor = _input.AnthaPalette.AnthaColors[i]
		//PlaceHolder for the LHComponent solution to be pipetted into the well
		var solution *wtype.LHComponent
		//PlaceHolder for the total well volume
		var totalVolume []wunit.Volume

		//extract RGBA values
		r, g, b, a := currentAnthaColor.Color.RGBA()

		//getting the volume of RGB needed
		vr := wunit.NewVolume((float64(r)/float64(maxRGB))*ColorVolume.RawValue(), ColorVolume.Unit().PrefixedSymbol())
		vg := wunit.NewVolume((float64(g)/float64(maxRGB))*ColorVolume.RawValue(), ColorVolume.Unit().PrefixedSymbol())
		vb := wunit.NewVolume((float64(b)/float64(maxRGB))*ColorVolume.RawValue(), ColorVolume.Unit().PrefixedSymbol())

		fmt.Println(vr, vg, vb, a)

		//Check if the values for each colors is lower than the lowerThreshold defined
		if r <= uint32(_input.LowerThreshold) && g <= uint32(_input.LowerThreshold) && b <= uint32(_input.LowerThreshold) {
			continue
		} else {
			//----------------------------------------------------------------------------------
			//Pipetting red
			//----------------------------------------------------------------------------------

			//Initiating the LHComponents to pipette with volume information
			rSample := mixer.Sample(currentAnthaColor.Component, vr)

			// if no other components are going to be added afterward, change the liquidtype to Mix else change to pipette above
			//TODO: why do you multiply by 4
			if wunit.AddVolumes(totalVolume).EqualTo(wunit.MultiplyVolume(ColorVolume, 4)) {
				rSample.Type = wtype.LTMegaMix
			} else {
				rSample.Type = wtype.LTDISPENSEABOVE
			}

			//since this is the first color to be pipetted, we use MixNamed to instantiate the LHComponent
			solution = execute.MixNamed(_ctx, PalettePlate.Type, "", "Palette", rSample)

			//----------------------------------------------------------------------------------
			//Pipetting green
			//----------------------------------------------------------------------------------

			//Initiating the LHComponents to pipette with volume information
			gSample := mixer.Sample(currentAnthaColor.Component, vg)

			// if no other components are going to be added afterward, change the liquidtype to Mix else change to pipette above
			//TODO: why do you multiply by 4
			if wunit.AddVolumes(totalVolume).EqualTo(wunit.MultiplyVolume(ColorVolume, 4)) {
				gSample.Type = wtype.LTMegaMix
			} else {
				gSample.Type = wtype.LTDISPENSEABOVE
			}

			//if the solution is already instantiated, we add it with the Mix() command. Else we make a new one with MixNamed()
			if solution != nil {
				solution = execute.Mix(_ctx, solution, gSample)
			} else {
				solution = execute.MixNamed(_ctx, PalettePlate.Type, "", "Palette", gSample)
			}

			//----------------------------------------------------------------------------------
			//Pipetting Blue
			//----------------------------------------------------------------------------------
			//Initiating the LHComponents to pipette with volume information
			bSample := mixer.Sample(currentAnthaColor.Component, vb)

			// if no other components are going to be added afterward, change the liquidtype to Mix else change to pipette above
			//TODO: why do you multiply by 4
			if wunit.AddVolumes(totalVolume).EqualTo(wunit.MultiplyVolume(ColorVolume, 4)) {
				bSample.Type = wtype.LTMegaMix
			} else {
				bSample.Type = wtype.LTDISPENSEABOVE
			}

			//if the solution is already instantiated, we add it with the Mix() command. Else we make a new one with MixNamed()
			if solution != nil {
				solution = execute.Mix(_ctx, solution, bSample)
			} else {
				solution = execute.MixNamed(_ctx, PalettePlate.Type, "", "Palette", bSample)
			}

			//adding the final created LHComponent to the AnthaColor (since it has the added mixing information)
			_input.AnthaPalette.AnthaColors[i].Component = solution

		}
	}

	//returning the AnthaPalette with updated LHComponents
	_output.MixedAnthaPalette = _input.AnthaPalette

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
	LowerThreshold int
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
				{Name: "LowerThreshold", Desc: "ColorVolume\t\tVolume\t\t\t\t//Volume of each color on the palette you want to generate\nPalettePlate\twtype.LHPlate\t\t//Plate on which the palette will be generated\n\nRGB value below which we do not create a color (we consider it nil). The number is between 0 and 65535.\n", Kind: "Parameters"},
				{Name: "MixedAnthaPalette", Desc: "The palette with physical location information added to the LHComponents.\n", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
