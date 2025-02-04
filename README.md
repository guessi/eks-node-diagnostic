# eks-node-diagnostic

[![GitHub Actions](https://github.com/guessi/eks-node-diagnostic/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/guessi/eks-node-diagnostic/actions/workflows/go.yml)
[![GoDoc](https://godoc.org/github.com/guessi/eks-node-diagnostic?status.svg)](https://godoc.org/github.com/guessi/eks-node-diagnostic)
[![Go Report Card](https://goreportcard.com/badge/github.com/guessi/eks-node-diagnostic)](https://goreportcard.com/report/github.com/guessi/eks-node-diagnostic)
[![GitHub release](https://img.shields.io/github/release/guessi/eks-node-diagnostic.svg)](https://github.com/guessi/eks-node-diagnostic/releases/latest)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/guessi/eks-node-diagnostic)](https://github.com/guessi/eks-node-diagnostic/blob/main/go.mod)

Friendly `NodeDiagnostic` generator with no Python/SDK required

## 🤔 Why we need this? what it is trying to resolve?

Back on December 3, 2024, Amazon EKS announced [Amazon EKS Auto Mode](https://aws.amazon.com/blogs/containers/getting-started-with-amazon-eks-auto-mode/) and how to use [NodeDiagnostic](https://docs.aws.amazon.com/eks/latest/userguide/auto-get-logs.html) for troubleshooting. However, not all computers have a pre-set Python environment, and it may be difficult for people who are not familiar with Python, let alone solving package dependencies, version conflicts, and virtual environment setting issues.

To remove these obstacles, [eks-node-diagnostic](https://github.com/guessi/eks-node-diagnostic) comes to solve this complex problem, aiming to simplify the entire process into executing binaries, easy to install without having to deal with Python dependencies.

## 🔢 Prerequisites

* An existing Amazon EKS cluster with the node monitoring agent.
* An existing Amazon S3 bucket for storing node logs generated by `NodeDiagnostic`.
* An IAM Role/User with `s3:PutObject` permission (to generate presigned S3 url).

## 👀 Key differences with [official guidance](https://docs.aws.amazon.com/eks/latest/userguide/auto-get-logs.html)

* Single executable binary only, no Python and no AWS SDK required.
* Run anywhere, compatible with Linux (amd64/arm64), Windows (amd64 only), macOS (amd64/arm64).
* Friendly setup with [Homebrew](https://brew.sh/) for mac users.

## 🚀 Quick start

```bash
$ eks-node-diagnostic --help
```

Apply in batch:

```bash
$ cat config.yaml
---
region: us-east-1
expiredSeconds: 300
bucketName: node-diagnostic-EXAMPLE
nodes:
- i-EXAMPLE1111111111
- i-EXAMPLE2222222222
- i-EXAMPLE3333333333
...
```

```bash
$ eks-node-diagnostic -c config.yaml | kubectl apply -f -
nodediagnostic.eks.amazonaws.com/i-EXAMPLE1111111111 created
nodediagnostic.eks.amazonaws.com/i-EXAMPLE2222222222 created
nodediagnostic.eks.amazonaws.com/i-EXAMPLE3333333333 created
...
```

Apply one-by-one slowly:

```bash
$ eks-node-diagnostic -r ${REGION} -n ${NODE} -b ${BUCKET} | kubectl apply -f -
nodediagnostic.eks.amazonaws.com/i-EXAMPLE created
```

## :accessibility: FAQ

Where can I find log archive generated by `NodeDiagnostic`?

* Log archive generated by `NodeDiagnostic` would be placed at path below,

    > `s3://{{ BUCKET }}/node-diagnostic/log__{{ REGION }}__{{ NODE }}__{{ TIMESTAMP }}.tar.gz`

How do I report an issue or submit a feature request?

* Please go for project's [issue page](https://github.com/guessi/eks-node-diagnostic/issues) and describe your idea in detail.

## 👷 Install

### For macOS users (Recommended)

```bash
brew tap guessi/tap && brew update && brew install eks-node-diagnostic
```

### Manually setup (Linux, Windows, macOS)

<details><!-- markdownlint-disable-line -->
<summary>Click to expand!</summary><!-- markdownlint-disable-line -->

#### For Linux users

```bash
curl -fsSL https://github.com/guessi/eks-node-diagnostic/releases/latest/download/eks-node-diagnostic-Linux-$(uname -m).tar.gz -o - | tar zxvf -
mv -vf ./eks-node-diagnostic /usr/local/bin/eks-node-diagnostic
```

#### For macOS users

```bash
curl -fsSL https://github.com/guessi/eks-node-diagnostic/releases/latest/download/eks-node-diagnostic-Darwin-$(uname -m).tar.gz -o - | tar zxvf -
mv -vf ./eks-node-diagnostic /usr/local/bin/eks-node-diagnostic
```

#### For Windows users

```powershell
$SRC = 'https://github.com/guessi/eks-node-diagnostic/releases/latest/download/eks-node-diagnostic-Windows-x86_64.tar.gz'
$DST = 'C:\Temp\eks-node-diagnostic-Windows-x86_64.tar.gz'
Invoke-RestMethod -Uri $SRC -OutFile $DST
```

</details>

## ⚖️ License

[Apache-2.0](LICENSE)
