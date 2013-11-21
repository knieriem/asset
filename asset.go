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

	"bitbucket.org/kardianos/osext"
	"code.google.com/p/go.tools/godoc/vfs"
	"code.google.com/p/go.tools/godoc/vfs/zipfs"
)

const (
	assetDirName = "assets"
)

// A file system containing the asset files
var FileSystem vfs.FileSystem

func init() {
	exe, err := osext.Executable()
	if err != nil {
		log.Fatal(err)
	}

	ns := vfs.NameSpace{}

	fi, err := os.Stat(assetDirName)
	if err == nil {
		if fi.IsDir() {
			ns.Bind("/", vfs.OS(assetDirName), "/", vfs.BindReplace)
			log.Println("asset: found local directory")
		}
	}

	FileSystem = &ns

	zr, err := zip.OpenReader(exe)
	if err != nil {
		log.Println("asset: could not open embedded zip file")
		return
	}

	ns.Bind("/", zipfs.New(zr, "-"), "/assets", vfs.BindAfter)
}

// FileString reads an asset file an returns its contents as a string.
func FileString(name string) (content string, err error) {
	f, err := FileSystem.Open(name)
	if err != nil {
		return
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	content = string(b)
	return
}
