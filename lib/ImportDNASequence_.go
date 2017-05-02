// Protocol for creating a DNASequence from a sequence file format. // Supported formats: .gdx .fasta .gb
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
)

// Input parameters for this protocol

//Supported file formats formats: .gdx .fasta .gb
// optional parameter to rename the dna sequence. If left empty the name specified in the file is used
//Optional specification if DNA is of type "Plasmid, Linear or SingleStranded". If left empty the SequenceType is assigned from the file.
//If true, sequence is searched for ORF's

// Data which is returned from this protocol

//Return DNA sequence as type wtype.DNASequence
//Status for user
//Warnings for user

// Physical inputs to this protocol

// Physical outputs from this protocol

func _ImportDNASequenceRequirements() {

}

// Actions to perform before protocol itself
func _ImportDNASequenceSetup(_ctx context.Context, _input *ImportDNASequenceInput) {
}

// Core process of the protocol: steps to be performed for each input
func _ImportDNASequenceSteps(_ctx context.Context, _input *ImportDNASequenceInput, _output *ImportDNASequenceOutput) {
	//Use function DNAFileToDNASequence to read in file. The function determines
	//which file type was given as input and returns the DNA sequence as type wtype.DNAsequence
	seqs, err := parser.DNAFileToDNASequence(_input.SequenceFile)
	if err != nil {
		execute.Errorf(_ctx, "The file could not be imported. Please check if file format supported or if file empty")
	}
	if len(seqs) == 1 && err == nil {
		_output.DNA = seqs[0]

		if _input.RenameSequence != "" {
			_output.DNA.Nm = _input.RenameSequence
		}

		//Check if input "Plasmid, Linear, SingleStranded" matches reader file format. Outputs error to user.
		if _input.OverrideSequenceType == "" {
			return
		} else if _input.OverrideSequenceType != "" {
			if _input.OverrideSequenceType == "Plasmid" {
				if _output.DNA.Plasmid == false {
					fmt.Println("Warning: Sequence not specified as 'Plasmid' in file")
					_output.DNA.Plasmid = true
					_output.DNA.Singlestranded = false
				}
			} else if _input.OverrideSequenceType == "SingleStranded" {
				if _output.DNA.Singlestranded == false {
					fmt.Println("Warning: Sequence not specified as 'SingleStranded' in file")
					_output.DNA.Plasmid = false
					_output.DNA.Singlestranded = true
				}
			} else if _input.OverrideSequenceType == "Linear" {
				if _output.DNA.Singlestranded == true || _output.DNA.Plasmid == true {
					fmt.Println("Warning: Sequence not specified as 'Linear' in file")
					_output.DNA.Plasmid = false
					_output.DNA.Singlestranded = false
				}
			} else {
				execute.Errorf(_ctx, "Unknown DNA type specification. Please use Plasmid, SingleStranded or Linear as SequenceType")
			}
		}

	} else {
		execute.Errorf(_ctx, "Multiple Sequences are not supported. Please check format.")
	}

	//Finds all ORFs in imported DNA sequence
	if _input.CheckForORFs {
		orfs := sequences.FindallORFs(_output.DNA.Seq)

		if len(_output.DNA.Features) == 0 {
			features := sequences.ORFs2Features(orfs)

			_output.DNA = wtype.Annotate(_output.DNA, features)
			_output.Status = fmt.Sprintln(text.Print("ORFs: ", orfs))

		}
	}
	_output.Status = fmt.Sprintln(text.Print("DNA_Seq: ", _output.DNA)) + _output.Status
	_output.Warnings = err

}

// Actions to perform after steps block to analyze data
func _ImportDNASequenceAnalysis(_ctx context.Context, _input *ImportDNASequenceInput, _output *ImportDNASequenceOutput) {

}

func _ImportDNASequenceValidation(_ctx context.Context, _input *ImportDNASequenceInput, _output *ImportDNASequenceOutput) {

}
func _ImportDNASequenceRun(_ctx context.Context, input *ImportDNASequenceInput) *ImportDNASequenceOutput {
	output := &ImportDNASequenceOutput{}
	_ImportDNASequenceSetup(_ctx, input)
	_ImportDNASequenceSteps(_ctx, input, output)
	_ImportDNASequenceAnalysis(_ctx, input, output)
	_ImportDNASequenceValidation(_ctx, input, output)
	return output
}

func ImportDNASequenceRunSteps(_ctx context.Context, input *ImportDNASequenceInput) *ImportDNASequenceSOutput {
	soutput := &ImportDNASequenceSOutput{}
	output := _ImportDNASequenceRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ImportDNASequenceNew() interface{} {
	return &ImportDNASequenceElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ImportDNASequenceInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ImportDNASequenceRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ImportDNASequenceInput{},
			Out: &ImportDNASequenceOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type ImportDNASequenceElement struct {
	inject.CheckedRunner
}

type ImportDNASequenceInput struct {
	CheckForORFs         bool
	OverrideSequenceType string
	RenameSequence       string
	SequenceFile         wtype.File
}

type ImportDNASequenceOutput struct {
	DNA      wtype.DNASequence
	Status   string
	Warnings error
}

type ImportDNASequenceSOutput struct {
	Data struct {
		DNA      wtype.DNASequence
		Status   string
		Warnings error
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "ImportDNASequence",
		Constructor: ImportDNASequenceNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol for creating a DNASequence from a sequence file format. // Supported formats: .gdx .fasta .gb\n",
			Path: "src/github.com/antha-lang/elements/an/Data/DNAImport/ImportDNASequence.an",
			Params: []component.ParamDesc{
				{Name: "CheckForORFs", Desc: "If true, sequence is searched for ORF's\n", Kind: "Parameters"},
				{Name: "OverrideSequenceType", Desc: "Optional specification if DNA is of type \"Plasmid, Linear or SingleStranded\". If left empty the SequenceType is assigned from the file.\n", Kind: "Parameters"},
				{Name: "RenameSequence", Desc: "optional parameter to rename the dna sequence. If left empty the name specified in the file is used\n", Kind: "Parameters"},
				{Name: "SequenceFile", Desc: "Supported file formats formats: .gdx .fasta .gb\n", Kind: "Parameters"},
				{Name: "DNA", Desc: "Return DNA sequence as type wtype.DNASequence\n", Kind: "Data"},
				{Name: "Status", Desc: "Status for user\n", Kind: "Data"},
				{Name: "Warnings", Desc: "Warnings for user\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
