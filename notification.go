package udnssdk

import (
	"fmt"
	"log"
	"time"
)

// NotificationsService manages Probes
type NotificationsService struct {
	client *Client
}

// NotificationDTO manages notifications
type NotificationDTO struct {
	Email       string                   `json:"email"`
	PoolRecords []NotificationPoolRecord `json:"poolRecords"`
}

// NotificationPoolRecord does things unknown
type NotificationPoolRecord struct {
	PoolRecord   string              `json:"poolRecord"`
	Notification NotificationInfoDTO `json:"notification"`
}

// NotificationInfoDTO does things unknown
type NotificationInfoDTO struct {
	Probe     bool `json:"probe"`
	Record    bool `json:"record"`
	Scheduled bool `json:"scheduled"`
}

// NotificationListDTO does things unknown
type NotificationListDTO struct {
	Notifications []NotificationDTO `json:"notifications"`
	Queryinfo     QueryInfo         `json:"queryInfo"`
	Resultinfo    ResultInfo        `json:"resultInfo"`
}

// NotificationKey collects the identifiers of an Notification
type NotificationKey struct {
	Zone string
	Type string
	Name string
	GUID string
}

// RRSetKey generates the RRSetKey for the NotificationKey
func (n *NotificationKey) RRSetKey() *RRSetKey {
	return &RRSetKey{
		Zone: n.Zone,
		Type: n.Type,
		Name: n.Name,
	}
}

// URI generates the URI for a probe
func (n *NotificationKey) URI() string {
	return fmt.Sprintf("%s/%s", n.RRSetKey().NotificationsURI(), n.GUID)
}

// Select requests all notifications by RRSetKey and optional query, using pagination and error handling
func (s *NotificationsService) Select(r RRSetKey, query string) ([]NotificationDTO, *Response, error) {
	// TODO: Sane Configuration for timeouts / retries
	maxerrs := 5
	waittime := 5 * time.Second

	// init accumulators
	pis := []NotificationDTO{}
	errcnt := 0
	offset := 0

	for {
		reqNotifications, ri, res, err := s.SelectWithOffset(r, query, offset)
		if err != nil {
			if res.StatusCode >= 500 {
				errcnt = errcnt + 1
				if errcnt < maxerrs {
					time.Sleep(waittime)
					continue
				}
			}
			return pis, res, err
		}

		log.Printf("[DEBUG] ResultInfo: %+v\n", ri)
		for _, pi := range reqNotifications {
			pis = append(pis, pi)
		}
		if ri.ReturnedCount+ri.Offset >= ri.TotalCount {
			return pis, res, nil
		}
		offset = ri.ReturnedCount + ri.Offset
		continue
	}
}

// SelectWithOffset requests list of notifications by RRSetKey, query and offset, also returning list metadata, the actual response, or an error
func (s *NotificationsService) SelectWithOffset(r RRSetKey, query string, offset int) ([]NotificationDTO, ResultInfo, *Response, error) {
	var tld NotificationListDTO

	uri := r.NotificationsQueryURI(query, offset)
	res, err := s.client.get(uri, &tld)

	log.Printf("DEBUG - ResultInfo: %+v\n", tld.Resultinfo)
	pis := []NotificationDTO{}
	for _, pi := range tld.Notifications {
		pis = append(pis, pi)
	}
	return pis, tld.Resultinfo, res, err
}

// Find requests a notification by NotificationKey,returning the actual response, or an error
func (s *NotificationsService) Find(n NotificationKey) (NotificationInfoDTO, *Response, error) {
	var t NotificationInfoDTO
	res, err := s.client.get(n.URI(), &t)
	return t, res, err
}

// Create requests creation of an event by RRSetKey, with provided NotificationInfoDTO, returning actual response or an error
func (s *NotificationsService) Create(r RRSetKey, ev NotificationInfoDTO) (*Response, error) {
	return s.client.post(r.NotificationsURI(), ev, nil)
}

// Update requests update of an event by NotificationKey, with provided NotificationInfoDTO, returning the actual response or an error
func (s *NotificationsService) Update(n NotificationKey, ev NotificationInfoDTO) (*Response, error) {
	return s.client.put(n.URI(), ev, nil)
}

// Delete requests deletion of an event by NotificationKey, returning the actual response or an error
func (s *NotificationsService) Delete(n NotificationKey) (*Response, error) {
	return s.client.delete(n.URI(), nil)
}
