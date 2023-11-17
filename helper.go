package algopack

import (
	"io"
	"log"
	"os"
)

func UniqueData[K uint | uint32 | uint64, T any](data []T, getKey func(d T) K) []T {
	// 去重标签映射
	uniqueData := make(map[K]T)
	for _, row := range data {
		key := getKey(row)
		if _, ok := uniqueData[key]; !ok {
			uniqueData[key] = row
		}
	}

	var result []T
	for _, value := range uniqueData {
		result = append(result, value)
	}

	return result
}

func GenTrainLabelNames(labels []Label) []TrainLabelName {
	var labelNames []TrainLabelName
	for i, label := range labels {
		n := TrainLabelName{
			LabelId: label.Id,
			Index:   i + 1,
			Name:    label.Name,
		}
		labelNames = append(labelNames, n)
	}
	return labelNames
}

func MkdirAll(paths []string) {
	for _, dirPath := range paths {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func CreateFile(filePath string) (*os.File, error) {
	// 使用 os.Stat() 检查文件是否存在
	_, err := os.Stat(filePath)

	// 如果文件不存在
	if os.IsNotExist(err) {
		// 使用 os.Create() 创建文件
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		return file, nil
	} else if err != nil {
		// 如果发生其他错误，例如权限问题等
		return nil, err
	}
	return nil, err
}

func FileCopy(sourceFilePath string, destinationFilePath string) error {

	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destinationFilePath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func MapKeys[K comparable, V any](m map[K]V) []K {
	values := make([]K, 0)
	for k, _ := range m {
		values = append(values, k)
	}
	return values
}

func MapValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0)
	for _, v := range m {
		values = append(values, v)
	}
	return values
}
