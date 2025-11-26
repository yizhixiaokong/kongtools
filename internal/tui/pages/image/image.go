package image

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"kongtools/internal/tui/messages"
	"kongtools/internal/tui/styles"
)

// 自定义消息类型
type imageMsg struct{ output string }
type imageErrMsg struct{ err error }
type imageDownloadedMsg struct{ path string }

// 页面状态
type pageState int

const (
	stateInput pageState = iota
	stateLoading
	stateViewing
)

type ImagePage struct {
	width            int
	height           int
	textInput        textinput.Model
	viewport         viewport.Model
	imageOutput      string
	err              error
	keys             imageKeyMap
	currentSizeIndex int
	currentInput     string // 当前输入的 URL 或路径
	tempFilePath     string // 缓存下载的临时文件路径
	state            pageState
}

type previewSize struct {
	label string
	value string
}

var previewSizes = []previewSize{
	{"Small", "40x20"},
	{"Medium", "80x40"},
	{"Large", "120x60"},
}

type imageKeyMap struct {
	Enter     key.Binding
	Back      key.Binding
	Quit      key.Binding
	CycleSize key.Binding
	Clear     key.Binding
	Left      key.Binding
	Right     key.Binding
}

func (k imageKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.CycleSize, k.Clear, k.Back, k.Quit}
}

func (k imageKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter, k.CycleSize, k.Clear},
		{k.Left, k.Right, k.Back, k.Quit},
	}
}

var imageKeys = imageKeyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "获取图片"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "返回"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "退出"),
	),
	CycleSize: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "切换尺寸"),
	),
	Clear: key.NewBinding(
		key.WithKeys("c", "backspace"),
		key.WithHelp("c", "清除/重置"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "向左滚动"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "向右滚动"),
	),
}

func NewImagePage() *ImagePage {
	ti := textinput.New()
	ti.Placeholder = "输入图片 URL 或本地路径"
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 50

	vp := viewport.New(0, 0)

	return &ImagePage{
		textInput:        ti,
		viewport:         vp,
		imageOutput:      "",
		err:              nil,
		keys:             imageKeys,
		currentSizeIndex: 0,
		currentInput:     "",
		tempFilePath:     "",
		state:            stateInput,
	}
}

// cleanup 清理临时文件
func (m *ImagePage) cleanup() {
	if m.tempFilePath != "" {
		os.Remove(m.tempFilePath)
		m.tempFilePath = ""
	}
}

func (m *ImagePage) Init() tea.Cmd {
	return textinput.Blink
}

func (m *ImagePage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textInput.Width = msg.Width - 4
		m.updateViewportHeight()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Back):
			// 直接返回主菜单
			return m, func() tea.Msg {
				return messages.SwitchPageMsg{Page: "main"}
			}

		case key.Matches(msg, m.keys.Quit):
			m.cleanup()
			return m, tea.Quit

		case key.Matches(msg, m.keys.Clear):
			if m.state == stateViewing {
				m.state = stateInput
				m.imageOutput = ""
				m.viewport.SetContent("")
				m.textInput.Focus()
				m.viewport.SetXOffset(0)
				return m, nil
			}

		case key.Matches(msg, m.keys.Left):
			if m.state == stateViewing {
				m.viewport.ScrollLeft(2)
				return m, nil
			}

		case key.Matches(msg, m.keys.Right):
			if m.state == stateViewing {
				m.viewport.ScrollRight(2)
				return m, nil
			}

		case key.Matches(msg, m.keys.CycleSize):
			if m.state == stateViewing {
				m.currentSizeIndex = (m.currentSizeIndex + 1) % len(previewSizes)
				// 如果当前有内容，尝试重新加载
				if m.currentInput != "" {
					size := previewSizes[m.currentSizeIndex].value
					m.imageOutput = "正在重新转换..."
					m.viewport.SetContent(m.imageOutput)
					m.state = stateLoading
					m.viewport.SetXOffset(0) // 重置横向滚动

					// 确定使用哪个路径
					var path string
					if m.tempFilePath != "" {
						path = m.tempFilePath
					} else {
						path = resolveLocalPath(m.currentInput)
					}

					return m, func() tea.Msg { return convertImage(path, size) }
				}
			}
			return m, nil

		case key.Matches(msg, m.keys.Enter):
			if m.state == stateInput {
				// 输入完成后触发下载和转换
				input := strings.TrimSpace(m.textInput.Value())
				if input == "" {
					return m, nil
				}

				// 清理旧的临时文件
				m.cleanup()
				m.currentInput = input
				m.imageOutput = "正在处理..."
				m.viewport.SetContent(m.imageOutput)
				m.err = nil
				m.state = stateLoading
				m.viewport.SetXOffset(0)
				m.updateViewportHeight()
				m.textInput.Blur()

				if isUrl(input) {
					return m, func() tea.Msg { return downloadImage(input) }
				}

				path := resolveLocalPath(input)
				size := previewSizes[m.currentSizeIndex].value
				return m, func() tea.Msg { return convertImage(path, size) }
			}
		}

	case imageDownloadedMsg:
		m.tempFilePath = msg.path
		size := previewSizes[m.currentSizeIndex].value
		return m, func() tea.Msg { return convertImage(m.tempFilePath, size) }

	case imageMsg:
		m.imageOutput = msg.output
		m.viewport.SetContent(msg.output)
		m.viewport.GotoTop()
		m.viewport.SetXOffset(0)
		m.state = stateViewing
		return m, nil

	case imageErrMsg:
		m.err = msg.err
		m.imageOutput = ""
		m.viewport.SetContent("")
		m.state = stateInput // 出错后返回输入模式
		m.textInput.Focus()
		m.updateViewportHeight()
		return m, nil
	}

	if m.state == stateInput {
		m.textInput, cmd = m.textInput.Update(msg)
	}

	var cmdVP tea.Cmd
	if m.state == stateViewing {
		m.viewport, cmdVP = m.viewport.Update(msg)
	}

	return m, tea.Batch(cmd, cmdVP)
}

func (m *ImagePage) View() string {
	header := m.headerView()
	return lipgloss.JoinVertical(lipgloss.Left, header, m.viewport.View())
}

func (m *ImagePage) headerView() string {
	var sb strings.Builder

	// 标题
	title := styles.TitleStyle.Render("🖼️  Image Viewer")
	sb.WriteString(lipgloss.PlaceHorizontal(m.width, lipgloss.Center, title) + "\n\n")

	// 尺寸信息
	size := previewSizes[m.currentSizeIndex]
	var helpText string
	if m.state == stateViewing {
		helpText = fmt.Sprintf("当前尺寸: %s (%s) - 按 Tab 切换, c/Backspace 清除", size.label, size.value)
	} else {
		helpText = fmt.Sprintf("当前尺寸: %s (%s)", size.label, size.value)
	}
	sizeInfo := styles.InfoStyle.Render(helpText)
	sb.WriteString(lipgloss.PlaceHorizontal(m.width, lipgloss.Center, sizeInfo) + "\n\n")

	// 输入框 (仅在输入模式或加载模式显示)
	if m.state == stateInput || m.state == stateLoading {
		sb.WriteString(lipgloss.PlaceHorizontal(m.width, lipgloss.Center, m.textInput.View()) + "\n\n")
	}

	// 状态提示
	if m.state == stateLoading {
		loading := styles.InfoStyle.Render("⏳ " + m.imageOutput)
		sb.WriteString(lipgloss.PlaceHorizontal(m.width, lipgloss.Center, loading) + "\n")
	}

	// 错误提示
	if m.err != nil {
		errView := styles.ErrorStyle.Render("Error: " + m.err.Error())
		sb.WriteString(lipgloss.PlaceHorizontal(m.width, lipgloss.Center, errView) + "\n")
	}

	return sb.String()
}

func (m *ImagePage) updateViewportHeight() {
	headerHeight := lipgloss.Height(m.headerView())
	m.viewport.Width = m.width
	m.viewport.Height = m.height - headerHeight
	if m.viewport.Height < 0 {
		m.viewport.Height = 0
	}
}

func (m *ImagePage) Title() string {
	return "图片预览"
}

func (m *ImagePage) Help() help.KeyMap {
	return m.keys
}

func (m *ImagePage) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.textInput.Width = width - 10 // 留一些边距
	m.updateViewportHeight()
}

// isUrl 判断是否为 URL
func isUrl(input string) bool {
	return strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://")
}

// resolveLocalPath 解析本地路径（处理 ~）
func resolveLocalPath(input string) string {
	if strings.HasPrefix(input, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, input[1:])
		}
	}
	return input
}

// downloadImage 下载图片
func downloadImage(url string) tea.Msg {
	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "bubbletea-img-*.jpg")
	if err != nil {
		return imageErrMsg{err}
	}
	defer tmpFile.Close()

	// 下载图片
	resp, err := http.Get(url)
	if err != nil {
		os.Remove(tmpFile.Name())
		return imageErrMsg{err}
	}
	defer resp.Body.Close()

	// 保存到临时文件
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		os.Remove(tmpFile.Name())
		return imageErrMsg{err}
	}

	return imageDownloadedMsg{path: tmpFile.Name()}
}

// convertImage 调用 chafa 转换图片
func convertImage(filePath string, size string) tea.Msg {
	// 检查 chafa 是否安装
	if _, err := exec.LookPath("chafa"); err != nil {
		return imageErrMsg{fmt.Errorf("chafa 未安装，请先安装: sudo apt install chafa")}
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); err != nil {
		return imageErrMsg{fmt.Errorf("文件不存在或无法访问: %v", err)}
	}

	// 调用 chafa 转换
	// 使用指定的尺寸
	cmd := exec.Command("chafa", "-s", size, filePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return imageErrMsg{fmt.Errorf("chafa error: %v", err)}
	}

	return imageMsg{string(output)}
}
