package udnssdk

import (
	"fmt"
	"log"
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

// DirectionalPoolPath generates the URI for directional pools by account, type & slug ID
func DirectionalPoolPath(acct, typ, slugID string) string {
	if slugID == "" {
		return fmt.Sprintf("accounts/%s/dirgroups/%s", acct, typ)
	}
	return fmt.Sprintf("accounts/%s/dirgroups/%s/%s", acct, typ, slugID)
}

// DirectionalPoolQueryPath generates the URI for directional pools by account, type, query & offset
func DirectionalPoolQueryPath(account, typ, query string, offset int) string {
	uri := DirectionalPoolPath(account, typ, "")

	if query != "" {
		uri = fmt.Sprintf("%s?sort=NAME&query=%s&offset=%d", uri, query, offset)
	} else {
		uri = fmt.Sprintf("%s?offset=%d", uri, offset)
	}

	return uri
}

// DirectionalIPPoolQueryPath generates the URI for a directional IP pool by account, query & offset
func DirectionalIPPoolQueryPath(account, query string, offset int) string {
	return DirectionalPoolQueryPath(account, "ip", query, offset)
}

// DirectionalGeoPoolQueryPath generates the URI for a directional geo pool by account, query & offset
func DirectionalGeoPoolQueryPath(account, query string, offset int) string {
	return DirectionalPoolQueryPath(account, "geo", query, offset)
}

// GetDirectionalGeoPool requests a geo directional-pool by name & account
func (s *DirectionalPoolsService) GetDirectionalGeoPool(name, acct string) (AccountLevelGeoDirectionalGroupDTO, *Response, error) {
	uri := DirectionalPoolPath(acct, "geo", name)
	var t AccountLevelGeoDirectionalGroupDTO
	res, err := s.client.get(uri, &t)
	return t, res, err
}

// GetDirectionalIPPool requests a IP directional-pool by name & account
func (s *DirectionalPoolsService) GetDirectionalIPPool(name, acct string) (AccountLevelIPDirectionalGroupDTO, *Response, error) {
	uri := DirectionalPoolPath(acct, "ip", name)
	var t AccountLevelIPDirectionalGroupDTO
	res, err := s.client.get(uri, &t)
	return t, res, err
}

// CreateDirectionalGeoPool requests creation of a geo direcctional-pool by name & account, given a directional-pool
func (s *DirectionalPoolsService) CreateDirectionalGeoPool(name, acct string, dp AccountLevelGeoDirectionalGroupDTO) (*Response, error) {
	uri := DirectionalPoolPath(acct, "geo", name)
	var ignored interface{}
	return s.client.post(uri, dp, &ignored)
}

// CreateDirectionalIPPool requests creation of an IP directional-pool by name & account, given a directional-pool
func (s *DirectionalPoolsService) CreateDirectionalIPPool(name, acct string, dp AccountLevelIPDirectionalGroupDTO) (*Response, error) {
	uri := DirectionalPoolPath(acct, "ip", name)
	var ignored interface{}
	return s.client.post(uri, dp, &ignored)
}

// UpdateDirectionalGeoPool requests update of a geo directional-pool by name & account, given a directional-pool
func (s *DirectionalPoolsService) UpdateDirectionalGeoPool(name, acct string, dp AccountLevelGeoDirectionalGroupDTO) (*Response, error) {
	uri := DirectionalPoolPath(acct, "geo", name)
	var ignored interface{}
	return s.client.put(uri, dp, &ignored)
}

// UpdateDirectionalIPPool requests update of an IP directional-pool by name & account, given a directional-pool
func (s *DirectionalPoolsService) UpdateDirectionalIPPool(name, acct string, dp AccountLevelIPDirectionalGroupDTO) (*Response, error) {
	uri := DirectionalPoolPath(acct, "ip", name)
	var ignored interface{}
	return s.client.put(uri, dp, &ignored)
}

// ListAllDirectionalGeoPools requests all geo directional-pools, by query and account, providing pagination and error handling
func (s *DirectionalPoolsService) ListAllDirectionalGeoPools(query, account string) ([]AccountLevelGeoDirectionalGroupDTO, error) {
	// TODO: Sane Configuration for timeouts / retries
	maxerrs := 5
	waittime := 5 * time.Second

	// init accumulators
	dtos := []AccountLevelGeoDirectionalGroupDTO{}
	errcnt := 0
	offset := 0

	for {
		reqDtos, ri, res, err := s.ListDirectionalGeoPools(query, account, offset)
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

// ListDirectionalGeoPools requests list of geo directional-pools, by query & account, and an offset, returning the directional-group, the list-metadata, the actual response, or an error
func (s *DirectionalPoolsService) ListDirectionalGeoPools(query, account string, offset int) ([]AccountLevelGeoDirectionalGroupDTO, ResultInfo, *Response, error) {
	var tld AccountLevelGeoDirectionalGroupListDTO

	uri := DirectionalGeoPoolQueryPath(account, query, offset)
	res, err := s.client.get(uri, &tld)

	pis := []AccountLevelGeoDirectionalGroupDTO{}
	for _, pi := range tld.GeoGroups {
		pis = append(pis, pi)
	}
	return pis, tld.Resultinfo, res, err
}

// ListAllDirectionalIPPools requests all IP directional-pools, using pagination and error handling
func (s *DirectionalPoolsService) ListAllDirectionalIPPools(query, account string) ([]AccountLevelIPDirectionalGroupDTO, error) {
	// TODO: Sane Configuration for timeouts / retries
	maxerrs := 5
	waittime := 5 * time.Second

	// init accumulators
	gs := []AccountLevelIPDirectionalGroupDTO{}
	errcnt := 0
	offset := 0

	for {
		reqIPGroups, ri, res, err := s.ListDirectionalIPPools(query, account, offset)
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

// ListDirectionalIPPools requests all IP directional-pools, by query & account, and an offset, returning the list of IP groups, list metadata & the actual response, or an error
func (s *DirectionalPoolsService) ListDirectionalIPPools(query, account string, offset int) ([]AccountLevelIPDirectionalGroupDTO, ResultInfo, *Response, error) {
	var tld AccountLevelIPDirectionalGroupListDTO

	uri := DirectionalIPPoolQueryPath(account, query, offset)
	res, err := s.client.get(uri, &tld)

	pis := []AccountLevelIPDirectionalGroupDTO{}
	for _, pi := range tld.IPGroups {
		pis = append(pis, pi)
	}

	return pis, tld.Resultinfo, res, err
}

// DeleteDirectionalGeoPool deletes a geo directional-pool
func (s *DirectionalPoolsService) DeleteDirectionalGeoPool(dp, acct string) (*Response, error) {
	path := DirectionalPoolPath(acct, "geo", dp)
	return s.client.delete(path, nil)
}

// DeleteDirectionalIPPool deletes an IP directional-pool
func (s *DirectionalPoolsService) DeleteDirectionalIPPool(dp, acct string) (*Response, error) {
	path := DirectionalPoolPath(acct, "ip", dp)
	return s.client.delete(path, nil)
}
