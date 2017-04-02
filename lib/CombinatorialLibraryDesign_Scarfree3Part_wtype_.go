// This protocol is intended to design a combinatorial library of all combinations of lists of options for 3 parts plus vectors.
// Overhangs are added to complement the adjacent parts and leave no scar according to a specified TypeIIs Restriction Enzyme.
// If assembly simulation fails after overhangs are added. In order to help the user
// diagnose the reason, a report of the part overhangs
// is returned to the user along with a list of cut sites in each part.
package lib

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/export"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/oligos"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"path/filepath"
)

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// parts to order
// parts + vector map ready for feeding into downstream AutoAssembly element

// desired sequence to end up with after assembly

// Input Requirement specification
func _CombinatorialLibraryDesign_Scarfree3Part_wtypeRequirements() {
	// e.g. are MoClo types valid?
}

// Conditions to run on startup
func _CombinatorialLibraryDesign_Scarfree3Part_wtypeSetup(_ctx context.Context, _input *CombinatorialLibraryDesign_Scarfree3Part_wtypeInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _CombinatorialLibraryDesign_Scarfree3Part_wtypeSteps(_ctx context.Context, _input *CombinatorialLibraryDesign_Scarfree3Part_wtypeInput, _output *CombinatorialLibraryDesign_Scarfree3Part_wtypeOutput) {
	_output.StatusMap = make(map[string]string)
	_output.PartswithOverhangsMap = make(map[string][]wtype.DNASequence) // parts to order
	_output.Assemblies = make(map[string][]wtype.DNASequence)
	_output.PassMap = make(map[string]bool)
	_output.SeqsMap = make(map[string]wtype.DNASequence) // desired sequence to end up with after assembly
	_output.Sequences = make([]wtype.DNASequence, 0)
	_output.Parts = make([][]wtype.DNASequence, 0)
	SequencingPrimers := make([][]wtype.DNASequence, 0)
	_output.EndreportMap = make(map[string]string)
	_output.PositionReportMap = make(map[string][]string)
	_output.StatusMap = make(map[string]string)
	_output.PrimerMap = make(map[string]oligos.Primer)

	var counter int = 1

	for j := range _input.Vectors {
		for k := range _input.Part1s {
			for l := range _input.Part2s {
				for m := range _input.Part3s {

					key := _input.ProjectName + "_" + _input.Vectors[j].Name() + "_" + _input.Part1s[k].Name() + "_" + _input.Part2s[l].Name() + "_" + _input.Part3s[m].Name()
					assembly := Scarfree_siteremove_orfcheck_wtypeRunSteps(_ctx, &Scarfree_siteremove_orfcheck_wtypeInput{Constructname: key,
						Seqsinorder: []wtype.DNASequence{_input.Part1s[k], _input.Part2s[l], _input.Part3s[m]},
						Enzymename:  _input.EnzymeName,
						Vector:      _input.Vectors[j],
						OtherEnzymeSitesToRemove:      _input.SitesToRemove,
						ORFstoConfirm:                 _input.ORFStoconfirm, // enter each as amino acid sequence
						RemoveproblemRestrictionSites: _input.RemoveproblemRestrictionSites,
						EndsAlreadyadded:              _input.EndsAlreadyadded,
						ExporttoFastaFile:             _input.FolderPerConstruct,
						BlastSeqswithNoName:           _input.BlastSearchSeqs},
					)
					key = key                                                             //+ Vectors[j]
					_output.PartswithOverhangsMap[key] = assembly.Data.PartswithOverhangs // parts to order
					_output.Assemblies[key] = assembly.Data.PartsAndVector                // parts + vector to be fed into assembly element
					_output.PassMap[key] = assembly.Data.Simulationpass
					_output.EndreportMap[key] = assembly.Data.Endreport
					_output.PositionReportMap[key] = assembly.Data.PositionReport
					_output.SeqsMap[key] = assembly.Data.NewDNASequence
					_output.Sequences = append(_output.Sequences, assembly.Data.NewDNASequence)
					_output.Parts = append(_output.Parts, assembly.Data.PartswithOverhangs)
					_output.StatusMap[key] = assembly.Data.Status

					// for each vector we'll also design sequencing primers

					// check binding. Exact matches only.
					sites := sequences.FindSeqsinSeqs(assembly.Data.NewDNASequence.Sequence(), []string{assembly.Data.Insert.Sequence()})

					if len(sites) != 1 {
						_output.Warnings[key] = fmt.Errorf("Found %d Insert %s sites in full assembled sequence %s", len(sites[0].Positions), assembly.Data.Insert.Sequence(), assembly.Data.NewDNASequence.Sequence())
					}

					primer := PrimerDesign_ColonyPCR_wtypeRunSteps(_ctx, &PrimerDesign_ColonyPCR_wtypeInput{FullDNASeq: assembly.Data.NewDNASequence,
						Maxtemp:                                  wunit.NewTemperature(72, "C"),
						Mintemp:                                  wunit.NewTemperature(50, "C"),
						Maxgc:                                    0.7,
						Minlength:                                15,
						Maxlength:                                35,
						Seqstoavoid:                              []string{},
						PermittednucleotideOverlapBetweenPrimers: 10, // number of nucleotides which primers can overlap by
						RegionSequence:                           assembly.Data.Insert,
						FlankTargetSequence:                      true},
					)

					// rename primers
					primer.Data.FWDPrimer.Nm = primer.Data.FWDPrimer.Nm + _input.ProjectName + _input.Vectors[j].Nm + "_FWD"
					primer.Data.REVPrimer.Nm = primer.Data.REVPrimer.Nm + _input.ProjectName + _input.Vectors[j].Nm + "_REV"

					_output.PrimerMap[key+"_FWD"] = primer.Data.FWDPrimer
					_output.PrimerMap[key+"_REV"] = primer.Data.REVPrimer
					SequencingPrimers = append(SequencingPrimers, []wtype.DNASequence{primer.Data.FWDPrimer.DNASequence, primer.Data.REVPrimer.DNASequence})
					counter++
				}
			}
		}

		// export sequence to fasta
		if _input.FolderPerProject {

			var err error
			// export simulated sequences to file
			_output.AssembledSequences, _, err = export.FastaSerial(export.LOCAL, filepath.Join(_input.ProjectName, "AssembledSequences"), _output.Sequences)

			if err != nil {
				execute.Errorf(_ctx, "Error exporting sequence file for %s: %s", _input.ProjectName, err.Error())
			}
			// add fasta file for each set of parts with overhangs
			labels := []string{"Part1s", "Part2s", "Part3s"}

			refactoredparts := make(map[string][]wtype.DNASequence)

			newparts := make([]wtype.DNASequence, 0)

			for _, parts := range _output.Parts {

				for j := range parts {
					newparts = refactoredparts[labels[j]]
					newparts = append(newparts, parts[j])
					refactoredparts[labels[j]] = newparts
				}
			}

			for key, value := range refactoredparts {

				duplicateremoved := search.RemoveDuplicateSequences(value)

				file, _, err := export.FastaSerial(export.LOCAL, filepath.Join(_input.ProjectName, key), duplicateremoved)

				if err != nil {
					execute.Errorf(_ctx, "Error exporting parts to order file for %s %s: %s", _input.ProjectName, key, err.Error())
				}

				_output.PartsToOrder = append(_output.PartsToOrder, file)
			}

			// add fasta file for each set of primers
			labels = []string{"FWDPrimers", "REVPrimers"}

			refactoredparts = make(map[string][]wtype.DNASequence)

			newparts = make([]wtype.DNASequence, 0)

			for _, parts := range SequencingPrimers {

				for j := range parts {
					newparts = refactoredparts[labels[j]]
					newparts = append(newparts, parts[j])
					refactoredparts[labels[j]] = newparts
				}
			}

			for key, value := range refactoredparts {

				duplicateremoved := search.RemoveDuplicateSequences(value)

				primerFile, _, err := export.FastaSerial(export.LOCAL, filepath.Join(_input.ProjectName, key), duplicateremoved)

				if err != nil {
					execute.Errorf(_ctx, "Error exporting primers to order file for %s %s: %s", _input.ProjectName, key, err.Error())
				}

				_output.PrimersToOrder = append(_output.PrimersToOrder, primerFile)
			}

		}

	}
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _CombinatorialLibraryDesign_Scarfree3Part_wtypeAnalysis(_ctx context.Context, _input *CombinatorialLibraryDesign_Scarfree3Part_wtypeInput, _output *CombinatorialLibraryDesign_Scarfree3Part_wtypeOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _CombinatorialLibraryDesign_Scarfree3Part_wtypeValidation(_ctx context.Context, _input *CombinatorialLibraryDesign_Scarfree3Part_wtypeInput, _output *CombinatorialLibraryDesign_Scarfree3Part_wtypeOutput) {
}
func _CombinatorialLibraryDesign_Scarfree3Part_wtypeRun(_ctx context.Context, input *CombinatorialLibraryDesign_Scarfree3Part_wtypeInput) *CombinatorialLibraryDesign_Scarfree3Part_wtypeOutput {
	output := &CombinatorialLibraryDesign_Scarfree3Part_wtypeOutput{}
	_CombinatorialLibraryDesign_Scarfree3Part_wtypeSetup(_ctx, input)
	_CombinatorialLibraryDesign_Scarfree3Part_wtypeSteps(_ctx, input, output)
	_CombinatorialLibraryDesign_Scarfree3Part_wtypeAnalysis(_ctx, input, output)
	_CombinatorialLibraryDesign_Scarfree3Part_wtypeValidation(_ctx, input, output)
	return output
}

func CombinatorialLibraryDesign_Scarfree3Part_wtypeRunSteps(_ctx context.Context, input *CombinatorialLibraryDesign_Scarfree3Part_wtypeInput) *CombinatorialLibraryDesign_Scarfree3Part_wtypeSOutput {
	soutput := &CombinatorialLibraryDesign_Scarfree3Part_wtypeSOutput{}
	output := _CombinatorialLibraryDesign_Scarfree3Part_wtypeRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func CombinatorialLibraryDesign_Scarfree3Part_wtypeNew() interface{} {
	return &CombinatorialLibraryDesign_Scarfree3Part_wtypeElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &CombinatorialLibraryDesign_Scarfree3Part_wtypeInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _CombinatorialLibraryDesign_Scarfree3Part_wtypeRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &CombinatorialLibraryDesign_Scarfree3Part_wtypeInput{},
			Out: &CombinatorialLibraryDesign_Scarfree3Part_wtypeOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type CombinatorialLibraryDesign_Scarfree3Part_wtypeElement struct {
	inject.CheckedRunner
}

type CombinatorialLibraryDesign_Scarfree3Part_wtypeInput struct {
	BlastSearchSeqs               bool
	EndsAlreadyadded              bool
	EnzymeName                    string
	FolderPerConstruct            bool
	FolderPerProject              bool
	ORFStoconfirm                 []string
	Part1s                        []wtype.DNASequence
	Part2s                        []wtype.DNASequence
	Part3s                        []wtype.DNASequence
	ProjectName                   string
	RemoveproblemRestrictionSites bool
	SitesToRemove                 []string
	Vectors                       []wtype.DNASequence
}

type CombinatorialLibraryDesign_Scarfree3Part_wtypeOutput struct {
	AssembledSequences    wtype.File
	Assemblies            map[string][]wtype.DNASequence
	EndreportMap          map[string]string
	Parts                 [][]wtype.DNASequence
	PartsToOrder          []wtype.File
	PartswithOverhangsMap map[string][]wtype.DNASequence
	PassMap               map[string]bool
	PositionReportMap     map[string][]string
	PrimerMap             map[string]oligos.Primer
	PrimersToOrder        []wtype.File
	SeqsMap               map[string]wtype.DNASequence
	Sequences             []wtype.DNASequence
	StatusMap             map[string]string
	Warnings              map[string]error
}

type CombinatorialLibraryDesign_Scarfree3Part_wtypeSOutput struct {
	Data struct {
		AssembledSequences    wtype.File
		Assemblies            map[string][]wtype.DNASequence
		EndreportMap          map[string]string
		Parts                 [][]wtype.DNASequence
		PartsToOrder          []wtype.File
		PartswithOverhangsMap map[string][]wtype.DNASequence
		PassMap               map[string]bool
		PositionReportMap     map[string][]string
		PrimerMap             map[string]oligos.Primer
		PrimersToOrder        []wtype.File
		SeqsMap               map[string]wtype.DNASequence
		Sequences             []wtype.DNASequence
		StatusMap             map[string]string
		Warnings              map[string]error
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "CombinatorialLibraryDesign_Scarfree3Part_wtype",
		Constructor: CombinatorialLibraryDesign_Scarfree3Part_wtypeNew,
		Desc: component.ComponentDesc{
			Desc: "This protocol is intended to design a combinatorial library of all combinations of lists of options for 3 parts plus vectors.\nOverhangs are added to complement the adjacent parts and leave no scar according to a specified TypeIIs Restriction Enzyme.\nIf assembly simulation fails after overhangs are added. In order to help the user\ndiagnose the reason, a report of the part overhangs\nis returned to the user along with a list of cut sites in each part.\n",
			Path: "src/github.com/antha-lang/elements/an/Data/DNA/TypeIISAssembly_design/CombinatorialDesign/scarfree/CombinatorialLibraryDesign_Scarfree_wtype.an",
			Params: []component.ParamDesc{
				{Name: "BlastSearchSeqs", Desc: "", Kind: "Parameters"},
				{Name: "EndsAlreadyadded", Desc: "", Kind: "Parameters"},
				{Name: "EnzymeName", Desc: "", Kind: "Parameters"},
				{Name: "FolderPerConstruct", Desc: "", Kind: "Parameters"},
				{Name: "FolderPerProject", Desc: "", Kind: "Parameters"},
				{Name: "ORFStoconfirm", Desc: "", Kind: "Parameters"},
				{Name: "Part1s", Desc: "", Kind: "Parameters"},
				{Name: "Part2s", Desc: "", Kind: "Parameters"},
				{Name: "Part3s", Desc: "", Kind: "Parameters"},
				{Name: "ProjectName", Desc: "", Kind: "Parameters"},
				{Name: "RemoveproblemRestrictionSites", Desc: "", Kind: "Parameters"},
				{Name: "SitesToRemove", Desc: "", Kind: "Parameters"},
				{Name: "Vectors", Desc: "", Kind: "Parameters"},
				{Name: "AssembledSequences", Desc: "", Kind: "Data"},
				{Name: "Assemblies", Desc: "parts + vector map ready for feeding into downstream AutoAssembly element\n", Kind: "Data"},
				{Name: "EndreportMap", Desc: "", Kind: "Data"},
				{Name: "Parts", Desc: "", Kind: "Data"},
				{Name: "PartsToOrder", Desc: "", Kind: "Data"},
				{Name: "PartswithOverhangsMap", Desc: "parts to order\n", Kind: "Data"},
				{Name: "PassMap", Desc: "", Kind: "Data"},
				{Name: "PositionReportMap", Desc: "", Kind: "Data"},
				{Name: "PrimerMap", Desc: "", Kind: "Data"},
				{Name: "PrimersToOrder", Desc: "", Kind: "Data"},
				{Name: "SeqsMap", Desc: "desired sequence to end up with after assembly\n", Kind: "Data"},
				{Name: "Sequences", Desc: "", Kind: "Data"},
				{Name: "StatusMap", Desc: "", Kind: "Data"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
