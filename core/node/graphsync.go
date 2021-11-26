package node

import (
	"context"
	"fmt"

	"github.com/ipfs/go-graphsync"
	gsimpl "github.com/ipfs/go-graphsync/impl"
	"github.com/ipfs/go-graphsync/network"
	"github.com/ipfs/go-graphsync/storeutil"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	libp2p "github.com/libp2p/go-libp2p-core"
	peer "github.com/libp2p/go-libp2p-core/peer"

	"github.com/ipld/go-ipld-prime"
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
		fmt.Println(requestData.Root(), requestData.ID())
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
