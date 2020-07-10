package getnada

type resDomains struct {
	Name string `json:"name"`
}

type resInboxes struct {
	Last  int    `json:"last"`
	Total string `json:"total"`
	Msgs  []Mail `json:"msgs"`
}

type Mail struct {
	UID       string `json:"uid"`
	FromName  string `json:"f"`
	Subject   string `json:"s"`
	FromEmail string `json:"fe"`
	HTML      string `json:"html"`
}
