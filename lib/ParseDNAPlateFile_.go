// Protocol for parsing the contents of a DNAPlate file in xlsx format from Thermo.
// The products of parsing this plate can the be wired directly into the ResuspendDNA or PairOligos elements.
package lib

import

// we need to import the wtype package to use the LHComponent type
// the mixer package is required to use the Sample function
(
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/doe"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/factory"
	"strings"
)

// Input parameters for this protocol (data)

// supported format is xlsx file from Thermo
// default should be JMP

// using the fwd primer name, returns the rev primer name

func _ParseDNAPlateFileRequirements() {
}

func _ParseDNAPlateFileSetup(_ctx context.Context, _input *ParseDNAPlateFileInput) {
}

func _ParseDNAPlateFileSteps(_ctx context.Context, _input *ParseDNAPlateFileInput, _output *ParseDNAPlateFileOutput) {

	if _input.SequenceInfoFileformat == "" {
		_input.SequenceInfoFileformat = "JMP"
	}

	var ReplaceMap = map[string]string{
		"01": "1",
		"02": "2",
		"03": "3",
		"04": "4",
		"05": "5",
		"06": "6",
		"07": "7",
		"08": "8",
		"09": "9",
		"_":  "",
	}

	// headers
	NameHeader := "Seq Name"
	seqHeader := "Seq"
	MassHeader := "Yield_ug"
	MWHeader := "MW"
	WellHeader := "Customer Well"
	PlateNameHeader := "Plate"
	volHeader := "Volume_ul"
	concHeader := "Conc_ug/ul"

	// initialise maps
	_output.PartMassMap = make(map[string]wunit.Mass)
	_output.PartMolecularWeightMap = make(map[string]float64)
	_output.PartLocationsMap = make(map[string]string)
	_output.PartPlateMap = make(map[string]string)
	headersfound := make([]string, 0)
	_output.FwdOligotoRevOligo = make(map[string]string)
	_output.PartsList = make(map[string]*wtype.LHComponent)

	// get contents from file
	fileContents, err := _input.SequenceInfoFile.ReadAll()

	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	dnaparts, err := doe.RunsFromDesignPreResponsesContents(fileContents, []string{"Length", "MW", "Tm", "Yield"}, _input.SequenceInfoFileformat)

	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	// code for parsing the data from the xl file into the strings, this searches the file in direction i followed by j
	for i, partinfo := range dnaparts {

		var partname string
		var partmass float64
		var partwell string
		var partmw float64
		var platename string
		var partSeq string
		var partConc float64
		var partVol float64

		for j := range partinfo.Setpoints {

			//First creates an array of part names
			if partinfo.Factordescriptors[j] == NameHeader {

				if name, found := partinfo.Setpoints[j].(string); found {

					if name == "" || name == "BLANK" {
						fmt.Print("Skipping ", name)
					} else {
						partname = name
						_output.Partnames = append(_output.Partnames, name)
					}

				} else {
					execute.Errorf(_ctx, fmt.Sprint("wrong type", partinfo.Factordescriptors[j], partinfo.Setpoints[j]))
				}

				if i == 0 {
					headersfound = append(headersfound, NameHeader)
				}

			}

			//second create an array of plasmid masses
			if partinfo.Factordescriptors[j] == MassHeader {

				if mass, found := partinfo.Setpoints[j].(float64); found {
					partmass = mass
				} else if mass, found := partinfo.Setpoints[j].(string); found {
					if mass == "" {
						partmass = 0.0 // mw
						// empty so skip
					}
				} else {
					execute.Errorf(_ctx, fmt.Sprint("wrong type", partinfo.Factordescriptors[j], partinfo.Setpoints[j]))
				}

				if i == 0 {
					headersfound = append(headersfound, MassHeader)
				}
			}

			if partinfo.Factordescriptors[j] == MWHeader {

				if mw, found := partinfo.Setpoints[j].(int); found {
					partmw = float64(mw)
				} else if mw, found := partinfo.Setpoints[j].(float64); found {
					partmw = mw
				} else if mw, found := partinfo.Setpoints[j].(string); found {
					if mw == "" {
						partmw = 0.0 // mw
						// empty so skip
					}
				} else {
					execute.Errorf(_ctx, fmt.Sprint("wrong type", partinfo.Factordescriptors[j], partinfo.Setpoints[j]))
				}

				if i == 0 {
					headersfound = append(headersfound, MWHeader)
				}
			}

			if partinfo.Factordescriptors[j] == seqHeader {

				if seq, found := partinfo.Setpoints[j].(string); found {
					partSeq = seq
				} else {
					execute.Errorf(_ctx, fmt.Sprint("wrong type", partinfo.Factordescriptors[j], partinfo.Setpoints[j]))
				}

				if i == 0 {
					headersfound = append(headersfound, seqHeader)
				}
			}

			if partinfo.Factordescriptors[j] == volHeader {

				if vol, found := partinfo.Setpoints[j].(float64); found {
					partVol = vol
				} else if vol, found := partinfo.Setpoints[j].(string); found {
					if vol == "" {
						partVol = 0.0 // mw
						// empty so skip
					}
				} else {
					execute.Errorf(_ctx, fmt.Sprint("wrong type", partinfo.Factordescriptors[j], partinfo.Setpoints[j]))
				}

				if i == 0 {
					headersfound = append(headersfound, volHeader)
				}
			}

			if partinfo.Factordescriptors[j] == concHeader {

				if conc, found := partinfo.Setpoints[j].(float64); found {
					partConc = conc
				} else if conc, found := partinfo.Setpoints[j].(string); found {
					if conc == "" {
						partConc = 0.0 // mw
						// empty so skip
					}
				} else {
					execute.Errorf(_ctx, fmt.Sprint("wrong type", partinfo.Factordescriptors[j], partinfo.Setpoints[j]))
				}

				if i == 0 {
					headersfound = append(headersfound, concHeader)
				}
			}

			if partinfo.Factordescriptors[j] == WellHeader {

				if well, found := partinfo.Setpoints[j].(string); found {

					for key, value := range ReplaceMap {

						if strings.Contains(well, key) {
							well = strings.Replace(well, key, value, 1)

							break
						}
					}
					partwell = well
				} else {
					execute.Errorf(_ctx, fmt.Sprint("wrong type", partinfo.Factordescriptors[j], partinfo.Setpoints[j]))
				}

				if i == 0 {
					headersfound = append(headersfound, WellHeader)
				}
			}

			//third creates an array of part lengths in bp
			if partinfo.Factordescriptors[j] == PlateNameHeader {

				if plate, found := partinfo.Setpoints[j].(string); found {

					platename = plate
				} else {
					execute.Errorf(_ctx, fmt.Sprint("wrong type", partinfo.Factordescriptors[j], partinfo.Setpoints[j]))
				}
				if i == 0 {
					headersfound = append(headersfound, PlateNameHeader)
				}
			}

			//internal check if there are not 4 headers (as we know there should be 4) return an error telling us which ones were found and which were not
			/*
				if len(headersfound)!= 4 {
					Errorf(fmt.Sprint("Only found these headers in input file: ", headersfound))
				}
			*/
		}

		if partname == "" || partname == "BLANK" {
			fmt.Print("Skipping ", partname)
		} else {
			_output.PartLocationsMap[partname] = partwell
			_output.PartMolecularWeightMap[partname] = partmw
			_output.PartMassMap[partname] = wunit.NewMass(partmass, "ug")
			_output.PartPlateMap[partname] = platename

			part := factory.GetComponentByType("dna_part")
			part.CName = partname

			if partVol > 0.0 {
				part.SetVolume(wunit.NewVolume(partVol, "ul"))
			}

			if partConc > 0.0 {
				part.SetConcentration(wunit.NewConcentration(partConc, "ug/ul"))
			}

			_output.PartsList[partname] = part

			var seq wtype.DNASequence
			seq.Nm = partname
			seq.Seq = partSeq

			_output.PartSeqs[partname] = seq
		}

	}

	for _, partname := range _output.Partnames {

		if !strings.Contains(partname, "_Revers") {
			for _, partname2 := range _output.Partnames {
				if strings.Contains(partname2, "_Revers") && strings.Contains(partname2, partname) {
					_output.FwdOligotoRevOligo[partname] = partname2
					break
				}
			}
		}

	}
	_output.OligoPairs = len(_output.FwdOligotoRevOligo)

	_output.HeadersFound = headersfound
}

func _ParseDNAPlateFileAnalysis(_ctx context.Context, _input *ParseDNAPlateFileInput, _output *ParseDNAPlateFileOutput) {
}

func _ParseDNAPlateFileValidation(_ctx context.Context, _input *ParseDNAPlateFileInput, _output *ParseDNAPlateFileOutput) {
}
func _ParseDNAPlateFileRun(_ctx context.Context, input *ParseDNAPlateFileInput) *ParseDNAPlateFileOutput {
	output := &ParseDNAPlateFileOutput{}
	_ParseDNAPlateFileSetup(_ctx, input)
	_ParseDNAPlateFileSteps(_ctx, input, output)
	_ParseDNAPlateFileAnalysis(_ctx, input, output)
	_ParseDNAPlateFileValidation(_ctx, input, output)
	return output
}

func ParseDNAPlateFileRunSteps(_ctx context.Context, input *ParseDNAPlateFileInput) *ParseDNAPlateFileSOutput {
	soutput := &ParseDNAPlateFileSOutput{}
	output := _ParseDNAPlateFileRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ParseDNAPlateFileNew() interface{} {
	return &ParseDNAPlateFileElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ParseDNAPlateFileInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ParseDNAPlateFileRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ParseDNAPlateFileInput{},
			Out: &ParseDNAPlateFileOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type ParseDNAPlateFileElement struct {
	inject.CheckedRunner
}

type ParseDNAPlateFileInput struct {
	SequenceInfoFile       wtype.File
	SequenceInfoFileformat string
}

type ParseDNAPlateFileOutput struct {
	FwdOligotoRevOligo     map[string]string
	HeadersFound           []string
	OligoPairs             int
	PartLocationsMap       map[string]string
	PartMassMap            map[string]wunit.Mass
	PartMolecularWeightMap map[string]float64
	PartPlateMap           map[string]string
	PartSeqs               map[string]wtype.DNASequence
	Partnames              []string
	PartsList              map[string]*wtype.LHComponent
}

type ParseDNAPlateFileSOutput struct {
	Data struct {
		FwdOligotoRevOligo     map[string]string
		HeadersFound           []string
		OligoPairs             int
		PartLocationsMap       map[string]string
		PartMassMap            map[string]wunit.Mass
		PartMolecularWeightMap map[string]float64
		PartPlateMap           map[string]string
		PartSeqs               map[string]wtype.DNASequence
		Partnames              []string
	}
	Outputs struct {
		PartsList map[string]*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "ParseDNAPlateFile",
		Constructor: ParseDNAPlateFileNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol for parsing the contents of a DNAPlate file in xlsx format from Thermo.\nThe products of parsing this plate can the be wired directly into the ResuspendDNA or PairOligos elements.\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/ResuspendDNA/ParseDNAInputFile.an",
			Params: []component.ParamDesc{
				{Name: "SequenceInfoFile", Desc: "supported format is xlsx file from Thermo\n", Kind: "Parameters"},
				{Name: "SequenceInfoFileformat", Desc: "default should be JMP\n", Kind: "Parameters"},
				{Name: "FwdOligotoRevOligo", Desc: "using the fwd primer name, returns the rev primer name\n", Kind: "Data"},
				{Name: "HeadersFound", Desc: "", Kind: "Data"},
				{Name: "OligoPairs", Desc: "", Kind: "Data"},
				{Name: "PartLocationsMap", Desc: "", Kind: "Data"},
				{Name: "PartMassMap", Desc: "", Kind: "Data"},
				{Name: "PartMolecularWeightMap", Desc: "", Kind: "Data"},
				{Name: "PartPlateMap", Desc: "", Kind: "Data"},
				{Name: "PartSeqs", Desc: "", Kind: "Data"},
				{Name: "Partnames", Desc: "", Kind: "Data"},
				{Name: "PartsList", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
