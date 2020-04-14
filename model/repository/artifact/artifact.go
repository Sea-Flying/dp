package artifact

type Class struct {
	Name        string
	Group       string
	Profile     string
	Kind        string
	RepoName    string
	CreatedTime string
	GitUrl      string
	CiUrl       string
	Description string
}

type Entity struct {
	ClassName     string
	Version       string
	Group         string
	GeneratedTime string
	Profile       string
	RepoName      string
	ClassKind     string
	Uploader      string
	Url           string
	Checksum      string
}

type Repo struct {
	Name         string
	Kind         string
	ArtifactKind string
	Url          string
	BaseUrl      string
	Descprition  string
	CreatedTime  string
}
