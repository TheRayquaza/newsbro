package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Feed represents both meta and normal RSS feed entries
type Feed struct {
	Name        string   `yaml:"name"`
	DisplayName string   `yaml:"display_name,omitempty"`
	Description string   `yaml:"description,omitempty"`
	Link        string   `yaml:"link,omitempty"`
	Parents     []string `yaml:"parents,omitempty"`
}

// Root structure of your feeds.yaml
type FeedsFile struct {
	RSSMetaFeeds []Feed `yaml:"rss_meta_feeds"`
	RSSFeeds     []Feed `yaml:"rss_feeds"`
}

// slugify replicates your Python slugify logic
func slugify(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")

	re := regexp.MustCompile(`[^a-z0-9-]`)
	name = re.ReplaceAllString(name, "")
	return name
}

func refactorFeeds(feeds []Feed) []Feed {
	for i, f := range feeds {
		if f.DisplayName == "" {
			f.DisplayName = f.Name
		}
		origName := f.Name
		f.Name = slugify(f.Name)

		if len(f.Parents) > 0 {
			for j, p := range f.Parents {
				f.Parents[j] = slugify(p)
			}
		}

		fmt.Printf("Refactored: %-25s → %-25s\n", origName, f.Name)
		feeds[i] = f
	}
	return feeds
}

func main() {
	inputFile := "rss.yaml"
	outputFile := "feeds_refactored.yaml"

	// Read YAML file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		panic(fmt.Errorf("failed to read %s: %w", inputFile, err))
	}

	var feeds FeedsFile
	if err := yaml.Unmarshal(data, &feeds); err != nil {
		panic(fmt.Errorf("failed to parse YAML: %w", err))
	}

	// Refactor both meta and normal feeds
	feeds.RSSMetaFeeds = refactorFeeds(feeds.RSSMetaFeeds)
	feeds.RSSFeeds = refactorFeeds(feeds.RSSFeeds)

	// Write back
	out, err := yaml.Marshal(&feeds)
	if err != nil {
		panic(fmt.Errorf("failed to marshal YAML: %w", err))
	}

	if err := os.WriteFile(outputFile, out, 0644); err != nil {
		panic(fmt.Errorf("failed to write %s: %w", outputFile, err))
	}

	fmt.Printf("\n✅ Refactored feeds written to %s\n", outputFile)
}
