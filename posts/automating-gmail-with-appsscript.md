---
title: "Automating Gmail with AppsScript"
description: "Smarter email automation with Google AppsScript."
created: 2020-06-01
---

We're all drowning in email. The problem isn't the blatant spam --- it's easy to
unsubscribe from most of that and block the rest. The real killer is the email
that's *sort of* interesting: discussion on projects you're peripherally
involved in, weekly newsletters, chatter on GitHub issues you filed months ago,
and all the other mail you skim when you have the time and skip if you're busy.

Triaging these emails manually can be an exhausting game of whack-a-mole. But
if you're a Gmail user, there's hope --- you can use
[AppsScript](https://www.google.com/script/start/) to automate big parts of the
job. Best of all, it's free! In this post, I'll walk you through creating your
first AppsScript project and give you a taste of what you can accomplish. (Hat
tip to [Prashant Varanasi](https://github.com/prashantv), who first introduced
me to AppsScript.)

## Setup

[Create a new AppsScript project](https://www.google.com/script/start), then
add a small function and save the project. This code doesn't do anything yet,
but we'll add to it later on.

```js
function processMail() {
}
```

AppsScript includes a cron-like triggers service, but you have to be careful:
Google
[limits](https://developers.google.com/apps-script/guides/services/quotas) the
amount of time your script spends running, and it also limits the number of
Gmail operations you can perform. For me, running my email management script
every 15 minutes keeps my inbox nicely groomed without blowing through my
quotas. Set up a time-based trigger for the `processMail` function right from
the editor by going to "Edit," then "Current project's triggers," and finally
"Add Trigger" in the bottom right.

## Simple: mark archived threads read

The simplest bits of my email automation work on one thread at a time, with no
further context required: for example, I use a small function to make sure that
all my archived email is marked read. (It's puzzling to me that this isn't the
default, since it's the only way to make the badge count on Gmail's mobile apps
useful.)

```js
function processMail() {
  const rules = [
    markArchivedRead
  ];
  for (let rule of rules) {
    rule();
  }
}

function markArchivedRead() {
  return eachThread("gmail: mark archived threads read",
                    "-in:inbox is:unread",
                    function(thread) { thread.markRead(); });
}

function eachThread(operation, query, f) {
  const threads = GmailApp.search(query).slice(0, 100);
  if (threads.length <= 0) {
    Logger.log("%s: 0 threads match query %s", operation, query);
    return;
  }
  const n = threads.length;
  Logger.log("%s: %s threads match query %s", operation, n, query);
  for (let thread of threads) {
    f(thread);
  }
  Logger.log("%s: processed %s threads, done!", n, operation);
}
```

I have 5--10 small tasks like this active most of the time, and the `eachThread`
helper keeps each of them nice and short. Note that `eachThread` limits itself
to processing 100 threads per invocation, keeping execution time short and
capping the number of Gmail operations consumed.

I use similar per-thread functions for a variety of simple tasks, most of which
either groom my archived mail or trim my inbox when I start falling behind.
Examples include:

* Un-starring everything in my archive. I use stars to protect inbox threads
  from further automation; for reference material, I use a label.
* Retroactively applying labels to my archive, since the Gmail UI hangs when
  trying to modify more than a few thousand threads.
* Archiving promotions (using Google's automatic categorization) more than
  three days old, unless they're starred.
* Deleting high-volume notifications after they've been archived for a year.
  Email from exception-tracking systems like Sentry is a good candidate for
  time-based deletion: it's valuable information when it's fresh, but it's not
  worth keeping gigabytes of it around forever.
* For a while, I automatically archived long-running threads with lots of
  participants, no emails from me, and no mention of teams or projects I was
  responsible for. This took a lot of maintenance and had a fair number of
  false positives, but kept my sanity intact while I sat on some *extremely*
  bureaucratic engineering committees.

## Moderate: limit inbox size

Despite my best intentions, newsletters and other non-essential threads often
accumulate in my inbox. It's painful to clear this backlog out by hand, because
I actually want to read most of it: I end up agonizing over whether I've got
time to read *just one more* interesting article or thread, opening a million
browser tabs, and burning hours of time better spent elsewhere.

Instead, I use filters to label interesting-but-optional mail as it arrives. (I
have a *lot* of filters, which I recently started managing with
[`gmailctl`](https://github.com/mbrt/gmailctl).) If more than a hundred of
these emails pile up in my inbox, I archive the older messages until only 50
remain. I also tag the auto-archived messages, so I know that I haven't read
them if they show up in search later on.

```js
function limitInbox() {
  const max = 100;
  const op = "gmail: limit inbox";
  const purged = GmailApp.getUserLabelByName("optional/purged");
  const threads = GmailApp.search("in:inbox label:optional");
  Logger.log("%s: %s optional threads", op, threads.length);
  if (threads.length <= max) {
    Logger.log("%s: done!", op);
    return;
  }
  let n = 0;
  for (let thread of threads.slice(Math.floor(max/2))) {
    if (!thread.hasStarredMessages()) {
     thread.addLabel(purged);
     thread.moveToArchive();
     n++;
    }
  }
  if (!onVacation()) {
    pushSMS(`auto-archived ${n} non-essential emails!`);
  }
  Logger.log("%s: archived %s threads, done!", n);
}

function onVacation() {
  const email = Session.getEffectiveUser().getEmail();
  const cal = CalendarApp.getCalendarById(email);
  for (let event of cal.getEventsForDay(new Date())) {
    let t = event.getTitle();
    if (t.includes("OOO") || t.includes("PTO") {
      return true;
    }
  }
  return false;
}

function pushSMS(msg) {
    // T-Mobile, Sprint, Verizon, and AT&T all support email-to-SMS
    GmailApp.sendEmail('1234567890@vtext.com', 'AppsScript', msg);
}
```

With this de-bulking script active, going on vacation or getting busy for a
week doesn't leave me with an hour-long inbox cleanup chore. It's surprisingly
liberating.

I use the `onVacation` and `pushSMS` functions regularly: the first lets me
toggle vacation-only behavior with minimal effort, and the second notifies me
if my scripts are running amok.

## Complex: reduce notification spam

The most complex portions of my AppsScript project selectively archive
notifications. Code review systems like Phabricator and GitHub, exception
trackers like Sentry, and many RSS-like subscriptions send tons of
notifications. Often, I'm only interested in the oldest or newest unread
notification for each item.

For example, I love reading trashy, RPG-inspired web novels on [Royal
Road](https://www.royalroad.com/). They send me an email each time a new
chapter gets published in a book I'm following, but I only catch up on my
trashy reading a few times a week. Rather than letting all the notifications
sit in my inbox, I'd rather keep only the oldest email for each book.

```js
function queueLitRPG() {
  const op = "gmail: queue litRPG";
  const threads = GmailApp
    .search('in:inbox from:royalroad.com subject:"New Chapter of"');
  let unread = {};
  for (let thread of threads) {
    const book = thread
      .getFirstMessageSubject()
      .replace(/New Chapter of/, '')
      .trim();
    const chapter = {
      date: thread.getLastMessageDate(),
      thread: thread
    };
    if (unread[book] == undefined) {
      unread[book] = [chapter];
    } else {
      unread[book].push(chapter);
    }
  }
  for (const [book, chapters] of Object.entries(unread)) {
    // sort most recent first
    const sorted = chapters
      .slice()
      .sort((a, b) => b.date - a.date);
    // keep the oldest
    for (let chapter of sorted.slice(0, -1)) {
      chapter.thread.moveToArchive();
    }
    Logger.log("%s: done with %s", op, book);
  }
}
```

I use a similar approach to:

* Keep only the latest Sentry notification for each exception.
* Keep only the latest GitHub and Phabricator notifications that new commits
  have been pushed to a pull request that I'm reviewing.
* Archive Phabricator diff notifications (like GitHub pull requests) sent to a
  group if someone else has already started a review. Phab makes this very easy
  by [stamping lots of
  information](https://secure.phabricator.com/book/phabricator/article/mail_rules/#stamps-and-gmail)
  into email headers. I wish GitHub supported something similar.

## Calendars: the final frontier

I haven't worked much with the calendar support in AppsScript yet, but there's
so much low-hanging fruit. I'd love to try:

* Declining meetings if they don't leave me time to eat.
* Scheduling travel time around meetings when necessary.
* Emailing me a weekly summary of where my time went and who I met with.
