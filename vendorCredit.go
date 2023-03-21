package quickbooks

import (
	"encoding/json"
	"errors"
	"strconv"
)

type VendorCredit struct {
	Id        string         `json:"Id,omitempty"`
	TotalAmt  json.Number    `json:",omitempty"`
	TxnDate   *Date          `json:",omitempty"`
	VendorRef *ReferenceType `json:",omitempty"`
	Line      []Line

	SyncToken               string         `json:",omitempty"`
	CurrencyRef             *ReferenceType `json:",omitempty"`
	APAccountRef            *ReferenceType `json:",omitempty"`
	SalesTermRef            *ReferenceType `json:",omitempty"`
	LinkedTxn               []LinkedTxn    `json:",omitempty"`
	TransactionLocationType string         `json:",omitempty"`
	DueDate                 Date           `json:",omitempty"`
	MetaData                MetaData       `json:",omitempty"`
	DocNumber               string
	PrivateNote             string         `json:",omitempty"`
	TxnTaxDetail            *TxnTaxDetail  `json:",omitempty"`
	ExchangeRate            json.Number    `json:",omitempty"`
	DepartmentRef           *ReferenceType `json:",omitempty"`
	IncludeInAnnualTPAR     bool           `json:",omitempty"`
	HomeBalance             json.Number    `json:",omitempty"`
	RecurDataRef            *ReferenceType `json:",omitempty"`
	Balance                 json.Number    `json:",omitempty"`
}

// CreateVendorCredit creates the given VendorCredit on the QuickBooks server, returning
// the resulting VendorCredit object.
func (c *Client) CreateVendorCredit(vendorCredit *VendorCredit) (*VendorCredit, error) {
	var resp struct {
		VendorCredit VendorCredit
		Time         Date
	}

	if err := c.post("vendorCredit", vendorCredit, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.VendorCredit, nil
}

// DeleteVendorCredit deletes the vendorCredit
func (c *Client) DeleteVendorCredit(vendorCredit *VendorCredit) error {
	if vendorCredit.Id == "" || vendorCredit.SyncToken == "" {
		return errors.New("missing id/sync token")
	}

	return c.post("vendorCredit", vendorCredit, nil, map[string]string{"operation": "delete"})
}

// FindVendorCredits gets the full list of VendorCredits in the QuickBooks account.
func (c *Client) FindVendorCredits() ([]VendorCredit, error) {
	var resp struct {
		QueryResponse struct {
			VendorCredits []VendorCredit `json:"VendorCredit"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM VendorCredit", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no vendorCredits could be found")
	}

	vendorCredits := make([]VendorCredit, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM VendorCredit ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.VendorCredits == nil {
			return nil, errors.New("no vendorCredits could be found")
		}

		vendorCredits = append(vendorCredits, resp.QueryResponse.VendorCredits...)
	}

	return vendorCredits, nil
}

// FindVendorCreditById finds the vendorCredit by the given id
func (c *Client) FindVendorCreditById(id string) (*VendorCredit, error) {
	var resp struct {
		VendorCredit VendorCredit
		Time         Date
	}

	if err := c.get("vendorCredit/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.VendorCredit, nil
}

// QueryVendorCredits accepts an SQL query and returns all vendorCredits found using it
func (c *Client) QueryVendorCredits(query string) ([]VendorCredit, error) {
	var resp struct {
		QueryResponse struct {
			VendorCredits []VendorCredit `json:"VendorCredit"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.VendorCredits == nil {
		return nil, errors.New("could not find any vendorCredits")
	}

	return resp.QueryResponse.VendorCredits, nil
}

// UpdateVendorCredit updates the vendorCredit
func (c *Client) UpdateVendorCredit(vendorCredit *VendorCredit) (*VendorCredit, error) {
	if vendorCredit.Id == "" {
		return nil, errors.New("missing vendorCredit id")
	}

	existingVendorCredit, err := c.FindVendorCreditById(vendorCredit.Id)
	if err != nil {
		return nil, err
	}

	vendorCredit.SyncToken = existingVendorCredit.SyncToken

	payload := struct {
		*VendorCredit
		Sparse bool `json:"sparse"`
	}{
		VendorCredit: vendorCredit,
		Sparse:       true,
	}

	var vendorCreditData struct {
		VendorCredit VendorCredit
		Time         Date
	}

	if err = c.post("vendorCredit", payload, &vendorCreditData, nil); err != nil {
		return nil, err
	}

	return &vendorCreditData.VendorCredit, err
}
