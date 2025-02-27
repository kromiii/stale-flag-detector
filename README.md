# stale-flag-detector

A CLI tool to detect stale feature flags in your unleash project.

```bash
% ./stale-flag-detector
Stale flags:
- unleash-ai-example-stale
- another-stale-flag
```

You can also use the tool in GitHub Actions as follows:

```yaml
- id: stale-flag-detector
  uses: kromiii/stale-flag-detector@v1
  with:
    unleash-api-endpoint: ${{ secrets.UNLEASH_API_ENDPOINT }}
    unleash-api-token: ${{ secrets.UNLEASH_API_TOKEN }}

- id: create-issue
  uses: JasonEtco/create-an-issue@v2
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
    STALE_FLAGS: ${{ steps.stale-flag-detector.outputs.flags }}
```

## Description

This tool connects to an Unleash API and identifies feature flags that have become stale based on their creation date and type. It can be useful for maintaining clean feature flag configurations and identifying flags that may need attention or removal.

## Features

- Detects stale flags based on configurable lifetimes for different flag types
- Option to exclude potentially stale flags
- Can output results as a grep-compatible regex

## Installation

Clone the repository and build the project:

```
git clone https://github.com/kromiii/stale-flag-detector.git
cd stale-flag-detector
go build
```

## Configuration

The tool uses environment variables for configuration. Set the following variables:

- UNLEASH_API_ENDPOINT: The URL of your Unleash API
- UNLEASH_API_TOKEN: Your Unleash API token (**require admin scope**)
- UNLEASH_PROJECT_ID: (Optional) The Unleash project ID (defaults to "default")
- RELEASE_FLAG_LIFETIME: (Optional) Lifetime for release flags in days (defaults to 40)
- EXPERIMENT_FLAG_LIFETIME: (Optional) Lifetime for experiment flags in days (defaults to 40)
- OPERATIONAL_FLAG_LIFETIME: (Optional) Lifetime for operational flags in days (defaults to 7)
- PERMISSION_FLAG_LIFETIME: (Optional) Lifetime for permission flags (set to "permanent" for no expiry)

## Usage

Run the tool with the following command:

./stale-flag-detector [options]

Options:

- --exclude-potentially-stale-flags: Exclude potentially stale flags from the results
- --output-format: Specifies the output format (defaults to "markdown-unordered-list")
  - "markdown-unordered-list": Output as Markdown unordered list
  - "markdown-task-list": Output as Markdown task list
  - "regex": Output as a grep-compatible regex

## Example

```bash
# Output as Markdown unordered list (default)
% ./stale-flag-detector
Stale flags:
- unleash-ai-example-stale
- another-stale-flag

# Output as Markdown task list
% ./stale-flag-detector --output-format=markdown-task-list
Stale flags:
- [ ] unleash-ai-example-stale
- [ ] another-stale-flag

# Output as regex for grep
% ./stale-flag-detector --output-format=regex
(unleash-ai-example-stale|another-stale-flag)
```

## Contributing

Contributions are welcome. Please open an issue or submit a pull request.

## License

MIT License
