package response

//NAME
//READY
//STATUS
//RESTARTS
//AGE
//IP
//NODE

type PodListItem struct {
	Name    string `json:"name"`
	Ready   string `json:"ready"`
	Status  string `json:"status"`
	Restart int32  `json:"restart"`
	Age     int64  `json:"age"`
	IP      string `json:"ip"`
	Node    string `json:"node"`
}
