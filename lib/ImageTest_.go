// Protocol ImageTest performs something.
package lib

import

// Place golang packages to import here
(
	"context"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Parameters to this protocol

// name of image file or if using URL use this field to set the desired filename
// select this if getting the image from a URL
// enter URL link to the image file here if applicable

// Expected outputs

// Output data of this protocol

// Physical inputs to this protocol

// Physical outputs to this protocol

// Conditions to run on startup
func _ImageTestSetup(_ctx context.Context, _input *ImageTestInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _ImageTestSteps(_ctx context.Context, _input *ImageTestInput, _output *ImageTestOutput) {

	// run SerialDilution_ForConcentration element
	result := PipetteImageRunSteps(_ctx, &PipetteImageInput{VolumePerWell: _input.VolumePerWell,
		Imagefilename:         _input.Imagefilename, // name of image file or if using URL use this field to set the desired filename
		UseURL:                _input.UseURL,        // select this if getting the image from a URL
		URL:                   _input.URL,           // enter URL link to the image file here if applicable
		Palettename:           _input.Palettename,
		OnlythisColour:        _input.OnlythisColour,
		Notthiscolour:         _input.Notthiscolour,
		UVimage:               _input.UVimage,
		Rotate:                _input.Rotate,
		AutoRotate:            _input.AutoRotate,
		Subset:                _input.Subset,
		UseLiquidClass:        _input.UseLiquidClass,
		Subsetnames:           _input.Subsetnames,
		CheckResizeAlgorithms: _input.CheckResizeAlgorithms,

		ComponentType: _input.ComponentType,
		OutPlate:      _input.OutPlate},
	)

	if result.Data.Numberofpixels != _input.ExpectedNumberOfPixels {
		execute.Errorf(_ctx, "Pipette Image test Fail for image %s: Expected pixels %d: got %d", _input.Imagefilename, _input.ExpectedNumberOfPixels, result.Data.Numberofpixels)
	}

}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _ImageTestAnalysis(_ctx context.Context, _input *ImageTestInput, _output *ImageTestOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _ImageTestValidation(_ctx context.Context, _input *ImageTestInput, _output *ImageTestOutput) {

}
func _ImageTestRun(_ctx context.Context, input *ImageTestInput) *ImageTestOutput {
	output := &ImageTestOutput{}
	_ImageTestSetup(_ctx, input)
	_ImageTestSteps(_ctx, input, output)
	_ImageTestAnalysis(_ctx, input, output)
	_ImageTestValidation(_ctx, input, output)
	return output
}

func ImageTestRunSteps(_ctx context.Context, input *ImageTestInput) *ImageTestSOutput {
	soutput := &ImageTestSOutput{}
	output := _ImageTestRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ImageTestNew() interface{} {
	return &ImageTestElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ImageTestInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ImageTestRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ImageTestInput{},
			Out: &ImageTestOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type ImageTestElement struct {
	inject.CheckedRunner
}

type ImageTestInput struct {
	AutoRotate             bool
	CheckResizeAlgorithms  bool
	ComponentType          *wtype.LHComponent
	ExpectedNumberOfPixels int
	Imagefilename          string
	Notthiscolour          string
	OnlythisColour         string
	OutPlate               *wtype.LHPlate
	Palettename            string
	Rotate                 bool
	Subset                 bool
	Subsetnames            []string
	URL                    string
	UVimage                bool
	UseLiquidClass         string
	UseURL                 bool
	VolumePerWell          wunit.Volume
}

type ImageTestOutput struct {
}

type ImageTestSOutput struct {
	Data struct {
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "ImageTest",
		Constructor: ImageTestNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol ImageTest performs something.\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/PipetteImage/PipetteImage/ImageTest/ImageTest.an",
			Params: []component.ParamDesc{
				{Name: "AutoRotate", Desc: "", Kind: "Parameters"},
				{Name: "CheckResizeAlgorithms", Desc: "", Kind: "Parameters"},
				{Name: "ComponentType", Desc: "", Kind: "Inputs"},
				{Name: "ExpectedNumberOfPixels", Desc: "Expected outputs\n", Kind: "Parameters"},
				{Name: "Imagefilename", Desc: "name of image file or if using URL use this field to set the desired filename\n", Kind: "Parameters"},
				{Name: "Notthiscolour", Desc: "", Kind: "Parameters"},
				{Name: "OnlythisColour", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Palettename", Desc: "", Kind: "Parameters"},
				{Name: "Rotate", Desc: "", Kind: "Parameters"},
				{Name: "Subset", Desc: "", Kind: "Parameters"},
				{Name: "Subsetnames", Desc: "", Kind: "Parameters"},
				{Name: "URL", Desc: "enter URL link to the image file here if applicable\n", Kind: "Parameters"},
				{Name: "UVimage", Desc: "", Kind: "Parameters"},
				{Name: "UseLiquidClass", Desc: "", Kind: "Parameters"},
				{Name: "UseURL", Desc: "select this if getting the image from a URL\n", Kind: "Parameters"},
				{Name: "VolumePerWell", Desc: "", Kind: "Parameters"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
