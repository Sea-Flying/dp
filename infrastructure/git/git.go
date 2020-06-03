package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os"
	. "voyageone.com/dp/infrastructure/model/global"
	"voyageone.com/dp/infrastructure/utils"
	"voyageone.com/dp/scheduler/model/repository"
	"voyageone.com/dp/scheduler/service"
)

func GitCloneHclTemplate(job repository.DPJob, gitTargetDir string) error {
	var nomadtemplate = repository.NomadTemplate{
		Name: job.NomadTemplateName,
	}
	getNomadTplErr := service.GetNomadTemplateByName(&nomadtemplate)
	if getNomadTplErr != nil {
		return getNomadTplErr
	}
	gitRepoDir := gitTargetDir + "/" + job.NomadTemplateName
	var r *git.Repository
	var gitOpenDirErr error
	var httpAuth = http.BasicAuth{
		Username: DPConfig.Gitlab.Username,
		Password: DPConfig.Gitlab.Token,
	}
	//判断git仓库的文件夹是否存在
	if !utils.FileExists(gitRepoDir) {
		mkdirError := os.MkdirAll(gitRepoDir, os.ModePerm|os.ModeDir)
		if mkdirError != nil {
			return mkdirError
		}
		goto GitClone
	}
	//判断git仓库(.git)是否已经存在
	r, gitOpenDirErr = git.PlainOpen(gitRepoDir)
	if gitOpenDirErr != nil {
		goto GitClone
	} else {
		goto GitPull
	}
GitClone:
	{
		_, gitCloneErr := git.PlainClone(gitRepoDir, false, &git.CloneOptions{
			URL:      nomadtemplate.GitUrl,
			Auth:     &httpAuth,
			Progress: DPLogger.Writer(),
		})
		return gitCloneErr
	}
GitPull:
	{
		w, _ := r.Worktree()
		gitPullErr := w.Pull(&git.PullOptions{
			RemoteName: "origin",
			Auth:       &httpAuth,
			Force:      true,
		})
		if gitPullErr == nil || gitPullErr == git.NoErrAlreadyUpToDate {
			return nil
		} else {
			return gitPullErr
		}
	}
}
