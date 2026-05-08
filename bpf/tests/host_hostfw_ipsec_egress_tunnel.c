// SPDX-License-Identifier: (GPL-2.0-only OR BSD-2-Clause)
/* Copyright Authors of Cilium */

/* Egress hostFW skip on encrypted recirculation — tunnel routing variant.
 *
 * Mirrors host_hostfw_ipsec_egress.c with TUNNEL_MODE / ENCAP_IFINDEX
 * defined so that the cil_to_netdev tunnel-specific conditional compilation
 * (e.g. the bpf/lib/ipsec.h:199-218 tunnel-mode branch in
 * ipsec_maybe_redirect_to_encrypt) is exercised. The skip is a direct
 * ctx_is_encrypt() mark check at the top of the egress hostFW block, so
 * it must fire identically in both modes — this file proves that.
 */

#define TUNNEL_MODE			1
#define ENCAP_IFINDEX			42

#define ENABLE_IPSEC			1

#include "host_hostfw_ipsec_egress.h"
