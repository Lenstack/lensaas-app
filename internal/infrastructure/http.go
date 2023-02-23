package infrastructure

import (
	"go.uber.org/zap"
	"net"
	"net/http"
)

type HttpServer struct {
	port string
}

func NewHttpServer(port string, handlers http.Handler, logger *zap.Logger) *HttpServer {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Sugar().Error(err)
		return nil
	}

	logger.Sugar().Info("Server is running on port: " + port)
	err = http.Serve(listen, handlers)
	if err != nil {
		logger.Sugar().Error(err)
		return nil
	}
	return &HttpServer{port}
}
