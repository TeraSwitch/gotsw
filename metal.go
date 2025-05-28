package gotsw

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"time"
)

// Metal represents a bare metal service
type Metal struct {
	ID          int64   `json:"id"`
	Created     string  `json:"created"`
	Deleted     *string `json:"deleted"`
	ObjectType  string  `json:"objectType"`
	ProjectID   int64   `json:"projectId"`
	DisplayName string  `json:"displayName"`

	// Service-specific fields
	RegionID   string     `json:"regionId"`
	Status     Status     `json:"status"`
	PowerState PowerState `json:"powerState"`

	// Hardware configuration
	TierID   string `json:"tierId"`
	MemoryGB int32  `json:"memoryGb"`
	ImageID  string `json:"imageId"`
	StorageDevices map[string]MetalStorageDevice `json:"storageDevices"`

	// Network configuration
	IPAddresses []netip.Addr `json:"ipAddresses"`

	// Pricing
	MonthlyPrice float64 `json:"monthlyPrice"`
	HourlyPrice  float64 `json:"hourlyPrice"`

	// Additional metadata
	Tags   []string            `json:"tags"`
	Events []ProvisioningEvent `json:"events"`
}

// Status represents the current status of a service
type Status string

const (
	StatusPending    Status = "Pending"
	StatusActive     Status = "Active"
	StatusSuspended  Status = "Suspended"
	StatusTerminated Status = "Terminated"
	StatusError      Status = "Error"
)

// PowerState represents the power state of a service
type PowerState string

const (
	PowerStateOff       PowerState = "Off"
	PowerStateOn        PowerState = "On"
	PowerStateRebooting PowerState = "Rebooting"
	PowerStateUnknown   PowerState = "Unknown"
)

// ProvisioningEvent represents an event during service provisioning
type ProvisioningEvent struct {
	Priority  int32      `json:"priority"`
	Body      string     `json:"body,omitempty"`
	Timestamp time.Time  `json:"timestamp"`
	State     EventState `json:"state"`
}

// EventState represents the state of a provisioning event
type EventState string

const (
	EventStatePending    EventState = "Pending"
	EventStateInProgress EventState = "InProgress"
	EventStateComplete   EventState = "Complete"
	EventStateError      EventState = "Error"
)

// Result represents a standard API response with metadata and results
type Result[T any] struct {
	Success          bool         `json:"success"`
	Message          string       `json:"message,omitempty"`
	ValidationErrors []any        `json:"validationErrors,omitempty"`
	Metadata         ListMetadata `json:"metadata"`
	Result           T            `json:"result,omitempty"`
}

// ListMetadata contains metadata about list responses
type ListMetadata struct {
	TotalCount int32 `json:"total_count"`
	Limit      int32 `json:"limit"`
	Skip       int32 `json:"skip"`
}

// Now we can update our response types to use the generic Result
type (
	MetalResponse     Result[Metal]
	ListMetalResponse Result[[]Metal]
)

type ListMetalOptions struct {
	Skip  int32
	Limit int32

	Status    Status
	Region    string
	Tier      string
	Tag       string
	ProjectID int64
	TierType  MetalTierType
}

func (o *ListMetalOptions) ToQueryParams() []RequestOption {
	allOpts := []RequestOption{}
	if o.Skip > 0 {
		allOpts = append(allOpts, WithQueryParam("Skip", fmt.Sprint(o.Skip)))
	}
	if o.Limit > 0 {
		allOpts = append(allOpts, WithQueryParam("Limit", fmt.Sprint(o.Limit)))
	}
	if o.Status != "" {
		allOpts = append(allOpts, WithQueryParam("Status", string(o.Status)))
	}
	if o.Region != "" {
		allOpts = append(allOpts, WithQueryParam("Region", o.Region))
	}
	if o.Tier != "" {
		allOpts = append(allOpts, WithQueryParam("Tier", o.Tier))
	}
	if o.Tag != "" {
		allOpts = append(allOpts, WithQueryParam("Tag", o.Tag))
	}
	if o.ProjectID > 0 {
		allOpts = append(allOpts, WithQueryParam("ProjectId", fmt.Sprint(o.ProjectID)))
	}
	if o.TierType != "" {
		allOpts = append(allOpts, WithQueryParam("MetalTierType", string(o.TierType)))
	}
	return allOpts
}

// ListMetal retrieves a list of metal services with optional filtering
func (c *Client) ListMetal(ctx context.Context, opts ListMetalOptions) (*ListMetalResponse, error) {
	resp := &ListMetalResponse{}

	httpResp, err := c.Request(ctx, http.MethodGet, "Metal", nil, opts.ToQueryParams()...)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateBareMetalRequest represents the request parameters for creating a new metal service
type CreateBareMetalRequest struct {
	Quantity       int               `json:"quantity,omitempty"`
	DisplayName    string            `json:"displayName"`
	RegionID       string            `json:"regionId"`
	TierID         string            `json:"tierId"`
	MemoryGB       int               `json:"memoryGb"`
	ImageID        string            `json:"imageId"`
	Tags           []string          `json:"tags"`
	TemplateID     *int              `json:"templateId,omitempty"`
	SSHKeyIDs      []int             `json:"sshKeyIds"`
	Disks          map[string]string `json:"disks"`
	IPXEUrl        *string           `json:"ipxeUrl"`
	UserData       *string           `json:"userData"`
	Password       *string           `json:"password,omitempty"`
	ReservePricing bool              `json:"reservePricing"`
	Partitions     []Partition       `json:"partitions,omitempty"`
	RaidArrays     []RaidArray       `json:"raidArrays,omitempty"`
}

type Partition struct {
	Name      string `json:"name"`
	Device    string `json:"device"`
	SizeBytes *int64 `json:"sizeBytes,omitempty"`
}

type RaidArray struct {
	Name       string     `json:"name"`
	Type       RaidType   `json:"type"`
	Members    []string   `json:"members"`
	FileSystem FileSystem `json:"filesystem"`
	MountPoint string     `json:"mountPoint"`
}

type RaidType string

const (
	RaidTypeNone    RaidType = "None"
	RaidTypeRaid0   RaidType = "Raid0"
	RaidTypeRaid1   RaidType = "Raid1"
	RaidTypeUnknown RaidType = "Unknown"
)

type FileSystem string

const (
	FileSystemBtrfs       FileSystem = "Btrfs"
	FileSystemExt2        FileSystem = "Ext2"
	FileSystemExt4        FileSystem = "Ext4"
	FileSystemFat32       FileSystem = "Fat32"
	FileSystemRamfs       FileSystem = "Ramfs"
	FileSystemSwap        FileSystem = "Swap"
	FileSystemTmpfs       FileSystem = "Tmpfs"
	FileSystemUnformatted FileSystem = "Unformatted"
	FileSystemUnknown     FileSystem = "Unknown"
	FileSystemVfat        FileSystem = "Vfat"
	FileSystemXfs         FileSystem = "Xfs"
	FileSystemZfsroot     FileSystem = "Zfsroot"
)

// CreateMetalService creates a new metal service
func (c *Client) CreateMetalService(ctx context.Context, projectID int64, req *CreateBareMetalRequest) (*MetalResponse, error) {
	resp := &MetalResponse{}
	httpResp, err := c.Request(ctx, http.MethodPost, "Metal", req,
		WithQueryParam("projectId", fmt.Sprint(projectID)),
	)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// // ListMetalTemplates retrieves all metal service templates
// func (c *Client) ListMetalTemplates(ctx context.Context) (*MetalTemplateIEnumerableApiResponse, error) {
// 	resp := &MetalTemplateIEnumerableApiResponse{}
// 	httpResp, err := c.Request(ctx, http.MethodGet, "Metal/templates", nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer httpResp.Body.Close()

// 	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
// 		return nil, err
// 	}
// 	return resp, nil
// }

// MetalTier represents a hardware configuration tier for metal services
type MetalTier struct {
	ID                 string                          `json:"id"`             // ID of the tier
	CPU                string                          `json:"cpu"`            // The CPU model for the tier
	CPUDescription     string                          `json:"cpuDescription"` // Description of CPU (cores/threads)
	ExternalIdentifier string                          `json:"externalIdentifier"`
	Hidden             bool                            `json:"hidden"`
	Availability       map[string]*ServiceAvailability `json:"availability"`
	MemoryOptions      []MemoryOption                  `json:"memoryOptions"`  // Available memory configurations
	DriveSlots         []DriveSlot                     `json:"driveSlots"`     // Available drive slots
	NetworkOptions     []NetworkOption                 `json:"networkOptions"` // Available network configurations
	MonthlyPrice       float64                         `json:"monthlyPrice"`   // The monthly price for the tier
	HourlyPrice        float64                         `json:"hourlyPrice"`    // The hourly price for the tier
	MemoryOptionSetID  int                             `json:"memoryOptionSetId"`
	DriveSlotSetID     int                             `json:"driveSlotSetId"`
	NetworkOptionSetID int                             `json:"networkOptionSetId"`
	TierType           MetalTierType                   `json:"tierType"`
}

// ServiceAvailability represents the availability of a service in a region
type ServiceAvailability struct {
	MaxQuantity int `json:"maxQuantity"`
}

// MemoryOption represents a memory configuration option
type MemoryOption struct {
	GB           int     `json:"gb"`           // Amount of memory in Gigabytes
	MonthlyPrice float64 `json:"monthlyPrice"` // Monthly price for this memory option
	HourlyPrice  float64 `json:"hourlyPrice"`  // Hourly price for this memory option
	Default      bool    `json:"default"`      // If this is the default memory option
}

// DriveSlot represents a slot for a storage drive
type DriveSlot struct {
	ID       string               `json:"id"`       // ID in format like nvme0n1
	Default  string               `json:"default"`  // Default value for this drive slot
	Required bool                 `json:"required"` // Whether a drive is required in this slot
	Options  []MetalStorageDevice `json:"options"`  // Available drives for this slot
}

// MetalStorageDevice describes a drive which can be installed in a drive slot
type MetalStorageDevice struct {
	Name         string       `json:"name"`         // Name of the drive (e.g. "1.92t" or "960g")
	Default      bool         `json:"default"`      // True if this is the default drive for the tier
	Type         StorageType  `json:"type"`         // Type of storage (HDD, SSD, NVME)
	CapacityGB   int          `json:"capacityGb"`   // The capacity of the drive in Gigabytes
	Details      DriveDetails `json:"details"`      // Additional drive details
	MonthlyPrice float64      `json:"monthlyPrice"` // The monthly price for the specific drive
	HourlyPrice  float64      `json:"hourlyPrice"`  // The hourly price for the specific drive
	IsBossDrive  bool         `json:"isBossDrive"`  // True if this drive is a boss drive
	ID           string       `json:"id,omitempty"`
}

// DriveDetails contains additional information about a drive
type DriveDetails struct {
	DeviceName string `json:"deviceName,omitempty"`
	Serial     string `json:"serial,omitempty"`
}

// StorageType describes the storage type of the drive
type StorageType string

const (
	StorageTypeHDD     StorageType = "HDD"
	StorageTypeSSD     StorageType = "SSD"
	StorageTypeNVME    StorageType = "NVME"
	StorageTypeUnknown StorageType = "Unknown"
)

// NetworkOption represents a network configuration option
type NetworkOption struct {
	SpeedGbps    int     `json:"speedGbps"`    // Speed in Gbps
	MonthlyPrice float64 `json:"monthlyPrice"` // Monthly price for this network option
	HourlyPrice  float64 `json:"hourlyPrice"`  // Hourly price for this network option
	Default      bool    `json:"default"`      // If this is the default option
	IsBonded     bool    `json:"isBonded"`     // If this network option is bonded
}

// MetalTierType represents the type of metal tier
type MetalTierType string

const (
	MetalTierTypeCompute MetalTierType = "Compute"
	MetalTierTypeGPU     MetalTierType = "GPU"
)

// Update the response type for ListMetalTiers
type MetalTierResponse Result[[]MetalTier]

// ListMetalTiers retrieves all metal service tiers with optional type filter
func (c *Client) ListMetalTiers(ctx context.Context, tierType MetalTierType) (*MetalTierResponse, error) {
	resp := &MetalTierResponse{}
	opts := []RequestOption{}
	if tierType != "" {
		opts = append(opts, WithQueryParam("metalTierType", string(tierType)))
	}

	httpResp, err := c.Request(ctx, http.MethodGet, "Metal/tiers", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, errors.New(resp.Message)
	}

	return resp, nil
}

// GetMetalService retrieves a single metal service by ID
func (c *Client) GetMetalService(ctx context.Context, id int64) (*MetalResponse, error) {
	resp := &MetalResponse{}
	httpResp, err := c.Request(ctx, http.MethodGet, fmt.Sprintf("Metal/%d", id), nil)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, errors.New(resp.Message)
	}

	return resp, nil
}

type ReinstallMetalRequest struct {
	DisplayName     string           `json:"displayName,omitempty"`
	ImageID         string           `json:"imageId,omitempty"`
	SSHKeyIDs       []int64           `json:"sshKeyIds,omitempty"`
	Password        string           `json:"password,omitempty"`
	UserData        string           `json:"userData,omitempty"`
	IPXEUrl         string           `json:"ipxeUrl,omitempty"`
	Partitions      []Partition       `json:"partitions,omitempty"`
	RaidArrays      []RaidArray       `json:"raidArrays,omitempty"`
}

// ReinstallMetalService reinstalls a metal service by ID
func (c *Client) ReinstallMetalService(ctx context.Context, id int64, req *ReinstallMetalRequest) (*MetalResponse, error) {
	resp := &MetalResponse{}
	httpResp, err := c.Request(ctx, http.MethodPost, fmt.Sprintf("Metal/%d/Reinstall", id), req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

type PowerCommand string

const (
	PowerCommandPowerOff PowerCommand = "PowerOff"
	PowerCommandPowerOn  PowerCommand = "PowerOn"
)

// SendPowerCommand sends a power command to a metal service
func (c *Client) SendPowerCommand(ctx context.Context, id int64, command PowerCommand) (*MetalResponse, error) {
	resp := &MetalResponse{}
	httpResp, err := c.Request(ctx, http.MethodPost, fmt.Sprintf("Metal/%d/PowerCommand", id), nil, WithQueryParam("command", fmt.Sprint(command)))
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, errors.New(resp.Message)
	}

	return resp, nil
}

// LogMessage represents a log entry for a metal service
type LogMessage struct {
	Timestamp string `json:"timestamp"`
	Name      string `json:"name,omitempty"`
	Message   string `json:"message,omitempty"`
}

// LogMessageResponse represents a response containing log messages
type LogMessageResponse Result[[]LogMessage]

// GetMetalLogs retrieves logs for a metal service
func (c *Client) GetMetalLogs(ctx context.Context, id int64) (*LogMessageResponse, error) {
	resp := &LogMessageResponse{}
	httpResp, err := c.Request(ctx, http.MethodGet, fmt.Sprintf("Metal/%d/Logs", id), nil)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, errors.New(resp.Message)
	}

	return resp, nil
}

// MetalConfiguration represents the availability configuration for metal services
type MetalConfiguration struct {
	Disks    map[string]string `json:"disks"`    // Dictionary of disk names and sizes
	MemoryGB int               `json:"memoryGb"` // Amount of memory in GB
	Tier     MetalTier         `json:"tier"`     // The tier configuration
	Quantity int               `json:"quantity"` // Available quantity
}

// MetalConfigurationResponse represents a response containing metal configuration data
type MetalConfigurationResponse Result[[]MetalConfiguration]

// Update the GetMetalAvailability function to use the new response type
func (c *Client) GetMetalAvailability(ctx context.Context, projectId int64, regionId string) (*MetalConfigurationResponse, error) {
	allOpts := []RequestOption{
		WithQueryParam("ProjectId", fmt.Sprint(projectId)),
		WithQueryParam("Region", regionId),
	}

	resp := &MetalConfigurationResponse{}
	httpResp, err := c.Request(ctx, http.MethodGet, "Metal/Availability", nil, allOpts...)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, errors.New(resp.Message)
	}

	return resp, nil
}

// RenameMetalService renames a metal service
func (c *Client) RenameMetalService(ctx context.Context, id int64, name string) (*Result[struct{}], error) {
	resp := &Result[struct{}]{}
	httpResp, err := c.Request(ctx, http.MethodPost, fmt.Sprintf("Metal/%d/rename", id), map[string]string{"name": name})
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, errors.New(resp.Message)
	}

	return resp, nil
}
