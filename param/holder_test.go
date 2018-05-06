package param

// func TestChangeGlobalInflation(t *testing.T) {
// 	ctx, gm := setupTest(t)
//
// 	cases := []struct {
// 		contentCreatorAllocation sdk.Rat
// 		developerAllocation      sdk.Rat
// 		infraAllocation          sdk.Rat
// 		validatorAllocation      sdk.Rat
// 	}{
// 		{sdk.NewRat(1, 100), sdk.NewRat(50, 100), sdk.NewRat(20, 100), sdk.NewRat(29, 100)},
// 	}
//
// 	for _, cs := range cases {
// 		err := gm.ChangeGlobalInflationParam(
// 			ctx, cs.infraAllocation, cs.contentCreatorAllocation,
// 			cs.developerAllocation, cs.validatorAllocation)
// 		assert.Nil(t, err)
// 		allocation, err := gm.storage.GetGlobalAllocationParam(ctx)
// 		assert.Nil(t, err)
// 		assert.Equal(t, cs.contentCreatorAllocation, allocation.ContentCreatorAllocation)
// 		assert.Equal(t, cs.developerAllocation, allocation.DeveloperAllocation)
// 		assert.Equal(t, cs.validatorAllocation, allocation.ValidatorAllocation)
// 		assert.Equal(t, cs.infraAllocation, allocation.InfraAllocation)
// 	}
// }

// func TestChangeInfraInternalInflation(t *testing.T) {
// 	ctx, gm := setupTest(t)
//
// 	cases := []struct {
// 		storageAllocation sdk.Rat
// 		CDNAllocation     sdk.Rat
// 	}{
// 		{sdk.NewRat(1, 100), sdk.NewRat(99, 100)},
// 		{sdk.ZeroRat, sdk.OneRat},
// 	}
//
// 	for _, cs := range cases {
// 		err := gm.ChangeInfraInternalInflationParam(ctx, cs.storageAllocation, cs.CDNAllocation)
// 		assert.Nil(t, err)
// 		allocation, err := gm.paramHolder.GetInfraInternalAllocationParam(ctx)
// 		assert.Nil(t, err)
// 		assert.Equal(t, cs.storageAllocation, allocation.StorageAllocation)
// 		assert.Equal(t, cs.CDNAllocation, allocation.CDNAllocation)
// 	}
// }
