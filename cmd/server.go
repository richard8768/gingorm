package cmd

import (
	"fmt"
	"gin_demo/internal/config"
	zapLogger "gin_demo/internal/logger"
	middlewares "gin_demo/internal/middleware"
	"gin_demo/internal/router"
	"gin_demo/internal/util"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the application",
	Run: func(cmd *cobra.Command, args []string) {
		if err := util.InitTrans("zh"); err != nil {
			fmt.Println("Init translator failed.")
		}
		port := strconv.Itoa(config.GetHttpPort())
		mode := config.GetRunMode()
		gin.SetMode(mode)
		logger, _ := zapLogger.SetupLogger()
		r := gin.New()
		//r.Use(gin.Logger(), gin.Recovery())
		r.Use(zapLogger.GinLogger(logger), zapLogger.GinRecovery(logger, true))
		r = middlewares.SetupMiddlewares(r)
		r = router.SetupRouter(r)

		s := endless.NewServer(":"+port, r)
		s.ReadHeaderTimeout = 20 * time.Second
		s.WriteTimeout = 20 * time.Second
		s.MaxHeaderBytes = 1 << 20

		s.BeforeBegin = func(addr string) {
			log.Printf("Actual pid is %d", syscall.Getpid())
		}

		// 记录 pid 到文件
		if err := recordPID(); err != nil {
			log.Fatalf("Failed to write PID file: %v", err)
		}

		log.Println("Starting server...")
		if err := s.ListenAndServe(); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Stopping application...")
		pid, err := getPID()
		if err != nil {
			log.Fatalf("Failed to write PID file: %v", err)
		}

		// 发送 SIGTERM 信号优雅地关闭进程
		if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
			log.Fatalf("Failed to write PID file: %v", err)
		}

		log.Println("Server stopped gracefully")
		removePIDFile() // 确保在退出时删除 PID 文件

	},
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the application",
	Run: func(cmd *cobra.Command, args []string) {
		// 在这里重启你的服务或任何其他重启逻辑，通常先停止后启动
		log.Println("Restarting the application...")
		fmt.Println("Restarting application...")
		pid, err := getPID()
		if err != nil {
			log.Fatalf("Failed to write PID file: %v", err)
		}

		if err := syscall.Kill(pid, syscall.SIGHUP); err != nil {
			log.Fatalf("Failed to write PID file: %v", err)
		}
	},
}

func getPID() (int, error) {
	file, err := os.Open(pidFilePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var pid int
	_, err = fmt.Fscanf(file, "%d", &pid)
	if err != nil {
		return 0, err
	}
	return pid, nil
}

// 记录 pid 到指定文件
func recordPID() error {
	pid := syscall.Getpid()
	file, err := os.Create(pidFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%d\n", pid)
	if err != nil {
		return err
	}
	return nil
}

func removePIDFile() {
	if err := os.Remove(pidFilePath); err != nil {
		log.Fatalf("Failed to remove PID file: %v", err)
	}
}
