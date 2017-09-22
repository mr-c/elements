// Protocol PlotSpreadsheet creates a plot from a spreadsheet
package lib

import (
	"bytes"
	"strconv"
	"strings"
	"unicode"

	"context"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
	"github.com/tealeg/xlsx"
)

// Input parameters for this protocol (data)

// Which sheet to load from XlsFile. 0 is the first sheet.

// Address range (e.g., A1:A8, B7:G7) of X values in sheet.

// Address ranges of Y values in sheet (e.g., ["B4:B16", "C4:C16"])

// Output plot format. Supported values are eps, jpg, jpeg, pdf, png,
// tif and tiff. Default is png.

// Plot width with units. Default is "10cm".

// Plot height with units. Default is "10cm".

// Data which is returned from this protocol, and data types

// Generated plot

func _PlotSpreadsheetRequirements() {

}

// Conditions to run on startup
func _PlotSpreadsheetSetup(_ctx context.Context, _input *PlotSpreadsheetInput) {

}

// The core process for this protocol, with the steps to be performed for every
// input
func _PlotSpreadsheetSteps(_ctx context.Context, _input *PlotSpreadsheetInput, _output *PlotSpreadsheetOutput) {
	// Convert spreadsheet address to index
	index := func(address string) (col, row int) {
		idx := strings.IndexFunc(address, unicode.IsDigit)
		colStr := address
		rowStr := ""
		if idx >= 0 {
			colStr = address[:idx]
			rowStr = address[idx:]
		}
		row, _ = strconv.Atoi(rowStr)
		col = wutil.AlphaToNum(strings.ToUpper(colStr))
		if row > 0 {
			row = row - 1
		}
		if col > 0 {
			col = col - 1
		}
		return
	}

	// Read cell range from sheet
	readCells := func(sheet *xlsx.Sheet, cellRange string) []float64 {
		c := strings.SplitN(cellRange, ":", 2)
		if len(c) < 2 {
			c = []string{"A1", "A1"}
		}
		startCol, startRow := index(c[0])
		endCol, endRow := index(c[1])

		if startCol >= endCol {
			startCol, endCol = endCol, startCol
		}
		if startRow >= endRow {
			startRow, endRow = endRow, startRow
		}

		total := (endRow - startRow + 1) * (endCol - startCol + 1)
		vs := make([]float64, total)
		idx := 0
		for i := startRow; i < endRow+1; i++ {
			for j := startCol; j < endCol+1; j++ {
				v, _ := sheet.Rows[i].Cells[j].Float()
				if len(sheet.Rows) <= i {
					execute.Errorf(_ctx, "invalid cell %s%d", wutil.NumToAlpha(j+1), i+1)
				}
				if len(sheet.Rows[i].Cells) <= j {
					execute.Errorf(_ctx, "invalid cell %s%d", wutil.NumToAlpha(j+1), i+1)
				}
				vs[idx] = v
				idx++
			}
		}
		return vs
	}

	data, _ := _input.XlsFile.ReadAll()
	file, err := xlsx.OpenBinary(data)
	if err != nil {
		execute.Errorf(_ctx, "cannot read spreadsheet: %s", err)
	}

	sheet := file.Sheets[_input.Sheet]

	xvalues := readCells(sheet, _input.XRange)
	var yvalues [][]float64
	for idx, yrange := range _input.YRanges {
		ys := readCells(sheet, yrange)
		if lx, ly := len(xvalues), len(ys); lx != ly {
			execute.Errorf(_ctx, "number of x values (%d) does not equal number of y values (%d) for %dth y range",
				lx, ly, idx)
		}
		yvalues = append(yvalues, ys)
	}

	var points []interface{}
	for _, ys := range yvalues {
		xys := make(plotter.XYs, len(ys))
		for i, n := 0, len(ys); i < n; i++ {
			xys[i].X = xvalues[i]
			xys[i].Y = ys[i]
		}
		points = append(points, xys)
	}

	// now plot the graph
	plt, err := plot.New()
	if err != nil {
		execute.Errorf(_ctx, "cannot make plot: %s", err)
	}
	plotutil.AddScatters(plt, points...)

	width, err := vg.ParseLength(_input.PlotWidth)
	if err != nil {
		width = 10 * vg.Centimeter
	}

	height, err := vg.ParseLength(_input.PlotWidth)
	if err != nil {
		height = 10 * vg.Centimeter
	}

	if len(_input.PlotFormat) == 0 {
		_input.PlotFormat = "png"
	}

	w, err := plt.WriterTo(width, height, _input.PlotFormat)
	if err != nil {
		execute.Errorf(_ctx, "cannot write plot: %s", err)
	}

	var out bytes.Buffer
	w.WriteTo(&out)

	_output.Plot.Name = "plot." + _input.PlotFormat
	_output.Plot.WriteAll(out.Bytes())
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PlotSpreadsheetAnalysis(_ctx context.Context, _input *PlotSpreadsheetInput, _output *PlotSpreadsheetOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PlotSpreadsheetValidation(_ctx context.Context, _input *PlotSpreadsheetInput, _output *PlotSpreadsheetOutput) {

}
func _PlotSpreadsheetRun(_ctx context.Context, input *PlotSpreadsheetInput) *PlotSpreadsheetOutput {
	output := &PlotSpreadsheetOutput{}
	_PlotSpreadsheetSetup(_ctx, input)
	_PlotSpreadsheetSteps(_ctx, input, output)
	_PlotSpreadsheetAnalysis(_ctx, input, output)
	_PlotSpreadsheetValidation(_ctx, input, output)
	return output
}

func PlotSpreadsheetRunSteps(_ctx context.Context, input *PlotSpreadsheetInput) *PlotSpreadsheetSOutput {
	soutput := &PlotSpreadsheetSOutput{}
	output := _PlotSpreadsheetRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PlotSpreadsheetNew() interface{} {
	return &PlotSpreadsheetElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PlotSpreadsheetInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PlotSpreadsheetRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PlotSpreadsheetInput{},
			Out: &PlotSpreadsheetOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wtype.FALSE
	_ = wunit.Make_units
)

type PlotSpreadsheetElement struct {
	inject.CheckedRunner
}

type PlotSpreadsheetInput struct {
	PlotFormat string
	PlotHeight string
	PlotWidth  string
	Sheet      int
	XRange     string
	XlsFile    wtype.File
	YRanges    []string
}

type PlotSpreadsheetOutput struct {
	Plot wtype.File
}

type PlotSpreadsheetSOutput struct {
	Data struct {
		Plot wtype.File
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PlotSpreadsheet",
		Constructor: PlotSpreadsheetNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol PlotSpreadsheet creates a plot from a spreadsheet\n",
			Path: "src/github.com/antha-lang/elements/an/PlotSpreadsheet/element.an",
			Params: []component.ParamDesc{
				{Name: "PlotFormat", Desc: "Output plot format. Supported values are eps, jpg, jpeg, pdf, png,\ntif and tiff. Default is png.\n", Kind: "Parameters"},
				{Name: "PlotHeight", Desc: "Plot height with units. Default is \"10cm\".\n", Kind: "Parameters"},
				{Name: "PlotWidth", Desc: "Plot width with units. Default is \"10cm\".\n", Kind: "Parameters"},
				{Name: "Sheet", Desc: "Which sheet to load from XlsFile. 0 is the first sheet.\n", Kind: "Parameters"},
				{Name: "XRange", Desc: "Address range (e.g., A1:A8, B7:G7) of X values in sheet.\n", Kind: "Parameters"},
				{Name: "XlsFile", Desc: "", Kind: "Parameters"},
				{Name: "YRanges", Desc: "Address ranges of Y values in sheet (e.g., [\"B4:B16\", \"C4:C16\"])\n", Kind: "Parameters"},
				{Name: "Plot", Desc: "Generated plot\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
