package sb_executil

import (
	"bytes"
	"context"
	"fmt"
	"main/internal/domain"
	"main/internal/dto"
	"time"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"
)

func ExecCreate(ctx context.Context, apiClient *client.Client, containerID string, cmd []string) (*dto.SandboxExecResponse, error) {
	exeCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	execID, err := apiClient.ExecCreate(exeCtx, containerID, client.ExecCreateOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	})
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.ExecAttach(exeCtx, execID.ID, client.ExecAttachOptions{})
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	copyDone := make(chan error, 1)

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	go func() {
		_, copyErr := stdcopy.StdCopy(&stdoutBuf, &stderrBuf, resp.Reader)
		copyDone <- copyErr
	}()

	select {
	case err := <-copyDone:
		if err != nil {
			return nil, err
		}
	case <-exeCtx.Done():
		resp.Close()
		inspect, err := apiClient.ExecInspect(context.Background(), execID.ID, client.ExecInspectOptions{})
		if err == nil && inspect.PID != 0 {
			killCmd := []string{"kill", "-9", fmt.Sprintf("%d", inspect.PID)}
			killExec, _ := apiClient.ExecCreate(context.Background(), containerID, client.ExecCreateOptions{Cmd: killCmd})
			apiClient.ExecStart(ctx, killExec.ID, client.ExecStartOptions{})
		}
		return nil, domain.InvalidRequestError("execution time out", nil)
	}
	inspect, err := apiClient.ExecInspect(context.Background(), execID.ID, client.ExecInspectOptions{})
	if err != nil {
		return nil, err
	}

	return &dto.SandboxExecResponse{
		Stdout:   stdoutBuf.String(),
		Stderr:   stderrBuf.String(),
		ExitCode: inspect.ExitCode,
	}, nil
}
