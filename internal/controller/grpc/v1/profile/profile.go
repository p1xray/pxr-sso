package profile

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/wrappers"
	ssoprofilepb "github.com/p1xray/pxr-sso-protos/gen/go/profile"
	"github.com/p1xray/pxr-sso/internal/controller"
	"github.com/p1xray/pxr-sso/internal/controller/grpc/response"
	"github.com/p1xray/pxr-sso/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	emptyID = 0
)

type serverAPI struct {
	ssoprofilepb.UnimplementedSsoProfileServer
	profile controller.UserProfile
}

// RegisterProfileServer registers the implementation of the API service with the gRPC server.
func RegisterProfileServer(gRPC *grpc.Server, profile controller.UserProfile) {
	ssoprofilepb.RegisterSsoProfileServer(gRPC, &serverAPI{profile: profile})
}

func (s *serverAPI) GetProfile(
	ctx context.Context,
	req *ssoprofilepb.GetProfileRequest,
) (*ssoprofilepb.GetProfileResponse, error) {
	if err := validateGetProfileRequest(req); err != nil {
		return nil, err
	}

	userProfile, err := s.profile.Execute(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			return nil, response.NotFoundError("user not found")
		}

		return nil, response.InternalError("failed to get user profile")
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
		UserId:        userProfile.ID,
		Username:      userProfile.Username,
		Fio:           userProfile.FullName,
		DateOfBirth:   dateOfBirthPb,
		Gender:        genderPb,
		AvatarFileKey: avatarFileKeyPb,
	}, nil
}

func validateGetProfileRequest(req *ssoprofilepb.GetProfileRequest) error {
	if req.GetUserId() == emptyID {
		return response.InvalidArgumentError("user id is empty")
	}

	return nil
}
