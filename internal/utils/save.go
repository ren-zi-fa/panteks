package utils

import (
	"os"
	"path/filepath"
)

func SaveToHTML(data []byte, output ...string) (*string, int, error) {
	var filePath string

	if len(output) > 0 && output[0] != "" {
		filePath = output[0]
	} else {
		filePath = "output.html"
	}
   
	info, err := os.Stat(filePath)
	if err == nil && info.IsDir() {
		filePath = filepath.Join(filePath, "output.html")
	}

	
	folder := filepath.Dir(filePath)
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return nil, 0, err
	}


	file, err := os.Create(filePath)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()


	n, err := file.Write(data)
	if err != nil {
		return nil, 0, err
	}

	
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, n, err
	}

	return &absPath, n, nil
}

func SaveToTXT(data []byte, output ...string) (*string, int, error) {
	var filePath string

	if len(output) > 0 && output[0] != "" {
		filePath = output[0]
	} else {
		filePath = "output.txt"
	}
   
	info, err := os.Stat(filePath)
	if err == nil && info.IsDir() {
		filePath = filepath.Join(filePath, "output.txt")
	}

	
	folder := filepath.Dir(filePath)
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return nil, 0, err
	}


	file, err := os.Create(filePath)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()


	n, err := file.Write(data)
	if err != nil {
		return nil, 0, err
	}

	
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, n, err
	}

	return &absPath, n, nil
}







