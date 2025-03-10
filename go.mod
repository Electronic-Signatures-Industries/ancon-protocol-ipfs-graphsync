module github.com/ipfs/go-ipfs

require (
	bazil.org/fuse v0.0.0-20200117225306-7b5117fecadc
	contrib.go.opencensus.io/exporter/prometheus v0.3.0
	github.com/blang/semver/v4 v4.0.0
	github.com/cheggaaa/pb v1.0.29
	github.com/coreos/go-systemd/v22 v22.3.2
	github.com/dustin/go-humanize v1.0.0
	github.com/elgris/jsondiff v0.0.0-20160530203242-765b5c24c302
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gabriel-vasile/mimetype v1.1.2
	github.com/go-bindata/go-bindata/v3 v3.1.3
	github.com/hashicorp/go-multierror v1.1.1
	github.com/ipfs/go-bitswap v0.4.0
	github.com/ipfs/go-block-format v0.0.3
	github.com/ipfs/go-blockservice v0.1.7
	github.com/ipfs/go-cid v0.0.7
	github.com/ipfs/go-cidutil v0.0.2
	github.com/ipfs/go-datastore v0.4.6
	github.com/ipfs/go-detect-race v0.0.1
	github.com/ipfs/go-ds-badger v0.2.7
	github.com/ipfs/go-ds-flatfs v0.4.5
	github.com/ipfs/go-ds-leveldb v0.4.2
	github.com/ipfs/go-ds-measure v0.1.0
	github.com/ipfs/go-fetcher v1.5.0
	github.com/ipfs/go-filestore v0.0.3
	github.com/ipfs/go-fs-lock v0.0.7
	github.com/ipfs/go-graphsync v0.9.3
	github.com/ipfs/go-ipfs-blockstore v0.1.6
	github.com/ipfs/go-ipfs-chunker v0.0.5
	github.com/ipfs/go-ipfs-cmds v0.6.0
	github.com/ipfs/go-ipfs-config v0.16.0
	github.com/ipfs/go-ipfs-exchange-interface v0.0.1
	github.com/ipfs/go-ipfs-exchange-offline v0.0.1
	github.com/ipfs/go-ipfs-files v0.0.9
	github.com/ipfs/go-ipfs-keystore v0.0.2
	github.com/ipfs/go-ipfs-pinner v0.1.2
	github.com/ipfs/go-ipfs-posinfo v0.0.1
	github.com/ipfs/go-ipfs-provider v0.6.1
	github.com/ipfs/go-ipfs-routing v0.1.0
	github.com/ipfs/go-ipfs-util v0.0.2
	github.com/ipfs/go-ipld-cbor v0.0.5
	github.com/ipfs/go-ipld-format v0.2.0
	github.com/ipfs/go-ipld-git v0.1.1
	github.com/ipfs/go-ipld-legacy v0.1.0
	github.com/ipfs/go-ipns v0.1.2
	github.com/ipfs/go-log v1.0.5
	github.com/ipfs/go-merkledag v0.4.0
	github.com/ipfs/go-metrics-interface v0.0.1
	github.com/ipfs/go-metrics-prometheus v0.0.2
	github.com/ipfs/go-mfs v0.1.2
	github.com/ipfs/go-namesys v0.3.1
	github.com/ipfs/go-path v0.1.2
	github.com/ipfs/go-pinning-service-http-client v0.1.0
	github.com/ipfs/go-unixfs v0.2.5
	github.com/ipfs/go-unixfsnode v1.1.3
	github.com/ipfs/go-verifcid v0.0.1
	github.com/ipfs/interface-go-ipfs-core v0.5.1
	github.com/ipfs/tar-utils v0.0.1
	github.com/ipld/go-car v0.3.2
	github.com/ipld/go-codec-dagpb v1.3.0
	github.com/ipld/go-ipld-prime v0.12.3
	github.com/jbenet/go-random v0.0.0-20190219211222-123a90aedc0c
	github.com/jbenet/go-temp-err-catcher v0.1.0
	github.com/jbenet/goprocess v0.1.4
	github.com/libp2p/go-doh-resolver v0.3.1
	github.com/libp2p/go-libp2p v0.15.0
	github.com/libp2p/go-libp2p-circuit v0.4.0
	github.com/libp2p/go-libp2p-connmgr v0.2.4
	github.com/libp2p/go-libp2p-core v0.9.0
	github.com/libp2p/go-libp2p-discovery v0.5.1
	github.com/libp2p/go-libp2p-http v0.2.1
	github.com/libp2p/go-libp2p-kad-dht v0.13.1
	github.com/libp2p/go-libp2p-kbucket v0.4.7
	github.com/libp2p/go-libp2p-loggables v0.1.0
	github.com/libp2p/go-libp2p-mplex v0.4.1
	github.com/libp2p/go-libp2p-noise v0.2.2
	github.com/libp2p/go-libp2p-peerstore v0.2.8
	github.com/libp2p/go-libp2p-pubsub v0.5.4
	github.com/libp2p/go-libp2p-pubsub-router v0.4.0
	github.com/libp2p/go-libp2p-quic-transport v0.12.0
	github.com/libp2p/go-libp2p-record v0.1.3
	github.com/libp2p/go-libp2p-routing-helpers v0.2.3
	github.com/libp2p/go-libp2p-swarm v0.5.3
	github.com/libp2p/go-libp2p-testing v0.4.2
	github.com/libp2p/go-libp2p-tls v0.2.0
	github.com/libp2p/go-libp2p-yamux v0.5.4
	github.com/libp2p/go-socket-activation v0.1.0
	github.com/libp2p/go-tcp-transport v0.2.8
	github.com/libp2p/go-ws-transport v0.5.0
	github.com/miekg/dns v1.1.43
	github.com/mitchellh/go-homedir v1.1.0
	github.com/multiformats/go-multiaddr v0.4.0
	github.com/multiformats/go-multiaddr-dns v0.3.1
	github.com/multiformats/go-multibase v0.0.3
	github.com/multiformats/go-multicodec v0.3.0
	github.com/multiformats/go-multihash v0.0.15
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/statsd_exporter v0.21.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.0
	github.com/whyrusleeping/go-sysinfo v0.0.0-20190219211824-4a357d4b90b1
	github.com/whyrusleeping/multiaddr-filter v0.0.0-20160516205228-e903e4adabd7
	go.opencensus.io v0.23.0
	go.uber.org/fx v1.13.1
	go.uber.org/zap v1.19.0
	golang.org/x/crypto v0.0.0-20210813211128-0a44fdfbc16e
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210816183151-1e6c022a8912
)

go 1.16
