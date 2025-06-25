package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type Service struct {
	client *client.Client
}

type ContainerConfig struct {
	Image       string
	RepoPath    string
	WorkDir     string
	Command     []string
	Env         []string
	SetupScript string
}

func NewService() (*Service, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}
	return &Service{client: cli}, nil
}

func (s *Service) StartContainer(ctx context.Context, config ContainerConfig) (string, error) {
	err := s.pullImage(ctx, config.Image)
	if err != nil {
		return "", fmt.Errorf("failed to pull image: %w", err)
	}

	workDir := "/workspace"
	if config.WorkDir != "" {
		workDir = config.WorkDir
	}

	mounts := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: config.RepoPath,
			Target: workDir,
		},
	}

	containerCommand := config.Command

	if config.SetupScript != "" {
		if _, err := os.Stat(config.SetupScript); err == nil {
			mounts = append(mounts, mount.Mount{
				Type:     mount.TypeBind,
				Source:   config.SetupScript,
				Target:   "/setup.sh",
				ReadOnly: true,
			})

			setupCommand := []string{
				"/bin/bash", "-c",
				"echo 'Running setup script...' && chmod +x /setup.sh && /setup.sh && echo 'Setup complete. Running main command...' && " +
					fmt.Sprintf("exec %s", containerCommand[0]),
			}
			if len(containerCommand) > 1 {
				for i, arg := range containerCommand[1:] {
					if i == 0 {
						setupCommand[2] += " " + arg
					} else {
						setupCommand[2] += " '" + arg + "'"
					}
				}
			}
			containerCommand = setupCommand
		}
	}

	containerConfig := &container.Config{
		Image:      config.Image,
		WorkingDir: workDir,
		Env:        config.Env,
		Cmd:        containerCommand,
		Tty:        true,
		OpenStdin:  true,
	}

	hostConfig := &container.HostConfig{
		Mounts:     mounts,
		AutoRemove: true,
	}

	networkConfig := &network.NetworkingConfig{}

	resp, err := s.client.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, nil, "")
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	if err := s.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	return resp.ID, nil
}

func (s *Service) pullImage(ctx context.Context, imageName string) error {
	reader, err := s.client.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()

	_, err = io.Copy(os.Stdout, reader)
	return err
}

func (s *Service) WaitForContainer(ctx context.Context, containerID string) error {
	statusCh, errCh := s.client.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("error waiting for container: %w", err)
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return fmt.Errorf("container exited with status code: %d", status.StatusCode)
		}
	}
	return nil
}

func (s *Service) GetContainerLogs(ctx context.Context, containerID string) (io.ReadCloser, error) {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
	}
	return s.client.ContainerLogs(ctx, containerID, options)
}

func (s *Service) StopContainer(ctx context.Context, containerID string) error {
	return s.client.ContainerStop(ctx, containerID, container.StopOptions{})
}

func (s *Service) RemoveContainer(ctx context.Context, containerID string) error {
	return s.client.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
}

func (s *Service) ListContainers(ctx context.Context) ([]types.Container, error) {
	return s.client.ContainerList(ctx, container.ListOptions{All: true})
}

func (s *Service) Close() error {
	return s.client.Close()
}
