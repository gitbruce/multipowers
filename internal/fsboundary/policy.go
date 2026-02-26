package fsboundary

import (
	"fmt"
	"path/filepath"
	"strings"
)

func ValidateWritePath(targetPath, projectRoot string) error {
	ap, _ := filepath.Abs(targetPath)
	pr, _ := filepath.Abs(projectRoot)
	if strings.HasPrefix(ap, pr) {
		return nil
	}
	return fmt.Errorf("path outside project boundary: %s", ap)
}
