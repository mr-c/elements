// Example of how to look up molecule properties from pubchem.
// A map of molar concentrations required to make up 1 Mol/l of each compound is also returned.
package lib

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Pubchem"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Name of compound or array of multiple compounds

// Set Concentrations for compounds or set default. If no concentration is set the molar concentration will be used

// molecule type is returned consisting of name, formula, molecular weight and chemical ID (CID)

// Names of compounds

// A map of molar concentrations required to make up 1 Mol/l of each compound is returned.
// The compound name is used as the key. If any duplicate compound names are returned and errow will be generated.

// return any warnings if the default concentrations have been used for any molecule

func _LookUpMoleculesRequirements() {
}
func _LookUpMoleculesSetup(_ctx context.Context, _input *LookUpMoleculesInput) {
}
func _LookUpMoleculesSteps(_ctx context.Context, _input *LookUpMoleculesInput, _output *LookUpMoleculesOutput) {

	_output.MolarConcentrations = make(map[string]wunit.Concentration)

	for _, molecule := range _input.Compoundlist {
		// method of making molecule from name
		moleculeProperties, err := pubchem.MakeMolecule(molecule)

		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}

		_output.CompoundProperties = append(_output.CompoundProperties, moleculeProperties)
		_output.CompoundNames = append(_output.CompoundNames, moleculeProperties.Name)

		// now set conc
		var moleculeConc wunit.Concentration
		var found bool

		// check if molecule conc is set in OverRideDefaultConcentration list
		if moleculeConc, found = _input.OverRideDefaultConcentration[molecule]; found {
			_output.Warnings = append(_output.Warnings, fmt.Sprintf(`Override Concentration for molecule %s specified so using specified concentration %s instead of default`, molecule, moleculeConc.ToString()))
			// else check if a default is specified
		} else if moleculeConc, found = _input.OverRideDefaultConcentration["default"]; found {
			_output.Warnings = append(_output.Warnings, fmt.Sprintf(`No concentration for molecule %s specified but "default" was specified so using default %s`, molecule, moleculeConc.ToString()))
			// else use 1 Mol/l conc as default
		} else {
			moleculeConc = wunit.NewConcentration(1, "M/l")
			_output.Warnings = append(_output.Warnings, fmt.Sprintf("No concentration for molecule %s specified so using default Molar concentration %s", molecule, moleculeConc.ToString()))

		}

		if _, found := _output.MolarConcentrations[moleculeProperties.Name]; found {
			execute.Errorf(_ctx, "Duplicate compound found for %s", moleculeProperties.Name)
		} else {
			_output.MolarConcentrations[moleculeProperties.Name] = moleculeProperties.GramPerL(moleculeConc)
		}
	}

}
func _LookUpMoleculesAnalysis(_ctx context.Context, _input *LookUpMoleculesInput, _output *LookUpMoleculesOutput) {

}

func _LookUpMoleculesValidation(_ctx context.Context, _input *LookUpMoleculesInput, _output *LookUpMoleculesOutput) {

}
func _LookUpMoleculesRun(_ctx context.Context, input *LookUpMoleculesInput) *LookUpMoleculesOutput {
	output := &LookUpMoleculesOutput{}
	_LookUpMoleculesSetup(_ctx, input)
	_LookUpMoleculesSteps(_ctx, input, output)
	_LookUpMoleculesAnalysis(_ctx, input, output)
	_LookUpMoleculesValidation(_ctx, input, output)
	return output
}

func LookUpMoleculesRunSteps(_ctx context.Context, input *LookUpMoleculesInput) *LookUpMoleculesSOutput {
	soutput := &LookUpMoleculesSOutput{}
	output := _LookUpMoleculesRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func LookUpMoleculesNew() interface{} {
	return &LookUpMoleculesElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &LookUpMoleculesInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _LookUpMoleculesRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &LookUpMoleculesInput{},
			Out: &LookUpMoleculesOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type LookUpMoleculesElement struct {
	inject.CheckedRunner
}

type LookUpMoleculesInput struct {
	Compoundlist                 []string
	OverRideDefaultConcentration map[string]wunit.Concentration
}

type LookUpMoleculesOutput struct {
	CompoundNames       []string
	CompoundProperties  []pubchem.Molecule
	MolarConcentrations map[string]wunit.Concentration
	Warnings            []string
}

type LookUpMoleculesSOutput struct {
	Data struct {
		CompoundNames       []string
		CompoundProperties  []pubchem.Molecule
		MolarConcentrations map[string]wunit.Concentration
		Warnings            []string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "LookUpMolecules",
		Constructor: LookUpMoleculesNew,
		Desc: component.ComponentDesc{
			Desc: "Example of how to look up molecule properties from pubchem.\nA map of molar concentrations required to make up 1 Mol/l of each compound is also returned.\n",
			Path: "src/github.com/antha-lang/elements/an/Data/LookUpMolecule/LookUpMolecule.an",
			Params: []component.ParamDesc{
				{Name: "Compoundlist", Desc: "Name of compound or array of multiple compounds\n", Kind: "Parameters"},
				{Name: "OverRideDefaultConcentration", Desc: "Set Concentrations for compounds or set default. If no concentration is set the molar concentration will be used\n", Kind: "Parameters"},
				{Name: "CompoundNames", Desc: "Names of compounds\n", Kind: "Data"},
				{Name: "CompoundProperties", Desc: "molecule type is returned consisting of name, formula, molecular weight and chemical ID (CID)\n", Kind: "Data"},
				{Name: "MolarConcentrations", Desc: "A map of molar concentrations required to make up 1 Mol/l of each compound is returned.\nThe compound name is used as the key. If any duplicate compound names are returned and errow will be generated.\n", Kind: "Data"},
				{Name: "Warnings", Desc: "return any warnings if the default concentrations have been used for any molecule\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
