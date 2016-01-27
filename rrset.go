package udnssdk

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type RRSetsService struct {
	client *Client
}

//type stringProfile StringProfile
type Metaprofile struct {
	Context string `json:"@context"`
}

func (sp *StringProfile) UnmarshalJSON(b []byte) (err error) {
	sp.Profile = string(b)
	return nil
}
func (sp *StringProfile) MarshalJSON() ([]byte, error) {
	if sp.Profile != "" {
		return []byte(sp.Profile), nil
	} else {
		return json.Marshal(nil)
	}
}
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
func (sp *StringProfile) GoString() string {
	return sp.Profile
}
func (sp *StringProfile) String() string {
	return sp.Profile
}

type StringProfile struct {
	Profile string `json:"profile,omitempty"`
}

type RRSet struct {
	OwnerName string         `json:"ownerName"`
	RRType    string         `json:"rrtype"`
	TTL       int            `json:"ttl"`
	RData     []string       `json:"rdata"`
	Profile   *StringProfile `json:"profile,omitempty"`
}

type RRSetListDTO struct {
	ZoneName   string     `json:"zoneName"`
	Rrsets     []RRSet    `json:"rrsets"`
	Queryinfo  QueryInfo  `json:"queryInfo"`
	Resultinfo ResultInfo `json:"resultInfo"`
}
type rrsetWrapper struct {
	RRSet RRSet `json:"rrset"`
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

// List the zone rrsets.
//
func (s *RRSetsService) GetRRSets(zone string, rrsetName, rrsetType string) ([]RRSet, *Response, error) {
	//v := url.Values{}

	if rrsetType == "" {
		rrsetType = "ANY"
	}
	reqStr := rrsetPath(zone, rrsetType, rrsetName)
	var rrsld RRSetListDTO
	//wrappedRRSets := []RRSet{}
	rrsets := []RRSet{}
	offset := 0

	// TODO: Sane Configuration for timeouts / retries
	timeout := 5
	waittime := 5 * time.Second
	errcnt := 0
	for true {

		res, err := s.client.get(fmt.Sprintf("%s?offset=%d", reqStr, offset), &rrsld)
		if err != nil {
			if res.StatusCode >= 500 {
				errcnt = errcnt + 1
				if errcnt < timeout {
					time.Sleep(waittime)
					continue
				}
			}
			return rrsets, res, err

		}
		log.Printf("ResultInfo: %+v\n", rrsld.Resultinfo)
		for _, rrset := range rrsld.Rrsets {
			rrsets = append(rrsets, rrset)
		}
		if rrsld.Resultinfo.ReturnedCount+rrsld.Resultinfo.Offset >= rrsld.Resultinfo.TotalCount {
			return rrsets, res, nil
		} else {
			offset = rrsld.Resultinfo.ReturnedCount + rrsld.Resultinfo.Offset
			continue
		}
	}
	return rrsets, nil, nil
}

// CreateRRSet creates a zone rrset.
func (s *RRSetsService) CreateRRSet(zone string, rrsetAttributes RRSet) (*Response, error) {
	path := rrsetPath(zone, rrsetAttributes.RRType, rrsetAttributes.OwnerName)
	var retval interface{}
	res, err := s.client.post(path, rrsetAttributes, &retval)
	//log.Printf("CreateRRSet Retval: %+v", retval)
	return res, err
}

// UpdateRRSet updates a zone rrset.
func (s *RRSetsService) UpdateRRSet(zone string, rrsetAttributes RRSet) (*Response, error) {
	path := rrsetPath(zone, rrsetAttributes.RRType, rrsetAttributes.OwnerName)
	var retval interface{}

	res, err := s.client.put(path, rrsetAttributes, &retval)
	//log.Printf("UpdateRRSet Retval: %+v", retval)

	return res, err
}

// DeleteRRSet deletes a zone rrset.
//
func (s *RRSetsService) DeleteRRSet(zone string, rrsetAttributes RRSet) (*Response, error) {
	path := rrsetPath(zone, rrsetAttributes.RRType, rrsetAttributes.OwnerName)

	return s.client.delete(path, nil)
}

// UpdateIP updates the IP of specific A rrset.
//
// This is not part of the standard API. However,
// this is useful for Dynamic DNS (DDNS or DynDNS).
/*
func (rrset *RRSet) UpdateIP(client *Client, IP string) error {
  newdata := []string{IP}
  newRRSet := RRSet{RData: newdata, OwnerName: rrset.OwnerName}
	_, _, err := client.Zones.UpdateRRSet(rrset.ZoneId, rrset.Id, newRRSet)
	return err
}
*/
