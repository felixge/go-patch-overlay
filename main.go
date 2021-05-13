package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()
	o := Overlay{PatchDir: flag.Arg(0), OverlayDir: flag.Arg(1)}
	jsonPath, err := o.Generate()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", jsonPath)
	return nil
}

type Overlay struct {
	PatchDir   string
	OverlayDir string
	Goroot     string
}

type overlayJSON struct {
	Replace map[string]string
}

func (o Overlay) Generate() (string, error) {
	j := overlayJSON{Replace: map[string]string{}}
	if err := o.generate(&j); err != nil {
		return "", err
	}
	jsonPath := filepath.Join(o.OverlayDir, "overlay.json")
	jsonData, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(jsonPath, jsonData, 0644)
	return jsonPath, err
}

func (o *Overlay) generate(j *overlayJSON) error {
	if o.OverlayDir == "" {
		tmpDir, err := ioutil.TempDir("", "go-patch-overlay")
		if err != nil {
			return err
		}
		o.OverlayDir = tmpDir
	}

	var err error
	o.OverlayDir, err = filepath.Abs(o.OverlayDir)
	if err != nil {
		return err
	}

	if o.Goroot == "" {
		o.Goroot = runtime.GOROOT()
	}

	if err := os.RemoveAll(o.OverlayDir); err != nil {
		return err
	}
	patches, err := filepath.Glob(filepath.Join(o.PatchDir, "*.patch"))
	if err != nil {
		return err
	}
	sort.Strings(patches)

	for _, patch := range patches {
		if err := o.applyPatch(patch, j); err != nil {
			return err
		}
	}
	return nil
}

func (o Overlay) applyPatch(pathPath string, j *overlayJSON) error {
	patchData, err := ioutil.ReadFile(pathPath)
	if err != nil {
		return err
	}
	files, _, err := gitdiff.Parse(bytes.NewReader(patchData))
	if err != nil {
		return err
	}
	for _, file := range files {
		overlayPath := filepath.Join(o.OverlayDir, file.NewName)
		srcPath := filepath.Join(o.Goroot, file.OldName)
		if err := os.MkdirAll(filepath.Dir(overlayPath), 0755); err != nil {
			return err
		}
		if _, err := os.Stat(overlayPath); os.IsNotExist(err) {
			if err := copyFile(srcPath, overlayPath); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		beforeData, err := ioutil.ReadFile(overlayPath)
		if err != nil {
			return err
		}
		afterData := &bytes.Buffer{}
		if err := gitdiff.NewApplier(bytes.NewReader(beforeData)).ApplyFile(afterData, file); err != nil {
			return err
		}
		if err := ioutil.WriteFile(overlayPath, afterData.Bytes(), 0644); err != nil {
			return err
		}
		j.Replace[srcPath] = overlayPath
	}
	return nil
}

func copyFile(src, dst string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dst, input, 0644)
}
