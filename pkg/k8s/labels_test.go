// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package k8s

import (
	"testing"

	"github.com/cilium/hive/hivetest"
	"github.com/stretchr/testify/require"

	slim_corev1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/api/core/v1"
	slim_metav1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/apis/meta/v1"
)

func TestGetPodMetadata(t *testing.T) {
	ns := &slim_corev1.Namespace{
		ObjectMeta: slim_metav1.ObjectMeta{
			Name: "default",
			Labels: map[string]string{
				"kubernetes.io/metadata.name": "default",
				"namespace-level-key":         "namespace-level-value",
			},
		},
	}

	pod := &slim_corev1.Pod{
		ObjectMeta: slim_metav1.ObjectMeta{
			Namespace: "default",
			Labels: map[string]string{
				"app":                         "test",
				"io.kubernetes.pod.namespace": "default",
			},
			Annotations: map[string]string{},
		},
	}

	expectedLabels := map[string]string{
		"app":                          "test",
		"io.cilium.k8s.policy.cluster": "",
		"io.kubernetes.pod.namespace":  "default",
		"io.cilium.k8s.namespace.labels.kubernetes.io/metadata.name": "default",
		"io.cilium.k8s.namespace.labels.namespace-level-key":         "namespace-level-value",
	}

	t.Run("normal scenario", func(t *testing.T) {
		pod := pod.DeepCopy()

		_, labels := GetPodMetadata(hivetest.Logger(t), ns, pod)
		require.Equal(t, expectedLabels, labels)
	})

	t.Run("pod labels contains cilium owned label", func(t *testing.T) {
		t.Run("override namespace labels", func(t *testing.T) {
			pod := pod.DeepCopy()
			pod.Labels["io.cilium.k8s.namespace.labels.namespace-level-key"] = "override-namespace-level-value"

			_, labels := GetPodMetadata(hivetest.Logger(t), ns, pod)
			require.Equal(t, expectedLabels, labels)
		})

		t.Run("add one more namespace labels", func(t *testing.T) {
			pod := pod.DeepCopy()
			pod.Labels["io.cilium.k8s.namespace.labels.another-namespace-key"] = "another-namespace-level-value"

			_, labels := GetPodMetadata(hivetest.Logger(t), ns, pod)
			require.Equal(t, expectedLabels, labels)
		})
	})
}
