
# TODO lists

## Features

- [x] Add k8s command exec system
- [x] Add env variable manage in sources
- [x] Add diferent flow with simple source
- [ ] Add patch system post merge
- [x] Consume combi config from k8s ConfigMap/Secrets

## Supported Formats

- [x] Add support to nginx conf files with custom parser
- [x] Add support to libconfig conf files with custom parser
- [ ] Add support to hcl files
- [x] Add support to json files with golang parser
- [x] Add support to yaml files with golang parser

## Code

- [x] Change merge order beetween global and specific configs
- [x] Add source and target structure to consume configs from different sources
- [x] Add type DaemonT with attached flow functions
- [ ] Refactor the code to clean it and add comments
- [x] Change to custom parser in libconfig kind instead use library
- [x] Change to custom parser code in nginx encoder
