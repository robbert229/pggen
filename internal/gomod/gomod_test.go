package gomod

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolvePackage(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{
			path: "Foo.go",
			want: "github.com/robbert229/pggen/internal/gomod",
		},
		{
			path: "../Foo.go",
			want: "github.com/robbert229/pggen/internal",
		},
		{
			path: "./Foo.go",
			want: "github.com/robbert229/pggen/internal/gomod",
		},
		{
			path: "blah/qux/Foo.go",
			want: "github.com/robbert229/pggen/internal/gomod/blah/qux",
		},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got, err := GuessPackage(tt.path)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}

}
