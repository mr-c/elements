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
// the type system allows the antha compiler to catch many types of common errors before the programme
// is run the antha type system extends this to biological types such as volumes here.
// functions require inputs of particular types to be adhered to

// Data which is returned from this protocol, and data types

// Antha inherits all standard primitives valid in golang;
// for example the string type shown here used to return a textual message

// Physical Inputs to this protocol with types

// the LHComponent is the principal liquidhandling type in antha
// the * signifies that this is a pointer to the component rather than the component itself
// most key antha functions such as Sample and Mix use *LHComponent rather than LHComponent
// since the type is imported from the wtype package we need to use  *wtype.LHComponent rather
// than simply *LHComponent

// Physical outputs from this protocol with types

// An output LHComponent variable is created called Sample

func _SampleARequirements() {
}

// Conditions to run on startup
func _SampleASetup(_ctx context.Context, _input *SampleAInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _SampleASteps(_ctx context.Context, _input *SampleAInput, _output *SampleAOutput) {

	// the Sample function is imported from the mixer library
	// in the mixer library the function signature can be found, here it is:
	// func Sample(l *wtype.LHComponent, v wunit.Volume) *wtype.LHComponent {
	// The function signature  shows that the function requires a *LHComponent and a Volume and returns
	// an *LHComponent
	var sampleA *wtype.LHComponent
	sampleA = mixer.Sample(_input.Solution, _input.SampleVolume)

	// To maintain good practice in coding all variables within the steps section should be lower case
	// when that variable is to become an output (or is an input) the first letter is capitalised as shown below.
	// This is to do with variable visibility, a lower case variable is only visible to the code within this
	// element (private), a capitalised variable is public and can be seen/accessed by other elements (public)
	_output.SampleA = sampleA

	// we can also create data outputs as a string like this
	_output.Status = _input.SampleVolume.ToString() + " of " + _input.Solution.CName + " sampled"
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _SampleAAnalysis(_ctx context.Context, _input *SampleAInput, _output *SampleAOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _SampleAValidation(_ctx context.Context, _input *SampleAInput, _output *SampleAOutput) {
}
func _SampleARun(_ctx context.Context, input *SampleAInput) *SampleAOutput {
	output := &SampleAOutput{}
	_SampleASetup(_ctx, input)
	_SampleASteps(_ctx, input, output)
	_SampleAAnalysis(_ctx, input, output)
	_SampleAValidation(_ctx, input, output)
	return output
}

func SampleARunSteps(_ctx context.Context, input *SampleAInput) *SampleASOutput {
	soutput := &SampleASOutput{}
	output := _SampleARun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SampleANew() interface{} {
	return &SampleAElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SampleAInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SampleARun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SampleAInput{},
			Out: &SampleAOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type SampleAElement struct {
	inject.CheckedRunner
}

type SampleAInput struct {
	SampleVolume wunit.Volume
	Solution     *wtype.LHComponent
}

type SampleAOutput struct {
	SampleA *wtype.LHComponent
	Status  string
}

type SampleASOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
		SampleA *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "SampleA",
		Constructor: SampleANew,
		Desc: component.ComponentDesc{
			Desc: "Example protocol demonstrating the use of the Sample function\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson2_Sample/JAJALesson2/2A_Sample/A_Sample.an",
			Params: []component.ParamDesc{
				{Name: "SampleVolume", Desc: "antha, like golang is a strongly typed language in which the type of a variable must be declared.\nIn this case we're creating a variable called SampleVolume which is of type Volume;\nthe type system allows the antha compiler to catch many types of common errors before the programme\nis run the antha type system extends this to biological types such as volumes here.\nfunctions require inputs of particular types to be adhered to\n", Kind: "Parameters"},
				{Name: "Solution", Desc: "the LHComponent is the principal liquidhandling type in antha\nthe * signifies that this is a pointer to the component rather than the component itself\nmost key antha functions such as Sample and Mix use *LHComponent rather than LHComponent\nsince the type is imported from the wtype package we need to use  *wtype.LHComponent rather\nthan simply *LHComponent\n", Kind: "Inputs"},
				{Name: "SampleA", Desc: "An output LHComponent variable is created called Sample\n", Kind: "Outputs"},
				{Name: "Status", Desc: "Antha inherits all standard primitives valid in golang;\nfor example the string type shown here used to return a textual message\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
