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

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

func _MultichannelTestProtocol2Requirements() {
}

// Conditions to run on startup
func _MultichannelTestProtocol2Setup(_ctx context.Context, _input *MultichannelTestProtocol2Input) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _MultichannelTestProtocol2Steps(_ctx context.Context, _input *MultichannelTestProtocol2Input, _output *MultichannelTestProtocol2Output) {
	for k := 0; k < len(_input.Vols); k++ {
		var s *wtype.LHComponent
		if !_input.Vols[k].IsZero() {
			s = mixer.Sample(_input.Parts[k], _input.Vols[k])
			execute.MixTo(_ctx, _input.OutputPlateType, _input.Wells[k], 1, s)
		}
	}
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MultichannelTestProtocol2Analysis(_ctx context.Context, _input *MultichannelTestProtocol2Input, _output *MultichannelTestProtocol2Output) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _MultichannelTestProtocol2Validation(_ctx context.Context, _input *MultichannelTestProtocol2Input, _output *MultichannelTestProtocol2Output) {
}
func _MultichannelTestProtocol2Run(_ctx context.Context, input *MultichannelTestProtocol2Input) *MultichannelTestProtocol2Output {
	output := &MultichannelTestProtocol2Output{}
	_MultichannelTestProtocol2Setup(_ctx, input)
	_MultichannelTestProtocol2Steps(_ctx, input, output)
	_MultichannelTestProtocol2Analysis(_ctx, input, output)
	_MultichannelTestProtocol2Validation(_ctx, input, output)
	return output
}

func MultichannelTestProtocol2RunSteps(_ctx context.Context, input *MultichannelTestProtocol2Input) *MultichannelTestProtocol2SOutput {
	soutput := &MultichannelTestProtocol2SOutput{}
	output := _MultichannelTestProtocol2Run(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MultichannelTestProtocol2New() interface{} {
	return &MultichannelTestProtocol2Element{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MultichannelTestProtocol2Input{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MultichannelTestProtocol2Run(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MultichannelTestProtocol2Input{},
			Out: &MultichannelTestProtocol2Output{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type MultichannelTestProtocol2Element struct {
	inject.CheckedRunner
}

type MultichannelTestProtocol2Input struct {
	OutputPlateType string
	Parts           []*wtype.LHComponent
	Vols            []wunit.Volume
	Wells           []string
}

type MultichannelTestProtocol2Output struct {
}

type MultichannelTestProtocol2SOutput struct {
	Data struct {
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MultichannelTestProtocol2",
		Constructor: MultichannelTestProtocol2New,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/elements/an/Test/MultichannelTestProtocol/MultichannelTestProtocol2.an",
			Params: []component.ParamDesc{
				{Name: "OutputPlateType", Desc: "", Kind: "Inputs"},
				{Name: "Parts", Desc: "", Kind: "Inputs"},
				{Name: "Vols", Desc: "", Kind: "Parameters"},
				{Name: "Wells", Desc: "", Kind: "Parameters"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
