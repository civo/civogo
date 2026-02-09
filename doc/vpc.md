# VPC API Documentation

VPC (Virtual Private Cloud) is the product umbrella containing networking-related features in Civo.

## Overview

The VPC methods provide a unified namespace for all VPC-related features:

- **Networks** - Private networks for instances to connect to
- **Firewalls** - Firewall rules for network security
- **DNS** - DNS domains and records management
- **Load Balancers** - Load balancer configuration
- **Reserved IPs** - Reserved IP address management

These methods call the `/v2/vpc/*` API endpoints. The original methods (e.g., `ListNetworks`, `ListFirewalls`) remain unchanged and continue to work. SDK users can migrate to VPC methods at their own pace.

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/civo/civogo"
)

func main() {
    client, _ := civogo.NewClient("your-api-key", "LON1")

    // List all VPC networks
    networks, _ := client.ListVPCNetworks()
    for _, network := range networks {
        fmt.Printf("Network: %s (%s)\n", network.Label, network.ID)
    }
}
```

---

## VPC Networks

Private networks for instances to connect to.

### ListVPCNetworks

Lists all private networks.

```go
networks, err := client.ListVPCNetworks()
if err != nil {
    log.Fatal(err)
}

for _, network := range networks {
    fmt.Printf("ID: %s, Name: %s, Default: %v\n", network.ID, network.Label, network.Default)
}
```

### GetVPCNetwork

Gets a specific network by ID.

```go
network, err := client.GetVPCNetwork("network-id")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Network: %s, CIDR: %s, Status: %s\n", network.Label, network.CIDR, network.Status)
```

### GetDefaultVPCNetwork

Finds the default private network for an account.

```go
defaultNetwork, err := client.GetDefaultVPCNetwork()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Default Network: %s\n", defaultNetwork.Label)
```

### FindVPCNetwork

Finds a network by either part of the ID or part of the name/label.

```go
network, err := client.FindVPCNetwork("production")
if err != nil {
    if errors.Is(err, civogo.ZeroMatchesError) {
        fmt.Println("No network found")
    }
    if errors.Is(err, civogo.MultipleMatchesError) {
        fmt.Println("Multiple networks matched, be more specific")
    }
    log.Fatal(err)
}

fmt.Printf("Found: %s\n", network.Label)
```

### NewVPCNetwork

Creates a new private network with just a label.

```go
result, err := client.NewVPCNetwork("my-network")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created network: %s (ID: %s)\n", result.Label, result.ID)
```

### CreateVPCNetwork

Creates a new network with full configuration options.

```go
ipv4Enabled := true
config := civogo.NetworkConfig{
    Label:       "my-custom-network",
    IPv4Enabled: &ipv4Enabled,
    CIDRv4:      "10.0.0.0/24",
}

result, err := client.CreateVPCNetwork(config)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created network: %s\n", result.ID)
```

### UpdateVPCNetwork

Updates an existing network.

```go
config := civogo.NetworkConfig{
    Label: "updated-network-name",
}

result, err := client.UpdateVPCNetwork("network-id", config)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Updated network: %s\n", result.Label)
```

### RenameVPCNetwork

Renames an existing private network.

```go
result, err := client.RenameVPCNetwork("new-name", "network-id")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Renamed to: %s\n", result.Label)
```

### DeleteVPCNetwork

Deletes a private network.

```go
resp, err := client.DeleteVPCNetwork("network-id")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Result: %s\n", resp.Result)
```

---

## VPC Firewalls

Firewall management for network security.

### ListVPCFirewalls

Returns all firewalls owned by the calling API account.

```go
firewalls, err := client.ListVPCFirewalls()
if err != nil {
    log.Fatal(err)
}

for _, fw := range firewalls {
    fmt.Printf("Firewall: %s, Rules: %d\n", fw.Name, fw.RulesCount)
}
```

### FindVPCFirewall

Finds a firewall by either part of the ID or part of the name.

```go
firewall, err := client.FindVPCFirewall("web-firewall")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found: %s (ID: %s)\n", firewall.Name, firewall.ID)
```

### NewVPCFirewall

Creates a new firewall.

```go
createRules := true
config := &civogo.FirewallConfig{
    Name:        "my-firewall",
    NetworkID:   "network-id",
    Region:      "LON1",
    CreateRules: &createRules, // Creates default rules
}

result, err := client.NewVPCFirewall(config)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created firewall: %s\n", result.ID)
```

### RenameVPCFirewall

Renames a firewall.

```go
config := &civogo.FirewallConfig{
    Name: "new-firewall-name",
}

resp, err := client.RenameVPCFirewall("firewall-id", config)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Result: %s\n", resp.Result)
```

### DeleteVPCFirewall

Deletes a firewall.

```go
resp, err := client.DeleteVPCFirewall("firewall-id")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Result: %s\n", resp.Result)
```

---

## VPC Firewall Rules

Individual rules within firewalls.

### ListVPCFirewallRules

Gets all rules for a firewall.

```go
rules, err := client.ListVPCFirewallRules("firewall-id")
if err != nil {
    log.Fatal(err)
}

for _, rule := range rules {
    fmt.Printf("Rule: %s, Protocol: %s, Ports: %s\n", rule.ID, rule.Protocol, rule.Ports)
}
```

### FindVPCFirewallRule

Finds a firewall rule by ID or part of the ID.

```go
rule, err := client.FindVPCFirewallRule("firewall-id", "rule-id")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found rule: %s\n", rule.ID)
```

### NewVPCFirewallRule

Creates a new rule within a firewall.

```go
config := &civogo.FirewallRuleConfig{
    FirewallID: "firewall-id",
    Protocol:   "tcp",
    StartPort:  "443",
    EndPort:    "443",
    Cidr:       []string{"0.0.0.0/0"},
    Direction:  "ingress",
    Action:     "allow",
    Label:      "HTTPS",
}

rule, err := client.NewVPCFirewallRule(config)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created rule: %s\n", rule.ID)
```

### DeleteVPCFirewallRule

Deletes a firewall rule.

```go
resp, err := client.DeleteVPCFirewallRule("firewall-id", "rule-id")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Result: %s\n", resp.Result)
```

---

## VPC DNS Domains

DNS domain management.

### ListVPCDNSDomains

Returns all DNS domains owned by the calling API account.

```go
domains, err := client.ListVPCDNSDomains()
if err != nil {
    log.Fatal(err)
}

for _, domain := range domains {
    fmt.Printf("Domain: %s (ID: %s)\n", domain.Name, domain.ID)
}
```

### FindVPCDNSDomain

Finds a domain by either part of the ID or part of the name.

```go
domain, err := client.FindVPCDNSDomain("example.com")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found: %s\n", domain.Name)
```

### GetVPCDNSDomain

Returns the DNS domain that matches the name.

```go
domain, err := client.GetVPCDNSDomain("example.com")
if err != nil {
    if errors.Is(err, civogo.ErrDNSDomainNotFound) {
        fmt.Println("Domain not found")
    }
    log.Fatal(err)
}

fmt.Printf("Domain ID: %s\n", domain.ID)
```

### CreateVPCDNSDomain

Registers a new DNS domain.

```go
domain, err := client.CreateVPCDNSDomain("example.com")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created domain: %s (ID: %s)\n", domain.Name, domain.ID)
```

### UpdateVPCDNSDomain

Updates a DNS domain.

```go
domain := &civogo.DNSDomain{ID: "domain-id", Name: "old-name.com"}

updated, err := client.UpdateVPCDNSDomain(domain, "new-name.com")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Updated to: %s\n", updated.Name)
```

### DeleteVPCDNSDomain

Deletes a DNS domain.

```go
domain := &civogo.DNSDomain{ID: "domain-id"}

resp, err := client.DeleteVPCDNSDomain(domain)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Result: %s\n", resp.Result)
```

---

## VPC DNS Records

DNS record management within domains.

### ListVPCDNSRecords

Returns all records associated with a domain.

```go
records, err := client.ListVPCDNSRecords("domain-id")
if err != nil {
    log.Fatal(err)
}

for _, record := range records {
    fmt.Printf("Record: %s -> %s (Type: %s)\n", record.Name, record.Value, record.Type)
}
```

### GetVPCDNSRecord

Returns a specific DNS record.

```go
record, err := client.GetVPCDNSRecord("domain-id", "record-id")
if err != nil {
    if errors.Is(err, civogo.ErrDNSRecordNotFound) {
        fmt.Println("Record not found")
    }
    log.Fatal(err)
}

fmt.Printf("Record: %s\n", record.Name)
```

### CreateVPCDNSRecord

Creates a new DNS record.

```go
config := &civogo.DNSRecordConfig{
    Type:     civogo.DNSRecordTypeA,
    Name:     "www",
    Value:    "1.2.3.4",
    TTL:      600,
}

record, err := client.CreateVPCDNSRecord("domain-id", config)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created record: %s\n", record.ID)
```

### UpdateVPCDNSRecord

Updates a DNS record.

```go
record := &civogo.DNSRecord{
    ID:          "record-id",
    DNSDomainID: "domain-id",
}

config := &civogo.DNSRecordConfig{
    Type:  civogo.DNSRecordTypeA,
    Name:  "www",
    Value: "5.6.7.8",
}

updated, err := client.UpdateVPCDNSRecord(record, config)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Updated value: %s\n", updated.Value)
```

### DeleteVPCDNSRecord

Deletes a DNS record.

```go
record := &civogo.DNSRecord{
    ID:          "record-id",
    DNSDomainID: "domain-id",
}

resp, err := client.DeleteVPCDNSRecord(record)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Result: %s\n", resp.Result)
```

---

## VPC Load Balancers

Load balancer management.

### ListVPCLoadBalancers

Returns all load balancers owned by the calling API account.

```go
loadbalancers, err := client.ListVPCLoadBalancers()
if err != nil {
    log.Fatal(err)
}

for _, lb := range loadbalancers {
    fmt.Printf("LB: %s, Public IP: %s, State: %s\n", lb.Name, lb.PublicIP, lb.State)
}
```

### GetVPCLoadBalancer

Returns a specific load balancer.

```go
lb, err := client.GetVPCLoadBalancer("loadbalancer-id")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Load Balancer: %s, Algorithm: %s\n", lb.Name, lb.Algorithm)
```

### FindVPCLoadBalancer

Finds a load balancer by either part of the ID or part of the name.

```go
lb, err := client.FindVPCLoadBalancer("web-lb")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found: %s\n", lb.Name)
```

### CreateVPCLoadBalancer

Creates a new load balancer.

```go
config := &civogo.LoadBalancerConfig{
    Name:      "my-loadbalancer",
    Region:    "LON1",
    NetworkID: "network-id",
    Algorithm: "round_robin",
    Backends: []civogo.LoadBalancerBackendConfig{
        {
            IP:         "10.0.0.1",
            Protocol:   "http",
            SourcePort: 80,
            TargetPort: 8080,
        },
    },
}

lb, err := client.CreateVPCLoadBalancer(config)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created LB: %s (Public IP: %s)\n", lb.ID, lb.PublicIP)
```

### UpdateVPCLoadBalancer

Updates a load balancer.

```go
config := &civogo.LoadBalancerUpdateConfig{
    Name:      "updated-lb-name",
    Region:    "LON1",
    Algorithm: "least_connections",
}

lb, err := client.UpdateVPCLoadBalancer("loadbalancer-id", config)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Updated: %s\n", lb.Name)
```

### DeleteVPCLoadBalancer

Deletes a load balancer.

```go
resp, err := client.DeleteVPCLoadBalancer("loadbalancer-id")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Result: %s\n", resp.Result)
```

---

## VPC Reserved IPs

Reserved IP address management.

### ListVPCIPs

Returns all reserved IPs in the region.

```go
ips, err := client.ListVPCIPs()
if err != nil {
    log.Fatal(err)
}

for _, ip := range ips.Items {
    fmt.Printf("IP: %s, Name: %s\n", ip.IP, ip.Name)
}
```

### GetVPCIP

Finds a reserved IP by ID.

```go
ip, err := client.GetVPCIP("ip-id")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("IP: %s, Name: %s\n", ip.IP, ip.Name)
```

### FindVPCIP

Finds a reserved IP by name, IP address, or ID.

```go
ip, err := client.FindVPCIP("1.2.3.4")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found: %s (ID: %s)\n", ip.IP, ip.ID)
```

### NewVPCIP

Creates a new reserved IP.

```go
config := &civogo.CreateIPRequest{
    Name:   "my-reserved-ip",
    Region: "LON1",
}

ip, err := client.NewVPCIP(config)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created IP: %s (Address: %s)\n", ip.ID, ip.IP)
```

### UpdateVPCIP

Updates a reserved IP.

```go
config := &civogo.UpdateIPRequest{
    Name:   "updated-ip-name",
    Region: "LON1",
}

ip, err := client.UpdateVPCIP("ip-id", config)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Updated: %s\n", ip.Name)
```

### AssignVPCIP

Assigns a reserved IP to a Civo resource.

```go
resp, err := client.AssignVPCIP("ip-id", "instance-id", "instance", "LON1")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Result: %s\n", resp.Result)
```

### UnassignVPCIP

Unassigns a reserved IP from a Civo resource.

```go
resp, err := client.UnassignVPCIP("ip-id", "LON1")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Result: %s\n", resp.Result)
```

### DeleteVPCIP

Deletes a reserved IP.

```go
resp, err := client.DeleteVPCIP("ip-id")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Result: %s\n", resp.Result)
```

---

## Migration Guide

To migrate from original methods to VPC methods, simply replace the method names:

| Original Method | VPC Method |
|-----------------|------------|
| `ListNetworks()` | `ListVPCNetworks()` |
| `GetNetwork()` | `GetVPCNetwork()` |
| `FindNetwork()` | `FindVPCNetwork()` |
| `NewNetwork()` | `NewVPCNetwork()` |
| `CreateNetwork()` | `CreateVPCNetwork()` |
| `UpdateNetwork()` | `UpdateVPCNetwork()` |
| `RenameNetwork()` | `RenameVPCNetwork()` |
| `DeleteNetwork()` | `DeleteVPCNetwork()` |
| `ListFirewalls()` | `ListVPCFirewalls()` |
| `FindFirewall()` | `FindVPCFirewall()` |
| `NewFirewall()` | `NewVPCFirewall()` |
| `RenameFirewall()` | `RenameVPCFirewall()` |
| `DeleteFirewall()` | `DeleteVPCFirewall()` |
| `ListDNSDomains()` | `ListVPCDNSDomains()` |
| `CreateDNSDomain()` | `CreateVPCDNSDomain()` |
| `ListLoadBalancers()` | `ListVPCLoadBalancers()` |
| `CreateLoadBalancer()` | `CreateVPCLoadBalancer()` |
| `ListIPs()` | `ListVPCIPs()` |
| `NewIP()` | `NewVPCIP()` |

Both old and new methods work - no breaking changes.

---

## Error Handling

VPC methods use the same error types as the original methods:

```go
network, err := client.FindVPCNetwork("missing-network")
if err != nil {
    if errors.Is(err, civogo.ZeroMatchesError) {
        fmt.Println("No network found matching the search term")
    } else if errors.Is(err, civogo.MultipleMatchesError) {
        fmt.Println("Multiple networks found, please be more specific")
    } else {
        fmt.Printf("Error: %v\n", err)
    }
}
```

Common errors:
- `ZeroMatchesError` - No resources found matching the search
- `MultipleMatchesError` - Multiple resources matched, need more specific search
- `IDisEmptyError` - Required ID parameter was empty
- `ErrDNSDomainNotFound` - DNS domain not found
- `ErrDNSRecordNotFound` - DNS record not found
