package e2bserver

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"

	cs "github.com/acompany-develop/QuickMPC-BTS/src/BeaverTripleService/ConfigStore"
	logger "github.com/acompany-develop/QuickMPC-BTS/src/BeaverTripleService/Log"
	tg "github.com/acompany-develop/QuickMPC-BTS/src/BeaverTripleService/TripleGenerator"
	pb "github.com/acompany-develop/QuickMPC-BTS/src/Proto/EngineToBts"
)

type server struct {
	pb.UnimplementedEngineToBtsServer
}

// モック時に置き換わる関数
var GetPartyIdFromIp = func(reqIpAddrAndPort string) (uint32, error) {
	arr := strings.Split(reqIpAddrAndPort, ":")
	if len(arr) != 2 {
		errText := fmt.Sprintf("requestのIpAddessの形式が異常: %s", reqIpAddrAndPort)
		logger.Error(errText)
		return 0, fmt.Errorf(errText)
	}
	reqIpAddr, _ := arr[0], arr[1]

	var partyId uint32
	for _, party := range cs.Conf.RequestPartyList {
		if reqIpAddr == party.IpAddress {
			partyId = party.PartyId
			break
		}
	}
	if partyId == 0 {
		errText := fmt.Sprintf("PartyList[%s, %s, %s]に存在しないIPからのリクエスト: %s", cs.Conf.RequestPartyList[0].IpAddress, cs.Conf.RequestPartyList[1].IpAddress, cs.Conf.RequestPartyList[2].IpAddress, reqIpAddr)
		logger.Error(errText)
		return 0, fmt.Errorf(errText)
	}

	return partyId, nil
}

func (s *server) GetTriples(ctx context.Context, in *pb.GetTriplesRequest) (*pb.GetTriplesResponse, error) {
	p, _ := peer.FromContext(ctx)
	reqIpAddrAndPort := p.Addr.String()
	partyId, err := GetPartyIdFromIp(reqIpAddrAndPort)
	if err != nil {
		return nil, err
	}
	logger.Infof("Ip %s, jobId: %d, partyId: %d\n", reqIpAddrAndPort, in.GetJobId(), partyId)
	triples, err := tg.GetTriples(in.GetJobId(), partyId, in.GetAmount())
	if err != nil {
		return nil, err
	}

	return &pb.GetTriplesResponse{
		Triples: triples,
	}, nil
}

// requestを受け取った際の共通処理
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 定期実行されるhealth checkでlogは必要ないため即時return
	if info.FullMethod == "/grpc.health.v1.Health/Check" {
		return handler(ctx, req)
	}

	logger.Infof("received: %s", info.FullMethod)
	// 処理を実行する
	res, err := handler(ctx, req)

	// エラー時にログとしてrequest，responseを出力する
	if err != nil {
		logger.Errorf("request: {%v}\tresponse: {%v}\n", req, res)
	}

	logger.Infof("send: %s", info.FullMethod)
	return res, err
}

func RunServer() {
	listenIp := fmt.Sprintf("0.0.0.0:%d", cs.Conf.Port)
	lis, err := net.Listen("tcp", listenIp)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	enforcementPolicyMinTime := 5
	s := grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             (time.Duration(enforcementPolicyMinTime) * time.Second),
				PermitWithoutStream: true,
			},
		),
		grpc.UnaryInterceptor(unaryInterceptor),
	)

	pb.RegisterEngineToBtsServer(s, &server{})
	grpcHealthServer := health.NewServer()
	healthpb.RegisterHealthServer(s, grpcHealthServer)
	reflection.Register(s)
	logger.Info("a2dbg Server listening on: ", listenIp)
	if err := s.Serve(lis); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}
