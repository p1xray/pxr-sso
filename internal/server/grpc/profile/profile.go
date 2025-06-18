package profileserver

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/wrappers"
	ssoprofilepb "github.com/p1xray/pxr-sso-protos/gen/go/profile"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"pxr-sso/internal/logic/service"
	"pxr-sso/internal/server"
)

const (
	emptyValue = 0
)

type serverAPI struct {
	ssoprofilepb.UnimplementedSsoProfileServer
	profile server.ProfileService
}

// Register registers the implementation of the API service with the gRPC server.
func Register(gRPC *grpc.Server, profileService server.ProfileService) {
	ssoprofilepb.RegisterSsoProfileServer(gRPC, &serverAPI{profile: profileService})
}

func (s *serverAPI) GetProfile(
	ctx context.Context,
	req *ssoprofilepb.GetProfileRequest,
) (*ssoprofilepb.GetProfileResponse, error) {
	if err := validateGetProfileRequest(req); err != nil {
		return nil, err
	}

	userProfile, err := s.profile.UserProfile(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, server.NotFoundError("user not found")
		}

		return nil, server.InternalError("failed to get user profile")
	}

	var dateOfBirthPb *timestamppb.Timestamp
	if userProfile.DateOfBirth != nil {
		dateOfBirthPb = timestamppb.New(*userProfile.DateOfBirth)
	}

	genderPb := ssoprofilepb.Gender_GENDER_UNSPECIFIED
	if userProfile.Gender != nil {
		genderPb = ssoprofilepb.Gender(*userProfile.Gender)
	}

	var avatarFileKeyPb *wrappers.StringValue
	if userProfile.AvatarFileKey != nil {
		avatarFileKeyPb = &wrappers.StringValue{Value: *userProfile.AvatarFileKey}
	}

	return &ssoprofilepb.GetProfileResponse{
		UserId:        userProfile.UserID,
		Username:      userProfile.Username,
		Fio:           userProfile.FIO,
		DateOfBirth:   dateOfBirthPb,
		Gender:        genderPb,
		AvatarFileKey: avatarFileKeyPb,
	}, nil
}

func validateGetProfileRequest(req *ssoprofilepb.GetProfileRequest) error {
	if req.GetUserId() == emptyValue {
		return server.InvalidArgumentError("user id is empty")
	}

	return nil
}
