package proposal

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestGetProposalPassParam(t *testing.T) {
	ctx, _, pm, _, _, _, _ := setupTest(t, 0)

	proposalParam, _ := pm.paramHolder.GetProposalParam(ctx)
	testCases := []struct {
		testName      string
		proposalType  types.ProposalType
		wantError     sdk.Error
		wantPassRatio sdk.Rat
		wantPassVotes types.Coin
	}{
		{testName: "test pass param for changeParamProposal",
			proposalType:  types.ChangeParam,
			wantError:     nil,
			wantPassRatio: proposalParam.ChangeParamPassRatio,
			wantPassVotes: proposalParam.ChangeParamPassVotes,
		},

		{testName: "test pass param for contenCensorshipProposal",
			proposalType:  types.ContentCensorship,
			wantError:     nil,
			wantPassRatio: proposalParam.ContentCensorshipPassRatio,
			wantPassVotes: proposalParam.ContentCensorshipPassVotes,
		},

		{testName: "test pass param for protocolUpgradeProposal",
			proposalType:  types.ProtocolUpgrade,
			wantError:     nil,
			wantPassRatio: proposalParam.ProtocolUpgradePassRatio,
			wantPassVotes: proposalParam.ProtocolUpgradePassVotes,
		},

		{testName: "test wrong proposal type",
			proposalType:  23,
			wantError:     ErrWrongProposalType(),
			wantPassRatio: proposalParam.ProtocolUpgradePassRatio,
			wantPassVotes: proposalParam.ProtocolUpgradePassVotes,
		},
	}
	for _, tc := range testCases {
		ratio, votes, err := pm.GetProposalPassParam(ctx, tc.proposalType)
		assert.Equal(t, tc.wantError, err)
		if tc.wantError != nil {
			continue
		}
		assert.Equal(t, tc.wantPassRatio, ratio)
		assert.Equal(t, tc.wantPassVotes, votes)
	}

}
