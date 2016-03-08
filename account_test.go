package udnssdk

import (
	"testing"
)

func Test_Online_ListAccountsOfUser(t *testing.T) {
	if !enableIntegrationTests {
		t.SkipNow()
	}

	testClient, err := NewClient(testUsername, testPassword, testBaseURL)
	if err != nil {
		t.Fatal(err)
	}

	accounts, resp, err := testClient.Accounts.Select()
	t.Logf("Accounts: %+v \n", accounts)
	t.Logf("Response: %+v\n", resp.Response)
	testAccounts = accounts
	if err != nil {
		if resp.Response.StatusCode == 404 {
			t.SkipNow()
		}
		t.Fatal(err)
	}
}
