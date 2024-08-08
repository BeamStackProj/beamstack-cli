package objects

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
)

func applyYAML(dynamicClient dynamic.Interface, path string) error {
	var data []byte
	var err error

	if isURL(path) {
		data, err = downloadFile(path)
		if err != nil {
			return err
		}
	} else {
		data, err = os.ReadFile(path)
		if err != nil {
			return err
		}
	}

	yamlDecoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(data), 100)
	for {
		var rawObj runtime.RawExtension
		if err := yamlDecoder.Decode(&rawObj); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		obj := &unstructured.Unstructured{}
		if err := json.Unmarshal(rawObj.Raw, obj); err != nil {
			return err
		}

		gvk := obj.GroupVersionKind()
		resourceClient := dynamicClient.Resource(
			schema.GroupVersionResource{
				Group:    gvk.Group,
				Version:  gvk.Version,
				Resource: strings.ToLower(gvk.Kind) + "s",
			},
		).Namespace(obj.GetNamespace())

		existing, err := resourceClient.Get(context.Background(), obj.GetName(), metav1.GetOptions{})

		if errors.IsNotFound(err) {
			// Create the resource
			_, err = resourceClient.Create(context.Background(), obj, metav1.CreateOptions{})
			if err != nil {
				return err
			}
			// fmt.Printf("Created %s %s/%s\n", gvk.Kind, obj.GetNamespace(), obj.GetName())
		} else if err != nil {
			return err
		} else {
			return nil //TODO: fix patch bug
			// Update the resource
			existingJSON, err := json.Marshal(existing.Object)
			if err != nil {
				return err
			}

			rawObjJSON, err := json.Marshal(rawObj.Object)
			if err != nil {
				return err
			}
			patch, err := strategicpatch.CreateTwoWayMergePatch(existingJSON, rawObjJSON, obj)
			if err != nil {
				return err
			}
			_, err = resourceClient.Patch(context.Background(), obj.GetName(), types.StrategicMergePatchType, patch, metav1.PatchOptions{})
			if err != nil {
				return err
			}
			fmt.Printf("Updated %s %s/%s\n", gvk.Kind, obj.GetNamespace(), obj.GetName())
		}
	}

	return nil
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func SplitAPIVersion(apiVersion string) (string, string, error) {
	parts := strings.Split(apiVersion, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid apiVersion format: %s", apiVersion)
	}
	return parts[0], parts[1], nil
}

func toUnstructured(obj interface{}) (*unstructured.Unstructured, error) {
	unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}
	return &unstructured.Unstructured{Object: unstructuredMap}, nil
}
