package udnssdk

import (
	"fmt"
	"log"
	"time"
)

/* Manages Events */

type EventInfoDTO struct {
	Id         string    `json:"id"`
	PoolRecord string    `json:"poolRecord"`
	EventType  string    `json:"type"`
	Start      time.Time `json:"start"`
	Repeat     string    `json:"repeat"`
	End        time.Time `json:"end"`
	Notify     string    `json:"notify"`
}

type EventInfoListDTO struct {
	Events     []EventInfoDTO `json:"events"`
	Queryinfo  QueryInfo      `json:"queryInfo"`
	Resultinfo ResultInfo     `json:"resultInfo"`
}

func EventPath(zone, typ, name, guid string) string {
	if guid == "" {
		return fmt.Sprintf("zones/%s/rrsets/%s/%s/events", zone, typ, name)
	} else {

		return fmt.Sprintf("zones/%s/rrsets/%s/%s/events/%s", zone, typ, name, guid)
	}
}

// List Event
func (s *SBTCService) ListEvents(query, name, typ, zone string) ([]EventInfoDTO, *Response, error) {
	offset := 0

	reqStr := EventPath(zone, typ, name, "")
	log.Printf("DEBUG - ListEvents: %s\n", reqStr)
	var tld EventInfoListDTO
	//wrappedEvents := []Event{}

	res, err := s.client.get(reqStr, &tld)
	pis := []EventInfoDTO{}
	if query != "" {
		reqStr = fmt.Sprintf("%s?sort=NAME&query=%s&offset=", reqStr, query)
	} else {
		reqStr = fmt.Sprintf("%s?offset=", reqStr)
	}
	// TODO: Sane Configuration for timeouts / retries
	timeout := 5
	waittime := 5 * time.Second
	errcnt := 0
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
		for _, pi := range tld.Events {
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

// Get a Event.
func (s *SBTCService) GetEvent(name, typ, zone, guid string) (EventInfoDTO, *Response, error) {
	reqStr := EventPath(zone, typ, name, guid)
	var t EventInfoDTO
	res, err := s.client.get(reqStr, &t)
	return t, res, err
}

// Create a Event
func (s *SBTCService) CreateEvent(name, typ, zone string, ev EventInfoDTO) (*Response, error) {
	reqStr := EventPath(zone, typ, name, "")
	var retval interface{}
	res, err := s.client.post(reqStr, ev, &retval)

	return res, err
}

// Update
func (s *SBTCService) UpdateEvent(name, typ, zone, guid string, ev EventInfoDTO) (*Response, error) {
	reqStr := EventPath(zone, typ, name, guid)
	var retval interface{}
	res, err := s.client.put(reqStr, ev, &retval)

	return res, err
}

// DeleteEvent
//
func (s *SBTCService) DeleteEvent(name, typ, zone, guid string) (*Response, error) {
	path := EventPath(zone, typ, name, guid)
	return s.client.delete(path, nil)
}
