package internalgrpc

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/pb"
	memorystorage "github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func getFreePort(t *testing.T) string {
	t.Helper()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer l.Close()
	return l.Addr().String()
}

func startGRPCServer(t *testing.T) (pb.EventServiceClient, func()) {
	t.Helper()

	grpcAddr := getFreePort(t)

	logg := logger.New("error")
	st := memorystorage.New()
	ap := app.New(logg, st)
	srv := New(ap, logg)

	go func() {
		if err := srv.Start(grpcAddr); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
		}
	}()

	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	client := pb.NewEventServiceClient(conn)

	cleanup := func() {
		srv.Stop()
		conn.Close()
	}

	return client, cleanup
}

func TestCreateAndListDayGRPC(t *testing.T) {
	client, cleanup := startGRPCServer(t)
	defer cleanup()

	ctx := context.Background()

	base := time.Date(2025, 7, 3, 12, 0, 0, 0, time.UTC)

	event := &pb.Event{
		Id:           "e1",
		Title:        "test grpc",
		StartTime:    timestamppb.New(base),
		Duration:     durationpb.New(time.Hour),
		UserId:       "u1",
		Description:  "desc",
		NotifyBefore: durationpb.New(30 * time.Minute),
	}

	// --- Create Event ---
	_, err := client.CreateEvent(ctx, &pb.CreateEventRequest{Event: event})
	require.NoError(t, err)

	// --- List Day ---
	resp, err := client.ListDay(ctx, &pb.ListDayRequest{
		UserId: "u1",
		Date:   timestamppb.New(base),
	})
	require.NoError(t, err)
	require.Len(t, resp.Events, 1)
	require.Equal(t, "e1", resp.Events[0].Id)
}
