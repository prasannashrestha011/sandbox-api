package sb_executil

import (
	"bytes"
	"context"
	"main/internal/dto"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"
)

func ExecCreate(ctx context.Context, apiClient *client.Client, containerID string, cmd []string) (*dto.SandboxExecResponse, error) {
	execID, err := apiClient.ExecCreate(ctx, containerID, client.ExecCreateOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	})
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.ExecAttach(ctx, execID.ID, client.ExecAttachOptions{})
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	_, err = stdcopy.StdCopy(&stdoutBuf, &stderrBuf, resp.Reader)
	if err != nil {
		return nil, err
	}
	inspect, err := apiClient.ExecInspect(ctx, execID.ID, client.ExecInspectOptions{})
	if err != nil {
		return nil, err
	}

	return &dto.SandboxExecResponse{
		Stdout:   stdoutBuf.String(),
		Stderr:   stderrBuf.String(),
		ExitCode: inspect.ExitCode,
	}, nil
}
