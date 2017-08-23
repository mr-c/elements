package lib

import (
	"context"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol (data)

//= 50.(uL)

//= 2 (hours)

//Shakerspeed Rate

//Plateoutdilution float64

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _Transformation_complete_noIncubationRequirements() {
}

// Conditions to run on startup
func _Transformation_complete_noIncubationSetup(_ctx context.Context, _input *Transformation_complete_noIncubationInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _Transformation_complete_noIncubationSteps(_ctx context.Context, _input *Transformation_complete_noIncubationInput, _output *Transformation_complete_noIncubationOutput) {

	var compCells *wtype.LHComponent

	_input.CompetentCells.SetVolume(_input.MinimumCompetentCellVolume)

	// find all wells on plate
	allWells := _input.CompetentCellPlate.AllWellPositions(wtype.BYCOLUMN)

	var cellsFound bool

	// range through plate looking to see if the competent cells component is present
	for i := range allWells {
		if _input.CompetentCellPlate.WellMap()[allWells[i]].WContents.CName == _input.CompetentCells.CName && _input.CompetentCells.Volume().GreaterThan(_input.MinimumCompetentCellVolume) {
			compCells = _input.CompetentCellPlate.WellMap()[allWells[i]].WContents
			cellsFound = true
			break
		}
	}
	if !cellsFound {
		// range through plate until an empty one is found
		for i := range allWells {
			if _input.CompetentCellPlate.WellMap()[allWells[i]].Empty() {
				_input.CompetentCells.SetVolume(_input.MinimumCompetentCellVolume)
				_input.CompetentCellPlate.WellMap()[allWells[i]].WContents = _input.CompetentCells
				compCells = _input.CompetentCellPlate.WellMap()[allWells[i]].WContents
				cellsFound = true
				break
			}
		}
	}

	if !cellsFound {
		execute.Errorf(_ctx, "No %s Cells found on plate with greater than %s volume and no empty positions", _input.CompetentCells.CName, _input.MinimumCompetentCellVolume.ToString())
	} else {
		execute.SetInputPlate(_ctx, _input.CompetentCellPlate)
	}

	readycompetentcellsComp := compCells // readycompetentcells IS now a LHComponent
	/*
		readycompetentcellsComp := Incubate (readycompetentcells,Preplasmidtemp, Preplasmidtime, false) // we can incubate an LHComponent so this is fine
	*/

	DNAsample := mixer.Sample(_input.Reaction, _input.Reactionvolume)

	transformedcellsComp := execute.Mix(_ctx, readycompetentcellsComp, DNAsample)

	//transformedcellsComp := Incubate (transformedcells, Postplasmidtemp, Postplasmidtime,false)

	recoverymixture := mixer.Sample(_input.Recoverymedium, _input.Recoveryvolume)

	recoverymix2Comp := execute.Mix(_ctx, transformedcellsComp, recoverymixture)

	//recoverymix2Comp := Incubate (recoverymix2,  Recoverytemp, Recoverytime, true)

	plateout := mixer.Sample(recoverymix2Comp, _input.Plateoutvolume)
	platedculture := execute.MixInto(_ctx, _input.AgarPlate, "", plateout)

	_output.Platedculture = platedculture

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Transformation_complete_noIncubationAnalysis(_ctx context.Context, _input *Transformation_complete_noIncubationInput, _output *Transformation_complete_noIncubationOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _Transformation_complete_noIncubationValidation(_ctx context.Context, _input *Transformation_complete_noIncubationInput, _output *Transformation_complete_noIncubationOutput) {
}
func _Transformation_complete_noIncubationRun(_ctx context.Context, input *Transformation_complete_noIncubationInput) *Transformation_complete_noIncubationOutput {
	output := &Transformation_complete_noIncubationOutput{}
	_Transformation_complete_noIncubationSetup(_ctx, input)
	_Transformation_complete_noIncubationSteps(_ctx, input, output)
	_Transformation_complete_noIncubationAnalysis(_ctx, input, output)
	_Transformation_complete_noIncubationValidation(_ctx, input, output)
	return output
}

func Transformation_complete_noIncubationRunSteps(_ctx context.Context, input *Transformation_complete_noIncubationInput) *Transformation_complete_noIncubationSOutput {
	soutput := &Transformation_complete_noIncubationSOutput{}
	output := _Transformation_complete_noIncubationRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Transformation_complete_noIncubationNew() interface{} {
	return &Transformation_complete_noIncubationElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Transformation_complete_noIncubationInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Transformation_complete_noIncubationRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Transformation_complete_noIncubationInput{},
			Out: &Transformation_complete_noIncubationOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type Transformation_complete_noIncubationElement struct {
	inject.CheckedRunner
}

type Transformation_complete_noIncubationInput struct {
	AgarPlate                  *wtype.LHPlate
	CompetentCellPlate         *wtype.LHPlate
	CompetentCells             *wtype.LHComponent
	MinimumCompetentCellVolume wunit.Volume
	Plateoutvolume             wunit.Volume
	Postplasmidtemp            wunit.Temperature
	Postplasmidtime            wunit.Time
	Preplasmidtemp             wunit.Temperature
	Preplasmidtime             wunit.Time
	Reaction                   *wtype.LHComponent
	Reactionvolume             wunit.Volume
	Recoverymedium             *wtype.LHComponent
	Recoverytemp               wunit.Temperature
	Recoverytime               wunit.Time
	Recoveryvolume             wunit.Volume
}

type Transformation_complete_noIncubationOutput struct {
	Platedculture *wtype.LHComponent
}

type Transformation_complete_noIncubationSOutput struct {
	Data struct {
	}
	Outputs struct {
		Platedculture *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Transformation_complete_noIncubation",
		Constructor: Transformation_complete_noIncubationNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/Transformation/Transformation_noIncubation.an",
			Params: []component.ParamDesc{
				{Name: "AgarPlate", Desc: "", Kind: "Inputs"},
				{Name: "CompetentCellPlate", Desc: "", Kind: "Inputs"},
				{Name: "CompetentCells", Desc: "", Kind: "Inputs"},
				{Name: "MinimumCompetentCellVolume", Desc: "= 50.(uL)\n", Kind: "Parameters"},
				{Name: "Plateoutvolume", Desc: "", Kind: "Parameters"},
				{Name: "Postplasmidtemp", Desc: "", Kind: "Parameters"},
				{Name: "Postplasmidtime", Desc: "", Kind: "Parameters"},
				{Name: "Preplasmidtemp", Desc: "", Kind: "Parameters"},
				{Name: "Preplasmidtime", Desc: "", Kind: "Parameters"},
				{Name: "Reaction", Desc: "", Kind: "Inputs"},
				{Name: "Reactionvolume", Desc: "", Kind: "Parameters"},
				{Name: "Recoverymedium", Desc: "", Kind: "Inputs"},
				{Name: "Recoverytemp", Desc: "", Kind: "Parameters"},
				{Name: "Recoverytime", Desc: "= 2 (hours)\n", Kind: "Parameters"},
				{Name: "Recoveryvolume", Desc: "", Kind: "Parameters"},
				{Name: "Platedculture", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
