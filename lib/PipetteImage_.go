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
	"github.com/antha-lang/antha/microArch/factory"
	goimage "image"
	"image/color"
	"sort"
	"strings"
)

// Input parameters for this protocol (data)

//
//Desired name for the output image file
// Image File
//The URL from which to download the image
//Boolean to signal that you want to use an image from a URL, not a file
//Name of the palette set you want to use (look at the map in the source code)
//use this key to select only one color
//use this key to remove this color
//use this key to use fluorescent proteins
//set to true to rotate the image to fit the plate.
//set to true to rotate the image to fit the plate.
//set to true to use a few colors from the selected palette
//key to set which liquidPolicy to use
//name of the subset of colors to use
//iterate through potential resizeAlgorithm to resize the image

// Data which is returned from this protocol, and data types

//the image resized to fit the plate
//the images for each resize algorithms
//the number of pixels in the final image
//the keys of the IDs of the LHcomponents for every color

// Physical Inputs to this protocol with types

//Component type for the paint LHComponents. Set to "paint" if none given.
//the type of plate to which the image is printed

// Physical outputs from this protocol with types

//The LHComponent for each pixel.

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

	solutions := make([]*wtype.LHComponent, 0)

	counter := 0
	_output.UniqueComponents = make([]string, 0)

	// get components from factory
	componentmap := make(map[string]*wtype.LHComponent, 0)

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
			componenttopick = factory.GetComponentByType("Paint")
		}

		componenttopick.CName = componentname
		componentmap[componentname] = componenttopick

	}

	//---------------------------------------------------------------------
	//Pipetting
	//---------------------------------------------------------------------

	// loop through the position to colour map pipetting the correct coloured protein into each well
	for locationkey, colour := range positiontocolourmap {

		component := componentmap[colourtostringmap[colour]]
		// make sure liquid class is appropriate for cell culture in case this is not set elsewhere
		component.Type, _ = wtype.LiquidTypeFromString(_input.UseLiquidClass) //wtype.LTCulture

		// if the option to only print a single colour is not selected then the pipetting actions for all colours (apart from if not this colour is not empty) will follow
		if _input.OnlythisColour != "" {

			if image.Colourcomponentmap[colour] == _input.OnlythisColour {

				_output.UniqueComponents = append(_output.UniqueComponents, component.CName)

				pixelSample := mixer.Sample(component, _input.VolumePerWell)
				solution := execute.MixTo(_ctx, _input.OutPlate.Type, locationkey, 1, pixelSample)

				solutions = append(solutions, solution)

				counter = counter + 1
			}

		} else {

			if component.CName != _input.Notthiscolour && component.CName != "transparent" {

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
				{Name: "AutoRotate", Desc: "set to true to rotate the image to fit the plate.\n", Kind: "Parameters"},
				{Name: "CheckResizeAlgorithms", Desc: "iterate through potential resizeAlgorithm to resize the image\n", Kind: "Parameters"},
				{Name: "ComponentType", Desc: "Component type for the paint LHComponents. Set to \"paint\" if none given.\n", Kind: "Inputs"},
				{Name: "ImageFileName", Desc: "Desired name for the output image file\n", Kind: "Parameters"},
				{Name: "InputFile", Desc: "Image File\n", Kind: "Parameters"},
				{Name: "Notthiscolour", Desc: "use this key to remove this color\n", Kind: "Parameters"},
				{Name: "OnlythisColour", Desc: "use this key to select only one color\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "the type of plate to which the image is printed\n", Kind: "Inputs"},
				{Name: "Palettename", Desc: "Name of the palette set you want to use (look at the map in the source code)\n", Kind: "Parameters"},
				{Name: "Rotate", Desc: "set to true to rotate the image to fit the plate.\n", Kind: "Parameters"},
				{Name: "Subset", Desc: "set to true to use a few colors from the selected palette\n", Kind: "Parameters"},
				{Name: "Subsetnames", Desc: "name of the subset of colors to use\n", Kind: "Parameters"},
				{Name: "URL", Desc: "The URL from which to download the image\n", Kind: "Parameters"},
				{Name: "UVimage", Desc: "use this key to use fluorescent proteins\n", Kind: "Parameters"},
				{Name: "UseLiquidClass", Desc: "key to set which liquidPolicy to use\n", Kind: "Parameters"},
				{Name: "UseURL", Desc: "Boolean to signal that you want to use an image from a URL, not a file\n", Kind: "Parameters"},
				{Name: "VolumePerWell", Desc: "", Kind: "Parameters"},
				{Name: "Numberofpixels", Desc: "the number of pixels in the final image\n", Kind: "Data"},
				{Name: "Pixels", Desc: "The LHComponent for each pixel.\n", Kind: "Outputs"},
				{Name: "ResizedImage", Desc: "the image resized to fit the plate\n", Kind: "Data"},
				{Name: "ResizedImages", Desc: "the images for each resize algorithms\n", Kind: "Data"},
				{Name: "UniqueComponents", Desc: "the keys of the IDs of the LHcomponents for every color\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
