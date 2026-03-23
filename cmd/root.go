// Package cmd /*
package cmd

import (
	"fmt"
	"gin_demo/internal/config"
	"os"

	"github.com/spf13/cobra"
)

var AppConfig *config.AllConfig
var pidFilePath = "gingorm.pid" // 指定 pid 文件路径

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "main.go",
	Short: "manage application",
	Long:  ` you can start stop or restart the application `,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "help" {
			return nil
		}
		// 加载配置文件
		var err error
		AppConfig, err = config.InitLoadConfig()
		if err != nil {
			return fmt.Errorf("加载配置文件失败: %w", err)
		}
		mysqlSetupErr := config.SetupDB(AppConfig.DataBase)
		if mysqlSetupErr != nil {
			panic(fmt.Errorf("init mysql failed, err:%v\n", mysqlSetupErr))
			return nil
		}

		redisSetupErr := config.SetupRedis(AppConfig.Redis)
		if redisSetupErr != nil {
			panic(fmt.Errorf("init redis failed, err:%v\n", redisSetupErr))
			return nil
		}

		return nil
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(restartCmd)
	rootCmd.AddCommand(genCmd)
}
