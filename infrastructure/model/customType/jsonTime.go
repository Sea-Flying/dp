package customType

import "time"

type dpJsonTime time.Time

const DpTimeLayout = "2006-01-02 15:04:05.000"

func (t *dpJsonTime) UnmarshalJSON(data []byte) error {
	parsed, err := time.Parse(`"`+DpTimeLayout+`"`, string(data))
	if err != nil {
		return err
	}
	*t = dpJsonTime(parsed)
	return nil
}

func (t dpJsonTime) MarshalJSON() ([]byte, error) {
	str := time.Time(t).Format(DpTimeLayout)
	return []byte(str), nil
}

func (t dpJsonTime) String() string {
	return time.Time(t).Format(DpTimeLayout)
}
