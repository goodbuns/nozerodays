# nozerodays

reminder bot reminds you if you haven't made a github contribution by 8pm PST


## Testing Locally
Make sure you've set the necessary environment variables. (GITHUB_USERNAME,
GITHUB_ACCESS_TOKEN, ORGANIZATIONS (string of whitelisted organizations,
separated by spaces), WEBHOOK_URL, LOCATION)
```
make run
```

## Deploying to production
Make sure that the following secret exists within the kubernetes cluster:
```
kubectl create secret generic nozerodays --from-literal=username=$USERNAME --from-literal=github-access-token=$GITHUB_ACCESS_TOKEN --from-literal=webhook-url=$WEBHOOK_URL
```

Setting up docker github packages registry credentials:
```
kubectl create secret docker-registry github-packages --docker-server="docker.pkg.github.com" --docker-username=$USERNAME --docker-password=$GITHUB_ACCESS_TOKEN --docker-email=$YOUR_EMAIL
```
