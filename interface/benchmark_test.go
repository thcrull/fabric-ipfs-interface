package bench

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/thcrull/fabric-ipfs-interface/interface/fabric/wrapper"
	"github.com/thcrull/fabric-ipfs-interface/interface/ipfs/wrapper"
	pb "github.com/thcrull/fabric-ipfs-interface/weightpb"
)

// -------------------------------
// Helpers
// -------------------------------

var aux = 0
var participantId = 10
var aggregatorId = 20

func readVec(path string) ([]float64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(data)%8 != 0 {
		return nil, fmt.Errorf("file %s invalid size", path)
	}

	n := len(data) / 8
	out := make([]float64, n)
	for i := range out {
		bits := binary.LittleEndian.Uint64(data[i*8:])
		out[i] = math.Float64frombits(bits)
	}
	return out, nil
}

func listDataFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	files := []string{}
	for _, e := range entries {
		if !e.IsDir() &&
			strings.HasPrefix(e.Name(), "data_") &&
			strings.HasSuffix(e.Name(), ".bin") {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}
	return files, nil
}

// -------------------------------
// Benchmark Core
// -------------------------------

func runEpoch(
	ctx context.Context,
	meta *fabric_client.MetadataService,
	ipfs *ipfs_client.IpfsClient,
	vec []float64,
	cids *[]string,
) error {

	model := &pb.WeightModel{Values: vec}

	// Add → Pin
	cid, err := ipfs.AddAndPinFile(ctx, model)
	if err != nil {
		return err
	}

	// Track CID for later clean-up
	*cids = append(*cids, cid)

	// Add metadata
	if err := meta.AddParticipantModelMetadata(participantId, aux, cid, "homomorphic-hash"); err != nil {
		return err
	}
	aux++

	return nil
}

func runBenchmark(b *testing.B, file string, epochs int) {
	// -------------------------------
	// Load vector once
	// -------------------------------
	vec, err := readVec(file)
	if err != nil {
		b.Fatalf("read vec: %v", err)
	}

	// -------------------------------
	// Setup Fabric & IPFS
	// -------------------------------
	meta, err := fabric_client.NewMetadataService("../config/admin.yaml")
	if err != nil {
		b.Fatalf("metadata: %v", err)
	}
	defer meta.Close()

	ipfs, err := ipfs_client.NewIpfsClient("../config/admin.yaml")
	if err != nil {
		b.Fatalf("ipfs: %v", err)
	}

	err = meta.AddParticipant(participantId, "enc-key", "hom-shared", "comm-key")
	if err != nil {
		b.Fatalf("add participant: %v", err)
	}

	err = meta.AddAggregator(aggregatorId, map[string]string{fmt.Sprintf("%d", participantId): "comm-key"})
	if err != nil {
		b.Fatalf("add aggregator: %v", err)
	}

	// Track all CIDs created during benchmark
	var createdCids []string

	ctx := context.Background()
	start := time.Now() // measure wall-clock time

	// -------------------------------
	// Actual measured part
	// -------------------------------
	b.ResetTimer()

	for i := 0; i < epochs; i++ {
		if err := runEpoch(ctx, meta, ipfs, vec, &createdCids); err != nil {
			b.Fatalf("epoch %d: %v", i, err)
		}
	}

	b.StopTimer()

	elapsed := time.Since(start).Seconds()
	b.Logf("Total time: %.4f sec", elapsed)

	// -------------------------------
	// Clean-up
	// -------------------------------

	// Delete all participant metadata
	err = meta.DeleteAggregator(aggregatorId)
	if err != nil {
		b.Fatalf("delete aggregator: %v", err)
	}
	err = meta.DeleteParticipant(participantId)
	if err != nil {
		b.Fatalf("delete participant: %v", err)
	}
	err = meta.DeleteAllAggregatorModelMetadata()
	if err != nil {
		b.Fatalf("delete aggregator metadata: %v", err)
	}
	err = meta.DeleteAllParticipantModelMetadata()
	if err != nil {
		b.Fatalf("delete participant metadata: %v", err)
	}

	// Delete all created IPFS pins (double safety)
	for _, cid := range createdCids {
		_ = ipfs.UnpinFile(ctx, cid)
	}
}

// -------------------------------
// Dynamic Benchmark Registration
// -------------------------------

func BenchmarkFabricIPFS(b *testing.B) {
	files, err := listDataFiles("../data")
	if err != nil {
		b.Fatalf("list files: %v", err)
	}

	for _, f := range files {
		f := f // capture

		b.Run(filepath.Base(f)+"_1epoch", func(b *testing.B) {
			b.ReportAllocs()
			runBenchmark(b, f, 1)
		})

		b.Run(filepath.Base(f)+"_100epoch", func(b *testing.B) {
			b.ReportAllocs()
			runBenchmark(b, f, 100)
		})
	}
}
