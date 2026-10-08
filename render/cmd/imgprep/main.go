// Command imgprep ensures every image in a directory has a PNG, WebP, AVIF,
// and JXL sibling, generating only whichever of those four are missing.
//
// Images are grouped by their path without extension (e.g. "2026/foo.jpg"
// and "2026/foo.webp" belong to the same group "2026/foo"). For each group,
// the highest-fidelity file already present (preferring PNG, then
// JPEG/JPG, then WebP, then AVIF, then JXL) is used as the source for
// anything missing. WebP, AVIF, and JXL sources are decoded to a temporary
// PNG first, since the encoders used here don't read each other's formats
// directly.
package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"image/png"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const (
	extPNG  = ".png"
	extWebP = ".webp"
	extAVIF = ".avif"
	extJXL  = ".jxl"
)

// group holds the path (if any) of each format already present for one
// image, keyed by its path without extension.
type group struct {
	png, jpg, webp, avif, jxl string
}

func main() {
	dir := flag.String("dir", "./static", "directory to scan for images (recursive)")
	flag.Parse()

	if err := run(*dir); err != nil {
		fmt.Fprintf(os.Stderr, "imgprep: %v\n", err)
		os.Exit(1)
	}
}

func run(dir string) error {
	groups := map[string]*group{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		ext := strings.ToLower(filepath.Ext(path))
		key := strings.TrimSuffix(path, filepath.Ext(path))
		g, ok := groups[key]
		if !ok {
			g = &group{}
		}
		switch ext {
		case extPNG:
			g.png = path
		case ".jpg", ".jpeg":
			g.jpg = path
		case extWebP:
			g.webp = path
		case extAVIF:
			g.avif = path
		case extJXL:
			g.jxl = path
		default:
			return nil
		}
		groups[key] = g
		return nil
	})
	if err != nil {
		return err
	}

	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		if err := processGroup(key, groups[key]); err != nil {
			return fmt.Errorf("%s: %w", key, err)
		}
	}
	return nil
}

// processGroup fills in whichever of png/webp/avif/jxl are missing from g,
// deriving them from the best already-present format.
func processGroup(key string, g *group) error {
	if g.png != "" && g.webp != "" && g.avif != "" && g.jxl != "" {
		return nil
	}

	direct, origSource, isTemp, err := directSource(key, g)
	if err != nil {
		return err
	}
	tempConsumed := false
	if isTemp {
		defer func() {
			if !tempConsumed {
				os.Remove(direct)
			}
		}()
	}

	if g.png == "" {
		dst := key + extPNG
		if err := refuseIfExists(dst); err != nil {
			return err
		}
		if isTemp {
			if err := os.Rename(direct, dst); err != nil {
				return fmt.Errorf("derive %s: %w", dst, err)
			}
			tempConsumed = true
			direct = dst
			isTemp = false
		} else if err := jpegToPNG(direct, dst); err != nil {
			return fmt.Errorf("derive %s: %w", dst, err)
		}
		fmt.Printf("generated %s (from %s)\n", dst, origSource)
		g.png = dst
	}

	if g.webp == "" {
		dst := key + extWebP
		if err := refuseIfExists(dst); err != nil {
			return err
		}
		if err := runTool("cwebp", "-quiet", direct, "-o", dst); err != nil {
			return fmt.Errorf("derive %s: %w", dst, err)
		}
		fmt.Printf("generated %s (from %s)\n", dst, direct)
		g.webp = dst
	}

	if g.avif == "" {
		dst := key + extAVIF
		if err := refuseIfExists(dst); err != nil {
			return err
		}
		if err := runTool("avifenc", direct, dst); err != nil {
			return fmt.Errorf("derive %s: %w", dst, err)
		}
		fmt.Printf("generated %s (from %s)\n", dst, direct)
		g.avif = dst
	}

	if g.jxl == "" {
		dst := key + extJXL
		if err := refuseIfExists(dst); err != nil {
			return err
		}
		if err := runTool("cjxl", direct, dst); err != nil {
			return fmt.Errorf("derive %s: %w", dst, err)
		}
		fmt.Printf("generated %s (from %s)\n", dst, direct)
		g.jxl = dst
	}

	return nil
}

// refuseIfExists errors out if path already exists, so a logic bug or race
// can never cause a generated file to silently clobber something already in
// the repo (the external encoders and os.Rename otherwise overwrite their
// destination without complaint).
func refuseIfExists(path string) error {
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("refusing to overwrite existing file %s", path)
	} else if !os.IsNotExist(err) {
		return err
	}
	return nil
}

// directSource returns a path that avifenc/cwebp/cjxl can read natively (a
// PNG or JPEG), plus origSource, the actual file that content came from (for
// logging). If the only format present is WebP, AVIF, or JXL, it's decoded
// to a temporary PNG first, and isTemp reports that the caller owns cleaning
// it up (or consuming it as the group's missing PNG).
func directSource(key string, g *group) (path, origSource string, isTemp bool, err error) {
	switch {
	case g.png != "":
		return g.png, g.png, false, nil
	case g.jpg != "":
		return g.jpg, g.jpg, false, nil
	case g.webp != "":
		tmp := key + ".imgprep-tmp.png"
		if err := runTool("dwebp", g.webp, "-o", tmp); err != nil {
			return "", "", false, err
		}
		return tmp, g.webp, true, nil
	case g.avif != "":
		tmp := key + ".imgprep-tmp.png"
		if err := runTool("avifdec", g.avif, tmp); err != nil {
			return "", "", false, err
		}
		return tmp, g.avif, true, nil
	case g.jxl != "":
		tmp := key + ".imgprep-tmp.png"
		if err := runTool("djxl", g.jxl, tmp); err != nil {
			return "", "", false, err
		}
		return tmp, g.jxl, true, nil
	default:
		return "", "", false, fmt.Errorf("no source image found")
	}
}

func jpegToPNG(src, dst string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	img, err := jpeg.Decode(f)
	if err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	defer out.Close()

	return png.Encode(out, img)
}

func runTool(name string, args ...string) error {
	c := exec.Command(name, args...)
	out, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %w: %s", name, err, strings.TrimSpace(string(out)))
	}
	return nil
}
