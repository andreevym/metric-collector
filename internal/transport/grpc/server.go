package grpc

import (
	"context"
	"fmt"
	"github.com/andreevym/metric-collector/internal/controller"
	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/storage/store"
	"github.com/andreevym/metric-collector/internal/transport/grpc/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
)

type Server struct {
	proto.UnimplementedMetricCollectorServer

	grpcServer    *grpc.Server
	metricStorage store.Storage
	dbClient      store.Client
	controller    controller.Controller
	secretKey     string
	cryptoKey     string
	trustedSubnet string
	address       string
}

func NewGrpcServer(
	dbClient store.Client,
	metricStorage store.Storage,
	secretKey string,
	cryptoKey string,
	trustedSubnet string,
	address string,
) *Server {
	controller := controller.NewController(metricStorage, dbClient)
	return &Server{
		metricStorage: metricStorage,
		dbClient:      dbClient,
		secretKey:     secretKey,
		cryptoKey:     cryptoKey,
		trustedSubnet: trustedSubnet,
		address:       address,
		controller:    controller,
	}
}

func (s Server) Run() error {
	listen, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("run grpc server: %w", err)
	}
	s.grpcServer = grpc.NewServer()
	proto.RegisterMetricCollectorServer(s.grpcServer, &Server{})
	logger.Logger().Info("listening grpc server", zap.String("address", s.address))
	if err := s.grpcServer.Serve(listen); err != nil {
		return fmt.Errorf("start grpc server: %w", err)
	}
	return nil
}

func (s Server) Shutdown() error {
	logger.Logger().Info("shutting down grpc server")
	s.grpcServer.GracefulStop()
	return nil
}

func (s Server) Ping(context.Context, *proto.PingRequest) (*proto.PingResponse, error) {
	err := s.controller.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping error: %w", err)
	}

	return &proto.PingResponse{}, nil
}
func (s Server) Updates(ctx context.Context, updatesRequest *proto.UpdatesRequest) (*proto.UpdatesResponse, error) {
	metrics := make([]*store.Metric, 0, len(updatesRequest.Metrics))
	for _, metric := range updatesRequest.Metrics {
		metrics = append(metrics, &store.Metric{
			ID:    metric.Id,
			MType: metric.Type,
			Delta: &metric.Delta,
			Value: &metric.Value,
		})
	}
	err := s.controller.Updates(ctx, metrics)
	if err != nil {
		return nil, fmt.Errorf("update metrics error: %w", err)
	}
	return &proto.UpdatesResponse{}, nil
}

func (s Server) Update(ctx context.Context, r *proto.UpdateRequest) (*proto.UpdateResponse, error) {
	m := &store.Metric{
		ID:    r.Id,
		MType: r.Type,
		Value: &r.Value,
		Delta: &r.Delta,
	}

	respMetric, err := s.controller.Update(ctx, m)
	if err != nil {
		return nil, fmt.Errorf("failed to update metrics: %w", err)
	}

	updateResponse := &proto.UpdateResponse{
		Id:   respMetric.ID,
		Type: respMetric.MType,
	}
	if respMetric.Delta != nil {
		updateResponse.Delta = *respMetric.Delta
	}
	if respMetric.Value != nil {
		updateResponse.Value = *respMetric.Value
	}
	return updateResponse, nil
}

func (s Server) Value(ctx context.Context, r *proto.ValueRequest) (*proto.ValueResponse, error) {
	metric := s.controller.Value(ctx, r.Id, r.MetricType)
	if metric == nil {
		return nil, status.Error(codes.NotFound, "metric not found")
	}

	m := &proto.Metric{
		Id:   metric.ID,
		Type: metric.MType,
	}
	if metric.Delta != nil {
		m.Delta = *metric.Delta
	}
	if metric.Value != nil {
		m.Value = *metric.Value
	}
	return &proto.ValueResponse{
		Metric: m,
	}, nil
}
