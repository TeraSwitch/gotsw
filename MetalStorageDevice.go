package gotsw

type MetalStorageDevice struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	CapacityBytes uint64 `json:"capacityBytes"`
}
