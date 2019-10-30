package linux

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/pkg/errors"
)

func runCommand(client *Client, sudo bool, command string, stdinContent string) (string, string, error) {
	if sudo && client.useSudo {
		command = fmt.Sprintf("sudo %s", command)
	}
	session, err := client.connection.NewSession()
	if err != nil {
		return "", "", errors.Wrap(err, "Failed to create session")
	}
	stdin, err := session.StdinPipe()
	defer stdin.Close()
	if err != nil {
		return "", "", errors.Wrap(err, "Unable to setup stdin for session")
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		return "", "", errors.Wrap(err, "Unable to setup stderr for session")
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		return "", "", errors.Wrap(err, "Unable to setup stdout for session")
	}

	log.Printf("Running command %s", command)

	var stdoutOutput, stderrOutput []byte
	err = session.Start(command)
	if err != nil {
		return "", "", errors.Wrap(err, fmt.Sprintf("Unable to start command %s", command))
	}
	if stdinContent != "" {
		stdin.Write([]byte(stdinContent))
		stdin.Close()
	}
	err = session.Wait()

	if err != nil {
		stderrOutput, err2 := ioutil.ReadAll(stderr)
		if err2 != nil {
			log.Printf("Unable to read stderr for command: %v", err)
		}
		log.Printf("Stderr output: %s", strings.TrimSpace(string(stderrOutput)))

		return string(stdoutOutput), string(stderrOutput), errors.Wrap(err, fmt.Sprintf("Error running command %s", command))
	}
	stdoutOutput, err = ioutil.ReadAll(stdout)
	if err != nil {
		return string(stdoutOutput), string(stderrOutput), errors.Wrap(err, fmt.Sprintf("Unable to read stdout for command %s", command))
	}

	return string(stdoutOutput), string(stderrOutput), nil
}
