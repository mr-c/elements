// This protocol is intended to design a combinatorial library of all combinations of a list of Vectors, Promoters,
// RBSs, CDSs and Terminators according to an assembly standard.
// A list of sequencing primers to order will also be returned.
package lib

import (
	"context"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/oligos"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol (data)

// Custom
// level of assembly standard options are: Level0, Level1

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// parts to order
// parts + vector map ready for feeding into downstream AutoAssembly element

// desired sequence to end up with after assembly

// Input Requirement specification
func _CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapRequirements() {
	// e.g. are MoClo types valid?
}

// Conditions to run on startup
func _CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapSetup(_ctx context.Context, _input *CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapSteps(_ctx context.Context, _input *CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapInput, _output *CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapOutput) {
	_output.StatusMap = make(map[string]string)
	_output.PartswithOverhangsMap = make(map[string][]wtype.DNASequence) // parts to order
	_output.Assemblies = make(map[string][]wtype.DNASequence)
	_output.PassMap = make(map[string]bool)
	_output.SeqsMap = make(map[string]wtype.DNASequence) // desired sequence to end up with after assembly
	_output.EndreportMap = make(map[string]string)
	_output.PositionReportMap = make(map[string][]string)
	_output.StatusMap = make(map[string]string)
	_output.PrimerMap = make(map[string]oligos.Primer)
	SequencingPrimers := make([][]wtype.DNASequence, 0)

	var counter int = 1

	for j := range _input.Vectors {
		for k := range _input.PROs {
			for l := range _input.RBSs {
				for m := range _input.CDSs {
					for n := range _input.TERs {
						key := _input.ProjectName + _input.Vectors[j].Nm + "_" + _input.PROs[k].Nm + "_" + _input.RBSs[l].Nm + "_" + _input.CDSs[m].Nm

						assembly := AssemblyStandard_siteremove_orfcheck_wtypeRunSteps(_ctx, &AssemblyStandard_siteremove_orfcheck_wtypeInput{Constructname: key,
							Seqsinorder:                   []wtype.DNASequence{_input.PROs[k], _input.RBSs[l], _input.CDSs[m], _input.TERs[n]},
							AssemblyStandard:              _input.Standard,
							Level:                         _input.StandardLevel, // of assembly standard
							Vector:                        _input.Vectors[j],
							PartMoClotypesinorder:         []string{"Pro", "5U + NT1", "CDS1", "3U + Ter"},
							OtherEnzymeSitesToRemove:      _input.SitesToRemove,
							ORFstoConfirm:                 []string{}, // enter each as amino acid sequence
							RemoveproblemRestrictionSites: true,
							EndsAlreadyadded:              false,
							ExporttoFastaFile:             true,
							BlastSeqswithNoName:           _input.BlastSearchSeqs},
						)
						key = key                                                             //+ Vectors[j]
						_output.PartswithOverhangsMap[key] = assembly.Data.PartswithOverhangs // parts to order
						_output.Assemblies[key] = assembly.Data.PartsAndVector                // parts + vector to be fed into assembly element
						_output.PassMap[key] = assembly.Data.Simulationpass
						_output.EndreportMap[key] = assembly.Data.Endreport
						_output.PositionReportMap[key] = assembly.Data.PositionReport
						_output.SeqsMap[key] = assembly.Data.NewDNASequence
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
						SequencingPrimers = append(SequencingPrimers, []wtype.DNASequence{primer.Data.FWDPrimer.DNASequence, primer.Data.REVPrimer.DNASequence})
						counter++
					}
				}
			}
		}
	}
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapAnalysis(_ctx context.Context, _input *CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapInput, _output *CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapValidation(_ctx context.Context, _input *CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapInput, _output *CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapOutput) {
}
func _CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapRun(_ctx context.Context, input *CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapInput) *CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapOutput {
	output := &CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapOutput{}
	_CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapSetup(_ctx, input)
	_CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapSteps(_ctx, input, output)
	_CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapAnalysis(_ctx, input, output)
	_CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapValidation(_ctx, input, output)
	return output
}

func CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapRunSteps(_ctx context.Context, input *CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapInput) *CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapSOutput {
	soutput := &CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapSOutput{}
	output := _CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapNew() interface{} {
	return &CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapInput{},
			Out: &CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapElement struct {
	inject.CheckedRunner
}

type CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapInput struct {
	BlastSearchSeqs bool
	CDSs            []wtype.DNASequence
	PROs            []wtype.DNASequence
	ProjectName     string
	RBSs            []wtype.DNASequence
	SitesToRemove   []string
	Standard        string
	StandardLevel   string
	TERs            []wtype.DNASequence
	Vectors         []wtype.DNASequence
}

type CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapOutput struct {
	Assemblies            map[string][]wtype.DNASequence
	EndreportMap          map[string]string
	PartswithOverhangsMap map[string][]wtype.DNASequence
	PassMap               map[string]bool
	PositionReportMap     map[string][]string
	PrimerMap             map[string]oligos.Primer
	SeqsMap               map[string]wtype.DNASequence
	StatusMap             map[string]string
}

type CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapSOutput struct {
	Data struct {
		Assemblies            map[string][]wtype.DNASequence
		EndreportMap          map[string]string
		PartswithOverhangsMap map[string][]wtype.DNASequence
		PassMap               map[string]bool
		PositionReportMap     map[string][]string
		PrimerMap             map[string]oligos.Primer
		SeqsMap               map[string]wtype.DNASequence
		StatusMap             map[string]string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_map",
		Constructor: CombinatorialLibraryDesign_PRO_RBS_CDS_wtype_mapNew,
		Desc: component.ComponentDesc{
			Desc: "This protocol is intended to design a combinatorial library of all combinations of a list of Vectors, Promoters,\nRBSs, CDSs and Terminators according to an assembly standard.\nA list of sequencing primers to order will also be returned.\n",
			Path: "src/github.com/antha-lang/elements/an/Data/DNA/TypeIISAssembly_design/CombinatorialLibraryDesign4part_wtype.an",
			Params: []component.ParamDesc{
				{Name: "BlastSearchSeqs", Desc: "", Kind: "Parameters"},
				{Name: "CDSs", Desc: "", Kind: "Parameters"},
				{Name: "PROs", Desc: "", Kind: "Parameters"},
				{Name: "ProjectName", Desc: "", Kind: "Parameters"},
				{Name: "RBSs", Desc: "", Kind: "Parameters"},
				{Name: "SitesToRemove", Desc: "", Kind: "Parameters"},
				{Name: "Standard", Desc: "Custom\n", Kind: "Parameters"},
				{Name: "StandardLevel", Desc: "level of assembly standard options are: Level0, Level1\n", Kind: "Parameters"},
				{Name: "TERs", Desc: "", Kind: "Parameters"},
				{Name: "Vectors", Desc: "", Kind: "Parameters"},
				{Name: "Assemblies", Desc: "parts + vector map ready for feeding into downstream AutoAssembly element\n", Kind: "Data"},
				{Name: "EndreportMap", Desc: "", Kind: "Data"},
				{Name: "PartswithOverhangsMap", Desc: "parts to order\n", Kind: "Data"},
				{Name: "PassMap", Desc: "", Kind: "Data"},
				{Name: "PositionReportMap", Desc: "", Kind: "Data"},
				{Name: "PrimerMap", Desc: "", Kind: "Data"},
				{Name: "SeqsMap", Desc: "desired sequence to end up with after assembly\n", Kind: "Data"},
				{Name: "StatusMap", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
