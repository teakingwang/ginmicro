package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/teakingwang/ginmicro/pkg/errs"
	"github.com/teakingwang/ginmicro/pkg/resp"
	"google.golang.org/grpc/status"
)

// CustomErrorInterceptor 用于封装 gRPC-Gateway 错误响应
func CustomErrorInterceptor(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	// 将 gRPC error 转为状态码和 message
	s := status.Convert(err)
	code := s.Code()
	msg := s.Message()

	// 可以从 s.Details 里获取更详细的信息

	// 可以扩展为自定义 code 映射
	appCode := errs.CodeServerError
	if code.String() == "NotFound" {
		appCode = errs.CodeNotFound
	}

	// 封装响应
	res := resp.HTTPResponse{
		Code:    appCode,
		Message: msg,
		Data:    nil,
	}

	w.Header().Set("Content-Type", marshaler.ContentType(nil))
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
