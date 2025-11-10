func loadConfig() error {
	// Chercher le fichier de config
	configPaths := []string{
		"yaml/custom-ca-name.yaml",
		"./yaml/custom-ca-name.yaml",
		"../yaml/custom-ca-name.yaml",
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
```

## Structure mise à jour
```
yor-custom-tagging/
├── yaml/
│   ├── tags.yaml             # Tags YAML standard
│   └── custom-ca-name.yaml   # Config du plugin custom
├── plugin/
│   └── yor_auto_unique_id.go
└── ...
