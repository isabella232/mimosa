const HOSTS_COLLECTION = [
  {
    "name": "i-host-1234",
    "hostname": "ec11.compute-1.amazonaws.com",
    "ip": "54.166.212.236",
    "source": "mount doom",
    "state": "running",
    "timestamp": " 2019-11-07T14:39:00Z",
    "tasks": {
      "431687819d0085067de627c7d74def727cc9dee8": {
        "name": "puppetlabs/package",
        "status": "success",
        "timestamp": " 2019-11-11T14:39:00Z",
        "resultid": "9def8bca087006c43c3e2501ac98bf2546fe250d"
      },
      "c43c3e2501ac98bf2546fe250d9def8bca087006": {
        "name": "puppetlabs/service",
        "status": "failure",
        "timestamp": " 2019-11-12T14:39:00Z",
        "resultid": "6c43c3e2501ac98bf2546fe250d9def8bca08700"
      }
    }
  },
  {
    "name": "i-smart-abcd",
    "hostname": "ec22.compute-1.amazonaws.com",
    "ip": "54.166.212.237",
    "source": "isengard",
    "state": "terminated",
    "timestamp": " 2017-10-07T14:39:00Z",
    "tasks": {
      "431687819d0085067de627c7d74def727cc9dee8": {
        "name": "puppetlabs/package",
        "status": "failure",
        "timestamp": " 2017-10-08T14:39:00Z",
        "resultid": "9def8bca087006c43c3e2501ac98bf2546fe250d"
      },
      "c43c3e2501ac98bf2546fe250d9def8bca087006": {
        "name": "puppetlabs/service",
        "status": "success",
        "timestamp": " 2019-11-10T14:39:00Z",
        "resultid": "6c43c3e2501ac98bf2546fe250d9def8bca08700"
      }
    }
  },
  {
    "name": "i-spy-xyzp",
    "hostname": "ec33.compute-1.amazonaws.com",
    "ip": "54.166.212.238",
    "source": "hobbiton",
    "state": "stopped",
    "timestamp": " 2018-11-28T14:39:00Z",
    "tasks": {
      "431687819d0085067de627c7d74def727cc9dee8": {
        "name": "puppetlabs/package",
        "status": "success",
        "timestamp": " 2019-11-28T14:39:00Z",
        "resultid": "9def8bca087006c43c3e2501ac98bf2546fe250d"
      },
      "c43c3e2501ac98bf2546fe250d9def8bca087006": {
        "name": "puppetlabs/service",
        "status": "success",
        "timestamp": " 2019-11-29T14:39:00Z",
        "resultid": "6c43c3e2501ac98bf2546fe250d9def8bca08700"
      }
    }
  },
  {
    "name": "i-oh-mi-lmnop",
    "hostname": "ec44.compute-1.amazonaws.com",
    "ip": "54.166.212.239",
    "source": "bree",
    "state": "running",
    "timestamp": " 2019-07-01T14:39:00Z",
    "tasks": {
      "431687819d0085067de627c7d74def727cc9dee8": {
        "name": "puppetlabs/package",
        "status": "failure",
        "timestamp": " 2019-12-02T14:39:00Z",
        "resultid": "9def8bca087006c43c3e2501ac98bf2546fe250d"
      },
      "c43c3e2501ac98bf2546fe250d9def8bca087006": {
        "name": "puppetlabs/service",
        "status": "failure",
        "timestamp": " 2019-11-03T14:39:00Z",
        "resultid": "6c43c3e2501ac98bf2546fe250d9def8bca08700"
      }
    }
  }
]

export default HOSTS_COLLECTION