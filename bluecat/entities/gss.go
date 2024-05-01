package entities

// GSSApplication
type GSSApplication struct {
	BAMBase       `json:"-"`
	Configuration string        `json:"configuration,omitempty"`
	View          string        `json:"view,omitempty"`
	Zone          string        `json:"zone,omitempty"`
	AbsoluteName  string        `json:"absolute_name,omitempty"`
	Fallback      []interface{} `json:"fallback,omitempty"`
	TTL           int           `json:"ttl,omitempty"`
	Properties    string        `json:"properties,omitempty"`
	HealthCheck   []interface{} `json:"health_check_type,omitempty"`
	SearchOrder   []interface{} `json:"search_order,omitempty"`
	ApplicationId int           `json:"id,omitempty"`
}

// GSSAnswer
type GSSAnswer struct {
	BAMBase       `json:"-"`
	ApplicationId int           `json:"application_id,omitempty"`
	Addresses     []interface{} `json:"addresses,omitempty"`
	Region        string        `json:"region,omitempty"`
	Name          string        `json:"name,omitempty"`
	AnswerId      int           `json:"id,omitempty"`
	Type          string        `json:"type,omitempty"`
}

// GSSSearchOrder
type GSSSearchOrder struct {
	BAMBase       `json:"-"`
	Nodes         []interface{} `json:"nodes,omitempty"`
	Links         []interface{} `json:"links,omitempty"`
	Name          string        `json:"name,omitempty"`
	SearchOrderId int           `json:"id,omitempty"`
}
