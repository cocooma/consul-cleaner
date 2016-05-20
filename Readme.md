# About
This is a bit of a rework and extending functionalities of the existing consul-watchdog cli tool.

New features:

- aws host discovery
- consul host discovery
- list services
- list services in specific "any unknown passing warning critical"
- list node status
- list target host/s
- deregister service/s in specific state "any unknown passing warning critical"

# AWS auth
The best way to configure credentials on a development machine is to use the ~/.aws/credentials file, which might look like:

```
[default]
aws_access_key_id = AKID1234567890
aws_secret_access_key = MY-SECRET-KEY
```
Alternatively, you can set the following environment variables:

```
AWS_ACCESS_KEY_ID=AKID1234567890
AWS_SECRET_ACCESS_KEY=MY-SECRET-KEY
```

# Usage
### Cli help
```
Usage of consul-cleaner

  -ar, --awsRegion=eu-west-1       AWS Region. Default: eu-west-1
  -d, --dryrun                     Dryrun
  -drsrv, --deregisterService      Deregister service. Use it with -serviceState
  -fl, --forceLeaveNode            Force leave consul node. Use it with -nodeStatusCode
  -hd, --hostDiscovery=aws         Host discovery. 'consul' or 'aws' or 'stdin'
  -lasrv, --listAllServices        List of all services
  -lchk, --listChecks              List checks
  -lns, --listNodeStatus           List nodes status
  -lsrvis, --listServiceInState    List of services in specific state. Use it with -serviceState
  -lth, --listTargetHosts          List target hosts
  -nsc, --nodeStatusCode=4         Consul node status code. Default: 4
  -p, --port=8500                  Consul members endpoint port. Default: 8500
  -ss, --serviceState=critical     State of the service you wish to deregister. Default: critical
  -t, --tag                        AWS tag
  -tv, --tagValue                  AWS tag value
  -u, --url=localhost              Consul member endpoint. Default: localhost
```


###Example usage
#### passing host/s from stdin
```
echo "11.88.62.55 13.55.22.3" | consul-cleaner -hd stdin -lasrv
```

#### using the consul host discovery
```
consul-cleaner -hd consul -u whateverhostname -p 80 -lasrv
```

#### using the aws host discovery
```
consul-cleaner -hd aws -ar eu-west-1 -t whatevertag -tv whatevertagkey -lasrv
```

# Build the tool
You need to have go installed on the system where you wish to compile it.
For more information in regards go installation please check https://golang.org/doc/install

```
git clone https://github.com/cocooma/consul-cleaner.git
cd devops-tools/consul-cleaner
go get ./...
go build consul-cleaner

```
It should produce an executable called consul-cleaner in your current directory.
