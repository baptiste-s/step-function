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
	"gopkg.in/yaml.v2"
)

// Configuration du plugin
type PluginConfig struct {
	CacibName struct {
		Env  string `yaml:"env"`
		Team string `yaml:"team"`
	} `yaml:"cacib_name"`
}

var config PluginConfig

func init() {
	log.SetOutput(os.Stderr)
	log.Println("[PLUGIN] =========== Plugin UniqueID initialisé ===========")
	
	// Charger la config
	if err := loadConfig(); err != nil {
		log.Printf("[PLUGIN] WARNING: Erreur chargement config: %v, utilisation des variables d'env\n", err)
	}
}

func loadConfig() error {
	// Chercher le fichier de config
	configPaths := []string{
		"config/custom-ca-name.yaml",
		"./config/custom-ca-name.yaml",
		"../config/custom-ca-name.yaml",
		os.Getenv("YOR_PLUGIN_CONFIG"),
	}
	
	var configFile string
	for _, path := range configPaths {
		if path == "" {
			continue
		}
		if _, err := os.Stat(path); err == nil {
			configFile = path
			break
		}
	}
	
	if configFile == "" {
		return fmt.Errorf("fichier de config introuvable")
	}
	
	log.Printf("[PLUGIN] Chargement de la config depuis: %s\n", configFile)
	
	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}
	
	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}
	
	log.Printf("[PLUGIN] Config chargée: env=%s, team=%s\n", config.CacibName.Env, config.CacibName.Team)
	return nil
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
	
	// Priorité : config file > variables d'env > valeur par défaut
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
	
	// Chercher le yor_trace dans les tags existants et nouveaux
	var yorTrace string
	
	// Chercher dans les tags existants
	for _, tag := range block.GetExistingTags() {
		if tag.GetKey() == "yor_trace" {
			yorTrace = tag.GetValue()
			break
		}
	}
	
	// Si pas trouvé, chercher dans les nouveaux tags
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
		log.Printf("[PLUGIN] yor_trace trouvé: %s, utilisation des 16 premiers chars: %s\n", yorTrace, shortID)
	} else {
		// Fallback : MD5 du resourceID si yor_trace pas disponible
		log.Printf("[PLUGIN] WARNING: yor_trace non trouvé, fallback sur MD5 du resourceID\n")
		resourceID := block.GetResourceID()
		hash := md5.Sum([]byte(resourceID))
		hashString := hex.EncodeToString(hash[:])
		shortID = hashString[:16]
	}
	
	// Construire la valeur finale
	value := fmt.Sprintf("%s-%s-%s", env, team, shortID)
	
	log.Printf("[PLUGIN] Tag généré: %s = %s\n", t.Key, value)
	
	return &tags.Tag{Key: t.Key, Value: value}, nil
}

func (t *UniqueIDTag) GetPriority() int {
	return 1
}

func (t *UniqueIDTag) GetDescription() string {
	return "Custom Unique ID tag based on yor_trace"
}

var ExtraTagGroups = []interface{}{
	&UniqueIDTagGroup{},
}
