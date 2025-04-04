package combi

import (
	"fmt"
	"path/filepath"

	"combi/api/v1alpha5"
	"combi/internal/logger"
	"combi/internal/sets/credentials"
	"combi/internal/sets/sources"
	"combi/internal/utils"
)

// setup TODO
func (c *CombiT) setup(cfg any) (err error) {
	switch cfg := cfg.(type) {
	case v1alpha5.CombiT:
		{
			err = c.v1alpha5Setup(cfg)
		}
	default:
		{
			err = fmt.Errorf("unsupported apiVersion")
		}
	}
	return err
}

// v1alpha4Setup TODO
func (c *CombiT) v1alpha5Setup(cfg v1alpha5.CombiT) (err error) {
	c.log = logger.NewLogger(logger.GetLevel(cfg.Conf.Logger.Level))
	c.syncTime = cfg.Conf.SyncTime

	c.creds, err = credentials.NewSet()
	if err != nil {
		return err
	}
	for _, cv := range cfg.Conf.Credentials {
		err = c.creds.Add(credentials.OptionsT{
			Name: cv.Name,
			Type: cv.Type,
			SshKey: credentials.OptionsSshKeyT{
				User:     cv.SshKey.User,
				SshKey:   cv.SshKey.SshKeyFile,
				Password: cv.SshKey.Password,
			},
			Kube: credentials.OptionsKubeT{
				InCluster:      cv.K8s.InCluster,
				KubeconfigPath: cv.K8s.KubeconfigPath,
				MasterUrl:      cv.K8s.MasterUrl,
			},
		})
		if err != nil {
			return err
		}
	}

	c.srcs, err = sources.NewSet()
	if err != nil {
		return err
	}
	for _, sv := range cfg.Conf.Sources {
		workDir := ""
		workDir, err = utils.GenHashString(sv.Name, sv.Type, sv.Encoder)
		if err != nil {
			return err
		}
		err = c.srcs.Add(sources.OptionsT{
			Name:    sv.Name,
			SrcType: sv.Type,
			EncType: sv.Encoder,
			WorkDir: filepath.Join(cfg.Conf.WorkingDir, workDir),
			CredRef: c.creds.Get(sv.Credential),

			File: sv.File,
			Git: sources.OptionsGitT{
				Url:      sv.Git.SshUrl,
				Branch:   sv.Git.Branch,
				Filepath: sv.Git.File,
			},
			K8s: sources.OptionsK8sT{
				Kind:      sv.K8s.Kind,
				Namespace: sv.K8s.Namespace,
				Name:      sv.K8s.Name,
				Key:       sv.K8s.Key,
			},
		})
		if err != nil {
			return err
		}
	}

	// cfg.Conf.TmpFiles
	return nil
}
