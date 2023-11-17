package test

import (
	"github.com/kordar/algopack"
	"log"
	"testing"
)

func TestYolov5ServiceTest(t *testing.T) {
	yolov5Pack := algopack.NewYolov5Pack()
	params := algopack.Params{
		Labels:    []algopack.Label{{Id: 1, Name: "AAA"}, {Id: 2, Name: "BBB"}},
		Params:    nil,
		Relations: map[uint64]uint64{},
		BasePath:  "/Users/mac/Documents/test/10",
		YamlName:  "data.yaml",
	}
	resParams := yolov5Pack.BeforeExecute(params)
	labelCount, err := yolov5Pack.Execute(resParams, algopack.DemoImage{
		Url:  "/Users/mac/Pictures/bda1dbb041c9c7bf6c7134cfb6d54512.png",
		Sign: 1,
		W:    1023,
		H:    800,
	}, []algopack.Annotation{
		{ImgId: 1, X1: 0, X2: 450, Y1: 3, Y2: 300, Label: algopack.Label{Id: 1, Name: "AAAA"}},
		{ImgId: 1, X1: 0, X2: 450, Y1: 3, Y2: 300, Label: algopack.Label{Id: 2, Name: "BBBB"}},
		{ImgId: 1, X1: 0, X2: 23, Y1: 3, Y2: 3, Label: algopack.Label{Id: 1, Name: "AAAA"}},
	}, false)
	log.Printf("ccc = %+v, eee = %+v", labelCount, err)
}

func TestYolov5Pack(t *testing.T) {
	log.Printf("%d\n", 0%1)
}
