package handlers

import (
	"context"

	"github.com/VysMax/analyze-cfg/analysis"
	pb "github.com/VysMax/analyze-cfg/gen/proto"
	"github.com/VysMax/analyze-cfg/models"
)

type Server struct {
	pb.UnimplementedCfgAnalyzerServer
}

func (s *Server) Analyze(ctx context.Context, req *pb.AnalyzeRequest) (*pb.AnalyzeResponse, error) {
	cfg := convertFromPb(req)

	var problems analysis.Problems
	if err := problems.AnalyzeCfg(cfg); err != nil {
		return nil, err
	}

	pbProblems := make([]*pb.Problem, len(problems))
	for i, problem := range problems {
		pbProblems[i] = &pb.Problem{
			Filename:       "from GRPC",
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
		cfg.Server.TlsVerify = req.Server.TlsVerify
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
