// This protocol will load DNA samples on an E-GEL for DNA analysis, and can take input from other protocols which exports an array of LHComponents.
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

//optionally define name of experiment
//select true if DNA samples already contain Loading Dye
//define number of technical replicates

//default is Load policy but can be overriden by specifying here

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// This is the input that comes as the output from an upstream element (AutoAssembly.an)

//WaterSolution //Chemspiderlink // not correct link but similar desirable
// gel
// plate to mix samples if required

// Physical outputs from this protocol with types

//Gel

// No special requirements on inputs
func _DNA_gelRequirements() {

}

// Condititions run on startup
// Including configuring an controls required, and the blocking level needed
// for them (in this case, per plate of samples processed)
func _DNA_gelSetup(_ctx context.Context, _input *DNA_gelInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _DNA_gelSteps(_ctx context.Context, _input *DNA_gelInput, _output *DNA_gelOutput) {

	//set up some arrays to fill and LHComponent variables for the DNA samples
	var loadedsamples []*wtype.LHComponent

	//setup variable for error reporting
	var err error

	//specify default mixing policy
	if _input.MixingPolicy == "" {
		_input.MixingPolicy = "Load"
	}

	//get well positions of DNA Gel from plate library ensuring the list is by row rather than by column
	var wells []string = _input.DNAGel.AllWellPositions(wtype.BYROW)

	var DNAGelloadmix *wtype.LHComponent
	var loadedsample *wtype.LHComponent

	//begin counter at first well position as E-GEL must be run upside down
	var counter int = len(wells) - 1

	//work out and copy sample volume
	samplevolume := (wunit.CopyVolume(_input.DNAGelRunVolume))

	//work out volume of water required
	samplevolume.Subtract(_input.WaterVol)

	//assign water to specific liquid handling load type
	_input.Water.Type = wtype.LTloadwater

	for j := 0; j < _input.Replicates; j++ {

		for i := range _input.Reactions {

			//sample water at specified water volume
			waterSample := mixer.Sample(_input.Water, _input.WaterVol)

			//update position to correspond to counter
			position := wells[counter]

			//get well coordinates from correct position
			wellcoords := wtype.MakeWellCoordsA1(position)

			//if it is the last column, add a ladder sample
			fmt.Println(_input.DNAGel.WlsX, wellcoords.X)
			if wellcoords.X == _input.DNAGel.WlsX-1 {

				_input.Ladder.Type, err = wtype.LiquidTypeFromString(_input.MixingPolicy)
				if err != nil {
					execute.Errorf(_ctx, "Error in specifying MixingPolicy %s for DNA Gel: %s", _input.MixingPolicy, err.Error())
				}

				//perform liquid handling for addiiton of ladder sample
				laddersample := execute.MixInto(_ctx, _input.DNAGel,
					position,
					mixer.SampleForTotalVolume(_input.Water, _input.DNAGelRunVolume),
					mixer.Sample(_input.Ladder, _input.LadderVolume),
				)

				//add ladder to array of loaded samples
				loadedsamples = append(loadedsamples, laddersample)

				//decrease counter by 1, as pipetting Gel backwards
				counter--

			}

			// refresh position in case ladder was added
			position = wells[counter]

			sampletotest := _input.Reactions[i]

			// load sample

			// add loading dye if necessary
			if !_input.LoadingDyeInSample {

				_input.LoadingDye.Type, err = wtype.LiquidTypeFromString("NeedToMix")
				if err != nil {
					execute.Errorf(_ctx, err.Error())
				}

				DNAGelloadmixsolution := execute.MixInto(_ctx, _input.MixPlate,
					"",
					mixer.Sample(sampletotest, samplevolume),
					mixer.Sample(_input.LoadingDye, _input.LoadingDyeVolume),
				)
				DNAGelloadmix = DNAGelloadmixsolution
			} else {
				DNAGelloadmix = sampletotest
			}

			// Ensure  sample will be dispensed appropriately:

			// replacing following line with temporary hard code whilst developing protocol:
			DNAGelloadmix.Type, err = wtype.LiquidTypeFromString(_input.MixingPolicy)
			if err != nil {
				execute.Errorf(_ctx, "Error in specifying MixingPolicy %s for DNA Gel: %s", _input.MixingPolicy, err.Error())
			}

			loadedsample = execute.MixInto(_ctx, _input.DNAGel,
				position,
				waterSample,
				mixer.Sample(DNAGelloadmix, samplevolume),
			)

			loadedsamples = append(loadedsamples, loadedsample)

			counter--

		}

	}

	_output.LoadedSamples = loadedsamples
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _DNA_gelAnalysis(_ctx context.Context, _input *DNA_gelInput, _output *DNA_gelOutput) {

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _DNA_gelValidation(_ctx context.Context, _input *DNA_gelInput, _output *DNA_gelOutput) {

}
func _DNA_gelRun(_ctx context.Context, input *DNA_gelInput) *DNA_gelOutput {
	output := &DNA_gelOutput{}
	_DNA_gelSetup(_ctx, input)
	_DNA_gelSteps(_ctx, input, output)
	_DNA_gelAnalysis(_ctx, input, output)
	_DNA_gelValidation(_ctx, input, output)
	return output
}

func DNA_gelRunSteps(_ctx context.Context, input *DNA_gelInput) *DNA_gelSOutput {
	soutput := &DNA_gelSOutput{}
	output := _DNA_gelRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func DNA_gelNew() interface{} {
	return &DNA_gelElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &DNA_gelInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _DNA_gelRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &DNA_gelInput{},
			Out: &DNA_gelOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type DNA_gelElement struct {
	inject.CheckedRunner
}

type DNA_gelInput struct {
	DNAGel             *wtype.LHPlate
	DNAGelRunVolume    wunit.Volume
	Ladder             *wtype.LHComponent
	LadderVolume       wunit.Volume
	LoadingDye         *wtype.LHComponent
	LoadingDyeInSample bool
	LoadingDyeVolume   wunit.Volume
	MixPlate           *wtype.LHPlate
	MixingPolicy       string
	ProjectName        string
	Reactions          []*wtype.LHComponent
	Replicates         int
	Water              *wtype.LHComponent
	WaterVol           wunit.Volume
}

type DNA_gelOutput struct {
	Error         error
	LoadedSamples []*wtype.LHComponent
}

type DNA_gelSOutput struct {
	Data struct {
		Error error
	}
	Outputs struct {
		LoadedSamples []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "DNA_gel",
		Constructor: DNA_gelNew,
		Desc: component.ComponentDesc{
			Desc: "This protocol will load DNA samples on an E-GEL for DNA analysis, and can take input from other protocols which exports an array of LHComponents.\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/DNA_gel/DNA_gel.an",
			Params: []component.ParamDesc{
				{Name: "DNAGel", Desc: "gel\n", Kind: "Inputs"},
				{Name: "DNAGelRunVolume", Desc: "", Kind: "Parameters"},
				{Name: "Ladder", Desc: "", Kind: "Inputs"},
				{Name: "LadderVolume", Desc: "", Kind: "Parameters"},
				{Name: "LoadingDye", Desc: "WaterSolution //Chemspiderlink // not correct link but similar desirable\n", Kind: "Inputs"},
				{Name: "LoadingDyeInSample", Desc: "select true if DNA samples already contain Loading Dye\n", Kind: "Parameters"},
				{Name: "LoadingDyeVolume", Desc: "", Kind: "Parameters"},
				{Name: "MixPlate", Desc: "plate to mix samples if required\n", Kind: "Inputs"},
				{Name: "MixingPolicy", Desc: "default is Load policy but can be overriden by specifying here\n", Kind: "Parameters"},
				{Name: "ProjectName", Desc: "optionally define name of experiment\n", Kind: "Parameters"},
				{Name: "Reactions", Desc: "This is the input that comes as the output from an upstream element (AutoAssembly.an)\n", Kind: "Inputs"},
				{Name: "Replicates", Desc: "define number of technical replicates\n", Kind: "Parameters"},
				{Name: "Water", Desc: "", Kind: "Inputs"},
				{Name: "WaterVol", Desc: "", Kind: "Parameters"},
				{Name: "Error", Desc: "", Kind: "Data"},
				{Name: "LoadedSamples", Desc: "Gel\n", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
