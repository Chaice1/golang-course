package collectorapp

import (
	"log"
	"net"

	collectorusecase "github.com/Chaice1/golang-course/task2/collector/internal/usecase"

	collectorpb "github.com/Chaice1/golang-course/task2/gen"

	collectorclient "github.com/Chaice1/golang-course/task2/collector/internal/adapter/client"
	collectorhandler "github.com/Chaice1/golang-course/task2/collector/internal/controller"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunApp(config Config) {

	server := grpc.NewServer()
	GithubClient := collectorclient.GitHubApiClient{}
	CollectorUsecase := collectorusecase.NewCollectorService(&GithubClient)
	CollectorHandler := collectorhandler.NewHandler(CollectorUsecase)

	collectorpb.RegisterCollectorServer(server, CollectorHandler)

	reflection.Register(server)

	lis, err := net.Listen("tcp", config.GRPCport)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}

}
