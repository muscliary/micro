package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

const (
	RTColorscheme = 0
	RTSyntax      = 1
	RTHelp        = 2
	RTPlugin      = 3
	NumTypes      = 4 // How many filetypes are there
)

type RTFiletype byte

// RuntimeFile allows the program to read runtime data like colorschemes or syntax files
type RuntimeFile interface {
	// Name returns a name of the file without paths or extensions
	Name() string
	// Data returns the content of the file.
	Data() ([]byte, error)
}

// allFiles contains all available files, mapped by filetype
var allFiles [NumTypes][]RuntimeFile

// some file on filesystem
type realFile string

// some asset file
type assetFile string

// some file on filesystem but with a different name
type namedFile struct {
	realFile
	name string
}

// a file with the data stored in memory
type memoryFile struct {
	name string
	data []byte
}

func (mf memoryFile) Name() string {
	return mf.name
}
func (mf memoryFile) Data() ([]byte, error) {
	return mf.data, nil
}

func (rf realFile) Name() string {
	fn := filepath.Base(string(rf))
	return fn[:len(fn)-len(filepath.Ext(fn))]
}

func (rf realFile) Data() ([]byte, error) {
	return ioutil.ReadFile(string(rf))
}

func (af assetFile) Name() string {
	fn := path.Base(string(af))
	return fn[:len(fn)-len(path.Ext(fn))]
}

func (af assetFile) Data() ([]byte, error) {
	return Asset(string(af))
}

func (nf namedFile) Name() string {
	return nf.name
}

// AddRuntimeFile registers a file for the given filetype
func AddRuntimeFile(fileType RTFiletype, file RuntimeFile) {
	allFiles[fileType] = append(allFiles[fileType], file)
}

// AddRuntimeFilesFromDirectory registers each file from the given directory for
// the filetype which matches the file-pattern
func AddRuntimeFilesFromDirectory(fileType RTFiletype, directory, pattern string) {
	files, _ := ioutil.ReadDir(directory)
	for _, f := range files {
		if ok, _ := filepath.Match(pattern, f.Name()); !f.IsDir() && ok {
			fullPath := filepath.Join(directory, f.Name())
			AddRuntimeFile(fileType, realFile(fullPath))
		}
	}
}

// AddRuntimeFilesFromAssets registers each file from the given asset-directory for
// the filetype which matches the file-pattern
func AddRuntimeFilesFromAssets(fileType RTFiletype, directory, pattern string) {
	files, err := AssetDir(directory)
	if err != nil {
		return
	}
	for _, f := range files {
		if ok, _ := path.Match(pattern, f); ok {
			AddRuntimeFile(fileType, assetFile(path.Join(directory, f)))
		}
	}
}

// FindRuntimeFile finds a runtime file of the given filetype and name
// will return nil if no file was found
func FindRuntimeFile(fileType RTFiletype, name string) RuntimeFile {
	for _, f := range ListRuntimeFiles(fileType) {
		if f.Name() == name {
			return f
		}
	}
	return nil
}

// ListRuntimeFiles lists all known runtime files for the given filetype
func ListRuntimeFiles(fileType RTFiletype) []RuntimeFile {
	return allFiles[fileType]
}

// InitRuntimeFiles initializes all assets file and the config directory
func InitRuntimeFiles() {
	add := func(fileType RTFiletype, dir, pattern string) {
		AddRuntimeFilesFromDirectory(fileType, filepath.Join(configDir, dir), pattern)
		AddRuntimeFilesFromAssets(fileType, path.Join("runtime", dir), pattern)
	}

	add(RTColorscheme, "colorschemes", "*.micro")
	add(RTSyntax, "syntax", "*.yaml")
	add(RTHelp, "help", "*.md")

	// Search configDir for plugin-scripts
	files, _ := ioutil.ReadDir(filepath.Join(configDir, "plugins"))
	for _, f := range files {
		realpath, _ := filepath.EvalSymlinks(filepath.Join(configDir, "plugins", f.Name()))
		realpathStat, _ := os.Stat(realpath)
		if realpathStat.IsDir() {
			scriptPath := filepath.Join(configDir, "plugins", f.Name(), f.Name()+".lua")
			if _, err := os.Stat(scriptPath); err == nil {
				AddRuntimeFile(RTPlugin, realFile(scriptPath))
			}
		}
	}

	if files, err := AssetDir("runtime/plugins"); err == nil {
		for _, f := range files {
			scriptPath := path.Join("runtime/plugins", f, f+".lua")
			if _, err := AssetInfo(scriptPath); err == nil {
				AddRuntimeFile(RTPlugin, assetFile(scriptPath))
			}
		}
	}
}

// PluginReadRuntimeFile allows plugin scripts to read the content of a runtime file
func PluginReadRuntimeFile(fileType RTFiletype, name string) string {
	if file := FindRuntimeFile(fileType, name); file != nil {
		if data, err := file.Data(); err == nil {
			return string(data)
		}
	}
	return ""
}

// PluginListRuntimeFiles allows plugins to lists all runtime files of the given type
func PluginListRuntimeFiles(fileType RTFiletype) []string {
	files := ListRuntimeFiles(fileType)
	result := make([]string, len(files))
	for i, f := range files {
		result[i] = f.Name()
	}
	return result
}

// PluginAddRuntimeFile adds a file to the runtime files for a plugin
func PluginAddRuntimeFile(plugin string, filetype RTFiletype, filePath string) {
	fullpath := filepath.Join(configDir, "plugins", plugin, filePath)
	if _, err := os.Stat(fullpath); err == nil {
		AddRuntimeFile(filetype, realFile(fullpath))
	} else {
		fullpath = path.Join("runtime", "plugins", plugin, filePath)
		AddRuntimeFile(filetype, assetFile(fullpath))
	}
}

// PluginAddRuntimeFilesFromDirectory adds files from a directory to the runtime files for a plugin
func PluginAddRuntimeFilesFromDirectory(plugin string, filetype RTFiletype, directory, pattern string) {
	fullpath := filepath.Join(configDir, "plugins", plugin, directory)
	if _, err := os.Stat(fullpath); err == nil {
		AddRuntimeFilesFromDirectory(filetype, fullpath, pattern)
	} else {
		fullpath = path.Join("runtime", "plugins", plugin, directory)
		AddRuntimeFilesFromAssets(filetype, fullpath, pattern)
	}
}

// PluginAddRuntimeFileFromMemory adds a file to the runtime files for a plugin from a given string
func PluginAddRuntimeFileFromMemory(plugin string, filetype RTFiletype, filename, data string) {
	AddRuntimeFile(filetype, memoryFile{filename, []byte(data)})
}
