package udnssdk

import (
	"fmt"
	"log"
	"time"
)

/* Manages Notifications */
type NotificationDTO struct {
	Email       string                   `json:"email"`
	PoolRecords []NotificationPoolRecord `json:"poolRecords"`
}

type NotificationPoolRecord struct {
	PoolRecord   string              `json:"poolRecord"`
	Notification NotificationInfoDTO `json:"notification"`
}

type NotificationInfoDTO struct {
	Probe     bool `json:"probe"`
	Record    bool `json:"record"`
	Scheduled bool `json:"scheduled"`
}

type NotificationListDTO struct {
	Notifications []NotificationDTO `json:"notifications"`
	Queryinfo     QueryInfo         `json:"queryInfo"`
	Resultinfo    ResultInfo        `json:"resultInfo"`
}

func NotificationPath(zone, typ, name, guid string) string {
	if guid == "" {
		return fmt.Sprintf("zones/%s/rrsets/%s/%s/notifications", zone, typ, name)
	} else {

		return fmt.Sprintf("zones/%s/rrsets/%s/%s/notifications/%s", zone, typ, name, guid)
	}
}

// List Notification
func (s *SBTCService) ListNotifications(query, name, typ, zone string) ([]NotificationDTO, *Response, error) {
	offset := 0

	reqStr := NotificationPath(zone, typ, name, "")
	log.Printf("DEBUG - ListNotifications: %s\n", reqStr)
	var tld NotificationListDTO
	//wrappedNotifications := []Notification{}

	res, err := s.client.get(reqStr, &tld)
	pis := []NotificationDTO{}
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
		for _, pi := range tld.Notifications {
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

// Get a Notification.
func (s *SBTCService) GetNotification(name, typ, zone, guid string) (NotificationInfoDTO, *Response, error) {
	reqStr := NotificationPath(zone, typ, name, guid)
	var t NotificationInfoDTO
	res, err := s.client.get(reqStr, &t)
	return t, res, err
}

// Create a Notification
func (s *SBTCService) CreateNotification(name, typ, zone string, ev NotificationInfoDTO) (*Response, error) {
	reqStr := NotificationPath(zone, typ, name, "")
	var retval interface{}
	res, err := s.client.post(reqStr, ev, &retval)

	return res, err
}

// Update
func (s *SBTCService) UpdateNotification(name, typ, zone, guid string, ev NotificationInfoDTO) (*Response, error) {
	reqStr := NotificationPath(zone, typ, name, guid)
	var retval interface{}
	res, err := s.client.put(reqStr, ev, &retval)

	return res, err
}

// DeleteNotification
//
func (s *SBTCService) DeleteNotification(name, typ, zone, guid string) (*Response, error) {
	path := NotificationPath(zone, typ, name, guid)
	return s.client.delete(path, nil)
}
