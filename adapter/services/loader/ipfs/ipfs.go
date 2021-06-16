package ipfs

import (
	"context"
	"fmt"
	config "github.com/ipfs/go-ipfs-config"
	files "github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	iCore "github.com/ipfs/interface-go-ipfs-core"
	iCorePath "github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/libp2p/go-libp2p-core/peer"
	peerStore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func setupPlugins(externalPluginsPath string) error {
	plugins, err := loader.NewPluginLoader(filepath.Join(externalPluginsPath, "plugins"))

	if err != nil {
		return fmt.Errorf("error loading plugins: %s", err)
	}
	if err := plugins.Initialize(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}
	if err := plugins.Inject(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	return nil
}

func createTempRepo() (string, error) {
	repoPath, err := ioutil.TempDir("", "ipfs-shell")

	if err != nil {
		return "", fmt.Errorf("failed to get temp dir: %s", err)
	}

	cfg, err := config.Init(ioutil.Discard, 2048)

	if err != nil {
		return "", err
	}

	err = fsrepo.Init(repoPath, cfg)

	if err != nil {
		return "", fmt.Errorf("failed to init ephemeral node: %s", err)
	}

	return repoPath, nil
}

func createNode(ctx context.Context, repoPath string) (iCore.CoreAPI, error) {
	repo, err := fsrepo.Open(repoPath)

	if err != nil {
		return nil, err
	}

	nodeOptions := &core.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTOption,
		Repo:    repo,
	}

	node, err := core.NewNode(ctx, nodeOptions)

	if err != nil {
		return nil, err
	}

	return coreapi.NewCoreAPI(node)
}

func spawnEphemeral(ctx context.Context) (iCore.CoreAPI, error) {
	if err := setupPlugins(""); err != nil {
		return nil, err
	}

	repoPath, err := createTempRepo()

	if err != nil {
		return nil, fmt.Errorf("failed to create temp repo: %s", err)
	}

	return createNode(ctx, repoPath)
}

func connectToPeers(ctx context.Context, ipfs iCore.CoreAPI, peers []string) error {
	var wg sync.WaitGroup

	peerInfos := make(map[peer.ID]*peerStore.PeerInfo, len(peers))

	for _, addrStr := range peers {
		addr, err := ma.NewMultiaddr(addrStr)

		if err != nil {
			return err
		}

		pii, err := peerStore.InfoFromP2pAddr(addr)

		if err != nil {
			return err
		}

		pi, ok := peerInfos[pii.ID]

		if !ok {
			pi = &peerStore.PeerInfo{ID: pii.ID}
			peerInfos[pi.ID] = pi
		}

		pi.Addrs = append(pi.Addrs, pii.Addrs...)
	}

	wg.Add(len(peerInfos))

	for _, peerInfo := range peerInfos {
		go func(peerInfo *peerStore.PeerInfo) {
			defer wg.Done()
			err := ipfs.Swarm().Connect(ctx, *peerInfo)
			if err != nil {
				log.Printf("ipfs peer error %s: %s", peerInfo.ID, err)
			}
		}(peerInfo)
	}

	wg.Wait()
	return nil
}

var ipfsClient iCore.CoreAPI
var ctx = context.Background()

func getClient() iCore.CoreAPI {
	if ipfsClient != nil {
		return ipfsClient
	}

	ipfs, err := spawnEphemeral(ctx)
	ipfsClient = ipfs

	if err != nil {
		fmt.Println(fmt.Sprintf("fail to create ipfs instance: %s, retrying", err))
		return getClient()
	}

	return ipfs
}

func Init() {
	ipfs := getClient()

	bootstrapNodes := []string{
		// IPFS Bootstrapper nodes.
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",

		// IPFS Cluster Pinning nodes
		"/ip4/138.201.67.219/tcp/4001/p2p/QmUd6zHcbkbcs7SMxwLs48qZVX3vpcM8errYS7xEczwRMA",
		"/ip4/138.201.67.219/udp/4001/quic/p2p/QmUd6zHcbkbcs7SMxwLs48qZVX3vpcM8errYS7xEczwRMA",
		"/ip4/138.201.67.220/tcp/4001/p2p/QmNSYxZAiJHeLdkBg38roksAR9So7Y5eojks1yjEcUtZ7i",
		"/ip4/138.201.67.220/udp/4001/quic/p2p/QmNSYxZAiJHeLdkBg38roksAR9So7Y5eojks1yjEcUtZ7i",
		"/ip4/138.201.68.74/tcp/4001/p2p/QmdnXwLrC8p1ueiq2Qya8joNvk3TVVDAut7PrikmZwubtR",
		"/ip4/138.201.68.74/udp/4001/quic/p2p/QmdnXwLrC8p1ueiq2Qya8joNvk3TVVDAut7PrikmZwubtR",
		"/ip4/94.130.135.167/tcp/4001/p2p/QmUEMvxS2e7iDrereVYc5SWPauXPyNwxcy9BXZrC1QTcHE",
		"/ip4/94.130.135.167/udp/4001/quic/p2p/QmUEMvxS2e7iDrereVYc5SWPauXPyNwxcy9BXZrC1QTcHE",
	}

	go connectToPeers(ctx, ipfs, bootstrapNodes)

}

func Get(ipfsPath string) ([]byte, error) {
	fmt.Println(fmt.Sprintf("downloading from ipfs: %s", ipfsPath))

	ipfsId := strings.Split(ipfsPath, "/")[0]
	ipfs := getClient()

	path := iCorePath.New(ipfsPath)
	file, err := ipfs.Unixfs().Get(ctx, path)

	if err != nil {
		return nil, err
	}

	outputPath := fmt.Sprintf("./%s", ipfsId)
	err = files.WriteTo(file, outputPath)

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(outputPath)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := os.Remove(outputPath)

		if err != nil {
			fmt.Println(fmt.Sprintf("error with removing temp IPFS file: %s", err))
		}
	}()

	return data, nil
}
