// The Aliquot element will transfer a defined liquid at defined volumes a specified number of times into a chosen plate type.
// The user has the option to premix the solution to be aliquoted if the input solution tends to sediment or separate when left to stand
// (e.g. a suspension of cells in media) or has recently been thawed. Upstream elements that produce solutions as outputs can be wired
// into the Solution parameter of this element for aliquoting. If the solution already exists in your lab or has been made manually but a definition for this solution does
// not exist in the Antha library then the Add_Solution element can be used to define this solution with the output from the
// element wired into the Solution parameter of this element.
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

// This parameter represents the volume of solution that you have in the lab available to be aliquoted. It does not represent the total volume to be aliquoted or the volume of liquid that will be used.

// This parameter dictates the final volume each aliquot will have.

// This parameter states the number of aliquots that will be made from the input Solution.

// This parameter is an optional field. If the solution to be aliquoted has components that may sink to the bottom of the solution then select this option for the solution to be premixed prior to transfer.

// This parameter is an optional field. If you want to change the name of the input Solution for traceability then do so. If not the default name will be given as the chosen input Solution LHComponent name.

// This parameter is an optional field. If set to true then the aliquots will be transferred to a specific named plate such that if two instances of this element are run in parallel the aliquots from both will be put on the same output plate rather than two separate output plates.

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// This Physical input will have associated properties to determine how the liquid should be handled, e.g. is your Solution water or is it Glycerol. If your physical liquid does not exist in the Antha LHComponent library then create a new one on the fly with the Add_Solution element and wire the output into this input. Alternatively wire a solution made by another element into this input to be alliquoted.

// This parameter alows you to specify the type of plate you are aliquoting your Solution into. Choose from one of the available plate options from the Antha plate library.

// Physical outputs from this protocol with types

// This is a list of the resulting aliquots that have been made by the element.

func _AliquotRequirements() {

}

// Conditions to run on startup
func _AliquotSetup(_ctx context.Context, _input *AliquotInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _AliquotSteps(_ctx context.Context, _input *AliquotInput, _output *AliquotOutput) {

	number := _input.SolutionVolume.SIValue() / _input.VolumePerAliquot.SIValue()
	possiblenumberofAliquots, _ := wutil.RoundDown(number)
	if possiblenumberofAliquots < _input.NumberofAliquots {
		execute.Errorf(_ctx, "Not enough solution for this many aliquots")
	}

	//check if maxvolume of outplate is higher than specified aliquot volume
	if _input.OutPlate.Welltype.MaxVolume().LessThanRounded(_input.VolumePerAliquot, 5) {
		execute.Errorf(_ctx, "Aliquot volume specified (%s) too high for well capacity (%s) of current plate (%s)", _input.VolumePerAliquot.ToString(), _input.OutPlate.Welltype.MaxVolume(), _input.OutPlate.Name())
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
func _AliquotAnalysis(_ctx context.Context, _input *AliquotInput, _output *AliquotOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _AliquotValidation(_ctx context.Context, _input *AliquotInput, _output *AliquotOutput) {

}
func _AliquotRun(_ctx context.Context, input *AliquotInput) *AliquotOutput {
	output := &AliquotOutput{}
	_AliquotSetup(_ctx, input)
	_AliquotSteps(_ctx, input, output)
	_AliquotAnalysis(_ctx, input, output)
	_AliquotValidation(_ctx, input, output)
	return output
}

func AliquotRunSteps(_ctx context.Context, input *AliquotInput) *AliquotSOutput {
	soutput := &AliquotSOutput{}
	output := _AliquotRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AliquotNew() interface{} {
	return &AliquotElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AliquotInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AliquotRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AliquotInput{},
			Out: &AliquotOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type AliquotElement struct {
	inject.CheckedRunner
}

type AliquotInput struct {
	ChangeSolutionName string
	NumberofAliquots   int
	OptimisePlateUsage bool
	OutPlate           *wtype.LHPlate
	PreMix             bool
	Solution           *wtype.LHComponent
	SolutionVolume     wunit.Volume
	VolumePerAliquot   wunit.Volume
}

type AliquotOutput struct {
	Aliquots []*wtype.LHComponent
}

type AliquotSOutput struct {
	Data struct {
	}
	Outputs struct {
		Aliquots []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Aliquot",
		Constructor: AliquotNew,
		Desc: component.ComponentDesc{
			Desc: "The Aliquot element will transfer a defined liquid at defined volumes a specified number of times into a chosen plate type.\nThe user has the option to premix the solution to be aliquoted if the input solution tends to sediment or separate when left to stand\n(e.g. a suspension of cells in media) or has recently been thawed. Upstream elements that produce solutions as outputs can be wired\ninto the Solution parameter of this element for aliquoting. If the solution already exists in your lab or has been made manually but a definition for this solution does\nnot exist in the Antha library then the Add_Solution element can be used to define this solution with the output from the\nelement wired into the Solution parameter of this element.\n",
			Path: "src/github.com/antha-lang/elements/starter/Aliquot/Aliquot.an",
			Params: []component.ParamDesc{
				{Name: "ChangeSolutionName", Desc: "This parameter is an optional field. If you want to change the name of the input Solution for traceability then do so. If not the default name will be given as the chosen input Solution LHComponent name.\n", Kind: "Parameters"},
				{Name: "NumberofAliquots", Desc: "This parameter states the number of aliquots that will be made from the input Solution.\n", Kind: "Parameters"},
				{Name: "OptimisePlateUsage", Desc: "This parameter is an optional field. If set to true then the aliquots will be transferred to a specific named plate such that if two instances of this element are run in parallel the aliquots from both will be put on the same output plate rather than two separate output plates.\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "This parameter alows you to specify the type of plate you are aliquoting your Solution into. Choose from one of the available plate options from the Antha plate library.\n", Kind: "Inputs"},
				{Name: "PreMix", Desc: "This parameter is an optional field. If the solution to be aliquoted has components that may sink to the bottom of the solution then select this option for the solution to be premixed prior to transfer.\n", Kind: "Parameters"},
				{Name: "Solution", Desc: "This Physical input will have associated properties to determine how the liquid should be handled, e.g. is your Solution water or is it Glycerol. If your physical liquid does not exist in the Antha LHComponent library then create a new one on the fly with the Add_Solution element and wire the output into this input. Alternatively wire a solution made by another element into this input to be alliquoted.\n", Kind: "Inputs"},
				{Name: "SolutionVolume", Desc: "This parameter represents the volume of solution that you have in the lab available to be aliquoted. It does not represent the total volume to be aliquoted or the volume of liquid that will be used.\n", Kind: "Parameters"},
				{Name: "VolumePerAliquot", Desc: "This parameter dictates the final volume each aliquot will have.\n", Kind: "Parameters"},
				{Name: "Aliquots", Desc: "This is a list of the resulting aliquots that have been made by the element.\n", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
