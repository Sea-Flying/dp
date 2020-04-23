package service

import (
	"bytes"
	"path"
	"text/template"
	"voyageone.com/dp/scheduler/model"
)

func RenderJobTemplate(j model.DPJob, tpath string) (jobHcl string, err error) {
	name := path.Base(tpath)
	t := template.Must(template.New(name).Delims("<<", ">>").ParseFiles(tpath))
	buf := bytes.Buffer{}
	if err := t.Execute(&buf, j.NomadTemplateParams); err != nil {
		return "", err
	} else {
		return buf.String(), nil
	}
}
