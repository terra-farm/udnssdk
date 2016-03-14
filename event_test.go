package udnssdk

import (
	"testing"
)

func Test_ListEvents(t *testing.T) {
	if !enableProbeTests {
		t.SkipNow()
	}
	r := RRSetKey{
		Zone: testProbeDomain,
		Type: testProbeType,
		Name: testProbeName,
	}
	events, err := testClient.Events.Select(r, "")
	t.Logf("Events: %+v \n", events)
	if err != nil {
		t.Fatal(err)
	}
}
