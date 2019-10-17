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
)

const (
	exportVersion = 1
	importVersion = 1
)

// GlobalManager - a event manager module, it schedules event, execute events
// and store errors of executed events.
type GlobalManager struct {
	storage     model.GlobalStorage
	paramHolder param.ParamHolder

	// events
	hourly  types.BCEventExec
	daily   types.BCEventExec
	monthly types.BCEventExec
	yearly  types.BCEventExec
}

// NewGlobalManager - return the global manager
func NewGlobalManager(key sdk.StoreKey, holder param.ParamHolder, cdc *codec.Codec,
	hourly types.BCEventExec,
	daily types.BCEventExec,
	monthly types.BCEventExec,
	yearly types.BCEventExec,
) GlobalManager {
	return GlobalManager{
		storage:     model.NewGlobalStorage(key, cdc),
		paramHolder: holder,
		hourly:      hourly,
		daily:       daily,
		monthly:     monthly,
		yearly:      yearly,
	}
}

func (gm GlobalManager) InitGenesis(ctx sdk.Context) {
	// will be updated on the first OnBeginBlock.
	gm.storage.SetGlobalTime(ctx, &model.GlobalTime{})
}

// OnBeginBlock - update internal time related fields and execute
// blockchain scheduled events.
func (gm GlobalManager) OnBeginBlock(ctx sdk.Context) {
	blockTime := ctx.BlockHeader().Time.Unix()
	globalTime := gm.storage.GetGlobalTime(ctx)
	if globalTime.ChainStartTime == 0 {
		globalTime.ChainStartTime = blockTime
		globalTime.LastBlockTime = blockTime
	}

	pastMinutes := globalTime.PastMinutes
	nowMinutes := (blockTime - globalTime.ChainStartTime) / 60
	for next := pastMinutes + 1; next <= nowMinutes; next++ {
		gm.execBeginBlockEventsAt(ctx, next)
	}
	globalTime.PastMinutes = nowMinutes
	gm.storage.SetGlobalTime(ctx, globalTime)
}

func (gm GlobalManager) execBeginBlockEventsAt(ctx sdk.Context, pastMinutes int64) {
	if pastMinutes%60 == 0 && gm.hourly != nil {
		gm.hourly(ctx)
	}
	if pastMinutes%linotypes.MinutesPerDay == 0 && gm.daily != nil {
		gm.daily(ctx)
	}
	if pastMinutes%linotypes.MinutesPerMonth == 0 && gm.monthly != nil {
		gm.monthly(ctx)
	}
	if pastMinutes%linotypes.MinutesPerYear == 0 && gm.yearly != nil {
		gm.yearly(ctx)
	}
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

// ExecuteEvents - execute events, log errors to storage.
func (gm GlobalManager) ExecuteEvents(ctx sdk.Context, exec types.EventExec) {
	currentTime := ctx.BlockHeader().Time.Unix()
	lastBlockTime := gm.storage.GetGlobalTime(ctx).LastBlockTime
	for i := lastBlockTime; i < currentTime; i++ {
		events := gm.storage.GetTimeEventList(ctx, i)
		for _, event := range events.Events {
			cachedCtx, write := ctx.CacheContext()
			err := exec(cacheCtx, event)
			if err == nil {
				write()
			} else {
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
}

// GetTimeEventListAtTime - get time event list at given time
func (gm GlobalManager) GetTimeEventListAtTime(ctx sdk.Context, unixTime int64) linotypes.TimeEventList {
	eventList := gm.storage.GetTimeEventList(ctx, unixTime)
	return *eventList
}

// GetLastBlockTime - get last block time from KVStore
func (gm *GlobalManager) GetLastBlockTime(ctx sdk.Context) int64 {
	globalTime, err := gm.storage.GetGlobalTime(ctx)
	if err != nil {
		panic("Global Time Not Initialized at Genesis")
	}
	return globalTime.LastBlockTime
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

func (gm *GlobalManager) ExportToFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error {
	// state := &model.GlobalTablesIR{
	// 	Version: exportVersion,
	// }
	// storeMap := gm.storage.PartialStoreMap(ctx)

	// // export events
	// storeMap[string(model.TimeEventListSubStore)].Iterate(func(key []byte, val interface{}) bool {
	// 	ts, err := strconv.ParseInt(string(key), 10, 64)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	events := val.(*types.TimeEventList)
	// 	state.GlobalTimeEventLists = append(state.GlobalTimeEventLists, model.GlobalTimeEventsIR{
	// 		UnixTime:      ts,
	// 		TimeEventList: *events,
	// 	})
	// 	return false
	// })

	// // export stakes
	// storeMap[string(model.LinoStakeStatSubStore)].Iterate(func(key []byte, val interface{}) bool {
	// 	day, err := strconv.ParseInt(string(key), 10, 64)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	stakeStats := val.(*model.LinoStakeStat)
	// 	state.GlobalStakeStats = append(state.GlobalStakeStats, model.GlobalStakeStatDayIR{
	// 		Day:       day,
	// 		StakeStat: model.LinoStakeStatIR(*stakeStats),
	// 	})
	// 	return false
	// })

	// meta, err := gm.storage.GetGlobalMeta(ctx)
	// if err != nil {
	// 	return err
	// }
	// state.Meta = model.GlobalMetaIR(*meta)

	// pool, err := gm.storage.GetInflationPool(ctx)
	// if err != nil {
	// 	return err
	// }
	// state.InflationPool = model.InflationPoolIR{
	// 	DeveloperInflationPool: pool.DeveloperInflationPool,
	// 	ValidatorInflationPool: pool.ValidatorInflationPool,
	// }

	// consumption, err := gm.storage.GetConsumptionMeta(ctx)
	// if err != nil {
	// 	return err
	// }
	// state.ConsumptionMeta = model.ConsumptionMetaIR(*consumption)

	// tps, err := gm.storage.GetTPS(ctx)
	// if err != nil {
	// 	return err
	// }
	// state.TPS = model.TPSIR(*tps)

	// globalt, err := gm.storage.GetGlobalTime(ctx)
	// if err != nil {
	// 	return err
	// }
	// state.Time = model.GlobalTimeIR(*globalt)

	// return utils.Save(filepath, cdc, state)
}

func (gm *GlobalManager) ImportFromFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error {
	// rst, err := utils.Load(filepath, cdc, func() interface{} { return &model.GlobalTablesIR{} })
	// if err != nil {
	// 	return err
	// }
	// table := rst.(*model.GlobalTablesIR)

	// if table.Version != importVersion {
	// 	return fmt.Errorf("unsupported import version: %d", table.Version)
	// }

	// // import events
	// for _, v := range table.GlobalTimeEventLists {
	// 	err := gm.storage.SetTimeEventList(ctx, v.UnixTime, &v.TimeEventList)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// // import table.GlobalStakeStats
	// for _, v := range table.GlobalStakeStats {
	// 	stat := model.LinoStakeStat(v.StakeStat)
	// 	err := gm.storage.SetLinoStakeStat(ctx, v.Day, &stat)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// meta := model.GlobalMeta(table.Meta)
	// err = gm.storage.SetGlobalMeta(ctx, &meta)
	// if err != nil {
	// 	return err
	// }

	// pool := model.InflationPool{
	// 	DeveloperInflationPool: table.InflationPool.DeveloperInflationPool,
	// 	ValidatorInflationPool: table.InflationPool.ValidatorInflationPool,
	// }
	// err = gm.storage.SetInflationPool(ctx, &pool)
	// if err != nil {
	// 	return err
	// }

	// consumption := model.ConsumptionMeta(table.ConsumptionMeta)
	// err = gm.storage.SetConsumptionMeta(ctx, &consumption)
	// if err != nil {
	// 	return err
	// }

	// t := model.GlobalTime(table.Time)
	// err = gm.storage.SetGlobalTime(ctx, &t)
	// if err != nil {
	// 	return err
	// }

	// return nil
}

// XXX(yumin): if we want to add back the following codes, or handle the param change in the
// same way as before, add back:
// in global's storage cdc.
// or we can find a better way to do it.
// get and set params
// TODO add more change methods
// func (gm *GlobalManager) ChangeGlobalInflationParam(ctx sdk.Context, InfraAllocation sdk.Dec,
// 	ContentCreatorAllocation sdk.Dec, DeveloperAllocation sdk.Dec, ValidatorAllocation sdk.Dec) sdk.Error {
// 	allocation, err := gm.paramHolder.GetGlobalAllocationParam(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	allocation.ContentCreatorAllocation = ContentCreatorAllocation
// 	allocation.DeveloperAllocation = DeveloperAllocation
// 	allocation.InfraAllocation = InfraAllocation
// 	allocation.ValidatorAllocation = ValidatorAllocation
//
// 	if err := gm.paramHolder.SetGlobalAllocationParam(ctx, allocation); err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func (gm *GlobalManager) ChangeInfraInternalInflationParam(
// 	ctx sdk.Context, StorageAllocation sdk.Dec, CDNAllocation sdk.Dec) sdk.Error {
// 	allocation, err := gm.storage.GetInfraInternalAllocationParam(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	allocation.CDNAllocation = CDNAllocation
// 	allocation.StorageAllocation = StorageAllocation
// 	if err := gm.storage.SetInfraInternalAllocationParam(ctx, allocation); err != nil {
// 		return err
// 	}
// 	return nil
// }
