package quickbooks

import (
	"encoding/json"
	"errors"
	"strconv"
)

type CreditMemo struct {
	Id         string `json:"Id,omitempty"`
	DocNumber  string
	CustomMemo string         `json:",omitempty"`
	TxnDate    *Date          `json:",omitempty"`
	TotalAmt   json.Number    `json:",omitempty"`
	CustomRef  *ReferenceType `json:",omitempty"`
	Line       []Line

	SyncToken               string         `json:",omitempty"`
	CurrencyRef             *ReferenceType `json:",omitempty"`
	APAccountRef            *ReferenceType `json:",omitempty"`
	SalesTermRef            *ReferenceType `json:",omitempty"`
	LinkedTxn               []LinkedTxn    `json:",omitempty"`
	TransactionLocationType string         `json:",omitempty"`
	DueDate                 Date           `json:",omitempty"`
	MetaData                MetaData       `json:",omitempty"`
	TxnTaxDetail            *TxnTaxDetail  `json:",omitempty"`
	ExchangeRate            json.Number    `json:",omitempty"`
	DepartmentRef           *ReferenceType `json:",omitempty"`
	IncludeInAnnualTPAR     bool           `json:",omitempty"`
	HomeBalance             json.Number    `json:",omitempty"`
	RecurDataRef            *ReferenceType `json:",omitempty"`
	Balance                 json.Number    `json:",omitempty"`
}

// CreateCreditMemo creates the given CreditMemo on the QuickBooks server, returning
// the resulting CreditMemo object.
func (c *Client) CreateCreditMemo(creditMemo *CreditMemo) (*CreditMemo, error) {
	var resp struct {
		CreditMemo CreditMemo
		Time       Date
	}

	if err := c.post("creditMemo", creditMemo, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.CreditMemo, nil
}

// DeleteCreditMemo deletes the creditMemo
func (c *Client) DeleteCreditMemo(creditMemo *CreditMemo) error {
	if creditMemo.Id == "" || creditMemo.SyncToken == "" {
		return errors.New("missing id/sync token")
	}

	return c.post("creditMemo", creditMemo, nil, map[string]string{"operation": "delete"})
}

// FindCreditMemos gets the full list of CreditMemos in the QuickBooks account.
func (c *Client) FindCreditMemos() ([]CreditMemo, error) {
	var resp struct {
		QueryResponse struct {
			CreditMemos   []CreditMemo `json:"CreditMemo"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM CreditMemo", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no creditMemos could be found")
	}

	creditMemos := make([]CreditMemo, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM CreditMemo ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.CreditMemos == nil {
			return nil, errors.New("no creditMemos could be found")
		}

		creditMemos = append(creditMemos, resp.QueryResponse.CreditMemos...)
	}

	return creditMemos, nil
}

// FindCreditMemoById finds the creditMemo by the given id
func (c *Client) FindCreditMemoById(id string) (*CreditMemo, error) {
	var resp struct {
		CreditMemo CreditMemo
		Time       Date
	}

	if err := c.get("creditMemo/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.CreditMemo, nil
}

// QueryCreditMemos accepts an SQL query and returns all creditMemos found using it
func (c *Client) QueryCreditMemos(query string) ([]CreditMemo, error) {
	var resp struct {
		QueryResponse struct {
			CreditMemos   []CreditMemo `json:"CreditMemo"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.CreditMemos == nil {
		return nil, errors.New("could not find any creditMemos")
	}

	return resp.QueryResponse.CreditMemos, nil
}

// UpdateCreditMemo updates the creditMemo
func (c *Client) UpdateCreditMemo(creditMemo *CreditMemo) (*CreditMemo, error) {
	if creditMemo.Id == "" {
		return nil, errors.New("missing creditMemo id")
	}

	existingCreditMemo, err := c.FindCreditMemoById(creditMemo.Id)
	if err != nil {
		return nil, err
	}

	creditMemo.SyncToken = existingCreditMemo.SyncToken

	payload := struct {
		*CreditMemo
		Sparse bool `json:"sparse"`
	}{
		CreditMemo: creditMemo,
		Sparse:     true,
	}

	var creditMemoData struct {
		CreditMemo CreditMemo
		Time       Date
	}

	if err = c.post("creditMemo", payload, &creditMemoData, nil); err != nil {
		return nil, err
	}

	return &creditMemoData.CreditMemo, err
}
