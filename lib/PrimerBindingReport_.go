// This element will produce a csv report of primer properties and binding sites to a list of input sequences
package lib

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/export"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/oligos"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"path/filepath"
	"sort"
	"strconv"
)

// Input parameters for this protocol

// Data which is returned from this protocol

// Physical inputs to this protocol

// Physical outputs from this protocol

func _PrimerBindingReportRequirements() {

}

// Actions to perform before protocol itself
func _PrimerBindingReportSetup(_ctx context.Context, _input *PrimerBindingReportInput) {

}

// order outputs
func mapToSlice(seqsMap map[string]wtype.DNASequence) []wtype.DNASequence {
	var names []string
	nameToSeq := make(map[string]wtype.DNASequence)

	for _, seq := range seqsMap {

		if seq, found := nameToSeq[seq.Name()]; !found {
			nameToSeq[seq.Name()] = seq
			names = append(names, seq.Name())
		} else {
			originalNm := seq.Name()
			counter := 2
			for {
				seq.Nm = originalNm + strconv.Itoa(counter)
				if seq, found := nameToSeq[seq.Name()]; !found {
					nameToSeq[seq.Name()] = seq
					names = append(names, seq.Name())
					break
				} else {
					counter++
				}
			}
			break
		}
	}

	sort.Strings(names)

	var seqsSlice []wtype.DNASequence

	for _, nm := range names {
		seqsSlice = append(seqsSlice, nameToSeq[nm])
	}

	return seqsSlice
}

// Core process of the protocol: steps to be performed for each input
func _PrimerBindingReportSteps(_ctx context.Context, _input *PrimerBindingReportInput, _output *PrimerBindingReportOutput) {

	//seqsSlice := mapToSlice(Sequences)

	// check each sequence for binding to other sequences:

	var report [][]string

	var header []string = []string{"Sequence", "Primer", "Primer Sequence", "Binding Sites", "Melting Temp", "GC Content"}

	report = append(report, header)

	for _, sequence := range _input.Sequences {

		for _, primer := range _input.AllPrimers {

			bindingsites := oligos.CheckNonSpecificBinding(sequence, primer.DNASequence)

			output := []string{sequence.Nm, primer.Nm, primer.Sequence(), fmt.Sprint(bindingsites), primer.MeltingTemp.ToString(), fmt.Sprint(primer.GCContent*100) + "%"}
			report = append(report, output)

		}

	}

	if filepath.Ext(_input.OutPutFileName) != "csv" {
		_input.OutPutFileName = _input.OutPutFileName + ".csv"
	}

	primerBindingReport, err := export.CSV(report, _input.OutPutFileName)

	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	_output.PrimerBindingReport = primerBindingReport

}

// Actions to perform after steps block to analyze data
func _PrimerBindingReportAnalysis(_ctx context.Context, _input *PrimerBindingReportInput, _output *PrimerBindingReportOutput) {

}
func _PrimerBindingReportValidation(_ctx context.Context, _input *PrimerBindingReportInput, _output *PrimerBindingReportOutput) {

}
func _PrimerBindingReportRun(_ctx context.Context, input *PrimerBindingReportInput) *PrimerBindingReportOutput {
	output := &PrimerBindingReportOutput{}
	_PrimerBindingReportSetup(_ctx, input)
	_PrimerBindingReportSteps(_ctx, input, output)
	_PrimerBindingReportAnalysis(_ctx, input, output)
	_PrimerBindingReportValidation(_ctx, input, output)
	return output
}

func PrimerBindingReportRunSteps(_ctx context.Context, input *PrimerBindingReportInput) *PrimerBindingReportSOutput {
	soutput := &PrimerBindingReportSOutput{}
	output := _PrimerBindingReportRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PrimerBindingReportNew() interface{} {
	return &PrimerBindingReportElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PrimerBindingReportInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PrimerBindingReportRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PrimerBindingReportInput{},
			Out: &PrimerBindingReportOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type PrimerBindingReportElement struct {
	inject.CheckedRunner
}

type PrimerBindingReportInput struct {
	AllPrimers     map[string]oligos.Primer
	OutPutFileName string
	Sequences      map[string]wtype.DNASequence
}

type PrimerBindingReportOutput struct {
	PrimerBindingReport wtype.File
}

type PrimerBindingReportSOutput struct {
	Data struct {
		PrimerBindingReport wtype.File
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PrimerBindingReport",
		Constructor: PrimerBindingReportNew,
		Desc: component.ComponentDesc{
			Desc: "This element will produce a csv report of primer properties and binding sites to a list of input sequences\n",
			Path: "src/github.com/antha-lang/elements/an/Data/DNA/PrimerDesign/PrimerBindingReport/PrimerBindingReport.an",
			Params: []component.ParamDesc{
				{Name: "AllPrimers", Desc: "", Kind: "Parameters"},
				{Name: "OutPutFileName", Desc: "", Kind: "Parameters"},
				{Name: "Sequences", Desc: "", Kind: "Parameters"},
				{Name: "PrimerBindingReport", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
