package main

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

var ErrNoSuchGameConfig = errors.New("no such game config")
var ErrLoopFound = errors.New("loop found")

type Config struct {
	Name  string `yaml:"name"`
	Admin string `yaml:"admin"`
	MOTD  string `yaml:"motd"`

	Voting *VotingConfig `yaml:"voting"`
}

type VotingConfig struct {
	Configs []VotingGameConfig `yaml:"configs"`
}

// GameConfigLines returns a slice of strings containing formatted xVoting
// handler `GameConfig=` values.
func (v *VotingConfig) GameConfigStrings() ([]string, error) {
	configs, err := v.ExtendedGameConfigs()
	if err != nil {
		return nil, err
	}

	strings := make([]string, 0, len(configs))
	for _, cfg := range configs {
		strings = append(strings, cfg.GameConfigString())
	}

	return strings, nil
}

func (v *VotingConfig) ExtendedGameConfigs() ([]VotingGameConfig, error) {
	extended := make([]VotingGameConfig, 0)
	for _, c := range v.Configs {
		cfg, err := v.ExtendedGameConfig(c.ID)
		if err != nil {
			return nil, err
		}

		extended = append(extended, cfg)
	}

	return extended, nil
}

// ExtendedGameConfig returns a VotingGameConfig value filled with any
// extensions and sets the `extends` field to an empty string.
func (v *VotingConfig) ExtendedGameConfig(name string) (VotingGameConfig, error) {
	configs := make([]VotingGameConfig, 0, 10)
	visited := make(map[string]struct{})

	curr := name
	for curr != "" {
		if _, seen := visited[curr]; seen {
			return VotingGameConfig{}, fmt.Errorf("error with %s: %w", name, ErrLoopFound)
		}

		var config VotingGameConfig
		var found bool

		for _, c := range v.Configs {
			if c.ID == curr {
				config = c
				found = true
				break
			}
		}

		if !found {
			return VotingGameConfig{}, fmt.Errorf("error with %s: %w", curr, ErrNoSuchGameConfig)
		}

		visited[curr] = struct{}{}
		configs = append(configs, config)
		curr = config.Extends
	}

	// Build config starting from the root
	result := VotingGameConfig{
		ID:       configs[0].ID,
		Mutators: make([]string, 0),
		Options:  make(map[string]string),
	}

	for i := len(configs) - 1; i >= 0; i-- {
		c := configs[i]
		if c.Name != "" {
			result.Name = c.Name
		}

		if c.Game != "" {
			result.Game = c.Game
		}

		if c.Prefix != "" {
			result.Prefix = c.Prefix
		}

		// Only append mutators
		for _, mut := range c.Mutators {
			if !slices.Contains(result.Mutators, mut) {
				result.Mutators = append(result.Mutators, mut)
			}
		}

		// Merge options
		for opt, value := range c.Options {
			result.Options[opt] = value
		}
	}

	return result, nil
}

type VotingGameConfig struct {
	Extends  string            `yaml:"extends"`
	ID       string            `yaml:"id"`
	Name     string            `yaml:"name"`
	Game     string            `yaml:"game"`
	Prefix   string            `yaml:"prefix"`
	Mutators []string          `yaml:"mutators"`
	Options  map[string]string `yaml:"options"`
}

// GameConfigString generates the approprate string to set as the value in
// UT2004.ini
func (c VotingGameConfig) GameConfigString() string {
	return fmt.Sprintf(
		`(GameClass="%s",Prefix="%s",Acronym="%s",GameName="%s",Mutators="%s",Options="%s")`,
		c.Game,
		c.Prefix,
		c.ID,
		c.Name,
		strings.Join(c.Mutators, ","),
		strings.Join(c.optionsAsParams(), "?"),
	)
}

func (c VotingGameConfig) optionsAsParams() []string {
	opts := make([]string, 0, len(c.Options))
	for k, v := range c.Options {
		opts = append(opts, fmt.Sprintf("%s=%s", k, v))
	}
	return opts
}

// AppendParams appends the options to the given string
func (c VotingGameConfig) AppendParams(s string) string {
	parts := []string{s}

	if c.Game != "" {
		parts = append(parts, "Game="+c.Game)
	}

	if len(c.Mutators) > 0 {
		parts = append(parts, "Mutator="+strings.Join(c.Mutators, ","))
	}

	opts := c.optionsAsParams()
	for _, opt := range opts {
		parts = append(parts, opt)
	}

	return strings.Join(parts, "?")
}
