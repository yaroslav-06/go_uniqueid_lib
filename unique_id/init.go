package uniqueid

import (
	"fmt"
	"math/rand"
	"redirect_logger/internal/db"
)

const name = "generator"

type Generator struct {
	ctx           *db.DbCtx
	element_count uint64
}

/*
Returns new id generator,
needs a unique DbCtx (not used by any other Generator)
*/
func NewGenerator(ctx *db.DbCtx) *Generator {
	gr := &Generator{}
	gr.ctx = ctx
	gr.element_count = uint64(1 << 63)
	return gr
}

func (gr *Generator) DoesIdExists(id string) bool {
	vl, err := gr.ctx.MapFieldExists(name, id)
	return err == nil && vl
}

func (gr *Generator) GetNewId() (string, error) {
	vl, err := gr.ctx.GetMapSize(name)
	if err != nil {
		return "", fmt.Errorf("can't get map size for new id: %w", err)
	}
	if vl >= gr.element_count {
		return "", fmt.Errorf("can't generate any new ids limit of %d reached", gr.element_count)
	}
	for {
		new_id := fmt.Sprint(uint64(rand.Int63()) % gr.element_count)
		fmt.Printf("trying: %s\n", new_id)
		if !gr.DoesIdExists(new_id) {
			gr.ctx.SaveMapField(name, new_id, "t")
			return new_id, nil
		}
	}
}
