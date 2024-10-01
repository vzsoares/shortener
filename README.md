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

- configure aws credentials

local:

1. Clone the repo
1. Install node dependencies with `yarn install`
1. Copy `.env.example` to `.env.dev` and `.env.prod` and fill the values
    > Some values only exist after running Terraform apply
2. Configure your runner in `nx.json tasksRunnerOptions`
3. ⚠️Configure the Terraform provider by changing or removing the `provider.tf` file in each environment.
4. export AWS_PROFILE={your-profile}
5. Deploy **once** with `nx run shortener:first-deploy:dev` and `nx run shortener:first-deploy:prod`
6. Get the cloudfront distribution id and put it on your env
7. Finally deploy with `nx run shortener:deploy:dev` and `nx run shortener:first:prod`

github actions:
- configure a `prod` and `dev` environment with:
  - secrets:
    - AWS_ACCESS_KEY_ID
    - AWS_PROFILE
    - AWS_REGION
    - AWS_SECRET_ACCESS_KEY
  - variables:
    > same as .env
    - API_BASE_URL
    - API_BASE_URL_DOMAIN
    - ARTIFACTS_BUCKET_NAME
    - DYNAMO_URL_TABLE_NAME
    - FRONT_BASE_URL
    - FRONT_BASE_URL_DOMAIN
    - FRONT_BUCKET_NAME
    - FRONT_CLOUDFRONT_DISTRIBUTION_ID
    - GATEWAY_API_NAME
    - NX_RUNNER
    - STAGE

## Usage ☃️

See all available commands bellow and use with: `nx run {project}:{task}:{environment}`.
Ex: `nx run engine:serve:dev` to start a local api.

| Projects    | Tasks                                                   |
|--------------|----------------------------------------------------------|
| engine        | build, lint, serve, test, tidy                           |
| front        | build, lint, publish, serve, test, tidy                  |
| public-api   | build, lint, serve, test, tidy                           |
| utils        | lint, test, tidy                                         |
| shortener    | deploy, first-deploy                                     |

## Architecture 🎨

See in [/docs](/docs)

| [![deployment](./static/hero.jpg)](./docs/deployment.md)                   | [![sequence](./static/sequence_ex.png)](./docs/sequence.md)                   |
| -------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------- |
| [![entity](./static/entity_ex.png)](./docs/entity.md)                | [![usecase](./static/usecase_ex.png)](./docs/use-case.md)                |
| [![swagger_p](./static/swagger_public_ex.png)](./docs/swagger-public-api.yml)                    | [![swagger_e](./static/swagger_engine_ex.png)](./docs/swagger-engine.yml)                    |



