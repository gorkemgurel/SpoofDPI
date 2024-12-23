package util

import (
	"fmt"
	"regexp"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

type Config struct {
	Addr            string
	Port            int
	DnsAddr         string
	DnsPort         int
	DnsIPv4Only     bool
	EnableDoh       bool
	Debug           bool
	Silent          bool
	SystemProxy     bool
	Timeout         int
	WindowSize      int
	AllowedPatterns []*regexp.Regexp
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		config = new(Config)
	}
	return config
}

func (c *Config) Load(args *Args) error {
	c.Addr = args.Addr
	c.Port = int(args.Port)
	c.DnsAddr = args.DnsAddr
	c.DnsPort = int(args.DnsPort)
	c.DnsIPv4Only = args.DnsIPv4Only
	c.Debug = args.Debug
	c.EnableDoh = args.EnableDoh
	c.Silent = args.Silent
	c.SystemProxy = args.SystemProxy
	c.Timeout = int(args.Timeout)
	c.WindowSize = int(args.WindowSize)
	patterns, err := parseAllowedPattern(args.AllowedPattern, args.PatternFile)
	if err != nil {
		return fmt.Errorf("parsing patterns: %w", err)
	}
	c.AllowedPatterns = patterns

	return nil
}

func parseAllowedPattern(patterns StringArray, filePath string) ([]*regexp.Regexp, error) {
	patternSet := make(map[*regexp.Regexp]struct{})
	
	filePatterns, err := loadPatternsFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("loading patterns from file: %w", err)
	}
	for _, pattern := range filePatterns {
		patternSet[pattern] = struct{}{}
	}
	
	for _, rawPattern := range patterns {
		pattern := regexp.MustCompile(rawPattern)
		patternSet[pattern] = struct{}{}
	}

	allowedPatterns := make([]*regexp.Regexp, len(patternSet))
	writeI := 0
	for pattern := range patternSet {
		allowedPatterns[writeI] = pattern
		writeI++
	}
	return allowedPatterns, nil
}

func PrintColoredBanner() {
	cyan := putils.LettersFromStringWithStyle("Spoof", pterm.NewStyle(pterm.FgCyan))
	purple := putils.LettersFromStringWithStyle("DPI", pterm.NewStyle(pterm.FgLightMagenta))
	pterm.DefaultBigText.WithLetters(cyan, purple).Render()

	pterm.DefaultBulletList.WithItems([]pterm.BulletListItem{
		{Level: 0, Text: "ADDR    : " + fmt.Sprint(config.Addr)},
		{Level: 0, Text: "PORT    : " + fmt.Sprint(config.Port)},
		{Level: 0, Text: "DNS     : " + fmt.Sprint(config.DnsAddr)},
		{Level: 0, Text: "DEBUG   : " + fmt.Sprint(config.Debug)},
	}).Render()

	pterm.DefaultBasicText.Println("Press 'CTRL + c' to quit")
}
