// This protocol is intended to design assembly parts using a specified enzyme.
// overhangs are added to complement the adjacent parts and leave no scar.
// parts can be entered as genbank (.gb) files, sequences or biobrick IDs
// If assembly simulation fails after overhangs are added. In order to help the user
// diagnose the reason, a report of the part overhangs
// is returned to the user along with a list of cut sites in each part.
package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/export"
	"github.com/antha-lang/antha/antha/anthalib/wtype"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"

	"context"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"path/filepath"
	"strconv"
	"strings"
)

// Input parameters for this protocol (data)

//MoClo
// of assembly standard

// labels e.g. pro = promoter

// enter each as amino acid sequence

// Option to add Level 1 adaptor sites to the Promoters and terminators to support hierarchical assembly
// If Custom design the valid options currently supported are: "Device1","Device2", "Device3".
// If left empty no adaptor sequence is added.

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// and level 1 ends added if MakeLevel1Device is selected
// parts to order
// parts to order + vector
// desired sequence to end up with after assembly
// sequence of the assembled insert. Useful for sequencing validation and Primer design

// Input Requirement specification
func _AssemblyStandard_TypeIIsDesignRequirements() {
	// e.g. are MoClo types valid?
}

// Conditions to run on startup
func _AssemblyStandard_TypeIIsDesignSetup(_ctx context.Context, _input *AssemblyStandard_TypeIIsDesignInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _AssemblyStandard_TypeIIsDesignSteps(_ctx context.Context, _input *AssemblyStandard_TypeIIsDesignInput, _output *AssemblyStandard_TypeIIsDesignOutput) {

	// set warnings reported back to user to none initially
	warnings := make([]string, 0)

	// declare some temporary variables to be used later
	var warning string
	var err error

	// get these directly from Parameters
	var partsinorder = _input.Seqsinorder
	var vectordata = _input.Vector

	_output.Status = "all parts available"

	// export this as data output
	_output.OriginalParts = partsinorder
	// check parts for restriction sites first and remove if the user has chosen to
	enz, found := enzymes.Enzymelookup[_input.AssemblyStandard][_input.Level]

	if !found {
		execute.Errorf(_ctx, "AssemblyStandard ", _input.AssemblyStandard, " level ", _input.Level, " not found")
	}

	// get properties of other enzyme sites to remove
	removetheseenzymes := make([]wtype.RestrictionEnzyme, 0)
	removetheseenzymes = append(removetheseenzymes, enz.RestrictionEnzyme)

	for _, enzyme := range _input.OtherEnzymeSitesToRemove {

		enzyTypeII, err := lookup.EnzymeLookup(enzyme)

		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}

		removetheseenzymes = append(removetheseenzymes, enzyTypeII)
	}

	// check number of sites per part and return if > 0!
	var report []string
	var siteFound bool
	for _, part := range partsinorder {

		info := enzymes.Restrictionsitefinder(part, removetheseenzymes)

		for i := range info {
			sitepositions := enzymes.SitepositionString(info[i])

			if len(sitepositions) > 0 {
				siteFound = true
			}
			sitepositions = fmt.Sprint(part.Nm+" "+info[i].Enzyme.Name+" positions:", sitepositions)
			report = append(report, sitepositions)
		}
	}
	if siteFound {

		errormessage := fmt.Sprintf("Found problem restriction sites in 1 or more parts: %s", report)

		if !_input.RemoveproblemRestrictionSites && !_input.EndsAlreadyadded {
			execute.Errorf(_ctx, errormessage)
		} else {
			warnings = append(warnings, errormessage)
		}
	}

	warning = text.Print("RemoveproblemRestrictionSites =", _input.RemoveproblemRestrictionSites)
	warnings = append(warnings, warning)
	if _input.RemoveproblemRestrictionSites && !_input.EndsAlreadyadded {
		newparts := make([]wtype.DNASequence, 0)
		warning = "Starting process or removing restrictionsite"
		warnings = append(warnings, warning)

		for _, part := range partsinorder {
			if part.Seq != "" {
				info := enzymes.Restrictionsitefinder(part, removetheseenzymes)

				for _, anysites := range info {
					if anysites.Sitefound {
						warning = "problem " + anysites.Enzyme.Name + " site found in " + part.Nm
						warnings = append(warnings, warning)
						orf, orftrue := sequences.FindBiggestORF(part.Seq)
						warning = fmt.Sprintln(anysites.Enzyme.Name+" site found in orf ", part.Nm, " ", orftrue, " site positions ", anysites.Positions("ALL"), "orf between", orf.StartPosition, " and ", orf.EndPosition /*orf.DNASeq[orf.StartPosition:orf.EndPosition]*/)
						warnings = append(warnings, warning)
						if orftrue && len(orf.ProtSeq) > 20 {
							allsitestoavoid := make([]string, 0)
							allsitestoavoid = append(allsitestoavoid, anysites.Recognitionsequence, sequences.RevComp(anysites.Recognitionsequence))
							orfcoordinates := sequences.MakeStartendPair(orf.StartPosition, orf.EndPosition)
							for _, position := range anysites.Positions("ALL") {
								if orf.StartPosition < position && position < orf.EndPosition {
									originalcodon := ""
									codonoption := ""
									originalPart := part.Dup()
									part, originalcodon, codonoption, err = sequences.ReplaceCodoninORF(originalPart, orfcoordinates, position, allsitestoavoid)
									warning = fmt.Sprintln("sites to avoid: ", allsitestoavoid[0], allsitestoavoid[1])
									warnings = append(warnings, warning)
									warnings = append(warnings, "For Part Sequence: "+originalPart.Seq+" position: "+strconv.Itoa(position)+" original codon to replace: "+originalcodon+" replaced with replacementcodon: "+codonoption)
									if err != nil {
										warning := fmt.Sprint("removal of "+anysites.Enzyme.Name+" site from orf "+orf.DNASeq, "in part "+part.Nm+" failed! improve your algorithm! "+err.Error())
										warnings = append(warnings, warning)
										execute.Errorf(_ctx, warning)
									}
								} else if !_input.OnlyRemovesitesinORFs {
									allsitestoavoid := make([]string, 0)
									part, err = sequences.RemoveSite(part, anysites.Enzyme, allsitestoavoid)
									if err != nil {
										warning = text.Print("Failed to remove "+anysites.Enzyme.Name+" site from position "+strconv.Itoa(position)+". Position found to be outside of orf: "+orf.DNASeq, ". Error: "+err.Error())
										warnings = append(warnings, warning)
										execute.Errorf(_ctx, warning)
									}
								}
							}
						} else if !_input.OnlyRemovesitesinORFs {
							allsitestoavoid := make([]string, 0)
							temppart, err := sequences.RemoveSite(part, anysites.Enzyme, allsitestoavoid)

							if err != nil {
								warning := fmt.Sprintf("Faiure to automatically remove %s site at positions %+v in part %s: %s! improve algorithm  or remove manually until you do!: %s", anysites.Enzyme.Name, anysites.Positions("ALL"), part.Nm, part.Seq, err.Error())
								warnings = append(warnings, warning)
								execute.Errorf(_ctx, warning)

							}
							warning = fmt.Sprintln("modified "+temppart.Nm+"new seq: ", temppart.Seq, "original seq: ", part.Seq)
							warnings = append(warnings, warning)
							part = temppart

						}
					}

				}
				newparts = append(newparts, part)

			} else {
				newparts = append(newparts, part)
			}

		}
		partsinorder = newparts
	}
	// export the parts list with sites removed
	_output.PartsWithSitesRemoved = partsinorder
	// make vector into an antha type DNASequence

	if _input.MakeLevel1Device != "" {

		standard, found := enzymes.EndlinksString[_input.AssemblyStandard]

		if !found {
			execute.Errorf(_ctx, "No assembly standard %s found", _input.AssemblyStandard)
		}

		level1, found := standard["Level1"]

		if !found {
			execute.Errorf(_ctx, "No Level1 found for standard %s", _input.AssemblyStandard)
		}

		overhangs, found := level1[_input.MakeLevel1Device]

		if !found {
			execute.Errorf(_ctx, "No overhangs found for %s in standard %s", _input.MakeLevel1Device, _input.AssemblyStandard)
		}

		if len(overhangs) != 2 {
			execute.Errorf(_ctx, "found %d overhangs for %s in standard %s, expecting %d", len(overhangs), _input.MakeLevel1Device, _input.AssemblyStandard, 2)

		}

		if overhangs[0] == "" {
			execute.Errorf(_ctx, "blunt 5' overhang found for %s in standard %s, expecting %d", _input.MakeLevel1Device, _input.AssemblyStandard, 2)
		}

		_output.PartsWithSitesRemoved[0], err = enzymes.AddL1UAdaptor(_output.PartsWithSitesRemoved[0], _input.AssemblyStandard, "Level1", _input.MakeLevel1Device, _input.ReverseLevel1Orientation)
		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}

		_output.PartsWithSitesRemoved[len(_output.PartsWithSitesRemoved)-1], err = enzymes.AddL1DAdaptor(_output.PartsWithSitesRemoved[len(_output.PartsWithSitesRemoved)-1], _input.AssemblyStandard, "Level1", _input.MakeLevel1Device, _input.ReverseLevel1Orientation)
		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}

	}
	//lookup restriction enzyme
	restrictionenzyme, err := lookup.TypeIIsLookup(enz.Name)
	if err != nil {
		warnings = append(warnings, text.Print("Error", err.Error()))
	}

	//  Add overhangs for scarfree assembly based on part seqeunces only, i.e. no Assembly standard
	fmt.Println("warnings:", warnings)

	if _input.EndsAlreadyadded {
		_output.PartswithOverhangs = partsinorder
	} else {
		//  Add standard overhangs using chosen assembly standard
		_output.PartswithOverhangs, err = enzymes.MakeStandardTypeIIsassemblyParts(partsinorder, _input.AssemblyStandard, _input.Level, _input.PartMoClotypesinorder)

		if err != nil {
			warnings = append(warnings, text.Print("Error", err.Error()))
			execute.Errorf(_ctx, err.Error())
		}

	}

	// Check that assembly is feasible with designed parts by simulating assembly of the sequences with the chosen enzyme
	assembly := enzymes.Assemblyparameters{_input.Constructname, restrictionenzyme.Name, vectordata, _output.PartswithOverhangs}
	status, numberofassemblies, _, newDNASequence, err := enzymes.Assemblysimulator(assembly)

	if err != nil {
		execute.Errorf(_ctx, "Error simulating assembly of %s: %s ", _input.Constructname, err.Error())
	}

	_output.Insert, err = assembly.Insert()

	if err != nil {
		execute.Errorf(_ctx, "Error calculating insert from assembly: %s. Sites at positions: %s ", err.Error(), siteReport(partsinorder, removetheseenzymes))
	}

	// export parts + vector as one array
	for _, part := range _output.PartswithOverhangs {
		_output.PartsAndVector = append(_output.PartsAndVector, part)
	}

	// now add vector
	_output.PartsAndVector = append(_output.PartsAndVector, vectordata)

	endreport := "Endreport only run in the event of assembly simulation failure"
	//sites := "Restriction mapper only run in the event of assembly simulation failure"
	newDNASequence.Nm = _input.Constructname
	_output.NewDNASequence = newDNASequence
	if err == nil && numberofassemblies == 1 {

		_output.Simulationpass = true
	} // else {

	warnings = append(warnings, status)
	// perform mock digest to test fragement overhangs (fragments are hidden by using _, )
	_, stickyends5, stickyends3 := enzymes.TypeIIsdigest(vectordata, restrictionenzyme)

	allends := make([]string, 0)
	ends := ""

	ends = text.Print(vectordata.Nm+" 5 Prime end: ", stickyends5)
	allends = append(allends, ends)
	ends = text.Print(vectordata.Nm+" 3 Prime end: ", stickyends3)
	allends = append(allends, ends)

	for _, part := range _output.PartswithOverhangs {
		_, stickyends5, stickyends3 := enzymes.TypeIIsdigest(part, restrictionenzyme)
		ends = text.Print(part.Nm+" 5 Prime end: ", stickyends5)
		allends = append(allends, ends)
		ends = text.Print(part.Nm+" 3 Prime end: ", stickyends3)
		allends = append(allends, ends)
	}
	endreport = strings.Join(allends, " ")
	if !_output.Simulationpass {
		warnings = append(warnings, endreport)
	}
	//	}

	// check number of sites per part !

	sites := make([]int, 0)
	multiple := make([]string, 0)
	for _, part := range _output.PartswithOverhangs {

		info := enzymes.Restrictionsitefinder(part, removetheseenzymes)

		for i := range info {
			sitepositions := enzymes.SitepositionString(info[i])

			sites = append(sites, info[i].Numberofsites)
			sitepositions = text.Print(part.Nm+" "+info[i].Enzyme.Name+" positions:", sitepositions)
			multiple = append(multiple, sitepositions)
		}
	}

	for _, orf := range _input.ORFstoConfirm {
		if sequences.LookforSpecificORF(_output.NewDNASequence.Seq, orf) == false {
			warning = text.Print("orf not present: ", orf)
			warnings = append(warnings, warning)
			_output.ORFmissing = true
		}
	}

	if len(warnings) == 0 {
		warnings = append(warnings, "none")
	}
	_output.Warnings = fmt.Errorf(strings.Join(warnings, ";"))
	_output.Endreport = endreport
	_output.PositionReport = multiple

	partsummary := make([]string, 0)
	for _, part := range _output.PartswithOverhangs {
		partsummary = append(partsummary, text.Print(part.Nm, part.Seq))
	}
	partsummary = append(partsummary, text.Print("Vector:"+vectordata.Nm, vectordata.Seq))
	partstoorder := text.Print("PartswithOverhangs: ", partsummary)

	// Print status
	if _output.Status != "all parts available" {
		_output.Status = fmt.Sprintln(_output.Status)
	} else {
		_output.Status = fmt.Sprintln(
			text.Print("simulator status: ", status),
			text.Print("Endreport after digestion: ", endreport),
			text.Print("Sites per part for "+enz.Name, sites),
			text.Print("Positions: ", multiple),
			text.Print("Warnings:", _output.Warnings.Error()),
			text.Print("Simulationpass=", _output.Simulationpass),
			text.Print("NewDNASequence: ", _output.NewDNASequence),
			text.Print("Any Orfs to confirm missing from new DNA sequence:", _output.ORFmissing),
			partstoorder,
		)
		// export data to file
		//anthapath.ExporttoFile("Report"+"_"+Constructname+".txt",[]byte(Status))
		//anthapath.ExportTextFile("Report"+"_"+Constructname+".txt",Status)
		fmt.Println(_output.Status)
	}

	// export sequence to fasta
	if _input.ExporttoFastaFile && _output.Simulationpass {
		exportedsequences := make([]wtype.DNASequence, 0)
		// add dna sequence produced
		exportedsequences = append(exportedsequences, _output.NewDNASequence)
		if len(exportedsequences) == 0 {
			execute.Errorf(_ctx, "No Sequences!")
		}
		// export to file
		_output.AssembledSequenceFile, _, err = export.FastaSerial(export.LOCAL, filepath.Join(_input.Constructname, "AssemblyProduct"), exportedsequences)
		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}
		// reset
		exportedsequences = make([]wtype.DNASequence, 0)
		// add all parts with overhangs
		for _, part := range _output.PartswithOverhangs {
			exportedsequences = append(exportedsequences, part)
		}
		_output.PartsToOrderFile, _, err = export.FastaSerial(export.LOCAL, filepath.Join(_input.Constructname, "Parts"), exportedsequences)
		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AssemblyStandard_TypeIIsDesignAnalysis(_ctx context.Context, _input *AssemblyStandard_TypeIIsDesignInput, _output *AssemblyStandard_TypeIIsDesignOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AssemblyStandard_TypeIIsDesignValidation(_ctx context.Context, _input *AssemblyStandard_TypeIIsDesignInput, _output *AssemblyStandard_TypeIIsDesignOutput) {
}
func _AssemblyStandard_TypeIIsDesignRun(_ctx context.Context, input *AssemblyStandard_TypeIIsDesignInput) *AssemblyStandard_TypeIIsDesignOutput {
	output := &AssemblyStandard_TypeIIsDesignOutput{}
	_AssemblyStandard_TypeIIsDesignSetup(_ctx, input)
	_AssemblyStandard_TypeIIsDesignSteps(_ctx, input, output)
	_AssemblyStandard_TypeIIsDesignAnalysis(_ctx, input, output)
	_AssemblyStandard_TypeIIsDesignValidation(_ctx, input, output)
	return output
}

func AssemblyStandard_TypeIIsDesignRunSteps(_ctx context.Context, input *AssemblyStandard_TypeIIsDesignInput) *AssemblyStandard_TypeIIsDesignSOutput {
	soutput := &AssemblyStandard_TypeIIsDesignSOutput{}
	output := _AssemblyStandard_TypeIIsDesignRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AssemblyStandard_TypeIIsDesignNew() interface{} {
	return &AssemblyStandard_TypeIIsDesignElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AssemblyStandard_TypeIIsDesignInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AssemblyStandard_TypeIIsDesignRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AssemblyStandard_TypeIIsDesignInput{},
			Out: &AssemblyStandard_TypeIIsDesignOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type AssemblyStandard_TypeIIsDesignElement struct {
	inject.CheckedRunner
}

type AssemblyStandard_TypeIIsDesignInput struct {
	AssemblyStandard              string
	BlastSeqswithNoName           bool
	Constructname                 string
	EndsAlreadyadded              bool
	ExporttoFastaFile             bool
	Level                         string
	MakeLevel1Device              string
	ORFstoConfirm                 []string
	OnlyRemovesitesinORFs         bool
	OtherEnzymeSitesToRemove      []string
	PartMoClotypesinorder         []string
	RemoveproblemRestrictionSites bool
	ReverseLevel1Orientation      bool
	Seqsinorder                   []wtype.DNASequence
	Vector                        wtype.DNASequence
}

type AssemblyStandard_TypeIIsDesignOutput struct {
	AssembledSequenceFile wtype.File
	Endreport             string
	Insert                wtype.DNASequence
	NewDNASequence        wtype.DNASequence
	ORFmissing            bool
	OriginalParts         []wtype.DNASequence
	PartsAndVector        []wtype.DNASequence
	PartsToOrderFile      wtype.File
	PartsWithSitesRemoved []wtype.DNASequence
	PartswithOverhangs    []wtype.DNASequence
	PositionReport        []string
	Simulationpass        bool
	Status                string
	Warnings              error
}

type AssemblyStandard_TypeIIsDesignSOutput struct {
	Data struct {
		AssembledSequenceFile wtype.File
		Endreport             string
		Insert                wtype.DNASequence
		NewDNASequence        wtype.DNASequence
		ORFmissing            bool
		OriginalParts         []wtype.DNASequence
		PartsAndVector        []wtype.DNASequence
		PartsToOrderFile      wtype.File
		PartsWithSitesRemoved []wtype.DNASequence
		PartswithOverhangs    []wtype.DNASequence
		PositionReport        []string
		Simulationpass        bool
		Status                string
		Warnings              error
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AssemblyStandard_TypeIIsDesign",
		Constructor: AssemblyStandard_TypeIIsDesignNew,
		Desc: component.ComponentDesc{
			Desc: "This protocol is intended to design assembly parts using a specified enzyme.\noverhangs are added to complement the adjacent parts and leave no scar.\nparts can be entered as genbank (.gb) files, sequences or biobrick IDs\nIf assembly simulation fails after overhangs are added. In order to help the user\ndiagnose the reason, a report of the part overhangs\nis returned to the user along with a list of cut sites in each part.\n",
			Path: "src/github.com/antha-lang/elements/an/Data/DNA/TypeIISAssembly_design/MoClo/AssemblyStandard_TypeIIsDesign.an",
			Params: []component.ParamDesc{
				{Name: "AssemblyStandard", Desc: "MoClo\n", Kind: "Parameters"},
				{Name: "BlastSeqswithNoName", Desc: "", Kind: "Parameters"},
				{Name: "Constructname", Desc: "", Kind: "Parameters"},
				{Name: "EndsAlreadyadded", Desc: "", Kind: "Parameters"},
				{Name: "ExporttoFastaFile", Desc: "", Kind: "Parameters"},
				{Name: "Level", Desc: "of assembly standard\n", Kind: "Parameters"},
				{Name: "MakeLevel1Device", Desc: "Option to add Level 1 adaptor sites to the Promoters and terminators to support hierarchical assembly\nIf Custom design the valid options currently supported are: \"Device1\",\"Device2\", \"Device3\".\nIf left empty no adaptor sequence is added.\n", Kind: "Parameters"},
				{Name: "ORFstoConfirm", Desc: "enter each as amino acid sequence\n", Kind: "Parameters"},
				{Name: "OnlyRemovesitesinORFs", Desc: "", Kind: "Parameters"},
				{Name: "OtherEnzymeSitesToRemove", Desc: "", Kind: "Parameters"},
				{Name: "PartMoClotypesinorder", Desc: "labels e.g. pro = promoter\n", Kind: "Parameters"},
				{Name: "RemoveproblemRestrictionSites", Desc: "", Kind: "Parameters"},
				{Name: "ReverseLevel1Orientation", Desc: "", Kind: "Parameters"},
				{Name: "Seqsinorder", Desc: "", Kind: "Parameters"},
				{Name: "Vector", Desc: "", Kind: "Parameters"},
				{Name: "AssembledSequenceFile", Desc: "", Kind: "Data"},
				{Name: "Endreport", Desc: "", Kind: "Data"},
				{Name: "Insert", Desc: "sequence of the assembled insert. Useful for sequencing validation and Primer design\n", Kind: "Data"},
				{Name: "NewDNASequence", Desc: "desired sequence to end up with after assembly\n", Kind: "Data"},
				{Name: "ORFmissing", Desc: "", Kind: "Data"},
				{Name: "OriginalParts", Desc: "", Kind: "Data"},
				{Name: "PartsAndVector", Desc: "parts to order + vector\n", Kind: "Data"},
				{Name: "PartsToOrderFile", Desc: "", Kind: "Data"},
				{Name: "PartsWithSitesRemoved", Desc: "and level 1 ends added if MakeLevel1Device is selected\n", Kind: "Data"},
				{Name: "PartswithOverhangs", Desc: "parts to order\n", Kind: "Data"},
				{Name: "PositionReport", Desc: "", Kind: "Data"},
				{Name: "Simulationpass", Desc: "", Kind: "Data"},
				{Name: "Status", Desc: "", Kind: "Data"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
