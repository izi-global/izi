package version

import (
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"text/template"

	"time"

	iziLogger "github.com/izi-global/izi/logger"
)

// RuntimeInfo holds information about the current runtime.
type RuntimeInfo struct {
	GoVersion    string
	GOOS         string
	GOARCH       string
	NumCPU       int
	GOPATH       string
	GOROOT       string
	Compiler     string
	IZIVersion   string
	IZIGoVersion string
}

// InitBanner loads the banner and prints it to output
// All errors are ignored, the application will not
// print the banner in case of error.
func InitBanner(out io.Writer, in io.Reader) {
	if in == nil {
		iziLogger.Log.Fatal("The input is nil")
	}

	banner, err := ioutil.ReadAll(in)
	if err != nil {
		iziLogger.Log.Fatalf("Error while trying to read the banner: %s", err)
	}

	show(out, string(banner))
}

func show(out io.Writer, content string) {
	t, err := template.New("banner").
		Funcs(template.FuncMap{"Now": Now}).
		Parse(content)

	if err != nil {
		iziLogger.Log.Fatalf("Cannot parse the banner template: %s", err)
	}

	err = t.Execute(out, RuntimeInfo{
		GetGoVersion(),
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumCPU(),
		os.Getenv("GOPATH"),
		runtime.GOROOT(),
		runtime.Compiler,
		version,
		GetIZIGoVersion(),
	})
	if err != nil {
		iziLogger.Log.Error(err.Error())
	}
}

// Now returns the current local time in the specified layout
func Now(layout string) string {
	return time.Now().Format(layout)
}
