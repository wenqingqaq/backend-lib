package metrics

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
	"io"

	"context"
	kerr "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	ks "github.com/go-kratos/kratos/v2/transport/http/status"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/grpclog"
	nethttp "net/http"
)

var (
	BuildInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "build_info",
		Help: "Build information",
	}, []string{"version", "commit", "date"})

	Up = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "up",
		Help: "service up",
	}, []string{"service"})
)

func NewMetricsServer() *Server {
	addr := ":7070"
	path := "/metrics"

	srv := NewServer(
		Address(addr),
		Path(path),
		Gatherer(prometheus.DefaultGatherer),
	)
	log.Infof("[Metrics] server listening on: [::]%s", addr)

	return srv
}

func InitMetrics(name, version, gitCommit, BuildTs string) {
	BuildInfo.WithLabelValues(version, gitCommit, BuildTs).Set(1)
	Up.WithLabelValues(name).Set(1)
	prometheus.MustRegister(BuildInfo)
	prometheus.MustRegister(Up)
}

func ErrorHandler(_ context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
	var resp = new(kerr.Status)
	const fallback = `{"code": 500, "message": "failed to marshal error message"}`

	s := status.Convert(err)
	pb := s.Proto()

	details := pb.GetDetails()
	if len(details) > 0 {
		info := &errdetails.ErrorInfo{}
		err = details[0].UnmarshalTo(info)
		if err = details[0].UnmarshalTo(info); err == nil {
			resp.Reason = info.Reason
			resp.Metadata = info.Metadata
		}
	}

	resp.Message = pb.GetMessage()
	resp.Code = int32(ks.FromGRPCCode(s.Code()))

	contentType := marshaler.ContentType(pb)
	writer.Header().Set("Content-Type", contentType)

	buf, merr := marshaler.Marshal(resp)
	if merr != nil {
		grpclog.Infof("Failed to marshal error message %q: %v", s, merr)
		writer.WriteHeader(nethttp.StatusInternalServerError)
		if _, err = io.WriteString(writer, fallback); err != nil {
			grpclog.Infof("Failed to write response: %v", err)
		}
		return
	}
	st := ks.FromGRPCCode(s.Code())
	writer.WriteHeader(st)
	if _, err := writer.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}
}
