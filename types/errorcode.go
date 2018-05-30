package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ABCI Response Codes
	// Base SDK reserves 0 ~ 99.
	// Coin errors reserve 100 ~ 199.
	// Lino authentication errors reserve 200 ~ 299.
	// Lino register handler errors reserve 300 ~ 309.
	CodeInvalidUsername   sdk.CodeType = 301
	CodeAccRegisterFailed sdk.CodeType = 302
	CodeUsernameNotFound  sdk.CodeType = 303

	// Lino account handler errors reserve 310 ~ 399
	CodeAccountStorageFail sdk.CodeType = 310
	CodeAccountManagerFail sdk.CodeType = 311
	CodeAccountHandlerFail sdk.CodeType = 312
	CodeInvalidMsg         sdk.CodeType = 313
	CodeInvalidMemo        sdk.CodeType = 314

	// Lino post handler errors reserve 400 ~ 499
	CodePostMarshalError   sdk.CodeType = 400
	CodePostUnmarshalError sdk.CodeType = 401
	CodePostNotFound       sdk.CodeType = 402
	CodePostCreateError    sdk.CodeType = 403
	CodePostLikeError      sdk.CodeType = 404
	CodePostDonateError    sdk.CodeType = 405
	CodePostManagerError   sdk.CodeType = 406
	CodePostHandlerError   sdk.CodeType = 407
	CodePostMsgError       sdk.CodeType = 408
	CodePostStorageError   sdk.CodeType = 409

	// validator errors reserve 500 ~ 599
	CodeValidatorHandlerFailed sdk.CodeType = 500
	CodeValidatorManagerFailed sdk.CodeType = 501

	// Event errors reserve 600 ~ 699
	CodeGlobalStorageGenesisError sdk.CodeType = 600
	CodeGlobalStorageError        sdk.CodeType = 601
	CodeGlobalManagerError        sdk.CodeType = 602

	// Vote errors reserve 700 ~ 799
	CodeVoteHandlerFailed sdk.CodeType = 700
	CodeVoteManagerFailed sdk.CodeType = 701

	// Infra errors reserve 800 ~ 899
	CodeInfraProviderHandlerFailed sdk.CodeType = 800
	CodeInfraProviderManagerFailed sdk.CodeType = 801
	CodeInfraInvalidMsg            sdk.CodeType = 802

	// Developer errors reserve 900 ~ 999
	CodeDeveloperHandlerFailed sdk.CodeType = 900
	CodeDeveloperManagerFailed sdk.CodeType = 901

	// Param errors reserve 1000 ~ 1099
	CodeParamStoreError         sdk.CodeType = 1000
	CodeParamHolderGenesisError sdk.CodeType = 1001

	// Proposal errors reserve 1100 ~ 1199
	CodeProposalStoreError   sdk.CodeType = 1100
	CodeProposalManagerError sdk.CodeType = 1101
	CodeProposalEventError   sdk.CodeType = 1102
	CodeProposalHandlerError sdk.CodeType = 1103
	CodeProposalMsgError     sdk.CodeType = 1104
)
