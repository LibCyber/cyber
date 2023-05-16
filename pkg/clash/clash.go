package clash

import (
	clashConfig "github.com/Dreamacro/clash/config"
	clashLog "github.com/Dreamacro/clash/log"
	clashT "github.com/Dreamacro/clash/tunnel"

	"gopkg.in/yaml.v3"
)

type ClashRawConfig struct {
	Port               int               `yaml:"port"`
	SocksPort          int               `yaml:"socks-port"`
	RedirPort          int               `yaml:"redir-port,omitempty"`
	TProxyPort         int               `yaml:"tproxy-port,omitempty"`
	MixedPort          int               `yaml:"mixed-port,omitempty"`
	Authentication     []string          `yaml:"authentication,omitempty"`
	AllowLan           bool              `yaml:"allow-lan"`
	BindAddress        string            `yaml:"bind-address,omitempty"`
	Mode               clashT.TunnelMode `yaml:"mode"`
	LogLevel           clashLog.LogLevel `yaml:"log-level"`
	IPv6               bool              `yaml:"ipv6,omitempty"`
	ExternalController string            `yaml:"external-controller"`
	ExternalUI         string            `yaml:"external-ui,omitempty"`
	Secret             string            `yaml:"secret,omitempty"`
	Interface          string            `yaml:"interface-name,omitempty"`
	RoutingMark        int               `yaml:"routing-mark,omitempty"`

	ProxyProvider map[string]map[string]any `yaml:"proxy-providers,omitempty"`
	Hosts         map[string]string         `yaml:"hosts,omitempty"`
	DNS           clashConfig.RawDNS        `yaml:"dns,omitempty"`
	Experimental  clashConfig.Experimental  `yaml:"experimental,omitempty"`
	Profile       clashConfig.Profile       `yaml:"profile,omitempty"`
	Proxy         []map[string]any          `yaml:"proxies"`
	ProxyGroup    []map[string]any          `yaml:"proxy-groups"`
	Rule          []string                  `yaml:"rules"`

	Tun *TunConfig `yaml:"tun,omitempty"`
}

type TunConfig struct {
	Enable              bool     `yaml:"enable"`
	Stack               string   `yaml:"stack"`
	AutoRoute           bool     `yaml:"auto-route"`
	AutoDetectInterface bool     `yaml:"auto-detect-interface"`
	DnsHijack           []string `yaml:"dns-hijack"`
}

func UnmarshalClashRawConfig(buf []byte) (*ClashRawConfig, error) {
	// config with default value
	rawCfg := &ClashRawConfig{
		AllowLan: false,
		//BindAddress:    "*",
		Mode:           clashT.Rule,
		Authentication: []string{},
		LogLevel:       clashLog.INFO,
		Hosts:          map[string]string{},
		Rule:           []string{},
		Proxy:          []map[string]any{},
		ProxyGroup:     []map[string]any{},
		//DNS: &clashConfig.RawDNS{
		//	Enable:      false,
		//	UseHosts:    true,
		//	FakeIPRange: "198.18.0.1/16",
		//	FallbackFilter: clashConfig.RawFallbackFilter{
		//		GeoIP:     true,
		//		GeoIPCode: "CN",
		//		IPCIDR:    []string{},
		//	},
		//	DefaultNameserver: []string{
		//		"114.114.114.114",
		//		"8.8.8.8",
		//	},
		//},
		//Profile: &clashConfig.Profile{
		//	StoreSelected: true,
		//},
	}

	if err := yaml.Unmarshal(buf, rawCfg); err != nil {
		return nil, err
	}

	return rawCfg, nil
}

//goland:noinspection GoUnusedExportedFunction
func CheckIfConfigValid(content string) bool {
	_, err := UnmarshalClashRawConfig([]byte(content))
	if err != nil {
		return false
	}
	return true
}

// MergeRules 合并追加新规则
func MergeRules(oldRules []string, newRules []string) []string {
	var rules = make([]string, 0, len(oldRules))
	for _, rule := range oldRules {
		rules = append(rules, rule)
	}
	for _, rule := range newRules {
		if !contains(oldRules, rule) {
			rules = append(rules, rule)
		}
	}
	return rules
}

func contains(rules []string, rule string) bool {
	for _, r := range rules {
		if r == rule {
			return true
		}
	}
	return false
}
