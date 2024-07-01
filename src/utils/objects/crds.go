package objects

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"encoding/json"

	"k8s.io/apimachinery/pkg/util/yaml"

	"bytes"

	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/dynamic"

	"github.com/BeamStackProj/beamstack-cli/src/utils"
)

func CreateObject(path string) error {

	config := utils.GetKubeConfig()
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Example: Apply a YAML file from a URL
	err = applyYAML(dynamicClient, path)
	if err != nil {
		return fmt.Errorf("error applying yaml from url %v", err)
	}
	return nil
}

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

		existing, err := resourceClient.Get(context.Background(), obj.GetName(), v1.GetOptions{})

		if errors.IsNotFound(err) {
			// Create the resource
			_, err = resourceClient.Create(context.Background(), obj, v1.CreateOptions{})
			if err != nil {
				return err
			}
			fmt.Printf("Created %s %s/%s\n", gvk.Kind, obj.GetNamespace(), obj.GetName())
		} else if err != nil {
			return err
		} else {
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
			fmt.Println(obj)
			fmt.Println(patch)
			_, err = resourceClient.Patch(context.Background(), obj.GetName(), types.StrategicMergePatchType, patch, v1.PatchOptions{})
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
