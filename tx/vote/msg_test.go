package vote

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/tx/vote/model"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestVoteMsg(t *testing.T) {
	cases := []struct {
		voteMsg     VoteMsg
		expectError sdk.Error
	}{
		{NewVoteMsg("user1", 1, true), nil},
		{NewVoteMsg("", 1, true), ErrInvalidUsername()},
	}

	for _, cs := range cases {
		result := cs.voteMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestVoterDepositMsg(t *testing.T) {
	cases := []struct {
		voterDepositMsg VoterDepositMsg
		expectError     sdk.Error
	}{
		{NewVoterDepositMsg("user1", "1"), nil},
		{NewVoterDepositMsg("", "1"), ErrInvalidUsername()},
		{NewVoterDepositMsg("user1", "-1"), sdk.ErrInvalidCoins("LNO can't be less than lower bound")},
	}

	for _, cs := range cases {
		result := cs.voterDepositMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestVoterDepositMsgPermission(t *testing.T) {
	msg := NewVoterDepositMsg("user1", "1")
	permissionLevel := msg.Get(types.PermissionLevel)
	permission, ok := permissionLevel.(int)
	assert.Equal(t, ok, true)
	assert.Equal(t, permission, types.Active)
}

func TestVoterWithdrawMsg(t *testing.T) {
	cases := []struct {
		voterWithdrawMsg VoterWithdrawMsg
		expectError      sdk.Error
	}{
		{NewVoterWithdrawMsg("user1", "1"), nil},
		{NewVoterWithdrawMsg("", "1"), ErrInvalidUsername()},
		{NewVoterWithdrawMsg("user1", "-1"), sdk.ErrInvalidCoins("LNO can't be less than lower bound")},
	}

	for _, cs := range cases {
		result := cs.voterWithdrawMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestVoterRevokeMsg(t *testing.T) {
	cases := []struct {
		voterRevokeMsg VoterRevokeMsg
		expectError    sdk.Error
	}{
		{NewVoterRevokeMsg("user1"), nil},
		{NewVoterRevokeMsg(""), ErrInvalidUsername()},
	}

	for _, cs := range cases {
		result := cs.voterRevokeMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestDelegateMsg(t *testing.T) {
	cases := []struct {
		delegateMsg DelegateMsg
		expectError sdk.Error
	}{
		{NewDelegateMsg("user1", "user2", "1"), nil},
		{NewDelegateMsg("", "user2", "1"), ErrInvalidUsername()},
		{NewDelegateMsg("user1", "", "1"), ErrInvalidUsername()},
		{NewDelegateMsg("", "", "1"), ErrInvalidUsername()},
		{NewDelegateMsg("user1", "user2", "-1"), sdk.ErrInvalidCoins("LNO can't be less than lower bound")},
	}

	for _, cs := range cases {
		result := cs.delegateMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestDelegateMsgPermission(t *testing.T) {
	msg := NewDelegateMsg("user1", "user2", "1")
	permissionLevel := msg.Get(types.PermissionLevel)
	permission, ok := permissionLevel.(int)
	assert.Equal(t, ok, true)
	assert.Equal(t, permission, types.Active)
}

func TestRevokeDelegationMsg(t *testing.T) {
	cases := []struct {
		revokeDelegationMsg RevokeDelegationMsg
		expectError         sdk.Error
	}{
		{NewRevokeDelegationMsg("user1", "user2"), nil},
		{NewRevokeDelegationMsg("", "user2"), ErrInvalidUsername()},
		{NewRevokeDelegationMsg("user1", ""), ErrInvalidUsername()},
		{NewRevokeDelegationMsg("", ""), ErrInvalidUsername()},
	}

	for _, cs := range cases {
		result := cs.revokeDelegationMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestDelegatorWithdrawMsg(t *testing.T) {
	cases := []struct {
		delegatorWithdrawMsg DelegatorWithdrawMsg
		expectError          sdk.Error
	}{
		{NewDelegatorWithdrawMsg("user1", "user2", "1"), nil},
		{NewDelegatorWithdrawMsg("", "", "1"), ErrInvalidUsername()},
		{NewDelegatorWithdrawMsg("user1", "user2", "-1"), sdk.ErrInvalidCoins("LNO can't be less than lower bound")},
	}

	for _, cs := range cases {
		result := cs.delegatorWithdrawMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestCreateProposalMsg(t *testing.T) {
	des1 := model.ChangeParameterDescription{
		InfraAllocation:          sdk.Rat{20, 100},
		ContentCreatorAllocation: sdk.Rat{55, 100},
		DeveloperAllocation:      sdk.Rat{20, 100},
		ValidatorAllocation:      sdk.Rat{5, 100},
		CDNAllocation:            sdk.Rat{5, 100},
		StorageAllocation:        sdk.Rat{95, 100},
	}

	des2 := model.ChangeParameterDescription{
		InfraAllocation:          sdk.Rat{20, 100},
		ContentCreatorAllocation: sdk.Rat{55, 100},
		DeveloperAllocation:      sdk.Rat{25, 100},
		ValidatorAllocation:      sdk.Rat{5, 100},
		CDNAllocation:            sdk.Rat{5, 100},
		StorageAllocation:        sdk.Rat{95, 100},
	}

	des3 := model.ChangeParameterDescription{
		InfraAllocation:          sdk.Rat{20, 100},
		ContentCreatorAllocation: sdk.Rat{55, 100},
		DeveloperAllocation:      sdk.Rat{20, 100},
		ValidatorAllocation:      sdk.Rat{5, 100},
		CDNAllocation:            sdk.Rat{15, 100},
		StorageAllocation:        sdk.Rat{95, 100},
	}
	cases := []struct {
		createProposalMsg CreateProposalMsg
		expectError       sdk.Error
	}{
		{NewCreateProposalMsg("user1", des1), nil},
		{NewCreateProposalMsg("user1", des2), ErrIllegalParameter()},
		{NewCreateProposalMsg("user1", des3), ErrIllegalParameter()},
		{NewCreateProposalMsg("", des1), ErrInvalidUsername()},
	}

	for _, cs := range cases {
		result := cs.createProposalMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}
