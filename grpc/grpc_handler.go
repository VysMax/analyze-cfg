package grpcserver

import (
	"context"

	pb "github.com/VysMax/analyze-cfg/gen/proto"
	"github.com/VysMax/analyze-cfg/models"
	"github.com/VysMax/analyze-cfg/usecase"
)

type Server struct {
	pb.UnimplementedCfgAnalyzerServer
	service *usecase.Service
}

func NewHandler(s *usecase.Service) *Server {
	return &Server{service: s}
}

func (s *Server) Analyze(ctx context.Context, req *pb.AnalyzeRequest) (*pb.AnalyzeResponse, error) {
	cfg := convertFromPb(req)

	problems, err := s.service.AnalyzeCfg(cfg)
	if err != nil {
		return nil, err
	}

	pbProblems := make([]*pb.Problem, len(problems))
	for i, problem := range problems {
		pbProblems[i] = &pb.Problem{
			Path:           problem.Path,
			Description:    problem.Description,
			Recommendation: problem.Recommendation,
			Severity:       problem.Severity,
		}
	}

	return &pb.AnalyzeResponse{Problems: pbProblems}, nil
}

func convertFromPb(req *pb.AnalyzeRequest) *models.Config {
	cfg := &models.Config{}

	if req.Server != nil {
		cfg.Server.Host = req.Server.Host
		if req.Server.TlsVerify != nil {
			cfg.Server.TlsVerify = &req.Server.TlsVerify.Value
		}

	}

	if req.Database != nil {
		cfg.Database.Password = req.Database.Password
	}

	if req.Storage != nil {
		cfg.Storage.Path = req.Storage.Path
		cfg.Storage.Permissions = req.Storage.Permissions
		cfg.Storage.DigestAlgorithm = req.Storage.DigestAlgorithm
	}

	if req.Log != nil {
		cfg.Log.Level = req.Log.Level
	}
	return cfg
}
