// Protocol Add_Solution allows for a new LHComponent (liquid handling component) description to be made when it does not exist in the LHComponent library.
// The element takes a user defined name, stock concentration and LHPolicy to apply to the NewSolution variable. The NewSolution variable must be based off
// of a TemplateComponent that already exists in the LHComponent library. The NewSolution output can be wired into elements as an input so that new LHComponents
// dont need to be made and populated into the library before an element can be used
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
	"github.com/antha-lang/antha/microArch/factory"
)

// Parameters to this protocol

// Name of new LHComponent, if empty defaults to TemplateComponent name
// Stock concentration being used, if empty defaults to TemplateComponent concentration, if there is no concentration associated with TemplateComponent it will not set a concentration
// If empty defaults to LHPolicy of TemplateComponent LHComponent

// Output data of this protocol

// Outputs a string to the terminal window saying what the new LHComponent is called, which LHcomponent it is based off of, the concentration of this component and the LHPolicy that should be used when handling this component.
// Outputs the NewSolution name as a string

// Physical inputs to this protocol

// Physical outputs to this protocol

// This is the NewSolution output that can be wired into another element and be used straight away without having to input it into the LHComponent library

// Conditions to run on startup
func _Add_SolutionSetup(_ctx context.Context, _input *Add_SolutionInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _Add_SolutionSteps(_ctx context.Context, _input *Add_SolutionInput, _output *Add_SolutionOutput) {

	//Initialise variable err with type error
	var err error

	TemplateComponent := factory.GetComponentByType("water")

	//Store the TemplateComponent name into a variable for text output at the end of the run in Status
	templateComponentName := TemplateComponent.CName

	//Initialise NewComponent variable with the LHComponent properties of the TemplateComponent

	_output.NewSolution = TemplateComponent

	//Sets the name of the NewSolution to the user specified Name parameter.
	//If name parameter is empty set to Name of TemplateComponent
	if _input.Name != "" {
		_output.NewSolution.CName = _input.Name
	} else {
		_output.NewSolution.CName = TemplateComponent.CName
	}

	//Sets the Concentration of the NewSolution to user specified StockConcentration
	//If left blank it defaults to the TemplateComponent concentration
	//If left blank and there is no concentration associated with TemplateComponent it will not set a concentration
	//Instead will provide a string of 0mM to be output in variable Status
	var NewSolutionConc string
	if _input.StockConcentration.RawValue() > 0.0 {
		_output.NewSolution.SetConcentration(_input.StockConcentration)
		NewSolutionConc = _output.NewSolution.Concentration().ToString()
	} else if TemplateComponent.HasConcentration() {
		_output.NewSolution.SetConcentration(TemplateComponent.Concentration())
		NewSolutionConc = _output.NewSolution.Concentration().ToString()
	} else {
		NewSolutionConc = "0mM"
	}

	//Sets the user defined LHPolicy in UseLHPolicy to use for the NewSolution
	//If an unkown LHPolicy is provided by the user an error will be generated
	//If left blank it defaults to the TemplateComponent LHPolicy type
	if _input.UseLHPolicy != "" {
		_output.NewSolution.Type, err = wtype.LiquidTypeFromString(_input.UseLHPolicy)
		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}
	} else {
		_output.NewSolution.Type = TemplateComponent.Type
	}

	//Provides a string output describing the new LHComponent
	_output.Status = _output.NewSolution.CName + " LHComponent created based on " + templateComponentName + " LHComponent, with a concentration of " + NewSolutionConc + " using the  " + _output.NewSolution.GetType() + " LHPolicy"

	//Outputs the new LHComponent name in Data
	_output.NewSolutionName = _output.NewSolution.CName
}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _Add_SolutionAnalysis(_ctx context.Context, _input *Add_SolutionInput, _output *Add_SolutionOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _Add_SolutionValidation(_ctx context.Context, _input *Add_SolutionInput, _output *Add_SolutionOutput) {

}
func _Add_SolutionRun(_ctx context.Context, input *Add_SolutionInput) *Add_SolutionOutput {
	output := &Add_SolutionOutput{}
	_Add_SolutionSetup(_ctx, input)
	_Add_SolutionSteps(_ctx, input, output)
	_Add_SolutionAnalysis(_ctx, input, output)
	_Add_SolutionValidation(_ctx, input, output)
	return output
}

func Add_SolutionRunSteps(_ctx context.Context, input *Add_SolutionInput) *Add_SolutionSOutput {
	soutput := &Add_SolutionSOutput{}
	output := _Add_SolutionRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Add_SolutionNew() interface{} {
	return &Add_SolutionElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Add_SolutionInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Add_SolutionRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Add_SolutionInput{},
			Out: &Add_SolutionOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type Add_SolutionElement struct {
	inject.CheckedRunner
}

type Add_SolutionInput struct {
	Name               string
	StockConcentration wunit.Concentration
	UseLHPolicy        wtype.PolicyName
}

type Add_SolutionOutput struct {
	NewSolution     *wtype.LHComponent
	NewSolutionName string
	Status          string
}

type Add_SolutionSOutput struct {
	Data struct {
		NewSolutionName string
		Status          string
	}
	Outputs struct {
		NewSolution *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Add_Solution",
		Constructor: Add_SolutionNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol Add_Solution allows for a new LHComponent (liquid handling component) description to be made when it does not exist in the LHComponent library.\nThe element takes a user defined name, stock concentration and LHPolicy to apply to the NewSolution variable. The NewSolution variable must be based off\nof a TemplateComponent that already exists in the LHComponent library. The NewSolution output can be wired into elements as an input so that new LHComponents\ndont need to be made and populated into the library before an element can be used\n",
			Path: "src/github.com/antha-lang/elements/starter/Add_Solution.an",
			Params: []component.ParamDesc{
				{Name: "Name", Desc: "Name of new LHComponent, if empty defaults to TemplateComponent name\n", Kind: "Parameters"},
				{Name: "StockConcentration", Desc: "Stock concentration being used, if empty defaults to TemplateComponent concentration, if there is no concentration associated with TemplateComponent it will not set a concentration\n", Kind: "Parameters"},
				{Name: "UseLHPolicy", Desc: "If empty defaults to LHPolicy of TemplateComponent LHComponent\n", Kind: "Parameters"},
				{Name: "NewSolution", Desc: "This is the NewSolution output that can be wired into another element and be used straight away without having to input it into the LHComponent library\n", Kind: "Outputs"},
				{Name: "NewSolutionName", Desc: "Outputs the NewSolution name as a string\n", Kind: "Data"},
				{Name: "Status", Desc: "Outputs a string to the terminal window saying what the new LHComponent is called, which LHcomponent it is based off of, the concentration of this component and the LHPolicy that should be used when handling this component.\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
