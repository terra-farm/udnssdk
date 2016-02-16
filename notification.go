package udnssdk

import (
	"fmt"
	"log"
	"time"
)

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

// NotificationPath generates a URI by zone, type & guid
func notificationPath(zone, typ, name, guid string) string {
	if guid == "" {
		return fmt.Sprintf("zones/%s/rrsets/%s/%s/notifications", zone, typ, name)
	}
	return fmt.Sprintf("zones/%s/rrsets/%s/%s/notifications/%s", zone, typ, name, guid)
}

// ListAllNotifications finds all notification by name, type & zone, with optional query
func (s *SBTCService) ListAllNotifications(query, name, typ, zone string) ([]NotificationDTO, *Response, error) {
	// TODO: Sane Configuration for timeouts / retries
	maxerrs := 5
	waittime := 5 * time.Second

	// init accumulators
	pis := []NotificationDTO{}
	errcnt := 0
	offset := 0

	for {
		reqNotifications, ri, res, err := s.ListNotifications(query, name, typ, zone, offset)
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

func notificationQueryPath(zone, typ, name, query string, offset int) string {
	uri := notificationPath(zone, typ, name, "")
	if query != "" {
		uri = fmt.Sprintf("%s?sort=NAME&query=%s&offset=%d", uri, query, offset)
	} else {
		uri = fmt.Sprintf("%s?offset=%d", uri, offset)
	}
	return uri
}

// ListNotifications for things
func (s *SBTCService) ListNotifications(query, name, typ, zone string, offset int) ([]NotificationDTO, ResultInfo, *Response, error) {
	var tld NotificationListDTO

	uri := notificationQueryPath(zone, typ, name, query, offset)
	res, err := s.client.get(uri, &tld)

	log.Printf("DEBUG - ResultInfo: %+v\n", tld.Resultinfo)
	pis := []NotificationDTO{}
	for _, pi := range tld.Notifications {
		pis = append(pis, pi)
	}
	return pis, tld.Resultinfo, res, err
}

// GetNotification returns a notification by name, type, zone & guid
func (s *SBTCService) GetNotification(name, typ, zone, email string) (NotificationDTO, *Response, error) {
	reqStr := notificationPath(zone, typ, name, email)
	var t NotificationDTO
	res, err := s.client.get(reqStr, &t)
	return t, res, err
}

// CreateNotification creates a notification by name, type & zone, with the NotificationDTO ev
func (s *SBTCService) CreateNotification(name, typ, zone, email string, ev NotificationDTO) (*Response, error) {
	reqStr := notificationPath(zone, typ, name, "")
	var ignored interface{}
	return s.client.post(reqStr, ev, &ignored)
}

// UpdateNotification updates a notification by name, type, zone & email, with NotificationDTO ev
func (s *SBTCService) UpdateNotification(name, typ, zone, email string, ev NotificationDTO) (*Response, error) {
	reqStr := notificationPath(zone, typ, name, email)
	var ignored interface{}
	return s.client.put(reqStr, ev, &ignored)
}

// DeleteNotification deletes a notification by name, type, zone & email
func (s *SBTCService) DeleteNotification(name, typ, zone, email string) (*Response, error) {
	path := notificationPath(zone, typ, name, email)
	return s.client.delete(path, nil)
}
