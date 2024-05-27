package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"

	"loan/x/loan/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ApproveLoan(goCtx context.Context, msg *types.MsgApproveLoan) (*types.MsgApproveLoanResponse, error) {
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

	lender, _ := sdk.AccAddressFromBech32(msg.Creator)
	borrower, _ := sdk.AccAddressFromBech32(loan.Borrower)
	amount, err := sdk.ParseCoinsNormalized(loan.Amount)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidLengthLoan, "Cannot parse coins in loan amount")
	}

	err = k.bankKeeper.SendCoins(ctx, lender, borrower, amount)
	if err != nil {
		return nil, err
	}
	loan.Lender = msg.Creator
	loan.State = "approved"
	k.SetLoan(ctx, loan)

	return &types.MsgApproveLoanResponse{}, nil
}
