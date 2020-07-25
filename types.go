package gncloud

import "encoding/json"

type Response struct {
	Ocs Ocs `json:"ocs"`
}

type Meta struct {
	Status     string `json:"status"`
	Statuscode int    `json:"statuscode"`
	Message    string `json:"message"`
}

type Ocs struct {
	Meta Meta            `json:"meta"`
	Data json.RawMessage `json:"data"`
}
