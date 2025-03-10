package grpc

import (
	"context"

	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	jsoniter "github.com/json-iterator/go"

	"git.solsynth.dev/hypernet/nexus/pkg/proto"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/database"
)

func (v *Server) BroadcastEvent(ctx context.Context, in *proto.EventInfo) (*proto.EventResponse, error) {
	switch in.GetEvent() {
	case "deletion":
		data := nex.DecodeMap(in.GetData())
		resType, ok := data["type"].(string)
		if !ok {
			break
		}
		switch resType {
		case "account":
			var data struct {
				ID int `json:"id"`
			}
			if err := jsoniter.Unmarshal(in.GetData(), &data); err != nil {
				break
			}
			tx := database.C.Begin()
			for _, model := range database.AutoMaintainRange {
				switch model.(type) {
				default:
					tx.Delete(model, "account_id = ?", data.ID)
				}
			}
			tx.Commit()
		}
	}

	return &proto.EventResponse{}, nil
}
