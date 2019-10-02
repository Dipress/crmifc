package postgres

import (
	"context"
	"testing"

	"github.com/dipress/crmifc/internal/role"
)

func TestRoleCreate(t *testing.T) {
	t.Parallel()

	db, teardown := postgresDB(t)
	defer teardown()

	r := NewRepository(db)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nr := role.NewRole{
		Name: "Admin",
	}

	var rol role.Role
	err := r.CreateRole(ctx, &nr, &rol)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if rol.ID == 0 {
		t.Error("expected to parse returned id")
	}

}

func TestRoleFind(t *testing.T) {
	t.Log("with initialized repository")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		r := NewRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		nr := role.NewRole{
			Name: "Manager",
		}

		var rol role.Role
		err := r.CreateRole(ctx, &nr, &rol)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		t.Log("\ttest:0\tshould find the role into the database")
		{
			_, err := r.FindRole(ctx, rol.ID)
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

		r := NewRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		nr := role.NewRole{
			Name: "Manager",
		}

		var rol role.Role
		err := r.CreateRole(ctx, &nr, &rol)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		t.Log("\ttest:0\tshould update the role into the database")
		{
			rol.Name = "Member"
			err := r.UpdateRole(ctx, 1, &rol)
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

		r := NewRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		nr := role.NewRole{
			Name: "Admin",
		}

		var rol role.Role
		err := r.CreateRole(ctx, &nr, &rol)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		t.Log("\ttest:0\tshould delete the role into the database")
		{
			err := r.DeleteRole(ctx, 1)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}

func TestListRoles(t *testing.T) {
	t.Parallel()
	t.Log("with initialized repository")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		r := NewRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		nr1 := role.NewRole{
			Name: "Admin",
		}

		nr2 := role.NewRole{
			Name: "Manager",
		}

		var rol1 role.Role
		err := r.CreateRole(ctx, &nr1, &rol1)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		var rol2 role.Role
		err = r.CreateRole(ctx, &nr2, &rol2)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		t.Log("\ttest:0\tshould show list of roles")
		{
			var roles role.Roles
			err := r.ListRoles(ctx, &roles)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(roles.Roles) != 2 {
				t.Error("expected to slice of two roles")
			}
		}
	}
}
