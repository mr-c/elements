// this protocol will set up a specified number of reactions one component at a time, i.e. in the following format:
// add component 1 into reaction 1 location,
// add component 1 into reaction 2 location,
// add component 1 into reaction n location,
// add component 2 into reaction 1 location,
// add component 2 into reaction 2 location,
// add component 2 into reaction n location,
// add component x into reaction 1 location,
// add component x into reaction 2 location,
// add component x into reaction n location,
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

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _ReactionSetUpRequirements() {
}

// Conditions to run on startup
func _ReactionSetUpSetup(_ctx context.Context, _input *ReactionSetUpInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _ReactionSetUpSteps(_ctx context.Context, _input *ReactionSetUpInput, _output *ReactionSetUpOutput) {

	reactions := make([]*wtype.LHComponent, 0)

	for i := 0; i < _input.NumberofReactions; i++ {

		bufferSample := mixer.SampleForTotalVolume(_input.Buffer, _input.TotalVolume)

		buffer := execute.MixInto(_ctx, _input.OutPlate, "", bufferSample)

		subSample := mixer.Sample(_input.Substrate, _input.SubstrateVolume)

		// this will Mix subSample into buffer
		subplusbuffer := execute.Mix(_ctx, buffer, subSample)

		enzSample := mixer.Sample(_input.Enzyme, _input.EnzymeVolume)

		// by separating each reaction into 3 mix commands when the reactions are pipetted out in parallel at runtime
		// the scheduler will try to optimise each mix step and will hence look for where components are shared between
		// parallel threads to take advantage of things like tip reuse and multipipetting (if the liquid class permits these)
		reaction := execute.Mix(_ctx, subplusbuffer, enzSample)

		reactions = append(reactions, reaction)

	}
	_output.Reactions = reactions

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _ReactionSetUpAnalysis(_ctx context.Context, _input *ReactionSetUpInput, _output *ReactionSetUpOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _ReactionSetUpValidation(_ctx context.Context, _input *ReactionSetUpInput, _output *ReactionSetUpOutput) {
}
func _ReactionSetUpRun(_ctx context.Context, input *ReactionSetUpInput) *ReactionSetUpOutput {
	output := &ReactionSetUpOutput{}
	_ReactionSetUpSetup(_ctx, input)
	_ReactionSetUpSteps(_ctx, input, output)
	_ReactionSetUpAnalysis(_ctx, input, output)
	_ReactionSetUpValidation(_ctx, input, output)
	return output
}

func ReactionSetUpRunSteps(_ctx context.Context, input *ReactionSetUpInput) *ReactionSetUpSOutput {
	soutput := &ReactionSetUpSOutput{}
	output := _ReactionSetUpRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ReactionSetUpNew() interface{} {
	return &ReactionSetUpElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ReactionSetUpInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ReactionSetUpRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ReactionSetUpInput{},
			Out: &ReactionSetUpOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type ReactionSetUpElement struct {
	inject.CheckedRunner
}

type ReactionSetUpInput struct {
	Buffer            *wtype.LHComponent
	Enzyme            *wtype.LHComponent
	EnzymeVolume      wunit.Volume
	NumberofReactions int
	OutPlate          *wtype.LHPlate
	Substrate         *wtype.LHComponent
	SubstrateVolume   wunit.Volume
	TotalVolume       wunit.Volume
}

type ReactionSetUpOutput struct {
	Reactions []*wtype.LHComponent
	Status    string
}

type ReactionSetUpSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
		Reactions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "ReactionSetUp",
		Constructor: ReactionSetUpNew,
		Desc: component.ComponentDesc{
			Desc: "this protocol will set up a specified number of reactions one component at a time, i.e. in the following format:\nadd component 1 into reaction 1 location,\nadd component 1 into reaction 2 location,\nadd component 1 into reaction n location,\nadd component 2 into reaction 1 location,\nadd component 2 into reaction 2 location,\nadd component 2 into reaction n location,\nadd component x into reaction 1 location,\nadd component x into reaction 2 location,\nadd component x into reaction n location,\n",
			Path: "src/github.com/antha-lang/elements/an/AnthaAcademy/Lesson3_MixPart2/A_Assaysetup_componentbycomponent.an",
			Params: []component.ParamDesc{
				{Name: "Buffer", Desc: "", Kind: "Inputs"},
				{Name: "Enzyme", Desc: "", Kind: "Inputs"},
				{Name: "EnzymeVolume", Desc: "", Kind: "Parameters"},
				{Name: "NumberofReactions", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Substrate", Desc: "", Kind: "Inputs"},
				{Name: "SubstrateVolume", Desc: "", Kind: "Parameters"},
				{Name: "TotalVolume", Desc: "", Kind: "Parameters"},
				{Name: "Reactions", Desc: "", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
