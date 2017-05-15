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
	"github.com/antha-lang/antha/microArch/factory"
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

//InPlate *wtype.LHPlate

// Physical outputs from this protocol with types

func _AssemblePalette_OneByOne_RGB_transform_2Requirements() {

}

// Conditions to run on startup
func _AssemblePalette_OneByOne_RGB_transform_2Setup(_ctx context.Context, _input *AssemblePalette_OneByOne_RGB_transform_2Input) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _AssemblePalette_OneByOne_RGB_transform_2Steps(_ctx context.Context, _input *AssemblePalette_OneByOne_RGB_transform_2Input, _output *AssemblePalette_OneByOne_RGB_transform_2Output) {

	//------------------------------------------------------------
	//Globals
	//------------------------------------------------------------

	var (
		ReactionTemp                wunit.Temperature = wunit.NewTemperature(25, "C")
		ReactionTime                wunit.Time        = wunit.NewTime(35, "min")
		OutputLocation              string            = ""
		CompetentCellPlateWell      string            = ""
		RecoveryPlateWell           string            = ""
		RecoveryTemp                wunit.Temperature = wunit.NewTemperature(37.0, "C")
		RecoveryTime                wunit.Time        = wunit.NewTime(2.0, "h")
		TransformationVolume        wunit.Volume      = wunit.NewVolume(2.0, "ul")
		PostPlasmidTemp             wunit.Temperature = wunit.NewTemperature(2.0, "C")
		PostPlasmidTime             wunit.Time        = wunit.NewTime(5, "min")
		CompetentCellTransferVolume wunit.Volume      = wunit.NewVolume(20.0, "ul")
		RecoveryPlateNumber         int               = 1

		PlatewithRecoveryMedia  *wtype.LHPlate = factory.GetPlateByType("pcrplate_skirted_riser")
		PlateWithCompetentCells *wtype.LHPlate = factory.GetPlateByType("pcrplate_skirted_riser")
	)

	redname := _input.Red.CName
	greenname := _input.Green.CName
	bluename := _input.Blue.CName

	var imgFile *wtype.File
	var imgBase *goimage.NRGBA

	//--------------------------------------------------------------
	//Fetching image
	//--------------------------------------------------------------

	// if image is from url, download
	if _input.UseURL {
		//downloading image
		imgFile, err := download.File(_input.URL, _input.Imagefilename)
		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}

		//opening the image file
		imgBase, err := image.OpenFile(imgFile)
		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}
	}

	//--------------------------------------------------------------
	//Image Processing
	//--------------------------------------------------------------

	if _input.PosterizeImage {
		imgBase, err := image.Posterize(imgBase, _input.PosterizeLevels)
		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}
	}

	//--------------------------------------------------------------
	//Choosing Palette
	//--------------------------------------------------------------

	// make palette of colours from image
	chosencolourpalette := image.MakeSmallPalleteFromImage(imgBase, _input.PlateWithMasterMix, _input.Rotate)

	// make a map of colour to well coordinates
	positiontocolourmap, _ := image.ImagetoPlatelayout(imgBase, _input.PlateWithMasterMix, &chosencolourpalette, _input.Rotate, _input.AutoRotate)

	// remove duplicates
	positiontocolourmap = image.RemoveDuplicatesValuesfromMap(positiontocolourmap)

	//fmt.Println("positions", positiontocolourmap)

	solutions := make([]*wtype.LHComponent, 0)
	colourtoComponentMap := make(map[string]*wtype.LHComponent)

	counter := 0
	//wellpositionarray := PlateWithMasterMix.AllWellPositions(wtype.BYCOLUMN)

	for _, colour := range positiontocolourmap {

		colourindex := chosencolourpalette.Index(colour)

		if colour != nil {
			components := make([]*wtype.LHComponent, 0)

			r, g, b, _ := colour.RGBA()

			//var maxuint8 uint8 = 255

			if r == 0 && g == 0 && b == 0 {

				continue

			} else {

				//OutputLocation 			= wellpositionarray[counter]
				//CompetentCellPlateWell 	= wellpositionarray[counter]
				//RecoveryPlateWell 		= wellpositionarray[counter]

				counter = counter + 1

				/*
							MasterMix.Type,err = wtype.LiquidTypeFromString(LHPolicyName)

							if err != nil {
					                 Errorf("cannot find liquid type: %s", err)
					    	}

							mmxSample:=mixer.Sample(MasterMix,MasterMixVolume)

							components = append(components, mmxSample)
				*/
				_input.Red.CName = fmt.Sprint(redname, "_RBSstrength_", uint8(r))
				_input.Red.Type = wtype.LTPostMix

				redSample := mixer.Sample(_input.Red, _input.VolumeForeachColourPlasmid)
				components = append(components, redSample)

				_input.Green.CName = fmt.Sprint(greenname, "_RBSstrength_", uint8(g))
				_input.Green.Type = wtype.LTPostMix

				greenSample := mixer.Sample(_input.Green, _input.VolumeForeachColourPlasmid)

				components = append(components, greenSample)

				_input.Blue.CName = fmt.Sprint(bluename, "_RBSstrength_", uint8(b))
				_input.Blue.Type = wtype.LTPostMix

				blueSample := mixer.Sample(_input.Blue, _input.VolumeForeachColourPlasmid)

				components = append(components, blueSample)

				reaction := execute.MixTo(_ctx, _input.PalettePlate.Type, OutputLocation, 1, components...)

				dnaSample := mixer.Sample(reaction, TransformationVolume)

				dnaSample.Type = wtype.LTCulture

				execute.Incubate(_ctx, dnaSample, ReactionTemp, ReactionTime, false)

				transformation := execute.MixTo(_ctx, PlateWithCompetentCells.Type, CompetentCellPlateWell, 1, dnaSample)

				transformation.Type = wtype.LTCulture

				execute.Incubate(_ctx, transformation, PostPlasmidTemp, PostPlasmidTime, false)

				transformationSample := mixer.Sample(transformation, CompetentCellTransferVolume)

				solution := execute.MixTo(_ctx, PlatewithRecoveryMedia.Type, RecoveryPlateWell, RecoveryPlateNumber, transformationSample)

				// incubate the reaction mixture
				// commented out pending changes to incubate
				execute.Incubate(_ctx, solution, RecoveryTemp, RecoveryTime, true)

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
func _AssemblePalette_OneByOne_RGB_transform_2Analysis(_ctx context.Context, _input *AssemblePalette_OneByOne_RGB_transform_2Input, _output *AssemblePalette_OneByOne_RGB_transform_2Output) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AssemblePalette_OneByOne_RGB_transform_2Validation(_ctx context.Context, _input *AssemblePalette_OneByOne_RGB_transform_2Input, _output *AssemblePalette_OneByOne_RGB_transform_2Output) {

}
func _AssemblePalette_OneByOne_RGB_transform_2Run(_ctx context.Context, input *AssemblePalette_OneByOne_RGB_transform_2Input) *AssemblePalette_OneByOne_RGB_transform_2Output {
	output := &AssemblePalette_OneByOne_RGB_transform_2Output{}
	_AssemblePalette_OneByOne_RGB_transform_2Setup(_ctx, input)
	_AssemblePalette_OneByOne_RGB_transform_2Steps(_ctx, input, output)
	_AssemblePalette_OneByOne_RGB_transform_2Analysis(_ctx, input, output)
	_AssemblePalette_OneByOne_RGB_transform_2Validation(_ctx, input, output)
	return output
}

func AssemblePalette_OneByOne_RGB_transform_2RunSteps(_ctx context.Context, input *AssemblePalette_OneByOne_RGB_transform_2Input) *AssemblePalette_OneByOne_RGB_transform_2SOutput {
	soutput := &AssemblePalette_OneByOne_RGB_transform_2SOutput{}
	output := _AssemblePalette_OneByOne_RGB_transform_2Run(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AssemblePalette_OneByOne_RGB_transform_2New() interface{} {
	return &AssemblePalette_OneByOne_RGB_transform_2Element{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AssemblePalette_OneByOne_RGB_transform_2Input{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AssemblePalette_OneByOne_RGB_transform_2Run(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AssemblePalette_OneByOne_RGB_transform_2Input{},
			Out: &AssemblePalette_OneByOne_RGB_transform_2Output{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type AssemblePalette_OneByOne_RGB_transform_2Element struct {
	inject.CheckedRunner
}

type AssemblePalette_OneByOne_RGB_transform_2Input struct {
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

type AssemblePalette_OneByOne_RGB_transform_2Output struct {
	Colours              []*wtype.LHComponent
	ColourtoComponentMap map[string]*wtype.LHComponent
	Numberofcolours      int
	Palette              color.Palette
}

type AssemblePalette_OneByOne_RGB_transform_2SOutput struct {
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
	if err := addComponent(component.Component{Name: "AssemblePalette_OneByOne_RGB_transform_2",
		Constructor: AssemblePalette_OneByOne_RGB_transform_2New,
		Desc: component.ComponentDesc{
			Desc: "Generates instructions to make a pallette of all colours in an image\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/PipetteImage/AssemblePalette_transform.an",
			Params: []component.ParamDesc{
				{Name: "AutoRotate", Desc: "", Kind: "Parameters"},
				{Name: "Blue", Desc: "", Kind: "Inputs"},
				{Name: "Green", Desc: "", Kind: "Inputs"},
				{Name: "Imagefilename", Desc: "name of image file or if using URL use this field to set the desired filename\n", Kind: "Parameters"},
				{Name: "PalettePlate", Desc: "", Kind: "Inputs"},
				{Name: "PlateWithMasterMix", Desc: "InPlate *wtype.LHPlate\n", Kind: "Inputs"},
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
