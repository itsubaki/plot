package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	filename := os.Args[1]
	file := Must(os.Open(filename))
	defer file.Close()

	reader := csv.NewReader(file)
	records := Must(reader.ReadAll())

	x, y := make([]float64, 0), make([]float64, 0)
	for _, r := range records {
		x = append(x, Must(strconv.ParseFloat(r[0], 64)))
		y = append(y, Must(strconv.ParseFloat(r[1], 64)))
	}

	if err := Save(x, y, filename+".png"); err != nil {
		panic(err)
	}
}

func Must[T any](a T, err error) T {
	if err != nil {
		panic(err)
	}

	return a
}

func Save(x, y []float64, filename string) error {
	xys := make(plotter.XYs, 0)
	for i := range x {
		xys = append(xys, plotter.XY{
			X: x[i],
			Y: y[i],
		})
	}

	line, err := plotter.NewLine(xys)
	if err != nil {
		return fmt.Errorf("plotter newline: %v", err)
	}

	p := plot.New()
	p.Add(line)

	if err := p.Save(8*vg.Inch, 4*vg.Inch, filename); err != nil {
		return fmt.Errorf("save: %v", err)
	}

	return nil
}
