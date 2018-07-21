package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// See https://github.com/cosmos/cosmos-sdk/issues/766
	LinoErrorCodeSpace = 11

	// Lino common errors reserve 100 ~ 149
	CodeInvalidUsername     sdk.CodeType = 100
	CodeAccountNotFound     sdk.CodeType = 101
	CodeFailedToMarshal     sdk.CodeType = 102
	CodeFailedToUnmarshal   sdk.CodeType = 103
	CodeIllegalWithdraw     sdk.CodeType = 104
	CodeInsufficientDeposit sdk.CodeType = 105
	CodeInvalidCoin         sdk.CodeType = 106
	CodePostNotFound        sdk.CodeType = 107
	CodeDeveloperNotFound   sdk.CodeType = 108
	CodeInvalidCoins        sdk.CodeType = 109

	// Lino authenticate errors reserve 150 ~ 199
	CodeIncorrectStdTxType   sdk.CodeType = 150
	CodeNoSignatures         sdk.CodeType = 151
	CodeUnknownMsgType       sdk.CodeType = 152
	CodeWrongNumberOfSigners sdk.CodeType = 153
	CodeInvalidSequence      sdk.CodeType = 154
	CodeUnverifiedBytes      sdk.CodeType = 155

	// ABCI Response Codes
	CodeGenesisFailed sdk.CodeType = 200

	// // Lino register handler errors reserve 300 ~ 309.
	// CodeAccRegisterFailed sdk.CodeType = 302
	// CodeUsernameNotFound  sdk.CodeType = 303

	// Lino account errors reserve 300 ~ 399
	CodeRewardNotFound                     sdk.CodeType = 300
	CodeAccountMetaNotFound                sdk.CodeType = 301
	CodeAccountInfoNotFound                sdk.CodeType = 302
	CodeAccountBankNotFound                sdk.CodeType = 303
	CodePendingStakeQueueNotFound          sdk.CodeType = 304
	CodeGrantPubKeyNotFound                sdk.CodeType = 305
	CodeFailedToMarshalAccountInfo         sdk.CodeType = 306
	CodeFailedToMarshalAccountBank         sdk.CodeType = 307
	CodeFailedToMarshalAccountMeta         sdk.CodeType = 308
	CodeFailedToMarshalFollowerMeta        sdk.CodeType = 309
	CodeFailedToMarshalFollowingMeta       sdk.CodeType = 310
	CodeFailedToMarshalReward              sdk.CodeType = 311
	CodeFailedToMarshalPendingStakeQueue   sdk.CodeType = 312
	CodeFailedToMarshalGrantPubKey         sdk.CodeType = 313
	CodeFailedToMarshalRelationship        sdk.CodeType = 314
	CodeFailedToMarshalBalanceHistory      sdk.CodeType = 315
	CodeFailedToUnmarshalAccountInfo       sdk.CodeType = 316
	CodeFailedToUnmarshalAccountBank       sdk.CodeType = 317
	CodeFailedToUnmarshalAccountMeta       sdk.CodeType = 318
	CodeFailedToUnmarshalReward            sdk.CodeType = 319
	CodeFailedToUnmarshalPendingStakeQueue sdk.CodeType = 320
	CodeFailedToUnmarshalGrantPubKey       sdk.CodeType = 321
	CodeFailedToUnmarshalRelationship      sdk.CodeType = 322
	CodeFailedToUnmarshalBalanceHistory    sdk.CodeType = 323
	CodeFolloweeNotFound                   sdk.CodeType = 324
	CodeFollowerNotFound                   sdk.CodeType = 325
	CodeReceiverNotFound                   sdk.CodeType = 326
	CodeSenderNotFound                     sdk.CodeType = 327
	CodeReferrerNotFound                   sdk.CodeType = 328
	CodeAddSavingCoinWithFullStake         sdk.CodeType = 329
	CodeAddSavingCoin                      sdk.CodeType = 330
	CodeInvalidMemo                        sdk.CodeType = 331
	CodeInvalidJSONMeta                    sdk.CodeType = 332
	CodeCheckRecoveryKey                   sdk.CodeType = 333
	CodeCheckTransactionKey                sdk.CodeType = 334
	CodeCheckGrantMicropaymentKey          sdk.CodeType = 335
	CodeCheckGrantPostKey                  sdk.CodeType = 336
	CodeCheckAuthenticatePubKeyOwner       sdk.CodeType = 337
	CodeGrantKeyExpired                    sdk.CodeType = 338
	CodeGrantKeyNoLeftTimes                sdk.CodeType = 339
	CodeGrantKeyMismatch                   sdk.CodeType = 340
	CodeMicropaymentGrantKeyMismatch       sdk.CodeType = 341
	CodePostGrantKeyMismatch               sdk.CodeType = 342
	CodeGetRecoveryKey                     sdk.CodeType = 343
	CodeGetTransactionKey                  sdk.CodeType = 344
	CodeGetMicropaymentKey                 sdk.CodeType = 345
	CodeGetPostKey                         sdk.CodeType = 346
	CodeGetSavingFromBank                  sdk.CodeType = 347
	CodeGetSequence                        sdk.CodeType = 348
	CodeGetLastReportOrUpvoteAt            sdk.CodeType = 349
	CodeUpdateLastReportOrUpvoteAt         sdk.CodeType = 350
	CodeGetFrozenMoneyList                 sdk.CodeType = 351
	CodeIncreaseSequenceByOne              sdk.CodeType = 352
	CodeGrantTimesExceedsLimitation        sdk.CodeType = 353
	CodeUnsupportGrantLevel                sdk.CodeType = 354
	CodeRevokePermissionLevelMismatch      sdk.CodeType = 355
	CodeCheckUserTPSCapacity               sdk.CodeType = 356
	CodeAccountTPSCapacityNotEnough        sdk.CodeType = 357
	CodeAccountSavingCoinNotEnough         sdk.CodeType = 358
	CodeAccountAlreadyExists               sdk.CodeType = 359
	CodeRegisterFeeInsufficient            sdk.CodeType = 360
	CodeFailedToMarshalRewardHistory       sdk.CodeType = 361
	CodeFailedToUnmarshalRewardHistory     sdk.CodeType = 362

	// Lino post errors reserve 400 ~ 499
	CodePostMetaNotFound                     sdk.CodeType = 400
	CodePostLikeNotFound                     sdk.CodeType = 401
	CodePostReportOrUpvoteNotFound           sdk.CodeType = 402
	CodePostCommentNotFound                  sdk.CodeType = 403
	CodePostViewNotFound                     sdk.CodeType = 404
	CodePostDonationNotFound                 sdk.CodeType = 405
	CodeFailedToMarshalPostInfo              sdk.CodeType = 406
	CodeFailedToMarshalPostMeta              sdk.CodeType = 407
	CodeFailedToMarshalPostLike              sdk.CodeType = 408
	CodeFailedToMarshalPostReportOrUpvote    sdk.CodeType = 409
	CodeFailedToMarshalPostComment           sdk.CodeType = 410
	CodeFailedToMarshalPostView              sdk.CodeType = 411
	CodeFailedToMarshalPostDonations         sdk.CodeType = 412
	CodeFailedToUnmarshalPostInfo            sdk.CodeType = 413
	CodeFailedToUnmarshalPostMeta            sdk.CodeType = 414
	CodeFailedToUnmarshalPostLike            sdk.CodeType = 415
	CodeFailedToUnmarshalPostReportOrUpvote  sdk.CodeType = 416
	CodeFailedToUnmarshalPostComment         sdk.CodeType = 417
	CodeFailedToUnmarshalPostView            sdk.CodeType = 418
	CodeFailedToUnmarshalPostDonations       sdk.CodeType = 419
	CodePostAlreadyExist                     sdk.CodeType = 420
	CodeInvalidPostRedistributionSplitRate   sdk.CodeType = 421
	CodeDonatePostIsDeleted                  sdk.CodeType = 422
	CodeCannotDonateToSelf                   sdk.CodeType = 423
	CodeMicropaymentExceedsLimitation        sdk.CodeType = 424
	CodeProcessSourceDonation                sdk.CodeType = 425
	CodeProcessDonation                      sdk.CodeType = 426
	CodeUpdatePostIsDeleted                  sdk.CodeType = 427
	CodeReportOrUpvoteTooOften               sdk.CodeType = 428
	CodeReportOrUpvoteAlreadyExist           sdk.CodeType = 429
	CodeNoPostID                             sdk.CodeType = 430
	CodePostIDTooLong                        sdk.CodeType = 431
	CodeNoAuthor                             sdk.CodeType = 432
	CodeNoUsername                           sdk.CodeType = 433
	CodeCommentAndRepostConflict             sdk.CodeType = 434
	CodePostTitleExceedMaxLength             sdk.CodeType = 435
	CodePostContentExceedMaxLength           sdk.CodeType = 436
	CodeRedistributionSplitRateLengthTooLong sdk.CodeType = 437
	CodeIdentifierLengthTooLong              sdk.CodeType = 438
	CodeURLLengthTooLong                     sdk.CodeType = 439
	CodeTooManyURL                           sdk.CodeType = 440
	CodePostLikeNoUsername                   sdk.CodeType = 441
	CodePostLikeWeightOverflow               sdk.CodeType = 442
	CodePostLikeInvalidTarget                sdk.CodeType = 443
	CodeInvalidTarget                        sdk.CodeType = 444
	CodeCreatePostSourceInvalid              sdk.CodeType = 445
	CodeGetSourcePost                        sdk.CodeType = 446

	// Lino validator errors reserve 500 ~ 599
	CodeValidatorNotFound              sdk.CodeType = 500
	CodeValidatorListNotFound          sdk.CodeType = 501
	CodeFailedToMarshalValidator       sdk.CodeType = 502
	CodeFailedToMarshalValidatorList   sdk.CodeType = 503
	CodeFailedToUnmarshalValidator     sdk.CodeType = 504
	CodeFailedToUnmarshalValidatorList sdk.CodeType = 505
	CodeUnbalancedAccount              sdk.CodeType = 506
	CodeValidatorPubKeyAlreadyExist    sdk.CodeType = 507

	// Lino global errors reserve 600 ~ 699
	CodeInfraInflationCoinConversion     sdk.CodeType = 600
	CodeContentCreatorCoinConversion     sdk.CodeType = 601
	CodeDeveloperCoinConversion          sdk.CodeType = 602
	CodeValidatorCoinConversion          sdk.CodeType = 603
	CodeGlobalMetaNotFound               sdk.CodeType = 604
	CodeInflationPoolNotFound            sdk.CodeType = 605
	CodeGlobalConsumptionMetaNotFound    sdk.CodeType = 606
	CodeGlobalTPSNotFound                sdk.CodeType = 607
	CodeFailedToMarshalTimeEventList     sdk.CodeType = 608
	CodeFailedToMarshalGlobalMeta        sdk.CodeType = 609
	CodeFailedToMarshalInflationPoll     sdk.CodeType = 610
	CodeFailedToMarshalConsumptionMeta   sdk.CodeType = 611
	CodeFailedToMarshalTPS               sdk.CodeType = 612
	CodeFailedToUnmarshalTimeEventList   sdk.CodeType = 613
	CodeFailedToUnmarshalGlobalMeta      sdk.CodeType = 614
	CodeFailedToUnmarshalInflationPool   sdk.CodeType = 615
	CodeFailedToUnmarshalConsumptionMeta sdk.CodeType = 616
	CodeFailedToUnmarshalTPS             sdk.CodeType = 617
	CodeRegisterExpiredEvent             sdk.CodeType = 618

	// Vote errors reserve 700 ~ 799
	CodeVoterNotFound                  sdk.CodeType = 700
	CodeVoteNotFound                   sdk.CodeType = 701
	CodeReferenceListNotFound          sdk.CodeType = 702
	CodeDelegationNotFound             sdk.CodeType = 703
	CodeFailedToMarshalVoter           sdk.CodeType = 704
	CodeFailedToMarshalVote            sdk.CodeType = 705
	CodeFailedToMarshalDelegation      sdk.CodeType = 706
	CodeFailedToMarshalReferenceList   sdk.CodeType = 707
	CodeFailedToUnmarshalVoter         sdk.CodeType = 708
	CodeFailedToUnmarshalVote          sdk.CodeType = 709
	CodeFailedToUnmarshalDelegation    sdk.CodeType = 710
	CodeFailedToUnmarshalReferenceList sdk.CodeType = 711
	CodeValidatorCannotRevoke          sdk.CodeType = 712
	CodeVoteAlreadyExist               sdk.CodeType = 713

	// Lino infra errors reserve 800 ~ 899
	CodeInfraProviderNotFound              sdk.CodeType = 800
	CodeInfraProviderListNotFound          sdk.CodeType = 801
	CodeFailedToMarshalInfraProvider       sdk.CodeType = 802
	CodeFailedToMarshalInfraProviderList   sdk.CodeType = 803
	CodeFailedToUnmarshalInfraProvider     sdk.CodeType = 804
	CodeFailedToUnmarshalInfraProviderList sdk.CodeType = 805
	CodeInvalidUsage                       sdk.CodeType = 806

	// Lino developer errors reserve 900 ~ 999
	CodeDeveloperListNotFound          sdk.CodeType = 900
	CodeFailedToMarshalDeveloper       sdk.CodeType = 901
	CodeFailedToMarshalDeveloperList   sdk.CodeType = 902
	CodeFailedToUnmarshalDeveloper     sdk.CodeType = 903
	CodeFailedToUnmarshalDeveloperList sdk.CodeType = 904
	CodeDeveloperAlreadyExist          sdk.CodeType = 905
	CodeInsufficientDeveloperDeposit   sdk.CodeType = 906
	CodeInvalidAuthenticateApp         sdk.CodeType = 907
	CodeInvalidValidityPeriod          sdk.CodeType = 908
	CodeGrantPermissionTooHigh         sdk.CodeType = 909
	CodeInvalidGrantTimes              sdk.CodeType = 910

	// Param errors reserve 1000 ~ 1099
	CodeParamHolderGenesisError                       sdk.CodeType = 1000
	CodeDeveloperParamNotFound                        sdk.CodeType = 1001
	CodeValidatorParamNotFound                        sdk.CodeType = 1002
	CodeCoinDayParamNotFound                          sdk.CodeType = 1003
	CodeBandwidthParamNotFound                        sdk.CodeType = 1004
	CodeAccountParamNotFound                          sdk.CodeType = 1005
	CodeVoteParamNotFound                             sdk.CodeType = 1006
	CodeProposalParamNotFound                         sdk.CodeType = 1007
	CodeGlobalAllocationParamNotFound                 sdk.CodeType = 1008
	CodeInfraAllocationParamNotFound                  sdk.CodeType = 1009
	CodePostParamNotFound                             sdk.CodeType = 1010
	CodeInvalidaParameter                             sdk.CodeType = 1011
	CodeEvaluateOfContentValueParamNotFound           sdk.CodeType = 1012
	CodeFailedToUnmarshalGlobalAllocationParam        sdk.CodeType = 1013
	CodeFailedToUnmarshalPostParam                    sdk.CodeType = 1014
	CodeFailedToUnmarshalValidatorParam               sdk.CodeType = 1015
	CodeFailedToUnmarshalEvaluateOfContentValueParam  sdk.CodeType = 1016
	CodeFailedToUnmarshalInfraInternalAllocationParam sdk.CodeType = 1017
	CodeFailedToUnmarshalDeveloperParam               sdk.CodeType = 1018
	CodeFailedToUnmarshalVoteParam                    sdk.CodeType = 1019
	CodeFailedToUnmarshalProposalParam                sdk.CodeType = 1020
	CodeFailedToUnmarshalCoinDayParam                 sdk.CodeType = 1021
	CodeFailedToUnmarshalBandwidthParam               sdk.CodeType = 1022
	CodeFailedToUnmarshalAccountParam                 sdk.CodeType = 1023
	CodeFailedToMarshalGlobalAllocationParam          sdk.CodeType = 1024
	CodeFailedToMarshalPostParam                      sdk.CodeType = 1025
	CodeFailedToMarshalValidatorParam                 sdk.CodeType = 1026
	CodeFailedToMarshalEvaluateOfContentValueParam    sdk.CodeType = 1027
	CodeFailedToMarshalInfraInternalAllocationParam   sdk.CodeType = 1028
	CodeFailedToMarshalDeveloperParam                 sdk.CodeType = 1029
	CodeFailedToMarshalVoteParam                      sdk.CodeType = 1030
	CodeFailedToMarshalProposalParam                  sdk.CodeType = 1031
	CodeFailedToMarshalCoinDayParam                   sdk.CodeType = 1032
	CodeFailedToMarshalBandwidthParam                 sdk.CodeType = 1033
	CodeFailedToMarshalAccountParam                   sdk.CodeType = 1034

	// Proposal errors reserve 1100 ~ 1199
	CodeOngoingProposalNotFound         sdk.CodeType = 1100
	CodeCensorshipPostNotFound          sdk.CodeType = 1101
	CodeProposalNotFound                sdk.CodeType = 1102
	CodeProposalListNotFound            sdk.CodeType = 1103
	CodeNextProposalIDNotFound          sdk.CodeType = 1104
	CodeFailedToMarshalProposal         sdk.CodeType = 1105
	CodeFailedToMarshalProposalList     sdk.CodeType = 1106
	CodeFailedToMarshalNextProposalID   sdk.CodeType = 1107
	CodeFailedToUnmarshalProposal       sdk.CodeType = 1108
	CodeFailedToUnmarshalProposalList   sdk.CodeType = 1109
	CodeFailedToUnmarshalNextProposalID sdk.CodeType = 1110
	CodeCensorshipPostIsDeleted         sdk.CodeType = 1111
	CodeNotOngoingProposal              sdk.CodeType = 1112
	CodeIncorrectProposalType           sdk.CodeType = 1113
	CodeInvalidPermlink                 sdk.CodeType = 1114
	CodeInvalidLink                     sdk.CodeType = 1115
	CodeIllegalParameter                sdk.CodeType = 1116
)
