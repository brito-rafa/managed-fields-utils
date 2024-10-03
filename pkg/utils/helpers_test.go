package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDetectExternalManagersHPA(t *testing.T) {
	testCases := []struct {
		desc                    string
		managedFields           []metav1.ManagedFieldsEntry
		wasOverwritten          bool
		expectedExternalManager string
		originalManager         string
	}{
		{
			desc:                    "no managed fields",
			managedFields:           []metav1.ManagedFieldsEntry{},
			wasOverwritten:          false,
			expectedExternalManager: "",
			originalManager:         "original-manager",
		},
		{
			desc: "different external managers with different fields",
			managedFields: []metav1.ManagedFieldsEntry{
				{
					APIVersion: "autoscaling/v1",
					FieldsType: "FieldsV1",
					FieldsV1:   HPAManagedFieldsMetaAndSpec(),
					Manager:    "kubectl-client-side-apply",
					Operation:  "Update",
					Time:       &metav1.Time{Time: MustParseTime("2044-06-17T19:56:27Z")},
				},
				{
					APIVersion: "autoscaling/v2",
					FieldsType: "FieldsV1",
					FieldsV1:   HPAManagedFieldsSpecMaxReplica(),
					Manager:    "original-manager",
					Operation:  "Update",
					Time:       &metav1.Time{Time: MustParseTime("2044-06-18T00:20:30Z")},
				},
				{
					APIVersion:  "autoscaling/v2",
					FieldsType:  "FieldsV1",
					FieldsV1:    HPAManagedFieldsStatus(),
					Manager:     "kube-controller-manager",
					Operation:   "Update",
					Subresource: "status",
					Time:        &metav1.Time{Time: MustParseTime("2044-06-18T21:01:10Z")},
				},
			},
			wasOverwritten:          false,
			expectedExternalManager: "",
			originalManager:         "original-manager",
		},
		{
			desc: "a single manager",
			managedFields: []metav1.ManagedFieldsEntry{
				{
					APIVersion: "autoscaling/v1",
					FieldsType: "FieldsV1",
					FieldsV1:   HPAManagedFieldsSpecMetrics(),
					Manager:    "kubectl-client-side-apply",
					Operation:  "Update",
					Time:       &metav1.Time{Time: MustParseTime("2044-06-17T19:56:27Z")},
				},
			},
			wasOverwritten:          false,
			expectedExternalManager: "",
			originalManager:         "original-manager",
		},
		{
			desc: "external manager overwriting original-manager",
			managedFields: []metav1.ManagedFieldsEntry{
				{
					APIVersion: "autoscaling/v1",
					FieldsType: "FieldsV1",
					FieldsV1:   HPAManagedFieldsSpecMetrics(),
					Manager:    "kubectl-client-side-apply",
					Operation:  "Update",
					Time:       &metav1.Time{Time: MustParseTime("2044-06-17T19:56:27Z")},
				},
				{
					APIVersion: "autoscaling/v1",
					FieldsType: "FieldsV1",
					FieldsV1:   HPAManagedFieldsSpecMetrics(),
					Manager:    "original-manager",
					Operation:  "Update",
					Time:       &metav1.Time{Time: MustParseTime("2024-06-17T19:56:27Z")},
				},
			},
			wasOverwritten:          true,
			expectedExternalManager: "kubectl-client-side-apply",
			originalManager:         "original-manager",
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%q", tc.desc), func(t *testing.T) {
			wasOverwritten, manager := DetectExternalManager(tc.originalManager, tc.managedFields)
			assert.Equal(t, tc.wasOverwritten, wasOverwritten)
			assert.Equal(t, tc.expectedExternalManager, manager)
		})
	}
}

func TestDetectExternalManagersAppsV1(t *testing.T) {
	testCases := []struct {
		desc                    string
		managedFields           []metav1.ManagedFieldsEntry
		wasOverwritten          bool
		expectedExternalManager string
		originalManager         string
	}{
		{
			desc:                    "no managed fields",
			managedFields:           []metav1.ManagedFieldsEntry{},
			wasOverwritten:          false,
			expectedExternalManager: "",
			originalManager:         "original-manager",
		},
		{
			desc: "external manager and no original-manager as manager",
			managedFields: []metav1.ManagedFieldsEntry{
				{
					APIVersion: "apps/v1",
					FieldsType: "FieldsV1",
					FieldsV1:   AppsV1ManagedFieldsMetaAndSpec(),
					Manager:    "kubectl-client-side-apply",
					Operation:  "Update",
					Time:       &metav1.Time{Time: MustParseTime("2044-06-17T19:56:27Z")},
				},
			},
			wasOverwritten:          false,
			expectedExternalManager: "",
			originalManager:         "original-manager",
		},
		{
			desc: "original-manager being overwritten by external manager",
			managedFields: []metav1.ManagedFieldsEntry{
				{
					APIVersion: "apps/v1",
					FieldsType: "FieldsV1",
					FieldsV1:   AppsV1ManagedFieldsMetaAndSpecRequests(),
					Manager:    "kubectl-client-side-apply",
					Operation:  "Update",
					Time:       &metav1.Time{Time: MustParseTime("2044-06-17T19:56:27Z")},
				},
				{
					APIVersion: "apps/v1",
					FieldsType: "FieldsV1",
					FieldsV1:   AppsV1ManagedFieldsMetaAndSpecRequests(),
					Manager:    "original-manager",
					Operation:  "Update",
					Time:       &metav1.Time{Time: MustParseTime("2024-06-17T19:56:27Z")},
				},
			},
			wasOverwritten:          true,
			expectedExternalManager: "kubectl-client-side-apply",
			originalManager:         "original-manager",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%q", tc.desc), func(t *testing.T) {
			wasOverwritten, manager := DetectExternalManager(tc.originalManager, tc.managedFields)
			assert.Equal(t, tc.wasOverwritten, wasOverwritten)
			assert.Equal(t, tc.expectedExternalManager, manager)
		})
	}

}

func TestFieldsV1ToJSONPaths(t *testing.T) {
	testCases := []struct {
		desc              string
		managedFieldsV1   *metav1.FieldsV1
		expectedJSONPaths []string
		regex             bool
	}{
		{
			desc:            "meta with one annotation",
			managedFieldsV1: ManagedFieldsMetaSmall(),
			expectedJSONPaths: []string{
				"/metadata/annotations/nm.kubernetes/utan",
			},
			regex: false,
		},
		{
			desc:            "hpa with annotation and spec",
			managedFieldsV1: HPAManagedFieldsMetaAndSpec(),
			expectedJSONPaths: []string{
				"/metadata/annotations/.",
				"/metadata/annotations/janitor/expires",
				"/metadata/annotations/kubectl.kubernetes.io/last-applied-configuration",
				"/metadata/annotations/nm.kubernetes/deploy_date",
				"/metadata/annotations/nm.kubernetes/deployer_email",
				"/metadata/annotations/nm.kubernetes/git_url",
				"/metadata/annotations/nm.kubernetes/utan",
				"/metadata/labels/.",
				"/metadata/labels/app",
				"/metadata/labels/app_tertiary",
				"/metadata/labels/caas-test-deleteme",
				"/metadata/labels/sha",
				"/spec/minReplicas",
				"/spec/scaleTargetRef",
				"/spec/targetCPUUtilizationPercentage",
			},
			regex: false,
		},
		{
			desc:            "hpa with annotation and spec with regex",
			managedFieldsV1: HPAManagedFieldsMetaAndSpec(),
			expectedJSONPaths: []string{
				"\\/metadata\\/annotations\\/*.*",
				"\\/metadata\\/annotations\\/janitor\\/expires",
				"\\/metadata\\/annotations\\/kubectl.kubernetes.io\\/last-applied-configuration",
				"\\/metadata\\/annotations\\/nm.kubernetes\\/deploy_date",
				"\\/metadata\\/annotations\\/nm.kubernetes\\/deployer_email",
				"\\/metadata\\/annotations\\/nm.kubernetes\\/git_url",
				"\\/metadata\\/annotations\\/nm.kubernetes\\/utan",
				"\\/metadata\\/labels\\/*.*",
				"\\/metadata\\/labels\\/app",
				"\\/metadata\\/labels\\/app_tertiary",
				"\\/metadata\\/labels\\/caas-test-deleteme",
				"\\/metadata\\/labels\\/sha",
				"\\/spec\\/minReplicas",
				"\\/spec\\/scaleTargetRef",
				"\\/spec\\/targetCPUUtilizationPercentage",
			},
			regex: true,
		},
		{
			desc:            "appsv1 with annotation and container key",
			managedFieldsV1: AppsV1ManagedFieldsMetaAndSpecWithContainerArgument(),
			expectedJSONPaths: []string{
				"/metadata/annotations/kubernetes.io/change-cause",
				"/metadata/annotations/stormforge.io/last-updated",
				"/metadata/annotations/stormforge.io/recommendation-url",
				"/spec/template/spec/containers/[{\"name\":\"nginx\"}]/args",
				"/spec/template/spec/containers/[{\"name\":\"nginx\"}]/command",
			},
			regex: false,
		},
		{
			desc:            "appsv1 with annotation and container key with regex",
			managedFieldsV1: AppsV1ManagedFieldsMetaAndSpecWithContainerArgument(),
			expectedJSONPaths: []string{
				"\\/metadata\\/annotations\\/kubernetes.io\\/change-cause",
				"\\/metadata\\/annotations\\/stormforge.io\\/last-updated",
				"\\/metadata\\/annotations\\/stormforge.io\\/recommendation-url",
				"\\/spec\\/template\\/spec\\/containers\\/*.*\\/args",
				"\\/spec\\/template\\/spec\\/containers\\/*.*\\/command",
			},
			regex: true,
		},
		{
			desc:            "hpa status with key and dot with regex",
			managedFieldsV1: HPAManagedFieldsStatus(),
			expectedJSONPaths: []string{
				"\\/status\\/conditions\\/*.*",
				"\\/status\\/conditions\\/*.*\\/*.*",
				"\\/status\\/conditions\\/*.*\\/lastTransitionTime",
				"\\/status\\/conditions\\/*.*\\/message",
				"\\/status\\/conditions\\/*.*\\/reason",
				"\\/status\\/conditions\\/*.*\\/status",
				"\\/status\\/conditions\\/*.*\\/type",
				"\\/status\\/currentMetrics",
				"\\/status\\/currentReplicas",
				"\\/status\\/desiredReplicas",
				"\\/status\\/lastScaleTime",
			},
			regex: true,
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%q", tc.desc), func(t *testing.T) {

			jsonPaths, err := FieldsV1ToJSONPaths(tc.managedFieldsV1)
			if tc.regex {
				jsonPaths, err = FieldsV1ToJSONPaths(tc.managedFieldsV1, tc.regex)

			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedJSONPaths, jsonPaths)
		})
	}

}
