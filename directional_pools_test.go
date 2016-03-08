package udnssdk

import (
	"encoding/json"
	"testing"
)

func Test_ListAllDirectionPoolsGeoNoQuery(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
	if testAccounts == nil {
		t.Logf("No Accounts Present, skipping...")
		t.SkipNow()
	}


	testClient, err := NewClient(testUsername, testPassword, testBaseURL)
	if err != nil {
		t.Fatal(err)
	}

	accountName := testAccounts[0].AccountName
	p := GeoDirectionalPoolKey{Account: AccountKey(accountName)}
	dpools, err := testClient.DirectionalPools.Geos().Select(p, "")
	t.Logf("Geo Pools: %v \n", dpools)

	if err != nil {
		t.Fatal(err)
	}
}

func Test_ListAllDirectionPoolsGeoQuery(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
	if testAccounts == nil {
		t.Logf("No Accounts Present, skipping...")
		t.SkipNow()
	}

	testClient, err := NewClient(testUsername, testPassword, testBaseURL)
	if err != nil {
		t.Fatal(err)
	}

	accountName := testAccounts[0].AccountName
	p := GeoDirectionalPoolKey{Account: AccountKey(accountName)}
	dpools, err := testClient.DirectionalPools.Geos().Select(p, testQuery)

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Geo Pools: %v \n", dpools)
}

func Test_ListAllDirectionalPoolsIPNoQuery(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
	if testAccounts == nil {
		t.Logf("No Accounts Present, skipping...")
		t.SkipNow()
	}

	testClient, err := NewClient(testUsername, testPassword, testBaseURL)
	if err != nil {
		t.Fatal(err)
	}

	accountName := testAccounts[0].AccountName
	p := IPDirectionalPoolKey{Account: AccountKey(accountName)}
	dpools, err := testClient.DirectionalPools.IPs().Select(p, "")

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("IP Pools: %v \n", dpools)
}

func Test_ListAllDirectionalPoolsIPQuery(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
	if testAccounts == nil {
		t.Logf("No Accounts Present, skipping...")
		t.SkipNow()
	}

	testClient, err := NewClient(testUsername, testPassword, testBaseURL)
	if err != nil {
		t.Fatal(err)
	}

	accountName := testAccounts[0].AccountName
	p := IPDirectionalPoolKey{Account: AccountKey(accountName)}
	dpools, err := testClient.DirectionalPools.IPs().Select(p, testQuery)
	t.Logf("IP Pools: %v \n", dpools)

	if err != nil {
		t.Fatal(err)
	}
}

func Test_Create_DirectionalPoolIP(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
	if testAccounts == nil {
		t.Logf("No Accounts Present, skipping...")
		t.SkipNow()
	}

	testClient, err := NewClient(testUsername, testPassword, testBaseURL)
	if err != nil {
		t.Fatal(err)
	}

	accountName := testAccounts[0].AccountName
	t.Logf("Creating %s with %+v\n", testIPDPool.Name, testIPDPool)
	p := IPDirectionalPoolKey{
		Account: AccountKey(accountName),
		ID:      testIPDPool.Name,
	}
	resp, err := testClient.DirectionalPools.IPs().Create(p, testIPDPool)
	t.Logf("Response: %+v\n", resp.Response)

	if err != nil {
		t.Fatal(err)
	}
}

func Test_Get_DirectionalPoolIP(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}
	if testAccounts == nil {
		t.Logf("No Accounts Present, skipping...")
		t.SkipNow()
	}

	testClient, err := NewClient(testUsername, testPassword, testBaseURL)
	if err != nil {
		t.Fatal(err)
	}

	accountName := testAccounts[0].AccountName
	p := GeoDirectionalPoolKey{
		Account: AccountKey(accountName),
		ID:      testIPDPool.Name,
	}
	dp, resp, err := testClient.DirectionalPools.Geos().Find(p)

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
	if !enableIntegrationTests {
		t.SkipNow()
	}
	if testAccounts == nil {
		t.Logf("No Accounts Present, skipping...")
		t.SkipNow()
	}

	testClient, err := NewClient(testUsername, testPassword, testBaseURL)
	if err != nil {
		t.Fatal(err)
	}

	accountName := testAccounts[0].AccountName
	p := GeoDirectionalPoolKey{
		Account: AccountKey(accountName),
		ID:      testIPDPool.Name,
	}
	resp, err := testClient.DirectionalPools.Geos().Delete(p)

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
