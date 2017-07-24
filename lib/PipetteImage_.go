// Generates instructions to pipette out a defined image onto a defined plate using a defined palette of coloured bacteria
package lib

import (
	"context"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	goimage "image"
	"image/color"
	"sort"
	"strings"
)

// Input parameters for this protocol (data)

//Desired name for the output image file
// Image File

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _PipetteImageRequirements() {

}

// Conditions to run on startup
func _PipetteImageSetup(_ctx context.Context, _input *PipetteImageInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _PipetteImageSteps(_ctx context.Context, _input *PipetteImageInput, _output *PipetteImageOutput) {

	//--------------------------------------------------------------
	//Globals
	//--------------------------------------------------------------

	var imgBase *goimage.NRGBA
	var err error

	//--------------------------------------------------------------
	//Opening image
	//--------------------------------------------------------------

	//opening the image file
	imgBase, err = image.OpenFile(_input.InputFile)
	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	//--------------------------------------------------------------
	//Choosing palette
	//--------------------------------------------------------------

	// check that palette name is valid
	_, ok := image.AvailablePalettes()[_input.Palettename]

	if !ok {
		var validpalettes []string

		for key := range image.AvailablePalettes() {
			validpalettes = append(validpalettes, key)
		}

		sort.Strings(validpalettes)

		execute.Errorf(_ctx, "Palette", _input.Palettename, "not available. Valid entries are: ", strings.Join(validpalettes, ","))
	}

	// make sub pallete if necessary
	var chosencolourpalette color.Palette

	if _input.Subset {
		chosencolourpalette = image.MakeSubPallette(_input.Palettename, _input.Subsetnames)
	} else {
		chosencolourpalette = image.AvailablePalettes()[_input.Palettename]
	}

	//--------------------------------------------------------------
	//Image processing
	//--------------------------------------------------------------

	if _input.CheckResizeAlgorithms {
		_output.ResizedImages = image.CheckAllResizealgorithms(imgBase, _input.OutPlate, _input.Rotate, image.AllResampleFilters)
	}
	// resize image to fit dimensions of plate and change each pixel to match closest colour from chosen palette
	// the output of this is a map of well positions to colours needed
	positiontocolourmap, imgBase := image.ImagetoPlatelayout(imgBase, _input.OutPlate, &chosencolourpalette, _input.Rotate, _input.AutoRotate)

	colourtostringmap := image.AvailableComponentmaps()[_input.Palettename]

	// if the image will be printed using fluorescent proteins, 2 previews will be generated for the image (i) under UV light (ii) under visible light

	if _input.UVimage {
		uvmap := image.AvailableComponentmaps()[_input.Palettename]
		visiblemap := image.Visibleequivalentmaps()[_input.Palettename]

		if _input.Subset {
			uvmap = image.MakeSubMapfromMap(colourtostringmap, _input.Subsetnames)
			visiblemap = image.MakeSubMapfromMap(colourtostringmap, _input.Subsetnames)
		}
		image.PrintFPImagePreview(imgBase, _input.OutPlate, _input.Rotate, visiblemap, uvmap)
	}

	// get components from factory
	componentmap := make(map[string]*wtype.LHComponent, 0)

	if _input.Subset {
		colourtostringmap = image.MakeSubMapfromMap(colourtostringmap, _input.Subsetnames)
	}

	for colourname := range colourtostringmap {

		componentname := colourtostringmap[colourname]

		// use template component instead
		var componenttopick *wtype.LHComponent

		if _input.ComponentType != nil {
			componenttopick = _input.ComponentType
		} else {
			componenttopick = execute.NewComponent(_ctx, "Paint")
		}
		componenttopick.CName = componentname

		componentmap[componentname] = componenttopick

	}

	solutions := make([]*wtype.LHComponent, 0)

	counter := 0
	_output.UniqueComponents = make([]string, 0)

	//---------------------------------------------------------------------
	//Pipetting
	//---------------------------------------------------------------------

	// loop through the position to colour map pipeting the correct coloured protein into each well
	for locationkey, colour := range positiontocolourmap {

		component := componentmap[colourtostringmap[colour]]

		// make sure liquid class is appropriate for cell culture in case this is not set elsewhere
		component.Type, _ = wtype.LiquidTypeFromString(_input.UseLiquidClass) //wtype.LTCulture

		// if the option to only print a single colour is not selected then the pipetting actions for all colours (apart from if not this colour is not empty) will follow
		if _input.OnlythisColour != "" /*&& !strings.Contains(locationkey,"x")&& !strings.Contains(locationkey,"X")*/ {

			if image.Colourcomponentmap[colour] == _input.OnlythisColour {

				_output.UniqueComponents = append(_output.UniqueComponents, component.CName)

				counter = counter + 1

				pixelSample := mixer.Sample(component, _input.VolumePerWell)

				solution := execute.MixTo(_ctx, _input.OutPlate.Type, locationkey, 1, pixelSample)

				solutions = append(solutions, solution)
			}

		} else {
			if component.CName != _input.Notthiscolour && component.CName != "transparent" /*&& !strings.Contains(locationkey,"x")&& !strings.Contains(locationkey,"X")*/ {

				_output.UniqueComponents = append(_output.UniqueComponents, component.CName)

				counter = counter + 1

				component.Type, _ = wtype.LiquidTypeFromString(_input.UseLiquidClass)
				pixelSample := mixer.Sample(component, _input.VolumePerWell)

				solution := execute.MixTo(_ctx, _input.OutPlate.Type, locationkey, 1, pixelSample)

				solutions = append(solutions, solution)
			}
		}
	}

	_output.UniqueComponents = search.RemoveDuplicates(_output.UniqueComponents)
	_output.Pixels = solutions

	_output.Numberofpixels = len(_output.Pixels)

	//--------------------------------------------------------------
	//Exporting resulting images
	//--------------------------------------------------------------

	_output.ResizedImage, err = image.Export(imgBase, _input.ImageFileName)
	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PipetteImageAnalysis(_ctx context.Context, _input *PipetteImageInput, _output *PipetteImageOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PipetteImageValidation(_ctx context.Context, _input *PipetteImageInput, _output *PipetteImageOutput) {

}
func _PipetteImageRun(_ctx context.Context, input *PipetteImageInput) *PipetteImageOutput {
	output := &PipetteImageOutput{}
	_PipetteImageSetup(_ctx, input)
	_PipetteImageSteps(_ctx, input, output)
	_PipetteImageAnalysis(_ctx, input, output)
	_PipetteImageValidation(_ctx, input, output)
	return output
}

func PipetteImageRunSteps(_ctx context.Context, input *PipetteImageInput) *PipetteImageSOutput {
	soutput := &PipetteImageSOutput{}
	output := _PipetteImageRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PipetteImageNew() interface{} {
	return &PipetteImageElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PipetteImageInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PipetteImageRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PipetteImageInput{},
			Out: &PipetteImageOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type PipetteImageElement struct {
	inject.CheckedRunner
}

type PipetteImageInput struct {
	AutoRotate            bool
	CheckResizeAlgorithms bool
	ComponentType         *wtype.LHComponent
	ImageFileName         string
	InputFile             wtype.File
	Notthiscolour         string
	OnlythisColour        string
	OutPlate              *wtype.LHPlate
	Palettename           string
	Rotate                bool
	Subset                bool
	Subsetnames           []string
	URL                   string
	UVimage               bool
	UseLiquidClass        wtype.PolicyName
	UseURL                bool
	VolumePerWell         wunit.Volume
}

type PipetteImageOutput struct {
	Numberofpixels   int
	Pixels           []*wtype.LHComponent
	ResizedImage     wtype.File
	ResizedImages    []*goimage.NRGBA
	UniqueComponents []string
}

type PipetteImageSOutput struct {
	Data struct {
		Numberofpixels   int
		ResizedImage     wtype.File
		ResizedImages    []*goimage.NRGBA
		UniqueComponents []string
	}
	Outputs struct {
		Pixels []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PipetteImage",
		Constructor: PipetteImageNew,
		Desc: component.ComponentDesc{
			Desc: "Generates instructions to pipette out a defined image onto a defined plate using a defined palette of coloured bacteria\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/PipetteImage/PipetteImage/PipetteImage/PipetteImage.an",
			Params: []component.ParamDesc{
				{Name: "AutoRotate", Desc: "", Kind: "Parameters"},
				{Name: "CheckResizeAlgorithms", Desc: "", Kind: "Parameters"},
				{Name: "ComponentType", Desc: "", Kind: "Inputs"},
				{Name: "ImageFileName", Desc: "Desired name for the output image file\n", Kind: "Parameters"},
				{Name: "InputFile", Desc: "Image File\n", Kind: "Parameters"},
				{Name: "Notthiscolour", Desc: "", Kind: "Parameters"},
				{Name: "OnlythisColour", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Palettename", Desc: "", Kind: "Parameters"},
				{Name: "Rotate", Desc: "", Kind: "Parameters"},
				{Name: "Subset", Desc: "", Kind: "Parameters"},
				{Name: "Subsetnames", Desc: "", Kind: "Parameters"},
				{Name: "URL", Desc: "", Kind: "Parameters"},
				{Name: "UVimage", Desc: "", Kind: "Parameters"},
				{Name: "UseLiquidClass", Desc: "", Kind: "Parameters"},
				{Name: "UseURL", Desc: "", Kind: "Parameters"},
				{Name: "VolumePerWell", Desc: "", Kind: "Parameters"},
				{Name: "Numberofpixels", Desc: "", Kind: "Data"},
				{Name: "Pixels", Desc: "", Kind: "Outputs"},
				{Name: "ResizedImage", Desc: "", Kind: "Data"},
				{Name: "ResizedImages", Desc: "", Kind: "Data"},
				{Name: "UniqueComponents", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
