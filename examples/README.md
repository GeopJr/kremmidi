<h1 align="center">Frontends</h1>
<h4 align="center">Kremmidi consumers.</h4>

# What is this?

Frontends are the apps that consume the database that Kremmidi populates. They are not meant to be complete by any means but they act as examples or POCs on how to consume the database.

# Guides

- [discordjs.guide](https://discordjs.guide)
- [Slack's bolt tutorial](https://slack.dev/bolt-js/tutorial/getting-started)

# Notes

- No PRs will be accepted on expanding them.
- Please don't host them as is.
- Matrix does not 100% work as it needs further development to handle encrypted rooms and whatnot.

# Screenshots

The screenshots showcase the examples in action. For the sake of not posting huge screenshots, some part messages were deleted (usually there are 15-30 parts).

## Discord

<p align="center">
    <img width=768" src="https://i.imgur.com/ze4Irr9.png" alt="A Discord chat between a User and a bot (BASIL):
    Reply to: User used /platforms.
    BASIL BOT Today at 12:08
    Available platforms: android-aarch64, android-x86, android-armv7, android-x86_64, linux32, linux64, osx64, win32, win64
    Only you can see this • Dismiss message
    Reply to: User used /get.
    BASIL BOT Today at 12:09
    I'm about to send the browser in parts. Download them all.
    Only you can see this • Dismiss message
    BASIL BOT Today at 12:09
    Attachment file: tor-browser-11.5.3-android-armv7-multi.apk.000
    5.00 MB
    Only you can see this • Dismiss message
    BASIL BOT oday at 12:09
    Attachment file: tor-browser-11.5.3-android-armv7-multi.apk.001
    5.00 MB
    Only you can see this • Dismiss message
    BASIL BOT Today at 12:09
    Open the following file in your browser and drag-n-drop all the parts to bake them into one.
    -Preview of onion-rings.html-
    Only you can see this • Dismiss message" />
</p>

## Slack

<p align="center">
    <img width=768" src="https://i.imgur.com/KbFep9W.png" alt="A Slack chat between a User and a bot:
    USER 3:40 PM
    !platforms
    test APP 3.40 PM
    Available platforms: android-aarch64 , android-x86, android-armv7, android-x86_64, linux32, linux64) osx64, win32, win64" />
    <img width=768" src="https://i.imgur.com/BPtWmJ0.png" alt="A Slack chat between a User and a bot:
    USER 3:40 PM
    !platforms
    test APP 3.40 PM
    Available platforms: android-aarch64 , android-x86, android-armv7, android-x86_64, linux32, linux64) osx64, win32, win64
    USER 4:20 PM
    !get android-aarch64.
    test APP 4:20 PM
    I'll start sending the binary in 18 parts. Afterwards I'll send onion-rings with instructions! Please wait for all parts to finish uploading...
    -Attachment tor-browser-11.5.3-android-aarch64-multi.apk-
    Please download onion-rings.html and open it in your browser. Afterwards follow the instructions on screen and drag-n-drop all parts into it to get the original binary.
    -Preview of onion-rings.html-" />
</p>

## Matrix

<p align="center">
    <img width=768" src="https://i.imgur.com/xgWvDZK.png" alt="A Matrix chat between a User and a bot (Kremmidi):
    User: !android-armv7
    kremmidi: I'll start sending the binary in *17* parts. Afterwards I'll send onion-rings with instructions! Please wait for all parts to finish uploading...
    kremmidi: -Attachment tor-browser-11.5.3-android-armv7-multi.apk.000-
    kremmidi: -Attachment tor-browser-11.5.3-android-armv7-multi.apk.001-
    kremmidi: -Attachment tor-browser-11.5.3-android-armv7-multi.apk.002-
    kremmidi: -Attachment tor-browser-11.5.3-android-armv7-multi.apk.003-
    kremmidi: Please download onion-rings.html and open it in your browser. Afterwards follow the instructions on screen and drag-n-drop all parts into it to get the original binary.
    kremmidi: -Attachment onion-rings.html-" />
</p>