// This protocol allows an input plate to be set up with a series of solutions on it according to the user's specification.
package lib

import

// Place golang packages to import here
(
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Parameters to this protocol

// optional parameter to specify the well positions for a particular solution
// This will conflict with ByRow if both are used together so an error will be returned in that situation.

// If selected the solutions will be added by row (A1, A2, A3 ...), the default is by column (A1, B1, C1 ...).
// This will conflict with the SpecifyWellPositions if both are used together so an error will be returned in that situation.

// Specify whether replicate samples are to be set up

// Specify the volumes per sample. If left blank the components must have a volume otherwise an error will occur.
// If the component has a volume and a volume is specified here the specified volume takes precedent.

// Output data of this protocol

// Physical inputs to this protocol

// input solutions to add to the plate.

// type of plate to add solutions to.

// Physical outputs to this protocol

// Conditions to run on startup
func _Setup_Input_PlateSetup(_ctx context.Context, _input *Setup_Input_PlateInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _Setup_Input_PlateSteps(_ctx context.Context, _input *Setup_Input_PlateInput, _output *Setup_Input_PlateOutput) {

	_output.PlateWithSolutions = _input.PlateType.Dup()

	var newlistofsolutions []*wtype.LHComponent

	for i := range _input.Solutions {
		newlistofsolutions = append(newlistofsolutions, _input.Solutions[i])
	}

	// first go through all solutions which have a well position explicitely specified
	for solutionName, wells := range _input.SpecifyWellPositions {

		// find solution in the list of solutions
		solution, err := findSolutionByName(_input.Solutions, solutionName)

		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}

		// find if volume has been specified
		vol, found := _input.Volumes[solutionName]

		if !found {
			vol, found = _input.Volumes["default"]
			if !found {
				if solution.Volume().RawValue() == 0.0 {
					execute.Errorf(_ctx, "Solution %s has no volume and no volume or default volume specified in Volumes parameter", solution.Name())
				} else {
					vol = solution.Volume()
				}
			}
		}

		solution.SetVolume(vol)

		replicates, found := _input.Replicates[solutionName]

		if !found {
			replicates, found = _input.Replicates["default"]
			if !found {
				replicates = 1
			}
		}

		for i := 0; i < replicates; i++ {

			// get all plate components and return into both a slice and a map
			for _, wellcontents := range wells {

				if _output.PlateWithSolutions.WellMap()[wellcontents].Empty() {

					_output.PlateWithSolutions.WellMap()[wellcontents].WContents = solution
					_output.SolutionsOnPlate = append(_output.SolutionsOnPlate, _output.PlateWithSolutions.WellMap()[wellcontents].WContents)
					break
				} else {
					execute.Errorf(_ctx, "Solution %s specified to add to location %s but a sample %s is already present at that position.", solution.Name(), wellcontents, _output.PlateWithSolutions.WellMap()[wellcontents].WContents.Name())
				}
			}
			solution = solution.Dup()
		}
		newlistofsolutions, err = removeSolutionFromList(newlistofsolutions, solution.Name())

		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}
	}

	// now add the rest of the solutions
	for _, solution := range newlistofsolutions {

		vol, found := _input.Volumes[solution.Name()]

		if !found {
			vol, found = _input.Volumes["default"]
			if !found {
				if solution.Volume().RawValue() == 0.0 {
					execute.Errorf(_ctx, "Solution %s has no volume and no volume or default volume specified in Volumes parameter", solution.Name())
				} else {
					vol = solution.Volume()
				}
			}
		}

		solution.SetVolume(vol)

		replicates, found := _input.Replicates[solution.Name()]

		if !found {
			replicates, found = _input.Replicates["default"]
			if !found {
				replicates = 1
			}
		}

		for i := 0; i < replicates; i++ {

			// get all plate components and return into both a slice and a map
			for _, wellcontents := range _output.PlateWithSolutions.AllWellPositions(_input.ByRow) {

				if _output.PlateWithSolutions.WellMap()[wellcontents].Empty() {

					_output.PlateWithSolutions.WellMap()[wellcontents].WContents = solution
					_output.SolutionsOnPlate = append(_output.SolutionsOnPlate, _output.PlateWithSolutions.WellMap()[wellcontents].WContents)
					break
				}
			}
		}
	}

	execute.SetInputPlate(_ctx, _output.PlateWithSolutions)
}

func findSolutionByName(solutions []*wtype.LHComponent, solutionName string) (solution *wtype.LHComponent, err error) {
	for _, sol := range solutions {
		if sol.Name() == solutionName {
			return sol, nil
		}
	}
	return solution, fmt.Errorf("solution %s not found in solutions", solutionName)
}

func removeSolutionFromList(solutions []*wtype.LHComponent, solutionName string) (newsolutions []*wtype.LHComponent, err error) {
	var exclude int
	for i, sol := range solutions {
		if sol.Name() == solutionName {
			exclude = i
			break
		}
	}

	for i, sol := range solutions {

		if i != exclude {
			newsolutions = append(newsolutions, sol)
		}
	}

	if len(solutions) == len(newsolutions) {
		err = fmt.Errorf("solution %s not found in solutions", solutionName)
	}

	return newsolutions, err
}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _Setup_Input_PlateAnalysis(_ctx context.Context, _input *Setup_Input_PlateInput, _output *Setup_Input_PlateOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _Setup_Input_PlateValidation(_ctx context.Context, _input *Setup_Input_PlateInput, _output *Setup_Input_PlateOutput) {

}
func _Setup_Input_PlateRun(_ctx context.Context, input *Setup_Input_PlateInput) *Setup_Input_PlateOutput {
	output := &Setup_Input_PlateOutput{}
	_Setup_Input_PlateSetup(_ctx, input)
	_Setup_Input_PlateSteps(_ctx, input, output)
	_Setup_Input_PlateAnalysis(_ctx, input, output)
	_Setup_Input_PlateValidation(_ctx, input, output)
	return output
}

func Setup_Input_PlateRunSteps(_ctx context.Context, input *Setup_Input_PlateInput) *Setup_Input_PlateSOutput {
	soutput := &Setup_Input_PlateSOutput{}
	output := _Setup_Input_PlateRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Setup_Input_PlateNew() interface{} {
	return &Setup_Input_PlateElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Setup_Input_PlateInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Setup_Input_PlateRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Setup_Input_PlateInput{},
			Out: &Setup_Input_PlateOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type Setup_Input_PlateElement struct {
	inject.CheckedRunner
}

type Setup_Input_PlateInput struct {
	ByRow                bool
	PlateType            *wtype.LHPlate
	Replicates           map[string]int
	Solutions            []*wtype.LHComponent
	SpecifyWellPositions map[string][]string
	Volumes              map[string]wunit.Volume
}

type Setup_Input_PlateOutput struct {
	PlateWithSolutions *wtype.LHPlate
	SolutionsOnPlate   []*wtype.LHComponent
	Sum                float64
}

type Setup_Input_PlateSOutput struct {
	Data struct {
		Sum float64
	}
	Outputs struct {
		PlateWithSolutions *wtype.LHPlate
		SolutionsOnPlate   []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Setup_Input_Plate",
		Constructor: Setup_Input_PlateNew,
		Desc: component.ComponentDesc{
			Desc: "This protocol allows an input plate to be set up with a series of solutions on it according to the user's specification.\n",
			Path: "src/github.com/antha-lang/elements/an/Utility/Setup_Input_Plate/Setup_Input_Plate.an",
			Params: []component.ParamDesc{
				{Name: "ByRow", Desc: "If selected the solutions will be added by row (A1, A2, A3 ...), the default is by column (A1, B1, C1 ...).\nThis will conflict with the SpecifyWellPositions if both are used together so an error will be returned in that situation.\n", Kind: "Parameters"},
				{Name: "PlateType", Desc: "type of plate to add solutions to.\n", Kind: "Inputs"},
				{Name: "Replicates", Desc: "Specify whether replicate samples are to be set up\n", Kind: "Parameters"},
				{Name: "Solutions", Desc: "input solutions to add to the plate.\n", Kind: "Inputs"},
				{Name: "SpecifyWellPositions", Desc: "optional parameter to specify the well positions for a particular solution\nThis will conflict with ByRow if both are used together so an error will be returned in that situation.\n", Kind: "Parameters"},
				{Name: "Volumes", Desc: "Specify the volumes per sample. If left blank the components must have a volume otherwise an error will occur.\nIf the component has a volume and a volume is specified here the specified volume takes precedent.\n", Kind: "Parameters"},
				{Name: "PlateWithSolutions", Desc: "", Kind: "Outputs"},
				{Name: "SolutionsOnPlate", Desc: "", Kind: "Outputs"},
				{Name: "Sum", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
