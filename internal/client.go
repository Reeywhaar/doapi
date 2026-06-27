package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const doAPIBase = "https://api.digitalocean.com/v2"

type Domain struct {
	Name string `json:"name"`
}

type Links struct {
	Pages struct {
		Next string `json:"next"`
	} `json:"pages"`
}

type DomainsResponse struct {
	Domains []Domain `json:"domains"`
	Links   Links    `json:"links"`
}

type DomainRecord struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
	Data string `json:"data"`
	TTL  int    `json:"ttl"`
}

type RecordsResponse struct {
	DomainRecords []DomainRecord `json:"domain_records"`
	Links         Links          `json:"links"`
}

type DomainWithRecords struct {
	Name    string         `json:"name"`
	Records []DomainRecord `json:"records"`
}

type Client struct {
	token string
}

func NewClient(token string) *Client {
	return &Client{token: token}
}

func (c *Client) doRequest(method, url string, body []byte) ([]byte, int, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	return data, resp.StatusCode, err
}

// perPage is the maximum page size the DigitalOcean API allows. Using the
// largest page size minimizes the number of round-trips needed to paginate.
const perPage = 200

func (c *Client) listDomains() ([]Domain, error) {
	var domains []Domain
	url := fmt.Sprintf("%s/domains?per_page=%d", doAPIBase, perPage)
	for url != "" {
		data, status, err := c.doRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		if status != http.StatusOK {
			return nil, fmt.Errorf("DO API error %d: %s", status, data)
		}
		var doResp DomainsResponse
		if err := json.Unmarshal(data, &doResp); err != nil {
			return nil, fmt.Errorf("failed to parse domains response")
		}
		domains = append(domains, doResp.Domains...)
		url = doResp.Links.Pages.Next
	}
	return domains, nil
}

func (c *Client) listRecords(domain string) ([]DomainRecord, error) {
	var records []DomainRecord
	url := fmt.Sprintf("%s/domains/%s/records?per_page=%d", doAPIBase, domain, perPage)
	for url != "" {
		data, status, err := c.doRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		if status != http.StatusOK {
			return nil, fmt.Errorf("DO API error %d: %s", status, data)
		}
		var recResp RecordsResponse
		if err := json.Unmarshal(data, &recResp); err != nil {
			return nil, fmt.Errorf("failed to parse records response")
		}
		records = append(records, recResp.DomainRecords...)
		url = recResp.Links.Pages.Next
	}
	return records, nil
}

func (c *Client) ListDomains() ([]DomainWithRecords, error) {
	domains, err := c.listDomains()
	if err != nil {
		return nil, err
	}
	result := make([]DomainWithRecords, 0, len(domains))
	for _, d := range domains {
		records, err := c.listRecords(d.Name)
		if err != nil {
			records = []DomainRecord{}
		}
		result = append(result, DomainWithRecords{Name: d.Name, Records: records})
	}
	return result, nil
}

const defaultTTL = 1800

func (c *Client) CreateRecord(domain, typ, name, data string, ttl int) ([]byte, int, error) {
	if ttl <= 0 {
		ttl = defaultTTL
	}
	body, _ := json.Marshal(map[string]any{
		"type": typ,
		"name": name,
		"data": data,
		"ttl":  ttl,
	})
	return c.doRequest("POST", doAPIBase+"/domains/"+domain+"/records", body)
}

func (c *Client) DeleteRecord(domain string, id int) (int, error) {
	url := fmt.Sprintf("%s/domains/%s/records/%d", doAPIBase, domain, id)
	_, status, err := c.doRequest("DELETE", url, nil)
	return status, err
}
