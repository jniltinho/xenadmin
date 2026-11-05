package rules

import (
	"encoding/xml"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"

	"xcpng-center/internal/messages"
)

// DefaultISCSIPort is Util.DEFAULT_ISCSI_PORT.
const DefaultISCSIPort = 3260

// IQN is IScsiIqnInfo: a target IQN the scan found, at the address and port to use for it.
type IQN struct {
	Index   int    `json:"index"`
	IQN     string `json:"iqn"`
	Address string `json:"address"`
	Port    int    `json:"port"`
	Text    string `json:"text"` // the combo's "iqn... (address:port)"
}

// LUN is ISCSIInfo: a LUN of the target.
type LUN struct {
	ScsiID string `json:"scsi_id"`
	LunID  int    `json:"lun_id"`
	Vendor string `json:"vendor"`
	Serial string `json:"serial"`
	Size   int64  `json:"size"`
	Text   string `json:"text"` // the combo's "LUN 1: serial: 4.0 GB (vendor)"
}

// ParseIQNs is ISCSIPopulateIQNsAction on the XML of an SR_BACKEND_FAILURE_96: the IQNs sorted by name (CA-40066),
// "*" last; an address in brackets for IPv6 (multi-homing, CA-11607).
func ParseIQNs(doc string) ([]IQN, error) {
	// the targets' element names vary (TGT, target...): read every child of <iscsi-target-iqns>
	var raw struct {
		Lists []struct {
			Inner string `xml:",innerxml"`
		} `xml:"iscsi-target-iqns"`
	}
	if err := xml.Unmarshal([]byte("<r>"+noDeclaration(doc)+"</r>"), &raw); err != nil {
		return nil, err
	}
	var out []IQN
	for _, l := range raw.Lists {
		var targets struct {
			Items []struct {
				Fields []struct {
					XMLName xml.Name
					Text    string `xml:",chardata"`
				} `xml:",any"`
			} `xml:",any"`
		}
		if err := xml.Unmarshal([]byte("<t>"+l.Inner+"</t>"), &targets); err != nil {
			return nil, err
		}
		for _, item := range targets.Items {
			i := IQN{Index: -1, Port: DefaultISCSIPort}
			for _, f := range item.Fields {
				v := strings.TrimSpace(f.Text)
				switch strings.ToLower(f.XMLName.Local) {
				case "index":
					i.Index, _ = strconv.Atoi(v)
				case "ipaddress":
					i.Address = iscsiAddress(v)
				case "port":
					if p, err := strconv.Atoi(v); err == nil {
						i.Port = p
					}
				case "targetiqn":
					i.IQN = v
				}
			}
			i.Text = fmt.Sprintf("%s (%s:%d)", i.IQN, i.Address, i.Port)
			out = append(out, i)
		}
	}
	sort.SliceStable(out, func(a, b int) bool {
		if (out[a].IQN == "*") != (out[b].IQN == "*") {
			return out[b].IQN == "*"
		}
		return naturalLess(out[a].IQN, out[b].IQN)
	})
	return out, nil
}

// ParseLUNs is ISCSIPopulateLunsAction on the XML of an SR_BACKEND_FAILURE_107 or _87: the LUNs by LUN id.
func ParseLUNs(doc string) ([]LUN, error) {
	var raw struct {
		Targets []struct {
			LUNs []struct {
				Fields []struct {
					XMLName xml.Name
					Text    string `xml:",chardata"`
				} `xml:",any"`
			} `xml:",any"`
		} `xml:"iscsi-target"`
	}
	if err := xml.Unmarshal([]byte("<r>"+noDeclaration(doc)+"</r>"), &raw); err != nil {
		return nil, err
	}
	var out []LUN
	for _, t := range raw.Targets {
		for _, node := range t.LUNs {
			l := LUN{LunID: -1, Size: -1}
			for _, f := range node.Fields {
				v := strings.TrimSpace(f.Text)
				switch strings.ToLower(f.XMLName.Local) {
				case "lunid":
					l.LunID, _ = strconv.Atoi(v)
				case "vendor":
					l.Vendor = v
				case "serial":
					l.Serial = v
				case "size":
					if n, err := strconv.ParseInt(v, 10, 64); err == nil {
						l.Size = n
					}
				case "scsiid":
					l.ScsiID = v
				}
			}
			l.Text = LUNText(l)
			out = append(out, l)
		}
	}
	sort.SliceStable(out, func(a, b int) bool { return out[a].LunID < out[b].LunID })
	return out, nil
}

// noDeclaration drops the <?xml ...?> line, which may only start a document (the SM's XML is wrapped to read it).
func noDeclaration(doc string) string {
	doc = strings.TrimSpace(doc)
	if strings.HasPrefix(doc, "<?xml") {
		if i := strings.Index(doc, "?>"); i >= 0 {
			return doc[i+2:]
		}
	}
	return doc
}

// iscsiAddress is ParseIPAddress: an IP address as the SM wants it (IPv6 in brackets), anything else unchanged.
func iscsiAddress(a string) string {
	host := a
	if h, _, err := net.SplitHostPort(a); err == nil {
		host = h
	}
	ip := net.ParseIP(strings.Trim(host, "[]"))
	switch {
	case ip == nil:
		return a
	case ip.To4() == nil:
		return "[" + ip.String() + "]"
	}
	return ip.String()
}

// LUNText is SrRules.LunLabel: "LUN 1: serial: 4.0 GB (vendor)".
func LUNText(l LUN) string {
	s := "LUN"
	if l.LunID != -1 {
		s += " " + strconv.Itoa(l.LunID)
	}
	if l.Serial != "" {
		s += ": " + l.Serial
	}
	if l.Size >= 0 {
		s += ": " + DiskSizeText(l.Size)
	}
	if l.Vendor != "" {
		s += " (" + l.Vendor + ")"
	}
	return s
}

// DiskSizeText is Util.DiskSizeString: binary units, one decimal.
func DiskSizeText(bytes int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	v, i := float64(bytes), 0
	for v >= 1024 && i < len(units)-1 {
		v, i = v/1024, i+1
	}
	if i == 0 {
		return fmt.Sprintf("%d B", bytes)
	}
	return fmt.Sprintf("%.1f %s", v, units[i])
}

// ISCSIDeviceConfig is LVMoISCSI.GetDeviceConfig: the IQN's own address and port, the LUN by SCSI id (by serial when
// it has no LUN id), CHAP when given.
func ISCSIDeviceConfig(iqn IQN, lun LUN, chapUser, chapPassword string) map[string]string {
	d := map[string]string{"target": iqn.Address, "port": strconv.Itoa(iqn.Port), "targetIQN": iqn.IQN}
	if lun.LunID == -1 {
		d["LUNSerial"] = lun.Serial
	} else {
		d["SCSIid"] = lun.ScsiID
	}
	if chapUser != "" {
		d["chapuser"], d["chappassword"] = chapUser, chapPassword
	}
	return d
}

// ISCSIDescription is SrRules.IscsiDescription: "iSCSI SR [address (iqn; LUN ...)]".
func ISCSIDescription(iqn IQN, lun LUN) string {
	return messages.Format(messages.Get("NEWSR_ISCSI_DESCRIPTION"), iqn.Address, iqn.IQN, LUNText(lun))
}

// ISCSIError is SrRules.IscsiError for a scan's XAPI failure code ("" keeps the server's own text).
func ISCSIError(code string) string {
	switch code {
	case "SR_BACKEND_FAILURE_68":
		return messages.Get("LOGGING_IN_TO_THE_ISCSI_TARGET_FAILED")
	case "SR_BACKEND_FAILURE_140":
		return messages.Get("INVALID_HOST")
	case "SR_BACKEND_FAILURE_141":
		return messages.Get("SR_UNABLE_TO_CONNECT_TO_SCSI_TARGET")
	}
	return ""
}

// naturalLess is StringUtility.NaturalCompare: digit runs compare as numbers ("lun2" < "lun10").
func naturalLess(a, b string) bool {
	for a != "" && b != "" {
		ra, rb := []rune(a), []rune(b)
		if isASCIIDigit(ra[0]) && isASCIIDigit(rb[0]) { // digits() reads ASCII digits: the same test, or no progress
			na, ia := digits(a)
			nb, ib := digits(b)
			if na != nb {
				return na < nb
			}
			a, b = a[ia:], b[ib:]
			continue
		}
		if ra[0] != rb[0] {
			return ra[0] < rb[0]
		}
		a, b = string(ra[1:]), string(rb[1:])
	}
	return len(a) < len(b)
}

func isASCIIDigit(r rune) bool { return r >= '0' && r <= '9' }

func digits(s string) (int64, int) {
	i := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}
	n, _ := strconv.ParseInt(s[:i], 10, 64)
	return n, i
}
