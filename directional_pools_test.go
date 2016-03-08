package udnssdk

import (
	"encoding/json"
	"testing"
)

func Test_ListAllDirectionPoolsGeoNoQuery(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
	if !enableDirectionalPoolTests {
		t.SkipNow()
	}
	if testAccounts == nil {
		t.Logf("No Accounts Present, skipping...")
		t.SkipNow()
	}

	accountName := testAccounts[0].AccountName
	dpools, err := testClient.DirectionalPools.ListAllDirectionalGeoPools("", accountName)

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Geo Pools: %v \n", dpools)
}

func Test_ListAllDirectionPoolsGeoQuery(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
	if !enableDirectionalPoolTests {
		t.SkipNow()
	}
	if testAccounts == nil {
		t.Logf("No Accounts Present, skipping...")
		t.SkipNow()
	}

	accountName := testAccounts[0].AccountName
	dpools, err := testClient.DirectionalPools.ListAllDirectionalGeoPools(testQuery, accountName)

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Geo Pools: %v \n", dpools)
}

func Test_ListAllDirectionalPoolsIPNoQuery(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
	if !enableDirectionalPoolTests {
		t.SkipNow()
	}
	if testAccounts == nil {
		t.Logf("No Accounts Present, skipping...")
		t.SkipNow()
	}

	accountName := testAccounts[0].AccountName
	dpools, err := testClient.DirectionalPools.ListAllDirectionalIPPools("", accountName)

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("IP Pools: %v \n", dpools)
}

func Test_ListAllDirectionalPoolsIPQuery(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
	if !enableDirectionalPoolTests {
		t.SkipNow()
	}
	if testAccounts == nil {
		t.Logf("No Accounts Present, skipping...")
		t.SkipNow()
	}

	accountName := testAccounts[0].AccountName
	dpools, err := testClient.DirectionalPools.ListAllDirectionalIPPools(testQuery, accountName)

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("IP Pools: %v \n", dpools)
}

func Test_Create_DirectionalPoolIP(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
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

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Response: %+v\n", resp.Response)
}

func Test_Get_DirectionalPoolIP(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
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
	t.Logf("Test Get IP DPool Group (%s, %s)\n", testIPDPool.Name, testIPDPool)
	dp, resp, err := testClient.DirectionalPools.GetDirectionalIPPool(testIPDPool.Name, accountName)

	if err != nil {
		t.Logf("GetDirectionalPoolIP Error: %+v\n", err)
		if resp.Response.StatusCode == 404 {
			return
		}
		t.Fatal(err)
	}
	t.Logf("Response: %+v\n", resp.Response)
	t.Logf("DPool: %+v\n", dp)
	dp2, er := json.Marshal(dp)
	t.Logf("DPool Marshalled back: %s - %+v\n", string(dp2), er)
}

func Test_Delete_DirectionalPoolIP(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
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
	t.Logf("Delete IP DPool Group (%s, %s)\n", testIPDPool.Name, testIPDPool)
	resp, err := testClient.DirectionalPools.DeleteDirectionalIPPool(testIPDPool.Name, accountName)

	if err != nil {
		t.Logf("DeleteDirectionalPoolIP Error: %+v\n", err)
		if resp.Response.StatusCode == 404 {
			return
		}
		t.Fatal(err)
	}
	t.Logf("Response: %+v\n", resp.Response)
}
