package postgres

import (
	"context"
	"testing"

	"github.com/dipress/crmifc/internal/category"
)

func TestCreateCategory(t *testing.T) {
	t.Parallel()

	r := NewRepository(db)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nc := category.NewCategory{
		Name: "Real IP",
	}

	var cat category.Category
	err := r.CreateCategory(ctx, &nc, &cat)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if cat.ID == 0 {
		t.Error("expected to parse returned id")
	}
}
