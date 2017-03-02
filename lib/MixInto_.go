// Aliquot a solution into a specified plate.
// optionally premix the solution before aliquoting
package lib

import (
	"context"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol (data)

// optional field. Select if the solution to be aliquoted should be premixed prior to transer
// optional field to change the name of the component

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// this time we're specifying what plate we're using

// Physical outputs from this protocol with types

func _MixIntoRequirements() {

}

// Conditions to run on startup
func _MixIntoSetup(_ctx context.Context, _input *MixIntoInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _MixIntoSteps(_ctx context.Context, _input *MixIntoInput, _output *MixIntoOutput) {

	number := _input.SolutionVolume.SIValue() / _input.VolumePerAliquot.SIValue()
	possiblenumberofAliquots, _ := wutil.RoundDown(number)
	if possiblenumberofAliquots < _input.NumberofAliquots {
		execute.Errorf(_ctx, "Not enough solution for this many aliquots")
	}

	// if PreMix is selected change liquid type accordingly
	if _input.PreMix {
		_input.Solution.Type = wtype.LTPreMix
	}

	// if a solution name is given change the name
	if _input.ChangeSolutionName != "" {
		_input.Solution.CName = _input.ChangeSolutionName
	}

	aliquots := make([]*wtype.LHComponent, 0)

	for i := 0; i < _input.NumberofAliquots; i++ {
		if _input.Solution.TypeName() == "dna" {
			_input.Solution.Type = wtype.LTDoNotMix
		}
		aliquotSample := mixer.Sample(_input.Solution, _input.VolumePerAliquot)

		// the MixInto command is used instead of Mix to specify the plate
		// MixInto allows you to specify the exact plate to MixInto (i.e. rather than just a plate type. e.g. barcode 123214234)
		// the three input fields to the MixInto command represent
		// 1. the plate
		// 2. well location as a  string e.g. "A1" (in this case leaving it blank "" will leave the well location up to the scheduler),
		// 3. the sample or array of samples to be mixed
		var aliquot *wtype.LHComponent
		if _input.OptimisePlateUsage {
			aliquot = execute.MixNamed(_ctx, _input.OutPlate.Type, "", "AliquotPlate", aliquotSample)
		} else {
			aliquot = execute.MixInto(_ctx, _input.OutPlate, "", aliquotSample)
		}
		if aliquot != nil {
			aliquots = append(aliquots, aliquot)
		}
	}
	_output.Aliquots = aliquots
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MixIntoAnalysis(_ctx context.Context, _input *MixIntoInput, _output *MixIntoOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _MixIntoValidation(_ctx context.Context, _input *MixIntoInput, _output *MixIntoOutput) {

}
func _MixIntoRun(_ctx context.Context, input *MixIntoInput) *MixIntoOutput {
	output := &MixIntoOutput{}
	_MixIntoSetup(_ctx, input)
	_MixIntoSteps(_ctx, input, output)
	_MixIntoAnalysis(_ctx, input, output)
	_MixIntoValidation(_ctx, input, output)
	return output
}

func MixIntoRunSteps(_ctx context.Context, input *MixIntoInput) *MixIntoSOutput {
	soutput := &MixIntoSOutput{}
	output := _MixIntoRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MixIntoNew() interface{} {
	return &MixIntoElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MixIntoInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MixIntoRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MixIntoInput{},
			Out: &MixIntoOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type MixIntoElement struct {
	inject.CheckedRunner
}

type MixIntoInput struct {
	ChangeSolutionName string
	NumberofAliquots   int
	OptimisePlateUsage bool
	OutPlate           *wtype.LHPlate
	PreMix             bool
	Solution           *wtype.LHComponent
	SolutionVolume     wunit.Volume
	VolumePerAliquot   wunit.Volume
}

type MixIntoOutput struct {
	Aliquots []*wtype.LHComponent
}

type MixIntoSOutput struct {
	Data struct {
	}
	Outputs struct {
		Aliquots []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MixInto",
		Constructor: MixIntoNew,
		Desc: component.ComponentDesc{
			Desc: "Aliquot a solution into a specified plate.\noptionally premix the solution before aliquoting\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/AnthaLangAcademy/Lesson3_Mix_Loops/MixTestForSid/MixInto.an",
			Params: []component.ParamDesc{
				{Name: "ChangeSolutionName", Desc: "optional field to change the name of the component\n", Kind: "Parameters"},
				{Name: "NumberofAliquots", Desc: "", Kind: "Parameters"},
				{Name: "OptimisePlateUsage", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "this time we're specifying what plate we're using\n", Kind: "Inputs"},
				{Name: "PreMix", Desc: "optional field. Select if the solution to be aliquoted should be premixed prior to transer\n", Kind: "Parameters"},
				{Name: "Solution", Desc: "", Kind: "Inputs"},
				{Name: "SolutionVolume", Desc: "", Kind: "Parameters"},
				{Name: "VolumePerAliquot", Desc: "", Kind: "Parameters"},
				{Name: "Aliquots", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
