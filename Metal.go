package gotsw

type Metal struct {
	Id          uint64    `json:"id"`
	ProjectId   uint64    `json:"projectId"`
	Status      string    `json:"status"`
	RegionId    string    `json:"regionId"`
	TierId      string    `json:"tierId"`
	DisplayName string    `json:"displayName"`
	PowerState  string    `json:"powerState"`
	CurrentTask string    `json:"currentTask"`
	ImageId     string    `json:"imageId"`
	IpAddresses []string  `json:"ipAddresses"`
	Region      Region    `json:"region"`
	Tier        MetalTier `json:"tier"`
	Image       Image     `json:"image"`
}
