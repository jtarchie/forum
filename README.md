# Forum

Use `rqlite` to provide an edge forum software experience. The idea is the forum
(and software) just needs to be eventually consistent. In the case of a forum,
it means it will be _eventually read_.

## Development

```bash
brew bundle
task
```

## Deployment

These will be instructions on how to deploy to [`fly`](https://fly.io).

```bash
fly launch
fly deploy
```
