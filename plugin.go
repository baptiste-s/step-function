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

func (d *UniqueIDTagGroup) GetDefaultTags() []tags.ITag {
	log.Println("[PLUGIN] GetDefaultTags appelé")
	defaultTags := []tags.ITag{
		&UniqueIDTag{},
	}
	for _, tag := range defaultTags {
		tag.Init()
	}
	return defaultTags
}

func (d *UniqueIDTagGroup) InitTagGroup(_ string, skippedTags []string, explicitlySpecifiedTags []string, options ...tagging.InitTagGroupOption) {
	log.Printf("[PLUGIN] InitTagGroup appelé - skippedTags: %v, explicitlySpecifiedTags: %v\n", skippedTags, explicitlySpecifiedTags)
	
	opt := tagging.InitTagGroupOptions{
		TagPrefix: "",
	}
	for _, fn := range options {
		fn(&opt)
	}
	
	d.SkippedTags = skippedTags
	d.SpecifiedTags = explicitlySpecifiedTags
	d.SetTags(d.GetDefaultTags())
	
	log.Println("[PLUGIN] InitTagGroup terminé")
}

func (d *UniqueIDTagGroup) CreateTagsForBlock(block structure.IBlock) error {
	log.Printf("[PLUGIN] CreateTagsForBlock appelé pour: %s (type: %s)\n", 
		block.GetResourceID(), block.GetResourceType())
	
	err := d.UpdateBlockTags(block, block)
	if err != nil {
		log.Printf("[PLUGIN] ERREUR dans CreateTagsForBlock: %v\n", err)
	} else {
		log.Printf("[PLUGIN] CreateTagsForBlock OK pour %s\n", block.GetResourceID())
	}
	return err
}

type UniqueIDTag struct {
	tags.Tag
}

func (t *UniqueIDTag) Init() {
	log.Println("[PLUGIN] UniqueIDTag.Init appelé")
	t.Key = "carma-name"
	log.Printf("[PLUGIN] Tag key défini: %s\n", t.Key)
}

func (t *UniqueIDTag) CalculateValue(data interface{}) (tags.ITag, error) {
	log.Println("[PLUGIN] ===== CalculateValue appelé =====")
	
	block, ok := data.(structure.IBlock)
	if !ok {
		log.Println("[PLUGIN] ERREUR: Impossible de convertir data en IBlock")
		return nil, fmt.Errorf("failed to convert data to IBlock")
	}

	env := os.Getenv("YOR_ENV")
	if env == "" {
		env = "unknown"
		log.Println("[PLUGIN] WARNING: YOR_ENV non défini, utilisation de 'unknown'")
	}
	
	team := os.Getenv("YOR_TEAM")
	if team == "" {
		team = "unknown"
		log.Println("[PLUGIN] WARNING: YOR_TEAM non défini, utilisation de 'unknown'")
	}

	resourceID := block.GetResourceID()
	log.Printf("[PLUGIN] ResourceID: %s, Env: %s, Team: %s\n", resourceID, env, team)
	
	hash := md5.Sum([]byte(resourceID))
	hashString := hex.EncodeToString(hash[:])
	value := fmt.Sprintf("%s-%s-%s", env, team, hashString)
	
	log.Printf("[PLUGIN] Tag calculé: %s = %s\n", t.Key, value)
	
	// IMPORTANT: Créer un nouveau tag avec la valeur
	newTag := &UniqueIDTag{}
	newTag.Key = t.Key
	newTag.Value = value
	
	return newTag, nil
}

// Export du TagGroup
var ExternalTagGroups []tagging.ITagGroup

func init() {
	log.Println("[PLUGIN] Initialisation de ExternalTagGroups")
	ExternalTagGroups = []tagging.ITagGroup{
		&UniqueIDTagGroup{},
	}
	log.Printf("[PLUGIN] ExternalTagGroups configuré avec %d groupe(s)\n", len(ExternalTagGroups))
}
