package algopack

import (
	"errors"
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type Yolov5Pack struct {
}

func NewYolov5Pack() *Yolov5Pack {
	return &Yolov5Pack{}
}

func (yol Yolov5Pack) BeforeExecute(params Params) Params {
	trainLabelNames := GenTrainLabelNames(params.Labels)
	params.TrainLabelName = trainLabelNames
	params.TrainLabelIdIndex = GenLabelNameIndex(trainLabelNames)
	yol.GenerateConfigFile(params)
	return params
}

func (yol Yolov5Pack) Execute(params Params, image ImageExec, annotations []Annotation, toVal bool, args ...any) (map[uint64]int, error) {
	//TODO implement me
	labelIdCount := map[uint64]int{}
	if er := image.BeforeExec(); er != nil {
		return labelIdCount, er
	}

	content := ""
	imgName := image.GetImgName()
	for _, annotation := range annotations {
		// 获取标注关联labelId的实际labelId和该id对应的index值
		realLabelId := GetRealLabelId(annotation.Label.Id, params.Relations)
		labelIndex, ok := params.TrainLabelIdIndex[realLabelId]
		if !ok {
			log.Printf("未查询到有效标签\n")
			continue
		}

		dw := float64(1) / float64(image.Width())
		dh := float64(1) / float64(image.Height())
		x := float64(annotation.X1+annotation.X2) / 2.0
		y := float64(annotation.Y1+annotation.Y2) / 2.0
		w := float64(annotation.X2 - annotation.X1)
		h := float64(annotation.Y2 - annotation.Y1)
		x = x * dw
		w = w * dw
		y = y * dh
		h = h * dh

		if x > 1 || y > 1 || w > 1 || h > 1 {
			log.Printf("标注坐标错误:%+v\n", annotation)
			return labelIdCount, errors.New(fmt.Sprintf("%s:[x=%f,y=%f,w=%f,h=%f]标注坐标错误", imgName, x, y, w, h))
		}

		content += fmt.Sprintf("%d %f %f %f %f\n", labelIndex, x, y, w, h)
		labelIdCount[realLabelId] += 1
	}

	ext := filepath.Ext(imgName)
	txtFileName := strings.ReplaceAll(imgName, ext, ".txt")

	var imgFilepath string
	var txtFilepath string
	index := ""
	if args[0] != nil {
		index = args[0].(string)
	}
	if toVal {
		imgFilepath = path.Join(params.BasePath, index, "images", "val", imgName)
		txtFilepath = path.Join(params.BasePath, index, "labels", "val", txtFileName)
	} else {
		imgFilepath = path.Join(params.BasePath, index, "images", "train", imgName)
		txtFilepath = path.Join(params.BasePath, index, "labels", "train", txtFileName)
	}

	if err := image.Copy(imgFilepath); err != nil {
		return labelIdCount, err
	}

	file, err := CreateFile(txtFilepath)
	if err != nil {
		return labelIdCount, errors.New(fmt.Sprintf("%s:txt文件创建异常=%v", imgName, err))
	}

	defer file.Close()

	_, err2 := file.Write([]byte(content))
	if err2 != nil {
		return labelIdCount, errors.New(fmt.Sprintf("%s:txt文件写入异常=%v", imgName, err2))
	}

	return labelIdCount, nil
}

func (yol Yolov5Pack) GenerateConfigFile(params Params) {

	if len(params.Subs) == 0 {
		yol.g(params)
		return
	}

	yol.sg(params)
}

func (yol Yolov5Pack) g(params Params) {
	configFilepath := path.Join(params.BasePath, "config")
	imgTrainFilepath := path.Join(params.BasePath, "images", "train")
	imgValFilepath := path.Join(params.BasePath, "images", "val")
	MkdirAll([]string{
		configFilepath,
		path.Join(params.BasePath, "images"),
		imgTrainFilepath,
		imgValFilepath,
		path.Join(params.BasePath, "labels"),
		path.Join(params.BasePath, "labels", "train"),
		path.Join(params.BasePath, "labels", "val"),
	})

	yamlPath := path.Join(configFilepath, params.YamlName)
	file, err := CreateFile(yamlPath)
	if err != nil {
		return
	}

	defer file.Close()
	nameSize := len(params.TrainLabelName)
	content := "train: " + imgTrainFilepath + "\n"
	content += "val: " + imgValFilepath + "\n"
	content += "nc: " + strconv.Itoa(nameSize) + "\n"
	content += "names: ["
	if len(params.TrainLabelName) > 0 {
		for _, labelName := range params.TrainLabelName {
			content += "'" + labelName.Name + "',"
		}
		content = content[:len(content)-1]
	}
	content += "]"
	_, err2 := file.Write([]byte(content))
	if err2 != nil {
		log.Panicln(err2)
	}
}

func (yol Yolov5Pack) sg(params Params) {
	configFilepath := path.Join(params.BasePath, "config")

	trainPaths := make([]string, 0)
	valPaths := make([]string, 0)
	for _, index := range params.Subs {
		imgTrainFilepath := path.Join(params.BasePath, index, "images", "train")
		imgValFilepath := path.Join(params.BasePath, index, "images", "val")
		MkdirAll([]string{
			configFilepath,
			path.Join(params.BasePath, index, "images"),
			imgTrainFilepath,
			imgValFilepath,
			path.Join(params.BasePath, index, "labels"),
			path.Join(params.BasePath, index, "labels", "train"),
			path.Join(params.BasePath, index, "labels", "val"),
		})
		trainPaths = append(trainPaths, imgTrainFilepath)
		valPaths = append(valPaths, imgValFilepath)
	}

	yamlPath := path.Join(configFilepath, params.YamlName)
	file, err := CreateFile(yamlPath)
	if err != nil {
		return
	}

	defer file.Close()
	nameSize := len(params.TrainLabelName)
	content := "train: " + fmt.Sprintf("[%s]", strings.Join(trainPaths, ", ")) + "\n"
	content += "val: " + fmt.Sprintf("[%s]", strings.Join(valPaths, ", ")) + "\n"
	content += "nc: " + strconv.Itoa(nameSize) + "\n"
	content += "names: ["
	if len(params.TrainLabelName) > 0 {
		for _, labelName := range params.TrainLabelName {
			content += "'" + labelName.Name + "',"
		}
		content = content[:len(content)-1]
	}
	content += "]"
	_, err2 := file.Write([]byte(content))
	if err2 != nil {
		log.Panicln(err2)
	}
}
