package main

import (
  "context"
  "github.com/raksul/app-mesh-example/server/proto/echo"
  "google.golang.org/grpc"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/credentials"
  "google.golang.org/grpc/reflection"
  "google.golang.org/grpc/status"
  health "google.golang.org/grpc/health/grpc_health_v1"
  "log"
  "net"
  "os"
)

type server struct{}

const version = "0.9"

func (s *server) Check(ctx context.Context, in *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
  log.Printf("Received Check request: %v", in)
  return &health.HealthCheckResponse{Status: health.HealthCheckResponse_SERVING}, nil
}

func (s *server) Watch(in *health.HealthCheckRequest, _ health.Health_WatchServer) error {
  log.Printf("Received Watch request: %v", in)
  return status.Error(codes.Unimplemented, "unimplemented")
}

func (*server) Echo(ctx context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {
  log.Printf("Echo was called:")
  message := "Hello, " + req.GetName() + "-san! (Said " + ipAddress() + ", Version " + version + ")"
  res := &echo.EchoResponse { Message: message }
  return res, nil
}

func ipAddress() string {
  addrs, err := net.InterfaceAddrs()
  if err != nil {
    return "UNKNOWN"
  }
  for _, addr := range addrs {
    if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
      if ipnet.IP.To4() != nil {
        return ipnet.IP.String()
      }
    }
  }
  return addrs[0].String()
}

func startServer(p string, ssl bool) *grpc.Server {
  log.Println("Echo Server Version " + version)
  log.Printf("Starting on port: %v, ssl: %v\n", p, ssl)
  var opts []grpc.ServerOption
  if ssl {
    certFile := "ssl/server.crt"
    keyFile := "ssl/server.pem"
    creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
    if sslErr != nil {
      log.Fatalf("Error while loading TLS certificate: %v", sslErr)
      return nil
    }
    opts = append(opts, grpc.Creds(creds))
  }
  return grpc.NewServer(opts...)
}

func getEnv(name string, defaultValue string) string {
  value, exists := os.LookupEnv(name)
  if !exists {
    return defaultValue
  }
  return value
}

func main() {
  port := getEnv("PORT", "50051")
  lis, err := net.Listen("tcp", ":" + port)
  if err != nil {
    log.Fatalf("Failed to listen: %v", err)
  }

  s := startServer(port, false)

  cs := server{}

  echo.RegisterEchoServiceServer(s, &cs)
  health.RegisterHealthServer(s, &cs)
  reflection.Register(s)

  e := s.Serve(lis)
  if e != nil {
    log.Fatalf("Failed to start server: %v", e)
  }
}
