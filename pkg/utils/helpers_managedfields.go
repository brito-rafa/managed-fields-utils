package utils

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ManagedFieldsMetaSmall() *metav1.FieldsV1 {

	return &metav1.FieldsV1{Raw: []byte(`
	{ 
		"f:metadata": {
			"f:annotations": {
				"f:nm.kubernetes/utan":                               {}
			}
		}
	}
	`)}
}

func HPAManagedFieldsMetaAndSpec() *metav1.FieldsV1 {

	return &metav1.FieldsV1{Raw: []byte(`
	{ 
		"f:metadata": {
			"f:annotations": {
				".":                 {},
				"f:janitor/expires": {},
				"f:kubectl.kubernetes.io/last-applied-configuration": {},
				"f:nm.kubernetes/deploy_date":                        {},
				"f:nm.kubernetes/deployer_email":                     {},
				"f:nm.kubernetes/git_url":                            {},
				"f:nm.kubernetes/utan":                               {}
			},
			"f:labels": {
				".":                    {},
				"f:app":                {},
				"f:app_tertiary":       {},
				"f:caas-test-deleteme": {},
				"f:sha":                {}
			}
		},
		"f:spec": {
			"f:minReplicas":                    {},
			"f:scaleTargetRef":                 {},
			"f:targetCPUUtilizationPercentage": {}
		}
	}
	`)}
}

func HPAManagedFieldsSpecMaxReplica() *metav1.FieldsV1 {

	return &metav1.FieldsV1{Raw: []byte(`
	{
		"f:metadata": {
			"f:labels": {
				"f:k8slens-edit-resource-version": {}
			}
		},
		"f:spec": {
			"f:maxReplicas": {}
		}
	}
	`)}
}

func HPAManagedFieldsStatus() *metav1.FieldsV1 {
	return &metav1.FieldsV1{Raw: []byte(`
	{
		"f:status": {
			"f:conditions": {
				".": {},
				"k:{\"type\":\"AbleToScale\"}": {
					".":                    {},
					"f:lastTransitionTime": {},
					"f:message":            {},
					"f:reason":             {},
					"f:status":             {},
					"f:type":               {}
				},
				"k:{\"type\":\"ScalingActive\"}": {
					".":                    {},
					"f:lastTransitionTime": {},
					"f:message":            {},
					"f:reason":             {},
					"f:status":             {},
					"f:type":               {}
				},
				"k:{\"type\":\"ScalingLimited\"}": {
					".":                    {},
					"f:lastTransitionTime": {},
					"f:message":            {},
					"f:reason":             {},
					"f:status":             {},
					"f:type":               {}
				}
			},
			"f:currentMetrics":  {},
			"f:currentReplicas": {},
			"f:desiredReplicas": {},
			"f:lastScaleTime":   {}
		}
	}
	`)}
}

func HPAManagedFieldsSpecMetrics() *metav1.FieldsV1 {
	return &metav1.FieldsV1{Raw: []byte(` 
	{
		"f:metrics": {}
	}
	`)}
}

func AppsV1ManagedFieldsMetaAndSpecWithoutContainers() *metav1.FieldsV1 {
	return &metav1.FieldsV1{Raw: []byte(`
	{
		"f:metadata": {
			"f:annotations": {
				".": {}
			},
			"f:labels": {
				".": {}
			}
		},
		"f:spec": {
			"f:replicas": {},
			"f:progressDeadlineSeconds": {},
			"f:template": {
				"f:metadata": {
					"f:labels": {
						".": {}
					}
				},
				"f:spec": {
					"f:containers": {}
				}
			}
		}
	}
	`)}
}

func AppsV1ManagedFieldsMetaAndSpecWithContainerArgument() *metav1.FieldsV1 {
	return &metav1.FieldsV1{Raw: []byte(`
{
  "f:metadata": {
    "f:annotations": {
      "f:kubernetes.io/change-cause": {},
      "f:stormforge.io/last-updated": {},
      "f:stormforge.io/recommendation-url": {}
    }
  },
  "f:spec": {
    "f:template": {
      "f:spec": {
        "f:containers": {
          "k:{\"name\":\"nginx\"}": {
            "f:args": {},
            "f:command": {}
          }
        }
      }
    }
  }
}
	`)}
}

func AppsV1ManagedFieldsMetaAndSpec() *metav1.FieldsV1 {
	return &metav1.FieldsV1{Raw: []byte(`
{
  "f:metadata": {
    "f:annotations": {
      "f:kubernetes.io/change-cause": {},
      "f:stormforge.io/last-updated": {},
      "f:stormforge.io/recommendation-url": {}
    }
  },
  "f:spec": {
    "f:template": {
      "f:spec": {
        "f:containers": {
          "k:{\"name\":\"nginx\"}": {
            "f:args": {},
            "f:command": {},
			"f:resources": {
				"f:requests": {},
				"f:limits": {}
			}
          }
        }
      }
    }
  }
}
	`)}
}

func AppsV1ManagedFieldsMetaAndSpecRequests() *metav1.FieldsV1 {
	return &metav1.FieldsV1{Raw: []byte(`
{
  "f:spec": {
    "f:template": {
      "f:spec": {
        "f:containers": {
          "k:{\"name\":\"nginx\"}": {
			"f:resources": {
				"f:requests": {}
			}
          }
        }
      }
    }
  }
}
	`)}
}

func AppsV1ManagedFieldsMetaAndSpecLimits() *metav1.FieldsV1 {
	return &metav1.FieldsV1{Raw: []byte(`
{
  "f:spec": {
    "f:template": {
      "f:spec": {
        "f:containers": {
          "k:{\"name\":\"nginx\"}": {
			"f:resources": {
				"f:limits": {}
			}
          }
        }
      }
    }
  }
}
	`)}
}
