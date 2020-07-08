package tormdr

type Config struct {
	TorMDRBinaryPath    string
	DataDirectory       string
	HardwareAccel       bool
	KeepalivePeriod     int
	UseSocks5Proxy      bool
	Socks5ProxyAddress  string
	Socks5ProxyUserName string
	Socks5ProxyPassword string
	UseObfs4Proxy       bool
	Obfs4ProxyPath      string
	Obfs4Bridges        []string
}

var defaultArgs = []string{"",
	"RunAsDaemon", "0",
	"ClientOnly", "1",
	"AvoidDiskWrites", "1",
	"FetchHidServDescriptors", "0",
	"FetchServerDescriptors", "1",
	"FetchUselessDescriptors", "0",
	"UseEntryGuards", "0",
	"NumEntryGuards", "0",
	"UseGuardFraction", "0",
	"DownloadExtraInfo", "0",
	"UseMicrodescriptors", "1",
	"ClientUseIPv4", "1",
	"ClientUseIPv6", "0",
	"DirCache", "0",
	"NewCircuitPeriod", "90 days",
	"MaxCircuitDirtiness", "30 days",
	"EnforceDistinctSubnets", "0",
}
