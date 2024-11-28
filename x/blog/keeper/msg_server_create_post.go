package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"blog/x/blog/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreatePost(goCtx context.Context, msg *types.MsgCreatePost) (*types.MsgCreatePostResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Capture start time
	startTime := time.Now()

	var post = types.Post{
		Creator: msg.Creator,
		Title:   msg.Title,
		Body:    msg.Body,
	}
	id := k.AppendPost(
		ctx, post,
	)

	// Extract transaction metadata
	txHash := ctx.TxBytes()                           // Raw transaction bytes
	timestamp := ctx.BlockTime().Format(time.RFC3339) // Transaction timestamp
	gasWanted := ctx.GasMeter().Limit()               // Gas limit
	gasUsed := ctx.GasMeter().GasConsumed()           // Gas used

	timeTaken := time.Since(startTime).Milliseconds()

	// Create a log entry (example in JSON format)
	logEntry := fmt.Sprintf(`{
        "txHash": "%s",
        "timestamp": "%s",
        "gasWanted": %d,
        "gasUsed": %d,
        "msgSize": %d,
		"title": "%s",
		"body": "%s",
		"user": "%s",
		"timeTakenMs": %d
	}`, hex.EncodeToString(txHash), timestamp, gasWanted, gasUsed, len(ctx.TxBytes()), msg.Title, msg.Body, msg.Creator, timeTaken)

	// Write the log to the chain’s logger
	logFile, err := os.OpenFile("/Users/niketshwetabh/Documents/ignite-cli/blog/logs/tx_logs.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		ctx.Logger().Error("Failed to open log file", "error", err)
	}
	defer logFile.Close()

	logFile.WriteString(logEntry + "\n")

	return &types.MsgCreatePostResponse{
		Id: id,
	}, nil
}
