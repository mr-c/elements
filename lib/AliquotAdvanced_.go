// The Aliquot element will transfer a defined liquid at defined volumes a specified number of times into a chosen plate type.
// The user has the option to premix the solution to be aliquoted if the input solution tends to sediment or separate when left to stand
// (e.g. a suspension of cells in media) or has recently been thawed. Upstream elements that produce solutions as outputs can be wired
// into the Solution parameter of this element for aliquoting. If the solution already exists in your lab or has been made manually but a definition for this solution does
// not exist in the Antha library then the Add_Solution element can be used to define this solution with the output from the
// element wired into the Solution parameter of this element.
package lib

import (
	"context"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"strconv"
)

// Input parameters for this protocol (data)

// This parameter represents the volume of solution that you have in the lab available to be aliquoted. It does not represent the total volume to be aliquoted or the volume of liquid that will be used.

// This parameter dictates the final volume each aliquot will have.

// This parameter states the number of aliquots that will be made from the input Solution.

// This parameter states whether the aliquots should be made by row or column.

// This parameter sets the number of replica plates to perform aliquots to. Default number of plates is 1.

// This parameter is an optional field. If the solution to be aliquoted has components that may sink to the bottom of the solution then select this option for the solution to be premixed prior to transfer.

// This parameter is an optional field. If you want to change the name of the input Solution for traceability then do so. If not the default name will be given as the chosen input Solution LHComponent name.

// This parameter is an optional field. It states the number of wells that have already been used in the output plate and will start making aliquots from this position onwards. If there is more than one replica plate all plates would have the same number of wells already used.

// This name will be used as the identifier for the specific plate you are making your aliquots to appended with the replical plate number. If left blank it will default to AliquotPlate. If running more than one instance of the Aliquot element in parallel and you want the dilutions to all be made on the same plate be sure this name is the same across all instance parameter sets, and be sure to wire WellsAlreadyUsed into WellsUsed between the different instances.

// Data which is returned from this protocol, and data types

// This data output is a count of how many wells have been used in the output plate.

// Physical Inputs to this protocol with types

// This Physical input will have associated properties to determine how the liquid should be handled, e.g. is your Solution water or is it Glycerol. If your physical liquid does not exist in the Antha LHComponent library then create a new one on the fly with the Add_Solution element and wire the output into this input. Alternatively wire a solution made by another element into this input to be aliquoted.

// This parameter allows you to specify the type of plate you are aliquoting your Solution into. Choose from one of the available plate options from the Antha plate library.

// Physical outputs from this protocol with types

// This is a list of the resulting aliquots that have been made by the element.

func _AliquotAdvancedRequirements() {

}

// Conditions to run on startup
func _AliquotAdvancedSetup(_ctx context.Context, _input *AliquotAdvancedInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _AliquotAdvancedSteps(_ctx context.Context, _input *AliquotAdvancedInput, _output *AliquotAdvancedOutput) {

	// We need to make sure that we have enough solution after subtracting the residual volume of solution left in the input plate.
	// In future this will be calculated explicitly, but here we are estimating it as 10% extra for simplicity.
	residualVolAllowance := 0.10
	residualVol := wunit.MultiplyVolume(_input.SolutionVolume, residualVolAllowance)

	// Calculate the volume needed based on the number of aliquots, number of replica plates, and aliquot amount specified. This is only used for error messages.
	minVolume := wunit.MultiplyVolume(_input.VolumePerAliquot, float64(_input.NumberofAliquots*_input.NumberOfReplicaPlates))
	volumeNeeded := wunit.MultiplyVolume(minVolume, (1 / (1 - residualVolAllowance)))

	// Here we're doing some maths to work out what the possible number of aliquots is that we can make given the volume specified and the volume of solution we have.
	// We round this number down to the nearest number of aliquots.
	number := (_input.SolutionVolume.SIValue() - residualVol.SIValue()) / _input.VolumePerAliquot.SIValue()
	possiblenumberofAliquots, _ := wutil.RoundDown(number)

	// The total number of aliquots to be made is the number specified by the user for each of the Replica Plates being made.
	if possiblenumberofAliquots < (_input.NumberofAliquots * _input.NumberOfReplicaPlates) {
		execute.Errorf(_ctx, "Not enough solution for this many aliquots. You have specified %s, but %s is required based on the parameters you have specified and a 10 percent allowance for residual volume left in the input plate.", _input.SolutionVolume.ToString(), volumeNeeded.ToString())
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

	// This code allows the user to specify how the aliquots should be made, by row or by column.
	allwellpositions := _input.OutPlate.AllWellPositions(_input.ByRow)

	// Make a slice to contain all the liquid handling instructions for all of the aliquots.
	aliquots := make([]*wtype.LHComponent, 0)

	// This code checks to make sure the number of replica plates is greater than 0.
	if _input.NumberOfReplicaPlates < 1 {
		execute.Errorf(_ctx, "Number of replica plates must be greater than 0")
	}

	// create a counter variable to be able to keep track of the wells used each iteration of the aliquot loop.
	// var counter int

	// If wells have already been used and wired to this element from a previous element append them to the WellsUsed slice
	_output.WellsUsed = append(_output.WellsUsed, _input.WellsAlreadyUsed...)

	// This loop allows the user to specify the number of replica plates of aliquots they want.
	for platenumber := 1; platenumber < (_input.NumberOfReplicaPlates + 1); platenumber++ {

		// This loop cycles through the number of aliquots required
		for i := 0; i < _input.NumberofAliquots; i++ {

			// Here we are setting the AliquotPlateName to default to AliquotPlate if left blank by the user
			if _input.AliquotPlateName == "" {
				_input.AliquotPlateName = "AliquotPlate"
			}

			// Here we are renaming the Aliquot plate to account for replica plates. e.g. if there are no replicas then plate name will be set to "AliquotPlate" if there are 2 replicas then plate name will be set to "AliquotPlate1" and "AliquotPlate2", etc.
			platename := (_input.AliquotPlateName + strconv.Itoa(platenumber))

			// This statement changes the liquid handling policy if the solution being aliquoted is DNA to avoid cross contamination.
			if _input.Solution.TypeName() == "dna" {
				_input.Solution.Type = wtype.LTDoNotMix
			}

			// This line is using the Sample function from the mixer library to create a liquid handling component from a specified physical solution and the user defined aliquot volume
			aliquotSample := mixer.Sample(_input.Solution, _input.VolumePerAliquot)

			// Create a variable to contain all the liquid handling instructions for each aliquot
			var aliquot *wtype.LHComponent

			well, _ := findNextFreeWell(allwellpositions, _output.WellsUsed)

			// The MixNamed command here cycles through the well positions of the chosen plate type and plate name for each aliquot.
			// the MixNamed command is used instead of Mix to specify the plate type (e.g. "greiner384" or "pcrplate_skirted")
			// the four input fields to the MixNamed command represent
			// 1. the platetype as a string: commonly the input to the Antha element will actually be an LHPlate rather than a string so the type field can be accessed with OutPlate.Type
			// 2. well location as a  string e.g. "A1" (in this instance determined by a counter and the plate type or leaving it blank "" will leave the well location up to the scheduler),
			// 3. the plate name as a string,
			// 4. the sample or array of samples to be mixed; in the case of an array you'd normally feed this in as samples...
			aliquot = execute.MixNamed(_ctx, _input.OutPlate.Type, well, platename, aliquotSample)

			// store all of the used wells in the WellsUsed slice.
			_output.WellsUsed = append(_output.WellsUsed, well)

			// If the liquid handling instructions for a single aliquot were made with no errors then add to the aliquots slice.
			if aliquot != nil {
				aliquots = append(aliquots, aliquot)
			}

		}

	}
	_output.Aliquots = aliquots
}

// Check the wells that have already been used in this plate to be sure not to pipette on top of a previous solution. This is predominantly used for when two aliquot elements are used in serial to make sure the second set of aliquots are not pipette onto the first set.
func findNextFreeWell(allWellPositions []string, wellsUsed []string) (well string, err error) {

	for _, well := range allWellPositions {
		// If a well position is found to already have been used then add one to our counter that specifies the next well to use. See step 2 of the following comments.
		if !search.InSlice(well, wellsUsed) {
			return well, nil
		}
	}
	return
}

/*
func inSlice(a string, bs []string)bool{
	for _,b := range bs {
		if b == a {
			return true
		}
	}
	return false
}
*/

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AliquotAdvancedAnalysis(_ctx context.Context, _input *AliquotAdvancedInput, _output *AliquotAdvancedOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _AliquotAdvancedValidation(_ctx context.Context, _input *AliquotAdvancedInput, _output *AliquotAdvancedOutput) {

}
func _AliquotAdvancedRun(_ctx context.Context, input *AliquotAdvancedInput) *AliquotAdvancedOutput {
	output := &AliquotAdvancedOutput{}
	_AliquotAdvancedSetup(_ctx, input)
	_AliquotAdvancedSteps(_ctx, input, output)
	_AliquotAdvancedAnalysis(_ctx, input, output)
	_AliquotAdvancedValidation(_ctx, input, output)
	return output
}

func AliquotAdvancedRunSteps(_ctx context.Context, input *AliquotAdvancedInput) *AliquotAdvancedSOutput {
	soutput := &AliquotAdvancedSOutput{}
	output := _AliquotAdvancedRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AliquotAdvancedNew() interface{} {
	return &AliquotAdvancedElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AliquotAdvancedInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AliquotAdvancedRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AliquotAdvancedInput{},
			Out: &AliquotAdvancedOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type AliquotAdvancedElement struct {
	inject.CheckedRunner
}

type AliquotAdvancedInput struct {
	AliquotPlateName      string
	ByRow                 bool
	ChangeSolutionName    string
	NumberOfReplicaPlates int
	NumberofAliquots      int
	OutPlate              *wtype.LHPlate
	PreMix                bool
	Solution              *wtype.LHComponent
	SolutionVolume        wunit.Volume
	VolumePerAliquot      wunit.Volume
	WellsAlreadyUsed      []string
}

type AliquotAdvancedOutput struct {
	Aliquots  []*wtype.LHComponent
	WellsUsed []string
}

type AliquotAdvancedSOutput struct {
	Data struct {
		WellsUsed []string
	}
	Outputs struct {
		Aliquots []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AliquotAdvanced",
		Constructor: AliquotAdvancedNew,
		Desc: component.ComponentDesc{
			Desc: "The Aliquot element will transfer a defined liquid at defined volumes a specified number of times into a chosen plate type.\nThe user has the option to premix the solution to be aliquoted if the input solution tends to sediment or separate when left to stand\n(e.g. a suspension of cells in media) or has recently been thawed. Upstream elements that produce solutions as outputs can be wired\ninto the Solution parameter of this element for aliquoting. If the solution already exists in your lab or has been made manually but a definition for this solution does\nnot exist in the Antha library then the Add_Solution element can be used to define this solution with the output from the\nelement wired into the Solution parameter of this element.\n",
			Path: "src/github.com/antha-lang/elements/starter/Aliquot/AliquotAdvanced.an",
			Params: []component.ParamDesc{
				{Name: "AliquotPlateName", Desc: "This name will be used as the identifier for the specific plate you are making your aliquots to appended with the replical plate number. If left blank it will default to AliquotPlate. If running more than one instance of the Aliquot element in parallel and you want the dilutions to all be made on the same plate be sure this name is the same across all instance parameter sets, and be sure to wire WellsAlreadyUsed into WellsUsed between the different instances.\n", Kind: "Parameters"},
				{Name: "ByRow", Desc: "This parameter states whether the aliquots should be made by row or column.\n", Kind: "Parameters"},
				{Name: "ChangeSolutionName", Desc: "This parameter is an optional field. If you want to change the name of the input Solution for traceability then do so. If not the default name will be given as the chosen input Solution LHComponent name.\n", Kind: "Parameters"},
				{Name: "NumberOfReplicaPlates", Desc: "This parameter sets the number of replica plates to perform aliquots to. Default number of plates is 1.\n", Kind: "Parameters"},
				{Name: "NumberofAliquots", Desc: "This parameter states the number of aliquots that will be made from the input Solution.\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "This parameter allows you to specify the type of plate you are aliquoting your Solution into. Choose from one of the available plate options from the Antha plate library.\n", Kind: "Inputs"},
				{Name: "PreMix", Desc: "This parameter is an optional field. If the solution to be aliquoted has components that may sink to the bottom of the solution then select this option for the solution to be premixed prior to transfer.\n", Kind: "Parameters"},
				{Name: "Solution", Desc: "This Physical input will have associated properties to determine how the liquid should be handled, e.g. is your Solution water or is it Glycerol. If your physical liquid does not exist in the Antha LHComponent library then create a new one on the fly with the Add_Solution element and wire the output into this input. Alternatively wire a solution made by another element into this input to be aliquoted.\n", Kind: "Inputs"},
				{Name: "SolutionVolume", Desc: "This parameter represents the volume of solution that you have in the lab available to be aliquoted. It does not represent the total volume to be aliquoted or the volume of liquid that will be used.\n", Kind: "Parameters"},
				{Name: "VolumePerAliquot", Desc: "This parameter dictates the final volume each aliquot will have.\n", Kind: "Parameters"},
				{Name: "WellsAlreadyUsed", Desc: "This parameter is an optional field. It states the number of wells that have already been used in the output plate and will start making aliquots from this position onwards. If there is more than one replica plate all plates would have the same number of wells already used.\n", Kind: "Parameters"},
				{Name: "Aliquots", Desc: "This is a list of the resulting aliquots that have been made by the element.\n", Kind: "Outputs"},
				{Name: "WellsUsed", Desc: "This data output is a count of how many wells have been used in the output plate.\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
