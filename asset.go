// Package asset provides access to assets that have been stored into
// the binary executable which has been turned into a
// self-extracting zip archive.
//
// If the current directory contains a sub directory `assets', files present
// there will override corresponding files in the zip archive.
package asset

import (
	"archive/zip"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/kardianos/osext"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/zipfs"

	"golang.org/x/tools/godoc/vfs/httpfs"
	"net/http"

	"github.com/knieriem/vfsutil"
)

const (
	assetDirName = "assets"
)

// A file system containing the asset files
var FS vfs.FileSystem

var ns = vfs.NameSpace{}

var exeDir string

func init() {
	exe, err := osext.Executable()
	if err != nil {
		log.Fatal(err)
	}

	exeDir = filepath.Dir(exe)

	assetDir := filepath.Join(exeDir, assetDirName)
	fi, err := os.Stat(assetDir)
	if err == nil {
		if fi.IsDir() {
			ns.Bind("/", vfsutil.LabeledOS(assetDir, "asset dir"), "/", vfs.BindReplace)
			log.Println("asset: found local directory")
		}
	}

	FS = &ns

	zr, err := zip.OpenReader(exe)
	if err != nil {
		return
	}

	ns.Bind("/", vfsutil.LabeledFS(zipfs.New(zr, "-"), "builtin"), "/assets", vfs.BindAfter)
}

func BindExeDir() {
	ns.Bind("/", vfsutil.LabeledOS(exeDir, ".exe dir"), "/", vfs.BindBefore)
}

func BindExeSubDir(name string) {
	ns.Bind("/", vfsutil.LabeledOS(filepath.Join(exeDir, name), ".exe dir"), "/", vfs.BindBefore)
}

func BindBefore(dir string) {
	ns.Bind("/", vfs.OS(dir), "/", vfs.BindBefore)
}

// FileString reads an asset file an returns its contents as a string.
func FileString(name string) (content string, err error) {
	f, err := FS.Open(name)
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

func Open(name string) (vfs.ReadSeekCloser, error) {
	return FS.Open(name)
}

func HttpFS(root string) http.FileSystem {
	ns := vfs.NameSpace{}
	ns.Bind("/", FS, root, vfs.BindReplace)
	return httpfs.New(&ns)
}

func Stat(path string) (fi os.FileInfo, err error) {
	return FS.Stat(path)
}

func ReadDir(path string) (fi []os.FileInfo, err error) {
	return FS.ReadDir(path)
}
