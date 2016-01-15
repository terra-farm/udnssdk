package udnssdk

import (
	"fmt"
)

// ZonesService handles communication with the Zone related blah blah
type AccountsService struct {
	client *Client
}

type Account struct {
	AccountName           string `json:"accountName"`
	AccountHolderUserName string `json:"accountHolderUserName"`
	OwnerUserName         string `json:"ownerUserName"`
	NumberOfUsers         int    `json:"numberOfUsers"`
	NumberOfGroups        int    `json:"numberOfGroups"`
	AccountType           string `json:"accountType"`
}

type AccountListDTO struct {
	Accounts   []Account  `json:"accounts"`
	Resultinfo ResultInfo `json:"resultInfo"`
}
type accountWrapper struct {
	Account Account `json:"account"`
}

// accountPath links to the account url.
func accountPath(accountName string) string {
	path := "accounts"
	if accountName != "" {
		path = fmt.Sprintf("accounts/%s", accountName)
	}
	return path
}

/*
func accountPath(tid string) string {
	return fmt.Sprintf("accounts/%s", tid)
}
*/

func (s *AccountsService) GetAccountsOfUser() ([]Account, *Response, error) {
	reqStr := accountPath("")
	var ald AccountListDTO
	res, err := s.client.get(reqStr, &ald)
	if err != nil {
		return []Account{}, res, err
	}
	accts := []Account{}
	for _, t := range ald.Accounts {
		accts = append(accts, t)
	}
	return accts, res, nil
}

/*
// TODO:  Implement Zones
func (s *AccountsService) GetZonesOfAccount(accountName string) ([]Account, *Response, error) {
	reqStr := fmt.Sprintf("%s/zones", accountPath(accountName))
	var ald AccountListDTO
	fmt.Printf("In GetZonesOfAccount(%s)..  ReqStr: %s\n", accountName, reqStr)
	res, err := s.client.get(reqStr, &ald)
	if err != nil {
		return []Account{}, res, err
	}
	accts := []Account{}
	for _, t := range ald.Accounts {
		accts = append(accts, t)
	}
	fmt.Printf("Exiting GetZonesOfAccount(%s)..\n", accountName)
	return accts, res, nil
}
*/
// Get the status of a account.
func (s *AccountsService) GetAccountStatus(tid string) (Account, *Response, error) {
	reqStr := accountPath(tid)
	var t Account
	res, err := s.client.get(reqStr, &t)
	if err != nil {
		return t, res, err
	}
	return t, res, err
}

// HTTP BS to dance around bad program structure
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

// List accounts
//
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
//
func (s *AccountsService) DeleteAccount(tid string) (*Response, error) {
	path := accountPath(tid)
	return s.client.delete(path, nil)
}
