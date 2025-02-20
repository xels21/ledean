package main

import (
	"ledean/picscaler/parameter"
	"ledean/picscaler/picscaler"
)

func main() {
	parm := parameter.GetParameter()
	parm.PrintParameter()
	picScaler := picscaler.NewPicScaler(parm.InDir, parm.OutDir, parm.Name, parm.PixelCount, parm.AsByte)
	picScaler.Scale()
	picScaler.CreateController()
}
