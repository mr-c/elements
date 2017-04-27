//Append_Prepend_DNASequence will take in an array of DNA sequences and Append and/or Prepend extra base pairs to the DNA sequences.
//If Append or Prepend is left blank, no extra base pairs will be added.
//Modified Sequences are outputted in a FASTA file.
package lib

import

// Place golang packages to import here
(
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/export"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"path/filepath"
	"strings"
)

// Parameters to this protocol

//this will be used to name the output FASTA file and folder
//input DNA sequences
//string of DNA, if empty nothing is added
//string of DNA, if empty nothing is added

// Output data of this protocol

//output modified DNA sequences
//error messages reported back to the user
//output Fasta file

// Physical inputs to this protocol

// Physical outputs to this protocol

// Conditions to run on startup
func _Append_Prepend_DNASequenceSetup(_ctx context.Context, _input *Append_Prepend_DNASequenceInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _Append_Prepend_DNASequenceSteps(_ctx context.Context, _input *Append_Prepend_DNASequenceInput, _output *Append_Prepend_DNASequenceOutput) {

	//setup warnings slice to append errors to
	warnings := make([]string, 0)

	//range through the InputSequences
	for _, editedSequence := range _input.InputSequences {

		//check if the InputSequence is a plasmid, and return warning message if so
		if editedSequence.Plasmid {
			plasmidError := fmt.Errorf("Warning: The input DNA sequence %s is listed as a plasmid and should not be Appended/Prepended, please proceed with caution or fix", editedSequence.Nm)
			fmt.Println(plasmidError)
		}

		//check input sequences for illegal (non-nucleotide) characters and return error message if so
		passed, illegals, _ := sequences.Illegalnucleotides(editedSequence)

		if !passed {
			var newstatus = make([]string, 0)
			for _, illegal := range illegals {

				newstatus = append(newstatus, "part: "+editedSequence.Nm+" "+editedSequence.Seq+": contains illegalnucleotides:"+illegal.ToString())
			}
			warnings = append(warnings, strings.Join(newstatus, ""))
			fmt.Errorf(strings.Join(newstatus, ""))
		}

		//Append and Prepend the given additional bp to the input sequence
		editedSequence.Append(_input.AddSuffix)
		editedSequence.Prepend(_input.AddPrefix)

		//append modified sequences to the ModifiedSequences array
		_output.ModifiedSequences = append(_output.ModifiedSequences, editedSequence)
	}

	//add the ModifiedSequences to a FASTA file in new folder with ProjectName
	outputFile, _, err := export.FastaSerial(export.LOCAL, filepath.Join(_input.ProjectName, "AssemblyProduct"), _output.ModifiedSequences)
	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}
	_output.ModifiedSequenceFile = outputFile

	//add all warnings to the Warnings output and report to user
	_output.Warnings = fmt.Errorf(strings.Join(warnings, ";"))

}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _Append_Prepend_DNASequenceAnalysis(_ctx context.Context, _input *Append_Prepend_DNASequenceInput, _output *Append_Prepend_DNASequenceOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _Append_Prepend_DNASequenceValidation(_ctx context.Context, _input *Append_Prepend_DNASequenceInput, _output *Append_Prepend_DNASequenceOutput) {

}
func _Append_Prepend_DNASequenceRun(_ctx context.Context, input *Append_Prepend_DNASequenceInput) *Append_Prepend_DNASequenceOutput {
	output := &Append_Prepend_DNASequenceOutput{}
	_Append_Prepend_DNASequenceSetup(_ctx, input)
	_Append_Prepend_DNASequenceSteps(_ctx, input, output)
	_Append_Prepend_DNASequenceAnalysis(_ctx, input, output)
	_Append_Prepend_DNASequenceValidation(_ctx, input, output)
	return output
}

func Append_Prepend_DNASequenceRunSteps(_ctx context.Context, input *Append_Prepend_DNASequenceInput) *Append_Prepend_DNASequenceSOutput {
	soutput := &Append_Prepend_DNASequenceSOutput{}
	output := _Append_Prepend_DNASequenceRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Append_Prepend_DNASequenceNew() interface{} {
	return &Append_Prepend_DNASequenceElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Append_Prepend_DNASequenceInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Append_Prepend_DNASequenceRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Append_Prepend_DNASequenceInput{},
			Out: &Append_Prepend_DNASequenceOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type Append_Prepend_DNASequenceElement struct {
	inject.CheckedRunner
}

type Append_Prepend_DNASequenceInput struct {
	AddPrefix      string
	AddSuffix      string
	InputSequences []wtype.DNASequence
	ProjectName    string
}

type Append_Prepend_DNASequenceOutput struct {
	ModifiedSequenceFile wtype.File
	ModifiedSequences    []wtype.DNASequence
	Warnings             error
}

type Append_Prepend_DNASequenceSOutput struct {
	Data struct {
		ModifiedSequenceFile wtype.File
		ModifiedSequences    []wtype.DNASequence
		Warnings             error
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Append_Prepend_DNASequence",
		Constructor: Append_Prepend_DNASequenceNew,
		Desc: component.ComponentDesc{
			Desc: "Append_Prepend_DNASequence will take in an array of DNA sequences and Append and/or Prepend extra base pairs to the DNA sequences.\nIf Append or Prepend is left blank, no extra base pairs will be added.\nModified Sequences are outputted in a FASTA file.\n",
			Path: "src/github.com/antha-lang/elements/an/Data/DNA/EditSequence/Append_Prepend_DNASequence/Append_Prepend_DNASequence.an",
			Params: []component.ParamDesc{
				{Name: "AddPrefix", Desc: "string of DNA, if empty nothing is added\n", Kind: "Parameters"},
				{Name: "AddSuffix", Desc: "string of DNA, if empty nothing is added\n", Kind: "Parameters"},
				{Name: "InputSequences", Desc: "input DNA sequences\n", Kind: "Parameters"},
				{Name: "ProjectName", Desc: "this will be used to name the output FASTA file and folder\n", Kind: "Parameters"},
				{Name: "ModifiedSequenceFile", Desc: "output Fasta file\n", Kind: "Data"},
				{Name: "ModifiedSequences", Desc: "output modified DNA sequences\n", Kind: "Data"},
				{Name: "Warnings", Desc: "error messages reported back to the user\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
