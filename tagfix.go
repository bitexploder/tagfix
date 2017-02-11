package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mikkyang/id3-go"
	"regexp"
)

type Options struct {
	Mp3Dir    string
	IdxPat    string
	NamePat   string
	TitleTemp string
	Quiet     bool
	Trial     bool
}

func expandHome(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir + "/"

	if len(path) >= 2 && path[:2] == "~/" {
		path = strings.Replace(path, "~/", dir, 1)
	}

	return path
}

func extract1(pattern string, target string) (string, error) {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(target)

	//fmt.Println("matches: ", matches, " ", len(matches))
	if len(matches) < 2 {
		return "", errors.New(fmt.Sprintf("too few matches (%+v)", matches))
	}
	if len(matches) > 2 {
		return "", errors.New(fmt.Sprintf("too many matches (%+v)", matches))
	}
	return matches[1], nil
}

func flags() Options {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r, "\n")
			flag.PrintDefaults()
			os.Exit(1)
		}
	}()

	var o Options

	var mp3dir = flag.String("mp3dir", "", "directory of mp3 files")
	var idxpat = flag.String("idxpat", "([0-9]+)", "regex pattern to extract file index string")
	var namepat = flag.String("namepat", "- ([A-Za-z0-9 ]+) Ch.*", "regex pattern to extract name of item")
	var titletemp = flag.String("titletemp", "{{ .Index }} - {{ .Name }}", "a golang text template for the id3 title to be set to")
	var quiet = flag.Bool("quiet", false, "don't make a lot of noise")
	var trial = flag.Bool("trial", false, "try things out without modifying the mp3 file")

	flag.Parse()

	if *mp3dir == "" {
		panic("mp3 directory required")
	}
	o.Mp3Dir = expandHome(*mp3dir)
	o.IdxPat = *idxpat
	o.NamePat = *namepat
	o.TitleTemp = *titletemp
	o.Quiet = *quiet
	o.Trial = *trial
	return o
}

func main() {
	opts := flags()
	//fmt.Printf("Options: %+v\n", opts)
	files, _ := filepath.Glob(opts.Mp3Dir + "/*.mp3")

	for _, fnameFull := range files {
		var mp3f *id3.File
		var err error

		fname := filepath.Base(fnameFull)
		if !opts.Trial {
			mp3f, err = id3.Open(fnameFull)
			if err != nil {
				panic(err)
			}
		}

		title := make(map[string]string)
		if idx, err := extract1(opts.IdxPat, fname); err == nil {
			title["Index"] = idx
		} else {
			fmt.Printf("idx pat err: %s. Skipping '%s'\n", err.Error(), fname)
			continue
		}

		if n, err := extract1(opts.NamePat, fname); err == nil {
			title["Name"] = n
		} else {
			fmt.Printf("title pat err: %s. Skipping '%s'\n", err.Error(), fname)
			continue
		}

		tmpl, err := template.New("title").Parse(opts.TitleTemp)
		if err != nil {
			panic(err)
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, title)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Title to set: %s for '%s'\n", buf.String(), fname)
		if !opts.Trial {
			mp3f.SetTitle(buf.String())
			mp3f.Close()
		}

	}

	// fmt.Printf("Files: %+v\n", files)
}
