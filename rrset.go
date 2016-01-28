package udnssdk

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

// RRSetsService provides access to RRSet resources
type RRSetsService struct {
	client *Client
}

// Here is the big 'Profile' mess that should get refactored to a more managable place

//type stringProfile StringProfile
type Metaprofile struct {
	Context     string `json:"@context"`
	realprofile interface{}
}

// UnmarshalJSON does what it says on the tin
func (sp *StringProfile) UnmarshalJSON(b []byte) (err error) {
	sp.Profile = string(b)
	return nil
}

// MarshalJSON does what it says on the tin
func (sp *StringProfile) MarshalJSON() ([]byte, error) {
	if sp.Profile != "" {
		return []byte(sp.Profile), nil
	}
	return json.Marshal(nil)
}

// GetType does unknown things
func (sp *StringProfile) GetType() string {
	if sp.Profile == "" {
		return ""
	}
	var mp Metaprofile
	err := json.Unmarshal([]byte(sp.Profile), &mp)
	if err != nil {
		log.Printf("Error getting profile type: %+v\n", err)
		return ""
	}
	return mp.Context
}

// GoString returns the StringProfile's Profile.
func (sp *StringProfile) GoString() string {
	return sp.Profile
}

// String returns the StringProfile's Profile.
func (sp *StringProfile) String() string {
	return sp.Profile
}

// StringProfile wraps a Profile string
type StringProfile struct {
	Profile string `json:"profile,omitempty"`
}
type RDPoolProfile struct {
	Context     string `json:"@context"`
	Order       string `json:"order"`
	Description string `json:"description"`
}

type GeoInfo struct {
	Name           string   `json:"name"`
	IsAccountLevel bool     `json:"isAccountLevel,omitempty"`
	Codes          []string `json:"codes"`
}
type IpInfo struct {
	Name           string      `json:"name"`
	IsAccountLevel bool        `json:"isAccountLevel,omitempty"`
	Ips            []IPAddrDTO `json:"ips"`
}
type DPRDataInfo struct {
	AllNonConfigured bool    `json:"allNonConfigured,omitempty"`
	IpInfo           IpInfo  `json:"ipInfo,omitempty"`
	GeoInfo          GeoInfo `json:"geoInfo,omitempty"`
}
type DirPoolProfile struct {
	Context         string        `json:"@context"`
	Description     string        `json:"description"`
	ConflictResolve string        `json:"conflictResolve,omitempty"`
	RDataInfo       []DPRDataInfo `json:"rdataInfo"`
	NoResponse      DPRDataInfo   `json:"noResponse"`
}
type SBRDataInfo struct {
	State         string `json:"state"`
	RunProbes     bool   `json:"runProbes,omitempty"`
	Priority      int    `json:"priority"`
	FailoverDelay int    `json:"failoverDelay,omitempty"`
	Threshold     int    `json:"threshold"`
	Weight        int    `json:"weight"`
}
type BackupRecord struct {
	RData         string `json:"rdata"`
	FailoverDelay int    `json:"failoverDelay,omitempty"`
}
type SBPoolProfile struct {
	Context       string         `json:"@context"`
	Description   string         `json:"description"`
	RunProbes     bool           `json:"runProbes,omitempty"`
	ActOnProbes   bool           `json:"actOnProbes,omitempty"`
	Order         string         `json:"order,omitempty"`
	MaxActive     int            `json:"maxActive,omitempty"`
	MaxServed     int            `json:"maxServed,omitempty"`
	RDataInfo     []SBRDataInfo  `json:"rdataInfo"`
	BackupRecords []BackupRecord `json:"backupRecords"`
}
type TCPoolProfile struct {
	Context      string        `json:"@context"`
	Description  string        `json:"description"`
	RunProbes    bool          `json:"runProbes,omitempty"`
	ActOnProbes  bool          `json:"actOnProbes,omitempty"`
	MaxToLB      int           `json:"maxToLB,omitempty"`
	RDataInfo    []SBRDataInfo `json:"rdataInfo"`
	BackupRecord BackupRecord  `json:"backupRecord"`
}

func (sp *StringProfile) GetProfileObject() interface{} {
	typ := sp.GetType()
	if typ == "" {
		return nil
	}
	tmp := strings.Split(typ, "/")
	x := tmp[len(tmp)-1]
	switch x {
	case "DirPool.jsonschema":
		var dpp DirPoolProfile
		err := json.Unmarshal([]byte(sp.Profile), &dpp)
		if err != nil {
			log.Printf("Could not Unmarshal the DirPoolProfile.\n")
			return nil
		}
		return dpp
	case "RDPool.jsonschema":
		var rdp RDPoolProfile
		err := json.Unmarshal([]byte(sp.Profile), &rdp)
		if err != nil {
			log.Printf("Could not Unmarshal the RDPoolProfile.\n")
			return nil
		}
		return rdp
	case "SBPool.jsonschema":
		var sbp SBPoolProfile
		err := json.Unmarshal([]byte(sp.Profile), &sbp)
		if err != nil {
			log.Printf("Could not Unmarshal the SBPoolProfile.\n")
			return nil
		}
		return sbp
	case "TCPool.jsonschema":
		var tcp TCPoolProfile
		err := json.Unmarshal([]byte(sp.Profile), &tcp)
		if err != nil {
			log.Printf("Could not Unmarshal the TCPoolProfile.\n")
			return nil
		}
		return tcp
	default:
		log.Printf("ERROR - Fall through on GetProfileObject - %s.\n", x)
		return fmt.Errorf("Fallthrough on GetProfileObject type %s\n", x)
	}

}

type RRSet struct {
	OwnerName string         `json:"ownerName"`
	RRType    string         `json:"rrtype"`
	TTL       int            `json:"ttl"`
	RData     []string       `json:"rdata"`
	Profile   *StringProfile `json:"profile,omitempty"`
}

// RRSetListDTO wraps a list of RRSet resources
type RRSetListDTO struct {
	ZoneName   string     `json:"zoneName"`
	Rrsets     []RRSet    `json:"rrsets"`
	Queryinfo  QueryInfo  `json:"queryInfo"`
	Resultinfo ResultInfo `json:"resultInfo"`
}

// rrsetPath generates the resource path for given rrset that belongs to a zone.
func rrsetPath(zone string, rrtype interface{}, rrset interface{}) string {
	path := fmt.Sprintf("zones/%s/rrsets", zone)
	if rrtype != nil {
		path += fmt.Sprintf("/%v", rrtype)
		if rrset != nil {
			path += fmt.Sprintf("/%v", rrset)
		}
	}
	return path
}

func rrsetQueryPath(zone, rrsetName, rrsetType string, offset int) string {
	if rrsetType == "" {
		rrsetType = "ANY"
	}
	reqStr := rrsetPath(zone, rrsetType, rrsetName)
	return fmt.Sprintf("%s?offset=%d", reqStr, offset)
}

// ListAllRRSets will list the zone rrsets, paginating through all available results
func (s *RRSetsService) ListAllRRSets(zone string, rrsetName, rrsetType string) ([]RRSet, error) {
	// TODO: Sane Configuration for timeouts / retries
	maxerrs := 5
	waittime := 5 * time.Second

	rrsets := []RRSet{}
	errcnt := 0
	offset := 0

	for {
		reqRrsets, ri, res, err := s.ListRRSets(zone, rrsetName, rrsetType, offset)
		if err != nil {
			if res.StatusCode >= 500 {
				errcnt = errcnt + 1
				if errcnt < maxerrs {
					time.Sleep(waittime)
					continue
				}
			}
			return rrsets, err
		}

		log.Printf("ResultInfo: %+v\n", ri)
		for _, rrset := range reqRrsets {
			rrsets = append(rrsets, rrset)
		}
		if ri.ReturnedCount+ri.Offset >= ri.TotalCount {
			return rrsets, nil
		}
		offset = ri.ReturnedCount + ri.Offset
		continue
	}
}

// ListRRSets requests zone rrsets by zone, rrsetName, rrsetType & optional offset
func (s *RRSetsService) ListRRSets(zone, rrsetName, rrsetType string, offset int) ([]RRSet, ResultInfo, *Response, error) {
	var rrsld RRSetListDTO

	uri := rrsetQueryPath(zone, rrsetName, rrsetType, offset)
	res, err := s.client.get(uri, &rrsld)

	rrsets := []RRSet{}
	for _, rrset := range rrsld.Rrsets {
		rrsets = append(rrsets, rrset)
	}
	return rrsets, rrsld.Resultinfo, res, err
}

// CreateRRSet creates a zone rrset.
func (s *RRSetsService) CreateRRSet(zone string, rrsetAttributes RRSet) (*Response, error) {
	uri := rrsetPath(zone, rrsetAttributes.RRType, rrsetAttributes.OwnerName)
	var ignored interface{}
	return s.client.post(uri, rrsetAttributes, &ignored)
}

// UpdateRRSet updates a zone rrset.
func (s *RRSetsService) UpdateRRSet(zone string, rrsetAttributes RRSet) (*Response, error) {
	uri := rrsetPath(zone, rrsetAttributes.RRType, rrsetAttributes.OwnerName)
	var ignored interface{}
	return s.client.put(uri, rrsetAttributes, &ignored)
}

// DeleteRRSet deletes a zone rrset.
func (s *RRSetsService) DeleteRRSet(zone string, rrsetAttributes RRSet) (*Response, error) {
	uri := rrsetPath(zone, rrsetAttributes.RRType, rrsetAttributes.OwnerName)
	return s.client.delete(uri, nil)
}
