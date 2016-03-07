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
