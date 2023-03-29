<a href="https://opensource.newrelic.com/oss-category/#community-project"><picture><source media="(prefers-color-scheme: dark)" srcset="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/dark/Community_Project.png"><source media="(prefers-color-scheme: light)" srcset="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/Community_Project.png"><img alt="New Relic Open Source community project banner." src="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/Community_Project.png"></picture></a>

# NRI Conviva Integration

![GitHub forks](https://img.shields.io/github/forks/newrelic/nri-conviva?style=social)
![GitHub stars](https://img.shields.io/github/stars/newrelic/nri-conviva?style=social)
![GitHub watchers](https://img.shields.io/github/watchers/newrelic/nri-conviva?style=social)

![GitHub all releases](https://img.shields.io/github/downloads/newrelic/nri-conviva/total)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/newrelic/nri-conviva)
![GitHub last commit](https://img.shields.io/github/last-commit/newrelic/nri-conviva)
![GitHub Release Date](https://img.shields.io/github/release-date/newrelic/nri-conviva)

![GitHub issues](https://img.shields.io/github/issues/newrelic/nri-conviva)
![GitHub issues closed](https://img.shields.io/github/issues-closed/newrelic/nri-conviva)
![GitHub pull requests](https://img.shields.io/github/issues-pr/newrelic/nri-conviva)
![GitHub pull requests closed](https://img.shields.io/github/issues-pr-closed/newrelic/nri-conviva)

The project provides a New Relic Infrastructure integration that uses the
[Conviva v3 Metrics API](https://developer.conviva.com/docs/metrics-api-v3/3e38d9ead39fc-metrics-v3-api-user-guide-beta)
to pull metrics from Conviva and push them into New Relic as dimensional
metrics.

## Installation

As this is a custom integration it must be installed manually. All directions
assume a standard Infrastructure installation.

### Linux
1. Download `nri-conviva` from the [GitHub Release directory](https://github.com/newrelic/nri-conviva/releases)
2. Place `nri-conviva` in `/var/db/newrelic-infra/custom-integrations`
3. Copy the [sample configuration](conviva-config.yml.sample) to
`/etc/newrelic-infra/integrations.d/`

### Windows

> TODO

## Usage

The New Relic Infrastructure integration for Conviva supports the following
Conviva v3 Metrics API concepts.

* [Single metrics and metric groups](https://developer.conviva.com/docs/metrics-api-v3/3e38d9ead39fc-metrics-v3-api-user-guide-beta#single-metrics-and-metric-groups)
* [Multiple singular metrics](https://developer.conviva.com/docs/metrics-api-v3/3e38d9ead39fc-metrics-v3-api-user-guide-beta#retrieving-multiple-singular-metrics-in-a-request)
* [Dimensions](https://developer.conviva.com/docs/metrics-api-v3/3e38d9ead39fc-metrics-v3-api-user-guide-beta#dimensions)
* [Time ranges](https://developer.conviva.com/docs/metrics-api-v3/3e38d9ead39fc-metrics-v3-api-user-guide-beta#time-range)
* [Interval granularity](https://developer.conviva.com/docs/metrics-api-v3/3e38d9ead39fc-metrics-v3-api-user-guide-beta#interval-granularity)
* [Filtering](https://developer.conviva.com/docs/metrics-api-v3/3e38d9ead39fc-metrics-v3-api-user-guide-beta#filtering-by-dimensions)

### Configuration

The New Relic Infrastructure integration for Conviva configuration file follows
[the standard configuration format](https://docs.newrelic.com/docs/infrastructure/host-integrations/infrastructure-integrations-sdk/specifications/host-integrations-standard-configuration-format/) for New Relic Infrastructure
integrations. A [sample](./conviva-config.yml.sample) is included that shows
examples of all options. A description of the options is as follows.

#### Environment variables

The following configuration options are supported in the [`env`](https://docs.newrelic.com/docs/infrastructure/host-integrations/infrastructure-integrations-sdk/specifications/host-integrations-standard-configuration-format/#env)
section of the configuration file.

| Variable Name | Description | Default |
| --- | --- | --- |
| CLIENT_ID | The Conviva v3 API [client ID](https://developer.conviva.com/docs/metrics-v2-api/2aa3a044c485b-metrics-v2-api-user-guide#authentication) | The OS environment variable named `CONVIVA_CLIENT_ID` |
| CLIENT_SECRET | The Conviva v3 API [client secret](https://developer.conviva.com/docs/metrics-v2-api/2aa3a044c485b-metrics-v2-api-user-guide#authentication) | The OS environment variable named `CONVIVA_CLIENT_SECRET` |
| CONFIG_PATH | Path to the [Conviva collector configuration file](https://docs.newrelic.com/docs/infrastructure/host-integrations/infrastructure-integrations-sdk/specifications/host-integrations-standard-configuration-format/#config) | `${config.path}` |

Note that the `CONFIG_PATH` variable will be auto-populated with the path to the
temporary file created by the agent containing the contents of the `config`
section. If you prefer to keep the integration specific configuration separate
from the standard integration configuration, you may pass a static path in the
`CONFIG_PATH` variable.

#### Conviva collector configuration

The Conviva collector configuration is specified in the
[`config`](https://docs.newrelic.com/docs/infrastructure/host-integrations/infrastructure-integrations-sdk/specifications/host-integrations-standard-configuration-format/#config)
section of the configuration file. The collector configuration specifies the
endpoint URL, the default time range and interval granularity, and the metrics
that the integration should collect on each invocation. The following
configuration options are supported.

| Variable Name | Description | Default |
| --- | --- | --- |
| apiV3Url | The Conviva v3 API endpoint | https://api.conviva.com/insights/3.0/metrics |
| startOffset | An offset from the current time for the start of the query time range, specified as a [Go duration string](https://pkg.go.dev/time#ParseDuration) | 20m |
| endOffset | An offset from the current time for the end of the query time range, specified as a [Go duration string](https://pkg.go.dev/time#ParseDuration) | 10m |
| granularity | The time interval granularity for the query, specified in [ISO 8601 format](https://en.wikipedia.org/wiki/ISO_8601#Durations) | PT1M |
| metrics | The array of metric definitions specifying the metrics to collect | [] |

** Time Range Note **

The Conviva v3 Metrics API can return inconsistent results for "recent" data. In
our testing, querying for data in the last 10 minutes at 1 minute granularity
was problematic. This is why the default `endOffset` is `10m`. Conviva is aware
of these issues and may address them in the future. With the default values for
`startOffset`, `endOffset`, and `granularity`, we were able to see consistent
results.

The metrics that the Conviva collector should collect are specfied as a list of
metric definitions in the `metrics` configuration option. The following options
are supported in a metric definition.

| Variable Name | Description | Default |
| --- | --- | --- |
| metric | The name of a single metric to collect | |
| metricGroup | The name of a metric group to collect | |
| names | A list of multiple metric names to collect in a single query | |
| dimensions | A list of group by dimension names to collect for the metric or metric group | [] |
| filters | A set of filtering dimensions to filter results by where each filter is specified as a key:value pair where the key is a dimension name and the value is a list of values to include | {} |
| startOffset | A query specific override for `startOffset` | |
| endOffset | A query specific override for `endOffset` | |
| granularity | A query specific override for `granularity` | |

** Dimensions note **

Dimensions can only be queried one at a time. This has two ramifications.

1. There is a 1:1 relationship between a "group by" dimension and an API call.
   In other words, the more dimensions you specify in the `dimensions` array,
   the more API calls have to be made.
2. There is a 1:1 relationship between a Conviva metric name + dimension and the
   dimensional metric that is created for it in New Relic. This means that it is
   not currently possible to do queries like
   `SELECT average(conviva.bitrate) FROM Metric WHERE browser_name = 'Chrome' *AND* geo_country_code = 'us`.

** Filters note **

[Logical OR filtering](https://developer.conviva.com/docs/metrics-api-v3/3e38d9ead39fc-metrics-v3-api-user-guide-beta#logical-or-filtering-of-the-same-dimension)
is supported by specifying multiple values for a single filtering dimension.
While it is possible to specify multiple values for multiple filtering
dimensions in the configuration, logical `OR` filtering only works for filtering
of the same dimension. For complex logic, a saved filter is required. Currently,
querying with saved filters is not supported.

## Building

Golang is required to build the integration. We recommend Golang 1.18 or higher.

After cloning this repository, go to the directory of the NGINX integration and
build it:

```bash
$ make build
```

The command above executes the tests for the NGINX integration and builds an
executable file called `nri-conviva` under the `bin` directory. 

To start the integration, run `nri-conviva`:

```bash
$ ./bin/nri-conviva
```

If you want to know more about usage of `./bin/nri-conviva`, pass the `-help`
parameter:

```bash
$ ./bin/nri-conviva -help
```

External dependencies are managed through the [go modules](https://blog.golang.org/using-go-modules).
All the external dependencies and its versions are listed in the `go.mod` file.
The vendor folder is not required anymore.

## Testing

To run the tests execute:

```bash
$ make test
```

## Support

> TODO

## Privacy

At New Relic we take your privacy and the security of your information
seriously, and are committed to protecting your information. We must emphasize
the importance of not sharing personal data in public forums, and ask all users
to scrub logs and diagnostic information for sensitive information, whether
personal, proprietary, or otherwise.

We define “Personal Data” as any information relating to an identified or
identifiable individual, including, for example, your name, phone number, post
code or zip code, Device ID, IP address, and email address.

For more information, review [New Relic’s General Data Privacy Notice](https://newrelic.com/termsandconditions/privacy).

## Contribute

We encourage your contributions to improve [project name]! Keep in mind that
when you submit your pull request, you'll need to sign the CLA via the
click-through using CLA-Assistant. You only have to sign the CLA one time per
project.

If you have any questions, or to execute our corporate CLA (which is required
if your contribution is on behalf of a company), drop us an email at
opensource@newrelic.com.

**A note about vulnerabilities**

As noted in our [security policy](../../security/policy), New Relic is committed
to the privacy and security of our customers and their data. We believe that
providing coordinated disclosure by security researchers and engaging with the
security community are important means to achieve our security goals.

If you believe you have found a security vulnerability in this project or any of
New Relic's products or websites, we welcome and greatly appreciate you
reporting it to New Relic through [HackerOne](https://hackerone.com/newrelic).

If you would like to contribute to this project, review [these guidelines](./CONTRIBUTING.md).

To all contributors, we thank you!  Without your contribution, this project
would not be what it is today.  We also host a community project page dedicated
to [Project Name](<LINK TO https://opensource.newrelic.com/projects/... PAGE>).

## License

The [New Relic Integration for Conviva] is licensed under the [Apache 2.0](http://apache.org/licenses/LICENSE-2.0.txt)
License.

> [If applicable: The [project name] also uses source code from third-party libraries. You can find full details on which libraries are used and the terms under which they are licensed in the third-party notices document.]
