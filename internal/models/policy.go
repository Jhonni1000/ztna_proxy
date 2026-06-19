package models

import "time"

type Policy struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	SourceCIDR string `json:"source_cidr"`
	TargetURL  string `json:"target_url"`
	RequireJWT bool   `json:"require_jwt"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}