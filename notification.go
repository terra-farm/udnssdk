package udnssdk

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// NotificationsService manages Probes
type NotificationsService struct {
	client *Client
}

// Notification manages notifications
type Notification struct {
	Email       string                   `json:"email"`
	PoolRecords []NotificationPoolRecord `json:"poolRecords"`
}

// NotificationPoolRecord does things unknown
type NotificationPoolRecord struct {
	PoolRecord   string           `json:"poolRecord"`
	Notification NotificationInfo `json:"notification"`
}

// NotificationInfo does things unknown
type NotificationInfo struct {
	Probe     bool `json:"probe"`
	Record    bool `json:"record"`
	Scheduled bool `json:"scheduled"`
}

// NotificationList does things unknown
type NotificationList struct {
	Notifications []Notification `json:"notifications"`
	Queryinfo     QueryInfo      `json:"queryInfo"`
	Resultinfo    ResultInfo     `json:"resultInfo"`
}

// NotificationKey collects the identifiers of an Notification
type NotificationKey struct {
	Zone  string
	Type  string
	Name  string
	Email string
}

// RRSetKey generates the RRSetKey for the NotificationKey
func (k NotificationKey) RRSetKey() RRSetKey {
	return RRSetKey{
		Zone: k.Zone,
		Type: k.Type,
		Name: k.Name,
	}
}

// URI generates the URI for a probe
func (k NotificationKey) URI() string {
	return fmt.Sprintf("%s/%s", k.RRSetKey().NotificationsURI(), k.Email)
}

// Select requests all notifications by RRSetKey and optional query, using pagination and error handling
func (s *NotificationsService) Select(k RRSetKey, query string) ([]Notification, *http.Response, error) {
	// TODO: Sane Configuration for timeouts / retries
	maxerrs := 5
	waittime := 5 * time.Second

	// init accumulators
	pis := []Notification{}
	errcnt := 0
	offset := 0

	for {
		reqNotifications, ri, res, err := s.SelectWithOffset(k, query, offset)
		if err != nil {
			if res != nil && res.StatusCode >= 500 {
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
func (s *NotificationsService) SelectWithOffset(k RRSetKey, query string, offset int) ([]Notification, ResultInfo, *http.Response, error) {
	var tld NotificationList

	uri := k.NotificationsQueryURI(query, offset)
	res, err := s.client.get(uri, &tld)

	log.Printf("DEBUG - ResultInfo: %+v\n", tld.Resultinfo)
	pis := []Notification{}
	for _, pi := range tld.Notifications {
		pis = append(pis, pi)
	}
	return pis, tld.Resultinfo, res, err
}

// Find requests a notification by NotificationKey,returning the actual response, or an error
func (s *NotificationsService) Find(k NotificationKey) (Notification, *http.Response, error) {
	var t Notification
	res, err := s.client.get(k.URI(), &t)
	return t, res, err
}

// Create requests creation of an event by RRSetKey, with provided NotificationInfo, returning actual response or an error
func (s *NotificationsService) Create(k NotificationKey, n Notification) (*http.Response, error) {
	return s.client.post(k.URI(), n, nil)
}

// Update requests update of an event by NotificationKey, with provided NotificationInfoO, returning the actual response or an error
func (s *NotificationsService) Update(k NotificationKey, n Notification) (*http.Response, error) {
	return s.client.put(k.URI(), n, nil)
}

// Delete requests deletion of an event by NotificationKey, returning the actual response or an error
func (s *NotificationsService) Delete(k NotificationKey) (*http.Response, error) {
	return s.client.delete(k.URI(), nil)
}
