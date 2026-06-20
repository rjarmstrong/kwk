---
date: 2016-03-08T21:07:13+01:00
title: FAQ
weight: 0
---

## Is this free software?

Yes kwk is free for basic accounts and will always be. We will need to introduce paid accounts to help pay for hosting and development.  
## Where are my snippets stored?

They are replicated in near real-time in three secure data centers around the world. Edited snippets may still be in your kwk local cache, to remove this run `rm -fr ~/.kwk` (You will need to sign-in again).

## What security do you have in place?

* SHA256 checksums are used for all snippets and are validated prior to executing snippets and when viewing snippets.
* Our server always chooses the transport cipher so we can upgrade it at any time. TLS 1.2 is always used.
* With our auto-update feature you will have the latest version of kwk.
* End-to-end encryption is coming soon.
* Other user's snippets cannot be run (except for /kwk/... ), you first need to clone them.
* sudo still applies when you run snippets.

## Where is the auto completion for commands?

There is none, we may add this in future but for the time being we felt it was better to focus on using intuitive commands and minimizing keystrokes.

> HINT: Another cool way to save time is that many commands don't need a specific order
> `kwk r hello` and `kwk hello r` are equivalent. Meaning you don't have to back track your cursor after hitting the up arrow.

## When will you have Windows support?

Soon, we do have some work already done and it is high priority. We also currently don't support running dotnet core scripts since it has very specific folder structure requirements to compile and run.

## The colors don't look great on my terminal?

That's not good at all! Please be a sport and <a href="https://github.com/kwk-cli/kwk-cli-issues/issues" target="_blank">log an issue</a>.  
  
> Note: The current theme is designed for dark terminals with 256 colors.

## My script is showing an error when running but the code is valid?

Scripts can only run if the correct runtime is installed! Please ensure the version of the runtime is correct too.

## Can I Contribute?

We are working at making parts of the CLI open source, until then please contribute <a href="https://github.com/kwk-cli/kwk-cli-issues/issues" target="_blank">on the issue log</a> or posting to our Twitter account.
