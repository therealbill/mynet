---
name: Hugo S3 Deployment
description: >-
  This skill should be used when the user asks about deploying Hugo to AWS S3,
  S3 static site hosting, Hugo with CloudFront, the hugo deploy command,
  configuring Hugo's deployment section, GitHub Actions with AWS OIDC for Hugo,
  S3 bucket static website configuration, CloudFront cache invalidation after
  Hugo deploy, or choosing between GitHub Pages and S3 for Hugo hosting.
  Also applies when GitHub Pages limitations (auth, custom headers, size)
  require an alternative deployment target.
version: 1.0.0
---

## Overview

AWS S3 + CloudFront is an alternative to GitHub Pages for hosting Hugo sites. Choose S3 when you need: custom authentication, multiple sites from one repo, custom HTTP headers, sites larger than GitHub Pages limits (1GB), or deployment to a private CDN.

Hugo has a built-in `hugo deploy` command that syncs the built site to S3 with intelligent caching and content-type handling.

## When to Choose S3 over GitHub Pages

| Factor | GitHub Pages | S3 + CloudFront |
|--------|-------------|-----------------|
| Cost | Free | ~$1-5/month for small sites |
| Setup | Minimal | Moderate (S3, CloudFront, IAM) |
| Custom domain | Yes (with limits) | Yes (full control) |
| HTTPS | Automatic | Via CloudFront |
| Auth/access control | Public only | IAM, signed URLs, WAF |
| Custom headers | No | Yes |
| Size limit | 1GB | Unlimited |
| Multiple sites | One per repo | Unlimited |
| Build location | GitHub Actions only | Any CI |

## Hugo Deploy Configuration

Add to `hugo.toml`:

```toml
[deployment]
  [[deployment.targets]]
    name = "production"
    URL = "s3://my-site-bucket?region=us-east-1"

  [[deployment.matchers]]
    pattern = "^.+\\.(js|css|svg|ttf|woff|woff2)$"
    cacheControl = "max-age=31536000, immutable"
    gzip = true

  [[deployment.matchers]]
    pattern = "^.+\\.(png|jpg|jpeg|gif|webp)$"
    cacheControl = "max-age=31536000, immutable"
    gzip = false

  [[deployment.matchers]]
    pattern = "^.+\\.(html|xml|json)$"
    cacheControl = "max-age=300"
    gzip = true
```

Deploy locally:

```bash
hugo deploy --target production
```

## S3 Bucket Setup

Create an S3 bucket configured for static website hosting:

```bash
# Create bucket
aws s3 mb s3://my-site-bucket --region us-east-1

# Enable static website hosting
aws s3 website s3://my-site-bucket \
  --index-document index.html \
  --error-document 404.html

# Set bucket policy for public access (if not using CloudFront OAI)
aws s3api put-bucket-policy --bucket my-site-bucket --policy '{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Principal": "*",
    "Action": "s3:GetObject",
    "Resource": "arn:aws:s3:::my-site-bucket/*"
  }]
}'
```

## GitHub Actions with OIDC (Recommended)

OIDC federation eliminates long-lived AWS credentials. GitHub Actions authenticates directly with AWS using short-lived tokens.

### IAM Setup

Create an OIDC provider and IAM role (one-time setup):

```bash
# Create OIDC provider for GitHub Actions
aws iam create-open-id-connect-provider \
  --url https://token.actions.githubusercontent.com \
  --client-id-list sts.amazonaws.com \
  --thumbprint-list 6938fd4d98bab03faadb97b34396831e3780aea1
```

Create IAM role trust policy (`trust-policy.json`):

```json
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Principal": {
      "Federated": "arn:aws:iam::ACCOUNT_ID:oidc-provider/token.actions.githubusercontent.com"
    },
    "Action": "sts:AssumeRoleWithWebIdentity",
    "Condition": {
      "StringEquals": {
        "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
      },
      "StringLike": {
        "token.actions.githubusercontent.com:sub": "repo:USERNAME/REPO:ref:refs/heads/master"
      }
    }
  }]
}
```

Create the role with S3 and CloudFront permissions:

```bash
aws iam create-role \
  --role-name hugo-deploy \
  --assume-role-policy-document file://trust-policy.json

aws iam put-role-policy \
  --role-name hugo-deploy \
  --policy-name hugo-deploy-policy \
  --policy-document '{
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Action": ["s3:PutObject", "s3:DeleteObject", "s3:ListBucket"],
        "Resource": [
          "arn:aws:s3:::my-site-bucket",
          "arn:aws:s3:::my-site-bucket/*"
        ]
      },
      {
        "Effect": "Allow",
        "Action": "cloudfront:CreateInvalidation",
        "Resource": "arn:aws:cloudfront::ACCOUNT_ID:distribution/DISTRIBUTION_ID"
      }
    ]
  }'
```

### GitHub Actions Workflow

Create `.github/workflows/hugo-s3-deploy.yml`:

```yaml
name: Deploy Hugo site to S3

on:
  push:
    branches: [master, main]
    paths:
      - 'content/**'
      - 'layouts/**'
      - 'static/**'
      - 'assets/**'
      - 'data/**'
      - 'themes/**'
      - 'hugo.toml'
      - 'go.mod'
      - 'go.sum'
      - '.github/workflows/hugo-s3-deploy.yml'
  workflow_dispatch:

permissions:
  id-token: write
  contents: read

env:
  HUGO_VERSION: '0.142.0'
  AWS_REGION: 'us-east-1'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Install Hugo CLI
        run: |
          wget -O ${{ runner.temp }}/hugo.deb https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_extended_${HUGO_VERSION}_linux-amd64.deb
          sudo dpkg -i ${{ runner.temp }}/hugo.deb

      - name: Checkout
        uses: actions/checkout@v4
        with:
          submodules: recursive
          fetch-depth: 0

      - name: Configure AWS credentials (OIDC)
        uses: aws-actions/configure-aws-credentials@v4
        with:
          role-to-assume: arn:aws:iam::ACCOUNT_ID:role/hugo-deploy
          aws-region: ${{ env.AWS_REGION }}

      - name: Cache Hugo modules
        uses: actions/cache@v4
        with:
          path: |
            ~/.cache/hugo_cache
            /tmp/hugo_cache
          key: ${{ runner.os }}-hugo-${{ hashFiles('go.sum') }}
          restore-keys: |
            ${{ runner.os }}-hugo-

      - name: Build with Hugo
        run: hugo --minify

      - name: Deploy to S3
        run: hugo deploy --target production --maxDeletes 100

      - name: Invalidate CloudFront cache
        run: |
          aws cloudfront create-invalidation \
            --distribution-id ${{ vars.CLOUDFRONT_DISTRIBUTION_ID }} \
            --paths "/*"
```

Store `CLOUDFRONT_DISTRIBUTION_ID` as a repository variable (Settings > Secrets and variables > Actions > Variables), not a secret — it is not sensitive.

## Environment-Based Deployment

Use separate buckets for staging and production:

```toml
[deployment]
  [[deployment.targets]]
    name = "staging"
    URL = "s3://my-site-staging?region=us-east-1"

  [[deployment.targets]]
    name = "production"
    URL = "s3://my-site-production?region=us-east-1"
```

```bash
hugo deploy --target staging
hugo deploy --target production
```

## Common Issues

| Issue | Fix |
|-------|-----|
| Access Denied on deploy | Check IAM role trust policy matches repo/branch |
| OIDC auth fails | Verify `id-token: write` permission in workflow |
| Old content after deploy | Run CloudFront invalidation |
| MIME types wrong | Hugo deploy handles this; check `deployment.matchers` config |
| 403 on site access | Check S3 bucket policy or CloudFront OAI configuration |
