package nimvle

import (
	"fmt"
	"strings"

	"github.com/neovim/go-client/nvim"
)

const (
	getWindowList = `vivid#getWindowList()`
)

// Nimvle provides nvim methods around the nvim.Nvim interface.
type Nimvle struct {
	v          *nvim.Nvim
	pluginName string
}

// New makes a new Nvimutils object for the specified nvim.Nvim.
func New(v *nvim.Nvim, name string) *Nimvle {
	return &Nimvle{
		v:          v,
		pluginName: name,
	}
}

// Log is a wrapped `:echom`
func (n *Nimvle) Log(message interface{}) error {
	return n.v.Command("echom '" + n.pluginName + ": " + fmt.Sprintf("%v", message) + "'")
}

// CurrentBufferFilenameExtension obtains the file name extension of the current buffer.
func (n *Nimvle) CurrentBufferFilenameExtension() (string, error) {
	buf, err := n.v.CurrentBuffer()
	if err != nil {
		return "", err
	}

	bufferName, err := n.v.BufferName(buf)
	if err != nil {
		return "", err
	}

	dotName := strings.Split(bufferName, ".")[len(strings.Split(bufferName, "."))-1]

	return dotName, nil
}

// GetContentFromCurrentBuffer obtains the content of the current buffer.
func (n *Nimvle) GetContentFromCurrentBuffer() (string, error) {
	buf, err := n.v.CurrentBuffer()
	if err != nil {
		return "", err
	}

	lines, err := n.v.BufferLines(buf, 0, -1, true)
	if err != nil {
		return "", err
	}

	var content string
	for i, c := range lines {
		content += string(c)
		if i < len(lines)-1 {
			content += "\n"
		}
	}

	return content, nil
}

// GetContentFromBuffer obtains the content of the buffer.
func (n *Nimvle) GetContentFromBuffer(buf nvim.Buffer) (string, error) {
	lines, err := n.v.BufferLines(buf, 0, -1, true)
	if err != nil {
		return "", err
	}

	var content string
	for i, c := range lines {
		content += string(c)
		if i < len(lines)-1 {
			content += "\n"
		}
	}

	return content, nil
}

func (n *Nimvle) setContentToBuffer(buf nvim.Buffer, lines []string) error {
	var byteContent [][]byte
	for _, c := range lines {
		byteContent = append(byteContent, []byte(c))
	}

	return n.v.SetBufferLines(buf, 0, -1, true, byteContent)
}

// SetContentToBuffer writes content to buffer.
func (n *Nimvle) SetContentToBuffer(buf nvim.Buffer, content string) error {
	lines := strings.Split(content, "\n")

	return n.setContentToBuffer(buf, lines)
}

// SetStringerContentToBuffer writes stringer content to buffer.
func (n *Nimvle) SetStringerContentToBuffer(buf nvim.Buffer, str fmt.Stringer) error {
	content := fmt.Sprint(str)
	lines := strings.Split(content, "\n")

	return n.setContentToBuffer(buf, lines)
}

// GetWindowList obtains window list. It depends on `NoahOrberg/vivid.vim`
func (n *Nimvle) GetWindowList() (map[string]int, error) {
	res := make(map[string]int)

	if err := n.v.Eval(getWindowList, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// SplitOpenBuffer divides and opens buffer.
func (n *Nimvle) SplitOpenBuffer(buf nvim.Buffer) error {
	var bwin nvim.Window
	var win nvim.Window

	b := n.v.NewBatch()
	b.CurrentWindow(&bwin)
	b.Command(fmt.Sprintf("sb %d", buf))
	b.CurrentWindow(&win)
	b.SetWindowHeight(win, 15)
	if err := b.Execute(); err != nil {
		return err
	}

	return n.v.SetCurrentWindow(bwin)
}

// NewScratchBuffer creates new scratch buffer.
func (n *Nimvle) NewScratchBuffer(bufferName, ft string) (*nvim.Buffer, error) {
	var scratchBuf nvim.Buffer
	var bwin nvim.Window
	var win nvim.Window

	b := n.v.NewBatch()
	b.CurrentWindow(&bwin)
	b.Command("silent! execute 'new' '" + bufferName + "'")
	b.CurrentBuffer(&scratchBuf)
	b.SetBufferOption(scratchBuf, "buftype", "nofile")
	b.SetBufferOption(scratchBuf, "bufhidden", "hide")
	b.Command("setlocal noswapfile")
	b.Command(
		fmt.Sprintf("set filetype=%v", ft))
	b.Command("setlocal nobuflisted")
	b.SetBufferOption(scratchBuf, "undolevels", -1)
	b.CurrentWindow(&win)
	b.SetWindowHeight(win, 15)

	if err := b.Execute(); err != nil {
		return nil, err
	}

	if err := n.v.SetCurrentWindow(bwin); err != nil {
		return nil, err
	}

	return &scratchBuf, nil
}

// ShowScratchBuffer shows selected buffer.
func (n *Nimvle) ShowScratchBuffer(scratch nvim.Buffer, str fmt.Stringer) error {
	var opened bool

	if err := n.SetStringerContentToBuffer(scratch, str); err != nil {
		return err
	}

	winls, err := n.GetWindowList()
	if err != nil {
		return err
	}

	if !opened {
		for _, bufname := range winls {
			if nvim.Buffer(bufname) == scratch {
				opened = true
				break
			}
		}
	}

	if !opened {
		n.SplitOpenBuffer(scratch)
	}

	return nil
}

// Input input in cmdline.
func (n *Nimvle) Input(ask string) (string, error) {
	var input string
	if err := n.v.Eval(`input("`+ask+`: ")`, &input); err != nil {
		return "", err
	}

	return input, nil
}

// Get Variable (return `interface{}`)
func (n *Nimvle) GetVar(varname string) (interface{}, error) {
	var v interface{}
	if err := n.v.Var(varname, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Eval expression
func (n *Nimvle) Eval(exp string) (interface{}, error) {
	var res interface{}
	if err := n.v.Eval(exp, &res); err != nil {
		return nil, err
	}

	return res, nil
}
