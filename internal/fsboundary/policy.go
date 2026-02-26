package fsboundary

import (
	"fmt"
	"path/filepath"
	"strings"
)

func inPath(child, parent string) bool {
	c, _ := filepath.Abs(child)
	p, _ := filepath.Abs(parent)
	return c == p || strings.HasPrefix(c, p+string(filepath.Separator))
}

func ValidateWritePath(targetPath, projectRoot string) error {
	if inPath(targetPath, projectRoot) {
		return nil
	}
	return fmt.Errorf("path outside project boundary: %s", targetPath)
}

func ValidateArtifactPath(targetPath, projectRoot string) error {
	allowed := filepath.Join(projectRoot, ".multipowers")
	if inPath(targetPath, allowed) {
		return nil
	}
	return fmt.Errorf("artifact must stay under .multipowers: %s", targetPath)
}
