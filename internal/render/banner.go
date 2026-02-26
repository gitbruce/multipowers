package render

import "fmt"

func Activated(workflow string) string {
	return fmt.Sprintf("🐙 CLAUDE OCTOPUS ACTIVATED - %s", workflow)
}
