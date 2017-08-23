// Protocol Transfer_multi will transfer a series of components to a specified destination plate
// Liquid names may be specified which will override the component names
// if no liquid names are specified the Starting Solution names are preserved
// either all liquidnames must be specified or none
package lib

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol (data)

// Liquid names may be specified which will override the component names
// if no liquid names are specified the Starting Solution names are preserved
// either all liquidnames must be specified or none

// One liquid volume is specified for all transfers

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// The solutions to be transferred

// one outplate is specified for all transfers

// Physical outputs from this protocol with types

func _Transfer_multiRequirements() {

}

// Conditions to run on startup
func _Transfer_multiSetup(_ctx context.Context, _input *Transfer_multiInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _Transfer_multiSteps(_ctx context.Context, _input *Transfer_multiInput, _output *Transfer_multiOutput) {

	_output.Status = make(map[string]string)

	if len(_input.Liquidnames) != len(_input.Startingsolutions) && len(_input.Liquidnames) > 0 {
		execute.Errorf(_ctx, "Unequal length of liquid names specified compared to starting solutions, either make these the same length or keep Liquidnames empty")
	}

	var outPlateWells []string = _input.OutPlate.AllWellPositions(wtype.BYCOLUMN)

	for i, liquid := range _input.Startingsolutions {

		if len(_input.Liquidnames) != 0 {
			liquid.CName = _input.Liquidnames[i]
		}

		sample := mixer.Sample(liquid, _input.LiquidVolume)

		_output.FinalSolutions = append(_output.FinalSolutions, execute.MixInto(_ctx, _input.OutPlate, outPlateWells[i], sample))

		// if a liquid name has already been specified the name will be appended
		if _, found := _output.Status[liquid.CName]; found {
			originalName := liquid.CName
			var number int = 2
			for {
				liquid.CName = originalName + fmt.Sprint(number)
				if _, found := _output.Status[liquid.CName]; found {
					number++
					liquid.CName = originalName + fmt.Sprint(number)
				} else {
					break
				}
			}
		}
		_output.Status[liquid.CName] = _input.LiquidVolume.ToString() + " of " + liquid.CName + " was mixed into " + _input.OutPlate.Type
	}
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Transfer_multiAnalysis(_ctx context.Context, _input *Transfer_multiInput, _output *Transfer_multiOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _Transfer_multiValidation(_ctx context.Context, _input *Transfer_multiInput, _output *Transfer_multiOutput) {

}
func _Transfer_multiRun(_ctx context.Context, input *Transfer_multiInput) *Transfer_multiOutput {
	output := &Transfer_multiOutput{}
	_Transfer_multiSetup(_ctx, input)
	_Transfer_multiSteps(_ctx, input, output)
	_Transfer_multiAnalysis(_ctx, input, output)
	_Transfer_multiValidation(_ctx, input, output)
	return output
}

func Transfer_multiRunSteps(_ctx context.Context, input *Transfer_multiInput) *Transfer_multiSOutput {
	soutput := &Transfer_multiSOutput{}
	output := _Transfer_multiRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Transfer_multiNew() interface{} {
	return &Transfer_multiElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Transfer_multiInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Transfer_multiRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Transfer_multiInput{},
			Out: &Transfer_multiOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type Transfer_multiElement struct {
	inject.CheckedRunner
}

type Transfer_multiInput struct {
	LiquidVolume      wunit.Volume
	Liquidnames       []string
	OutPlate          *wtype.LHPlate
	Startingsolutions []*wtype.LHComponent
}

type Transfer_multiOutput struct {
	FinalSolutions []*wtype.LHComponent
	Status         map[string]string
}

type Transfer_multiSOutput struct {
	Data struct {
		Status map[string]string
	}
	Outputs struct {
		FinalSolutions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Transfer_multi",
		Constructor: Transfer_multiNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol Transfer_multi will transfer a series of components to a specified destination plate\nLiquid names may be specified which will override the component names\nif no liquid names are specified the Starting Solution names are preserved\neither all liquidnames must be specified or none\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/Transfer/Transfer_multi.an",
			Params: []component.ParamDesc{
				{Name: "LiquidVolume", Desc: "One liquid volume is specified for all transfers\n", Kind: "Parameters"},
				{Name: "Liquidnames", Desc: "Liquid names may be specified which will override the component names\nif no liquid names are specified the Starting Solution names are preserved\neither all liquidnames must be specified or none\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "one outplate is specified for all transfers\n", Kind: "Inputs"},
				{Name: "Startingsolutions", Desc: "The solutions to be transferred\n", Kind: "Inputs"},
				{Name: "FinalSolutions", Desc: "", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
