# Snitch
<img src="https://github.com/ezeoleaf/snitch/blob/main/img/transparent-snitch.png" width="150">

Snitch is a Slack bot that fetches and publish PRs with pendings reviews. It could publish the messages in channels or via DM

## Disclaimer

I hold no liability for what you do with this bot or what happens to you by using this bot.

## Usage

### Configuring the bot

Before running the bot, you must first set it up so it can connect to Github and Slack API.

To do this, you will need to setup the following environment variables:
```
- ENTERPRISE_GITHUB (Only to true if you are using an enterprise version of Github)
- GITHUB_ADDRES (Only if you are using an enterprise version of github. The URL should be like https://github.{company}.io/)
- GITHUB_TOKEN (Your generated token with **repo** and **user** permissions granted)
- HTTP_ADDRESS (The address of the server where your bot is running)
- SLACK_API_TOKEN (A bot token generated for an Slack application)
```

For generating Github access token you can follow this [guide](https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/creating-a-personal-access-token)

For creating a Slack application you can follow this [guide](https://api.slack.com/authentication/basics#creating)

After creating the Slack application, you need to create two slash commands:
- **repo-prs**: (This will receive a Github repository and returns all the PRs and Issues in it)
- **user-prs**: (This will receive a Github username and returns all the PRs and Issues that the user has pending of review)

![image](https://user-images.githubusercontent.com/10358977/147968887-6cc0530c-8c32-4792-bf4f-252a5defa9e5.png)

## How does it look like?

When using:

- `/repo-prs ezeoleaf/snitch`
![image](https://user-images.githubusercontent.com/10358977/147970481-ddbc62ca-1141-4c46-b4be-fdcd76ea778d.png)

- `/user-prs ezeoleaf`
![image](https://user-images.githubusercontent.com/10358977/147970549-e39363c2-baaf-443e-bd73-98a9edee1f7d.png)

## Have questions? Need help with the bot?

If you're having issues with or have questions about the bot, [file an issue](https://github.com/ezeoleaf/snitch/issues) in this repository so anyone can get back to you.

Or feel free to contact me <ezeoleaf@gmail.com> :)
