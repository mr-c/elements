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
// functions require inputs of particular types to be adhered t

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

func _SampleExercise3Requirements() {

}

// Conditions to run on startup
func _SampleExercise3Setup(_ctx context.Context, _input *SampleExercise3Input) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _SampleExercise3Steps(_ctx context.Context, _input *SampleExercise3Input, _output *SampleExercise3Output) {

	// the Sample function is imported from the mixer library
	// in the mixer library the function signature can be found, here it is:
	// func Sample(l *wtype.LHComponent, v wunit.Volume) *wtype.LHComponent {
	// The function signature  shows that the function requires a *LHComponent and a Volume and returns an *LHComponent
	dilutedsampleA := make([]*wtype.LHComponent, 0)
	dilutedsampleB := make([]*wtype.LHComponent, 0)
	var solutiona *wtype.LHComponent
	var solutionb *wtype.LHComponent

	diluentsampleA := mixer.SampleForTotalVolume(_input.Diluent, _input.FinalVolume)
	diluentsampleB := mixer.SampleForTotalVolume(_input.Diluent, _input.FinalVolume)

	dilutedsampleA = append(dilutedsampleA, diluentsampleA)
	dilutedsampleB = append(dilutedsampleB, diluentsampleB)

	solutiona = mixer.Sample(_input.InputSolution, _input.SampleVolumeA)
	solutionb = mixer.Sample(_input.InputSolution, _input.SampleVolumeB)

	dilutedsampleA = append(dilutedsampleA, solutiona)
	dilutedsampleB = append(dilutedsampleB, solutionb)

	// The Sample function is not sufficient to generate liquid handling instructions alone,
	// We would need a Mix command to instruct where to put the sample

	// we can also create data outputs as a string like this
	_output.StatusA = _input.SampleVolumeA.ToString() + " of " + _input.InputSolution.CName + " sampled"
	_output.StatusB = _input.SampleVolumeB.ToString() + " of " + _input.InputSolution.CName + " sampled"

	// To maintain good practice in coding all variables within the steps section should be lower case
	// when that variable is to become an output (or is an input) the first letter is capitalised as shown below

	_output.DilutedSampleA = dilutedsampleA
	_output.DilutedSampleB = dilutedsampleB

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _SampleExercise3Analysis(_ctx context.Context, _input *SampleExercise3Input, _output *SampleExercise3Output) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _SampleExercise3Validation(_ctx context.Context, _input *SampleExercise3Input, _output *SampleExercise3Output) {

}
func _SampleExercise3Run(_ctx context.Context, input *SampleExercise3Input) *SampleExercise3Output {
	output := &SampleExercise3Output{}
	_SampleExercise3Setup(_ctx, input)
	_SampleExercise3Steps(_ctx, input, output)
	_SampleExercise3Analysis(_ctx, input, output)
	_SampleExercise3Validation(_ctx, input, output)
	return output
}

func SampleExercise3RunSteps(_ctx context.Context, input *SampleExercise3Input) *SampleExercise3SOutput {
	soutput := &SampleExercise3SOutput{}
	output := _SampleExercise3Run(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SampleExercise3New() interface{} {
	return &SampleExercise3Element{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SampleExercise3Input{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SampleExercise3Run(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SampleExercise3Input{},
			Out: &SampleExercise3Output{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type SampleExercise3Element struct {
	inject.CheckedRunner
}

type SampleExercise3Input struct {
	Diluent       *wtype.LHComponent
	FinalVolume   wunit.Volume
	InputSolution *wtype.LHComponent
	SampleVolumeA wunit.Volume
	SampleVolumeB wunit.Volume
}

type SampleExercise3Output struct {
	DilutedSampleA []*wtype.LHComponent
	DilutedSampleB []*wtype.LHComponent
	StatusA        string
	StatusB        string
}

type SampleExercise3SOutput struct {
	Data struct {
		StatusA string
		StatusB string
	}
	Outputs struct {
		DilutedSampleA []*wtype.LHComponent
		DilutedSampleB []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "SampleExercise3",
		Constructor: SampleExercise3New,
		Desc: component.ComponentDesc{
			Desc: "Example protocol demonstrating the use of the Sample function\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson2_Sample/JAJALesson2/SampleExercise3/SampleExercise3.an",
			Params: []component.ParamDesc{
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "FinalVolume", Desc: "antha, like golang is a strongly typed language in which the type of a variable must be declared.\nIn this case we're creating a variable called SampleVolume which is of type Volume;\nthe type system allows the antha compiler to catch many types of common errors before the programme is run\nthe antha type system extends this to biological types such as volumes here.\nfunctions require inputs of particular types to be adhered t\n", Kind: "Parameters"},
				{Name: "InputSolution", Desc: "the LHComponent is the principal liquidhandling type in antha\nthe * signifies that this is a pointer to the component rather than the component itself\nmost key antha functions such as Sample and Mix use *LHComponent rather than LHComponent\nsince the type is imported from the wtype package we need to use  *wtype.LHComponent rather than simply *LHComponent\n", Kind: "Inputs"},
				{Name: "SampleVolumeA", Desc: "", Kind: "Parameters"},
				{Name: "SampleVolumeB", Desc: "", Kind: "Parameters"},
				{Name: "DilutedSampleA", Desc: "An output LHComponent variable is created called Sample\n", Kind: "Outputs"},
				{Name: "DilutedSampleB", Desc: "", Kind: "Outputs"},
				{Name: "StatusA", Desc: "Antha inherits all standard primitives valid in golang;\nfor example the string type shown here used to return a textual message\n", Kind: "Data"},
				{Name: "StatusB", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
