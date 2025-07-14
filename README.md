# ginmicro
microservice with gin
# generate proto
protoc -I. --go_out=. --go-grpc_out=. --grpc-gateway_out=. api/user/user.proto
protoc -I. --go_out=. --go-grpc_out=. --grpc-gateway_out=. api/order/order.proto
# Makefile
make proto	自动生成 user.pb.go / user_grpc.pb.go
make build	编译所有服务（二进制输出到 bin/）
make clean	清除 bin/
make run-user	运行 user-service 服务
# grpcurl
grpcurl -v -plaintext \                     
-proto api/order/order.proto \
-import-path api \
-import-path third_party \
-d '{"id": c}' \
localhost:50052 \
order.OrderService.GetOrder

grpcurl -v -plaintext -proto api/order/order.proto -import-path api -import-path third_party -d '{"id": 100000000000000001}' localhost:50052 order.OrderService.GetOrder
