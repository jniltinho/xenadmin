package app

import (
	xenapi "go/xenapi"

	"xcpng-center/internal/rules"
	"xcpng-center/internal/xapi"
)

// NewSRKind is a type the New Storage wizard offers, with its default name in this pool (SrRules.DefaultName).
type NewSRKind struct {
	Kind        rules.SRKind `json:"kind"`
	Name        string       `json:"name"`
	Blurb       string       `json:"blurb"`
	DefaultName string       `json:"default_name"`
	Iso         bool         `json:"iso"`
}

// ReadNewSRKinds is the wizard's Type page: NFS and SMB storage, then the NFS and SMB ISO libraries.
func ReadNewSRKinds(c *xapi.Connection) ([]NewSRKind, error) {
	srs, err := xenapi.SR.GetAllRecords(c.Session)
	if err != nil {
		return nil, xapi.Sanitize(err)
	}
	taken := map[string]bool{}
	for _, sr := range srs {
		taken[sr.NameLabel] = true
	}
	var out []NewSRKind
	for _, k := range rules.SRKinds {
		out = append(out, NewSRKind{Kind: k, Name: rules.SRTypeName(k), Blurb: rules.SRBlurb(k), DefaultName: rules.SRDefaultName(k, taken), Iso: rules.SRIsIso(k)})
	}
	return out, nil
}
