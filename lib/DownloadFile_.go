// Downloads a file from a specied URL
package lib

import

// Place golang packages to import here
(
	"context"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/download"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Parameters to this protocol

// name of image file or if using URL use this field to set the desired filename
// enter URL link to the image file here if applicable

// Output data of this protocol

// Physical inputs to this protocol

// Physical outputs to this protocol

// Conditions to run on startup
func _DownloadFileSetup(_ctx context.Context, _input *DownloadFileInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _DownloadFileSteps(_ctx context.Context, _input *DownloadFileInput, _output *DownloadFileOutput) {

	_output.OutPutFile, _output.Error = download.File(_input.URL, _input.Imagefilename)
	if _output.Error != nil {
		execute.Errorf(_ctx, _output.Error.Error())
	}

}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _DownloadFileAnalysis(_ctx context.Context, _input *DownloadFileInput, _output *DownloadFileOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _DownloadFileValidation(_ctx context.Context, _input *DownloadFileInput, _output *DownloadFileOutput) {

}
func _DownloadFileRun(_ctx context.Context, input *DownloadFileInput) *DownloadFileOutput {
	output := &DownloadFileOutput{}
	_DownloadFileSetup(_ctx, input)
	_DownloadFileSteps(_ctx, input, output)
	_DownloadFileAnalysis(_ctx, input, output)
	_DownloadFileValidation(_ctx, input, output)
	return output
}

func DownloadFileRunSteps(_ctx context.Context, input *DownloadFileInput) *DownloadFileSOutput {
	soutput := &DownloadFileSOutput{}
	output := _DownloadFileRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func DownloadFileNew() interface{} {
	return &DownloadFileElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &DownloadFileInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _DownloadFileRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &DownloadFileInput{},
			Out: &DownloadFileOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type DownloadFileElement struct {
	inject.CheckedRunner
}

type DownloadFileInput struct {
	Imagefilename string
	URL           string
}

type DownloadFileOutput struct {
	Error      error
	OutPutFile wtype.File
}

type DownloadFileSOutput struct {
	Data struct {
		Error      error
		OutPutFile wtype.File
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "DownloadFile",
		Constructor: DownloadFileNew,
		Desc: component.ComponentDesc{
			Desc: "Downloads a file from a specied URL\n",
			Path: "src/github.com/antha-lang/elements/an/Data/DownloadFile/DownloadFile.an",
			Params: []component.ParamDesc{
				{Name: "Imagefilename", Desc: "name of image file or if using URL use this field to set the desired filename\n", Kind: "Parameters"},
				{Name: "URL", Desc: "enter URL link to the image file here if applicable\n", Kind: "Parameters"},
				{Name: "Error", Desc: "", Kind: "Data"},
				{Name: "OutPutFile", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
