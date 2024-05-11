# Voice recorder bot 

## prerequisites
    
- aws account

- designed aws s3 bucket

- discord app token

## installation

```
git clone https://github.com/Crampustallin/discord_bot.git
```

- get dependencies

when you cloned the repository move to the project's directory and run the command below

```
go mod tidy
```

- environment variables list

```
AWS_REGION # your aws' region

AWS_BUCKET_NAME # your aws bucket name

DISCORD_BOT_TOKEN # your discord app's token
```

after you finished setting env variables run the command below

```
go run ./cmd
```

## The design choices

I'm convinced that the logic of objects shouldn't be depend on each other. So I tried my maximum to seperate them and maintain abstraction.

- code structure

```
discord_bot/
        cmd/
            main.go
        internals/
            aws/
                aws.go # the main logic for aws bucket
            bot/
                tools/
                    tools.go # voice recording logic
                bot.go 
                handlers.go # here you can find all event listeners
                commands.go # in this file one can find functions for different commands for discord user's interaction with the bot
            
```

The bot interacts with aws bucket through Storage interface. It helps for scalability and further modifications of the bot.

## The bot usage

- !url [key] 

The command gets presigned url for [key] object from the s3 bucket.
**key** is required

- !list 

Gets list of all objects in the bucket
