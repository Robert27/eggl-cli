package kill

import "strings"

func trimOutput(data []byte) string {
	return strings.TrimSpace(string(data))
}
