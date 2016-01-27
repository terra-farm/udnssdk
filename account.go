package udnssdk

import (
	"fmt"
)

// AccountsService provides access to account resources
type AccountsService struct {
	client *Client
}

// Account represents responses from the service
type Account struct {
	AccountName           string `json:"accountName"`
	AccountHolderUserName string `json:"accountHolderUserName"`
	OwnerUserName         string `json:"ownerUserName"`
	NumberOfUsers         int    `json:"numberOfUsers"`
	NumberOfGroups        int    `json:"numberOfGroups"`
	AccountType           string `json:"accountType"`
}

// AccountListDTO represents a account index response
type AccountListDTO struct {
	Accounts   []Account  `json:"accounts"`
	Resultinfo ResultInfo `json:"resultInfo"`
}

// accountPath links to the account url.
func accountPath(accountName string) string {
	path := "accounts"
	if accountName != "" {
		path = fmt.Sprintf("accounts/%s", accountName)
	}
	return path
}

// GetAccountsOfUser gets all the accounts of user
func (s *AccountsService) GetAccountsOfUser() ([]Account, *Response, error) {
	var ald AccountListDTO
	uri := accountPath("")
	res, err := s.client.get(uri, &ald)

	accts := []Account{}
	for _, t := range ald.Accounts {
		accts = append(accts, t)
	}
	return accts, res, err
}

/*
// TODO:  Implement Zones
func (s *AccountsService) GetZonesOfAccount(accountName string) ([]Account, *Response, error) {
	reqStr := fmt.Sprintf("%s/zones", accountPath(accountName))
	var ald AccountListDTO
	log.Printf("In GetZonesOfAccount(%s)..  ReqStr: %s\n", accountName, reqStr)
	res, err := s.client.get(reqStr, &ald)
	if err != nil {
		return []Account{}, res, err
	}
	accts := []Account{}
	for _, t := range ald.Accounts {
		accts = append(accts, t)
	}
	log.Printf("Exiting GetZonesOfAccount(%s)..\n", accountName)
	return accts, res, nil
}
*/

// GetAccountStatus Get the status of a account.
func (s *AccountsService) GetAccountStatus(tid string) (Account, *Response, error) {
	reqStr := accountPath(tid)
	var t Account
	res, err := s.client.get(reqStr, &t)
	if err != nil {
		return t, res, err
	}
	return t, res, err
}

// GetAccountResultByURI gets account result by URI
func (s *AccountsService) GetAccountResultByURI(uri string) (*Response, error) {
	req, err := s.client.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	res, err := s.client.HttpClient.Do(req)

	if err != nil {
		return &Response{Response: res}, err
	}
	return &Response{Response: res}, err
}

// ListAccounts lists accounts
func (s *AccountsService) ListAccounts(query string, offset, limit int) ([]Account, *Response, error) {
	// TODO: Soooo... This function does not handle pagination of Accounts....
	//v := url.Values{}

	reqStr := "accounts"
	var tld AccountListDTO
	//wrappedAccounts := []Account{}

	res, err := s.client.get(reqStr, &tld)
	if err != nil {
		return []Account{}, res, err
	}

	accounts := []Account{}
	for _, t := range tld.Accounts {
		accounts = append(accounts, t)
	}

	return accounts, res, nil
}

// DeleteAccount deletes a account.
func (s *AccountsService) DeleteAccount(tid string) (*Response, error) {
	path := accountPath(tid)
	return s.client.delete(path, nil)
}
