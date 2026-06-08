package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func main() {
	var scatter, swapxy bool
	var w, h int
	var yMin string
	flag.BoolVar(&scatter, "scatter", false, "draw points only")
	flag.BoolVar(&swapxy, "swap-xy", false, "swap the X and Y axes")
	flag.IntVar(&w, "width", 8, "")
	flag.IntVar(&h, "height", 4, "")
	flag.StringVar(&yMin, "y-min", "", "")
	flag.Parse()

	if flag.NArg() < 1 {
		panic("usage: plot [--scatter] <csv-file>")
	}

	path := flag.Arg(0)
	file := Must(os.Open(path))
	defer file.Close()

	reader := csv.NewReader(file)
	records := Must(reader.ReadAll())

	x, y := make([]float64, 0, len(records)), make([]float64, 0, len(records))
	for _, r := range records {
		x = append(x, Must(strconv.ParseFloat(r[0], 64)))
		y = append(y, Must(strconv.ParseFloat(r[1], 64)))
	}

	if swapxy {
		x, y = y, x
	}

	out := strings.TrimSuffix(path, filepath.Ext(path)) + ".png"
	switch {
	case scatter:
		if err := SaveAsScatter(x, y, w, h, out); err != nil {
			panic(err)
		}
	default:
		if err := Save(x, y, w, h, yMin, out); err != nil {
			panic(err)
		}
	}
}

func Save(x, y []float64, w, h int, yMin, filename string) error {
	xys := make(plotter.XYs, 0, len(x))
	for i := range x {
		xys = append(xys, plotter.XY{
			X: x[i],
			Y: y[i],
		})
	}

	p := plot.New()
	line, err := plotter.NewLine(xys)
	if err != nil {
		return fmt.Errorf("new line: %v", err)
	}

	line.Color = color.RGBA{R: 0, G: 120, B: 255, A: 255}
	p.Add(line)
	p.Add(plotter.NewGrid())

	if len(yMin) > 0 {
		ymin, err := strconv.Atoi(yMin)
		if err != nil {
			return err
		}

		p.Y.Min = float64(ymin)
	}

	wInch := font.Length(w) * vg.Inch
	hInch := font.Length(h) * vg.Inch
	if err := p.Save(wInch, hInch, filename); err != nil {
		return fmt.Errorf("save: %v", err)
	}

	return nil
}

func SaveAsScatter(x, y []float64, w, h int, filename string) error {
	xys := make(plotter.XYs, 0, len(x))
	for i := range x {
		xys = append(xys, plotter.XY{
			X: x[i],
			Y: y[i],
		})
	}

	p := plot.New()
	scatter, err := plotter.NewScatter(xys)
	if err != nil {
		return fmt.Errorf("new scatter: %v", err)
	}

	scatter.GlyphStyle.Shape = draw.CircleGlyph{}
	scatter.GlyphStyle.Color = color.RGBA{R: 0, G: 120, B: 255, A: 255}
	scatter.GlyphStyle.Radius = vg.Points(2)

	p.X.Min = 0
	p.X.Max = 2 * math.Pi
	p.Y.Min = 0
	p.Y.Max = math.Pi
	p.X.Tick.Marker = plot.ConstantTicks(Ticks2Pi())
	p.Y.Tick.Marker = plot.ConstantTicks(Ticks2Pi())
	p.Add(scatter)
	p.Add(plotter.NewGrid())

	wInch := font.Length(w) * vg.Inch
	hInch := font.Length(h) * vg.Inch
	if err := p.Save(wInch, hInch, filename); err != nil {
		return fmt.Errorf("save: %v", err)
	}

	return nil
}

func Ticks2Pi() []plot.Tick {
	return []plot.Tick{
		{Value: 0, Label: "0"},
		{Value: 1 * math.Pi / 4, Label: "pi/4"},
		{Value: 2 * math.Pi / 4, Label: "2pi/4"},
		{Value: 3 * math.Pi / 4, Label: "3pi/4"},
		{Value: 4 * math.Pi / 4, Label: "pi"},
		{Value: 5 * math.Pi / 4, Label: "5pi/4"},
		{Value: 6 * math.Pi / 4, Label: "6pi/4"},
		{Value: 7 * math.Pi / 4, Label: "7pi/4"},
		{Value: 8 * math.Pi / 4, Label: "2pi"},
	}
}

func Must[T any](a T, err error) T {
	if err != nil {
		panic(err)
	}

	return a
}
