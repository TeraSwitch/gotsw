package gotsw

type Partition struct {
	Name       string `json:"name"`
	Device     string `json:"device"`
	SizeBytes  uint64 `json:"sizeBytes"`
	FileSystem string `json:"fileSystem"`
	MountPoint string `json:"mountPoint"`
}
