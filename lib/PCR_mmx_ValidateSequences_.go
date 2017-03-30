// Perform a single pcr reaction per element and validate that the primers will be expected to bind once each to the template sequence. Exact primer matches only.
// Thermocycle conditions are calculated from the input sequences and polymerase name.
// Valid Polymerases for calculation of properties are "Q5Polymerase" and "Taq".
package lib

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/oligos"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"math"
	"strings"
)

// Input parameters for this protocol (data)

// PCRprep parameters:

// DNA sequence of template from which amplicon will be amplified

// Forward primer sequence

// Reverse primer sequence

// Volume of mastermix to add to the reaction

// Select this if the primers have already been added to the mastermix.
// If this is selected no primers will be added to any reactions.
// Should only be used if all reactions share the same primers.

// Select this if the polymerase has already been added to the mastermix.

// Volume of forward primer to add. Will only be used if PrimersalreadyAddedtoMasterMix is not selected.

// Volume of reverse primer to add. Will only be used if PrimersalreadyAddedtoMasterMix is not selected.

// Volume of template to add

// Volume of polymerase enzyme to add per reaction. Will only be used if PolymeraseAlreadyaddedtoMastermix is not selected.

// Optionally specify a specific well position here or leave blank for Antha to decide

// Reaction parameters: (could be a entered as thermocycle parameters type possibly?)

// Initial denaturation time prior to starting cycles

// Denaturation time per cycle

// Annealing time per cycle

// Time that extension conditions will be held for following completion of all cycles.

// Degrees C below lowest MeltingTemp to set annealing Temperature.

// Data which is returned from this protocol, and data types

// Sequence of the expected product of the PCR reaction

// All exact binding site of the fwd primer found in the template

// All exact binding site of the rev primer found in the template

// Melting temperature calculated for the forward primer

// Melting temperature calculated for the reverse primer

// Annealing temperature used based upon calculated primer melting temperatures and AnnealingTempOffset.

// Extension time calculated based upon Polymerase properties and length of amplicon

// Extension time calculated based upon Polymerase properties

// Melting temperature calculated based on lowest of primer melting temperatures.

// Physical Inputs to this protocol with types

// Actual FWD primer component to use. Will only be used if PrimersalreadyAddedtoMasterMix is not selected.

// Actual REV primer component to use. Will only be used if PrimersalreadyAddedtoMasterMix is not selected.

// Actual Master mix to use

// Valid options are Q5Polymerase and Taq. Will only be used if PolymeraseAlreadyaddedtoMastermix is not selected.

// Actual Template component to use

// Type of plate to set up reaction on. Recommended plate is pcrplate

// Physical outputs from this protocol with types

// The output PCR reaction

func _PCR_mmx_ValidateSequencesRequirements() {

}

// Conditions to run on startup
func _PCR_mmx_ValidateSequencesSetup(_ctx context.Context, _input *PCR_mmx_ValidateSequencesInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _PCR_mmx_ValidateSequencesSteps(_ctx context.Context, _input *PCR_mmx_ValidateSequencesInput, _output *PCR_mmx_ValidateSequencesOutput) {

	// rename components
	_input.Template.CName = _input.TemplateSequence.Name()
	_input.FwdPrimer.CName = _input.FwdPrimerSeq.Name()
	_input.RevPrimer.CName = _input.RevPrimerSeq.Name()

	// check for non-specific binding. Exact matches only.
	_output.FwdPrimerSites = sequences.FindSeqsinSeqs(_input.TemplateSequence.Sequence(), []string{_input.FwdPrimerSeq.Sequence()})

	_output.RevPrimerSites = sequences.FindSeqsinSeqs(_input.TemplateSequence.Sequence(), []string{_input.RevPrimerSeq.Sequence()})

	if len(_output.FwdPrimerSites) != 1 || len(_output.RevPrimerSites) != 1 {

		errordescription := fmt.Sprint(
			fmt.Sprint("Unexpected number of primer binding sites found in template"),
			text.Print("FwdPrimerSitesfound:", fmt.Sprint(_output.FwdPrimerSites)),
			text.Print("RevPrimerSitesfound:", fmt.Sprint(_output.RevPrimerSites)),
		)

		execute.Errorf(_ctx, errordescription)
	}

	fwdposition := _output.FwdPrimerSites[0].Positions[0]

	revposition := _output.RevPrimerSites[0].Positions[0]

	var startposition int
	var endposition int

	if !_output.FwdPrimerSites[0].Reverse && _output.RevPrimerSites[0].Reverse && fwdposition < revposition {
		startposition = fwdposition
		endposition = revposition
	} else if _output.FwdPrimerSites[0].Reverse && !_output.RevPrimerSites[0].Reverse && fwdposition < revposition {
		startposition = revposition
		endposition = fwdposition
	}

	// work out what the pcr product will be
	_output.Amplicon = oligos.DNAregion(_input.TemplateSequence, startposition, endposition)

	// Make a mastermix
	mmxSample := mixer.Sample(_input.MasterMix, _input.MasterMixVolume)

	// pipette out to make mastermix
	mastermix := execute.MixInto(_ctx, _input.OutPlate, _input.OptionalWellPosition, mmxSample)

	// rest samples to zero
	samples := make([]*wtype.LHComponent, 0)

	// add primers

	if !_input.PrimersalreadyAddedtoMasterMix {
		FwdPrimerSample := mixer.Sample(_input.FwdPrimer, _input.FwdPrimerVol)
		samples = append(samples, FwdPrimerSample)
		RevPrimerSample := mixer.Sample(_input.RevPrimer, _input.RevPrimerVol)
		samples = append(samples, RevPrimerSample)
	}

	// add template
	templateSample := mixer.Sample(_input.Template, _input.Templatevolume)
	samples = append(samples, templateSample)

	for j := range samples {
		mastermix = execute.Mix(_ctx, mastermix, samples[j])
	}
	reaction := mastermix

	// this needs to go after an initial denaturation!
	if !_input.PolymeraseAlreadyaddedtoMastermix {
		polySample := mixer.Sample(_input.PCRPolymerase, _input.PolymeraseVolume)

		reaction = execute.Mix(_ctx, reaction, polySample)
	}

	// thermocycle parameters called from enzyme lookup:
	polymerase := _input.PCRPolymerase.CName

	polymeraseproperties, polymerasefound := enzymes.DNApolymerasetemps[polymerase]

	if !polymerasefound {
		validoptions := make([]string, 0)
		for polymerasename := range enzymes.DNApolymerasetemps {
			validoptions = append(validoptions, polymerasename)
		}

		execute.Errorf(_ctx, "No Properties for %s found. Valid options are: %s", polymerase, strings.Join(validoptions, ","))
	}

	var found bool

	_output.ExtensionTemp, found = polymeraseproperties["extensiontemp"]
	if !found {
		execute.Errorf(_ctx, "No extension temp found for polymerase %s", polymerase)
	}
	_output.MeltingTemp, found = polymeraseproperties["meltingtemp"]
	if !found {
		execute.Errorf(_ctx, "No melting temp found for polymerase %s", polymerase)
	}

	var err error

	_output.ExtensionTime, err = enzymes.CalculateExtensionTime(_input.PCRPolymerase, _output.Amplicon)

	if err != nil {
		execute.Errorf(_ctx, "Can't calculate extension time of polymerase: %s", err.Error())
	}

	// work out annealing temperature
	_output.FwdPrimerMeltingTemp = oligos.BasicMeltingTemp(_input.FwdPrimerSeq)

	_output.RevPrimerMeltingTemp = oligos.BasicMeltingTemp(_input.RevPrimerSeq)

	// check which primer has the lowest melting temperature
	lowest := math.Min(_output.FwdPrimerMeltingTemp.SIValue(), _output.RevPrimerMeltingTemp.SIValue())

	// start PCR AnnealingTempOffset degrees below lowest melting temp
	_output.AnnealingTemp = wunit.NewTemperature(lowest-_input.AnnealingTempOffset.SIValue(), "C")

	// initial Denaturation

	r1 := execute.Incubate(_ctx, reaction, _output.MeltingTemp, _input.InitDenaturationTime, false)

	for i := 0; i < _input.NumberOfCycles; i++ {

		// Denature

		r1 = execute.Incubate(_ctx, r1, _output.MeltingTemp, _input.DenaturationTime, false)

		// Anneal
		r1 = execute.Incubate(_ctx, r1, _output.AnnealingTemp, _input.AnnealingTime, false)

		// Extend
		r1 = execute.Incubate(_ctx, r1, _output.ExtensionTemp, _output.ExtensionTime, false)

	}
	// Final Extension
	r1 = execute.Incubate(_ctx, r1, _output.ExtensionTemp, _input.FinalExtensionTime, false)

	// all done
	_output.Reaction = r1

	_output.Reaction.CName = _input.ReactionName
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PCR_mmx_ValidateSequencesAnalysis(_ctx context.Context, _input *PCR_mmx_ValidateSequencesInput, _output *PCR_mmx_ValidateSequencesOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PCR_mmx_ValidateSequencesValidation(_ctx context.Context, _input *PCR_mmx_ValidateSequencesInput, _output *PCR_mmx_ValidateSequencesOutput) {
}
func _PCR_mmx_ValidateSequencesRun(_ctx context.Context, input *PCR_mmx_ValidateSequencesInput) *PCR_mmx_ValidateSequencesOutput {
	output := &PCR_mmx_ValidateSequencesOutput{}
	_PCR_mmx_ValidateSequencesSetup(_ctx, input)
	_PCR_mmx_ValidateSequencesSteps(_ctx, input, output)
	_PCR_mmx_ValidateSequencesAnalysis(_ctx, input, output)
	_PCR_mmx_ValidateSequencesValidation(_ctx, input, output)
	return output
}

func PCR_mmx_ValidateSequencesRunSteps(_ctx context.Context, input *PCR_mmx_ValidateSequencesInput) *PCR_mmx_ValidateSequencesSOutput {
	soutput := &PCR_mmx_ValidateSequencesSOutput{}
	output := _PCR_mmx_ValidateSequencesRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PCR_mmx_ValidateSequencesNew() interface{} {
	return &PCR_mmx_ValidateSequencesElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PCR_mmx_ValidateSequencesInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PCR_mmx_ValidateSequencesRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PCR_mmx_ValidateSequencesInput{},
			Out: &PCR_mmx_ValidateSequencesOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type PCR_mmx_ValidateSequencesElement struct {
	inject.CheckedRunner
}

type PCR_mmx_ValidateSequencesInput struct {
	AnnealingTempOffset               wunit.Temperature
	AnnealingTime                     wunit.Time
	DenaturationTime                  wunit.Time
	FinalExtensionTime                wunit.Time
	FwdPrimer                         *wtype.LHComponent
	FwdPrimerSeq                      wtype.DNASequence
	FwdPrimerVol                      wunit.Volume
	InitDenaturationTime              wunit.Time
	MasterMix                         *wtype.LHComponent
	MasterMixVolume                   wunit.Volume
	NumberOfCycles                    int
	OptionalWellPosition              string
	OutPlate                          *wtype.LHPlate
	PCRPolymerase                     *wtype.LHComponent
	PolymeraseAlreadyaddedtoMastermix bool
	PolymeraseVolume                  wunit.Volume
	PrimersalreadyAddedtoMasterMix    bool
	ReactionName                      string
	RevPrimer                         *wtype.LHComponent
	RevPrimerSeq                      wtype.DNASequence
	RevPrimerVol                      wunit.Volume
	Template                          *wtype.LHComponent
	TemplateSequence                  wtype.DNASequence
	Templatevolume                    wunit.Volume
}

type PCR_mmx_ValidateSequencesOutput struct {
	Amplicon             wtype.DNASequence
	AnnealingTemp        wunit.Temperature
	ExtensionTemp        wunit.Temperature
	ExtensionTime        wunit.Time
	FwdPrimerMeltingTemp wunit.Temperature
	FwdPrimerSites       []search.Thingfound
	MeltingTemp          wunit.Temperature
	Reaction             *wtype.LHComponent
	RevPrimerMeltingTemp wunit.Temperature
	RevPrimerSites       []search.Thingfound
}

type PCR_mmx_ValidateSequencesSOutput struct {
	Data struct {
		Amplicon             wtype.DNASequence
		AnnealingTemp        wunit.Temperature
		ExtensionTemp        wunit.Temperature
		ExtensionTime        wunit.Time
		FwdPrimerMeltingTemp wunit.Temperature
		FwdPrimerSites       []search.Thingfound
		MeltingTemp          wunit.Temperature
		RevPrimerMeltingTemp wunit.Temperature
		RevPrimerSites       []search.Thingfound
	}
	Outputs struct {
		Reaction *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PCR_mmx_ValidateSequences",
		Constructor: PCR_mmx_ValidateSequencesNew,
		Desc: component.ComponentDesc{
			Desc: "Perform a single pcr reaction per element and validate that the primers will be expected to bind once each to the template sequence. Exact primer matches only.\nThermocycle conditions are calculated from the input sequences and polymerase name.\nValid Polymerases for calculation of properties are \"Q5Polymerase\" and \"Taq\".\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/PCR/PCR_mmx_ValidateSequences.an",
			Params: []component.ParamDesc{
				{Name: "AnnealingTempOffset", Desc: "Degrees C below lowest MeltingTemp to set annealing Temperature.\n", Kind: "Parameters"},
				{Name: "AnnealingTime", Desc: "Annealing time per cycle\n", Kind: "Parameters"},
				{Name: "DenaturationTime", Desc: "Denaturation time per cycle\n", Kind: "Parameters"},
				{Name: "FinalExtensionTime", Desc: "Time that extension conditions will be held for following completion of all cycles.\n", Kind: "Parameters"},
				{Name: "FwdPrimer", Desc: "Actual FWD primer component to use. Will only be used if PrimersalreadyAddedtoMasterMix is not selected.\n", Kind: "Inputs"},
				{Name: "FwdPrimerSeq", Desc: "Forward primer sequence\n", Kind: "Parameters"},
				{Name: "FwdPrimerVol", Desc: "Volume of forward primer to add. Will only be used if PrimersalreadyAddedtoMasterMix is not selected.\n", Kind: "Parameters"},
				{Name: "InitDenaturationTime", Desc: "Initial denaturation time prior to starting cycles\n", Kind: "Parameters"},
				{Name: "MasterMix", Desc: "Actual Master mix to use\n", Kind: "Inputs"},
				{Name: "MasterMixVolume", Desc: "Volume of mastermix to add to the reaction\n", Kind: "Parameters"},
				{Name: "NumberOfCycles", Desc: "", Kind: "Parameters"},
				{Name: "OptionalWellPosition", Desc: "Optionally specify a specific well position here or leave blank for Antha to decide\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "Type of plate to set up reaction on. Recommended plate is pcrplate\n", Kind: "Inputs"},
				{Name: "PCRPolymerase", Desc: "Valid options are Q5Polymerase and Taq. Will only be used if PolymeraseAlreadyaddedtoMastermix is not selected.\n", Kind: "Inputs"},
				{Name: "PolymeraseAlreadyaddedtoMastermix", Desc: "Select this if the polymerase has already been added to the mastermix.\n", Kind: "Parameters"},
				{Name: "PolymeraseVolume", Desc: "Volume of polymerase enzyme to add per reaction. Will only be used if PolymeraseAlreadyaddedtoMastermix is not selected.\n", Kind: "Parameters"},
				{Name: "PrimersalreadyAddedtoMasterMix", Desc: "Select this if the primers have already been added to the mastermix.\nIf this is selected no primers will be added to any reactions.\nShould only be used if all reactions share the same primers.\n", Kind: "Parameters"},
				{Name: "ReactionName", Desc: "PCRprep parameters:\n", Kind: "Parameters"},
				{Name: "RevPrimer", Desc: "Actual REV primer component to use. Will only be used if PrimersalreadyAddedtoMasterMix is not selected.\n", Kind: "Inputs"},
				{Name: "RevPrimerSeq", Desc: "Reverse primer sequence\n", Kind: "Parameters"},
				{Name: "RevPrimerVol", Desc: "Volume of reverse primer to add. Will only be used if PrimersalreadyAddedtoMasterMix is not selected.\n", Kind: "Parameters"},
				{Name: "Template", Desc: "Actual Template component to use\n", Kind: "Inputs"},
				{Name: "TemplateSequence", Desc: "DNA sequence of template from which amplicon will be amplified\n", Kind: "Parameters"},
				{Name: "Templatevolume", Desc: "Volume of template to add\n", Kind: "Parameters"},
				{Name: "Amplicon", Desc: "Sequence of the expected product of the PCR reaction\n", Kind: "Data"},
				{Name: "AnnealingTemp", Desc: "Annealing temperature used based upon calculated primer melting temperatures and AnnealingTempOffset.\n", Kind: "Data"},
				{Name: "ExtensionTemp", Desc: "Extension time calculated based upon Polymerase properties\n", Kind: "Data"},
				{Name: "ExtensionTime", Desc: "Extension time calculated based upon Polymerase properties and length of amplicon\n", Kind: "Data"},
				{Name: "FwdPrimerMeltingTemp", Desc: "Melting temperature calculated for the forward primer\n", Kind: "Data"},
				{Name: "FwdPrimerSites", Desc: "All exact binding site of the fwd primer found in the template\n", Kind: "Data"},
				{Name: "MeltingTemp", Desc: "Melting temperature calculated based on lowest of primer melting temperatures.\n", Kind: "Data"},
				{Name: "Reaction", Desc: "The output PCR reaction\n", Kind: "Outputs"},
				{Name: "RevPrimerMeltingTemp", Desc: "Melting temperature calculated for the reverse primer\n", Kind: "Data"},
				{Name: "RevPrimerSites", Desc: "All exact binding site of the rev primer found in the template\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}

/*type Polymerase struct {
	wtype.LHComponent
	Rate_BPpers float64
	Fidelity_errorrate float64 // could dictate how many colonies are checked in validation!
	Extensiontemp Temperature
	Hotstart bool
	StockConcentration Concentration // this is normally in U?
	TargetConcentration Concentration
	// this is also a glycerol solution rather than a watersolution!
}
*/
