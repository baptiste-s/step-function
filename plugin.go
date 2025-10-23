package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"github.com/bridgecrewio/yor/src/common/structure"
	"github.com/bridgecrewio/yor/src/common/tagging"
	"github.com/bridgecrewio/yor/src/common/tagging/tags"
)

type UniqueIDTagGroup struct {
	tagging.TagGroup
}

func (d *UniqueIDTagGroup) CreateTagsForBlock(block structure.IBlock) error {
	return d.UpdateBlockTags(block, block)
}

func (d *UniqueIDTagGroup) GetDefaultTags() []tags.ITag {
	return []tags.ITag{
		&UniqueIDTag{},
	}
}

func (d *UniqueIDTagGroup) InitTagGroup(_ string, skippedTags []string, explicitlySpecifiedTags []string, options ...tagging.InitTagGroupOption) {
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
	t.Key = "carma-name"
}

func (t *UniqueIDTag) CalculateValue(data interface{}) (tags.ITag, error) {
	block, ok := data.(structure.IBlock)
	if !ok {
		return nil, fmt.Errorf("failed to convert data to IBlock")
	}

	// Lire env et team depuis les variables d'environnement
	env := os.Getenv("YOR_ENV")
	if env == "" {
		env = "unknown"
	}

	team := os.Getenv("YOR_TEAM")
	if team == "" {
		team = "unknown"
	}

	resourceID := block.GetResourceID()

	// Calculer le MD5
	hash := md5.Sum([]byte(resourceID))
	hashString := hex.EncodeToString(hash[:])

	// Construire la valeur finale (TOUT est calculé, pas d'expression Terraform)
	value := fmt.Sprintf("%s-%s-%s", env, team, hashString)

	return &tags.Tag{Key: t.Key, Value: value}, nil
}

var ExtraTagGroups = []interface{}{&UniqueIDTagGroup{}}
