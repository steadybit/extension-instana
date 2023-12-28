<img src="./logo.png" height="100" align="right" alt="Instana logo">

# Steadybit extension-instana

A [Steadybit](https://www.steadybit.com/) extension for [Instana](https://www.ibm.com/de-de/products/instana).

Learn about the capabilities of this extension in our [Reliability Hub](https://hub.steadybit.com/extension/com.steadybit.extension_instana).

## Configuration

| Environment Variable            | Helm value | Meaning                                                                                                                                               | Required | Default |
|---------------------------------|------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|----------|---------|
| `STEADYBIT_EXTENSION_BASE_URL`  |            | The Instana Base Url, like `https://$UNIT-$TENANT.instana.io`                                                                                         | yes      |         |
| `STEADYBIT_EXTENSION_API_TOKEN` |            | The Instana [API Token](https://www.ibm.com/docs/en/instana-observability/current?topic=apis-web-rest-api#tokens), see the required permissions below | yes      |         |

The extension supports all environment variables provided
by [steadybit/extension-kit](https://github.com/steadybit/extension-kit#environment-variables).

## Installation

### Using Docker

```sh
docker run \
  --rm \
  -p 8090 \
  --name steadybit-extension-instana \
  ghcr.io/steadybit/extension-instana:latest
```

### Using Helm in Kubernetes

```sh
helm repo add steadybit-extension-instana https://steadybit.github.io/extension-instana
helm repo update
helm upgrade steadybit-extension-instana \
    --install \
    --wait \
    --timeout 5m0s \
    --create-namespace \
    --namespace steadybit-extension \
    steadybit-extension-instana/steadybit-extension-instana
```

## Register the extension

Make sure to register the extension at the steadybit platform. Please refer to
the [documentation](https://docs.steadybit.com/integrate-with-steadybit/extensions/extension-installation) for more information.
