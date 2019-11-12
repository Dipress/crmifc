package postgres

import (
	"context"
	"testing"

	"github.com/dipress/crmifc/internal/role"
)

func TestRoleCreate(t *testing.T) {
	t.Log("with initialized repository")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		r := NewRoleRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		t.Log("\ttest:0\tshould create the role into the database")
		{

			nr := role.NewRole{
				Name: "Admin",
			}

			var rol role.Role
			err := r.Create(ctx, &nr, &rol)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if rol.ID == 0 {
				t.Error("expected to parse returned id")
			}
		}
	}

}

func TestRoleFind(t *testing.T) {
	t.Log("with initialized repository")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		r := NewRoleRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		nr := role.NewRole{
			Name: "Manager",
		}

		var rol role.Role
		err := r.Create(ctx, &nr, &rol)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		t.Log("\ttest:0\tshould find the role into the database")
		{
			_, err := r.Find(ctx, rol.ID)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}

func TestRoleUpdate(t *testing.T) {
	t.Log("with initialized repository")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		r := NewRoleRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		nr := role.NewRole{
			Name: "Manager",
		}

		var rol role.Role
		err := r.Create(ctx, &nr, &rol)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		t.Log("\ttest:0\tshould update the role into the database")
		{
			rol.Name = "Member"
			err := r.Update(ctx, 1, &rol)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}

func TestRoleDelete(t *testing.T) {
	t.Log("with initialized repository")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		r := NewRoleRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		nr := role.NewRole{
			Name: "Admin",
		}

		var rol role.Role
		err := r.Create(ctx, &nr, &rol)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		t.Log("\ttest:0\tshould delete the role into the database")
		{
			err := r.Delete(ctx, 1)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}

func TestListRoles(t *testing.T) {
	t.Log("with initialized repository")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		r := NewRoleRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		nr1 := role.NewRole{
			Name: "Admin",
		}

		nr2 := role.NewRole{
			Name: "Manager",
		}

		var rol1 role.Role
		err := r.Create(ctx, &nr1, &rol1)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		var rol2 role.Role
		err = r.Create(ctx, &nr2, &rol2)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		t.Log("\ttest:0\tshould show list of roles")
		{
			var roles role.Roles
			err := r.List(ctx, &roles)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(roles.Roles) != 2 {
				t.Error("expected to slice of two roles")
			}
		}
	}
}
