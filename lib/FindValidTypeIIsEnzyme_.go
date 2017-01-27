// The protocol is intended to return a list of valid TypeIIs enzymes which are compatible with a set of sequences and optionally the vector sequence.
// A list of ApprovedEnzymes enzymes can be specified. If no enzyme
// from the list is feasible to use (i.e. due to the presence of existing restriction sites in a part)
// all typeIIs enzymes will be screened to find feasible backup options.
package lib

import (
	"context"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/igem"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"strconv"
	"strings"
)

// Input parameters for this protocol (data)

// list sequences as strings
// optionally specify the vector sequence.
// list preferred enzyme names

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// Input Requirement specification
func _FindValidTypeIIsEnzymeRequirements() {
	// e.g. are MoClo types valid?
}

// Conditions to run on startup
func _FindValidTypeIIsEnzymeSetup(_ctx context.Context, _input *FindValidTypeIIsEnzymeInput) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _FindValidTypeIIsEnzymeSteps(_ctx context.Context, _input *FindValidTypeIIsEnzymeInput, _output *FindValidTypeIIsEnzymeOutput) {

	// make an empty array of DNA Sequences ready to fill
	partsinorder := make([]wtype.DNASequence, 0)

	for i, part := range _input.Seqsinorder {
		if strings.Contains(part, "BBa_") {
			part = igem.GetSequence(part)
		}
		partDNA := wtype.MakeLinearDNASequence("Part "+strconv.Itoa(i), part)

		partsinorder = append(partsinorder, partDNA)
	}

	if _input.Vector != "" {
		// add vector
		// make vector into an antha type DNASequence
		vectordata := wtype.MakePlasmidDNASequence("Vector", _input.Vector)
		partsinorder = append(partsinorder, vectordata)
	}
	// Find all possible typeIIs enzymes we could use for these sequences (i.e. non cutters of all parts)
	possibilities := lookup.FindEnzymeNamesofClass("TypeIIs")

	for _, possibility := range possibilities {
		// check number of sites per part !
		var sitefound bool

		enz, err := lookup.EnzymeLookup(possibility)

		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}

		for _, part := range partsinorder {

			info := enzymes.Restrictionsitefinder(part, []wtype.RestrictionEnzyme{enz})
			if len(info) != 0 {
				if info[0].Sitefound == true {
					sitefound = true
				}
			}
		}
		if !sitefound {
			_output.BackupEnzymeNames = append(_output.BackupEnzymeNames, possibility)
			_output.BackupEnzymes = append(_output.BackupEnzymes, enz)
		}
	}

	for _, enzyme := range _input.ApprovedEnzymes {

		var sitefound bool

		// check number of sites per part !
		enz, err := lookup.EnzymeLookup(enzyme)

		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}

		for _, part := range partsinorder {

			info := enzymes.Restrictionsitefinder(part, []wtype.RestrictionEnzyme{enz})
			if len(info) != 0 {
				if info[0].Sitefound {
					sitefound = true
				}
			}
		}
		if !sitefound {
			_output.ValidApprovedEnzymes = append(_output.ValidApprovedEnzymes, enzyme)
		}
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _FindValidTypeIIsEnzymeAnalysis(_ctx context.Context, _input *FindValidTypeIIsEnzymeInput, _output *FindValidTypeIIsEnzymeOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _FindValidTypeIIsEnzymeValidation(_ctx context.Context, _input *FindValidTypeIIsEnzymeInput, _output *FindValidTypeIIsEnzymeOutput) {
}
func _FindValidTypeIIsEnzymeRun(_ctx context.Context, input *FindValidTypeIIsEnzymeInput) *FindValidTypeIIsEnzymeOutput {
	output := &FindValidTypeIIsEnzymeOutput{}
	_FindValidTypeIIsEnzymeSetup(_ctx, input)
	_FindValidTypeIIsEnzymeSteps(_ctx, input, output)
	_FindValidTypeIIsEnzymeAnalysis(_ctx, input, output)
	_FindValidTypeIIsEnzymeValidation(_ctx, input, output)
	return output
}

func FindValidTypeIIsEnzymeRunSteps(_ctx context.Context, input *FindValidTypeIIsEnzymeInput) *FindValidTypeIIsEnzymeSOutput {
	soutput := &FindValidTypeIIsEnzymeSOutput{}
	output := _FindValidTypeIIsEnzymeRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func FindValidTypeIIsEnzymeNew() interface{} {
	return &FindValidTypeIIsEnzymeElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &FindValidTypeIIsEnzymeInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _FindValidTypeIIsEnzymeRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &FindValidTypeIIsEnzymeInput{},
			Out: &FindValidTypeIIsEnzymeOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type FindValidTypeIIsEnzymeElement struct {
	inject.CheckedRunner
}

type FindValidTypeIIsEnzymeInput struct {
	ApprovedEnzymes []string
	Constructname   string
	Seqsinorder     []string
	Vector          string
}

type FindValidTypeIIsEnzymeOutput struct {
	BackupEnzymeNames    []string
	BackupEnzymes        []wtype.RestrictionEnzyme
	ValidApprovedEnzymes []string
}

type FindValidTypeIIsEnzymeSOutput struct {
	Data struct {
		BackupEnzymeNames    []string
		BackupEnzymes        []wtype.RestrictionEnzyme
		ValidApprovedEnzymes []string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "FindValidTypeIIsEnzyme",
		Constructor: FindValidTypeIIsEnzymeNew,
		Desc: component.ComponentDesc{
			Desc: "The protocol is intended to return a list of valid TypeIIs enzymes which are compatible with a set of sequences and optionally the vector sequence.\nA list of ApprovedEnzymes enzymes can be specified. If no enzyme\nfrom the list is feasible to use (i.e. due to the presence of existing restriction sites in a part)\nall typeIIs enzymes will be screened to find feasible backup options.\n",
			Path: "src/github.com/antha-lang/elements/an/Data/DNA/TypeIISAssembly_design/FindValidTypeIIsEnzyme.an",
			Params: []component.ParamDesc{
				{Name: "ApprovedEnzymes", Desc: "list preferred enzyme names\n", Kind: "Parameters"},
				{Name: "Constructname", Desc: "", Kind: "Parameters"},
				{Name: "Seqsinorder", Desc: "list sequences as strings\n", Kind: "Parameters"},
				{Name: "Vector", Desc: "optionally specify the vector sequence.\n", Kind: "Parameters"},
				{Name: "BackupEnzymeNames", Desc: "", Kind: "Data"},
				{Name: "BackupEnzymes", Desc: "", Kind: "Data"},
				{Name: "ValidApprovedEnzymes", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
