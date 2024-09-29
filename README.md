# Shortener

![GitHub watchers](https://img.shields.io/github/watchers/vzsoares/shortener?style=for-the-badge)
![GitHub forks](https://img.shields.io/github/forks/vzsoares/shortener?style=for-the-badge)
![GitHub stars](https://img.shields.io/github/stars/vzsoares/shortener?style=for-the-badge)
![GitHub release](https://img.shields.io/github/v/release/vzsoares/shortener?style=for-the-badge)
![GitHub license](https://img.shields.io/github/license/vzsoares/shortener?style=for-the-badge)
[![Connect on linkedin](https://img.shields.io/badge/Connect-grey?style=for-the-badge&logo=linkedin)](https://www.linkedin.com/in/vinicius-zenha/)

<p align="center">
<img src="./static/hero.jpg" width="500px">
</p>

Deploy your own shortener service in the aws free tier. Made in Go with minimal dependencies. Raw HTML JS and _tailwind_ front. Main engine to integrate with other internal services.

## Features 📃

-   Fully serverless pay per request
-   Terraform infra
-   Github actions
-   Generic main engine
    -   Easily integrate using a api-key
-   Bff implementation example

### Requirements 🛠️

-   make
-   aws-cli
-   terraform
-   go
-   node
-   yarn

## Setup 🦩

Make sure to have all requirements.

1. Clone the repo
1. Install node dependencies with `yarn install`
1. Copy `.env.example` to `.env.dev` and `.env.prod` and fill the values
    > Some values only exist after running Terraform apply
2. Configure your runner in `nx.json tasksRunnerOptions`
2. Configure or remove the `terraform remoteBackend bucket` in each environment
3. Deploy with `nx run shortener:deploy:dev` and `nx run shortener:deploy:prod`

> Work in progress...

## Usage ☃️

As this project uses nx you can see all available commands with `nx show projects` then `nx show project {projectName}`

> Work in progress...

## Architecture 🎨

See in [/docs](/docs)

| [![deployment](./static/hero.jpg)](./docs/deployment.md)                   | [![sequence](./static/sequence_ex.png)](./docs/sequence.md)                   |
| -------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------- |
| [![entity](./static/entity_ex.png)](./docs/entity.md)                | [![usecase](./static/usecase_ex.png)](./docs/use-case.md)                |
| [![swagger_p](./static/swagger_public_ex.png)](./docs/swagger-public-api.yml)                    | [![swagger_e](./static/swagger_engine_ex.png)](./docs/swagger-engine.yml)                    |



