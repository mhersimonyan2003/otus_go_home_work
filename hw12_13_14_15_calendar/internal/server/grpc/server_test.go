package grpc

import (
	"context"
	"net"
	"testing"
	"time"

	pb "github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/api"
	"github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	storage := memorystorage.New()
	logger := logger.New(logger.Debug)
	s := grpc.NewServer()
	pb.RegisterCalendarServiceServer(s, NewGRPCServer(app.New(storage), logger))
	go func() {
		if err := s.Serve(lis); err != nil {
			logger.Error("Server exited with error: " + err.Error())
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestGRPCServer(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewCalendarServiceClient(conn)

	StartTime := timestamppb.New(time.Now())
	EndTime := timestamppb.New(time.Now().Add(time.Hour))

	tests := []struct {
		name      string
		method    string
		req       interface{}
		wantErr   bool
		wantLen   int
		wantEvent *pb.Event
	}{
		{
			name:    "AddEventSuccess",
			method:  "AddEvent",
			req:     &pb.Event{Id: "1", Title: "Test Event", StartTime: StartTime, EndTime: EndTime, Details: "Details"},
			wantErr: false,
		},
		{
			name:   "AddEventFailure",
			method: "AddEvent",
			req: &pb.Event{
				Id: "1", Title: "Test Event With the same Id", StartTime: StartTime,
				EndTime: EndTime, Details: "Details",
			},
			wantErr: true,
		},
		{
			name:   "UpdateEventSuccess",
			method: "UpdateEvent",
			req: &pb.Event{
				Id: "1", Title: "Updated Event", StartTime: StartTime,
				EndTime: EndTime, Details: "Updated Details",
			},
			wantErr: false,
		},
		{
			name:    "UpdateEventFailure",
			method:  "UpdateEvent",
			req:     &pb.Event{Id: "", Title: "", StartTime: nil, EndTime: nil, Details: ""},
			wantErr: true,
		},
		{
			name:    "GetEventSuccess",
			method:  "GetEvent",
			req:     &pb.GetEventRequest{Id: "1"},
			wantErr: false,
			wantEvent: &pb.Event{
				Id: "1", Title: "Updated Event", StartTime: StartTime,
				EndTime: EndTime, Details: "Updated Details",
			},
		},
		{
			name:    "DeleteEventSuccess",
			method:  "DeleteEvent",
			req:     &pb.DeleteEventRequest{Id: "1"},
			wantErr: false,
		},
		{
			name:    "DeleteEventFailure",
			method:  "DeleteEvent",
			req:     &pb.DeleteEventRequest{Id: ""},
			wantErr: true,
		},
		{
			name:   "ListEventsSuccess",
			method: "ListEvents",
			req: &pb.ListEventsRequest{
				Start: timestamppb.New(time.Now().Add(-time.Hour)),
				End:   timestamppb.New(time.Now().Add(time.Hour)),
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name:      "GetEventFailure",
			method:    "GetEvent",
			req:       &pb.GetEventRequest{Id: ""},
			wantErr:   true,
			wantEvent: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.method {
			case "AddEvent":
				_, err := client.AddEvent(ctx, tt.req.(*pb.Event))
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			case "UpdateEvent":
				_, err := client.UpdateEvent(ctx, tt.req.(*pb.Event))
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			case "DeleteEvent":
				_, err := client.DeleteEvent(ctx, tt.req.(*pb.DeleteEventRequest))
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			case "ListEvents":
				res, err := client.ListEvents(ctx, tt.req.(*pb.ListEventsRequest))
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Len(t, res.Events, tt.wantLen)
				}
			case "GetEvent":
				res, err := client.GetEvent(ctx, tt.req.(*pb.GetEventRequest))
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.True(t, proto.Equal(tt.wantEvent, res), "expected and actual events are not equal")
				}
			}
		})
	}
}
