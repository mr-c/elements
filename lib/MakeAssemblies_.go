// Assemble multiple assemblies using TypeIIs construct assembly
package lib

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Reaction volume
// Volumes corresponding to input parts
// Names corresonding to input parts
// Vector volume
// Buffer volume
// ATP volume
// Restriction enzyme volume
// Ligase volume

// Reaction temperature
// Reaction time
// Inactivation temperature
// Inactivation time

// Prefix for reaction names

// Input parts, one per assembly
// Vector to use
// Restriction enzyme to use
// Buffer to use
// Water to use
// Ligase to use
// ATP to use
// Output plate

// List of assembled parts

func _MakeAssembliesSetup(_ctx context.Context, _input *MakeAssembliesInput) {}

func _MakeAssembliesSteps(_ctx context.Context, _input *MakeAssembliesInput, _output *MakeAssembliesOutput) {
	for k := range _input.Parts {
		result := TypeIISConstructAssemblyRunSteps(_ctx, &TypeIISConstructAssemblyInput{ReactionVolume: _input.ReactionVolume,
			PartVols:           _input.PartVols[k],
			PartNames:          _input.PartNames[k],
			VectorVol:          _input.VectorVol,
			BufferVol:          _input.BufferVol,
			AtpVol:             _input.AtpVol,
			ReVol:              _input.ReVol,
			LigVol:             _input.LigVol,
			ReactionTemp:       _input.ReactionTemp,
			ReactionTime:       _input.ReactionTime,
			OutputReactionName: fmt.Sprintf("%s%d", _input.OutputReactionName, k),

			Parts:             _input.Parts[k],
			Vector:            _input.Vector,
			RestrictionEnzyme: _input.RestrictionEnzyme,
			Buffer:            _input.Buffer,
			Water:             _input.Water,
			Ligase:            _input.Ligase,
			Atp:               _input.Atp,
			OutPlate:          _input.OutPlate},
		)
		_output.Reactions = append(_output.Reactions, result.Outputs.Reaction)
	}
}

func _MakeAssembliesAnalysis(_ctx context.Context, _input *MakeAssembliesInput, _output *MakeAssembliesOutput) {
}

func _MakeAssembliesValidation(_ctx context.Context, _input *MakeAssembliesInput, _output *MakeAssembliesOutput) {
}
func _MakeAssembliesRun(_ctx context.Context, input *MakeAssembliesInput) *MakeAssembliesOutput {
	output := &MakeAssembliesOutput{}
	_MakeAssembliesSetup(_ctx, input)
	_MakeAssembliesSteps(_ctx, input, output)
	_MakeAssembliesAnalysis(_ctx, input, output)
	_MakeAssembliesValidation(_ctx, input, output)
	return output
}

func MakeAssembliesRunSteps(_ctx context.Context, input *MakeAssembliesInput) *MakeAssembliesSOutput {
	soutput := &MakeAssembliesSOutput{}
	output := _MakeAssembliesRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MakeAssembliesNew() interface{} {
	return &MakeAssembliesElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MakeAssembliesInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MakeAssembliesRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MakeAssembliesInput{},
			Out: &MakeAssembliesOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type MakeAssembliesElement struct {
	inject.CheckedRunner
}

type MakeAssembliesInput struct {
	Atp                *wtype.LHComponent
	AtpVol             wunit.Volume
	Buffer             *wtype.LHComponent
	BufferVol          wunit.Volume
	InactivationTemp   wunit.Temperature
	InactivationTime   wunit.Time
	LigVol             wunit.Volume
	Ligase             *wtype.LHComponent
	OutPlate           *wtype.LHPlate
	OutputReactionName string
	PartNames          [][]string
	PartVols           [][]wunit.Volume
	Parts              [][]*wtype.LHComponent
	ReVol              wunit.Volume
	ReactionTemp       wunit.Temperature
	ReactionTime       wunit.Time
	ReactionVolume     wunit.Volume
	RestrictionEnzyme  *wtype.LHComponent
	Vector             *wtype.LHComponent
	VectorVol          wunit.Volume
	Water              *wtype.LHComponent
}

type MakeAssembliesOutput struct {
	Reactions []*wtype.LHComponent
}

type MakeAssembliesSOutput struct {
	Data struct {
	}
	Outputs struct {
		Reactions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MakeAssemblies",
		Constructor: MakeAssembliesNew,
		Desc: component.ComponentDesc{
			Desc: "Assemble multiple assemblies using TypeIIs construct assembly\n",
			Path: "src/github.com/antha-lang/elements/an/MakeAssemblies/element.an",
			Params: []component.ParamDesc{
				{Name: "Atp", Desc: "ATP to use\n", Kind: "Inputs"},
				{Name: "AtpVol", Desc: "ATP volume\n", Kind: "Parameters"},
				{Name: "Buffer", Desc: "Buffer to use\n", Kind: "Inputs"},
				{Name: "BufferVol", Desc: "Buffer volume\n", Kind: "Parameters"},
				{Name: "InactivationTemp", Desc: "Inactivation temperature\n", Kind: "Parameters"},
				{Name: "InactivationTime", Desc: "Inactivation time\n", Kind: "Parameters"},
				{Name: "LigVol", Desc: "Ligase volume\n", Kind: "Parameters"},
				{Name: "Ligase", Desc: "Ligase to use\n", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "Output plate\n", Kind: "Inputs"},
				{Name: "OutputReactionName", Desc: "Prefix for reaction names\n", Kind: "Parameters"},
				{Name: "PartNames", Desc: "Names corresonding to input parts\n", Kind: "Parameters"},
				{Name: "PartVols", Desc: "Volumes corresponding to input parts\n", Kind: "Parameters"},
				{Name: "Parts", Desc: "Input parts, one per assembly\n", Kind: "Inputs"},
				{Name: "ReVol", Desc: "Restriction enzyme volume\n", Kind: "Parameters"},
				{Name: "ReactionTemp", Desc: "Reaction temperature\n", Kind: "Parameters"},
				{Name: "ReactionTime", Desc: "Reaction time\n", Kind: "Parameters"},
				{Name: "ReactionVolume", Desc: "Reaction volume\n", Kind: "Parameters"},
				{Name: "RestrictionEnzyme", Desc: "Restriction enzyme to use\n", Kind: "Inputs"},
				{Name: "Vector", Desc: "Vector to use\n", Kind: "Inputs"},
				{Name: "VectorVol", Desc: "Vector volume\n", Kind: "Parameters"},
				{Name: "Water", Desc: "Water to use\n", Kind: "Inputs"},
				{Name: "Reactions", Desc: "List of assembled parts\n", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
