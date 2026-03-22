// SPDX-License-Identifier: (GPL-2.0-only OR BSD-2-Clause)
/* Copyright Authors of Cilium */

/* Test that host firewall ingress policy is enforced on packets arriving
 * from the network when IPsec is enabled. This validates that decrypted
 * IPsec packets go through do_netdev() where host firewall runs.
 */

#include <bpf/ctx/skb.h>
#include "common.h"
#include "pktgen.h"

/* Enable code paths under test */
#define ENABLE_IPV4			1
#define ENABLE_HOST_FIREWALL		1
#define ENABLE_IPSEC			1

#define REMOTE_POD_SEC_IDENTITY		112233

#define NODE_IP				v4_node_one
#define NODE_PORT			bpf_htons(8080)

#define REMOTE_POD_IP			v4_pod_one_on_node_two
#define REMOTE_POD_PORT			bpf_htons(50000)

#define REMOTE_NODE_IP			v4_node_two

#define ENCRYPT_KEY			0xFF

static volatile const __u8 *node_mac = mac_one;
static volatile const __u8 *remote_mac = mac_two;

#include "lib/bpf_host.h"

ASSIGN_CONFIG(bool, enable_conntrack_accounting, true)

#include "lib/endpoint.h"
#include "lib/ipcache.h"
#include "lib/ipsec.h"
#include "lib/node.h"
#include "lib/policy.h"

/* Test 1: A packet from a remote pod (simulating already-decrypted IPsec
 * traffic) arriving at from-netdev and destined to the host node IP.
 * Host firewall ingress policy should be enforced.
 *
 * With a default-deny policy and an allow rule for the remote pod identity,
 * the packet should be allowed through.
 */
PKTGEN("tc", "hostfw_ipsec_ingress_01_allow")
int hostfw_ipsec_ingress_01_allow_pktgen(struct __ctx_buff *ctx)
{
	struct pktgen builder;
	struct tcphdr *tcp;

	pktgen__init(&builder, ctx);

	tcp = pktgen__push_ipv4_tcp_packet(&builder,
					   (__u8 *)remote_mac, (__u8 *)node_mac,
					   REMOTE_POD_IP, NODE_IP,
					   REMOTE_POD_PORT, NODE_PORT);
	if (!tcp)
		return TEST_ERROR;

	pktgen__finish(&builder);

	return 0;
}

SETUP("tc", "hostfw_ipsec_ingress_01_allow")
int hostfw_ipsec_ingress_01_allow_setup(struct __ctx_buff *ctx)
{
	/* Set up host endpoint in the endpoint map */
	endpoint_v4_add_entry(NODE_IP, 0, 0, ENDPOINT_F_HOST, HOST_ID,
			      0, (__u8 *)node_mac, (__u8 *)node_mac);

	/* Add ipcache entries */
	ipcache_v4_add_entry(NODE_IP, 0, HOST_ID, 0, 0);
	ipcache_v4_add_entry(REMOTE_POD_IP, 0, REMOTE_POD_SEC_IDENTITY,
			     REMOTE_NODE_IP, ENCRYPT_KEY);
	ipcache_v4_add_world_entry();

	/* Set up IPsec state */
	ipsec_set_encrypt_state(ENCRYPT_KEY);
	node_v4_add_entry(REMOTE_NODE_IP, 123, ENCRYPT_KEY);

	/* Add host firewall ingress policy: allow from REMOTE_POD_SEC_IDENTITY */
	policy_add_ingress_allow_l3_l4_entry(REMOTE_POD_SEC_IDENTITY, 0, 0, 0);

	/* Simulate packet arriving from network (no special marks - this is
	 * what a decrypted IPsec packet looks like after do_decrypt() clears
	 * the mark on the second pass through cil_from_netdev).
	 */
	return netdev_receive_packet(ctx);
}

CHECK("tc", "hostfw_ipsec_ingress_01_allow")
int hostfw_ipsec_ingress_01_allow_check(const struct __ctx_buff *ctx)
{
	void *data, *data_end;
	__u32 *status_code;

	test_init();

	data = (void *)(long)ctx_data(ctx);
	data_end = (void *)(long)ctx->data_end;

	if (data + sizeof(__u32) > data_end)
		test_fatal("status code out of bounds");

	status_code = data;

	/* Packet should be allowed through (CTX_ACT_OK means passed to stack) */
	assert(*status_code == CTX_ACT_OK);

	test_finish();
}

/* Test 2: Same scenario but with a deny-all ingress policy.
 * The packet from the remote pod should be dropped.
 */
PKTGEN("tc", "hostfw_ipsec_ingress_02_deny")
int hostfw_ipsec_ingress_02_deny_pktgen(struct __ctx_buff *ctx)
{
	struct pktgen builder;
	struct tcphdr *tcp;

	pktgen__init(&builder, ctx);

	tcp = pktgen__push_ipv4_tcp_packet(&builder,
					   (__u8 *)remote_mac, (__u8 *)node_mac,
					   REMOTE_POD_IP, NODE_IP,
					   REMOTE_POD_PORT, NODE_PORT);
	if (!tcp)
		return TEST_ERROR;

	pktgen__finish(&builder);

	return 0;
}

SETUP("tc", "hostfw_ipsec_ingress_02_deny")
int hostfw_ipsec_ingress_02_deny_setup(struct __ctx_buff *ctx)
{
	/* Set up host endpoint */
	endpoint_v4_add_entry(NODE_IP, 0, 0, ENDPOINT_F_HOST, HOST_ID,
			      0, (__u8 *)node_mac, (__u8 *)node_mac);

	/* Add ipcache entries */
	ipcache_v4_add_entry(NODE_IP, 0, HOST_ID, 0, 0);
	ipcache_v4_add_entry(REMOTE_POD_IP, 0, REMOTE_POD_SEC_IDENTITY,
			     REMOTE_NODE_IP, ENCRYPT_KEY);
	ipcache_v4_add_world_entry();

	/* Set up IPsec state */
	ipsec_set_encrypt_state(ENCRYPT_KEY);
	node_v4_add_entry(REMOTE_NODE_IP, 123, ENCRYPT_KEY);

	/* Add host firewall deny-all ingress policy */
	policy_add_ingress_deny_all_entry();

	return netdev_receive_packet(ctx);
}

CHECK("tc", "hostfw_ipsec_ingress_02_deny")
int hostfw_ipsec_ingress_02_deny_check(const struct __ctx_buff *ctx)
{
	void *data, *data_end;
	__u32 *status_code;

	test_init();

	data = (void *)(long)ctx_data(ctx);
	data_end = (void *)(long)ctx->data_end;

	if (data + sizeof(__u32) > data_end)
		test_fatal("status code out of bounds");

	status_code = data;

	/* Packet should be dropped by host firewall policy */
	assert(*status_code == CTX_ACT_DROP);

	test_finish();
}
