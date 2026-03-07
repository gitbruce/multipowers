package tracks

import (
	"os"
	"path/filepath"
)

func IsLinkedWorktreeCheckout(projectDir string) (bool, error) {
	gitPath := filepath.Join(projectDir, ".git")
	info, err := os.Stat(gitPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return !info.IsDir(), nil
}
