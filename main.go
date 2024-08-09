package main

import (
	"DataSetOperator/utils/calcutils"
	"DataSetOperator/utils/fileutils"
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

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

// 创建一个空白的间隔矩形
func createSpacer(height float32) fyne.CanvasObject {
	// 创建一个透明的矩形
	spacer := canvas.NewRectangle(color.Transparent)

	// 设置矩形的最小尺寸
	spacer.SetMinSize(fyne.NewSize(0, height))
	return spacer
}

func main() {
	// 创建一个新的应用程序实例
	a := app.New()

	// 设置应用程序的主题
	a.Settings().SetTheme(&myTheme{})

	// 创建一个新的窗口
	w := a.NewWindow("SetIntersectUnionDiff")

	// 创建标签用于显示文件 1 和文件 2 的路径
	file1Label := widget.NewLabel("文件 1: 未选择")
	file2Label := widget.NewLabel("文件 2: 未选择")

	// 定义变量用于存储文件路径
	var file1Path, file2Path string

	// 创建文件对话框的函数
	createFileDialog := func(callback func(string)) *dialog.FileDialog {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				// 如果发生错误，显示错误对话框
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				return
			}
			// 调用回调函数并传递文件路径
			callback(reader.URI().Path())
		}, w)

		// 设置文件对话框的尺寸
		fd.Resize(fyne.NewSize(800, 600))
		return fd
	}

	// 创建按钮用于选择文件 1
	selectFile1Button := widget.NewButton("选择文件 1", func() {
		// 打开文件对话框选择文件 1
		fd := createFileDialog(func(path string) {
			file1Path = path
			// 更新文件 1 标签的文本
			file1Label.SetText(fmt.Sprintf("文件 1: %s", file1Path))
		})
		fd.Show()
	})

	// 创建按钮用于选择文件 2
	selectFile2Button := widget.NewButton("选择文件 2", func() {
		// 打开文件对话框选择文件 2
		fd := createFileDialog(func(path string) {
			file2Path = path
			// 更新文件 2 标签的文本
			file2Label.SetText(fmt.Sprintf("文件 2: %s", file2Path))
		})
		fd.Show()
	})

	operations := []string{"交集", "并集", "差集"}

	// 创建选择操作的下拉框
	operationSelect := widget.NewSelect(operations, nil)
	operationSelect.PlaceHolder = "请选择操作"

	// 创建一个文本框用于显示结果
	resultEntry := widget.NewMultiLineEntry()
	resultEntry.Disable()
	resultEntry.Wrapping = fyne.TextWrapBreak

	// 创建一个按钮用于计算结果
	calculateButton := widget.NewButton("计算", func() {
		if file1Path == "" || file2Path == "" {
			dialog.ShowError(fmt.Errorf("请选择两个文件"), w)
			return
		}

		set1, err1 := fileutils.ReadFile(file1Path)
		if err1 != nil {
			dialog.ShowError(err1, w)
			return
		}

		set2, err2 := fileutils.ReadFile(file2Path)
		if err2 != nil {
			dialog.ShowError(err2, w)
			return
		}

		var result map[string]struct{}

		switch operationSelect.Selected {
		case "交集":
			result = calcutils.Intersection(set1, set2)
		case "并集":
			result = calcutils.Union(set1, set2)
		case "差集":
			result = calcutils.Difference(set1, set2)
		default:
			dialog.ShowError(fmt.Errorf("无效的操作"), w)
			return
		}

		// 将结果写入文本框
		var sb strings.Builder
		for k := range result {
			sb.WriteString(k)
			sb.WriteString("\n")
		}
		resultEntry.SetText(sb.String())
	})

	// 创建一个按钮用于复制结果
	copyButton := widget.NewButton("复制结果", func() {
		w.Clipboard().SetContent(resultEntry.Text)
		dialog.ShowInformation("复制成功", "结果已复制到剪贴板", w)
	})

	// 创建一个滚动容器用于显示结果
	scrollContainer := container.NewScroll(resultEntry)

	// 设置滚动容器的最小尺寸
	scrollContainer.SetMinSize(fyne.NewSize(600, 300))

	// 创建主界面布局
	content := container.NewBorder(
		// 顶部容器，包含文件选择按钮、标签和操作选择
		container.NewVBox(
			selectFile1Button,
			file1Label,
			createSpacer(10),
			selectFile2Button,
			file2Label,
			createSpacer(10),
			operationSelect,
			createSpacer(10),
			calculateButton,
		),
		// 底部容器，包含复制按钮
		container.NewHBox(copyButton), // 将复制按钮放在底部
		nil,                           // 左侧容器为空
		nil,                           // 右侧容器为空
		// 中间容器，包含滚动容器
		scrollContainer,
	)

	// 创建带内边距的容器
	paddedContent := container.NewPadded(content)

	// 设置窗口内容
	w.SetContent(paddedContent)

	// 设置窗口尺寸
	w.Resize(fyne.NewSize(800, 600))

	// 显示并运行窗口
	w.ShowAndRun()
}
