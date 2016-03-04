package udnssdk

import (
	"encoding/json"
	"testing"
)

func Test_ListAllDirectionPoolsGeoNoQuery(t *testing.T) {
	if !enableDirectionalPoolTests {
		t.SkipNow()
	}
	if testAccounts == nil {
		t.Logf("No Accounts Present, skipping...")
		t.SkipNow()
	}
	accountName := testAccounts[0].AccountName
	dpools, err := testClient.DirectionalPools.ListAllDirectionalGeoPools("", accountName)
	t.Logf("Geo Pools: %v \n", dpools)
	if err != nil {
		t.Fatal(err)
	}

	t.SkipNow()
}

func Test_ListAllDirectionPoolsGeoQuery(t *testing.T) {
	if !enableDirectionalPoolTests {
		t.SkipNow()
	}
	if testAccounts == nil {
		t.Logf("No Accounts Present, skipping...")
		t.SkipNow()
	}
	accountName := testAccounts[0].AccountName
	dpools, err := testClient.DirectionalPools.ListAllDirectionalGeoPools(testQuery, accountName)
	t.Logf("Geo Pools: %v \n", dpools)
	if err != nil {
		t.Fatal(err)
	}

	t.SkipNow()
}

func Test_ListAllDirectionalPoolsIPNoQuery(t *testing.T) {
	if !enableDirectionalPoolTests {
		t.SkipNow()
	}
	if testAccounts == nil {
		t.Logf("No Accounts Present, skipping...")
		t.SkipNow()
	}
	accountName := testAccounts[0].AccountName
	dpools, err := testClient.DirectionalPools.ListAllDirectionalIPPools("", accountName)
	t.Logf("IP Pools: %v \n", dpools)
	if err != nil {
		t.Fatal(err)
	}

	t.SkipNow()
}

func Test_ListAllDirectionalPoolsIPQuery(t *testing.T) {
	if !enableDirectionalPoolTests {
		t.SkipNow()
	}
	if testAccounts == nil {
		t.Logf("No Accounts Present, skipping...")
		t.SkipNow()
	}
	accountName := testAccounts[0].AccountName
	dpools, err := testClient.DirectionalPools.ListAllDirectionalIPPools(testQuery, accountName)
	t.Logf("IP Pools: %v \n", dpools)
	if err != nil {
		t.Fatal(err)
	}

	t.SkipNow()
}

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
