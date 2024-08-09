package main

import (
	"bufio"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
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

	selectFile1Button := widget.NewButton("Select File 1", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				return
			}
			file1Path = reader.URI().Path()
			file1Label.SetText(fmt.Sprintf("File 1: %s", file1Path))
		}, w)
	})

	selectFile2Button := widget.NewButton("Select File 2", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				return
			}
			file2Path = reader.URI().Path()
			file2Label.SetText(fmt.Sprintf("File 2: %s", file2Path))
		}, w)
	})

	operations := []string{"Intersection", "Union", "Difference"}
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
		case "Intersection":
			result = intersection(set1, set2)
		case "Union":
			result = union(set1, set2)
		case "Difference":
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
	w.Resize(fyne.NewSize(400, 300))
	w.ShowAndRun()
}
