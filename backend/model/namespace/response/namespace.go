package response

type NameSpace struct {
	Name              string `json:"name"`
	CreationTimeStamp int64  `json:"creationTimeStamp"`
	Status            string `json:"status"`
}
