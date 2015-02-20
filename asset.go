// Package asset provides access to assets that have been stored into
// the binary executable which has been turned into a
// self-extracting zip archive.
//
// If the current directory contains a sub directory `assets', files present
// there will override tcorresponding files in the zip archive.
package asset

import (
	"archive/zip"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"bitbucket.org/kardianos/osext"
	"code.google.com/p/go.tools/godoc/vfs"
	"code.google.com/p/go.tools/godoc/vfs/zipfs"

	"code.google.com/p/go.tools/godoc/vfs/httpfs"
	"net/http"
)

const (
	assetDirName = "assets"
)

// A file system containing the asset files
var FS vfs.FileSystem

func init() {
	exe, err := osext.Executable()
	if err != nil {
		log.Fatal(err)
	}

	ns := vfs.NameSpace{}

	exeDir := filepath.Dir(exe)

	assetDir := filepath.Join(exeDir, assetDirName)
	fi, err := os.Stat(assetDir)
	if err == nil {
		if fi.IsDir() {
			ns.Bind("/", vfs.OS(assetDir), "/", vfs.BindReplace)
			log.Println("asset: found local directory")
		}
	}
	ns.Bind("/", vfs.OS(exeDir), "/", vfs.BindBefore)

	FS = &ns

	zr, err := zip.OpenReader(exe)
	if err != nil {
		return
	}

	ns.Bind("/", zipfs.New(zr, "-"), "/assets", vfs.BindAfter)
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
