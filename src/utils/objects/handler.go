package objects

import (
	"context"
	"fmt"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"

	"github.com/BeamStackProj/beamstack-cli/src/types"
	"github.com/BeamStackProj/beamstack-cli/src/utils"
)

// HandleResource waits for a specified Kubernetes resource to reach a given condition and signals completion via a channel.
//
// Parameters:
// - resource: The type of Kubernetes resource to monitor. Valid values include "Deployment", "CustomResourceDefinition", "Pod", "Service", and others.
// - namespace: The namespace where the resource is located. Use an empty string for cluster-wide resources.
// - condition: The condition to wait for. Valid values include "Established" for CustomResourceDefinitions, "Available" for Deployments and Pods, etc.
// - chanel: A channel used to signal progress. The channel is closed when the condition is met or an error occurs.
func HandleResource(resource string, namespace string, condition string, chanel chan types.ProgCount) {
	defer close(chanel)
	config := utils.GetKubeConfig()
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	gvr, err := getResourceGVR(discoveryClient, resource)
	if err != nil {
		panic(err.Error())
	}
	waitForResourceCondition(dynamicClient, gvr, namespace, condition, chanel)
}

func waitForResourceCondition(client dynamic.Interface, gvr schema.GroupVersionResource, namespace string, condition string, chanel chan types.ProgCount) {
	finishedItems := []unstructured.Unstructured{}

	for {
		var (
			resourceList *unstructured.UnstructuredList
			err          error
		)

		if namespace != "" {
			resourceList, err = client.Resource(gvr).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
		} else {
			resourceList, err = client.Resource(gvr).List(context.TODO(), metav1.ListOptions{})
		}
		if err != nil {
			panic(err.Error())
		}
		if len(resourceList.Items) == len(finishedItems) {
			return
		}

		chanel <- types.ProgCount{OnInit: true, Count: len(resourceList.Items)}

		for _, item := range resourceList.Items {
			err := wait.PollUntilContextTimeout(context.Background(), time.Second*2, time.Minute*10, true, func(context.Context) (bool, error) {
				var (
					res *unstructured.Unstructured
					err error
				)

				if namespace != "" {
					res, err = client.Resource(gvr).Namespace(namespace).Get(context.TODO(), item.GetName(), metav1.GetOptions{})
				} else {
					res, err = client.Resource(gvr).Get(context.TODO(), item.GetName(), metav1.GetOptions{})
				}

				if err != nil {
					return false, err
				}

				conditions, found, err := unstructuredNestedSlice(res.Object, "status", "conditions")
				if err != nil || !found {
					return false, nil
				}

				for _, cond := range conditions {
					if condMap, ok := cond.(map[string]interface{}); ok {
						if condType, found := condMap["type"].(string); found && condType == condition {
							if condStatus, found := condMap["status"].(string); found && condStatus == "True" {
								finishedItems = append(finishedItems, item)
								chanel <- types.ProgCount{OnInit: false, Count: 1}
								return true, nil
							}
						}
					}
				}
				return false, nil
			})
			if err != nil {
				fmt.Printf("Error waiting for %s %s to be %s: %v\n", gvr.Resource, item.GetName(), condition, err)
			}
		}
		time.Sleep(2 * time.Second)
	}
}

func getResourceGVR(discoveryClient *discovery.DiscoveryClient, resource string) (schema.GroupVersionResource, error) {
	apiResourceLists, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	for _, apiResourceList := range apiResourceLists {
		for _, apiResource := range apiResourceList.APIResources {
			if strings.EqualFold(apiResource.Kind, resource) {
				groupVersion, err := schema.ParseGroupVersion(apiResourceList.GroupVersion)
				if err != nil {
					return schema.GroupVersionResource{}, err
				}
				return schema.GroupVersionResource{
					Group:    groupVersion.Group,
					Version:  groupVersion.Version,
					Resource: apiResource.Name,
				}, nil
			}
		}
	}
	return schema.GroupVersionResource{}, fmt.Errorf("resource %s not found", resource)
}

func unstructuredNestedSlice(obj map[string]interface{}, fields ...string) ([]interface{}, bool, error) {
	v, found, err := unstructured.NestedFieldNoCopy(obj, fields...)
	if !found || err != nil {
		return nil, found, err
	}
	if s, ok := v.([]interface{}); ok {
		return s, true, nil
	}
	return nil, false, fmt.Errorf("expected []interface{} but got %T", v)
}
