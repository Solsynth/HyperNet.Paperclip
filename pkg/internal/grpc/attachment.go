package grpc

import (
	"context"

	"git.solsynth.dev/hypernet/paperclip/pkg/internal/database"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/models"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/services"
	"git.solsynth.dev/hypernet/paperclip/pkg/proto"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (v *Server) GetAttachment(ctx context.Context, request *proto.GetAttachmentRequest) (*proto.GetAttachmentResponse, error) {
	tx := database.C
	if request.Id != nil {
		tx = tx.Where("id = ?", request.Id)
	} else if request.Rid != nil {
		tx = tx.Where("rid = ?", request.Rid)
	} else {
		return nil, status.Error(codes.InvalidArgument, "you must provide id or random id")
	}

	if request.UserId != nil {
		tx = tx.Where("account_id = ?", request.UserId)
	}

	var attachment models.Attachment
	if err := tx.First(&attachment).Error; err != nil {
		return nil, status.Error(codes.NotFound, "attachment not found")
	}

	return &proto.GetAttachmentResponse{
		Attachment: lo.ToPtr(attachment).ToAttachmentInfo(),
	}, nil
}

func (v *Server) ListAttachment(ctx context.Context, request *proto.ListAttachmentRequest) (*proto.ListAttachmentResponse, error) {
	tx := database.C
	if len(request.Id) == 0 && len(request.Rid) == 0 {
		return nil, status.Error(codes.InvalidArgument, "you must provide at least one id or random id")
	}
	if len(request.Id) > 0 {
		tx = tx.Where("id IN ?", request.Id)
	}
	if len(request.Rid) > 0 {
		tx = tx.Where("rid IN ?", request.Rid)
	}

	attachments := make([]models.Attachment, 0)
	err := tx.Find(&attachments).Error
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.ListAttachmentResponse{
		Attachments: lo.Map(attachments, func(v models.Attachment, _ int) *proto.AttachmentInfo {
			return v.ToAttachmentInfo()
		}),
	}, nil
}

func (v *Server) UpdateVisibility(ctx context.Context, request *proto.UpdateVisibilityRequest) (*proto.UpdateVisibilityResponse, error) {
	log.Debug().Any("request", request).Msg("Update attachment visibility via grpc...")

	tx := database.C
	if len(request.Id) == 0 && len(request.Rid) == 0 {
		return nil, status.Error(codes.InvalidArgument, "you must provide at least one id or random id")
	}
	if len(request.Id) > 0 {
		tx = tx.Where("id IN ?", request.Id)
	}
	if len(request.Rid) > 0 {
		tx = tx.Where("rid IN ?", request.Rid)
	}

	if request.UserId != nil {
		tx = tx.Where("account_id = ?", request.UserId)
	}

	var rowsAffected int64
	if err := tx.Updates(&models.Attachment{IsIndexable: request.IsIndexable}).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	} else {
		rowsAffected = tx.RowsAffected
	}

	return &proto.UpdateVisibilityResponse{
		Count: int32(rowsAffected),
	}, nil
}

func (v *Server) UpdateUsage(ctx context.Context, request *proto.UpdateUsageRequest) (*proto.UpdateUsageResponse, error) {
	tx := database.C
	if len(request.Id) == 0 && len(request.Rid) == 0 {
		return nil, status.Error(codes.InvalidArgument, "you must provide at least one id or random id")
	}
	if len(request.Id) > 0 {
		tx = tx.Where("id IN ?", request.Id)
	}
	if len(request.Rid) > 0 {
		tx = tx.Where("rid IN ?", request.Rid)
	}

	if rows, err := services.CountAttachmentUsage(tx, int(request.GetDelta())); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	} else {
		return &proto.UpdateUsageResponse{
			Count: int32(rows),
		}, nil
	}
}

func (v *Server) DeleteAttachment(ctx context.Context, request *proto.DeleteAttachmentRequest) (*proto.DeleteAttachmentResponse, error) {
	tx := database.C
	if len(request.Id) == 0 && len(request.Rid) == 0 {
		return nil, status.Error(codes.InvalidArgument, "you must provide at least one id or random id")
	}
	if len(request.Id) > 0 {
		tx = tx.Where("id IN ?", request.Id)
	}
	if len(request.Rid) > 0 {
		tx = tx.Where("rid IN ?", request.Rid)
	}

	if request.UserId != nil {
		tx = tx.Where("account_id = ?", request.UserId)
	}

	var rowsAffected int64
	if err := tx.Delete(&models.Attachment{}).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	} else {
		rowsAffected = tx.RowsAffected
	}

	return &proto.DeleteAttachmentResponse{
		Count: int32(rowsAffected),
	}, nil
}
