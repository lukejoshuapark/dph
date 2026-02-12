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

Four flags are supported:

|Name|Required|Notes|
|----|--------|-----|
|f   |No      |The file to load flag definitions from.  The default is `flags.yml`.|
|d   |No      |Whether to do a dry run i.e. no actual mutations in PostHog.  The default is `false`.|
|r   |No      |Whether to operate in reverse - when enabled, flags that exist in PostHog will be populated in to the flag definition file.|
|s   |No      |Whether to operate in "safe" mode - when enabled, flags that exist in PostHog that don't exist in the flag definition file will _not_ be deleted.|

## Environment Variables

The following environment variables are used but not all are required.

|Name|Required|Notes|
|----|--------|-----|
|DPH_ORGANIZATION_ID|Yes|The PostHog Organization ID that the project belongs to.|
|DPH_PROJECT_ID|Yes|The PostHog Project ID of the project to manage feature flags in.|
|DPH_PERSONAL_API_KEY|Yes|A PostHog personal API key that has full feature flag scopes.|
|DPH_API_BASE_URL|No|The base URL of the PostHog API.  The default is `https://us.posthog.com`.|
|DPH_GOOGLE_CHAT_WEBHOOK_URL|No|A webhook URL for a Google Chat space to send notifications to.  Notifications will always be sent if this is present.|

## File Schema

The file schema used to define flags is very simple.

```yml
flags:
  FirstFlagKey:
    description: The key of this flag is FirstFlagKey.
  SecondFlagKey:
    description: The key of this flag is SecondFlagKey.
    exclude:
      - Production
```

The `exclude` property is a list of project names.  If the name of the PostHog project specified using `DPH_PROJECT_ID`
matches one of the names in the exclude list, `dph` determines that this flag should not exist for the project, and will
not create it/delete it if it exists.
