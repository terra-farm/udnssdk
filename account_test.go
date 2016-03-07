package udnssdk

import (
	"testing"
)

func Test_ListAccountsOfUser(t *testing.T) {

	if !enableAccountTests {
		t.SkipNow()
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
