// SPDX-License-Identifier: (GPL-2.0-only OR BSD-2-Clause)
/* Copyright Authors of Cilium */

/* cil_to_host IPsec recirc + NodePort smoke test — tunnel routing variant. */

#define TUNNEL_MODE			1
#define ENCAP_IFINDEX			42

#define ENABLE_IPSEC			1

#include "host_hostfw_ipsec_to_host.h"
