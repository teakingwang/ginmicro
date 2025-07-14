package middleware

import (
	"encoding/json"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/teakingwang/ginmicro/pkg/errs"
	"github.com/teakingwang/ginmicro/pkg/resp"
	"google.golang.org/protobuf/proto"
)

type wrappedJSONMarshaler struct {
	runtime.JSONPb
}

func NewWrappedMarshaler() *wrappedJSONMarshaler {
	return &wrappedJSONMarshaler{JSONPb: runtime.JSONPb{}}
}

func (m *wrappedJSONMarshaler) Marshal(v interface{}) ([]byte, error) {
	// 如果是 proto.Message，封装成标准结构
	if msg, ok := v.(proto.Message); ok {
		// 原始 data
		rawData, err := m.JSONPb.Marshal(msg)
		if err != nil {
			return nil, err
		}
		var data interface{}
		_ = json.Unmarshal(rawData, &data) // 反序列化一次，变成 map

		wrapped := resp.HTTPResponse{
			Code:    errs.CodeSuccess,
			Message: errs.Message(errs.CodeSuccess),
			Data:    data,
		}
		return json.Marshal(wrapped)
	}

	// 不是 proto，走默认
	return m.JSONPb.Marshal(v)
}

func (m *wrappedJSONMarshaler) ContentType(_ interface{}) string {
	return "application/json"
}
