package types

var BaseKubePrometheusValues = map[string]interface{}{
	"fullnameOverride": "prometheus",
	"defaultRules": map[string]interface{}{
		"create": true,
		"rules": map[string]interface{}{
			"alertmanager":                true,
			"etcd":                        true,
			"configReloaders":             true,
			"general":                     true,
			"k8s":                         true,
			"kubeApiserverAvailability":   true,
			"kubeApiserverBurnrate":       true,
			"kubeApiserverHistogram":      true,
			"kubeApiserverSlos":           true,
			"kubelet":                     true,
			"kubeProxy":                   true,
			"kubePrometheusGeneral":       true,
			"kubePrometheusNodeRecording": true,
			"kubernetesApps":              true,
			"kubernetesResources":         true,
			"kubernetesStorage":           true,
			"kubernetesSystem":            true,
			"kubeScheduler":               true,
			"kubeStateMetrics":            true,
			"network":                     true,
			"node":                        true,
			"nodeExporterAlerting":        true,
			"nodeExporterRecording":       true,
			"prometheus":                  true,
			"prometheusOperator":          true,
		},
	},
	"alertmanager": map[string]interface{}{
		"fullnameOverride": "alertmanager",
		"enabled":          true,
		"ingress": map[string]interface{}{
			"enabled": false,
		},
	},
	"grafana": map[string]interface{}{
		"enabled":                   true,
		"fullnameOverride":          "grafana",
		"forceDeployDatasources":    false,
		"forceDeployDashboards":     false,
		"defaultDashboardsEnabled":  true,
		"defaultDashboardsTimezone": "utc",
		"serviceMonitor": map[string]interface{}{
			"enabled": true,
		},
		"admin": map[string]interface{}{
			"existingSecret": "grafana-admin-credentials",
			"userKey":        "admin-user",
			"passwordKey":    "admin-password",
		},
	},
	"kubeApiServer": map[string]interface{}{
		"enabled": true,
	},
	"kubelet": map[string]interface{}{
		"enabled": true,
		"serviceMonitor": map[string]interface{}{
			"metricRelabelings": []interface{}{
				map[string]interface{}{
					"action":       "replace",
					"sourceLabels": []interface{}{"node"},
					"targetLabel":  "instance",
				},
			},
		},
	},
	"kubeControllerManager": map[string]interface{}{
		"enabled":   true,
		"endpoints": []interface{}{},
	},
	"coreDns": map[string]interface{}{
		"enabled": true,
	},
	"kubeDns": map[string]interface{}{
		"enabled": false,
	},
	"kubeEtcd": map[string]interface{}{
		"enabled":   true,
		"endpoints": []interface{}{},
		"service": map[string]interface{}{
			"enabled":    true,
			"port":       2381,
			"targetPort": 2381,
		},
	},
	"kubeScheduler": map[string]interface{}{
		"enabled":   true,
		"endpoints": []interface{}{},
	},
	"kubeProxy": map[string]interface{}{
		"enabled":   true,
		"endpoints": []interface{}{},
	},
	"kubeStateMetrics": map[string]interface{}{
		"enabled": true,
	},
	"kube-state-metrics": map[string]interface{}{
		"fullnameOverride": "kube-state-metrics",
		"selfMonitor": map[string]interface{}{
			"enabled": true,
		},
		"prometheus": map[string]interface{}{
			"monitor": map[string]interface{}{
				"enabled": true,
				"relabelings": []interface{}{
					map[string]interface{}{
						"action":       "replace",
						"regex":        "(.*)",
						"replacement":  "$1",
						"sourceLabels": []interface{}{"__meta_kubernetes_pod_node_name"},
						"targetLabel":  "kubernetes_node",
					},
				},
			},
		},
	},
	"nodeExporter": map[string]interface{}{
		"enabled": true,
		"serviceMonitor": map[string]interface{}{
			"relabelings": []interface{}{
				map[string]interface{}{
					"action":       "replace",
					"regex":        "(.*)",
					"replacement":  "$1",
					"sourceLabels": []interface{}{"__meta_kubernetes_pod_node_name"},
					"targetLabel":  "kubernetes_node",
				},
			},
		},
	},
	"prometheus-node-exporter": map[string]interface{}{
		"fullnameOverride": "node-exporter",
		"podLabels": map[string]interface{}{
			"jobLabel": "node-exporter",
		},
		"extraArgs": []interface{}{
			"--collector.filesystem.mount-points-exclude=^/(dev|proc|sys|var/lib/docker/.+|var/lib/kubelet/.+)($|/)",
			"--collector.filesystem.fs-types-exclude=^(autofs|binfmt_misc|bpf|cgroup2?|configfs|debugfs|devpts|devtmpfs|fusectl|hugetlbfs|iso9660|mqueue|nsfs|overlay|proc|procfs|pstore|rpc_pipefs|securityfs|selinuxfs|squashfs|sysfs|tracefs)$",
		},
		"service": map[string]interface{}{
			"portName": "http-metrics",
		},
		"prometheus": map[string]interface{}{
			"monitor": map[string]interface{}{
				"enabled": true,
				"relabelings": []interface{}{
					map[string]interface{}{
						"action":       "replace",
						"regex":        "(.*)",
						"replacement":  "$1",
						"sourceLabels": []interface{}{"__meta_kubernetes_pod_node_name"},
						"targetLabel":  "kubernetes_node",
					},
				},
			},
		},
		"resources": map[string]interface{}{
			"requests": map[string]interface{}{
				"memory": "512Mi",
				"cpu":    "250m",
			},
			"limits": map[string]interface{}{
				"memory": "2048Mi",
			},
		},
	},
	"prometheusOperator": map[string]interface{}{
		"enabled": true,
		"prometheusConfigReloader": map[string]interface{}{
			"resources": map[string]interface{}{
				"requests": map[string]interface{}{
					"cpu":    "200m",
					"memory": "50Mi",
				},
				"limits": map[string]interface{}{
					"memory": "200Mi",
				},
			},
		},
	},
	"prometheus": map[string]interface{}{
		"enabled": true,
		"prometheusSpec": map[string]interface{}{
			"replicas":                                1,
			"replicaExternalLabelName":                "replica",
			"ruleSelectorNilUsesHelmValues":           false,
			"serviceMonitorSelectorNilUsesHelmValues": false,
			"podMonitorSelectorNilUsesHelmValues":     false,
			"probeSelectorNilUsesHelmValues":          false,
			"retention":                               "6h",
			"enableAdminAPI":                          true,
			"walCompression":                          true,
		},
	},
	"thanosRuler": map[string]interface{}{
		"enabled": false,
	},
}

func updateEndpoints(values map[string]interface{}, key string, nodeEndpoints []interface{}) {
	if component, ok := values[key].(map[string]interface{}); ok {
		component["endpoints"] = nodeEndpoints
		values[key] = component
	}
}

func GetKubePrometheusValues(nodeEndpoints []interface{}) *map[string]interface{} {
	keysToUpdate := []string{"kubeControllerManager", "kubeEtcd", "kubeScheduler", "kubeProxy"}
	for _, key := range keysToUpdate {
		updateEndpoints(BaseKubePrometheusValues, key, nodeEndpoints)
	}
	return &BaseKubePrometheusValues
}
