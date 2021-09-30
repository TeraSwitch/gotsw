package gotsw

type SshKey struct {
	Id        string `json:"id"`
	ProjectId string `json:"projectId"`
	Key       string `json:"key"`
}
