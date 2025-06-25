package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/abhinavchadaga/claude-code-background-agent/internal/docker"
	"github.com/spf13/cobra"
)

var containerListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List running agent containers",
	RunE:    listContainers,
}

var containerStopCmd = &cobra.Command{
	Use:   "stop [CONTAINER_ID]",
	Short: "Stop & remove a container",
	Args:  cobra.ExactArgs(1),
	RunE:  stopContainer,
}

func init() {
	containerCmd.AddCommand(containerListCmd)
	containerCmd.AddCommand(containerStopCmd)
}

func listContainers(cmd *cobra.Command, args []string) error {
	dockerService, err := docker.NewService()
	if err != nil {
		return fmt.Errorf("failed to connect to Docker: %w", err)
	}
	defer dockerService.Close()

	ctx := context.Background()
	containers, err := dockerService.ListContainers(ctx)
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	fmt.Printf("%-12s %-30s %-15s %-10s\n", "CONTAINER ID", "IMAGE", "STATUS", "NAMES")
	for _, container := range containers {
		shortID := container.ID[:12]
		image := container.Image
		status := container.Status
		names := strings.Join(container.Names, ", ")

		if len(names) > 0 && names[0] == '/' {
			names = names[1:]
		}

		fmt.Printf("%-12s %-30s %-15s %-10s\n", shortID, image, status, names)
	}

	return nil
}

func stopContainer(cmd *cobra.Command, args []string) error {
	containerID := args[0]

	dockerService, err := docker.NewService()
	if err != nil {
		return fmt.Errorf("failed to connect to Docker: %w", err)
	}
	defer dockerService.Close()

	ctx := context.Background()

	fmt.Printf("Stopping container: %s\n", containerID)
	if err := dockerService.StopContainer(ctx, containerID); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	fmt.Printf("Removing container: %s\n", containerID)
	if err := dockerService.RemoveContainer(ctx, containerID); err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}

	fmt.Printf("Container %s stopped and removed successfully\n", containerID)
	return nil
}
