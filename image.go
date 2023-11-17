package algopack

import (
	"fmt"
	"net/url"
	"path/filepath"
)

type ImageExec interface {
	Id() uint64
	Width() int
	Height() int
	URI() string
	GetImgName() string
	Copy(targetFilePath string) error
	BeforeExec() error
}

type AnnotationGroup struct {
	GroupId       string
	GroupColor    string
	ParentGroupId string
}

type Annotation struct {
	ImgId  uint64
	X1     int
	X2     int
	Y1     int
	Y2     int
	Label  Label
	Score  string // 置信度得分
	Shape  string // 标识标注类型: rectangle表示标注矩形，polygon标注多边形
	Group  AnnotationGroup
	Points string
}

type DemoImage struct {
	Url  string
	Sign int
	W    int
	H    int
}

func (d DemoImage) Id() uint64 {
	return uint64(d.Sign)
}

func (d DemoImage) Width() int {
	return d.W
}

func (d DemoImage) Height() int {
	return d.H
}

func (d DemoImage) URI() string {
	return d.Url
}

func (d DemoImage) BeforeExec() error {
	return nil
}

func (d DemoImage) GetImgName() string {
	parsedURL, err := url.Parse(d.Url)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return ""
	}
	return filepath.Base(parsedURL.Path)
}

func (d DemoImage) Copy(targetFilePath string) error {
	//TODO implement me
	return FileCopy(d.Url, targetFilePath)
}
