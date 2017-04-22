package udnssdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ProbeType wraps the possible types of a ProbeInfo
type ProbeType string

// Here lie all the possible ProbeType values
const (
	DNSProbeType      ProbeType = "DNS"
	FTPProbeType      ProbeType = "FTP"
	HTTPProbeType     ProbeType = "HTTP"
	PingProbeType     ProbeType = "PING"
	SMTPProbeType     ProbeType = "SMTP"
	SMTPSENDProbeType ProbeType = "SMTP_SEND"
	TCPProbeType      ProbeType = "TCP"
)

// ProbeInfo wraps a probe response
type ProbeInfo struct {
	ID         string        `json:"id,omitempty"`
	PoolRecord string        `json:"poolRecord,omitempty"`
	ProbeType  ProbeType     `json:"type"`
	Interval   string        `json:"interval"`
	Agents     []string      `json:"agents"`
	Threshold  int           `json:"threshold"`
	Details    *ProbeDetails `json:"details"`
}

// ProbeDetailsLimit wraps a probe
type ProbeDetailsLimit struct {
	Warning  int `json:"warning"`
	Critical int `json:"critical"`
	Fail     int `json:"fail"`
}

// ProbeDetails wraps the details of a probe
type ProbeDetails struct {
	data   []byte
	Detail interface{} `json:"detail,omitempty"`
	typ    ProbeType
}

// GetData returns the data because I'm working around something.
func (s *ProbeDetails) GetData() []byte {
	return s.data
}

// Populate does magical things with json unmarshalling to unroll the Probe into
// an appropriate datatype.  These are helper structures and functions for testing
// and direct API use.  In the Terraform implementation, we will use Terraforms own
// warped schema structure to handle the marshalling and unmarshalling.
func (s *ProbeDetails) Populate(t ProbeType) (err error) {
	s.typ = t
	d, err := s.GetDetailsObject(t)
	if err != nil {
		return err
	}
	s.Detail = d
	return nil
}

// GetDetailsObject extracts the appropriate details object from a ProbeDetails with the given ProbeType
func (s *ProbeDetails) GetDetailsObject(t ProbeType) (interface{}, error) {
	switch t {
	case DNSProbeType:
		return s.DNSProbeDetails()
	case FTPProbeType:
		return s.FTPProbeDetails()
	case HTTPProbeType:
		return s.HTTPProbeDetails()
	case PingProbeType:
		return s.PingProbeDetails()
	case SMTPProbeType:
		return s.SMTPProbeDetails()
	case SMTPSENDProbeType:
		return s.SMTPSENDProbeDetails()
	case TCPProbeType:
		return s.TCPProbeDetails()
	default:
		return nil, fmt.Errorf("Invalid ProbeType: %#v", t)
	}
}

// DNSProbeDetails returns the ProbeDetails data deserialized as a DNSProbeDetails
func (s *ProbeDetails) DNSProbeDetails() (DNSProbeDetails, error) {
	var d DNSProbeDetails
	err := json.Unmarshal(s.data, &d)
	return d, err
}

// FTPProbeDetails returns the ProbeDetails data deserialized as a FTPProbeDetails
func (s *ProbeDetails) FTPProbeDetails() (FTPProbeDetails, error) {
	var d FTPProbeDetails
	err := json.Unmarshal(s.data, &d)
	return d, err
}

// HTTPProbeDetails returns the ProbeDetails data deserialized as a HTTPProbeDetails
func (s *ProbeDetails) HTTPProbeDetails() (HTTPProbeDetails, error) {
	var d HTTPProbeDetails
	err := json.Unmarshal(s.data, &d)
	return d, err
}

// PingProbeDetails returns the ProbeDetails data deserialized as a PingProbeDetails
func (s *ProbeDetails) PingProbeDetails() (PingProbeDetails, error) {
	var d PingProbeDetails
	err := json.Unmarshal(s.data, &d)
	return d, err
}

// SMTPProbeDetails returns the ProbeDetails data deserialized as a SMTPProbeDetails
func (s *ProbeDetails) SMTPProbeDetails() (SMTPProbeDetails, error) {
	var d SMTPProbeDetails
	err := json.Unmarshal(s.data, &d)
	return d, err
}

// SMTPSENDProbeDetails returns the ProbeDetails data deserialized as a SMTPSENDProbeDetails
func (s *ProbeDetails) SMTPSENDProbeDetails() (SMTPSENDProbeDetails, error) {
	var d SMTPSENDProbeDetails
	err := json.Unmarshal(s.data, &d)
	return d, err
}

// TCPProbeDetails returns the ProbeDetails data deserialized as a TCPProbeDetails
func (s *ProbeDetails) TCPProbeDetails() (TCPProbeDetails, error) {
	var d TCPProbeDetails
	err := json.Unmarshal(s.data, &d)
	return d, err
}

// UnmarshalJSON does what it says on the tin
func (s *ProbeDetails) UnmarshalJSON(b []byte) (err error) {
	s.data = b
	return nil
}

// MarshalJSON does what it says on the tin
func (s *ProbeDetails) MarshalJSON() ([]byte, error) {
	var err error
	if s.Detail != nil {
		return json.Marshal(s.Detail)
	}
	if len(s.data) != 0 {
		return s.data, err
	}
	return json.Marshal(nil)
}

// GoString returns a string representation of the ProbeDetails internal data
func (s *ProbeDetails) GoString() string {
	return string(s.data)
}
func (s *ProbeDetails) String() string {
	return string(s.data)
}

// Transaction wraps a transaction response
type Transaction struct {
	Method          string                       `json:"method"`
	URL             string                       `json:"url"`
	TransmittedData string                       `json:"transmittedData,omitempty"`
	FollowRedirects bool                         `json:"followRedirects,omitempty"`
	Limits          map[string]ProbeDetailsLimit `json:"limits"`
}

// HTTPProbeDetails wraps HTTP probe details
type HTTPProbeDetails struct {
	Transactions []Transaction      `json:"transactions"`
	TotalLimits  *ProbeDetailsLimit `json:"totalLimits,omitempty"`
}

// PingProbeDetails wraps Ping probe details
type PingProbeDetails struct {
	Packets    int                          `json:"packets,omitempty"`
	PacketSize int                          `json:"packetSize,omitempty"`
	Limits     map[string]ProbeDetailsLimit `json:"limits"`
}

// FTPProbeDetails wraps FTP probe details
type FTPProbeDetails struct {
	Port        int                          `json:"port,omitempty"`
	PassiveMode bool                         `json:"passiveMode,omitempty"`
	Username    string                       `json:"username,omitempty"`
	Password    string                       `json:"password,omitempty"`
	Path        string                       `json:"path"`
	Limits      map[string]ProbeDetailsLimit `json:"limits"`
}

// TCPProbeDetails wraps TCP probe details
type TCPProbeDetails struct {
	Port      int                          `json:"port,omitempty"`
	ControlIP string                       `json:"controlIP,omitempty"`
	Limits    map[string]ProbeDetailsLimit `json:"limits"`
}

// SMTPProbeDetails wraps SMTP probe details
type SMTPProbeDetails struct {
	Port   int                          `json:"port,omitempty"`
	Limits map[string]ProbeDetailsLimit `json:"limits"`
}

// SMTPSENDProbeDetails wraps SMTP SEND probe details
type SMTPSENDProbeDetails struct {
	Port    int                          `json:"port,omitempty"`
	From    string                       `json:"from"`
	To      string                       `json:"to"`
	Message string                       `json:"message,omitempty"`
	Limits  map[string]ProbeDetailsLimit `json:"limits"`
}

// DNSProbeDetails wraps DNS probe details
type DNSProbeDetails struct {
	Port       int                          `json:"port,omitempty"`
	TCPOnly    bool                         `json:"tcpOnly,omitempty"`
	RecordType string                       `json:"type,omitempty"`
	OwnerName  string                       `json:"ownerName,omitempty"`
	Limits     map[string]ProbeDetailsLimit `json:"limits"`
}

// ProbeList wraps a list of probes
type ProbeList struct {
	Probes     []ProbeInfo `json:"probes"`
	Queryinfo  QueryInfo   `json:"queryInfo"`
	Resultinfo ResultInfo  `json:"resultInfo"`
}

// ProbesService manages Probes
type ProbesService struct {
	client *Client
}

// ProbeKey collects the identifiers of a Probe
type ProbeKey struct {
	Zone string
	Name string
	ID   string
}

// RRSetKey generates the RRSetKey for the ProbeKey
func (k ProbeKey) RRSetKey() RRSetKey {
	return RRSetKey{
		Zone: k.Zone,
		Type: "A", // Only A records have probes
		Name: k.Name,
	}
}

// URI generates the URI for a probe
func (k ProbeKey) URI() string {
	return fmt.Sprintf("%s/%s", k.RRSetKey().ProbesURI(), k.ID)
}

// Select returns all probes by a RRSetKey, with an optional query
func (s *ProbesService) Select(k RRSetKey, query string) ([]ProbeInfo, *http.Response, error) {
	var pld ProbeList

	// This API does not support pagination.
	uri := k.ProbesQueryURI(query)
	res, err := s.client.get(uri, &pld)

	ps := []ProbeInfo{}
	if err == nil {
		for _, t := range pld.Probes {
			ps = append(ps, t)
		}
	}
	return ps, res, err
}

// Find returns a probe from a ProbeKey
func (s *ProbesService) Find(k ProbeKey) (ProbeInfo, *http.Response, error) {
	var t ProbeInfo
	res, err := s.client.get(k.URI(), &t)
	return t, res, err
}

// Create creates a probe with a RRSetKey using the ProbeInfo dp
func (s *ProbesService) Create(k RRSetKey, dp ProbeInfo) (*http.Response, error) {
	return s.client.post(k.ProbesURI(), dp, nil)
}

// Update updates a probe given a ProbeKey with the ProbeInfo dp
func (s *ProbesService) Update(k ProbeKey, dp ProbeInfo) (*http.Response, error) {
	return s.client.put(k.URI(), dp, nil)
}

// Delete deletes a probe by its ProbeKey
func (s *ProbesService) Delete(k ProbeKey) (*http.Response, error) {
	return s.client.delete(k.URI(), nil)
}
