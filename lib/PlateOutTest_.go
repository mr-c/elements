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
	"github.com/antha-lang/antha/microArch/factory"
)

// Parameters to this protocol

//optionally specify the number of agar plates to begin counting from (Default = 1)
//specify if dilution of the transformed cells is required, and the level of dilution (Default = 1 in which the sample will not be diluted). Dilution will be performed with the Diluent (Default = LB)
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
//the liquid with which to dilute the transformed cells (Default = LB)
//the transformed cells (Default = neb5compcells).

// Physical outputs to this protocol

//the plated cultures are outputted as an array which can be fed into other protocols in the Antha workflow
//the number of plates used

// Conditions to run on startup
func _PlateOutTestSetup(_ctx context.Context, _input *PlateOutTestInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _PlateOutTestSteps(_ctx context.Context, _input *PlateOutTestInput, _output *PlateOutTestOutput) {

	//set platenumber variable (Default = 1) that will count up number of plates used
	var platenumber int = _input.AgarPlateNumber

	if _input.AgarPlateNumber == 0 {
		_input.AgarPlateNumber = 1
	}

	//set counter variable to count up number of wells used, and set at the number of wells already used in the output agar plate (Default = 0)
	var counter int = _input.WellsAlreadyUsed

	//set up error variable
	var err error

	//set up a slice to add the plate out reactions to
	var plateOutVolumes []*wtype.LHComponent

	//attribute specified liquidpolicy to the plate out reaction (Default = plateout)
	_input.TransformedCells.Type, err = wtype.LiquidTypeFromString(_input.PlateOutLiquidPolicy)

	if _input.PlateOutLiquidPolicy == "" {
		_input.PlateOutLiquidPolicy = "plateout"
	}

	//get plate dimenson and well info for specified agarplate from plate library
	var wellpositionarray []string = _input.AgarPlate.AllWellPositions(wtype.BYCOLUMN)

	if err != nil {
		execute.Errorf(_ctx, "Error in specifying Liquid Policy %s for Plate Out: %s", _input.PlateOutLiquidPolicy, err.Error())
	}

	//create loop for processing through specified number of replicates
	for j := 0; j < _input.NumberofReplicates; j++ {

		//detect next well location accessing array slice using counter as pointer
		nextwell := wellpositionarray[counter]

		//check if dilution is required and calculate required dilution, performing mix command with speicifed Diluent (Default = LB)
		var nilComponent *wtype.LHComponent

		if _input.Diluent == nilComponent {
			_input.Diluent = factory.GetComponentByType("LB")
		}

		if _input.Dilution > 1 {
			dilutedSample := mixer.SampleForTotalVolume(_input.Diluent, _input.PlateOutVolume)
			plateOutVolumes = append(plateOutVolumes, dilutedSample)
			_input.PlateOutVolume = wunit.DivideVolume(_input.PlateOutVolume, float64(_input.Dilution))
		}

		//aspirate transformed cells at specified volumes
		plateOutSample := mixer.Sample(_input.TransformedCells, _input.PlateOutVolume)

		//append transformed cell volumes to plate out volumes array
		plateOutVolumes = append(plateOutVolumes, plateOutSample)

		//perform mix actions with the plate out volume reactions from above into specified plate and location
		platedCulture := execute.MixNamed(_ctx, _input.AgarPlate.Type, nextwell, fmt.Sprint("TransformedPlate", platenumber), plateOutVolumes...)

		//append plated out cultures to output array
		_output.PlatedCultures = append(_output.PlatedCultures, platedCulture)

	}

	//increase counter for next iteration and add additonal plate if needed
	if counter+1 == len(wellpositionarray) {
		platenumber++
		counter = 0
	} else {
		counter++
	}
}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _PlateOutTestAnalysis(_ctx context.Context, _input *PlateOutTestInput, _output *PlateOutTestOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _PlateOutTestValidation(_ctx context.Context, _input *PlateOutTestInput, _output *PlateOutTestOutput) {

}
func _PlateOutTestRun(_ctx context.Context, input *PlateOutTestInput) *PlateOutTestOutput {
	output := &PlateOutTestOutput{}
	_PlateOutTestSetup(_ctx, input)
	_PlateOutTestSteps(_ctx, input, output)
	_PlateOutTestAnalysis(_ctx, input, output)
	_PlateOutTestValidation(_ctx, input, output)
	return output
}

func PlateOutTestRunSteps(_ctx context.Context, input *PlateOutTestInput) *PlateOutTestSOutput {
	soutput := &PlateOutTestSOutput{}
	output := _PlateOutTestRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PlateOutTestNew() interface{} {
	return &PlateOutTestElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PlateOutTestInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PlateOutTestRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PlateOutTestInput{},
			Out: &PlateOutTestOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type PlateOutTestElement struct {
	inject.CheckedRunner
}

type PlateOutTestInput struct {
	AgarPlate            *wtype.LHPlate
	AgarPlateNumber      int
	Diluent              *wtype.LHComponent
	Dilution             int
	IncubationTemp       wunit.Temperature
	IncubationTime       wunit.Time
	NumberofReplicates   int
	PlateOutLiquidPolicy string
	PlateOutVolume       wunit.Volume
	TransformedCells     *wtype.LHComponent
	WellsAlreadyUsed     int
}

type PlateOutTestOutput struct {
	AgarPlatesUsed   int
	PlatedCultures   []*wtype.LHComponent
	TransformedPlate int
}

type PlateOutTestSOutput struct {
	Data struct {
		AgarPlatesUsed int
	}
	Outputs struct {
		PlatedCultures   []*wtype.LHComponent
		TransformedPlate int
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PlateOutTest",
		Constructor: PlateOutTestNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol PlateOut dispenses a liquid input (i.e Transformed Cells) at a user-definable volume onto an output plate of the users choice.\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/PlateOut/PlateOut.an",
			Params: []component.ParamDesc{
				{Name: "AgarPlate", Desc: "the output plate type, which can be any plate within the Antha library (Default = falcon6wellAgar)\n", Kind: "Inputs"},
				{Name: "AgarPlateNumber", Desc: "optionally specify the number of agar plates to begin counting from (Default = 1)\n", Kind: "Parameters"},
				{Name: "Diluent", Desc: "the liquid with which to dilute the transformed cells (Default = LB)\n", Kind: "Inputs"},
				{Name: "Dilution", Desc: "specify if dilution of the transformed cells is required, and the level of dilution (Default = 1 in which the sample will not be diluted). Dilution will be performed with the Diluent (Default = LB)\n", Kind: "Parameters"},
				{Name: "IncubationTemp", Desc: "set Incubation temperature\n", Kind: "Parameters"},
				{Name: "IncubationTime", Desc: "set Incubation time\n", Kind: "Parameters"},
				{Name: "NumberofReplicates", Desc: "specify number of technical replicates to plate out\n", Kind: "Parameters"},
				{Name: "PlateOutLiquidPolicy", Desc: "optionally specify the liquid handling policy to use when plating out (Default = PlateOut). Can change\n", Kind: "Parameters"},
				{Name: "PlateOutVolume", Desc: "optionally specify the plate out volume. If Dilution is required, this volume will be made up to with the transformed cells and the diluent\n", Kind: "Parameters"},
				{Name: "TransformedCells", Desc: "the transformed cells (Default = neb5compcells).\n", Kind: "Inputs"},
				{Name: "WellsAlreadyUsed", Desc: "specify if some wells have already been used in the Agar Plate (i.e. if a plate is being used for multiple tranformations, or an overlay)\n", Kind: "Parameters"},
				{Name: "AgarPlatesUsed", Desc: "returns number of output AgarPlates used\n", Kind: "Data"},
				{Name: "PlatedCultures", Desc: "the plated cultures are outputted as an array which can be fed into other protocols in the Antha workflow\n", Kind: "Outputs"},
				{Name: "TransformedPlate", Desc: "the number of plates used\n", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
