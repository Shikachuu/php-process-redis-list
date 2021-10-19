package cmd

import (
	"github.com/Shikachuu/php-process-redis-list/pkg/queue"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

func RootCommand() *cobra.Command {
	var (
		version           bool
		redisHost         string
		redisPassword     string
		redisDb           int
		redisList         string
		numberOfWorkers   int
		commandForWorkers string
	)

	cmd := &cobra.Command{
		Use:   "qpp",
		Short: "Start the QPP to spawn processes based on queue messages",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if version {
				printVersion()
				return nil
			}
			q, err := queue.NewRedisQueue(redisHost, redisPassword, redisDb, redisList)
			if err != nil {
				return err
			}
			qmc := make(chan string)
			go func() {
				_ = q.Listen(qmc)
			}()
			select {
			case msg := <-qmc:
				cmd := exec.Command(commandForWorkers)
				cmd.Stdin = strings.NewReader(msg)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Env = os.Environ()
				go func() {
					_ = cmd.Start()
					_ = cmd.Wait()
				}()
			}
			return nil
		},
	}
	// Flags
	cmd.Flags().StringVar(&redisHost, "redis.host", "127.0.0.1:6379", "Set the address for the redis host")
	cmd.Flags().StringVar(&redisPassword, "redis.password", "", "Set the password for the redis host")
	cmd.Flags().IntVar(&redisDb, "redis.db", 0, "Set the DB for the redis host")
	cmd.Flags().StringVar(&redisList, "redis.list", "", "Set the list that the app listens on")
	cmd.Flags().IntVar(&numberOfWorkers, "workers.max", 1, "Set the max number of running processes")
	cmd.Flags().StringVar(&commandForWorkers, "workers.command", "", "Set the command that the workers will execute")

	cmd.Flags().BoolVarP(&version, "version", "v", false, "Show the version information")
	// Commands
	cmd.AddCommand(VersionCommand())

	return cmd
}
