package udnssdk

import (
	"testing"
)

func Test_ListProbes(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
	if !enableProbeTests {
		t.SkipNow()
	}

	testClient, err := NewClient(testUsername, testPassword, testBaseURL)
	if err != nil {
		t.Fatal(err)
	}

	probes, resp, err := testClient.SBTCService.ListProbes("", testProbeName, testProbeType, testProbeDomain)

	if err != nil {
		if resp.Response.StatusCode == 404 {
			t.Logf("ERROR - %+v", err)
			t.SkipNow()
		}
		t.Fatal(err)
	}
	t.Logf("Probes: %+v \n", probes)
	t.Logf("Response: %+v\n", resp.Response)
	for i, e := range probes {
		t.Logf("DEBUG - Probe %d Data - %s\n", i, e.Details.data)
		err = e.Details.Populate(e.ProbeType)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("DEBUG - Populated Probe: %+v\n", e)
	}
}

func Test_GetProbeAlerts(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
	if !enableProbeTests {
		t.SkipNow()
	}

	testClient, err := NewClient(testUsername, testPassword, testBaseURL)
	if err != nil {
		t.Fatal(err)
	}

	probes, err := testClient.SBTCService.ListAllProbeAlerts(testProbeName, testProbeType, testProbeDomain)

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Probe Alerts: %+v \n", probes)
}

/* TODO: A full probe test suite.  I'm not really even sure I understand how this
 * works well enough to write one yet.  What is the correct order of operations?
 */

func Test_ListEvents(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
	if !enableProbeTests {
		t.SkipNow()
	}

	testClient, err := NewClient(testUsername, testPassword, testBaseURL)
	if err != nil {
		t.Fatal(err)
	}

	events, err := testClient.SBTCService.ListAllEvents("", testProbeName, testProbeType, testProbeDomain)

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Events: %+v \n", events)
}

// TODO: Write a full Event test suite.  We do not use these at my firm.

func Test_ListNotifications(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
	if !enableProbeTests {
		t.SkipNow()
	}

	testClient, err := NewClient(testUsername, testPassword, testBaseURL)
	if err != nil {
		t.Fatal(err)
	}

	events, resp, err := testClient.SBTCService.ListAllNotifications("", testProbeName, testProbeType, testProbeDomain)

	if err != nil {
		if resp.Response.StatusCode == 404 {
			t.Logf("ERROR - %+v", err)
			t.SkipNow()
		}
		t.Fatal(err)
	}
	t.Logf("Notifications: %+v \n", events)
	t.Logf("Response: %+v\n", resp.Response)
}

// TODO: Write a full Notification test suite.  We do use these.
