---
title: "Deploy to AWS S3"
description: "How to deploy a Hugo site to AWS S3 with CloudFront using the /hugo-deploy command and OIDC authentication"
weight: 4
---

# Deploy to AWS S3

This guide covers setting up automated deployment of a Hugo site to AWS S3 with CloudFront CDN using the `/hugo-deploy` command and the hugo-s3-deployment skill.

## When to use this

Choose S3 over GitHub Pages when you need:

- Custom authentication or access control (IAM, signed URLs, WAF)
- Custom HTTP headers or cache control
- Sites larger than the GitHub Pages 1GB limit
- Multiple sites deployed from one repository
- Deployment from any CI system (not just GitHub Actions)

For simpler deployments without these requirements, see {{< ref "howto/deploy-to-github-pages" >}}.

## Prerequisites

- An existing Hugo site in a GitHub repository
- An AWS account with permissions to create S3 buckets, CloudFront distributions, and IAM roles
- AWS CLI installed locally (for initial setup)

## Steps

### 1. Create the S3 bucket

Create an S3 bucket configured for static website hosting:

```bash
aws s3 mb s3://your-site-bucket --region us-east-1

aws s3 website s3://your-site-bucket \
  --index-document index.html \
  --error-document 404.html
```

### 2. Set up CloudFront (recommended)

CloudFront provides HTTPS, global CDN distribution, and custom domain support. Create a CloudFront distribution pointing to the S3 bucket. Note the distribution ID -- you will need it for cache invalidation.

### 3. Set up OIDC authentication

OIDC federation allows GitHub Actions to authenticate with AWS using short-lived tokens instead of stored credentials.

Create the GitHub Actions OIDC provider in AWS:

```bash
aws iam create-open-id-connect-provider \
  --url https://token.actions.githubusercontent.com \
  --client-id-list sts.amazonaws.com \
  --thumbprint-list 6938fd4d98bab03faadb97b34396831e3780aea1
```

Create an IAM role with a trust policy that restricts access to your repository and branch. The role needs `s3:PutObject`, `s3:DeleteObject`, and `s3:ListBucket` permissions on the bucket, plus `cloudfront:CreateInvalidation` on the distribution.

### 4. Add Hugo deployment configuration

Add the deployment target to `hugo.toml`:

```toml
[deployment]
  [[deployment.targets]]
    name = "production"
    URL = "s3://your-site-bucket?region=us-east-1"

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

### 5. Run the `/hugo-deploy` command

```
/hugo-deploy s3
```

The command generates `.github/workflows/hugo-s3-deploy.yml` using the `s3-deploy.yml.tmpl` template.

### 6. Configure repository variables

Add the CloudFront distribution ID as a repository variable (not a secret, since it is not sensitive):

1. Go to your repository Settings on GitHub
2. Navigate to Secrets and variables, then Actions
3. Under Variables, add `CLOUDFRONT_DISTRIBUTION_ID` with your distribution ID

Update the IAM role ARN in the workflow file. Replace `ACCOUNT_ID` with your AWS account ID:

```yaml
role-to-assume: arn:aws:iam::123456789012:role/hugo-deploy
```

### 7. Review the path filter

Ensure all mounted docs directories appear in the workflow's path filter:

```yaml
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
  - 'plugin-a/docs/**'
  - 'plugin-b/docs/**'
```

### 8. Push and verify

Commit and push the workflow:

```bash
git add .github/workflows/hugo-s3-deploy.yml hugo.toml
git commit -m "Add Hugo S3 deployment workflow"
git push
```

Monitor the workflow in the GitHub Actions tab. After successful deployment, verify the site loads through CloudFront.

## Deployment lifecycle

The generated workflow performs these steps in order:

1. Install Hugo Extended edition
2. Check out the repository
3. Authenticate with AWS using OIDC
4. Cache Hugo modules
5. Build the site with `hugo --minify`
6. Sync to S3 with `hugo deploy --target production --maxDeletes 100`
7. Invalidate the CloudFront cache

The `--maxDeletes 100` flag prevents accidental mass deletion. If you need to remove more than 100 files in a single deployment, increase this value or run manually.

## Troubleshooting

- **Access Denied on deploy**: Verify the IAM role trust policy matches your repository name and branch. Check that `id-token: write` is in the workflow permissions.
- **OIDC authentication fails**: Confirm the OIDC provider is created in the correct AWS account and the trust policy's `sub` condition matches your repository.
- **Old content after deploy**: The CloudFront invalidation step should clear the cache. If content is still stale, verify the invalidation completed in the CloudFront console.
- **MIME types incorrect**: Hugo's deploy command sets content types automatically. If types are wrong, check the `deployment.matchers` configuration in `hugo.toml`.
- **403 on site access**: Check the S3 bucket policy allows public read access, or verify the CloudFront Origin Access Identity (OAI) configuration.

## Related

- {{< ref "howto/deploy-to-github-pages" >}} -- Simpler deployment to GitHub Pages
- {{< ref "reference/commands" >}} -- Full `/hugo-deploy` command specification
- {{< ref "reference/templates" >}} -- The `s3-deploy.yml.tmpl` template specification
- {{< ref "explanation/design-decisions" >}} -- Why separate deployment skills for GitHub Pages and S3
