Note: This repository is public and can be found at: https://github.com/igoooor/traefik-maintenance-plugin
Traefik pod will git pull from that public repo. The bitbucket copy is just there for reference 

# Traefik Maintenance Plugin by Conteo

## Configuration

The following declaration (given here in YAML) defines a plugin:

```yaml
# Static configuration

experimental:
  maintenance:
    moduleName: github.com/igoooor/traefik-maintenance-plugin
    version: "v0.0.1" # Grep the latest version 

```

Here is an example of a file provider dynamic configuration (given here in YAML), where the interesting part is the http.middlewares section:

```yaml
# Dynamic configuration

http:
  middlewares:
    maintenance: # Middleware name
      plugin:
        maintenance: # Plugin name
          enabled: true
```