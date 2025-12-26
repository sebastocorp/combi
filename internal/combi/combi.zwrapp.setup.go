package combi

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"combi/api/v1alpha5"
	"combi/internal/encoders"
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

	if !slices.Contains(
		[]string{encoders.TypeJSON, encoders.TypeYAML,
			encoders.TypeLIBCONFIG, encoders.TypeNGINX},
		cfg.Conf.Target.Encoder) {
		return fmt.Errorf("encoder type '%s' in target not suported", cfg.Conf.Target.Encoder)
	}
	c.target.encType = cfg.Conf.Target.Encoder

	c.target.build.bType = cfg.Conf.Target.Build.Type
	switch c.target.build.bType {
	case TargetBuildTypeSOURCE:
		{
			c.target.build.src = cfg.Conf.Target.Build.Source
		}
	case TargetBuildTypeTEMPLATE:
		{
			var tmplBytes []byte
			tmplBytes, err = os.ReadFile(cfg.Conf.Target.Build.Template)
			if err != nil {
				err = fmt.Errorf("error reading target template file: %w", err)
				return err
			}

			c.target.build.tmpl, err = tmpl.NewTemplate("result", string(tmplBytes))
			if err != nil {
				err = fmt.Errorf("error creating target template: %w", err)
				return err
			}
		}
	default:
		{
			return fmt.Errorf("build type '%s' in target not supported", cfg.Conf.Target.Build.Type)
		}
	}
	c.target.build.file = cfg.Conf.Target.Build.File
	c.target.build.mode = fs.FileMode(cfg.Conf.Target.Build.Mode)

	err = os.MkdirAll(filepath.Dir(c.target.build.file), utils.DirModePerm)
	if err != nil {
		err = fmt.Errorf("error creating target file directory: %w", err)
		return err
	}

	err = os.MkdirAll(cfg.Conf.WorkingDir, utils.DirModePerm)
	if err != nil {
		err = fmt.Errorf("error creating working directory: %w", err)
		return err
	}

	// Sets definitions

	c.creds, err = credentials.NewSet()
	if err != nil {
		err = fmt.Errorf("error creating credential set: %w", err)
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
			err = fmt.Errorf("error adding '%s' credential in set: %w", cv.Name, err)
			return err
		}
	}

	c.srcs, err = sources.NewSet()
	if err != nil {
		err = fmt.Errorf("error creating source set: %w", err)
		return err
	}
	for _, sv := range cfg.Conf.Sources {
		var workDirName string
		workDirName, err = utils.GenHashString(strings.Join([]string{sv.Name, sv.Type, sv.Encoder}, "."))
		if err != nil {
			err = fmt.Errorf("error generating '%s' source workdir name: %w", sv.Name, err)
			return err
		}

		workDirPath := filepath.Join(cfg.Conf.WorkingDir, workDirName)
		for _, dirv := range []string{"sync"} {
			err = os.MkdirAll(filepath.Join(workDirPath, dirv), utils.DirModePerm)
			if err != nil {
				err = fmt.Errorf("error creating '%s' source workdir directory: %w", sv.Name, err)
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
			err = fmt.Errorf("error adding '%s' source in set: %w", sv.Name, err)
			return err
		}
	}

	c.target.cons, err = conditions.NewSet()
	if err != nil {
		err = fmt.Errorf("error creating condition set: %w", err)
		return err
	}
	for _, cv := range cfg.Conf.Target.Conditions {
		err = c.target.cons.Add(conditions.OptionsT{
			Name:      cv.Name,
			Mandatory: cv.Mandatory,
			Tmpl:      cv.Template,
			Expect:    cv.Expect,
		})
		if err != nil {
			err = fmt.Errorf("error adding '%s' condition in set: %w", cv.Name, err)
			return err
		}
	}

	c.target.acts, err = actions.NewSet()
	if err != nil {
		err = fmt.Errorf("error creating action set: %w", err)
		return err
	}
	for _, av := range cfg.Conf.Target.Actions {
		err = c.target.acts.Add(actions.OptionsT{
			Name:    av.Name,
			Type:    av.Type,
			On:      av.On,
			CredRef: c.creds.Get(av.Credential),
			K8s: actions.OptionsK8sT{
				Namespace: av.K8s.Namespace,
				Pod:       av.K8s.Pod,
				Container: av.K8s.Container,
			},
			Cmd: actions.OptionsCmdT{
				Cmd: av.Cmd,
			},
		})
		if err != nil {
			err = fmt.Errorf("error adding '%s' action in set: %w", av.Name, err)
			return err
		}
	}

	return nil
}
