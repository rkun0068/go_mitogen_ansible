package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/spf13/cobra"

	"github.com/rkun0068/go_mitogen_ansible/output"
)

var (
	// Define command line parameters related to `ansible-playbook`
	hostsFile string
	// inventory string
	extraVars string
	tags      string
)

// Define counter for different results
var (
	successCount int32
	failureCount int32
)

// playbookCmd is a subcommand that simulates the ansible-playbook command
var playbookCmd = &cobra.Command{
	Use:   "playbook [yaml_file] -f [hosts_file]",
	Short: "Run Ansible playbooks",
	Long:  "Run Ansible playbooks with various options, including Mitogen and concurrency support.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get positional parameters playbook
		playbookArgs := args
		if len(playbookArgs) != 1 {
			fmt.Println("Error: Exactly one playbook must be specified")
			os.Exit(1)
		}
		// Read hosts file and split IP addresses
		hosts, err := readHostsFile(hostsFile)
		if err != nil {
			fmt.Println("func readHostsFile error: %v\n", err)
			os.Exit(1)
		}
		// Use WaitGroup to synchronize the goroutines
		var wg sync.WaitGroup
		// Execute playbook for each host in parallel using goroutines
		startTime := time.Now()
		for _, host := range hosts {
			wg.Add(1)
			go func(host string) {
				defer wg.Done()

				command := buildCommand(playbookArgs, host)
				err := executeCommand(command, host)
				if err != nil {
					fmt.Printf("Error executing command for host %s: %v\n", host, err)
					atomic.AddInt32(&failureCount, 1)
				} else {
					atomic.AddInt32(&successCount, 1)
				}

			}(host)
		}
		wg.Wait()
		duration := time.Since(startTime)
		total := successCount + failureCount
		fmt.Printf("\nExecution Summary:\n")
		fmt.Printf("Total: %d, Success: %d, Failure: %d, Duration: %s\n", total, successCount, failureCount, duration)
	},
}

func buildCommand(playbookArgs []string, host string) string {
	var command []string
	command = append(command, "ansible-playbook")

	if extraVars != "" {
		command = append(command, "-e", fmt.Sprintf("%q", extraVars))
	}

	if tags != "" {
		command = append(command, "-t", tags)
	}
	command = append(command, "-i", fmt.Sprintf("%q,", host))

	command = append(command, playbookArgs...)

	return strings.Join(command, " ")
}

func init() {
	// Add common command line flags to the playbook subcommand
	// playbookCmd.Flags().StringVarP(&inventory, "inventory", "i", "", "Specify inventory file path or comma-separated host list")
	playbookCmd.Flags().StringVarP(&hostsFile, "hosts-file", "f", "", "Specify hosts ip file path")
	playbookCmd.Flags().StringVarP(&extraVars, "extra-vars", "e", "", "Set additional variables as key=value or YAML/JSON")
	playbookCmd.Flags().StringVarP(&tags, "tags", "t", "", "Only run plays and tasks tagged with these values")
	rootCmd.AddCommand(playbookCmd)
}

// Read hosts file and return the list of hosts
func readHostsFile(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	var hosts []string
	for _, line := range lines {
		host := strings.TrimSpace(line)
		if host != "" {
			hosts = append(hosts, host)
		}
	}
	return hosts, nil
}

// Execute the command
func executeCommand(command string, host string) error {
	logFile, logger, err := output.CreateLogFile(host)
	if err != nil {
		return fmt.Errorf("func executeCommand error: %v", err)
	}

	defer output.CloseLogFile(logFile)

	//	execute command
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = logger.Writer()
	cmd.Stderr = logger.Writer()

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("func executeCommand error: %v", err)
	}

	return nil
}
