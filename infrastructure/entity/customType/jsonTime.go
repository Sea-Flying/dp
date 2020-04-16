package customType

import "time"

type DpJsonTime time.Time

const DpTimeLayout = "2006-01-02 15:04:05.000"

func (t *DpJsonTime) UnmarshalJSON(data []byte) error {
	parsed, err := time.Parse(`"`+DpTimeLayout+`"`, string(data))
	if err != nil {
		return err
	}
	*t = DpJsonTime(parsed)
	return nil
}

func (t DpJsonTime) MarshalJSON() ([]byte, error) {
	str := time.Time(t).Format(DpTimeLayout)
	return []byte(str), nil
}

func (t DpJsonTime) String() string {
	return time.Time(t).Format(DpTimeLayout)
}
