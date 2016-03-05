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

// ProbeKey collects the identifiers of a Probe
type ProbeKey struct {
	Zone string
	Type string
	Name string
	GUID string
}

// RRSetKey generates the RRSetKey for the ProbeKey
func (p *ProbeKey) RRSetKey() *RRSetKey {
	return &RRSetKey{
		Zone: p.Zone,
		Type: p.Type,
		Name: p.Name,
	}
}

// URI generates the URI for a probe
func (p *ProbeKey) URI() string {
	return fmt.Sprintf("%s/%s", p.RRSetKey().ProbesURI(), p.GUID)
}

// ProbesService manages Probes
type ProbesService struct {
	client *Client
}

// Probes allows access to the Probes API
func (s *SBTCService) Probes() *ProbesService {
	return &ProbesService{client: s.client}
}

// Select returns all probes by a RRSetKey, with an optional query
func (s *ProbesService) Select(r RRSetKey, query string) ([]ProbeInfoDTO, *Response, error) {
	var pld ProbeListDTO

	// This API does not support pagination.
	uri := r.ProbesQueryURI(query)
	res, err := s.client.get(uri, &pld)

	ps := []ProbeInfoDTO{}
	if err == nil {
		for _, t := range pld.Probes {
			ps = append(ps, t)
		}
	}
	return ps, res, err
}

// Find returns a probe from a ProbeKey
func (s *ProbesService) Find(p ProbeKey) (ProbeInfoDTO, *Response, error) {
	var t ProbeInfoDTO
	res, err := s.client.get(p.URI(), &t)
	return t, res, err
}

// Create creates a probe with a RRSetKey using the ProbeInfoDTO dp
func (s *ProbesService) Create(r RRSetKey, dp ProbeInfoDTO) (*Response, error) {
	var ignored interface{}
	return s.client.post(r.ProbesURI(), dp, &ignored)
}

// Update updates a probe given a ProbeKey with the ProbeInfoDTO dp
func (s *ProbesService) Update(p ProbeKey, dp ProbeInfoDTO) (*Response, error) {
	return s.client.put(p.URI(), dp, nil)
}

// Delete deletes a probe by its ProbeKey
func (s *ProbesService) Delete(p ProbeKey) (*Response, error) {
	return s.client.delete(p.URI(), nil)
}

// AlertsService manages Alerts
type AlertsService struct {
	client *Client
}

// Alerts allows access to the Alerts API
func (s *SBTCService) Alerts() *AlertsService {
	return &AlertsService{client: s.client}
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

// Select returns all probe alerts with a RRSetKey
func (s *AlertsService) Select(z RRSetKey) ([]ProbeAlertDataDTO, error) {
	// TODO: Sane Configuration for timeouts / retries
	maxerrs := 5
	waittime := 5 * time.Second

	// init accumulators
	as := []ProbeAlertDataDTO{}
	offset := 0
	errcnt := 0

	for {
		reqAlerts, ri, res, err := s.SelectWithOffset(z, offset)
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

// SelectWithOffset returns the probe alerts with a RRSetKey, accepting an offset
func (s *AlertsService) SelectWithOffset(z RRSetKey, offset int) ([]ProbeAlertDataDTO, ResultInfo, *Response, error) {
	var ald ProbeAlertDataListDTO

	uri := z.AlertsQueryURI(offset)
	res, err := s.client.get(uri, &ald)

	as := []ProbeAlertDataDTO{}
	for _, a := range ald.Alerts {
		as = append(as, a)
	}
	return as, ald.Resultinfo, res, err
}
