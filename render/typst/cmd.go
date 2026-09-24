package typst

import (
	"os/exec"
)

// CompileHTML compiles typstInput to htmlOutput and blocks until done.
func CompileHTML(typstPath string, htmlPath string) error {
	c := exec.Command("typst", "compile", "--format", "html", typstPath, htmlPath)
	c.Env = append(c.Env, "TYPST_FEATURES=html")
	return c.Run()
}
