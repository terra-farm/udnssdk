package udnssdk

import (
	"testing"
)

func Test_GetProbeAlerts(t *testing.T) {
	if !enableProbeTests {
		t.SkipNow()
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
