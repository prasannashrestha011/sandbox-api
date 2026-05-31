package sb_executil

import (
	"bytes"
	"context"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"
)

func ExecCreate(ctx context.Context, apiClient *client.Client, containerID string, cmd []string) (string, error) {
	execID, err := apiClient.ExecCreate(ctx, containerID, client.ExecCreateOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	})
	if err != nil {
		return "", err
	}

	resp, err := apiClient.ExecAttach(ctx, execID.ID, client.ExecAttachOptions{})
	if err != nil {
		return "", err
	}
	defer resp.Close()

	var buf bytes.Buffer

	stdcopy.StdCopy(&buf, &buf, resp.Reader)
	println(buf.String())
	return buf.String(), nil
}
