package udnssdk

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type RRSetsService struct {
	client *Client
}

// Here is the big 'Profile' mess that should get refactored to a more managable place

//type stringProfile StringProfile
type Metaprofile struct {
	Context     string `json:"@context"`
	realprofile interface{}
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
