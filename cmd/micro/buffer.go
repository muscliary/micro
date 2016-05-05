package main

import (
	"crypto/md5"
	"github.com/vinzmay/go-rope"
	"io/ioutil"
	"strings"
)

// Buffer stores the text for files that are loaded into the text editor
// It uses a rope to efficiently store the string and contains some
// simple functions for saving and wrapper functions for modifying the rope
type Buffer struct {
	// Stores the text of the buffer
	r *rope.Rope

	// Path to the file on disk
	path string
	// Name of the buffer on the status line
	name string

	// This is the text stored every time the buffer is saved to check if the buffer is modified
	savedText           [16]byte
	netInsertions       int
	dirtySinceLastCheck bool

	// Provide efficient and easy access to text and lines so the rope String does not
	// need to be constantly recalculated
	// These variables are updated in the update() function
	lines    []string
	numLines int

	// Syntax highlighting rules
	rules []SyntaxRule
	// The buffer's filetype
	filetype string
}

// NewBuffer creates a new buffer from `txt` with path and name `path`
func NewBuffer(txt, path string) *Buffer {
	b := new(Buffer)
	if txt == "" {
		b.r = new(rope.Rope)
	} else {
		b.r = rope.New(txt)
	}
	b.path = path
	b.name = path
	b.savedText = md5.Sum([]byte(txt))

	b.Update()
	b.UpdateRules()

	return b
}

// UpdateRules updates the syntax rules and filetype for this buffer
// This is called when the colorscheme changes
func (b *Buffer) UpdateRules() {
	b.rules, b.filetype = GetRules(b)
}

func (b *Buffer) String() string {
	text := ""
	if b.r.Len() != 0 {
		text = b.r.String()
	}
	return text
}

// Update fetches the string from the rope and updates the `text` and `lines` in the buffer
func (b *Buffer) Update() {
	b.lines = strings.Split(b.String(), "\n")
	b.numLines = len(b.lines)
}

// Save saves the buffer to its default path
func (b *Buffer) Save() error {
	return b.SaveAs(b.path)
}

// SaveAs saves the buffer to a specified path (filename), creating the file if it does not exist
func (b *Buffer) SaveAs(filename string) error {
	b.UpdateRules()
	data := []byte(b.String())
	err := ioutil.WriteFile(filename, data, 0644)
	if err == nil {
		b.savedText = md5.Sum(data)
		b.netInsertions = 0
	}
	return err
}

// IsDirty returns whether or not the buffer has been modified compared to the one on disk
func (b *Buffer) IsDirty() bool {
	if !b.dirtySinceLastCheck {
		return false
	}
	if b.netInsertions == 0 {
		isDirty := b.savedText != md5.Sum([]byte(b.String()))
		b.dirtySinceLastCheck = isDirty
		return isDirty
	}
	return true
}

// Insert a string into the rope
func (b *Buffer) Insert(idx int, value string) {
	b.dirtySinceLastCheck = true
	b.netInsertions += len(value)
	b.r = b.r.Insert(idx, value)
	b.Update()
}

// Remove a slice of the rope from start to end (exclusive)
// Returns the string that was removed
func (b *Buffer) Remove(start, end int) string {
	b.dirtySinceLastCheck = true
	b.netInsertions -= end - start
	if start < 0 {
		start = 0
	}
	if end > b.Len() {
		end = b.Len()
	}
	removed := b.r.Report(start+1, end-start)
	// The rope implenentation I am using wants indicies starting at 1 instead of 0
	start++
	end++
	b.r = b.r.Delete(start, end-start)
	b.Update()
	return removed
}

// Len gives the length of the buffer
func (b *Buffer) Len() int {
	return b.r.Len()
}
