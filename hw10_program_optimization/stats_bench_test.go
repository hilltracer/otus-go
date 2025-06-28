package hw10programoptimization

import (
	"archive/zip"
	"testing"

	"github.com/stretchr/testify/require"
)

// go test -tags bench -run=^$ -bench=BenchmarkGetDomainStat -benchtime=20x \
// -cpuprofile=cpu.pprof -memprofile=mem.pprof -memprofilerate=1.
func BenchmarkGetDomainStat(b *testing.B) {
	zr, err := zip.OpenReader("testdata/users.dat.zip")
	require.NoError(b, err)
	defer zr.Close()

	require.Equal(b, 1, len(zr.File))

	b.ResetTimer() // reset timer to exclude setup time

	for i := 0; i < b.N; i++ {
		data, err := zr.File[0].Open()
		require.NoError(b, err)

		_, err = GetDomainStat(data, "biz")
		require.NoError(b, err)

		_ = data.Close()
	}
}
