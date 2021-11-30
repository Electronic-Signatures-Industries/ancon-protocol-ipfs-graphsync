package node

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-fetcher"
	"github.com/ipfs/go-graphsync"
	gsimpl "github.com/ipfs/go-graphsync/impl"
	"github.com/ipfs/go-graphsync/network"
	"github.com/ipfs/go-graphsync/storeutil"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	libp2p "github.com/libp2p/go-libp2p-core"
	peer "github.com/libp2p/go-libp2p-core/peer"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
	"github.com/ipld/go-ipld-prime/traversal"
	"github.com/ipld/go-ipld-prime/traversal/selector"
	"github.com/ipld/go-ipld-prime/traversal/selector/builder"

	"github.com/ipld/go-ipld-prime/node/basicnode"
	ipldselector "github.com/ipld/go-ipld-prime/traversal/selector"
	"go.uber.org/fx"

	"github.com/ipfs/go-ipfs/core/node/helpers"
)

// Graphsync constructs a graphsync
func Graphsync(lc fx.Lifecycle, mctx helpers.MetricsCtx, host libp2p.Host, bs blockstore.GCBlockstore) graphsync.GraphExchange {
	ctx := helpers.LifecycleCtx(mctx, lc)

	network := network.NewFromLibp2pHost(host)
	exchange := gsimpl.New(ctx, network,
		storeutil.LinkSystemForBlockstore(bs),
	)

	exchange.RegisterIncomingResponseHook(
		func(p peer.ID, responseData graphsync.ResponseData, hookActions graphsync.IncomingResponseHookActions) {
			fmt.Println(responseData.Status().String(), responseData.RequestID())
		})

	exchange.RegisterIncomingRequestHook(func(p peer.ID, requestData graphsync.RequestData, hookActions graphsync.IncomingRequestHookActions) {
		// var has bool
		// receivedRequestData, has = requestData.Extension(td.extensionName)
		// if !has {
		// 	hookActions.TerminateWithError(errors.New("Missing extension"))
		// } else {
		// 	hookActions.SendExtensionData(td.extensionResponse)
		// }
		hookActions.ValidateRequest()

		has, _ := bs.Has(requestData.Root())
		if !has {
			FetchBlock(ctx, exchange, p, cidlink.Link{Cid: requestData.Root()})
		}
		hookActions.UseLinkTargetNodePrototypeChooser(basicnode.Chooser)
		fmt.Println(requestData.Root(), requestData.ID(), requestData.IsCancel())
	})

	finalResponseStatusChan := make(chan graphsync.ResponseStatusCode, 1)
	exchange.RegisterCompletedResponseListener(func(p peer.ID, request graphsync.RequestData, status graphsync.ResponseStatusCode) {
		select {
		case finalResponseStatusChan <- status:
		default:
		}
	})
	return exchange
}

var selectAll ipld.Node = func() ipld.Node {
	ssb := builder.NewSelectorSpecBuilder(basicnode.Prototype.Any)
	return ssb.ExploreRecursive(
		ipldselector.RecursionLimitDepth(100), // default max
		ssb.ExploreAll(ssb.ExploreRecursiveEdge()),
	).Node()
}()

func FetchBlock(ctx context.Context, gs graphsync.GraphExchange, p peer.ID, c ipld.Link) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	resps, errs := gs.Request(ctx, p, c, selectAll)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case _, ok := <-resps:
			if !ok {
				resps = nil
			}
		case err, ok := <-errs:
			if !ok {
				// done.
				return nil
			}
			if err != nil {
				return fmt.Errorf("got an unexpected error: %s", err)
			}
		}
	}
}
func PushBlock(ctx context.Context, gs graphsync.GraphExchange, p peer.ID, c ipld.Link) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	resps, errs := gs.Request(ctx, p, c, selectAll)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case _, ok := <-resps:
			if !ok {
				resps = nil
			}
		case err, ok := <-errs:
			if !ok {
				// done.
				return nil
			}
			if err != nil {
				return fmt.Errorf("got an unexpected error: %s", err)
			}
		}
	}
}

var matchAllSelector ipld.Node

func init() {
	ssb := builder.NewSelectorSpecBuilder(basicnode.Prototype.Any)
	matchAllSelector = ssb.ExploreRecursive(selector.RecursionLimitNone(), ssb.ExploreUnion(
		ssb.Matcher(),
		ssb.ExploreAll(ssb.ExploreRecursiveEdge()),
	)).Node()
}

// Block fetches a schemaless node graph corresponding to single block by link.
func Block(ctx context.Context, f fetcher.Fetcher, link ipld.Link) (ipld.Node, error) {
	prototype, err := f.PrototypeFromLink(link)
	if err != nil {
		return nil, err
	}
	return f.BlockOfType(ctx, link, prototype)
}

// BlockMatching traverses a schemaless node graph starting with the given link using the given selector and possibly crossing
// block boundaries. Each matched node is sent to the FetchResult channel.
func BlockMatching(ctx context.Context, f fetcher.Fetcher, root ipld.Link, match ipld.Node, cb fetcher.FetchCallback) error {
	return f.BlockMatchingOfType(ctx, root, match, nil, cb)
}

// BlockAll traverses all nodes in the graph linked by root. The nodes will be untyped and send over the results
// channel.
func BlockAll(ctx context.Context, f fetcher.Fetcher, root ipld.Link, cb fetcher.FetchCallback) error {
	return f.BlockMatchingOfType(ctx, root, matchAllSelector, nil, cb)
}

// Fetcher is an interface for reading from a dag. Reads may be local or remote, and may employ data exchange
// protocols like graphsync and bitswap
type Fetcher interface {
	// NodeMatching traverses a node graph starting with the provided root node using the given selector node and
	// possibly crossing block boundaries. Each matched node is passed as FetchResult to the callback. Errors returned
	// from callback will halt the traversal. The sequence of events is: NodeMatching begins, the callback is called zero
	// or more times with a FetchResult, then NodeMatching returns.
	NodeMatching(ctx context.Context, root ipld.Node, selector ipld.Node, cb FetchCallback) error

	// BlockOfType fetches a node graph of the provided type corresponding to single block by link.
	BlockOfType(ctx context.Context, link ipld.Link, nodePrototype ipld.NodePrototype) (ipld.Node, error)

	// BlockMatchingOfType traverses a node graph starting with the given root link using the given selector node and
	// possibly crossing block boundaries. The nodes will be typed using the provided prototype. Each matched node is
	// passed as a FetchResult to the callback. Errors returned from callback will halt the traversal.
	// The sequence of events is: BlockMatchingOfType begins, the callback is called zero or more times with a
	// FetchResult, then BlockMatchingOfType returns.
	BlockMatchingOfType(
		ctx context.Context,
		root ipld.Link,
		selector ipld.Node,
		nodePrototype ipld.NodePrototype,
		cb FetchCallback) error

	// Uses the given link to pick a prototype to build the linked node.
	PrototypeFromLink(link ipld.Link) (ipld.NodePrototype, error)
}

// FetchResult is a single node read as part of a dag operation called on a fetcher
type FetchResult struct {
	Node          ipld.Node
	Path          ipld.Path
	LastBlockPath ipld.Path
	LastBlockLink ipld.Link
}

// FetchCallback is called for each node traversed during a fetch
type FetchCallback func(result FetchResult) error

// Factory is anything that can create new sessions of the fetcher
type Factory interface {
	NewSession(ctx context.Context) Fetcher
}

type fetcherSession struct {
	linkSystem   ipld.LinkSystem
	protoChooser traversal.LinkTargetNodePrototypeChooser
}

// FetcherConfig2 defines a configuration object from which Fetcher instances are constructed
type FetcherConfig2 struct {
	blockService     blockservice.BlockService
	NodeReifier      ipld.NodeReifier
	PrototypeChooser traversal.LinkTargetNodePrototypeChooser
}

// NewFetcherConfig creates a FetchConfig from which session may be created and nodes retrieved.
func NewFetcherConfig(blockService blockservice.BlockService) FetcherConfig2 {
	return FetcherConfig2{
		blockService:     blockService,
		PrototypeChooser: DefaultPrototypeChooser,
	}
}

// NewSession creates a session from which nodes may be retrieved.
// The session ends when the provided context is canceled.
func (fc FetcherConfig2) NewSession(ctx context.Context) fetcher.Fetcher {
	return fc.FetcherWithSession(ctx, blockservice.NewSession(ctx, fc.blockService))
}

func (fc FetcherConfig2) FetcherWithSession(ctx context.Context, s *blockservice.Session) fetcher.Fetcher {
	ls := cidlink.DefaultLinkSystem()
	// while we may be loading blocks remotely, they are already hash verified by the time they load
	// into ipld-prime
	ls.TrustedStorage = true
	ls.StorageReadOpener = blockOpener(ctx, s)
	ls.NodeReifier = fc.NodeReifier

	protoChooser := fc.PrototypeChooser
	return &fetcherSession{linkSystem: ls, protoChooser: protoChooser}
}

// WithReifier derives a different fetcher factory from the same source but
// with a chosen NodeReifier for pathing semantics.
func (fc FetcherConfig2) WithReifier(nr ipld.NodeReifier) fetcher.Factory {
	return FetcherConfig2{
		blockService:     fc.blockService,
		NodeReifier:      nr,
		PrototypeChooser: fc.PrototypeChooser,
	}
}

// interface check
var _ fetcher.Factory = FetcherConfig2{}

// BlockOfType fetches a node graph of the provided type corresponding to single block by link.
func (f *fetcherSession) BlockOfType(ctx context.Context, link ipld.Link, ptype ipld.NodePrototype) (ipld.Node, error) {
	return f.linkSystem.Load(ipld.LinkContext{}, link, ptype)
}

func (f *fetcherSession) nodeMatching(ctx context.Context, initialProgress traversal.Progress, node ipld.Node, match ipld.Node, cb fetcher.FetchCallback) error {
	matchSelector, err := selector.ParseSelector(match)
	if err != nil {
		return err
	}
	return initialProgress.WalkMatching(node, matchSelector, func(prog traversal.Progress, n ipld.Node) error {
		return cb(fetcher.FetchResult{
			Node:          n,
			Path:          prog.Path,
			LastBlockPath: prog.LastBlock.Path,
			LastBlockLink: prog.LastBlock.Link,
		})
	})
}

func (f *fetcherSession) blankProgress(ctx context.Context) traversal.Progress {
	return traversal.Progress{
		Cfg: &traversal.Config{
			LinkSystem:                     f.linkSystem,
			LinkTargetNodePrototypeChooser: f.protoChooser,
		},
	}
}

func (f *fetcherSession) NodeMatching(ctx context.Context, node ipld.Node, match ipld.Node, cb fetcher.FetchCallback) error {
	return f.nodeMatching(ctx, f.blankProgress(ctx), node, match, cb)
}

func (f *fetcherSession) BlockMatchingOfType(ctx context.Context, root ipld.Link, match ipld.Node,
	_ ipld.NodePrototype, cb fetcher.FetchCallback) error {

	// retrieve first node
	prototype, err := f.PrototypeFromLink(root)
	if err != nil {
		return err
	}
	node, err := f.BlockOfType(ctx, root, prototype)
	if err != nil {
		return err
	}

	progress := f.blankProgress(ctx)
	progress.LastBlock.Link = root
	return f.nodeMatching(ctx, progress, node, match, cb)
}

func (f *fetcherSession) PrototypeFromLink(lnk ipld.Link) (ipld.NodePrototype, error) {
	return f.protoChooser(lnk, ipld.LinkContext{})
}

// DefaultPrototypeChooser supports choosing the prototype from the link and falling
// back to a basicnode.Any builder
var DefaultPrototypeChooser = func(lnk ipld.Link, lnkCtx ipld.LinkContext) (ipld.NodePrototype, error) {
	if tlnkNd, ok := lnkCtx.LinkNode.(schema.TypedLinkNode); ok {
		return tlnkNd.LinkTargetNodePrototype(), nil
	}
	return basicnode.Prototype.Any, nil
}

func blockOpener(ctx context.Context, bs *blockservice.Session) ipld.BlockReadOpener {
	return func(_ ipld.LinkContext, lnk ipld.Link) (io.Reader, error) {
		cidLink, ok := lnk.(cidlink.Link)
		if !ok {
			return nil, fmt.Errorf("invalid link type for loading: %v", lnk)
		}

		blk, err := bs.GetBlock(ctx, cidLink.Cid)
		if err != nil {
			return nil, err
		}

		return bytes.NewReader(blk.RawData()), nil
	}
}
