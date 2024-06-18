package initialize

import (
	"fmt"
	"os"

	"github.com/Beamflow/beamflow-cli/src/types"
	"github.com/Beamflow/beamflow-cli/src/utils"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize beamstack on cluster",
	Long:  `Initialize Beamstack in a new Kubernetes environment, setting up essential configurations and prerequisites.`,
	Run:   runInit,
}

var (
	ConfigFile      string = ""
	Name            string = ""
	DefaultOperator string = "flink"
	FlinkVersion    string = "latest"
	SparkVersion    string = "latest"
	monitoring      bool   = false
	Flink           bool   = false
	Spark           bool   = false
	operators       types.Operator
)

func init() {

	InitCmd.Flags().StringVarP(&Name, "name", "n", Name, "Name of profile. will be randomly generated if not provided.")
	InitCmd.Flags().StringVarP(&ConfigFile, "config", "c", ConfigFile, "Path to configuration file.")
	InitCmd.Flags().StringVarP(&DefaultOperator, "default-operator", "d", DefaultOperator, "Default operator.")
	InitCmd.Flags().StringVarP(&FlinkVersion, "flink-version", "f", FlinkVersion, "Flink version to be installed. Ignored if Flink is not specified for installation.")
	InitCmd.Flags().StringVarP(&SparkVersion, "spark-version", "s", SparkVersion, "Spark Version to be installed. Ignored if Spark is not specified for installation.")
	InitCmd.Flags().BoolVarP(&Flink, "flink", "F", Flink, "If specified, flink is installed.")
	InitCmd.Flags().BoolVarP(&Spark, "spark", "S", Spark, "If specified, Spark is installed.")
}

func runInit(cmd *cobra.Command, args []string) {
	fmt.Println("Initializing cluster ! !")
	currentContext, err := utils.GetCurrentContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current context: %v\n", err)
		return
	}
	contextsStringMap := viper.GetStringMapString("contexts")

	if _, ok := contextsStringMap[currentContext]; ok {
		fmt.Println("Current cluster already initialized")
		fmt.Print("Do you want reinitialize? (Y/N) ")

		var userInput string
		_, err := fmt.Scanln(&userInput)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			return
		}

		switch userInput {
		case "Y", "y":
			fmt.Println("ReInitializing...")
		case "N", "n":
			fmt.Println("Aborting...")
			return
		default:
			fmt.Println("Invalid input. Aborting...")
			return
		}
	}

	var Profile types.Profiles

	if ConfigFile == "" {

		if Name == "" {
			Name = uuid.NewString()
		}

		if !Flink && !Spark {
			Flink = true
		}

		if Flink {
			flinkDefault := false
			if DefaultOperator == "flink" {
				flinkDefault = true
			}
			operators.Flink = &types.OperatorDetails{
				Version:   FlinkVersion,
				IsDefault: flinkDefault,
			}
		}

		if Spark {
			sparkDefault := false
			if DefaultOperator == "spark" {
				sparkDefault = true
			}
			operators.Spark = &types.OperatorDetails{
				Version:   SparkVersion,
				IsDefault: sparkDefault,
			}
		}

		Profile = types.Profiles{
			Name:       Name,
			Operators:  operators,
			Monitoring: nil,
			Packages:   []types.Package{},
		}

	} else {
		Profile, err = utils.GetProfile(ConfigFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
	}

	contextsStringMap[currentContext] = Profile.Name

	viper.Set("contexts", contextsStringMap)

	if err := viper.WriteConfig(); err != nil {
		fmt.Printf("Error writing config file: %v\n", err)
	}

}
