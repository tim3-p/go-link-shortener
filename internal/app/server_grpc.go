package app

import (
	"context"
	"errors"
	"net"

	"github.com/jackc/pgx/v4"
	"github.com/tim3-p/go-link-shortener/internal/configs"
	"github.com/tim3-p/go-link-shortener/internal/models"
	"github.com/tim3-p/go-link-shortener/internal/pkg"
	"github.com/tim3-p/go-link-shortener/internal/storage"
	pb "github.com/tim3-p/go-link-shortener/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Application repo
type GrpcAppHandler struct {
	pb.UnimplementedShortenerServer
	storage storage.Repository
}

// App Handler constructor
func NewGrpcAppHandler(s storage.Repository) *GrpcAppHandler {
	return &GrpcAppHandler{storage: s}
}

func (h *GrpcAppHandler) DBPing(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	conn, err := pgx.Connect(context.Background(), configs.EnvConfig.DatabaseDSN)
	if err != nil {
		return &pb.Empty{}, err
	}
	defer conn.Close(context.Background())
	return &pb.Empty{}, nil
}

func (h *GrpcAppHandler) StatsHandler(ctx context.Context, in *pb.Empty) (*pb.StatsResponse, error) {

	if configs.EnvConfig.TrustedSubnet == "" {
		return nil, status.Errorf(codes.PermissionDenied, "You don't have access to this handler")
	}

	_, IPNet, err := net.ParseCIDR(configs.EnvConfig.TrustedSubnet)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "You don't have access to this handler")
	}

	if !IPNet.Contains(net.ParseIP(ctx.Value("X-Real-IP").(string))) {
		return nil, status.Errorf(codes.PermissionDenied, "You subnet don't have access to this handler")
	}

	urls, users, err := h.storage.GetStats()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error %v", err)
	}

	return &pb.StatsResponse{
		UrlsCount:  int32(urls),
		UsersCount: int32(users),
	}, nil
}

func (h *GrpcAppHandler) GetHandler(ctx context.Context, in *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	resp, err := h.storage.Get(in.Value, userIDVar)
	if err != nil {

		if errors.Is(err, storage.ErrURLDeleted) {
			return nil, status.Errorf(codes.Internal, "URL deleted")
		}

		return nil, status.Errorf(codes.NotFound, "ID not found")
	}
	return &pb.SimpleResponse{
		Value: resp,
	}, nil
}

func (h *GrpcAppHandler) PostHandler(ctx context.Context, in *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	urlHash := pkg.HashURL([]byte(in.Value))

	err := h.storage.Add(urlHash, in.Value, userIDVar)
	_, err = pkg.CheckDBError(err)

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.SimpleResponse{
		Value: configs.EnvConfig.BaseURL + "/" + urlHash,
	}, nil
}

func (h *GrpcAppHandler) UserUrls(ctx context.Context, in *pb.Empty) (*pb.UserUrlsResponse, error) {

	mapRes, err := h.storage.GetUserURLs(userIDVar)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if len(mapRes) == 0 {
		return nil, status.Error(codes.NotFound, "Not found")
	}

	res := &pb.UserUrlsResponse{}

	for key, element := range mapRes {

		res.URLs = append(res.URLs, &pb.UserUrlsResponse_URL{
			ShortURL:  configs.EnvConfig.BaseURL + "/" + key,
			OriginURL: element,
		})
	}

	return res, nil
}

func (h *GrpcAppHandler) ShortenBatchHandler(ctx context.Context, in *pb.ShortenBatchRequest) (*pb.ShortenBatchResponse, error) {
	var req []models.ShortenBatchRequest
	resp := &pb.ShortenBatchResponse{}

	for _, v := range in.BatchURL {
		req = append(req, models.ShortenBatchRequest{
			CorrelationID: v.CorrelationID,
			OriginalURL:   v.OriginalURL,
		})
	}

	for _, value := range req {
		urlHash := pkg.HashURL([]byte(value.OriginalURL))
		h.storage.Add(urlHash, string(value.OriginalURL), userIDVar)

		resp.BatchURL = append(resp.BatchURL, &pb.ShortenBatchResponse_URL{
			CorrelationID: value.CorrelationID,
			ShortURL:      configs.EnvConfig.BaseURL + "/" + urlHash,
		})
	}

	return resp, nil
}

func (h *GrpcAppHandler) DeleteBatchHandler(ctx context.Context, in *pb.DeleteBatchHandlerRequest) (*pb.Empty, error) {

	TChan <- &models.Task{
		URLs:   in.ID,
		UserID: userIDVar,
	}

	return &pb.Empty{}, nil
}
