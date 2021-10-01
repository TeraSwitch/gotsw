package gotsw

type SshKey struct {
	Id        uint64 `json:"id"`
	ProjectId string `json:"projectId"`
	Key       string `json:"key"`
}
