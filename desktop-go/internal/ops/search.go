package ops

import (
	"errors"

	xenapi "go/xenapi"

	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// SaveSearch is Search.Save: the search's XML in the (first) pool's gui_config under xensearch.Search-<uuid>, the old
// entry removed first so saving replaces it.
func SaveSearch(c *xapi.Connection, uuid, text string) error {
	if uuid == "" {
		return errors.New("the search has no uuid")
	}
	pools, err := xenapi.Pool.GetAll(c.Session)
	if err != nil {
		return xapi.Sanitize(err)
	}
	if len(pools) == 0 {
		return errors.New("The search could not be saved to the pool.")
	}
	key := rules.SearchKeyPrefix + uuid
	_ = xenapi.Pool.RemoveFromGuiConfig(c.Session, pools[0], key) // absent: nothing to remove
	return xapi.Sanitize(xenapi.Pool.AddToGuiConfig(c.Session, pools[0], key, text))
}

// DeleteSearch is Search.Delete on one connection: the search removed from each of its pools.
func DeleteSearch(c *xapi.Connection, uuid string) error {
	if uuid == "" {
		return errors.New("the search has no uuid")
	}
	pools, err := xenapi.Pool.GetAll(c.Session)
	if err != nil {
		return xapi.Sanitize(err)
	}
	for _, p := range pools {
		if err := xenapi.Pool.RemoveFromGuiConfig(c.Session, p, rules.SearchKeyPrefix+uuid); err != nil {
			return xapi.Sanitize(err)
		}
	}
	return nil
}
