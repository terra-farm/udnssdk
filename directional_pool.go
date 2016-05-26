package udnssdk

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"time"
)

// DirectionalPoolsService manages 'account level' 'geo' and 'ip' groups for directional-pools
type DirectionalPoolsService struct {
	client *Client
}

// DirectionalPool wraps an account-level directional-groups response from a index request
type DirectionalPool struct {
	DirectionalPoolID         string `json:"taskId"`
	DirectionalPoolStatusCode string `json:"taskStatusCode"`
	Message                   string `json:"message"`
	ResultURI                 string `json:"resultUri"`
}

// AccountLevelGeoDirectionalGroupDTO wraps an account-level, geo directonal-group response
type AccountLevelGeoDirectionalGroupDTO struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Codes       []string `json:"codes"`
}

// IPAddrDTO wraps an IP address range or CIDR block
type IPAddrDTO struct {
	Start   string `json:"start,omitempty"`
	End     string `json:"end,omitempty"`
	CIDR    string `json:"cidr,omitempty"`
	Address string `json:"address,omitempty"`
}

// Min returns the lowest net.IP member of an IPAddrDTO, returning nil on error
func (i IPAddrDTO) Min() net.IP {
	if i.Start != "" {
		return net.ParseIP(i.Start)
	} else if i.Address != "" {
		return net.ParseIP(i.Address)
	} else if i.CIDR != "" {
		// on error, ip is nil just like ParseIP
		ip, _, _ := net.ParseCIDR(i.CIDR)
		return ip
	}
	return nil
}

// ByIPRange implements sort.Interface for []IPAddrDTO based on
// the Min function.
type ByIPRange []IPAddrDTO

func (is ByIPRange) Len() int           { return len(is) }
func (is ByIPRange) Swap(i, j int)      { is[i], is[j] = is[j], is[i] }
func (is ByIPRange) Less(i, j int) bool { return bytes.Compare(is[i].Min(), is[j].Min()) == -1 }

// AccountLevelIPDirectionalGroupDTO wraps an account-level, IP directional-group response
type AccountLevelIPDirectionalGroupDTO struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	IPs         []IPAddrDTO `json:"ips"`
}

// DirectionalPoolListDTO wraps a list of account-level directional-groups response from a index request
type DirectionalPoolListDTO struct {
	DirectionalPools []DirectionalPool `json:"tasks"`
	Queryinfo        QueryInfo         `json:"queryInfo"`
	Resultinfo       ResultInfo        `json:"resultInfo"`
}

// AccountLevelGeoDirectionalGroupListDTO wraps a list of account-level, geo directional-groups response from a index request
type AccountLevelGeoDirectionalGroupListDTO struct {
	AccountName string                               `json:"zoneName"`
	GeoGroups   []AccountLevelGeoDirectionalGroupDTO `json:"geoGroups"`
	Queryinfo   QueryInfo                            `json:"queryInfo"`
	Resultinfo  ResultInfo                           `json:"resultInfo"`
}

// AccountLevelIPDirectionalGroupListDTO wraps an account-level, IP directional-group response
type AccountLevelIPDirectionalGroupListDTO struct {
	AccountName string                              `json:"zoneName"`
	IPGroups    []AccountLevelIPDirectionalGroupDTO `json:"ipGroups"`
	Queryinfo   QueryInfo                           `json:"queryInfo"`
	Resultinfo  ResultInfo                          `json:"resultInfo"`
}

// Geos allows access to the Geo DirectionalPools API
func (s *DirectionalPoolsService) Geos() *GeoDirectionalPoolsService {
	return &GeoDirectionalPoolsService{client: s.client}
}

// IPs allows access to the IP DirectionalPools API
func (s *DirectionalPoolsService) IPs() *IPDirectionalPoolsService {
	return &IPDirectionalPoolsService{client: s.client}
}

// DirectionalPoolKey collects the identifiers of a DirectionalPool
type DirectionalPoolKey struct {
	Account AccountKey
	Type    string
	Name    string
}

// URI generates the URI for directional pools by account, type & slug ID
func (k DirectionalPoolKey) URI() string {
	if k.Name == "" {
		return fmt.Sprintf("%s/dirgroups/%s", k.Account.URI(), k.Type)
	}
	return fmt.Sprintf("%s/dirgroups/%s/%s", k.Account.URI(), k.Type, k.Name)
}

// QueryURI generates the URI for directional pools by account, type, query & offset
func (k DirectionalPoolKey) QueryURI(query string, offset int) string {
	uri := k.URI()

	if query != "" {
		uri = fmt.Sprintf("%s?sort=NAME&query=%s&offset=%d", uri, query, offset)
	} else {
		uri = fmt.Sprintf("%s?offset=%d", uri, offset)
	}

	return uri
}

// GeoDirectionalPoolKey collects the identifiers of an DirectionalPool with type Geo
type GeoDirectionalPoolKey struct {
	Account AccountKey
	Name    string
}

// DirectionalPoolKey generates the DirectionalPoolKey for the GeoDirectionalPoolKey
func (k GeoDirectionalPoolKey) DirectionalPoolKey() DirectionalPoolKey {
	return DirectionalPoolKey{
		Account: k.Account,
		Type:    "geo",
		Name:    k.Name,
	}
}

// URI generates the URI for a GeoDirectionalPool
func (k GeoDirectionalPoolKey) URI() string {
	return k.DirectionalPoolKey().URI()
}

// QueryURI generates the GeoDirectionalPool URI with query
func (k GeoDirectionalPoolKey) QueryURI(query string, offset int) string {
	return k.DirectionalPoolKey().QueryURI(query, offset)
}

// GeoDirectionalPoolsService manages 'geo' groups for directional-pools
type GeoDirectionalPoolsService struct {
	client *Client
}

// Select requests all geo directional-pools, by query and account, providing pagination and error handling
func (s *GeoDirectionalPoolsService) Select(k GeoDirectionalPoolKey, query string) ([]AccountLevelGeoDirectionalGroupDTO, error) {
	// TODO: Sane Configuration for timeouts / retries
	maxerrs := 5
	waittime := 5 * time.Second

	// init accumulators
	dtos := []AccountLevelGeoDirectionalGroupDTO{}
	errcnt := 0
	offset := 0

	for {
		reqDtos, ri, res, err := s.SelectWithOffset(k, query, offset)
		if err != nil {
			if res.StatusCode >= 500 {
				errcnt = errcnt + 1
				if errcnt < maxerrs {
					time.Sleep(waittime)
					continue
				}
			}
			return dtos, err
		}

		log.Printf("[DEBUG] ResultInfo: %+v\n", ri)
		for _, d := range reqDtos {
			dtos = append(dtos, d)
		}
		if ri.ReturnedCount+ri.Offset >= ri.TotalCount {
			return dtos, nil
		}
		offset = ri.ReturnedCount + ri.Offset
		continue
	}
}

// SelectWithOffset requests list of geo directional-pools, by query & account, and an offset, returning the directional-group, the list-metadata, the actual response, or an error
func (s *GeoDirectionalPoolsService) SelectWithOffset(k GeoDirectionalPoolKey, query string, offset int) ([]AccountLevelGeoDirectionalGroupDTO, ResultInfo, *Response, error) {
	var tld AccountLevelGeoDirectionalGroupListDTO

	res, err := s.client.get(k.QueryURI(query, offset), &tld)

	pis := []AccountLevelGeoDirectionalGroupDTO{}
	for _, pi := range tld.GeoGroups {
		pis = append(pis, pi)
	}
	return pis, tld.Resultinfo, res, err
}

// Find requests a geo directional-pool by name & account
func (s *GeoDirectionalPoolsService) Find(k GeoDirectionalPoolKey) (AccountLevelGeoDirectionalGroupDTO, *Response, error) {
	var t AccountLevelGeoDirectionalGroupDTO
	res, err := s.client.get(k.URI(), &t)
	return t, res, err
}

// Create requests creation of a DirectionalPool by DirectionalPoolKey given a directional-pool
func (s *GeoDirectionalPoolsService) Create(k GeoDirectionalPoolKey, val interface{}) (*Response, error) {
	return s.client.post(k.URI(), val, nil)
}

// Update requests update of a DirectionalPool by DirectionalPoolKey given a directional-pool
func (s *GeoDirectionalPoolsService) Update(k GeoDirectionalPoolKey, val interface{}) (*Response, error) {
	return s.client.put(k.URI(), val, nil)
}

// Delete requests deletion of a DirectionalPool
func (s *GeoDirectionalPoolsService) Delete(k GeoDirectionalPoolKey) (*Response, error) {
	return s.client.delete(k.URI(), nil)
}

// IPDirectionalPoolKey collects the identifiers of an DirectionalPool with type IP
type IPDirectionalPoolKey struct {
	Account AccountKey
	Name    string
}

// DirectionalPoolKey generates the DirectionalPoolKey for the IPDirectionalPoolKey
func (k IPDirectionalPoolKey) DirectionalPoolKey() DirectionalPoolKey {
	return DirectionalPoolKey{
		Account: k.Account,
		Type:    "ip",
		Name:    k.Name,
	}
}

// URI generates the IPDirectionalPool query URI
func (k IPDirectionalPoolKey) URI() string {
	return k.DirectionalPoolKey().URI()
}

// QueryURI generates the IPDirectionalPool URI with query
func (k IPDirectionalPoolKey) QueryURI(query string, offset int) string {
	return k.DirectionalPoolKey().QueryURI(query, offset)
}

// IPDirectionalPoolsService manages 'geo' groups for directional-pools
type IPDirectionalPoolsService struct {
	client *Client
}

// Select requests all IP directional-pools, using pagination and error handling
func (s *IPDirectionalPoolsService) Select(k IPDirectionalPoolKey, query string) ([]AccountLevelIPDirectionalGroupDTO, error) {
	// TODO: Sane Configuration for timeouts / retries
	maxerrs := 5
	waittime := 5 * time.Second

	// init accumulators
	gs := []AccountLevelIPDirectionalGroupDTO{}
	errcnt := 0
	offset := 0

	for {
		reqIPGroups, ri, res, err := s.SelectWithOffset(k, query, offset)
		if err != nil {
			if res.StatusCode >= 500 {
				errcnt = errcnt + 1
				if errcnt < maxerrs {
					time.Sleep(waittime)
					continue
				}
			}
			return gs, err
		}

		log.Printf("ResultInfo: %+v\n", ri)
		for _, g := range reqIPGroups {
			gs = append(gs, g)
		}
		if ri.ReturnedCount+ri.Offset >= ri.TotalCount {
			return gs, nil
		}
		offset = ri.ReturnedCount + ri.Offset
		continue
	}
}

// SelectWithOffset requests all IP directional-pools, by query & account, and an offset, returning the list of IP groups, list metadata & the actual response, or an error
func (s *IPDirectionalPoolsService) SelectWithOffset(k IPDirectionalPoolKey, query string, offset int) ([]AccountLevelIPDirectionalGroupDTO, ResultInfo, *Response, error) {
	var tld AccountLevelIPDirectionalGroupListDTO

	res, err := s.client.get(k.QueryURI(query, offset), &tld)

	pis := []AccountLevelIPDirectionalGroupDTO{}
	for _, pi := range tld.IPGroups {
		pis = append(pis, pi)
	}

	return pis, tld.Resultinfo, res, err
}

// Find requests a directional-pool by name & account
func (s *IPDirectionalPoolsService) Find(k IPDirectionalPoolKey) (AccountLevelIPDirectionalGroupDTO, *Response, error) {
	var t AccountLevelIPDirectionalGroupDTO
	res, err := s.client.get(k.URI(), &t)
	return t, res, err
}

// Create requests creation of a DirectionalPool by DirectionalPoolKey given a directional-pool
func (s *IPDirectionalPoolsService) Create(k IPDirectionalPoolKey, val interface{}) (*Response, error) {
	return s.client.post(k.URI(), val, nil)
}

// Update requests update of a DirectionalPool by DirectionalPoolKey given a directional-pool
func (s *IPDirectionalPoolsService) Update(k IPDirectionalPoolKey, val interface{}) (*Response, error) {
	return s.client.put(k.URI(), val, nil)
}

// Delete deletes an  directional-pool
func (s *IPDirectionalPoolsService) Delete(k IPDirectionalPoolKey) (*Response, error) {
	return s.client.delete(k.URI(), nil)
}
