package gotsw

type MetalRecommendedStorageLayout struct {
	Description string      `json:"description"`
	Partitions  []Partition `json:"partitions"`
	RaidArrays  []RaidArray `json:"raidArrays"`
}
