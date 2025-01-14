<img src="./logo.png" height="100" align="right" alt="Instana logo">

# Steadybit extension-instana

A [Steadybit](https://www.steadybit.com/) extension for [Instana](https://www.ibm.com/de-de/products/instana).

Learn about the capabilities of this extension in our [Reliability Hub](https://hub.steadybit.com/extension/com.steadybit.extension_instana).

## Configuration

| Environment Variable            | Helm value | Meaning                                                                                                                                               | Required | Default |
|---------------------------------|------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|----------|---------|
| `STEADYBIT_EXTENSION_BASE_URL`  |            | The Instana Base Url, like `https://$UNIT-$TENANT.instana.io`                                                                                         | yes      |         |
| `STEADYBIT_EXTENSION_API_TOKEN` |            | The Instana [API Token](https://www.ibm.com/docs/en/instana-observability/current?topic=apis-web-rest-api#tokens), see the required permissions below | yes      |         |

The extension supports all environment variables provided by [steadybit/extension-kit](https://github.com/steadybit/extension-kit#environment-variables).

When installed as linux package this configuration is in`/etc/steadybit/extension-instana`.

## Permissions

The extension requires the following scopes:
- "Configuration of Events, Alerts and Smart Alerts for Applications, websites and mobile apps" - `canConfigureCustomAlerts` (if you want to use the "Create Maintenance Window" action)

## Installation

### Kubernetes

Detailed information about agent and extension installation in kubernetes can also be found in
our [documentation](https://docs.steadybit.com/install-and-configure/install-agent/install-on-kubernetes).

#### Recommended (via agent helm chart)

All extensions provide a helm chart that is also integrated in the
[helm-chart](https://github.com/steadybit/helm-charts/tree/main/charts/steadybit-agent) of the agent.

You must provide additional values to activate this extension.

```
--set extension-instana.enabled=true \
--set extension-instana.instana.baseUrl={{YOUR_BASE_URL}} \
--set extension-instana.instana.apiToken={{YOUR_API_TOKEN}} \
```

Additional configuration options can be found in
the [helm-chart](https://github.com/steadybit/extension-instana/blob/main/charts/steadybit-extension-instana/values.yaml) of the
extension.

#### Alternative (via own helm chart)

If you need more control, you can install the extension via its
dedicated [helm-chart](https://github.com/steadybit/extension-instana/blob/main/charts/steadybit-extension-instana).

```bash
helm repo add steadybit-extension-instana https://steadybit.github.io/extension-instana
helm repo update
helm upgrade steadybit-extension-instana \
    --install \
    --wait \
    --timeout 5m0s \
    --create-namespace \
    --namespace steadybit-agent \
    --set instana.baseUrl={{YOUR_BASE_URL}} \
    --set instana.apiToken={{YOUR_API_TOKEN}} \
    steadybit-extension-instana/steadybit-extension-instana
```

### Linux Package

Please use
our [agent-linux.sh script](https://docs.steadybit.com/install-and-configure/install-agent/install-on-linux-hosts)
to install the extension on your Linux machine. The script will download the latest version of the extension and install
it using the package manager.

After installing, configure the extension by editing `/etc/steadybit/extension-instana` and then restart the service.

## Extension registration

Make sure that the extension is registered with the agent. In most cases this is done automatically. Please refer to
the [documentation](https://docs.steadybit.com/install-and-configure/install-agent/extension-registration) for more
information about extension registration and how to verify.
