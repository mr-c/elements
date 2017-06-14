// Protocol PlateOutReactionInput takes in an array of TransformedCells (i.e. recovered cells) from another element (e.g. AutTransformation_multi) and performs a plate out reaction onto plates of the users choice
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

//Optionally specify the number of agar plates to begin counting from (Default = 1)

//Set Incubation temperature if using an associated Incubator

//Set Incubation time if using an associated Incubator

//Specify number of technical replicates to plate out

//Optionally specify the liquid handling policy to use when plating out (Default = PlateOut)

//Specify the plate out volumes. This is a map of each reaction name to a series of volumes of which the reaction will be plated out.
//If a "default" volume is specified, this will be applied to all reactions which do not have specified volumes.
//If Dilution is required, this volume will be made up to with the transformed cells and the diluent.

//Optionally specify if some wells have already been used in the Agar Plate
//(i.e. if a plate has been previously used for tranformations, or an overlay)

// Output data of this protocol

// Physical inputs to this protocol

//The output plate type, which can be any plate within the Antha library (Default = falcon6wellAgar)
//An omniwell may be used for plating out up to 96 spots, but a 96 well plate image must be selected in Antha (e.g. pcrplate_skirted)

//The transformed cells that can be inputed from another protocol (e.g. AutTransformation_multi)

// Physical outputs to this protocol

//the plated cultures are outputted as an array which can be fed into other protocols in the Antha workflow

// Conditions to run on startup
func _AutoPlateOut_MultiSetup(_ctx context.Context, _input *AutoPlateOut_MultiInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _AutoPlateOut_MultiSteps(_ctx context.Context, _input *AutoPlateOut_MultiInput, _output *AutoPlateOut_MultiOutput) {
	//Setup counter to track WellsAlreadyUsed
	var counter int = _input.WellsAlreadyUsed
	var platecounter int = _input.AgarPlateNumber

	//Range through the inputted array and perform the PlateOutTest protocol
	for _, plateout := range _input.TransformedCells {

		var reactionname string = plateout.CName

		var volumes []wunit.Volume

		//Check if volumes speicifed for each reaction and assign default if necessary
		if value, found := _input.PlateOutVolumes[reactionname]; found {
			volumes = value
		} else if value, found := _input.PlateOutVolumes["default"]; found {
			volumes = value
		} else {
			execute.Errorf(_ctx, "No plate out volumes set for %s. Please set these", reactionname)
		}

		//Range through the plate out volumes
		for _, plateoutvols := range volumes {

			//Run PlateOut element
			result := PlateOutTestRunSteps(_ctx, &PlateOutTestInput{AgarPlateNumber: platecounter,
				IncubationTemp:       _input.IncubationTemp,
				IncubationTime:       _input.IncubationTime,
				NumberofReplicates:   _input.NumberofReplicates,
				PlateOutVolume:       plateoutvols,
				PlateOutLiquidPolicy: _input.PlateOutLiquidPolicy,
				WellsAlreadyUsed:     counter,

				TransformedCells: plateout,
				AgarPlate:        _input.AgarPlate},
			)

			//Append outputted plated cultures
			for _, plateoutorder := range result.Outputs.PlatedCultures {
				_output.PlatedCultures = append(_output.PlatedCultures, plateoutorder)
			}

			//Increase counters
			counter = result.Outputs.WellsUsed
			platecounter = result.Outputs.TransformedPlateNumber
		}
	}
}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _AutoPlateOut_MultiAnalysis(_ctx context.Context, _input *AutoPlateOut_MultiInput, _output *AutoPlateOut_MultiOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _AutoPlateOut_MultiValidation(_ctx context.Context, _input *AutoPlateOut_MultiInput, _output *AutoPlateOut_MultiOutput) {

}
func _AutoPlateOut_MultiRun(_ctx context.Context, input *AutoPlateOut_MultiInput) *AutoPlateOut_MultiOutput {
	output := &AutoPlateOut_MultiOutput{}
	_AutoPlateOut_MultiSetup(_ctx, input)
	_AutoPlateOut_MultiSteps(_ctx, input, output)
	_AutoPlateOut_MultiAnalysis(_ctx, input, output)
	_AutoPlateOut_MultiValidation(_ctx, input, output)
	return output
}

func AutoPlateOut_MultiRunSteps(_ctx context.Context, input *AutoPlateOut_MultiInput) *AutoPlateOut_MultiSOutput {
	soutput := &AutoPlateOut_MultiSOutput{}
	output := _AutoPlateOut_MultiRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AutoPlateOut_MultiNew() interface{} {
	return &AutoPlateOut_MultiElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AutoPlateOut_MultiInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AutoPlateOut_MultiRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AutoPlateOut_MultiInput{},
			Out: &AutoPlateOut_MultiOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type AutoPlateOut_MultiElement struct {
	inject.CheckedRunner
}

type AutoPlateOut_MultiInput struct {
	AgarPlate            *wtype.LHPlate
	AgarPlateNumber      int
	IncubationTemp       wunit.Temperature
	IncubationTime       wunit.Time
	NumberofReplicates   int
	PlateOutLiquidPolicy wtype.PolicyName
	PlateOutVolumes      map[string][]wunit.Volume
	TransformedCells     []*wtype.LHComponent
	WellsAlreadyUsed     int
}

type AutoPlateOut_MultiOutput struct {
	PlatedCultures []*wtype.LHComponent
}

type AutoPlateOut_MultiSOutput struct {
	Data struct {
	}
	Outputs struct {
		PlatedCultures []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AutoPlateOut_Multi",
		Constructor: AutoPlateOut_MultiNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol PlateOutReactionInput takes in an array of TransformedCells (i.e. recovered cells) from another element (e.g. AutTransformation_multi) and performs a plate out reaction onto plates of the users choice\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/PlateOut/AutoPlateOut_Multi.an",
			Params: []component.ParamDesc{
				{Name: "AgarPlate", Desc: "The output plate type, which can be any plate within the Antha library (Default = falcon6wellAgar)\nAn omniwell may be used for plating out up to 96 spots, but a 96 well plate image must be selected in Antha (e.g. pcrplate_skirted)\n", Kind: "Inputs"},
				{Name: "AgarPlateNumber", Desc: "Optionally specify the number of agar plates to begin counting from (Default = 1)\n", Kind: "Parameters"},
				{Name: "IncubationTemp", Desc: "Set Incubation temperature if using an associated Incubator\n", Kind: "Parameters"},
				{Name: "IncubationTime", Desc: "Set Incubation time if using an associated Incubator\n", Kind: "Parameters"},
				{Name: "NumberofReplicates", Desc: "Specify number of technical replicates to plate out\n", Kind: "Parameters"},
				{Name: "PlateOutLiquidPolicy", Desc: "Optionally specify the liquid handling policy to use when plating out (Default = PlateOut)\n", Kind: "Parameters"},
				{Name: "PlateOutVolumes", Desc: "Specify the plate out volumes. This is a map of each reaction name to a series of volumes of which the reaction will be plated out.\nIf a \"default\" volume is specified, this will be applied to all reactions which do not have specified volumes.\nIf Dilution is required, this volume will be made up to with the transformed cells and the diluent.\n", Kind: "Parameters"},
				{Name: "TransformedCells", Desc: "The transformed cells that can be inputed from another protocol (e.g. AutTransformation_multi)\n", Kind: "Inputs"},
				{Name: "WellsAlreadyUsed", Desc: "Optionally specify if some wells have already been used in the Agar Plate\n(i.e. if a plate has been previously used for tranformations, or an overlay)\n", Kind: "Parameters"},
				{Name: "PlatedCultures", Desc: "the plated cultures are outputted as an array which can be fed into other protocols in the Antha workflow\n", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
