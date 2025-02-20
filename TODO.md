
# TODO lists

## Features

- [ ] Add k8s command exec system
- [ ] Add patch system post merge
- [x] Consume combi config from k8s ConfigMap/Secrets

## Supports

- [ ] Add support to nginx conf files with custom parser
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
