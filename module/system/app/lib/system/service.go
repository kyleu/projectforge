package system

import (
	"context"
	"time"

	"github.com/samber/lo"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"

	"{{{ .Package }}}/app/util"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Status(ctx context.Context, logger util.Logger) (*Status, error) {
	logger.Debugf("calculating system status...")
	return &Status{
		CPU:     s.cpuStatus(ctx),
		Memory:  s.memStatus(ctx),
		Process: s.processStatus(ctx),
		Disk:    s.diskStatus(ctx),
		Host:    s.hostStatus(ctx),
		Extra:   s.extraStatus(ctx),
	}, nil
}

func (s *Service) cpuStatus(ctx context.Context) *CPUStatus {
	cc, _ := cpu.CountsWithContext(ctx, false)
	cl, _ := cpu.CountsWithContext(ctx, true)
	ci, _ := cpu.InfoWithContext(ctx)
	pc, _ := cpu.PercentWithContext(ctx, 0, true)
	pt, _ := cpu.TimesWithContext(ctx, true)
	ld, _ := load.AvgWithContext(ctx)
	return &CPUStatus{Count: cc, Logical: cl, Info: ci, Percent: pc, Times: util.ArrayFromAny[any](pt), Load: ld}
}

func (s *Service) memStatus(ctx context.Context) *MemoryStatus {
	vm, _ := mem.VirtualMemoryWithContext(ctx)
	sm, _ := mem.SwapMemoryWithContext(ctx)
	return &MemoryStatus{Virtual: vm, Swap: sm}
}

func (s *Service) hostStatus(ctx context.Context) *HostStatus {
	bt, _ := host.BootTimeWithContext(ctx)
	hi, _ := host.InfoWithContext(ctx)
	return &HostStatus{BootTime: time.Unix(int64(bt), 0), Info: hi}
}

func (s *Service) diskStatus(ctx context.Context) *DiskStatus {
	pts, _ := disk.PartitionsWithContext(ctx, false)
	us, _ := util.AsyncCollect(pts, func(x disk.PartitionStat) (any, error) {
		return disk.UsageWithContext(ctx, x.Mountpoint)
	})
	return &DiskStatus{Partitions: util.ArrayFromAny[any](pts), Usage: us}
}

func (s *Service) processStatus(ctx context.Context) *ProcessStatus {
	procs, _ := process.ProcessesWithContext(ctx)
	ret := lo.Map(procs, func(x *process.Process, _ int) int32 {
		return x.Pid
	})
	return &ProcessStatus{Processes: util.ArraySorted(ret)}
}

func (s *Service) extraStatus(ctx context.Context) util.ValueMap {
	return util.ValueMap{}
}
