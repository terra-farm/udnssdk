package udnssdk

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

/* Manages Probes */

type ProbesService struct {
	client *Client
}

type ProbeInfoDTO struct {
	Id         string           `json:"id"`
	PoolRecord string           `json:"poolRecord"`
	ProbeType  string           `json:"type"`
	Interval   string           `json:"interval"`
	Agents     []string         `json:"agents"`
	Threshold  int              `json:"threshold"`
	Details    *ProbeDetailsDTO `json:"details"`
}

type ProbeDetailsLimitDTO struct {
	Warning  int `json:"warning"`
	Critical int `json:"critical"`
	Fail     int `json:"fail"`
}

/* This Has To Be Magic! */
type ProbeDetailsDTO struct {
	data   []byte
	Detail interface{} `json:"detail,omitempty"`
	typ    string
}

func (s *ProbeDetailsDTO) Populate(typ string) (err error) {
	log.Printf("DEBUG - ProbeDetailsDTO.Populate(\"%s\")\n", typ)
	switch strings.ToUpper(typ) {
	case "HTTP":
		var pp HTTPProbeDetailsDTO
		err = json.Unmarshal(s.data, &pp)
		s.typ = typ
		s.Detail = pp
		return err
	case "PING":
		var pp PingProbeDetailsDTO
		err = json.Unmarshal(s.data, &pp)
		s.typ = typ
		s.Detail = pp
		return err
	case "FTP":
		var pp FTPProbeDetailsDTO
		err = json.Unmarshal(s.data, &pp)
		s.typ = typ
		s.Detail = pp
		return err
	case "TCP":
		var pp TCPProbeDetailsDTO
		err = json.Unmarshal(s.data, &pp)
		s.typ = typ
		s.Detail = pp
		return err
	case "SMTP":
		var pp SMTPProbeDetailsDTO
		err = json.Unmarshal(s.data, &pp)
		s.typ = typ
		s.Detail = pp
		return err
	case "SMTP_SEND":
		var pp SMTPSENDProbeDetailsDTO
		err = json.Unmarshal(s.data, &pp)
		s.typ = typ
		s.Detail = pp
		return err
	case "DNS":
		var pp DNSProbeDetailsDTO
		err = json.Unmarshal(s.data, &pp)
		s.typ = typ
		s.Detail = pp
		return err
	default:
		return fmt.Errorf("ERROR - ProbeDetailsDTO.Populate(\"%s\") - Fall through!\n", typ)
	}
}

func (s *ProbeDetailsDTO) UnmarshalJSON(b []byte) (err error) {
	s.data = b
	return nil
}
func (s *ProbeDetailsDTO) MarshalJSON() ([]byte, error) {
	var err error
	log.Printf("In marshal\n\n")
	if s.Detail != nil {
		d, e := json.Marshal(s.Detail)
		log.Printf("DEBUG - Detail Marshal: %+v Err: %+v\n", string(d), e)
		return d, e
	}
	if len(s.data) != 0 {
		return s.data, err
	} else {
		return json.Marshal(nil)
	}
}
func (s *ProbeDetailsDTO) GoString() string {
	return string(s.data)
}
func (s *ProbeDetailsDTO) String() string {
	return string(s.data)
}

type Transaction struct {
	Method          string                          `json:"method"`
	Url             string                          `json:"url"`
	TransmittedData string                          `json:"transmittedData,omitempty"`
	FollowRedirects bool                            `json:"followRedirects,omitempty"`
	Limits          map[string]ProbeDetailsLimitDTO `json:"limits"`
}
type HTTPProbeDetailsDTO struct {
	Transactions []Transaction         `json:"transactions"`
	TotalLimits  *ProbeDetailsLimitDTO `json:"totalLimits,omitempty"`
}
type PingProbeDetailsDTO struct {
	Packets    int                             `json:"packets,omitempty"`
	PacketSize int                             `json:"packetSize,omitempty"`
	Limits     map[string]ProbeDetailsLimitDTO `json:"limits"`
}
type FTPProbeDetailsDTO struct {
	Port        int                             `json:"port,omitempty"`
	PassiveMode bool                            `json:"passiveMode,omitempty"`
	Username    string                          `json:"username,omitempty"`
	Password    string                          `json:"password,omitempty"`
	Path        string                          `json:"path"`
	Limits      map[string]ProbeDetailsLimitDTO `json:"limits"`
}

type TCPProbeDetailsDTO struct {
	Port      int                             `json:"port,omitempty"`
	ControlIP string                          `json:"controlip,omitempty"`
	Limits    map[string]ProbeDetailsLimitDTO `json:"limits"`
}
type SMTPProbeDetailsDTO struct {
	Port   int                             `json:"port,omitempty"`
	Limits map[string]ProbeDetailsLimitDTO `json:"limits"`
}
type SMTPSENDProbeDetailsDTO struct {
	Port    int                             `json:"port,omitempty"`
	From    string                          `json:"from"`
	To      string                          `json:"from"`
	Message string                          `json:"from,omitempty"`
	Limits  map[string]ProbeDetailsLimitDTO `json:"limits"`
}
type DNSProbeDetailsDTO struct {
	Port       int                             `json:"port,omitempty"`
	TcpOnly    bool                            `json:"tcpOnly,omitempty"`
	RecordType string                          `json:"type,omitempty"`
	OwnerName  string                          `json:"ownerName,omitempty"`
	Limits     map[string]ProbeDetailsLimitDTO `json:"limits"`
}
type ProbeListDTO struct {
	Probes     []ProbeInfoDTO `json:"probes"`
	Queryinfo  QueryInfo      `json:"queryInfo"`
	Resultinfo ResultInfo     `json:"resultInfo"`
}

func ProbePath(zone, typ, name, guid string) string {
	if guid == "" {
		return fmt.Sprintf("zones/%s/rrsets/%s/%s/probes", zone, typ, name)
	} else {

		return fmt.Sprintf("zones/%s/rrsets/%s/%s/probes/%s", zone, typ, name, guid)
	}
}

// Get a Probe.
func (s *ProbesService) GetProbe(name, typ, zone, guid string) (ProbeInfoDTO, *Response, error) {
	reqStr := ProbePath(zone, typ, name, guid)
	var t ProbeInfoDTO
	res, err := s.client.get(reqStr, &t)
	return t, res, err
}

// Create a Probe
func (s *ProbesService) CreateProbe(name, typ, zone string, dp ProbeInfoDTO) (*Response, error) {
	reqStr := ProbePath(zone, typ, name, "")
	var retval interface{}
	res, err := s.client.post(reqStr, dp, &retval)

	return res, err
}

// Update
func (s *ProbesService) UpdateProbe(name, typ, zone, guid string, dp ProbeInfoDTO) (*Response, error) {
	reqStr := ProbePath(zone, typ, name, guid)
	var retval interface{}
	res, err := s.client.put(reqStr, dp, &retval)

	return res, err
}

// List
func (s *ProbesService) ListProbes(query, name, typ, zone string) ([]ProbeInfoDTO, *Response, error) {
	// TODO: Soooo... This function does not handle pagination of Probes....
	//v := url.Values{}

	reqStr := ProbePath(zone, typ, name, "")
	if query != "" {
		reqStr = fmt.Sprintf("%s?sort=NAME&query=%s", reqStr, query)
	}
	log.Printf("DEBUG - ListProbes: %s\n", reqStr)
	var tld ProbeListDTO
	//wrappedProbes := []Probe{}

	res, err := s.client.get(reqStr, &tld)
	dps := []ProbeInfoDTO{}

	if err == nil {
		for _, t := range tld.Probes {
			dps = append(dps, t)
		}
	}
	return dps, res, err
}

// DeleteProbe
//
func (s *ProbesService) DeleteProbe(name, typ, zone, guid string) (*Response, error) {
	path := ProbePath(zone, typ, name, guid)
	return s.client.delete(path, nil)
}
