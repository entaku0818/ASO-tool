package model

import "time"

type Keyword struct {
	ID        string    `json:"id"`
	AppID     string    `json:"app_id"`
	Keyword   string    `json:"keyword"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateKeywordRequest struct {
	AppID   string `json:"app_id"`
	Keyword string `json:"keyword"`
	Country string `json:"country,omitempty"`
}

func (r *CreateKeywordRequest) Validate() error {
	if r.AppID == "" {
		return ErrAppIDRequired
	}
	if r.Keyword == "" {
		return ErrKeywordRequired
	}
	if r.Country == "" {
		r.Country = "JP"
	}
	return nil
}
