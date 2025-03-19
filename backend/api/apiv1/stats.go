package apiv1;



type Query_Submit struct {
	Key string `json:"Key"`
}

type Result_Submit struct {
	Uptime uint64 `json:"uptime"`
	Rank uint64 `json:"rank"`
}
