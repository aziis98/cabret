package operation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"

	"github.com/aziis98/cabret"
)

func init() {
	registerType("program", &Program{})
}

// Program is a [cabret.ItemOperation] that passes the incoming item as input to the given command.
// Options are the following:
//
//	command: <shell command>
//	io: <format> # optional, by default is "raw"
//
// The io format can be one of the following
//
//   - "raw" will just pass the item data to the program
//   - "json" will pass the whole item as JSON to the given program, this should be useful for
//     making external plugins compatible with cabret.
type Program struct {
	IOFormat     string
	ShellCommand string
}

type Format interface {
	Input(item cabret.Content) (stdin io.Reader, err error)
	Output(item cabret.Content, stdout io.Reader) (*cabret.Content, error)
}

var ioProgramFormats = map[string]Format{}

func init() {
	ioProgramFormats["json"] = JsonFormat{}
	ioProgramFormats["raw"] = RawFormat{}
}

type JsonFormat struct{}

func (JsonFormat) Input(item cabret.Content) (stdin io.Reader, err error) {
	buf := &bytes.Buffer{}

	if err := json.NewEncoder(buf).Encode(map[string]any{
		"type":     item.Type,
		"metadata": item.Metadata,
		"data":     item.Data,
	}); err != nil {
		return nil, err
	}

	return buf, nil
}

func (JsonFormat) Output(item cabret.Content, stdout io.Reader) (*cabret.Content, error) {
	var result struct {
		Type     string         `json:"type"`
		Metadata map[string]any `json:"metadata"`
		Data     string         `json:"data"`
	}

	if err := json.NewDecoder(stdout).Decode(&result); err != nil {
		return nil, err
	}

	return &cabret.Content{
		Type:     result.Type,
		Metadata: result.Metadata,
		Data:     []byte(result.Data),
	}, nil
}

type RawFormat struct{}

func (RawFormat) Input(item cabret.Content) (stdin io.Reader, err error) {
	return bytes.NewReader(item.Data), nil
}

func (RawFormat) Output(item cabret.Content, stdout io.Reader) (*cabret.Content, error) {
	data, err := io.ReadAll(stdout)
	if err != nil {
		return nil, err
	}

	item.Data = data
	return &item, nil
}

func (op *Program) Configure(options map[string]any) error {
	var err error
	op.IOFormat, err = getKey(options, "io", "raw")
	if err != nil {
		return err
	}
	op.ShellCommand, err = getKey[string](options, "command")
	if err != nil {
		return err
	}

	return nil
}

func (op *Program) ProcessItem(item cabret.Content) (*cabret.Content, error) {
	ioFmt, ok := ioProgramFormats[op.IOFormat]
	if !ok {
		return nil, fmt.Errorf(`unknown io format "%s"`, op.IOFormat)
	}

	r, err := ioFmt.Input(item)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("sh", "-c", op.ShellCommand)
	cmd.Stdin = r
	var buf bytes.Buffer
	cmd.Stdout = &buf

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return ioFmt.Output(item, &buf)
}
