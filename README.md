<p align="center">
    <picture>
        <source media="(prefers-color-scheme: dark)" srcset="https://github.com/privateterraformregistry/privateterraformregistry/raw/main/assets/ptrhero-dark.png">
        <source media="(prefers-color-scheme: light)" srcset="https://github.com/privateterraformregistry/privateterraformregistry/raw/main/assets/ptrhero.jpg)">
        <img src="https://github.com/privateterraformregistry/privateterraformregistry/raw/main/assets/ptrhero.jpg" alt="Caddy" width="450">
    </picture>
</p>

<hr>

<p align="center">
    <a href="#"><img src="https://github.com/privateterraformregistry/privateterraformregistry/actions/workflows/go.yml/badge.svg" /></a>
</p>

<p align="center">
<strong>This project is currently in development and should not be used in production.</strong>
</p>

## What is it?

A private registry you can publish your Terraform modules too. Initially created to manage terraform modules in a mono repo, the roadmap for this project now includes a S3 backend and implementing the terraform login protocol.

## Quickstart

todo.

### Add a Terraform Module

To add a module to the registry, you must compress the directory containing your module.

```
tar -czf file.tar.gz tfmodule
```

Once you have compressed your module, you can publish it to the private terraform registry by sending an HTTP request to the registry endpoint.

Example of POSTing a module to the private terraform registry:
```
curl -X POST \
    -F 'module=@file.tar.gz' \
    https://your.url/modules/hashicorp/consul/aws/1.1.0
```

## Development

Running ```./run.sh``` will run the registry. 

[!warning] Terraform CLI operates over HTTPS, to use HTTPS locally you can use a tool such as NGROK to front your private registry server.

## Testing

Running ```./run.sh``` will run the tests for the registry.
