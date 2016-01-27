package udnssdk

import (
	"fmt"
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

// ListDirectionalGeoPools requests list of geo directional-pools, by query & account, and an offset, returning the directional-group, the list-metadata, the actual response, or an error
func (s *DirectionalPoolsService) ListDirectionalGeoPools(query, account string) ([]AccountLevelGeoDirectionalGroupDTO, *Response, error) {
	uri := DirectionalPoolPath(account, "geo", "")
	if query != "" {
		uri = fmt.Sprintf("%s?sort=NAME&query=%s", uri, query)
	}
	fmt.Printf("ListDirectionalPools: %s\n", uri)
	var tld AccountLevelGeoDirectionalGroupListDTO

	res, err := s.client.get(uri, &tld)
	pis := []AccountLevelGeoDirectionalGroupDTO{}
	if query != "" {
		uri = fmt.Sprintf("%s?sort=NAME&query=%s&offset=", uri, query)
	} else {
		uri = fmt.Sprintf("%s?offset=", uri)
	}
	// TODO: Sane Configuration for timeouts / retries
	timeout := 5
	waittime := 5 * time.Second
	errcnt := 0
	offset := 0
	for true {

		res, err := s.client.get(fmt.Sprintf("%s%d", uri, offset), &tld)
		if err != nil {
			if res.StatusCode >= 500 {
				errcnt = errcnt + 1
				if errcnt < timeout {
					time.Sleep(waittime)
					continue
				}
			}
			return pis, res, err
		}
		fmt.Printf("ResultInfo: %+v\n", tld.Resultinfo)
		for _, pi := range tld.GeoGroups {
			pis = append(pis, pi)
		}
		if tld.Resultinfo.ReturnedCount+tld.Resultinfo.Offset >= tld.Resultinfo.TotalCount {
			return pis, res, nil
		}
		offset = tld.Resultinfo.ReturnedCount + tld.Resultinfo.Offset
		continue
	}
	return pis, res, err
}

// ListDirectionalIPPools requests all IP directional-pools, by query & account, and an offset, returning the list of IP groups, list metadata & the actual response, or an error
func (s *DirectionalPoolsService) ListDirectionalIPPools(query, account string) ([]AccountLevelIPDirectionalGroupDTO, *Response, error) {
	uri := DirectionalPoolPath(account, "ip", "")
	if query != "" {
		uri = fmt.Sprintf("%s?sort=NAME&query=%s", uri, query)
	}
	fmt.Printf("ListDirectionalPools: %s\n", uri)
	var tld AccountLevelIPDirectionalGroupListDTO

	res, err := s.client.get(uri, &tld)
	pis := []AccountLevelIPDirectionalGroupDTO{}
	if query != "" {
		uri = fmt.Sprintf("%s?sort=NAME&query=%s&offset=", uri, query)
	} else {
		uri = fmt.Sprintf("%s?offset=", uri)
	}
	// TODO: Sane Configuration for timeouts / retries
	timeout := 5
	waittime := 5 * time.Second
	errcnt := 0
	offset := 0
	for true {

		res, err := s.client.get(fmt.Sprintf("%s%d", uri, offset), &tld)
		if err != nil {
			if res.StatusCode >= 500 {
				errcnt = errcnt + 1
				if errcnt < timeout {
					time.Sleep(waittime)
					continue
				}
			}
			return pis, res, err
		}
		fmt.Printf("ResultInfo: %+v\n", tld.Resultinfo)
		for _, pi := range tld.IPGroups {
			pis = append(pis, pi)
		}
		if tld.Resultinfo.ReturnedCount+tld.Resultinfo.Offset >= tld.Resultinfo.TotalCount {
			return pis, res, nil
		}
		offset = tld.Resultinfo.ReturnedCount + tld.Resultinfo.Offset
		continue
	}
	return pis, res, err
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
