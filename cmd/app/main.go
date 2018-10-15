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

	pb "github.com/nathanows/elegant-monolith/_protos"
	"github.com/nathanows/elegant-monolith/internal/company"
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

	var db *sqlx.DB
	{
		var err error
		connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable ", "USER", "DBNAME")
		db, err = sqlx.Connect("postgres", connStr)
		if err != nil {
			logger.Log("database", "elegant_monolith", "during", "connect", "err", err)
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
		grpcListener, err := net.Listen("tcp", ":8080")
		if err != nil {
			logger.Log("transport", "gRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "gRPC", "addr", ":8080")
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
