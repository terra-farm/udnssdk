package udnssdk

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	testUsername = os.Getenv("ULTRADNS_USERNAME")
	testPassword = os.Getenv("ULTRADNS_PASSWORD")
	testDomain   = os.Getenv("ULTRADNS_DOMAIN")
	testHostname = os.Getenv("ULTRADNS_TEST_HOSTNAME")
	testIP1      = os.Getenv("ULTRADNS_TEST_IP1")
	testIP2      = os.Getenv("ULTRADNS_TEST_IP2")
	testBaseURL  = os.Getenv("ULTRADNS_BASEURL")
	testClient   *Client
)

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())

	if testUsername == "" {
		fmt.Printf("Please configure ULTRADNS_USERNAME.\n")
		os.Exit(1)
	}
	if testPassword == "" {
		fmt.Printf("Please configure ULTRADNS_PASSWORD.\n")
		os.Exit(1)
	}
	if testDomain == "" {
		fmt.Printf("Please configure ULTRADNS_DOMAIN.\n")
		os.Exit(1)
	}
	if testHostname == "" {
		fmt.Printf("Please configure ULTRADNS_TEST_HOSTNAME.\n")
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

	os.Exit(m.Run())
}

func Test_CreateClient(t *testing.T) {
	fmt.Printf("Creating Client...\n")
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
	rrsets, resp, err := testClient.RRSets.GetRRSets(testDomain, testHostname, "ANY")

	t.Logf("GetRRSets(%s, %s, \"ANY\")", testDomain, testHostname)
	t.Logf("RRSets: %+v\n", rrsets)
	t.Logf("Response: %+v\n", resp.Response)
	if err != nil {
		t.Fatal(err)
	}
}

// Create Test
func Test_Create_RRSets(t *testing.T) {
	t.Logf("Creating %s with %s\n", testHostname, testIP1)
	rr1 := &RRSet{OwnerName: testHostname, RRType: "A", TTL: 300, RData: []string{testIP1}}
	resp, err := testClient.RRSets.CreateRRSet(testDomain, *rr1)
	t.Logf("Response: %+v\n", resp.Response)
	if err != nil {
		t.Fatal(err)
	}
}

// Another Get  Test if it matchs the Ip in IP1

//Update Test
func Test_Update_RRSets(t *testing.T) {
	t.Logf("Updating %s to %s\n", testHostname, testIP2)
	rr2 := &RRSet{OwnerName: testHostname, RRType: "A", TTL: 300, RData: []string{testIP2}}
	resp, err := testClient.RRSets.UpdateRRSet(testDomain, *rr2)
	t.Logf("Response: %+v\n", resp.Response)
	if err != nil {
		t.Fatal(err)
	}
}

// Another Get Test if it matches the Ip in IP2
func Test_GetRRSetsMid(t *testing.T) {
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
}

// Delete Test
func Test_Delete_RRSets(t *testing.T) {
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
	rrsets, resp, err := testClient.RRSets.GetRRSets(testDomain, testHostname, "ANY")

	t.Logf("GetRRSets(%s, %s, \"ANY\")", testDomain, testHostname)
	t.Logf("RRSets: %+v\n", rrsets)
	t.Logf("Response: %+v\n", resp.Response)
	if err != nil {
		t.Fatal(err)
	}
}
