// example protocol which allows a primitive method for searching the igem registry
// for parts with specified functions or a specified status (e.g. A = available or "Works", or results != none)
// see the igem package ("github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/igem")
// and igem website for more details about how to make the most of this http://parts.igem.org/Registry_API
package lib

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/igem"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"strconv"
	"strings"
)

// Input parameters for this protocol (data)

//    IGem Part Name classifications. Valid options are :
//    "GENERIC":          "BBa_B"	Generic basic parts such as Terminators, DNA, and Ribosome Binding Site
//    "PROTEINCODING":    "BBa_C"	Protein coding parts
//    "REPORTER":         "BBa_E"	Reporter parts
//    "SIGNALLING":       "BBa_F"	Signalling parts
//    "PRIMER":           "BBa_G"	Primer parts
//    "IAPPROJECT":       "BBa_I"	IAP 2003, 2004 project parts
//    "IGEMPROJECT":      "BBa_J"	iGEM project parts
//    "TAG":              "BBa_M"	Tag parts
//    "PROTEINGENERATOR": "BBa_P"	Protein Generator parts
//    "INVERTER":         "BBa_Q"	Inverter parts
//    "REGULATORY":       "BBa_R"	Regulatory parts
//    "INTERMEDIATE":     "BBa_S"	Intermediate parts
//    "CELLSTRAIN":       "BBa_V"	Cell strain parts

// e.g. strong, arsenic, fluorescent, alkane, logic gate

// This should be set to true

// only return parts marked as available in registry

// only return parts marked as working in registry

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// i.e. map[description]list of parts matching description
// i.e. map[biobrickID]description

// Input Requirement specification
func _FindIGemPartsThatRequirements() {

}

// Conditions to run on startup
func _FindIGemPartsThatSetup(_ctx context.Context, _input *FindIGemPartsThatInput) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _FindIGemPartsThatSteps(_ctx context.Context, _input *FindIGemPartsThatInput, _output *FindIGemPartsThatOutput) {

	Parttypes := []string{_input.Parttype}

	var BackupParts []string
	var WorkingBackupParts []string

	// initialise some variables for use later
	var parts [][]string
	OriginalPartMap := make(map[string][]string)
	_output.PartMap = make(map[string][]string)
	_output.BiobrickDescriptions = make(map[string]string)
	var highestrating int
	var parttypemap map[string]string
	partstatus := ""

	if _input.OnlyreturnAvailableParts {
		partstatus = "A"
	}

	// first we'll parse the igem registry based on the short description contained in the fasta header for each part sequence
	for _, parttype := range Parttypes {

		var subparts []string
		var err error

		subparts, parttypemap, err = igem.FilterRegistry(parttype, []string{partstatus}, _input.ExactTypeOnly)

		if err != nil {
			execute.Errorf(_ctx, "Error filtering igem registry: %s", err.Error())
		}

		parts = append(parts, subparts)
		OriginalPartMap[parttype+"_"+partstatus] = subparts
		_output.PartMap[parttype+"_"+partstatus] = subparts
	}

	othercriteria := ""
	if _input.OnlyreturnWorkingparts {
		othercriteria = "WORKS"
	}

	for desc, subparts := range OriginalPartMap {

		partdetails := igem.LookUp(subparts)

		// now we can get detailed information of all of those records to interrogate further
		// this can be slow if there are many parts to check
		// Parts will be added if they contain either description not both
		for i := range _input.Partdescriptions {

			for _, subpart := range subparts {

				// check if key words are in description and that status == "WORKS" if only working parts are desired
				if !_input.MatchAllDescriptions && strings.Contains(strings.ToUpper(partdetails.Description(subpart)), strings.ToUpper(_input.Partdescriptions[i])) &&
					strings.Contains(strings.ToUpper(partdetails.Results(subpart)), strings.ToUpper(othercriteria)) {

					if !search.InSlice(subpart, BackupParts) {
						BackupParts = append(BackupParts, subpart)
					}
					// ensure the highest rated part is returned
					rating, err := strconv.Atoi(partdetails.Rating(subpart))

					if err == nil && rating > highestrating {
						_output.HighestRatedMatch = subpart

						seq := partdetails.Sequence(_output.HighestRatedMatch)

						_output.HighestRatedMatchDNASequence = wtype.MakeLinearDNASequence(_output.HighestRatedMatch, seq)
						highestrating = rating
					}
				} else if _input.MatchAllDescriptions && search.Containsallthings((partdetails.Description(subpart)), _input.Partdescriptions) &&
					strings.Contains(partdetails.Results(subpart), othercriteria) {

					if !search.InSlice(subpart, BackupParts) {
						BackupParts = append(BackupParts, subpart)
					}
					// ensure the highest rated part is returned
					rating, err := strconv.Atoi(partdetails.Rating(subpart))

					if err == nil && rating > highestrating {
						_output.HighestRatedMatch = subpart

						seq := partdetails.Sequence(_output.HighestRatedMatch)

						_output.HighestRatedMatchDNASequence = wtype.MakeLinearDNASequence(_output.HighestRatedMatch, seq)
						highestrating = rating
					}
				}
				if !_input.MatchAllDescriptions && strings.Contains(strings.ToUpper(partdetails.Description(subpart)), strings.ToUpper(_input.Partdescriptions[i])) &&
					strings.Contains(partdetails.Results(subpart), "WORKS") {
					if !search.InSlice(subpart, WorkingBackupParts) {
						WorkingBackupParts = append(WorkingBackupParts, subpart)
					}
				} else if _input.MatchAllDescriptions && search.Containsallthings((partdetails.Description(subpart)), _input.Partdescriptions) &&
					strings.Contains(partdetails.Results(subpart), "WORKS") {
					if !search.InSlice(subpart, WorkingBackupParts) {
						WorkingBackupParts = append(WorkingBackupParts, subpart)
					}
				}
				// add to look up table to report back to user
				if _input.MatchAllDescriptions {
					var partdesc string

					for _, descriptor := range _input.Partdescriptions {
						partdesc = partdesc + "_" + descriptor
					}
					_output.PartMap[desc+"_"+partdesc] = BackupParts
					_output.PartMap[desc+"_"+partdesc+"+WORKS"] = WorkingBackupParts

				} else if !_input.MatchAllDescriptions {

					_output.PartMap[desc+"_"+_input.Partdescriptions[i]] = BackupParts
					_output.PartMap[desc+"_"+_input.Partdescriptions[i]+"+WORKS"] = WorkingBackupParts
				}

			}
			for _, part := range WorkingBackupParts {
				_output.Partslist = append(_output.Partslist, part)
			}

			// remove duplicates
			_output.Partslist = search.RemoveDuplicates(_output.Partslist)

			// reset
			//FulllistBackupParts = BackupParts
			BackupParts = make([]string, 0)
			WorkingBackupParts = make([]string, 0)

			//i = i + 1
			if _input.MatchAllDescriptions {
				// don't need to loop through each description if we're matching all
				continue
			}
		}
		for _, subpartarray := range _output.PartMap {
			for _, subpart := range subpartarray {
				if partdetails.Description(subpart) != "" {
					_output.BiobrickDescriptions[subpart] = partdetails.Description(subpart)
				} else {
					_output.BiobrickDescriptions[subpart] = parttypemap[subpart]
				}
			}
		}
	}

	_output.HighestRatedMatchScore = highestrating

	// print in pretty format on terminal
	for key, value := range _output.PartMap {
		fmt.Println(text.Print(key, value))
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _FindIGemPartsThatAnalysis(_ctx context.Context, _input *FindIGemPartsThatInput, _output *FindIGemPartsThatOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _FindIGemPartsThatValidation(_ctx context.Context, _input *FindIGemPartsThatInput, _output *FindIGemPartsThatOutput) {
}
func _FindIGemPartsThatRun(_ctx context.Context, input *FindIGemPartsThatInput) *FindIGemPartsThatOutput {
	output := &FindIGemPartsThatOutput{}
	_FindIGemPartsThatSetup(_ctx, input)
	_FindIGemPartsThatSteps(_ctx, input, output)
	_FindIGemPartsThatAnalysis(_ctx, input, output)
	_FindIGemPartsThatValidation(_ctx, input, output)
	return output
}

func FindIGemPartsThatRunSteps(_ctx context.Context, input *FindIGemPartsThatInput) *FindIGemPartsThatSOutput {
	soutput := &FindIGemPartsThatSOutput{}
	output := _FindIGemPartsThatRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func FindIGemPartsThatNew() interface{} {
	return &FindIGemPartsThatElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &FindIGemPartsThatInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _FindIGemPartsThatRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &FindIGemPartsThatInput{},
			Out: &FindIGemPartsThatOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type FindIGemPartsThatElement struct {
	inject.CheckedRunner
}

type FindIGemPartsThatInput struct {
	ExactTypeOnly            bool
	MatchAllDescriptions     bool
	OnlyreturnAvailableParts bool
	OnlyreturnWorkingparts   bool
	Partdescriptions         []string
	Parttype                 string
}

type FindIGemPartsThatOutput struct {
	BiobrickDescriptions         map[string]string
	HighestRatedMatch            string
	HighestRatedMatchDNASequence wtype.DNASequence
	HighestRatedMatchScore       int
	PartMap                      map[string][]string
	Partslist                    []string
	Warnings                     error
}

type FindIGemPartsThatSOutput struct {
	Data struct {
		BiobrickDescriptions         map[string]string
		HighestRatedMatch            string
		HighestRatedMatchDNASequence wtype.DNASequence
		HighestRatedMatchScore       int
		PartMap                      map[string][]string
		Partslist                    []string
		Warnings                     error
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "FindIGemPartsThat",
		Constructor: FindIGemPartsThatNew,
		Desc: component.ComponentDesc{
			Desc: "example protocol which allows a primitive method for searching the igem registry\nfor parts with specified functions or a specified status (e.g. A = available or \"Works\", or results != none)\nsee the igem package (\"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/igem\")\nand igem website for more details about how to make the most of this http://parts.igem.org/Registry_API\n",
			Path: "src/github.com/antha-lang/elements/an/Data/DNA/FindIGemPartsThat/FindIGemPartsThat.an",
			Params: []component.ParamDesc{
				{Name: "ExactTypeOnly", Desc: "This should be set to true\n", Kind: "Parameters"},
				{Name: "MatchAllDescriptions", Desc: "", Kind: "Parameters"},
				{Name: "OnlyreturnAvailableParts", Desc: "only return parts marked as available in registry\n", Kind: "Parameters"},
				{Name: "OnlyreturnWorkingparts", Desc: "only return parts marked as working in registry\n", Kind: "Parameters"},
				{Name: "Partdescriptions", Desc: "e.g. strong, arsenic, fluorescent, alkane, logic gate\n", Kind: "Parameters"},
				{Name: "Parttype", Desc: "   IGem Part Name classifications. Valid options are :\n   \"GENERIC\":          \"BBa_B\"\tGeneric basic parts such as Terminators, DNA, and Ribosome Binding Site\n   \"PROTEINCODING\":    \"BBa_C\"\tProtein coding parts\n   \"REPORTER\":         \"BBa_E\"\tReporter parts\n   \"SIGNALLING\":       \"BBa_F\"\tSignalling parts\n   \"PRIMER\":           \"BBa_G\"\tPrimer parts\n   \"IAPPROJECT\":       \"BBa_I\"\tIAP 2003, 2004 project parts\n   \"IGEMPROJECT\":      \"BBa_J\"\tiGEM project parts\n   \"TAG\":              \"BBa_M\"\tTag parts\n   \"PROTEINGENERATOR\": \"BBa_P\"\tProtein Generator parts\n   \"INVERTER\":         \"BBa_Q\"\tInverter parts\n   \"REGULATORY\":       \"BBa_R\"\tRegulatory parts\n   \"INTERMEDIATE\":     \"BBa_S\"\tIntermediate parts\n   \"CELLSTRAIN\":       \"BBa_V\"\tCell strain parts\n", Kind: "Parameters"},
				{Name: "BiobrickDescriptions", Desc: "i.e. map[biobrickID]description\n", Kind: "Data"},
				{Name: "HighestRatedMatch", Desc: "", Kind: "Data"},
				{Name: "HighestRatedMatchDNASequence", Desc: "", Kind: "Data"},
				{Name: "HighestRatedMatchScore", Desc: "", Kind: "Data"},
				{Name: "PartMap", Desc: "i.e. map[description]list of parts matching description\n", Kind: "Data"},
				{Name: "Partslist", Desc: "", Kind: "Data"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
