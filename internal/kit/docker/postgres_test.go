package docker

import (
	"testing"

	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
)

const exec = `
CREATE TABLE example (
	id	CHAR (36) PRIMARY KEY
);
`

func Test_NewPostgres(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	t.Log("with dockertest pool")
	{
		pool, err := dockertest.NewPool("")
		assert.Nil(t, err)

		t.Log("\ttest: 0\t should start postgres image and return valid connection and resource to purge")
		{
			pd, err := NewPostgres(pool)
			assert.Nil(t, err)

			_, err = pd.DB.Exec(exec)
			assert.Nil(t, err)
			err = pool.Purge(pd.Resource)
			assert.Nil(t, err)
		}
	}
}
