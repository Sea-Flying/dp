package git

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os"
	"strings"
	"sync"
	. "voyageone.com/dp/infrastructure/model/global"
	"voyageone.com/dp/infrastructure/utils"
	"voyageone.com/dp/scheduler/model/repository"
)

func CloneHclTemplate(job repository.DPJob, gitTargetDir string) (string, error) {
	var nomadTemplate = repository.NomadTemplate{
		Name: job.NomadTemplateName,
	}
	getNomadTplErr := nomadTemplate.GetByName()
	if getNomadTplErr != nil {
		return "", getNomadTplErr
	}
	splits := strings.Split(nomadTemplate.GitUrl, "/")
	floderName := strings.Split(splits[len(splits)-1], ".git")[0]
	gitRepoDir := gitTargetDir + "/" + floderName
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
			return gitRepoDir, mkdirError
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
		if _, ok := GitMutex[gitRepoDir]; !ok {
			GitMutex[gitRepoDir] = &sync.Mutex{}
		}
		GitMutex[gitRepoDir].Lock()
		defer GitMutex[gitRepoDir].Unlock()
		_, gitCloneErr := git.PlainClone(gitRepoDir, false, &git.CloneOptions{
			URL:      nomadTemplate.GitUrl,
			Auth:     &httpAuth,
			Progress: DPLogger.Writer(),
		})
		return gitRepoDir, gitCloneErr
	}
GitPull:
	{
		w, _ := r.Worktree()
		if _, ok := GitMutex[gitRepoDir]; !ok {
			GitMutex[gitRepoDir] = &sync.Mutex{}
		}
		GitMutex[gitRepoDir].Lock()
		defer GitMutex[gitRepoDir].Unlock()
		gitPullErr := w.Pull(&git.PullOptions{
			RemoteName: "origin",
			Auth:       &httpAuth,
			Force:      true,
		})
		if gitPullErr == nil || errors.Is(gitPullErr, git.NoErrAlreadyUpToDate) {
			return gitRepoDir, nil
		} else {
			return gitRepoDir, gitPullErr
		}
	}
}
