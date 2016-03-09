package udnssdk

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func Test_ListAllRRSetsPre(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
	if testClient == nil {
		t.Fatalf("TestClient Not Defined?\n")
	}
	if !enableRRSetTests {
		t.SkipNow()
	}

	t.Logf("ListAllRRSets(%s, %s, \"ANY\")", testDomain, testHostname)
	rrsets, err := testClient.RRSets.ListAllRRSets(testDomain, testHostname, "ANY")

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("RRSets: %+v\n", rrsets)
}

func Test_ListRRSets(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
	if !enableRRSetTests {
		t.SkipNow()
	}

	t.Logf("ListAllRRSets(%s, \"\", \"\")", testDomain)
	rrsets, err := testClient.RRSets.ListAllRRSets(testDomain, "", "")

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("RRSets: %+v\n", rrsets)
	t.Logf("Checking for profiles...\n")
	for _, rr := range rrsets {
		if rr.Profile != nil {
			typ := rr.Profile.GetType()
			if typ == "" {
				t.Fatalf("Could not get type for profile %+v\n", rr.Profile)
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
	if !enableIntegrationTests {
		t.SkipNow()
	}
	if !enableRRSetTests {
		t.SkipNow()
	}
	if !enableChanges {
		t.SkipNow()
	}

	t.Logf("Creating %s with %s\n", testHostname, testIP1)
	rr1 := &RRSet{OwnerName: testHostname, RRType: "A", TTL: 300, RData: []string{testIP1}, Profile: &StringProfile{Profile: testProfile}}
	resp, err := testClient.RRSets.CreateRRSet(testDomain, *rr1)

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Response: %+v\n", resp.Response)
}

// Another Get  Test if it matchs the Ip in IP1
func Test_ListAllRRSetsMid1(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
	if !enableRRSetTests {
		t.SkipNow()
	}

	t.Logf("ListAllRRSets(%s, %s, \"ANY\")", testDomain, testHostname)
	rrsets, err := testClient.RRSets.ListAllRRSets(testDomain, testHostname, "ANY")

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("RRSets: %+v\n", rrsets)
	// Do the test v IP1 here
	if rrsets[0].RData[0] != testIP1 {
		t.Fatalf("RData[0]\"%s\" != testIP1\"%s\"", rrsets[0].RData[0], testIP1)
	}
}

// Update Test
func Test_Update_RRSets(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
	if !enableRRSetTests {
		t.SkipNow()
	}
	if !enableChanges {
		t.SkipNow()
	}

	t.Logf("Updating %s to %s\n", testHostname, testIP2)
	rr2 := &RRSet{OwnerName: testHostname, RRType: "A", TTL: 300, RData: []string{testIP2}, Profile: &StringProfile{Profile: testProfile2}}
	resp, err := testClient.RRSets.UpdateRRSet(testDomain, *rr2)

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Response: %+v\n", resp.Response)
}

// Another Get Test if it matches the Ip in IP2
func Test_ListAllRRSetsMid(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
	if !enableRRSetTests {
		t.SkipNow()
	}

	t.Logf("ListAllRRSets(%s, %s, \"ANY\")", testDomain, testHostname)
	rrsets, err := testClient.RRSets.ListAllRRSets(testDomain, testHostname, "ANY")

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("RRSets: %+v\n", rrsets)
	// Do the test v IP2 here
	if rrsets[0].RData[0] != testIP2 {
		t.Fatalf("RData[0]\"%s\" != testIP2\"%s\"", rrsets[0].RData[0], testIP2)
	}
	t.Logf("Profile Check: %+v", rrsets[0].Profile.GetProfileObject())
}

// Delete Test
func Test_Delete_RRSets(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
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

	t.Logf("ListAllRRSets(%s, %s, \"ANY\")", testDomain, testHostname)
	rrsets, err := testClient.RRSets.ListAllRRSets(testDomain, testHostname, "ANY")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("RRSets: %+v\n", rrsets)
	for i, e := range rrsets {
		t.Logf("Deleting %s  - ( %d ) %+v \n", testHostname, i, e)
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

func Test_ListAllRRSetsPost(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
	if !enableRRSetTests {
		t.SkipNow()
	}

	t.Logf("ListAllRRSets(%s, %s, \"ANY\")", testDomain, testHostname)
	rrsets, err := testClient.RRSets.ListAllRRSets(testDomain, testHostname, "ANY")

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("RRSets: %+v\n", rrsets)
}
