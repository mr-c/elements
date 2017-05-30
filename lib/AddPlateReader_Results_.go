// Protocol to parse plate reader results and match up with a plate set up by the accuracy test.
// Some processing is carried out to:
// A: Plot expected results (based on mathematically diluting the stock concentration) vs actual (measured concentrations from beer-lambert law, A = εcl)
// B: Plot volume by correctness factor (Actual conc / Expected conc)
// C: Plot Actual conc vs correctness factor
// D: Plot run order vs correctness factor
// E: Calculate R2
// F: Calculate Coefficent of variance for each pipetting volume
// G: Validate results against success thresholds for R2 and %CV
// Additional optional features will return
// (1) the wavelength with optimal signal to noise for an aborbance spectrum
// (2) Comparision with manual pipetting steps
package lib

import (
	"bytes"
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Pubchem"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/buffers"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/doe"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/platereader"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/plot"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/factory"
	"github.com/montanaflynn/stats"
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Input parameters for this protocol (data)

// input file containing the Plate reader results exported from Mars
// Design file for the executed experiment containing the corresponding plate and well locations

// i.e. the sheet position in the plate reader results excel file; starting from 0
// current supported formats are "JMP" and "DX"
// set the desired name for the output file, if this is blank it will append the design file name with _output

//  Wavelength to use for calculations, should match up with extinction coefficient for molecule of interest
// extinction coefficient for target Molecule at the specified wavelength; e.g. 20330 for tartrazine at 472nm

// This should match the label in the header for each column in the plate reader result file, e.g. "Abs Spectrum"

/// wells of the blank sample locations on the plate

// whether the scan should be used to return the wavelength with maximum signal to noise found
// well used for finding wavelength with optimal signal to noise. This is ignored if FindOptWavelength is set to false

// name your response

//  Option to compare to manual pipetting
// if comparing to manual pipetting set the wells to use for each concentration here
// volume of diluent per well
// if true the StockVol represents the total volume per well instead of a fixed volume which the test solution was added to

// validation requirements
// set a threshold above which R2 will pass; 0 = 0%, 1 = 100%; e.g. 0.7 = 70%
// set a threshold below which CV will pass; 0 = 0%, 1 = 100%; e.g. 0.2 = 20%

// Option to override moecular weight value of a mpolecule

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// The name of the molecule to analyse. This will be used to find matching solutions in the design file and to look up the molecular weight.
// Currently only one solution name can be run at a time.

// Physical outputs from this protocol with types

func _AddPlateReader_ResultsRequirements() {
}

// Conditions to run on startup
func _AddPlateReader_ResultsSetup(_ctx context.Context, _input *AddPlateReader_ResultsInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _AddPlateReader_ResultsSteps(_ctx context.Context, _input *AddPlateReader_ResultsInput, _output *AddPlateReader_ResultsOutput) {

	if _input.OutputFilename == "" {

		_, filename := filepath.Split(_input.DesignFile.Name)

		_input.OutputFilename = strings.Split(filename, ".")[0] + "_output" + "_" + _input.Molecule.CName + fmt.Sprint(time.Now().Format("20060102150405")) + ".xlsx"
	} else {
		_input.OutputFilename = strings.Split(_input.OutputFilename, ".")[0] + "_" + _input.Molecule.CName + ".xlsx"
	}

	var actualconcentrations = make(map[string]wunit.Concentration)
	_output.ResponsetoManualValuesmap = make(map[string][]float64)

	var Molecularweight float64

	if molecularWeight, found := _input.OverrideMolecularWeight[_input.Molecule.CName]; found {
		Molecularweight = molecularWeight
	} else if molecularWeight, found := _input.OverrideMolecularWeight["default"]; found {
		Molecularweight = molecularWeight
	} else {
		molecule, err := pubchem.MakeMolecule(_input.Molecule.CName)

		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}

		Molecularweight = molecule.MolecularWeight
	}

	data, err := _input.MarsResultsFileXLSX.ReadAll()

	if err != nil {
		execute.Errorf(_ctx, "Error reading Mars file %s", err.Error())
	}

	marsdata, err := parser.ParseMarsXLSXBinary(data, _input.SheetNumber)
	if err != nil {
		_output.Errors = append(_output.Errors, err.Error())
		execute.Errorf(_ctx, err.Error())
	}

	// range through pairing up wells from mars output and doe design

	var runs []doe.Run

	// find out int factors from liquidhandling policies
	policyitemmap := wtype.MakePolicyItems()
	intfactors := make([]string, 0)

	for key, val := range policyitemmap {

		if val.Type.Name() == "int" {
			intfactors = append(intfactors, key)
		}
	}

	designFileData, err := _input.DesignFile.ReadAll()

	if err != nil {
		execute.Errorf(_ctx, "Error reading Design file %s", err.Error())
	}

	if _input.DesignFiletype == "DX" {
		runs, err = doe.RunsFromDXDesignContents(designFileData, intfactors)
		if err != nil {
			panic(err)
		}
	} else if _input.DesignFiletype == "JMP" {
		runs, err = doe.RunsFromJMPDesignContents(designFileData, []int{}, []int{}, intfactors)
		if err != nil {
			panic(err)
		}
	}
	_output.BlankValues = make([]float64, 0)

	for i := range _input.Blanks {

		blankValue, err := marsdata.ReadingsAsAverage(_input.Blanks[i], 1, _input.Wavelength, _input.ReadingTypeinMarsFile)

		if err != nil {
			execute.Errorf(_ctx, fmt.Sprint("blank sample not found. ", err.Error()))
		}

		_output.BlankValues = append(_output.BlankValues, blankValue)
	}

	runswithresponses := make([]doe.Run, 0)

	for k, run := range runs {

		// values for r2 to reset each run

		//xvalues := make([]float64, 0)
		//yvalues := make([]float64, 0)

		// add origin
		//xvalues = append(xvalues, 0.0)
		//yvalues = append(yvalues, 0.0)

		// Only get infor for the solution in question
		if runSolution, err := run.GetAdditionalInfo("Solution"); err == nil && runSolution == _input.Molecule.CName {

			var samples []string
			var manualsamples []string
			var ManualValues = make([]float64, 0)
			var manual float64
			var absorbance wtype.Absorbance
			var manualabsorbance wtype.Absorbance
			//var actualconcreplicates = make([]float64, 0)
			var manualCorrectnessFactorValues = make([]float64, 0)
			var correctnessFactorValues = make([]float64, 0)

			experimentalvolumeinterface, err := runs[k].GetAdditionalInfo("Volume")

			experimentalvolumestr := experimentalvolumeinterface.(string)

			//experimentalvolumestr = strings.TrimSpace(experimentalvolumestr)

			var volandunit []string

			if strings.Count(experimentalvolumestr, " ") == 1 {
				volandunit = strings.Split(experimentalvolumestr, " ")
			} else if strings.Count(experimentalvolumestr, "ul") == 1 && strings.HasSuffix(experimentalvolumestr, "ul") {
				volandunit = []string{strings.Trim(experimentalvolumestr, "ul"), "ul"}
			}

			vol, err := strconv.ParseFloat(strings.TrimSpace(volandunit[0]), 64)

			if err != nil {
				execute.Errorf(_ctx, err.Error())
			}

			experimentalvolume := wunit.NewVolume(vol, strings.TrimSpace(volandunit[1]))

			actualconcentrations[experimentalvolume.ToString()] = buffers.DiluteBasedonMolecularWeight(Molecularweight, _input.StockconcinMperL, experimentalvolume, _input.Diluent.CName, wunit.SubtractVolumes(_input.Stockvol, []wunit.Volume{experimentalvolume}))

			//locationHeaders := ResponsetoLocationMap[response]

			//  manual pipetting well
			if wellsmap, ok := _input.VolumeToManualwells[experimentalvolumestr]; _input.ManualComparison && ok {

				manualwell := wellsmap[0] // 1st well of array only

				manual, err = marsdata.ReadingsAsAverage(manualwell, 1, _input.Wavelength, _input.ReadingTypeinMarsFile)

				if err != nil {
					execute.Errorf(_ctx, err.Error())
				}

				manualsamples = _input.VolumeToManualwells[experimentalvolumestr]

				for i := range manualsamples {
					manualvalue, err := marsdata.ReadingsAsAverage(manualsamples[i], 1, _input.Wavelength, _input.ReadingTypeinMarsFile)

					if err != nil {
						execute.Errorf(_ctx, err.Error())
					}

					ManualValues = append(ManualValues, manualvalue)
				}

				_output.ResponsetoManualValuesmap[experimentalvolumestr] = ManualValues

				run = doe.AddNewResponseFieldandValue(run, _input.Responsecolumntofill+" Manual Raw average "+strconv.Itoa(_input.Wavelength), manual)

			} else if _input.ManualComparison {
				run = doe.AddNewResponseFieldandValue(run, _input.Responsecolumntofill+" Manual Raw average "+strconv.Itoa(_input.Wavelength), 0.0)
			}

			// then per replicate ...

			//for i, locationheader := range locationHeaders {
			well, err := runs[k].GetAdditionalInfo("Location")
			if err != nil {
				panic(err)
			}

			// check optimal difference for each well
			if _input.FindOptWavelength {
				_output.MeasuredOptimalWavelength, err = marsdata.FindOptimalWavelength(_input.WellForScanAnalysis, _input.Blanks[0], "Raw Data")

				if err != nil {
					execute.Errorf(_ctx, fmt.Sprint("Error found with well for scan analysis: ", err.Error()))
				}
			}

			rawaverage, err := marsdata.ReadingsAsAverage(well.(string), 1, _input.Wavelength, _input.ReadingTypeinMarsFile)

			run = doe.AddNewResponseFieldandValue(run, _input.Responsecolumntofill+" Raw average "+strconv.Itoa(_input.Wavelength), rawaverage)

			// blank correct

			samples = []string{well.(string)}

			blankcorrected, err := marsdata.BlankCorrect(samples, _input.Blanks, _input.Wavelength, _input.ReadingTypeinMarsFile)

			run = doe.AddNewResponseFieldandValue(run, _input.Responsecolumntofill+" BlankCorrected "+strconv.Itoa(_input.Wavelength), blankcorrected)

			// path length correct
			var volused wunit.Volume
			if _input.StockEqualsTotalVolPerWell {
				volused = _input.Stockvol
			} else {
				volused = wunit.AddVolumes([]wunit.Volume{_input.Stockvol, experimentalvolume})
			}
			pathlength, err := platereader.EstimatePathLength(factory.GetPlateByType(_input.PlateType.Type), volused)

			if err != nil {
				panic(err)
			}

			run = doe.AddNewResponseFieldandValue(run, _input.Responsecolumntofill+" pathlength "+strconv.Itoa(_input.Wavelength), pathlength.ToString())

			absorbance.Reading = blankcorrected

			pathlengthcorrect := platereader.PathlengthCorrect(pathlength, absorbance)

			run = doe.AddNewResponseFieldandValue(run, _input.Responsecolumntofill+" Pathlength corrected "+strconv.Itoa(_input.Wavelength), pathlengthcorrect.Reading)

			// molar absorbtivity of tartazine at 472nm is 20330
			// http://www.biochrom.co.uk/faq/8/119/what-is-the-limit-of-detection-of-the-zenyth-200.html

			actualconc := platereader.Concentration(pathlengthcorrect, _input.Extinctioncoefficient)

			run = doe.AddNewResponseFieldandValue(run, _input.Responsecolumntofill+"ActualConc", actualconc.SIValue())

			// calculate correctness factor based on expected conc

			expectedconc := actualconcentrations[experimentalvolume.ToString()]
			correctnessfactor := actualconc.SIValue() / expectedconc.SIValue()

			run = doe.AddNewResponseFieldandValue(run, _input.Responsecolumntofill+" ExpectedConc "+strconv.Itoa(_input.Wavelength), expectedconc.SIValue())

			// if Infinity or Not a number set to 0
			if math.IsInf(correctnessfactor, 0) || math.IsNaN(correctnessfactor) {
				correctnessfactor = 0.0
			}

			run = doe.AddNewResponseFieldandValue(run, _input.Responsecolumntofill+" CorrectnessFactor "+strconv.Itoa(_input.Wavelength), correctnessfactor)

			correctnessFactorValues = append(correctnessFactorValues, correctnessfactor)

			// add comparison to manually pipetted wells
			if _, ok := _input.VolumeToManualwells[experimentalvolumestr]; _input.ManualComparison && ok {
				manualblankcorrected, err := marsdata.BlankCorrect(manualsamples, _input.Blanks, _input.Wavelength, _input.ReadingTypeinMarsFile)
				if err != nil {
					execute.Errorf(_ctx, err.Error())
				}
				manualabsorbance.Reading = manualblankcorrected
				manualpathlengthcorrect := platereader.PathlengthCorrect(pathlength, manualabsorbance)
				manualactualconc := platereader.Concentration(manualpathlengthcorrect, _input.Extinctioncoefficient)
				manualcorrectnessfactor := actualconc.SIValue() / manualactualconc.SIValue()
				manualCorrectnessFactorValues = append(manualCorrectnessFactorValues, manualcorrectnessfactor)

				run = doe.AddNewResponseFieldandValue(run, _input.Responsecolumntofill+"ManualActualConc", manualactualconc.SIValue())
				run = doe.AddNewResponseFieldandValue(run, _input.Responsecolumntofill+" ManualCorrectnessFactor "+strconv.Itoa(_input.Wavelength), manualcorrectnessfactor)
			} else if _input.ManualComparison {
				run = doe.AddNewResponseFieldandValue(run, _input.Responsecolumntofill+"ManualActualConc", 0.0)
				run = doe.AddNewResponseFieldandValue(run, _input.Responsecolumntofill+" ManualCorrectnessFactor "+strconv.Itoa(_input.Wavelength), 0.0)

			}

			run = doe.AddNewResponseFieldandValue(run, "Runorder", k)

			runswithresponses = append(runswithresponses, run)

		}

	}

	xlsxfile := doe.XLSXFileFromRuns(runswithresponses, _input.OutputFilename, _input.DesignFiletype)

	_output.Runs = runswithresponses

	var buffer bytes.Buffer

	xlsxfile.Write(&buffer)

	_output.OutPutDesignFile.Name = _input.OutputFilename

	_output.OutPutDesignFile.WriteAll(buffer.Bytes())

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AddPlateReader_ResultsAnalysis(_ctx context.Context, _input *AddPlateReader_ResultsInput, _output *AddPlateReader_ResultsOutput) {

	_output.Errors = make([]string, 0)

	xvalues := make([]float64, 0)
	yvalues := make([]float64, 0)

	// add origin
	xvalues = append(xvalues, 0.0)
	yvalues = append(yvalues, 0.0)

	fmt.Println("in analysis")

	if len(_output.Runs) == 0 {
		execute.Errorf(_ctx, "no runs")
	}
	// 1. now calculate r2 and plot results
	for i, runwithresponses := range _output.Runs {
		// values for r2 to reset each run

		// get response value and check if it's a float64 type
		expectedconc, err := runwithresponses.GetResponseValue("Absorbance ExpectedConc " + strconv.Itoa(_input.Wavelength))

		if err != nil {
			_output.Errors = append(_output.Errors, err.Error())
		}

		expectedconcfloat, floattrue := expectedconc.(float64)
		// if float64 is true
		if floattrue {
			xvalues = append(xvalues, expectedconcfloat)
		} else {
			execute.Errorf(_ctx, "Run"+fmt.Sprint(i, runwithresponses)+" ExpectedConc:"+fmt.Sprint(expectedconcfloat))
		}

		// get response value and check if it's a float64 type
		actualconc, err := runwithresponses.GetResponseValue("AbsorbanceActualConc")

		if err != nil {
			fmt.Println(err.Error())
			_output.Errors = append(_output.Errors, err.Error())
		}

		actualconcfloat, floattrue := actualconc.(float64)

		if floattrue {
			yvalues = append(yvalues, actualconcfloat)
		} else {
			fmt.Println(err.Error())
			execute.Errorf(_ctx, " ActualConc:"+fmt.Sprint(actualconcfloat))
		}

	}

	_output.R2, _output.Variance, _output.Formula = plot.Rsquared("Expected Conc", xvalues, "Actual Conc", yvalues)
	//run.AddResponseValue("R2", rsquared)

	xygraph, err := plot.Plot(xvalues, [][]float64{yvalues})
	if err != nil {
		_output.Errors = append(_output.Errors, err.Error())
	}

	plot.AddAxesTitles(xygraph, "Expected Conc M/l", "Measured Conc M/l")

	xygraph.Title.Text = _input.Molecule.CName + ": Expected vs Measured Concentration"

	filenameandextension := strings.Split(_input.OutputFilename, ".")

	_output.ActualVsExpectedPlot, err = plot.Export(xygraph, "20cm", "20cm", filenameandextension[0]+"_plot"+".png")

	if err != nil {
		_output.Errors = append(_output.Errors, err.Error())
		execute.Errorf(_ctx, err.Error())
	}

	// reset
	xvalues = make([]float64, 0)
	yvalues = make([]float64, 0)

	// add origin
	xvalues = append(xvalues, 0.0)
	yvalues = append(yvalues, 0.0)

	// 2. now plot correctnessfactor
	for i, runwithresponses := range _output.Runs {
		// values for r2 to reset each run

		// get response value and check if it's a float64 type
		expectedconc, err := runwithresponses.GetResponseValue("Absorbance ExpectedConc " + strconv.Itoa(_input.Wavelength))

		if err != nil {
			_output.Errors = append(_output.Errors, err.Error())
		}

		expectedconcfloat, floattrue := expectedconc.(float64)
		// if float64 is true
		if floattrue {
			xvalues = append(xvalues, expectedconcfloat)
		} else {
			execute.Errorf(_ctx, "Run"+fmt.Sprint(i, runwithresponses)+" ExpectedConc:"+fmt.Sprint(expectedconcfloat))
		}

		// get response value and check if it's a float64 type
		correctness, err := runwithresponses.GetResponseValue("Absorbance CorrectnessFactor " + strconv.Itoa(_input.Wavelength))

		if err != nil {
			fmt.Println(err.Error())
			_output.Errors = append(_output.Errors, err.Error())
		}

		correctnessfloat, floattrue := correctness.(float64)

		if floattrue {
			yvalues = append(yvalues, correctnessfloat)
		} else {
			fmt.Println(err.Error())
			execute.Errorf(_ctx, " Absorbance CorrectnessFactor:"+fmt.Sprint(correctnessfloat))
		}

	}

	_output.R2_CorrectnessFactor, _, _ = plot.Rsquared("Expected Conc", xvalues, "Correctness Factor", yvalues)
	//run.AddResponseValue("R2", rsquared)

	correctnessgraph, err := plot.Plot(xvalues, [][]float64{yvalues})

	if err != nil {
		_output.Errors = append(_output.Errors, err.Error())
	}

	plot.AddAxesTitles(correctnessgraph, "Target Conc M/l", "Correctness Factor (Measured Conc / Expected Conc)")

	correctnessgraph.Title.Text = _input.Molecule.CName + ": Correctness Factor vs Target Concentration"

	_output.CorrectnessFactorPlot, err = plot.Export(correctnessgraph, "20cm", "20cm", filenameandextension[0]+"_correctnessfactor"+".png")

	if err != nil {
		_output.Errors = append(_output.Errors, err.Error())
		execute.Errorf(_ctx, err.Error())

	}

	// reset
	xvalues = make([]float64, 0)
	yvalues = make([]float64, 0)

	// add origin
	xvalues = append(xvalues, 0.0)
	yvalues = append(yvalues, 0.0)

	// 3. now look for systematic errors
	for i, runwithresponses := range _output.Runs {
		// values for r2 to reset each run

		// get response value and check if it's a float64 type
		runorder, err := runwithresponses.GetResponseValue("Runorder")

		if err != nil {
			_output.Errors = append(_output.Errors, err.Error())
		}

		runorderint, inttrue := runorder.(int)
		// if int is true
		if inttrue {
			xvalues = append(xvalues, float64(runorderint))
		} else {
			execute.Errorf(_ctx, "Run"+fmt.Sprint(i, runwithresponses)+" Run Order:"+fmt.Sprint(runorderint), " not an int")
		}

		// get response value and check if it's a float64 type
		actualconc, err := runwithresponses.GetResponseValue("AbsorbanceActualConc")

		if err != nil {
			fmt.Println(err.Error())
			_output.Errors = append(_output.Errors, err.Error())
		}

		actualconcfloat, floattrue := actualconc.(float64)

		if floattrue {
			yvalues = append(yvalues, actualconcfloat)
		} else {
			fmt.Println(err.Error())
			execute.Errorf(_ctx, " ActualConc:"+fmt.Sprint(actualconcfloat))
		}
	}

	runorderconcgraph, err := plot.Plot(xvalues, [][]float64{yvalues})

	if err != nil {
		_output.Errors = append(_output.Errors, err.Error())
	}

	plot.Export(runorderconcgraph, "10cm", "10cm", filenameandextension[0]+"_runorder"+".png")

	// reset
	xvalues = make([]float64, 0)
	yvalues = make([]float64, 0)

	// add origin
	xvalues = append(xvalues, 0.0)
	yvalues = append(yvalues, 0.0)

	// 4.  now look for systematic errors with correctness factor
	for i, runwithresponses := range _output.Runs {
		// values for r2 to reset each run

		// get response value and check if it's a float64 type
		runorder, err := runwithresponses.GetResponseValue("Runorder")

		if err != nil {
			_output.Errors = append(_output.Errors, err.Error())
		}

		runorderint, inttrue := runorder.(int)
		// if int is true
		if inttrue {
			xvalues = append(xvalues, float64(runorderint))
		} else {
			execute.Errorf(_ctx, "Run"+fmt.Sprint(i, runwithresponses)+" Run Order:"+fmt.Sprint(runorderint), " not an int")
		}

		// get response value and check if it's a float64 type
		correctness, err := runwithresponses.GetResponseValue("Absorbance CorrectnessFactor " + strconv.Itoa(_input.Wavelength))

		if err != nil {
			fmt.Println(err.Error())
			_output.Errors = append(_output.Errors, err.Error())
		}

		correctnessfloat, floattrue := correctness.(float64)

		if floattrue {
			yvalues = append(yvalues, correctnessfloat)
		} else {
			fmt.Println(err.Error())
			execute.Errorf(_ctx, " Absorbance CorrectnessFactor:"+fmt.Sprint(correctnessfloat))
		}

	}

	runordercorrectnessgraph, err := plot.Plot(xvalues, [][]float64{yvalues})

	if err != nil {
		_output.Errors = append(_output.Errors, err.Error())
	}

	plot.Export(runordercorrectnessgraph, "10cm", "10cm", filenameandextension[0]+"_runorder_correctnessfactor"+".png")

	// 5. workout CV for each volume
	replicateactualconcmap := make(map[string][]float64)
	_output.VolumeToActualConc = make(map[string]Dataset)
	replicatevalues := make([]float64, 0)

	replicatecorrectnessmap := make(map[string][]float64)
	correctnessvalues := make([]float64, 0)
	_output.VolumeToCorrectnessFactor = make(map[string]Dataset)

	//counter := 0

	// make map of replicate values for Actual Conc
	for _, runwithresponses := range _output.Runs {

		volstr, err := runwithresponses.GetAdditionalInfo("Volume")

		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}

		/*
			repstr, err := runwithresponses.GetAdditionalInfo("Replicate")

			if err != nil {
				Errorf(err.Error())
			}
		*/
		actualconc, err := runwithresponses.GetResponseValue("AbsorbanceActualConc")

		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}

		/*rep, err := strconv.Atoi(repstr.(string))

		if err != nil {
			Errorf(err.Error())
		}
		*/

		// Actual Conc map
		if _, found := replicateactualconcmap[volstr.(string)]; found /*&& rep == counter*/ {
			replicatevalues = replicateactualconcmap[volstr.(string)]
			replicatevalues = append(replicatevalues, actualconc.(float64))
			replicateactualconcmap[volstr.(string)] = replicatevalues
			replicatevalues = make([]float64, 0)
			//counter++
		} else if _, found := replicateactualconcmap[volstr.(string)]; !found {
			replicatevalues = append(replicatevalues, actualconc.(float64))
			replicateactualconcmap[volstr.(string)] = replicatevalues
			replicatevalues = make([]float64, 0)
			//counter++
		}

		// get response value and check if it's a float64 type
		correctness, err := runwithresponses.GetResponseValue("Absorbance CorrectnessFactor " + strconv.Itoa(_input.Wavelength))

		if err != nil {
			fmt.Println(err.Error())
			_output.Errors = append(_output.Errors, err.Error())
		}

		correctnessfloat, floattrue := correctness.(float64)

		if !floattrue {
			fmt.Println(err.Error())
			execute.Errorf(_ctx, " Correctnessfloat not float but:"+fmt.Sprint(correctnessfloat))
		}

		// correctness factor map
		if _, found := replicatecorrectnessmap[volstr.(string)]; found /*&& rep == counter*/ {
			correctnessvalues = replicatecorrectnessmap[volstr.(string)]
			correctnessvalues = append(correctnessvalues, correctnessfloat)
			replicatecorrectnessmap[volstr.(string)] = correctnessvalues
			correctnessvalues = make([]float64, 0)
			//counter++
		} else if _, found := replicatecorrectnessmap[volstr.(string)]; !found {
			correctnessvalues = append(correctnessvalues, correctnessfloat)
			replicatecorrectnessmap[volstr.(string)] = correctnessvalues
			correctnessvalues = make([]float64, 0)
			//counter++
		}

	}

	// process into datasets
	for key, values := range replicateactualconcmap {

		var dataset Dataset
		// process replicates into mean and cv
		dataset.Name = key + "_AbsorbanceActualConc"
		dataset.Mean, _ = stats.Mean(values)
		dataset.StdDev, _ = stats.StdDevS(values)
		dataset.Values = values

		dataset.CV = dataset.StdDev / dataset.Mean * float64(100)
		_output.VolumeToActualConc[key] = dataset

	}

	// process into datasets
	for key, values := range replicatecorrectnessmap {

		var dataset Dataset
		// process replicates into mean and cv
		dataset.Name = key + "_CorrectnessFactor"
		dataset.Mean, _ = stats.Mean(values)
		dataset.StdDev, _ = stats.StdDevS(values)
		dataset.Values = values

		dataset.CV = dataset.StdDev / dataset.Mean * float64(100)

		// if CV == Infinity or Not a number set to -1.0
		if math.IsInf(dataset.CV, 0) || math.IsNaN(dataset.CV) {
			dataset.CV = -1.0
		}

		_output.VolumeToCorrectnessFactor[key] = dataset

	}

	if _input.ManualComparison {

		// reset
		xvalues = make([]float64, 0)
		yvalues = make([]float64, 0)

		// add origin
		xvalues = append(xvalues, 0.0)
		yvalues = append(yvalues, 0.0)

		// 2. now plot correctnessfactor
		for i, runwithresponses := range _output.Runs {
			// values for r2 to reset each run

			// get response value and check if it's a float64 type
			expectedconc, err := runwithresponses.GetResponseValue("Absorbance ExpectedConc " + strconv.Itoa(_input.Wavelength))

			if err != nil {
				_output.Errors = append(_output.Errors, err.Error())
			}

			expectedconcfloat, floattrue := expectedconc.(float64)
			// if float64 is true
			if floattrue {
				xvalues = append(xvalues, expectedconcfloat)
			} else {
				execute.Errorf(_ctx, "Run"+fmt.Sprint(i, runwithresponses)+" ExpectedConc:"+fmt.Sprint(expectedconcfloat))
			}

			// get response value and check if it's a float64 type
			correctness, err := runwithresponses.GetResponseValue("Absorbance ManualCorrectnessFactor " + strconv.Itoa(_input.Wavelength))

			if err != nil {
				fmt.Println(err.Error())
				_output.Errors = append(_output.Errors, err.Error())
			}

			correctnessfloat, floattrue := correctness.(float64)

			if floattrue {
				yvalues = append(yvalues, correctnessfloat)
			} else {
				fmt.Println(err.Error())
				execute.Errorf(_ctx, "Manual Absorbance CorrectnessFactor:"+fmt.Sprint(correctnessfloat))
			}

		}

		_output.R2_CorrectnessFactor, _, _ = plot.Rsquared("Expected Conc", xvalues, "Manual Correctness Factor", yvalues)
		//run.AddResponseValue("R2", rsquared)

		correctnessgraph, err := plot.Plot(xvalues, [][]float64{yvalues})

		if err != nil {
			_output.Errors = append(_output.Errors, err.Error())
		}

		plot.Export(correctnessgraph, "10cm", "10cm", filenameandextension[0]+"_Manualcorrectnessfactor"+".png")

	}

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AddPlateReader_ResultsValidation(_ctx context.Context, _input *AddPlateReader_ResultsInput, _output *AddPlateReader_ResultsOutput) {

	_output.CVpass = true

	if _output.R2 > _input.R2threshold {
		_output.R2Pass = true
	} else {
		_output.Errors = append(_output.Errors, fmt.Sprint("R2 threshold of ", _input.R2threshold, " not met; R2 value = ", _output.R2))
	}

	for key, dataset := range _output.VolumeToActualConc {

		if dataset.CV > _input.CVthreshold {
			_output.CVpass = false
			_output.Errors = append(_output.Errors, fmt.Sprint(key, " coefficient of variance above ", _input.CVthreshold, " percent threshold; CV value = ", dataset.CV))
		}
	}

}

type Dataset struct {
	Name   string
	Values []float64
	Mean   float64
	StdDev float64
	CV     float64
}

func _AddPlateReader_ResultsRun(_ctx context.Context, input *AddPlateReader_ResultsInput) *AddPlateReader_ResultsOutput {
	output := &AddPlateReader_ResultsOutput{}
	_AddPlateReader_ResultsSetup(_ctx, input)
	_AddPlateReader_ResultsSteps(_ctx, input, output)
	_AddPlateReader_ResultsAnalysis(_ctx, input, output)
	_AddPlateReader_ResultsValidation(_ctx, input, output)
	return output
}

func AddPlateReader_ResultsRunSteps(_ctx context.Context, input *AddPlateReader_ResultsInput) *AddPlateReader_ResultsSOutput {
	soutput := &AddPlateReader_ResultsSOutput{}
	output := _AddPlateReader_ResultsRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AddPlateReader_ResultsNew() interface{} {
	return &AddPlateReader_ResultsElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AddPlateReader_ResultsInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AddPlateReader_ResultsRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AddPlateReader_ResultsInput{},
			Out: &AddPlateReader_ResultsOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type AddPlateReader_ResultsElement struct {
	inject.CheckedRunner
}

type AddPlateReader_ResultsInput struct {
	Blanks                     []string
	CVthreshold                float64
	DesignFile                 wtype.File
	DesignFiletype             string
	Diluent                    *wtype.LHComponent
	Extinctioncoefficient      float64
	FindOptWavelength          bool
	ManualComparison           bool
	MarsResultsFileXLSX        wtype.File
	Molecule                   *wtype.LHComponent
	OutputFilename             string
	OverrideMolecularWeight    map[string]float64
	PlateType                  *wtype.LHPlate
	R2threshold                float64
	ReadingTypeinMarsFile      string
	Responsecolumntofill       string
	SheetNumber                int
	StockEqualsTotalVolPerWell bool
	StockconcinMperL           wunit.Concentration
	Stockvol                   wunit.Volume
	VolumeToManualwells        map[string][]string
	Wavelength                 int
	WellForScanAnalysis        string
}

type AddPlateReader_ResultsOutput struct {
	ActualVsExpectedPlot      wtype.File
	BlankValues               []float64
	CV                        float64
	CVpass                    bool
	CorrectnessFactorPlot     wtype.File
	Errors                    []string
	Formula                   string
	MeasuredOptimalWavelength int
	OutPutDesignFile          wtype.File
	R2                        float64
	R2Pass                    bool
	R2_CorrectnessFactor      float64
	ResponsetoManualValuesmap map[string][]float64
	Runs                      []doe.Run
	Variance                  float64
	VolumeToActualConc        map[string]Dataset
	VolumeToCorrectnessFactor map[string]Dataset
}

type AddPlateReader_ResultsSOutput struct {
	Data struct {
		ActualVsExpectedPlot      wtype.File
		BlankValues               []float64
		CV                        float64
		CVpass                    bool
		CorrectnessFactorPlot     wtype.File
		Errors                    []string
		Formula                   string
		MeasuredOptimalWavelength int
		OutPutDesignFile          wtype.File
		R2                        float64
		R2Pass                    bool
		R2_CorrectnessFactor      float64
		ResponsetoManualValuesmap map[string][]float64
		Runs                      []doe.Run
		Variance                  float64
		VolumeToActualConc        map[string]Dataset
		VolumeToCorrectnessFactor map[string]Dataset
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AddPlateReader_Results",
		Constructor: AddPlateReader_ResultsNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol to parse plate reader results and match up with a plate set up by the accuracy test.\nSome processing is carried out to:\nA: Plot expected results (based on mathematically diluting the stock concentration) vs actual (measured concentrations from beer-lambert law, A = εcl)\nB: Plot volume by correctness factor (Actual conc / Expected conc)\nC: Plot Actual conc vs correctness factor\nD: Plot run order vs correctness factor\nE: Calculate R2\nF: Calculate Coefficent of variance for each pipetting volume\nG: Validate results against success thresholds for R2 and %CV\nAdditional optional features will return\n(1) the wavelength with optimal signal to noise for an aborbance spectrum\n(2) Comparision with manual pipetting steps\n",
			Path: "src/github.com/antha-lang/elements/an/Utility/AccuracyTest/AddPlateReaderResults.an",
			Params: []component.ParamDesc{
				{Name: "Blanks", Desc: "/ wells of the blank sample locations on the plate\n", Kind: "Parameters"},
				{Name: "CVthreshold", Desc: "set a threshold below which CV will pass; 0 = 0%, 1 = 100%; e.g. 0.2 = 20%\n", Kind: "Parameters"},
				{Name: "DesignFile", Desc: "Design file for the executed experiment containing the corresponding plate and well locations\n", Kind: "Parameters"},
				{Name: "DesignFiletype", Desc: "current supported formats are \"JMP\" and \"DX\"\n", Kind: "Parameters"},
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "Extinctioncoefficient", Desc: "extinction coefficient for target Molecule at the specified wavelength; e.g. 20330 for tartrazine at 472nm\n", Kind: "Parameters"},
				{Name: "FindOptWavelength", Desc: "whether the scan should be used to return the wavelength with maximum signal to noise found\n", Kind: "Parameters"},
				{Name: "ManualComparison", Desc: " Option to compare to manual pipetting\n", Kind: "Parameters"},
				{Name: "MarsResultsFileXLSX", Desc: "input file containing the Plate reader results exported from Mars\n", Kind: "Parameters"},
				{Name: "Molecule", Desc: "The name of the molecule to analyse. This will be used to find matching solutions in the design file and to look up the molecular weight.\nCurrently only one solution name can be run at a time.\n", Kind: "Inputs"},
				{Name: "OutputFilename", Desc: "set the desired name for the output file, if this is blank it will append the design file name with _output\n", Kind: "Parameters"},
				{Name: "OverrideMolecularWeight", Desc: "Option to override moecular weight value of a mpolecule\n", Kind: "Parameters"},
				{Name: "PlateType", Desc: "", Kind: "Inputs"},
				{Name: "R2threshold", Desc: "validation requirements\n\nset a threshold above which R2 will pass; 0 = 0%, 1 = 100%; e.g. 0.7 = 70%\n", Kind: "Parameters"},
				{Name: "ReadingTypeinMarsFile", Desc: "This should match the label in the header for each column in the plate reader result file, e.g. \"Abs Spectrum\"\n", Kind: "Parameters"},
				{Name: "Responsecolumntofill", Desc: "name your response\n", Kind: "Parameters"},
				{Name: "SheetNumber", Desc: "i.e. the sheet position in the plate reader results excel file; starting from 0\n", Kind: "Parameters"},
				{Name: "StockEqualsTotalVolPerWell", Desc: "if true the StockVol represents the total volume per well instead of a fixed volume which the test solution was added to\n", Kind: "Parameters"},
				{Name: "StockconcinMperL", Desc: "", Kind: "Parameters"},
				{Name: "Stockvol", Desc: "volume of diluent per well\n", Kind: "Parameters"},
				{Name: "VolumeToManualwells", Desc: "if comparing to manual pipetting set the wells to use for each concentration here\n", Kind: "Parameters"},
				{Name: "Wavelength", Desc: " Wavelength to use for calculations, should match up with extinction coefficient for molecule of interest\n", Kind: "Parameters"},
				{Name: "WellForScanAnalysis", Desc: "well used for finding wavelength with optimal signal to noise. This is ignored if FindOptWavelength is set to false\n", Kind: "Parameters"},
				{Name: "ActualVsExpectedPlot", Desc: "", Kind: "Data"},
				{Name: "BlankValues", Desc: "", Kind: "Data"},
				{Name: "CV", Desc: "", Kind: "Data"},
				{Name: "CVpass", Desc: "", Kind: "Data"},
				{Name: "CorrectnessFactorPlot", Desc: "", Kind: "Data"},
				{Name: "Errors", Desc: "", Kind: "Data"},
				{Name: "Formula", Desc: "", Kind: "Data"},
				{Name: "MeasuredOptimalWavelength", Desc: "", Kind: "Data"},
				{Name: "OutPutDesignFile", Desc: "", Kind: "Data"},
				{Name: "R2", Desc: "", Kind: "Data"},
				{Name: "R2Pass", Desc: "", Kind: "Data"},
				{Name: "R2_CorrectnessFactor", Desc: "", Kind: "Data"},
				{Name: "ResponsetoManualValuesmap", Desc: "", Kind: "Data"},
				{Name: "Runs", Desc: "", Kind: "Data"},
				{Name: "Variance", Desc: "", Kind: "Data"},
				{Name: "VolumeToActualConc", Desc: "", Kind: "Data"},
				{Name: "VolumeToCorrectnessFactor", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
