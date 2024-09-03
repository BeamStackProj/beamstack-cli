package create

import (
	"fmt"

	"github.com/BeamStackProj/beamstack-cli/src/objects"
	"github.com/BeamStackProj/beamstack-cli/src/types"
	"github.com/BeamStackProj/beamstack-cli/src/utils"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	esLongDesc = utils.LongDesc(`
		Create a elastic search with specified requirments.
		`)

	ElasticSearchVersion string = "8.15.0"
	Nodes                uint16 = 1
	esNamespace          string = "default"
)

// infoCmd represents the info command
var ElasticSearchCmd = &cobra.Command{
	Use:   "elasticsearch [NAME]",
	Short: "create elasticsearch cluster",
	Long:  esLongDesc,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("es command requires exactly one argument: cluster Name. Provided %d arguments", len(args))
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		spec := types.EsSpec{
			Version: ElasticSearchVersion,
			NodeSets: []types.EsNodeSet{
				{
					Name:  args[0],
					Count: Nodes,
					Config: map[string]interface{}{
						"node.store.allow_mmap": false,
					},
				},
			},
		}

		err := objects.CreateDynamicResource(
			metav1.TypeMeta{
				APIVersion: "elasticsearch.k8s.elastic.co/v1",
				Kind:       "Elasticsearch",
			},
			metav1.ObjectMeta{
				Name:      args[0],
				Namespace: esNamespace,
			},
			spec,
			"elasticsearches",
		)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Elastic search created")
		fmt.Printf("incluster url: http://%s-es-http.default.svc.cluster.local:9200 \n", args[0])
		fmt.Printf("Elasticsearch Password can be retrieved using kubectl \n kubectl get secret %s-es-elastic-user -n %s -o go-template='{{.data.elastic | base64decode}}'\n", args[0], esNamespace)
	},
}

func init() {
	ElasticSearchCmd.Flags().StringVar(&ElasticSearchVersion, "version", ElasticSearchVersion, "elastic search version")
	ElasticSearchCmd.Flags().Uint16Var(&Nodes, "nodes", Nodes, "number of elastic search Nodes")
	ElasticSearchCmd.Flags().StringVar(&esNamespace, "namespace", esNamespace, "namespace to deploy create elasticsearch cluster on")
}
