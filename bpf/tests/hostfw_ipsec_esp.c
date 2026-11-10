// SPDX-License-Identifier: (GPL-2.0-only OR BSD-2-Clause)
/* Copyright Authors of Cilium */
/* This file was not reviewed by a human before submission */

#define ENABLE_IPV4 1
#define ENABLE_IPV6

#include <bpf/ctx/skb.h>
#include "common.h"
#include "pktgen.h"

#include "lib/bpf_host.h"

#include "lib/ipcache.h"
#include "lib/node.h"
#include "lib/policy.h"
#include "scapy.h"

ASSIGN_CONFIG(bool, enable_ipv4, true)
ASSIGN_CONFIG(bool, enable_ipsec, true)
ASSIGN_CONFIG(bool, enable_host_firewall, true)
ASSIGN_CONFIG(bool, enable_extended_ip_protocols, true)

/* packet defined in ./scapy/ipsec_from_netdev_pkt_defs.py */
const __u8 v4_ipsec[] = {
	SCAPY_BUF_BYTES(v4_ipsec)
};

/* packet defined in ./scapy/ipsec_from_netdev_pkt_defs.py */
const __u8 v6_ipsec[] = {
	SCAPY_BUF_BYTES(v6_ipsec)
};

static __always_inline int pktgen(struct __ctx_buff *ctx, bool ipv4)
{
	struct pktgen builder;

	pktgen__init(&builder, ctx);

	if (ipv4)
		scapy_push_data(&builder, v4_ipsec, sizeof(v4_ipsec));
	else
		scapy_push_data(&builder, v6_ipsec, sizeof(v6_ipsec));

	pktgen__finish(&builder);
	return 0;
}

/* The ESP packet is sent from node_one to node_two. */
static __always_inline void
setup_nodes(bool ipv4, __be32 v4_host, __be32 v4_remote,
	    const __u8 *v6_host, const __u8 *v6_remote)
{
	if (ipv4) {
		ipcache_v4_add_entry(v4_host, 0, HOST_ID, 0, 0);
		ipcache_v4_add_entry(v4_remote, 0, REMOTE_NODE_ID, 0, 0);
		node_v4_add_entry(v4_remote, REMOTE_NODE_ID, 3);
	} else {
		ipcache_v6_add_entry((union v6addr *)v6_host, 0, HOST_ID, 0, 0);
		ipcache_v6_add_entry((union v6addr *)v6_remote, 0, REMOTE_NODE_ID, 0, 0);
		node_v6_add_entry((union v6addr *)v6_remote, REMOTE_NODE_ID, 3);
	}
}

static __always_inline int
ingress_setup(struct __ctx_buff *ctx, bool ipv4, bool allow)
{
	setup_nodes(ipv4, v4_node_two, v4_node_one, v6_node_two, v6_node_one);

	policy_delete_ingress_all_entry();
	if (allow)
		policy_add_ingress_allow_all_entry();

	return netdev_receive_packet(ctx);
}

static __always_inline int
ingress_check(const struct __ctx_buff *ctx, bool ipv4, bool allow)
{
	void *data, *data_end;
	__u32 *status_code;

	test_init();

	data = (void *)(long)ctx->data;
	data_end = (void *)(long)ctx->data_end;

	if (data + sizeof(*status_code) > data_end)
		test_fatal("status code out of bounds");

	status_code = data;

	if (!allow) {
		assert(*status_code == CTX_ACT_DROP);
		test_finish();
	}

	assert(*status_code == CTX_ACT_OK);
	assert(ctx_is_decrypt(ctx));

	if (ipv4) {
		ASSERT_CTX_BUF_OFF("v4_ipsec_pkt_ok", "Ether", ctx, sizeof(__u32),
				   v4_ipsec, sizeof(v4_ipsec));

		struct ipv4_ct_tuple tuple = {
			.daddr   = v4_node_one,
			.saddr   = v4_node_two,
			.nexthdr = IPPROTO_ESP,
			.flags   = TUPLE_F_IN,
		};

		if (!map_lookup_elem(get_ct_map4(&tuple), &tuple))
			test_fatal("no ingress CT entry for ESP");
	} else {
		ASSERT_CTX_BUF_OFF("v6_ipsec_pkt_ok", "Ether", ctx, sizeof(__u32),
				   v6_ipsec, sizeof(v6_ipsec));
	}

	policy_delete_ingress_all_entry();

	test_finish();
}

static __always_inline int
egress_setup(struct __ctx_buff *ctx, bool ipv4, bool allow)
{
	setup_nodes(ipv4, v4_node_one, v4_node_two, v6_node_one, v6_node_two);

	policy_delete_egress_all_entry();
	if (allow)
		policy_add_egress_allow_all_entry();

	set_identity_mark(ctx, 0, MARK_MAGIC_ENCRYPT);

	return netdev_send_packet(ctx);
}

static __always_inline int
egress_check(const struct __ctx_buff *ctx, bool ipv4, bool allow)
{
	void *data, *data_end;
	__u32 *status_code;

	test_init();

	data = (void *)(long)ctx->data;
	data_end = (void *)(long)ctx->data_end;

	if (data + sizeof(*status_code) > data_end)
		test_fatal("status code out of bounds");

	status_code = data;

	if (!allow) {
		assert(*status_code == CTX_ACT_DROP);
		test_finish();
	}

	assert(*status_code == CTX_ACT_OK);

	if (ipv4) {
		struct ipv4_ct_tuple tuple = {
			.daddr   = v4_node_one,
			.saddr   = v4_node_two,
			.nexthdr = IPPROTO_ESP,
			.flags   = TUPLE_F_OUT,
		};

		if (!map_lookup_elem(get_ct_map4(&tuple), &tuple))
			test_fatal("no egress CT entry for ESP");
	}

	policy_delete_egress_all_entry();

	test_finish();
}

PKTGEN(PROG_TYPE, "hostfw_ipsec_esp_1_ingress_v4_deny")
int ingress_v4_deny_pktgen(struct __ctx_buff *ctx)
{
	return pktgen(ctx, true);
}

SETUP(PROG_TYPE, "hostfw_ipsec_esp_1_ingress_v4_deny")
int ingress_v4_deny_setup(struct __ctx_buff *ctx)
{
	return ingress_setup(ctx, true, false);
}

CHECK(PROG_TYPE, "hostfw_ipsec_esp_1_ingress_v4_deny")
int ingress_v4_deny_check(const struct __ctx_buff *ctx)
{
	return ingress_check(ctx, true, false);
}

PKTGEN(PROG_TYPE, "hostfw_ipsec_esp_2_ingress_v4_allow")
int ingress_v4_allow_pktgen(struct __ctx_buff *ctx)
{
	return pktgen(ctx, true);
}

SETUP(PROG_TYPE, "hostfw_ipsec_esp_2_ingress_v4_allow")
int ingress_v4_allow_setup(struct __ctx_buff *ctx)
{
	return ingress_setup(ctx, true, true);
}

CHECK(PROG_TYPE, "hostfw_ipsec_esp_2_ingress_v4_allow")
int ingress_v4_allow_check(const struct __ctx_buff *ctx)
{
	return ingress_check(ctx, true, true);
}

PKTGEN(PROG_TYPE, "hostfw_ipsec_esp_3_ingress_v6_deny")
int ingress_v6_deny_pktgen(struct __ctx_buff *ctx)
{
	return pktgen(ctx, false);
}

SETUP(PROG_TYPE, "hostfw_ipsec_esp_3_ingress_v6_deny")
int ingress_v6_deny_setup(struct __ctx_buff *ctx)
{
	return ingress_setup(ctx, false, false);
}

CHECK(PROG_TYPE, "hostfw_ipsec_esp_3_ingress_v6_deny")
int ingress_v6_deny_check(const struct __ctx_buff *ctx)
{
	return ingress_check(ctx, false, false);
}

PKTGEN(PROG_TYPE, "hostfw_ipsec_esp_4_ingress_v6_allow")
int ingress_v6_allow_pktgen(struct __ctx_buff *ctx)
{
	return pktgen(ctx, false);
}

SETUP(PROG_TYPE, "hostfw_ipsec_esp_4_ingress_v6_allow")
int ingress_v6_allow_setup(struct __ctx_buff *ctx)
{
	return ingress_setup(ctx, false, true);
}

CHECK(PROG_TYPE, "hostfw_ipsec_esp_4_ingress_v6_allow")
int ingress_v6_allow_check(const struct __ctx_buff *ctx)
{
	return ingress_check(ctx, false, true);
}

PKTGEN(PROG_TYPE, "hostfw_ipsec_esp_5_egress_v4_deny")
int egress_v4_deny_pktgen(struct __ctx_buff *ctx)
{
	return pktgen(ctx, true);
}

SETUP(PROG_TYPE, "hostfw_ipsec_esp_5_egress_v4_deny")
int egress_v4_deny_setup(struct __ctx_buff *ctx)
{
	return egress_setup(ctx, true, false);
}

CHECK(PROG_TYPE, "hostfw_ipsec_esp_5_egress_v4_deny")
int egress_v4_deny_check(const struct __ctx_buff *ctx)
{
	return egress_check(ctx, true, false);
}

PKTGEN(PROG_TYPE, "hostfw_ipsec_esp_6_egress_v4_allow")
int egress_v4_allow_pktgen(struct __ctx_buff *ctx)
{
	return pktgen(ctx, true);
}

SETUP(PROG_TYPE, "hostfw_ipsec_esp_6_egress_v4_allow")
int egress_v4_allow_setup(struct __ctx_buff *ctx)
{
	return egress_setup(ctx, true, true);
}

CHECK(PROG_TYPE, "hostfw_ipsec_esp_6_egress_v4_allow")
int egress_v4_allow_check(const struct __ctx_buff *ctx)
{
	return egress_check(ctx, true, true);
}
