package udnssdk

import (
	"fmt"
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
	Start   string `json:"start"`
	End     string `json:"end"`
	CIDR    string `json:"cidr"`
	Address string `json:"address"`
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

// GetDirectionalPoolStatus Get the status of a task.
func (s *DirectionalPoolsService) GetDirectionalPoolStatus(tid string) (DirectionalPool, *Response, error) {
	reqStr := taskPath(tid)
	var t DirectionalPool
	res, err := s.client.get(reqStr, &t)
	if err != nil {
		return t, res, err
	}
	return t, res, err
}

// ListDirectionalPools requests all directional-pools, by query, type & account, returning the list of IP groups & the actual response, or an error
func (s *DirectionalPoolsService) ListDirectionalPools(query, dptype, account string) ([]DirectionalPool, *Response, error) {
	// TODO: Soooo... This function does not handle pagination of DirectionalPools....
	//v := url.Values{}

	reqStr := DirectionalPoolPath(account, dptype, "")
	if query != "" {
		reqStr = fmt.Sprintf("%s?q=%s", reqStr, query)
	}
	fmt.Printf("ListDirectionalPools: %s\n", reqStr)
	var tld DirectionalPoolListDTO
	//wrappedDirectionalPools := []DirectionalPool{}

	res, err := s.client.get(reqStr, &tld)
	if err != nil {
		return []DirectionalPool{}, res, err
	}

	dps := []DirectionalPool{}
	for _, t := range tld.DirectionalPools {
		dps = append(dps, t)
	}

	return dps, res, nil
}

// DeleteDirectionalPool deletes a task.
func (s *DirectionalPoolsService) DeleteDirectionalPool(tid string) (*Response, error) {
	path := taskPath(tid)
	return s.client.delete(path, nil)
}
