// Protocol for creating a DNASequence from a sequence file format. // Supported format: .fasta
package lib

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"path/filepath"
)

// Input parameters for this protocol

//Only supported file format: .fasta
//Optional for the user to override specification of SequenceType using the name of the sequence specified in the file as the key. Each sequence name should be assigned to a SequenceType of the following list: Plasmid, Linear, SingleStranded. If no entry for a sequence is specified the value in the sequence file is used. Alternatively a "default" value can be specified which will apply to all sequences with no entry.
//If true, sequence is searched for ORF's

// Data which is returned from this protocol

//Return DNA sequence as type wtype.DNASequence
//Status for user
//Warnings for user

// Physical inputs to this protocol

// Physical outputs from this protocol

func _ImportMultipleDNASequenceRequirements() {

}

// Actions to perform before protocol itself
func _ImportMultipleDNASequenceSetup(_ctx context.Context, _input *ImportMultipleDNASequenceInput) {
}

// Core process of the protocol: steps to be performed for each input
func _ImportMultipleDNASequenceSteps(_ctx context.Context, _input *ImportMultipleDNASequenceInput, _output *ImportMultipleDNASequenceOutput) {
	//Use function DNAFileToDNASequence to read in file. The function determines
	//which file type was given as input and returns the DNA sequence as type wtype.DNAsequence

	if filepath.Ext(_input.SequenceFile.Name) != ".fasta" {
		execute.Errorf(_ctx, "The file format of %s is not supported. Please use file of format .fasta for this element.", _input.SequenceFile.Name)
	}

	seqs, err := parser.DNAFileToDNASequence(_input.SequenceFile)

	if err != nil {
		execute.Errorf(_ctx, "The file %s could not be imported. Error: %s ", _input.SequenceFile.Name, err.Error())
	}

	if err == nil {
		_output.DNA = seqs
		l := len(_output.DNA) - 1

		if _input.CheckForORFs {
			for i := 0; i <= l; i++ {
				//Finds all ORFs in imported DNA sequence
				orfs := sequences.FindallORFs(_output.DNA[i].Seq)
				if len(_output.DNA[i].Features) == 0 {
					features := sequences.ORFs2Features(orfs)
					_output.DNA[i] = wtype.Annotate(_output.DNA[i], features)
				}
			}
		}

		if len(_input.OverrideSequenceType) > 0 {

			for i := 0; i <= l; i++ {
				name := _output.DNA[i].Nm
				if seqType, found := _input.OverrideSequenceType[name]; found {
					if seqType == "Plasmid" {
						_output.DNA[i].Plasmid = true
						_output.DNA[i].Singlestranded = false
					} else if seqType == "Linear" {
						_output.DNA[i].Plasmid = false
						_output.DNA[i].Singlestranded = false
					} else if seqType == "SingleStranded" {
						_output.DNA[i].Plasmid = false
						_output.DNA[i].Singlestranded = true
					} else {
						execute.Errorf(_ctx, "SequenceType %s is unknown for %s. Please specify as Plasmid, Linear or SingleStranded", seqType, name)
					}
				} else if !found {
					if seqType, found := _input.OverrideSequenceType["default"]; found {
						if seqType == "Plasmid" {
							_output.DNA[i].Plasmid = true
							_output.DNA[i].Singlestranded = false
						} else if seqType == "Linear" {
							_output.DNA[i].Plasmid = false
							_output.DNA[i].Singlestranded = false
						} else if seqType == "SingleStranded" {
							_output.DNA[i].Plasmid = false
							_output.DNA[i].Singlestranded = true
						} else {
							err = fmt.Errorf("SequenceType %s is unknown. Please specify type as Plasmid, Linear or SingleStranded.", seqType)
						}
					} else {
						fmt.Printf("For %s no SequenceType was found in OverrideSequenceType. No default value was declared so %s's sequenceType was assigned from file. Please check if %s is spelled correctly in map, if SequenceType should be reassigned. ", name, name, name)
					}

				}
			}
		}

		_output.Status = fmt.Sprintln(
			text.Print("DNA_Seq: ", _output.DNA),
		)

		_output.Warnings = err

	}
}

// Actions to perform after steps block to analyze data
func _ImportMultipleDNASequenceAnalysis(_ctx context.Context, _input *ImportMultipleDNASequenceInput, _output *ImportMultipleDNASequenceOutput) {

}

func _ImportMultipleDNASequenceValidation(_ctx context.Context, _input *ImportMultipleDNASequenceInput, _output *ImportMultipleDNASequenceOutput) {

}
func _ImportMultipleDNASequenceRun(_ctx context.Context, input *ImportMultipleDNASequenceInput) *ImportMultipleDNASequenceOutput {
	output := &ImportMultipleDNASequenceOutput{}
	_ImportMultipleDNASequenceSetup(_ctx, input)
	_ImportMultipleDNASequenceSteps(_ctx, input, output)
	_ImportMultipleDNASequenceAnalysis(_ctx, input, output)
	_ImportMultipleDNASequenceValidation(_ctx, input, output)
	return output
}

func ImportMultipleDNASequenceRunSteps(_ctx context.Context, input *ImportMultipleDNASequenceInput) *ImportMultipleDNASequenceSOutput {
	soutput := &ImportMultipleDNASequenceSOutput{}
	output := _ImportMultipleDNASequenceRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ImportMultipleDNASequenceNew() interface{} {
	return &ImportMultipleDNASequenceElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ImportMultipleDNASequenceInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ImportMultipleDNASequenceRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ImportMultipleDNASequenceInput{},
			Out: &ImportMultipleDNASequenceOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type ImportMultipleDNASequenceElement struct {
	inject.CheckedRunner
}

type ImportMultipleDNASequenceInput struct {
	CheckForORFs         bool
	OverrideSequenceType map[string]string
	SequenceFile         wtype.File
}

type ImportMultipleDNASequenceOutput struct {
	DNA      []wtype.DNASequence
	Status   string
	Warnings error
}

type ImportMultipleDNASequenceSOutput struct {
	Data struct {
		DNA      []wtype.DNASequence
		Status   string
		Warnings error
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "ImportMultipleDNASequence",
		Constructor: ImportMultipleDNASequenceNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol for creating a DNASequence from a sequence file format. // Supported format: .fasta\n",
			Path: "src/github.com/antha-lang/elements/an/Data/DNAImport/ImportMultipleDNASequence.an",
			Params: []component.ParamDesc{
				{Name: "CheckForORFs", Desc: "If true, sequence is searched for ORF's\n", Kind: "Parameters"},
				{Name: "OverrideSequenceType", Desc: "Optional for the user to override specification of SequenceType using the name of the sequence specified in the file as the key. Each sequence name should be assigned to a SequenceType of the following list: Plasmid, Linear, SingleStranded. If no entry for a sequence is specified the value in the sequence file is used. Alternatively a \"default\" value can be specified which will apply to all sequences with no entry.\n", Kind: "Parameters"},
				{Name: "SequenceFile", Desc: "Only supported file format: .fasta\n", Kind: "Parameters"},
				{Name: "DNA", Desc: "Return DNA sequence as type wtype.DNASequence\n", Kind: "Data"},
				{Name: "Status", Desc: "Status for user\n", Kind: "Data"},
				{Name: "Warnings", Desc: "Warnings for user\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
