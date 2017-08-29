// Protocol WavelengthParser performs something.
package lib

import

// Place golang packages to import here
(
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/export"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/platereader/dataset"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/platereader/dataset/parse"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Parameters to this protocol

//optional

// Output data of this protocol

// Physical inputs to this protocol

// Physical outputs to this protocol

// Conditions to run on startup
func _WavelengthParserSetup(_ctx context.Context, _input *WavelengthParserInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _WavelengthParserSteps(_ctx context.Context, _input *WavelengthParserInput, _output *WavelengthParserOutput) {

	var err error
	var warnings []error

	if _input.Wavelength == 0 {
		_input.Wavelength = 600
	}

	if _input.OutputFileName == "" {
		_input.OutputFileName = "Output.xlsx"
	}

	data, err := _input.MarsResultsFileXLSX.ReadAll()

	if err != nil {
		execute.Errorf(_ctx, "Error in reading MARS file %s: %s", _input.MarsResultsFileXLSX.Name, err.Error())
	}

	var plateReaderData dataset.AbsorbanceData

	if _input.PlateReaderFileType == "" {
		_input.PlateReaderFileType = "Mars"
	}

	if _input.PlateReaderFileType == "Mars" {
		plateReaderData, err = parse.ParseMarsXLSXBinary(data, _input.SheetNumber)
		if err != nil {
			_output.Warnings = append(_output.Warnings, wtype.NewWarning(err.Error()))
			execute.Errorf(_ctx, err.Error())
		}
	} else if _input.PlateReaderFileType == "SpectraMax" {
		plateReaderData, err = parse.ParseSpectraMaxData(data)
		if err != nil {
			_output.Warnings = append(_output.Warnings, wtype.NewWarning(err.Error()))
			execute.Errorf(_ctx, err.Error())
		}
	}

	//if len input plate zero

	//if err != nil {
	//Errorf("Error in reading InputPlateLayout file %s: %s", InputPlateLayout.Name, err.Error())
	//}

	componentsMap := make(map[string]*wtype.LHComponent)

	for _, well := range _input.InputPlateLayout.AllWellPositions(wtype.BYCOLUMN) {
		if !_input.InputPlateLayout.WellMap()[well].Empty() {
			component := _input.InputPlateLayout.WellMap()[well].WContents //is empty for error
			componentsMap[well] = component
		}
	}

	blankValues := make([]float64, 0)
	var blankValue float64

	if len(_input.BlankWells) == 0 && _input.BlankCorrect {
		execute.Errorf(_ctx, "Cannot Blank correct if no blank well is specified")
	}

	if len(_input.BlankWells) > 0 && _input.BlankCorrect {
		for i := range _input.BlankWells {
			blankValue, err = plateReaderData.AbsorbanceReading(_input.BlankWells[i], _input.Wavelength, _input.ReadingTypeinMarsFile)
			if err != nil {
				execute.Errorf(_ctx, fmt.Sprint("Blank sample not found at position ", _input.BlankWells[i], ": ", err.Error()))
			}
		}
		blankValues = append(blankValues, blankValue)
	}

	for well := range componentsMap {

		if len(_input.BlankWells) > 0 && _input.BlankCorrect {
			blankCorrected, err := plateReaderData.BlankCorrect([]string{well}, _input.BlankWells, _input.Wavelength, _input.ReadingTypeinMarsFile)
			if err != nil {
				warnings = append(warnings, err)
				execute.Errorf(_ctx, err.Error())
			}
			_output.BlankedData[well] = blankCorrected
		}
		rawData, err := plateReaderData.AbsorbanceReading(well, _input.Wavelength, _input.ReadingTypeinMarsFile)
		if err != nil {
			warnings = append(warnings, err)
			execute.Errorf(_ctx, err.Error())
		}
		_output.RawData[well] = rawData
	}
	var csv [][]string
	var headers []string = []string{"Well", "SampleName", "BlankedReading"}
	csv = append(csv, headers)
	for well, data := range _output.BlankedData {
		var data []string = []string{well, componentsMap[well].CName, fmt.Sprint(data)}
		csv = append(csv, data)
	}
	_output.OutPutFile, err = export.CSV(csv, _input.OutputFileName)
	if err != nil {
		warnings = append(warnings, err)
		execute.Errorf(_ctx, err.Error())
	}
}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _WavelengthParserAnalysis(_ctx context.Context, _input *WavelengthParserInput, _output *WavelengthParserOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _WavelengthParserValidation(_ctx context.Context, _input *WavelengthParserInput, _output *WavelengthParserOutput) {

}
func _WavelengthParserRun(_ctx context.Context, input *WavelengthParserInput) *WavelengthParserOutput {
	output := &WavelengthParserOutput{}
	_WavelengthParserSetup(_ctx, input)
	_WavelengthParserSteps(_ctx, input, output)
	_WavelengthParserAnalysis(_ctx, input, output)
	_WavelengthParserValidation(_ctx, input, output)
	return output
}

func WavelengthParserRunSteps(_ctx context.Context, input *WavelengthParserInput) *WavelengthParserSOutput {
	soutput := &WavelengthParserSOutput{}
	output := _WavelengthParserRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func WavelengthParserNew() interface{} {
	return &WavelengthParserElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &WavelengthParserInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _WavelengthParserRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &WavelengthParserInput{},
			Out: &WavelengthParserOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type WavelengthParserElement struct {
	inject.CheckedRunner
}

type WavelengthParserInput struct {
	BlankCorrect          bool
	BlankWells            []string
	InputPlateLayout      *wtype.LHPlate
	MarsResultsFileXLSX   wtype.File
	OutputFileName        string
	PlateReaderFileType   string
	ReadingTypeinMarsFile string
	SheetNumber           int
	UseScript             int
	Wavelength            int
}

type WavelengthParserOutput struct {
	BlankedData map[string]float64
	OutPutFile  wtype.File
	RawData     map[string]float64
	Warnings    []wtype.Warning
}

type WavelengthParserSOutput struct {
	Data struct {
		BlankedData map[string]float64
		OutPutFile  wtype.File
		RawData     map[string]float64
		Warnings    []wtype.Warning
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "WavelengthParser",
		Constructor: WavelengthParserNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol WavelengthParser performs something.\n",
			Path: "src/github.com/antha-lang/elements/an/Utility/WavelengthParser/WavelengthParser.an",
			Params: []component.ParamDesc{
				{Name: "BlankCorrect", Desc: "", Kind: "Parameters"},
				{Name: "BlankWells", Desc: "", Kind: "Parameters"},
				{Name: "InputPlateLayout", Desc: "", Kind: "Inputs"},
				{Name: "MarsResultsFileXLSX", Desc: "", Kind: "Parameters"},
				{Name: "OutputFileName", Desc: "", Kind: "Parameters"},
				{Name: "PlateReaderFileType", Desc: "", Kind: "Parameters"},
				{Name: "ReadingTypeinMarsFile", Desc: "", Kind: "Parameters"},
				{Name: "SheetNumber", Desc: "", Kind: "Parameters"},
				{Name: "UseScript", Desc: "optional\n", Kind: "Parameters"},
				{Name: "Wavelength", Desc: "", Kind: "Parameters"},
				{Name: "BlankedData", Desc: "", Kind: "Data"},
				{Name: "OutPutFile", Desc: "", Kind: "Data"},
				{Name: "RawData", Desc: "", Kind: "Data"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
