package grpc

import chatv1 "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/chat_v1"

type ChatImplementation struct {
	chatv1.UnimplementedChatServiceServer
	chatService ChatService
}

type ChatService interface{}

func NewGrpcChatImplementation(chatService ChatService) *ChatImplementation {
	return &ChatImplementation{
		chatService: chatService,
	}
}

// func (impl *AccessImplementation) Check(ctx context.Context, request *accessv1.CheckRequest) (*emptypb.Empty, error) {

// 	payload, ok := ctx.Value(interceptor.ContextKeyUserClaims).(jwt.Claims)

// 	if !ok {
// 		return nil, status.Errorf(codes.Internal, "missing required token")
// 	}

// 	audience, err := payload.GetAudience()

// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "extract role error")
// 	}

// 	role, err := strconv.Atoi(audience[0])

// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "extract role error")
// 	}

// 	access := impl.authService.Check(ctx, request.GetEndpointAddress(), domain.UserRole(role))

// 	if !access {
// 		return &emptypb.Empty{}, status.Errorf(codes.PermissionDenied, "access denied")
// 	}

// 	return &emptypb.Empty{}, nil
// }
