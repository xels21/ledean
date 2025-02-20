package parameter

import (
	"encoding/json"
	"flag"
	"path/filepath"
)

type Parameter struct {
	InDir      string `json:"inDir"`
	OutDir     string `json:"outDir"`
	PixelCount int    `json:"pixelCount"`
	Name       string `json:"name"`
	AsByte     bool   `json:"asByte"`
}

func GetParameter() *Parameter {
	var parm Parameter
	flag.StringVar(&parm.InDir, "in", ".", "Path to directory of to converting pictures")
	flag.StringVar(&parm.OutDir, "out", "", "Path to output directory for the converted pictures")
	flag.IntVar(&parm.PixelCount, "pixelCount", 50, "Amount of pixel in one column")
	flag.StringVar(&parm.Name, "name", "picture", "Name for package and out dir")
	flag.BoolVar(&parm.AsByte, "asByte", false, "Defines whether output should be generated as byte array (string)")
	flag.Parse()

	if parm.OutDir == "" {
		parm.OutDir = filepath.Join(parm.InDir, "gen_"+parm.Name)
	}

	return &parm
}

func (self *Parameter) PrintParameter() {
	parmAsJSON, _ := json.Marshal(self)
	println("Starting with:")
	println(string(parmAsJSON))
}
