package udnssdk

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

// SBTCService manages Probes
type SBTCService struct {
	client *Client
}

// ProbeAlertDataDTO wraps a probe alert response
type ProbeAlertDataDTO struct {
	PoolRecord      string    `json:"poolRecord"`
	ProbeType       string    `json:"probeType"`
	ProbeStatus     string    `json:"probeStatus"`
	AlertDate       time.Time `json:"alertDate"`
	FailoverOccured bool      `json:"failoverOccured"`
	OwnerName       string    `json:"ownerName"`
	Status          string    `json:"status"`
}

// ProbeAlertDataListDTO wraps the response for an index of probe alerts
type ProbeAlertDataListDTO struct {
	Alerts     []ProbeAlertDataDTO `json:"alerts"`
	Queryinfo  QueryInfo           `json:"queryInfo"`
	Resultinfo ResultInfo          `json:"resultInfo"`
}

// ProbeInfoDTO wraps a probe response
type ProbeInfoDTO struct {
	ID         string           `json:"id"`
	PoolRecord string           `json:"poolRecord"`
	ProbeType  string           `json:"type"`
	Interval   string           `json:"interval"`
	Agents     []string         `json:"agents"`
	Threshold  int              `json:"threshold"`
	Details    *ProbeDetailsDTO `json:"details"`
}

// ProbeDetailsLimitDTO wraps a probe
type ProbeDetailsLimitDTO struct {
	Warning  int `json:"warning"`
	Critical int `json:"critical"`
	Fail     int `json:"fail"`
}

// ProbeDetailsDTO wraps the details of a probe
type ProbeDetailsDTO struct {
	data   []byte
	Detail interface{} `json:"detail,omitempty"`
	typ    string
}

// GetData returns the data because I'm working around something.
func (s *ProbeDetailsDTO) GetData() []byte {
	return s.data
}

// Populate does magical things with json unmarshalling to unroll the Probe into
// an appropriate datatype.  These are helper structures and functions for testing
// and direct API use.  In the Terraform implementation, we will use Terraforms own
// warped schema structure to handle the marshalling and unmarshalling.
func (s *ProbeDetailsDTO) Populate(typ string) (err error) {
	// TODO: actually document
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

// UnmarshalJSON does what it says on the tin
func (s *ProbeDetailsDTO) UnmarshalJSON(b []byte) (err error) {
	s.data = b
	return nil
}

// MarshalJSON does what it says on the tin
func (s *ProbeDetailsDTO) MarshalJSON() ([]byte, error) {
	var err error
	if s.Detail != nil {
		return json.Marshal(s.Detail)
	}
	if len(s.data) != 0 {
		return s.data, err
	}
	return json.Marshal(nil)
}

// GoString returns a string representation of the ProbeDetailsDTO internal data
func (s *ProbeDetailsDTO) GoString() string {
	return string(s.data)
}
func (s *ProbeDetailsDTO) String() string {
	return string(s.data)
}

// Transaction wraps a transaction response
type Transaction struct {
	Method          string                          `json:"method"`
	URL             string                          `json:"url"`
	TransmittedData string                          `json:"transmittedData,omitempty"`
	FollowRedirects bool                            `json:"followRedirects,omitempty"`
	Limits          map[string]ProbeDetailsLimitDTO `json:"limits"`
}

// HTTPProbeDetailsDTO wraps HTTP probe details
type HTTPProbeDetailsDTO struct {
	Transactions []Transaction         `json:"transactions"`
	TotalLimits  *ProbeDetailsLimitDTO `json:"totalLimits,omitempty"`
}

// PingProbeDetailsDTO wraps Ping probe details
type PingProbeDetailsDTO struct {
	Packets    int                             `json:"packets,omitempty"`
	PacketSize int                             `json:"packetSize,omitempty"`
	Limits     map[string]ProbeDetailsLimitDTO `json:"limits"`
}

// FTPProbeDetailsDTO wraps FTP probe details
type FTPProbeDetailsDTO struct {
	Port        int                             `json:"port,omitempty"`
	PassiveMode bool                            `json:"passiveMode,omitempty"`
	Username    string                          `json:"username,omitempty"`
	Password    string                          `json:"password,omitempty"`
	Path        string                          `json:"path"`
	Limits      map[string]ProbeDetailsLimitDTO `json:"limits"`
}

// TCPProbeDetailsDTO wraps TCP probe details
type TCPProbeDetailsDTO struct {
	Port      int                             `json:"port,omitempty"`
	ControlIP string                          `json:"controlIP,omitempty"`
	Limits    map[string]ProbeDetailsLimitDTO `json:"limits"`
}

// SMTPProbeDetailsDTO wraps SMTP probe details
type SMTPProbeDetailsDTO struct {
	Port   int                             `json:"port,omitempty"`
	Limits map[string]ProbeDetailsLimitDTO `json:"limits"`
}

// SMTPSENDProbeDetailsDTO wraps SMTP SEND probe details
type SMTPSENDProbeDetailsDTO struct {
	Port    int                             `json:"port,omitempty"`
	From    string                          `json:"from"`
	To      string                          `json:"to"`
	Message string                          `json:"message,omitempty"`
	Limits  map[string]ProbeDetailsLimitDTO `json:"limits"`
}

// DNSProbeDetailsDTO wraps DNS probe details
type DNSProbeDetailsDTO struct {
	Port       int                             `json:"port,omitempty"`
	TCPOnly    bool                            `json:"tcpOnly,omitempty"`
	RecordType string                          `json:"type,omitempty"`
	OwnerName  string                          `json:"ownerName,omitempty"`
	Limits     map[string]ProbeDetailsLimitDTO `json:"limits"`
}

// ProbeListDTO wraps a list of probes
type ProbeListDTO struct {
	Probes     []ProbeInfoDTO `json:"probes"`
	Queryinfo  QueryInfo      `json:"queryInfo"`
	Resultinfo ResultInfo     `json:"resultInfo"`
}

// ProbePath generates the URI path for a probe
func ProbePath(zone, typ, name, guid string) string {
	if guid == "" {
		return fmt.Sprintf("zones/%s/rrsets/%s/%s/probes", zone, typ, name)
	}
	return fmt.Sprintf("zones/%s/rrsets/%s/%s/probes/%s", zone, typ, name, guid)
}

// AlertPath generates the URI path for an alert
func AlertPath(zone, typ, name string) string {
	return fmt.Sprintf("zones/%s/rrsets/%s/%s/alerts", zone, typ, name)
}

func probeQueryPath(zone, typ, name, query string) string {
	uri := ProbePath(zone, typ, name, "")
	if query != "" {
		uri = fmt.Sprintf("%s?sort=NAME&query=%s", uri, query)
	}
	return uri
}

func probeAlertQueryPath(zone, typ, name string, offset int) string {
	baseURI := AlertPath(zone, typ, name)
	return fmt.Sprintf("%s?offset=%d", baseURI, offset)
}

// ListAllProbeAlerts returns all probe alerts with name, type & zone
func (s *SBTCService) ListAllProbeAlerts(name, typ, zone string) ([]ProbeAlertDataDTO, error) {
	// TODO: Sane Configuration for timeouts / retries
	maxerrs := 5
	waittime := 5 * time.Second

	// init accumulators
	as := []ProbeAlertDataDTO{}
	offset := 0
	errcnt := 0

	for {
		reqAlerts, ri, res, err := s.ListProbeAlerts(name, typ, zone, offset)
		if err != nil {
			if res.StatusCode >= 500 {
				errcnt = errcnt + 1
				if errcnt < maxerrs {
					time.Sleep(waittime)
					continue
				}
			}
			return as, err
		}

		log.Printf("ResultInfo: %+v\n", ri)
		for _, a := range reqAlerts {
			as = append(as, a)
		}
		if ri.ReturnedCount+ri.Offset >= ri.TotalCount {
			return as, nil
		}
		offset = ri.ReturnedCount + ri.Offset
		continue
	}
}

// ListProbeAlerts returns the probe alerts with name, type & zone, accepting an offset
func (s *SBTCService) ListProbeAlerts(name, typ, zone string, offset int) ([]ProbeAlertDataDTO, ResultInfo, *Response, error) {
	var ald ProbeAlertDataListDTO

	uri := probeAlertQueryPath(zone, typ, name, offset)
	res, err := s.client.get(uri, &ald)

	as := []ProbeAlertDataDTO{}
	for _, a := range ald.Alerts {
		as = append(as, a)
	}
	return as, ald.Resultinfo, res, err
}

// GetProbe returns a probe with name, type, zone & guid
func (s *SBTCService) GetProbe(name, typ, zone, guid string) (ProbeInfoDTO, *Response, error) {
	var t ProbeInfoDTO
	uri := ProbePath(zone, typ, name, guid)
	res, err := s.client.get(uri, &t)
	return t, res, err
}

// CreateProbe creates a probe with name, type & zone using the ProbeInfoDTO dp
func (s *SBTCService) CreateProbe(name, typ, zone string, dp ProbeInfoDTO) (*Response, error) {
	uri := ProbePath(zone, typ, name, "")
	var ignored interface{}
	return s.client.post(uri, dp, &ignored)
}

// UpdateProbe updates a probe given a name, type, zone & guid with the ProbeInfoDTO dp
func (s *SBTCService) UpdateProbe(name, typ, zone, guid string, dp ProbeInfoDTO) (*Response, error) {
	uri := ProbePath(zone, typ, name, guid)
	var ignored interface{}
	return s.client.put(uri, dp, &ignored)
}

// ListProbes returns all probes by name, type & zone, with an optional query
func (s *SBTCService) ListProbes(query, name, typ, zone string) ([]ProbeInfoDTO, *Response, error) {
	var pld ProbeListDTO

	// This API does not support pagination.
	uri := probeQueryPath(zone, typ, name, query)
	res, err := s.client.get(uri, &pld)

	ps := []ProbeInfoDTO{}
	if err == nil {
		for _, t := range pld.Probes {
			ps = append(ps, t)
		}
	}
	return ps, res, err
}

// DeleteProbe deletes a probe by its name, type, zone & guid
func (s *SBTCService) DeleteProbe(name, typ, zone, guid string) (*Response, error) {
	uri := ProbePath(zone, typ, name, guid)
	return s.client.delete(uri, nil)
}
