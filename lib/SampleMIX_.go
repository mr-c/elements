// Example protocol demonstrating the use of the Sample function
package lib

import // this is the name of the protocol that will be called in a workflow or other antha element

// we need to import the wtype package to use the LHComponent type
// the mixer package is required to use the Sample function
(
	"context"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol (data)

// antha, like golang is a strongly typed language in which the type of a variable must be declared.
// In this case we're creating a variable called SampleVolume which is of type Volume;
// the type system allows the antha compiler to catch many types of common errors before the programme is run
// the antha type system extends this to biological types such as volumes here.
// functions require inputs of particular types to be adhered to

// Data which is returned from this protocol, and data types

// Antha inherits all standard primitives valid in golang;
//for example the string type shown here used to return a textual message

// Physical Inputs to this protocol with types

// the LHComponent is the principal liquidhandling type in antha
// the * signifies that this is a pointer to the component rather than the component itself
// most key antha functions such as Sample and Mix use *LHComponent rather than LHComponent
// since the type is imported from the wtype package we need to use  *wtype.LHComponent rather than simply *LHComponent

// Physical outputs from this protocol with types

// An output LHComponent variable is created called Sample

func _SampleMIXRequirements() {

}

// Conditions to run on startup
func _SampleMIXSetup(_ctx context.Context, _input *SampleMIXInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _SampleMIXSteps(_ctx context.Context, _input *SampleMIXInput, _output *SampleMIXOutput) {

	// the Sample function is imported from the mixer library
	// in the mixer library the function signature can be found, here it is:
	// func Sample(l *wtype.LHComponent, v wunit.Volume) *wtype.LHComponent {
	// The function signature  shows that the function requires a *LHComponent and a Volume and returns an *LHComponent
	var sample *wtype.LHComponent
	sample = mixer.Sample(_input.Solution, _input.SampleVolume)
	samplemix := execute.Mix(_ctx, sample)

	// The Sample function is not sufficient to generate liquid handling instructions alone,
	// We would need a Mix command to instruct where to put the sample

	// we can also create data outputs as a string like this
	_output.Status = _input.SampleVolume.ToString() + " of " + _input.Solution.CName + " sampled and mixed"

	// To maintain good practice in coding all variables within the steps section should be lower case
	// When that variable is to become an output (or is an input) the first letter is capitalised as shown below (or in CamelCase)

	_output.SampleMix = samplemix
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _SampleMIXAnalysis(_ctx context.Context, _input *SampleMIXInput, _output *SampleMIXOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _SampleMIXValidation(_ctx context.Context, _input *SampleMIXInput, _output *SampleMIXOutput) {

}
func _SampleMIXRun(_ctx context.Context, input *SampleMIXInput) *SampleMIXOutput {
	output := &SampleMIXOutput{}
	_SampleMIXSetup(_ctx, input)
	_SampleMIXSteps(_ctx, input, output)
	_SampleMIXAnalysis(_ctx, input, output)
	_SampleMIXValidation(_ctx, input, output)
	return output
}

func SampleMIXRunSteps(_ctx context.Context, input *SampleMIXInput) *SampleMIXSOutput {
	soutput := &SampleMIXSOutput{}
	output := _SampleMIXRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SampleMIXNew() interface{} {
	return &SampleMIXElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SampleMIXInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SampleMIXRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SampleMIXInput{},
			Out: &SampleMIXOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type SampleMIXElement struct {
	inject.CheckedRunner
}

type SampleMIXInput struct {
	SampleVolume wunit.Volume
	Solution     *wtype.LHComponent
}

type SampleMIXOutput struct {
	SampleMix *wtype.LHComponent
	Status    string
}

type SampleMIXSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
		SampleMix *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "SampleMIX",
		Constructor: SampleMIXNew,
		Desc: component.ComponentDesc{
			Desc: "Example protocol demonstrating the use of the Sample function\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson2_Sample/JAJALesson2/2A_SampleMIX/A_SampleMIX.an",
			Params: []component.ParamDesc{
				{Name: "SampleVolume", Desc: "antha, like golang is a strongly typed language in which the type of a variable must be declared.\nIn this case we're creating a variable called SampleVolume which is of type Volume;\nthe type system allows the antha compiler to catch many types of common errors before the programme is run\nthe antha type system extends this to biological types such as volumes here.\nfunctions require inputs of particular types to be adhered to\n", Kind: "Parameters"},
				{Name: "Solution", Desc: "the LHComponent is the principal liquidhandling type in antha\nthe * signifies that this is a pointer to the component rather than the component itself\nmost key antha functions such as Sample and Mix use *LHComponent rather than LHComponent\nsince the type is imported from the wtype package we need to use  *wtype.LHComponent rather than simply *LHComponent\n", Kind: "Inputs"},
				{Name: "SampleMix", Desc: "An output LHComponent variable is created called Sample\n", Kind: "Outputs"},
				{Name: "Status", Desc: "Antha inherits all standard primitives valid in golang;\nfor example the string type shown here used to return a textual message\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
