// SPDX-License-Identifier: (GPL-2.0-only OR BSD-2-Clause)
/* Copyright Authors of Cilium */

/* BPF tests for the egress hostFW skip on encrypted recirculation: verify
 * that the egress hostFW block in cil_to_netdev is bypassed for the
 * post-XFRM IPsec recirculation pass (MARK_MAGIC_ENCRYPT) while remaining
 * active for genuine host-originated traffic.
 *
 * Included by:
 *   host_hostfw_ipsec_egress.c        — native routing
 *   host_hostfw_ipsec_egress_tunnel.c — tunnel routing (TUNNEL_MODE)
 *
 * See bpf/bpf_host.c (cil_to_netdev, ENABLE_HOST_FIREWALL block) and
 * bpf/lib/host_firewall.h:343 (egress lookup short-circuit) for the
 * mechanism this test guards.
 */

#include <bpf/ctx/skb.h>
#include "common.h"
#include "pktgen.h"

#define ENABLE_IPV4			1
#define ENABLE_IPV6			1
#define ENABLE_HOST_FIREWALL		1

#define LOCAL_NODE_IP			v4_node_one
#define REMOTE_NODE_IP			v4_node_two
#define LOCAL_NODE_IP6			((const union v6addr *)v6_node_one)
#define REMOTE_NODE_IP6			((const union v6addr *)v6_node_two)

#define POD_SRC_IP			v4_pod_one
#define POD_DST_IP			v4_pod_one_on_node_two
#define POD_SRC_IDENTITY		(CIDR_IDENTITY_RANGE_START - 1)
#define POD_DST_IDENTITY		(CIDR_IDENTITY_RANGE_START - 2)

static volatile const __u8 *node_mac = mac_one;
static volatile const __u8 *peer_mac = mac_two;

#include "lib/bpf_host.h"

#include "lib/endpoint.h"
#include "lib/ipcache.h"
#include "lib/policy.h"

ASSIGN_CONFIG(bool, enable_conntrack_accounting, true)

/*
 * Test 1: hostfw_ipsec_egress_no_bogus_verdict (IPv4)
 *
 * Build an ESP packet shaped like the post-XFRM recirculation: outer src
 * is the local node IP, dst is a remote node IP, mark is MARK_MAGIC_ENCRYPT.
 * No host-policy ESP rule is installed. Without the encrypt-mark skip,
 * the hostFW egress block at bpf/bpf_host.c would evaluate
 *   { src = HOST_ID, dst = remote-node, proto = ESP, port = 0 }
 * and produce a default-deny verdict (or, pre-#44459, a DROP_INVALID
 * because ESP is not in the L4Proto enum). With the skip in place, the
 * `if (ctx_is_encrypt(ctx)) goto skip_host_firewall;` shortcut bypasses
 * that evaluation, so the packet leaves cil_to_netdev with CTX_ACT_OK
 * and no egress policy CT entry is created against the bogus tuple.
 */
PKTGEN("tc", "hostfw_ipsec_egress_v4_no_bogus_verdict")
int hostfw_ipsec_egress_v4_no_bogus_verdict_pktgen(struct __ctx_buff *ctx)
{
	struct pktgen builder;
	struct iphdr *l3;

	pktgen__init(&builder, ctx);

	l3 = pktgen__push_ipv4_packet(&builder, (__u8 *)node_mac, (__u8 *)peer_mac,
				      LOCAL_NODE_IP, REMOTE_NODE_IP);
	if (!l3)
		return TEST_ERROR;
	l3->protocol = IPPROTO_ESP;

	pktgen__finish(&builder);
	return 0;
}

SETUP("tc", "hostfw_ipsec_egress_v4_no_bogus_verdict")
int hostfw_ipsec_egress_v4_no_bogus_verdict_setup(struct __ctx_buff *ctx)
{
	endpoint_v4_add_entry(LOCAL_NODE_IP, 0, 0, ENDPOINT_F_HOST, HOST_ID,
			      0, (__u8 *)node_mac, (__u8 *)node_mac);
	ipcache_v4_add_entry(LOCAL_NODE_IP, 0, HOST_ID, 0, 0);
	ipcache_v4_add_entry(REMOTE_NODE_IP, 0, REMOTE_NODE_ID, 0, 0);
	ipcache_v4_add_world_entry();

	/* Default-deny posture on the host: would drop ESP from HOST_ID if
	 * the egress hostFW block actually ran. The encrypt-mark skip
	 * should bypass it.
	 */
	policy_add_egress_deny_all_entry();

	/* Stamp MARK_MAGIC_ENCRYPT to mimic the post-XFRM recirculated skb.
	 * The lower 16 bits select the "encrypt" branch via ctx_is_encrypt();
	 * upper bits (key/node_id) are not material for the hostFW skip but
	 * are set here for realism.
	 */
	ctx->mark = MARK_MAGIC_ENCRYPT;

	return netdev_send_packet(ctx);
}

CHECK("tc", "hostfw_ipsec_egress_v4_no_bogus_verdict")
int hostfw_ipsec_egress_v4_no_bogus_verdict_check(const struct __ctx_buff *ctx)
{
	void *data, *data_end;
	__u32 *status_code;

	test_init();

	endpoint_v4_del_entry(LOCAL_NODE_IP);
	policy_delete_egress_all_entry();

	data = (void *)(long)ctx_data(ctx);
	data_end = (void *)(long)ctx->data_end;

	if (data + sizeof(__u32) > data_end)
		test_fatal("status code out of bounds");

	status_code = data;
	assert(*status_code == CTX_ACT_OK);

	/* No egress CT entry for the spurious HOST_ID -> remote-node-IP/ESP
	 * tuple should exist: hostFW egress was skipped by the encrypt-mark
	 * guard.
	 */
	struct ipv4_ct_tuple tuple = {
		.daddr   = REMOTE_NODE_IP,
		.saddr   = LOCAL_NODE_IP,
		.dport   = 0,
		.sport   = 0,
		.nexthdr = IPPROTO_ESP,
		.flags   = TUPLE_F_OUT,
	};
	struct ct_entry *ct_entry = map_lookup_elem(get_ct_map4(&tuple), &tuple);

	if (ct_entry)
		test_fatal("hostFW egress unexpectedly created a CT entry on encrypted recirculation");

	test_finish();
}

/*
 * Test 2: hostfw_ipsec_egress_v4_host_originated_still_evaluated
 *
 * Scope-limit guard. A genuinely host-originated packet (mark =
 * MARK_MAGIC_HOST, not MARK_MAGIC_ENCRYPT) must still flow through the
 * hostFW egress evaluation. We use ICMP because it stays out of the
 * enable_extended_ip_protocols path and exercises the standard L4
 * lookup.
 */
PKTGEN("tc", "hostfw_ipsec_egress_v4_host_originated_drop")
int hostfw_ipsec_egress_v4_host_originated_drop_pktgen(struct __ctx_buff *ctx)
{
	struct pktgen builder;
	struct iphdr *l3;

	pktgen__init(&builder, ctx);

	l3 = pktgen__push_ipv4_packet(&builder, (__u8 *)node_mac, (__u8 *)peer_mac,
				      LOCAL_NODE_IP, REMOTE_NODE_IP);
	if (!l3)
		return TEST_ERROR;
	l3->protocol = IPPROTO_ICMP;

	pktgen__finish(&builder);
	return 0;
}

SETUP("tc", "hostfw_ipsec_egress_v4_host_originated_drop")
int hostfw_ipsec_egress_v4_host_originated_drop_setup(struct __ctx_buff *ctx)
{
	endpoint_v4_add_entry(LOCAL_NODE_IP, 0, 0, ENDPOINT_F_HOST, HOST_ID,
			      0, (__u8 *)node_mac, (__u8 *)node_mac);
	ipcache_v4_add_entry(LOCAL_NODE_IP, 0, HOST_ID, 0, 0);
	ipcache_v4_add_entry(REMOTE_NODE_IP, 0, REMOTE_NODE_ID, 0, 0);
	ipcache_v4_add_world_entry();

	policy_add_egress_deny_all_entry();

	/* MARK_MAGIC_HOST: host-originated, NOT encrypted recirculation.
	 * The encrypt-mark guard must not match — packet must still hit
	 * hostFW egress and be denied by the empty allow set.
	 */
	set_identity_mark(ctx, 0, MARK_MAGIC_HOST);

	return netdev_send_packet(ctx);
}

CHECK("tc", "hostfw_ipsec_egress_v4_host_originated_drop")
int hostfw_ipsec_egress_v4_host_originated_drop_check(const struct __ctx_buff *ctx)
{
	void *data, *data_end;
	__u32 *status_code;

	test_init();

	endpoint_v4_del_entry(LOCAL_NODE_IP);
	policy_delete_egress_all_entry();

	data = (void *)(long)ctx_data(ctx);
	data_end = (void *)(long)ctx->data_end;

	if (data + sizeof(__u32) > data_end)
		test_fatal("status code out of bounds");

	status_code = data;
	assert(*status_code == CTX_ACT_DROP);

	test_finish();
}

/*
 * Test 3: hostfw_ipsec_egress_v4_clear_pass_short_circuits
 *
 * Baseline guard documenting the §2.3 trace: a clear pod-to-pod packet
 * entering cil_to_netdev with a pod source identity hits the egress
 * lookup short-circuit at host_firewall.h:343 (because src is not
 * HOST_ID and ipcache_srcid is also not HOST_ID), so no hostFW verdict
 * is emitted on the clear inner pass. This establishes that the
 * encrypt-mark skip is suppressing a single bogus verdict on the
 * encrypted recirculation pass, not "skipping a duplicate" of an
 * earlier pass.
 */
PKTGEN("tc", "hostfw_ipsec_egress_v4_clear_pass")
int hostfw_ipsec_egress_v4_clear_pass_pktgen(struct __ctx_buff *ctx)
{
	struct pktgen builder;
	struct iphdr *l3;

	pktgen__init(&builder, ctx);

	l3 = pktgen__push_ipv4_packet(&builder, (__u8 *)node_mac, (__u8 *)peer_mac,
				      POD_SRC_IP, POD_DST_IP);
	if (!l3)
		return TEST_ERROR;
	l3->protocol = IPPROTO_TCP;

	pktgen__finish(&builder);
	return 0;
}

SETUP("tc", "hostfw_ipsec_egress_v4_clear_pass")
int hostfw_ipsec_egress_v4_clear_pass_setup(struct __ctx_buff *ctx)
{
	ipcache_v4_add_entry(POD_SRC_IP, 0, POD_SRC_IDENTITY, 0, 0);
	ipcache_v4_add_entry(POD_DST_IP, 0, POD_DST_IDENTITY, REMOTE_NODE_IP, 0);

	/* Default-deny on host policies. The clear pod-to-pod pass should
	 * NOT hit this — it must short-circuit at the egress lookup
	 * (host_firewall.h:343) because src is a pod identity, not HOST_ID.
	 */
	policy_add_egress_deny_all_entry();

	set_identity_mark(ctx, POD_SRC_IDENTITY, MARK_MAGIC_IDENTITY);

	return netdev_send_packet(ctx);
}

CHECK("tc", "hostfw_ipsec_egress_v4_clear_pass")
int hostfw_ipsec_egress_v4_clear_pass_check(const struct __ctx_buff *ctx)
{
	void *data, *data_end;
	__u32 *status_code;

	test_init();

	policy_delete_egress_all_entry();

	data = (void *)(long)ctx_data(ctx);
	data_end = (void *)(long)ctx->data_end;

	if (data + sizeof(__u32) > data_end)
		test_fatal("status code out of bounds");

	status_code = data;

	/* Clear pod-to-pod traffic: hostFW egress short-circuits in the
	 * lookup, so the deny-all policy does not fire and the packet
	 * exits with CTX_ACT_OK.
	 */
	assert(*status_code == CTX_ACT_OK);

	test_finish();
}

/*
 * Test 4: hostfw_ipsec_egress_v6_no_bogus_verdict
 *
 * IPv6 sibling of Test 1. The encrypt-mark guard uses ctx_is_encrypt()
 * which is protocol-agnostic, so the same skip must apply here.
 */
PKTGEN("tc", "hostfw_ipsec_egress_v6_no_bogus_verdict")
int hostfw_ipsec_egress_v6_no_bogus_verdict_pktgen(struct __ctx_buff *ctx)
{
	struct pktgen builder;
	struct ipv6hdr *l3;

	pktgen__init(&builder, ctx);

	l3 = pktgen__push_ipv6_packet(&builder, (__u8 *)node_mac, (__u8 *)peer_mac,
				      (__u8 *)LOCAL_NODE_IP6, (__u8 *)REMOTE_NODE_IP6);
	if (!l3)
		return TEST_ERROR;
	l3->nexthdr = IPPROTO_ESP;

	pktgen__finish(&builder);
	return 0;
}

SETUP("tc", "hostfw_ipsec_egress_v6_no_bogus_verdict")
int hostfw_ipsec_egress_v6_no_bogus_verdict_setup(struct __ctx_buff *ctx)
{
	ipcache_v6_add_entry(LOCAL_NODE_IP6, 0, HOST_ID, 0, 0);
	ipcache_v6_add_entry(REMOTE_NODE_IP6, 0, REMOTE_NODE_ID, 0, 0);

	policy_add_egress_deny_all_entry();

	ctx->mark = MARK_MAGIC_ENCRYPT;

	return netdev_send_packet(ctx);
}

CHECK("tc", "hostfw_ipsec_egress_v6_no_bogus_verdict")
int hostfw_ipsec_egress_v6_no_bogus_verdict_check(const struct __ctx_buff *ctx)
{
	void *data, *data_end;
	__u32 *status_code;

	test_init();

	policy_delete_egress_all_entry();

	data = (void *)(long)ctx_data(ctx);
	data_end = (void *)(long)ctx->data_end;

	if (data + sizeof(__u32) > data_end)
		test_fatal("status code out of bounds");

	status_code = data;
	assert(*status_code == CTX_ACT_OK);

	test_finish();
}
