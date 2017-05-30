// Perform accuracy test protocol using a series of concentrations as set points
package lib

import (
	"bytes"
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/doe"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	goimage "image"
	"strconv"
)

// Input parameters for this protocol (data)

// Total volunme per well

// Option to over ride the stock concentrations of the solutions being tested.

// Concentration Set points to make in the accuracy test for each solution.

// Number of replicates of all runs (apart from blanks which is specified separately.

// Select if printing as image

// ImageFile to use if printing as an image

// Number of blanks to be added to end of run.

// Option to override the LHPolicy of the test solution with the LHPolicy specified here.

// Option to override the LHPolicy of the test solution with the LHPolicy specified here.

// Name of the output file

// If selected, for each well, all contents will be added to a well before moving on to the next well.
// If not selected, diluent will be added to all wells followed by each solution.

// Specify the dilution factor to use when an intermediate dilution needs to be made to achieve a target concentration.
// If left blank the dilution factor will be calculated based upon the dilution necessary.

// optional parameter allowing pipetting to resume on partially filled plate

// Specify a minimum volume below which a dilution will need to be made.
// The default value is 0.5ul

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _AccuracyTestRequirements() {
}

// Conditions to run on startup
func _AccuracyTestSetup(_ctx context.Context, _input *AccuracyTestInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _AccuracyTestSteps(_ctx context.Context, _input *AccuracyTestInput, _output *AccuracyTestOutput) {

	//--------------------------------------------------------------------
	// Global variables declarations
	//--------------------------------------------------------------------

	var rotate = false
	var autorotate = true
	var wellpositionarray = make([]string, 0)
	_output.Runtowelllocationmap = make(map[string]string)
	_output.Blankwells = make([]string, 0)
	counter := _input.WellsUsed
	var DilutionFactor float64

	var minVolume wunit.Volume

	if _input.MinVolume.EqualTo(wunit.NewVolume(0.0, "ul")) {
		minVolume = wunit.NewVolume(0.5, "ul")
	} else {
		minVolume = _input.MinVolume
	}

	var platenum = 1
	var runs = make([]doe.Run, 1)
	var newruns = make([]doe.Run, 0)

	var err error
	_output.Errors = make([]error, 0)

	// work out plate layout based on picture or just in order
	if _input.Printasimage {

		// image placeholder variables
		var imgBase *goimage.NRGBA

		//--------------------------------------------------------------------
		//Fetching image
		//--------------------------------------------------------------------

		//open Image file
		imgBase, err = image.OpenFile(_input.ImageFile)
		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}

		//--------------------------------------------------------------------
		//Determine which palette to use
		//--------------------------------------------------------------------

		chosencolourpalette := image.AvailablePalettes()["Palette1"]
		positiontocolourmap, _ := image.ImagetoPlatelayout(imgBase, _input.OutPlate, &chosencolourpalette, rotate, autorotate)

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

	// if none is specified, use lhpolicy of first solution
	var referencepolicy wtype.LHPolicy
	var found bool

	if _input.OverrideLHPolicyForTestSolutions == "" {
		lhPolicy := wtype.PolicyName(_input.TestSols[0].TypeName())
		referencepolicy, found = liquidhandling.GetPolicyByName(lhPolicy)

		if found == false {
			execute.Errorf(_ctx, "policy %s not found", lhPolicy.String())
			_output.Errors = append(_output.Errors, fmt.Errorf("policy ", lhPolicy, " not found"))
		}
	} else {
		referencepolicy, found = liquidhandling.GetPolicyByName(_input.OverrideLHPolicyForTestSolutions)

		if found == false {
			execute.Errorf(_ctx, "policy %s not found", _input.OverrideLHPolicyForTestSolutions.String())
			_output.Errors = append(_output.Errors, fmt.Errorf("policy ", _input.OverrideLHPolicyForTestSolutions, " not found"))
		}
	}

	referencekeys := make([]string, 0)
	for key := range referencepolicy {
		referencekeys = append(referencekeys, key)
	}

	// calculate target volumes from concentrations
	var TestSolVolumes = make([]wunit.Volume, len(_input.TestSolConcs))

	for l := 0; l < len(_input.TestSolConcs); l++ {

		for k := 0; k < len(_input.TestSols); k++ {

			stockConc, found := _input.OverrideStockConcentrations[_input.TestSols[k].CName]

			if !found && _input.TestSols[k].HasConcentration() {
				stockConc = _input.TestSols[k].Concentration()
			} else {
				execute.Errorf(_ctx, "No Stock concentration found for %s. Please choose a component with a concentration or override the concentration in OverrideStockConcentrations", _input.TestSols[k].CName)
			}

			vol, err := wunit.VolumeForTargetConcentration(_input.TestSolConcs[l], stockConc, _input.TotalVolume)

			if err != nil {
				execute.Errorf(_ctx, err.Error())
			}

			TestSolVolumes[l] = vol

			for j := 0; j < _input.NumberofReplicates; j++ {

				for i := 0; i < len(runs); i++ {

					var diluted bool
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
					}
					// keep default policy for diluent

					// diluent first

					// change lhpolicy if desired
					if _input.OverrideLHPolicyforDiluent != "" {
						_input.Diluent.Type, err = wtype.LiquidTypeFromString(_input.OverrideLHPolicyforDiluent)
						if err != nil {
							_output.Errors = append(_output.Errors, err)
						}
					}

					var bufferSample *wtype.LHComponent
					var Dilution *wtype.LHComponent

					if TestSolVolumes[l].GreaterThan(wunit.NewVolume(0.0, "ul")) && TestSolVolumes[l].LessThan(minVolume) {

						if _input.SpecifyDilutionFactor == 0.0 {
							DilutionFactor = 4.0 * minVolume.RawValue() / TestSolVolumes[l].RawValue()
							DilutionFactor, err = wutil.Roundto(DilutionFactor, 2)

							if err != nil {
								execute.Errorf(_ctx, err.Error())
							}
						} else {
							DilutionFactor = _input.SpecifyDilutionFactor
						}

						// add diluent to dilution plate ready for dilution
						dilutedSampleBuffer := mixer.Sample(_input.Diluent, wunit.SubtractVolumes(_input.TotalVolume, []wunit.Volume{wunit.DivideVolume(_input.TotalVolume, DilutionFactor)}))
						Dilution = execute.MixNamed(_ctx, _input.OutPlate.Type, wellpositionarray[counter], fmt.Sprint("DilutionPlate", platenum), dilutedSampleBuffer)

						// add same volume to destination plate ready for dilutedsolution
						bufferSample = mixer.Sample(_input.Diluent, wunit.SubtractVolumes(_input.TotalVolume, []wunit.Volume{wunit.MultiplyVolume(TestSolVolumes[l], DilutionFactor)}))

					} else {

						bufferSample = mixer.Sample(_input.Diluent, wunit.SubtractVolumes(_input.TotalVolume, []wunit.Volume{TestSolVolumes[l]})) //SampleForTotalVolume(Diluent, TotalVolume)

					}

					if _input.PipetteOnebyOne {
						eachreaction = append(eachreaction, bufferSample)
					} else {

						solution = execute.MixTo(_ctx, _input.OutPlate.Type, wellpositionarray[counter], platenum, bufferSample)
					}
					// now test sample

					// change liquid class
					if _input.OverrideLHPolicyForTestSolutions != "" {
						_input.TestSols[k].Type, err = wtype.LiquidTypeFromString(_input.OverrideLHPolicyForTestSolutions)
						if err != nil {
							_output.Errors = append(_output.Errors, err)
						}
					}

					if TestSolVolumes[l].GreaterThan(minVolume) {

						//sample
						testSample := mixer.Sample(_input.TestSols[k], TestSolVolumes[l])

						if _input.PipetteOnebyOne {
							eachreaction = append(eachreaction, testSample)
							solution = execute.MixTo(_ctx, _input.OutPlate.Type, wellpositionarray[counter], platenum, eachreaction...)
						} else {
							// pipette out
							solution = execute.Mix(_ctx, solution, testSample)
						}

					} else if TestSolVolumes[l].GreaterThan(wunit.NewVolume(0.0, "ul")) && TestSolVolumes[l].LessThan(minVolume) {
						diluted = true

						DilutionFactor = 4.0 * minVolume.RawValue() / TestSolVolumes[l].RawValue()
						DilutionFactor, err = wutil.Roundto(DilutionFactor, 2)

						if err != nil {
							execute.Errorf(_ctx, err.Error())
						}

						//sample
						dilutionSample := mixer.Sample(_input.TestSols[k], wunit.DivideVolume(_input.TotalVolume, DilutionFactor))
						Dilution = execute.MixNamed(_ctx, _input.OutPlate.Type, wellpositionarray[counter], fmt.Sprint("DilutionPlate", platenum), dilutionSample)

						testSample := mixer.Sample(Dilution, wunit.MultiplyVolume(TestSolVolumes[l], DilutionFactor))

						if _input.PipetteOnebyOne {
							eachreaction = append(eachreaction, testSample)
							solution = execute.MixTo(_ctx, _input.OutPlate.Type, wellpositionarray[counter], platenum, eachreaction...)

						} else {
							// pipette out

							solution = execute.Mix(_ctx, solution, testSample)
						}

					}

					// get annotation info
					lhpolicy := wtype.LiquidTypeName(_input.TestSols[k].Type)

					volume := TestSolVolumes[l].ToString()
					conc := _input.TestSolConcs[l].ToString()

					solutionname := _input.TestSols[k].CName
					stockconc := stockConc.ToString()

					// add Solution Name
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Solution", solutionname)

					// add Solution Name
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Stock Conc", stockconc)

					// add Volume
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Volume", volume)

					// if diluted
					if diluted {
						run = doe.AddAdditionalHeaderandValue(run, "Additional", "PreDilutionFactor", DilutionFactor)
					} else {
						run = doe.AddAdditionalHeaderandValue(run, "Additional", "PreDilutionFactor", 0)
					}

					// add Concentration
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Concentration Set Point", conc)

					// add Replicate
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Replicate", strconv.Itoa(j+1))

					// full description
					description := volume + "_" + solutionname + "_replicate" + strconv.Itoa(j+1) + "_platenum" + strconv.Itoa(platenum)

					// add run to well position lookup table
					_output.Runtowelllocationmap[lhpolicy.String()+"_"+description] = wellpositionarray[counter]

					// add additional info for each run
					fmt.Println("len(runs)", len(runs), "counter", counter, "len(wellpositionarray)", len(wellpositionarray))

					run = doe.AddAdditionalHeaderandValue(run, "Additional", "PlateNumber", strconv.Itoa(platenum))
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Location", wellpositionarray[counter])

					// add setpoint printout to double check correct match up:
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "LHPolicy", lhpolicy)

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

					// print out LHPolicy info
					policy, found := liquidhandling.GetPolicyByName(lhpolicy)
					if !found {
						execute.Errorf(_ctx, "policy %s not found", lhpolicy.String())
						_output.Errors = append(_output.Errors, fmt.Errorf("policy ", lhpolicy, " not found"))
					}

					for _, key := range referencekeys {
						run = doe.AddAdditionalHeaderandValue(run, "Additional", "LHPolicy"+"_"+key, policy[key])
					}

					reactions = append(reactions, solution)
					newruns = append(newruns, run)

					counter = counter + 1

				}

			}
		}
	}

	// add blanks after

	for n := 0; n < platenum; n++ {

		for m := 0; m < _input.NumberofBlanks; m++ {

			var run doe.Run
			// use default policy for blank

			bufferSample := mixer.Sample(_input.Diluent, _input.TotalVolume)

			// add blanks to last column of plate
			well := wutil.NumToAlpha(_input.OutPlate.WlsY-m) + strconv.Itoa(_input.OutPlate.WlsX)

			reaction := execute.MixTo(_ctx, _input.OutPlate.Type, well, n+1, bufferSample)

			_output.Blankwells = append(_output.Blankwells, well)

			// get annotation info
			lhpolicy := wtype.LiquidTypeName(_input.Diluent.Type)

			volume := _input.TotalVolume.ToString()
			conc := "N/A"

			solutionname := _input.Diluent.CName
			stockconc := _input.Diluent.Concentration().ToString()

			// add Solution Name
			run = doe.AddAdditionalHeaderandValue(run, "Additional", "Solution", solutionname)

			// add Solution Name
			run = doe.AddAdditionalHeaderandValue(run, "Additional", "Stock Conc", stockconc)

			// add Volume
			run = doe.AddAdditionalHeaderandValue(run, "Additional", "Volume", volume)

			// if diluted

			run = doe.AddAdditionalHeaderandValue(run, "Additional", "PreDilutionFactor", 0)

			// add Concentration
			run = doe.AddAdditionalHeaderandValue(run, "Additional", "Concentration Set Point", conc)

			// add Replicate
			run = doe.AddAdditionalHeaderandValue(run, "Additional", "Replicate", strconv.Itoa(m+1))

			// add run to well position lookup table
			_output.Runtowelllocationmap["Blank"+strconv.Itoa(m+1)+" platenum"+strconv.Itoa(n+1)] = well

			// add additional info for each run
			fmt.Println("len(runs)", len(runs), "counter", counter, "len(wellpositionarray)", len(wellpositionarray))

			run = doe.AddAdditionalHeaderandValue(run, "Additional", "PlateNumber", strconv.Itoa(platenum))
			run = doe.AddAdditionalHeaderandValue(run, "Additional", "Location", wellpositionarray[counter])

			// add setpoint printout to double check correct match up:
			run = doe.AddAdditionalHeaderandValue(run, "Additional", "LHPolicy", lhpolicy)

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

			// print out LHPolicy info
			policy, found := liquidhandling.GetPolicyByName(lhpolicy)
			if !found {
				execute.Errorf(_ctx, "policy %s not found", lhpolicy.String())
				_output.Errors = append(_output.Errors, fmt.Errorf("policy ", lhpolicy, " not found"))
			}

			for _, key := range referencekeys {
				run = doe.AddAdditionalHeaderandValue(run, "Additional", "LHPolicy"+"_"+key, policy[key])
			}

			reactions = append(reactions, reaction)
			newruns = append(newruns, run)

			counter++

		}

	}

	// export overall DOE design file showing all well locations for all conditions
	xlsxfile := doe.XLSXFileFromRuns(newruns, _input.OutputFilename, "JMP")

	var out bytes.Buffer

	err = xlsxfile.Write(&out)

	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	_output.ExportedFile.Name = _input.OutputFilename
	_output.ExportedFile.WriteAll(out.Bytes())

	_output.Reactions = reactions
	_output.Runcount = len(_output.Reactions)
	_output.Pixelcount = len(wellpositionarray)
	_output.Runs = runs
	_output.Wellpositionarray = wellpositionarray
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AccuracyTestAnalysis(_ctx context.Context, _input *AccuracyTestInput, _output *AccuracyTestOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AccuracyTestValidation(_ctx context.Context, _input *AccuracyTestInput, _output *AccuracyTestOutput) {
}
func _AccuracyTestRun(_ctx context.Context, input *AccuracyTestInput) *AccuracyTestOutput {
	output := &AccuracyTestOutput{}
	_AccuracyTestSetup(_ctx, input)
	_AccuracyTestSteps(_ctx, input, output)
	_AccuracyTestAnalysis(_ctx, input, output)
	_AccuracyTestValidation(_ctx, input, output)
	return output
}

func AccuracyTestRunSteps(_ctx context.Context, input *AccuracyTestInput) *AccuracyTestSOutput {
	soutput := &AccuracyTestSOutput{}
	output := _AccuracyTestRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AccuracyTestNew() interface{} {
	return &AccuracyTestElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AccuracyTestInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AccuracyTestRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AccuracyTestInput{},
			Out: &AccuracyTestOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type AccuracyTestElement struct {
	inject.CheckedRunner
}

type AccuracyTestInput struct {
	Diluent                          *wtype.LHComponent
	ImageFile                        wtype.File
	MinVolume                        wunit.Volume
	NumberofBlanks                   int
	NumberofReplicates               int
	OutPlate                         *wtype.LHPlate
	OutputFilename                   string
	OverrideLHPolicyForTestSolutions wtype.PolicyName
	OverrideLHPolicyforDiluent       wtype.PolicyName
	OverrideStockConcentrations      map[string]wunit.Concentration
	PipetteOnebyOne                  bool
	Printasimage                     bool
	SpecifyDilutionFactor            float64
	TestSolConcs                     []wunit.Concentration
	TestSols                         []*wtype.LHComponent
	TotalVolume                      wunit.Volume
	WellsUsed                        int
}

type AccuracyTestOutput struct {
	Blankwells           []string
	Errors               []error
	ExportedFile         wtype.File
	Pixelcount           int
	Reactions            []*wtype.LHComponent
	Runcount             int
	Runs                 []doe.Run
	Runtowelllocationmap map[string]string
	Wellpositionarray    []string
}

type AccuracyTestSOutput struct {
	Data struct {
		Blankwells           []string
		Errors               []error
		ExportedFile         wtype.File
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
	if err := addComponent(component.Component{Name: "AccuracyTest",
		Constructor: AccuracyTestNew,
		Desc: component.ComponentDesc{
			Desc: "Perform accuracy test protocol using a series of concentrations as set points\n",
			Path: "src/github.com/antha-lang/elements/an/Utility/AccuracyTest/AccuracyTest_Conc.an",
			Params: []component.ParamDesc{
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "ImageFile", Desc: "ImageFile to use if printing as an image\n", Kind: "Parameters"},
				{Name: "MinVolume", Desc: "Specify a minimum volume below which a dilution will need to be made.\nThe default value is 0.5ul\n", Kind: "Parameters"},
				{Name: "NumberofBlanks", Desc: "Number of blanks to be added to end of run.\n", Kind: "Parameters"},
				{Name: "NumberofReplicates", Desc: "Number of replicates of all runs (apart from blanks which is specified separately.\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "OutputFilename", Desc: "Name of the output file\n", Kind: "Parameters"},
				{Name: "OverrideLHPolicyForTestSolutions", Desc: "Option to override the LHPolicy of the test solution with the LHPolicy specified here.\n", Kind: "Parameters"},
				{Name: "OverrideLHPolicyforDiluent", Desc: "Option to override the LHPolicy of the test solution with the LHPolicy specified here.\n", Kind: "Parameters"},
				{Name: "OverrideStockConcentrations", Desc: "Option to over ride the stock concentrations of the solutions being tested.\n", Kind: "Parameters"},
				{Name: "PipetteOnebyOne", Desc: "If selected, for each well, all contents will be added to a well before moving on to the next well.\nIf not selected, diluent will be added to all wells followed by each solution.\n", Kind: "Parameters"},
				{Name: "Printasimage", Desc: "Select if printing as image\n", Kind: "Parameters"},
				{Name: "SpecifyDilutionFactor", Desc: "Specify the dilution factor to use when an intermediate dilution needs to be made to achieve a target concentration.\nIf left blank the dilution factor will be calculated based upon the dilution necessary.\n", Kind: "Parameters"},
				{Name: "TestSolConcs", Desc: "Concentration Set points to make in the accuracy test for each solution.\n", Kind: "Parameters"},
				{Name: "TestSols", Desc: "", Kind: "Inputs"},
				{Name: "TotalVolume", Desc: "Total volunme per well\n", Kind: "Parameters"},
				{Name: "WellsUsed", Desc: "optional parameter allowing pipetting to resume on partially filled plate\n", Kind: "Parameters"},
				{Name: "Blankwells", Desc: "", Kind: "Data"},
				{Name: "Errors", Desc: "", Kind: "Data"},
				{Name: "ExportedFile", Desc: "", Kind: "Data"},
				{Name: "Pixelcount", Desc: "", Kind: "Data"},
				{Name: "Reactions", Desc: "", Kind: "Outputs"},
				{Name: "Runcount", Desc: "", Kind: "Data"},
				{Name: "Runs", Desc: "", Kind: "Data"},
				{Name: "Runtowelllocationmap", Desc: "", Kind: "Data"},
				{Name: "Wellpositionarray", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
