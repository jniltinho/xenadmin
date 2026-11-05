package rules

import (
	"strconv"
	"strings"

	"xcpng-center/internal/messages"
)

// CompareVersions is Helpers.ProductVersionCompare on dotted versions ("8.3.0" < "8.10"): -1, 0 or 1.
func CompareVersions(a, b string) int {
	x, y := strings.Split(a, "."), strings.Split(b, ".")
	for i := 0; i < len(x) || i < len(y); i++ {
		var p, q int
		if i < len(x) {
			p, _ = strconv.Atoi(x[i])
		}
		if i < len(y) {
			q, _ = strconv.Atoi(y[i])
		}
		if p != q {
			if p < q {
				return -1
			}
			return 1
		}
	}
	return 0
}

// TargetText is MigrateRules.MenuText: "xcp-lab-2 (Current server)" (MAINWINDOW_CONTEXT_REASON), or the name alone.
func TargetText(name, reason string) string {
	if reason == "" {
		return name
	}
	return messages.Format(messages.Get("MAINWINDOW_CONTEXT_REASON"), name, reason)
}
