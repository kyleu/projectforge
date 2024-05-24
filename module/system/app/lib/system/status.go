package system

import (
	"time"
)

type Status struct {
	CPU     *CPUStatus     `json:"cpu,omitempty"`
	Memory  *MemoryStatus  `json:"memory,omitempty"`
	Process *ProcessStatus `json:"process,omitempty"`
	Disk    *DiskStatus    `json:"disk,omitempty"`
	Host    *HostStatus    `json:"host,omitempty"`
	Extra   any            `json:"extra,omitempty"`
}

type CPUStatus struct {
	Count   int       `json:"count,omitempty"`
	Logical int       `json:"logical,omitempty"`
	Info    any       `json:"info,omitempty"`
	Percent []float64 `json:"percent,omitempty"`
	Times   []any     `json:"times,omitempty"`
	Load    any       `json:"load,omitempty"`
}

type MemoryStatus struct {
	Virtual any `json:"virtual,omitempty"`
	Swap    any `json:"swap,omitempty"`
}

type DiskStatus struct {
	Partitions []any `json:"partitions,omitempty"`
	Usage      []any `json:"usage,omitempty"`
}

type ProcessStatus struct {
	Processes any `json:"x,omitempty"`
}

type HostStatus struct {
	BootTime time.Time `json:"bootTime,omitempty"`
	Info     any       `json:"info,omitempty"`
}
