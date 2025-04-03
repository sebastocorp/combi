package combi

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"combi/api/v1alpha4"
	"combi/api/v1alpha5"
	"combi/internal/encoding"
	"combi/internal/logger"
	"combi/internal/sets/actions"
	"combi/internal/sets/conditions"
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
	case v1alpha4.CombiConfigT:
		{
			err = c.v1alpha4Setup(cfg)
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
		err = c.srcs.Add(sources.OptionsT{
			Name: sv.Name,
			Type: sv.Type,
			Cred: c.creds.Get(sv.Credential),

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

// v1alpha4Setup TODO
func (c *CombiT) v1alpha4Setup(cfg v1alpha4.CombiConfigT) error {
	c.log = logger.NewLogger(logger.GetLevel(cfg.Settings.Logger.Level))
	c.syncTime = cfg.Settings.SyncTime
	c.encoder = encoding.GetEncoder(cfg.Kind)

	// Target setup

	err := os.MkdirAll(cfg.Settings.Target.Path, fs.FileMode(cfg.Settings.Target.Mode))
	if err != nil {
		return err
	}

	c.target.filepath = filepath.Join(cfg.Settings.Target.Path, cfg.Settings.Target.File)
	c.target.mode = fs.FileMode(cfg.Settings.Target.Mode)

	// Sources setup

	for si := range cfg.Sources {
		var hashKey string
		hashKey, err = utils.GenHashString(cfg.Sources[si].Type, cfg.Sources[si].Name)
		if err != nil {
			return err
		}

		srcpath := filepath.Join(cfg.Settings.TmpFiles.Path, hashKey)
		err = os.MkdirAll(srcpath, fs.FileMode(cfg.Settings.TmpFiles.Mode))
		if err != nil {
			return err
		}

		err = c.srcs.Add(sources.OptionsT{
			Name: cfg.Sources[si].Name,
			Type: cfg.Sources[si].Type,
			Path: srcpath,

			Raw:  cfg.Sources[si].Raw,
			File: cfg.Sources[si].File,
			K8s: sources.OptionsK8sT{
				InCluster:      cfg.Sources[si].K8s.Context.InCluster,
				ConfigFilepath: cfg.Sources[si].K8s.Context.ConfigFilepath,
				MasterUrl:      cfg.Sources[si].K8s.Context.MasterUrl,
				Kind:           cfg.Sources[si].K8s.Kind,
				Namespace:      cfg.Sources[si].K8s.Namespace,
				Name:           cfg.Sources[si].K8s.Name,
				Key:            cfg.Sources[si].K8s.Key,
			},
			Git: sources.OptionsGitT{
				SshKeyFilepath: cfg.Sources[si].Git.SshKeyFilepath,
				Url:            cfg.Sources[si].Git.SshUrl,
				Branch:         cfg.Sources[si].Git.Branch,
				Filepath:       cfg.Sources[si].Git.Filepath,
			},
		})
		if err != nil {
			return err
		}
	}

	c.cons, err = conditions.NewSet()
	if err != nil {
		return err
	}
	for ci := range cfg.Behavior.Conditions {
		err = c.cons.Add(conditions.OptionsT{
			Name:      cfg.Behavior.Conditions[ci].Name,
			Mandatory: cfg.Behavior.Conditions[ci].Mandatory,
			Tmpl:      cfg.Behavior.Conditions[ci].Template,
			Expect:    cfg.Behavior.Conditions[ci].Expect,
		})
		if err != nil {
			return err
		}
	}

	c.acts, err = actions.NewActionSet()
	if err != nil {
		return err
	}
	for ai := range cfg.Behavior.Actions {
		err = c.acts.CreateAdd(actions.OptionsT{
			Name: cfg.Behavior.Actions[ai].Name,
			On:   cfg.Behavior.Actions[ai].On,
			Cmd:  cfg.Behavior.Actions[ai].Cmd,
			K8s: actions.OptionsK8sT{
				Namespace:      cfg.Behavior.Actions[ai].K8s.Namespace,
				Pod:            cfg.Behavior.Actions[ai].K8s.Pod,
				Container:      cfg.Behavior.Actions[ai].K8s.Container,
				InCluster:      cfg.Behavior.Actions[ai].K8s.Context.InCluster,
				ConfigFilepath: cfg.Behavior.Actions[ai].K8s.Context.ConfigFilepath,
				MasterUrl:      cfg.Behavior.Actions[ai].K8s.Context.MasterUrl,
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}
