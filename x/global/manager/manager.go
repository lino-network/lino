package manager

import (
	"fmt"
	"strconv"

	codec "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/param"
	linotypes "github.com/lino-network/lino/types"

	"github.com/lino-network/lino/utils"
	"github.com/lino-network/lino/x/global/model"
	"github.com/lino-network/lino/x/global/types"

	accmn "github.com/lino-network/lino/x/account/manager"
)

const (
	exportVersion = 2
	importVersion = 2
)

// GlobalManager - a event manager module, it schedules event, execute events
// and store errors of executed events.
type GlobalManager struct {
	storage     model.GlobalStorage
	paramHolder param.ParamKeeper

	// events
	hourly  linotypes.BCEventExec
	daily   linotypes.BCEventExec
	monthly linotypes.BCEventExec
	yearly  linotypes.BCEventExec
}

// NewGlobalManager - return the global manager
func NewGlobalManager(key sdk.StoreKey, keeper param.ParamKeeper, cdc *codec.Codec,
	hourly linotypes.BCEventExec,
	daily linotypes.BCEventExec,
	monthly linotypes.BCEventExec,
	yearly linotypes.BCEventExec,
) GlobalManager {
	return GlobalManager{
		storage:     model.NewGlobalStorage(key, cdc),
		paramHolder: keeper,
		hourly:      hourly,
		daily:       daily,
		monthly:     monthly,
		yearly:      yearly,
	}
}

func (gm GlobalManager) InitGenesis(ctx sdk.Context) {
	// will be updated on the first OnBeginBlock.
	gm.storage.SetGlobalTime(ctx, &model.GlobalTime{
		ChainStartTime: ctx.BlockTime().Unix(),
		LastBlockTime:  ctx.BlockTime().Unix(),
		PastMinutes:    0,
	})
}

// OnBeginBlock - update internal time related fields and execute
// blockchain scheduled events.
func (gm GlobalManager) OnBeginBlock(ctx sdk.Context) {
	blockTime := ctx.BlockHeader().Time.Unix()
	globalTime := gm.storage.GetGlobalTime(ctx)
	if blockTime < globalTime.LastBlockTime {
		// our simulation tests do not follow tendermint's spec that
		// the BFT Time H2.Time > H1.Time, if H2 = H1 + 1.
		// precisely, we use a same time point all the time.
		// panic("Premise of BFT time is BROKEN")
		return
	}
	pastMinutes := globalTime.PastMinutes
	nowMinutes := (blockTime - globalTime.ChainStartTime) / 60
	for next := pastMinutes + 1; next <= nowMinutes; next++ {
		gm.execBCEventsAt(ctx, next)
	}
	globalTime.PastMinutes = nowMinutes
	gm.storage.SetGlobalTime(ctx, globalTime)
}

// execBCEventsAt - execute blockchain events.
func (gm GlobalManager) execBCEventsAt(ctx sdk.Context, pastMinutes int64) {
	if pastMinutes%60 == 0 && gm.hourly != nil {
		gm.appendBCErr(ctx, gm.hourly(ctx)...)
	}
	if pastMinutes%linotypes.MinutesPerDay == 0 && gm.daily != nil {
		gm.appendBCErr(ctx, gm.daily(ctx)...)
	}
	if pastMinutes%linotypes.MinutesPerMonth == 0 && gm.monthly != nil {
		gm.appendBCErr(ctx, gm.monthly(ctx)...)
	}
	if pastMinutes%linotypes.MinutesPerYear == 0 && gm.yearly != nil {
		gm.appendBCErr(ctx, gm.yearly(ctx)...)
	}
}

func (gm GlobalManager) appendBCErr(ctx sdk.Context, newErrs ...linotypes.BCEventErr) {
	errs := gm.storage.GetBCErrors(ctx)
	for _, e := range newErrs {
		ctx.Logger().Error(fmt.Sprintf("eventErr: %+v", e))
		errs = append(errs, e)
	}
	gm.storage.SetBCErrors(ctx, errs)
}

// OnEndBlock - update last block time.
func (gm GlobalManager) OnEndBlock(ctx sdk.Context) {
	globalTime := gm.storage.GetGlobalTime(ctx)
	globalTime.LastBlockTime = ctx.BlockHeader().Time.Unix()
	gm.storage.SetGlobalTime(ctx, globalTime)
}

func (gm GlobalManager) RegisterEventAtTime(ctx sdk.Context, unixTime int64, event linotypes.Event) sdk.Error {
	// XXX(yumin): events are executed at begin block, but not include
	// the current time. So event registered at this block time will be executed
	// in the next block.
	if unixTime < ctx.BlockHeader().Time.Unix() {
		return types.ErrRegisterExpiredEvent(unixTime)
	}

	// see if event is allowed or not.
	if !gm.storage.CanEncode(event) {
		return types.ErrRegisterInvalidEvent()
	}

	eventList := gm.storage.GetTimeEventList(ctx, unixTime)
	eventList.Events = append(eventList.Events, event)
	gm.storage.SetTimeEventList(ctx, unixTime, eventList)
	return nil
}

func (gm GlobalManager) runEventIsolated(ctx sdk.Context, exec linotypes.EventExec, event linotypes.Event) sdk.Error {
	cachedCtx, write := ctx.CacheContext()
	err := exec(cachedCtx, event)
	if err == nil {
		write()
		return nil
	}
	return err
}

// ExecuteEvents - execute events, log errors to storage, up to current time (exclusively).
func (gm GlobalManager) ExecuteEvents(ctx sdk.Context, exec linotypes.EventExec) {
	currentTime := ctx.BlockTime().Unix()
	lastBlockTime := gm.storage.GetGlobalTime(ctx).LastBlockTime
	for i := lastBlockTime; i < currentTime; i++ {
		events := gm.storage.GetTimeEventList(ctx, i)
		for _, event := range events.Events {
			err := gm.runEventIsolated(ctx, exec, event)
			if err != nil {
				ctx.Logger().Error(fmt.Sprintf(
					"ExecEventErr: %+v, code: %d", event, err.Code()))
				errs := gm.storage.GetEventErrors(ctx)
				errs = append(errs, model.EventError{
					Time:    i,
					Event:   event,
					ErrCode: err.Code(),
				})
				gm.storage.SetEventErrors(ctx, errs)
			}
		}
		gm.storage.RemoveTimeEventList(ctx, i)
	}

	// upgrade-3 unlock all.
	if ctx.BlockHeight() == linotypes.Upgrade5Update3 {
		gm.execFutureEvents(ctx, exec, func(ts int64, event linotypes.Event) bool {
			if ts < currentTime {
				return false
			}
			switch event.(type) {
			case accmn.ReturnCoinEvent:
				return true
			default:
				return false
			}
		})
	}
}

// resolveFutureEvents does not return err. Events are executed in an isolated env,
// so it's fine to ignore errors but leave them in the store.
func (gm GlobalManager) execFutureEvents(
	ctx sdk.Context, exec linotypes.EventExec,
	filter func(ts int64, event linotypes.Event) bool) {
	eventLists := make(map[int64]*linotypes.TimeEventList)
	store := gm.storage.PartialStoreMap(ctx)
	// change store will invalidate the iterator, so copy first.
	store[string(model.TimeEventListSubStore)].Iterate(func(key []byte, val interface{}) bool {
		ts, err := strconv.ParseInt(string(key), 10, 64)
		if err != nil {
			return false
		}
		eventLists[ts] = val.(*linotypes.TimeEventList)
		return false
	})

	for ts, eventList := range eventLists {
		if eventList == nil {
			continue
		}
		left := linotypes.TimeEventList{}
		for _, event := range eventList.Events {
			if filter(ts, event) {
				err := gm.runEventIsolated(ctx, exec, event)
				if err != nil {
					left.Events = append(left.Events, event)
				}
			} else {
				left.Events = append(left.Events, event)
			}
		}
		if len(left.Events) == 0 {
			gm.storage.RemoveTimeEventList(ctx, ts)
		} else {
			gm.storage.SetTimeEventList(ctx, ts, &left)
		}
	}
}

// GetLastBlockTime - get last block time from KVStore
func (gm GlobalManager) GetLastBlockTime(ctx sdk.Context) int64 {
	return gm.storage.GetGlobalTime(ctx).LastBlockTime
}

// GetPastDay - get start time from KVStore to calculate past day
func (gm GlobalManager) GetPastDay(ctx sdk.Context, unixTime int64) int64 {
	globalTime := gm.storage.GetGlobalTime(ctx)
	pastDay := (unixTime - globalTime.ChainStartTime) / (3600 * 24)
	if pastDay < 0 {
		return 0
	}
	return pastDay
}

func (gm GlobalManager) GetBCEventErrors(ctx sdk.Context) []linotypes.BCEventErr {
	return gm.storage.GetBCErrors(ctx)
}

func (gm GlobalManager) GetEventErrors(ctx sdk.Context) []model.EventError {
	return gm.storage.GetEventErrors(ctx)
}

func (gm GlobalManager) GetGlobalTime(ctx sdk.Context) model.GlobalTime {
	return *gm.storage.GetGlobalTime(ctx)
}

func (gm GlobalManager) ExportToFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error {
	state := &model.GlobalTablesIR{
		Version: exportVersion,
	}
	storeMap := gm.storage.PartialStoreMap(ctx)

	// export events
	storeMap[string(model.TimeEventListSubStore)].Iterate(func(key []byte, val interface{}) bool {
		ts, err := strconv.ParseInt(string(key), 10, 64)
		if err != nil {
			panic(err)
		}
		events := val.(*linotypes.TimeEventList)
		state.GlobalTimeEventLists = append(state.GlobalTimeEventLists, model.GlobalTimeEventsIR{
			UnixTime:      ts,
			TimeEventList: *events,
		})
		return false
	})

	globalt := gm.storage.GetGlobalTime(ctx)
	state.Time = model.GlobalTimeIR(*globalt)

	// errors are not export, because we are performing an upgrade, why not fix the errors?
	// EventErrorSubStore
	// BCErrorSubStore

	return utils.Save(filepath, cdc, state)
}

func (gm GlobalManager) ImportFromFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error {
	rst, err := utils.Load(filepath, cdc, func() interface{} { return &model.GlobalTablesIR{} })
	if err != nil {
		return err
	}
	table := rst.(*model.GlobalTablesIR)

	if table.Version != importVersion {
		return fmt.Errorf("unsupported import version: %d", table.Version)
	}

	// import events
	for _, v := range table.GlobalTimeEventLists {
		gm.storage.SetTimeEventList(ctx, v.UnixTime, &v.TimeEventList)
	}

	t := model.GlobalTime(table.Time)
	gm.storage.SetGlobalTime(ctx, &t)
	return nil
}
