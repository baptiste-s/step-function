package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
	"gopkg.in/yaml.v2"
)

// Configuration (optionnelle)
type PluginConfig struct {
	CacibName struct {
		Env  string `yaml:"env"`
		Team string `yaml:"team"`
	} `yaml:"cacib_name"`
}

var config PluginConfig

func init() {
	// Essayer de charger le fichier de config, sinon on utilise les variables d'env
	data, err := os.ReadFile("yaml/custom-ca-name.yaml")
	if err == nil {
		yaml.Unmarshal(data, &config)
	}
}

type UniqueIDTagGroup struct {
	tagging.TagGroup
}

func (d *UniqueIDTagGroup) CreateTagsForBlock(block structure.IBlock) error {
	return d.UpdateBlockTags(block, block)
}

func (d *UniqueIDTagGroup) GetDefaultTags() []tags.ITag {
	tag := &UniqueIDTag{}
	tag.Init()
	return []tags.ITag{tag}
}

func (d *UniqueIDTagGroup) InitTagGroup(dir string, skippedTags []string, explicitlySpecifiedTags []string, options ...tagging.InitTagGroupOption) {
	d.SkippedTags = skippedTags
	d.Dir = dir
	d.SetTags(d.GetDefaultTags())
}

type UniqueIDTag struct {
	tags.Tag
}

func (t *UniqueIDTag) Init() {
	t.Key = "cacib-name"
}

func (t *UniqueIDTag) CalculateValue(data interface{}) (tags.ITag, error) {
	block, ok := data.(structure.IBlock)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to IBlock")
	}
	
	// Lire env et team (config > env vars > défaut)
	env := config.CacibName.Env
	if env == "" {
		env = os.Getenv("YOR_ENV")
	}
	if env == "" {
		env = "unknown"
	}
	
	team := config.CacibName.Team
	if team == "" {
		team = os.Getenv("YOR_TEAM")
	}
	if team == "" {
		team = "unknown"
	}
	
	// Chercher yor_trace
	var yorTrace string
	for _, tag := range block.GetExistingTags() {
		if tag.GetKey() == "yor_trace" {
			yorTrace = tag.GetValue()
			break
		}
	}
	if yorTrace == "" {
		for _, tag := range block.GetNewTags() {
			if tag.GetKey() == "yor_trace" {
				yorTrace = tag.GetValue()
				break
			}
		}
	}
	
	// Prendre les 16 premiers caractères
	var shortID string
	if yorTrace != "" && len(yorTrace) >= 16 {
		shortID = yorTrace[:16]
	} else {
		// Fallback sur MD5 si pas de yor_trace
		resourceID := block.GetResourceID()
		hash := md5.Sum([]byte(resourceID))
		hashString := hex.EncodeToString(hash[:])
		shortID = hashString[:16]
	}
	
	value := fmt.Sprintf("%s-%s-%s", env, team, shortID)
	return &tags.Tag{Key: t.Key, Value: value}, nil
}

func (t *UniqueIDTag) GetPriority() int {
	return 1
}

func (t *UniqueIDTag) GetDescription() string {
	return "Custom Unique ID tag"
}

var ExtraTagGroups = []interface{}{&UniqueIDTagGroup{}}
