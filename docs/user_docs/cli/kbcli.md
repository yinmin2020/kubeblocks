## kbcli

KubeBlocks CLI

### Synopsis

```

=============================================
 __    __ _______   ______  __       ______ 
|  \  /  \       \ /      \|  \     |      \
| ▓▓ /  ▓▓ ▓▓▓▓▓▓▓\  ▓▓▓▓▓▓\ ▓▓      \▓▓▓▓▓▓
| ▓▓/  ▓▓| ▓▓__/ ▓▓ ▓▓   \▓▓ ▓▓       | ▓▓  
| ▓▓  ▓▓ | ▓▓    ▓▓ ▓▓     | ▓▓       | ▓▓  
| ▓▓▓▓▓\ | ▓▓▓▓▓▓▓\ ▓▓   __| ▓▓       | ▓▓  
| ▓▓ \▓▓\| ▓▓__/ ▓▓ ▓▓__/  \ ▓▓_____ _| ▓▓_ 
| ▓▓  \▓▓\ ▓▓    ▓▓\▓▓    ▓▓ ▓▓     \   ▓▓ \
 \▓▓   \▓▓\▓▓▓▓▓▓▓  \▓▓▓▓▓▓ \▓▓▓▓▓▓▓▓\▓▓▓▓▓▓

=============================================
A Command Line Interface for KubeBlocks
```

```
kbcli [flags]
```

### Options

```
      --as string                      Username to impersonate for the operation. User could be a regular user or a service account in a namespace.
      --as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --as-uid string                  UID to impersonate for the operation.
      --cache-dir string               Default cache directory (default "$HOME/.kube/cache")
      --certificate-authority string   Path to a cert file for the certificate authority
      --client-certificate string      Path to a client certificate file for TLS
      --client-key string              Path to a client key file for TLS
      --cluster string                 The name of the kubeconfig cluster to use
      --context string                 The name of the kubeconfig context to use
  -h, --help                           help for kbcli
      --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string              Path to the kubeconfig file to use for CLI requests.
      --match-server-version           Require server version to match client version
  -n, --namespace string               If present, the namespace scope for this CLI request
      --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
  -s, --server string                  The address and port of the Kubernetes API server
      --tls-server-name string         Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
      --token string                   Bearer token for authentication to the API server
      --user string                    The name of the kubeconfig user to use
```

### SEE ALSO

* [kbcli backup-config](kbcli_backup-config.md)	 - KubeBlocks backup config
* [kbcli bench](kbcli_bench.md)	 - Run a benchmark
* [kbcli cluster](kbcli_cluster.md)	 - Database cluster operation command
* [kbcli kubeblocks](kbcli_kubeblocks.md)	 - KubeBlocks operation commands
* [kbcli options](kbcli_options.md)	 - Print the list of flags inherited by all commands
* [kbcli playground](kbcli_playground.md)	 - Bootstrap a KubeBlocks in local host
* [kbcli version](kbcli_version.md)	 - Print the version information
