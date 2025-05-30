package manager

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/XMPlusDev/XMPlus/app/xdispatcher"

	"dario.cat/mergo"
	"github.com/r3labs/diff/v2"
	"github.com/xmplusdev/xray-core/v25/app/proxyman"
	"github.com/xmplusdev/xray-core/v25/app/stats"
	"github.com/xmplusdev/xray-core/v25/common/serial"
	"github.com/xmplusdev/xray-core/v25/core"
	"github.com/xmplusdev/xray-core/v25/infra/conf"

	"github.com/XMPlusDev/XMPlus/api"
	"github.com/XMPlusDev/XMPlus/api/xmplus"
	_ "github.com/XMPlusDev/XMPlus/cmd/distro/all"
	"github.com/XMPlusDev/XMPlus/service"
	"github.com/XMPlusDev/XMPlus/service/controller"
)

// Manager Structure
type Manager struct {
	statusLock    sync.Mutex
	managerConfig *Config
	Server        *core.Instance
	Service       []service.Service
	Running       bool
}

func New(managerConfig *Config) *Manager {
	m := &Manager{managerConfig: managerConfig}
	return m
}

func (m *Manager) loadCore(managerConfig *Config) *core.Instance {
	// Log Config
	coreLogConfig := &conf.LogConfig{}
	logConfig := getDefaultLogConfig()
	if managerConfig.LogConfig != nil {
		if _, err := diff.Merge(logConfig, managerConfig.LogConfig, logConfig); err != nil {
			log.Panicf("Read Log config failed: %s", err)
		}
	}
	coreLogConfig.LogLevel = logConfig.Level
	coreLogConfig.AccessLog = logConfig.AccessPath
	coreLogConfig.ErrorLog = logConfig.ErrorPath
	coreLogConfig.DNSLog = logConfig.DNSLog
	coreLogConfig.MaskAddress = logConfig.MaskAddress

	// DNS config
	coreDnsConfig := &conf.DNSConfig{}
	if managerConfig.DnsConfigPath != "" {
		if data, err := os.ReadFile(managerConfig.DnsConfigPath); err != nil {
			log.Panicf("Failed to read DNS config file at: %s", managerConfig.DnsConfigPath)
		} else {
			if err = json.Unmarshal(data, coreDnsConfig); err != nil {
				log.Panicf("Failed to unmarshal DNS config: %s", managerConfig.DnsConfigPath)
			}
		}
	}

	dnsConfig, err := coreDnsConfig.Build()
	if err != nil {
		log.Panicf("Failed to understand DNS config, Please check: https://xtls.github.io/config/dns.html for help: %s", err)
	}

	// Routing config
	coreRouterConfig := &conf.RouterConfig{}
	if managerConfig.RouteConfigPath != "" {
		if data, err := os.ReadFile(managerConfig.RouteConfigPath); err != nil {
			log.Panicf("Failed to read Routing config file at: %s", managerConfig.RouteConfigPath)
		} else {
			if err = json.Unmarshal(data, coreRouterConfig); err != nil {
				log.Panicf("Failed to unmarshal Routing config: %s", managerConfig.RouteConfigPath)
			}
		}
	}
	routeConfig, err := coreRouterConfig.Build()
	if err != nil {
		log.Panicf("Failed to understand Routing config  Please check: https://xtls.github.io/config/routing.html for help: %s", err)
	}

	// Custom Inbound config
	var coreCustomInboundConfig []conf.InboundDetourConfig
	if managerConfig.InboundConfigPath != "" {
		if data, err := os.ReadFile(managerConfig.InboundConfigPath); err != nil {
			log.Panicf("Failed to read Custom Inbound config file at: %s", managerConfig.OutboundConfigPath)
		} else {
			if err = json.Unmarshal(data, &coreCustomInboundConfig); err != nil {
				log.Panicf("Failed to unmarshal Custom Inbound config: %s", managerConfig.OutboundConfigPath)
			}
		}
	}
	var inBoundConfig []*core.InboundHandlerConfig
	for _, config := range coreCustomInboundConfig {
		oc, err := config.Build()
		if err != nil {
			log.Panicf("Failed to understand Inbound config, Please check: https://xtls.github.io/config/inbound.html for help: %s", err)
		}
		inBoundConfig = append(inBoundConfig, oc)
	}

	// Custom Outbound config
	var coreCustomOutboundConfig []conf.OutboundDetourConfig
	if managerConfig.OutboundConfigPath != "" {
		if data, err := os.ReadFile(managerConfig.OutboundConfigPath); err != nil {
			log.Panicf("Failed to read Custom Outbound config file at: %s", managerConfig.OutboundConfigPath)
		} else {
			if err = json.Unmarshal(data, &coreCustomOutboundConfig); err != nil {
				log.Panicf("Failed to unmarshal Custom Outbound config: %s", managerConfig.OutboundConfigPath)
			}
		}
	}
	var outBoundConfig []*core.OutboundHandlerConfig
	for _, config := range coreCustomOutboundConfig {
		oc, err := config.Build()
		if err != nil {
			log.Panicf("Failed to understand Outbound config, Please check: https://xtls.github.io/config/outbound.html for help: %s", err)
		}
		outBoundConfig = append(outBoundConfig, oc)
	}

	// Policy config
	levelPolicyConfig := parseConnectionConfig(managerConfig.ConnectionConfig)
	corePolicyConfig := &conf.PolicyConfig{}
	corePolicyConfig.Levels = map[uint32]*conf.Policy{0: levelPolicyConfig}
	policyConfig, _ := corePolicyConfig.Build()

	// Reverse config
	coreReverseConfig := &conf.ReverseConfig{}
	if managerConfig.ReverseConfigPath != "" {
		if data, err := os.ReadFile(managerConfig.ReverseConfigPath); err != nil {
			log.Panicf("Failed to read Reverse config file at: %s", managerConfig.ReverseConfigPath)
		} else {
			if err = json.Unmarshal(data, coreReverseConfig); err != nil {
				log.Panicf("Failed to unmarshal Reverse config: %s", managerConfig.ReverseConfigPath)
			}
		}
	}
	reverseConfig, err := coreReverseConfig.Build()
	if err != nil {
		log.Panicf("Failed to understand Reverse config, Please check: https://xtls.github.io/config/reverse.html for help: %s", err)
	}

	// API config
	coreAPIConfig := &conf.APIConfig{}
	var apiConfigMsg *serial.TypedMessage
	if managerConfig.ApiConfigPath != "" {
		if data, err := os.ReadFile(managerConfig.ApiConfigPath); err != nil {
			log.Panicf("Failed to read API config file at: %s", managerConfig.ApiConfigPath)
		} else {
			if err = json.Unmarshal(data, coreAPIConfig); err != nil {
				log.Panicf("Failed to unmarshal API config: %s", managerConfig.ApiConfigPath)
			}
		}
		apiConfig, err := coreAPIConfig.Build()
		if err != nil {
			log.Panicf("Failed to understand API config, Please check: https://xtls.github.io/config/api.html for help: %s", err)
		}
		apiConfigMsg = serial.ToTypedMessage(apiConfig)
	}

	// Build Core Config
	config := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(coreLogConfig.Build()),
			serial.ToTypedMessage(&xdispatcher.Config{}),
			serial.ToTypedMessage(&stats.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
			serial.ToTypedMessage(policyConfig),
			serial.ToTypedMessage(dnsConfig),
			serial.ToTypedMessage(routeConfig),
			serial.ToTypedMessage(reverseConfig),
		},
		Inbound:  inBoundConfig,
		Outbound: outBoundConfig,
	}
	if apiConfigMsg != nil {
		config.App = append(config.App, apiConfigMsg)
	}

	server, err := core.New(config)
	if err != nil {
		log.Panicf("failed to create instance: %s", err)
	}

	//log.Printf("Core Version: %s", core.Version())

	return server
}

// Start the manager
func (m *Manager) Start() {
	m.statusLock.Lock()
	defer m.statusLock.Unlock()
	// Load Core
	server := m.loadCore(m.managerConfig)
	if err := server.Start(); err != nil {
		log.Panicf("Failed to start instance: %s", err)
	}
	m.Server = server

	// Load Nodes config
	for _, nodeConfig := range m.managerConfig.NodesConfig {
		var apiClient api.API
		apiClient = xmplus.New(nodeConfig.ApiConfig)

		var controllerService service.Service
		// Register controller service
		controllerConfig := getDefaultControllerConfig()
		if nodeConfig.ControllerConfig != nil {
			if err := mergo.Merge(controllerConfig, nodeConfig.ControllerConfig, mergo.WithOverride); err != nil {
				log.Panicf("Read Controller Config Failed")
			}
		}
		controllerService = controller.New(server, apiClient, controllerConfig)
		m.Service = append(m.Service, controllerService)

	}

	// Start all the service
	for _, s := range m.Service {
		err := s.Start()
		if err != nil {
			log.Panicf("XMPlus fialed to start: %s", err)
		}
	}
	m.Running = true
	return
}

// Close the manager
func (m *Manager) Close() {
	m.statusLock.Lock()
	defer m.statusLock.Unlock()
	for _, s := range m.Service {
		err := s.Close()
		if err != nil {
			log.Panicf("XMPlus fialed to close: %s", err)
		}
	}
	m.Service = nil
	m.Server.Close()
	m.Running = false
	return
}

func parseConnectionConfig(c *ConnectionConfig) (policy *conf.Policy) {
	connectionConfig := getDefaultConnectionConfig()
	if c != nil {
		if _, err := diff.Merge(connectionConfig, c, connectionConfig); err != nil {
			log.Panicf("Read ConnectionConfig failed: %s", err)
		}
	}
	policy = &conf.Policy{
		StatsUserUplink:   true,
		StatsUserDownlink: true,
		Handshake:         &connectionConfig.Handshake,
		ConnectionIdle:    &connectionConfig.ConnIdle,
		UplinkOnly:        &connectionConfig.UplinkOnly,
		DownlinkOnly:      &connectionConfig.DownlinkOnly,
		BufferSize:        &connectionConfig.BufferSize,
	}

	return
}
