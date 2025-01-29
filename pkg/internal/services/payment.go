package services

import (
	"context"
	"fmt"
	"time"

	"git.solsynth.dev/hypernet/paperclip/pkg/internal/gap"
	wproto "git.solsynth.dev/hypernet/wallet/pkg/proto"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

const DiscountFileSize = 52428800 // 50 MiB

// PlaceOrder create a transaction if needed for user
// Pricing according here: https://kb.solsynth.dev/solar-network/wallet#file-uploads
func PlaceOrder(user uint, filesize int64, withDiscount bool) error {
	if filesize <= DiscountFileSize && withDiscount {
		// Discount included
		return nil
	}

	var amount float64
	if withDiscount {
		billableSize := filesize - DiscountFileSize
		amount = float64(billableSize) / 1024 / 1024 * 1
	} else if filesize > DiscountFileSize {
		amount = 50 + float64(filesize-DiscountFileSize)/1024/1024*5
	} else {
		amount = float64(filesize) / 1024 / 1024 * 1
	}

	if !withDiscount {
		amount += 10 // Service fee
	}

	conn, err := gap.Nx.GetClientGrpcConn("wa")
	if err != nil {
		return fmt.Errorf("unable to connect wallet: %v", err)
	}

	wc := wproto.NewPaymentServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	resp, err := wc.MakeTransactionWithAccount(ctx, &wproto.MakeTransactionWithAccountRequest{
		PayerAccountId: lo.ToPtr(uint64(user)),
		Amount:         amount,
		Remark:         "File Uploading Fee",
	})
	if err != nil {
		return err
	}

	log.Info().
		Uint64("transaction", resp.Id).Float64("amount", amount).Bool("discount", withDiscount).
		Msg("Order placed for charge file uploading fee...")

	return nil
}
