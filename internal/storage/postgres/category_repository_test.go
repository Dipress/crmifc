package postgres

import (
	"context"
	"testing"

	"github.com/dipress/crmifc/internal/category"
)

func TestCreateCategory(t *testing.T) {
	t.Log("with initialized repository")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		r := NewCategoryRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		t.Log("\ttest:0\tshould create the category into the database")
		{
			nc := category.NewCategory{
				Name: "Real IP",
			}

			var cat category.Category
			err := r.Create(ctx, &nc, &cat)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if cat.ID == 0 {
				t.Error("expected to parse returned id")
			}

		}
	}

}

func TestFindCategory(t *testing.T) {
	t.Log("with initialized repository")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		r := NewCategoryRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		nc := category.NewCategory{
			Name: "IP Fake",
		}

		var cat category.Category
		err := r.Create(ctx, &nc, &cat)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		t.Log("\ttest:0\tshould find the category into the database")
		{
			_, err := r.Find(ctx, cat.ID)
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

		r := NewCategoryRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		nc := category.NewCategory{
			Name: "Channels",
		}

		var cat category.Category
		err := r.Create(ctx, &nc, &cat)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		t.Log("\ttest:0\tshould update the category into the database")
		{
			cat.Name = "Region Channels"
			err := r.Update(ctx, 1, &cat)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}

func TestCategoryDelete(t *testing.T) {
	t.Log("with initialized repository")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		r := NewCategoryRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		nc := category.NewCategory{
			Name: "Airmax",
		}

		var cat category.Category
		err := r.Create(ctx, &nc, &cat)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		t.Log("\ttest:0\tshould delete the category into the database")
		{
			err := r.Delete(ctx, 1)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}

func TestListCategories(t *testing.T) {
	t.Log("with initialized repository")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		r := NewCategoryRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		nc1 := category.NewCategory{
			Name: "Meets",
		}

		nc2 := category.NewCategory{
			Name: "Tasks",
		}

		var cat1 category.Category
		err := r.Create(ctx, &nc1, &cat1)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		var cat2 category.Category
		err = r.Create(ctx, &nc2, &cat2)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		t.Log("\ttest:0\tshould show list of categories")
		{
			var categories category.Categories
			err := r.List(ctx, &categories)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(categories.Categories) != 2 {
				t.Error("expected to slice of two categories")
			}
		}
	}
}
