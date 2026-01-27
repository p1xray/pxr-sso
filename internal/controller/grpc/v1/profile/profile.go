package profile

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/wrappers"
	ssoprofilepb "github.com/p1xray/pxr-sso-protos/gen/go/profile"
	"github.com/p1xray/pxr-sso/internal/controller"
	"github.com/p1xray/pxr-sso/internal/controller/grpc/response"
	"github.com/p1xray/pxr-sso/internal/enum"
	"github.com/p1xray/pxr-sso/internal/usecase"
	"github.com/p1xray/pxr-sso/internal/usecase/profile/edit"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

const (
	emptyValue = 0
)

type serverAPI struct {
	ssoprofilepb.UnimplementedSsoProfileServer
	profile controller.UserProfile
	edit    controller.EditProfile
}

// RegisterProfileServer registers the implementation of the API service with the gRPC server.
func RegisterProfileServer(gRPC *grpc.Server, profile controller.UserProfile, edit controller.EditProfile) {
	api := &serverAPI{
		profile: profile,
		edit:    edit,
	}

	ssoprofilepb.RegisterSsoProfileServer(gRPC, api)
}

// GetProfile is a gRPC handler for getting user profile data.
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
	if req.GetUserId() == emptyValue {
		return response.InvalidArgumentError("user id is empty")
	}

	return nil
}

// EditProfile is a gRPC handler for editing user profile data.
func (s *serverAPI) EditProfile(
	ctx context.Context,
	req *ssoprofilepb.EditProfileRequest,
) (*ssoprofilepb.EditProfileResponse, error) {
	if err := validateEditProfileRequest(req); err != nil {
		return nil, err
	}

	var dateOfBirth *time.Time
	if req.GetDateOfBirth() != nil {
		dateOfBirthPbAsTime := req.GetDateOfBirth().AsTime()
		dateOfBirth = &dateOfBirthPbAsTime
	}

	var gender *enum.GenderEnum
	if req.GetGender() != emptyValue {
		genderEnum := enum.GenderEnum(req.GetGender().Number())
		gender = &genderEnum
	}

	var avatarFileKey *string
	if req.GetAvatarFileKey() != nil {
		avatarFileKeyPbString := req.GetAvatarFileKey().GetValue()
		avatarFileKey = &avatarFileKeyPbString
	}

	params := edit.Params{
		ID:            req.GetUserId(),
		FullName:      req.GetFullName(),
		DateOfBirth:   dateOfBirth,
		Gender:        gender,
		AvatarFileKey: avatarFileKey,
	}

	if err := s.edit.Execute(ctx, params); err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			return nil, response.NotFoundError("user not found")
		}

		return nil, response.InternalError("failed to edit user profile data")
	}

	return &ssoprofilepb.EditProfileResponse{Success: true}, nil
}

func validateEditProfileRequest(req *ssoprofilepb.EditProfileRequest) error {
	if req.GetUserId() == emptyValue {
		return response.InvalidArgumentError("user id is empty")
	}

	return nil
}
