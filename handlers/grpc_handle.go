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
	cfg := convertFromPb(req.Config)

	var problems analysis.Problems
	if err := problems.AnalyseCfg(cfg); err != nil {
		return nil, err
	}

	pbProblems := make([]*pb.Problem, len(problems))
	for i, problem := range problems {
		pbProblems[i] = &pb.Problem{
			Filename:       problem.Filename,
			Path:           problem.Path,
			Description:    problem.Description,
			Recommendation: problem.Recommendation,
			Severity:       problem.Severity,
		}
	}

	return &pb.AnalyzeResponse{Problems: pbProblems}, nil
}

func convertFromPb(pbCfg *pb.Config) *models.Config {
	cfg := &models.Config{}

	if pbCfg.Server != nil {
		cfg.Server.Host = pbCfg.Server.Host
		cfg.Server.TlsVerify = pbCfg.Server.TlsVerify
	}

	if pbCfg.Database != nil {
		cfg.Database.Password = pbCfg.Database.Password
	}

	if pbCfg.Storage != nil {
		cfg.Storage.Path = pbCfg.Storage.Path
		cfg.Storage.Permissions = pbCfg.Storage.Permissions
		cfg.Storage.DigestAlgorithm = pbCfg.Storage.DigestAlgorithm
	}

	if pbCfg.Log != nil {
		cfg.Log.Level = pbCfg.Log.Level
	}
	return cfg
}
