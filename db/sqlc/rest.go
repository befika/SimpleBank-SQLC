package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
)

const DefaultPageSize = 10

type Sort []string

type QueryParams struct {
	Sort         string `json:"sort" form:"sort"`
	Filter       string `json:"filter" form:"filter"`
	Page         string `json:"page" form:"page"`
	PerPage      string `json:"per_page" form:"per_page"`
	LinkOperator string `json:"linkOperator" form:"linkOperator"`
}

type FilterParams struct {
	Sort Sort `json:"sort"`
	//	Filter  []Filter `json:"filter"`
	Page         int64    `json:"page"`
	PerPage      int64    `json:"per_page"`
	Filter       []Filter `json:"filter"`
	LinkOperator string   `json:"linkOperator"`
	Total        int64    `json:"total"`
	NoSort       bool     `json:"-"`
	NoLimit      bool     `json:"-"`
}

var OperationKey = "operation"

type Filter struct {
	ColumnField   string `json:"columnField"`
	OperatorValue string `json:"operatorValue"`
	Value         string `json:"value"`
}

type Filter2Search struct {
	Columns []string `json:"columns"`
	Value   string   `json:"value"`
}

func (q QueryParams) Get() (*FilterParams, error) {
	res := &FilterParams{}
	// res.Filter = []Filter{}
	if q.Sort != "" {
		decSort, err := url.QueryUnescape(q.Sort)
		if err != nil {
			log.Println("error while unescaping sort query param")
			return nil, err
		}
		q.Sort = decSort

		err = json.Unmarshal([]byte(q.Sort), &res.Sort)
		if err != nil {
			log.Println("error while unmarshalling sort json")
			return nil, err
		}
	}

	if q.Page == "" {
		res.Page = 1
	} else {

		p, err := toInt64(q.Page)
		if err != nil {
			p = 1
		}

		res.Page = p
	}

	fil := []Filter{}
	if q.Filter != "" {
		err := json.Unmarshal([]byte(q.Filter), &fil)
		if err != nil {
			return nil, errors.New("error while unescaping filter query param")
		}
	}

	if q.PerPage == "" {
		res.PerPage = DefaultPageSize
	} else {

		pp, err := toInt64(q.PerPage)
		if err != nil {
			pp = DefaultPageSize
		}
		res.PerPage = pp
	}
	if q.LinkOperator != "" {
		res.LinkOperator = q.LinkOperator
	}
	res.Filter = fil
	return res, nil
}

func ValidateOperationValue(operation string) error {
	if operation != "lth" {
		return errors.New("invalid operation")
	}
	return nil
}

func toInt64(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (f *FilterParams) ToQuery() url.Values {
	q := url.Values{}
	if f.Page != 0 {
		q.Add("page", fmt.Sprintf("%d", f.Page))
	}

	if f.PerPage != 0 {
		q.Add("per_page", fmt.Sprintf("%d", f.PerPage))
	}

	if len(f.Sort) > 1 {
		q.Add("sort", fmt.Sprintf(`["%v","%v"]`, f.Sort[0], f.Sort[1]))
	} else if len(f.Sort) == 1 {
		q.Add("sort", fmt.Sprintf(`["%v","ASC"]`, f.Sort[0]))
	}
	return q
}
