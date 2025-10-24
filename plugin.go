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
	log.Printf("[PLUGIN] CreateTagsForBlock appelé pour: %s\n", block.GetResourceID())
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

// IMPORTANT : Cette méthode est nécessaire pour que Yor reconnaisse le TagGroup
func (d *UniqueIDTagGroup) GetTagGroupByTags(tagsToApply []string) []tags.ITag {
	log.Println("[PLUGIN] GetTagGroupByTags appelé avec:", tagsToApply)
	return d.GetDefaultTags()
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
	
	t.Value = value
	return t, nil
}

func init() {
	log.Println("[PLUGIN] =========== Plugin UniqueID chargé ===========")
}

// CRITICAL: Le nom doit être exactement "ExternalTagGroups" (avec un 's')
var ExternalTagGroups []tagging.ITagGroup

func init() {
	ExternalTagGroups = []tagging.ITagGroup{
		&UniqueIDTagGroup{},
	}
	log.Println("[PLUGIN] ExternalTagGroups initialisé avec", len(ExternalTagGroups), "tag group(s)")
}
