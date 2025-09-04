package web

import (
	_"bytes"
	_"fmt"
	"os"
	"path/filepath"

	_"github.com/PuerkitoBio/goquery"
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

// func SaveToTXT(data []byte, selector ...string, output ...string) (*string, int, error) {
// 	// Tentukan nama file output
// 	fileName := "output.txt"
// 	if len(output) > 0 && output[0] != "" {
// 		fileName = output[0]
// 	}

// 	// Pastikan folder ada
// 	dir := filepath.Dir(fileName)
// 	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
// 		return nil, 0, fmt.Errorf("failed to create folder: %w", err)
// 	}

// 	// Parsing HTML
// 	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to parse HTML: %w", err)
// 	}

// 	// Ambil konten sesuai selector
// 	var result string
// 	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
// 		text := s.Text()
// 		if text != "" {
// 			result += text + "\n"
// 		}
// 	})

// 	if result == "" {
// 		return nil, 0, fmt.Errorf("selector '%s' not found or empty", selector)
// 	}

// 	// Simpan ke file
// 	if err := os.WriteFile(fileName, []byte(result), 0644); err != nil {
// 		return nil, 0, fmt.Errorf("failed to write file: %w", err)
// 	}

// 	return &fileName, len(result), nil
// }





