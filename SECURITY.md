# Security Policy

Hexa Policy Orchestrator follows similar security policy as other CNCF projects.

## Reporting a Vulnerability

Please report any security bugs privately to the maintainers listed in the [MAINTAINERS](MAINTAINERS.md) file. We will
fix the issue and coordinate a release date, acknowledging your effort.

## PostgreSQL

Hexa currently uses PostgreSQL to store provider service account credentials. Please ensure your PostgreSQL instance and
client connection is secure. More information can be found below -
* [encryption-options](https://www.postgresql.org/docs/8.1/encryption-options.html)
* [ssl-tcp](https://www.postgresql.org/docs/current/ssl-tcp.html)

## Hawk

Hexa currently uses [Hawk](https://github.com/mozilla/hawk) for HTTP authentication. The credentials used within this
repository are for testing. Please create new credentials for your environment.

## Common Weakness Enumeration (CWE)

```bash
rm -rf .codeql
CODEQL_EXTRACTOR_GO_BUILD_TRACING=on codeql database create .codeql --language=go
codeql database analyze .codeql --off-heap-ram=0 --format=csv --output=codeql-results.csv ../codeql-go/ql/src/codeql-suites/go-code-scanning.qls
```
