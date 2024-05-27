package keeper

import (
	"context"

	"loan/x/loan/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RequestLoan(goCtx context.Context, msg *types.MsgRequestLoan) (*types.MsgRequestLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	loan := types.Loan{
		Amount:     msg.GetAmount(),
		Fee:        msg.GetFee(),
		Collateral: msg.GetCollateral(),
		Deadline:   msg.GetDeadline(),
		State:      "requested",
		Borrower:   msg.GetCreator(),
	}

	borrower, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	collateral, err := sdk.ParseCoinsNormalized(loan.Collateral)
	if err != nil {
		panic(err)
	}

	if sdkError := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		borrower,
		types.ModuleName,
		collateral,
	); sdkError != nil {
		return nil, sdkError
	}

	k.AppendLoan(ctx, loan)

	return &types.MsgRequestLoanResponse{}, nil
}
