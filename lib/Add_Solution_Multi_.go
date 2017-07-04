// Protocol Add_Solution_Multi allows for making a slice of new LHComponents (liquid handling component) when they do not exist in the LHComponent library.
// The element recursively calls the Add_Solution element which takes a user defined name, stock concentration and LHPolicy to apply to the NewSolution variable.
// The NewSolution output can be wired into elements as an input so that new LHComponents
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

// list of desired names for new LHComponents, if empty returns an error

// Stock concentration being used,
// if a "default" is specified then that will be used as the default for all entries with no value
// if no default is specified, no concentration is set

// If empty this defaults to PostMix which mixes 3 times after dispensing.

// Output data of this protocol

// Outputs status to return to user on any substitutions made, what the new LHComponent is called, which LHcomponent it is based off of, the concentration of this component and the LHPolicy that should be used when handling this component.

// Outputs the NewSolution names

// Physical inputs to this protocol

// Physical outputs to this protocol

// This is the list of NewSolutions output that can be wired into another element and be used straight away without having to input it into the LHComponent library

// Conditions to run on startup
func _Add_Solution_MultiSetup(_ctx context.Context, _input *Add_Solution_MultiInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _Add_Solution_MultiSteps(_ctx context.Context, _input *Add_Solution_MultiInput, _output *Add_Solution_MultiOutput) {

	if len(_input.Names) == 0 {
		execute.Errorf(_ctx, "No Names specified for new components")
	}

	// initialise default conc as empty, if not found in map, no concentration will be set unless a "default" key is used in the StockConcentrations map
	var defaultConc wunit.Concentration

	if _, found := _input.StockConcentrations["default"]; found {
		defaultConc = _input.StockConcentrations["default"]
	}

	// if the length of the map is 1 this lhpolicy will be used for all components
	// if empty the lhpolicy of the Template Component is used
	var defaultLHPolicy wtype.PolicyName

	if policy, found := _input.UseLHPolicy["default"]; found {
		defaultLHPolicy = policy
	} else {
		defaultLHPolicy = ""
	}

	// initialise map for appending with results
	_output.Status = make(map[string]string)

	// range through component names
	for _, name := range _input.Names {

		var status string
		var stockConc wunit.Concentration
		var lhpolicy wtype.PolicyName
		var found bool

		// check if a concentration is specified
		if stockConc, found = _input.StockConcentrations[name]; !found {
			stockConc = defaultConc
			status = "No concentration specified for " + name + "; "
		}

		// check if an LHPolicy is specified
		if lhpolicy, found = _input.UseLHPolicy[name]; !found {
			lhpolicy = defaultLHPolicy
			status = status + "No lhpolicy specified so using default PostMix policy; "
		}

		// run Add_Solution element
		result := Add_SolutionRunSteps(_ctx, &Add_SolutionInput{Name: name,
			StockConcentration: stockConc,
			UseLHPolicy:        lhpolicy},
		)

		// append outputs
		_output.NewSolutions = append(_output.NewSolutions, result.Outputs.NewSolution)
		_output.NewSolutionNames = append(_output.NewSolutionNames, result.Data.NewSolutionName)
		_output.Status[name] = status + result.Data.Status
	}
	// done
}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _Add_Solution_MultiAnalysis(_ctx context.Context, _input *Add_Solution_MultiInput, _output *Add_Solution_MultiOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _Add_Solution_MultiValidation(_ctx context.Context, _input *Add_Solution_MultiInput, _output *Add_Solution_MultiOutput) {

}
func _Add_Solution_MultiRun(_ctx context.Context, input *Add_Solution_MultiInput) *Add_Solution_MultiOutput {
	output := &Add_Solution_MultiOutput{}
	_Add_Solution_MultiSetup(_ctx, input)
	_Add_Solution_MultiSteps(_ctx, input, output)
	_Add_Solution_MultiAnalysis(_ctx, input, output)
	_Add_Solution_MultiValidation(_ctx, input, output)
	return output
}

func Add_Solution_MultiRunSteps(_ctx context.Context, input *Add_Solution_MultiInput) *Add_Solution_MultiSOutput {
	soutput := &Add_Solution_MultiSOutput{}
	output := _Add_Solution_MultiRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Add_Solution_MultiNew() interface{} {
	return &Add_Solution_MultiElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Add_Solution_MultiInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Add_Solution_MultiRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Add_Solution_MultiInput{},
			Out: &Add_Solution_MultiOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type Add_Solution_MultiElement struct {
	inject.CheckedRunner
}

type Add_Solution_MultiInput struct {
	Names               []string
	StockConcentrations map[string]wunit.Concentration
	UseLHPolicy         map[string]wtype.PolicyName
}

type Add_Solution_MultiOutput struct {
	NewSolutionNames []string
	NewSolutions     []*wtype.LHComponent
	Status           map[string]string
}

type Add_Solution_MultiSOutput struct {
	Data struct {
		NewSolutionNames []string
		Status           map[string]string
	}
	Outputs struct {
		NewSolutions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Add_Solution_Multi",
		Constructor: Add_Solution_MultiNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol Add_Solution_Multi allows for making a slice of new LHComponents (liquid handling component) when they do not exist in the LHComponent library.\nThe element recursively calls the Add_Solution element which takes a user defined name, stock concentration and LHPolicy to apply to the NewSolution variable.\nThe NewSolution output can be wired into elements as an input so that new LHComponents\ndont need to be made and populated into the library before an element can be used\n",
			Path: "src/github.com/antha-lang/elements/starter/Add_Solution_Multi.an",
			Params: []component.ParamDesc{
				{Name: "Names", Desc: "list of desired names for new LHComponents, if empty returns an error\n", Kind: "Parameters"},
				{Name: "StockConcentrations", Desc: "Stock concentration being used,\nif a \"default\" is specified then that will be used as the default for all entries with no value\nif no default is specified, no concentration is set\n", Kind: "Parameters"},
				{Name: "UseLHPolicy", Desc: "If empty this defaults to PostMix which mixes 3 times after dispensing.\n", Kind: "Parameters"},
				{Name: "NewSolutionNames", Desc: "Outputs the NewSolution names\n", Kind: "Data"},
				{Name: "NewSolutions", Desc: "This is the list of NewSolutions output that can be wired into another element and be used straight away without having to input it into the LHComponent library\n", Kind: "Outputs"},
				{Name: "Status", Desc: "Outputs status to return to user on any substitutions made, what the new LHComponent is called, which LHcomponent it is based off of, the concentration of this component and the LHPolicy that should be used when handling this component.\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
