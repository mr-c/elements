// This protocol will load DNA samples on an E-GEL for DNA analysis. The loading dye can also be added to the samples if selected.
//A global volume will be loaded for all samples and can take input from other protocols which exports an array of LHComponents.
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

//default is Load policy but can be overriden by specifying here (i.e. for viscous samples requiring slow dispensing)
//specify the volume of DNA ladder to add
//select true if DNA samples already contain Loading Dye. If this is selected, the step to add loading dye will be skipped
//default is NeedToMix but can be overriden by specifying here (i.e. for hard to mix reaction samples)
//specify the volume of loading dye to add to each sample
//If selected, loading dye will be mixed with sample in input plate (instead of mixing in a seperate plate)
//define number of technical replicates
//specify the total running volume of the E-Gel in use
//specify volume of water to add to prime the E-Gel

// Data which is returned from this protocol, and data types

//error reporting

// Physical Inputs to this protocol with types

// E-GEL type. (Current valid options are the 48 and 96 well precast E-GELs from Thermo-Fisher)
//DNA ladder
//loading dye to mix with samples
// plate to mix samples if required
//Specifies the samples to load. These may be set here using the NewLHComponents element or fed in from a previous element such as AutoPCR_multi.
//water

// Physical outputs from this protocol with types

//samples outputted as an array which can be wired into downstream protocols

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
	if _input.GelLoadingMixingPolicy == "" {
		_input.GelLoadingMixingPolicy = "Load"
	}

	if _input.LoadingDyeMixingPolicy == "" {
		_input.LoadingDyeMixingPolicy = "NeedToMix"
	}

	//get well positions of DNA Gel from plate library ensuring the list is by row rather than by column
	var wells []string = _input.DNAGel.AllWellPositions(wtype.BYROW)

	//setup liquid handling component variables
	var DNAGelloadmix *wtype.LHComponent
	var loadedsample *wtype.LHComponent

	//begin counter at first well position as E-GEL must be run upside down
	var counter int = len(wells) - 1

	//work out and copy sample volume
	samplevolume := (wunit.CopyVolume(_input.TotalVolumePerWell))

	//work out volume of water required
	samplevolume.Subtract(_input.WaterVolume)

	//assign water to specific liquid handling load type
	_input.Water.Type = wtype.LTloadwater

	//calculate and loop through specified number of replicates
	for j := 0; j < _input.Replicates; j++ {

		//range through the reactions input array and perform specified actions
		for i := range _input.Reactions {

			//sample water at specified water volume
			waterSample := mixer.Sample(_input.Water, _input.WaterVolume)

			//update position to correspond to counter
			position := wells[counter]

			//get well coordinates from correct position
			wellcoords := wtype.MakeWellCoordsA1(position)

			//add ladder

			//if it is the last column, add a ladder sample
			if wellcoords.X == _input.DNAGel.WlsX-1 {

				//attribute specified mixinpolicy to the DNA ladder
				_input.Ladder.Type, err = wtype.LiquidTypeFromString(_input.GelLoadingMixingPolicy)
				if err != nil {
					execute.Errorf(_ctx, "Error in specifying GelLoadingMixingPolicy %s for DNA Gel: %s", _input.GelLoadingMixingPolicy, err.Error())
				}

				//perform liquid handling for addiiton of ladder sample
				laddersample := execute.MixInto(_ctx, _input.DNAGel,
					position,
					mixer.SampleForTotalVolume(_input.Water, _input.TotalVolumePerWell),
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

				//attribute specified mixinpolicy to the LoadingDye
				_input.LoadingDye.Type, err = wtype.LiquidTypeFromString(_input.LoadingDyeMixingPolicy)
				if err != nil {
					execute.Errorf(_ctx, "Error in specifying LoadingDyeMixingPolicy %s for DNA Gel: %s", _input.LoadingDyeMixingPolicy, err.Error())
				}

				//perform liquid handling for addiiton and mixing of the loading dye
				var DNAGelloadmixsolution *wtype.LHComponent

				// determine if OptimisePlateUsage selected and if so, perform mix on input plate, else perform mix on seperate plate
				if _input.OptimisePlateUsage == true {
					sample := execute.Mix(_ctx, mixer.Sample(sampletotest, samplevolume))
					DNAGelloadmixsolution = execute.Mix(_ctx, sample,
						mixer.Sample(_input.LoadingDye,
							_input.LoadingDyeVolume))
				} else {
					DNAGelloadmixsolution = execute.MixInto(_ctx, _input.MixPlate,
						"",
						mixer.Sample(sampletotest, samplevolume),
						mixer.Sample(_input.LoadingDye, _input.LoadingDyeVolume))
				}

				DNAGelloadmix = DNAGelloadmixsolution
			} else {
				DNAGelloadmix = sampletotest
			}

			//attribute specified mixinpolicy to the samples
			DNAGelloadmix.Type, err = wtype.LiquidTypeFromString(_input.GelLoadingMixingPolicy)
			if err != nil {
				execute.Errorf(_ctx, "Error in specifying GelLoadingMixingPolicy %s for DNA Gel: %s", _input.GelLoadingMixingPolicy, err.Error())
			}

			//load the DNA samples (either mixed with loading dye or pre-mixed) to the E-GEL
			loadedsample = execute.MixInto(_ctx, _input.DNAGel,
				position,
				waterSample,
				mixer.Sample(DNAGelloadmix, samplevolume),
			)
			//add the loaded samples to the loadedsamples array
			loadedsamples = append(loadedsamples, loadedsample)

			//decrease counter by 1 as loading the E-Gel backwards becuase of position constraints
			counter--

		}

	}
	//update output variable LoadedSamples with the output of the protocol
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
	DNAGel                 *wtype.LHPlate
	GelLoadingMixingPolicy string
	Ladder                 *wtype.LHComponent
	LadderVolume           wunit.Volume
	LoadingDye             *wtype.LHComponent
	LoadingDyeInSample     bool
	LoadingDyeMixingPolicy string
	LoadingDyeVolume       wunit.Volume
	MixPlate               *wtype.LHPlate
	OptimisePlateUsage     bool
	Reactions              []*wtype.LHComponent
	Replicates             int
	TotalVolumePerWell     wunit.Volume
	Water                  *wtype.LHComponent
	WaterVolume            wunit.Volume
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
			Desc: "This protocol will load DNA samples on an E-GEL for DNA analysis. The loading dye can also be added to the samples if selected.\nA global volume will be loaded for all samples and can take input from other protocols which exports an array of LHComponents.\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/DNA_gel/DNA_gel.an",
			Params: []component.ParamDesc{
				{Name: "DNAGel", Desc: "E-GEL type. (Current valid options are the 48 and 96 well precast E-GELs from Thermo-Fisher)\n", Kind: "Inputs"},
				{Name: "GelLoadingMixingPolicy", Desc: "default is Load policy but can be overriden by specifying here (i.e. for viscous samples requiring slow dispensing)\n", Kind: "Parameters"},
				{Name: "Ladder", Desc: "DNA ladder\n", Kind: "Inputs"},
				{Name: "LadderVolume", Desc: "specify the volume of DNA ladder to add\n", Kind: "Parameters"},
				{Name: "LoadingDye", Desc: "loading dye to mix with samples\n", Kind: "Inputs"},
				{Name: "LoadingDyeInSample", Desc: "select true if DNA samples already contain Loading Dye. If this is selected, the step to add loading dye will be skipped\n", Kind: "Parameters"},
				{Name: "LoadingDyeMixingPolicy", Desc: "default is NeedToMix but can be overriden by specifying here (i.e. for hard to mix reaction samples)\n", Kind: "Parameters"},
				{Name: "LoadingDyeVolume", Desc: "specify the volume of loading dye to add to each sample\n", Kind: "Parameters"},
				{Name: "MixPlate", Desc: "plate to mix samples if required\n", Kind: "Inputs"},
				{Name: "OptimisePlateUsage", Desc: "If selected, loading dye will be mixed with sample in input plate (instead of mixing in a seperate plate)\n", Kind: "Parameters"},
				{Name: "Reactions", Desc: "Specifies the samples to load. These may be set here using the NewLHComponents element or fed in from a previous element such as AutoPCR_multi.\n", Kind: "Inputs"},
				{Name: "Replicates", Desc: "define number of technical replicates\n", Kind: "Parameters"},
				{Name: "TotalVolumePerWell", Desc: "specify the total running volume of the E-Gel in use\n", Kind: "Parameters"},
				{Name: "Water", Desc: "water\n", Kind: "Inputs"},
				{Name: "WaterVolume", Desc: "specify volume of water to add to prime the E-Gel\n", Kind: "Parameters"},
				{Name: "Error", Desc: "error reporting\n", Kind: "Data"},
				{Name: "LoadedSamples", Desc: "samples outputted as an array which can be wired into downstream protocols\n", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
