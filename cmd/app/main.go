package app

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres driver
	"github.com/oklog/oklog/pkg/group"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/nathanows/elegant-monolith/_protos/companyusers"
	"github.com/nathanows/elegant-monolith/internal/company"
	"github.com/nathanows/elegant-monolith/pkg/conf"
)

var rootCmd = &cobra.Command{
	Use:   "elegant-monolith",
	Short: "The elegant-monolith is a sample go-kit api.",
	Run:   run,
}

// Execute is the entry point for spf13/cobra run from the projects main.go
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	var logger log.Logger
	{
		logger = log.NewJSONLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var config *Config
	{
		config = &Config{}
		parser, err := conf.NewParser(cmd, config, conf.EnvPrefix("EM"))
		if err != nil {
			logger.Log("config", "parser_init_err", "during", "NewParser", "err", err)
			os.Exit(1)
		}
		if err := parser.LoadConfig(); err != nil {
			logger.Log("config", "load_err", "during", "LoadConfig", "err", err)
			os.Exit(1)
		}
	}

	var db *sqlx.DB
	{
		var err error
		db, err = sqlx.Connect("postgres", config.DatabaseConfig.BuildDbConnectionStr())
		if err != nil {
			logger.Log("database", config.DatabaseConfig.Database, "during", "connect", "err", err)
			os.Exit(1)
		}
		defer db.Close()
	}

	var (
		repository = company.NewRepository(db)
		service    = company.NewService(logger, repository)
		endpoints  = company.NewEndpointSet(service, logger)
		grpcServer = company.NewGRPCServer(endpoints, logger)
	)

	var g group.Group
	{
		port := fmt.Sprintf(":%d", config.Port)
		grpcListener, err := net.Listen("tcp", port)
		if err != nil {
			logger.Log("transport", "gRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "gRPC", "port", port)
			baseServer := grpc.NewServer()
			pb.RegisterCompanySvcServer(baseServer, grpcServer)
			reflection.Register(baseServer)
			return baseServer.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}
	{
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	logger.Log("exit", g.Run())
}
