package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// adminNamePattern matches admin names of the form "First Last <Email>"
var adminNamePattern = regexp.MustCompile(`^\s*([^<]*)(\s*<\s*([^>]*)\s*>)?\s*$`)

// EnrichMap returns back the given map configured to use the first game config
// in the voting handler. If the map already has options then it is returned
// unchanged.
func (c Config) EnrichMap(m string) (string, error) {
	if strings.ContainsRune(m, '?') {
		return m, nil
	}

	configs, err := c.Voting.ExtendedGameConfigs()
	if err != nil {
		return "", err
	}

	if len(configs) == 0 {
		return m, nil
	}

	return configs[0].AppendParams(m), nil
}

func (c Config) Transform(w io.Writer, r io.Reader) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	out := bufio.NewWriter(w)
	defer out.Flush()

	var adminName, adminEmail string
	admin, err := Evaluate(c.Admin)
	if err != nil {
		return err
	}

	parsed := adminNamePattern.FindStringSubmatch(admin)
	if parsed != nil {
		adminName = parsed[1]
		adminEmail = parsed[3]
	}

	section := ""
	pendingAdds := make([]string, 0, 10)

	var votingConfigsReplaced bool

	for scanner.Scan() {
		line := scanner.Text()

		writeLine := func() {
			out.WriteString(line + "\n")
		}

		if s, isSection := sectionName(line); isSection {
			section = s
			pendingAdds = pendingAdds[:0]
			writeLine()
			continue
		}

		key, _, found := strings.Cut(line, "=")
		if !found {
			writeLine()
			continue
		}

		if strings.EqualFold(section, "URL") {
			switch strings.ToLower(key) {
			case "port":
				value, err := Evaluate(c.Port)
				if err != nil {
					return err
				}

				port, err := strconv.ParseUint(value, 10, 64)
				if err != nil {
					return err
				}

				fmt.Fprintf(out, "%s=%d\n", key, port)
			default:
				writeLine()
			}

			continue
		}

		if strings.EqualFold(section, "Engine.GameReplicationInfo") {
			switch strings.ToLower(key) {
			case "servername":
				value, err := Evaluate(c.Name)
				if err != nil {
					return err
				}

				fmt.Fprintf(out, "%s=%s\n", key, value)
			case "adminname":
				fmt.Fprintf(out, "%s=%s\n", key, adminName)
			case "adminemail":
				fmt.Fprintf(out, "%s=%s\n", key, adminEmail)
			case "messageoftheday":
				value, err := Evaluate(c.MOTD)
				if err != nil {
					return err
				}

				fmt.Fprintf(out, "%s=%s\n", key, value)
			default:
				writeLine()
			}

			continue
		}

		if strings.EqualFold(section, "xVoting.xVotingHandler") {
			switch strings.ToLower(key) {
			case "gameconfig":
				if !votingConfigsReplaced {
					configs, err := c.Voting.ExtendedGameConfigs()
					if err != nil {
						return err
					}

					for _, c := range configs {
						fmt.Fprintf(out, "%s=%s\n", key, c.GameConfigString())
					}

					// Ignore remaining configs from input
					votingConfigsReplaced = true
				}
			default:
				writeLine()
			}

			continue
		}

		writeLine()
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func sectionName(line string) (string, bool) {
	line = strings.TrimSpace(line)

	if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
		return line[1 : len(line)-1], true
	}

	return "", false
}
