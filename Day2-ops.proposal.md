This is pool related day-2 ops.

Below day2 ops are considered:
- expansion of current pool wrt capacity and IOPS
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
- Disk CR state changing between 'Active' and 'InActive', i.e., on and off reachability of a disk in cStor pool during run time is handled by cStor

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

Sample SPC yaml is as follows. Except redundantGroupSize, everything else is picked from current SPC yaml.
```
apiVersion: openebs.io/v1alpha1
kind: StoragePoolClaim
metadata:
  name: pool1
  annotations:
    openebs.io/cas-type: cstor
    #(Optional) Use the following to enforce the limits on CPU and RAM to be allocated/used 
    # by the cStor Pool containers. AuxResourceLimits can be used to specify 
    # CPU and RAM limits for the cStor Pool Management side cars
    #If not specified, the default requests and limits will be assigned.
    cas.openebs.io/config: |
      - name: PoolResourceRequests
        value: |-
            memory: 1Gi
            cpu: 100m
      - name: PoolResourceLimits
        value: |-
            memory: 2Gi
      - name: AuxResourceLimits
        value: |-
            memory: 0.5Gi 
            cpu: 50m
spec:
  name: pool1
  #Specify whether cstor pool should be created with disk or sparse files. 
  type: disk
  #Admin can specify the maximum number of cStor pools to be created with this name. 
  maxPools: 3 (default no-limit)
  poolSpec:
    #Define the type of pool to be created. Default is striped. The other supported type is "mirror"
    poolType: striped
    #Tells the number of disks in disk group to create pool1 of type 'striped'
    redundantGroupSize: 1
    #Pools can be configured with different features. 
    #An example feature could be to enable/disable over provisioning.
    overProvisioning: false
    #(Optional - Phase2) Define the required capacity and the max limit
    #capacity: 
    # requests:
    #   storage: 100Gi
    # increment:
    #   storage: 100Gi
    # limits:
    #   storage: 1PBi
  #(Optional) Specify the exact disks or type of disks on which pools 
  # should be created. Disks here refer to the Disk CRs added 
  # by the node-disk-manager. The disks could refer to local disks 
  # or disks from external storage providers like EBS, GPD or SAN.
  disks:
    # Specify exact disks. 
    # If 3 striped pools of single disk, then provide 3 disks from different nodes.
    # If 3 striped pools of two disk, then provide 6 disks with 2 from 3 different nodes.
    diskList:
      - disk-0c84c169ab2f398b92914f56dad41f81
      - disk-66a74896b61c60dcdaf7c7a76fde0ebb
      - disk-b34b3f97840872da9aa0bac1edc9578a
    #(Optiona - Phase2) Specify the type of disks to select from different nodes
    # disks will be selected to specify the requested capacity
    diskTypes:
      - type: block
        capacity:
          minStorage: 10Gi
          maxStorage: 10Ti
  #(Optional - Phase2) Specify the nodes or a list of nodes where the pool has to be created. 
  # by providing a list of node labels
  #nodeSelector: 
    nodetype: storage
Status:
  Phase: Init/Created/Deleted/Errored
Events:

```
SPC watcher will watch for updates to SPC yaml, creates CStorPool CRs.
This CStorPool will be marked for a node where pool need to be created. Only one CSP CR will be marked for a node.
As CSP CR represents a pool on a node, there can't be more than `maxPools` number of CSP CRs for a particular SPC.
Details of CStorPool CR are given in next sections.

Below is sample disk CR:
```

```
API server who watches for updates to above yaml will do following:

(This section mostly relates to pools creation/deletion)


Workflow for adding pools to SPC:
- Increasing maxPools in SPC will allow to create new CSP CRs if already `maxPools` number of are not created. Set maxPools from 3 to 4.
`maxPools: 4`
- Add the disk CRs related to new node
	- If maxPools number of CSP CRs are created, check the workflow for adding pool to SPC where few CSP CRs are marked as 'Deleted'
	- Otherwise, watcher creates CSP CRs based on the given disk CRs of a particular node.
	- Find the disks of a node on which pool of given poolType can be created, and, where CSP is not already created. Create a CSP for that node with all the available disks.


Workflow for removing a pool from SPC:
- Delete all the disk CRS related to a node from which pool need to be destroyed. API server considers it as deleting the pool.
	- Find that all the disk CRs related to a particular node are removed
	- set the status of particular node's CSP as 'Deleted', unsets everything in particular node's CSP (including node UID, CSP UID), deletes the deployment yamls.
	- This shouldn't unset pool structure as this can be needed for node replacement with old disks.

This CSP CR which is marked as 'Deleted' can still be considered for volume creation, if no other pools are available to create volume. Once this CSP CR gets to 'Created' state, volumes as well gets created on it.


Workflow for adding pools with new nodes to SPC, and, few CSP CRs are marked as 'Deleted':
- Add the disk CRs related to new node
	- If already `maxPools` number of CSP CRs are available and are in 'Created', there is nothing to do

	This is the case where `maxPools` number of CSP CRs are created but few of them are marked as 'Deleted'.
	- Find disk CRs related to node on which CSP is not available yet.
	- Find CSP CR which is marked for 'Deleted'
	- set the status of CSP as 'Created', sets its node UID to newly found node UID, create deployment yamls, set the pool structure.

Workflow for adding pools with old nodes to SPC (node replacement with old disks):
	In this case, a CSP CR will be available which will be in 'Deleted' state. Also, disk CRs that are part of this CSP will match with the disk CRs that are added in below workflow.
- remove the disk CRs that are part of old node, so that, workflow of removing pool from SPC will get kicked in.
- make sure that CSP is marked as 'Deleted' and deployments are deleted
- NDM detects the same disks in new node with same disk CR names
- add those disk CRS (which are removed in step 1) in SPC. This causes import of the pool on these disks of new node. (Disk CRs of old node need to be removed, and same disk CRs need to be added becz API server will not see any change in the SPC configuration. No visible change will be there in SPC becz NDM will give same disk CR for same disk on different node. But, disk CR will be updated with node UID. If API server can see the update on disk CR, then, this workflow will not be required, and it can be made automatic.) (This have dependency on cStor to make sure that pool can be imported.)

Workflow for adding disk CRs to SPC (pool capacity expansion):
- Add the disk CRs of a node on which CSP is available.
	- Finds the newly added disk on which CSP is available
	- Makes sure that all the disks in CSP strucure are in 'Online' state. If any disks are in 'Offline' state, check the next workflows.
	- Add a new group to pool structure of this CSP to increase the total capacity

Workflow for removing disks from SPC for temporary purpose like F/W upgrade:
- Remove the disk CRS of a node that need to be removed from pool
	- Find that disk CRs of a node that are removed and its CSP
	- Mark the disks as 'Offline' in CSP (this should stop pool reading/writing to this disk, and even during import time)
	If multiple disk CRs are removed, and, if it ends up in a case where entrire group is not usable, IOs can get suspended.


Workflow for removing few disks from SPC for permanent purpose like faulted drives:
- Do delete of disk CR from etcd
- Remove those disk CRs from yaml
	- Find that disk CRs of a node that are removed and its CSP
	- Mark the disks as 'Offline' in CSP (this should stop pool reading/writing to this disk, and even during import time)
	If multiple disk CRs are removed, and, if it ends up in a case where entrire group is not usable, change the pool structure given in CSP, so that, new pool will be created on remaining nodes.

Workflow for adding disk CRs to SPC in which few disk CRs are in 'Offline' state (replacing faulted disks):
- Add the disk CRs of a node on which CSP is available.
	As few disk CRs are in 'Offline' state, this newly added disk CRs will replace disk CRs that are in 'Offline' state in CSP
	- Find the newly added disk on which CSP is available
	- Find the disk CRs that are in 'Offline' state in CSP and can be replaced with newly added disks. Here, priority need to be given to find for the newly added disk CRs that are in 'Offline' state.
	- Replace those 'Offline' disk CRs with newly added disk CRs
	- Set the state of newly added disk CRs as 'Online'
	In above, covered the case where a disk CR has been removed and same added back. Missing data will be written to it by automatic resilvering.


CSP yaml:
API server creates this intent for each specific node by using nodeUID.
This contains intent about node for which this CSP applies, pool structure along with disk CRs.
This will have an UniqueID which represents this pool.
This will have link to SPC UID.

Below is sample:
```
apiVersion: openebs.io/v1alpha1
kind: CStorPool
metadata:
  #Name is auto generated using the prefix of StoragePoolClaim name and 
  # nodename hash
  name: pool1-84eb2e
  #Following uid will be auto generated when the CR is created.
  uid: 7b99e406-1260-11e8-aa43-00505684eb2e
  labels:
    "kubernetes.io/hostname": "node-host-label"
    "kubernetes.io/hostuid": "node-host-uid"
    openebs.io/storage-pool-claim: pool1
    openebs.io/storage-pool-claim-uid: spc-uid
spec:
  disks:
    #Disks that are actually used for creating the cstor pool are listed here. 
    diskStructure: 
      - group1:
          - disk-0c84c169ab2f398b92914f56dad41f81
            - status: online
          - disk-0c84c169ab2f398b92914f56dad41f82
            - status: offline
      - group2:
          - disk-1c84c169ab2f398b92914f56dad41f81
            - status: online
          - disk-1c84c169ab2f398b92914f56dad41f82
            - status: online
  poolSpec: 
    #Pool features as passed from the SPC.
    #Defines the type of pool as passed from the SPC. stripe or mirror. 
    poolType: "mirror"
    #overProvisioning: false       
# status is updated by the cstor-pool-mgmt to reflect the current status of the pool. 
# The valid values are : init, online, offline
status:
  phase: init
  pooluid: uid
```

pool-mgmt verifies that it picked the CSP that matches its node UID.
pool-mgmt verifies that all its disk CRs are existing, and their node UID matching with the node ID on which pool-mgmt resides.
pool-mgmt reads the devices's devlink id of all disk CRs (not to worry about those marked as 'Offline').
pool-mgmt first attempts import by using cache file and also tries to import using the pool UID and the directory in which disks resides (get the directory path from devlink)(should we consider of importing only when status.phase is NOT init? If status.phase is init, can we considering recreating the pool to avoid the case of importing with wrong disks? Otherwise, clearing label or format by NDM would help to address this?)
If import is not successful, pool-mgmt creates the pool using the pool structure in CSP.
If import is successful, pool-mgmt verifies the disk structure in CSP with the disk structure of pool available, and, any errors in disk status.

If any change in structure leading to other than disk operations, then, it will end up in creating pool.
- If disk CRs are added and no disks are in 'Offline' state, this will be considered as capacity expansion, and, 'zpool add' will be done.
- If disk CRs are added and few disks are in 'Offline' state, 
- If disk CRs that are 'Offline' and doesn't even exists, 




