package models

import (
	"github.com/stretchr/testify/assert"
	test_data2 "mgufrone.dev/job-alerts/pkg/test_data"
	"testing"
)

func TestJob_TransformFail(t *testing.T) {
	t.Parallel()
	var j Job
	jb, err := j.Transform()
	assert.Nil(t, jb)
	assert.NotNil(t, err)
}
func TestJob_FromDomainOk(t *testing.T) {
	t.Parallel()
	jb := test_data2.ValidJob()
	var j Job
	j.FromDomain(jb)
	jb2, err := j.Transform()
	assert.Equal(t, jb2, jb)
	assert.Nil(t, err)
}
