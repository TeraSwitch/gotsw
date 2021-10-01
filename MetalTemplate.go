package gotsw

type MetalTemplate struct {
	Id           uint64               `json:"id"`
	ProjectId    uint64               `json:"projectId"`
	DisplayName  string               `json:"displayName"`
	CreateModels []MetalCreateRequest `json:"createModel"`
	CloudInit    string               `json:"cloudInit"`
}
