// This protocol is intended to design a combinatorial library of all combinations of a list of Vectors, Promoters,
// RBSs, CDSs and Terminators according to an assembly standard ensuring compatibility with level 1 design.
// Level 1 adaptor sites (containing the correct restriction site are expected to be included in the promoter and terminator parts.
// These level 1 sites can be designed such that a series of level 1 parts may be joined together in a second assembly reaction.
// A list of sequencing primers to order will also be returned.
package lib

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/export"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/oligos"
	"github.com/antha-lang/antha/antha/anthalib/wtype"

	"context"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"path/filepath"
)

// Input parameters for this protocol (data)

// Custom design, may support MoClo, EcoFlex and GoldenBraid.

// default is "level0"

// Option to add Level 1 adaptor sites to the Promoters and terminators to support hierarchical assembly
// If Custom design the valid options currently supported are: "Device1","Device2", "Device3".
// If left empty no adaptor sequence is added.

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// parts to order
// parts + vector map ready for feeding into downstream AutoAssembly element

// desired sequence to end up with after assembly

// Input Requirement specification
func _CombinatorialLibraryDesign_From_Assembly_Standard_5PartRequirements() {
	// e.g. are MoClo types valid?
}

// Conditions to run on startup
func _CombinatorialLibraryDesign_From_Assembly_Standard_5PartSetup(_ctx context.Context, _input *CombinatorialLibraryDesign_From_Assembly_Standard_5PartInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _CombinatorialLibraryDesign_From_Assembly_Standard_5PartSteps(_ctx context.Context, _input *CombinatorialLibraryDesign_From_Assembly_Standard_5PartInput, _output *CombinatorialLibraryDesign_From_Assembly_Standard_5PartOutput) {
	_output.StatusMap = make(map[string]string)
	_output.PartswithOverhangsMap = make(map[string][]wtype.DNASequence) // parts to order
	_output.Assemblies = make(map[string][]wtype.DNASequence)
	_output.PassMap = make(map[string]bool)
	_output.SeqsMap = make(map[string]wtype.DNASequence) // desired sequence to end up with after assembly
	_output.EndreportMap = make(map[string]string)
	_output.PositionReportMap = make(map[string][]string)
	_output.StatusMap = make(map[string]string)
	_output.PrimerMap = make(map[string]oligos.Primer)

	_output.Sequences = make([]wtype.DNASequence, 0)
	_output.Parts = make([][]wtype.DNASequence, 0)
	_output.SequencingPrimers = make([][]wtype.DNASequence, 0)

	var standardLabels []string

	for _, level := range _input.StandardLabels {
		standardLabels = append(standardLabels, level)
	}

	var counter int = 1

	if _input.StandardLevel == "" {
		_input.StandardLevel = "Level0"
	}

	for j := range _input.Vectors {
		for k := range _input.Part1s {
			for l := range _input.Part2s {
				for m := range _input.Part3s {
					for n := range _input.Part4s {
						key := _input.ProjectName + _input.Vectors[j].Nm + "_" + _input.Part1s[k].Nm + "_" + _input.Part2s[l].Nm + "_" + _input.Part3s[m].Nm + _input.Part4s[n].Nm
						assembly := AssemblyStandard_TypeIIsDesignRunSteps(_ctx, &AssemblyStandard_TypeIIsDesignInput{Constructname: key,
							Seqsinorder:                   []wtype.DNASequence{_input.Part1s[k], _input.Part2s[l], _input.Part3s[m], _input.Part4s[n]},
							AssemblyStandard:              _input.Standard,
							Level:                         _input.StandardLevel, // of assembly standard
							Vector:                        _input.Vectors[j],
							PartMoClotypesinorder:         standardLabels,
							OtherEnzymeSitesToRemove:      _input.SitesToRemove,
							ORFstoConfirm:                 []string{}, // enter each as amino acid sequence
							RemoveproblemRestrictionSites: true,
							OnlyRemovesitesinORFs:         false,
							EndsAlreadyadded:              _input.EndsAlreadyAdded,
							ExporttoFastaFile:             _input.FolderPerConstruct,
							BlastSeqswithNoName:           _input.BlastSearchSeqs,
							MakeLevel1Device:              _input.MakeLevel1Device,
							ReverseLevel1Orientation:      _input.ReverseLevel1Orientation},
						)
						key = key                                                             //+ Vectors[j]
						_output.PartswithOverhangsMap[key] = assembly.Data.PartswithOverhangs // parts to order
						_output.Assemblies[key] = assembly.Data.PartsAndVector                // parts + vector to be fed into assembly element
						_output.Parts = append(_output.Parts, assembly.Data.PartswithOverhangs)
						_output.PassMap[key] = assembly.Data.Simulationpass
						_output.EndreportMap[key] = assembly.Data.Endreport
						_output.PositionReportMap[key] = assembly.Data.PositionReport
						_output.SeqsMap[key] = assembly.Data.NewDNASequence
						_output.Sequences = append(_output.Sequences, assembly.Data.NewDNASequence)
						_output.StatusMap[key] = assembly.Data.Status

						// for each vector we'll also design sequencing primers

						primer := PrimerDesign_ColonyPCR_wtypeRunSteps(_ctx, &PrimerDesign_ColonyPCR_wtypeInput{FullDNASeq: assembly.Data.NewDNASequence,
							Maxtemp:                                  wunit.NewTemperature(72, "C"),
							Mintemp:                                  wunit.NewTemperature(50, "C"),
							Maxgc:                                    0.7,
							Minlength:                                12,
							Maxlength:                                30,
							Seqstoavoid:                              []string{},
							PermittednucleotideOverlapBetweenPrimers: 10,                   // number of nucleotides which primers can overlap by
							RegionSequence:                           assembly.Data.Insert, // first part
							FlankTargetSequence:                      true},
						)

						// rename primers
						primer.Data.FWDPrimer.Nm = primer.Data.FWDPrimer.Nm + _input.ProjectName + _input.Vectors[j].Nm + "_FWD"
						primer.Data.REVPrimer.Nm = primer.Data.REVPrimer.Nm + _input.ProjectName + _input.Vectors[j].Nm + "_REV"

						_output.PrimerMap[key+"_FWD"] = primer.Data.FWDPrimer
						_output.PrimerMap[key+"_REV"] = primer.Data.REVPrimer
						_output.SequencingPrimers = append(_output.SequencingPrimers, []wtype.DNASequence{primer.Data.FWDPrimer.DNASequence, primer.Data.REVPrimer.DNASequence})
						counter++
					}
				}
			}
		}
	}

	// export sequence to fasta

	var err error
	// export simulated sequences to file
	_output.AssembledSequences, _, err = export.FastaSerial(export.LOCAL, filepath.Join(_input.ProjectName, "AssembledSequences"), _output.Sequences)

	if err != nil {
		execute.Errorf(_ctx, "Error exporting sequence file for %s: %s", _input.ProjectName, err.Error())
	}
	// add fasta file for each set of parts with overhangs

	labels := standardLabels

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

	for _, parts := range _output.SequencingPrimers {

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

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _CombinatorialLibraryDesign_From_Assembly_Standard_5PartAnalysis(_ctx context.Context, _input *CombinatorialLibraryDesign_From_Assembly_Standard_5PartInput, _output *CombinatorialLibraryDesign_From_Assembly_Standard_5PartOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _CombinatorialLibraryDesign_From_Assembly_Standard_5PartValidation(_ctx context.Context, _input *CombinatorialLibraryDesign_From_Assembly_Standard_5PartInput, _output *CombinatorialLibraryDesign_From_Assembly_Standard_5PartOutput) {
}
func _CombinatorialLibraryDesign_From_Assembly_Standard_5PartRun(_ctx context.Context, input *CombinatorialLibraryDesign_From_Assembly_Standard_5PartInput) *CombinatorialLibraryDesign_From_Assembly_Standard_5PartOutput {
	output := &CombinatorialLibraryDesign_From_Assembly_Standard_5PartOutput{}
	_CombinatorialLibraryDesign_From_Assembly_Standard_5PartSetup(_ctx, input)
	_CombinatorialLibraryDesign_From_Assembly_Standard_5PartSteps(_ctx, input, output)
	_CombinatorialLibraryDesign_From_Assembly_Standard_5PartAnalysis(_ctx, input, output)
	_CombinatorialLibraryDesign_From_Assembly_Standard_5PartValidation(_ctx, input, output)
	return output
}

func CombinatorialLibraryDesign_From_Assembly_Standard_5PartRunSteps(_ctx context.Context, input *CombinatorialLibraryDesign_From_Assembly_Standard_5PartInput) *CombinatorialLibraryDesign_From_Assembly_Standard_5PartSOutput {
	soutput := &CombinatorialLibraryDesign_From_Assembly_Standard_5PartSOutput{}
	output := _CombinatorialLibraryDesign_From_Assembly_Standard_5PartRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func CombinatorialLibraryDesign_From_Assembly_Standard_5PartNew() interface{} {
	return &CombinatorialLibraryDesign_From_Assembly_Standard_5PartElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &CombinatorialLibraryDesign_From_Assembly_Standard_5PartInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _CombinatorialLibraryDesign_From_Assembly_Standard_5PartRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &CombinatorialLibraryDesign_From_Assembly_Standard_5PartInput{},
			Out: &CombinatorialLibraryDesign_From_Assembly_Standard_5PartOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type CombinatorialLibraryDesign_From_Assembly_Standard_5PartElement struct {
	inject.CheckedRunner
}

type CombinatorialLibraryDesign_From_Assembly_Standard_5PartInput struct {
	BlastSearchSeqs          bool
	EndsAlreadyAdded         bool
	FolderPerConstruct       bool
	MakeLevel1Device         string
	Part1s                   []wtype.DNASequence
	Part2s                   []wtype.DNASequence
	Part3s                   []wtype.DNASequence
	Part4s                   []wtype.DNASequence
	ProjectName              string
	ReverseLevel1Orientation bool
	SitesToRemove            []string
	Standard                 enzymes.AssemblyStandard
	StandardLabels           [4]string
	StandardLevel            string
	Vectors                  []wtype.DNASequence
}

type CombinatorialLibraryDesign_From_Assembly_Standard_5PartOutput struct {
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
	SequencingPrimers     [][]wtype.DNASequence
	StatusMap             map[string]string
}

type CombinatorialLibraryDesign_From_Assembly_Standard_5PartSOutput struct {
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
		SequencingPrimers     [][]wtype.DNASequence
		StatusMap             map[string]string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "CombinatorialLibraryDesign_From_Assembly_Standard_5Part",
		Constructor: CombinatorialLibraryDesign_From_Assembly_Standard_5PartNew,
		Desc: component.ComponentDesc{
			Desc: "This protocol is intended to design a combinatorial library of all combinations of a list of Vectors, Promoters,\nRBSs, CDSs and Terminators according to an assembly standard ensuring compatibility with level 1 design.\nLevel 1 adaptor sites (containing the correct restriction site are expected to be included in the promoter and terminator parts.\nThese level 1 sites can be designed such that a series of level 1 parts may be joined together in a second assembly reaction.\nA list of sequencing primers to order will also be returned.\n",
			Path: "src/github.com/antha-lang/elements/an/Data/DNA/TypeIISAssembly_design/CombinatorialDesign/MoClo/Hierarchical/CombinatorialLibraryDesign_From_Assembly_Standard_5Part.an",
			Params: []component.ParamDesc{
				{Name: "BlastSearchSeqs", Desc: "", Kind: "Parameters"},
				{Name: "EndsAlreadyAdded", Desc: "", Kind: "Parameters"},
				{Name: "FolderPerConstruct", Desc: "", Kind: "Parameters"},
				{Name: "MakeLevel1Device", Desc: "Option to add Level 1 adaptor sites to the Promoters and terminators to support hierarchical assembly\nIf Custom design the valid options currently supported are: \"Device1\",\"Device2\", \"Device3\".\nIf left empty no adaptor sequence is added.\n", Kind: "Parameters"},
				{Name: "Part1s", Desc: "", Kind: "Parameters"},
				{Name: "Part2s", Desc: "", Kind: "Parameters"},
				{Name: "Part3s", Desc: "", Kind: "Parameters"},
				{Name: "Part4s", Desc: "", Kind: "Parameters"},
				{Name: "ProjectName", Desc: "", Kind: "Parameters"},
				{Name: "ReverseLevel1Orientation", Desc: "", Kind: "Parameters"},
				{Name: "SitesToRemove", Desc: "", Kind: "Parameters"},
				{Name: "Standard", Desc: "Custom design, may support MoClo, EcoFlex and GoldenBraid.\n", Kind: "Parameters"},
				{Name: "StandardLabels", Desc: "", Kind: "Parameters"},
				{Name: "StandardLevel", Desc: "default is \"level0\"\n", Kind: "Parameters"},
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
				{Name: "SequencingPrimers", Desc: "", Kind: "Data"},
				{Name: "StatusMap", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
