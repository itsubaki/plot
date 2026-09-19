package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	var xMin, xMax, yMin, yMax string
	flag.StringVar(&xMin, "x-min", "", "")
	flag.StringVar(&xMax, "x-max", "", "")
	flag.StringVar(&yMin, "y-min", "", "")
	flag.StringVar(&yMax, "y-max", "", "")
	flag.Parse()

	if flag.NArg() < 1 {
		panic("usage: plot <csv-file>")
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

	out := strings.TrimSuffix(path, filepath.Ext(path)) + ".png"
	if err := Save(x, y, xMin, xMax, yMin, yMax, out); err != nil {
		panic(err)
	}
}

func Save(x, y []float64, xMin, xMax, yMin, yMax, filename string) error {
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

	if err := set(xMin, &p.X.Min); err != nil {
		return err
	}

	if err := set(xMax, &p.X.Max); err != nil {
		return err
	}

	if err := set(yMin, &p.Y.Min); err != nil {
		return err
	}

	if err := set(yMax, &p.Y.Max); err != nil {
		return err
	}

	if err := p.Save(8*vg.Inch, 4*vg.Inch, filename); err != nil {
		return fmt.Errorf("save: %v", err)
	}

	return nil
}

func Must[T any](a T, err error) T {
	if err != nil {
		panic(err)
	}

	return a
}

func set(s string, dst *float64) error {
	if s == "" {
		return nil
	}

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}

	*dst = v
	return nil
}
