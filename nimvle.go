package nimvle

import (
	"fmt"
	"strings"

	"github.com/neovim/go-client/nvim"
)

const (
	getWindowList = `vivid#getWindowList()`
)

type Nimvle struct {
	v          *nvim.Nvim
	pluginName string
}

func New(v *nvim.Nvim, name string) *Nimvle {
	return &Nimvle{
		v:          v,
		pluginName: name,
	}
}

func (n *Nimvle) Log(message interface{}) error {
	return n.v.Command("echom '" + pluginName + ": " + fmt.Sprintf("%v", message) + "'")
}

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

func (n *Nimvle) SetContentToBuffer(buf nvim.Buffer, content string) error {
	var byteContent [][]byte

	tmp := strings.Split(content, "\n")
	for _, c := range tmp {
		byteContent = append(byteContent, []byte(c))
	}

	return n.v.SetBufferLines(buf, 0, -1, true, byteContent)
}

func (n *Nimvle) GetWindowList() (map[string]int, error) {
	res := make(map[string]int)

	if err := n.v.Eval(getWindowList, &res); err != nil {
		return nil, err
	}

	return res, nil
}

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

func (n *Nimvle) NewScratchBuffer(bufferName string) (*nvim.Buffer, error) {
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

// ScratchBufferを別ウィンドウで開いていればいいが、開かれていない場合などの処理
func (n *Nimvle) showScratchBuffer(scratch nvim.Buffer, str fmt.Stringer) error {
	var opened bool
	var scratch *nvim.Buffer
	var err error

	n.SetContentToBuffer(scratch, str.String())

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

func (n *Nimvle) Input(ask string) (string, error) {
	var input string
	if err := n.v.Eval(`input("`+ask+`: ")`, &input); err != nil {
		return "", err
	}

	return input, nil
}
