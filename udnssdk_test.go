package udnssdk

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	testUsername    = os.Getenv("ULTRADNS_USERNAME")
	testPassword    = os.Getenv("ULTRADNS_PASSWORD")
	testDomain      = os.Getenv("ULTRADNS_DOMAIN")
	testHostname    = os.Getenv("ULTRADNS_TEST_HOSTNAME")
	testIP1         = os.Getenv("ULTRADNS_TEST_IP1")
	testIP2         = os.Getenv("ULTRADNS_TEST_IP2")
	testBaseURL     = os.Getenv("ULTRADNS_BASEURL")
	testQuery       = os.Getenv("ULTRADNS_TEST_QUERY")
	testProbeType   = os.Getenv("ULTRADNS_TEST_PROBE_TYPE")
	testProbeName   = os.Getenv("ULTRADNS_TEST_PROBE_NAME")
	testProbeDomain = os.Getenv("ULTRADNS_TEST_PROBE_DOMAIN")

	testIPDPoolName    = "testipdpool"
	testIPDPoolAddress = "127.0.0.1"
	testIPDPoolDescr   = "A Test IP Directional Pool Group"
	testIPAddrDTO      = IPAddrDTO{Address: "127.0.0.1"}
	testIPDPool        = AccountLevelIPDirectionalGroupDTO{Name: "testippool", Description: "An IP Test Pool", Ips: []IPAddrDTO{IPAddrDTO{Address: "127.0.0.1"}}}
	testGeoDPool       = AccountLevelGeoDirectionalGroupDTO{Name: "testgeopool", Description: "A test geo pool", Codes: []string{"US, UK"}}
	testGeoDPoolName   = "testgeodpool"
	testGeoDPoolDescr  = "A Test Geo Directional Pool Group"
	testGeoDPoolCodes  = []string{"US", "UK"}

	envenableAccountTests         = os.Getenv("ULTRADNS_ENABLE_ACCOUNT_TESTS")
	envenableRRSetTests           = os.Getenv("ULTRADNS_ENABLE_RRSET_TESTS")
	envenableProbeTests           = os.Getenv("ULTRADNS_ENABLE_PROBE_TESTS")
	envenableChanges              = os.Getenv("ULTRADNS_ENABLE_CHANGES")
	envenableDirectionalPoolTests = os.Getenv("ULTRADNS_ENABLE_DPOOL_TESTS")
	enableAccountTests            = true
	enableRRSetTests              = true
	enableProbeTests              = true
	enableChanges                 = true
	enableDirectionalPoolTests    = false
	testProfile                   = `{"@context": "http://schemas.ultradns.com/RDPool.jsonschema", "order": "ROUND_ROBIN","description": "T. migratorius"}`
	testProfile2                  = `{"@context": "http://schemas.ultradns.com/RDPool.jsonschema", "order": "RANDOM","description": "T. migratorius"}`

	testClient   *Client
	testAccounts []Account
)

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())

	if testUsername == "" {
		log.Printf("Please configure ULTRADNS_USERNAME.\n")
		os.Exit(1)
	}
	if testPassword == "" {
		log.Printf("Please configure ULTRADNS_PASSWORD.\n")
		os.Exit(1)
	}
	if testDomain == "" {
		log.Printf("Please configure ULTRADNS_DOMAIN.\n")
		os.Exit(1)
	}
	if testHostname == "" {
		log.Printf("Please configure ULTRADNS_TEST_HOSTNAME.\n")
		os.Exit(1)
	}
	if testBaseURL == "" {
		testBaseURL = DefaultTestBaseURL
	}

	if testIP1 == "" {
		testIP1 = "54.86.13.225"
	}
	if testIP2 == "" {
		testIP2 = fmt.Sprintf("54.86.13.%d", (rand.Intn(254) + 1))
	}
	if testQuery == "" {
		testQuery = "nexus"
	}

	if testProbeName == "" || testProbeType == "" {
		testProbeName = "nexus2"
		testProbeType = "A"
	}
	if testProbeDomain == "" {
		testProbeDomain = testDomain
	}
	if envenableAccountTests == "false" || envenableAccountTests == "0" {
		enableAccountTests = false
	} else if envenableAccountTests == "true" || envenableAccountTests == "1" {
		enableAccountTests = true
	}

	if envenableRRSetTests == "false" || envenableRRSetTests == "0" {
		enableRRSetTests = false
	} else if envenableRRSetTests == "true" || envenableRRSetTests == "1" {
		enableRRSetTests = true
	}
	// TODO: I need a better way of handling this.
	/*
		if envenableFUDGETests == "false" || envenableFUDGETests == "0" {
			enableFUDGETests = false
		} else if envenableFUDGETests == "true" || envenableFUDGETests == "1" {
			enableFUDGETests = true
		}
	*/

	if envenableProbeTests == "false" || envenableProbeTests == "0" {
		enableProbeTests = false
	} else if envenableProbeTests == "true" || envenableProbeTests == "1" {
		enableProbeTests = true
	}

	if envenableChanges == "false" || envenableChanges == "0" {
		enableChanges = false
	} else if envenableChanges == "true" || envenableChanges == "1" {
		enableChanges = true
	}

	if envenableDirectionalPoolTests == "false" || envenableDirectionalPoolTests == "0" {
		enableDirectionalPoolTests = false
	} else if envenableDirectionalPoolTests == "true" || envenableDirectionalPoolTests == "1" {
		enableDirectionalPoolTests = true
	}

	testAccounts = nil
	os.Exit(m.Run())
}

func Test_CreateClient(t *testing.T) {
	log.Printf("Creating Client...\n")
	var err error
	testClient, err = NewClient(testUsername, testPassword, testBaseURL)

	if testClient == nil || err != nil {
		t.Fail()
		log.Fatalf("Could not create client - %+v\n", err)
		os.Exit(1)
	}
	t.Logf("Client created successfully.\n")
}

func Test_GetRRSetsPre(t *testing.T) {
	if testClient == nil {
		t.Fatalf("TestClient Not Defined?\n")
	}
	if !enableRRSetTests {
		t.SkipNow()
	}

	rrsets, resp, err := testClient.RRSets.GetRRSets(testDomain, testHostname, "ANY")

	t.Logf("GetRRSets(%s, %s, \"ANY\")", testDomain, testHostname)
	t.Logf("RRSets: %+v\n", rrsets)
	t.Logf("Response: %+v\n", resp.Response)
	if err != nil {
		if resp.Response.StatusCode == 404 {
			t.SkipNow()
		}
		t.Fatal(err)
	}
}

func Test_ListRRSets(t *testing.T) {
	if !enableRRSetTests {
		t.SkipNow()
	}
	rrsets, resp, err := testClient.RRSets.GetRRSets(testDomain, "", "")
	t.Logf("GetRRSets(%s, \"\", \"\")", testDomain)
	t.Logf("RRSets: %+v\n", rrsets)
	t.Logf("Response: %+v\n", resp.Response)
	if err != nil {
		if resp.Response.StatusCode == 404 {
			t.SkipNow()
		}
		t.Fatal(err)
	}
	t.Logf("Checking for profiles...\n")
	for _, rr := range rrsets {
		if rr.Profile != nil {
			typ := rr.Profile.GetType()
			if typ == "" {
				t.Fatal("Could not get type for profile %+v\n", rr.Profile)
			}
			t.Logf("Found Profile %s for %s\n", rr.Profile.GetType(), rr.OwnerName)
			st, er := json.Marshal(rr.Profile)
			t.Logf("Marshal the profile to JSON: %s / %+v", string(st), er)
			t.Logf("Check the Magic Profile: %+v\n", rr.Profile.GetProfileObject())
		}
	}
}

// Create Test
func Test_Create_RRSets(t *testing.T) {

	if !enableRRSetTests {
		t.SkipNow()

	}
	if !enableChanges {
		t.SkipNow()

	}
	t.Logf("Creating %s with %s\n", testHostname, testIP1)
	rr1 := &RRSet{OwnerName: testHostname, RRType: "A", TTL: 300, RData: []string{testIP1}, Profile: &StringProfile{Profile: testProfile}}
	resp, err := testClient.RRSets.CreateRRSet(testDomain, *rr1)
	t.Logf("Response: %+v\n", resp.Response)
	if err != nil {
		t.Fatal(err)
	}
}

// Another Get  Test if it matchs the Ip in IP1

func Test_GetRRSetsMid1(t *testing.T) {

	if !enableRRSetTests {
		t.SkipNow()

	}

	rrsets, resp, err := testClient.RRSets.GetRRSets(testDomain, testHostname, "ANY")

	t.Logf("GetRRSets(%s, %s, \"ANY\")", testDomain, testHostname)
	t.Logf("RRSets: %+v\n", rrsets)
	t.Logf("Response: %+v\n", resp.Response)
	if err != nil {
		t.Fatal(err)
	}
	// Do the test v IP1 here
	if rrsets[0].RData[0] != testIP1 {
		t.Fatalf("RData[0]\"%s\" != testIP1\"%s\"", rrsets[0].RData[0], testIP1)
	}
}

//Update Test
func Test_Update_RRSets(t *testing.T) {

	if !enableRRSetTests {
		t.SkipNow()

	}
	if !enableChanges {
		t.SkipNow()

	}

	t.Logf("Updating %s to %s\n", testHostname, testIP2)
	rr2 := &RRSet{OwnerName: testHostname, RRType: "A", TTL: 300, RData: []string{testIP2}, Profile: &StringProfile{Profile: testProfile2}}
	resp, err := testClient.RRSets.UpdateRRSet(testDomain, *rr2)
	t.Logf("Response: %+v\n", resp.Response)
	if err != nil {
		t.Fatal(err)
	}
} // Another Get Test if it matches the Ip in IP2
func Test_GetRRSetsMid(t *testing.T) {

	if !enableRRSetTests {
		t.SkipNow()

	}

	rrsets, resp, err := testClient.RRSets.GetRRSets(testDomain, testHostname, "ANY")

	t.Logf("GetRRSets(%s, %s, \"ANY\")", testDomain, testHostname)
	t.Logf("RRSets: %+v\n", rrsets)
	t.Logf("Response: %+v\n", resp.Response)
	if err != nil {
		t.Fatal(err)
	}
	// Do the test v IP2 here
	if rrsets[0].RData[0] != testIP2 {
		t.Fatalf("RData[0]\"%s\" != testIP2\"%s\"", rrsets[0].RData[0], testIP2)
	}
	t.Logf("Profile Check: %+v", rrsets[0].Profile.GetProfileObject())
}

// Delete Test
func Test_Delete_RRSets(t *testing.T) {

	if !enableRRSetTests {
		t.SkipNow()

	}
	if !enableChanges {
		t.SkipNow()

	}
	if testHostname == "" || testHostname[0] == '*' || testHostname[0] == '@' || testHostname == "www" || testHostname[0] == '.' {
		t.Fatalf("NO testHostname DEFINED!  DANGER")
		os.Exit(1)
	}
	t.Logf("Deleting %s\n", testHostname)
	t.Logf("Get RRSet for %s\n", testHostname)
	rrsets, resp, err := testClient.RRSets.GetRRSets(testDomain, testHostname, "ANY")
	t.Logf("GetRRSets(%s, %s, \"ANY\")", testDomain, testHostname)
	t.Logf("RRSets: %+v\n", rrsets)
	t.Logf("Response: %+v\n", resp.Response)
	if err != nil {
		t.Fatal(err)
	}
	for i, e := range rrsets {
		t.Logf("Deleting %s  - ( %d ) %+v \n", testHostname, i, e)
		/*		if e.OwnerName != testHostname {
				t.Logf("e.OwnerName(%s) != testHostname(%s).. Resetting..\n", e.OwnerName, testHostname)
				e.OwnerName = testHostname
				t.Logf("NewE: %+v\n", e)
			} */
		if strings.Index(e.RRType, " ") != -1 {
			t.Logf("Stripping RRType\n")
			x := strings.Fields(e.RRType)[0]
			e.RRType = x
		}
		resp, err := testClient.RRSets.DeleteRRSet(testDomain, e)
		t.Logf("Response: %+v\n", resp.Response)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func Test_GetRRSetsPost(t *testing.T) {

	if !enableRRSetTests {
		t.SkipNow()

	}
	rrsets, resp, err := testClient.RRSets.GetRRSets(testDomain, testHostname, "ANY")

	t.Logf("GetRRSets(%s, %s, \"ANY\")", testDomain, testHostname)
	t.Logf("RRSets: %+v\n", rrsets)
	t.Logf("Response: %+v\n", resp.Response)
	if err != nil {
		if resp.Response.StatusCode == 404 {
			return
		}
		t.Fatal(err)
	}
}

func Test_ListTasks(t *testing.T) {
	tasks, resp, err := testClient.Tasks.ListTasks("")
	t.Logf("Tasks: %+v \n", tasks)
	t.Logf("Response: %+v\n", resp.Response)
	if err != nil {
		if resp.Response.StatusCode == 404 {
			t.SkipNow()
		}
		t.Fatal(err)
	}
}

func Test_ListAccountsOfUser(t *testing.T) {

	if !enableAccountTests {
		t.SkipNow()
	}
	accounts, resp, err := testClient.Accounts.GetAccountsOfUser()
	t.Logf("Accounts: %+v \n", accounts)
	t.Logf("Response: %+v\n", resp.Response)
	testAccounts = accounts
	if err != nil {
		if resp.Response.StatusCode == 404 {
			t.SkipNow()
		}
		t.Fatal(err)
	}
}

/*
// TODO: Implement Zones
func TestListZonesOfAccount(t *testing.T) {
	if testAccounts == nil {
		t.Logf("No Accounts Present, skipping...")
		t.SkipNow()
	}
	zones, resp, err := testClient.Accounts.GetAccountsOfUser()
	t.Logf("Zones: %v \n", zones)
	t.Logf("Response: %+v\n", resp.Response)
	testZones = zones
	if err != nil {
		if resp.Response.StatusCode == 404 {
			t.SkipNow()
		}
		t.Fatal(err)
	}
}
*/

func Test_ListDirectionPoolsGeoNoQuery(t *testing.T) {
	if !enableDirectionalPoolTests {
		t.SkipNow()
	}
	if testAccounts == nil {
		t.Logf("No Accounts Present, skipping...")
		t.SkipNow()
	}
	accountName := testAccounts[0].AccountName
	dpools, resp, err := testClient.DirectionalPools.ListDirectionalGeoPools("", accountName)
	t.Logf("Geo Pools: %v \n", dpools)
	t.Logf("Response: %+v\n", resp.Response)
	if err != nil {
		if resp.Response.StatusCode == 404 {
			t.Logf("Error: %+v", err)
			t.SkipNow()
		}
		t.Fatal(err)
	}

	t.SkipNow()
}
func Test_ListDirectionPoolsGeoQuery(t *testing.T) {
	if !enableDirectionalPoolTests {
		t.SkipNow()
	}
	if testAccounts == nil {
		t.Logf("No Accounts Present, skipping...")
		t.SkipNow()
	}
	accountName := testAccounts[0].AccountName
	dpools, resp, err := testClient.DirectionalPools.ListDirectionalGeoPools(testQuery, accountName)
	t.Logf("Geo Pools: %v \n", dpools)
	t.Logf("Response: %+v\n", resp.Response)
	if err != nil {
		if resp.Response.StatusCode == 404 {
			t.Logf("Error: %+v", err)
			t.SkipNow()
		}
		t.Fatal(err)
	}

	t.SkipNow()
}

func Test_ListDirectionalPoolsIPNoQuery(t *testing.T) {
	if !enableDirectionalPoolTests {
		t.SkipNow()
	}
	if testAccounts == nil {
		t.Logf("No Accounts Present, skipping...")
		t.SkipNow()
	}
	accountName := testAccounts[0].AccountName
	dpools, resp, err := testClient.DirectionalPools.ListDirectionalIPPools("", accountName)
	t.Logf("IP Pools: %v \n", dpools)
	t.Logf("Response: %+v\n", resp.Response)
	if err != nil {
		if resp.Response.StatusCode == 404 {
			t.Logf("Error: %+v", err)
			t.SkipNow()
		}
		t.Fatal(err)
	}

	t.SkipNow()
}
func Test_ListDirectionalPoolsIPQuery(t *testing.T) {
	if !enableDirectionalPoolTests {
		t.SkipNow()
	}
	if testAccounts == nil {
		t.Logf("No Accounts Present, skipping...")
		t.SkipNow()
	}
	accountName := testAccounts[0].AccountName
	dpools, resp, err := testClient.DirectionalPools.ListDirectionalIPPools(testQuery, accountName)
	t.Logf("IP Pools: %v \n", dpools)
	t.Logf("Response: %+v\n", resp.Response)
	if err != nil {
		if resp.Response.StatusCode == 404 {
			t.Logf("Error: %+v", err)
			t.SkipNow()
		}
		t.Fatal(err)
	}

	t.SkipNow()
}

// Create Test
func Test_Create_DirectionalPoolIP(t *testing.T) {

	if !enableDirectionalPoolTests {
		t.SkipNow()

	}
	if !enableChanges {
		t.SkipNow()

	}
	if testAccounts == nil {
		t.Logf("No Accounts Present, skipping...")
		t.SkipNow()
	}

	accountName := testAccounts[0].AccountName
	t.Logf("Creating %s with %+v\n", testIPDPool.Name, testIPDPool)
	resp, err := testClient.DirectionalPools.CreateDirectionalIPPool(testIPDPool.Name, accountName, testIPDPool)
	t.Logf("Response: %+v\n", resp.Response)
	if err != nil {
		t.Fatal(err)
	}
}
func Test_Get_DirectionalPoolIP(t *testing.T) {
	if !enableDirectionalPoolTests {
		t.SkipNow()

	}
	if !enableChanges {
		t.SkipNow()

	}
	if testAccounts == nil {
		t.Logf("No Accounts Present, skipping...")
		t.SkipNow()
	}

	accountName := testAccounts[0].AccountName

	dp, resp, err := testClient.DirectionalPools.GetDirectionalIPPool(testIPDPool.Name, accountName)

	t.Logf("Test Get IP DPool Group (%s, %s)\n", testIPDPool.Name, testIPDPool)
	t.Logf("Response: %+v\n", resp.Response)
	t.Logf("DPool: %+v\n", dp)
	if err != nil {
		t.Logf("GetDirectionalPoolIP Error: %+v\n", err)
		if resp.Response.StatusCode == 404 {
			return
		}
		t.Fatal(err)
	}
	dp2, er := json.Marshal(dp)
	t.Logf("DPool Marshalled back: %s - %+v\n", string(dp2), er)

}

func Test_Delete_DirectionalPoolIP(t *testing.T) {
	if !enableDirectionalPoolTests {
		t.SkipNow()

	}
	if !enableChanges {
		t.SkipNow()

	}
	if testAccounts == nil {
		t.Logf("No Accounts Present, skipping...")
		t.SkipNow()
	}

	accountName := testAccounts[0].AccountName
	resp, err := testClient.DirectionalPools.DeleteDirectionalIPPool(testIPDPool.Name, accountName)

	t.Logf("Delete IP DPool Group (%s, %s)\n", testIPDPool.Name, testIPDPool)
	t.Logf("Response: %+v\n", resp.Response)
	if err != nil {
		t.Logf("DeleteDirectionalPoolIP Error: %+v\n", err)
		if resp.Response.StatusCode == 404 {
			return
		}
		t.Fatal(err)
	}
}
func Test_ListProbes(t *testing.T) {
	if !enableProbeTests {
		t.SkipNow()
	}
	probes, resp, err := testClient.SBTCService.ListProbes("", testProbeName, testProbeType, testProbeDomain)
	t.Logf("Probes: %+v \n", probes)
	t.Logf("Response: %+v\n", resp.Response)
	if err != nil {
		if resp.Response.StatusCode == 404 {
			t.Logf("ERROR - %+v", err)
			t.SkipNow()
		}
		t.Fatal(err)
	}
	for i, e := range probes {
		t.Logf("DEBUG - Probe %d Data - %s\n", i, e.Details.data)
		err = e.Details.Populate(e.ProbeType)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("DEBUG - Populated Probe: %+v\n", e)
		/*
			st, er := json.Marshal(e)
			if er != nil {
				t.Errorf("ERROR - Serialization 1 failed! %+v", er)
			}
			t.Logf("DEBUG - Testing Serialization 1: %s\n", string(st))

			st, er = json.Marshal(e.Details)
			if er != nil {
				t.Errorf("ERROR - Serialization 2 failed! %+v", er)
			}
			t.Logf("DEBUG - Testing Serialization 2 : %+v -> %s\n", e.Details, string(st))
		*/
	}
}

func Test_GetProbeAlerts(t *testing.T) {
	if !enableProbeTests {
		t.SkipNow()
	}
	probes, resp, err := testClient.SBTCService.GetProbeAlerts(testProbeName, testProbeType, testProbeDomain)
	t.Logf("Probe Alertss: %+v \n", probes)
	t.Logf("Response: %+v\n", resp.Response)
	if err != nil {
		if resp.Response.StatusCode == 404 {
			t.Logf("ERROR - %+v", err)
			t.SkipNow()
		}
		t.Fatal(err)
	}
	/*
		for i, e := range probes {
			t.Logf("DEBUG - Probe Alert %d Data - %+v\n", i, e)
		}
	*/

}

/* TODO: A full probe test suite.  I'm not really even sure I understand how this
 * works well enough to write one yet.  What is the correct order of operations?
 */

func Test_ListEvents(t *testing.T) {
	if !enableProbeTests {
		t.SkipNow()
	}
	events, resp, err := testClient.SBTCService.ListEvents("", testProbeName, testProbeType, testProbeDomain)
	t.Logf("Events: %+v \n", events)
	t.Logf("Response: %+v\n", resp.Response)
	if err != nil {
		if resp.Response.StatusCode == 404 {
			t.Logf("ERROR - %+v", err)
			t.SkipNow()
		}
		t.Fatal(err)
	}
}

// TODO: Write a full Event test suite.  We do not use these at my firm.

func Test_ListNotifications(t *testing.T) {
	if !enableProbeTests {
		t.SkipNow()
	}
	events, resp, err := testClient.SBTCService.ListNotifications("", testProbeName, testProbeType, testProbeDomain)
	t.Logf("Notifications: %+v \n", events)
	t.Logf("Response: %+v\n", resp.Response)
	if err != nil {
		if resp.Response.StatusCode == 404 {
			t.Logf("ERROR - %+v", err)
			t.SkipNow()
		}
		t.Fatal(err)
	}
}

// TODO: Write a full Notification test suite.  We do use these.
