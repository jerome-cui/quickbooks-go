package quickbooks

import (
	"encoding/json"
	"errors"
	"strconv"
)

type Purchase struct {
	Id          string      `json:"Id,omitempty"`
	TxnTypeId   json.Number `json:",omitempty"`
	DocNumber   string
	PrivateNote string         `json:",omitempty"`
	PaymentType string         `json:",omitempty"`
	TxnDate     *Date          `json:",omitempty"`
	DueDate     Date           `json:",omitempty"`
	AccountRef  *ReferenceType `json:",omitempty"`
	Line        []Line

	SyncToken    string         `json:",omitempty"`
	CurrencyRef  *ReferenceType `json:",omitempty"`
	APAccountRef *ReferenceType `json:",omitempty"`
	SalesTermRef *ReferenceType `json:",omitempty"`
	LinkedTxn    []LinkedTxn    `json:",omitempty"`
	// GlobalTaxCalculation
	TotalAmt                json.Number    `json:",omitempty"`
	TransactionLocationType string         `json:",omitempty"`
	MetaData                MetaData       `json:",omitempty"`
	TxnTaxDetail            *TxnTaxDetail  `json:",omitempty"`
	ExchangeRate            json.Number    `json:",omitempty"`
	DepartmentRef           *ReferenceType `json:",omitempty"`
	IncludeInAnnualTPAR     bool           `json:",omitempty"`
	HomeBalance             json.Number    `json:",omitempty"`
	RecurDataRef            *ReferenceType `json:",omitempty"`
	Balance                 json.Number    `json:",omitempty"`
}

// CreatePurchase creates the given Purchase on the QuickBooks server, returning
// the resulting Purchase object.
func (c *Client) CreatePurchase(purchase *Purchase) (*Purchase, error) {
	var resp struct {
		Purchase Purchase
		Time     Date
	}

	if err := c.post("purchase", purchase, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Purchase, nil
}

// DeletePurchase deletes the purchase
func (c *Client) DeletePurchase(purchase *Purchase) error {
	if purchase.Id == "" || purchase.SyncToken == "" {
		return errors.New("missing id/sync token")
	}

	return c.post("purchase", purchase, nil, map[string]string{"operation": "delete"})
}

// FindPurchases gets the full list of Purchases in the QuickBooks account.
func (c *Client) FindPurchases() ([]Purchase, error) {
	var resp struct {
		QueryResponse struct {
			Purchases     []Purchase `json:"Purchase"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM Purchase", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no purchases could be found")
	}

	purchases := make([]Purchase, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM Purchase ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.Purchases == nil {
			return nil, errors.New("no purchases could be found")
		}

		purchases = append(purchases, resp.QueryResponse.Purchases...)
	}

	return purchases, nil
}

// FindPurchaseById finds the purchase by the given id
func (c *Client) FindPurchaseById(id string) (*Purchase, error) {
	var resp struct {
		Purchase Purchase
		Time     Date
	}

	if err := c.get("purchase/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Purchase, nil
}

// QueryPurchases accepts an SQL query and returns all purchases found using it
func (c *Client) QueryPurchases(query string) ([]Purchase, error) {
	var resp struct {
		QueryResponse struct {
			Purchases     []Purchase `json:"Purchase"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.Purchases == nil {
		return nil, errors.New("could not find any purchases")
	}

	return resp.QueryResponse.Purchases, nil
}

// UpdatePurchase updates the purchase
func (c *Client) UpdatePurchase(purchase *Purchase) (*Purchase, error) {
	if purchase.Id == "" {
		return nil, errors.New("missing purchase id")
	}

	existingPurchase, err := c.FindPurchaseById(purchase.Id)
	if err != nil {
		return nil, err
	}

	purchase.SyncToken = existingPurchase.SyncToken

	payload := struct {
		*Purchase
		Sparse bool `json:"sparse"`
	}{
		Purchase: purchase,
		Sparse:   true,
	}

	var purchaseData struct {
		Purchase Purchase
		Time     Date
	}

	if err = c.post("purchase", payload, &purchaseData, nil); err != nil {
		return nil, err
	}

	return &purchaseData.Purchase, err
}
