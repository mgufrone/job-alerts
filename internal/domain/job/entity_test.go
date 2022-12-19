package job

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJob_ID(t *testing.T) {
	t.Parallel()
	validUUID := uuid.NewString()
	testCases := []struct {
		in          string
		out         string
		shouldError bool
	}{
		{"", "", true},
		{"0", "", true},
		{"-1", "", true},

		{"1", "1", true},
		{"12121212alksjdka", "12121212alksjdka", true},
		{"1234-54353-32323", "1234-54353-32323", true},
		{validUUID, validUUID, false},
	}
	for _, c := range testCases {
		j := &Entity{}
		i, _ := uuid.Parse(c.in)
		err := j.SetID(i)
		if c.shouldError {
			require.NotNil(t, err, c)
			continue
		}
		assert.Nil(t, err, c)
		assert.Equal(t, c.out, j.ID().String())
	}
}

func TestJob_Role(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		in          string
		out         string
		shouldError bool
	}{
		{"", "", true},
		{"😁", "", true},

		{"No weird chars åßå", "No Weird Chars", false},
		{"some role", "Some Role", false},
		{"some-dev-role", "Some-Dev-Role", false},
		{"DevOps Engineer", "DevOps Engineer", false},
		{"Technical Support Engineer â\u0080\u0094 GPU Cloud", "Technical Support Engineer GPU Cloud", false},
	}
	for _, c := range testCases {
		j := &Entity{}
		err := j.SetRole(c.in)
		if c.shouldError {
			assert.Equal(t, c.out, j.Role())
			require.NotNil(t, err, c)
			continue
		}
		assert.Nil(t, err, c)
		assert.Equal(t, c.out, j.Role())
	}
}

func TestJob_CompanyName(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		in          string
		out         string
		shouldError bool
	}{
		{"", "", true},
		{"Some company", "Some company", false},
		{"bizdev", "bizdev", false},
		{"Bank Biasalah", "Bank Biasalah", false},
	}
	for _, c := range testCases {
		j := &Entity{}
		err := j.SetCompanyName(c.in)
		if c.shouldError {
			require.NotNil(t, err, c)
			continue
		}
		assert.Nil(t, err, c)
		assert.Equal(t, c.out, j.CompanyName())
	}
}

func TestJob_Source(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		in          string
		out         string
		shouldError bool
	}{
		{"", "", true},
		{"Some company", "Some company", false},
		{"bizdev", "bizdev", false},
		{"Bank Biasalah", "Bank Biasalah", false},
	}
	for _, c := range testCases {
		j := &Entity{}
		err := j.SetSource(c.in)
		if c.shouldError {
			require.NotNil(t, err, c)
			continue
		}
		assert.Nil(t, err, c)
		assert.Equal(t, c.out, j.Source())
	}
}

func TestJob_CompanyURL(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		in          string
		out         string
		shouldError bool
	}{
		{"invalid url", "", true},
		{"https:something.com/somewhere", "", true},
		{"tls://securedomain.tls/somewhere/somehow", "", true},
		{"", "", false},
		{"-", "-", false},
		{"https://somdomain.tld/somepath/somecompany-name", "https://somdomain.tld/somepath/somecompany-name", false},
	}
	for _, c := range testCases {
		j := &Entity{}
		err := j.SetCompanyURL(c.in)
		if c.shouldError {
			require.NotNil(t, err, c)
			continue
		}
		assert.Nil(t, err, c)
		assert.Equal(t, c.out, j.CompanyURL())
	}
}

func TestJob_Description(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		in          string
		out         string
		shouldError bool
	}{
		{"", "", true},
		{"-", "-", false},
		{"invalid url", "invalid url", false},
		{"https://somdomain.tld/somepath/somecompany-name", "https://somdomain.tld/somepath/somecompany-name", false},
	}
	for _, c := range testCases {
		j := &Entity{}
		err := j.SetDescription(c.in)
		if c.shouldError {
			require.NotNil(t, err, c)
			continue
		}
		assert.Nil(t, err, c)
		assert.Equal(t, c.out, j.Description())
	}
}

func TestNewJob(t *testing.T) {
	type inStruct struct {
		id, role, companyName, companyURL, description, jobURL, source, location string
		tags                                                                     []string
	}
	testCases := []struct {
		in          inStruct
		out         *Entity
		shouldError bool
	}{
		{
			inStruct{
				"", "", "", "", "", "", "", "", nil,
			},
			nil,
			true,
		},
		{
			inStruct{
				uuid.Nil.String(), "", "", "", "", "", "", "", nil,
			},
			nil,
			true,
		},
		{
			inStruct{
				uuid.Nil.String(), "role", "", "", "", "", "", "", nil,
			},
			nil,
			true,
		},
		{
			inStruct{
				uuid.Nil.String(), "role", "company", "", "", "", "", "", nil,
			},
			nil,
			true,
		},
		{
			inStruct{
				uuid.Nil.String(), "role", "company", "https://company.com/company/", "", "", "", "", nil,
			},
			nil,
			true,
		},
		{
			inStruct{
				uuid.Nil.String(), "role", "company", "https://company.com/company/", "something wrong", "", "", "", nil,
			},
			nil,
			true,
		},
		{
			inStruct{
				uuid.Nil.String(), "role", "company", "https://company.com/company/", "something wrong",
				"https://somecompany.name/jobs/something", "", "", nil,
			},
			nil,
			true,
		},
	}
	for _, c := range testCases {
		i, _ := uuid.Parse(c.in.id)
		j, err := NewJob(
			i,
			c.in.role,
			c.in.companyName, c.in.companyURL, c.in.description, c.in.jobURL, c.in.source, c.in.location, c.in.tags)
		if c.shouldError {
			require.NotNil(t, err, c)
			continue
		}
		assert.Nil(t, err, c)
		assert.Equal(t, c.out, j)
	}
}
