# Repository-Specific GitHub User Setup

This guide explains how to use a different GitHub user for this specific repository.

## Quick Setup

Run the provided setup script with your GitHub username and email:

```sh
./setup-git-user.sh "Your Name" "your.email@example.com"
```

## Manual Configuration

If you prefer to set it up manually:

```sh
# Set repository-specific user
git config --local user.name "Your Name"
git config --local user.email "your.email@example.com"

# Verify configuration
git config --local --get user.name
git config --local --get user.email
```

## SSH Key Configuration

If your different GitHub user requires a different SSH key:

1. Add your SSH key to the SSH agent:

```sh
# Start SSH agent if not running
eval "$(ssh-agent -s)"

# Add your key
ssh-add ~/.ssh/your_different_github_key
```

2. Or create a repository-specific SSH configuration:

Create or edit `~/.ssh/config`:

```
Host github.com-alternate
    HostName github.com
    User git
    IdentityFile ~/.ssh/your_different_github_key
```

Then update your repository's remote URL:

```sh
git remote set-url origin git@github.com-alternate:username/hcloud_exporter.git
```

## HTTPS Authentication

If using HTTPS instead of SSH, you can store credentials:

```sh
# Store credentials for this repository
git config --local credential.helper store

# Or cache them temporarily (e.g., for 1 hour)
git config --local credential.helper 'cache --timeout=3600'
```

The next time you push, you'll be prompted for your username and password/token.

## GitHub CLI Configuration

If using GitHub CLI (gh), you can switch accounts:

```sh
gh auth login -h github.com -w
```

This will let you authenticate with a different account for operations in this repository.
