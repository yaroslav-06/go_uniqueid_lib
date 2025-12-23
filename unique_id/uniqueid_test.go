package uniqueid

import (
	"context"
	"fmt"
	"redirect_logger/internal/db"
	"testing"
)

func TestInit(t *testing.T) {
	ctx, err := db.GetParentContext("6379", "redirect_logger_test", context.Background())
	if err != nil {
		t.Errorf("can't get parent context")
	}
	NewGenerator(ctx.GetChild("test1"))
}

func TestNewId(t *testing.T) {
	ctx, err := db.GetParentContext("6379", "redirect_logger_test", context.Background())
	if err != nil {
		t.Errorf("can't get parent context")
	}
	gr := NewGenerator(ctx.GetChild("test1"))
	const tst_amount = 11
	gr.element_count = uint64(tst_amount)
	var vls [tst_amount]string
	for i := range tst_amount {
		vls[i], err = gr.GetNewId()
		sz, _ := gr.ctx.GetMapSize(name)
		fmt.Printf("vl: %s, sz: %d\n", vls[i], sz)
		if err != nil {
			t.Errorf("error generating new id: %s", err.Error())
		}
	}
	gsize, err := gr.ctx.GetMapSize(name)
	if err != nil {
		t.Errorf("couldn't get map size: %s", err.Error())
	}
	if gsize != uint64(tst_amount) {
		t.Errorf("the generator size doesn't match (%d instead of %d)", gsize, tst_amount)
	}

	_, err = gr.GetNewId()
	if err == nil {
		t.Errorf("somehow generated more id's than possible")
	}
	if err.Error() != fmt.Sprintf("can't generate any new ids limit of %d reached", tst_amount) {
		t.Errorf("the error message doesn't match: %s", err.Error())
	}
	gr.ctx.ClearField(name)
}
