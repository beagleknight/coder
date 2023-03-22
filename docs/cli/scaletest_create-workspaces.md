<!-- DO NOT EDIT | GENERATED CONTENT -->

# scaletest create-workspaces

Creates many users and workspace for each one and waits for them to be ready

## Usage

```console
coder scaletest create-workspaces
```

## Description

```console
After creation, waits for the workspaces to finish building and fully come online. Optionally runs a command inside each workspace, and connects to the workspace over WireGuard.

It is recommended that all rate limits are disabled on the server before running this scaletest. This test generates many login events which will be rate limited against the (most likely single) IP.
```

## Options

### --cleanup-concurrency

|             |                                                   |
| ----------- | ------------------------------------------------- |
| Environment | <code>$CODER_SCALETEST_CLEANUP_CONCURRENCY</code> |
| Default     | <code>1</code>                                    |

Number of concurrent cleanup jobs to run. 0 means unlimited.

### --cleanup-job-timeout

|             |                                                   |
| ----------- | ------------------------------------------------- |
| Environment | <code>$CODER_SCALETEST_CLEANUP_JOB_TIMEOUT</code> |
| Default     | <code>5m</code>                                   |

Timeout per job. Jobs may take longer to complete under higher concurrency limits.

### --cleanup-timeout

|             |                                               |
| ----------- | --------------------------------------------- |
| Environment | <code>$CODER_SCALETEST_CLEANUP_TIMEOUT</code> |
| Default     | <code>30m</code>                              |

Timeout for the entire cleanup run. 0 means unlimited.

### --concurrency

|             |                                           |
| ----------- | ----------------------------------------- |
| Environment | <code>$CODER_SCALETEST_CONCURRENCY</code> |
| Default     | <code>1</code>                            |

Number of concurrent jobs to run. 0 means unlimited.

### --connect-hold

|             |                                            |
| ----------- | ------------------------------------------ |
| Environment | <code>$CODER_SCALETEST_CONNECT_HOLD</code> |
| Default     | <code>30s</code>                           |

Time to hold the WireGuard connection open for.

### --connect-interval

|             |                                                |
| ----------- | ---------------------------------------------- |
| Environment | <code>$CODER_SCALETEST_CONNECT_INTERVAL</code> |
| Default     | <code>1s</code>                                |

How long to wait between making requests to the --connect-url once the connection is established.

### --connect-mode

|             |                                            |
| ----------- | ------------------------------------------ |
| Environment | <code>$CODER_SCALETEST_CONNECT_MODE</code> |
| Default     | <code>derp</code>                          |

Connection mode.

### --connect-timeout

|             |                                               |
| ----------- | --------------------------------------------- |
| Environment | <code>$CODER_SCALETEST_CONNECT_TIMEOUT</code> |
| Default     | <code>5s</code>                               |

Timeout for each request to the --connect-url.

### --connect-url

|             |                                           |
| ----------- | ----------------------------------------- |
| Environment | <code>$CODER_SCALETEST_CONNECT_URL</code> |

URL to connect to inside the the workspace over WireGuard. If not specified, no connections will be made over WireGuard.

### -c, --count

|             |                                     |
| ----------- | ----------------------------------- |
| Environment | <code>$CODER_SCALETEST_COUNT</code> |
| Default     | <code>1</code>                      |

Required: Number of workspaces to create.

### --job-timeout

|             |                                           |
| ----------- | ----------------------------------------- |
| Environment | <code>$CODER_SCALETEST_JOB_TIMEOUT</code> |
| Default     | <code>5m</code>                           |

Timeout per job. Jobs may take longer to complete under higher concurrency limits.

### --no-cleanup

|             |                                          |
| ----------- | ---------------------------------------- |
| Environment | <code>$CODER_SCALETEST_NO_CLEANUP</code> |

Do not clean up resources after the load test has finished. Useful for debugging.

### --no-plan

|             |                                       |
| ----------- | ------------------------------------- |
| Environment | <code>$CODER_SCALETEST_NO_PLAN</code> |

Skip the dry-run step to plan the workspace<br/>creation. This step ensures that the given<br/>parameters are valid for the given template.

### --no-wait-for-agents

|             |                                                  |
| ----------- | ------------------------------------------------ |
| Environment | <code>$CODER_SCALETEST_NO_WAIT_FOR_AGENTS</code> |

Do not wait for agents to start before marking<br/>the test as succeeded. This can be useful if you<br/>are running the test against a template that does<br/>not start the agent quickly.

### --output

|             |                                       |
| ----------- | ------------------------------------- |
| Environment | <code>$CODER_SCALETEST_OUTPUTS</code> |
| Default     | <code>text</code>                     |

Output format specs in the format "<format>[:<path>]". Not specifying a path will default to stdout. Available formats: text, json.

### --parameter

|             |                                          |
| ----------- | ---------------------------------------- |
| Environment | <code>$CODER_SCALETEST_PARAMETERS</code> |

Parameters to use for each workspace. Can be specified multiple times. Overrides any existing parameters with the same name from --parameters-file. Format: key=value.

### --parameters-file

|             |                                               |
| ----------- | --------------------------------------------- |
| Environment | <code>$CODER_SCALETEST_PARAMETERS_FILE</code> |

Path to a YAML file containing the parameters to use for each workspace.

### --run-command

|             |                                           |
| ----------- | ----------------------------------------- |
| Environment | <code>$CODER_SCALETEST_RUN_COMMAND</code> |

Command to run inside each workspace using reconnecting-pty (i.e. web terminal protocol). If not specified, no command will be run.

### --run-expect-output

|             |                                                 |
| ----------- | ----------------------------------------------- |
| Environment | <code>$CODER_SCALETEST_RUN_EXPECT_OUTPUT</code> |

Expect the command to output the given string (on a single line). If the command does not output the given string, it will be marked as failed.

### --run-expect-timeout

|             |                                                  |
| ----------- | ------------------------------------------------ |
| Environment | <code>$CODER_SCALETEST_RUN_EXPECT_TIMEOUT</code> |
| Default     | <code>false</code>                               |

Expect the command to timeout. If the command does not finish within the given --run-timeout, it will be marked as succeeded. If the command finishes before the timeout, it will be marked as failed.

### --run-log-output

|             |                                              |
| ----------- | -------------------------------------------- |
| Environment | <code>$CODER_SCALETEST_RUN_LOG_OUTPUT</code> |

Log the output of the command to the test logs. This should be left off unless you expect small amounts of output. Large amounts of output will cause high memory usage.

### --run-timeout

|             |                                           |
| ----------- | ----------------------------------------- |
| Environment | <code>$CODER_SCALETEST_RUN_TIMEOUT</code> |
| Default     | <code>5s</code>                           |

Timeout for the command to complete.

### -t, --template

|             |                                        |
| ----------- | -------------------------------------- |
| Environment | <code>$CODER_SCALETEST_TEMPLATE</code> |

Required: Name or ID of the template to use for workspaces.

### --timeout

|             |                                       |
| ----------- | ------------------------------------- |
| Environment | <code>$CODER_SCALETEST_TIMEOUT</code> |
| Default     | <code>30m</code>                      |

Timeout for the entire test run. 0 means unlimited.

### --trace

|             |                                     |
| ----------- | ----------------------------------- |
| Environment | <code>$CODER_SCALETEST_TRACE</code> |

Whether application tracing data is collected. It exports to a backend configured by environment variables. See: https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/protocol/exporter.md.

### --trace-coder

|             |                                           |
| ----------- | ----------------------------------------- |
| Environment | <code>$CODER_SCALETEST_TRACE_CODER</code> |

Whether opentelemetry traces are sent to Coder. We recommend keeping this disabled unless we advise you to enable it.

### --trace-honeycomb-api-key

|             |                                                       |
| ----------- | ----------------------------------------------------- |
| Environment | <code>$CODER_SCALETEST_TRACE_HONEYCOMB_API_KEY</code> |

Enables trace exporting to Honeycomb.io using the provided API key.

### --trace-propagate

|             |                                               |
| ----------- | --------------------------------------------- |
| Environment | <code>$CODER_SCALETEST_TRACE_PROPAGATE</code> |

Enables trace propagation to the Coder backend, which will be used to correlate server-side spans with client-side spans. Only enable this if the server is configured with the exact same tracing configuration as the client.