package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"bytes"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/bountylabs/log"
)

var path string
var repo string
var pkg string
var tsformat string
var versionstr string
var short bool

func init() {
	flag.StringVar(&path, "o", "version.go", "filename")
	flag.StringVar(&repo, "i", ".", "repository path")
	flag.StringVar(&pkg, "p", "version", "package")
	flag.StringVar(&tsformat, "tf", "", "timestamp format")
	flag.StringVar(&versionstr, "v", "", "version string")
	flag.BoolVar(&short, "s", false, "--s")
}

func stripchars(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}

// looks for git root in gitDir or above
func findgitroot(gitDir string) (string, error) {
	gitDir = filepath.Clean(gitDir)
	gitpath := gitDir + "/.git"

	for {
		_, err := os.Stat(gitpath)
		if err == nil {
			return gitpath, nil // found
		}

		dir := filepath.Dir(gitpath)
		// at filesystem root...not found
		if strings.HasSuffix(dir, "/") || strings.HasSuffix(dir, "\\") {
			return "", fmt.Errorf("No gitroot found at or above path %s", gitDir)
		}

		if os.IsNotExist(err) {
			gitpath = filepath.Clean(dir + "/../.git")
			continue
		}
	}
}

func main() {

	flag.Parse()

	gitrepo, err := findgitroot(repo)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
		return
	}

	//get commit hash (short or not)
	cmd := func() *exec.Cmd {
		if short {
			return exec.Command("git", "--git-dir", gitrepo, "rev-parse", "--short", "HEAD")
		}

		return exec.Command("git", "--git-dir", gitrepo, "rev-parse", "HEAD")
	}()

	cmdOut, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
		return
	}

	// generate timestamp
	tsUnix := strconv.FormatInt(time.Now().Unix(), 10)
	var tsFmt string
	if tsformat != "" {
		tsFmt = time.Now().Format(tsformat)
	}

	// map of values for template
	vals := make(map[string]string)
	vals["githash"] = stripchars(string(cmdOut), "\r\n ")
	vals["ts"] = tsUnix
	vals["tsformat"] = tsFmt
	vals["pkg"] = pkg
	vals["versionstr"] = versionstr

	// create template
	tmpl, err := template.New("ver").Parse(Template)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
		return
	}

	// execute template
	var b bytes.Buffer
	if err = tmpl.Execute(&b, vals); err != nil {
		log.Errorln(err)
		os.Exit(1)
		return
	}

	//write file
	if err = ioutil.WriteFile(path, b.Bytes(), 0644); err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
}

var Template = `package {{.pkg}}

const (
	GIT_COMMIT_HASH = "{{.githash}}"

	// Unix time (seconds since January 1, 1970 UTC)
	GENERATED = {{.ts}}
{{if .tsformat}}
	// human readable timestamp
	GENERATED_FMT = "{{.tsformat}}"{{end}}
{{if .versionstr}}
	VERSION = "{{.versionstr}}"{{end}}
)
`
