package gap

import (
	"fmt"
	"strings"
	"time"

	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/cachekit"
	"git.solsynth.dev/hypernet/nexus/pkg/proto"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"

	"github.com/spf13/viper"
)

var (
	Nx *nex.Conn
	Ca *cachekit.Conn
)

func InitializeToNexus() error {
	grpcBind := strings.SplitN(viper.GetString("grpc_bind"), ":", 2)
	httpBind := strings.SplitN(viper.GetString("bind"), ":", 2)

	outboundIp, _ := nex.GetOutboundIP()

	grpcOutbound := fmt.Sprintf("%s:%s", outboundIp, grpcBind[1])
	httpOutbound := fmt.Sprintf("%s:%s", outboundIp, httpBind[1])

	var err error
	Nx, err = nex.NewNexusConn(viper.GetString("nexus_addr"), &proto.ServiceInfo{
		Id:       viper.GetString("id"),
		Type:     "uc",
		Label:    "Paperclip",
		GrpcAddr: grpcOutbound,
		HttpAddr: lo.ToPtr("http://" + httpOutbound + "/api"),
	})
	if err == nil {
		go func() {
			err := Nx.RunRegistering()
			if err != nil {
				log.Error().Err(err).Msg("An error occurred while registering service...")
			}
		}()
	}

	if Ca, err = cachekit.NewConn(Nx, 3*time.Second); err != nil {
		return fmt.Errorf("failed to create cachekit connection: %v", err)
	}

	return err
}
