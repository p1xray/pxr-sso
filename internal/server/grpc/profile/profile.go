package profileserver

import (
	"context"
	ssoprofilepb "github.com/p1xray/pxr-sso-protos/gen/go/profile"
	"github.com/p1xray/pxr-sso/internal/server"
	"google.golang.org/grpc"
)

const (
	emptyValue = 0
)

type serverAPI struct {
	ssoprofilepb.UnimplementedSsoProfileServer
	// profile server.ProfileService TODO: implement profile use case
}

// Register registers the implementation of the API service with the gRPC server.
func Register(gRPC *grpc.Server) {
	ssoprofilepb.RegisterSsoProfileServer(gRPC, &serverAPI{})
}

func (s *serverAPI) GetProfile(
	ctx context.Context,
	req *ssoprofilepb.GetProfileRequest,
) (*ssoprofilepb.GetProfileResponse, error) {
	if err := validateGetProfileRequest(req); err != nil {
		return nil, err
	}

	//userProfile, err := s.profile.UserProfile(ctx, req.GetUserId())
	//if err != nil {
	//	if errors.Is(err, service.ErrUserNotFound) {
	//		return nil, server.NotFoundError("user not found")
	//	}
	//
	//	return nil, server.InternalError("failed to get user profile")
	//}
	//
	//var dateOfBirthPb *timestamppb.Timestamp
	//if userProfile.DateOfBirth != nil {
	//	dateOfBirthPb = timestamppb.New(*userProfile.DateOfBirth)
	//}
	//
	//genderPb := ssoprofilepb.Gender_GENDER_UNSPECIFIED
	//if userProfile.Gender != nil {
	//	genderPb = ssoprofilepb.Gender(*userProfile.Gender)
	//}
	//
	//var avatarFileKeyPb *wrappers.StringValue
	//if userProfile.AvatarFileKey != nil {
	//	avatarFileKeyPb = &wrappers.StringValue{Value: *userProfile.AvatarFileKey}
	//}
	//
	//return &ssoprofilepb.GetProfileResponse{
	//	UserId:        userProfile.UserID,
	//	Username:      userProfile.Username,
	//	Fio:           userProfile.FIO,
	//	DateOfBirth:   dateOfBirthPb,
	//	Gender:        genderPb,
	//	AvatarFileKey: avatarFileKeyPb,
	//}, nil

	return &ssoprofilepb.GetProfileResponse{}, nil
}

func validateGetProfileRequest(req *ssoprofilepb.GetProfileRequest) error {
	if req.GetUserId() == emptyValue {
		return server.InvalidArgumentError("user id is empty")
	}

	return nil
}
