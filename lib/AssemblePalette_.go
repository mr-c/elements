// Generates instructions to make a pallette of all colours in an image
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
	goimage "image"
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

//InPlate *LHPlate

// Physical outputs from this protocol with types

func _AssemblePaletteRequirements() {

}

// Conditions to run on startup
func _AssemblePaletteSetup(_ctx context.Context, _input *AssemblePaletteInput) {

}

type RBSdata struct {
	Seq      wtype.DNASequence
	Strength float64
}

type collection []RBSdata

func (rbscollection collection) Max() (rbs RBSdata) {

	//var maxSeq DNASequence
	var maxStrength float64
	for i := range rbscollection {
		if i == 0 {
			maxStrength = rbscollection[i].Strength
			//maxSeq = rbscollection[i].Seq
			rbs = rbscollection[i]
		} else if rbscollection[i].Strength > maxStrength {
			maxStrength = rbscollection[i].Strength
			//maxSeq = rbscollection[i].Seq
			rbs = rbscollection[i]
		}
	}
	return
}

// The core process for this protocol, with the steps to be performed
// for every input
func _AssemblePaletteSteps(_ctx context.Context, _input *AssemblePaletteInput, _output *AssemblePaletteOutput) {

	//-----------------------------------------------------------------------------
	//Globals
	//-----------------------------------------------------------------------------

	var rbsstrengthdata []RBSdata = []RBSdata{
		{wtype.MakeLinearDNASequence("rbs1", "gggcgcgc"), 0.0},
		{wtype.MakeLinearDNASequence("rbs2", "gggcgcgc"), 2.0},
		{wtype.MakeLinearDNASequence("rbs3", "gggcgcgc"), 5.0},
		{wtype.MakeLinearDNASequence("rbs4", "gggcgcgc"), 10.0},
	}

	redname := _input.Red.CName
	greenname := _input.Green.CName
	bluename := _input.Blue.CName

	fmt.Println(rbsstrengthdata)

	//image and error placeholders

	var imgFile wtype.File
	var imgBase *goimage.NRGBA
	var err error

	//----------------------------------------------------------------------------
	//Fetch image
	//----------------------------------------------------------------------------

	// if image is from url, download
	if _input.UseURL {

		//downloading image
		imgFile, err = download.File(_input.URL, _input.Imagefilename)
		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}

		//opening the image file
		imgBase, err = image.OpenFile(imgFile)
		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}
	}

	//----------------------------------------------------------------------------
	//Processing image
	//----------------------------------------------------------------------------

	if _input.PosterizeImage {
		imgBase, err = image.Posterize(imgBase, _input.PosterizeLevels)
		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}
	}

	//----------------------------------------------------------------------------
	//choosing palette
	//----------------------------------------------------------------------------

	// make palette of colours from image
	chosencolourpalette := image.MakeSmallPalleteFromImage(imgBase, _input.PlateWithMasterMix, _input.Rotate)

	// make a map of colour to well coordinates
	positiontocolourmap, _ := image.ImagetoPlatelayout(imgBase, _input.PlateWithMasterMix, &chosencolourpalette, _input.Rotate, _input.AutoRotate)

	// remove duplicates
	positiontocolourmap = image.RemoveDuplicatesValuesfromMap(positiontocolourmap)

	solutions := make([]*wtype.LHComponent, 0)
	colourtoComponentMap := make(map[string]*wtype.LHComponent)

	counter := 0

	for _, colour := range positiontocolourmap {

		colourindex := chosencolourpalette.Index(colour)

		if colour != nil {
			components := make([]*wtype.LHComponent, 0)

			r, g, b, _ := colour.RGBA()

			//var maxuint8 uint8 = 255

			if r == 0 && g == 0 && b == 0 {

				continue

			} else {

				counter = counter + 1

				_input.Red.CName = fmt.Sprint(redname, "_RBS_", uint8(r))
				_input.Red.Type = wtype.LTPostMix

				redSample := mixer.Sample(_input.Red, _input.VolumeForeachColourPlasmid)
				components = append(components, redSample)

				_input.Green.CName = fmt.Sprint(greenname, "_RBS_", uint8(g))
				_input.Green.Type = wtype.LTPostMix

				greenSample := mixer.Sample(_input.Green, _input.VolumeForeachColourPlasmid)

				components = append(components, greenSample)

				_input.Blue.CName = fmt.Sprint(bluename, "_RBS_", uint8(b))
				_input.Blue.Type = wtype.LTPostMix

				blueSample := mixer.Sample(_input.Blue, _input.VolumeForeachColourPlasmid)

				components = append(components, blueSample)

				solution := execute.MixInto(_ctx, _input.PalettePlate, "", components...)

				/*
					dnaSample := mixer.Sample(solution,wunit.NewVolume(1,"ul"))

					transformation := MixInto(factory.GetPlateByType("pcrplate_with_cooler"),"",dnaSample)

					transformationSample := mixer.Sample(transformation,wunit.NewVolume(20,"ul"))

					recovery := MixInto(factory.GetPlateByType("DSW96_riser"),"",transformationSample)

					solutions = append(solutions, recovery)
					colourtoComponentMap[strconv.Itoa(colourindex)] = recovery
				*/

				solutions = append(solutions, solution)
				colourtoComponentMap[strconv.Itoa(colourindex)] = solution

			}

		}
	}

	_output.Colours = solutions
	_output.Numberofcolours = len(chosencolourpalette)
	_output.Palette = chosencolourpalette
	_output.ColourtoComponentMap = colourtoComponentMap
	//fmt.Println("Unique Colours =",Numberofcolours,"from palette:", chosencolourpalette)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AssemblePaletteAnalysis(_ctx context.Context, _input *AssemblePaletteInput, _output *AssemblePaletteOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AssemblePaletteValidation(_ctx context.Context, _input *AssemblePaletteInput, _output *AssemblePaletteOutput) {

}
func _AssemblePaletteRun(_ctx context.Context, input *AssemblePaletteInput) *AssemblePaletteOutput {
	output := &AssemblePaletteOutput{}
	_AssemblePaletteSetup(_ctx, input)
	_AssemblePaletteSteps(_ctx, input, output)
	_AssemblePaletteAnalysis(_ctx, input, output)
	_AssemblePaletteValidation(_ctx, input, output)
	return output
}

func AssemblePaletteRunSteps(_ctx context.Context, input *AssemblePaletteInput) *AssemblePaletteSOutput {
	soutput := &AssemblePaletteSOutput{}
	output := _AssemblePaletteRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AssemblePaletteNew() interface{} {
	return &AssemblePaletteElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AssemblePaletteInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AssemblePaletteRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AssemblePaletteInput{},
			Out: &AssemblePaletteOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type AssemblePaletteElement struct {
	inject.CheckedRunner
}

type AssemblePaletteInput struct {
	AutoRotate                 bool
	Blue                       *wtype.LHComponent
	Green                      *wtype.LHComponent
	Imagefilename              string
	PalettePlate               *wtype.LHPlate
	PlateWithMasterMix         *wtype.LHPlate
	PosterizeImage             bool
	PosterizeLevels            int
	Red                        *wtype.LHComponent
	Rotate                     bool
	URL                        string
	UseURL                     bool
	VolumeForeachColourPlasmid wunit.Volume
}

type AssemblePaletteOutput struct {
	Colours              []*wtype.LHComponent
	ColourtoComponentMap map[string]*wtype.LHComponent
	Numberofcolours      int
	Palette              color.Palette
}

type AssemblePaletteSOutput struct {
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
	if err := addComponent(component.Component{Name: "AssemblePalette",
		Constructor: AssemblePaletteNew,
		Desc: component.ComponentDesc{
			Desc: "Generates instructions to make a pallette of all colours in an image\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/PipetteImage/AssemblePalette.an",
			Params: []component.ParamDesc{
				{Name: "AutoRotate", Desc: "", Kind: "Parameters"},
				{Name: "Blue", Desc: "", Kind: "Inputs"},
				{Name: "Green", Desc: "", Kind: "Inputs"},
				{Name: "Imagefilename", Desc: "name of image file or if using URL use this field to set the desired filename\n", Kind: "Parameters"},
				{Name: "PalettePlate", Desc: "", Kind: "Inputs"},
				{Name: "PlateWithMasterMix", Desc: "InPlate *LHPlate\n", Kind: "Inputs"},
				{Name: "PosterizeImage", Desc: "", Kind: "Parameters"},
				{Name: "PosterizeLevels", Desc: "", Kind: "Parameters"},
				{Name: "Red", Desc: "", Kind: "Inputs"},
				{Name: "Rotate", Desc: "", Kind: "Parameters"},
				{Name: "URL", Desc: "enter URL link to the image file here if applicable\n", Kind: "Parameters"},
				{Name: "UseURL", Desc: "select this if getting the image from a URL\n", Kind: "Parameters"},
				{Name: "VolumeForeachColourPlasmid", Desc: "", Kind: "Parameters"},
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
