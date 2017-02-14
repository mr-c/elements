// Protocol for transforming S.cerevisiae
package lib

import

// Place golang packages to import here
(
	"context"
	unitoperations "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/UnitOperations"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Parameters to this protocol

// cells per ml
// cells per ml
// w/v

// corresponding concentrations and set points for each dna transformation

// Boiling of single stranded carrier dna prior to transformation

// Initial culture

// Estimate is 4 hours to achieve TargetYeastConcentration

// Volume of competent cells

// transformation step

// heatshock

// Final resuspension

//Plate out

// Output data of this protocol

// Physical inputs to this protocol

// Growth

// OD measurement

// washing pellet

// transformation components:

// should be 50% w/v

// plate out

// Physical outputs to this protocol

// Conditions to run on startup
func _YeastTransformationSetup(_ctx context.Context, _input *YeastTransformationInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _YeastTransformationSteps(_ctx context.Context, _input *YeastTransformationInput, _output *YeastTransformationOutput) {

	var resuspendedcompetentcells *wtype.LHComponent

	// Pre warm media
	_input.YPD = execute.Incubate(_ctx, _input.YPD, wunit.NewTemperature(30, "C"), wunit.NewTime(10, "mins"), false)

	// sample pre-warmed media
	ypdSample := mixer.Sample(_input.YPD, _input.YPDVolume)

	// add pre-warmed media to growth plate
	starterculture := execute.MixInto(_ctx, _input.GrowthPlate, "", ypdSample)

	// add cells

	// first calculate volume to add based on concentration target and actual concentration
	cellvol, err := wunit.VolumeForTargetConcentration(_input.TargetStartingYeastConcentration, _input.YeastCells.Concentration(), _input.TotalInitialCultureVolume)

	if err != nil {
		execute.Errorf(_ctx, "NO!, problem with your conversion of concentraitons into volume: %s", err.Error())
	}

	cellsample := mixer.Sample(_input.YeastCells, cellvol)

	culturewithcells := execute.Mix(_ctx, starterculture, cellsample)

	incubatedculture := execute.Incubate(_ctx, culturewithcells, _input.InitialCultureTemp, _input.InitialCultureTime, true)

	// measure OD of yeast culture

	incubatedculture.Type = wtype.LTNeedToMix

	sampleforOD := mixer.Sample(incubatedculture, wunit.NewVolume(10, "ul"))

	absorbanceresults := AbsorbanceMeasurementRunSteps(_ctx, &AbsorbanceMeasurementInput{AbsorbanceWavelength: wunit.NewLength(600, "nm"),
		DilutionVolume:        wunit.NewVolume(90.0, "ul"),
		ExtinctionCoefficient: _input.ODYEASTExtinctionCoefficient,

		SampleForReading: sampleforOD,
		Diluent:          _input.ODDiluent,
		Plate:            _input.ODPlate},
	)

	actualconc := absorbanceresults.Data.ActualConcentration

	if actualconc.GreaterThan(_input.TargetTransformationYeastConcentration) {

		// centrifuge sample
		pellet, supernatant := unitoperations.Separate(culturewithcells)

		// discard supernatant
		wastesupernatant := mixer.SampleAll(supernatant)

		execute.MixInto(_ctx, _input.WastePlate, "", wastesupernatant)

		washsample := mixer.Sample(_input.PostSpinWashBuffer, _input.PostSpinWashVolume)

		washedpellet := execute.Mix(_ctx, pellet, washsample)

		// centrifuge sample
		pellet, supernatant = unitoperations.Separate(washedpellet)

		// discard supernatant
		wastesupernatant = mixer.SampleAll(supernatant)

		execute.MixInto(_ctx, _input.WastePlate, "", wastesupernatant)

		resuspensionsample := mixer.Sample(_input.PostSpinResuspensionBuffer, _input.PostWashResuspensionVolume)

		resuspendedcompetentcells = execute.Mix(_ctx, pellet, resuspensionsample)

	} else {
		execute.Errorf(_ctx, "Carry on incubating")
	}

	// boil sperm dna
	incubatedsscarrierdna := execute.Incubate(_ctx, _input.SSCarrierDNA, _input.SSDNABoilTemp, _input.SSDNABoilTime, false)

	// chill
	incubatedsscarrierdna = execute.Incubate(_ctx, incubatedsscarrierdna, _input.SSDNAPostBoilTemp, _input.SSDNAPostBoilTime, false)

	// calculate number of aliquots
	numberofaliquots, err := wutil.RoundDown(_input.PostWashResuspensionVolume.SIValue() / _input.CompetentCellVolumePerTransformation.SIValue())

	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	// Aliquot Competent cells
	for index := 0; index < numberofaliquots; index++ {

		aliquotSample := mixer.Sample(resuspendedcompetentcells, _input.CompetentCellVolumePerTransformation)

		aliquot := execute.Mix(_ctx, aliquotSample)

		_output.CompetentCells = append(_output.CompetentCells, aliquot)

	}
	// now transformation time
	if len(_output.CompetentCells) < len(_input.DNAToTransform) {
		execute.Errorf(_ctx, "Not enough competent cells for this many transformations")
	}

	numberofreactions := float64(len(_input.DNAToTransform) + 1)

	var mmxSamples []*wtype.LHComponent

	mmxSamples = append(mmxSamples, mixer.Sample(_input.LiAC, wunit.MultiplyVolume(_input.LiACVolume, numberofreactions)), mixer.Sample(_input.PEG3350, wunit.MultiplyVolume(_input.PEGVolume, numberofreactions)), mixer.Sample(incubatedsscarrierdna, wunit.MultiplyVolume(_input.SSCarrierDNAVolume, numberofreactions)))

	mastermix := execute.Mix(_ctx, mmxSamples...)

	for i, dnaSample := range _input.DNAToTransform {

		// centrifuge sample
		pellet, supernatant := unitoperations.Separate(_output.CompetentCells[i])

		// discard supernatant
		wastesupernatant := mixer.SampleAll(supernatant)

		// vortex and resuspend here

		compcellmix := execute.Mix(_ctx, pellet, mixer.Sample(mastermix, _input.MastermixVolumePerReaction))

		dnaVol, err := wunit.VolumeForTargetMass(_input.TargetDNAMassPerTransformation[i], _input.DNAStockConcentration[i])

		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}

		compcelldnamix := execute.Mix(_ctx, compcellmix, mixer.Sample(dnaSample, dnaVol))

		finalreaction := execute.Mix(_ctx, compcelldnamix, mixer.SampleForTotalVolume(_input.Water, _input.TotalTransformationVolume))

		incubatedreaction := execute.Incubate(_ctx, finalreaction, _input.HeatShockTemp, _input.HeatShockTime, false)

		// centrifuge sample
		pellet, supernatant = unitoperations.Separate(incubatedreaction)

		// discard supernatant
		wastesupernatant = mixer.SampleAll(supernatant)
		execute.MixInto(_ctx, _input.WastePlate, "", wastesupernatant)

		readytoplateout := execute.Mix(_ctx, pellet, mixer.Sample(_input.Water, _input.FinalResuspensionVolume))

		for _, platoutvol := range _input.PlateOutVolumes {
			plateoutSample := mixer.Sample(readytoplateout, platoutvol)

			plateoutculture := execute.MixInto(_ctx, _input.AgarPlate, "", plateoutSample)

			_output.PlatedTransformations = append(_output.PlatedTransformations, plateoutculture)
		}

	}

}

func _YeastTransformationRequirements() {

}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _YeastTransformationAnalysis(_ctx context.Context, _input *YeastTransformationInput, _output *YeastTransformationOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _YeastTransformationValidation(_ctx context.Context, _input *YeastTransformationInput, _output *YeastTransformationOutput) {

}
func _YeastTransformationRun(_ctx context.Context, input *YeastTransformationInput) *YeastTransformationOutput {
	output := &YeastTransformationOutput{}
	_YeastTransformationSetup(_ctx, input)
	_YeastTransformationSteps(_ctx, input, output)
	_YeastTransformationAnalysis(_ctx, input, output)
	_YeastTransformationValidation(_ctx, input, output)
	return output
}

func YeastTransformationRunSteps(_ctx context.Context, input *YeastTransformationInput) *YeastTransformationSOutput {
	soutput := &YeastTransformationSOutput{}
	output := _YeastTransformationRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func YeastTransformationNew() interface{} {
	return &YeastTransformationElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &YeastTransformationInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _YeastTransformationRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &YeastTransformationInput{},
			Out: &YeastTransformationOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type YeastTransformationElement struct {
	inject.CheckedRunner
}

type YeastTransformationInput struct {
	AgarPlate                              *wtype.LHPlate
	CentrifugationRCF                      wunit.Force
	CentrifugationTime                     wunit.Time
	CompetentCellVolumePerTransformation   wunit.Volume
	DNAStockConcentration                  []wunit.Concentration
	DNAToTransform                         []*wtype.LHComponent
	FinalResuspensionVolume                wunit.Volume
	GrowthPlate                            *wtype.LHPlate
	HeatShockTemp                          wunit.Temperature
	HeatShockTime                          wunit.Time
	InitialCultureTemp                     wunit.Temperature
	InitialCultureTime                     wunit.Time
	LiAC                                   *wtype.LHComponent
	LiACVolume                             wunit.Volume
	LiAcConcentration                      wunit.Concentration
	MastermixVolumePerReaction             wunit.Volume
	NumberOfWashes                         int
	ODDiluent                              *wtype.LHComponent
	ODPlate                                *wtype.LHPlate
	ODYEASTExtinctionCoefficient           float64
	PEG3350                                *wtype.LHComponent
	PEGConcentration                       wunit.Concentration
	PEGVolume                              wunit.Volume
	PlateOutVolumes                        []wunit.Volume
	PostSpinResuspensionBuffer             *wtype.LHComponent
	PostSpinWashBuffer                     *wtype.LHComponent
	PostSpinWashVolume                     wunit.Volume
	PostWashResuspensionVolume             wunit.Volume
	SSCarrierDNA                           *wtype.LHComponent
	SSCarrierDNAVolume                     wunit.Volume
	SSDNABoilTemp                          wunit.Temperature
	SSDNABoilTime                          wunit.Time
	SSDNAPostBoilTemp                      wunit.Temperature
	SSDNAPostBoilTime                      wunit.Time
	TargetDNAMassPerTransformation         []wunit.Mass
	TargetStartingYeastConcentration       wunit.Concentration
	TargetTransformationYeastConcentration wunit.Concentration
	TotalInitialCultureVolume              wunit.Volume
	TotalTransformationVolume              wunit.Volume
	WastePlate                             *wtype.LHPlate
	Water                                  *wtype.LHComponent
	YPD                                    *wtype.LHComponent
	YPDVolume                              wunit.Volume
	YeastCells                             *wtype.LHComponent
}

type YeastTransformationOutput struct {
	CompetentCells        []*wtype.LHComponent
	PlatedTransformations []*wtype.LHComponent
}

type YeastTransformationSOutput struct {
	Data struct {
	}
	Outputs struct {
		CompetentCells        []*wtype.LHComponent
		PlatedTransformations []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "YeastTransformation",
		Constructor: YeastTransformationNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol for transforming S.cerevisiae\n",
			Path: "src/github.com/antha-lang/elements/an/Yeast/YeastTransformation/YeastTransformation.an",
			Params: []component.ParamDesc{
				{Name: "AgarPlate", Desc: "plate out\n", Kind: "Inputs"},
				{Name: "CentrifugationRCF", Desc: "", Kind: "Parameters"},
				{Name: "CentrifugationTime", Desc: "", Kind: "Parameters"},
				{Name: "CompetentCellVolumePerTransformation", Desc: "Volume of competent cells\n", Kind: "Parameters"},
				{Name: "DNAStockConcentration", Desc: "corresponding concentrations and set points for each dna transformation\n", Kind: "Parameters"},
				{Name: "DNAToTransform", Desc: "", Kind: "Inputs"},
				{Name: "FinalResuspensionVolume", Desc: "Final resuspension\n", Kind: "Parameters"},
				{Name: "GrowthPlate", Desc: "", Kind: "Inputs"},
				{Name: "HeatShockTemp", Desc: "heatshock\n", Kind: "Parameters"},
				{Name: "HeatShockTime", Desc: "", Kind: "Parameters"},
				{Name: "InitialCultureTemp", Desc: "", Kind: "Parameters"},
				{Name: "InitialCultureTime", Desc: "Estimate is 4 hours to achieve TargetYeastConcentration\n", Kind: "Parameters"},
				{Name: "LiAC", Desc: "", Kind: "Inputs"},
				{Name: "LiACVolume", Desc: "", Kind: "Parameters"},
				{Name: "LiAcConcentration", Desc: "", Kind: "Parameters"},
				{Name: "MastermixVolumePerReaction", Desc: "transformation step\n", Kind: "Parameters"},
				{Name: "NumberOfWashes", Desc: "", Kind: "Parameters"},
				{Name: "ODDiluent", Desc: "OD measurement\n", Kind: "Inputs"},
				{Name: "ODPlate", Desc: "", Kind: "Inputs"},
				{Name: "ODYEASTExtinctionCoefficient", Desc: "", Kind: "Parameters"},
				{Name: "PEG3350", Desc: "should be 50% w/v\n", Kind: "Inputs"},
				{Name: "PEGConcentration", Desc: "w/v\n", Kind: "Parameters"},
				{Name: "PEGVolume", Desc: "", Kind: "Parameters"},
				{Name: "PlateOutVolumes", Desc: "Plate out\n", Kind: "Parameters"},
				{Name: "PostSpinResuspensionBuffer", Desc: "", Kind: "Inputs"},
				{Name: "PostSpinWashBuffer", Desc: "washing pellet\n", Kind: "Inputs"},
				{Name: "PostSpinWashVolume", Desc: "", Kind: "Parameters"},
				{Name: "PostWashResuspensionVolume", Desc: "", Kind: "Parameters"},
				{Name: "SSCarrierDNA", Desc: "transformation components:\n", Kind: "Inputs"},
				{Name: "SSCarrierDNAVolume", Desc: "", Kind: "Parameters"},
				{Name: "SSDNABoilTemp", Desc: "Boiling of single stranded carrier dna prior to transformation\n", Kind: "Parameters"},
				{Name: "SSDNABoilTime", Desc: "", Kind: "Parameters"},
				{Name: "SSDNAPostBoilTemp", Desc: "", Kind: "Parameters"},
				{Name: "SSDNAPostBoilTime", Desc: "", Kind: "Parameters"},
				{Name: "TargetDNAMassPerTransformation", Desc: "", Kind: "Parameters"},
				{Name: "TargetStartingYeastConcentration", Desc: "cells per ml\n", Kind: "Parameters"},
				{Name: "TargetTransformationYeastConcentration", Desc: "cells per ml\n", Kind: "Parameters"},
				{Name: "TotalInitialCultureVolume", Desc: "", Kind: "Parameters"},
				{Name: "TotalTransformationVolume", Desc: "", Kind: "Parameters"},
				{Name: "WastePlate", Desc: "", Kind: "Inputs"},
				{Name: "Water", Desc: "", Kind: "Inputs"},
				{Name: "YPD", Desc: "Growth\n", Kind: "Inputs"},
				{Name: "YPDVolume", Desc: "Initial culture\n", Kind: "Parameters"},
				{Name: "YeastCells", Desc: "", Kind: "Inputs"},
				{Name: "CompetentCells", Desc: "", Kind: "Outputs"},
				{Name: "PlatedTransformations", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
