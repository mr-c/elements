// Protocol TransformationEfficiency calculates transformation efficiency based on colony count and transformation parameters.
package lib

import

// Place golang packages to import here
(
	"context"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

type transformationEfficiency struct {
	CFUperugperml float64
	MetaData      metaData
}

type metaData struct {
	CFU            int
	PlateOutVolume wunit.Volume
	DNA            wtype.DNASequence
	DNAMass        wunit.Mass
	Assembly       map[string][]wtype.DNASequence
}

type relativeTransformationEfficiency struct {
	ControlEfficiency      transformationEfficiency
	ExperimentalEfficiency transformationEfficiency
}

func calcTransformationEfficiency(colonyCount int, plateoutVolumeinml wunit.Volume, plateOutdilutionX int, dnamassinug wunit.Mass) (cfuperugperml float64) {
	cfuperugperml = float64(colonyCount) * float64(plateOutdilutionX) / ((volumeToml(plateoutVolumeinml)) * massToug(dnamassinug))
	return
}

func newTransformationEfficiency(seq wtype.DNASequence, colonyCount int, plateoutVolumeinml wunit.Volume, plateOutdilutionX int, dnamassinug wunit.Mass) (efficiency transformationEfficiency) {
	cfuperugperml := calcTransformationEfficiency(colonyCount, plateoutVolumeinml, plateOutdilutionX, dnamassinug)

	efficiency = transformationEfficiency{
		CFUperugperml: cfuperugperml,
		MetaData: metaData{
			CFU:            colonyCount,
			PlateOutVolume: plateoutVolumeinml,
			DNAMass:        dnamassinug,
			DNA:            seq},
	}
	return efficiency
}

func volumeToml(volume wunit.Volume) (volumeinml float64) {
	volumeinml = volume.SIValue() * 1000
	return
}

func massToug(mass wunit.Mass) (massinug float64) {
	massinug = mass.SIValue() * 1000000000
	return
}

// Parameters to this protocol

// Output data of this protocol

// Physical inputs to this protocol

// Physical outputs to this protocol

// Conditions to run on startup
func _TransformationEfficiencySetup(_ctx context.Context, _input *TransformationEfficiencyInput) {

}

// The core process for this protocol. These steps are executed for each input.
func _TransformationEfficiencySteps(_ctx context.Context, _input *TransformationEfficiencyInput, _output *TransformationEfficiencyOutput) {

	var platoutDilutioninX int = 1

	// calculate transformation efficiency of experiemntal sample
	_output.EquivalentTransformationEfficiency = newTransformationEfficiency(_input.DNA, _input.ColonyCount, _input.PlateOutVolume, platoutDilutioninX, _input.DNAMass)

	// calculate transformation efficiency of control sample
	_output.ControlTransformationEfficiency = newTransformationEfficiency(_input.ControlDNA, _input.ControlColonyCount, _input.ControlPlateOutVolume, platoutDilutioninX, _input.ControlDNAMass)

	// Calculate relative transformation efficiency
	_output.RelativeTransformationEfficiency.ControlEfficiency = _output.ControlTransformationEfficiency

	_output.RelativeTransformationEfficiency.ExperimentalEfficiency = _output.EquivalentTransformationEfficiency

}

// Run after controls and a steps block are completed to post process any data
// and provide downstream results
func _TransformationEfficiencyAnalysis(_ctx context.Context, _input *TransformationEfficiencyInput, _output *TransformationEfficiencyOutput) {

}

// A block of tests to perform to validate that the sample was processed
// correctly. Optionally, destructive tests can be performed to validate
// results on a dipstick basis
func _TransformationEfficiencyValidation(_ctx context.Context, _input *TransformationEfficiencyInput, _output *TransformationEfficiencyOutput) {

}
func _TransformationEfficiencyRun(_ctx context.Context, input *TransformationEfficiencyInput) *TransformationEfficiencyOutput {
	output := &TransformationEfficiencyOutput{}
	_TransformationEfficiencySetup(_ctx, input)
	_TransformationEfficiencySteps(_ctx, input, output)
	_TransformationEfficiencyAnalysis(_ctx, input, output)
	_TransformationEfficiencyValidation(_ctx, input, output)
	return output
}

func TransformationEfficiencyRunSteps(_ctx context.Context, input *TransformationEfficiencyInput) *TransformationEfficiencySOutput {
	soutput := &TransformationEfficiencySOutput{}
	output := _TransformationEfficiencyRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func TransformationEfficiencyNew() interface{} {
	return &TransformationEfficiencyElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &TransformationEfficiencyInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _TransformationEfficiencyRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &TransformationEfficiencyInput{},
			Out: &TransformationEfficiencyOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type TransformationEfficiencyElement struct {
	inject.CheckedRunner
}

type TransformationEfficiencyInput struct {
	ColonyCount           int
	ControlColonyCount    int
	ControlDNA            wtype.DNASequence
	ControlDNAMass        wunit.Mass
	ControlPlateOutVolume wunit.Volume
	DNA                   wtype.DNASequence
	DNAMass               wunit.Mass
	PlateOutVolume        wunit.Volume
	TransformationName    string
}

type TransformationEfficiencyOutput struct {
	ControlTransformationEfficiency    transformationEfficiency
	EquivalentTransformationEfficiency transformationEfficiency
	Errors                             error
	RelativeTransformationEfficiency   relativeTransformationEfficiency
	Warnings                           error
}

type TransformationEfficiencySOutput struct {
	Data struct {
		ControlTransformationEfficiency    transformationEfficiency
		EquivalentTransformationEfficiency transformationEfficiency
		Errors                             error
		RelativeTransformationEfficiency   relativeTransformationEfficiency
		Warnings                           error
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "TransformationEfficiency",
		Constructor: TransformationEfficiencyNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol TransformationEfficiency calculates transformation efficiency based on colony count and transformation parameters.\n",
			Path: "src/github.com/antha-lang/elements/an/TransformationEfficiency/TransformationEfficiency.an",
			Params: []component.ParamDesc{
				{Name: "ColonyCount", Desc: "", Kind: "Parameters"},
				{Name: "ControlColonyCount", Desc: "", Kind: "Parameters"},
				{Name: "ControlDNA", Desc: "", Kind: "Parameters"},
				{Name: "ControlDNAMass", Desc: "", Kind: "Parameters"},
				{Name: "ControlPlateOutVolume", Desc: "", Kind: "Parameters"},
				{Name: "DNA", Desc: "", Kind: "Parameters"},
				{Name: "DNAMass", Desc: "", Kind: "Parameters"},
				{Name: "PlateOutVolume", Desc: "", Kind: "Parameters"},
				{Name: "TransformationName", Desc: "", Kind: "Parameters"},
				{Name: "ControlTransformationEfficiency", Desc: "", Kind: "Data"},
				{Name: "EquivalentTransformationEfficiency", Desc: "", Kind: "Data"},
				{Name: "Errors", Desc: "", Kind: "Data"},
				{Name: "RelativeTransformationEfficiency", Desc: "", Kind: "Data"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
