package typst

import (
	"os"
	"os/exec"
)

// CompileHTML compiles typstInput to htmlOutput and blocks until done.
func CompileHTML(typstPath string, htmlPath string) error {
	c := exec.Command("typst", "compile", "--format", "html", typstPath, htmlPath)
	c.Env = append(os.Environ(), "TYPST_FEATURES=html")
	c.Stderr = os.Stderr
	return c.Run()
}
