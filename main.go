package main

import (
	"bufio"
	"fmt"
	"image/color"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
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

type myTheme struct{}

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameDisabled {
		return color.White
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m myTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func main() {
	a := app.New()
	a.Settings().SetTheme(&myTheme{})
	w := a.NewWindow("SetIntersectUnionDiff")

	file1Label := widget.NewLabel("文件 1: 未选择")
	file2Label := widget.NewLabel("文件 2: 未选择")

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

	selectFile1Button := widget.NewButton("选择文件 1", func() {
		fd := createFileDialog(func(path string) {
			file1Path = path
			file1Label.SetText(fmt.Sprintf("文件 1: %s", file1Path))
		})
		fd.Show()
	})

	selectFile2Button := widget.NewButton("选择文件 2", func() {
		fd := createFileDialog(func(path string) {
			file2Path = path
			file2Label.SetText(fmt.Sprintf("文件 2: %s", file2Path))
		})
		fd.Show()
	})

	operations := []string{"交集", "并集", "差集"}
	operationSelect := widget.NewSelect(operations, nil)
	operationSelect.PlaceHolder = "请选择操作"

	resultText := canvas.NewText("", color.White)
	resultText.TextStyle.Monospace = true

	calculateButton := widget.NewButton("计算", func() {
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
			dialog.ShowError(fmt.Errorf("非法的操作类型"), w)
			return
		}

		// 使用 strings.Builder 来高效地构建字符串
		var sb strings.Builder
		for k := range result {
			sb.WriteString(k)
			sb.WriteString(string(os.PathListSeparator))
		}
		outputText := sb.String()

		// 移除最后一个多余的换行符
		if len(outputText) > 0 {
			outputText = outputText[:len(outputText)-1]
		}

		resultText.Text = outputText
		resultText.Refresh()
	})

	content := container.NewVBox(
		selectFile1Button,
		file1Label,
		selectFile2Button,
		file2Label,
		operationSelect,
		calculateButton,
		container.NewScroll(resultText), // 使用滚动容器包装resultText
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(800, 500))
	w.ShowAndRun()
}
