package main

// MIT License

// Copyright (c) 2018 Nick Potts

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kingpin"
	"github.com/jhoonb/archivex"
)

//AudacityProject dfd
type AudacityProject struct {
	ProjectName  string
	ParentFolder string
	Aup          os.FileInfo
	DataDir      string
}

// Compress grabs the project and shovels it into
func (ap AudacityProject) Compress(outputdir string) error {
	zf := archivex.TarFile{}
	a := filepath.Join(outputdir, ap.ProjectName+".tar.gz")

	if e := zf.Create(a); e != nil {
		return e
	}
	defer zf.Close()

	aupfile, err := os.Open(filepath.Join(ap.ParentFolder, ap.Aup.Name()))
	if err != nil {
		return err
	}
	defer aupfile.Close()
	if err := zf.Add(ap.Aup.Name(), aupfile, ap.Aup); err != nil {
		return err
	}

	if err := zf.AddAll(filepath.Join(ap.ParentFolder, ap.DataDir), true); err != nil {
		return err
	}

	return zf.Close()
}

//Locate finds all the AudacityProjects
func Locate(root string) []AudacityProject {
	r := []AudacityProject{}

	walkfxn := func(parent string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		name := info.Name()
		if info.Mode().IsRegular() && filepath.Ext(name) == ".aup" {
			projname := name[0 : len(name)-4]
			parentfolder, _ := filepath.Split(parent)
			datadir := projname + "_data"

			aupdir := path.Join(parentfolder, datadir) // this is audacity's data dir
			if fi, err := os.Stat(aupdir); err == nil && fi.IsDir() {
				r = append(r, AudacityProject{ProjectName: projname, ParentFolder: parentfolder, Aup: info, DataDir: datadir})
			}
			return nil
		}
		if info.IsDir() && strings.HasSuffix(info.Name(), "_data") {
			return filepath.SkipDir
		}
		return nil
	}
	filepath.Walk(root, walkfxn)
	return r
}

var (
	app    = kingpin.New("audacity-compresser", "A stupid tool to locate and compress audacity project data")
	root   = app.Arg("root", "Root Path to parse").Required().String()
	output = app.Flag("output", "Output dir.  Empty means place in the same directory as the project").Short('o').Default("").String()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	list := Locate(*root)
	for _, l := range list {
		fmt.Printf("Packing Audacity Project %q named at %q\n", l.ProjectName, l.ParentFolder)
		outdir := *output
		if outdir == "" {
			outdir = l.ParentFolder
		}
		if err := l.Compress(outdir); err != nil {
			fmt.Println("Unable to compress: ", err)
		}
	}
}
