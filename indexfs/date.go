package indexfs

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const DateRegexString = "[0-9]{4}-[0-9]{2}"

var DateRegex = regexp.MustCompile(fmt.Sprintf(`(?i)%s`, DateRegexString))

type (
	Year  uint16
	Month uint8
	Date  struct {
		Year  Year  `params:"year"`
		Month Month `params:"month"`
	}
)

var ErrNotDate = errors.New("input does not contain a date")

func (d *Date) UnmarshalText(b []byte) error {
	matches := DateRegex.FindAllIndex(b, -1)
	if len(matches) == 0 {
		return ErrNotDate
	}
	match := matches[len(matches)-1]
	b = b[match[0]:match[1]]

	year := b[:4]
	y, err := strconv.ParseUint(string(year), 10, 16)
	if err != nil {
		return err
	}

	month := b[len(b)-2:]
	m, err := strconv.ParseUint(string(month), 10, 8)
	if err != nil {
		return err
	} else if m < 1 || m > 12 {
		return errors.New("month out of range")
	}
	d.Year = Year(y)
	d.Month = Month(m)
	return nil
}

func (d *Date) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d *Date) String() string {
	if d == nil {
		return "0000-00"
	}
	return fmt.Sprintf("%04d-%02d", d.Year, d.Month)
}

func (d *Date) Time() time.Time {
	return time.Date(int(d.Year), time.Month(d.Month), 1, 0, 0, 0, 0, time.Local)
}

func (y Year) String() string {
	return fmt.Sprintf("%04d", y)
}

func (y Year) MarshalText() ([]byte, error) {
	return []byte(y.String()), nil
}

func (y *Year) UnmarshalText(b []byte) error {
	t, err := time.Parse("2006", string(b))
	if err != nil {
		return err
	}
	*y = Year(t.Year())
	return nil
}

func (m Month) String() string {
	return fmt.Sprintf("%d", m)
}

func (m Month) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}

func (m *Month) UnmarshalText(b []byte) error {
	t, err := time.Parse("1", string(b))
	if err != nil {
		return err
	}
	*m = Month(t.Month())
	return nil
}
