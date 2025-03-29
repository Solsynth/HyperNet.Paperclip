package filekit

import (
	"context"
	"time"

	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/cachekit"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/models"
	"git.solsynth.dev/hypernet/paperclip/pkg/proto"
	"github.com/goccy/go-json"
)

func GetAttachment(c *nex.Conn, rid string) (models.Attachment, error) {
	cacheConn, err := cachekit.NewConn(c, 3*time.Second)
	if err == nil {
		key := cachekit.FKey(cachekit.DAAttachment, rid)
		if attachment, err := cachekit.Get[models.Attachment](cacheConn, key); err == nil {
			return attachment, nil
		}
	}

	var attachment models.Attachment
	conn, err := c.GetClientGrpcConn("uc")
	if err != nil {
		return attachment, nil
	}

	pc := proto.NewAttachmentServiceClient(conn)
	resp, err := pc.GetAttachment(context.Background(), &proto.GetAttachmentRequest{
		Rid: &rid,
	})
	if err != nil {
		return attachment, err
	}

	if err := json.Unmarshal(resp.Attachment, &attachment); err != nil {
		return attachment, err
	}

	return attachment, nil
}

func ListAttachment(c *nex.Conn, rid []string) ([]models.Attachment, error) {
	var attachments []models.Attachment
	var missingRid []string
	cachedAttachments := make(map[string]models.Attachment)

	// Try to get attachments from cache
	cacheConn, err := cachekit.NewConn(c, 3*time.Second)
	if err == nil {
		for _, rid := range rid {
			key := cachekit.FKey(cachekit.DAAttachment, rid)
			if attachment, err := cachekit.Get[models.Attachment](cacheConn, key); err == nil {
				cachedAttachments[rid] = attachment
			} else {
				missingRid = append(missingRid, rid)
			}
		}
	}

	// If all attachments are found in cache, return them
	if len(missingRid) == 0 {
		for _, attachment := range cachedAttachments {
			attachments = append(attachments, attachment)
		}
		return attachments, nil
	}

	// Fetch missing attachments from the gRPC service
	conn, err := c.GetClientGrpcConn("uc")
	if err != nil {
		return attachments, err
	}

	pc := proto.NewAttachmentServiceClient(conn)
	resp, err := pc.ListAttachment(context.Background(), &proto.ListAttachmentRequest{
		Rid: missingRid,
	})
	if err != nil {
		return attachments, err
	}

	// Parse the fetched attachments
	for _, item := range resp.GetAttachments() {
		var attachment models.Attachment
		if err := json.Unmarshal(item, &attachment); err != nil {
			return attachments, err
		}
		attachments = append(attachments, attachment)
	}

	// Merge cached and fetched results
	for _, attachment := range cachedAttachments {
		attachments = append(attachments, attachment)
	}

	return attachments, nil
}

func UpdateVisibility(c *nex.Conn, request *proto.UpdateVisibilityRequest) error {
	conn, err := c.GetClientGrpcConn("uc")
	if err != nil {
		return nil
	}

	pc := proto.NewAttachmentServiceClient(conn)
	_, err = pc.UpdateVisibility(context.Background(), request)
	return err
}

func DeleteAttachment(c *nex.Conn, request *proto.DeleteAttachmentRequest) error {
	conn, err := c.GetClientGrpcConn("uc")
	if err != nil {
		return nil
	}

	pc := proto.NewAttachmentServiceClient(conn)
	_, err = pc.DeleteAttachment(context.Background(), request)
	return err
}

func CountAttachmentUsage(c *nex.Conn, request *proto.UpdateUsageRequest) error {
	conn, err := c.GetClientGrpcConn("uc")
	if err != nil {
		return nil
	}

	pc := proto.NewAttachmentServiceClient(conn)
	_, err = pc.UpdateUsage(context.Background(), request)
	return err
}
