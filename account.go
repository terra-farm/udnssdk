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
	return AccountKey(accountName).URI()
}

// GetAccountsOfUser gets all the accounts of user
func (s *AccountsService) GetAccountsOfUser() ([]Account, *Response, error) {
	return s.Select()
}

// GetAccountStatus Get the status of a account.
func (s *AccountsService) GetAccountStatus(tid string) (Account, *Response, error) {
	return s.Find(AccountKey(tid))
}

// ListAccounts lists accounts
func (s *AccountsService) ListAccounts(query string, offset, limit int) ([]Account, *Response, error) {
	return s.Select()
}

// DeleteAccount deletes a account.
func (s *AccountsService) DeleteAccount(tid string) (*Response, error) {
	return s.Delete(AccountKey(tid))
}

// ======== //

// AccountKey represents the string identifier of an Account
type AccountKey string

// URI generates the URI for an Account
func (a AccountKey) URI() string {
	uri := "accounts"
	if a != "" {
		uri = fmt.Sprintf("accounts/%s", a)
	}
	return uri
}

// AccountsURI generates the URI for Accounts collection
func AccountsURI() string {
	return "accounts"
}

// Select requests all Accounts of user
func (s *AccountsService) Select() ([]Account, *Response, error) {
	var ald AccountListDTO
	res, err := s.client.get(AccountsURI(), &ald)

	accts := []Account{}
	for _, t := range ald.Accounts {
		accts = append(accts, t)
	}
	return accts, res, err
}

// Find requests an Account by AccountKey
func (s *AccountsService) Find(a AccountKey) (Account, *Response, error) {
	var t Account
	res, err := s.client.get(a.URI(), &t)
	return t, res, err
}

// Delete requests deletion of an Account by AccountKey
func (s *AccountsService) Delete(a AccountKey) (*Response, error) {
	return s.client.delete(a.URI(), nil)
}
