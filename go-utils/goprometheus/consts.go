package goprometheus

import (
	"github.com/prometheus/common/model"
)

type MetricName = string
type MetricLabel = model.LabelName

type MetricQuery = struct {
	ProductCode string
	Group       string
	InstanceIds []int64
	TimeRange   string
}

var (
	MetricVersion                  = "version"
	MetricStartupTimestamp         = "startup_timestamp"
	MetricAliveTimestamp           = "alive_timestamp"
	MetricConfigUpdateTimestamp    = "config_update_timestamp"
	MetricRunTime                  = "run_time"
	MetricErrCode                  = "err_code"
	MetricTotalBandwidthIn         = "total_bandwidth_in"
	MetricTotalBandwidthOut        = "total_bandwidth_out"
	MetricMemberLevelUserCount     = "member_level_user_count"
	MetricRealMemberLevelUserCount = "real_member_level_user_count"
	MetricIpTcpConnectionsIn       = "ip_tcp_connections_in"
	MetricIpTcpConnectionsOut      = "ip_tcp_connections_out"
	MetricIpUdpConnectionsIn       = "ip_udp_connections_in"
	MetricIpUdpPortsOut            = "ip_udp_ports_out"
	MetricIpBandwidthIn            = "ip_bandwidth_in"
	MetricIpBandwidthOut           = "ip_bandwidth_out"
	MetricIpUserCount              = "ip_user_count"
	MetricNodeExporterBuildInfo    = "node_exporter_build_info"
	MetricNodeCpuSecondsTotal      = "node_cpu_seconds_total"
	MetricNodeCpuGuessSecondsTotal = "node_cpu_guess_seconds_total"
	MetricNodeMemTotal             = "node_memory_MemTotal_bytes"
	MetricNodeMemAvailable         = "node_memory_MemAvailable_bytes"
	MetricNodeDiskSize             = "node_filesystem_size_bytes"
	MetricNodeDiskFree             = "node_filesystem_free_bytes"
	MetricNodeDiskAvailable        = "node_filesystem_avail_bytes"
	MetricNodeDiskInode            = "node_filesystem_files"
	MetricNodeDiskInodeFree        = "node_filesystem_files_free"
	MetricNodeDiskReadonly         = "node_filesystem_readonly"
	MetricNodeDiskDeviceError      = "node_filesystem_device_error"
	MetricNodeTrafficIn            = "node_network_receive_bytes_total"
	MetricNodeTrafficOut           = "node_network_transmit_bytes_total"
	MetricUp                       = "up"

	// node_memory_Active_anon_bytes
	// node_memory_Active_bytes
	// node_memory_Active_file_bytes
	// node_memory_AnonHugePages_bytes
	// node_memory_AnonPages_bytes
	// node_memory_Bounce_bytes
	// node_memory_Buffers_bytes
	// node_memory_Cached_bytes
	// node_memory_CommitLimit_bytes
	// node_memory_Committed_AS_bytes
	// node_memory_DirectMap1G_bytes
	// node_memory_DirectMap2M_bytes
	// node_memory_DirectMap4k_bytes
	// node_memory_Dirty_bytes
	// node_memory_FileHugePages_bytes
	// node_memory_FilePmdMapped_bytes
	// node_memory_HardwareCorrupted_bytes
	// node_memory_HugePages_Free
	// node_memory_HugePages_Rsvd
	// node_memory_HugePages_Surp
	// node_memory_HugePages_Total
	// node_memory_Hugepagesize_bytes
	// node_memory_Hugetlb_bytes
	// node_memory_Inactive_anon_bytes
	// node_memory_Inactive_bytes
	// node_memory_Inactive_file_bytes
	// node_memory_KReclaimable_bytes
	// node_memory_KernelStack_bytes
	// node_memory_Mapped_bytes
	// node_memory_MemAvailable_bytes
	// node_memory_MemFree_bytes
	// node_memory_MemTotal_bytes
	// node_memory_Mlocked_bytes
	// node_memory_NFS_Unstable_bytes
	// node_memory_PageTables_bytes
	// node_memory_Percpu_bytes
	// node_memory_SReclaimable_bytes
	// node_memory_SUnreclaim_bytes
	// node_memory_ShmemHugePages_bytes
	// node_memory_ShmemPmdMapped_bytes
	// node_memory_Shmem_bytes
	// node_memory_Slab_bytes
	// node_memory_SwapCached_bytes
	// node_memory_SwapFree_bytes
	// node_memory_SwapTotal_bytes
	// node_memory_Unevictable_bytes
	// node_memory_VmallocChunk_bytes
	// node_memory_VmallocTotal_bytes
	// node_memory_VmallocUsed_bytes
	// node_memory_WritebackTmp_bytes
	// node_memory_Writeback_bytes
	// node_network_receive_bytes_total
	// node_network_receive_compressed_total
	// node_network_receive_drop_total
	// node_network_receive_errs_total
	// node_network_receive_fifo_total
	// node_network_receive_frame_total
	// node_network_receive_multicast_total
	// node_network_receive_nohandler_total
	// node_network_receive_packets_total
	// node_network_transmit_bytes_total
	// node_network_transmit_carrier_total
	// node_network_transmit_colls_total
	// node_network_transmit_compressed_total
	// node_network_transmit_drop_total
	// node_network_transmit_errs_total
	// node_network_transmit_fifo_total
	// node_network_transmit_packets_total
	// node_scrape_collector_duration_seconds
	// node_scrape_collector_success

)

var (
	MetricLabelName               = model.LabelName("__name__")
	MetricLabelServerType         = model.LabelName("server_type")
	MetricLabelJob                = model.LabelName("job")
	MetricLabelGroup              = model.LabelName("group")
	MetricLabelProductIds         = model.LabelName("product_ids")
	MetricLabelSupplierId         = model.LabelName("supplier_id")
	MetricLabelRetry              = model.LabelName("retry")
	MetricLabelInstanceId         = model.LabelName("instance_id")
	MetricLabelCpu                = model.LabelName("cpu")
	MetricLabelExportedInstanceId = model.LabelName("exported_instance_id")
	MetricLabelInstance           = model.LabelName("instance")
	MetricLabelMemberLevel        = model.LabelName("member_level")
	MetricLabelIp                 = model.LabelName("ip")
	MetricLabelProtocol           = model.LabelName("protocol")
	MetricLabelPort               = model.LabelName("port")
)
