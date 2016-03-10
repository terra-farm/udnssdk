package udnssdk

import (
	"fmt"
	"log"
	"time"
)

// EventInfoDTO wraps an event's info response
type EventInfoDTO struct {
	ID         string    `json:"id"`
	PoolRecord string    `json:"poolRecord"`
	EventType  string    `json:"type"`
	Start      time.Time `json:"start"`
	Repeat     string    `json:"repeat"`
	End        time.Time `json:"end"`
	Notify     string    `json:"notify"`
}

// EventInfoListDTO wraps a list of event info and list metadata, from an index request
type EventInfoListDTO struct {
	Events     []EventInfoDTO `json:"events"`
	Queryinfo  QueryInfo      `json:"queryInfo"`
	Resultinfo ResultInfo     `json:"resultInfo"`
}

// EventPath generates a URI by zone, type, name & guid
func EventPath(zone, typ, name, guid string) string {
	if guid == "" {
		return fmt.Sprintf("zones/%s/rrsets/%s/%s/events", zone, typ, name)
	}
	return fmt.Sprintf("zones/%s/rrsets/%s/%s/events/%s", zone, typ, name, guid)
}

// ListAllEvents requests all events, using pagination and error handling
func (s *SBTCService) ListAllEvents(query, name, typ, zone string) ([]EventInfoDTO, error) {
	// TODO: Sane Configuration for timeouts / retries
	maxerrs := 5
	waittime := 5 * time.Second

	// init accumulators
	pis := []EventInfoDTO{}
	offset := 0
	errcnt := 0

	for {
		reqEvents, ri, res, err := s.ListEvents(query, name, typ, zone, offset)
		if err != nil {
			if res.StatusCode >= 500 {
				errcnt = errcnt + 1
				if errcnt < maxerrs {
					time.Sleep(waittime)
					continue
				}
			}
			return pis, err
		}

		log.Printf("ResultInfo: %+v\n", ri)
		for _, pi := range reqEvents {
			pis = append(pis, pi)
		}
		if ri.ReturnedCount+ri.Offset >= ri.TotalCount {
			return pis, nil
		}
		offset = ri.ReturnedCount + ri.Offset
		continue
	}
}

func eventQueryPath(zone, typ, name, query string, offset int) string {
	uri := EventPath(zone, typ, name, "")
	if query != "" {
		uri = fmt.Sprintf("%s?sort=NAME&query=%s&offset=%d", uri, query, offset)
	} else {
		uri = fmt.Sprintf("%s?offset=%d", uri, offset)
	}
	return uri
}

// ListEvents requests list of events by query, name, type & zone, and offset, also returning list metadata, the actual response, or an error
func (s *SBTCService) ListEvents(query, name, typ, zone string, offset int) ([]EventInfoDTO, ResultInfo, *Response, error) {
	var tld EventInfoListDTO

	uri := eventQueryPath(zone, typ, name, query, offset)
	res, err := s.client.get(uri, &tld)

	pis := []EventInfoDTO{}
	for _, pi := range tld.Events {
		pis = append(pis, pi)
	}
	return pis, tld.Resultinfo, res, err
}

// GetEvent requests an event by name, type, zone & guid, also returning the actual response, or an error
func (s *SBTCService) GetEvent(name, typ, zone, guid string) (EventInfoDTO, *Response, error) {
	reqStr := EventPath(zone, typ, name, guid)
	var t EventInfoDTO
	res, err := s.client.get(reqStr, &t)
	return t, res, err
}

// CreateEvent requests creation of an event by name, type, zone, with provided event-info, returning actual response or an error
func (s *SBTCService) CreateEvent(name, typ, zone string, ev EventInfoDTO) (*Response, error) {
	reqStr := EventPath(zone, typ, name, "")
	return s.client.post(reqStr, ev, nil)
}

// UpdateEvent requests update of an event by name, type, zone & guid, withprovided event-info, returning the actual response or an error
func (s *SBTCService) UpdateEvent(name, typ, zone, guid string, ev EventInfoDTO) (*Response, error) {
	reqStr := EventPath(zone, typ, name, guid)
	return s.client.put(reqStr, ev, nil)
}

// DeleteEvent requests deletion of an event by name, type, zone & guid, returning the actual response or an error
func (s *SBTCService) DeleteEvent(name, typ, zone, guid string) (*Response, error) {
	path := EventPath(zone, typ, name, guid)
	return s.client.delete(path, nil)
}
