package config

import (
	serverConfig "eftep/pkg/config/server"
	"time"
)

var (
	DOWNLOAD_DIR           = "/tmp/eftepcli"
	MULTICAST_GROUPS       = serverConfig.MULTICAST_GROUPS
	DISCOVERY_BIND_PORT    = 5000
	DISCOVERY_BIND_IP_ADDR = [4]byte{0, 0, 0, 0}

	DISCOVERY_SERVER_PORT = serverConfig.DISCOVERY_PORT

	RECV_TIMEOUT_SECS = time.Second
)
