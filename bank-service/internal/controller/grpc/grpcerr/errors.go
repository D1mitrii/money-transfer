package grpcerr

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrParseUUID       = status.Error(codes.InvalidArgument, "incorrect format of acountUUID")
	ErrIncorrectAmount = status.Error(codes.InvalidArgument, "incorrect ammount")
	ErrAccountNotFound = status.Error(codes.NotFound, "account not found")
	ErrServiceLayer    = status.Error(codes.Internal, "service layer error")
)
