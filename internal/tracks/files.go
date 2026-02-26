package tracks

import "path/filepath"

func Dir(projectDir, id string) string {
	return filepath.Join(projectDir, ".multipowers", "tracks", id)
}
