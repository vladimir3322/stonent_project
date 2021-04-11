import {IPFS_SOURCE} from 'instances/raribleContract/constants';


export function getImageUrl(imageIPFSURL: string): string | null {
    const IPFSPathMatch = imageIPFSURL.match(/\/ipfs\/.*$/);

    if (!IPFSPathMatch) {
        return null;
    }

    const IPFSPath = IPFSPathMatch[0];

    if (!IPFSPath) {
        return null;
    }

    return `${IPFS_SOURCE}${IPFSPath}`;
}
