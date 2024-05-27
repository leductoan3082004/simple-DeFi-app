package keeper

import (
	"context"
	"fmt"

	"loan/x/loan/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CancelLoan(goCtx context.Context, msg *types.MsgCancelLoan) (*types.MsgCancelLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	loan, found := k.GetLoan(ctx, msg.GetId())
	if !found {
		return nil, fmt.Errorf("loan %d not found", msg.GetId())
	}

	if loan.GetState() != "requested" {
		return nil, fmt.Errorf("loan %d's state is not requested", msg.GetId())
	}

	if loan.Borrower != msg.GetCreator() {
		return nil, fmt.Errorf("loan %d's borrower is not the owner", msg.GetId())
	}

	// maybe do not care the error
	borrower, err := sdk.AccAddressFromBech32(msg.GetCreator())
	if err != nil {
		panic(err)
	}

	// maybe do not care the error
	collateral, err := sdk.ParseCoinsNormalized(loan.Collateral)
	if err != nil {
		panic(err)
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		borrower,
		collateral,
	); err != nil {
		return nil, err
	}

	loan.State = "canceled"
	k.SetLoan(ctx, loan)

	return &types.MsgCancelLoanResponse{}, nil
}
