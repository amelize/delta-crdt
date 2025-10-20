package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/amelize/delta-crdt/server/internal/config"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-yaml"
	"github.com/hashicorp/memberlist"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	rootCmd = &cobra.Command{
		Use:   "server",
		Short: "CRDT Example",
		Long:  `Server for crds based data`,
	}
)

func getConfig(logger *zap.Logger) config.Config {
	config := config.Config{}

	cobra.OnInitialize(func() {

	})

	address := rootCmd.PersistentFlags().StringP("address", "b", "", "Servers with ip 0.0.0.0:XXXX")
	servers := rootCmd.PersistentFlags().StringP("servers", "s", "", "Servers with ip 0.0.0.0:XXXX")
	listenPort := rootCmd.PersistentFlags().Int32P("port", "p", 0, "port 2 listen")
	clusterPort := rootCmd.PersistentFlags().Int32P("cluster-port", "c", 0, "port 2 listen")
	advListenPort := rootCmd.PersistentFlags().Int32P("advertise-port", "a", 0, "port 2 listen")

	rootCmd.Execute()

	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		logger.Warn("Config error", zap.Error(err))
	} else {
		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			logger.Fatal("Config error", zap.Error(err))
		}
	}

	if address != nil {
		config.BindAddress = *address
	}

	if servers != nil {
		config.Cluster.Servers = strings.Split(strings.ReplaceAll(*servers, " ", ""), ",")
	}

	if listenPort != nil {
		config.Port = *listenPort
	}

	if clusterPort != nil {
		config.ClusterPort = *clusterPort
	}

	if advListenPort != nil {
		config.AdvertisePort = *advListenPort
	}

	return config
}

func main() {

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	logger := zap.New(core)
	defer logger.Sync()

	config := getConfig(logger)

	router := gin.Default()
	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	router.Use(ginzap.RecoveryWithZap(logger, true))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router.Handler(),
	}

	memberListConfig := memberlist.DefaultLANConfig()

	memberListConfig.Name = fmt.Sprintf("node-%d", config.ClusterPort)

	memberListConfig.BindAddr = config.BindAddress
	memberListConfig.BindPort = int(config.ClusterPort)

	memberListConfig.AdvertiseAddr = config.BindAddress
	memberListConfig.AdvertisePort = int(config.AdvertisePort)

	list, err := memberlist.Create(memberListConfig)
	if err != nil {
		logger.Fatal("Failed to create memberlist: " + err.Error())
	}

	localNode := list.LocalNode()
	logger.Info("Local node", zap.String("node", string(localNode.Address())))

	_, err = list.Join(config.Cluster.Servers)
	if err != nil {
		logger.Warn("Failed to join cluster", zap.Error(err))
	}

	// Ask for members of the cluster
	for _, member := range list.Members() {
		fmt.Printf("Member: %s %s\n", member.Name, member.Addr)
	}

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome CRDT")
	})

	router.GET("/nodes", func(c *gin.Context) {
		var nodeList strings.Builder
		for _, member := range list.Members() {
			nodeList.WriteString(member.Address())
			nodeList.WriteString(",")
		}
		c.String(http.StatusOK, nodeList.String())
	})

	go func() {
		logger.Info("Starting server ", zap.Int32("port", config.Port), zap.String("address", config.BindAddress))

		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no params) by default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Warn("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown memberlist
	if err := list.Leave(time.Second); err != nil {
		logger.Warn("Server Shutdown - leave memberlist:", zap.Error(err))
	}

	// Shutdown memberlist
	if err := list.Shutdown(); err != nil {
		logger.Warn("Server Shutdown - memberlist:", zap.Error(err))
	}

	if err := srv.Shutdown(ctx); err != nil {
		logger.Warn("Server Shutdown:", zap.Error(err))
	}

	logger.Info("Server exiting")

}
