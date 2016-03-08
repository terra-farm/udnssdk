package udnssdk

import (
	"testing"
)

func Test_ProbesSelectProbes(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}

	testClient, err := NewClient(testUsername, testPassword, testBaseURL)
	if err != nil {
		t.Fatal(err)
	}

	r := RRSetKey{
		Zone: testProbeDomain,
		Type: testProbeType,
		Name: testProbeName,
	}
	probes, resp, err := testClient.Probes.Select(r, "")

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
	}
}

/* TODO: A full probe test suite.  I'm not really even sure I understand how this
 * works well enough to write one yet.  What is the correct order of operations?
 */

// TODO: Write a full Event test suite.  We do not use these at my firm.
func Test_ListNotifications(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}

	testClient, err := NewClient(testUsername, testPassword, testBaseURL)
	if err != nil {
		t.Fatal(err)
	}

	events, resp, err := testClient.SBTCService.ListAllNotifications("", testProbeName, testProbeType, testProbeDomain)
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
