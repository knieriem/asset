// Package asset provides access to assets that have been stored into
// the binary executable which has been turned into a
// self-extracting zip archive.
//
// If the current directory contains a sub directory `assets', files present
// there will override corresponding files in the zip archive.
package asset

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/knieriem/fsutil"
	"github.com/knieriem/osext"
)

const (
	assetDirName = "assets"
)

// A file system containing the asset files
var FS fs.FS

var ns = fsutil.NameSpace{}

var exeDir string

func init() {
	exe, err := osext.Executable()
	if err != nil {
		log.Fatal(err)
	}

	exeDir = filepath.Dir(exe)

	FS = &ns

	assetDir := filepath.Join(exeDir, assetDirName)
	fi, err := os.Stat(assetDir)
	if err == nil {
		if fi.IsDir() {
			log.Println("asset: found local directory")
			ns.Bind(".", os.DirFS(assetDir), withLabel("asset dir"))
			return
		}
	}

	zr, err := zip.OpenReader(exe)
	if err != nil {
		return
	}
	ns.Bind(".", zr, withLabel("builtin"))
}

// SetDefaultFS sets a file system to be used by global functions
// in case during init() a local asset directory could not be found and
// the executable does not have a zip file appended to it.
//
// This function allows maintaining the the package's previous behaviour
// while extending it to work with io/fs based file systems.
func SetDefaultFS(fsys fs.FS, root, source string) error {
	if len(ns.UnionFS) != 0 {
		return nil
	}
	sub, err := fs.Sub(fsys, root)
	if err != nil {
		return err
	}
	return ns.Bind(".", sub, withLabel(source))
}

func BindExeDir() {
	ns.Bind(".", os.DirFS(exeDir), withLabel(".exe dir"), fsutil.BindBefore())
}

func BindExeSubDir(name string) {
	ns.Bind(".", os.DirFS(filepath.Join(exeDir, name)), withLabel(".exe dir"), fsutil.BindBefore())
}

func BindBefore(dir string) {
	ns.Bind(".", os.DirFS(dir), fsutil.BindBefore())
}

func withLabel(val string) fsutil.BindOption {
	return fsutil.WithValue(fsutil.LabelKey, val)
}

// FileString reads an asset file and returns its contents as a string.
func FileString(name string) (content string, err error) {
	f, err := ns.Open(name)
	if err != nil {
		return
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	content = string(b)
	return
}

func Open(name string) (io.ReadSeekCloser, error) {
	f, err := ns.Open(name)
	if err != nil {
		return nil, err
	}
	seeker, ok := f.(io.Seeker)
	if ok {
		return &seekableFile{ReadCloser: f, Seeker: seeker}, nil
	}
	return nil, fmt.Errorf("file does not implement io.Seeker: %q", name)
}

type seekableFile struct {
	io.ReadCloser
	io.Seeker
}

//func HttpFS(root string) http.FileSystem {
//	ns := vfs.NameSpace{}
//	ns.Bind("/", FS, root, vfs.BindReplace)
//	return httpfs.New(&ns)
//}

func Stat(path string) (fi os.FileInfo, err error) {
	return fs.Stat(ns, path)
}

func ReadDir(path string) (fi []fs.DirEntry, err error) {
	return fs.ReadDir(ns, path)
}
