package types

type HubTemplate struct {
	Model
	Side              string `json:"side"`
	NetworkTemplateID uint   `json:"netid"`
}

type EndpointTemplate struct {
	Model
	NetworkTemplateID uint `json:"netid"`
}

type NetworkTemplate struct {
	Model
	Name     string             `json:"name"`
	Comment  string             `json:"comment"`
	Hubs     []HubTemplate      `json:"hubs"`
	Endpoint []EndpointTemplate `json:"endpoints"`
}
