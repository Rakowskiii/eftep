package config

var (
	IP_ADDR        = [4]byte{0, 0, 0, 0}
	EFTEP_PORT     = 8080
	DISCOVERY_PORT = 8081
	WORKDIR        = "/tmp/eftep"
	LOGFILE        = "/var/log/eftep.log"
)

var MULTICAST_GROUPS = [2][4]byte{
	{224, 0, 1, 1},
	{224, 0, 1, 2},
}
