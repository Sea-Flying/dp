package customType

type DPError string

func (e DPError) Error() string {
	return string(e)
}
