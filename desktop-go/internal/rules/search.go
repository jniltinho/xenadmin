package rules

import "strings"

// SearchKeyPrefix is Search.SearchPrefix: BrandManager.ExtensionSearch ("xensearch") + ".Search-", the pool's
// gui_config key of a saved search before its uuid.
const SearchKeyPrefix = "xensearch.Search-"

// ValidateSearch is SearchRules.Validate (the Search editor): a name and at least one type of object; "" when valid.
func ValidateSearch(name string, types int) string {
	switch {
	case strings.TrimSpace(name) == "":
		return "Enter a name for the search."
	case types == 0:
		return "Select at least one type of object to search for."
	}
	return ""
}
