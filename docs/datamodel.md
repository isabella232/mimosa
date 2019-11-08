## Mimosa Data Model

Firestore data model: https://firebase.google.com/docs/firestore/data-model

### Tenancy model

The tenancy model is built around two concepts:
* The "account" which represents a human user or an organisation
* The "workspace" which contains user data

Workspaces are the unit of tenancy and are isolated from each other.  Each account has a default workspace and can create additional workspaces to group data in a way that is meaningful to them. A workspace is owned by the account who created it. Accounts can assign permissions to other accounts on a per-workspace basis.

Permissions follow a concentric model:
* A "reader" has read access to the data inside a workspace but cannot run tasks or manage the workspace
* An "executor" can run tasks and has "reader" permissions
* An "admin" can manage sources and permissions associated with a workspace and has "executor" permissions

The workspace "owner" always has "admin" privileges.

#### Permissions example

* Alice, Bob and Cassie create accounts.
* Each has a single default workspace in which they discover inventory data and run tasks.
* Each workspace is completely isolated from the others.
* Alice creates a new worksapce called "projectX". She assigns "executor" permissions to Bob and "reader" permissions to Cassie.
* Bob creates a new worksapce called "projectY". He assigns "admin" permissions to Cassie.
* Permissions now look as follows (ignoring default workspaces):

    | Workspace |  Alice  |    Bob     |  Cassie  |
    | :-------- | :-----: | :--------: | :------: |
    | projectX  | *owner* | *executor* | *reader* |
    | projectY  |    -    |  *owner*   | *admin*  |

* ABC Company creates an account and workspaces for "team1" and "team2".
  * Alice is given "admin" permission for "team1" and "team2" workspaces
  * Bob is given "executor" permission for the "team1" workspace
* Permissions now look as follows (ignoring default workspaces):

    | Workspace |  Alice  |    Bob     |  Cassie  | ABC Company |
    | :-------- | :-----: | :--------: | :------: | :---------: |
    | projectX  | *owner* | *executor* | *reader* |      -      |
    | projectY  |    -    |  *owner*   | *admin*  |      -      |
    | team1     | *admin* | *executor* |    -     |   *owner*   |
    | team2     | *admin* |     -      |    -     |   *owner*   |


#### JWT custom claims

Permissions are captured as custom claims in the JWT allocated at login in the following format:
```json
{
    "owner": [ workspace IDs ... ],
    "admin": [ workspace IDs ... ],
    "executor": [ workspace IDs ... ],
    "reader": [ workspace IDs ... ]
}
```
For example:
```json
{
    "owner": [
        "Odud",
        "PCfF",
        "z4Tv"
    ],
    "admin": [
        "2tCw",
        "Lg91"
    ],
    "executor": [
        "t7c0",
        "YM9O",
        "vAgv",
        "odcC",
        "cO5w"
    ],
    "reader": [
        "cxj7"
    ]
}
```

#### Permission limits

1. An unlimited number of accounts can exist.
2. An unlimited number of accounts can be associated with a single workspace.
3. 16 million workspaces are supported system-wide in total.
4. A single account can be associated with up to 100 workspaces.

Limits have been chosen to reduce near-term engineering cost. Supporting unlimited workspaces and unlimited workspaces per account are both possible with additional engineering work and a scaling plan is in place.

### Document structure

Mimosa lays out documents as follows:

* `/users/<userid>`
* `/ws/<workspaceid>`
* `/ws/<workspaceid>/hosts/<hostid>`
* `/ws/<workspaceid>/tasks/<taskid>`
* `/ws/<workspaceid>/results/<resultid>`

Timestamps are RFC 3339 formatted strings.

#### User document

The document ID is obtained from the `uid` field in the JWT as shown [here](https://firebase.google.com/docs/auth/web/manage-users#get_a_users_profile.

The document contains information about the user's workspaces and has the following fields:

```json
{
    "workspaces": {
        "t7c0": "Team A",
        "YM9O": "Team A",
        "vAgv": "Team B",
        "odcC": "Team D",
        "cO5w": "Team E"
    }
}
```

#### Workspace document

The document ID is obtained from the JWT via the custom claim described above as shown [here](https://firebase.google.com/docs/auth/admin/custom-claims#access_custom_claims_on_the_client).

The document contains information about the workspace and has the following fields:

```json
{
    "name": "Team 1",
}
```

#### Host document

The document ID is determined by the source in a deterministic fashion based on the contents. Clients find it by listing the `hosts` subcollection in the chosen workspace.

The document contains information about the host and has the following fields:

```json
{
    "name": "i-064a45abfde8751ca",
    "hostname": "ec2-54-166-212-236.compute-1.amazonaws.com",
    "ip": "54.166.212.236",
    "source": "source-a1529a68-9e5c-4ac6-8fa7-0e43d5089b9d",
    "state": "running",
    "timestamp": " 2019-11-07T14:39:00Z",
    "tasks": {
        "431687819d0085067de627c7d74def727cc9dee8": {
            "name": "puppetlabs/package",
            "status": "success",
            "timestamp": " 2019-11-07T14:39:00Z",
            "resultid":"9def8bca087006c43c3e2501ac98bf2546fe250d"
     },
        "c43c3e2501ac98bf2546fe250d9def8bca087006": {
            "name": "puppetlabs/service",
            "status": "failure",
            "timestamp": " 2019-11-07T14:39:00Z",
            "resultid":"6c43c3e2501ac98bf2546fe250d9def8bca08700"
        }
    }
}
```

#### Task document

The document ID is allocated by Firestore. Clients find it by listing the `tasks` subcollection in the chosen workspace.

The document contains information about the task and has the following fields:

```json
{
    "name": "puppetlabs/package",
    "params": {
        "name": "openssl",
        "version": "1.2.3",
        "package_manager": "yum"
    },
    "note": "Upgrading a package",
    "timestamp": " 2019-11-07T14:39:00Z",
    "uid": "OdudPCfFz4TvOjuhPEDGl8IAv6s2",
    "displayname": "Alice Ackerman",
    "hosts": {
        "27c7d74def727cc9dee8431687819d0085067de6": {
            "hostname": "ec2-54-166-212-236.compute-1.amazonaws.com",
            "status": "success",
            "timestamp": " 2019-11-07T14:39:00Z",
            "resultid":"9def8bca087006c43c3e2501ac98bf2546fe250d"
      },
        "46fe250d9def8bca087006c43c3e2501ac98bf25": {
            "hostname": "ec2-55-166-212-236.compute-1.amazonaws.com",
            "status": "failure",
            "timestamp": " 2019-11-07T14:39:00Z",
            "resultid":"250d9def8bca087006c43c3e2501ac98bf2546fe"
        }
    }
}
```

#### Results document

The document ID is allocated by Firestore. Clients find it in the `hosts` field of a `task` document or the `tasks` field of a `host` document.

The document contains the raw output from Bolt for a single task on a single host. The format is determined by Bolt.

#### Security rules

```
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {

  	function isOwner(ws) {
      return ws in request.auth.token.owner;
   	}

   	function isAdmin(ws) {
      return ws in request.auth.token.admin;
   	}

   	function isExecutor(ws) {
      return ws in request.auth.token.executor;
   	}

    function isReader(ws) {
      return ws in request.auth.token.reader;
   	}

    match /ws/{ws} {
      allow read, write: if isAdmin(ws) || isOwner(ws)
      allow read: if isExecutor(ws) || isReader(ws)
  	}
  }
}
```
