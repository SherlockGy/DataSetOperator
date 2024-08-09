package main

import (
	"bufio"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"os"
)

func readFile(filename string) (map[string]struct{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	set := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		set[scanner.Text()] = struct{}{}
	}
	return set, scanner.Err()
}

func intersection(set1, set2 map[string]struct{}) map[string]struct{} {
	result := make(map[string]struct{})
	for k := range set1 {
		if _, ok := set2[k]; ok {
			result[k] = struct{}{}
		}
	}
	return result
}

func union(set1, set2 map[string]struct{}) map[string]struct{} {
	result := make(map[string]struct{})
	for k := range set1 {
		result[k] = struct{}{}
	}
	for k := range set2 {
		result[k] = struct{}{}
	}
	return result
}

func difference(set1, set2 map[string]struct{}) map[string]struct{} {
	result := make(map[string]struct{})
	for k := range set1 {
		if _, ok := set2[k]; !ok {
			result[k] = struct{}{}
		}
	}
	return result
}

func main() {
	a := app.New()
	w := a.NewWindow("SetIntersectUnionDiff")

	file1Label := widget.NewLabel("File 1: Not selected")
	file2Label := widget.NewLabel("File 2: Not selected")

	var file1Path, file2Path string

	createFileDialog := func(callback func(string)) *dialog.FileDialog {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				return
			}
			callback(reader.URI().Path())
		}, w)
		fd.Resize(fyne.NewSize(800, 600)) // 设置文件选择窗口的大小
		return fd
	}

	selectFile1Button := widget.NewButton("Select File 1", func() {
		fd := createFileDialog(func(path string) {
			file1Path = path
			file1Label.SetText(fmt.Sprintf("File 1: %s", file1Path))
		})
		fd.Show()
	})

	selectFile2Button := widget.NewButton("Select File 2", func() {
		fd := createFileDialog(func(path string) {
			file2Path = path
			file2Label.SetText(fmt.Sprintf("File 2: %s", file2Path))
		})
		fd.Show()
	})

	operations := []string{"交集", "并集", "差集"}
	operationSelect := widget.NewSelect(operations, nil)

	resultArea := widget.NewMultiLineEntry()
	resultArea.Disable()

	calculateButton := widget.NewButton("Calculate", func() {
		if file1Path == "" || file2Path == "" {
			dialog.ShowError(fmt.Errorf("please select both files"), w)
			return
		}

		set1, err1 := readFile(file1Path)
		if err1 != nil {
			dialog.ShowError(err1, w)
			return
		}

		set2, err2 := readFile(file2Path)
		if err2 != nil {
			dialog.ShowError(err2, w)
			return
		}

		var result map[string]struct{}

		switch operationSelect.Selected {
		case "交集":
			result = intersection(set1, set2)
		case "并集":
			result = union(set1, set2)
		case "差集":
			result = difference(set1, set2)
		default:
			dialog.ShowError(fmt.Errorf("invalid operation"), w)
			return
		}

		resultText := ""
		for k := range result {
			resultText += k + "\n"
		}
		resultArea.SetText(resultText)
	})

	content := container.NewVBox(
		selectFile1Button,
		file1Label,
		selectFile2Button,
		file2Label,
		operationSelect,
		calculateButton,
		resultArea,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(800, 500)) // 增大主窗口的尺寸
	w.ShowAndRun()
}
