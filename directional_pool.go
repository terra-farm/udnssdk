package udnssdk

import (
	"fmt"
	"time"
)

/* Directional Pools - This manages 'account level' 'geo' and 'ip' groups for
   Directional Pools.  */

type DirectionalPoolsService struct {
	client *Client
}

type AccountLevelGeoDirectionalGroupDTO struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Codes       []string `json:"codes"`
}
type IPAddrDTO struct {
	Start   string `json:"start,omitempty"`
	End     string `json:"end,omitempty"`
	Cidr    string `json:"cidr,omitempty"`
	Address string `json:"address,omitempty"`
}
type AccountLevelIPDirectionalGroupDTO struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Ips         []IPAddrDTO `json:"ips"`
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

		return fmt.Sprintf("accounts/%s/dirgroups/%s/%s", acct, typ, val)
	}
}

// Get a Direction Geo Pool.
func (s *DirectionalPoolsService) GetDirectionalGeoPool(name, acct string) (AccountLevelGeoDirectionalGroupDTO, *Response, error) {
	reqStr := DirectionalPoolPath(acct, "geo", name)
	var t AccountLevelGeoDirectionalGroupDTO
	res, err := s.client.get(reqStr, &t)
	return t, res, err
}

// Get a Direction Geo Pool.
func (s *DirectionalPoolsService) GetDirectionalIPPool(name, acct string) (AccountLevelIPDirectionalGroupDTO, *Response, error) {
	reqStr := DirectionalPoolPath(acct, "ip", name)
	var t AccountLevelIPDirectionalGroupDTO
	res, err := s.client.get(reqStr, &t)
	return t, res, err
}

// Create a Direction Pool
func (s *DirectionalPoolsService) CreateDirectionalGeoPool(name, acct string, dp AccountLevelGeoDirectionalGroupDTO) (*Response, error) {
	reqStr := DirectionalPoolPath(acct, "geo", name)
	var retval interface{}
	res, err := s.client.post(reqStr, dp, &retval)

	return res, err
}

// Create a Direction Pool
func (s *DirectionalPoolsService) CreateDirectionalIPPool(name, acct string, dp AccountLevelIPDirectionalGroupDTO) (*Response, error) {
	reqStr := DirectionalPoolPath(acct, "ip", name)
	var retval interface{}
	res, err := s.client.post(reqStr, dp, &retval)

	return res, err
}

// Update
func (s *DirectionalPoolsService) UpdateDirectionalGeoPool(name, acct string, dp AccountLevelGeoDirectionalGroupDTO) (*Response, error) {
	reqStr := DirectionalPoolPath(acct, "geo", name)
	var retval interface{}
	res, err := s.client.put(reqStr, dp, &retval)

	return res, err
}

// Update
func (s *DirectionalPoolsService) UpdateDirectionalIPPool(name, acct string, dp AccountLevelIPDirectionalGroupDTO) (*Response, error) {
	reqStr := DirectionalPoolPath(acct, "ip", name)
	var retval interface{}
	res, err := s.client.put(reqStr, dp, &retval)

	return res, err
}

// List Directional Pools
func (s *DirectionalPoolsService) ListDirectionalGeoPools(query, account string) ([]AccountLevelGeoDirectionalGroupDTO, *Response, error) {

	reqStr := DirectionalPoolPath(account, "geo", "")
	if query != "" {
		reqStr = fmt.Sprintf("%s?sort=NAME&query=%s", reqStr, query)
	}
	fmt.Printf("ListDirectionalPools: %s\n", reqStr)
	var tld AccountLevelGeoDirectionalGroupListDTO

	res, err := s.client.get(reqStr, &tld)
	pis := []AccountLevelGeoDirectionalGroupDTO{}
	if query != "" {
		reqStr = fmt.Sprintf("%s?sort=NAME&query=%s&offset=", reqStr, query)
	} else {
		reqStr = fmt.Sprintf("%s?offset=", reqStr)
	}
	// TODO: Sane Configuration for timeouts / retries
	timeout := 5
	waittime := 5 * time.Second
	errcnt := 0
	offset := 0
	for true {

		res, err := s.client.get(fmt.Sprintf("%s%d", reqStr, offset), &tld)
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
		} else {
			offset = tld.Resultinfo.ReturnedCount + tld.Resultinfo.Offset
			continue
		}
	}
	return pis, res, err
}

// List Directional Pools
func (s *DirectionalPoolsService) ListDirectionalIPPools(query, account string) ([]AccountLevelIPDirectionalGroupDTO, *Response, error) {
	// TODO: Soooo... This function does not handle pagination of DirectionalPools....
	//v := url.Values{}

	reqStr := DirectionalPoolPath(account, "ip", "")
	if query != "" {
		reqStr = fmt.Sprintf("%s?sort=NAME&query=%s", reqStr, query)
	}
	fmt.Printf("ListDirectionalPools: %s\n", reqStr)
	var tld AccountLevelIPDirectionalGroupListDTO

	res, err := s.client.get(reqStr, &tld)
	pis := []AccountLevelIPDirectionalGroupDTO{}
	if query != "" {
		reqStr = fmt.Sprintf("%s?sort=NAME&query=%s&offset=", reqStr, query)
	} else {
		reqStr = fmt.Sprintf("%s?offset=", reqStr)
	}
	// TODO: Sane Configuration for timeouts / retries
	timeout := 5
	waittime := 5 * time.Second
	errcnt := 0
	offset := 0
	for true {

		res, err := s.client.get(fmt.Sprintf("%s%d", reqStr, offset), &tld)
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
		for _, pi := range tld.IpGroups {
			pis = append(pis, pi)
		}
		if tld.Resultinfo.ReturnedCount+tld.Resultinfo.Offset >= tld.Resultinfo.TotalCount {
			return pis, res, nil
		} else {
			offset = tld.Resultinfo.ReturnedCount + tld.Resultinfo.Offset
			continue
		}
	}
	return pis, res, err
}

// Delete
//
func (s *DirectionalPoolsService) DeleteDirectionalGeoPool(dp, acct string) (*Response, error) {
	path := DirectionalPoolPath(acct, "geo", dp)
	return s.client.delete(path, nil)
}

// DeleteDirectionalPool deletes a task.
//
func (s *DirectionalPoolsService) DeleteDirectionalIPPool(dp, acct string) (*Response, error) {
	path := DirectionalPoolPath(acct, "ip", dp)

	return s.client.delete(path, nil)
}
