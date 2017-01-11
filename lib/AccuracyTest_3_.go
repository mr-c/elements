package lib

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/doe"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/download"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"strconv"
)

// Input parameters for this protocol (data)

// select this if getting the image from a URL
// enter URL link to the image file here if applicable

// optional parameter allowing pipetting to resume on partially filled plate

// Data which is returned from this protocol, and data types

//[]string //map[string]string

//NeatSamplewells []string

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _AccuracyTest_3Requirements() {
}

// Conditions to run on startup
func _AccuracyTest_3Setup(_ctx context.Context, _input *AccuracyTest_3Input) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _AccuracyTest_3Steps(_ctx context.Context, _input *AccuracyTest_3Input, _output *AccuracyTest_3Output) {

	// if dilution factor not set dilute by 10x
	if _input.DilutionFactor == 0 {
		_input.DilutionFactor = 10
	}
	minVolume := wunit.NewVolume(0.5, "ul")

	// declare some global variables for use later
	var rotate = false
	var autorotate = true
	var wellpositionarray = make([]string, 0)
	_output.Runtowelllocationmap = make(map[string]string)
	_output.Blankwells = make([]string, 0)
	counter := _input.WellsUsed

	var platenum = 1
	var runs = make([]doe.Run, 1)
	var newruns = make([]doe.Run, 0)
	var err error
	_output.Errors = make([]error, 0)
	// work out plate layout based on picture or just in order

	if _input.Printasimage {

		// if image is from url, download
		if _input.UseURL {
			err := download.File(_input.URL, _input.Imagefilename)
			if err != nil {
				execute.Errorf(_ctx, err.Error())
			}
		}

		chosencolourpalette := image.AvailablePalettes()["Palette1"]
		positiontocolourmap, _, _ := image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, &chosencolourpalette, rotate, autorotate)

		//Runtowelllocationmap = make([]string,0)

		for location, colour := range positiontocolourmap {
			R, G, B, A := colour.RGBA()

			if uint8(R) == 242 && uint8(G) == 243 && uint8(B) == 242 && uint8(A) == 255 {
				continue
			} else {
				wellpositionarray = append(wellpositionarray, location)
			}
		}

	} else {

		wellpositionarray = _input.OutPlate.AllWellPositions(wtype.BYCOLUMN)

	}
	reactions := make([]*wtype.LHComponent, 0)

	// use first policy as reference to ensure consistent range through map values
	referencepolicy, found := liquidhandling.GetPolicyByName(_input.LHPolicy)
	if found == false {
		execute.Errorf(_ctx, "policy "+_input.LHPolicy+" not found")
		_output.Errors = append(_output.Errors, fmt.Errorf("policy ", _input.LHPolicy, " not found"))
	}

	referencekeys := make([]string, 0)
	for key := range referencepolicy {
		referencekeys = append(referencekeys, key)
	}

	for l := 0; l < len(_input.TestSolVolumes); l++ {
		for k := 0; k < len(_input.TestSols); k++ {
			for j := 0; j < _input.NumberofReplicates; j++ {
				for i := 0; i < len(runs); i++ {

					var run doe.Run

					if counter == ((_input.OutPlate.WlsX * _input.OutPlate.WlsY) + _input.NumberofBlanks) {
						fmt.Println("plate full, counter = ", counter)
						platenum++
						counter = 0
					}

					var eachreaction []*wtype.LHComponent
					var solution *wtype.LHComponent

					if _input.PipetteOnebyOne {
						eachreaction = make([]*wtype.LHComponent, 0)
						//eachdilution = make([]*wtype.LHComponent, 0)
					}
					// keep default policy for diluent

					// diluent first

					// change lhpolicy if desired
					if _input.UseLHPolicyDoeforDiluent {
						_input.Diluent.Type, err = wtype.LiquidTypeFromString(_input.LHPolicy)
						if err != nil {
							_output.Errors = append(_output.Errors, err)
						}
					}

					var bufferSample *wtype.LHComponent
					var Dilution *wtype.LHComponent

					if _input.TestSolVolumes[l].LessThan(minVolume) {

						// add diluent to dilution plate ready for dilution
						dilutedSampleBuffer := mixer.Sample(_input.Diluent, wunit.SubtractVolumes(_input.TotalVolume, []wunit.Volume{wunit.MultiplyVolume(_input.TestSolVolumes[l], _input.DilutionFactor)}))
						Dilution = execute.MixNamed(_ctx, _input.OutPlate.Type, wellpositionarray[counter], fmt.Sprint("DilutionPlate", platenum), dilutedSampleBuffer)

						// add same volume to destination plate ready for dilutedsolution
						bufferSample = mixer.Sample(_input.Diluent, wunit.SubtractVolumes(_input.TotalVolume, []wunit.Volume{wunit.MultiplyVolume(_input.TestSolVolumes[l], _input.DilutionFactor)}))

					} else {

						bufferSample = mixer.Sample(_input.Diluent, wunit.SubtractVolumes(_input.TotalVolume, []wunit.Volume{_input.TestSolVolumes[l]})) //SampleForTotalVolume(Diluent, TotalVolume)

					}

					if _input.PipetteOnebyOne {
						eachreaction = append(eachreaction, bufferSample)
					} else {

						solution = execute.MixTo(_ctx, _input.OutPlate.Type, wellpositionarray[counter], platenum, bufferSample)
					}
					// now test sample

					// change liquid class
					if _input.UseLiquidPolicyForTestSolutions && _input.LHPolicy != "" {
						_input.TestSols[k].Type, err = wtype.LiquidTypeFromString(_input.LHPolicy)
						if err != nil {
							_output.Errors = append(_output.Errors, err)
						}
					}

					if _input.TestSolVolumes[l].GreaterThan(minVolume) {

						//sample
						testSample := mixer.Sample(_input.TestSols[k], _input.TestSolVolumes[l])

						if _input.PipetteOnebyOne {
							eachreaction = append(eachreaction, testSample)
							solution = execute.MixTo(_ctx, _input.OutPlate.Type, wellpositionarray[counter], platenum, eachreaction...)
						} else {
							// pipette out
							solution = execute.Mix(_ctx, solution, testSample)
						}

					} else if _input.TestSolVolumes[l].GreaterThan(wunit.NewVolume(0.0, "ul")) && _input.TestSolVolumes[l].LessThan(minVolume) {

						//sample
						dilutionSample := mixer.Sample(_input.TestSols[k], wunit.MultiplyVolume(_input.TestSolVolumes[l], _input.DilutionFactor))
						Dilution = execute.MixNamed(_ctx, _input.OutPlate.Type, wellpositionarray[counter], fmt.Sprint("DilutionPlate", platenum), dilutionSample)

						testSample := mixer.Sample(Dilution, wunit.MultiplyVolume(_input.TestSolVolumes[l], _input.DilutionFactor))

						if _input.PipetteOnebyOne {
							eachreaction = append(eachreaction, testSample)
							solution = execute.MixTo(_ctx, _input.OutPlate.Type, wellpositionarray[counter], platenum, eachreaction...)

						} else {
							// pipette out

							solution = execute.Mix(_ctx, solution, testSample)
						}

					}

					// get annotation info
					doerun := wtype.LiquidTypeName(_input.TestSols[k].Type)

					volume := _input.TestSolVolumes[l].ToString()

					solutionname := _input.TestSols[k].CName

					// add Solution Name
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Solution", solutionname)

					// add Volume
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Volume", volume)

					// add Replicate
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Replicate", strconv.Itoa(j+1))

					// full description
					description := volume + "_" + solutionname + "_replicate" + strconv.Itoa(j+1) + "_platenum" + strconv.Itoa(platenum)

					// add run to well position lookup table
					_output.Runtowelllocationmap[doerun+"_"+description] = wellpositionarray[counter]

					// add additional info for each run
					fmt.Println("len(runs)", len(runs), "counter", counter, "len(wellpositionarray)", len(wellpositionarray))

					run = doe.AddAdditionalHeaderandValue(run, "Additional", "PlateNumber", strconv.Itoa(platenum))
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Location", wellpositionarray[counter])

					// add setpoint printout to double check correct match up:
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "LHPolicy", doerun)

					// add plate info:
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Plate Type", _input.OutPlate.Type)

					// add plate ZStart:
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Plate WellZStart", _input.OutPlate.WellZStart)

					// add plate Height:
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Plate Height", _input.OutPlate.Height)

					// other plate offsets:
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Plate WellXOffset", _input.OutPlate.WellXOffset)

					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Plate WellYOffset", _input.OutPlate.WellYOffset)

					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Plate WellXStart", _input.OutPlate.WellXStart)

					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Plate WellYStart", _input.OutPlate.WellYStart)

					// add LHPolicy setpoint printout to double check correct match up:
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "LHPolicy", doerun)

					// print out LHPolicy info
					policy, found := liquidhandling.GetPolicyByName(doerun)
					if !found {
						execute.Errorf(_ctx, "policy "+doerun+" not found")
						_output.Errors = append(_output.Errors, fmt.Errorf("policy ", doerun, " not found"))
					}

					for _, key := range referencekeys {
						run = doe.AddAdditionalHeaderandValue(run, "Additional", "LHPolicy"+"_"+key, policy[key])
					}

					reactions = append(reactions, solution)
					newruns = append(newruns, run)

					counter++

				}

			}
		}
	}

	// export overall DOE design file showing all well locations for all conditions
	doe.XLSXFileFromRuns(newruns, _input.OutputFilename, _input.DXORJMP)

	// add blanks after

	for n := 0; n < platenum; n++ {
		for m := 0; m < _input.NumberofBlanks; m++ {

			// use defualt policy for blank

			bufferSample := mixer.Sample(_input.Diluent, _input.TotalVolume)
			//eachreaction = append(eachreaction,bufferSample)

			// add blanks to last column of plate
			well := wutil.NumToAlpha(_input.OutPlate.WlsY-m) + strconv.Itoa(_input.OutPlate.WlsX)

			reaction := execute.MixTo(_ctx, _input.OutPlate.Type, well, n+1, bufferSample)

			_output.Runtowelllocationmap["Blank"+strconv.Itoa(m+1)+" platenum"+strconv.Itoa(n+1)] = well

			_output.Blankwells = append(_output.Blankwells, well)

			reactions = append(reactions, reaction)
			counter++

		}

	}

	_output.Reactions = reactions
	_output.Runcount = len(_output.Reactions)
	_output.Pixelcount = len(wellpositionarray)
	_output.Runs = runs
	_output.Wellpositionarray = wellpositionarray
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AccuracyTest_3Analysis(_ctx context.Context, _input *AccuracyTest_3Input, _output *AccuracyTest_3Output) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AccuracyTest_3Validation(_ctx context.Context, _input *AccuracyTest_3Input, _output *AccuracyTest_3Output) {
}
func _AccuracyTest_3Run(_ctx context.Context, input *AccuracyTest_3Input) *AccuracyTest_3Output {
	output := &AccuracyTest_3Output{}
	_AccuracyTest_3Setup(_ctx, input)
	_AccuracyTest_3Steps(_ctx, input, output)
	_AccuracyTest_3Analysis(_ctx, input, output)
	_AccuracyTest_3Validation(_ctx, input, output)
	return output
}

func AccuracyTest_3RunSteps(_ctx context.Context, input *AccuracyTest_3Input) *AccuracyTest_3SOutput {
	soutput := &AccuracyTest_3SOutput{}
	output := _AccuracyTest_3Run(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AccuracyTest_3New() interface{} {
	return &AccuracyTest_3Element{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AccuracyTest_3Input{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AccuracyTest_3Run(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AccuracyTest_3Input{},
			Out: &AccuracyTest_3Output{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type AccuracyTest_3Element struct {
	inject.CheckedRunner
}

type AccuracyTest_3Input struct {
	DXORJMP                         string
	Diluent                         *wtype.LHComponent
	DilutionFactor                  float64
	Imagefilename                   string
	LHPolicy                        string
	NumberofBlanks                  int
	NumberofReplicates              int
	OutPlate                        *wtype.LHPlate
	OutputFilename                  string
	PipetteOnebyOne                 bool
	Printasimage                    bool
	TestSolVolumes                  []wunit.Volume
	TestSols                        []*wtype.LHComponent
	TotalVolume                     wunit.Volume
	URL                             string
	UseLHPolicyDoeforDiluent        bool
	UseLiquidPolicyForTestSolutions bool
	UseURL                          bool
	WellsUsed                       int
}

type AccuracyTest_3Output struct {
	Blankwells           []string
	Errors               []error
	Pixelcount           int
	Reactions            []*wtype.LHComponent
	Runcount             int
	Runs                 []doe.Run
	Runtowelllocationmap map[string]string
	Wellpositionarray    []string
}

type AccuracyTest_3SOutput struct {
	Data struct {
		Blankwells           []string
		Errors               []error
		Pixelcount           int
		Runcount             int
		Runs                 []doe.Run
		Runtowelllocationmap map[string]string
		Wellpositionarray    []string
	}
	Outputs struct {
		Reactions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AccuracyTest_3",
		Constructor: AccuracyTest_3New,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/elements/an/Utility/AccuracyTest_3.an",
			Params: []component.ParamDesc{
				{Name: "DXORJMP", Desc: "", Kind: "Parameters"},
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "DilutionFactor", Desc: "", Kind: "Parameters"},
				{Name: "Imagefilename", Desc: "", Kind: "Parameters"},
				{Name: "LHPolicy", Desc: "", Kind: "Parameters"},
				{Name: "NumberofBlanks", Desc: "", Kind: "Parameters"},
				{Name: "NumberofReplicates", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "OutputFilename", Desc: "", Kind: "Parameters"},
				{Name: "PipetteOnebyOne", Desc: "", Kind: "Parameters"},
				{Name: "Printasimage", Desc: "", Kind: "Parameters"},
				{Name: "TestSolVolumes", Desc: "", Kind: "Parameters"},
				{Name: "TestSols", Desc: "", Kind: "Inputs"},
				{Name: "TotalVolume", Desc: "", Kind: "Parameters"},
				{Name: "URL", Desc: "enter URL link to the image file here if applicable\n", Kind: "Parameters"},
				{Name: "UseLHPolicyDoeforDiluent", Desc: "", Kind: "Parameters"},
				{Name: "UseLiquidPolicyForTestSolutions", Desc: "", Kind: "Parameters"},
				{Name: "UseURL", Desc: "select this if getting the image from a URL\n", Kind: "Parameters"},
				{Name: "WellsUsed", Desc: "optional parameter allowing pipetting to resume on partially filled plate\n", Kind: "Parameters"},
				{Name: "Blankwells", Desc: "", Kind: "Data"},
				{Name: "Errors", Desc: "", Kind: "Data"},
				{Name: "Pixelcount", Desc: "", Kind: "Data"},
				{Name: "Reactions", Desc: "", Kind: "Outputs"},
				{Name: "Runcount", Desc: "", Kind: "Data"},
				{Name: "Runs", Desc: "", Kind: "Data"},
				{Name: "Runtowelllocationmap", Desc: "[]string //map[string]string\n", Kind: "Data"},
				{Name: "Wellpositionarray", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
