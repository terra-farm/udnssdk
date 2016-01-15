package udnssdk

import (
	"fmt"
)

// ZonesService handles communication with the Zone related blah blah
type DirectionalPoolsService struct {
	client *Client
}

type DirectionalPool struct {
	DirectionalPoolId         string `json:"taskId"`
	DirectionalPoolStatusCode string `json:"taskStatusCode"`
	Message                   string `json:"message"`
	ResultUri                 string `json:"resultUri"`
}
type AccountLevelGeoDirectionalGroupDTO struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Codes       []string `json:"codes"`
}
type IPAddrDTO struct {
	Start   string `json:"start"`
	End     string `json:"end"`
	Cidr    string `json:"cidr"`
	Address string `json:"address"`
}
type AccountLevelIPDirectionalGroupDTO struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Ips         []IPAddrDTO `json:"ips"`
}

type DirectionalPoolListDTO struct {
	DirectionalPools []DirectionalPool `json:"tasks"`
	Queryinfo        QueryInfo         `json:"queryInfo"`
	Resultinfo       ResultInfo        `json:"resultInfo"`
}

type AccountLevelGeoDirectionalGroupListDTO struct {
	AccountName string                               `json:"zoneName"`
	GeoGroups   []AccountLevelGeoDirectionalGroupDTO `json:"geoGroups"`
	Queryinfo   QueryInfo                            `json:"queryInfo"`
	Resultinfo  ResultInfo                           `json:"resultInfo"`
}

type AccountLevelIPDirectionalGroupListDTO struct {
	AccountName string                              `json:"zoneName"`
	IpGroups    []AccountLevelIPDirectionalGroupDTO `json:"ipGroups"`
	Queryinfo   QueryInfo                           `json:"queryInfo"`
	Resultinfo  ResultInfo                          `json:"resultInfo"`
}

/*
type taskWrapper struct {
	DirectionalPool DirectionalPool `json:"task"`
}
*/
func DirectionalPoolPath(acct, typ, val string) string {
	if val == "" {
		return fmt.Sprintf("accounts/%s/dirgroups/%s", acct, typ)
	} else {

		return fmt.Sprintf("accounts/%s/dirgroups/geo/%s/%s", acct, typ, val)
	}
}

// Get the status of a task.
func (s *DirectionalPoolsService) GetDirectionalPoolStatus(tid string) (DirectionalPool, *Response, error) {
	reqStr := taskPath(tid)
	var t DirectionalPool
	res, err := s.client.get(reqStr, &t)
	if err != nil {
		return t, res, err
	}
	return t, res, err
}

// List tasks
//
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
//
func (s *DirectionalPoolsService) DeleteDirectionalPool(tid string) (*Response, error) {
	path := taskPath(tid)
	return s.client.delete(path, nil)
}
