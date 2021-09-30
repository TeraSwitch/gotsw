package gotsw

type Image struct {
	Id                         string `json:"id"`
	DisplayName                string `json:"displayName"`
	OperatingSystemName        string `json:"operatingSystemName"`
	OperatingSystemVersion     string `json:"operatingSystemVersion"`
	DisableCustomizableStorage bool   `json:"disableCustomizableStorage"`
}
