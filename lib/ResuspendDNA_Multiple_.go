// Protocol for resuspending freeze dried DNA with a diluent
// The inputs of this are designed to be wired in from the ParseDNAPlate element which Parses a DNA plate file from Thermo.
// The protocol can also be used as a stand-alone element by specifying all the necessary paramters manually.
// Defaults can also be specified for most parameters to make the process easier.
package lib

import

// we need to import the wtype package to use the LHComponent type
// the mixer package is required to use the Sample function
(
	"context"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol (data)

// Target concentration to resuspend the DNA to.

// A list of part names

// A map of the mass of each part.
// A default can be specified which will apply to any part which does not have a value explicitely set.

// A map of the molecular weight of each part.
// A default can be specified which will apply to any part which does not have a value explicitely set.

// If selected then for any part that no well location is specified for in PartLocationsMap, the system will autoallocate a position for it on the input plate.

// A map of part name to well location (i.e. A1, B12). If a part is mssing an error will occur unless AutoAllocateMissingParts is selected.

// A map of which plate each part is located on.
// A default can be specified which will apply to any part which does not have a plate name explicitely set.

// If no Policy is specified the default policy will be MegaMix which mixes the sample 10 times.

// This can be wired directly into AutoAssembly_Conc_Multi to use some or all of the parts directly after resuspension.

func _ResuspendDNA_MultipleRequirements() {
}

func _ResuspendDNA_MultipleSetup(_ctx context.Context, _input *ResuspendDNA_MultipleInput) {
}

func _ResuspendDNA_MultipleSteps(_ctx context.Context, _input *ResuspendDNA_MultipleInput, _output *ResuspendDNA_MultipleOutput) {

	if _input.Projectname == "" {
		_input.Projectname = "ResuspendDNA"
	}

	// set up some empty slices to fill as we iterate through the reactions
	Reactions := make([]*wtype.LHComponent, 0)
	volumes := make([]wunit.Volume, 0)
	welllocations := make([]string, 0)

	_output.ResuspendedDNAMap = make(map[string]*wtype.LHComponent)
	_output.PartConcentrations = make(map[string]wunit.Concentration)

	for _, part := range _input.Parts {

		mass, found := _input.PartMassMap[part]

		if !found {
			mass, found = _input.PartMassMap["default"]

			if !found {
				execute.Errorf(_ctx, "Part %s not found in PartMassMap and no default set.", part)
			}
		}

		mw, found := _input.PartMolecularWeightMap[part]

		if !found {
			mw, found = _input.PartMolecularWeightMap["default"]

			if !found {
				execute.Errorf(_ctx, "Part %s not found in PartMolecularWeightMap and no default set.", part)
			}
		}

		well, found := _input.PartLocationsMap[part]

		if !found {

			if _input.AutoAllocateMissingParts {
				well = ""
			} else {
				execute.Errorf(_ctx, "Part %s not found in PartLocationsMap", part)
			}
		}

		plate, found := _input.PartPlateMap[part]

		if !found {
			plate, found = _input.PartPlateMap["default"]

			if !found {
				execute.Errorf(_ctx, "Part %s not found in PartPlateMap and no default set.", part)
			}
		}

		result := ResuspendDNARunSteps(_ctx, &ResuspendDNAInput{DNAMass: mass,
			TargetConc:           _input.TargetConc,
			MolecularWeight:      mw,
			Well:                 well,
			PlateName:            plate,
			OverRideLiquidPolicy: _input.OverRideLiquidPolicy,

			Diluent:  _input.Diluent,
			DNAPlate: _input.DNAPlate},
		)

		result.Outputs.ResuspendedDNA.CName = part

		resuspendedDNA := result.Outputs.ResuspendedDNA

		// convert concentration to g/l
		conc := _input.TargetConc.GramPerL(mw)

		resuspendedDNA.SetConcentration(conc)

		// add to output maps
		_output.ResuspendedDNAMap[part] = resuspendedDNA

		// add to output array
		_output.ResuspendedDNAArray = append(_output.ResuspendedDNAArray, resuspendedDNA)

		_output.PartConcentrations[part] = conc

		// add to slices to export as csv later
		Reactions = append(Reactions, resuspendedDNA)
		volumes = append(volumes, resuspendedDNA.Volume())
		welllocations = append(welllocations, well)
	}

	// once all values of loop have been completed, export the plate contents as a csv file, Not visible in UI at present! refactor exportCSV func.
	err := wtype.ExportPlateCSV(_input.Projectname+".csv", _input.DNAPlate, _input.Projectname+"outputPlate", welllocations, Reactions, volumes)
	_output.Errors = append(_output.Errors, err)

}

func _ResuspendDNA_MultipleAnalysis(_ctx context.Context, _input *ResuspendDNA_MultipleInput, _output *ResuspendDNA_MultipleOutput) {
}

func _ResuspendDNA_MultipleValidation(_ctx context.Context, _input *ResuspendDNA_MultipleInput, _output *ResuspendDNA_MultipleOutput) {
}
func _ResuspendDNA_MultipleRun(_ctx context.Context, input *ResuspendDNA_MultipleInput) *ResuspendDNA_MultipleOutput {
	output := &ResuspendDNA_MultipleOutput{}
	_ResuspendDNA_MultipleSetup(_ctx, input)
	_ResuspendDNA_MultipleSteps(_ctx, input, output)
	_ResuspendDNA_MultipleAnalysis(_ctx, input, output)
	_ResuspendDNA_MultipleValidation(_ctx, input, output)
	return output
}

func ResuspendDNA_MultipleRunSteps(_ctx context.Context, input *ResuspendDNA_MultipleInput) *ResuspendDNA_MultipleSOutput {
	soutput := &ResuspendDNA_MultipleSOutput{}
	output := _ResuspendDNA_MultipleRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ResuspendDNA_MultipleNew() interface{} {
	return &ResuspendDNA_MultipleElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ResuspendDNA_MultipleInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ResuspendDNA_MultipleRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ResuspendDNA_MultipleInput{},
			Out: &ResuspendDNA_MultipleOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type ResuspendDNA_MultipleElement struct {
	inject.CheckedRunner
}

type ResuspendDNA_MultipleInput struct {
	AutoAllocateMissingParts bool
	DNAPlate                 *wtype.LHPlate
	Diluent                  *wtype.LHComponent
	OverRideLiquidPolicy     wtype.PolicyName
	PartLocationsMap         map[string]string
	PartMassMap              map[string]wunit.Mass
	PartMolecularWeightMap   map[string]float64
	PartPlateMap             map[string]string
	Parts                    []string
	Projectname              string
	TargetConc               wunit.Concentration
}

type ResuspendDNA_MultipleOutput struct {
	Errors              []error
	PartConcentrations  map[string]wunit.Concentration
	PlateContents       wtype.File
	ResuspendedDNAArray []*wtype.LHComponent
	ResuspendedDNAMap   map[string]*wtype.LHComponent
}

type ResuspendDNA_MultipleSOutput struct {
	Data struct {
		Errors             []error
		PartConcentrations map[string]wunit.Concentration
		PlateContents      wtype.File
	}
	Outputs struct {
		ResuspendedDNAArray []*wtype.LHComponent
		ResuspendedDNAMap   map[string]*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "ResuspendDNA_Multiple",
		Constructor: ResuspendDNA_MultipleNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol for resuspending freeze dried DNA with a diluent\nThe inputs of this are designed to be wired in from the ParseDNAPlate element which Parses a DNA plate file from Thermo.\nThe protocol can also be used as a stand-alone element by specifying all the necessary paramters manually.\nDefaults can also be specified for most parameters to make the process easier.\n",
			Path: "src/github.com/antha-lang/elements/an/Liquid_handling/ResuspendDNA/ResuspendDNAFromPlate.an",
			Params: []component.ParamDesc{
				{Name: "AutoAllocateMissingParts", Desc: "If selected then for any part that no well location is specified for in PartLocationsMap, the system will autoallocate a position for it on the input plate.\n", Kind: "Parameters"},
				{Name: "DNAPlate", Desc: "", Kind: "Inputs"},
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "OverRideLiquidPolicy", Desc: "If no Policy is specified the default policy will be MegaMix which mixes the sample 10 times.\n", Kind: "Parameters"},
				{Name: "PartLocationsMap", Desc: "A map of part name to well location (i.e. A1, B12). If a part is mssing an error will occur unless AutoAllocateMissingParts is selected.\n", Kind: "Parameters"},
				{Name: "PartMassMap", Desc: "A map of the mass of each part.\nA default can be specified which will apply to any part which does not have a value explicitely set.\n", Kind: "Parameters"},
				{Name: "PartMolecularWeightMap", Desc: "A map of the molecular weight of each part.\nA default can be specified which will apply to any part which does not have a value explicitely set.\n", Kind: "Parameters"},
				{Name: "PartPlateMap", Desc: "A map of which plate each part is located on.\nA default can be specified which will apply to any part which does not have a plate name explicitely set.\n", Kind: "Parameters"},
				{Name: "Parts", Desc: "A list of part names\n", Kind: "Parameters"},
				{Name: "Projectname", Desc: "", Kind: "Parameters"},
				{Name: "TargetConc", Desc: "Target concentration to resuspend the DNA to.\n", Kind: "Parameters"},
				{Name: "Errors", Desc: "", Kind: "Data"},
				{Name: "PartConcentrations", Desc: "", Kind: "Data"},
				{Name: "PlateContents", Desc: "", Kind: "Data"},
				{Name: "ResuspendedDNAArray", Desc: "This can be wired directly into AutoAssembly_Conc_Multi to use some or all of the parts directly after resuspension.\n", Kind: "Outputs"},
				{Name: "ResuspendedDNAMap", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
