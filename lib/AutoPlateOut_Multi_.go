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

//optionally specify the number of agar plates to begin counting from (Default = 1)
//set Incubation temperature
//set Incubation time
//specify number of technical replicates to plate out
//optionally specify the liquid handling policy to use when plating out (Default = PlateOut). Can change
//specify the plate out volume. If Dilution is required, this volume will be made up to with the transformed cells and the diluent

//optionally specify if some wells have already been used in the Agar Plate (i.e. if a plate is being used for multiple tranformations, or an overlay)

// Output data of this protocol

// Physical inputs to this protocol

//the output plate type, which can be any plate within the Antha library (Default = falcon6wellAgar)
//the transformed cells that can be inputted from another protocol (e.g.  AutTransformation_multi)

// Physical outputs to this protocol

//the plated cultures are outputted as an array which can be fed into other protocols in the Antha workflow

// Conditions to run on startup
func _AutoPlateOut_MultiSetup(_ctx context.Context, _input *AutoPlateOut_MultiInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _AutoPlateOut_MultiSteps(_ctx context.Context, _input *AutoPlateOut_MultiInput, _output *AutoPlateOut_MultiOutput) {
	//setup counter to track WellsAlreadyUsed
	var counter int = _input.WellsAlreadyUsed
	var platecounter int = _input.AgarPlateNumber

	//range through the inputted array and perform the PlateOutTest protocol
	for _, plateout := range _input.TransformedCells {

		var reactionname string = plateout.CName

		var volumes []wunit.Volume

		if value, found := _input.PlateOutVolumes[reactionname]; found {
			volumes = value
		} else if value, found := _input.PlateOutVolumes["default"]; found {
			volumes = value
		} else {
			execute.Errorf(_ctx, "No plate out volumes set for %s. Please set these", reactionname)
		}

		for _, plateoutvols := range volumes {

			// Run PlateOut element
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
			for _, plateoutorder := range result.Outputs.PlatedCultures {
				_output.PlatedCultures = append(_output.PlatedCultures, plateoutorder)
			}

			//increase counter
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
	SpecifyReactionOrder []string
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
				{Name: "AgarPlate", Desc: "the output plate type, which can be any plate within the Antha library (Default = falcon6wellAgar)\n", Kind: "Inputs"},
				{Name: "AgarPlateNumber", Desc: "optionally specify the number of agar plates to begin counting from (Default = 1)\n", Kind: "Parameters"},
				{Name: "IncubationTemp", Desc: "set Incubation temperature\n", Kind: "Parameters"},
				{Name: "IncubationTime", Desc: "set Incubation time\n", Kind: "Parameters"},
				{Name: "NumberofReplicates", Desc: "specify number of technical replicates to plate out\n", Kind: "Parameters"},
				{Name: "PlateOutLiquidPolicy", Desc: "optionally specify the liquid handling policy to use when plating out (Default = PlateOut). Can change\n", Kind: "Parameters"},
				{Name: "PlateOutVolumes", Desc: "specify the plate out volume. If Dilution is required, this volume will be made up to with the transformed cells and the diluent\n", Kind: "Parameters"},
				{Name: "SpecifyReactionOrder", Desc: "", Kind: "Parameters"},
				{Name: "TransformedCells", Desc: "the transformed cells that can be inputted from another protocol (e.g.  AutTransformation_multi)\n", Kind: "Inputs"},
				{Name: "WellsAlreadyUsed", Desc: "optionally specify if some wells have already been used in the Agar Plate (i.e. if a plate is being used for multiple tranformations, or an overlay)\n", Kind: "Parameters"},
				{Name: "PlatedCultures", Desc: "the plated cultures are outputted as an array which can be fed into other protocols in the Antha workflow\n", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
