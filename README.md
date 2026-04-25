# vaultwatch

A CLI tool that monitors HashiCorp Vault secret expiration and sends alerts via configurable webhooks.

---

## Installation

```bash
go install github.com/youruser/vaultwatch@latest
```

Or build from source:

```bash
git clone https://github.com/youruser/vaultwatch.git
cd vaultwatch && go build -o vaultwatch .
```

---

## Usage

Configure your Vault address, token, and webhook endpoint in a config file:

```yaml
vault:
  address: "https://vault.example.com"
  token: "s.xxxxxxxxxxxxxxxx"

alerts:
  webhook_url: "https://hooks.slack.com/services/your/webhook/url"
  threshold_days: 7
```

Run vaultwatch against your Vault instance:

```bash
vaultwatch --config config.yaml
```

Monitor a specific secret path:

```bash
vaultwatch --config config.yaml --path secret/data/myapp
```

vaultwatch will poll Vault at a configurable interval and fire a webhook alert when any secret lease is within the expiration threshold.

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--config` | Path to config file | `./config.yaml` |
| `--path` | Vault secret path to monitor | all paths |
| `--interval` | Poll interval (e.g. `5m`, `1h`) | `15m` |
| `--dry-run` | Log alerts without sending webhooks | `false` |

---

## License

[MIT](LICENSE)