const RUN_COLLECTION = [
  {
    "name": "puppetlabs/package",
    "params": {
      "name": "openssl",
      "version": "1.2.3",
      "package_manager": "yummy"
    },
    "note": "Upgrading a package",
    "timestamp": " 2019-11-07T14:39:00Z",
    "uid": "OdudPCfFz4TvOjuhPEDGl8IAv6s2",
    "displayname": "Alice Ackerman",
    "hosts": {
      "46fe250d9def8bca087006c43c3e2501ac98bf25": {
        "hostname": "ec2-55-166-212-236.compute-1.amazonaws.com",
        "status": "failure",
        "timestamp": " 2019-11-07T14:39:00Z",
        "resultid": "250d9def8bca087006c43c3e2501ac98bf2546fe"
      }
    }
  }, 
  {
    "name": "puppetlabs/swan-dive",
    "params": {
      "name": "lake",
      "place": "somewhere"
    },
    "note": "Super cool diving move",
    "timestamp": " 2019-11-07T14:39:00Z",
    "uid": "OdudPCfFz4TvOjuhPEDGl8IAv6s2",
    "displayname": "Geralt of Rivia",
    "hosts": {
      "27c7d74def727cc9dee8431687819d0085067de6": {
        "hostname": "ec2-54-166-212-236.compute-1.amazonaws.com",
        "status": "success",
        "timestamp": " 2019-11-07T14:39:00Z",
        "resultid": "9def8bca087006c43c3e2501ac98bf2546fe250d"
      },
      "46fe250d9def8bca087006c43c3e2501ac98bf25": {
        "hostname": "ec2-55-166-212-236.compute-1.amazonaws.com",
        "status": "failure",
        "timestamp": " 2019-11-07T14:39:00Z",
        "resultid": "250d9def8bca087006c43c3e2501ac98bf2546fe"
      }
    }
  },
  {
    "name": "puppetlabs/dropshot",
    "params": {
      "name": "a-task-has-no-name",
      "package_manager": "wut"
    },
    "note": "Cheap move",
    "timestamp": " 2019-11-07T14:39:00Z",
    "uid": "OdudPCfFz4TvOjuhPEDGl8IAv6s2",
    "displayname": "Dovakhin",
    "hosts": {
      "27c7d74def727cc9dee8431687819d0085067de6": {
        "hostname": "ec2-54-166-212-236.compute-1.amazonaws.com",
        "status": "success",
        "timestamp": " 2019-11-07T14:39:00Z",
        "resultid": "9def8bca087006c43c3e2501ac98bf2546fe250d"
      },
      "46fe250d9def8bca087006c43c3e2501ac98bf25": {
        "hostname": "ec2-55-166-212-236.compute-1.amazonaws.com",
        "status": "success",
        "timestamp": " 2019-11-07T14:39:00Z",
        "resultid": "250d9def8bca087006c43c3e2501ac98bf2546fe"
      }
    }
  }
]

export default RUN_COLLECTION