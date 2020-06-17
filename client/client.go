package main

import (
  "context"
  "fmt"
  "github.com/yizumi/app-mesh-example/client/proto/echo"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials"
  "log"
  "net/http"
  "os"
)

func loadCertificate() (credentials.TransportCredentials, error) {
  certFile := "ssl/ca.crt"
  return credentials.NewClientTLSFromFile(certFile, "")
}

func buildClientOpts() (grpc.DialOption, error) {
  tls := false
  if tls {
    creds, sslErr := loadCertificate()
    if sslErr != nil {
      log.Fatalf("Failed to load certificate: %v", sslErr)
      return nil, sslErr
    }
    return grpc.WithTransportCredentials(creds), nil
  } else {
    return grpc.WithInsecure(), nil
  }
}

func getEnv(name string, defaultValue string) string {
  value, exists := os.LookupEnv(name)
  if !exists {
    return defaultValue
  }
  return value
}

var echoClient echo.EchoServiceClient;

func main() {
  hostname := getEnv("ECHO_HOST", "localhost:50051")
  port := getEnv("PORT", "8080")

  log.Println("Echo Client 0.1")
  log.Println("Echo Client connecting to " + hostname)

  opts, err := buildClientOpts()
  if err != nil {
    log.Fatalf("Failed to build options: %v", err)
    return
  }

  conn, err := grpc.Dial(hostname, opts)
  if err != nil {
    log.Fatalf("Failed to connect to the server: %v", err)
  }

  defer conn.Close()

  echoClient = echo.NewEchoServiceClient(conn)

  log.Println("Listening to port: " + port)
  http.HandleFunc("/", handle)
  log.Fatal(http.ListenAndServe(":" + port, nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
  message := doEcho(r.FormValue("name"))
  fmt.Fprintf(w, "Response from the server: %s", message)
}

func doEcho(name string) string {
  for ;; {
    res, err := echoClient.Echo(context.Background(), &echo.EchoRequest { Name: name })
    if err == nil {
      fmt.Println("Message received:", res.Message)
      return res.Message
    } else {
      fmt.Println("Doh! Something went rong: %v", err)
      return "(Error)"
    }
  }
}

