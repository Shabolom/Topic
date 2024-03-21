package models

type Topic struct {
	Name    string `json:"name"`
	About   string `json:"about"`
	Creator string `json:"creator"`
}
