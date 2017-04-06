// Adding given overhangs to a DNASequence objects
package lib

import

// Place golang packages to import here
(
	"context"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Parameters to this protocol

//String of DNA, if empty nothing is added
//String of DNA, if empty nothing is added

// Output data of this protocol

// Physical inputs to this protocol

// Physical outputs to this protocol

// Conditions to run on startup
func _EditSeqSetup(_ctx context.Context, _input *EditSeqInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _EditSeqSteps(_ctx context.Context, _input *EditSeqInput, _output *EditSeqOutput) {

	if _input.InputSeq.plasmid {
		_output.Warning = fmt.Errorf("Your sequence is circular, no ends to add to")
		execute.Errorf(_ctx, _output.Warning.Error())
	}

	_output.OutputSeq = _input.InputSeq.Dup()

	_output.OutputSeq.Append(_input.AddPrefix)
	_output.OutputSeq.PrePend(_input.AddSuffix)

	passed, illegals, wobble := sequences.Illegalnucleotides(_output.OutputSeq)

	if !passed && len(illegals) > 0 {
		_output.Warning = fmt.Errorf("DNA output sequence contains illegal characters %+v", illegals)
		execute.Errorf(_ctx, _output.Warning.Error())
	}

}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _EditSeqAnalysis(_ctx context.Context, _input *EditSeqInput, _output *EditSeqOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _EditSeqValidation(_ctx context.Context, _input *EditSeqInput, _output *EditSeqOutput) {

}
func _EditSeqRun(_ctx context.Context, input *EditSeqInput) *EditSeqOutput {
	output := &EditSeqOutput{}
	_EditSeqSetup(_ctx, input)
	_EditSeqSteps(_ctx, input, output)
	_EditSeqAnalysis(_ctx, input, output)
	_EditSeqValidation(_ctx, input, output)
	return output
}

func EditSeqRunSteps(_ctx context.Context, input *EditSeqInput) *EditSeqSOutput {
	soutput := &EditSeqSOutput{}
	output := _EditSeqRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func EditSeqNew() interface{} {
	return &EditSeqElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &EditSeqInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _EditSeqRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &EditSeqInput{},
			Out: &EditSeqOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type EditSeqElement struct {
	inject.CheckedRunner
}

type EditSeqInput struct {
	AddPrefix string
	AddSuffix string
	InputSeq  *wtype.DNASequence
}

type EditSeqOutput struct {
	OutputSeq *wtype.DNASequence
	Warning   error
}

type EditSeqSOutput struct {
	Data struct {
		OutputSeq *wtype.DNASequence
		Warning   error
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "EditSeq",
		Constructor: EditSeqNew,
		Desc: component.ComponentDesc{
			Desc: "Adding given overhangs to a DNASequence objects\n",
			Path: "src/github.com/antha-lang/elements/an/Data/DNA/EditSequence/EditSeq/EditSeq.an",
			Params: []component.ParamDesc{
				{Name: "AddPrefix", Desc: "String of DNA, if empty nothing is added\n", Kind: "Parameters"},
				{Name: "AddSuffix", Desc: "String of DNA, if empty nothing is added\n", Kind: "Parameters"},
				{Name: "InputSeq", Desc: "", Kind: "Parameters"},
				{Name: "OutputSeq", Desc: "", Kind: "Data"},
				{Name: "Warning", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
