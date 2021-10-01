package gotsw

type RaidArray struct {
	Name    string   `json:"name"`
	Type    string   `json:"raidType"`
	Members []string `json:"members"`

	Partitions []Partition `json:"partitions"`

	SizeBytes  uint64 `json:"sizeBytes"`
	FileSystem string `json:"fileSystem"`
	MountPoint string `json:"mountPoint"`
}
