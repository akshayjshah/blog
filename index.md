---
title: "Akshay Shah"
description: "Code, cooking, and caffeine."
hidden: true
hide_home: true
---

<img src="/static/headshot-2025-300x300.webp" alt="" class="profile-pic" height="300px" width="300px" fetchpriority="high">

Hi there! I'm Akshay, and you've stumbled across my oft-neglected blog. After
building infrastructure at Uber, Microsoft, and a pile of startups, I'm now the
Field CTO at [Antithesis](https://antithesis.com). We make the world's best
tools for testing distributed systems.

Shoot me an [email](mailto:akshay@akshayshah.org) and say hello, or
[read a bit more about me](/colophon/).

## Writing

{{ range .Posts }}
<div class="post-row">
  <span class="post-date">{{ .Created }}</span>
  <a href="{{ .Link }}">{{ .Title }}</a>{{ if .Via }} ({{ .Via}}){{ end }}
</div>
{{ end }}

## Cooking

{{ range .Recipes }}
<div class="post-row">
  <span class="post-date">{{ .Created }}</span>
  <a href="{{ .Link }}">{{ .Title }}</a>
</div>
{{ end }}
