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

type UniqueIDTagGroup struct {
	tagging.TagGroup
}

func (d *UniqueIDTagGroup) CreateTagsForBlock(block structure.IBlock) error {
	log.Println("[PLUGIN] CreateTagsForBlock appelé pour:", block.GetResourceID())
	return d.UpdateBlockTags(block, block)
}

func (d *UniqueIDTagGroup) GetDefaultTags() []tags.ITag {
	log.Println("[PLUGIN] GetDefaultTags appelé")
	return []tags.ITag{
		&UniqueIDTag{},
	}
}

func (d *UniqueIDTagGroup) InitTagGroup(_ string, skippedTags []string, explicitlySpecifiedTags []string, options ...tagging.InitTagGroupOption) {
	log.Println("[PLUGIN] InitTagGroup appelé")
	opt := tagging.InitTagGroupOptions{
		TagPrefix: "",
	}
	for _, fn := range options {
		fn(&opt)
	}
	d.SkippedTags = skippedTags
	d.SpecifiedTags = explicitlySpecifiedTags
	d.SetTags(d.GetDefaultTags())
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

// IMPORTANT : Essayez cette export
var TagGroup tagging.ITagGroup = &UniqueIDTagGroup{}
