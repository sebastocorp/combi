package combi

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"combi/api/v1alpha5"
	"combi/internal/logger"
	"combi/internal/sets/actions"
	"combi/internal/sets/conditions"
	"combi/internal/sets/credentials"
	"combi/internal/sets/sources"
	"combi/internal/tmpl"
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

	c.target.encType = cfg.Conf.Target.Encoder
	c.target.build.file = cfg.Conf.Target.Build.File
	c.target.build.mode = fs.FileMode(cfg.Conf.Target.Build.Mode)
	c.target.build.typep = cfg.Conf.Target.Build.Type
	c.target.build.src = cfg.Conf.Target.Build.Source
	c.target.build.tmpl, err = tmpl.NewTemplate("result", cfg.Conf.Target.Build.Template)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(c.target.build.file), utils.DirModePerm)
	if err != nil {
		return err
	}

	for _, dirv := range []string{"tmpl"} {
		err = os.MkdirAll(filepath.Join(cfg.Conf.WorkingDir, dirv), utils.DirModePerm)
		if err != nil {
			return err
		}
	}

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
		var workDirName string
		workDirName, err = utils.GenHashString(strings.Join([]string{sv.Name, sv.Type, sv.Encoder}, "."))
		if err != nil {
			return err
		}

		workDirPath := filepath.Join(cfg.Conf.WorkingDir, workDirName)
		for _, dirv := range []string{"sync"} {
			err = os.MkdirAll(filepath.Join(workDirPath, dirv), utils.DirModePerm)
			if err != nil {
				return err
			}
		}
		err = c.srcs.Add(sources.OptionsT{
			Name:    sv.Name,
			SrcType: sv.Type,
			EncType: sv.Encoder,
			WorkDir: workDirPath,
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

	c.target.cons, err = conditions.NewSet()
	if err != nil {
		return err
	}

	c.target.acts, err = actions.NewActionSet()
	if err != nil {
		return err
	}

	return nil
}
