package objects

import (
	"context"
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
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
// - channel: A channel used to signal progress. The channel is closed when the condition is met or an error occurs.
func HandleResources(resource string, namespace string, condition string, channel chan types.ProgCount) {
	defer close(channel)
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
	waitForResourceCondition(dynamicClient, gvr, namespace, condition, channel)
}

func HandleSpecificResource(gvr schema.GroupVersionResource, name, namespace, condition string, channel chan string) {
	defer close(channel)

	config := utils.GetKubeConfig()
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	waitForSpecificResourceCondition(dynamicClient, gvr, namespace, condition, name, channel)
}

// func waitForSpecificResourceCondition(client dynamic.Interface, gvr schema.GroupVersionResource, namespace string, condition string, name string, channel *chan string) error {
// 	var (
// 		resource   *unstructured.Unstructured
// 		err        error
// 		lastStatus string
// 		lastReason string
// 	)

// 	if namespace == "" {
// 		namespace = "default"
// 	}

// 	resource, err = client.Resource(gvr).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
// 	if err != nil {
// 		return err
// 	}

// 	for {
// 		err = wait.PollUntilContextTimeout(context.Background(), 2*time.Second, 10*time.Minute, true, func(ctx context.Context) (bool, error) {
// 			resource, err = client.Resource(gvr).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
// 			if err != nil {
// 				return false, err
// 			}

// 			conditions, found, err := unstructured.NestedSlice(resource.Object, "status", "conditions")
// 			if err != nil || !found {
// 				return false, nil
// 			}

// 			for _, cond := range conditions {
// 				if condMap, ok := cond.(map[string]interface{}); ok {
// 					if condType, found := condMap["type"].(string); found && condType == condition {
// 						if condStatus, found := condMap["status"].(string); found {
// 							reason := ""
// 							if condReason, found := condMap["reason"].(string); found {
// 								reason = condReason
// 							}
// 							if condStatus != lastStatus || reason != lastReason {
// 								lastStatus = condStatus
// 								lastReason = reason
// 								*channel <- fmt.Sprintf("Status: %s, Reason: %s", condStatus, reason)
// 							}
// 							if condStatus == "True" {
// 								return true, nil
// 							}
// 						}
// 					}
// 				}
// 			}
// 			return false, nil
// 		})
// 		if err != nil {
// 			fmt.Printf("Error waiting for %s %s to be %s: %v\n", gvr.Resource, resource.GetName(), condition, err)
// 		}
// 		time.Sleep(2 * time.Second)
// 	}
// }

func waitForSpecificResourceCondition(client dynamic.Interface, gvr schema.GroupVersionResource, namespace string, condition string, name string, channel chan string) error {
	var (
		resource   *unstructured.Unstructured
		err        error
		lastStatus string
		lastReason string
		lastEvents map[string]string = make(map[string]string)
	)

	if namespace == "" {
		namespace = "default"
	}

	resource, err = client.Resource(gvr).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	for {
		err = wait.PollUntilContextTimeout(context.Background(), 2*time.Second, 10*time.Minute, true, func(ctx context.Context) (bool, error) {
			resource, err = client.Resource(gvr).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				if errors.IsNotFound(err) {
					return false, fmt.Errorf("resource %s not found", name)
				}
				return false, err
			}

			conditions, found, err := unstructured.NestedSlice(resource.Object, "status", "conditions")
			if err != nil || !found {
				return false, nil
			}

			eventsUpdated := false

			for _, cond := range conditions {
				if condMap, ok := cond.(map[string]interface{}); ok {
					if condType, found := condMap["type"].(string); found && condType == condition {
						if condStatus, found := condMap["status"].(string); found {
							reason := ""
							if condReason, found := condMap["reason"].(string); found {
								reason = condReason
							}
							if condStatus != lastStatus || reason != lastReason {
								lastStatus = condStatus
								lastReason = reason
								eventsUpdated = true
								channel <- fmt.Sprintf("Condition: %s, Status: %s, Reason: %s", condType, condStatus, reason)
							}
							if condStatus == "True" {
								return true, nil
							}
						}
					}
				}
			}

			// Check for event changes
			events, found, err := unstructured.NestedSlice(resource.Object, "status", "events")
			if found && err == nil {
				for _, event := range events {
					if eventMap, ok := event.(map[string]interface{}); ok {
						eventReason := ""
						if reason, found := eventMap["reason"].(string); found {
							eventReason = reason
						}
						eventMessage := ""
						if message, found := eventMap["message"].(string); found {
							eventMessage = message
						}
						eventKey := fmt.Sprintf("%s: %s", eventReason, eventMessage)
						if _, found := lastEvents[eventKey]; !found {
							eventsUpdated = true
							lastEvents[eventKey] = eventMessage
							channel <- fmt.Sprintf("Event: %s, Message: %s", eventReason, eventMessage)
						}
					}
				}
			}

			if !eventsUpdated {
				return false, nil
			}
			return false, nil
		})
		if err != nil {
			return fmt.Errorf("error waiting for %s %s to be %s: %v", gvr.Resource, resource.GetName(), condition, err)
		}
		time.Sleep(2 * time.Second)
	}
}

func waitForResourceCondition(client dynamic.Interface, gvr schema.GroupVersionResource, namespace string, condition string, channel chan types.ProgCount) {
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

		channel <- types.ProgCount{OnInit: true, Count: len(resourceList.Items)}

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
								channel <- types.ProgCount{OnInit: false, Count: 1}
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
