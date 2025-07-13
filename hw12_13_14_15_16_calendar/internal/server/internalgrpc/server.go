package internalgrpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/pb"
	"github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ---- adapter interfaces ---------------------------------------------------

type Logger interface {
	Info(string)
	Error(string)
}

type Application interface {
	CreateFullEvent(ctx context.Context, e storage.Event) error
	UpdateEvent(ctx context.Context, e storage.Event) error
	DeleteEvent(ctx context.Context, id string) error

	ListDay(ctx context.Context, userID string, date time.Time) ([]storage.Event, error)
	ListWeek(ctx context.Context, userID string, weekStart time.Time) ([]storage.Event, error)
	ListMonth(ctx context.Context, userID string, monthStart time.Time) ([]storage.Event, error)
}

// ---- server ---------------------------------------------------------------

type Server struct {
	pb.UnimplementedEventServiceServer
	app    Application
	logger Logger
	srv    *grpc.Server
}

func New(app Application, logger Logger) *Server {
	unary := grpc.ChainUnaryInterceptor(loggingInterceptor(logger))
	s := &Server{
		app:    app,
		logger: logger,
		srv:    grpc.NewServer(unary),
	}
	pb.RegisterEventServiceServer(s.srv, s)
	reflection.Register(s.srv)
	return s
}

func (s *Server) Start(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("grpc listen: %w", err)
	}
	go func() {
		if err := s.srv.Serve(ln); err != nil {
			s.logger.Error("grpc serve: " + err.Error())
		}
	}()
	return nil
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}

// ---- service implementation ----------------------------------------------

func (s *Server) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.EventResponse, error) {
	if req == nil || req.Event == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	e := fromProto(req.Event)
	if err := s.app.CreateFullEvent(ctx, e); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.EventResponse{Event: req.Event}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.EventResponse, error) {
	if req == nil || req.Event == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	e := fromProto(req.Event)
	if err := s.app.UpdateEvent(ctx, e); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.EventResponse{Event: req.Event}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*emptypb.Empty, error) {
	if req == nil || req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id required")
	}
	if err := s.app.DeleteEvent(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) ListDay(ctx context.Context, req *pb.ListDayRequest) (*pb.EventsResponse, error) {
	t := req.GetDate().AsTime()
	evs, err := s.app.ListDay(ctx, req.GetUserId(), t)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.EventsResponse{Events: toProto(evs)}, nil
}

func (s *Server) ListWeek(ctx context.Context, req *pb.ListWeekRequest) (*pb.EventsResponse, error) {
	evs, err := s.app.ListWeek(ctx, req.GetUserId(), req.GetWeekStart().AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.EventsResponse{Events: toProto(evs)}, nil
}

func (s *Server) ListMonth(ctx context.Context, req *pb.ListMonthRequest) (*pb.EventsResponse, error) {
	evs, err := s.app.ListMonth(ctx, req.GetUserId(), req.GetMonthStart().AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.EventsResponse{Events: toProto(evs)}, nil
}

// ---- helpers --------------------------------------------------------------

func fromProto(p *pb.Event) storage.Event {
	return storage.Event{
		ID:           p.Id,
		Title:        p.Title,
		StartTime:    p.StartTime.AsTime(),
		Duration:     p.Duration.AsDuration(),
		Description:  p.Description,
		UserID:       p.UserId,
		NotifyBefore: p.NotifyBefore.AsDuration(),
	}
}

func toProto(src []storage.Event) []*pb.Event {
	out := make([]*pb.Event, 0, len(src))
	for _, e := range src {
		out = append(out, &pb.Event{
			Id:           e.ID,
			Title:        e.Title,
			StartTime:    timestamppb.New(e.StartTime),
			Duration:     durationpb.New(e.Duration),
			Description:  e.Description,
			UserId:       e.UserID,
			NotifyBefore: durationpb.New(e.NotifyBefore),
		})
	}
	return out
}

// ---- interceptors ---------------------------------------------------------

func loggingInterceptor(log Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		mtd := info.FullMethod

		ua := ""
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if v := md.Get("user-agent"); len(v) > 0 {
				ua = v[0]
			}
		}

		resp, err := handler(ctx, req)
		latency := time.Since(start)

		statusCode := codes.OK
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
		}

		log.Info(fmt.Sprintf("gRPC %s %s %d %dms \"%s\"", mtd, start.Format("02/Jan/2006:15:04:05 -0700"),
			statusCode, latency.Milliseconds(), ua))
		return resp, err
	}
}
