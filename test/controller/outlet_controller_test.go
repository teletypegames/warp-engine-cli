package controller_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

func TestOutletController_Index_ListsTheOutlets(t *testing.T) {
	svc := &mocksvc.MockOutletService{AllFunc: func() ([]domain.OutletModel, error) {
		return []domain.OutletModel{
			{Key: "wecp", Type: "wecp", Served: true, Description: "wecp (wecp)"},
			{Key: "assembly64", Type: "assembly64", Platforms: []string{"c64"}, Served: true, Description: "assembly64 (assembly64: c64)"},
		}, nil
	}}

	stdout, _, err := execute(cli.Router{OutletController: cli.OutletController{OutletService: svc}}, "outlet", "list")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "assembly64")
	assert.Contains(t, stdout, "c64")
}

func TestOutletController_Check_PrintsTheOfferAndFailsStrictOnASkip(t *testing.T) {
	var got_key string
	svc := &mocksvc.MockOutletService{CheckFunc: func(key string) (*domain.OutletCheckModel, error) {
		got_key = key
		return &domain.OutletCheckModel{Key: key, Type: "assembly64", Offer: true,
			Items:   []domain.OutletItemModel{{Name: "ariadne", Version: "1.0", Kind: "cartridge"}},
			Skipped: []domain.OutletSkipModel{{Name: "empty", Reason: "no cartridge"}}}, nil
	}}
	router := cli.Router{OutletController: cli.OutletController{OutletService: svc}}

	stdout, _, err := execute(router, "outlets", "check", "assembly64")
	assert.NoError(t, err)
	assert.Equal(t, "assembly64", got_key)
	assert.Contains(t, stdout, "ariadne")
	assert.Contains(t, stdout, "no cartridge")

	_, _, err = execute(router, "outlet", "check", "assembly64", "--strict")
	assert.Error(t, err)
}
