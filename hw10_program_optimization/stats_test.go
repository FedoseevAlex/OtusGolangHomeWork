// +build !bench

package hw10programoptimization

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
			 {"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
			 {"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
			 {"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
			 {"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("empty input", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(""), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("no email field", func(t *testing.T) {
		d := `{"Id":1}`
		result, err := GetDomainStat(bytes.NewBufferString(d), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("corrupted json", func(t *testing.T) {
		d := `{ Email: real_mail@orly.com`
		result, err := GetDomainStat(bytes.NewBufferString(d), "com")
		require.Error(t, err)
		require.Equal(t, DomainStat(nil), result)
	})

	t.Run("empty entries", func(t *testing.T) {
		d := `{}
			  {}
	  	      {}`
		result, err := GetDomainStat(bytes.NewBufferString(d), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("two email fields", func(t *testing.T) {
		d := `{"Email":"nulla@Linktype.com", "Email":"nulla@Linktype.ru"}`
		result, err := GetDomainStat(bytes.NewBufferString(d), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)

		result, err = GetDomainStat(bytes.NewBufferString(d), "ru")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"linktype.ru": 1}, result)
	})
}
