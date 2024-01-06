package commands

import (
	"deduplicate/version"
	"os"
	"runtime"
	"text/tabwriter"
	"text/template"

	"time"

	"github.com/spf13/cobra"
)

type versionOut struct {
	Program   string
	Version   string
	BuildTime string
	GitCommit string
	GoVersion string
	Os        string
	Arch      string
}

var versionTemplate = `
{{ .Program}}:
 Version:	{{.Version}}
 Built:	{{.BuildTime}}
 Go Version:	{{.GoVersion}}
 Git Commit:	{{.GitCommit}}
 Os:	{{.Os}}
 Arch:	{{.Arch}}
 `

func NewCmdVersion() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "version",
		Short:  "Prints out version",
		Long:   "Prints out version information",
		Hidden: false,
		Run: func(cmd *cobra.Command, args []string) {
			verstionDetails := versionOut{
				Program:   "deduplicate",
				Version:   version.Version,
				BuildTime: reformatDate(version.BuildTime),
				GitCommit: version.GitCommit,
				GoVersion: runtime.Version(),
				Os:        runtime.GOOS,
				Arch:      runtime.GOARCH,
			}
			t := template.New("versionOut")
			t, _ = t.Parse(versionTemplate)

			tab := tabwriter.NewWriter(os.Stdout, 20, 1, 1, ' ', 0)

			t.Execute(tab, verstionDetails)
			tab.Write([]byte("\n"))
			tab.Flush()
		},
	}
	return cmd
}
func reformatDate(buildTime string) string {
	t, errTime := time.Parse(time.RFC3339, buildTime)
	if errTime == nil {
		return t.Format(time.ANSIC)
	}
	return buildTime
}
