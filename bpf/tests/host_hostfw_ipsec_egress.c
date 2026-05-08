// SPDX-License-Identifier: (GPL-2.0-only OR BSD-2-Clause)
/* Copyright Authors of Cilium */

/* Egress hostFW skip on encrypted recirculation — native routing variant.
 *
 * The shared test logic lives in host_hostfw_ipsec_egress.h. The tunnel
 * routing variant lives in host_hostfw_ipsec_egress_tunnel.c.
 */

#define ENABLE_IPSEC			1

#include "host_hostfw_ipsec_egress.h"
