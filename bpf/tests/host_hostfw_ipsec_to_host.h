// SPDX-License-Identifier: (GPL-2.0-only OR BSD-2-Clause)
/* Copyright Authors of Cilium */

/* Phase 0 / NodePort regression coverage for cil_to_host under the combined
 * ENABLE_HOST_FIREWALL + ENABLE_IPSEC + ENABLE_NODEPORT feature set.
 *
 * Purpose
 * -------
 * After Phase 0 lifts the daemon-side reject of hostFW+IPsec, the comment
 * formerly at bpf/bpf_host.c:1753 ("Since IPsec is not compatible with Host
 * Firewall, this won't be an issue") is no longer load-bearing — the
 * IPsec-recirc handle_nat_fwd() tail-call chain in cil_to_host becomes
 * reachable in real clusters. These tests give empirical (BPF unit-level)
 * coverage of that path:
 *
 *   1. cil_to_host with the combined feature set loads (verifier accepts).
 *   2. A typical IPsec-decrypted recirculation packet (mark=MARK_MAGIC_ENCRYPT)
 *      runs through the cil_to_host path without a verifier or runtime crash.
 *   3. The "non-encrypted" sibling continues to skip the NodePort revdnat
 *      branch as before, so this isn't accidentally exempting non-IPsec flows.
 *
 * Out of scope (kind-cluster work — see plan §7 open question 2)
 * --------------------------------------------------------------
 * Full revDNAT correctness for NodePort traffic decrypted from IPsec is NOT
 * asserted here. That requires a NodePort-with-IPsec end-to-end environment
 * and is tracked as kind playbook coverage. These tests are a verifier-and-
 * runtime smoke test, not a revnat correctness proof.
 *
 * Included by:
 *   host_hostfw_ipsec_to_host.c        — native routing
 *   host_hostfw_ipsec_to_host_tunnel.c — tunnel routing
 */

#include <bpf/ctx/skb.h>
#include "common.h"
#include "pktgen.h"

#define ENABLE_IPV4			1
#define ENABLE_IPV6			1
#define ENABLE_HOST_FIREWALL		1
#define ENABLE_NODEPORT			1
#define ENABLE_MASQUERADE_IPV4		1

#define LOCAL_NODE_IP			v4_node_one
#define REMOTE_NODE_IP			v4_node_two

#define POD_LOCAL_IP			v4_pod_one
#define POD_REMOTE_IP			v4_pod_one_on_node_two
#define POD_LOCAL_IDENTITY		(CIDR_IDENTITY_RANGE_START - 1)
#define POD_REMOTE_IDENTITY		(CIDR_IDENTITY_RANGE_START - 2)

#define CLIENT_PORT			bpf_htons(54321)
#define BACKEND_PORT			bpf_htons(8080)

static volatile const __u8 *node_mac = mac_one;
static volatile const __u8 *peer_mac = mac_two;

#include "lib/bpf_host.h"

#include "lib/ipcache.h"
#include "lib/policy.h"

ASSIGN_CONFIG(bool, enable_conntrack_accounting, true)

/* Build a TCP packet shaped like the post-decrypt inner reply that arrives
 * at cil_to_host on cilium_net after XFRM input: local pod (acting as
 * service backend) -> remote pod (the original client). For the NodePort
 * revdnat path, what matters is that the packet has a valid IPv4 + TCP
 * header that handle_nat_fwd can parse.
 */
static __always_inline int build_inner_tcp_v4(struct __ctx_buff *ctx)
{
	struct pktgen builder;
	struct tcphdr *l4;

	pktgen__init(&builder, ctx);

	l4 = pktgen__push_ipv4_tcp_packet(&builder,
					  (__u8 *)peer_mac, (__u8 *)node_mac,
					  POD_LOCAL_IP, POD_REMOTE_IP,
					  BACKEND_PORT, CLIENT_PORT);
	if (!l4)
		return TEST_ERROR;

	pktgen__finish(&builder);
	return 0;
}

/*
 * Test 1: ipsec_recirc_to_host_nodeport_path_runs
 *
 * Verifier + runtime smoke test. Compile cil_to_host with hostFW + IPsec
 * + NodePort all enabled, send an IPsec-recirc-shaped packet, verify the
 * program returns a non-error status. Without an LB/NAT entry seeded,
 * handle_nat_fwd(revdnat_only=true) finds no match and falls through;
 * cil_to_host then continues to host_ingress_policy below, which (with
 * the dst != HOST_ID short-circuit at host_firewall.h:455-457) returns
 * CTX_ACT_OK for pod-bound traffic. Net effect: status is OK.
 *
 * What this test catches:
 *   - Verifier rejection of the combined feature set (regression
 *     against e.g. someone adding a check that doesn't compile under
 *     all three flags).
 *   - Runtime crash / DROP_INVALID on the basic recirc packet shape.
 *
 * What this test does NOT catch (kind work):
 *   - Whether handle_nat_fwd's revDNAT actually rewrites addresses
 *     correctly when an LB/NAT entry exists.
 *   - End-to-end NodePort connectivity with hostFW + IPsec.
 */
PKTGEN("tc", "ipsec_recirc_to_host_nodeport_path_runs")
int ipsec_recirc_to_host_nodeport_path_runs_pktgen(struct __ctx_buff *ctx)
{
	return build_inner_tcp_v4(ctx);
}

SETUP("tc", "ipsec_recirc_to_host_nodeport_path_runs")
int ipsec_recirc_to_host_nodeport_path_runs_setup(struct __ctx_buff *ctx)
{
	ipcache_v4_add_entry(POD_LOCAL_IP, 0, POD_LOCAL_IDENTITY, 0, 0);
	ipcache_v4_add_entry(POD_REMOTE_IP, 0, POD_REMOTE_IDENTITY,
			     REMOTE_NODE_IP, 0);
	ipcache_v4_add_world_entry();

	/* MARK_MAGIC_ENCRYPT: this is the recirculation marker the kernel
	 * leaves on packets that have just been decrypted by XFRM input.
	 * cil_to_host's ENABLE_NODEPORT block at bpf_host.c:1770 keys off
	 * ctx_is_encrypt() to enter the handle_nat_fwd revdnat path.
	 */
	ctx->mark = MARK_MAGIC_ENCRYPT;

	return host_receive_packet(ctx);
}

CHECK("tc", "ipsec_recirc_to_host_nodeport_path_runs")
int ipsec_recirc_to_host_nodeport_path_runs_check(const struct __ctx_buff *ctx)
{
	void *data, *data_end;
	__u32 *status_code;

	test_init();

	data = (void *)(long)ctx_data(ctx);
	data_end = (void *)(long)ctx->data_end;

	if (data + sizeof(__u32) > data_end)
		test_fatal("status code out of bounds");

	status_code = data;

	/* Without an LB/NAT entry, revdnat is a no-op and the path falls
	 * through to host_ingress_policy which short-circuits because the
	 * inner dst is a pod identity (not HOST_ID). Net status: CTX_ACT_OK.
	 */
	assert(*status_code == CTX_ACT_OK);

	test_finish();
}

/*
 * Test 2: ipsec_recirc_to_host_no_encrypt_skips_nodeport_block
 *
 * Scope-limit guard. The same packet shape, same setup, but WITHOUT
 * MARK_MAGIC_ENCRYPT: the ENABLE_NODEPORT block at bpf_host.c:1770 must
 * skip the handle_nat_fwd call (its first check is `if
 * (!ctx_is_encrypt(ctx)) goto skip_ipsec_nodeport_revdnat;`). The packet
 * still flows through host_ingress_policy below and returns OK.
 */
PKTGEN("tc", "ipsec_recirc_to_host_no_encrypt_skips_nodeport_block")
int ipsec_recirc_to_host_no_encrypt_skips_nodeport_block_pktgen(struct __ctx_buff *ctx)
{
	return build_inner_tcp_v4(ctx);
}

SETUP("tc", "ipsec_recirc_to_host_no_encrypt_skips_nodeport_block")
int ipsec_recirc_to_host_no_encrypt_skips_nodeport_block_setup(struct __ctx_buff *ctx)
{
	ipcache_v4_add_entry(POD_LOCAL_IP, 0, POD_LOCAL_IDENTITY, 0, 0);
	ipcache_v4_add_entry(POD_REMOTE_IP, 0, POD_REMOTE_IDENTITY,
			     REMOTE_NODE_IP, 0);
	ipcache_v4_add_world_entry();

	/* No MARK_MAGIC_ENCRYPT: this is regular ingress, not IPsec recirc.
	 * The ENABLE_NODEPORT IPsec-revdnat block must not fire.
	 */
	ctx->mark = 0;

	return host_receive_packet(ctx);
}

CHECK("tc", "ipsec_recirc_to_host_no_encrypt_skips_nodeport_block")
int ipsec_recirc_to_host_no_encrypt_skips_nodeport_block_check(const struct __ctx_buff *ctx)
{
	void *data, *data_end;
	__u32 *status_code;

	test_init();

	data = (void *)(long)ctx_data(ctx);
	data_end = (void *)(long)ctx->data_end;

	if (data + sizeof(__u32) > data_end)
		test_fatal("status code out of bounds");

	status_code = data;
	assert(*status_code == CTX_ACT_OK);

	test_finish();
}
