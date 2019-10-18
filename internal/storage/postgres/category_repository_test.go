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

func TestFindCategory(t *testing.T) {
	t.Log("with initialized repository")
	{
		t.Parallel()

		r := NewRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		nc := category.NewCategory{
			Name: "IP Fake",
		}

		var cat category.Category
		err := r.CreateCategory(ctx, &nc, &cat)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		t.Log("\ttest:0\tshould find the category into the database")
		{
			_, err := r.FindCategory(ctx, cat.ID)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

		}
	}
}

func TestCategoryUpdate(t *testing.T) {
	t.Log("with initialized repository")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		r := NewRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		nc := category.NewCategory{
			Name: "Channels",
		}

		var cat category.Category
		err := r.CreateCategory(ctx, &nc, &cat)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		t.Log("\ttest:0\tshould update the category into the database")
		{
			cat.Name = "Region Channels"
			err := r.UpdateCategory(ctx, 1, &cat)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}
