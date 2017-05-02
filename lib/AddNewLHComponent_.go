// Protocol AddNewLHComponent allows for a new LHComponent (liquid handling component) description to be made when it does not exist in the LHComponent library.
// The element takes a user defined name, stock concentration and LHPolicy to apply to the NewLHComponent variable. The NewLHComponent variable must be based off
// of a TemplateComponent that already exists in the LHComponent library. The NewLHComponent output can be wired into elements as an input so that new LHComponents
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
)

// Parameters to this protocol

// Name of new LHComponent, if empty defaults to TemplateComponent name
// Stock concentration being used, if empty defaults to TemplateComponent concentration, if there is no concentration associated with TemplateComponent it will not set a concentration
// If empty defaults to LHPolicy of TemplateComponent LHComponent

// Output data of this protocol

// Outputs a string to the terminal window saying what the new LHComponent is called, which LHcomponent it is based off of, the concentration of this component and the LHPolicy that should be used when handling this component.
// Outputs the NewLHComponent name as a string

// Physical inputs to this protocol

// This TemplateComponent must be specified in the parameters file or the element will have a run time error

// Physical outputs to this protocol

// This is the NewLHComponent output that can be wired into another element and be used straight away without having to input it into the LHComponent library

// Conditions to run on startup
func _AddNewLHComponentSetup(_ctx context.Context, _input *AddNewLHComponentInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _AddNewLHComponentSteps(_ctx context.Context, _input *AddNewLHComponentInput, _output *AddNewLHComponentOutput) {

	//Initialise variable err with type error
	var err error

	//Store the TemplateComponent name into a variable for text output at the end of the run in Status
	templateComponentName := _input.TemplateComponent.CName

	//Initialise NewComponent variable with the LHComponent properties of the TemplateComponent

	_output.NewLHComponent = _input.TemplateComponent

	//Sets the name of the NewLHComponent to the user specified Name parameter.
	//If name parameter is empty set to Name of TemplateComponent
	if _input.Name != "" {
		_output.NewLHComponent.CName = _input.Name
	} else {
		_output.NewLHComponent.CName = _input.TemplateComponent.CName
	}

	//Sets the Concentration of the NewLHComponent to user specified StockConcentration
	//If left blank it defaults to the TemplateComponent concentration
	//If left blank and there is no concentration associated with TemplateComponent it will not set a concentration
	//Instead will provide a string of 0mM to be output in variable Status
	var NewLHComponentConc string
	if _input.StockConcentration.RawValue() > 0.0 {
		_output.NewLHComponent.SetConcentration(_input.StockConcentration)
		NewLHComponentConc = _output.NewLHComponent.Concentration().ToString()
	} else if _input.TemplateComponent.HasConcentration() {
		_output.NewLHComponent.SetConcentration(_input.TemplateComponent.Concentration())
		NewLHComponentConc = _output.NewLHComponent.Concentration().ToString()
	} else {
		NewLHComponentConc = "0mM"
	}

	//Sets the user defined LHPolicy in UseLHPolicy to use for the NewLHComponent
	//If an unkown LHPolicy is provided by the user an error will be generated
	//If left blank it defaults to the TemplateComponent LHPolicy type
	if _input.UseLHPolicy != "" {
		_output.NewLHComponent.Type, err = wtype.LiquidTypeFromString(_input.UseLHPolicy)
		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}
	} else {
		_output.NewLHComponent.Type = _input.TemplateComponent.Type
	}

	//Provides a string output describing the new LHComponent
	_output.Status = _output.NewLHComponent.CName + " LHComponent created based on " + templateComponentName + " LHComponent, with a concentration of " + NewLHComponentConc + " using the  " + _output.NewLHComponent.GetType() + " LHPolicy"

	//Outputs the new LHComponent name in Data
	_output.NewLHComponentName = _output.NewLHComponent.CName
}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _AddNewLHComponentAnalysis(_ctx context.Context, _input *AddNewLHComponentInput, _output *AddNewLHComponentOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _AddNewLHComponentValidation(_ctx context.Context, _input *AddNewLHComponentInput, _output *AddNewLHComponentOutput) {

}
func _AddNewLHComponentRun(_ctx context.Context, input *AddNewLHComponentInput) *AddNewLHComponentOutput {
	output := &AddNewLHComponentOutput{}
	_AddNewLHComponentSetup(_ctx, input)
	_AddNewLHComponentSteps(_ctx, input, output)
	_AddNewLHComponentAnalysis(_ctx, input, output)
	_AddNewLHComponentValidation(_ctx, input, output)
	return output
}

func AddNewLHComponentRunSteps(_ctx context.Context, input *AddNewLHComponentInput) *AddNewLHComponentSOutput {
	soutput := &AddNewLHComponentSOutput{}
	output := _AddNewLHComponentRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AddNewLHComponentNew() interface{} {
	return &AddNewLHComponentElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AddNewLHComponentInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AddNewLHComponentRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AddNewLHComponentInput{},
			Out: &AddNewLHComponentOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type AddNewLHComponentElement struct {
	inject.CheckedRunner
}

type AddNewLHComponentInput struct {
	Name               string
	StockConcentration wunit.Concentration
	TemplateComponent  *wtype.LHComponent
	UseLHPolicy        string
}

type AddNewLHComponentOutput struct {
	NewLHComponent     *wtype.LHComponent
	NewLHComponentName string
	Status             string
}

type AddNewLHComponentSOutput struct {
	Data struct {
		NewLHComponentName string
		Status             string
	}
	Outputs struct {
		NewLHComponent *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AddNewLHComponent",
		Constructor: AddNewLHComponentNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol AddNewLHComponent allows for a new LHComponent (liquid handling component) description to be made when it does not exist in the LHComponent library.\nThe element takes a user defined name, stock concentration and LHPolicy to apply to the NewLHComponent variable. The NewLHComponent variable must be based off\nof a TemplateComponent that already exists in the LHComponent library. The NewLHComponent output can be wired into elements as an input so that new LHComponents\ndont need to be made and populated into the library before an element can be used\n",
			Path: "src/github.com/antha-lang/elements/starter/AddNewLHComponent.an",
			Params: []component.ParamDesc{
				{Name: "Name", Desc: "Name of new LHComponent, if empty defaults to TemplateComponent name\n", Kind: "Parameters"},
				{Name: "StockConcentration", Desc: "Stock concentration being used, if empty defaults to TemplateComponent concentration, if there is no concentration associated with TemplateComponent it will not set a concentration\n", Kind: "Parameters"},
				{Name: "TemplateComponent", Desc: "This TemplateComponent must be specified in the parameters file or the element will have a run time error\n", Kind: "Inputs"},
				{Name: "UseLHPolicy", Desc: "If empty defaults to LHPolicy of TemplateComponent LHComponent\n", Kind: "Parameters"},
				{Name: "NewLHComponent", Desc: "This is the NewLHComponent output that can be wired into another element and be used straight away without having to input it into the LHComponent library\n", Kind: "Outputs"},
				{Name: "NewLHComponentName", Desc: "Outputs the NewLHComponent name as a string\n", Kind: "Data"},
				{Name: "Status", Desc: "Outputs a string to the terminal window saying what the new LHComponent is called, which LHcomponent it is based off of, the concentration of this component and the LHPolicy that should be used when handling this component.\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
