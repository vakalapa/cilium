// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"

	ipsecfake "github.com/cilium/cilium/pkg/datapath/linux/ipsec/fake"
	"github.com/cilium/cilium/pkg/option"
	wgfake "github.com/cilium/cilium/pkg/wireguard/fake"
)

// nativeMode is a routing mode that disables IPsec's xfrm-output-mask probe
// (ProbeXfrmStateOutputMask is only called when tunneling is enabled), which
// keeps these tests purely in-memory.
const nativeMode = option.RoutingModeNative

// baseConfig returns an option.DaemonConfig that satisfies every IPsec-related
// validator other than the one under test, so individual tests can flip a
// single field at a time.
func baseConfig() *option.DaemonConfig {
	return &option.DaemonConfig{
		RoutingMode:                       nativeMode,
		EnableCiliumNodeCRD:               true,
		EnableL7Proxy:                     false,
		DNSProxyEnableTransparentMode:     false,
		EnableEncryptionStrictModeIngress: false,
	}
}

// TestValidateEncryptionDaemonConfig_HostFirewallIPsec verifies the Phase 0
// behavior: the cross-feature encryption validator must accept
// EnableHostFirewall=true with IPsec enabled in configurations the other
// IPsec validators allow, and must continue to reject IPsec combinations
// that other validators (WireGuard, strict ingress encryption,
// LocalRouterIPv*) gate.
func TestValidateEncryptionDaemonConfig_HostFirewallIPsec(t *testing.T) {
	tests := []struct {
		name        string
		mutate      func(*option.DaemonConfig)
		ipsec       ipsecfake.Config
		wireguard   wgfake.Config
		expectError string // substring match; "" means no error expected
	}{
		{
			name: "hostFW+IPsec is allowed (Phase 0 — was previously rejected)",
			mutate: func(c *option.DaemonConfig) {
				c.EnableHostFirewall = true
			},
			ipsec: ipsecfake.Config{EnableIPsec: true},
		},
		{
			name: "hostFW alone is allowed",
			mutate: func(c *option.DaemonConfig) {
				c.EnableHostFirewall = true
			},
		},
		{
			name:  "IPsec alone is allowed",
			ipsec: ipsecfake.Config{EnableIPsec: true},
		},
		{
			name: "neither is allowed",
		},
		{
			name: "hostFW+IPsec+endpoint-routes is allowed (gating left to IPsec validator suite)",
			mutate: func(c *option.DaemonConfig) {
				c.EnableHostFirewall = true
				c.EnableEndpointRoutes = true
			},
			ipsec: ipsecfake.Config{EnableIPsec: true},
		},

		// Cross-feature rejects must remain intact after Phase 0.
		{
			name:        "WireGuard+IPsec is still rejected",
			ipsec:       ipsecfake.Config{EnableIPsec: true},
			wireguard:   wgfake.Config{EnableWireguard: true},
			expectError: "WireGuard",
		},
		{
			// IPsec+strict-ingress is rejected at line 46 of daemon.go
			// before the strict-without-tunnel check at line 50 fires.
			// Keep RoutingMode=native so the earlier tunnel-probe netlink
			// path (line 40) is skipped in the test environment.
			name: "strict ingress encryption + IPsec is still rejected",
			mutate: func(c *option.DaemonConfig) {
				c.EnableEncryptionStrictModeIngress = true
			},
			ipsec:       ipsecfake.Config{EnableIPsec: true},
			expectError: "strict ingress encryption",
		},
		{
			name: "LocalRouterIPv4 + IPsec is still rejected",
			mutate: func(c *option.DaemonConfig) {
				c.LocalRouterIPv4 = "169.254.42.1"
			},
			ipsec:       ipsecfake.Config{EnableIPsec: true},
			expectError: "Cannot specify",
		},
		{
			name: "LocalRouterIPv6 + IPsec is still rejected",
			mutate: func(c *option.DaemonConfig) {
				c.LocalRouterIPv6 = "fe80::1"
			},
			ipsec:       ipsecfake.Config{EnableIPsec: true},
			expectError: "Cannot specify",
		},
		{
			name: "IPsec without CiliumNode CRD is still rejected",
			mutate: func(c *option.DaemonConfig) {
				c.EnableCiliumNodeCRD = false
			},
			ipsec:       ipsecfake.Config{EnableIPsec: true},
			expectError: "CiliumNode CRD",
		},
		{
			// Guards against the refactor silently dropping the DNS-proxy
			// transparent-mode requirement: with IPsec + L7 proxy enabled
			// and transparent mode disabled (and the insecure skip not
			// set), validation must reject.
			name: "IPsec + L7Proxy without DNS-proxy transparent mode is still rejected",
			mutate: func(c *option.DaemonConfig) {
				c.EnableL7Proxy = true
				c.DNSProxyEnableTransparentMode = false
			},
			ipsec:       ipsecfake.Config{EnableIPsec: true},
			expectError: "DNS proxy transparent mode",
		},
		{
			// And the inverse: transparent mode on must allow the combo.
			name: "IPsec + L7Proxy with DNS-proxy transparent mode is allowed",
			mutate: func(c *option.DaemonConfig) {
				c.EnableL7Proxy = true
				c.DNSProxyEnableTransparentMode = true
			},
			ipsec: ipsecfake.Config{EnableIPsec: true},
		},
		{
			// The IPsecConfig.DNSProxyInsecureSkipTransparentModeCheckEnabled
			// escape hatch must continue to bypass the check.
			name: "IPsec + L7Proxy with insecure-skip flag is allowed",
			mutate: func(c *option.DaemonConfig) {
				c.EnableL7Proxy = true
				c.DNSProxyEnableTransparentMode = false
			},
			ipsec: ipsecfake.Config{
				EnableIPsec: true,
				DNSProxyInsecureSkipTransparentModeCheck: true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := baseConfig()
			if tc.mutate != nil {
				tc.mutate(cfg)
			}

			err := validateEncryptionDaemonConfig(cfg, tc.ipsec, tc.wireguard)
			if tc.expectError == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectError)
			}
		})
	}
}

// TestValidateEncryptionDaemonConfig_HostFirewallIPsec_NoLegacyMessage guards
// against accidental reintroduction of the Phase 0 reject. If someone re-adds
// `IPSec cannot be used with the host firewall.`, this test fails fast
// instead of waiting for an end-to-end cluster failure to surface it.
func TestValidateEncryptionDaemonConfig_HostFirewallIPsec_NoLegacyMessage(t *testing.T) {
	cfg := baseConfig()
	cfg.EnableHostFirewall = true

	err := validateEncryptionDaemonConfig(cfg, ipsecfake.Config{EnableIPsec: true}, wgfake.Config{})
	require.NoError(t, err, "hostFW+IPsec must be accepted post Phase 0 (legacy reject must stay removed)")
}
