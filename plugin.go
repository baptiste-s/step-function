package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
)

func init() {
	log.SetOutput(os.Stderr)
	log.Println("[PLUGIN] =========== Plugin UniqueID initialisé ===========")
}

type UniqueIDTagGroup struct {
	tagging.TagGroup
}

// Implémentez toutes les méthodes requises explicitement
func (d *UniqueIDTagGroup) GetDefaultTags() []tags.ITag {
	log.Println("[PLUGIN] GetDefaultTags appelé")
	tag := &UniqueIDTag{}
	tag.Init()
	return []tags.ITag{tag}
}

func (d *UniqueIDTagGroup) InitTagGroup(dir string, skippedTags []string, explicitlySpecifiedTags []string, options ...tagging.InitTagGroupOption) {
	log.Printf("[PLUGIN] InitTagGroup appelé - skippedTags: %v\n", skippedTags)
	
	d.SkippedTags = skippedTags
	d.SpecifiedTags = explicitlySpecifiedTags
	d.Dir = dir
	d.SetTags(d.GetDefaultTags())
	
	log.Println("[PLUGIN] InitTagGroup terminé")
}

func (d *UniqueIDTagGroup) CreateTagsForBlock(block structure.IBlock) error {
	log.Printf("[PLUGIN] CreateTagsForBlock appelé pour: %s\n", block.GetResourceID())
	return d.UpdateBlockTags(block, block)
}

type UniqueIDTag struct {
	tags.Tag
}

func (t *UniqueIDTag) Init() {
	log.Println("[PLUGIN] UniqueIDTag.Init appelé")
	t.Key = "carma-name"
}

func (t *UniqueIDTag) CalculateValue(data interface{}) (tags.ITag, error) {
	log.Println("[PLUGIN] CalculateValue appelé")
	
	block, ok := data.(structure.IBlock)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to IBlock")
	}

	env := os.Getenv("YOR_ENV")
	if env == "" {
		env = "unknown"
	}
	
	team := os.Getenv("YOR_TEAM")
	if team == "" {
		team = "unknown"
	}

	resourceID := block.GetResourceID()
	hash := md5.Sum([]byte(resourceID))
	hashString := hex.EncodeToString(hash[:])
	value := fmt.Sprintf("%s-%s-%s", env, team, hashString)
	
	log.Printf("[PLUGIN] Tag généré: %s = %s\n", t.Key, value)
	
	return &tags.Tag{Key: t.Key, Value: value}, nil
}

func (t *UniqueIDTag) GetPriority() int {
	return 1
}

func (t *UniqueIDTag) GetDescription() string {
	return "Custom Unique ID tag"
}

// Export - essayez avec le pointeur concret
var ExtraTagGroups = []interface{}{
	&UniqueIDTagGroup{},
}
