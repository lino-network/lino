package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/lino-network/lino/types"
)

var (
	// len of 1000
	maxLengthUTF8Str = `
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌123`

	// len of 1001
	tooLongUTF8Str = `
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌1234`
)

type DeveloperMsgTestSuite struct {
	suite.Suite
}

func TestDeveloperMsgTestSuite(t *testing.T) {
	suite.Run(t, new(DeveloperMsgTestSuite))
}

func (suite *DeveloperMsgTestSuite) TestDeveloperRegisterMsgValidateBasic() {
	testCases := []struct {
		testName             string
		developerRegisterMsg DeveloperRegisterMsg
		expectError          sdk.Error
	}{
		{
			testName:             "normal case",
			developerRegisterMsg: NewDeveloperRegisterMsg("user1", "", "", ""),
			expectError:          nil,
		},
		{
			testName: "utf8 description",
			developerRegisterMsg: NewDeveloperRegisterMsg(
				"user1", "", maxLengthUTF8Str, ""),
			expectError: nil,
		},
		{
			testName: "utf8 app metadata",
			developerRegisterMsg: NewDeveloperRegisterMsg(
				"user1", "", "", "a:bb,你:👌 。"),
			expectError: nil,
		},
		{
			testName:             "invalid username",
			developerRegisterMsg: NewDeveloperRegisterMsg("", "10", "", ""),
			expectError:          ErrInvalidUsername(),
		},
		{
			testName: "website is too long",
			developerRegisterMsg: NewDeveloperRegisterMsg(
				"user1", string(make([]byte, types.MaximumLengthOfDeveloperWebsite+1)), "", ""),
			expectError: ErrInvalidWebsite(),
		},
		{
			testName: "description is too long",
			developerRegisterMsg: NewDeveloperRegisterMsg(
				"user1", "", string(make([]byte, types.MaximumLengthOfDeveloperDesctiption+1)), ""),
			expectError: ErrInvalidDescription(),
		},
		{
			testName: "utf8 description is too long",
			developerRegisterMsg: NewDeveloperRegisterMsg(
				"user1", "", tooLongUTF8Str, ""),
			expectError: ErrInvalidDescription(),
		},
		{
			testName: "app metadata is too long",
			developerRegisterMsg: NewDeveloperRegisterMsg(
				"user1", "", "", string(make([]byte, types.MaximumLengthOfAppMetadata+1))),
			expectError: ErrInvalidAppMetadata(),
		},
		{
			testName: "utf8 app metadata is too long",
			developerRegisterMsg: NewDeveloperRegisterMsg(
				"user1", "", "", tooLongUTF8Str),
			expectError: ErrInvalidAppMetadata(),
		},
	}

	for _, tc := range testCases {
		result := tc.developerRegisterMsg.ValidateBasic()
		suite.Equal(result, tc.expectError, "%s", tc.testName)
	}
}

func (suite *DeveloperMsgTestSuite) TestDeveloperUpdateMsgValidateBasic() {
	testCases := []struct {
		testName           string
		developerUpdateMsg DeveloperUpdateMsg
		expectError        sdk.Error
	}{
		{
			testName:           "normal case",
			developerUpdateMsg: NewDeveloperUpdateMsg("user1", "", "", ""),
			expectError:        nil,
		},
		{
			testName: "utf8 description",
			developerUpdateMsg: NewDeveloperUpdateMsg(
				"user1", "", maxLengthUTF8Str, ""),
			expectError: nil,
		},
		{
			testName: "uft8 app metadata",
			developerUpdateMsg: NewDeveloperUpdateMsg(
				"user1", "", "", maxLengthUTF8Str),
			expectError: nil,
		},
		{
			testName:           "invalid username",
			developerUpdateMsg: NewDeveloperUpdateMsg("", "", "", ""),
			expectError:        ErrInvalidUsername(),
		},
		{
			testName: "website is too long",
			developerUpdateMsg: NewDeveloperUpdateMsg(
				"user1", string(make([]byte, types.MaximumLengthOfDeveloperWebsite+1)), "", ""),
			expectError: ErrInvalidWebsite(),
		},
		{
			testName: "description is too long",
			developerUpdateMsg: NewDeveloperUpdateMsg(
				"user1", "", string(make([]byte, types.MaximumLengthOfDeveloperDesctiption+1)), ""),
			expectError: ErrInvalidDescription(),
		},
		{
			testName: "utf8 description is too long",
			developerUpdateMsg: NewDeveloperUpdateMsg(
				"user1", "", tooLongUTF8Str, ""),
			expectError: ErrInvalidDescription(),
		},
		{
			testName: "app metadata is too long",
			developerUpdateMsg: NewDeveloperUpdateMsg(
				"user1", "", "", string(make([]byte, types.MaximumLengthOfAppMetadata+1))),
			expectError: ErrInvalidAppMetadata(),
		},
		{
			testName: "utf8 app metadata is too long",
			developerUpdateMsg: NewDeveloperUpdateMsg(
				"user1", "", "", tooLongUTF8Str),
			expectError: ErrInvalidAppMetadata(),
		},
	}

	for _, tc := range testCases {
		result := tc.developerUpdateMsg.ValidateBasic()
		suite.Equal(tc.expectError, result, "%s: %v", tc.testName, result)
	}
}

// func TestDeveloperRevokeMsg(t *testing.T) {
// 	testCases := []struct {
// 		testName           string
// 		developerRevokeMsg DeveloperRevokeMsg
// 		expectError        sdk.Error
// 	}{
// 		{
// 			testName:           "normal case",
// 			developerRevokeMsg: NewDeveloperRevokeMsg("user1"),
// 			expectError:        nil,
// 		},
// 		{
// 			testName:           "invalid username",
// 			developerRevokeMsg: NewDeveloperRevokeMsg(""),
// 			expectError:        ErrInvalidUsername(),
// 		},
// 	}

// 	for _, tc := range testCases {
// 		result := tc.developerRevokeMsg.ValidateBasic()
// 		if !assert.Equal(t, result, tc.expectError) {
// 			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectError)
// 		}
// 	}
// }

func (suite *DeveloperMsgTestSuite) TestIDAIssueMsgValidateBasic() {
	testCases := []struct {
		testName    string
		msg         IDAIssueMsg
		expectedErr sdk.Error
	}{
		{
			testName: "invalid username1",
			msg: IDAIssueMsg{
				Username: "1",
				IDAPrice: 123,
			},
			expectedErr: ErrInvalidUsername(),
		},
		{
			testName: "invalid username2",
			msg: IDAIssueMsg{
				Username: "longlonglonglonglonglonglonglongusername",
				IDAPrice: 123,
			},
			expectedErr: ErrInvalidUsername(),
		},
		{
			testName: "invalid ida price 1",
			msg: IDAIssueMsg{
				Username: "lino",
				IDAPrice: 1001,
			},
			expectedErr: ErrInvalidIDAPrice(),
		},
		{
			testName: "invalid ida price 1",
			msg: IDAIssueMsg{
				Username: "lino",
				IDAPrice: 0,
			},
			expectedErr: ErrInvalidIDAPrice(),
		},
		{
			testName: "ok",
			msg: IDAIssueMsg{
				Username: "lino2",
				IDAPrice: 123,
			},
			expectedErr: nil,
		},
	}
	for _, tc := range testCases {
		suite.Equal(tc.expectedErr, tc.msg.ValidateBasic(), "%s", tc.testName)
	}
}

func (suite *DeveloperMsgTestSuite) TestIDAMintMsgValidateBasic() {
	testCases := []struct {
		testName    string
		msg         IDAMintMsg
		expectedErr sdk.Error
	}{
		{
			testName: "invalid username1",
			msg: IDAMintMsg{
				Username: "1",
				Amount:   "123",
			},
			expectedErr: ErrInvalidUsername(),
		},
		{
			testName: "invalid username2",
			msg: IDAMintMsg{
				Username: "longlonglonglonglonglonglonglongusername",
				Amount:   "123",
			},
			expectedErr: ErrInvalidUsername(),
		},
		{
			testName: "invalid amount1",
			msg: IDAMintMsg{
				Username: "user1",
				Amount:   "10000000000000000000000000000000000000000000000000000000000",
			},
			expectedErr: types.ErrInvalidCoins("LNO overflow"),
		},
		{
			testName: "invalid amount2",
			msg: IDAMintMsg{
				Username: "user1",
				Amount:   "-1",
			},
			expectedErr: types.ErrInvalidCoins("LNO can't be less than lower bound"),
		},
		{
			testName: "invalid amount3",
			msg: IDAMintMsg{
				Username: "user1",
				Amount:   "0",
			},
			expectedErr: types.ErrInvalidCoins("LNO can't be less than lower bound"),
		},
		{
			testName: "invalid amount4",
			msg: IDAMintMsg{
				Username: "user1",
				Amount:   "-10000000000000000000000000000000000000000000000000000000000",
			},
			expectedErr: types.ErrInvalidCoins("LNO can't be less than lower bound"),
		},
		{
			testName: "invalid amount5",
			msg: IDAMintMsg{
				Username: "user1",
				Amount:   "0x3242",
			},
			expectedErr: types.ErrInvalidCoins("Illegal LNO"),
		},
	}
	for _, tc := range testCases {
		suite.Equal(tc.expectedErr, tc.msg.ValidateBasic(), "%s", tc.testName)
	}
}

func (suite *DeveloperMsgTestSuite) TestIDATransferMsgValidateBasic() {
	testCases := []struct {
		testName    string
		msg         IDATransferMsg
		expectedErr sdk.Error
	}{
		{
			testName: "invalid app",
			msg: IDATransferMsg{
				App:    "1",
				From:   "user1",
				To:     "app1",
				Amount: "123",
				Signer: "app1",
			},
			expectedErr: ErrInvalidUsername(),
		},
		{
			testName: "invalid from",
			msg: IDATransferMsg{
				App:    "app1",
				From:   "",
				To:     "app1",
				Amount: "123",
				Signer: "app1",
			},
			expectedErr: ErrInvalidUsername(),
		},
		{
			testName: "invalid to",
			msg: IDATransferMsg{
				App:    "app1",
				From:   "user1",
				To:     "",
				Amount: "123",
				Signer: "app1",
			},
			expectedErr: ErrInvalidUsername(),
		},
		{
			testName: "invalid target 1",
			msg: IDATransferMsg{
				App:    "app1",
				From:   "user1",
				To:     "user2",
				Amount: "123",
				Signer: "app1",
			},
			expectedErr: ErrInvalidTransferTarget(),
		},
		{
			testName: "self transfer",
			msg: IDATransferMsg{
				App:    "app1",
				From:   "user1",
				To:     "user1",
				Amount: "123",
				Signer: "app1",
			},
			expectedErr: ErrIDATransferSelf(),
		},
		{
			testName: "invalid amount 1",
			msg: IDATransferMsg{
				App:    "app1",
				From:   "user1",
				To:     "app1",
				Amount: "0",
				Signer: "app1",
			},
			expectedErr: types.ErrInvalidIDAAmount(),
		},
		{
			testName: "invalid amount 2",
			msg: IDATransferMsg{
				App:    "app1",
				From:   "user1",
				To:     "app1",
				Amount: "-123",
				Signer: "app1",
			},
			expectedErr: types.ErrInvalidIDAAmount(),
		},
		{
			testName: "invalid amount 3",
			msg: IDATransferMsg{
				App:    "app1",
				From:   "user1",
				To:     "app1",
				Amount: "100000000000000000000000000000000000000000000000",
				Signer: "app1",
			},
			expectedErr: types.ErrInvalidIDAAmount(),
		},
		{
			testName: "invalid amount 4",
			msg: IDATransferMsg{
				App:    "app1",
				From:   "user1",
				To:     "app1",
				Amount: "0x011",
				Signer: "app1",
			},
			expectedErr: types.ErrInvalidIDAAmount(),
		},
		{
			testName: "ok1",
			msg: IDATransferMsg{
				App:    "app1",
				From:   "user1",
				To:     "app1",
				Amount: "123",
				Signer: "x",
			},
			expectedErr: ErrInvalidUsername(),
		},
		{
			testName: "ok user to app",
			msg: IDATransferMsg{
				App:    "app1",
				From:   "user1",
				To:     "app1",
				Amount: "123",
				Signer: "app1",
			},
			expectedErr: nil,
		},
		{
			testName: "ok app to user",
			msg: IDATransferMsg{
				App:    "app1",
				From:   "app1",
				To:     "user1",
				Amount: "123345",
				Signer: "app1",
			},
			expectedErr: nil,
		},
		{
			testName: "ok signer others",
			msg: IDATransferMsg{
				App:    "app1",
				From:   "app1",
				To:     "user1",
				Amount: "123345",
				Signer: "app1affiliated",
			},
			expectedErr: nil,
		},
	}
	for _, tc := range testCases {
		suite.Equal(tc.expectedErr, tc.msg.ValidateBasic(), "%s", tc.testName)
	}
}

func (suite *DeveloperMsgTestSuite) TestIDAAuthorizeMsgValidateBasic() {
	testCases := []struct {
		testName    string
		msg         IDAAuthorizeMsg
		expectedErr sdk.Error
	}{
		{
			testName: "invalid username",
			msg: IDAAuthorizeMsg{
				Username: "",
				App:      "app1",
				Activate: true,
			},
			expectedErr: ErrInvalidUsername(),
		},
		{
			testName: "invalid app",
			msg: IDAAuthorizeMsg{
				Username: "user1",
				App:      "x",
				Activate: true,
			},
			expectedErr: ErrInvalidUsername(),
		},
		{
			testName: "authrize self",
			msg: IDAAuthorizeMsg{
				Username: "app1",
				App:      "app1",
				Activate: true,
			},
			expectedErr: ErrInvalidIDAAuth(),
		},
		{
			testName: "ok",
			msg: IDAAuthorizeMsg{
				Username: "user1",
				App:      "app1",
				Activate: true,
			},
			expectedErr: nil,
		},
	}
	for _, tc := range testCases {
		suite.Equal(tc.expectedErr, tc.msg.ValidateBasic(), "%s", tc.testName)
	}
}

func (suite *DeveloperMsgTestSuite) TestUpdateAffiliatedMsgValidateBasic() {
	testCases := []struct {
		testName    string
		msg         UpdateAffiliatedMsg
		expectedErr sdk.Error
	}{
		{
			testName: "invalid app",
			msg: UpdateAffiliatedMsg{
				App:      "x",
				Username: "user1",
				Activate: true,
			},
			expectedErr: ErrInvalidUsername(),
		},
		{
			testName: "invalid user",
			msg: UpdateAffiliatedMsg{
				App:      "app1",
				Username: "x",
				Activate: true,
			},
			expectedErr: ErrInvalidUsername(),
		},
		{
			testName: "ok",
			msg: UpdateAffiliatedMsg{
				App:      "app1",
				Username: "user1",
				Activate: true,
			},
			expectedErr: nil,
		},
	}
	for _, tc := range testCases {
		suite.Equal(tc.expectedErr, tc.msg.ValidateBasic(), "%s", tc.testName)
	}
}

func (suite *DeveloperMsgTestSuite) TestMsgPermission() {
	testCases := []struct {
		testName         string
		msg              types.Msg
		expectPermission types.Permission
	}{
		{
			testName:         "developer register msg",
			msg:              NewDeveloperRegisterMsg("test", "", "", ""),
			expectPermission: types.TransactionPermission,
		},
		{
			testName:         "developer register msg",
			msg:              NewDeveloperUpdateMsg("test", "", "", ""),
			expectPermission: types.TransactionPermission,
		},
		{
			testName:         "developer revoke msg",
			msg:              NewDeveloperRevokeMsg("test"),
			expectPermission: types.TransactionPermission,
		},
		{
			testName:         "ida issue msg",
			msg:              IDAIssueMsg{Username: "user1"},
			expectPermission: types.TransactionPermission,
		},
		{
			testName:         "ida mint msg",
			msg:              IDAMintMsg{Username: "user1"},
			expectPermission: types.TransactionPermission,
		},
		{
			testName:         "ida transfer msg",
			msg:              IDATransferMsg{App: "app1"},
			expectPermission: types.AppOrAffiliatedPermission,
		},
		{
			testName:         "ida auth msg",
			msg:              IDAAuthorizeMsg{Username: "user1"},
			expectPermission: types.TransactionPermission,
		},
		{
			testName:         "update affiliated",
			msg:              UpdateAffiliatedMsg{App: "app1"},
			expectPermission: types.TransactionPermission,
		},
	}

	for _, tc := range testCases {
		permission := tc.msg.GetPermission()
		suite.Equal(tc.expectPermission, permission, "%s", tc.testName)
	}
}

func (suite *DeveloperMsgTestSuite) TestGetSigners() {
	testCases := []struct {
		testName      string
		msg           types.Msg
		expectSigners []types.AccountKey
	}{
		{
			testName:      "developer register msg",
			msg:           NewDeveloperRegisterMsg("test", "", "", ""),
			expectSigners: []types.AccountKey{"test"},
		},
		{
			testName:      "developer update msg",
			msg:           NewDeveloperUpdateMsg("test", "", "", ""),
			expectSigners: []types.AccountKey{"test"},
		},
		{
			testName:      "developer revoke msg",
			msg:           NewDeveloperRevokeMsg("test"),
			expectSigners: []types.AccountKey{"test"},
		},
		{
			testName:      "issue ida msg",
			msg:           IDAIssueMsg{Username: "test"},
			expectSigners: []types.AccountKey{"test"},
		},
		{
			testName:      "mint ida msg",
			msg:           IDAMintMsg{Username: "test"},
			expectSigners: []types.AccountKey{"test"},
		},
		{
			testName:      "ida transfer msg",
			msg:           IDATransferMsg{App: "app", Signer: "signer"},
			expectSigners: []types.AccountKey{"signer"},
		},
		{
			testName:      "ida authorize msg",
			msg:           IDAAuthorizeMsg{Username: "user"},
			expectSigners: []types.AccountKey{"user"},
		},
		{
			testName:      "update affiliated msg",
			msg:           UpdateAffiliatedMsg{App: "app"},
			expectSigners: []types.AccountKey{"app"},
		},
	}

	for _, tc := range testCases {
		suite.Equal(len(tc.expectSigners), len(tc.msg.GetSigners()), "%s", tc.testName)
		for i, signer := range tc.msg.GetSigners() {
			suite.Require().Equal(tc.expectSigners[i], types.AccountKey(signer), "%s", tc.testName)
		}
	}
}
