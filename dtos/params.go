package dtos

import "time"

type GetAllAssetsRequest struct {
	Id            string    `json:"id"`
	TypeForm      string    `json:"type_form"`
	Description   string    `json:"description"`
	Timestamp     time.Time `json:"timestamp"`
	InsertionType string    `json:"insertion_type"`
	Hash          string    `json:"hash"`
}

type PostAssetRequest struct {
	Id            string    `json:"id"`
	TypeForm      string    `json:"type_form"`
	Description   string    `json:"description"`
	Timestamp     time.Time `json:"timestamp"`
	InsertionType string    `json:"insertion_type"`
	Hash          string    `json:"hash"`
}

type AssetRequest struct {
	Id            string    `json:"id"`
	TypeForm      string    `json:"type_form"`
	Description   string    `json:"description"`
	Timestamp     time.Time `json:"timestamp"`
	InsertionType string    `json:"insertion_type"`
	Hash          string    `json:"hash"`
}

type PutAssetRequest struct {
	TypeForm      string    `json:"type_form"`
	Description   string    `json:"description"`
	Timestamp     time.Time `json:"timestamp"`
	InsertionType string    `json:"insertion_type"`
	Hash          string    `json:"hash"`
}

type Filter struct {
	Ids            []string        `json:"ids"`
	TypeForms      []string        `json:"type_forms"`
	InsertionTypes []string        `json:"insertion_types"`
	Hashs          []string        `json:"hashs"`
	TimeFilter     TimestampFilter `json:"time_filter"`
}

type TimestampFilter struct {
	Min time.Time `json:"min"`
	Max time.Time `json:"max"`
}
