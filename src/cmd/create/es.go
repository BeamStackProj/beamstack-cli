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
					Config: types.EsNodeConfig{
						NodeStoreAllowMMAP: false,
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
				Name:      "default",
				Namespace: "elastic-system",
			},
			spec,
			"elasticsearches",
		)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Elastic search created")
		fmt.Printf("incluster url: http://%s.default.svc.cluster.local:9200 \n", args[0])
		fmt.Printf("Es Password can be retrieved using kubectl \n kubectl get secret %s-es-elastic-user -o go-template='{{.data.elastic | base64decode}}'\n", args[0])
	},
}

func init() {
	ElasticSearchCmd.Flags().StringVar(&ElasticSearchVersion, "version", ElasticSearchVersion, "elastic search version")
	ElasticSearchCmd.Flags().Uint16Var(&Nodes, "nodes", Nodes, "number of elastic search Nodes")
}
