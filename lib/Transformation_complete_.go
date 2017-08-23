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

func _Transformation_completeRequirements() {
}

// Conditions to run on startup
func _Transformation_completeSetup(_ctx context.Context, _input *Transformation_completeInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _Transformation_completeSteps(_ctx context.Context, _input *Transformation_completeInput, _output *Transformation_completeOutput) {

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

	readycompetentcells := compCells

	readycompetentcellsComp := execute.Incubate(_ctx, readycompetentcells, _input.Preplasmidtemp, _input.Preplasmidtime, false) // we can incubate an LHComponent so this is fine

	DNAsample := mixer.Sample(_input.Reaction, _input.Reactionvolume)

	transformedcells := execute.Mix(_ctx, readycompetentcellsComp, DNAsample)

	transformedcellsComp := execute.Incubate(_ctx, transformedcells, _input.Postplasmidtemp, _input.Postplasmidtime, false)

	recoverymixture := mixer.Sample(_input.Recoverymedium, _input.Recoveryvolume)

	recoverymix2 := execute.Mix(_ctx, transformedcellsComp, recoverymixture)

	recoverymix2Comp := execute.Incubate(_ctx, recoverymix2, _input.Recoverytemp, _input.Recoverytime, true)

	plateout := mixer.Sample(recoverymix2Comp, _input.Plateoutvolume)
	platedculture := execute.MixInto(_ctx, _input.AgarPlate, "", plateout)

	_output.Platedculture = platedculture

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Transformation_completeAnalysis(_ctx context.Context, _input *Transformation_completeInput, _output *Transformation_completeOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _Transformation_completeValidation(_ctx context.Context, _input *Transformation_completeInput, _output *Transformation_completeOutput) {
}
func _Transformation_completeRun(_ctx context.Context, input *Transformation_completeInput) *Transformation_completeOutput {
	output := &Transformation_completeOutput{}
	_Transformation_completeSetup(_ctx, input)
	_Transformation_completeSteps(_ctx, input, output)
	_Transformation_completeAnalysis(_ctx, input, output)
	_Transformation_completeValidation(_ctx, input, output)
	return output
}

func Transformation_completeRunSteps(_ctx context.Context, input *Transformation_completeInput) *Transformation_completeSOutput {
	soutput := &Transformation_completeSOutput{}
	output := _Transformation_completeRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Transformation_completeNew() interface{} {
	return &Transformation_completeElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Transformation_completeInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Transformation_completeRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Transformation_completeInput{},
			Out: &Transformation_completeOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type Transformation_completeElement struct {
	inject.CheckedRunner
}

type Transformation_completeInput struct {
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

type Transformation_completeOutput struct {
	Platedculture *wtype.LHComponent
}

type Transformation_completeSOutput struct {
	Data struct {
	}
	Outputs struct {
		Platedculture *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Transformation_complete",
		Constructor: Transformation_completeNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/Transformation/Transformation_complete.an",
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
