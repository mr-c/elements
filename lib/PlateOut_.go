// Protocol PlateOut dispenses a liquid input (i.e Transformed Cells) at a user-definable volume onto an output plate of the users choice.
package lib

import

// Place golang packages to import here
(
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Parameters to this protocol

//optionally specify the number of agar plates to begin counting from (Default = 1)
//set Incubation temperature
//set Incubation time
//specify number of technical replicates to plate out
//optionally specify the liquid handling policy to use when plating out (Default = PlateOut). Can change
//optionally specify the plate out volume. If Dilution is required, this volume will be made up to with the transformed cells and the diluent
//specify if some wells have already been used in the Agar Plate (i.e. if a plate is being used for multiple tranformations, or an overlay)

// Output data of this protocol

//returns number of output AgarPlates used

// Physical inputs to this protocol

//the output plate type, which can be any plate within the Antha library (Default = falcon6wellAgar)
//the transformed cells (Default = neb5compcells).

// Physical outputs to this protocol

//the plated cultures are outputted as an array which can be fed into other protocols in the Antha workflow
//the number of plates used
//number of wells used

// Conditions to run on startup
func _PlateOutSetup(_ctx context.Context, _input *PlateOutInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _PlateOutSteps(_ctx context.Context, _input *PlateOutInput, _output *PlateOutOutput) {
	//set up some default values
	var defaultPlateOutPolicy wtype.PolicyName = "plateout"
	var defaultAgarPlateNumber int = 1
	var defaultNumberofReplicates int = 1
	var defaultWellsAlreadyUsed int = 0

	//set up error variable
	var warnings []string

	//set counter variable to count up number of wells used, and set at the number of wells already used in the output agar plate (Default = 0)
	var counter int = _input.WellsAlreadyUsed

	if _input.WellsAlreadyUsed < 0 {
		_input.WellsAlreadyUsed = defaultWellsAlreadyUsed
		wellsUsedError := fmt.Errorf("Invalid WellsAlreadyUsed specified, assinging to default: %d", defaultWellsAlreadyUsed)
		warnings = append(warnings, wellsUsedError.Error())
	}

	if _input.AgarPlateNumber <= 0 {
		_input.AgarPlateNumber = defaultAgarPlateNumber
		agarPlateNumError := fmt.Errorf("Invalid AgarPlateNumber specified, assinging to default: %d", defaultAgarPlateNumber)
		warnings = append(warnings, agarPlateNumError.Error())
	}

	//set platenumber variable (Default = 1) that will count up number of plates used
	var platenumber int = _input.AgarPlateNumber

	if _input.PlateOutLiquidPolicy == "" {
		_input.PlateOutLiquidPolicy = defaultPlateOutPolicy
		plateOutPolicyError := fmt.Errorf("Invalid PlateOutLiquidPolicy specified, assinging to default: %s", defaultPlateOutPolicy)
		warnings = append(warnings, plateOutPolicyError.Error())
	}

	//attribute specified liquidpolicy to the plate out reaction (Default = plateout)
	_input.TransformedCells.Type, _ = wtype.LiquidTypeFromString(_input.PlateOutLiquidPolicy)

	//get plate dimenson and well info for specified agarplate from plate library
	var wellpositionarray []string = _input.AgarPlate.AllWellPositions(wtype.BYCOLUMN)

	if _input.NumberofReplicates <= 0 {
		_input.NumberofReplicates = defaultNumberofReplicates
		numberOfReplicatesError := fmt.Errorf("Invalid NumberofReplicates specified, assigning to default: %s", defaultNumberofReplicates)
		warnings = append(warnings, numberOfReplicatesError.Error())

	}
	//create loop for processing through specified number of replicates
	for j := 0; j < _input.NumberofReplicates; j++ {

		//set up a slice to add the plate out reactions to
		var plateOutSamplesSlice []*wtype.LHComponent

		//set up variable for tracking well
		var nextwell string

		//detect next well location accessing array slice using counter as pointer
		nextwell = wellpositionarray[counter]

		//aspirate transformed cells at specified volumes
		plateOutSample := mixer.Sample(_input.TransformedCells, _input.PlateOutVolume)

		//append transformed cell volumes to plate out volumes array
		plateOutSamplesSlice = append(plateOutSamplesSlice, plateOutSample)

		//perform mix actions with the plate out volume reactions from above into specified plate and location
		platedCulture := execute.MixNamed(_ctx, _input.AgarPlate.Type, nextwell, fmt.Sprint("TransformedPlateNumber", platenumber), plateOutSamplesSlice...)

		//append plated out cultures to output array
		_output.PlatedCultures = append(_output.PlatedCultures, platedCulture)

		//increase counter for next iteration and add additonal plate if needed
		if counter+1 == len(wellpositionarray) {
			platenumber++
			counter = 0
		} else {
			counter++
		}

		//add WellLocationsUsed to output data slice
		_output.WellLocationsUsed = append(_output.WellLocationsUsed, nextwell)
	}

	//update counters and append warnings
	_output.WellsUsed = counter
	_output.TransformedPlateNumber = platenumber
	_output.Errors = warnings
}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _PlateOutAnalysis(_ctx context.Context, _input *PlateOutInput, _output *PlateOutOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _PlateOutValidation(_ctx context.Context, _input *PlateOutInput, _output *PlateOutOutput) {

}
func _PlateOutRun(_ctx context.Context, input *PlateOutInput) *PlateOutOutput {
	output := &PlateOutOutput{}
	_PlateOutSetup(_ctx, input)
	_PlateOutSteps(_ctx, input, output)
	_PlateOutAnalysis(_ctx, input, output)
	_PlateOutValidation(_ctx, input, output)
	return output
}

func PlateOutRunSteps(_ctx context.Context, input *PlateOutInput) *PlateOutSOutput {
	soutput := &PlateOutSOutput{}
	output := _PlateOutRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PlateOutNew() interface{} {
	return &PlateOutElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PlateOutInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PlateOutRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PlateOutInput{},
			Out: &PlateOutOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type PlateOutElement struct {
	inject.CheckedRunner
}

type PlateOutInput struct {
	AgarPlate            *wtype.LHPlate
	AgarPlateNumber      int
	IncubationTemp       wunit.Temperature
	IncubationTime       wunit.Time
	NumberofReplicates   int
	PlateOutLiquidPolicy wtype.PolicyName
	PlateOutVolume       wunit.Volume
	TransformedCells     *wtype.LHComponent
	WellsAlreadyUsed     int
}

type PlateOutOutput struct {
	AgarPlatesUsed         int
	Errors                 []string
	PlatedCultures         []*wtype.LHComponent
	TransformedPlateNumber int
	WellLocationsUsed      []string
	WellsUsed              int
}

type PlateOutSOutput struct {
	Data struct {
		AgarPlatesUsed    int
		Errors            []string
		WellLocationsUsed []string
	}
	Outputs struct {
		PlatedCultures         []*wtype.LHComponent
		TransformedPlateNumber int
		WellsUsed              int
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PlateOut",
		Constructor: PlateOutNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol PlateOut dispenses a liquid input (i.e Transformed Cells) at a user-definable volume onto an output plate of the users choice.\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/PlateOut/PlateOut.an",
			Params: []component.ParamDesc{
				{Name: "AgarPlate", Desc: "the output plate type, which can be any plate within the Antha library (Default = falcon6wellAgar)\n", Kind: "Inputs"},
				{Name: "AgarPlateNumber", Desc: "optionally specify the number of agar plates to begin counting from (Default = 1)\n", Kind: "Parameters"},
				{Name: "IncubationTemp", Desc: "set Incubation temperature\n", Kind: "Parameters"},
				{Name: "IncubationTime", Desc: "set Incubation time\n", Kind: "Parameters"},
				{Name: "NumberofReplicates", Desc: "specify number of technical replicates to plate out\n", Kind: "Parameters"},
				{Name: "PlateOutLiquidPolicy", Desc: "optionally specify the liquid handling policy to use when plating out (Default = PlateOut). Can change\n", Kind: "Parameters"},
				{Name: "PlateOutVolume", Desc: "optionally specify the plate out volume. If Dilution is required, this volume will be made up to with the transformed cells and the diluent\n", Kind: "Parameters"},
				{Name: "TransformedCells", Desc: "the transformed cells (Default = neb5compcells).\n", Kind: "Inputs"},
				{Name: "WellsAlreadyUsed", Desc: "specify if some wells have already been used in the Agar Plate (i.e. if a plate is being used for multiple tranformations, or an overlay)\n", Kind: "Parameters"},
				{Name: "AgarPlatesUsed", Desc: "returns number of output AgarPlates used\n", Kind: "Data"},
				{Name: "Errors", Desc: "", Kind: "Data"},
				{Name: "PlatedCultures", Desc: "the plated cultures are outputted as an array which can be fed into other protocols in the Antha workflow\n", Kind: "Outputs"},
				{Name: "TransformedPlateNumber", Desc: "the number of plates used\n", Kind: "Outputs"},
				{Name: "WellLocationsUsed", Desc: "", Kind: "Data"},
				{Name: "WellsUsed", Desc: "number of wells used\n", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
