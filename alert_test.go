package udnssdk

import (
	"testing"
)

func Test_GetProbeAlerts(t *testing.T) {
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
	alerts, err := testClient.Alerts.Select(r)
	t.Logf("Probe Alerts: %+v \n", alerts)
	if err != nil {
		t.Fatal(err)
	}
}
