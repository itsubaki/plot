package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"image/color"
	"math"
	"os"
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func main() {
	var sphere bool
	flag.BoolVar(&sphere, "sphere", false, "draw points only and clamp x-axis to 0..2pi")
	flag.Parse()

	if flag.NArg() < 1 {
		panic("usage: plot [--sphere] <csv-file>")
	}

	filename := flag.Arg(0)
	file := Must(os.Open(filename))
	defer file.Close()

	reader := csv.NewReader(file)
	records := Must(reader.ReadAll())

	x, y := make([]float64, 0, len(records)), make([]float64, 0, len(records))
	for _, r := range records {
		x = append(x, Must(strconv.ParseFloat(r[0], 64)))
		y = append(y, Must(strconv.ParseFloat(r[1], 64)))
	}

	switch {
	case sphere:
		if err := SaveAsSphere(x, y, filename+".png"); err != nil {
			panic(err)
		}
	default:
		if err := Save(x, y, filename+".png"); err != nil {
			panic(err)
		}
	}
}

func Save(x, y []float64, filename string) error {
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
	p.Add(line)

	if err := p.Save(8*vg.Inch, 4*vg.Inch, filename); err != nil {
		return fmt.Errorf("save: %v", err)
	}

	return nil
}

func SaveAsSphere(x, y []float64, filename string) error {
	xys := make(plotter.XYs, 0, len(x))
	for i := range x {
		xys = append(xys, plotter.XY{
			X: y[i],
			Y: x[i],
		})
	}

	p := plot.New()
	scatter, err := plotter.NewScatter(xys)
	if err != nil {
		return fmt.Errorf("new scatter: %v", err)
	}

	scatter.GlyphStyle.Color = color.RGBA{R: 30, G: 144, B: 255, A: 90}
	scatter.GlyphStyle.Radius = vg.Points(0.8)
	scatter.GlyphStyle.Shape = draw.CircleGlyph{}

	p.X.Min = 0
	p.X.Max = 2 * math.Pi
	p.Y.Min = 0
	p.Y.Max = math.Pi
	p.X.Tick.Marker = plot.ConstantTicks(Ticks2Pi())
	p.Y.Tick.Marker = plot.ConstantTicks(TicksPi())
	p.Add(scatter)

	if err := p.Save(8*vg.Inch, 4*vg.Inch, filename); err != nil {
		return fmt.Errorf("save: %v", err)
	}

	return nil
}

func TicksPi() []plot.Tick {
	return []plot.Tick{
		{Value: 0, Label: "0"},
		{Value: math.Pi / 4, Label: "pi/4"},
		{Value: math.Pi / 2, Label: "2pi/4"},
		{Value: 3 * math.Pi / 4, Label: "3pi/4"},
		{Value: math.Pi, Label: "pi"},
	}
}

func Ticks2Pi() []plot.Tick {
	return []plot.Tick{
		{Value: 0, Label: "0"},
		{Value: math.Pi / 4, Label: "pi/4"},
		{Value: math.Pi / 2, Label: "2pi/4"},
		{Value: 3 * math.Pi / 4, Label: "3pi/4"},
		{Value: math.Pi, Label: "pi"},
		{Value: 5 * math.Pi / 4, Label: "5pi/4"},
		{Value: 3 * math.Pi / 2, Label: "6pi/4"},
		{Value: 7 * math.Pi / 4, Label: "7pi/4"},
		{Value: 2 * math.Pi, Label: "2pi"},
	}
}

func Must[T any](a T, err error) T {
	if err != nil {
		panic(err)
	}

	return a
}
