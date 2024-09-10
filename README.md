# stale-flag-detector

A tool to detect stale feature flags in Unleash.

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
- --output-regex: Output the stale flags as a grep-compatible regex

## Example

./stale-flag-detector --output-regex

This will output a regex of all stale flags, which can be used with grep to search your codebase for usage of these flags.

## Contributing

Contributions are welcome. Please open an issue or submit a pull request.

## License

MIT License
