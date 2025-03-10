package filekit

import (
	"context"

	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"git.solsynth.dev/hypernet/paperclip/pkg/proto"
)

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
