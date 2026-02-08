---
name: hugo-deploy
description: Set up or trigger Hugo site deployment
arguments:
  - name: target
    description: "Deployment target: pages or s3 (default: pages)"
    required: false
---

# Deploy Hugo Site

Set up or trigger deployment to GitHub Pages or AWS S3.

## Process

1. **Check for existing deployment workflow:**
   ```bash
   ls .github/workflows/hugo-*.yml 2>/dev/null
   ```

2. **If no workflow exists:**
   - Ask user: GitHub Pages or S3?
   - Generate the appropriate GitHub Actions workflow
   - For S3: also generate the `[deployment]` section in `hugo.toml`
   - Commit the workflow file
   - Explain next steps (GitHub Pages source setting, or AWS OIDC setup)

3. **If workflow exists:**
   - Validate the workflow configuration
   - Check path filters include all mounted docs directories
   - Verify Hugo version is current
   - Report status and explain how to trigger deployment

## Usage

```
/hugo-deploy [target]
```

## Examples

```
/hugo-deploy          # Check existing or set up GitHub Pages
/hugo-deploy pages    # Set up GitHub Pages deployment
/hugo-deploy s3       # Set up S3 deployment with OIDC
```

## Post-Setup

### GitHub Pages

After creating the workflow:

1. Go to Settings > Pages
2. Set Source to "GitHub Actions"
3. Push a commit that changes content to trigger the first deploy

### S3

After creating the workflow:

1. Set up AWS OIDC provider and IAM role (see hugo-s3-deployment skill)
2. Add `CLOUDFRONT_DISTRIBUTION_ID` as a repository variable
3. Push a commit to trigger the first deploy
