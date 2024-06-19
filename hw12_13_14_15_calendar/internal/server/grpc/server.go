package grpc

import (
	"context"
	"fmt"
	"net"

	pb "github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/api"
	"github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedCalendarServiceServer
	app    *app.App
	logger logger.Logger
}

func NewGRPCServer(app *app.App, logger logger.Logger) *Server {
	return &Server{
		app:    app,
		logger: logger,
	}
}

func (s *Server) AddEvent(ctx context.Context, req *pb.Event) (*emptypb.Empty, error) {
	event := storage.Event{
		ID:        req.Id,
		Title:     req.Title,
		StartTime: req.StartTime.AsTime(),
		EndTime:   req.EndTime.AsTime(),
		Details:   req.Details,
	}

	if err := s.app.AddEvent(event); err != nil {
		s.logger.Error("Failed to add event: " + err.Error())
		return nil, err
	}

	s.logger.Info("Event added: " + req.Id)
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, req *pb.Event) (*emptypb.Empty, error) {
	event := storage.Event{
		ID:        req.Id,
		Title:     req.Title,
		StartTime: req.StartTime.AsTime(),
		EndTime:   req.EndTime.AsTime(),
		Details:   req.Details,
	}

	if err := s.app.UpdateEvent(event); err != nil {
		s.logger.Error("Failed to update event: " + err.Error())
		return nil, err
	}

	s.logger.Info("Event updated: " + req.Id)
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*emptypb.Empty, error) {
	if err := s.app.DeleteEvent(req.Id); err != nil {
		s.logger.Error("Failed to delete event: " + err.Error())
		return nil, err
	}

	s.logger.Info("Event deleted: " + req.Id)
	return &emptypb.Empty{}, nil
}

func (s *Server) ListEvents(ctx context.Context, req *pb.ListEventsRequest) (*pb.ListEventsResponse, error) {
	start := req.Start.AsTime()
	end := req.End.AsTime()

	events, err := s.app.ListEvents(start, end)
	if err != nil {
		s.logger.Error("Failed to list events: " + err.Error())
		return nil, err
	}

	pbEvents := make([]*pb.Event, len(events))
	for i, event := range events {
		pbEvents[i] = &pb.Event{
			Id:        event.ID,
			Title:     event.Title,
			StartTime: timestamppb.New(event.StartTime),
			EndTime:   timestamppb.New(event.EndTime),
			Details:   event.Details,
		}
	}

	s.logger.Info("Events listed")
	return &pb.ListEventsResponse{Events: pbEvents}, nil
}

func (s *Server) GetEvent(ctx context.Context, req *pb.GetEventRequest) (*pb.Event, error) {
	event, err := s.app.GetEventByID(req.Id)
	if err != nil {
		s.logger.Error("Failed to get event: " + err.Error())
		return nil, err
	}

	s.logger.Info("Event retrieved: " + req.Id)
	return &pb.Event{
		Id:        event.ID,
		Title:     event.Title,
		StartTime: timestamppb.New(event.StartTime),
		EndTime:   timestamppb.New(event.EndTime),
		Details:   event.Details,
	}, nil
}

func RunGRPCServer(logger logger.Logger, app *app.App, host string, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		logger.Error("Failed to listen: " + err.Error())
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCalendarServiceServer(grpcServer, NewGRPCServer(app, logger))

	logger.Info("Starting GRPC server on " + lis.Addr().String())
	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("Failed to serve: " + err.Error())
		return err
	}

	return nil
}
