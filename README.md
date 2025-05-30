# XMPlus

#### Config directory
```
/etc/XMPlus
```

### Onclick XMPlus backennd Install
```
bash <(curl -Ls https://raw.githubusercontent.com/XMPlusDev/XMPlus/scripts/install.sh)
```

### /etc/XMPlus/config.yml
```
Log:
  Level: warning # Log level: none, error, warning, info, debug 
  AccessPath: # /etc/XMPlus/access.Log
  ErrorPath: # /etc/XMPlus/error.log
  DNSLog: false # false or true
  MaskAddress: half # half, full, quarter
DnsConfigPath:  #/etc/XMPlus/dns.json
RouteConfigPath: # /etc/XMPlus/route.json
InboundConfigPath: # /etc/XMPlus/inbound.json
OutboundConfigPath: # /etc/XMPlus/outbound.json
ConnectionConfig:
  Handshake: 8 
  ConnIdle: 300 
  UplinkOnly: 0 
  DownlinkOnly: 0 
  BufferSize: 64
Nodes:
  -
    ApiConfig:
      ApiHost: "https://www.xyz.com"
      ApiKey: "123"
      NodeID: 1
      Timeout: 30 
      RuleListPath: # /etc/XMPlus/rulelist Path to local rulelist file
    ControllerConfig:
      EnableDNS: false # Use custom DNS config, Please ensure that you set the dns.json well
      DNSStrategy: AsIs # AsIs, UseIP, UseIPv4, UseIPv6
      CertConfig:
        Email: author@xmplus.dev                    # Required when Cert Mode is not none
        CertFile: /etc/XMPlus/node1.xmplus.dev.crt  # Required when Cert Mode is file
        KeyFile: /etc/XMPlus/node1.xmplus.dev.key   # Required when Cert Mode is file
        Provider: cloudflare                        # Required when Cert Mode is dns
        CertEnv:                                    # Required when Cert Mode is dns
          CLOUDFLARE_EMAIL:                         # Required when Cert Mode is dns
          CLOUDFLARE_API_KEY:                       # Required when Cert Mode is dns
      EnableFallback: false # Only support for Trojan and Vless
      FallBackConfigs:  # Support multiple fallbacks
        - SNI: # TLS SNI(Server Name Indication), Empty for any
          Alpn: # Alpn, Empty for any
          Path: # HTTP PATH, Empty for any
          Dest: 80 # Required, Destination of fallback, check https://xtls.github.io/config/features/fallback.html for details.
          ProxyProtocolVer: 0 # Send PROXY protocol version, 0 for disable
      IPLimit:
        Enable: false # Enable the global ip limit of a user 
        RedisNetwork: tcp # Redis protocol, tcp or unix
        RedisAddr: 127.0.0.1:6379 # Redis server address, or unix socket path
        RedisUsername: default # Redis username
        RedisPassword: YOURPASSWORD # Redis password
        RedisDB: 0 # Redis DB
        Timeout: 5 # Timeout for redis request
        Expiry: 60 # Expiry time (second) 
```

## XMPlus Panel Server configuration

### Network Settings

#### TCP
```
{
  "transport" : "tcp",
  "acceptProxyProtocol": false,
  "flow": "xtls-rprx-vision",
  "header": {
    "type": "none"
  },
  "socketSettings" : {
    "useSocket" : false,
    "DomainStrategy": "asis",
    "tcpKeepAliveInterval": 0,
    "tcpUserTimeout": 0,
    "tcpMaxSeg": 0,
    "tcpWindowClamp": 0,
    "tcpKeepAliveIdle": 0,
    "tcpMptcp": false
  }
}
```
#### TCP + HTTP
```
{
  "transport" : "tcp",
  "acceptProxyProtocol": false,
  "header": {
    "type": "http",
    "request": {
      "path": "/xmplus",
      "headers": {
        "Host": ["www.baidu.com", "www.taobao.com", "www.cloudflare.com"]
      }
    }
  },
  "socketSettings" : {
    "useSocket" : false,
    "DomainStrategy": "asis",
    "tcpKeepAliveInterval": 0,
    "tcpUserTimeout": 0,
    "tcpMaxSeg": 0,
    "tcpWindowClamp": 0,
    "tcpKeepAliveIdle": 0,
    "tcpMptcp": false
  }
}
```
####  WS
```
{
  "transport": "ws",
  "acceptProxyProtocol": false,
  "path": "/xmplus?ed=2560",
  "host": "hk1.xyz.com",
  "heartbeatperiod": 30,
  "custom_host": "fakedomain.com",
  "socketSettings" : {
    "useSocket" : false,
    "DomainStrategy": "asis",
    "tcpKeepAliveInterval": 0,
    "tcpUserTimeout": 0,
    "tcpMaxSeg": 0,
    "tcpWindowClamp": 0,
    "tcpKeepAliveIdle": 0,
    "tcpMptcp": false
  }
}
```

####  GRPC
```
{
  "transport" : "grpc",
  "acceptProxyProtocol": false,
  "serviceName": "xmplus",
  "authority": "hk1.xyz.com",
  "socketSettings" : {
    "useSocket" : false,
    "DomainStrategy": "asis",
    "tcpKeepAliveInterval": 0,
    "tcpUserTimeout": 0,
    "tcpMaxSeg": 0,
    "tcpWindowClamp": 0,
    "tcpKeepAliveIdle": 0,
    "tcpMptcp": false
  }
}
```


####  KCP
```
{
  "transport" : "kcp",
  "acceptProxyProtocol": false,
  "congestion": false,
  "header": {
    "type": "none"
  },
  "seed": "password",
  "socketSettings" : {
    "useSocket" : false,
    "DomainStrategy": "asis",
    "tcpKeepAliveInterval": 0,
    "tcpUserTimeout": 0,
    "tcpMaxSeg": 0,
    "tcpWindowClamp": 0,
    "tcpKeepAliveIdle": 0,
    "tcpMptcp": false
  }
}
```

####  HTTPUPGRADE
```
{
  "transport" : "httpupgrade",
  "acceptProxyProtocol": false,
  "host": "hk1.xyz.com",
  "path": "/xmplus?ed=2560",
  "custom_host": "fakedomain.com",
  "socketSettings" : {
    "useSocket" : false,
    "DomainStrategy": "asis",
    "tcpKeepAliveInterval": 0,
    "tcpUserTimeout": 0,
    "tcpMaxSeg": 0,
    "tcpWindowClamp": 0,
    "tcpKeepAliveIdle": 0,
    "tcpMptcp": false
  }
}
```

####  XHTTP / SPLITHTTP
```
{
  "transport" : "splithttp",
  "acceptProxyProtocol": false,
  "host": "hk1.xyz.com",
  "custom_host": "fakedomain.com",
  "path": "/",
  "noSSEHeader": false,
  "noGRPCHeader": true,
  "mode": "auto",
  "socketSettings" : {
    "useSocket" : false,
    "DomainStrategy": "asis",
    "tcpKeepAliveInterval": 0,
    "tcpUserTimeout": 0,
    "tcpMaxSeg": 0,
    "tcpWindowClamp": 0,
    "tcpKeepAliveIdle": 0,
    "tcpMptcp": false
  }
}
```

### Security Settings


#### TLS
```
{
  "serverName": "xmplus.dev",
  "rejectUnknownSni": false,
  "allowInsecure": false,
  "fingerprint": "chrome",
  "sni": "xmplus.dev",
  "curvepreferences": "X25519",
  "alpn": [
    "h2",
    "http/1.1"
  ],
  "serverNameToVerify" : ""
}
```
#### REALITY

`Generate Private and Public Keys :   xmplus x25519`

```
{
  "show" : false,
  "dest": "www.cloudflare.com:443",
  "privatekey" : "yBaw532IIUNuQWDTncozoBaLJmcd1JZzvsHUgVPxMk8",
  "minclientver":"",
  "maxclientver":"",
  "maxtimediff":0,
  "proxyprotocol":0,
  "shortids" : [
    "6ba85179e30d4fc2"
  ],
  "serverNames": [
    "www.cloudflare.com"
  ],
  "fingerprint": "chrome",
  "spiderx": "",
  "publickey": "7xhH4b_VkliBxGulljcyPOH-bYUA2dl-XAdZAsfhk04"
}
```

---

## API Feature Support (gRPC Management)

XMPlus now supports the xray-core API block for gRPC management, stats, and dynamic routing. To enable:

1. **Add `ApiConfigPath` to your `config.yml`:**
   ```yaml
   ApiConfigPath: ./files/api.json
   ```
2. **Create `api.json` (example):**
   ```json
   {
     "tag": "api",
     "services": ["HandlerService", "StatsService"]
   }
   ```
3. **Add an inbound for API (e.g., gRPC) in `inbound.json`:**
   ```json
   [
     {
       "tag": "api-in",
       "protocol": "grpc",
       "port": 10085,
       "listen": "127.0.0.1",
       "settings": {}
     }
   ]
   ```
4. **Add an outbound for API in `outbound.json`:**
   ```json
   [
     {
       "tag": "api",
       "protocol": "freedom",
       "settings": {}
     }
   ]
   ```
5. **Add a routing rule in `route.json`:**
   ```json
   {
     "rules": [
       {
         "inboundTag": ["api-in"],
         "outboundTag": "api",
         "type": "field"
       }
     ]
   }
   ```

**References:**
- [Xray-core API config documentation](https://xtls.github.io/config/api.html)
- [Xray-core Routing config documentation](https://xtls.github.io/config/routing.html)

This enables remote management and stats via gRPC (e.g., for use with XrayR/Xray API clients).

---

## API Feature Documentation

### What is the API Feature?
The API feature in XMPlus (compatible with xray-core) enables remote management, dynamic configuration, and statistics collection via gRPC. This is useful for advanced deployments, automation, monitoring, and integration with control panels or orchestration tools.

### When should I use it?
- If you want to manage inbounds/outbounds, routing rules, or query traffic stats at runtime.
- If you use a panel or automation tool that requires gRPC access to the core.
- If you want to script or remotely control your XMPlus instance.

### How does it work?
- The API block (see `api.json`) defines which gRPC services are enabled (e.g., HandlerService, StatsService).
- You must expose an inbound (usually gRPC) and outbound for the API, and route API traffic using a routing rule.
- Once enabled, you can use compatible clients or tools to connect to the API port and perform management actions.

### Example Use Cases
- Query real-time traffic stats for users or inbounds.
- Dynamically add/remove inbounds, outbounds, or routing rules without restarting XMPlus.
- Integrate with a web panel for live server management.

### Troubleshooting & FAQ
- **Q: My API port is not accessible?**
  - Ensure you have an inbound (e.g., gRPC) listening on the correct port and routed to the API outbound.
  - Check firewall and network settings.
- **Q: I get permission or service errors?**
  - Make sure the `services` array in `api.json` includes the required services (e.g., HandlerService, StatsService).
- **Q: How do I test the API?**
  - Use `grpcurl`, `xray api` commands, or compatible panel software to connect to the API port.
- **Q: Can I use HTTP instead of gRPC?**
  - No, the API feature is gRPC-based. Use a gRPC client.
- **Q: Where can I find more info?**
  - See the [Xray-core API documentation](https://xtls.github.io/config/api.html) and [Routing documentation](https://xtls.github.io/config/routing.html).

---
