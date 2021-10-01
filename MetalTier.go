package gotsw

type MetalTier struct {
	Id             string               `json:"id"`
	Cpu            string               `json:"cpu"`
	MemoryGb       int                  `json:"memoryGb"`
	StorageDevices []MetalStorageDevice `json:"storageDevices"`
}
