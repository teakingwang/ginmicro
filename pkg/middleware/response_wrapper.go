package middleware

import (
	"context"
	"encoding/json"
	"github.com/teakingwang/ginmicro/pkg/errs"
	"net/http"

	"github.com/teakingwang/ginmicro/pkg/resp"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

// CustomResponseInterceptor 用于统一封装响应格式（code/message/data）
func CustomResponseInterceptor(ctx context.Context, w http.ResponseWriter, m proto.Message) error {
	// 如果头部存在 "X-Bypass-Wrap"，跳过封装（方便某些特殊接口）
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if val := md.Get("X-Bypass-Wrap"); len(val) > 0 {
			return nil // 使用默认的 response
		}
	}

	response := resp.HTTPResponse{
		Code:    errs.CodeSuccess,
		Message: errs.Message(errs.CodeSuccess),
		Data:    m,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	return json.NewEncoder(w).Encode(response)
}
