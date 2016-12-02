// Validate that the primers will be expected to bind once each to the template sequence
// exact matches only
package lib

import (
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
	"golang.org/x/net/context"
	"strings"
)

// Input parameters for this protocol (data)

// PCRprep parameters:

// Leave blank for Antha to decide

// Reaction parameters: (could be a entered as thermocycle parameters type possibly?)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _PCR_ValidateSequencesRequirements() {

}

// Conditions to run on startup
func _PCR_ValidateSequencesSetup(_ctx context.Context, _input *PCR_ValidateSequencesInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _PCR_ValidateSequencesSteps(_ctx context.Context, _input *PCR_ValidateSequencesInput, _output *PCR_ValidateSequencesOutput) {

	// rename components

	_input.Template.CName = _input.Targetsequence.Name()
	_input.FwdPrimer.CName = _input.FwdPrimerSeq.Name()
	_input.RevPrimer.CName = _input.RevPrimerSeq.Name()

	_output.FwdPrimerSites = sequences.FindSeqsinSeqs(_input.Targetsequence.Sequence(), []string{_input.FwdPrimerSeq.Sequence()})

	_output.RevPrimerSites = sequences.FindSeqsinSeqs(_input.Targetsequence.Sequence(), []string{_input.RevPrimerSeq.Sequence()})

	if len(_output.FwdPrimerSites) != 1 || len(_output.RevPrimerSites) != 1 {

		errordescription := fmt.Sprint(
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

	// This is the pcr product
	_output.Amplicon = oligos.DNAregion(_input.Targetsequence, startposition, endposition)

	// Make a mastermix

	mmxSample := mixer.Sample(_input.MasterMix, _input.MasterMixVolume)

	// pipette out to make mastermix
	mastermix := execute.MixInto(_ctx, _input.OutPlate, _input.OptionalWellPosition, mmxSample)

	// rest samples to zero
	samples := make([]*wtype.LHComponent, 0)

	// if this is false do stuff inside {}

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

		execute.Errorf(_ctx, "No Properties for", polymerase, "found.", "Valid options are:", strings.Join(validoptions, ","))
	}

	var found bool

	_output.ExtensionTemp, found = polymeraseproperties["extensiontemp"]
	if !found {
		execute.Errorf(_ctx, "No extension temp found for polymerase ", polymerase)
	}
	_output.MeltingTemp, found = polymeraseproperties["meltingtemp"]
	if !found {
		execute.Errorf(_ctx, "No melting temp found for polymerase ", polymerase)
	}

	var err error

	_output.ExtensionTime, err = enzymes.CalculateExtensionTime(_input.PCRPolymerase, _output.Amplicon)

	if err != nil {
		execute.Errorf(_ctx, "Can't calculate extensiontime of polymerase.", err.Error())
	}

	// work out annealing temperature

	_output.Fwdprimermeltingtemp = oligos.BasicMeltingTemp(_input.FwdPrimerSeq)

	_output.Revprimermeltingtemp = oligos.BasicMeltingTemp(_input.RevPrimerSeq)

	// check which primer has the lowest melting temperature
	if _output.Fwdprimermeltingtemp.SIValue() < _output.Revprimermeltingtemp.SIValue() {
		// start PCR 3 degress below lowest melting temp
		_output.AnnealingTemp = wunit.NewTemperature(_output.Fwdprimermeltingtemp.SIValue()-3.0, "C")
	} else {
		// start PCR 3 degress below lowest melting temp
		_output.AnnealingTemp = wunit.NewTemperature(_output.Revprimermeltingtemp.SIValue()-3.0, "C")
	}

	// initial Denaturation

	r1 := execute.Incubate(_ctx, reaction, _output.MeltingTemp, _input.InitDenaturationtime, false)

	for i := 0; i < _input.Numberofcycles; i++ {

		// Denature

		r1 = execute.Incubate(_ctx, r1, _output.MeltingTemp, _input.Denaturationtime, false)

		// Anneal
		r1 = execute.Incubate(_ctx, r1, _output.AnnealingTemp, _input.AnnealingTime, false)

		//extensiontime := TargetTemplatelengthinBP/PCRPolymerase.RateBPpers // we'll get type issues here so leave it out for now

		// Extend
		r1 = execute.Incubate(_ctx, r1, _output.ExtensionTemp, _output.ExtensionTime, false)

	}
	// Final Extension
	r1 = execute.Incubate(_ctx, r1, _output.ExtensionTemp, _input.Finalextensiontime, false)

	// all done
	_output.Reaction = r1

	_output.Reaction.CName = _input.ReactionName
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PCR_ValidateSequencesAnalysis(_ctx context.Context, _input *PCR_ValidateSequencesInput, _output *PCR_ValidateSequencesOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PCR_ValidateSequencesValidation(_ctx context.Context, _input *PCR_ValidateSequencesInput, _output *PCR_ValidateSequencesOutput) {
}
func _PCR_ValidateSequencesRun(_ctx context.Context, input *PCR_ValidateSequencesInput) *PCR_ValidateSequencesOutput {
	output := &PCR_ValidateSequencesOutput{}
	_PCR_ValidateSequencesSetup(_ctx, input)
	_PCR_ValidateSequencesSteps(_ctx, input, output)
	_PCR_ValidateSequencesAnalysis(_ctx, input, output)
	_PCR_ValidateSequencesValidation(_ctx, input, output)
	return output
}

func PCR_ValidateSequencesRunSteps(_ctx context.Context, input *PCR_ValidateSequencesInput) *PCR_ValidateSequencesSOutput {
	soutput := &PCR_ValidateSequencesSOutput{}
	output := _PCR_ValidateSequencesRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PCR_ValidateSequencesNew() interface{} {
	return &PCR_ValidateSequencesElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PCR_ValidateSequencesInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PCR_ValidateSequencesRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PCR_ValidateSequencesInput{},
			Out: &PCR_ValidateSequencesOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PCR_ValidateSequencesElement struct {
	inject.CheckedRunner
}

type PCR_ValidateSequencesInput struct {
	AnnealingTime                     wunit.Time
	Denaturationtime                  wunit.Time
	Finalextensiontime                wunit.Time
	FwdPrimer                         *wtype.LHComponent
	FwdPrimerSeq                      wtype.DNASequence
	FwdPrimerVol                      wunit.Volume
	InitDenaturationtime              wunit.Time
	MasterMix                         *wtype.LHComponent
	MasterMixVolume                   wunit.Volume
	Numberofcycles                    int
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
	Targetsequence                    wtype.DNASequence
	Template                          *wtype.LHComponent
	Templatevolume                    wunit.Volume
}

type PCR_ValidateSequencesOutput struct {
	Amplicon             wtype.DNASequence
	AnnealingTemp        wunit.Temperature
	ExtensionTemp        wunit.Temperature
	ExtensionTime        wunit.Time
	FwdPrimerSites       []search.Thingfound
	Fwdprimermeltingtemp wunit.Temperature
	MeltingTemp          wunit.Temperature
	Reaction             *wtype.LHComponent
	RevPrimerSites       []search.Thingfound
	Revprimermeltingtemp wunit.Temperature
}

type PCR_ValidateSequencesSOutput struct {
	Data struct {
		Amplicon             wtype.DNASequence
		AnnealingTemp        wunit.Temperature
		ExtensionTemp        wunit.Temperature
		ExtensionTime        wunit.Time
		FwdPrimerSites       []search.Thingfound
		Fwdprimermeltingtemp wunit.Temperature
		MeltingTemp          wunit.Temperature
		RevPrimerSites       []search.Thingfound
		Revprimermeltingtemp wunit.Temperature
	}
	Outputs struct {
		Reaction *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PCR_ValidateSequences",
		Constructor: PCR_ValidateSequencesNew,
		Desc: component.ComponentDesc{
			Desc: "Validate that the primers will be expected to bind once each to the template sequence\nexact matches only\n",
			Path: "src/github.com/antha-lang/elements/starter/MakeMasterMix_PCR/PCR_primerbind.an",
			Params: []component.ParamDesc{
				{Name: "AnnealingTime", Desc: "", Kind: "Parameters"},
				{Name: "Denaturationtime", Desc: "", Kind: "Parameters"},
				{Name: "Finalextensiontime", Desc: "", Kind: "Parameters"},
				{Name: "FwdPrimer", Desc: "", Kind: "Inputs"},
				{Name: "FwdPrimerSeq", Desc: "", Kind: "Parameters"},
				{Name: "FwdPrimerVol", Desc: "", Kind: "Parameters"},
				{Name: "InitDenaturationtime", Desc: "", Kind: "Parameters"},
				{Name: "MasterMix", Desc: "", Kind: "Inputs"},
				{Name: "MasterMixVolume", Desc: "", Kind: "Parameters"},
				{Name: "Numberofcycles", Desc: "Reaction parameters: (could be a entered as thermocycle parameters type possibly?)\n", Kind: "Parameters"},
				{Name: "OptionalWellPosition", Desc: "Leave blank for Antha to decide\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "PCRPolymerase", Desc: "", Kind: "Inputs"},
				{Name: "PolymeraseAlreadyaddedtoMastermix", Desc: "", Kind: "Parameters"},
				{Name: "PolymeraseVolume", Desc: "", Kind: "Parameters"},
				{Name: "PrimersalreadyAddedtoMasterMix", Desc: "", Kind: "Parameters"},
				{Name: "ReactionName", Desc: "PCRprep parameters:\n", Kind: "Parameters"},
				{Name: "RevPrimer", Desc: "", Kind: "Inputs"},
				{Name: "RevPrimerSeq", Desc: "", Kind: "Parameters"},
				{Name: "RevPrimerVol", Desc: "", Kind: "Parameters"},
				{Name: "Targetsequence", Desc: "", Kind: "Parameters"},
				{Name: "Template", Desc: "", Kind: "Inputs"},
				{Name: "Templatevolume", Desc: "", Kind: "Parameters"},
				{Name: "Amplicon", Desc: "", Kind: "Data"},
				{Name: "AnnealingTemp", Desc: "", Kind: "Data"},
				{Name: "ExtensionTemp", Desc: "", Kind: "Data"},
				{Name: "ExtensionTime", Desc: "", Kind: "Data"},
				{Name: "FwdPrimerSites", Desc: "", Kind: "Data"},
				{Name: "Fwdprimermeltingtemp", Desc: "", Kind: "Data"},
				{Name: "MeltingTemp", Desc: "", Kind: "Data"},
				{Name: "Reaction", Desc: "", Kind: "Outputs"},
				{Name: "RevPrimerSites", Desc: "", Kind: "Data"},
				{Name: "Revprimermeltingtemp", Desc: "", Kind: "Data"},
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
