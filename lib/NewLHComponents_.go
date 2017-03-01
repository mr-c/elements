// Protocol NewLHComponents allows for making a slice of new LHComponents (liquid handling component) when they do not exist in the LHComponent library.
// The element recursively calls the AddNewLHComponent element which takes a user defined name, stock concentration and LHPolicy to apply to the NewLHComponent variable. The NewLHComponent variable must be based off
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
	"github.com/antha-lang/antha/microArch/factory"
)

// Parameters to this protocol

// list of desired names for new LHComponents, if empty returns an error

// Stock concentration being used,
// if empty this defaults to TemplateComponent concentration,
// if only 1 entry, the entry for that will be used as a default for all
// if a "default" is specified then that will be used as the default for all entries with no value
// if there is no concentration associated with TemplateComponent and no default is specified, no concentration is set

// If empty this defaults to LHPolicy of TemplateComponent LHComponent,
// if only 1 entry this policy is used for all
// if a "default" is specified this policy is used for all entries with no value

// Output data of this protocol

// Outputs status to return to user on any substitutions made, what the new LHComponent is called, which LHcomponent it is based off of, the concentration of this component and the LHPolicy that should be used when handling this component.

// Outputs the NewLHComponent names

// Physical inputs to this protocol

// This TemplateComponent must be specified in the parameters file or the element will have a run time error,
// if length 1, the same template is used for all
// if "default" is specified that will be used as default for all entries with no value

// Physical outputs to this protocol

// This is the list of NewLHComponents output that can be wired into another element and be used straight away without having to input it into the LHComponent library

// Conditions to run on startup
func _NewLHComponentsSetup(_ctx context.Context, _input *NewLHComponentsInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _NewLHComponentsSteps(_ctx context.Context, _input *NewLHComponentsInput, _output *NewLHComponentsOutput) {

	if len(_input.Names) == 0 {
		execute.Errorf(_ctx, "No Names specified for new components")
	}

	// initialise default conc as empty, if not found in map, no concentration will be set unless a "default" key is used in the StockConcentrations map
	var defaultConc wunit.Concentration

	if _, found := _input.StockConcentrations["default"]; found {
		defaultConc = _input.StockConcentrations["default"]
	}

	// if the length of the map is 1 this template will be used for all components
	// if empty the protocol will terminate with an error
	var defaultTemplate *wtype.LHComponent

	if len(_input.TemplateComponents) == 1 {
		for _, v := range _input.TemplateComponents {
			defaultTemplate = v
		}
	} else if _, found := _input.TemplateComponents["default"]; found {
		defaultTemplate = _input.TemplateComponents["default"]
	} else if len(_input.TemplateComponents) == 0 {
		execute.Errorf(_ctx, "No template components specified")
	}

	// if the length of the map is 1 this lhpolicy will be used for all components
	// if empty the lhpolicy of the Template Component is used
	var defaultLHPolicy string

	if len(_input.UseLHPolicy) == 1 {
		for _, v := range _input.UseLHPolicy {
			defaultLHPolicy = v
		}
	}

	// initialise map for appending with results
	_output.Status = make(map[string]string)

	// range through component names
	for _, name := range _input.Names {

		var status string
		var stockConc wunit.Concentration
		var template *wtype.LHComponent
		var lhpolicy string
		var found bool

		// check if a concentration is specified
		if stockConc, found = _input.StockConcentrations[name]; !found {
			stockConc = defaultConc
			status = "No concentration specified for " + name + "; "
		}

		// check if a template component is specified
		if template, found = _input.TemplateComponents[name]; !found {
			template = factory.GetComponentByType(defaultTemplate.CName)
			status = status + "No template specified so using default " + defaultTemplate.CName + "; "
		}

		// check if an LHPolicy is specified
		if lhpolicy, found = _input.UseLHPolicy[name]; !found {
			lhpolicy = defaultLHPolicy
			status = status + "No lhpolicy specified so using policy " + template.TypeName() + " from default component " + defaultTemplate.CName + "; "
		}

		// run AddNewLHComponent element
		result := AddNewLHComponentRunSteps(_ctx, &AddNewLHComponentInput{Name: name,
			StockConcentration: stockConc,
			UseLHPolicy:        lhpolicy,

			TemplateComponent: template},
		)

		// append outputs
		_output.NewLHComponents = append(_output.NewLHComponents, result.Outputs.NewLHComponent)
		_output.NewLHComponentNames = append(_output.NewLHComponentNames, result.Data.NewLHComponentName)
		_output.Status[name] = status + result.Data.Status
	}
	// done
}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _NewLHComponentsAnalysis(_ctx context.Context, _input *NewLHComponentsInput, _output *NewLHComponentsOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _NewLHComponentsValidation(_ctx context.Context, _input *NewLHComponentsInput, _output *NewLHComponentsOutput) {

}
func _NewLHComponentsRun(_ctx context.Context, input *NewLHComponentsInput) *NewLHComponentsOutput {
	output := &NewLHComponentsOutput{}
	_NewLHComponentsSetup(_ctx, input)
	_NewLHComponentsSteps(_ctx, input, output)
	_NewLHComponentsAnalysis(_ctx, input, output)
	_NewLHComponentsValidation(_ctx, input, output)
	return output
}

func NewLHComponentsRunSteps(_ctx context.Context, input *NewLHComponentsInput) *NewLHComponentsSOutput {
	soutput := &NewLHComponentsSOutput{}
	output := _NewLHComponentsRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func NewLHComponentsNew() interface{} {
	return &NewLHComponentsElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &NewLHComponentsInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _NewLHComponentsRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &NewLHComponentsInput{},
			Out: &NewLHComponentsOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type NewLHComponentsElement struct {
	inject.CheckedRunner
}

type NewLHComponentsInput struct {
	Names               []string
	StockConcentrations map[string]wunit.Concentration
	TemplateComponents  map[string]*wtype.LHComponent
	UseLHPolicy         map[string]string
}

type NewLHComponentsOutput struct {
	NewLHComponentNames []string
	NewLHComponents     []*wtype.LHComponent
	Status              map[string]string
}

type NewLHComponentsSOutput struct {
	Data struct {
		NewLHComponentNames []string
		Status              map[string]string
	}
	Outputs struct {
		NewLHComponents []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "NewLHComponents",
		Constructor: NewLHComponentsNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol NewLHComponents allows for making a slice of new LHComponents (liquid handling component) when they do not exist in the LHComponent library.\nThe element recursively calls the AddNewLHComponent element which takes a user defined name, stock concentration and LHPolicy to apply to the NewLHComponent variable. The NewLHComponent variable must be based off\nof a TemplateComponent that already exists in the LHComponent library. The NewLHComponent output can be wired into elements as an input so that new LHComponents\ndont need to be made and populated into the library before an element can be used\n",
			Path: "src/github.com/antha-lang/elements/an/Utility/NewLHComponents.an",
			Params: []component.ParamDesc{
				{Name: "Names", Desc: "list of desired names for new LHComponents, if empty returns an error\n", Kind: "Parameters"},
				{Name: "StockConcentrations", Desc: "Stock concentration being used,\nif empty this defaults to TemplateComponent concentration,\nif only 1 entry, the entry for that will be used as a default for all\nif a \"default\" is specified then that will be used as the default for all entries with no value\nif there is no concentration associated with TemplateComponent and no default is specified, no concentration is set\n", Kind: "Parameters"},
				{Name: "TemplateComponents", Desc: "This TemplateComponent must be specified in the parameters file or the element will have a run time error,\nif length 1, the same template is used for all\nif \"default\" is specified that will be used as default for all entries with no value\n", Kind: "Inputs"},
				{Name: "UseLHPolicy", Desc: "If empty this defaults to LHPolicy of TemplateComponent LHComponent,\nif only 1 entry this policy is used for all\nif a \"default\" is specified this policy is used for all entries with no value\n", Kind: "Parameters"},
				{Name: "NewLHComponentNames", Desc: "Outputs the NewLHComponent names\n", Kind: "Data"},
				{Name: "NewLHComponents", Desc: "This is the list of NewLHComponents output that can be wired into another element and be used straight away without having to input it into the LHComponent library\n", Kind: "Outputs"},
				{Name: "Status", Desc: "Outputs status to return to user on any substitutions made, what the new LHComponent is called, which LHcomponent it is based off of, the concentration of this component and the LHPolicy that should be used when handling this component.\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
