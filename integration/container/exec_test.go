package container // import "github.com/docker/docker/integration/container"

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/integration/internal/container"
	"github.com/docker/docker/integration/internal/request"
	"github.com/gotestyourself/gotestyourself/assert"
	is "github.com/gotestyourself/gotestyourself/assert/cmp"
	"github.com/gotestyourself/gotestyourself/skip"
)

func TestExec(t *testing.T) {
	skip.If(t, testEnv.OSType == "windows") // TODO Windows: FIXME. Probably needs to wait for container to be in running state.
	defer setupTest(t)()
	ctx := context.Background()
	client := request.NewAPIClient(t)

	cID := container.Run(t, ctx, client, container.WithTty(true), container.WithWorkingDir("/root"))

	id, err := client.ContainerExecCreate(ctx, cID,
		types.ExecConfig{
			WorkingDir:   "/tmp",
			Env:          strslice.StrSlice([]string{"FOO=BAR"}),
			AttachStdout: true,
			Cmd:          strslice.StrSlice([]string{"sh", "-c", "env"}),
		},
	)
	assert.NilError(t, err)

	resp, err := client.ContainerExecAttach(ctx, id.ID,
		types.ExecStartCheck{
			Detach: false,
			Tty:    false,
		},
	)
	assert.NilError(t, err)
	defer resp.Close()
	r, err := ioutil.ReadAll(resp.Reader)
	assert.NilError(t, err)
	out := string(r)
	assert.NilError(t, err)
	assert.Assert(t, is.Contains(out, "PWD=/tmp"), "exec command not running in expected /tmp working directory")
	assert.Assert(t, is.Contains(out, "FOO=BAR"), "exec command not running with expected environment variable FOO")
}
