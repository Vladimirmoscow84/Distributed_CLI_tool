package model

type GrepConfig struct {
	Pattern    string
	IgnoreCase bool
	ShowNumber bool
	Invert     bool
}

type GrepResult struct {
	Lines []string
}

type Shard struct {
	ID   int
	Data []byte
}

type ShardResult struct {
	ID    int
	Lines []string
}

type ClusterConfig struct {
	Quorum int
	Peers  []string
}
