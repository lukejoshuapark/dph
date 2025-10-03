# 🐗 dph

`dph` (Declarative PostHog) syncs the keys and descriptions of your feature
flags to match a provided file.

The life-cycle of your feature flags become coupled with the deployment of your
code, but enabling/disabling them is controlled via the UI.

## Usage

`dph` will:

- Create flags in PostHog that exist in the supplied file but do not exist in
PostHog.

- Delete flags in PostHog that do not exist in the supplied file.

- Update the description of flags where the description in PostHog does not
match the supplied file.

Basic usage is as simple as:

```bash
dph
```

Two flags are supported:

|Name|Required|Notes|
|----|--------|-----|
|f   |No      |The file to load flag definitions from.  The default is `flags.yml`.|
|d   |No      |Whether to do a dry run i.e. no actual mutations in PostHog.  The default is `false`.|

## Environment Variables

The following environment variables are used but not all are required.

|Name|Required|Notes|
|----|--------|-----|
|DPH_PROJECT_ID|Yes|The PostHog Project ID of the project to manage feature flags in.|
|DPH_PERSONAL_API_KEY|Yes|A PostHog personal API key that has full feature flag scopes.|
|DPH_API_BASE_URL|No|The base URL of the PostHog API.  The default is `https://us.posthog.com`|

## File Schema

The file schema used to define flags is very simple.

```yml
flags:
  FirstFlagKey:
    description: The key of this flag is FirstFlagKey.
  SecondFlagKey:
    description: The key of this flag is SecondFlagKey.
```
