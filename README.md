# Stylelia

Welcome to our entry to the 2021 Chef Hackathon by [Artur Kondas, @youshy](github.com/youshy/) and [Jason Field, @xorima](https://github.com/xorima/).
We each came at this with our own desires as outcomes. For Jason it was a desire to learn more about golang and to build something useful with the language and as a long term Chef contributor it was the perfect opportunity to do just that. For Artur it was a want to learn more about Chef and the ecosystem, along with a desire to build some cool new solutions to problems people are experiancing every day within the ecosystem.

## Why does this exist?

Within the Chef ecosystem there is a utility called [Cookstyle](https://docs.chef.io/workstation/cookstyle/), this is a derrivitive of [rubocop](https://github.com/rubocop/rubocop) and provides [static code analysis](https://www.perforce.com/blog/sca/what-static-analysis) to [Chef-Infra](https://www.chef.io/products/chef-infra) Cookbooks and other Ruby based chef tools.
This is great and enables chef-infra customers to have cookbooks that are inline with chef's recommended best practices, reducing misconfiguration and ensuring a smoother upgrade process to newer versions of Chef-Infra.

Stylelia builds on cookstyle by removing the need to run cookstyle manually on repositores for newer releases of cookstyle. It does this by checking out a targetted repo and running cookstyle on it, if there are changes it then opens a pull request to the upstream repository on GitHub, leaving the developers to review and merge the changes. This removes the concerns in a CI pipeline of being out of date with the latest cookstyle recommendations and means that when people go to write new business value in cookbooks they only need to focus on the code they are writing and not if the existing cookbook is inline with Chef's cookstyle recommendations.

The real point here is to remove the [toil](https://sre.google/sre-book/eliminating-toil/) of running cookstyle from the developers concerns. With a great test suite it could even be fully automated with reviews and merges so no human interaction would be required. Though these cookbook test suites and automated reviews are out of the scope of this solution.

Stylelia is built to run in [AWS](https://aws.amazon.com/) [Lambda](https://docs.aws.amazon.com/lambda/latest/dg/welcome.html) and use [caching](https://aws.amazon.com/caching/) to understand when the cookbook really needs processing.

The Cache stores the current [Default branch](https://git-scm.com/book/en/v2/Git-Branching-Branches-in-a-Nutshell) [Commit Sha](https://git-scm.com/book/en/v2/Git-Tools-Revision-Selection) and the version of cookstyle used at that time. We explicitly check the default branch as not every branch is named `main`. If there is a difference between the cache and the current state or the repository does not exist in the cache then cookstyle is run. If a difference is detected then a [pull request](https://docs.github.com/en/github/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/about-pull-requests) is created. If a pull request already exists we [rebase](https://git-scm.com/docs/git-rebase) the branch and update the Pull Request text to be a true reflection of the changes.

## What is Cookstyle

> Cookstyle is a code linting tool that helps you write better Chef Infra cookbooks by detecting and automatically correcting style, syntax, and logic mistakes in your code.
>
> Cookstyle is powered by the RuboCop linting engine. RuboCop ships with over three-hundred rules, or cops, designed to detect common Ruby coding mistakes and enforce a common coding style. We’ve customized Cookstyle with a subset of those cops that we believe are perfectly tailored for cookbook development. We also ship Chef-specific cops that catch common cookbook coding mistakes, cleanup portions of code that are no longer necessary, and detect deprecations that prevent cookbooks from running on the latest releases of Chef Infra Client.
>
> Cookstyle increases code quality by:
>
> - Enforcing style conventions and best practices.
> - Helping every member of a team author similarly structured code.
> - Maintaining uniformity in the source code.
> - Setting expectations for fellow (and future) project contributors.
> - Detecting deprecated code that creates errors after upgrading to a newer Chef Infra Client release.
> - Detecting common Chef Infra mistakes that cause code to fail or behave incorrectly.

For more information on cookstyle see the [Official Documentation](https://docs.chef.io/workstation/cookstyle/), the list of [Cops (What it is checking for)](https://docs.chef.io/workstation/cookstyle/cops/) or the [Github Repository](https://github.com/chef/cookstyle/)

---

## Setup

The below section covers setup and running Stylelia

## Prerequisits

The following tools must be installed on your machine ahead of developing or running the solution. Items which are required only for development are marked with `Dev Only` before their name.

- `Dev Only` git - [Setup Guide](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- `Dev Only` golang 1.17 - [Setup Guide](https://golang.org/doc/install)
- docker (Must be able to run Linux containers) - [Setup Guide](https://docs.docker.com/get-docker/)
- docker-compose - [Setup Guide](https://docs.docker.com/compose/install/)
- GitHub account - [Setup Guide](https://github.com/join)
- GitHub personal access token (With Full repo access) - [Setup Guide](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)
- A GitHub Orginisation - [Setup Guide](https://docs.github.com/en/organizations/collaborating-with-groups-in-organizations/creating-a-new-organization-from-scratch)
- A single cookbook in a repository within the above GitHub Orginisation that you can access with the personal access token and make branches on directly. You will also need it to have cookstyle failures. If you do not have one already setup please fork [snort](https://github.com/stylelia/snort/) which has been specially modified to have cookstyle failures in it.

You will also need the following ports available on your machine:

- 8081 - Redis Commander
- 6379 - Redis

## Developing

To develop this tool, first you will need to Clone out the repository, it is assumed at this point that you have [generated SSH Keys and configured github to work with SSH Keys](https://docs.github.com/en/authentication/connecting-to-github-with-ssh)

```bash
git clone git@github.com:stylelia/analyser.git
```

After you have cloned the repository you will need to install all of the required golang modules

```bash
go mod tidy
```

You will now need to bring up the backend cache system, which is [Redis](https://redis.io/) and is used in tests. This has all been configured with Docker so you will simply need to run:

```bash
docker-compose up
```

Once these are up you can then access a redis gui application called [Redis Commander](https://github.com/joeferner/redis-commander) on [http://127.0.0.1:8081/](http://127.0.0.1:8081/)

At this point you are ready to make changes to the code. Once you have made changes please ensure they work by running the tests. These should be run in either the `app/analyser`, `pkg/redis` or `pkg/inmemorycache` folders depending on where your changes have been made

The command assumes the current working directory is the root of the repository

```bash
cd <folder>
go test
```

If you choose to open Pull Requests on GitHub then an automated CI pipeline will run tests for all the golang packages listed above as well as [markdownlint](https://github.com/markdownlint/markdownlint) and [yaml lint](http://www.yamllint.com/)

### Running Locally

During development you might want to run the entire application locally. This is also how we recommend running the tool during the hackathon to remove the prerequists around [AWS accounts](https://aws.amazon.com/). We use a local [AWS Lambda docker container](https://github.com/lambci/docker-lambda) to simulate the desired end state environment with additional tools installed to simulate the [Lambda Layers](https://docs.aws.amazon.com/lambda/latest/dg/configuration-layers.html) that are used in the production environment.

It is assumed at this point you have already followed the instructions in Developing and have the following:

- Code checked out
- golang setup
- docker-compose up and running

Additionally at this point you will need the following environment variables set, the examples here are in bash, [This guide](https://www.schrodinger.com/kb/1842) should help with other Shells or Operating Systems

The items within `<>` should be replaced and filled in by yourself, for example `<Repo Name>` should become `TheNameOfMyRepo`

```bash
export REDIS_PASSWORD="MySecurePassword"
export ORGANISATION=<GitHub Orginisation Name>
export NAME=<Repository Name>
export GITHUB_TOKEN=<GitHub Token from token creation>
export GIT_EMAIL=<Your GitHub Email Address>
export GIT_USERNAME=<Your Real Name>
```

Once you have these environment variables set you are able to build and run the Stylelia. The first step is to build the Docker Container, then the go binary and finally run the container on the same network as docker-compose.

It is assumed you are in the root of the repository for these commands.

```bash
docker build . -t stylelia
go build
docker run --rm --network=analyser_redis -e REDIS_HOST="redis" -e REDIS_PORT="6379" -e REDIS_PASSWORD="${REDIS_PASSWORD}" -e ORGANISATION=${ORGANISATION} -e GITHUB_TOKEN="${GITHUB_TOKEN}" -e NAME=${NAME} -e GIT_EMAIL=${GIT_EMAIL} -e GIT_USERNAME=${GIT_USERNAME} -v "$PWD":/var/task:ro,delegated stylelia analyser '{"stylelia": "run"}'
```

The output of this should be `null` which shows that there is no error
Once this has run you should see the Pull Request in your repository. If for some reason you wish to remove the run from the cache you can login to redis-commander and delete the key (see Developing section for details on how to access)
An example Pull Request can be found [here](https://github.com/stylelia/snort/pull/4) you will also see that the commit message contains the same level of detail as the pull request.

## Production

Running in Production is kept out of this repository due to the propriatry nature of this tool and the hosting environment. This tool is designed to run in AWS Lambda and utilise the scale and price advantages that come with lambda's only run when needed nature.

## Future Plans

We see Stylelia as the start of a new way to remove the toil involved in keeping code up to date, in a similar way to how Depenabot handles dependancies. We plan to scale this tool so it will work from a [GitHub App](https://docs.github.com/en/developers/github-marketplace/creating-apps-for-github-marketplace).

We also want to see this tool support multiple static code analysis tools, for example [chefstyle](https://github.com/chef/chefstyle) another rubocop derrivitive tool.

Finally we want to see this automatically triggered when new commits are merged into the default branches or when new versions of the tools are released. This will reduce the feedback loop for these tools to run

## License

© 2021 Jason Field & Artur Kondas. All rights reserved.
