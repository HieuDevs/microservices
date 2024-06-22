package api

import (
	"context"
	"fmt"
	"log"
	"logger-service/src/data"
	"logger-service/src/logs"
	"net"

	"google.golang.org/grpc"
)

type LogService struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogService) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	log := req.GetLogEntry()

	logEntry := data.LogEntry{
		Name: log.Name,
		Data: log.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		return &logs.LogResponse{
			Result: "failed!",
		}, err
	}

	return &logs.LogResponse{
		Result: "logged!",
	}, nil
}

func (app *Config) GRPCLisener() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", app.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen gRPC: %v", err)
	}

	s := grpc.NewServer()
	logs.RegisterLogServiceServer(s, &LogService{Models: app.Models})
	log.Printf("gRPC server is running on port: %s", app.GRPCPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
