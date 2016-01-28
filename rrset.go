package udnssdk

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// RRSetsService provides access to RRSet resources
type RRSetsService struct {
	client *Client
}

// Metaprofile does unknown things
type Metaprofile struct {
	Context string `json:"@context"`
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

// RRSet wraps an RRSet resource
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
