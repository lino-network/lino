package account

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/lino-network/lino/app"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	acc "github.com/lino-network/lino/x/account"
	accModel "github.com/lino-network/lino/x/account/model"
	developerModel "github.com/lino-network/lino/x/developer/model"
	global "github.com/lino-network/lino/x/global"
	post "github.com/lino-network/lino/x/post"
	val "github.com/lino-network/lino/x/validator"
	vote "github.com/lino-network/lino/x/vote"
	voteModel "github.com/lino-network/lino/x/vote/model"
	abci "github.com/tendermint/tendermint/abci/types"
	crypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

var (
	GenesisAccount SimAcc
	AccountList    []*SimAcc = []*SimAcc{}
	VoterList      []*SimAcc = []*SimAcc{}
	StatisticParam Statistic

	RichStandard = types.NewCoinFromInt64(100000 * types.Decimals)

	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	PostList []PostStruct = []PostStruct{}

	percentageUserChargeAccount  = 4
	percentageUserTransferMoney  = 2
	percentageUserDepositAsVoter = 2
	percentageOfComment          = 10
	percentageUserPost           = 2
	percentageUserDonation       = 2

	TotalTransferCoin = types.NewCoinFromInt64(0)
	TotalDonationCoin = types.NewCoinFromInt64(0)
)

type PostStruct struct {
	Author types.AccountKey
	PostID string
}

type DonationStruct struct {
	User   types.AccountKey
	Author types.AccountKey
	PostID string
	Amount types.Coin
	Time   int64
}

type Statistic struct {
	CheckTxFaildTimes        int64
	DeliverTxFailedTimes     int64
	TransferToAccTimes       int64
	TransferToAddrTimes      int64
	ChargeAccountTimes       int64
	VoterDepositSuccessTimes int64
	VoterDepositFailedTimes  int64
	PostTimes                int64
	CommentTimes             int64
	DonateFailedTimes        int64
	DonateSuccessTimes       int64
}

func (st Statistic) String() string {
	return fmt.Sprintf(`CheckTxFaildTimes: %v, DeliverTxFailedTimes:%v
		TransferToAccTimes:%v, TransferToAddrTimes:%v
		ChargeAccountTimes:%v,
		VoterDepositSuccessTimes:%v, VoterDepositFailedTimes:%v,
		Post: %v, Comment: %v, DonateFailedTimes: %v, DonateSuccessTimes: %v`,
		st.CheckTxFaildTimes, st.DeliverTxFailedTimes,
		st.TransferToAccTimes, st.TransferToAddrTimes,
		st.ChargeAccountTimes,
		st.VoterDepositSuccessTimes, st.VoterDepositFailedTimes,
		st.PostTimes, st.CommentTimes, st.DonateFailedTimes, st.DonateSuccessTimes)
}

type SimAcc struct {
	AccountID        int64
	ExpectCoin       types.Coin
	ResetPrivKey     secp256k1.PrivKeySecp256k1
	TxPrivKey        secp256k1.PrivKeySecp256k1
	AppPrivKey       secp256k1.PrivKeySecp256k1
	Username         types.AccountKey
	Sequence         int64
	IsValidator      bool
	ValidatorPrivKey ed25519.PrivKeyEd25519
	IsVoter          bool
	IsInfraProvider  bool
	IsDeveloper      bool
	DelegateList     []*SimAcc
	ActualReward     types.Coin
}

func NewSimAcc(accountName string, initCoin types.Coin) *SimAcc {
	return &SimAcc{
		AccountID:    int64(len(AccountList)),
		ExpectCoin:   initCoin,
		ResetPrivKey: secp256k1.GenPrivKey(),
		TxPrivKey:    secp256k1.GenPrivKey(),
		AppPrivKey:   secp256k1.GenPrivKey(),
		Username:     types.AccountKey(accountName),
		Sequence:     0,
		IsValidator:  false,
		IsVoter:      false,
		DelegateList: []*SimAcc{},
		ActualReward: types.NewCoinFromInt64(0),
	}
}

func (simAcc *SimAcc) Action(lb *app.LinoBlockchain) {
	if rand.Intn(len(AccountList)) > len(AccountList)/100 {
		// only 1% of total user will take action at each block
		return
	}
	if rand.Intn(100) < percentageUserChargeAccount {
		StatisticParam.ChargeAccountTimes++
		chargeMoney := 10 + rand.Intn(1000)
		GenesisAccount.Transfer(lb, AccountList[rand.Intn(len(AccountList))], strconv.Itoa(chargeMoney))
	}
	if rand.Intn(100) < percentageUserTransferMoney {
		transferMoneyOut := 10 + rand.Intn(1000)
		simAcc.Transfer(lb, AccountList[rand.Intn(len(AccountList))], strconv.Itoa(transferMoneyOut))
	}

	if rand.Intn(100) < percentageUserDepositAsVoter {
		if simAcc.ExpectCoin.IsGT(RichStandard) {
			simAcc.DepositToBeVoter(lb)
		}
	}

	if rand.Intn(100) < percentageUserPost {
		simAcc.PublishPost(lb)
	}
}

func (simAcc *SimAcc) Donation(lb *app.LinoBlockchain, atTime int64) *post.RewardEvent {
	if rand.Intn(len(AccountList)) > len(AccountList)/1000 {
		// only 1% of total user will take action at each block
		return nil
	}
	amount := rand.Int63n(1000)
	if len(PostList) == 0 {
		return nil
	}
	donateTo := PostList[rand.Intn(len(PostList))]
	if rand.Intn(100) < percentageUserDonation {
		ph := param.NewParamHolder(lb.CapKeyParamStore)
		accManager := acc.NewAccountManager(lb.CapKeyAccountStore, ph)
		postManager := post.NewPostManager(lb.CapKeyPostStore, ph)
		globalManager := global.NewGlobalManager(lb.CapKeyGlobalStore, ph)
		ctx := lb.BaseApp.NewContext(false, abci.Header{Time: time.Unix(atTime, 0)})
		numOfConsumptionOnAuthor, _ := accManager.GetDonationRelationship(ctx, simAcc.Username, donateTo.Author)
		created, totalReward, _ := postManager.GetCreatedTimeAndReward(ctx, types.GetPermlink(donateTo.Author, donateTo.PostID))
		if atTime-created == 120 {
			return nil
		}
		//fmt.Println("donate from", simAcc.Username, " to ", donateTo.Author)
		donateMsg := post.NewDonateMsg(string(simAcc.Username),
			strconv.Itoa(int(amount)), string(donateTo.Author), donateTo.PostID, "", RandStringRunes(10))
		if !broadcastMsg(lb, donateMsg, simAcc.Sequence, simAcc.TxPrivKey) {
			StatisticParam.DonateFailedTimes++
			return nil
		}
		TotalDonationCoin = TotalDonationCoin.Plus(types.NewCoinFromInt64(amount * types.Decimals))
		StatisticParam.DonateSuccessTimes++
		simAcc.Sequence++

		simAcc.ExpectCoin = simAcc.ExpectCoin.Minus(types.NewCoinFromInt64(amount * types.Decimals))

		for _, acc := range AccountList {
			if acc.Username == donateTo.Author {
				directDeposit := types.RatToCoin(
					types.NewCoinFromInt64(amount * types.Decimals).ToRat().Mul(sdk.NewRat(95, 100)))
				acc.ActualReward = acc.ActualReward.Plus(directDeposit)
				acc.ExpectCoin = acc.ExpectCoin.Plus(directDeposit)
				break
			}
		}
		evaluate, _ := globalManager.EvaluateConsumption(ctx, types.NewCoinFromInt64(amount*types.Decimals), numOfConsumptionOnAuthor, created, totalReward)
		fmt.Println("in sim evaluate:", evaluate, ", amount:", types.NewCoinFromInt64(amount*types.Decimals), numOfConsumptionOnAuthor, created, atTime, totalReward)
		friction := types.RatToCoin(
			types.NewCoinFromInt64(amount * types.Decimals).ToRat().Mul(sdk.NewRat(5, 100)))
		fmt.Println("donate to ", donateTo.Author)
		return &post.RewardEvent{
			PostAuthor: donateTo.Author,
			PostID:     donateTo.PostID,
			Consumer:   simAcc.Username,
			Evaluate:   evaluate,
			Original:   types.NewCoinFromInt64(amount * types.Decimals),
			Friction:   friction,
			FromApp:    "",
		}
	}
	return nil
}

func (simAcc *SimAcc) DepositToBeVoter(lb *app.LinoBlockchain) bool {
	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	voteParam, _ := ph.GetVoteParam(ctx)
	amount := voteParam.VoterMinDeposit.Plus(types.NewCoinFromInt64(rand.Int63n(1000)))
	voteMsg := vote.NewVoterDepositMsg(string(simAcc.Username), strconv.Itoa(int(amount.ToInt64())))
	IsNewVoter := false
	if !simAcc.IsVoter {
		IsNewVoter = true
	}
	if !broadcastMsg(lb, voteMsg, simAcc.Sequence, simAcc.TxPrivKey) {
		StatisticParam.VoterDepositFailedTimes++
		return false
	}
	StatisticParam.VoterDepositSuccessTimes++
	simAcc.Sequence++
	if IsNewVoter {
		VoterList = append(VoterList, simAcc)
		simAcc.IsVoter = true
	}
	return true
}

func (simAcc *SimAcc) VoteOtherPeople(lb *app.LinoBlockchain) bool {
	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	voteParam, _ := ph.GetVoteParam(ctx)
	amount := voteParam.VoterMinDeposit.Plus(types.NewCoinFromInt64(rand.Int63n(1000)))
	voteMsg := vote.NewVoterDepositMsg(string(simAcc.Username), strconv.Itoa(int(amount.ToInt64())))
	if !broadcastMsg(lb, voteMsg, simAcc.Sequence, simAcc.TxPrivKey) {
		StatisticParam.VoterDepositFailedTimes++
		return false
	}
	StatisticParam.VoterDepositSuccessTimes++
	simAcc.Sequence++
	return true
}

func (simAcc *SimAcc) Transfer(
	lb *app.LinoBlockchain, to *SimAcc, amount types.LNO) bool {
	transferMsg := acc.NewTransferMsg(
		string(simAcc.Username), string(to.Username), amount, "")
	//fmt.Println("transfer to:", to.Username)
	if !broadcastMsg(lb, transferMsg, simAcc.Sequence, simAcc.TxPrivKey) {
		return false
	}
	coin, _ := types.LinoToCoin(amount)
	TotalTransferCoin = TotalTransferCoin.Plus(coin)
	simAcc.ExpectCoin = simAcc.ExpectCoin.Minus(coin)
	simAcc.Sequence++
	StatisticParam.TransferToAccTimes++
	to.ExpectCoin = to.ExpectCoin.Plus(coin)
	return true
}

func broadcastMsg(lb *app.LinoBlockchain, msg sdk.Msg, seq int64, privKey crypto.PrivKey) bool {
	tx := genTx(msg, seq, privKey)
	res := lb.Simulate(tx)
	if res.Code != sdk.ABCICodeOK {
		StatisticParam.CheckTxFaildTimes += 1
		//fmt.Print("creat account failed:", res)
		return false
	}
	res = lb.Deliver(tx)
	if res.Code != sdk.ABCICodeOK {
		StatisticParam.DeliverTxFailedTimes += 1
		return false
	}
	return true
}

func (simAcc *SimAcc) PublishPost(lb *app.LinoBlockchain) bool {
	var isComment bool
	var postMsg post.CreatePostMsg
	if rand.Intn(100) < percentageOfComment && len(PostList) > 0 {
		parentPost := PostList[rand.Intn(len(PostList))]
		postMsg = post.NewCreatePostMsg(
			string(simAcc.Username), RandStringRunes(8), RandStringRunes(8),
			RandStringRunes(50), string(parentPost.Author), parentPost.PostID, "", "", "0", nil)
	}

	postMsg = post.NewCreatePostMsg(
		string(simAcc.Username), RandStringRunes(8), RandStringRunes(8),
		RandStringRunes(50), "", "", "", "", "0", nil)

	tx := genTx(postMsg, simAcc.Sequence, simAcc.AppPrivKey)
	res := lb.Simulate(tx)
	if res.Code != sdk.ABCICodeOK {
		return false
	}
	res = lb.Deliver(tx)
	if res.Code != sdk.ABCICodeOK {
		return false
	}
	simAcc.Sequence++
	PostList = append(PostList,
		PostStruct{postMsg.Author, postMsg.PostID})
	if isComment {
		StatisticParam.CommentTimes++
	} else {
		StatisticParam.PostTimes++
	}
	return true
}

func (simAcc *SimAcc) CheckSelfBalance(lb *app.LinoBlockchain) bool {
	totalCoin := types.NewCoinFromInt64(0)
	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	accManager := acc.NewAccountManager(lb.CapKeyAccountStore, ph)
	accStorage := accModel.NewAccountStorage(lb.CapKeyAccountStore)
	saving, _ := accManager.GetSavingFromBank(ctx, simAcc.Username)
	totalCoin = totalCoin.Plus(saving)

	voterManager := vote.NewVoteManager(lb.CapKeyVoteStore, ph)
	voterStorage := voteModel.NewVoteStorage(lb.CapKeyVoteStore)
	validatorManager := val.NewValidatorManager(lb.CapKeyValStore, ph)
	developerStorage := developerModel.NewDeveloperStorage(lb.CapKeyDeveloperStore)
	if simAcc.IsVoter {
		coin, err := voterManager.GetVoterDeposit(ctx, simAcc.Username)
		if err != nil {
			panic(err)
		}
		totalCoin = totalCoin.Plus(coin)
	}
	if simAcc.IsValidator {
		coin, err := validatorManager.GetValidatorDeposit(ctx, simAcc.Username)
		if err != nil {
			panic(err)
		}
		totalCoin = totalCoin.Plus(coin)
	}
	if simAcc.IsDeveloper {
		developer, err := developerStorage.GetDeveloper(ctx, simAcc.Username)
		if err != nil {
			panic(err)
		}
		totalCoin = totalCoin.Plus(developer.Deposit)
	}
	for _, delegateVoter := range simAcc.DelegateList {
		delegation, _ := voterStorage.GetDelegation(ctx, delegateVoter.Username, simAcc.Username)
		totalCoin = totalCoin.Plus(delegation.Amount)
	}
	if !totalCoin.IsEqual(simAcc.ExpectCoin) {
		fmt.Println("total:", totalCoin, ", expect:", simAcc.ExpectCoin)
		return false
	}
	reward, _ := accStorage.GetReward(ctx, simAcc.Username)
	if !reward.TotalIncome.IsEqual(simAcc.ActualReward) {
		fmt.Println("actualReward:", reward.TotalIncome, ", expect reward:", simAcc.ActualReward)
		return false
	}
	return true
}

func CreateAccount(accountName string, lb *app.LinoBlockchain, numOfLino string) bool {
	coin, _ := types.LinoToCoin(numOfLino)
	simAcc := NewSimAcc(accountName, coin.Minus(types.NewCoinFromInt64(1*types.Decimals)))

	registerMsg := acc.NewRegisterMsg(string(GenesisAccount.Username), accountName, numOfLino, simAcc.ResetPrivKey.PubKey(), simAcc.TxPrivKey.PubKey(), simAcc.AppPrivKey.PubKey())
	tx := genTx(registerMsg, GenesisAccount.Sequence, GenesisAccount.TxPrivKey)
	res := lb.Simulate(tx)
	if res.Code != sdk.ABCICodeOK {
		fmt.Print("creat account failed:", res)
		return false
	}

	res = lb.Deliver(tx)
	if res.Code != sdk.ABCICodeOK {
		fmt.Print("creat account failed:", res)
		return false
	}
	GenesisAccount.Sequence++
	GenesisAccount.ExpectCoin = GenesisAccount.ExpectCoin.Minus(coin)
	AccountList = append(AccountList, simAcc)
	return true
}

func genTx(msg sdk.Msg, seq int64, priv crypto.PrivKey) auth.StdTx {
	bz, _ := priv.Sign(auth.StdSignBytes("Lino", 0, seq, auth.StdFee{}, []sdk.Msg{msg}, ""))
	sigs := []auth.StdSignature{{
		PubKey:    priv.PubKey(),
		Signature: bz,
		Sequence:  seq}}

	return auth.NewStdTx([]sdk.Msg{msg}, auth.StdFee{}, sigs, "")
}

// XXX reference the common declaration of this function
func subspace(prefix []byte) (start, end []byte) {
	end = make([]byte, len(prefix))
	copy(end, prefix)
	end[len(end)-1]++
	return prefix, end
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
