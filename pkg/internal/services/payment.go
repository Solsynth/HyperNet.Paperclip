package services

import (
	"context"
	"fmt"
	"time"

	"git.solsynth.dev/hypernet/paperclip/pkg/internal/database"
	"git.solsynth.dev/hypernet/paperclip/pkg/internal/models"

	"git.solsynth.dev/hypernet/paperclip/pkg/internal/gap"
	wproto "git.solsynth.dev/hypernet/wallet/pkg/proto"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

func GetLastDayUploadedBytes(user uint) (int64, error) {
	deadline := time.Now().Add(-24 * time.Hour)
	var totalSize int64
	if err := database.C.
		Model(&models.Attachment{}).
		Where("account_id = ?", user).
		Where("created_at >= ?", deadline).
		Select("SUM(size)").
		Scan(&totalSize).Error; err != nil {
		return totalSize, err
	}
	return totalSize, nil
}

// PlaceOrder create a transaction if needed for user
// Pricing according here: https://kb.solsynth.dev/solar-network/wallet#file-uploads
func PlaceOrder(user uint, filesize int64, withDiscount bool) error {
	currentBytes, _ := GetLastDayUploadedBytes(user)
	discountFileSize := viper.GetInt64("payment.discount")

	if currentBytes+filesize <= discountFileSize {
		// Discount included
		return nil
	}

	var amount float64
	if withDiscount {
		amount = float64(filesize) / 1024 / 1024 * 1
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
