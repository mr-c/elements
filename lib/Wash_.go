// Protocol Wash performs a wash step
package lib

import

// Place golang packages to import here
(
	"context"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Parameters to this protocol

// Output data of this protocol

// Physical inputs to this protocol

// Physical outputs to this protocol

// Conditions to run on startup
func _WashSetup(_ctx context.Context, _input *WashInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _WashSteps(_ctx context.Context, _input *WashInput, _output *WashOutput) {

	//setup some variables for use during the protocol
	var samples []*wtype.LHComponent
	var err error
	var mixPolicy wtype.PolicyName

	WastePlate := execute.NewPlate(_ctx, "reservoir")
	var wastedBufferFromSamples []*wtype.LHComponent

	//var samplesSlice []string = MakeSlice(SamplesToWash)

	_input.WashBuffer.SetVolume(wunit.NewVolume(4000, "ul"))
	_input.WashPlate.Welltype.Add(_input.WashBuffer)

	//determine if WashMixing is selected, and if so, assign LHPolicy to PostMix
	if _input.WashMixing {
		mixPolicy = "PostMix"
	} else {
		mixPolicy = "DoNotMix"
	}

	//get plate info from input plate and was
	var samplesWells []string = _input.SamplesPlate.AllWellPositions(wtype.BYCOLUMN)

	//loop through up to desired number of washes
	for j := 0; j < _input.NumberOfWashes; j++ {

		//range through the input samples and add wash buffer to each
		for i := range _input.SamplesToWash {

			//sample washbuffer at specified volume
			washBufferSample := mixer.Sample(_input.WashBuffer, _input.WashVolume)

			//assign LHpolicy to wash sample (PostMix or NeedToMix)
			washBufferSample.Type, err = wtype.LiquidTypeFromString(mixPolicy)

			if err != nil {
				execute.Errorf(_ctx, "LHPolicy %s invalid: %s", mixPolicy, err.Error())
			}

			//update position to correspond to counter
			position := samplesWells[i]

			//add wash buffer into samples
			washSamples := execute.MixNamed(_ctx, _input.SamplesPlate.Type, position, "WashPlate", washBufferSample)

			//add wash solutions to slice for subsequent removal
			samples = append(samples, washSamples)

		}

		//range through slice of washe solutions from previous loop to remove
		for k := range samples {

			//determine volume to be removed by adding WashVolume and adding 20ul excess
			newWashSolutionVolume := wunit.AddVolumes([]wunit.Volume{_input.WashVolume, wunit.NewVolume(0, "ul")})

			//position := samplesWells[k]

			//remove wash buffer at updated volume
			washBufferRemoval := mixer.Sample(samples[k], newWashSolutionVolume)

			wastedBufferFromSamples = append(wastedBufferFromSamples, washBufferRemoval)

		}

		//transfer used wash buffer to WastePlate
		wasteDisposal := execute.MixNamed(_ctx, WastePlate.Type, "A1", "WastePlate", wastedBufferFromSamples...)

		//setup slice to add wasted wash buffer
		var wastedBuffer []*wtype.LHComponent

		//add wasted wash buffer to slice
		wastedBuffer = append(wastedBuffer, wasteDisposal)

		//update outputs
		_output.WasteBuffer = wastedBuffer
	}

	//update outputs
	_output.Errors = err
	_output.ProcessedSamples = samples

}

func MakeSlice(v []*wtype.LHComponent) (Names []string) {
	for _, x := range v {
		Names = append(Names, x.CName)
	}
	return Names
}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _WashAnalysis(_ctx context.Context, _input *WashInput, _output *WashOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _WashValidation(_ctx context.Context, _input *WashInput, _output *WashOutput) {

}
func _WashRun(_ctx context.Context, input *WashInput) *WashOutput {
	output := &WashOutput{}
	_WashSetup(_ctx, input)
	_WashSteps(_ctx, input, output)
	_WashAnalysis(_ctx, input, output)
	_WashValidation(_ctx, input, output)
	return output
}

func WashRunSteps(_ctx context.Context, input *WashInput) *WashSOutput {
	soutput := &WashSOutput{}
	output := _WashRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func WashNew() interface{} {
	return &WashElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &WashInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _WashRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &WashInput{},
			Out: &WashOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type WashElement struct {
	inject.CheckedRunner
}

type WashInput struct {
	IncubationBetweenWash bool
	IncubationTemperature wunit.Temperature
	IncubationTime        wunit.Time
	NumberOfWashes        int
	SamplesPlate          *wtype.LHPlate
	SamplesToWash         []*wtype.LHComponent
	WashBuffer            *wtype.LHComponent
	WashMixing            bool
	WashPlate             *wtype.LHPlate
	WashVolume            wunit.Volume
}

type WashOutput struct {
	Errors           error
	ProcessedSamples []*wtype.LHComponent
	WasteBuffer      []*wtype.LHComponent
}

type WashSOutput struct {
	Data struct {
		Errors error
	}
	Outputs struct {
		ProcessedSamples []*wtype.LHComponent
		WasteBuffer      []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Wash",
		Constructor: WashNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol Wash performs a wash step\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/Wash/Wash.an",
			Params: []component.ParamDesc{
				{Name: "IncubationBetweenWash", Desc: "", Kind: "Parameters"},
				{Name: "IncubationTemperature", Desc: "", Kind: "Parameters"},
				{Name: "IncubationTime", Desc: "", Kind: "Parameters"},
				{Name: "NumberOfWashes", Desc: "", Kind: "Parameters"},
				{Name: "SamplesPlate", Desc: "", Kind: "Inputs"},
				{Name: "SamplesToWash", Desc: "", Kind: "Inputs"},
				{Name: "WashBuffer", Desc: "", Kind: "Inputs"},
				{Name: "WashMixing", Desc: "", Kind: "Parameters"},
				{Name: "WashPlate", Desc: "", Kind: "Inputs"},
				{Name: "WashVolume", Desc: "", Kind: "Parameters"},
				{Name: "Errors", Desc: "", Kind: "Data"},
				{Name: "ProcessedSamples", Desc: "", Kind: "Outputs"},
				{Name: "WasteBuffer", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
