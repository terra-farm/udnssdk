package udnssdk

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// ZonesService provides access to the zones resources
type ZonesService struct {
	client *Client
}

// Zone wraps a response (or item) from:
//   /zone         GET
//   /zone/{name}  GET
// AKA "Zone DTO"
type Zone struct {
	Properties         ZoneProperties         `json:"properties"`
	RestrictIPList     []ZoneRestrictIP       `json:"restrictIPList"`
	PrimaryNameServers ZonePrimaryNameServers `json:"primaryNameServers"`
	OriginalZoneName   string                 `json:"originalZoneName"`
	RegistrarInfo      ZoneRegistrarInfo
	TSig               ZoneTSig            `json:"tsig"`
	NotifyAddresses    []ZoneNotifyAddress `json:"notifyAddresses"`
}

// ZoneCreate wraps a requests to:
//   /zone         POST
//   /zone/{name}  PUT
// AKA "Zone Create DTO"
type ZoneCreate struct {
	Properties      ZoneProperties      `json:"properties"`
	PrimaryCreate   ZonePrimaryCreate   `json:"primaryCreateInfo"`
	SecondaryCreate ZoneSecondaryCreate `json:"secondaryCreateInfo"`
	AliasCreateInfo ZoneAliasCreate     `json:"aliasCreateInfo"`
}

// ZoneProperties wraps the properties value of Zone and represents attributes common to all zones
// All values will be ignored if present in a Create or Update
// AKA "Zone Properties DTO"
type ZoneProperties struct {
	Name                 ZoneKey   `json:"name"`
	AccountName          string    `json:"accountName"`
	Owner                string    `json:"owner"`
	Type                 string    `json:"type"` // "ALIAS", "PRIMARY", "SECONDARY"
	RecordCount          int       `json:"recordCount"`
	DNSSecStatus         string    `json:"dnssecStatus"`         // "SIGNED", "UNSIGNED"
	LastModifiedDateTime time.Time `json:"lastModifiedDateTime"` // ISO8601/RFC3399
}

// ZonePrimaryCreate wraps the primaryCreateInfo value of the ZoneCreate
// AKA "Primary Zone DTO"
type ZonePrimaryCreate struct {
	ForceImport      bool                `json:"forceImport"`
	CreateType       string              `json:"createType"` // "NEW", "COPY", "TRANSFER", "UPLOAD"
	NameServer       ZoneNameServerIP    `json:"nameServer"`
	OriginalZoneName string              `json:"originalZoneName"`
	RestrictIPList   []ZoneRestrictIP    `json:"restrictIPList"`
	TSig             ZoneTSig            `json:"tsig"`
	NotifyAddresses  []ZoneNotifyAddress `json:"notifyAddresses"`
	Inheirit         string              `json:"inheirit"` // "ALL", "NONE", "IP_RANGE", "NOTIFY_IP", "TSIG", joined by ","
}

// ZoneRestrictIP wraps the restrictIPList value of a Zone and ZonePrimaryCreate
// AKA "Restrict IP DTOs"
// cf IPAddrRange
type ZoneRestrictIP struct {
	Start   string `json:"startIP"`
	End     string `json:"endIP"`
	CIDR    string `json:"cidr"`
	Address string `json:"singleIP"`
	Comment string `json:"comment"`
}

// ZoneSecondaryCreate wraps the secondaryCreateInfo value of ZoneCreate
// AKA "Secondary Zone DTO"
type ZoneSecondaryCreate struct {
	PrimaryNameServers ZonePrimaryNameServers `json:"primaryNameServers"`
}

// ZoneAliasCreate wraps the aliasCreateInfo value of ZoneCreate
// AKA "Alias Zone DTO"
type ZoneAliasCreate struct {
	OriginalZoneName string `json:"originalZoneName"`
}

// ZonePrimaryNameServers wraps the primaryNameServers value of a Zone and ZoneCreate
// AKA "Name Server IP List DTO"
type ZonePrimaryNameServers struct {
	NameServerIPList ZoneNameServerIPList `json:"nameServerIpList"`
}

// ZoneNameServerIPList yes, this is really how the data is specified
type ZoneNameServerIPList struct {
	NameServerIP1 ZoneNameServerIP `json:"nameServerIp1"`
	NameServerIP2 ZoneNameServerIP `json:"nameServerIp2"`
	NameServerIP3 ZoneNameServerIP `json:"nameServerIp3"`
}

// ZoneNameServerIP wraps the nameServerIp* values or ZoneNameServerIPList
type ZoneNameServerIP struct {
	IP           string `json:"ip"`
	TSigKey      string `json:"tsigKey"`
	TSigKeyValue string `json:"tsigKeyValue"`
}

// ZoneRegistrarInfo wraps the registrarInfo of Zone
// AKA "Registrar Info DTO"
type ZoneRegistrarInfo struct {
	Registrar       string                   `json:"registrar"`
	WhoisExpiration time.Time                `json:"whoisExpiration"`
	NameServers     ZoneRegistrarNameServers `json:"nameServers"`
}

// ZoneRegistrarNameServers wraps the nameServers value of ZoneRegistrarInfo
type ZoneRegistrarNameServers struct {
	Ok        []string `json:"ok"`
	Unknown   []string `json:"unknown"`
	Missing   []string `json:"missing"`
	Incorrect []string `json:"incorrect"`
}

// ZoneTSig was the tsig value of Zone and ZoneCreate
// AKA "TSIG DTO"
type ZoneTSig struct {
	KeyName     string `json:"tsigKeyName"`
	KeyValue    string `json:"tsigKeyValue"`
	Description string `json:"description"`
	Algorithm   string `json:"tsigAlgorithm"`
}

// ZoneNotifyAddress wraps the notifyAddress value of Zone and ZonePrimaryCreate
// AKA "Notify Address DTO"
type ZoneNotifyAddress struct {
	NotifyAddress string `json:"notifyAddress"`
	Description   string `json:"description"`
}

// ZoneKey is the key for an UltraDNS zone
type ZoneKey string

// URI generates the URI for a task
func (z ZoneKey) URI() string {
	return fmt.Sprintf("zones/%s", z)
}

// ZonesQueryURI generates the query URI for the zone collection given a query and offset
func ZonesQueryURI(query string, offset int) string {
	sort := "NAME"
	limit := 100
	reverse := false
	if query != "" {
		return fmt.Sprintf("zones?q=%v&offset=%d&limit=%d&sort=%v&reverse=%v", query, offset, limit, sort, reverse)
	}
	return fmt.Sprintf("zones?offset=%d&limit=%d&sort=%v&reverse=%v", offset, limit, sort, reverse)
}

// Select requests all zones, with pagination
func (s *ZonesService) Select(query string) ([]Zone, error) {
	// TODO: Sane Configuration for timeouts / retries
	maxerrs := 5
	waittime := 5 * time.Second

	// init accumulators
	dtos := []Zone{}
	offset := 0
	errcnt := 0

	for {
		reqDtos, ri, res, err := s.SelectWithOffset(query, offset)
		if err != nil {
			if res != nil && res.StatusCode >= 500 {
				errcnt = errcnt + 1
				if errcnt < maxerrs {
					time.Sleep(waittime)
					continue
				}
			}
			return dtos, err
		}

		log.Printf("[DEBUG] ResultInfo: %+v\n", ri)
		for _, d := range reqDtos {
			dtos = append(dtos, d)
		}
		if ri.ReturnedCount+ri.Offset >= ri.TotalCount {
			return dtos, nil
		}
		offset = ri.ReturnedCount + ri.Offset
		continue
	}
}

// ZoneList wraps a response from:
//   /zone?q=   GET
// AKA "Zone List DTO"
type ZoneList struct {
	Zones      []Zone     `json:"zones"`
	QueryInfo  QueryInfo  `json:"queryInfo"`
	ResultInfo ResultInfo `json:"resultInfo"`
}

// SelectWithOffset request zones by query & offset, list them also returning list metadata, the actual response, or an error
func (s *ZonesService) SelectWithOffset(query string, offset int) ([]Zone, ResultInfo, *http.Response, error) {
	var zl ZoneList

	uri := ZonesQueryURI(query, offset)
	res, err := s.client.get(uri, &zl)

	zs := []Zone{}
	for _, z := range zl.Zones {
		zs = append(zs, z)
	}
	return zs, zl.ResultInfo, res, err
}

// Find Get the properties of a zone.
func (s *ZonesService) Find(z ZoneKey) (Zone, *http.Response, error) {
	var zv Zone
	res, err := s.client.get(z.URI(), &zv)
	return zv, res, err
}

// Create creates a zone with val
func (s *ZonesService) Create(z ZoneKey, val ZoneCreate) (*http.Response, error) {
	var ignored interface{}
	return s.client.post(z.URI(), val, &ignored)
}

// Update updates a Zone with the provided val
// Cannot be used to:
// - update an alias
// - specify primary name servers for a primary zone
// - specify IPs, TSig, Notify addresses for a secondary zone
func (s *ZonesService) Update(z ZoneKey, val ZoneCreate) (*http.Response, error) {
	var ignored interface{}
	return s.client.put(z.URI(), val, &ignored)
}

// Convert: secondary => primary
// Unalias: alias     => primary

// Delete requests deletions
func (s *ZonesService) Delete(z ZoneKey) (*http.Response, error) {
	return s.client.delete(z.URI(), nil)
}
