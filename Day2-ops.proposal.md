This is pool related day-2 ops.

Below day2 ops are considered:
- expansion of current pool
- disk replacement with current pool being intact

Below ones involve re-creating of new pools and new volumes:
- replacing old node with new node and new disks
- replacing old disks in a node with new disks of same node

Below ones involve importing of pool (not creating new pool):
- replacing old node with new node and old disks

Below ones are related to (upward) horizontal scaling of pool pods:
- Add new pool to SPC for later volume provisioning

Below ones are related to deleting pool pod (this doesn't handle about data consisteny):
(this is required in replacing with new node)
- Remove a pool pod from SPC

Not considered:
- Downscaling of pool pods which involves downscaling of volumes (?)

Assumptions or dependencies on other components:
cStor:
- Does automatic resilvering of pools due to momentary loss of disks in a pool
- Allows to import pools with cache file due to memontary loss of atleast one full redundant group and with changes in device's dev link-id
- Handles rebuilding of data in re-created volumes
- Allows to import pools even without cache file and with directory path of disks and with changes in device's dev link-id

NDM:
- un-usability of disk CR by changing its status to 'InActive' if the disk shouldn't be used for pools
- usability of disk CR by changing its status to 'Active' if the disk should be used for pools
- gives same disk CR for same disk even if the disk moves across different nodes (ephemeral disks also)
- give different disk CR for different disk even if the disk is in same node

SPC yaml:
Admin creates this intent.
This contains intent on the number of pools, nodes on which pool need to be available, disks that need to be part of pool, poolType, redundantGroupSize, overprovisiong (?)(this need to be part of volume)
Various fields are:
maxPools (change name to poolGroupSize) - number of pools, which is nothing but number of nodes
DiskList - is list of Disk CRs - gives info about nodes on which pool need to be created, and also, the disks that need to be part of this pool
poolType - can take striped/mirror/raidz/raidz2
redundantGroupSize - gives the number of disks that need to be part of diskGroup. Valid for raidz and raidz2
This will have an UniqueID which represents this SPC.

API server who watches for changes in this intent will do following:
- If maxPools number of CSP CRs are created, it won't add new CSP CRs
- If maxPools number of CSP CRs are not created, it creates one based on the given disk CRs of a particular node.
- If all the disk CRs related to a particular node are removed, API server considers it as deleting the pool from SPC. In this case, API server sets the status of particular node's CSP as 'Deleted', unsets everything in particular node's CSP (including node UID, CSP UID), deletes the deployment yamls. This shouldn't unset pool structure as this can be needed for node replacement and old disks.
- If maxPools number of CSP CRs are created but few of them are marked as 'Deleted', it finds for any disk CRs related to node on which CSP is not available yet. It resets the 'Deleted' state of CSP, sets its node UID to newly found node UID, creates deployment yamls. This is related to replacing old node with new node and old disks.
  - First priority need to be given to node replacement with old disks, by searching for disk CRs that matches to any CSP's structure in 'Deleted' state, rather than creating new pool on any new nodes. By doing this way, it just imports old pool on new node.
- Workflow for node replacement with old disks goes like:
  - remove the disk CRs that are part of old node
  - makes sure that CSP is marked as 'Deleted' and deployments are deleted
  - add the same disk CRs that are removed. This causes import of the pool on these disks. (Disk CRs of old node need to be removed, and same disk CRs need to be added becz API server will not see any change in the SPC configuration. No visible change will be there in SPC becz NDM will give same disk CR for same disk on different node. But, disk CR will be updated with node UID. If API server can see the update on disk CR, then, this workflow will not be required, and it can be made automatic.)


API server converts this to CSP yamls for each node.



CSP yaml:
API server creates this intent for each specific node by using nodeUID.
This contains intent about node for which this CSP applies, pool structure along with disk CRs
This will have an UniqueID which represents this pool.
This will have link to SPC UID.





